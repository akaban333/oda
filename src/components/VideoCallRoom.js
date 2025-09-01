import React, { useState, useEffect, useRef, useCallback } from 'react';
import { Icon } from '@iconify/react';
import { realtimeAPI } from '../utils/api';

const VideoCallRoom = ({ roomId, currentUser, onClose }) => {
  // WebSocket and real-time state
  const [ws, setWs] = useState(null);
  const [isConnected, setIsConnected] = useState(false);
  const [onlineUsers, setOnlineUsers] = useState([]);
  const [chatMessages, setChatMessages] = useState([]);
  const [newMessage, setNewMessage] = useState('');
  
  // Video call state
  const [isInCall, setIsInCall] = useState(false);
  const [isVideoEnabled, setIsVideoEnabled] = useState(true);
  const [isAudioEnabled, setIsAudioEnabled] = useState(true);
  const [isScreenSharing, setIsScreenSharing] = useState(false);
  const [remoteStreams, setRemoteStreams] = useState(new Map());
  
  // WebRTC state
  const localVideoRef = useRef(null);
  const localStreamRef = useRef(null);
  const peerConnectionsRef = useRef(new Map());
  const screenStreamRef = useRef(null);
  
  // Chat scroll reference
  const chatEndRef = useRef(null);

  // WebSocket connection
  useEffect(() => {
    const connectWebSocket = () => {
      try {
        const websocket = realtimeAPI.createWebSocket();
        
        websocket.onopen = () => {
          setIsConnected(true);
          setWs(websocket);
          
          // Join the room
          websocket.send(JSON.stringify({
            type: 'join_room',
            roomId: roomId,
            user: currentUser
          }));
        };

        websocket.onmessage = (event) => {
          const data = JSON.parse(event.data);
          handleWebSocketMessage(data);
        };

        websocket.onclose = () => {
          setIsConnected(false);
          setWs(null);
          // Attempt to reconnect after 3 seconds
          setTimeout(connectWebSocket, 3000);
        };

        websocket.onerror = (error) => {
          console.error('WebSocket error:', error);
          setIsConnected(false);
        };

      } catch (error) {
        console.error('Failed to create WebSocket connection:', error);
      }
    };

    connectWebSocket();

    return () => {
      if (ws) {
        ws.close();
      }
    };
  }, [roomId, currentUser]);

  // Handle WebSocket messages
  const handleWebSocketMessage = useCallback((data) => {
    switch (data.type) {
      case 'chat':
        setChatMessages(prev => [...prev, {
          id: Date.now(),
          user: data.user,
          message: data.message,
          timestamp: new Date()
        }]);
        break;
        
      case 'user_joined':
        setOnlineUsers(prev => {
          if (!prev.find(u => u.id === data.user.id)) {
            return [...prev, data.user];
          }
          return prev;
        });
        break;
        
      case 'user_left':
        setOnlineUsers(prev => prev.filter(u => u.id !== data.userId));
        break;
        
      case 'rtc_offer':
      case 'rtc_answer':
      case 'rtc_candidate':
        handleRTCSignaling(data);
        break;
        
      case 'call_started':
        setIsInCall(true);
        break;
        
      case 'call_ended':
        setIsInCall(false);
        endCall();
        break;
        
      default:
        console.log('Unknown message type:', data.type);
    }
  }, []);

  // Get user media
  const getUserMedia = useCallback(async (constraints) => {
    try {
      const stream = await navigator.mediaDevices.getUserMedia(constraints);
      localStreamRef.current = stream;
      
      // Ensure the video element is properly set up
      if (localVideoRef.current) {
        localVideoRef.current.srcObject = stream;
        // Force the video to load and play
        localVideoRef.current.load();
        localVideoRef.current.play().catch(e => console.log('Auto-play prevented:', e));
      }
      
      return stream;
    } catch (error) {
      console.error('Error getting user media:', error);
      throw error;
    }
  }, []);

  // Start video call
  const startCall = useCallback(async () => {
    try {
      const stream = await getUserMedia({ video: isVideoEnabled, audio: isAudioEnabled });
      setIsInCall(true);
      
      // Ensure video is visible immediately
      if (localVideoRef.current && stream) {
        localVideoRef.current.srcObject = stream;
        localVideoRef.current.style.display = 'block';
      }
      
      if (ws) {
        ws.send(JSON.stringify({
          type: 'start_call',
          roomId: roomId
        }));
      }
    } catch (error) {
      console.error('Failed to start call:', error);
      alert('Failed to access camera/microphone. Please check permissions.');
    }
  }, [getUserMedia, isVideoEnabled, isAudioEnabled, ws, roomId]);

  // End video call
  const endCall = useCallback(() => {
    console.log('Ending call and cleaning up media streams...');
    
    // Stop all tracks in the local stream
    if (localStreamRef.current) {
      const tracks = localStreamRef.current.getTracks();
      tracks.forEach(track => {
        console.log('Stopping track:', track.kind, track.id);
        track.stop();
      });
      localStreamRef.current = null;
    }
    
    // Stop screen sharing if active
    if (screenStreamRef.current) {
      const tracks = screenStreamRef.current.getTracks();
      tracks.forEach(track => {
        console.log('Stopping screen share track:', track.kind, track.id);
        track.stop();
      });
      screenStreamRef.current = null;
    }
    
    // Close peer connections
    peerConnectionsRef.current.forEach(pc => {
      pc.close();
    });
    peerConnectionsRef.current.clear();
    
    // Clear remote streams
    setRemoteStreams(new Map());
    
    // Reset call state
    setIsInCall(false);
    setIsScreenSharing(false);
    
    // Clear the video element's srcObject to release the camera
    if (localVideoRef.current) {
      localVideoRef.current.srcObject = null;
      localVideoRef.current.style.display = 'none';
    }
    
    // Notify server
    if (ws) {
      ws.send(JSON.stringify({
        type: 'end_call',
        roomId: roomId
      }));
    }
    
    console.log('Call ended and cleanup completed');
  }, [ws, roomId]);

  // Toggle video
  const toggleVideo = useCallback(() => {
    if (localStreamRef.current) {
      const videoTrack = localStreamRef.current.getVideoTracks()[0];
      if (videoTrack) {
        videoTrack.enabled = !videoTrack.enabled;
        setIsVideoEnabled(videoTrack.enabled);
        
        // Show/hide video element based on state
        if (localVideoRef.current) {
          localVideoRef.current.style.display = videoTrack.enabled ? 'block' : 'none';
        }
      }
    }
  }, []);

  // Toggle audio
  const toggleAudio = useCallback(() => {
    if (localStreamRef.current) {
      const audioTrack = localStreamRef.current.getAudioTracks()[0];
      if (audioTrack) {
        audioTrack.enabled = !audioTrack.enabled;
        setIsAudioEnabled(audioTrack.enabled);
      }
    }
  }, []);

  // Start screen sharing
  const startScreenShare = useCallback(async () => {
    try {
      const screenStream = await navigator.mediaDevices.getDisplayMedia({
        video: true,
        audio: true
      });
      
      screenStreamRef.current = screenStream;
      setIsScreenSharing(true);
      
      // Replace video track in peer connections
      const videoTrack = screenStream.getVideoTracks()[0];
      peerConnectionsRef.current.forEach(pc => {
        const sender = pc.getSenders().find(s => 
          s.track && s.track.kind === 'video'
        );
        if (sender) {
          sender.replaceTrack(videoTrack);
        }
      });
      
      videoTrack.onended = () => {
        stopScreenShare();
      };
      
    } catch (error) {
      console.error('Error starting screen share:', error);
    }
  }, []);

  // Stop screen sharing
  const stopScreenShare = useCallback(() => {
    if (screenStreamRef.current) {
      const tracks = screenStreamRef.current.getTracks();
      tracks.forEach(track => {
        console.log('Stopping screen share track:', track.kind, track.id);
        track.stop();
      });
      screenStreamRef.current = null;
    }
    
    setIsScreenSharing(false);
    
    // Switch back to camera if in call
    if (isInCall && localStreamRef.current) {
      const videoTrack = localStreamRef.current.getVideoTracks()[0];
      peerConnectionsRef.current.forEach(pc => {
        const sender = pc.getSenders().find(s => 
          s.track && s.track.kind === 'video'
        );
        if (sender && videoTrack) {
          sender.replaceTrack(videoTrack);
        }
      });
    }
  }, [isInCall]);

  // Handle RTC signaling
  const handleRTCSignaling = useCallback((data) => {
    // This is a simplified WebRTC implementation
    // In a production app, you'd want a more robust implementation
    console.log('RTC signaling received:', data);
  }, []);

  // Send chat message
  const sendMessage = useCallback(() => {
    if (newMessage.trim() && ws) {
      ws.send(JSON.stringify({
        type: 'chat',
        roomId: roomId,
        message: newMessage.trim()
      }));
      setNewMessage('');
    }
  }, [newMessage, ws, roomId]);

  // Auto-scroll chat to bottom
  useEffect(() => {
    chatEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [chatMessages]);

  // Load chat history on mount
  useEffect(() => {
    const loadChatHistory = async () => {
      try {
        const data = await realtimeAPI.getChatHistory(roomId);
        setChatMessages(data.messages || []);
      } catch (error) {
        console.error('Failed to load chat history:', error);
      }
    };
    
    loadChatHistory();
  }, [roomId]);

  // Cleanup function for component unmount
  useEffect(() => {
    return () => {
      console.log('VideoCallRoom component unmounting, cleaning up...');
      // Ensure all media streams are stopped when component unmounts
      if (localStreamRef.current) {
        const tracks = localStreamRef.current.getTracks();
        tracks.forEach(track => {
          console.log('Stopping track on unmount:', track.kind, track.id);
          track.stop();
        });
        localStreamRef.current = null;
      }
      
      if (screenStreamRef.current) {
        const tracks = screenStreamRef.current.getTracks();
        tracks.forEach(track => {
          console.log('Stopping screen share track on unmount:', track.kind, track.id);
          track.stop();
        });
        screenStreamRef.current = null;
      }
      
      // Close peer connections
      peerConnectionsRef.current.forEach(pc => pc.close());
      peerConnectionsRef.current.clear();
      
      // Clear video element
      if (localVideoRef.current) {
        localVideoRef.current.srcObject = null;
      }
    };
  }, []);

  return (
    <div className="fixed inset-0 bg-black bg-opacity-90 z-50 flex">
      {/* Video Area */}
      <div className="flex-1 relative">
        {/* Local video - always show when in call, even if video is disabled */}
        {isInCall && (
          <video
            ref={localVideoRef}
            autoPlay
            muted
            playsInline
            className="absolute top-4 right-4 w-48 h-36 bg-gray-800 rounded-lg border-2 border-white z-10"
            style={{ display: isVideoEnabled ? 'block' : 'none' }}
          />
        )}
        
        {/* Remote videos grid */}
        <div className="grid grid-cols-2 gap-2 h-full p-4">
          {Array.from(remoteStreams.entries()).map(([userId, stream]) => (
            <video
              key={userId}
              autoPlay
              playsInline
              className="w-full h-full bg-gray-800 rounded-lg"
              srcObject={stream}
            />
          ))}
          
          {/* Placeholder when no remote streams */}
          {remoteStreams.size === 0 && (
            <div className="col-span-2 flex items-center justify-center">
              <div className="text-white text-center">
                <Icon icon="mdi:video-off" className="w-16 h-16 mx-auto mb-4 opacity-50" />
                <p className="text-lg">No one else in the call</p>
                <p className="text-sm opacity-75">Invite others to join!</p>
              </div>
            </div>
          )}
        </div>
        
        {/* Call controls */}
        <div className="absolute bottom-4 left-1/2 transform -translate-x-1/2 flex gap-4">
          {!isInCall ? (
            <button
              onClick={startCall}
              className="bg-green-500 hover:bg-green-600 text-white px-6 py-3 rounded-full flex items-center gap-2"
            >
              <Icon icon="mdi:video" className="w-5 h-5" />
              Start Call
            </button>
          ) : (
            <>
              <button
                onClick={toggleVideo}
                className={`${isVideoEnabled ? 'bg-gray-600' : 'bg-red-500'} hover:bg-gray-700 text-white p-3 rounded-full`}
              >
                <Icon icon={isVideoEnabled ? "mdi:video" : "mdi:video-off"} className="w-5 h-5" />
              </button>
              
              <button
                onClick={toggleAudio}
                className={`${isAudioEnabled ? 'bg-gray-600' : 'bg-red-500'} hover:bg-gray-700 text-white p-3 rounded-full`}
              >
                <Icon icon={isAudioEnabled ? "mdi:microphone" : "mdi:microphone-off"} className="w-5 h-5" />
              </button>
              
              <button
                onClick={isScreenSharing ? stopScreenShare : startScreenShare}
                className={`${isScreenSharing ? 'bg-blue-500' : 'bg-gray-600'} hover:bg-gray-700 text-white p-3 rounded-full`}
              >
                <Icon icon="mdi:monitor-share" className="w-5 h-5" />
              </button>
              
              <button
                onClick={endCall}
                className="bg-red-500 hover:bg-red-600 text-white p-3 rounded-full"
              >
                <Icon icon="mdi:phone-hangup" className="w-5 h-5" />
              </button>
            </>
          )}
        </div>
      </div>
      
      {/* Chat Panel */}
      <div className="w-80 bg-white flex flex-col">
        {/* Chat header */}
        <div className="bg-gray-100 p-4 border-b flex justify-between items-center">
          <h3 className="font-bold text-gray-800">Room Chat</h3>
          <div className="flex items-center gap-2">
            <div className={`w-2 h-2 rounded-full ${isConnected ? 'bg-green-500' : 'bg-red-500'}`}></div>
            <span className="text-sm text-gray-600">
              {onlineUsers.length} online
            </span>
            <button
              onClick={onClose}
              className="text-gray-600 hover:text-gray-800"
            >
              <Icon icon="mdi:close" className="w-5 h-5" />
            </button>
          </div>
        </div>
        
        {/* Online users */}
        <div className="p-2 border-b bg-gray-50">
          <div className="flex flex-wrap gap-1">
            {onlineUsers.map(user => (
              <span
                key={user.id}
                className="bg-blue-100 text-blue-800 px-2 py-1 rounded-full text-xs"
              >
                {user.username || user.displayName}
              </span>
            ))}
          </div>
        </div>
        
        {/* Chat messages */}
        <div className="flex-1 overflow-y-auto p-4 space-y-2">
          {chatMessages.map(message => (
            <div key={message.id} className="mb-2">
              <div className="flex items-center gap-2 mb-1">
                <span className="font-semibold text-sm text-gray-700">
                  {message.user?.username || message.user?.displayName || 'Unknown'}
                </span>
                <span className="text-xs text-gray-500">
                  {new Date(message.timestamp).toLocaleTimeString()}
                </span>
              </div>
              <p className="text-gray-800 text-sm">{message.message}</p>
            </div>
          ))}
          <div ref={chatEndRef} />
        </div>
        
        {/* Chat input */}
        <div className="p-4 border-t">
          <div className="flex gap-2">
            <input
              type="text"
              value={newMessage}
              onChange={(e) => setNewMessage(e.target.value)}
              onKeyPress={(e) => e.key === 'Enter' && sendMessage()}
              placeholder="Type a message..."
              className="flex-1 border rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              disabled={!isConnected}
            />
            <button
              onClick={sendMessage}
              disabled={!isConnected || !newMessage.trim()}
              className="bg-blue-500 hover:bg-blue-600 disabled:bg-gray-300 text-white px-4 py-2 rounded-lg"
            >
              <Icon icon="mdi:send" className="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default VideoCallRoom; 