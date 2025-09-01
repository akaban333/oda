/* eslint-disable react-hooks/exhaustive-deps */
/* eslint-disable no-unused-vars */
import React, { useState, useEffect, useCallback, useMemo, useRef } from 'react';
import { Icon } from '@iconify/react';
import PropTypes from 'prop-types';
import { todosAPI, materialsAPI, notesAPI, authAPI, sessionsAPI, realtimeAPI, friendsAPI, roomsAPI, testBackendConnectivity } from '../utils/api';
import VideoCallRoom from './VideoCallRoom';

// Import new components
import RoomSidebar from './rooms/RoomSidebar';
import UserRoom from './rooms/UserRoom';
import SharedSpace from './rooms/SharedSpace';
import SharedRoomsView from './rooms/SharedRoomsView';
import RoomModals from './rooms/RoomModals';

// Import custom hooks
import { useRooms } from '../hooks/useRooms';
import { useSession } from '../hooks/useSession';
import { useTimer } from '../hooks/useTimer';

const RoomsInterface = ({ initialSection, initialXP }) => {
  // State to track active section - ensure userRoom is the default
  const [activeSection, setActiveSection] = useState('userRoom');
  const [sharedLink, setSharedLink] = useState('');
  const [menuExpanded, setMenuExpanded] = useState(false);
  
  // State for modals and interactive functionality
  const [showTodoModal, setShowTodoModal] = useState(false);
  const [showMaterialModal, setShowMaterialModal] = useState(false);
  const [showCreateRoomModal, setShowCreateRoomModal] = useState(false);
  const [createRoomError, setCreateRoomError] = useState(null);
  const [showJoinModal, setShowJoinModal] = useState(false);
  const [notifications, setNotifications] = useState([]);
  const [showNotifications, setShowNotifications] = useState(false);
  const [roomInvitations, setRoomInvitations] = useState([]);
  const [showRoomInvitations, setShowRoomInvitations] = useState(false);
  const [showRoomSettings, setShowRoomSettings] = useState(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
  const [showInviteFriendsModal, setShowInviteFriendsModal] = useState(false);
  const [newTodo, setNewTodo] = useState('');
  const [todos, setTodos] = useState([]);
  const [newMaterialType, setNewMaterialType] = useState('link');
  const [newMaterialUrl, setNewMaterialUrl] = useState('');
  const [uploadedFile, setUploadedFile] = useState(null);

  // User's files - now fetched from backend
  const [userMaterials, setUserMaterials] = useState([]);
  const [userNotes, setUserNotes] = useState([]);
  const [newNote, setNewNote] = useState('');
  
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [successMessage, setSuccessMessage] = useState(null);
  
  // XP and privileges state
  const [userXP, setUserXP] = useState(initialXP);
  const [maxRoomParticipants, setMaxRoomParticipants] = useState(4);
  const [sharedRoomsCount, setSharedRoomsCount] = useState(1);
  const [userPrivileges, setUserPrivileges] = useState(null);
  
  // Room management state
  const [showVideoCall, setShowVideoCall] = useState(false);
  const [joinRoomCode, setJoinRoomCode] = useState('');
  const [currentUser, setCurrentUser] = useState(null);
  
  // NEW: Real-time chat state
  const [chatMessages, setChatMessages] = useState([]);
  const [newChatMessage, setNewChatMessage] = useState('');
  const [isTyping, setIsTyping] = useState(false);
  const [typingUsers, setTypingUsers] = useState([]);
  const [onlineUsers, setOnlineUsers] = useState([]);
  const [wsConnection, setWsConnection] = useState(null);
  
  // NEW: Video call state
  const [localStream, setLocalStream] = useState(null);
  const [remoteStreams, setRemoteStreams] = useState(new Map());
  const [isVideoEnabled, setIsVideoEnabled] = useState(true);
  const [isAudioEnabled, setIsAudioEnabled] = useState(true);
  const [isScreenSharing, setIsScreenSharing] = useState(false);
  const [peerConnections, setPeerConnections] = useState(new Map());
  
  // NEW: Room management state
  const [editingRoom, setEditingRoom] = useState(null);
  const [roomToDelete, setRoomToDelete] = useState(null);
  
  // NEW: Refs for video elements
  const localVideoRef = useRef(null);
  const remoteVideoRefs = useRef(new Map());
  
  // NEW: Room creation modal state
  const [newRoomData, setNewRoomData] = useState({
    name: '',
    description: '',
    maxParticipants: 4
  });
  
  // NEW: Friends invitation state
  const [friendsList, setFriendsList] = useState([]);
  const [invitingToRoom, setInvitingToRoom] = useState(null);
  
  // NEW: Sharing functionality state
  const [sharingItem, setSharingItem] = useState(null);
  const [showShareModal, setShowShareModal] = useState(false);
  const [selectedRoomsForSharing, setSelectedRoomsForSharing] = useState([]);
  const [selectedMaterialsForSharing, setSelectedMaterialsForSharing] = useState([]);
  const [refreshTrigger, setRefreshTrigger] = useState(0);

  // Use custom hooks
  const {
    userRooms,
    setUserRooms,
    currentRoom,
    setCurrentRoom,
    refreshRooms,
    createRoom,
    joinRoom,
    leaveRoom,
    generateInvitationCode,
    handleEditRoom,
    handleDeleteRoom,
    inviteFriendToRoom,
    handleAcceptRoomInvitation
  } = useRooms(currentUser);

  const {
    sessionActive,
    setSessionActive,
    sessionXP,
    setSessionXP,
    inactiveTime,
    setInactiveTime,
    showSessionXPModal,
    setShowSessionXPModal,
    lastUserActivity,
    setLastUserActivity,
    startSession,
    endSession,
    activityHandler
  } = useSession(currentUser);

  const {
    isTimerRunning,
    setIsTimerRunning,
    workTime,
    setWorkTime,
    breakTime,
    setBreakTime,
    currentTime,
    setCurrentTime,
    isBreak,
    setIsBreak,
    showTimerNotification,
    setShowTimerNotification,
    timerNotificationMessage,
    setTimerNotificationMessage,
    timerInterval,
    setTimerInterval,
    toggleTimer,
    formatTime,
    closeTimerNotification
  } = useTimer(sessionActive);



  // Debounced error setter to prevent excessive re-renders
  const setErrorDebounced = useCallback((message) => {
    if (errorTimeout.current) {
      clearTimeout(errorTimeout.current);
    }
    errorTimeout.current = setTimeout(() => {
      setError(message);
    }, 100);
  }, []);

  // Memoized error display to prevent unnecessary re-renders
  const errorDisplay = useMemo(() => {
    if (!error) return null;
    return (
      <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
        <span className="block sm:inline">{error}</span>
        <button
          onClick={() => setError(null)}
          className="float-right font-bold text-red-700 hover:text-red-900"
        >
          ×
        </button>
      </div>
    );
  }, [error]);
  
  // Calculate max participants based on XP - memoized with useCallback
  const calculateMaxParticipants = useCallback((xp) => {
    // Base of 4 participants, then +1 for each 300 XP
    return 4 + Math.floor(xp / 300);
  }, []);

  // Calculate max rooms based on XP - memoized with useCallback
  const calculateMaxRooms = useCallback((xp) => {
    // Base of 1 room, then +1 for each 1000 XP
    return 1 + Math.floor(xp / 1000);
  }, []);

  // Handle room deletion from inside the room
  const handleDeleteRoomFromInside = useCallback((roomId) => {
    // Set the room to delete and show confirmation modal
    setRoomToDelete({ id: roomId, name: currentRoom?.name || 'Unknown Room' });
    setShowDeleteConfirm(true);
  }, [currentRoom?.name, setRoomToDelete, setShowDeleteConfirm]);

  // Handle leaving a room
  const handleLeaveRoom = useCallback(async (roomId) => {
    try {
      await leaveRoom(roomId);
      // After leaving, redirect to shared rooms view
      setActiveSection('sharedRooms');
      setCurrentRoom(null);
      setSuccessMessage('Successfully left the room');
      // Clear success message after 3 seconds
      setTimeout(() => setSuccessMessage(null), 3000);
    } catch (error) {
      setErrorDebounced('Failed to leave room: ' + error.message);
    }
  }, [leaveRoom, setActiveSection, setCurrentRoom, setErrorDebounced]);

  // Fetch notifications function - defined early to avoid initialization issues
  const fetchNotifications = useCallback(async () => {
    try {
      const data = await fetch('/api/v1/notifications/', {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('accessToken')}`,
          'Content-Type': 'application/json'
        }
      });
      if (data.ok) {
        const notificationsData = await data.json();
        // Filter out room invitations since they're handled separately
        const filteredNotifications = (notificationsData.notifications || []).filter(
          n => n.type !== 'room_invitation'
        );
        setNotifications(filteredNotifications);
      }
    } catch (error) {
      console.error('Error fetching notifications:', error);
    }
  }, []);

  // Fetch room invitations separately for shared rooms interface
  const fetchRoomInvitations = useCallback(async () => {
    try {
      const data = await fetch('/api/v1/notifications/', {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('accessToken')}`,
          'Content-Type': 'application/json'
        }
      });
      if (data.ok) {
        const notificationsData = await data.json();
        // Only get room invitations
        const invitations = (notificationsData.notifications || []).filter(
          n => n.type === 'room_invitation'
        );
        setRoomInvitations(invitations);
      }
    } catch (error) {
      console.error('Error fetching room invitations:', error);
    }
  }, []);

  // Fetch user data and permissions
  useEffect(() => {
    const fetchUserData = async () => {
      try {
        setLoading(true);
        setError(null);
        
        // Check if user is authenticated
        const token = localStorage.getItem('accessToken');
        if (!token) {
          console.log('No access token found, skipping data fetch');
          setLoading(false);
          return;
        }
        
        console.log('Fetching user data...');
        
        // Fetch user data in parallel
        const [userProfile, todosData, materialsData, notesData, privilegesData, roomsData] = await Promise.all([
          authAPI.getProfile(),
          todosAPI.getTodos(),
          materialsAPI.getMaterials(), 
          notesAPI.getNotes(),
          sessionsAPI.getPrivileges(),
          roomsAPI.getRooms()
        ]);
        
        console.log('Rooms API response:', roomsData);
        console.log('Rooms array:', roomsData?.rooms);
        
        // Update user XP and profile
        setUserXP(userProfile.user?.xp || initialXP);
        setCurrentUser(userProfile.user);
        
        // Update data
        setTodos(todosData.todos || []);
        setUserMaterials(materialsData.materials || []);
        setUserNotes(notesData.notes || []);
        setUserRooms(roomsData.rooms || []);
        
        console.log('Updated userRooms state:', roomsData.rooms || []);
        
        // Update privileges from backend
        if (privilegesData) {
          setUserPrivileges(privilegesData);
          setMaxRoomParticipants(privilegesData.maxParticipants || 4);
          setSharedRoomsCount(privilegesData.maxSharedRooms || 1);
        } else {
          // Fallback to calculated values
          setMaxRoomParticipants(calculateMaxParticipants(userProfile.xp || initialXP));
          setSharedRoomsCount(calculateMaxRooms(userProfile.xp || initialXP));
        }
        
      } catch (error) {
        console.error('Error fetching user data:', error);
        if (error.message.includes('401') || error.message.includes('Unauthorized')) {
          setError('Authentication failed. Please log in again.');
        } else if (error.message.includes('Failed to fetch') || error.message.includes('NetworkError')) {
          setError('Cannot connect to server. Please check your connection and try again.');
        } else {
          setError('Failed to load data. Please try again.');
        }
      } finally {
        setLoading(false);
      }
    };

    fetchUserData();
  }, [initialXP, calculateMaxParticipants, calculateMaxRooms]);

  // Fetch notifications on component mount
  useEffect(() => {
    const token = localStorage.getItem('accessToken');
    if (token) {
      fetchNotifications();
    }
  }, [fetchNotifications]);

  // Fetch room invitations on component mount
  useEffect(() => {
    const token = localStorage.getItem('accessToken');
    if (token) {
      fetchRoomInvitations();
    }
  }, [fetchRoomInvitations]);

  // NEW: Cleanup video streams on unmount
  useEffect(() => {
    return () => {
      if (localStream) {
        localStream.getTracks().forEach(track => track.stop());
      }
      if (remoteStreams.size > 0) {
        remoteStreams.forEach(stream => {
          stream.getTracks().forEach(track => track.stop());
        });
      }
    };
  }, [localStream, remoteStreams]);

  // Also fetch user profile when entering a room to ensure currentUser is set
  useEffect(() => {
    if (activeSection === 'sharedSpace' && currentRoom && !currentUser) {
      const fetchUserProfile = async () => {
        try {
          const userProfile = await authAPI.getProfile();
          setCurrentUser(userProfile.user);
        } catch (error) {
          console.error('Error fetching user profile:', error);
        }
      };
      fetchUserProfile();
    }
  }, [activeSection, currentRoom?.id, currentRoom?._id, currentUser?.id, currentUser?._id, currentUser]); // Include currentUser to prevent size changes

  // Refs for performance optimization
  const errorTimeout = useRef(null);

  // Memoized room filtering to prevent unnecessary recalculations
  const { userRoom, sharedRooms } = useMemo(() => {
    const userRoom = userRooms.find(room => room.type === 'personal');
    const sharedRooms = userRooms.filter(room => room.type === 'shared' || !room.type);
    return { userRoom, sharedRooms };
  }, [userRooms]);

  // Room management functions
  
  // Memoized room selection handler to prevent unnecessary re-renders
  const handleRoomSelection = useCallback(async (room) => {
    // If there's an active video call, end it before switching rooms
    if (showVideoCall) {
      setShowVideoCall(false);
    }
    
    try {
      // Enter the room first to add user as participant
      const roomId = room.id || room._id;
      await roomsAPI.enterRoom(roomId);
      
      // Refresh rooms to get updated participant count
      await refreshRooms();
      
      setCurrentRoom(room);
      setActiveSection('sharedSpace');
    } catch (error) {
      console.error('Error entering room:', error);
      // Still set the room even if entering fails
      setCurrentRoom(room);
      setActiveSection('sharedSpace');
    }
  }, [showVideoCall, refreshRooms]);

  // Add a new todo item - memoized with useCallback
  const addTodo = useCallback(async () => {
    console.log('addTodo function called with newTodo:', newTodo);
    
    if (newTodo.trim() !== '') {
      try {
        console.log('Attempting to create todo with data:', {
          title: newTodo,
          description: '',
          completed: false,
          roomId: '',
          priority: 2,
          assigneeIds: [],
          tags: []
        });
        
        const newTodoData = await todosAPI.createTodo({
          title: newTodo,
          description: '',
          dueDate: null,
          priority: 2, // Medium priority
          roomId: '', // Empty for now, will be set when room functionality is implemented
          assigneeIds: [],
          tags: []
        });
        
        console.log('Todo created successfully:', newTodoData);
        setTodos(prev => [...prev, newTodoData.todo]);
        setNewTodo('');
        setShowTodoModal(false);
        setError(null); // Clear any previous errors
        
        // Refresh shared data if we're in a shared space
        if (currentRoom && activeSection === 'sharedSpace') {
          setRefreshTrigger(prev => prev + 1);
        }
      } catch (error) {
        console.error('Error creating todo:', error);
        console.error('Error details:', {
          message: error.message,
          stack: error.stack,
          response: error.response
        });
        setError('Failed to create todo. Please try again.');
      }
    } else {
      console.log('addTodo called but newTodo is empty');
    }
  }, [newTodo, setTodos, setNewTodo, setShowTodoModal, setError, currentRoom, activeSection]);

  // Mark todo as done - memoized with useCallback
  const markTodoAsDone = useCallback(async (todoId) => {
    try {
      await todosAPI.completeTodo(todoId);
      setTodos(prev => prev.map(todo => 
        todo.id === todoId ? { ...todo, completed: !todo.completed } : todo
      ));
      
      // Refresh shared data if we're in a shared space
      if (currentRoom && activeSection === 'sharedSpace') {
        setRefreshTrigger(prev => prev + 1);
      }
    } catch (error) {
      console.error('Error updating todo:', error);
      setError('Failed to update todo. Please try again.');
    }
  }, [currentRoom, activeSection]);

  // Delete todo - memoized with useCallback
  const deleteTodo = useCallback(async (todoId) => {
    console.log('deleteTodo called with todoId:', todoId);
    try {
      await todosAPI.deleteTodo(todoId);
      setTodos(prev => prev.filter(todo => todo.id !== todoId));
      console.log('Todo deleted successfully');
      
      // Refresh shared rooms data to remove the deleted todo
      if (currentRoom) {
        refreshRooms();
        // Also refresh the current room's shared data if we're in a shared space
        if (activeSection === 'sharedSpace') {
          // Force a refresh of the shared data by triggering a re-render
          setRefreshTrigger(prev => prev + 1);
        }
      }
    } catch (error) {
      console.error('Error deleting todo:', error);
      console.error('Error details:', {
        message: error.message,
        stack: error.stack,
        response: error.response
      });
      setError('Failed to delete todo. Please try again.');
    }
  }, [currentRoom, refreshRooms, activeSection]);

  // Add material - memoized with useCallback
  const addMaterial = useCallback(async () => {
    console.log('addMaterial function called with newMaterialUrl:', newMaterialUrl, 'uploadedFile:', uploadedFile);
    
    if ((newMaterialType === 'link' && newMaterialUrl.trim() !== '') || 
        (newMaterialType === 'upload' && uploadedFile)) {
      try {
        let materialData;
        
        if (newMaterialType === 'link') {
          materialData = {
            name: `Material ${Date.now()}`,
            description: 'Uploaded material',
            roomId: '', // Empty for now, will be set when room functionality is implemented
            fileType: 'link',
            fileUrl: newMaterialUrl,
            fileSize: 0
          };
        } else {
          // For file uploads, create a blob URL for immediate access
          const blobUrl = URL.createObjectURL(uploadedFile);
          materialData = {
            name: uploadedFile.name,
            description: 'Uploaded material',
            roomId: '', // Empty for now, will be set when room functionality is implemented
            fileType: uploadedFile.type || 'file',
            fileUrl: blobUrl, // Use blob URL for immediate access
            fileSize: uploadedFile.size,
            originalFile: uploadedFile // Store the original file for potential backend upload
          };
        }
        
        console.log('Attempting to create material with data:', materialData);
        
        const newMaterialData = await materialsAPI.createMaterial(materialData);
        
        console.log('Material created successfully:', newMaterialData);
        setUserMaterials(prev => [...prev, newMaterialData.material]);
        setNewMaterialUrl('');
        setUploadedFile(null);
        setShowMaterialModal(false);
        setError(null); // Clear any previous errors
        
        // Refresh shared data if we're in a shared space
        if (currentRoom && activeSection === 'sharedSpace') {
          setRefreshTrigger(prev => prev + 1);
        }
      } catch (error) {
        console.error('Error creating material:', error);
        console.error('Error details:', {
          message: error.message,
          stack: error.stack,
          response: error.response
        });
        setError('Failed to create material. Please try again.');
      }
    } else {
      console.log('addMaterial called but no valid material data');
    }
  }, [newMaterialUrl, newMaterialType, uploadedFile, setUserMaterials, setNewMaterialUrl, setUploadedFile, setShowMaterialModal, setError, currentRoom, activeSection]);

  // Delete material - memoized with useCallback
  const deleteMaterial = useCallback(async (materialId) => {
    console.log('deleteMaterial called with materialId:', materialId);
    try {
      await materialsAPI.deleteMaterial(materialId);
      setUserMaterials(prev => prev.filter(material => material.id !== materialId));
      console.log('Material deleted successfully');
      
      // Refresh shared rooms data to remove the deleted material
      if (currentRoom) {
        refreshRooms();
        // Also refresh the current room's shared data if we're in a shared space
        if (activeSection === 'sharedSpace') {
          // Force a refresh of the shared data by triggering a re-render
          setRefreshTrigger(prev => prev + 1);
        }
      }
    } catch (error) {
      console.error('Error deleting material:', error);
      console.error('Error details:', {
        message: error.message,
        stack: error.stack,
        response: error.response
      });
      setError('Failed to delete material. Please try again.');
    }
  }, [currentRoom, refreshRooms, activeSection]);

  // Add note - memoized with useCallback
  const handleAddNote = useCallback(async () => {
    console.log('handleAddNote function called with newNote:', newNote);
    
    if (newNote.trim() !== '') {
      try {
        console.log('Attempting to create note with data:', {
          content: newNote,
          roomId: '',
          isShared: false
        });
        
        const newNoteData = await notesAPI.createNote({
          content: newNote,
          roomId: '', // Empty for now, will be set when room functionality is implemented
          isShared: false
        });
        
        console.log('Note created successfully:', newNoteData);
        setUserNotes(prev => [...prev, newNoteData.note]);
        setNewNote('');
        setError(null); // Clear any previous errors
        
        // Refresh shared data if we're in a shared space
        if (currentRoom && activeSection === 'sharedSpace') {
          setRefreshTrigger(prev => prev + 1);
        }
      } catch (error) {
        console.error('Error creating note:', error);
        console.error('Error details:', {
          message: error.message,
          stack: error.stack,
          response: error.response
        });
        setError('Failed to create note. Please try again.');
      }
    } else {
      console.log('handleAddNote called but newNote is empty');
    }
  }, [newNote, setUserNotes, setNewNote, setError, currentRoom, activeSection]);

  // Delete note - memoized with useCallback
  const deleteNote = useCallback(async (noteId) => {
    console.log('deleteNote called with noteId:', noteId);
    try {
      await notesAPI.deleteNote(noteId);
      setUserNotes(prev => prev.filter(note => note.id !== noteId));
      console.log('Note deleted successfully');
      
      // Refresh shared rooms data to remove the deleted note
      if (currentRoom) {
        refreshRooms();
        // Also refresh the current room's shared data if we're in a shared space
        if (activeSection === 'sharedSpace') {
          // Force a refresh of the shared data by triggering a re-render
          // This ensures the deleted note is removed from the UI immediately
          setRefreshTrigger(prev => prev + 1);
        }
      }
    } catch (error) {
      console.error('Error deleting note:', error);
      console.error('Error details:', {
        message: error.message,
        stack: error.stack,
        response: error.response
      });
      setError('Failed to delete note. Please try again.');
    }
  }, [currentRoom, refreshRooms, activeSection]);

  // NEW: Sharing functionality functions
  // Open share modal for an item
  const openShareModal = useCallback((item, itemType) => {
    setSharingItem({ ...item, type: itemType });
    setSelectedRoomsForSharing([]);
    setShowShareModal(true);
  }, []);

  // Share item to selected rooms
  const shareItemToRooms = useCallback(async () => {
    if (!sharingItem || selectedRoomsForSharing.length === 0) {
      setError('Please select at least one room to share in');
      return;
    }

    try {
      const { type, ...item } = sharingItem;
      
      // Share to each selected room
      for (const roomId of selectedRoomsForSharing) {
        if (type === 'todo') {
          // Check if todo already exists in this room
          const existingTodos = await todosAPI.getTodos();
          const alreadyShared = existingTodos.todos?.some(todo => 
            todo.roomId === roomId && todo.title === item.title
          );
          
          if (!alreadyShared) {
            await todosAPI.createTodo({
              ...item,
              roomId: roomId,
              isShared: true
            });
          }
        } else if (type === 'material') {
          // Check if material already exists in this room
          const existingMaterials = await materialsAPI.getRoomMaterials(roomId);
          const alreadyShared = existingMaterials.materials?.some(material => 
            material.roomId === roomId && material.name === item.name
          );
          
          if (!alreadyShared) {
            await materialsAPI.createMaterial({
              ...item,
              roomId: roomId,
              isShared: true
            });
          }
        } else if (type === 'note') {
          // Check if note already exists in this room
          const existingNotes = await notesAPI.getRoomNotes(roomId);
          const alreadyShared = existingNotes.notes?.some(note => 
            note.roomId === roomId && note.content === item.content
          );
          
          if (!alreadyShared) {
            await notesAPI.createNote({
              ...item,
              roomId: roomId,
              isShared: true
            });
          }
        }
      }

      setShowShareModal(false);
      setSharingItem(null);
      setSelectedRoomsForSharing([]);
      setError(null);
      
      // Refresh the rooms data to show the new shared items
      refreshRooms();
      
      // Also refresh the current room's shared data if we're in a shared space
      if (currentRoom && activeSection === 'sharedSpace') {
        // Force a refresh of the shared data by triggering a re-render
        setRefreshTrigger(prev => prev + 1);
      }
      
    } catch (error) {
      console.error('Error sharing item:', error);
      setError('Failed to share item. Please try again.');
    }
  }, [sharingItem, selectedRoomsForSharing, refreshRooms, currentRoom, activeSection]);

  // Toggle room selection for sharing
  const toggleRoomSelection = useCallback((roomId) => {
    setSelectedRoomsForSharing(prev => 
      prev.includes(roomId) 
        ? prev.filter(id => id !== roomId)
        : [...prev, roomId]
    );
  }, []);

  // NEW: Toggle material selection for sharing to a specific room
  const toggleMaterialSelection = useCallback((material, materialType) => {
    setSelectedMaterialsForSharing(prev => {
      const existingIndex = prev.findIndex(m => m.id === material.id && m.type === materialType);
      if (existingIndex >= 0) {
        // Remove if already selected
        return prev.filter((_, index) => index !== existingIndex);
      } else {
        // Add if not selected
        return [...prev, { ...material, type: materialType }];
      }
    });
  }, []);

      // NEW: Share files to the current room
  const shareMaterialsToRoom = useCallback(async () => {
    if (!currentRoom || selectedMaterialsForSharing.length === 0) {
              setError('Please select files to share and ensure you are in a shared room');
      return;
    }

    try {
      // Share each selected material to the current room
      for (const material of selectedMaterialsForSharing) {
        const { type, ...item } = material;
        
        if (type === 'todo') {
          try {
            // Check if todo already exists in this room
            const existingTodos = await todosAPI.getTodos();
            const alreadyShared = existingTodos.todos?.some(todo => 
              todo.roomId === (currentRoom.id || currentRoom._id) && todo.title === item.title
            );
            
            if (!alreadyShared) {
              await todosAPI.createTodo({
                ...item,
                roomId: currentRoom.id || currentRoom._id,
                isShared: true
              });
            }
          } catch (todoError) {
            console.error('Error sharing todo:', todoError);
            // Continue with other items even if this one fails
          }
        } else if (type === 'material') {
          try {
            // Check if material already exists in this room
            const existingMaterials = await materialsAPI.getRoomMaterials(currentRoom.id || currentRoom._id);
            const alreadyShared = existingMaterials.materials?.some(material => 
              material.roomId === (currentRoom.id || currentRoom._id) && material.name === item.name
            );
            
            if (!alreadyShared) {
              await materialsAPI.createMaterial({
                ...item,
                roomId: currentRoom.id || currentRoom._id,
                isShared: true
              });
            }
          } catch (materialError) {
            console.error('Error sharing material:', materialError);
            // Continue with other items even if this one fails
          }
        } else if (type === 'note') {
          try {
            // Check if note already exists in this room
            const existingNotes = await notesAPI.getRoomNotes(currentRoom.id || currentRoom._id);
            const alreadyShared = existingNotes.notes?.some(note => 
              note.roomId === (currentRoom.id || currentRoom._id) && note.content === item.content
            );
            
            if (!alreadyShared) {
              await notesAPI.createNote({
                ...item,
                roomId: currentRoom.id || currentRoom._id,
                isShared: true
              });
            }
          } catch (noteError) {
            console.error('Error sharing note:', noteError);
            // Continue with other items even if this one fails
          }
        }
      }

      // Clear selection and refresh room data
      setSelectedMaterialsForSharing([]);
      setError(null);
      
      // Refresh the room data to show the new shared items
      refreshRooms();
      
      // Also refresh the current room's shared data if we're in a shared space
      if (currentRoom && activeSection === 'sharedSpace') {
        // This will trigger a refresh of the shared data in SharedSpace
        // through the useEffect that depends on currentRoom
      }
      
      // Show success message
      setSuccessMessage('Files shared successfully!');
      setTimeout(() => setSuccessMessage(null), 3000);
      
    } catch (error) {
      console.error('Error sharing files:', error);
      if (error.message && error.message.includes('Room ID required')) {
        setError('Failed to share files: Invalid room. Please try refreshing the page.');
      } else {
        setError('Failed to share files. Please try again.');
      }
    }
      }, [currentRoom, selectedMaterialsForSharing, refreshRooms, activeSection]);

  // Function to show XP requirements modal - memoized with useCallback
  const showXPRequirement = useCallback((action) => {
    let requiredXP = 0;
    let message = '';
    
    if (action === 'add_participant') {
      requiredXP = (Math.floor(userXP / 300) + 1) * 300;
      message = `You need ${requiredXP} XP to add another participant to your room.`;
    } else if (action === 'new_room') {
      requiredXP = (Math.floor(userXP / 1000) + 1) * 1000;
      message = `You need ${requiredXP} XP to create another room.`;
    }
    
    // In a real app, this would show a modal component instead of an alert
    alert(`${message} Current XP: ${userXP}`);
  }, [userXP]);

  // NEW: Room creation functions
  const handleCreateRoom = useCallback(async () => {
    try {
      if (!newRoomData.name.trim()) {
        setCreateRoomError('Room name is required');
        return;
      }
      
      if (newRoomData.name.trim().length < 3) {
        setCreateRoomError('Room name must be at least 3 characters long');
        return;
      }
      
      // Check if user is authenticated
      const token = localStorage.getItem('accessToken');
      if (!token) {
        setCreateRoomError('You must be logged in to create a room. Please log in again.');
        return;
      }
      
      // Simple validation - allow up to 4 participants without XP requirement
      if (newRoomData.maxParticipants > 4) {
        setCreateRoomError('Maximum participants allowed is 4.');
        return;
      }
      
      const room = await createRoom(newRoomData);
      
      // Always add the room to state immediately
      setUserRooms(prev => [...prev, room]);
      
      // Also refresh from backend to ensure consistency
      try {
        const roomsData = await roomsAPI.getRooms();
        if (roomsData?.rooms) {
          setUserRooms(roomsData.rooms);
        }
      } catch (refreshError) {
        // Room is already in state from above, so we're good
      }
      
      // Set as current room
      setCurrentRoom(room);
      
      // Switch to the shared rooms interface to show the created room
      setActiveSection('sharedRooms');
      
      // Show success message
      setSuccessMessage(`Room "${room.name}" created successfully!`);
      setTimeout(() => setSuccessMessage(null), 3000);
      
      // Close modal and reset form
      setShowCreateRoomModal(false);
      setNewRoomData({ name: '', description: '', maxParticipants: 4 });
      setCreateRoomError(null);
    } catch (error) {
      setCreateRoomError(`Failed to create room: ${error.message || 'Unknown error'}`);
    }
  }, [newRoomData, createRoom, currentUser]);

  // NEW: Friends invitation functions
  const fetchFriends = useCallback(async () => {
    try {
      const data = await friendsAPI.getFriends();
      setFriendsList(data.friends || []);
    } catch (error) {
      console.error('Error fetching friends:', error);
      setError('Failed to load friends list');
    }
  }, []);

  const handleInviteFriends = useCallback((room) => {
    setInvitingToRoom(room);
    setShowInviteFriendsModal(true);
    setError(null); // Clear any previous errors
    // Only fetch friends if we don't have them already
    if (friendsList.length === 0) {
      fetchFriends();
    }
  }, [fetchFriends, friendsList.length]);

  // Toggle menu expansion - memoized with useCallback
  const toggleMenu = useCallback(() => {
    setMenuExpanded(!menuExpanded);
  }, [menuExpanded]);



  // NEW: WebSocket connection effect
  useEffect(() => {
    if (currentRoom && activeSection === 'sharedSpace') {
      const roomId = currentRoom.id || currentRoom._id;
      
      // Load chat history
      const loadChatHistory = async () => {
        try {
          const history = await realtimeAPI.getChatHistory(roomId);
          setChatMessages(history.messages || []);
        } catch (error) {
          console.error('Error loading chat history:', error);
        }
      };
      
      loadChatHistory();
      
      // Load online users
      const loadOnlineUsers = async () => {
        try {
          const users = await realtimeAPI.getOnlineUsers(roomId);
          setOnlineUsers(users.users || []);
        } catch (error) {
          console.error('Error loading online users:', error);
        }
      };
      
      loadOnlineUsers();
      
      // Establish WebSocket connection
      const token = localStorage.getItem('token');
      
      // Create WebSocket connection
      const wsUrl = `ws://localhost:8080/api/v1/realtime/ws?roomId=${roomId}`;
      const ws = new WebSocket(wsUrl);
      
      // Set up connection handler
      ws.onopen = () => {
        setWsConnection(ws);
      };
      
      ws.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data);
          
          switch (message.type) {
            case 'chat':
              setChatMessages(prev => [...prev, message]);
              break;
            case 'user_online':
              setOnlineUsers(prev => {
                const userExists = prev.find(u => u.userId === message.userId);
                if (!userExists) {
                  return [...prev, { userId: message.userId, username: message.username }];
                }
                return prev;
              });
              break;
            case 'user_offline':
              setOnlineUsers(prev => prev.filter(u => u.userId !== message.userId));
              break;
            case 'typing':
              // Handle typing indicators
              break;
            default:
              break;
          }
        } catch (error) {
          console.error('Error parsing WebSocket message:', error);
        }
      };
      
      ws.onerror = (error) => {
        console.error('WebSocket error:', error);
      };
      
      ws.onclose = (event) => {
        setWsConnection(null);
      };
      
      // Cleanup function for this effect
      return () => {
        if (ws.readyState === WebSocket.OPEN || ws.readyState === WebSocket.CONNECTING) {
          ws.close();
        }
      };
    }
    
    return () => {
      if (wsConnection) {
        console.log('Cleaning up existing WebSocket connection');
        wsConnection.close();
        setWsConnection(null);
      }
    };
  }, [currentRoom?.id, currentRoom?._id, activeSection]); // Only depend on room ID, not entire object

  // Cleanup video call when switching rooms or unmounting
  useEffect(() => {
    return () => {
      // If video call is active when component unmounts, ensure it's properly closed
      if (showVideoCall) {
        setShowVideoCall(false);
      }
    };
  }, [showVideoCall]);

  // Show loading state
  if (loading) {
    return (
      <div className="flex h-screen relative font-sans overflow-hidden">
        <div 
          className="fixed inset-0 z-0"
          style={{
            background: 'linear-gradient(135deg, #7681b3 0%, #8b96d1 50%, #a5b0e8 100%)',
            backgroundSize: 'cover',
            backgroundPosition: 'center',
            backgroundAttachment: 'fixed',
          }}
          aria-hidden="true"
        />
        <div className="flex-grow relative z-10 flex items-center justify-center">
          <div className="text-center text-white">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-white mx-auto mb-4"></div>
            <p className="text-lg">Loading Rooms Interface...</p>
          </div>
        </div>
      </div>
    );
  }

  // Show error state only for critical errors that prevent the interface from loading
  if (error && error.includes('Failed to load')) {
    return (
      <div className="flex h-screen relative font-sans overflow-hidden">
        <div 
          className="fixed inset-0 z-0"
          style={{
            background: 'linear-gradient(135deg, #7681b3 0%, #8b96d1 50%, #a5b0e8 100%)',
            backgroundSize: 'cover',
            backgroundPosition: 'center',
            backgroundAttachment: 'fixed',
          }}
          aria-hidden="true"
        />
        <div className="flex-grow relative z-10 flex items-center justify-center">
          <div className="text-center text-white bg-red-500 bg-opacity-80 p-6 rounded-lg max-w-md mx-4">
            <h2 className="text-xl font-bold mb-4">Error Loading Rooms</h2>
            <p className="mb-4">{error}</p>
            <button 
              onClick={() => window.location.reload()} 
              className="bg-white text-red-500 px-4 py-2 rounded hover:bg-gray-100"
            >
              Retry
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="flex h-screen relative font-sans overflow-hidden">
      {/* Full screen background image */}
      <div 
        className="fixed inset-0 z-0"
        style={{
          background: 'linear-gradient(135deg, #7681b3 0%, #8b96d1 50%, #a5b0e8 100%)',
          backgroundSize: 'cover',
          backgroundPosition: 'center',
          backgroundAttachment: 'fixed',
        }}
        aria-hidden="true"
      />
      
      {/* Side Menu - Collapsible */}
      <RoomSidebar
        menuExpanded={menuExpanded}
        toggleMenu={toggleMenu}
        activeSection={activeSection}
        setActiveSection={setActiveSection}
        currentRoom={currentRoom}
        sharedRooms={sharedRooms}
        handleRoomSelection={handleRoomSelection}
        refreshRooms={refreshRooms}
      />
      
      {/* Main content */}
      <div className="flex-grow relative z-10 overflow-hidden flex flex-col">
        {/* Notification area for success messages and non-critical errors */}
        {(successMessage || (error && !error.includes('Failed to load'))) && (
          <div className="fixed top-20 right-4 z-50 max-w-sm">
            {successMessage && (
              <div className="bg-green-500 text-white p-4 rounded-lg shadow-lg mb-2">
                <div className="flex items-center justify-between">
                  <div className="flex items-center">
                    <Icon icon="mdi:check-circle" className="w-5 h-5 mr-2" />
                    <span>{successMessage}</span>
                  </div>
                  <button 
                    onClick={() => setSuccessMessage(null)} 
                    className="ml-2 text-white hover:text-green-200"
                  >
                    <Icon icon="mdi:close" className="w-4 h-4" />
                  </button>
                </div>
              </div>
            )}
            {error && !error.includes('Failed to load') && (
              <div className="bg-red-500 text-white p-4 rounded-lg shadow-lg mb-2">
                <div className="flex items-center justify-between">
                  <div className="flex items-center">
                    <Icon icon="mdi:alert-circle" className="w-5 h-5 mr-2" />
                    <span>{error}</span>
                  </div>
                  <button 
                    onClick={() => setError(null)} 
                    className="ml-2 text-white hover:text-red-200"
                  >
                    <Icon icon="mdi:close" className="w-4 h-4" />
                  </button>
                </div>
              </div>
            )}
          </div>
        )}
        
        <main className="flex-grow mt-12">
          {activeSection === 'userRoom' && (
            <UserRoom 
              sessionActive={sessionActive} 
              sessionXP={sessionXP} 
              startSession={startSession} 
              endSession={endSession}
              todos={todos}
              markTodoAsDone={markTodoAsDone}
              deleteTodo={deleteTodo}
              userMaterials={userMaterials}
              deleteMaterial={deleteMaterial}
              userNotes={userNotes}
              deleteNote={deleteNote}
              handleAddNote={handleAddNote}
              newNote={newNote}
              setNewNote={setNewNote}
              // Timer state from context
              isTimerRunning={isTimerRunning}
              workTime={workTime}
              setWorkTime={setWorkTime}
              breakTime={breakTime}
              setBreakTime={setBreakTime}
              currentTime={currentTime}
              isBreak={isBreak}
              showTimerNotification={showTimerNotification}
              timerNotificationMessage={timerNotificationMessage}
              toggleTimer={toggleTimer}
              formatTime={formatTime}
              closeTimerNotification={closeTimerNotification}
              // Modal state setters
              setShowTodoModal={setShowTodoModal}
              setShowMaterialModal={setShowMaterialModal}
              // Room information and delete functionality
              currentRoom={currentRoom}
              currentUser={currentUser}
              onDeleteRoom={handleDeleteRoomFromInside}
            />
          )}
          {activeSection === 'sharedSpace' && currentRoom && (
            <SharedSpace 
              userXP={userXP} 
              maxRoomParticipants={maxRoomParticipants} 
              showXPRequirement={showXPRequirement} 
              todos={todos}
              markTodoAsDone={markTodoAsDone}
              deleteTodo={deleteTodo}
              userMaterials={userMaterials}
              deleteMaterial={deleteMaterial}
              userNotes={userNotes}
              deleteNote={deleteNote}
              handleAddNote={handleAddNote}
              newNote={newNote}
              setNewNote={setNewNote}
              setActiveSection={setActiveSection}
              currentRoom={currentRoom}
              currentUser={currentUser}
              wsConnection={wsConnection}
              chatMessages={chatMessages}
              setChatMessages={setChatMessages}
              typingUsers={typingUsers}
              setErrorDebounced={setErrorDebounced}
              selectedMaterialsForSharing={selectedMaterialsForSharing}
              toggleMaterialSelection={toggleMaterialSelection}
              shareMaterialsToRoom={shareMaterialsToRoom}
              refreshRooms={refreshRooms}
              refreshTrigger={refreshTrigger}
              onDeleteRoom={handleDeleteRoomFromInside}
              onLeaveRoom={handleLeaveRoom}
            />
          )}
          {activeSection === 'sharedRooms' && (
            <SharedRoomsView 
              userRooms={userRooms}
              currentUser={currentUser}
              setShowCreateRoomModal={setShowCreateRoomModal}
              setShowJoinModal={setShowJoinModal}
              refreshRooms={refreshRooms}
              showNotifications={showNotifications}
              setShowNotifications={setShowNotifications}
              fetchNotifications={fetchNotifications}
              notifications={notifications}
              showRoomInvitations={showRoomInvitations}
              setShowRoomInvitations={setShowRoomInvitations}
              fetchRoomInvitations={fetchRoomInvitations}
              roomInvitations={roomInvitations}
              testBackendConnectivity={testBackendConnectivity}
              roomsAPI={roomsAPI}
              setRoomToDelete={setRoomToDelete}
              setShowDeleteConfirm={setShowDeleteConfirm}
              setEditingRoom={setEditingRoom}
              setShowRoomSettings={setShowRoomSettings}
              handleInviteFriends={handleInviteFriends}
              handleJoinRoom={joinRoom}
              handleStartCall={() => setShowVideoCall(true)}
              handleRoomSelection={handleRoomSelection}
            />
          )}
        </main>
      </div>
      
      {/* Video Call Modal */}
      {showVideoCall && currentRoom && currentUser && (
        <VideoCallRoom
          roomId={currentRoom.id || currentRoom._id}
          currentUser={currentUser}
          onClose={() => {
            // Ensure video call is properly ended before closing
            setShowVideoCall(false);
          }}
        />
      )}

      {/* All Modals */}
      <RoomModals
        // Todo Modal
        showTodoModal={showTodoModal}
        setShowTodoModal={setShowTodoModal}
        newTodo={newTodo}
        setNewTodo={setNewTodo}
        addTodo={addTodo}
        
        // Material Modal
        showMaterialModal={showMaterialModal}
        setShowMaterialModal={setShowMaterialModal}
        newMaterialType={newMaterialType}
        setNewMaterialType={setNewMaterialType}
        newMaterialUrl={newMaterialUrl}
        setNewMaterialUrl={setNewMaterialUrl}
        uploadedFile={uploadedFile}
        setUploadedFile={setUploadedFile}
        addMaterial={addMaterial}
        
        // Session XP Modal
        showSessionXPModal={showSessionXPModal}
        setShowSessionXPModal={setShowSessionXPModal}
        sessionXP={sessionXP}
        
        // Video Call Modal
        showVideoCall={showVideoCall}
        setShowVideoCall={setShowVideoCall}
        currentRoom={currentRoom}
        currentUser={currentUser}
        VideoCallRoom={VideoCallRoom}
        
        // Join Room Modal
        showJoinModal={showJoinModal}
        setShowJoinModal={setShowJoinModal}
        joinRoomCode={joinRoomCode}
        setJoinRoomCode={setJoinRoomCode}
        joinRoom={joinRoom}
        
        // Create Room Modal
        showCreateRoomModal={showCreateRoomModal}
        setShowCreateRoomModal={setShowCreateRoomModal}
        newRoomData={newRoomData}
        setNewRoomData={setNewRoomData}
        createRoomError={createRoomError}
        setCreateRoomError={setCreateRoomError}
        handleCreateRoom={handleCreateRoom}
        maxRoomParticipants={maxRoomParticipants}
        userXP={userXP}
        
        // Room Settings Modal
        showRoomSettings={showRoomSettings}
        setShowRoomSettings={setShowRoomSettings}
        editingRoom={editingRoom}
        setEditingRoom={setEditingRoom}
        handleEditRoom={handleEditRoom}
        generateInvitationCode={generateInvitationCode}
        errorDisplay={errorDisplay}
        
        // Delete Room Confirmation Modal
        showDeleteConfirm={showDeleteConfirm}
        setShowDeleteConfirm={setShowDeleteConfirm}
        roomToDelete={roomToDelete}
        setRoomToDelete={setRoomToDelete}
        handleDeleteRoom={handleDeleteRoom}
        onDeleteSuccess={() => {
          // After successful deletion, redirect to shared rooms view
          setActiveSection('sharedRooms');
          setCurrentRoom(null);
          setSuccessMessage('Room deleted successfully');
          // Clear success message after 3 seconds
          setTimeout(() => setSuccessMessage(null), 3000);
        }}
        
        // Friends Invitation Modal
        showInviteFriendsModal={showInviteFriendsModal}
        setShowInviteFriendsModal={setShowInviteFriendsModal}
        invitingToRoom={invitingToRoom}
        setInvitingToRoom={setInvitingToRoom}
        friendsList={friendsList}
        inviteFriendToRoom={inviteFriendToRoom}
        userRooms={userRooms}
        
        // Notifications Modal
        showNotifications={showNotifications}
        setShowNotifications={setShowNotifications}
        notifications={notifications}
        
        // Room Invitations Modal
        showRoomInvitations={showRoomInvitations}
        setShowRoomInvitations={setShowRoomInvitations}
        roomInvitations={roomInvitations}
        handleAcceptRoomInvitation={handleAcceptRoomInvitation}
        
        // Share Modal
        showShareModal={showShareModal}
        setShowShareModal={setShowShareModal}
        sharingItem={sharingItem}
        setSharingItem={setSharingItem}
        selectedRoomsForSharing={selectedRoomsForSharing}
        toggleRoomSelection={toggleRoomSelection}
        shareItemToRooms={shareItemToRooms}
        sharedRooms={sharedRooms}
      />
    </div>
  );
};

// Add PropTypes for the main component
RoomsInterface.propTypes = {
  initialSection: PropTypes.oneOf(['userRoom', 'sharedSpace', 'sharedRooms']),
  initialXP: PropTypes.number
};

// Add default props
RoomsInterface.defaultProps = {
  initialSection: 'userRoom',
  initialXP: 250
};

export default RoomsInterface;