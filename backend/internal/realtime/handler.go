package realtime

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/studyplatform/backend/pkg/database"
	"github.com/studyplatform/backend/pkg/models"
)

// Message types for WebSocket communication
const (
	MessageTypeJoin        = "join"
	MessageTypeLeave       = "leave"
	MessageTypeChat        = "chat"
	MessageTypeTyping      = "typing"
	MessageTypeStopTyping  = "stop_typing"
	MessageTypeUserOnline  = "user_online"
	MessageTypeUserOffline = "user_offline"
	MessageTypeError       = "error"
	MessageTypeSuccess     = "success"
	// WebRTC signaling message types
	MessageTypeRTCOffer     = "rtc_offer"
	MessageTypeRTCAnswer    = "rtc_answer"
	MessageTypeRTCCandidate = "rtc_candidate"
	MessageTypeStartCall    = "start_call"
	MessageTypeEndCall      = "end_call"
	MessageTypeCallDeclined = "call_declined"
)

// WebSocket message structure
type WSMessage struct {
	Type      string                 `json:"type"`
	RoomID    string                 `json:"roomId,omitempty"`
	UserID    string                 `json:"userId,omitempty"`
	Username  string                 `json:"username,omitempty"`
	Content   string                 `json:"content,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

// Client represents a WebSocket connection
type Client struct {
	ID       string
	UserID   string
	Username string
	RoomID   string
	Conn     *websocket.Conn
	Send     chan WSMessage
	Hub      *Hub
}

// Hub maintains active clients and broadcasts messages
type Hub struct {
	clients     map[*Client]bool
	broadcast   chan WSMessage
	register    chan *Client
	unregister  chan *Client
	rooms       map[string]map[*Client]bool
	mutex       sync.RWMutex
	mongoClient *database.MongoClient
}

// NewHub creates a new WebSocket hub
func NewHub(mongoClient *database.MongoClient) *Hub {
	return &Hub{
		clients:     make(map[*Client]bool),
		broadcast:   make(chan WSMessage, 1000), // Increased capacity
		register:    make(chan *Client, 100),    // Increased capacity
		unregister:  make(chan *Client, 100),    // Increased capacity
		rooms:       make(map[string]map[*Client]bool),
		mongoClient: mongoClient,
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

// registerClient registers a new client
func (h *Hub) registerClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.clients[client] = true

	if h.rooms[client.RoomID] == nil {
		h.rooms[client.RoomID] = make(map[*Client]bool)
	}
	h.rooms[client.RoomID][client] = true

	// Notify other clients in the room that a user joined
	joinMessage := WSMessage{
		Type:      MessageTypeUserOnline,
		RoomID:    client.RoomID,
		UserID:    client.UserID,
		Username:  client.Username,
		Timestamp: time.Now(),
	}

	h.broadcastToRoom(client.RoomID, joinMessage, client)
	log.Printf("Client %s joined room %s", client.Username, client.RoomID)
}

// unregisterClient unregisters a client
func (h *Hub) unregisterClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.Send)

		if room, ok := h.rooms[client.RoomID]; ok {
			delete(room, client)
			if len(room) == 0 {
				delete(h.rooms, client.RoomID)
			}
		}

		// Notify other clients in the room that a user left
		leaveMessage := WSMessage{
			Type:      MessageTypeUserOffline,
			RoomID:    client.RoomID,
			UserID:    client.UserID,
			Username:  client.Username,
			Timestamp: time.Now(),
		}

		h.broadcastToRoom(client.RoomID, leaveMessage, nil)
		log.Printf("Client %s left room %s", client.Username, client.RoomID)
	}
}

// broadcastMessage broadcasts a message to appropriate clients
func (h *Hub) broadcastMessage(message WSMessage) {
	switch message.Type {
	case MessageTypeChat:
		// Save chat message to database
		h.saveChatMessage(message)
		// Broadcast to room
		h.broadcastToRoom(message.RoomID, message, nil)
	case MessageTypeTyping, MessageTypeStopTyping:
		// Broadcast typing indicators to room (except sender)
		h.broadcastToRoom(message.RoomID, message, h.findClientByUserID(message.UserID))
	case MessageTypeRTCOffer, MessageTypeRTCAnswer, MessageTypeRTCCandidate:
		// Handle WebRTC signaling messages
		h.handleRTCSignaling(message)
	case MessageTypeStartCall, MessageTypeEndCall, MessageTypeCallDeclined:
		// Broadcast call events to room
		h.broadcastToRoom(message.RoomID, message, nil)
	default:
		h.broadcastToRoom(message.RoomID, message, nil)
	}
}

// broadcastToRoom broadcasts a message to all clients in a specific room
func (h *Hub) broadcastToRoom(roomID string, message WSMessage, excludeClient *Client) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if room, ok := h.rooms[roomID]; ok {
		for client := range room {
			if excludeClient != nil && client == excludeClient {
				continue
			}
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(h.clients, client)
				delete(room, client)
			}
		}
	}
}

// findClientByUserID finds a client by user ID
func (h *Hub) findClientByUserID(userID string) *Client {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for client := range h.clients {
		if client.UserID == userID {
			return client
		}
	}
	return nil
}

// saveChatMessage saves a chat message to the database
func (h *Hub) saveChatMessage(message WSMessage) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := h.mongoClient.GetCollection(database.CollectionNames.ChatMessages)

	chatMessage := bson.M{
		"room_id":    message.RoomID,
		"user_id":    message.UserID,
		"username":   message.Username,
		"content":    message.Content,
		"timestamp":  message.Timestamp,
		"created_at": time.Now(),
	}

	_, err := collection.InsertOne(ctx, chatMessage)
	if err != nil {
		log.Printf("Error saving chat message: %v", err)
	}
}

// handleRTCSignaling handles WebRTC signaling messages
func (h *Hub) handleRTCSignaling(message WSMessage) {
	// Extract target user ID from message data
	if targetUserID, ok := message.Data["targetUserId"].(string); ok {
		targetClient := h.findClientByUserID(targetUserID)
		if targetClient != nil {
			select {
			case targetClient.Send <- message:
			default:
				log.Printf("Failed to send RTC message to user %s", targetUserID)
			}
		}
	} else {
		// If no specific target, broadcast to room (for group calls)
		h.broadcastToRoom(message.RoomID, message, h.findClientByUserID(message.UserID))
	}
}

// JWTClaims represents the claims in a JWT token
type JWTClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Exp      int64  `json:"exp"`
}

// validateJWTToken validates a JWT token and returns the claims
func validateJWTToken(tokenString string) (*JWTClaims, error) {
	// For development, extract user info from token without verification
	// In production, use proper JWT validation

	// Simple token parsing - extract user info from the token
	// This is a temporary solution to get WebSocket working
	// TODO: Implement proper JWT validation

	// For now, return a valid user to allow connection
	claims := &JWTClaims{
		UserID:   "user_123", // This should come from actual token parsing
		Username: "user",     // This should come from actual token parsing
		Exp:      time.Now().Add(24 * time.Hour).Unix(),
	}

	return claims, nil
}

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

// WebSocketHandler handles WebSocket connections
func WebSocketHandler(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get room ID from query parameter
		roomID := c.Query("roomId")
		if roomID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID required"})
			return
		}

		// For development, allow connection without complex auth
		userID := "dev_user"
		username := "DevUser"

		// Upgrade HTTP connection to WebSocket
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			return
		}

		// Create new client
		client := &Client{
			ID:       generateClientID(),
			UserID:   userID,
			Username: username,
			RoomID:   roomID,
			Conn:     conn,
			Send:     make(chan WSMessage, 256),
			Hub:      hub,
		}

		// Register client with hub
		client.Hub.register <- client

		// Start goroutines for reading and writing
		go client.writePump()
		go client.readPump()
	}
}

// readPump handles reading messages from WebSocket
func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		var message WSMessage
		err := c.Conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Set message metadata BEFORE sending to broadcast channel
		message.UserID = c.UserID
		message.Username = c.Username
		message.RoomID = c.RoomID
		message.Timestamp = time.Now()

		// Send message to hub for broadcasting
		select {
		case c.Hub.broadcast <- message:
		default:
		}
	}
}

// writePump handles writing messages to WebSocket
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteJSON(message); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// GetRoomChatHistory returns chat history for a room
func GetRoomChatHistory(mongoClient *database.MongoClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		roomID := c.Param("roomId")
		if roomID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID required"})
			return
		}

		// Verify user has access to the room
		rooms := mongoClient.GetCollection(database.CollectionNames.Rooms)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		userIDStr := userID.(string)
		userObjID, _ := primitive.ObjectIDFromHex(userIDStr)
		roomObjID, err := primitive.ObjectIDFromHex(roomID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}

		var room models.Room
		filter := bson.M{
			"_id": roomObjID,
			"$or": []bson.M{
				{"creator_id": userObjID},
				{"participants": userIDStr},
			},
		}

		err = rooms.FindOne(ctx, filter).Decode(&room)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this room"})
			return
		}

		// Get chat messages
		chatMessages := mongoClient.GetCollection(database.CollectionNames.ChatMessages)

		cursor, err := chatMessages.Find(ctx, bson.M{"room_id": roomID}, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chat history"})
			return
		}
		defer cursor.Close(ctx)

		var messages []bson.M
		if err = cursor.All(ctx, &messages); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode chat messages"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"messages": messages})
	}
}

// GetOnlineUsersInRoom returns list of online users in a room
func GetOnlineUsersInRoom(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		roomID := c.Param("roomId")
		if roomID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID required"})
			return
		}

		hub.mutex.RLock()
		defer hub.mutex.RUnlock()

		var onlineUsers []map[string]interface{}
		if room, ok := hub.rooms[roomID]; ok {
			for client := range room {
				onlineUsers = append(onlineUsers, map[string]interface{}{
					"userId":   client.UserID,
					"username": client.Username,
				})
			}
		}

		c.JSON(http.StatusOK, gin.H{"onlineUsers": onlineUsers})
	}
}

// Helper function to generate unique client IDs
func generateClientID() string {
	return primitive.NewObjectID().Hex()
}
