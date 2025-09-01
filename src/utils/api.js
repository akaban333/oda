// API utility functions for making authenticated requests to the backend
import logger from './logger';

// Base API configuration - will be overridden by environment variables in production
const API_BASE_URL = process.env.REACT_APP_API_BASE_URL || 'http://localhost:8080/api/v1';

// Helper function to get auth headers
const getAuthHeaders = () => {
  const token = localStorage.getItem('accessToken');
  return {
    'Content-Type': 'application/json',
    ...(token && { 'Authorization': `Bearer ${token}` })
  };
};

// Generic API request function
const apiRequest = async (endpoint, options = {}) => {
  const url = `${API_BASE_URL}${endpoint}`;
  const config = {
    headers: getAuthHeaders(),
    ...options
  };

  logger.log('Making API request:', {
    url,
    method: options.method || 'GET',
    headers: config.headers,
    body: options.body
  });
  
  // Add extra logging for DELETE requests
  if (options.method === 'DELETE') {
    logger.log('DELETE request details:', { url, method: options.method });
  }

  try {
    const response = await fetch(url, config);
    
    logger.log('API response:', {
      status: response.status,
      statusText: response.statusText,
      headers: Object.fromEntries(response.headers.entries())
    });
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      console.error('API error response:', errorData);
      throw new Error(errorData.error || `HTTP ${response.status}: ${response.statusText}`);
    }
    
    const data = await response.json();
    logger.log('API success response:', data);
    return data;
  } catch (error) {
    console.error('API request failed:', error);
    throw error;
  }
};

// Auth API functions
export const authAPI = {
  // Get current user profile
  getProfile: () => apiRequest('/auth/me'),
  
  // Update user profile
  updateProfile: (data) => apiRequest('/auth/me', {
    method: 'PUT',
    body: JSON.stringify(data)
  }),
  
  // Update user XP
  updateXP: (data) => apiRequest('/auth/xp', {
    method: 'PUT',
    body: JSON.stringify(data)
  }),
  
  // Logout
  logout: () => apiRequest('/auth/logout', { method: 'POST' }),
  
  // Refresh token
  refreshToken: () => apiRequest('/auth/refresh', { method: 'POST' })
};

// Rooms API functions
export const roomsAPI = {
  // Get user's rooms
  getRooms: () => apiRequest('/rooms/'),
  
  // Create new room
  createRoom: (data) => apiRequest('/rooms/', {
    method: 'POST',
    body: JSON.stringify(data)
  }),
  
  // Get room by ID
  getRoom: (id) => apiRequest(`/rooms/${id}`),
  
  // Update room
  updateRoom: (id, data) => apiRequest(`/rooms/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data)
  }),
  
  // Delete room
  deleteRoom: (id) => apiRequest(`/rooms/${id}`, {
    method: 'DELETE'
  }),
  
  // Join room by invitation code
  joinRoomByCode: (data) => apiRequest('/rooms/join', {
    method: 'POST',
    body: JSON.stringify(data)
  }),
  
  // Generate invitation code for room
  generateInvitationCode: (roomId) => apiRequest(`/rooms/${roomId}/generate-code`, {
    method: 'POST'
  }),
  
  // Invite user to room
  inviteUserToRoom: (roomId, targetUserId) => apiRequest(`/rooms/${roomId}/invite`, {
    method: 'POST',
    body: JSON.stringify({ targetUserId })
  }),
  
  // Accept room invitation
  acceptRoomInvitation: (roomId) => apiRequest(`/rooms/${roomId}/accept`, {
    method: 'POST'
  }),
  
  // Leave room
  leaveRoom: (roomId) => apiRequest(`/rooms/${roomId}/leave`, {
    method: 'POST'
  }),
  
  // Enter room (add user as participant)
  enterRoom: (roomId) => apiRequest(`/rooms/${roomId}/enter`, {
    method: 'POST'
  })
};

// Sessions API functions
export const sessionsAPI = {
  // Start session
  startSession: (data) => apiRequest('/sessions/start', {
    method: 'POST',
    body: JSON.stringify(data)
  }),
  
  // End session
  endSession: (id, data) => apiRequest(`/sessions/${id}/end`, {
    method: 'PUT',
    body: JSON.stringify(data)
  }),
  
  // Get session history
  getSessions: () => apiRequest('/sessions/'),
  
  // Get current session
  getCurrentSession: () => apiRequest('/sessions/current'),
  
  // Get session statistics
  getStats: () => apiRequest('/sessions/stats'),
  
  // Check XP privileges
  getPrivileges: () => apiRequest('/sessions/privileges?action=add_participant')
};

// Materials API functions
export const materialsAPI = {
  // Get user's materials
  getMaterials: () => apiRequest('/materials/'),
  
  // Get room materials
  getRoomMaterials: (roomId) => apiRequest(`/rooms/${roomId}/materials/`),
  
  // Create material
  createMaterial: (data) => apiRequest('/materials/', {
    method: 'POST',
    body: JSON.stringify(data)
  }),
  
  // Get specific material
  getMaterial: (id) => apiRequest(`/materials/${id}`),
  
  // Update material
  updateMaterial: (id, data) => apiRequest(`/materials/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data)
  }),
  
  // Delete material
  deleteMaterial: (id) => apiRequest(`/materials/${id}`, { method: 'DELETE' })
};

// Todos API functions
export const todosAPI = {
  // Get all todos
  getTodos: () => apiRequest('/todos/'),
  
  // Create new todo
  createTodo: (data) => apiRequest('/todos/', {
    method: 'POST',
    body: JSON.stringify(data)
  }),
  
  // Complete todo
  completeTodo: (todoId) => apiRequest(`/todos/${todoId}/complete`, {
    method: 'PUT'
  }),
  
  // Delete todo
  deleteTodo: (todoId) => apiRequest(`/todos/${todoId}`, {
    method: 'DELETE'
  }),
  
  // Update todo
  updateTodo: (todoId, data) => apiRequest(`/todos/${todoId}`, {
    method: 'PUT',
    body: JSON.stringify(data)
  })
};

// Notes API functions
export const notesAPI = {
  // Get user's notes
  getNotes: () => apiRequest('/notes/'),
  
  // Get room notes
  getRoomNotes: (roomId) => apiRequest(`/rooms/${roomId}/notes/`),
  
  // Create note
  createNote: (data) => apiRequest('/notes/', {
    method: 'POST',
    body: JSON.stringify(data)
  }),
  
  // Get specific note
  getNote: (id) => apiRequest(`/notes/${id}`),
  
  // Update note
  updateNote: (id, data) => apiRequest(`/notes/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data)
  }),
  
  // Delete note
  deleteNote: (id) => apiRequest(`/notes/${id}`, { method: 'DELETE' })
};

// Friends API functions
export const friendsAPI = {
  // Get friends list
  getFriends: () => apiRequest('/friends/'),
  
  // Send friend request
  sendFriendRequest: (data) => apiRequest('/friends/request', {
    method: 'POST',
    body: JSON.stringify(data)
  }),
  
  // Accept friend request
  acceptFriendRequest: (id) => apiRequest(`/friends/${id}/accept`, { method: 'PUT' }),
  
  // Reject friend request
  rejectFriendRequest: (id) => apiRequest(`/friends/${id}/reject`, { method: 'PUT' }),
  
  // Remove friend
  removeFriend: (id) => apiRequest(`/friends/${id}/remove`, { method: 'DELETE' })
};

// Users API functions
export const usersAPI = {
  // Search users
  searchUsers: (query) => apiRequest(`/users/search?search=${encodeURIComponent(query)}`),
  
  // Get user by ID
  getUser: (id) => apiRequest(`/users/${id}`)
};

// Posts API functions
export const postsAPI = {
  // Get all posts (for social feed)
  getPosts: () => apiRequest('/posts/'),
  
  // Get user's own posts
  getMyPosts: () => apiRequest('/posts/?mine=true'),
  
  // Create new post
  createPost: (data) => apiRequest('/posts/', {
    method: 'POST',
    body: JSON.stringify(data)
  }),
  
  // Delete post
  deletePost: (postId) => apiRequest(`/posts/${postId}`, {
    method: 'DELETE'
  }),
  
  // Like/unlike post
  likePost: (postId) => apiRequest(`/posts/${postId}/like`, {
    method: 'PUT',
    body: JSON.stringify({ postId: postId })
  }),
  
  // Create comment
  createComment: (postId, data) => apiRequest(`/posts/${postId}/comments`, {
    method: 'POST',
    body: JSON.stringify(data)
  }),
  
  // Like/unlike comment
  likeComment: (commentId) => apiRequest(`/posts/comments/${commentId}/like`, {
    method: 'PUT',
    body: JSON.stringify({ commentId: commentId })
  })
};

// Notifications API functions
export const notificationsAPI = {
  // Get user's notifications
  getNotifications: () => apiRequest('/notifications/'),
  
  // Mark notification as read (now deletes it)
  markAsRead: (id) => apiRequest(`/notifications/${id}/read`, {
    method: 'PUT'
  }),
  
  // Delete notification directly
  deleteNotification: (id) => apiRequest(`/notifications/${id}`, {
    method: 'DELETE'
  }),
  
  // Create notification (system use)
  createNotification: (data) => apiRequest('/notifications/', {
    method: 'POST',
    body: JSON.stringify(data)
  })
};

// Real-time API functions
export const realtimeAPI = {
  // Get chat history for a room
  getChatHistory: (roomId) => apiRequest(`/realtime/chat/${roomId}`),
  
  // Get online users in a room
  getOnlineUsers: (roomId) => apiRequest(`/realtime/online/${roomId}`),
  
  // Create WebSocket connection
  createWebSocket: () => {
    const token = localStorage.getItem('accessToken');
    const wsBaseUrl = process.env.REACT_APP_WS_BASE_URL || 'ws://localhost:8080/api/v1';
    const wsUrl = `${wsBaseUrl}/realtime/ws${token ? `?token=${token}` : ''}`;
    return new WebSocket(wsUrl);
  }
};

// Test function to check backend connectivity
export const testBackendConnectivity = async () => {
  try {
    logger.log('Testing backend connectivity...');
    
    // Test health endpoint
    const healthResponse = await fetch('http://localhost:8080/api/v1/health');
    logger.log('Health check response:', healthResponse.status, healthResponse.statusText);
    
    if (healthResponse.ok) {
      const healthData = await healthResponse.json();
      logger.log('Health check data:', healthData);
    }
    
    // Test simple test endpoint (no auth required)
    const simpleTestResponse = await fetch('http://localhost:8080/api/v1/simple-test');
    logger.log('Simple test response:', simpleTestResponse.status, simpleTestResponse.statusText);
    
    if (simpleTestResponse.ok) {
      const simpleTestData = await simpleTestResponse.json();
      logger.log('Simple test data:', simpleTestData);
    } else {
      const errorData = await simpleTestResponse.json().catch(() => ({}));
      logger.log('Simple test error:', errorData);
    }
    
    // Test debug rooms endpoint (no auth required)
    const debugResponse = await fetch('http://localhost:8080/api/v1/debug-rooms');
    logger.log('Debug rooms response:', debugResponse.status, debugResponse.statusText);
    
    if (debugResponse.ok) {
      const debugData = await debugResponse.json();
      logger.log('Debug rooms data:', debugData);
    } else {
      const errorData = await debugResponse.json().catch(() => ({}));
      logger.log('Debug rooms error:', errorData);
    }
    
    return true;
  } catch (error) {
    logger.error('Backend connectivity test failed:', error);
    return false;
  }
};

const api = {
  auth: authAPI,
  rooms: roomsAPI,
  sessions: sessionsAPI,
  materials: materialsAPI,
  todos: todosAPI,
  notes: notesAPI,
  friends: friendsAPI,
  users: usersAPI,
  posts: postsAPI,
  notifications: notificationsAPI,
  realtime: realtimeAPI
};

export default api; 