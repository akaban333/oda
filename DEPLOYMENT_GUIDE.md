# 🚀 Study Platform Deployment Guide

This guide covers deploying your Study Platform application with the new monitoring, rate limiting, and CI/CD features.

## 📋 Prerequisites

- Go 1.21+ installed
- Node.js 18+ installed
- MongoDB instance (local or cloud)
- Git repository with GitHub Actions enabled
- Docker (optional, for containerized deployment)

## 🏗️ Backend Deployment

### 1. Environment Configuration

Create environment files for different environments:

```bash
# .env.development
ENV=development
GATEWAY_PORT=8080
MONGODB_URI=mongodb://localhost:27017/studyplatform_dev
LOG_LEVEL=debug
JWT_SECRET=your-dev-jwt-secret
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173

# .env.production
ENV=production
GATEWAY_PORT=8080
MONGODB_URI=mongodb://your-production-mongo/studyplatform_prod
LOG_LEVEL=info
JWT_SECRET=your-production-jwt-secret
ALLOWED_ORIGINS=https://yourdomain.com
```

### 2. Build and Run

```bash
cd backend

# Install dependencies
go mod download

# Build the application
go build -o bin/api ./cmd/api

# Run the application
./bin/api
```

### 3. Docker Deployment (Optional)

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

```bash
# Build and run with Docker
docker build -t studyplatform-backend .
docker run -p 8080:8080 --env-file .env.production studyplatform-backend
```

## 🌐 Frontend Deployment

### 1. Build for Production

```bash
# Install dependencies
npm ci

# Build for production
npm run build

# The build folder contains your production-ready app
```

### 2. Deploy to Vercel (Recommended)

```bash
# Install Vercel CLI
npm i -g vercel

# Deploy
vercel --prod
```

### 3. Deploy to Netlify

```bash
# Install Netlify CLI
npm install -g netlify-cli

# Deploy
netlify deploy --prod --dir=build
```

### 4. Static Hosting (Nginx/Apache)

Copy the `build` folder to your web server directory and configure your web server.

## 📊 Monitoring & Health Checks

### 1. Health Check Endpoints

Your application now includes comprehensive health monitoring:

- **`/api/v1/health`** - Detailed system health status
- **`/api/v1/health/simple`** - Simple health check for load balancers
- **`/api/v1/metrics`** - Prometheus-compatible metrics
- **`/api/v1/admin/rate-limit-stats`** - Rate limiting statistics

### 2. Health Check Response Example

```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "uptime": "2h15m30s",
  "version": "1.0.0",
  "environment": "production",
  "services": {
    "database": {
      "status": "healthy",
      "message": "Database is responding normally",
      "last_check": "2024-01-15T10:30:00Z"
    },
    "memory": {
      "status": "healthy",
      "message": "Memory usage is normal",
      "last_check": "2024-01-15T10:30:00Z"
    }
  },
  "metrics": {
    "memory": {
      "alloc": 12345678,
      "total_alloc": 98765432,
      "sys": 45678901
    },
    "database": {
      "connection_status": "connected",
      "response_time": "2.5ms",
      "total_users": 150,
      "total_rooms": 25
    }
  }
}
```

### 3. Prometheus Metrics

The `/metrics` endpoint provides Prometheus-compatible metrics:

```prometheus
# HELP go_memory_alloc_bytes Current memory usage in bytes
# TYPE go_memory_alloc_bytes gauge
go_memory_alloc_bytes 12345678

# HELP go_database_users_total Total number of users
# TYPE go_database_users_total gauge
go_database_users_total 150
```

## 🛡️ Rate Limiting

### 1. Configuration

Rate limiting is automatically applied to all routes with these defaults:

- **60 requests per minute** per IP address
- **10 burst requests** allowed
- **5-minute blocking** for excessive requests

### 2. Custom Configuration

```go
config := &middleware.RateLimitConfig{
    RequestsPerMinute: 100,
    BurstSize:         20,
    WindowSize:        time.Minute,
}
rateLimiter := middleware.NewRateLimiter(config)
```

### 3. Rate Limit Response

When rate limited, clients receive:

```json
{
  "error": "Rate limit exceeded",
  "message": "Too many requests. Please try again later.",
  "retry_after": 180
}
```

## 🔄 CI/CD Pipeline

### 1. GitHub Actions Setup

The pipeline automatically runs on:
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop` branches

### 2. Pipeline Stages

1. **Backend Testing** - Go tests with MongoDB
2. **Frontend Testing** - React tests and build
3. **Security Scanning** - Trivy vulnerability scanner
4. **Code Quality** - Linting and formatting checks
5. **Deployment** - Staging (develop) and Production (main)

### 3. Environment Protection

- Production deployment requires manual approval
- All tests must pass before deployment
- Security scans are mandatory

## 🚨 Error Tracking & Logging

### 1. Structured Logging

All logs include structured fields:

```go
logger.Info("User login", 
    logger.Field("user_id", userID),
    logger.Field("ip_address", clientIP),
    logger.Field("user_agent", userAgent),
)
```

### 2. Error Aggregation

Errors are automatically categorized and tracked:

- **Database errors** - Connection issues, query failures
- **Authentication errors** - Invalid tokens, expired sessions
- **Validation errors** - Invalid input data
- **System errors** - Memory issues, goroutine leaks

### 3. Error Monitoring

Access error statistics via the error tracker:

```go
stats := errorTracker.GetErrorStats()
unresolved := errorTracker.GetUnresolvedErrors()
```

## 🔧 Production Checklist

### 1. Security

- [ ] Use strong JWT secrets
- [ ] Enable HTTPS/TLS
- [ ] Configure CORS properly
- [ ] Set up rate limiting
- [ ] Enable security headers

### 2. Monitoring

- [ ] Set up health check monitoring
- [ ] Configure Prometheus metrics collection
- [ ] Set up alerting for critical issues
- [ ] Monitor rate limiting statistics
- [ ] Track error rates and patterns

### 3. Performance

- [ ] Enable connection pooling for MongoDB
- [ ] Configure appropriate timeouts
- [ ] Set up caching where appropriate
- [ ] Monitor memory and CPU usage
- [ ] Profile database queries

### 4. Reliability

- [ ] Set up database backups
- [ ] Configure graceful shutdown
- [ ] Implement retry mechanisms
- [ ] Set up load balancing
- [ ] Plan for disaster recovery

## 📈 Scaling Considerations

### 1. Horizontal Scaling

- Use multiple backend instances behind a load balancer
- Ensure MongoDB is configured for replication
- Implement sticky sessions for WebSocket connections

### 2. Database Scaling

- Consider MongoDB sharding for large datasets
- Implement read replicas for read-heavy workloads
- Use connection pooling effectively

### 3. Caching Strategy

- Implement Redis for session storage
- Cache frequently accessed data
- Use CDN for static assets

## 🆘 Troubleshooting

### 1. Common Issues

**Health checks failing:**
- Check MongoDB connection
- Verify environment variables
- Check log files for errors

**Rate limiting too aggressive:**
- Adjust rate limit configuration
- Check for legitimate high-traffic scenarios
- Monitor rate limit statistics

**High memory usage:**
- Check for memory leaks in Go code
- Monitor goroutine count
- Review database query patterns

### 2. Debug Endpoints

- **`/api/v1/friends/debug`** - Database connection test
- **`/api/v1/admin/rate-limit-stats`** - Rate limiting status
- **`/api/v1/health`** - Comprehensive system status

### 3. Log Analysis

```bash
# View recent logs
tail -f backend.log

# Search for errors
grep "ERROR" backend.log

# Monitor specific user activity
grep "user_id:12345" backend.log
```

## 📚 Additional Resources

- [Go Best Practices](https://golang.org/doc/effective_go.html)
- [MongoDB Production Notes](https://docs.mongodb.com/manual/core/security-checklist/)
- [React Production Build](https://create-react-app.dev/docs/production-build/)
- [Prometheus Monitoring](https://prometheus.io/docs/introduction/overview/)

## 🆘 Support

For deployment issues:
1. Check the logs first
2. Verify environment configuration
3. Test health check endpoints
4. Review monitoring metrics
5. Check CI/CD pipeline status

---

**Happy Deploying! 🎉** 