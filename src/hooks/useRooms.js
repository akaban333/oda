import { useState, useCallback, useEffect } from 'react';
import { roomsAPI } from '../utils/api';

export const useRooms = (currentUser) => {
  const [userRooms, setUserRooms] = useState([]);
  const [currentRoom, setCurrentRoom] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Function to refresh rooms data
  const refreshRooms = useCallback(async () => {
    try {
      const roomsData = await roomsAPI.getRooms();
      if (roomsData?.rooms) {
        setUserRooms(roomsData.rooms);
      }
    } catch (error) {
      setError('Failed to refresh rooms. Please try again.');
    }
  }, []);

  // Fetch rooms on mount and when currentUser changes
  useEffect(() => {
    if (currentUser) {
      refreshRooms();
    }
  }, [currentUser, refreshRooms]);

  // Room management functions
  const createRoom = useCallback(async (roomData) => {
    try {
      const newRoom = await roomsAPI.createRoom(roomData);
      return newRoom.room || newRoom;
    } catch (error) {
      throw error;
    }
  }, []);

  const joinRoom = useCallback(async (invitationCode) => {
    try {
      const room = await roomsAPI.joinRoomByCode({ InvitationCode: invitationCode });
      
      // Check if room already exists to prevent duplicates
      const roomExists = userRooms.some(existingRoom => existingRoom.id === room.room.id);
      if (!roomExists) {
        setUserRooms(prev => [...prev, room.room]);
      }
      
      setCurrentRoom(room.room);
      return room.room;
    } catch (error) {
      console.error('Error joining room:', error);
      throw error;
    }
  }, [userRooms]);

  const leaveRoom = useCallback(async (roomId) => {
    try {
      await roomsAPI.leaveRoom(roomId);
      setUserRooms(prev => prev.filter(room => (room.id || room._id) !== roomId));
      if (currentRoom && (currentRoom.id || currentRoom._id) === roomId) {
        setCurrentRoom(null);
      }
    } catch (error) {
      console.error('Error leaving room:', error);
      throw error;
    }
  }, [currentRoom]);

  const generateInvitationCode = useCallback(async (roomId) => {
    try {
      const data = await roomsAPI.generateInvitationCode(roomId);
      return data.invitationCode;
    } catch (error) {
      console.error('Error generating invitation code:', error);
      throw error;
    }
  }, []);

  const handleEditRoom = useCallback(async (roomData) => {
    try {
      await roomsAPI.updateRoom(roomData._id, roomData);
      setUserRooms(prev => prev.map(room => 
        (room.id || room._id) === roomData.id ? { ...room, ...roomData } : room
      ));
      setError(null);
    } catch (error) {
      console.error('Error updating room:', error);
      setError('Failed to update room. Please try again.');
    }
  }, []);

  const handleDeleteRoom = useCallback(async (roomToDelete) => {
    if (!roomToDelete) return;
    
    try {
      await roomsAPI.deleteRoom(roomToDelete.id || roomToDelete._id);
      setUserRooms(prev => prev.filter(room => (room.id || room._id) !== (roomToDelete.id || roomToDelete._id)));
      
      if (currentRoom && (currentRoom.id || currentRoom._id) === (roomToDelete.id || roomToDelete._id)) {
        setCurrentRoom(null);
      }
      
      setError(null);
    } catch (error) {
      console.error('Error deleting room:', error);
      setError('Failed to delete room. Please try again.');
    }
  }, [currentRoom]);

  const inviteFriendToRoom = useCallback(async (friendId, roomId) => {
    try {
      // Check XP requirements for additional participants
      const room = userRooms.find(r => (r.id || r._id) === roomId);
      if (!room) {
        setError('Room not found');
        return;
      }

      const currentParticipants = room.participantCount || 0;
      const maxAllowed = room.maxParticipants || 4;
      
      if (currentParticipants >= maxAllowed) {
        setError('Room is at maximum capacity');
        return;
      }

      // Invite friend to room (friendId is uniqueId)
      await roomsAPI.inviteUserToRoom(roomId, friendId);
      
      // Refresh rooms data
      const roomsData = await roomsAPI.getRooms();
      setUserRooms(roomsData.rooms || []);
      
      setError(null);
    } catch (error) {
      console.error('Error inviting friend:', error);
      setError('Failed to invite friend to room');
    }
  }, [userRooms]);

  const handleAcceptRoomInvitation = useCallback(async (roomId) => {
    try {
      await roomsAPI.acceptRoomInvitation(roomId);
      // Refresh rooms and room invitations
      refreshRooms();
      setError(null);
    } catch (error) {
      console.error('Error accepting room invitation:', error);
      setError('Failed to accept room invitation');
    }
  }, [refreshRooms]);

  return {
    userRooms,
    setUserRooms,
    currentRoom,
    setCurrentRoom,
    loading,
    setLoading,
    error,
    setError,
    refreshRooms,
    createRoom,
    joinRoom,
    leaveRoom,
    generateInvitationCode,
    handleEditRoom,
    handleDeleteRoom,
    inviteFriendToRoom,
    handleAcceptRoomInvitation
  };
}; 