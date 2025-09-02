package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"

	internal_auth "github.com/studyplatform/backend/internal/auth"
	internal_material "github.com/studyplatform/backend/internal/material"
	internal_note "github.com/studyplatform/backend/internal/note"
	internal_notification "github.com/studyplatform/backend/internal/notification"
	internal_post "github.com/studyplatform/backend/internal/post"
	internal_realtime "github.com/studyplatform/backend/internal/realtime"
	internal_room "github.com/studyplatform/backend/internal/room"
	internal_session "github.com/studyplatform/backend/internal/session"
	internal_todo "github.com/studyplatform/backend/internal/todo"
	pkg_auth "github.com/studyplatform/backend/pkg/auth"
	"github.com/studyplatform/backend/pkg/database"
	"github.com/studyplatform/backend/pkg/logger"
	"github.com/studyplatform/backend/pkg/middleware"
	"github.com/studyplatform/backend/pkg/monitoring"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		// It's okay if .env doesn't exist, we'll use default values
		fmt.Println("Warning: .env file not found, using default values")
	}

	// Initialize logger
	logger.Init()
	defer logger.Close()

	// Get port from environment variable
	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "8080" // Default port
	}

	// Set Gin mode based on environment
	env := os.Getenv("ENV")
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Connect to MongoDB
	mongoClient, err := database.NewMongoClient()
	if err != nil {
		logger.Fatal("Failed to connect to MongoDB", logger.Field("error", err))
	}
	defer mongoClient.Close()

	// Run migration to add unique_id to all users
	if err := internal_auth.MigrateAddUniqueIDToUsers(mongoClient); err != nil {
		logger.Fatal("Migration failed", logger.Field("error", err))
	}

	// Run migration to fix password field names
	if err := internal_auth.MigratePasswordFields(mongoClient); err != nil {
		logger.Fatal("Password migration failed", logger.Field("error", err))
	}

	// Run migration to remove username uniqueness constraints
	if err := internal_auth.MigrateRemoveUsernameUniqueness(mongoClient); err != nil {
		logger.Fatal("Username uniqueness migration failed", logger.Field("error", err))
	}

	// Run migration to convert friend relationships to use UniqueIDs
	if err := internal_auth.MigrateFriendsToUniqueID(mongoClient); err != nil {
		logger.Fatal("Friends migration failed", logger.Field("error", err))
	}

	jwtManager := pkg_auth.NewManager()

	// Initialize WebSocket hub
	hub := internal_realtime.NewHub(mongoClient)
	go hub.Run()

	// Initialize health checker and monitoring
	version := os.Getenv("APP_VERSION")
	if version == "" {
		version = "1.0.0"
	}
	healthChecker := monitoring.NewHealthChecker(mongoClient, version, env)
	healthChecker.StartPeriodicChecks()

	// Initialize rate limiter
	rateLimiter := middleware.NewRateLimiter(nil) // Use default config
	defer rateLimiter.Close()

	// Create router
	router := gin.New()

	// Apply middleware
	middlewareManager := middleware.NewMiddleware()
	router.Use(gin.Recovery())
	router.Use(middleware.CORS()) // Move CORS to the top
	router.Use(middlewareManager.Logger())
	router.Use(middlewareManager.RequestID())

	// Apply rate limiting to all routes
	router.Use(rateLimiter.RateLimit())

	// Register routes
	registerRoutes(router, mongoClient, jwtManager, middlewareManager, hub, healthChecker, rateLimiter)

	// Create HTTP server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Run server in a goroutine
	go func() {
		logger.Info(fmt.Sprintf("Starting server on port %s", port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", logger.Field("error", err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", logger.Field("error", err))
	}

	logger.Info("Server exited properly")
}

func registerRoutes(router *gin.Engine, mongoClient *database.MongoClient, jwtManager *pkg_auth.Manager, middlewareManager *middleware.Middleware, hub *internal_realtime.Hub, healthChecker *monitoring.HealthChecker, rateLimiter *middleware.RateLimiter) {
	// API Version
	apiV1 := router.Group("/api/v1")

	// Enhanced health checks using monitoring system
	apiV1.GET("/health", healthChecker.HealthCheckHandler())
	apiV1.GET("/health/simple", healthChecker.SimpleHealthCheckHandler())
	apiV1.GET("/metrics", healthChecker.MetricsHandler())

	// Rate limiting stats (admin only)
	apiV1.GET("/admin/rate-limit-stats", func(c *gin.Context) {
		// In production, add admin authentication here
		c.JSON(http.StatusOK, rateLimiter.GetRateLimitStats())
	})

	// Friends health check
	apiV1.GET("/friends/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "friends_api_ok",
			"timestamp": time.Now().Unix(),
		})
	})

	// Friends debug endpoint (no auth required)
	apiV1.GET("/friends/debug", func(c *gin.Context) {
		// Test database connection
		users := mongoClient.GetCollection(database.CollectionNames.Users)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		count, err := users.CountDocuments(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Database connection failed",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":             "debug_ok",
			"database_connected": true,
			"total_users":        count,
			"timestamp":          time.Now().Unix(),
		})
	})

	// Auth routes
	authRoutes := apiV1.Group("/auth")
	{
		// Add a simple test endpoint
		authRoutes.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message":   "Auth routes are working",
				"timestamp": time.Now().Unix(),
			})
		})

		authRoutes.POST("/register", internal_auth.RegisterHandler(mongoClient, jwtManager))
		authRoutes.POST("/login", internal_auth.LoginHandler(mongoClient, jwtManager))
		authRoutes.POST("/logout", middlewareManager.Auth(), internal_auth.LogoutHandler(mongoClient, jwtManager))
		authRoutes.POST("/refresh", internal_auth.RefreshTokenHandler(mongoClient, jwtManager))

		// Add /me endpoint with Auth middleware
		authRoutes.GET("/me", middlewareManager.Auth(), internal_auth.MeHandler(mongoClient))
		// Add profile update endpoint
		authRoutes.PUT("/me", middlewareManager.Auth(), internal_auth.UpdateProfileHandler(mongoClient))
		// Add XP update endpoint
		authRoutes.PUT("/xp", middlewareManager.Auth(), internal_auth.UpdateXPHandler(mongoClient))
	}

	// Friends routes
	friendsRoutes := apiV1.Group("/friends")
	friendsRoutes.Use(middlewareManager.Auth())
	{
		friendsRoutes.POST("/request", internal_auth.SendFriendRequestHandler(mongoClient))
		friendsRoutes.PUT("/:id/accept", internal_auth.AcceptFriendRequestHandler(mongoClient))
		friendsRoutes.PUT("/:id/reject", internal_auth.RejectFriendRequestHandler(mongoClient))
		friendsRoutes.DELETE("/:id/remove", internal_auth.RemoveFriendHandler(mongoClient))
		friendsRoutes.GET("/", internal_auth.ListFriendsHandler(mongoClient))
		friendsRoutes.GET("/requests", internal_auth.ListFriendRequestsHandler(mongoClient))
	}
	// User search/profile routes
	userRoutes := apiV1.Group("/users")
	userRoutes.Use(middlewareManager.Auth())
	{
		userRoutes.GET("/", internal_auth.SearchUsersHandler(mongoClient))
		userRoutes.GET("", internal_auth.SearchUsersHandler(mongoClient))        // Handle both with and without trailing slash
		userRoutes.GET("/search", internal_auth.SearchUsersHandler(mongoClient)) // Alternative search endpoint
		userRoutes.GET("/:id", internal_auth.GetUserByIDHandler(mongoClient))    // This must be last to avoid catching search requests
	}

	// Room routes
	roomRoutes := apiV1.Group("/rooms")
	roomRoutes.Use(middlewareManager.Auth())
	{
		roomRoutes.GET("/", internal_room.ListRoomsHandler(mongoClient))
		roomRoutes.POST("/", internal_room.CreateRoomHandler(mongoClient))
		roomRoutes.POST("/join", internal_room.JoinRoomByCodeHandler(mongoClient))

		// Room-specific sub-routes must come BEFORE the general :id route
		roomRoutes.GET("/:id/notes/", internal_note.GetRoomNotesHandler(mongoClient))
		roomRoutes.GET("/:id/materials/", internal_material.GetRoomMaterialsHandler(mongoClient))
		roomRoutes.POST("/:id/generate-code", internal_room.GenerateInvitationCodeHandler(mongoClient))
		roomRoutes.POST("/:id/invite", internal_room.InviteUserToRoomHandler(mongoClient))     // Invite user to room
		roomRoutes.POST("/:id/accept", internal_room.AcceptRoomInvitationHandler(mongoClient)) // Accept room invitation
		roomRoutes.POST("/:id/leave", internal_room.LeaveRoomHandler(mongoClient))
		roomRoutes.POST("/:id/enter", internal_room.EnterRoomHandler(mongoClient)) // Enter/join a room

		// General room CRUD routes must come LAST
		roomRoutes.GET("/:id", internal_room.GetRoomHandler(mongoClient))
		roomRoutes.PUT("/:id", internal_room.UpdateRoomHandler(mongoClient))
		roomRoutes.DELETE("/:id", internal_room.DeleteRoomHandler(mongoClient))
	}

	// Test route for debugging (remove in production)
	apiV1.GET("/test-room", internal_room.TestRoomHandler(mongoClient))
	apiV1.GET("/debug-rooms", internal_room.DebugListAllRoomsHandler(mongoClient))
	apiV1.GET("/simple-test", internal_room.SimpleTestHandler(mongoClient))

	// Session routes
	sessionRoutes := apiV1.Group("/sessions")
	sessionRoutes.Use(middlewareManager.Auth())
	{
		sessionRoutes.POST("/start", internal_session.StartSessionHandler(mongoClient))
		sessionRoutes.POST("/end", internal_session.EndSessionHandler(mongoClient))
		sessionRoutes.GET("/", internal_session.ListSessionsHandler(mongoClient))
		sessionRoutes.POST("/:id/ping", internal_session.ActivityPingHandler(mongoClient))
		sessionRoutes.GET("/stats", internal_session.GetUserSessionStats(mongoClient))
		sessionRoutes.GET("/privileges", internal_session.CheckXPPrivileges(mongoClient))
	}

	// Material routes
	materialRoutes := apiV1.Group("/materials")
	materialRoutes.Use(middlewareManager.Auth())
	{
		materialRoutes.GET("/", internal_material.ListMaterialsHandler(mongoClient))
		materialRoutes.POST("/", internal_material.CreateMaterialHandler(mongoClient))
		materialRoutes.GET("/:id", internal_material.GetMaterialHandler(mongoClient))
		materialRoutes.PUT("/:id", internal_material.UpdateMaterialHandler(mongoClient))
		materialRoutes.DELETE("/:id", internal_material.DeleteMaterialHandler(mongoClient))
		materialRoutes.POST("/upload-url", internal_material.GenerateUploadURLHandler(mongoClient))
		materialRoutes.POST("/confirm-upload", internal_material.ConfirmUploadHandler(mongoClient))
		materialRoutes.POST("/:id/share", internal_material.ShareMaterialHandler(mongoClient))
	}

	// TODO routes
	todoRoutes := apiV1.Group("/todos")
	todoRoutes.Use(middlewareManager.Auth())
	{
		todoRoutes.GET("/", internal_todo.ListTodosHandler(mongoClient))
		todoRoutes.POST("/", internal_todo.CreateTodoHandler(mongoClient))
		todoRoutes.GET("/:id", internal_todo.GetTodoHandler(mongoClient))
		todoRoutes.PUT("/:id", internal_todo.UpdateTodoHandler(mongoClient))
		todoRoutes.PUT("/:id/complete", internal_todo.CompleteTodoHandler(mongoClient))
		todoRoutes.DELETE("/:id", internal_todo.DeleteTodoHandler(mongoClient))
	}

	// Note routes
	noteRoutes := apiV1.Group("/notes")
	noteRoutes.Use(middlewareManager.Auth())
	{
		noteRoutes.GET("/", internal_note.ListNotesHandler(mongoClient))
		noteRoutes.POST("/", internal_note.CreateNoteHandler(mongoClient))
		noteRoutes.GET("/:id", internal_note.GetNotesHandler(mongoClient))
		noteRoutes.PUT("/:id", internal_note.UpdateNoteHandler(mongoClient))
		noteRoutes.DELETE("/:id", internal_note.DeleteNoteHandler(mongoClient))
	}

	// Posts routes
	posts := apiV1.Group("/posts")
	{
		posts.GET("/", middlewareManager.Auth(), internal_post.ListPostsHandler(mongoClient))
		posts.POST("/", middlewareManager.Auth(), internal_post.CreatePostHandler(mongoClient))
		posts.PUT("/:id/like", middlewareManager.Auth(), internal_post.LikePostHandler(mongoClient))
		posts.DELETE("/:id", middlewareManager.Auth(), internal_post.DeletePostHandler(mongoClient))
		posts.POST("/:postId/comments", middlewareManager.Auth(), internal_post.CreateCommentHandler(mongoClient))
		posts.PUT("/comments/:commentId/like", middlewareManager.Auth(), internal_post.LikeCommentHandler(mongoClient))
	}

	// Notification routes
	notificationRoutes := apiV1.Group("/notifications")
	notificationRoutes.Use(middlewareManager.Auth())
	{
		notificationRoutes.GET("/", internal_notification.ListNotificationsHandler(mongoClient))
		notificationRoutes.PUT("/:id/read", internal_notification.MarkNotificationReadHandler(mongoClient))
		notificationRoutes.DELETE("/:id", internal_notification.DeleteNotificationHandler(mongoClient))
		notificationRoutes.POST("/", internal_notification.CreateNotificationHandler(mongoClient))
		notificationRoutes.DELETE("/clear-friend-requests", internal_notification.ClearFriendRequestNotificationsHandler(mongoClient))
		notificationRoutes.DELETE("/clear-all", internal_notification.ClearAllNotificationsHandler(mongoClient))
	}

	// WebSocket and real-time routes
	realtimeRoutes := apiV1.Group("/realtime")
	realtimeRoutes.Use(middlewareManager.Auth())
	{
		realtimeRoutes.GET("/chat/:roomId", internal_realtime.GetRoomChatHistory(mongoClient))
		realtimeRoutes.GET("/online/:roomId", internal_realtime.GetOnlineUsersInRoom(hub))
	}

	// WebSocket route - no auth middleware (handles auth in WebSocket handler)
	apiV1.GET("/realtime/ws", internal_realtime.WebSocketHandler(hub))
}
