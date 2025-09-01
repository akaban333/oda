# Study Platform Backend API Documentation

## Overview
This document provides comprehensive documentation for all available API endpoints in the Study Platform backend.

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
Most endpoints require JWT authentication. Include the JWT token in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

---

## Authentication Endpoints

### Register User
- **POST** `/auth/register`
- **Description**: Register a new user account
- **Body**:
  ```json
  {
    "username": "johndoe",
    "email": "john@example.com",
    "password": "securepassword",
    "firstName": "John",
    "lastName": "Doe"
  }
  ```

### Login
- **POST** `/auth/login`
- **Description**: Login user and get JWT tokens
- **Body**:
  ```json
  {
    "email": "john@example.com",
    "password": "securepassword"
  }
  ```

### Logout
- **POST** `/auth/logout`
- **Description**: Logout user and invalidate tokens
- **Headers**: Authorization required

### Refresh Token
- **POST** `/auth/refresh`
- **Description**: Refresh JWT access token

### Get Current User
- **GET** `/auth/me`
- **Description**: Get current authenticated user's profile
- **Headers**: Authorization required

### Update Profile
- **PUT** `/auth/me`
- **Description**: Update user profile
- **Headers**: Authorization required
- **Body**:
  ```json
  {
    "username": "newusername",
    "firstName": "John",
    "lastName": "Doe",
    "bio": "Updated bio",
    "avatarUrl": "https://example.com/avatar.jpg"
  }
  ```

### Update XP
- **PUT** `/auth/xp`
- **Description**: Update user XP (admin only)
- **Headers**: Authorization required

---

## Friends Management

### Send Friend Request
- **POST** `/friends/request`
- **Description**: Send a friend request using unique ID
- **Headers**: Authorization required
- **Body**:
  ```json
  {
    "uniqueId": "USER123456"
  }
  ```

### Accept Friend Request
- **PUT** `/friends/:id/accept`
- **Description**: Accept a friend request
- **Headers**: Authorization required

### Reject Friend Request
- **PUT** `/friends/:id/reject`
- **Description**: Reject a friend request
- **Headers**: Authorization required

### List Friends
- **GET** `/friends/`
- **Description**: Get list of friends
- **Headers**: Authorization required

### List Friend Requests
- **GET** `/friends/requests`
- **Description**: Get pending friend requests
- **Headers**: Authorization required

---

## User Search

### Search Users
- **GET** `/users/`
- **Description**: Search for users
- **Headers**: Authorization required
- **Query Parameters**: `q` (search query)

### Get User by ID
- **GET** `/users/:id`
- **Description**: Get user profile by ID
- **Headers**: Authorization required

---

## Room Management

### List Rooms
- **GET** `/rooms/`
- **Description**: Get user's rooms
- **Headers**: Authorization required

### Create Room
- **POST** `/rooms/`
- **Description**: Create a new study room
- **Headers**: Authorization required
- **Body**:
  ```json
  {
    "name": "Biology Study Group",
    "description": "Studying for final exams",
    "privacy": "private",
    "maxParticipants": 5,
    "tags": ["biology", "exam"]
  }
  ```

### Get Room Details
- **GET** `/rooms/:id`
- **Description**: Get specific room details
- **Headers**: Authorization required

### Update Room
- **PUT** `/rooms/:id`
- **Description**: Update room settings (owner only)
- **Headers**: Authorization required

### Delete Room
- **DELETE** `/rooms/:id`
- **Description**: Delete room (owner only)
- **Headers**: Authorization required

### Join Room by Invitation Code
- **POST** `/rooms/join`
- **Description**: Join a room using invitation code
- **Headers**: Authorization required
- **Body**:
  ```json
  {
    "invitationCode": "ABC12345"
  }
  ```

### Generate Invitation Code
- **POST** `/rooms/:id/generate-code`
- **Description**: Generate new invitation code for room
- **Headers**: Authorization required

### Leave Room
- **POST** `/rooms/:id/leave`
- **Description**: Leave a room
- **Headers**: Authorization required

---

## Session Management

### Start Session
- **POST** `/sessions/start`
- **Description**: Start a new study session
- **Headers**: Authorization required
- **Body**:
  ```json
  {
    "roomId": "room_id_here"
  }
  ```

### End Session
- **POST** `/sessions/end`
- **Description**: End a study session with XP calculation
- **Headers**: Authorization required
- **Body**:
  ```json
  {
    "sessionId": "session_id_here",
    "duration": 3600,
    "inactiveDuration": 300,
    "pomodoroCount": 2,
    "activityData": []
  }
  ```

### List Sessions
- **GET** `/sessions/`
- **Description**: Get user's session history
- **Headers**: Authorization required

### Activity Ping
- **POST** `/sessions/:id/ping`
- **Description**: Send activity ping to prevent session timeout
- **Headers**: Authorization required
- **Body**:
  ```json
  {
    "activityType": "keyboard",
    "timestamp": "2023-01-01T10:00:00Z"
  }
  ```

### Session Statistics
- **GET** `/sessions/stats`
- **Description**: Get comprehensive session statistics
- **Headers**: Authorization required

### Check XP Privileges
- **GET** `/sessions/privileges`
- **Description**: Check user privileges based on XP
- **Headers**: Authorization required
- **Query Parameters**: `action` (create_shared_room, add_participant)

---

## Materials Management

### List Materials
- **GET** `/materials/`
- **Description**: Get user's materials
- **Headers**: Authorization required

### Create Material
- **POST** `/materials/`
- **Description**: Create material record (legacy method)
- **Headers**: Authorization required

### Get Material
- **GET** `/materials/:id`
- **Description**: Get specific material
- **Headers**: Authorization required

### Update Material
- **PUT** `/materials/:id`
- **Description**: Update material metadata
- **Headers**: Authorization required

### Delete Material
- **DELETE** `/materials/:id`
- **Description**: Delete material
- **Headers**: Authorization required

### Generate Upload URL
- **POST** `/materials/upload-url`
- **Description**: Generate presigned URL for file upload
- **Headers**: Authorization required
- **Body**:
  ```json
  {
    "fileName": "document.pdf",
    "fileType": "pdf",
    "contentType": "application/pdf"
  }
  ```

### Confirm Upload
- **POST** `/materials/confirm-upload`
- **Description**: Confirm file upload and create material record
- **Headers**: Authorization required
- **Body**:
  ```json
  {
    "objectName": "timestamp_document.pdf",
    "name": "Biology Notes",
    "description": "Chapter 5 summary",
    "fileType": "pdf",
    "fileSize": 1024000,
    "roomId": "room_id_optional"
  }
  ```

### Share Material
- **POST** `/materials/:id/share`
- **Description**: Share material with another user
- **Headers**: Authorization required
- **Body**:
  ```json
  {
    "userId": "user_id_here",
    "permission": "view"
  }
  ```

---

## Todo Management

### List Todos
- **GET** `/todos/`
- **Description**: Get user's todos
- **Headers**: Authorization required

### Create Todo
- **POST** `/todos/`
- **Description**: Create a new todo
- **Headers**: Authorization required

### Get Todo
- **GET** `/todos/:id`
- **Description**: Get specific todo
- **Headers**: Authorization required

### Update Todo
- **PUT** `/todos/:id`
- **Description**: Update todo
- **Headers**: Authorization required

### Complete Todo
- **PUT** `/todos/:id/complete`
- **Description**: Mark todo as complete
- **Headers**: Authorization required

### Delete Todo
- **DELETE** `/todos/:id`
- **Description**: Delete todo
- **Headers**: Authorization required

---

## Notes Management

### List Notes
- **GET** `/notes/`
- **Description**: Get user's notes
- **Headers**: Authorization required

### Create Note
- **POST** `/notes/`
- **Description**: Create a new note
- **Headers**: Authorization required

### Get Note
- **GET** `/notes/:id`
- **Description**: Get specific note
- **Headers**: Authorization required

### Update Note
- **PUT** `/notes/:id`
- **Description**: Update note
- **Headers**: Authorization required

### Delete Note
- **DELETE** `/notes/:id`
- **Description**: Delete note
- **Headers**: Authorization required

---

## Social Features (Posts)

### List Posts
- **GET** `/posts/`
- **Description**: Get social feed posts
- **Headers**: Authorization required

### Create Post
- **POST** `/posts/`
- **Description**: Create a new post
- **Headers**: Authorization required
- **Body**:
  ```json
  {
    "content": "This is my post content"
  }
  ```

### Like Post
- **PUT** `/posts/:id/like`
- **Description**: Like/unlike a post
- **Headers**: Authorization required

### Delete Post
- **DELETE** `/posts/:id`
- **Description**: Delete own post
- **Headers**: Authorization required

### Create Comment
- **POST** `/posts/:postId/comments`
- **Description**: Add comment to post
- **Headers**: Authorization required

### Like Comment
- **PUT** `/posts/comments/:commentId/like`
- **Description**: Like/unlike a comment
- **Headers**: Authorization required

---

## Notifications

### List Notifications
- **GET** `/notifications/`
- **Description**: Get user notifications
- **Headers**: Authorization required

### Mark Notification as Read (Deletes it)
- **PUT** `/notifications/:id/read`
- **Description**: Mark notification as read and delete it
- **Headers**: Authorization required

### Delete Notification
- **DELETE** `/notifications/:id`
- **Description**: Delete a specific notification
- **Headers**: Authorization required

### Create Notification
- **POST** `/notifications/`
- **Description**: Create notification (system use)
- **Headers**: Authorization required

---

## Real-time Communication

### WebSocket Connection
- **GET** `/realtime/ws`
- **Description**: Establish WebSocket connection for real-time features
- **Headers**: Authorization required
- **Query Parameters**: `roomId` (required)

### Chat History
- **GET** `/realtime/chat/:roomId`
- **Description**: Get chat message history for a room
- **Headers**: Authorization required

### Online Users
- **GET** `/realtime/online/:roomId`
- **Description**: Get list of online users in a room
- **Headers**: Authorization required

---

## WebSocket Message Types

When connected to WebSocket, you can send/receive these message types:

### Chat Messages
```json
{
  "type": "chat",
  "content": "Hello everyone!",
  "roomId": "room_id_here"
}
```

### Typing Indicators
```json
{
  "type": "typing",
  "roomId": "room_id_here"
}
```

### Video Call Signaling
```json
{
  "type": "rtc_offer",
  "roomId": "room_id_here",
  "data": {
    "targetUserId": "user_id",
    "offer": "webrtc_offer_data"
  }
}
```

---

## XP System Rules

### XP Earning
- **2 XP per minute** of active study time
- **30 XP bonus** per completed Pomodoro timer
- **-1 XP penalty** per 5 minutes of inactivity

### XP Privileges
- **300 XP**: Can add 1 additional participant to rooms (max 6 total)
- **600 XP**: Can add 2 additional participants (max 7 total)
- **1000 XP**: Can create shared rooms
- **2000 XP**: Can create additional shared rooms

### Activity Detection
- Sessions automatically end after 10 minutes of inactivity
- Use `/sessions/:id/ping` to maintain active status

---

## Error Responses

All endpoints return consistent error responses:

```json
{
  "error": "Error message",
  "details": "Additional error details (optional)"
}
```

Common HTTP status codes:
- **200**: Success
- **201**: Created
- **400**: Bad Request
- **401**: Unauthorized
- **403**: Forbidden
- **404**: Not Found
- **409**: Conflict
- **500**: Internal Server Error 