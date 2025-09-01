import React, { useState, useEffect, useRef, useCallback } from 'react';
import { Icon } from '@iconify/react';
import PropTypes from 'prop-types';
import { materialsAPI, notesAPI, todosAPI, realtimeAPI } from '../../utils/api';

const SharedSpace = ({ 
  userXP, 
  maxRoomParticipants, 
  showXPRequirement, 
  todos,
  markTodoAsDone,
  deleteTodo,
  userMaterials,
  deleteMaterial,
  userNotes,
  deleteNote,
  handleAddNote,
  newNote,
  setNewNote,
  setActiveSection,
  currentRoom,
  currentUser,
  wsConnection,
  chatMessages,
  setChatMessages,
  typingUsers,
  setErrorDebounced,
  selectedMaterialsForSharing,
  toggleMaterialSelection,
  shareMaterialsToRoom,
  refreshRooms,
  refreshTrigger,
  onDeleteRoom,
  onLeaveRoom
}) => {
  // State for managing current tab in the side panel
  const [activeTab, setActiveTab] = useState('materials'); // materials, notes, todos
  // State for screen sharing
  const [screenSharing, setScreenSharing] = useState(false);
  // State for fullscreen mode
  const [isFullscreen, setIsFullscreen] = useState(false);
  
  // Shared data - now fetched from backend
  const [sharedMaterials, setSharedMaterials] = useState([]);
  const [sharedNotes, setSharedNotes] = useState([]);
  const [sharedTodos, setSharedTodos] = useState([]);
  const [loading, setLoading] = useState(true);
  
  // NEW: Real-time chat state
  const [newChatMessage, setNewChatMessage] = useState('');
  const [isTyping, setIsTyping] = useState(false);
  
  // Use ref to avoid infinite loop with setErrorDebounced
  const setErrorDebouncedRef = useRef(setErrorDebounced);
  setErrorDebouncedRef.current = setErrorDebounced;
  
  // NEW: Video call state
  const [isVideoCallActive, setIsVideoCallActive] = useState(false);
  const [localStream, setLocalStream] = useState(null);
  const [remoteStreams, setRemoteStreams] = useState(new Map());
  const [isVideoEnabled, setIsVideoEnabled] = useState(true);
  const [isAudioEnabled, setIsAudioEnabled] = useState(true);
  const [isScreenSharing, setIsScreenSharing] = useState(false);
  
  // NEW: Refs for video elements
  const localVideoRef = useRef(null);
  const remoteVideoRefs = useRef(new Map());
  
  // Function to refresh shared data - memoized with useCallback
  const refreshSharedData = useCallback(async () => {
    if (!currentRoom) return;
    
    try {
      setLoading(true);
      
      // Fetch room-specific data
      const roomId = currentRoom.id || currentRoom._id;
      console.log('SharedSpace: Fetching data for roomId:', roomId, 'currentRoom:', currentRoom);
      const [materialsData, notesData, todosData] = await Promise.all([
        materialsAPI.getRoomMaterials(roomId),
        notesAPI.getRoomNotes(roomId),
        todosAPI.getTodos() // TODO: Add getRoomTodos when backend supports it
      ]);
      
      setSharedMaterials(materialsData.materials || []);
      setSharedNotes(notesData.notes || []);
      setSharedTodos(todosData.todos || []);
      
    } catch (error) {
      console.error('Error refreshing shared data:', error);
      setErrorDebouncedRef.current('Failed to refresh shared data. Please try again.');
    } finally {
      setLoading(false);
    }
  }, [currentRoom?.id, currentRoom?._id]); // Only depend on the actual room ID, not the entire object
  
  // Single useEffect to handle data fetching and refreshing
  useEffect(() => {
    if (currentRoom) {
      // Debounce the refresh to prevent excessive API calls
      const timeoutId = setTimeout(() => {
        refreshSharedData();
      }, 1000); // Wait 1 second before making API calls
      
      return () => clearTimeout(timeoutId);
    }
  }, [currentRoom?.id, currentRoom?._id]); // Remove refreshSharedData dependency to break circular dependency
  
  // Cleanup video streams when component unmounts
  useEffect(() => {
    return () => {
      if (localStream) {
        localStream.getTracks().forEach(track => track.stop());
      }
      if (localVideoRef.current) {
        localVideoRef.current.srcObject = null;
      }
    };
  }, [localStream]);
  
  // Function to handle inviting friends to the shared space - memoized with useCallback
  const handleInvite = useCallback(() => {
    const currentParticipants = currentRoom?.participantCount || 0;
    const maxAllowed = currentRoom?.maxParticipants || 4;
    
    if (currentParticipants >= maxAllowed) {
      showXPRequirement('add_participant');
    } else {
      // Check XP requirements for additional participants
      const additionalParticipants = maxAllowed - 4; // Base is 4
      if (additionalParticipants > 0 && userXP < (additionalParticipants * 300)) {
        const requiredXP = additionalParticipants * 300;
        alert(`You need ${requiredXP} XP to invite more than 4 participants. Current: ${userXP} XP`);
      } else {
        // API integration point - open friends list to invite
        // In a real app, this would open a modal with friends list
        alert(`You can invite up to ${maxAllowed - currentParticipants} more participants. Opening friends list...`);
      }
    }
  }, [currentRoom, userXP, showXPRequirement]);
  
  // NEW: Toggle video call
  const toggleVideoCall = useCallback(async () => {
    if (!isVideoCallActive) {
      try {
        const stream = await navigator.mediaDevices.getUserMedia({ 
          video: true, 
          audio: true 
        });
        
        console.log('Got media stream:', stream);
        console.log('Stream tracks:', stream.getTracks().map(t => ({ kind: t.kind, id: t.id, enabled: t.enabled })));
        
        setLocalStream(stream);
        setIsVideoCallActive(true);
        
        // Wait for the next tick to ensure the video element is ready
        setTimeout(() => {
          if (localVideoRef.current) {
            console.log('Setting video srcObject:', stream);
            localVideoRef.current.srcObject = stream;
            
            // Ensure the video element is properly configured
            localVideoRef.current.play().catch(e => {
              console.log('Video play error (expected):', e);
            });
            
            // Add event listeners to debug video loading
            localVideoRef.current.onloadedmetadata = () => {
              console.log('Video metadata loaded');
            };
            localVideoRef.current.oncanplay = () => {
              console.log('Video can play');
            };
            localVideoRef.current.onerror = (e) => {
              console.error('Video error:', e);
            };
          } else {
            console.error('localVideoRef.current is null');
          }
        }, 100);
      } catch (error) {
        console.error('Error starting video call:', error);
        setErrorDebounced('Failed to start video call. Please check camera permissions.');
      }
    } else {
      // Stop video call and cleanup
      if (localStream) {
        localStream.getTracks().forEach(track => track.stop());
        setLocalStream(null);
      }
      
      // Clear video element
      if (localVideoRef.current) {
        localVideoRef.current.srcObject = null;
      }
      
      // Reset screen sharing state
      setIsScreenSharing(false);
      
      setIsVideoCallActive(false);
      setRemoteStreams(new Map());
    }
  }, [isVideoCallActive, localStream, setErrorDebounced]);

  // NEW: Toggle video/audio
  const toggleVideo = useCallback(() => {
    if (localStream) {
      const videoTrack = localStream.getVideoTracks()[0];
      if (videoTrack) {
        videoTrack.enabled = !videoTrack.enabled;
        setIsVideoEnabled(videoTrack.enabled);
      }
    }
  }, [localStream]);

  const toggleAudio = useCallback(() => {
    if (localStream) {
      const audioTrack = localStream.getAudioTracks()[0];
      if (audioTrack) {
        audioTrack.enabled = !audioTrack.enabled;
        setIsAudioEnabled(audioTrack.enabled);
      }
    }
  }, [localStream]);

  // NEW: Toggle fullscreen
  const toggleFullscreen = useCallback(() => {
    if (!document.fullscreenElement) {
      document.documentElement.requestFullscreen().then(() => {
        setIsFullscreen(true);
      }).catch(err => {
        console.error('Error entering fullscreen:', err);
      });
    } else {
      document.exitFullscreen().then(() => {
        setIsFullscreen(false);
      }).catch(err => {
        console.error('Error exiting fullscreen:', err);
      });
    }
  }, []);

  // NEW: Toggle screen sharing
  const toggleScreenSharing = useCallback(async () => {
    try {
      if (!isScreenSharing) {
        const screenStream = await navigator.mediaDevices.getDisplayMedia({ 
          video: true 
        });
        
        // Replace video track with screen track
        const videoTrack = localStream.getVideoTracks()[0];
        const screenTrack = screenStream.getVideoTracks()[0];
        
        localStream.removeTrack(videoTrack);
        localStream.addTrack(screenTrack);
        
        if (localVideoRef.current) {
          localVideoRef.current.srcObject = localStream;
        }
        
        setIsScreenSharing(true);
        
        // Stop screen sharing when user stops
        screenTrack.onended = () => {
          setIsScreenSharing(false);
          // Restore camera
          navigator.mediaDevices.getUserMedia({ video: true }).then(cameraStream => {
            const cameraTrack = cameraStream.getVideoTracks()[0];
            localStream.removeTrack(screenTrack);
            localStream.addTrack(cameraTrack);
            
            // Update video element
            if (localVideoRef.current) {
              localVideoRef.current.srcObject = localStream;
            }
          }).catch(error => {
            console.error('Error restoring camera after screen share:', error);
            // If camera restoration fails, stop the video call
            toggleVideoCall();
          });
        };
        
      } else {
        // Restore camera
        try {
          const cameraStream = await navigator.mediaDevices.getUserMedia({ video: true });
          const cameraTrack = cameraStream.getVideoTracks()[0];
          
          const currentVideoTrack = localStream.getVideoTracks()[0];
          if (currentVideoTrack) {
            localStream.removeTrack(currentVideoTrack);
            // Stop the screen sharing track
            currentVideoTrack.stop();
          }
          
          localStream.addTrack(cameraTrack);
          
          if (localVideoRef.current) {
            localVideoRef.current.srcObject = localStream;
          }
          
          setIsScreenSharing(false);
        } catch (error) {
          console.error('Error restoring camera:', error);
          setErrorDebounced('Failed to restore camera. Please try again.');
        }
      }
    } catch (error) {
      console.error('Error toggling screen sharing:', error);
      setErrorDebounced('Failed to toggle screen sharing.');
    }
  }, [isScreenSharing, localStream, setErrorDebounced]);
  
  // Toggle screen sharing - memoized with useCallback
  const toggleScreenSharingOld = useCallback(() => {
    // API integration point - start/stop screen sharing
    // In a real app, this would use the Web RTC API or a video call service
    setScreenSharing(prev => !prev);
    
    if (!screenSharing) {
      alert("Screen sharing started");
    } else {
      alert("Screen sharing stopped");
    }
  }, [screenSharing]);
  
  // Handle adding a new shared todo - memoized with useCallback
  const handleAddSharedTodo = useCallback(async (text) => {
    if (text.trim()) {
      try {
        const newTodoData = await todosAPI.createTodo({
          title: text,
          description: '',
          completed: false,
          roomId: '', // Empty for now, will be set when room functionality is implemented
          priority: 2, // Medium priority
          assigneeIds: [],
          tags: []
        });
        
        setSharedTodos(prev => [...prev, newTodoData.todo]);
      } catch (error) {
        console.error('Error creating shared todo:', error);
        setErrorDebounced('Failed to create shared todo. Please try again.');
      }
    }
  }, [setErrorDebounced]);
  
  // Handle marking shared todo as done - memoized with useCallback
  const handleMarkSharedTodoDone = useCallback(async (todoId) => {
    try {
      await todosAPI.completeTodo(todoId);
      setSharedTodos(prev => prev.map(todo => 
        todo.id === todoId ? { ...todo, completed: !todo.completed } : todo
      ));
    } catch (error) {
      console.error('Error updating shared todo:', error);
      setErrorDebounced('Failed to update shared todo. Please try again.');
    }
  }, [setErrorDebounced]);
  
  // Video controls - memoized with useCallback
  const openVideoCall = useCallback(() => {
    // This would open the video call modal
    console.log('Opening video call');
  }, []);

  // NEW: Send chat message
  const sendChatMessage = useCallback(async () => {
    if (!newChatMessage.trim()) {
      return;
    }
    
    if (!wsConnection) {
      setErrorDebounced('No WebSocket connection. Please refresh the page.');
      return;
    }
    
    if (!currentRoom) {
      return;
    }
    
    try {
      const message = {
        type: 'chat',
        content: newChatMessage.trim(),
        roomId: currentRoom.id || currentRoom._id,
        timestamp: new Date().toISOString()
      };
      
      wsConnection.send(JSON.stringify(message));
      setNewChatMessage('');
      
      // Add message to local state immediately for better UX
      setChatMessages(prev => [...prev, {
        ...message,
        id: Date.now(),
        userId: currentUser?._id,
        username: currentUser?.username
      }]);
    } catch (error) {
      console.error('Error sending message:', error);
      setErrorDebounced('Failed to send message. Please try again.');
    }
  }, [newChatMessage, wsConnection, currentRoom, currentUser, setChatMessages, setErrorDebounced]);

  // NEW: Handle typing
  const handleTyping = useCallback((isTyping) => {
    if (!wsConnection || !currentRoom) return;
    
    try {
      wsConnection.send(JSON.stringify({
        type: isTyping ? 'typing_start' : 'typing_stop',
        roomId: currentRoom.id || currentRoom._id
      }));
      setIsTyping(isTyping);
    } catch (error) {
      console.error('Error sending typing indicator:', error);
    }
  }, [wsConnection, currentRoom]);
  
  return (
                  <div className="h-full p-1" style={{ backgroundImage: 'url("/room-background-1.jpg")', backgroundSize: 'cover' }}>
      {/* Fullscreen Video Call Mode */}
      {isVideoCallActive ? (
        <div className="h-full bg-black flex flex-col">
          {/* Fullscreen Video Header */}
          <div className="flex justify-between items-center p-4 bg-black bg-opacity-80 text-white">
            <h2 className="text-xl font-bold">Video Call - {currentRoom?.name}</h2>
            <div className="flex items-center gap-4">
              <span className="text-sm">
                Participants: {currentRoom?.participantCount || 0}/{currentRoom?.maxParticipants || 4}
              </span>
              <button
                onClick={toggleFullscreen}
                className="bg-red-600 hover:bg-red-700 text-white px-3 py-2 rounded-lg"
              >
                Exit Fullscreen
              </button>
            </div>
          </div>
          
          {/* Fullscreen Video Grid */}
          <div className="flex-1 grid grid-cols-2 gap-4 p-4">
            {/* Local video - larger in fullscreen */}
            <div className="relative bg-gray-800 rounded-lg overflow-hidden">
              <video
                ref={localVideoRef}
                autoPlay
                muted
                playsInline
                className="w-full h-full object-cover"
              />
              <div className="absolute bottom-4 left-4 bg-black bg-opacity-50 text-white px-3 py-1 rounded">
                You
              </div>
            </div>
            
            {/* Remote videos - only show if there are actual participants */}
            {(currentRoom?.participantCount || 0) > 1 ? (
              currentRoom.participants?.filter(participantId => participantId !== currentUser?._id).map((participantId, i) => (
                <div key={participantId} className="relative bg-gray-800 rounded-lg overflow-hidden">
                  <div className="w-full h-full bg-gray-700 flex items-center justify-center">
                    <div className="text-center">
                      <div className="w-20 h-20 rounded-full bg-gray-500 mx-auto mb-2 flex items-center justify-center">
                        <Icon icon="mdi:account" className="text-white w-12 h-12" />
                      </div>
                      <p className="text-white text-lg">Participant {i + 1}</p>
                    </div>
                  </div>
                </div>
              ))
            ) : (
              // Show "No other participants" message when there are no participants
              <div className="flex items-center justify-center">
                <div className="text-center">
                  <div className="w-20 h-20 rounded-full bg-gray-500 mx-auto mb-2 flex items-center justify-center">
                    <Icon icon="mdi:account-group" className="text-white w-12 h-12" />
                  </div>
                  <p className="text-white text-lg">No other participants</p>
                  <p className="text-white text-sm opacity-75">Invite others to join!</p>
                </div>
              </div>
            )}
          </div>
          
          {/* Fullscreen Video Controls */}
          <div className="flex justify-center items-center gap-4 p-4 bg-black bg-opacity-80">
            <button
              onClick={toggleVideo}
              className={`rounded-full p-4 ${isVideoEnabled ? 'bg-blue-500 text-white' : 'bg-red-500 text-white'}`}
            >
              <Icon icon={isVideoEnabled ? "mdi:video" : "mdi:video-off"} className="w-6 h-6" />
            </button>
            <button
              onClick={toggleAudio}
              className={`rounded-full p-4 ${isAudioEnabled ? 'bg-blue-500 text-white' : 'bg-red-500 text-white'}`}
            >
              <Icon icon={isAudioEnabled ? "mdi:microphone" : "mdi:microphone-off"} className="w-6 h-6" />
            </button>
            <button
              onClick={toggleScreenSharing}
              className={`rounded-full p-4 ${isScreenSharing ? 'bg-green-500 text-white' : 'bg-gray-500 text-white'}`}
            >
              <Icon icon="mdi:monitor-share" className="w-6 h-6" />
            </button>
            <button
              onClick={() => {
                setIsVideoCallActive(false);
                setIsFullscreen(false);
              }}
              className="bg-red-600 hover:bg-red-700 text-white px-6 py-4 rounded-lg"
            >
              End Call
            </button>
          </div>
        </div>
      ) : (
        <>
          {/* Header with navigation */}
          <div className="flex justify-between items-center mb-2 p-2 bg-white bg-opacity-90 rounded-lg">
            <div className="flex items-center gap-4">
              <button
                onClick={() => setActiveSection('sharedRooms')}
                className="bg-blue-500 hover:bg-blue-600 text-white px-3 py-1 rounded-lg text-sm flex items-center gap-2"
              >
                <Icon icon="mdi:arrow-left" className="w-4 h-4" />
                Back to Rooms
              </button>
              <h2 className="text-[#1D154A] font-bold text-lg">Room: {currentRoom?.name || 'Unknown'}</h2>
            </div>
            <div className="flex items-center gap-4">
              <div className="text-sm text-[#1D154A]">
                Participants: {currentRoom?.participantCount || 0}/{currentRoom?.maxParticipants || 4}
              </div>
              

              
              {/* Room Owner Controls */}
              {currentRoom?.creatorId === currentUser?.uniqueId ? (
                // Delete Room Button - Only for room owner
                <button
                  onClick={() => onDeleteRoom && onDeleteRoom(currentRoom.id || currentRoom._id)}
                  className="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded-lg text-sm font-medium shadow-lg hover:shadow-xl hover:shadow-red-500/30 transition-all duration-200 hover:scale-105 flex items-center gap-2"
                  title="Delete this room"
                >
                  <Icon icon="mdi:trash" className="w-4 h-4" />
                  Delete Room
                </button>
              ) : (
                // Leave Room Button - For non-owners
                <button
                  onClick={() => onLeaveRoom && onLeaveRoom(currentRoom.id || currentRoom._id)}
                  className="bg-gray-500 hover:bg-gray-600 text-white px-3 py-1 rounded-lg text-sm font-medium shadow-lg hover:shadow-xl hover:shadow-gray-500/30 transition-all duration-200 hover:scale-105 flex items-center gap-2"
                  title="Leave this room"
                >
                  <Icon icon="mdi:exit-to-app" className="w-4 h-4" />
                  Leave Room
                </button>
              )}
            </div>
          </div>
          
          {/* Main Content Layout - 75% video / 25% tabs */}
          <div className="grid grid-cols-4 gap-2 h-[calc(100vh-160px)]">
            {/* Video Call Area - 75% of the width */}
            <div className="col-span-3 bg-white bg-opacity-75 rounded-lg shadow-sm border border-[#83A0BC] p-2 flex flex-col">
              <div className="flex justify-between items-center mb-1">
                <h3 className="text-[#1D154A] font-mono uppercase text-xs font-bold">VIDEO CALL</h3>
                <div className="flex items-center space-x-1">
                  <button
                    onClick={toggleScreenSharing}
                    className={`px-1.5 py-0.5 rounded-lg text-xs font-mono ${
                      isScreenSharing 
                        ? 'bg-[#D9A384] text-white' 
                        : 'bg-[#83A0BC] text-white'
                    }`}
                    aria-label={isScreenSharing ? "Stop screen sharing" : "Start screen sharing"}
                  >
                    {isScreenSharing ? "STOP SHARING" : "SHARE SCREEN"}
                  </button>
                  <button
                    onClick={toggleFullscreen}
                    className="text-[#1D154A] hover:text-blue-600 transition-colors"
                    aria-label="Toggle fullscreen"
                  >
                    <Icon icon={isFullscreen ? "mdi:fullscreen-exit" : "mdi:fullscreen"} className="w-4 h-4" />
                  </button>
                </div>
              </div>
              
              {/* Video call interface */}
              {isVideoCallActive ? (
                <div className="flex-1 flex flex-col">
                  {/* Local video */}
                  <div className="relative mb-2">
                    <video
                      ref={localVideoRef}
                      autoPlay
                      muted
                      playsInline
                      className="w-full h-32 object-cover rounded-lg border border-[#83A0BC]"
                    />
                    <div className="absolute bottom-2 left-2 bg-black bg-opacity-50 text-white text-xs px-2 py-1 rounded">
                      You
                    </div>
                  </div>
                  
                  {/* Remote videos grid */}
                  <div className="grid grid-cols-3 gap-2 flex-1">
                    {/* Show actual participants - only if there are any */}
                    {(currentRoom?.participantCount || 0) > 1 ? (
                      currentRoom.participants.map((participant, i) => (
                        <div key={participant.user_id || i} className="relative bg-[#F0F0F0] rounded-lg border border-[#83A0BC] flex flex-col items-center justify-center">
                          <div className="w-12 h-12 rounded-full bg-[#83A0BC] bg-opacity-80 flex items-center justify-center mb-1">
                            <Icon icon="mdi:account" className="text-white w-8 h-8" />
                          </div>
                          <p className="text-[#1D154A] text-xs font-mono">
                            {participant.user_id === currentUser?._id ? 'You' : `User ${i + 1}`}
                          </p>
                          
                          {/* Video controls */}
                          <div className="absolute bottom-1 left-0 right-0 flex justify-center space-x-1">
                            <button 
                              className="bg-white rounded-full p-0.5 hover:bg-gray-100 border border-[#83A0BC]"
                              aria-label="Toggle microphone"
                            >
                              <Icon icon="mdi:microphone" className="text-[#1D154A] w-3 h-3" />
                            </button>
                            <button 
                              className="bg-white rounded-full p-0.5 hover:bg-gray-100 border border-[#83A0BC]"
                              aria-label="Toggle video"
                            >
                              <Icon icon="mdi:video" className="text-[#1D154A] w-3 h-3" />
                            </button>
                          </div>
                        </div>
                      ))
                    ) : (
                      // Show "No participants" message when there are no participants
                      <div className="col-span-3 flex items-center justify-center">
                        <div className="text-center">
                          <div className="w-12 h-12 rounded-full bg-gray-300 flex items-center justify-center mb-1">
                            <Icon icon="mdi:account-group" className="text-gray-500 w-8 h-8" />
                          </div>
                          <p className="text-[#1D154A] text-xs font-mono text-gray-500">No participants</p>
                        </div>
                      </div>
                    )}
                    
                    {/* Show empty slots for remaining capacity */}
                    {Array.from({ length: Math.max(0, (currentRoom?.maxParticipants || 4) - (currentRoom?.participantCount || 1)) }).map((_, i) => (
                      <div key={`empty-${i}`} className="relative bg-[#F0F0F0] rounded-lg border border-[#83A0BC] flex flex-col items-center justify-center opacity-50">
                        <div className="w-12 h-12 rounded-full bg-gray-300 flex items-center justify-center mb-1">
                          <Icon icon="mdi:account-plus" className="text-gray-500 w-8 h-8" />
                        </div>
                        <p className="text-[#1D154A] text-xs font-mono text-gray-500">Empty Slot</p>
                      </div>
                    ))}
                  </div>
                  
                  {/* Video call controls */}
                  <div className="flex justify-center space-x-2 mt-2">
                    <button
                      onClick={toggleVideo}
                      className={`rounded-full p-2 ${isVideoEnabled ? 'bg-blue-500 text-white' : 'bg-red-500 text-white'}`}
                      aria-label={isVideoEnabled ? "Disable video" : "Enable video"}
                    >
                      <Icon icon={isVideoEnabled ? "mdi:video" : "mdi:video-off"} className="w-4 h-4" />
                    </button>
                    <button
                      onClick={toggleAudio}
                      className={`rounded-full p-2 ${isAudioEnabled ? 'bg-blue-500 text-white' : 'bg-red-500 text-white'}`}
                      aria-label={isAudioEnabled ? "Disable audio" : "Enable audio"}
                    >
                      <Icon icon={isAudioEnabled ? "mdi:microphone" : "mdi:microphone-off"} className="w-4 h-4" />
                    </button>
                    <button
                      onClick={toggleScreenSharing}
                      className={`rounded-full p-2 ${isScreenSharing ? 'bg-green-500 text-white' : 'bg-gray-500 text-white'}`}
                      aria-label={isScreenSharing ? "Stop screen sharing" : "Start screen sharing"}
                    >
                      <Icon icon="mdi:monitor-share" className="w-4 h-4" />
                    </button>
                    <button
                      onClick={toggleVideoCall}
                      className="bg-red-500 text-white rounded-full p-2 hover:bg-red-600"
                      aria-label="End video call"
                    >
                      <Icon icon="mdi:phone-hangup" className="w-4 h-4" />
                    </button>
                  </div>
                </div>
              ) : (
                <>
                  {/* Screen sharing UI */}
                  {screenSharing && (
                    <div className="flex-1 bg-[#F0F0F0] rounded-lg border border-[#83A0BC] flex flex-col items-center justify-center mb-1">
                      <div className="text-center">
                        <div className="w-24 h-16 mx-auto bg-[#83A0BC] bg-opacity-80 rounded flex items-center justify-center mb-1">
                          <Icon icon="mdi:monitor-screenshot" className="text-white w-8 h-8" />
                        </div>
                        <p className="text-[#1D154A] text-xs">You are sharing your screen</p>
                      </div>
                    </div>
                  )}
                  
                  {/* Video grid */}
                  <div className={`${screenSharing ? 'h-1/3' : 'flex-1'} grid grid-cols-3 gap-1`}>
                    {/* Show current user first */}
                    <div className="relative bg-[#F0F0F0] rounded-lg border border-[#83A0BC] flex flex-col items-center justify-center">
                      <div className="w-12 h-12 rounded-full bg-[#83A0BC] bg-opacity-80 flex items-center justify-center mb-1">
                        {currentUser?.avatarUrl ? (
                          <img 
                            src={currentUser.avatarUrl} 
                            alt="Your avatar" 
                            className="w-12 h-12 rounded-full object-cover"
                          />
                        ) : (
                          <Icon icon="mdi:account" className="text-white w-8 h-8" />
                        )}
                      </div>
                      <p className="text-[#1D154A] text-xs font-mono">
                        {currentUser?.username || 'You'}
                      </p>
                      
                      {/* Video controls */}
                      <div className="absolute bottom-1 left-0 right-0 flex justify-center space-x-1">
                        <button 
                          className="bg-white rounded-full p-0.5 hover:bg-gray-100 border border-[#83A0BC]"
                          aria-label="Toggle microphone"
                        >
                          <Icon icon="mdi:microphone" className="text-[#1D154A] w-3 h-3" />
                        </button>
                        <button 
                          className="bg-white rounded-full p-0.5 hover:bg-gray-100 border border-[#83A0BC]"
                          aria-label="Toggle video"
                        >
                          <Icon icon="mdi:video" className="text-[#1D154A] w-3 h-3" />
                        </button>
                      </div>
                    </div>
                    
                    {/* Show other participants - only if there are any */}
                    {(currentRoom?.participantCount || 0) > 1 ? (
                      currentRoom.participants.filter(participant => participant.userId !== currentUser?.uniqueId).map((participant, i) => (
                        <div key={participant.userId} className="relative bg-[#F0F0F0] rounded-lg border border-[#83A0BC] flex flex-col items-center justify-center">
                          <div className="w-12 h-12 rounded-full bg-[#83A0BC] bg-opacity-80 flex items-center justify-center mb-1">
                            {participant.avatarUrl ? (
                              <img 
                                src={participant.avatarUrl} 
                                alt={`${participant.username}'s avatar`}
                                className="w-12 h-12 rounded-full object-cover"
                              />
                            ) : (
                              <Icon icon="mdi:account" className="text-white w-8 h-8" />
                            )}
                          </div>
                          <p className="text-[#1D154A] text-xs font-mono">
                            {participant.username}
                          </p>
                          
                          {/* Video controls */}
                          <div className="absolute bottom-1 left-0 right-0 flex justify-center space-x-1">
                            <button 
                              className="bg-white rounded-full p-0.5 hover:bg-gray-100 border border-[#83A0BC]"
                              aria-label="Toggle microphone"
                            >
                              <Icon icon="mdi:microphone" className="text-[#1D154A] w-3 h-3" />
                            </button>
                            <button 
                              className="bg-blue-500 text-white rounded-full p-0.5 hover:bg-blue-600 border border-[#83A0BC]"
                              aria-label="Toggle video"
                            >
                              <Icon icon="mdi:video" className="text-[#1D154A] w-3 h-3" />
                            </button>
                          </div>
                        </div>
                      ))
                    ) : (
                      // Show "No other participants" message when there are no participants
                      <div className="col-span-2 flex items-center justify-center">
                        <div className="text-center">
                          <div className="w-12 h-12 rounded-full bg-gray-300 flex items-center justify-center mb-1">
                            <Icon icon="mdi:account-group" className="text-gray-500 w-8 h-8" />
                          </div>
                          <p className="text-[#1D154A] text-xs font-mono text-gray-500">No other participants</p>
                        </div>
                      </div>
                    )}
                    
                    {/* Show empty slots for remaining capacity */}
                    {Array.from({ length: Math.max(0, (currentRoom?.maxParticipants || 4) - (currentRoom?.participantCount || 1)) }).map((_, i) => (
                      <div key={`empty-${i}`} className="relative bg-[#F0F0F0] rounded-lg border border-[#83A0BC] flex flex-col items-center justify-center opacity-50">
                        <div className="w-12 h-12 rounded-full bg-gray-300 flex items-center justify-center mb-1">
                          <Icon icon="mdi:account-plus" className="text-gray-500 w-8 h-8" />
                        </div>
                        <p className="text-[#1D154A] text-xs font-mono text-gray-500">Empty Slot</p>
                      </div>
                    ))}
                  </div>
                  
                  {/* Video call controls */}
                  <div className="flex justify-center space-x-2 mt-1">
                    <button 
                      className="bg-green-500 text-white rounded-full p-2 hover:bg-green-600 flex items-center gap-2 px-4"
                      onClick={toggleVideoCall}
                      aria-label="Start video call"
                    >
                      <Icon icon="mdi:video" className="w-4 h-4" />
                      <span className="text-sm font-medium">Start Video Call</span>
                    </button>
                    <button 
                      className="bg-blue-500 text-white rounded-full p-2 hover:bg-blue-600"
                      onClick={toggleScreenSharing}
                      aria-label={screenSharing ? "Stop screen sharing" : "Start screen sharing"}
                    >
                      <Icon icon="mdi:monitor-share" className="w-4 h-4" />
                    </button>
                  </div>
                </>
              )}
            </div>
            
            {/* Side Panel - 25% of the width */}
            <div className="bg-white bg-opacity-75 rounded-lg shadow-sm border border-[#83A0BC] p-2 flex flex-col">
              <div className="flex mb-1">
                <button 
                  className={`flex-1 py-1 px-2 text-xs font-mono rounded-l-lg ${
                    activeTab === 'materials' 
                      ? 'bg-[#83A0BC] text-white' 
                      : 'bg-gray-200 text-[#1D154A] hover:bg-gray-300'
                  }`}
                  onClick={() => setActiveTab('materials')}
                  aria-label="View materials"
                        >
          Files
        </button>
                <button 
                  className={`flex-1 py-1 px-2 text-xs font-mono ${
                    activeTab === 'notes' 
                      ? 'bg-[#83A0BC] text-white' 
                      : 'bg-gray-200 text-[#1D154A] hover:bg-gray-300'
                  }`}
                  onClick={() => setActiveTab('notes')}
                  aria-label="View notes"
                >
                  Notes
                </button>
                <button 
                  className={`flex-1 py-1 px-2 text-xs font-mono rounded-r-lg ${
                    activeTab === 'todos' 
                      ? 'bg-[#83A0BC] text-white' 
                      : 'bg-gray-200 text-[#1D154A] hover:bg-gray-300'
                  }`}
                  onClick={() => setActiveTab('todos')}
                  aria-label="View todos"
                >
                  Todos
                </button>
              </div>
              
              {/* Share tab */}
              <button 
                className={`w-full py-1 px-2 text-xs font-mono rounded-lg mb-2 ${
                  activeTab === 'share' 
                    ? 'bg-[#83A0BC] text-white' 
                    : 'bg-gray-200 text-[#1D154A] hover:bg-gray-300'
                }`}
                onClick={() => setActiveTab('share')}
                            aria-label="Share files"
          >
            Share Files
              </button>
              
              {/* Chat tab */}
              <button 
                className={`w-full py-1 px-2 text-xs font-mono rounded-lg mb-2 ${
                  activeTab === 'chat' 
                    ? 'bg-[#83A0BC] text-white' 
                    : 'bg-gray-200 text-[#1D154A] hover:bg-gray-300'
                }`}
                onClick={() => setActiveTab('chat')}
                aria-label="View chat"
              >
                Chat
              </button>
              
              {/* Content based on active tab */}
              <div className="flex-1 overflow-auto hide-scrollbar">
                {activeTab === 'materials' && (
                  <div className="space-y-1">
                    {loading ? (
                      <div className="text-center py-4">
                        <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-[#1D154A] mx-auto mb-2"></div>
                        <p className="text-xs text-[#1D154A]">Loading files...</p>
                      </div>
                    ) : sharedMaterials.length === 0 ? (
                      <div className="text-center py-4">
                                    <p className="text-xs text-[#1D154A] opacity-60">No shared files yet</p>
            <p className="text-xs text-[#1D154A] opacity-40">Use the Share tab to add files</p>
                      </div>
                    ) : (
                      sharedMaterials.slice(0, 3).map((material) => (
                        <div key={material.id} className="bg-[#8E8E9A] bg-opacity-80 p-1 flex items-start border border-[#7A7A85] h-8">
                          <div className="flex items-start w-full">
                            <div className="bg-white p-0.5 flex items-center justify-center border border-[#7A7A85] h-6 w-6">
                              {material.type === 'pdf' && <span className="text-[#1D154A] font-bold text-xs">PDF</span>}
                              {material.type === 'figma' && (
                                <div className="flex gap-0.5">
                                  <div className="w-1.5 h-1.5 bg-[#1D154A] rounded-full"></div>
                                  <div className="w-1.5 h-1.5 bg-[#1D154A] rounded-full"></div>
                                  <div className="w-1.5 h-1.5 bg-[#1D154A] rounded-full flex items-center justify-center">
                                    <div className="w-0.5 h-0.5 bg-white rounded-full"></div>
                                  </div>
                                </div>
                              )}
                              {material.type === 'word' && <span className="text-[#1D154A] font-bold text-xs">W</span>}
                              {material.type === 'link' && <span className="text-[#1D154A] font-bold text-xs">🔗</span>}
                              {material.type === 'file' && <span className="text-[#1D154A] font-bold text-xs">📁</span>}
                              {!['pdf', 'figma', 'word', 'link', 'file'].includes(material.type) && <span className="text-[#1D154A] font-bold text-xs">F</span>}
                            </div>
                            <div className="flex-1 ml-1 text-white flex items-center">
                              <span className="text-xs font-mono truncate max-w-[90px] text-left">
                                {material.name}
                              </span>
                            </div>
                          </div>
                        </div>
                      ))
                    )}
                  </div>
                )}
                
                {activeTab === 'notes' && (
                  <div className="space-y-1">
                    {loading ? (
                      <div className="text-center py-4">
                        <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-[#1D154A] mx-auto mb-2"></div>
                        <p className="text-xs text-[#1D154A]">Loading notes...</p>
                      </div>
                    ) : sharedNotes.length === 0 ? (
                      <div className="text-center py-4">
                        <p className="text-xs text-[#1D154A] opacity-60">No shared notes yet</p>
                        <p className="text-xs text-[#1D154A] opacity-40">Use the Share tab to add notes</p>
                      </div>
                    ) : (
                      sharedNotes.slice(0, 3).map((note) => (
                        <div key={note.id} className="bg-[#83A0BC] bg-opacity-80 p-1 border border-[#6A8099]">
                          <p className="text-xs font-mono line-clamp-2">
                            {note.content}
                          </p>
                        </div>
                      ))
                    )}
                  </div>
                )}
                
                {activeTab === 'todos' && (
                  <div className="space-y-1">
                    {loading ? (
                      <div className="text-center py-4">
                        <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-[#1D154A] mx-auto mb-2"></div>
                        <p className="text-xs text-[#1D154A]">Loading todos...</p>
                      </div>
                    ) : sharedTodos.length === 0 ? (
                      <div className="text-center py-4">
                        <p className="text-xs text-[#1D154A] opacity-60">No shared todos yet</p>
                        <p className="text-xs text-[#1D154A] opacity-40">Use the Share tab to add todos</p>
                      </div>
                    ) : (
                      sharedTodos.slice(0, 3).map((todo) => (
                        <div key={todo.id} className="bg-[#C6AE82] bg-opacity-80 p-1 flex items-center border border-[#B09C70]">
                          <div className="flex items-start gap-1 flex-1">
                            <div className="w-3 h-3 flex-shrink-0 flex items-center justify-center mt-0.5">
                              <Icon icon="mdi:checkbox-outline" className="text-white w-3 h-3" />
                            </div>
                            <div className="flex-1 flex justify-between items-center">
                              <span className="text-white text-xs font-mono line-clamp-1">{todo.title}</span>
                              <button
                                onClick={() => handleMarkSharedTodoDone(todo.id)}
                                className="text-[10px] font-mono text-white bg-transparent whitespace-nowrap"
                                aria-label={`Mark ${todo.title} as done`}
                              >
                                Done
                              </button>
                            </div>
                          </div>
                        </div>
                      ))
                    )}
                  </div>
                )}

                {/* Share tab content */}
                {activeTab === 'share' && (
                  <div className="flex flex-col h-full">
                    <div className="mb-2">
                      <h4 className="text-xs font-mono text-[#1D154A] mb-2">Select files to share in this room:</h4>
                    </div>
                    
                                            {/* Personal Files Selection */}
                    <div className="flex-1 overflow-auto space-y-2 mb-2">
                      {/* Personal Todos */}
                      <div className="space-y-1">
                        <h5 className="text-xs font-mono text-[#1D154A] font-bold">Your Todos:</h5>
                        {todos.slice(0, 3).map((todo) => (
                          <label key={todo.id} className="flex items-center space-x-2 cursor-pointer">
                            <input
                              type="checkbox"
                              checked={selectedMaterialsForSharing.some(m => m.id === todo.id && m.type === 'todo')}
                              onChange={() => toggleMaterialSelection(todo, 'todo')}
                              className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                            />
                            <span className="text-xs text-[#1D154A] line-clamp-1">{todo.title}</span>
                          </label>
                        ))}
                      </div>
                      
                                              {/* Personal Files */}
                      <div className="space-y-1">
                        <h5 className="text-xs font-mono text-[#1D154A] font-bold">Your Files:</h5>
                        {userMaterials.slice(0, 3).map((material) => (
                          <label key={material.id} className="flex items-center space-x-2 cursor-pointer">
                            <input
                              type="checkbox"
                              checked={selectedMaterialsForSharing.some(m => m.id === material.id && m.type === 'material')}
                              onChange={() => toggleMaterialSelection(material, 'material')}
                              className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                            />
                            <span className="text-xs text-[#1D154A] line-clamp-1">{material.name}</span>
                          </label>
                        ))}
                      </div>
                      
                      {/* Personal Notes */}
                      <div className="space-y-1">
                        <h5 className="text-xs font-mono text-[#1D154A] font-bold">Your Notes:</h5>
                        {userNotes.slice(0, 3).map((note) => (
                          <label key={note.id} className="flex items-center space-x-2 cursor-pointer">
                            <label key={note.id} className="flex items-center space-x-2 cursor-pointer">
                              <input
                                type="checkbox"
                                checked={selectedMaterialsForSharing.some(m => m.id === note.id && m.type === 'note')}
                                onChange={() => toggleMaterialSelection(note, 'note')}
                                className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                              />
                              <span className="text-xs text-[#1D154A] line-clamp-1">{note.content}</span>
                            </label>
                          </label>
                        ))}
                      </div>
                    </div>
                    
                    {/* Share Button */}
                    <button
                      onClick={async () => {
                        await shareMaterialsToRoom();
                        refreshSharedData();
                      }}
                      disabled={selectedMaterialsForSharing.length === 0}
                      className="w-full bg-blue-500 hover:bg-blue-600 disabled:bg-gray-300 text-white py-2 px-3 rounded-lg text-xs transition-colors"
                      aria-label="Share selected files"
                    >
                      Share {selectedMaterialsForSharing.length} File{selectedMaterialsForSharing.length !== 1 ? 's' : ''} to Room
                    </button>
                  </div>
                )}

                {/* Chat tab content */}
                {activeTab === 'chat' && (
                  <div className="flex flex-col h-full">
                    {/* Chat messages */}
                    <div className="flex-1 overflow-auto space-y-2 mb-2">
                      {chatMessages.map((message) => (
                        <div key={message.id} className={`flex ${message.userId === currentUser?._id ? 'justify-end' : 'justify-start'}`}>
                          <div className={`max-w-[80%] p-2 rounded-lg text-xs ${
                            message.userId === currentUser?._id 
                              ? 'bg-blue-500 text-white' 
                              : 'bg-gray-200 text-gray-800'
                          }`}>
                            <div className="font-semibold text-[10px] mb-1">
                              {message.username || 'Unknown User'}
                            </div>
                            <div>{message.content}</div>
                            <div className="text-[10px] opacity-70 mt-1">
                              {new Date(message.timestamp).toLocaleTimeString()}
                            </div>
                          </div>
                        </div>
                      ))}
                      
                      {/* Typing indicators */}
                      {typingUsers.length > 0 && (
                        <div className="flex justify-start">
                          <div className="bg-gray-200 text-gray-800 p-2 rounded-lg text-xs">
                            <div className="flex items-center gap-1">
                              <span>Someone is typing</span>
                              <div className="flex gap-1">
                                <div className="w-1 h-1 bg-gray-500 rounded-full animate-bounce"></div>
                                <div className="w-1 h-1 bg-gray-500 rounded-full animate-bounce" style={{ animationDelay: '0.1s' }}></div>
                                <div className="w-1 h-1 bg-gray-500 rounded-full animate-bounce" style={{ animationDelay: '0.2s' }}></div>
                              </div>
                            </div>
                          </div>
                        </div>
                      )}
                    </div>
                    
                    {/* Chat input */}
                    <div className="flex gap-2">
                      <input
                        type="text"
                        value={newChatMessage}
                        onChange={(e) => setNewChatMessage(e.target.value)}
                        onKeyPress={(e) => e.key === 'Enter' && sendChatMessage()}
                        onFocus={() => handleTyping(true)}
                        onBlur={() => handleTyping(false)}
                        placeholder="Type a message..."
                        className="flex-1 border border-gray-300 rounded-lg px-2 py-1 text-xs focus:outline-none focus:ring-1 focus:ring-blue-500"
                        maxLength={500}
                      />
                      <button
                        onClick={sendChatMessage}
                        disabled={!newChatMessage.trim()}
                        className="bg-blue-500 text-white px-3 py-1 rounded-lg text-xs hover:bg-blue-600 disabled:bg-gray-300 transition-colors"
                        aria-label="Send message"
                      >
                        Send
                      </button>
                    </div>
                  </div>
                )}
              </div>
              
              {/* Invite button */}
              <div className="mt-1">
                <button 
                  onClick={handleInvite}
                  className="w-full bg-[#FFD169] text-[#1D154A] py-1 rounded-lg text-xs font-mono hover:bg-[#E5C090] transition-colors border border-[#E5C090]"
                  aria-label="Invite friends to shared space"
                >
                  INVITE FRIENDS
                </button>
              </div>
            </div>
          </div>
        </>
      )}
    </div>
  );
};

SharedSpace.propTypes = {
  userXP: PropTypes.number.isRequired,
  maxRoomParticipants: PropTypes.number.isRequired,
  showXPRequirement: PropTypes.func.isRequired,
  todos: PropTypes.array.isRequired,
  markTodoAsDone: PropTypes.func.isRequired,
  deleteTodo: PropTypes.func.isRequired,
  userMaterials: PropTypes.array.isRequired,
  deleteMaterial: PropTypes.func.isRequired,
  userNotes: PropTypes.array.isRequired,
  deleteNote: PropTypes.func.isRequired,
  handleAddNote: PropTypes.func.isRequired,
  newNote: PropTypes.string.isRequired,
  setNewNote: PropTypes.func.isRequired,
  setActiveSection: PropTypes.func.isRequired, // This is actually handleSectionChange now
  currentRoom: PropTypes.object,
  currentUser: PropTypes.object,
  wsConnection: PropTypes.object,
  chatMessages: PropTypes.array.isRequired,
  setChatMessages: PropTypes.func.isRequired,
  typingUsers: PropTypes.array.isRequired,
  setErrorDebounced: PropTypes.func.isRequired,
  selectedMaterialsForSharing: PropTypes.array.isRequired,
  toggleMaterialSelection: PropTypes.func.isRequired,
  shareMaterialsToRoom: PropTypes.func.isRequired,
  refreshRooms: PropTypes.func.isRequired,
  refreshTrigger: PropTypes.number.isRequired,
  onDeleteRoom: PropTypes.func,
  onLeaveRoom: PropTypes.func
};

export default SharedSpace; 