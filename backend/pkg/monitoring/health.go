package monitoring

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/studyplatform/backend/pkg/database"
	"github.com/studyplatform/backend/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
)

// SystemStatus represents the overall system health
type SystemStatus struct {
	Status      string            `json:"status"`
	Timestamp   time.Time         `json:"timestamp"`
	Uptime      string            `json:"uptime"`
	Version     string            `json:"version"`
	Environment string            `json:"environment"`
	Services    map[string]Status `json:"services"`
	Metrics     SystemMetrics     `json:"metrics"`
}

// Status represents individual service status
type Status struct {
	Status    string                 `json:"status"`
	Message   string                 `json:"message,omitempty"`
	Details   map[string]interface{} `json:"details,omitempty"`
	LastCheck time.Time              `json:"last_check"`
}

// SystemMetrics represents system performance metrics
type SystemMetrics struct {
	Memory     MemoryMetrics   `json:"memory"`
	CPU        CPUMetrics      `json:"cpu"`
	Goroutines int             `json:"goroutines"`
	Database   DatabaseMetrics `json:"database"`
}

// MemoryMetrics represents memory usage
type MemoryMetrics struct {
	Alloc      uint64 `json:"alloc"`
	TotalAlloc uint64 `json:"total_alloc"`
	Sys        uint64 `json:"sys"`
	NumGC      uint32 `json:"num_gc"`
	HeapAlloc  uint64 `json:"heap_alloc"`
	HeapSys    uint64 `json:"heap_sys"`
}

// CPUMetrics represents CPU usage
type CPUMetrics struct {
	NumCPU int `json:"num_cpu"`
}

// DatabaseMetrics represents database performance
type DatabaseMetrics struct {
	ConnectionStatus string        `json:"connection_status"`
	ResponseTime     time.Duration `json:"response_time"`
	TotalUsers       int64         `json:"total_users"`
	TotalRooms       int64         `json:"total_rooms"`
	TotalSessions    int64         `json:"total_sessions"`
}

// HealthChecker manages health checks
type HealthChecker struct {
	startTime time.Time
	version   string
	env       string
	mongo     *database.MongoClient
	mu        sync.RWMutex
	services  map[string]Status
}

// NewHealthChecker creates a new health checker instance
func NewHealthChecker(mongo *database.MongoClient, version, env string) *HealthChecker {
	return &HealthChecker{
		startTime: time.Now(),
		version:   version,
		env:       env,
		mongo:     mongo,
		services:  make(map[string]Status),
	}
}

// GetSystemStatus returns comprehensive system status
func (hc *HealthChecker) GetSystemStatus() SystemStatus {
	hc.mu.RLock()
	defer hc.mu.RUnlock()

	// Collect current metrics
	memStats := runtime.MemStats{}
	runtime.ReadMemStats(&memStats)

	// Check database status
	dbStatus := hc.checkDatabaseStatus()

	// Determine overall status
	overallStatus := "healthy"
	for _, service := range hc.services {
		if service.Status == "unhealthy" {
			overallStatus = "degraded"
		}
	}

	if dbStatus.ConnectionStatus == "disconnected" {
		overallStatus = "unhealthy"
	}

	return SystemStatus{
		Status:      overallStatus,
		Timestamp:   time.Now(),
		Uptime:      time.Since(hc.startTime).String(),
		Version:     hc.version,
		Environment: hc.env,
		Services:    hc.services,
		Metrics: SystemMetrics{
			Memory: MemoryMetrics{
				Alloc:      memStats.Alloc,
				TotalAlloc: memStats.TotalAlloc,
				Sys:        memStats.Sys,
				NumGC:      memStats.NumGC,
				HeapAlloc:  memStats.HeapAlloc,
				HeapSys:    memStats.HeapSys,
			},
			CPU: CPUMetrics{
				NumCPU: runtime.NumCPU(),
			},
			Goroutines: runtime.NumGoroutine(),
			Database:   dbStatus,
		},
	}
}

// checkDatabaseStatus checks MongoDB connection and performance
func (hc *HealthChecker) checkDatabaseStatus() DatabaseMetrics {
	start := time.Now()

	// Test database connection
	users := hc.mongo.GetCollection(database.CollectionNames.Users)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Count users
	userCount, err := users.CountDocuments(ctx, bson.M{})
	if err != nil {
		return DatabaseMetrics{
			ConnectionStatus: "disconnected",
			ResponseTime:     time.Since(start),
			TotalUsers:       0,
			TotalRooms:       0,
			TotalSessions:    0,
		}
	}

	// Count rooms
	rooms := hc.mongo.GetCollection(database.CollectionNames.Rooms)
	roomCount, _ := rooms.CountDocuments(ctx, bson.M{})

	// Count sessions
	sessions := hc.mongo.GetCollection(database.CollectionNames.Sessions)
	sessionCount, _ := sessions.CountDocuments(ctx, bson.M{})

	return DatabaseMetrics{
		ConnectionStatus: "connected",
		ResponseTime:     time.Since(start),
		TotalUsers:       userCount,
		TotalRooms:       roomCount,
		TotalSessions:    sessionCount,
	}
}

// UpdateServiceStatus updates the status of a specific service
func (hc *HealthChecker) UpdateServiceStatus(serviceName string, status Status) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	hc.services[serviceName] = status
}

// HealthCheckHandler returns a Gin handler for health checks
func (hc *HealthChecker) HealthCheckHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		status := hc.GetSystemStatus()
		c.JSON(http.StatusOK, status)
	}
}

// SimpleHealthCheckHandler returns a simple health check for load balancers
func (hc *HealthChecker) SimpleHealthCheckHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		status := hc.GetSystemStatus()
		if status.Status == "healthy" {
			c.JSON(http.StatusOK, gin.H{
				"status":    "ok",
				"timestamp": time.Now().Unix(),
			})
		} else {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":    "unhealthy",
				"timestamp": time.Now().Unix(),
			})
		}
	}
}

// MetricsHandler returns detailed metrics in Prometheus format
func (hc *HealthChecker) MetricsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		status := hc.GetSystemStatus()

		// Format as Prometheus metrics
		metrics := fmt.Sprintf(`# HELP go_memory_alloc_bytes Current memory usage in bytes
# TYPE go_memory_alloc_bytes gauge
go_memory_alloc_bytes %d

# HELP go_memory_total_alloc_bytes Total memory allocated in bytes
# TYPE go_memory_total_alloc_bytes counter
go_memory_total_alloc_bytes %d

# HELP go_goroutines Number of goroutines
# TYPE go_goroutines gauge
go_goroutines %d

# HELP go_database_users_total Total number of users
# TYPE go_database_users_total gauge
go_database_users_total %d

# HELP go_database_rooms_total Total number of rooms
# TYPE go_database_rooms_total gauge
go_database_rooms_total %d

# HELP go_database_sessions_total Total number of sessions
# TYPE go_database_sessions_total gauge
go_database_sessions_total %d

# HELP go_uptime_seconds Uptime in seconds
# TYPE go_uptime_seconds counter
go_uptime_seconds %f
`,
			status.Metrics.Memory.Alloc,
			status.Metrics.Memory.TotalAlloc,
			status.Metrics.Goroutines,
			status.Metrics.Database.TotalUsers,
			status.Metrics.Database.TotalRooms,
			status.Metrics.Database.TotalSessions,
			time.Since(hc.startTime).Seconds(),
		)

		c.Header("Content-Type", "text/plain")
		c.String(http.StatusOK, metrics)
	}
}

// StartPeriodicChecks starts background health monitoring
func (hc *HealthChecker) StartPeriodicChecks() {
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for range ticker.C {
			hc.performPeriodicChecks()
		}
	}()
}

// performPeriodicChecks runs periodic health checks
func (hc *HealthChecker) performPeriodicChecks() {
	// Check database
	dbStatus := hc.checkDatabaseStatus()
	status := "healthy"
	message := "Database is responding normally"

	if dbStatus.ConnectionStatus == "disconnected" {
		status = "unhealthy"
		message = "Database connection failed"
		logger.Error("Database health check failed", logger.Field("error", "connection timeout"))
	}

	hc.UpdateServiceStatus("database", Status{
		Status:    status,
		Message:   message,
		LastCheck: time.Now(),
		Details: map[string]interface{}{
			"response_time": dbStatus.ResponseTime.String(),
			"total_users":   dbStatus.TotalUsers,
			"total_rooms":   dbStatus.TotalRooms,
		},
	})

	// Check memory usage
	memStats := runtime.MemStats{}
	runtime.ReadMemStats(&memStats)

	memStatus := "healthy"
	memMessage := "Memory usage is normal"

	// Alert if memory usage is high (over 80% of system memory)
	if memStats.Sys > 0 && float64(memStats.Alloc)/float64(memStats.Sys) > 0.8 {
		memStatus = "warning"
		memMessage = "High memory usage detected"
		logger.Warn("High memory usage detected",
			logger.Field("alloc", memStats.Alloc),
			logger.Field("sys", memStats.Sys),
			logger.Field("percentage", float64(memStats.Alloc)/float64(memStats.Sys)*100),
		)
	}

	hc.UpdateServiceStatus("memory", Status{
		Status:    memStatus,
		Message:   memMessage,
		LastCheck: time.Now(),
		Details: map[string]interface{}{
			"alloc_bytes":       memStats.Alloc,
			"total_alloc_bytes": memStats.TotalAlloc,
			"sys_bytes":         memStats.Sys,
			"heap_alloc_bytes":  memStats.HeapAlloc,
		},
	})
}
