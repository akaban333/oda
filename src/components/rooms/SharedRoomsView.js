import React from 'react';
import { Icon } from '@iconify/react';
import PropTypes from 'prop-types';

const SharedRoomsView = ({
  userRooms,
  currentUser,
  setShowCreateRoomModal,
  setShowJoinModal,
  refreshRooms,
  showNotifications,
  setShowNotifications,
  fetchNotifications,
  notifications,
  showRoomInvitations,
  setShowRoomInvitations,
  fetchRoomInvitations,
  roomInvitations,
  testBackendConnectivity,
  roomsAPI,
  setRoomToDelete,
  setShowDeleteConfirm,
  setEditingRoom,
  setShowRoomSettings,
  handleInviteFriends,
  handleJoinRoom,
  handleStartCall,
  handleRoomSelection
}) => {
  return (
    <div className="h-full p-4" style={{ background: 'linear-gradient(135deg, #7681b3 0%, #8b96d1 50%, #a5b0e8 100%)', backgroundSize: 'cover' }}>
      {/* Header */}
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-white text-xl font-bold">Available Rooms</h2>
        <div className="flex gap-2">
          <button
            onClick={() => { setShowCreateRoomModal(true); }}
            className="bg-green-500 hover:bg-green-600 text-white px-4 py-2 rounded-lg flex items-center gap-2"
          >
            <Icon icon="mdi:plus" className="w-4 h-4" />
            Create Room
          </button>
          <button
            onClick={() => setShowJoinModal(true)}
            className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg flex items-center gap-2"
          >
            <Icon icon="mdi:plus" className="w-4 h-4" />
            Join Room
          </button>
          <button
            onClick={refreshRooms}
            className="bg-purple-500 hover:bg-purple-600 text-white px-4 py-2 rounded-lg flex items-center gap-2"
            title="Refresh rooms"
          >
            <Icon icon="mdi:refresh" className="w-4 h-4" />
            Refresh
          </button>

          <button
            onClick={() => {
              setShowNotifications(!showNotifications);
              if (!showNotifications) {
                fetchNotifications();
              }
            }}
            className="bg-yellow-500 hover:bg-yellow-600 text-white px-4 py-2 rounded-lg flex items-center gap-2 relative"
            title="Notifications (excluding room invitations)"
          >
            <Icon icon="mdi:bell" className="w-4 h-4" />
            {notifications.filter(n => !n.isRead && n.type !== 'room_invitation').length > 0 && (
              <span className="absolute -top-1 -right-1 bg-red-500 text-white text-xs rounded-full w-5 h-5 flex items-center justify-center">
                {notifications.filter(n => !n.isRead && n.type !== 'room_invitation').length}
              </span>
            )}
          </button>

          <button
            onClick={async () => {
              if (!showRoomInvitations) {
                fetchRoomInvitations();
              }
            }}
            className="bg-pink-500 hover:bg-pink-600 text-white px-4 py-2 rounded-lg flex items-center gap-2 relative"
            title="Room Invitations"
          >
            <Icon icon="mdi:account-plus" className="w-4 h-4" />
            {roomInvitations.length > 0 && (
              <span className="absolute -top-1 -right-1 bg-red-500 text-white text-xs rounded-full w-5 h-5 flex items-center justify-center">
                {roomInvitations.length}
              </span>
            )}
          </button>
        </div>
      </div>

      {/* Rooms Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {userRooms.map((room) => {
          const roomId = room.id || room._id;
          return (
            <div key={roomId} className="bg-white bg-opacity-20 backdrop-blur-md rounded-lg p-4 border border-white border-opacity-30 shadow-2xl shadow-black/30 hover:shadow-2xl hover:shadow-blue-400/20 transition-all duration-300 hover:scale-[1.02] hover:bg-white hover:bg-opacity-25">
              {/* Room Header */}
              <div className="flex justify-between items-center mb-3">
                <h3 className="text-white font-bold truncate text-lg drop-shadow-lg">{room.name}</h3>
                <div className="text-white text-sm bg-blue-500 bg-opacity-80 px-3 py-1.5 rounded-full shadow-lg border border-blue-400 border-opacity-50">
                  {room.participantCount || 0}/{room.maxParticipants || 4}
                </div>
              </div>

              {/* Room Description */}
              <p className="text-white text-opacity-80 text-sm mb-3 line-clamp-2">
                {room.description || 'No description'}
              </p>

              {/* Participants Preview */}
              <div className="flex items-center gap-2 mb-3">
                {/* Show participant count instead of fake avatars */}
                <div className="text-white text-sm">
                  <span className="font-medium">{room.participantCount || 0}</span> participant{(room.participantCount || 0) !== 1 ? 's' : ''}
                </div>
                {/* Show room owner info */}
                <div className="text-white text-xs bg-purple-500 bg-opacity-80 px-2.5 py-1 rounded-full shadow-md border border-purple-400 border-opacity-50">
                  Owner: {room.creatorId === currentUser?.uniqueId ? 'You' : (room.creatorUsername || 'Unknown')}
                </div>
              </div>

              {/* Action Buttons */}
              <div className="flex gap-2">
                <button
                  onClick={() => handleRoomSelection(room)}
                  className="flex-1 bg-blue-500 hover:bg-blue-600 text-white py-2 px-3 rounded-lg text-sm font-medium shadow-lg hover:shadow-xl hover:shadow-blue-500/30 transition-all duration-200 hover:scale-105"
                >
                  Enter Room
                </button>
                <button
                  onClick={() => handleStartCall(room)}
                  className="bg-green-500 hover:bg-green-600 text-white p-2 rounded-lg shadow-lg hover:shadow-xl hover:shadow-green-500/30 transition-all duration-200 hover:scale-105"
                >
                  <Icon icon="mdi:video" className="w-4 h-4" />
                </button>
                <button
                  onClick={() => handleInviteFriends(room)}
                  className="bg-blue-500 hover:bg-blue-600 text-white p-2 rounded-lg shadow-lg hover:shadow-xl hover:shadow-blue-500/30 transition-all duration-200 hover:scale-105"
                  aria-label="Invite friends to room"
                >
                  <Icon icon="mdi:account-plus" className="w-4 h-4" />
                </button>
                {room.creatorId === currentUser?.uniqueId && (
                  <>
                    <button
                      onClick={() => {
                        setEditingRoom(room);
                        setShowRoomSettings(true);
                      }}
                      className="bg-yellow-500 hover:bg-yellow-600 text-white p-2 rounded-lg shadow-lg hover:shadow-xl hover:shadow-yellow-500/30 transition-all duration-200 hover:scale-105"
                      aria-label="Edit room settings"
                    >
                      <Icon icon="mdi:cog" className="w-4 h-4" />
                    </button>
                    <button
                      onClick={() => {
                        setRoomToDelete(room);
                        setShowDeleteConfirm(true);
                      }}
                      className="bg-red-500 hover:bg-red-600 text-white p-2 rounded-lg shadow-lg hover:shadow-xl hover:shadow-red-500/30 transition-all duration-200 hover:scale-105"
                      aria-label="Delete room"
                    >
                      <Icon icon="mdi:trash" className="w-4 h-4" />
                    </button>
                  </>
                )}
              </div>
            </div>
          );
        })}
      </div>

      {userRooms.length === 0 && (
        <div className="text-center text-white text-opacity-60 mt-8">
          <Icon icon="mdi:account-group-outline" className="w-16 h-16 mx-auto mb-4" />
          <p className="text-lg mb-2">No rooms yet</p>
          <p className="text-sm">Create or join a room to start studying with others!</p>
        </div>
      )}
    </div>
  );
};

SharedRoomsView.propTypes = {
  userRooms: PropTypes.array.isRequired,
  currentUser: PropTypes.object,
  setShowCreateRoomModal: PropTypes.func.isRequired,
  setShowJoinModal: PropTypes.func.isRequired,
  refreshRooms: PropTypes.func.isRequired,
  showNotifications: PropTypes.bool.isRequired,
  setShowNotifications: PropTypes.func.isRequired,
  fetchNotifications: PropTypes.func.isRequired,
  notifications: PropTypes.array.isRequired,
  showRoomInvitations: PropTypes.bool.isRequired,
  setShowRoomInvitations: PropTypes.func.isRequired,
  fetchRoomInvitations: PropTypes.func.isRequired,
  roomInvitations: PropTypes.array.isRequired,
  testBackendConnectivity: PropTypes.func.isRequired,
  roomsAPI: PropTypes.object.isRequired,
  setRoomToDelete: PropTypes.func.isRequired,
  setShowDeleteConfirm: PropTypes.func.isRequired,
  setEditingRoom: PropTypes.func.isRequired,
  setShowRoomSettings: PropTypes.func.isRequired,
  handleInviteFriends: PropTypes.func.isRequired,
  handleJoinRoom: PropTypes.func.isRequired,
  handleStartCall: PropTypes.func.isRequired,
  handleRoomSelection: PropTypes.func.isRequired
};

export default SharedRoomsView; 