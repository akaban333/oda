import React from 'react';
import { Icon } from '@iconify/react';
import PropTypes from 'prop-types';

const RoomModals = ({
  // Todo Modal
  showTodoModal,
  setShowTodoModal,
  newTodo,
  setNewTodo,
  addTodo,
  
  // Material Modal
  showMaterialModal,
  setShowMaterialModal,
  newMaterialType,
  setNewMaterialType,
  newMaterialUrl,
  setNewMaterialUrl,
  uploadedFile,
  setUploadedFile,
  addMaterial,
  
  // Session XP Modal
  showSessionXPModal,
  setShowSessionXPModal,
  sessionXP,
  
  // Video Call Modal
  showVideoCall,
  setShowVideoCall,
  currentRoom,
  currentUser,
  VideoCallRoom,
  
  // Join Room Modal
  showJoinModal,
  setShowJoinModal,
  joinRoomCode,
  setJoinRoomCode,
  joinRoom,
  
  // Create Room Modal
  showCreateRoomModal,
  setShowCreateRoomModal,
  newRoomData,
  setNewRoomData,
  createRoomError,
  setCreateRoomError,
  handleCreateRoom,
  maxRoomParticipants,
  userXP,
  
  // Room Settings Modal
  showRoomSettings,
  setShowRoomSettings,
  editingRoom,
  setEditingRoom,
  handleEditRoom,
  generateInvitationCode,
  errorDisplay,
  
  // Delete Room Confirmation Modal
  showDeleteConfirm,
  setShowDeleteConfirm,
  roomToDelete,
  setRoomToDelete,
  handleDeleteRoom,
  onDeleteSuccess,
  
  // Friends Invitation Modal
  showInviteFriendsModal,
  setShowInviteFriendsModal,
  invitingToRoom,
  setInvitingToRoom,
  friendsList,
  inviteFriendToRoom,
  userRooms,
  userXP: userXPProp,
  
  // Notifications Modal
  showNotifications,
  setShowNotifications,
  notifications,
  
  // Room Invitations Modal
  showRoomInvitations,
  setShowRoomInvitations,
  roomInvitations,
  handleAcceptRoomInvitation,
  
  // Share Modal
  showShareModal,
  setShowShareModal,
  sharingItem,
  setSharingItem,
  selectedRoomsForSharing,
  toggleRoomSelection,
  shareItemToRooms,
  sharedRooms
}) => {
  return (
    <>
      {/* Todo Modal - Accessible dialog */}
      {showTodoModal && (
        <div 
          className="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
          role="dialog"
          aria-labelledby="add-todo-title"
          aria-describedby="add-todo-description"
          aria-modal="true"
        >
          <div className="bg-[#C6AE82] p-4 rounded-lg w-[350px] max-w-full">
            <h3 id="add-todo-title" className="text-white font-bold mb-3 uppercase text-sm">Add New Task</h3>
            <p id="add-todo-description" className="sr-only">Enter a new task in the field below and click Add Task to submit.</p>
            
            <input 
              type="text" 
              placeholder="Enter your task..." 
              className="w-full p-2 mb-3 border border-[#3A3E5B] text-sm"
              value={newTodo}
              onChange={(e) => setNewTodo(e.target.value)}
              onKeyPress={(e) => e.key === 'Enter' && addTodo()}
              maxLength={150}
              aria-label="Task description"
            />
            <div className="flex justify-end gap-2">
              <button 
                className="bg-white text-[#5A554B] px-3 py-1 rounded hover:bg-gray-100 text-sm"
                onClick={() => setShowTodoModal(false)}
                aria-label="Cancel adding task"
              >
                Cancel
              </button>
              <button 
                className="bg-white text-[#5A554B] px-3 py-1 rounded hover:bg-[#d5b581] text-sm"
                onClick={addTodo}
                disabled={!newTodo.trim()}
                aria-label="Add task"
              >
                Add Task
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Material Modal - Accessible dialog */}
      {showMaterialModal && (
        <div 
          className="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
          role="dialog"
          aria-labelledby="add-material-title"
          aria-describedby="add-material-description"
          aria-modal="true"
        >
          <div className="bg-[#8E8E9A] p-4 rounded-lg w-[400px] max-w-full">
            <h3 id="add-material-title" className="text-white font-bold mb-3 uppercase text-sm">Add New File</h3>
            <p id="add-material-description" className="sr-only">Choose between adding a file by link or uploading a file.</p>
            
            <div className="mb-3">
              <div className="flex gap-3 mb-2">
                <button 
                  className={`px-3 py-1 border border-[#3A3A5B] text-sm ${newMaterialType === 'link' ? 'bg-[#727279] text-white' : 'bg-white text-[#3A3A5B]'}`}
                  onClick={() => setNewMaterialType('link')}
                  aria-pressed={newMaterialType === 'link'}
                  aria-label="Add by link"
                >
                  Add by Link
                </button>
                <button 
                  className={`px-3 py-1 border border-[#3A3A5B] text-sm ${newMaterialType === 'upload' ? 'bg-[#727279] text-white' : 'bg-white text-[#3A3A5B]'}`}
                  onClick={() => setNewMaterialType('upload')}
                  aria-pressed={newMaterialType === 'upload'}
                  aria-label="Upload file"
                >
                  Upload File
                </button>
              </div>
              
              {newMaterialType === 'link' ? (
                <input 
                  type="text" 
                  placeholder="Enter file URL..." 
                  className="w-full p-2 border border-[#3A3A5B] text-sm"
                  value={newMaterialUrl}
                  onChange={(e) => setNewMaterialUrl(e.target.value)}
                  aria-label="Material URL"
                />
              ) : (
                <div className="border-2 border-dashed border-[#3A3A5B] p-3 text-center">
                  {uploadedFile ? (
                    <div>
                      <p className="text-white mb-2 text-sm">File selected: {uploadedFile.name}</p>
                      <button 
                        className="text-white underline text-sm"
                        onClick={() => setUploadedFile(null)}
                        aria-label="Remove selected file"
                      >
                        Remove
                      </button>
                    </div>
                  ) : (
                    <>
                      <p className="text-white mb-2 text-sm">Drag and drop a file or click to select</p>
                      <input 
                        type="file" 
                        className="hidden" 
                        id="file-upload"
                        onChange={(e) => e.target.files && setUploadedFile(e.target.files[0])}
                        aria-label="Upload file"
                      />
                      <label 
                        htmlFor="file-upload" 
                        className="bg-white text-[#3A3A5B] px-3 py-1 border border-[#3A3A5B] cursor-pointer inline-block text-sm"
                      >
                        Select File
                      </label>
                    </>
                  )}
                </div>
              )}
            </div>
            
            <div className="flex justify-end gap-2">
              <button 
                className="bg-white text-[#3A3A5B] px-3 py-1 rounded hover:bg-gray-100 text-sm"
                onClick={() => setShowMaterialModal(false)}
                aria-label="Cancel adding file"
              >
                Cancel
              </button>
              <button 
                className="bg-[#727279] text-white px-3 py-1 rounded hover:bg-[#62626a] text-sm"
                onClick={addMaterial}
                disabled={(newMaterialType === 'link' && !newMaterialUrl) || (newMaterialType === 'upload' && !uploadedFile)}
                                  aria-label="Add file"
              >
                Add File
              </button>
            </div>
          </div>
        </div>
      )}
      
      {/* Session XP Modal - Accessible dialog */}
      {showSessionXPModal && (
        <div 
          className="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
          role="dialog"
          aria-labelledby="session-complete-title"
          aria-describedby="session-complete-description"
          aria-modal="true"
        >
          <div className="bg-white p-4 rounded-lg shadow-xl max-w-xs w-full">
            <h3 id="session-complete-title" className="text-[#1D154A] text-lg font-bold mb-2">Session Complete!</h3>
            <p id="session-complete-description" className="text-gray-600 mb-2 text-sm">Your study session has ended.</p>
            
            <div className="bg-[#F9D2B4] p-3 rounded-lg mb-3">
              <div className="flex justify-between items-center mb-1">
                <span className="text-[#7B5834] text-sm">XP Earned:</span>
                <span className="text-[#7B5834] font-bold text-sm">{sessionXP} XP</span>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-[#7B5834] text-sm">Session Duration:</span>
                <span className="text-[#7B5834] text-sm">{Math.floor(sessionXP / 2)} minutes</span>
              </div>
            </div>
            
            <button 
              onClick={() => setShowSessionXPModal(false)} 
              className="w-full bg-[#83A0BC] text-white py-2 rounded-lg text-sm"
              aria-label="Close session summary"
            >
              Close
            </button>
          </div>
        </div>
      )}
      
      {/* Video Call Modal */}
      {showVideoCall && currentRoom && currentUser && (
        <VideoCallRoom
          roomId={currentRoom.id || currentRoom._id}
          currentUser={currentUser}
          onClose={() => setShowVideoCall(false)}
        />
      )}

      {/* Join Room Modal */}
      {showJoinModal && (
        <div className="fixed inset-0 bg-black bg-opacity-60 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 max-w-md w-full mx-4">
            <h3 className="text-xl font-bold mb-4">Join Room</h3>
            <p className="text-gray-600 mb-4">Enter the invitation code to join a room:</p>
            <input
              type="text"
              value={joinRoomCode}
              onChange={(e) => setJoinRoomCode(e.target.value)}
              placeholder="Enter invitation code"
              className="w-full border border-gray-300 rounded-lg px-3 py-2 mb-4 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <div className="flex gap-2">
              <button
                onClick={() => setShowJoinModal(false)}
                className="flex-1 bg-gray-500 hover:bg-gray-600 text-white py-2 px-4 rounded-lg"
              >
                Cancel
              </button>
              <button
                onClick={async () => {
                  try {
                    await joinRoom(joinRoomCode);
                    setJoinRoomCode('');
                    setShowJoinModal(false);
                  } catch (error) {
                    alert('Failed to join room. Please check the invitation code.');
                  }
                }}
                disabled={!joinRoomCode.trim()}
                className="flex-1 bg-blue-500 hover:bg-blue-600 disabled:bg-gray-300 text-white py-2 px-4 rounded-lg"
              >
                Join Room
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Create Room Modal */}
      {showCreateRoomModal && (
        <div className="fixed inset-0 bg-black bg-opacity-60 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 max-w-lg w-full mx-4 max-h-[90vh] overflow-y-auto">
            <h3 className="text-xl font-bold mb-4">Create New Room</h3>
            
            <div className="space-y-4">
              {/* Room Name */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Room Name *
                </label>
                <input
                  type="text"
                  value={newRoomData.name}
                  onChange={(e) => setNewRoomData(prev => ({ ...prev, name: e.target.value }))}
                  placeholder="Enter room name"
                  className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  maxLength={50}
                />
              </div>

              {/* Room Description */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Description
                </label>
                <textarea
                  value={newRoomData.description}
                  onChange={(e) => setNewRoomData(prev => ({ ...prev, description: e.target.value }))}
                  placeholder="Enter room description"
                  rows={3}
                  className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  maxLength={200}
                />
              </div>

              {/* Max Participants with XP-based limits */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Max Participants
                </label>
                <select
                  value={newRoomData.maxParticipants}
                  onChange={(e) => setNewRoomData(prev => ({ ...prev, maxParticipants: parseInt(e.target.value) }))}
                  className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  {Array.from({ length: maxRoomParticipants }, (_, i) => i + 1).map(num => (
                    <option key={num} value={num}>{num}</option>
                  ))}
                </select>
                <p className="text-xs text-gray-500 mt-1">
                  Your XP level allows up to {maxRoomParticipants} participants
                  {userXP < 300 && (
                    <span className="text-blue-600">
                      . Need 300 XP for 5+ participants, 600 XP for 6+ participants, etc.
                    </span>
                  )}
                </p>
              </div>
            </div>

            {/* Error display (scoped to create room) */}
            {createRoomError && (
              <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
                {createRoomError}
              </div>
            )}

            {/* Action buttons */}
            <div className="flex gap-2 mt-6">
              <button
                onClick={() => {
                  setShowCreateRoomModal(false);
                  setNewRoomData({ name: '', description: '', maxParticipants: 4 });
                  setCreateRoomError(null);
                }}
                className="flex-1 bg-gray-500 hover:bg-gray-600 text-white py-2 px-4 rounded-lg transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={handleCreateRoom}
                disabled={!newRoomData.name.trim()}
                className="flex-1 bg-blue-500 hover:bg-blue-600 disabled:bg-gray-300 text-white py-2 px-4 rounded-lg transition-colors"
              >
                Create Room
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Room Settings Modal */}
      {showRoomSettings && editingRoom && (
        <div className="fixed inset-0 bg-black bg-opacity-60 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 max-w-lg w-full mx-4 max-h-[90vh] overflow-y-auto">
            <h3 className="text-xl font-bold mb-4">Edit Room Settings</h3>
            
            <div className="space-y-4">
              {/* Room Name */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Room Name *
                </label>
                <input
                  type="text"
                  value={editingRoom.name}
                  onChange={(e) => setEditingRoom(prev => ({ ...prev, name: e.target.value }))}
                  placeholder="Enter room name"
                  className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  maxLength={50}
                />
              </div>

              {/* Room Description */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Description
                </label>
                <textarea
                  value={editingRoom.description}
                  onChange={(e) => setEditingRoom(prev => ({ ...prev, description: e.target.value }))}
                  placeholder="Enter room description"
                  rows={3}
                  className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  maxLength={200}
                />
              </div>

              {/* Max Participants with XP-based limits */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Max Participants
                </label>
                <select
                  value={editingRoom.maxParticipants}
                  onChange={(e) => setEditingRoom(prev => ({ ...prev, maxParticipants: parseInt(e.target.value) }))}
                  className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  {Array.from({ length: maxRoomParticipants }, (_, i) => i + 1).map(num => (
                    <option key={num} value={num}>{num}</option>
                  ))}
                </select>
                <p className="text-xs text-gray-500 mt-1">
                  Your XP level allows up to {maxRoomParticipants} participants
                </p>
              </div>

              {/* Invitation Code */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Invitation Code
                </label>
                <div className="flex gap-2">
                  <input
                    type="text"
                    value={editingRoom.invitationCode || 'Generating...'}
                    readOnly
                    className="flex-1 border border-gray-300 rounded-lg px-3 py-2 bg-gray-50"
                  />
                  <button
                    onClick={async () => {
                      try {
                        const code = await generateInvitationCode(editingRoom.id || editingRoom._id);
                        setEditingRoom(prev => ({ ...prev, invitationCode: code }));
                      } catch (error) {
                        console.error('Failed to generate invitation code:', error);
                      }
                    }}
                    className="bg-blue-500 text-white px-4 py-2 rounded-lg hover:bg-blue-600 transition-colors"
                  >
                    Generate
                  </button>
                </div>
              </div>
            </div>

            {/* Error display */}
            {errorDisplay}

            {/* Action buttons */}
            <div className="flex gap-2 mt-6">
              <button
                onClick={() => {
                  setShowRoomSettings(false);
                  setEditingRoom(null);
                }}
                className="flex-1 bg-gray-500 hover:bg-gray-600 text-white py-2 px-4 rounded-lg transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={() => handleEditRoom(editingRoom)}
                disabled={!editingRoom.name.trim()}
                className="flex-1 bg-blue-500 hover:bg-blue-600 disabled:bg-gray-300 text-white py-2 px-4 rounded-lg transition-colors"
              >
                Save Changes
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Delete Room Confirmation Modal */}
      {showDeleteConfirm && roomToDelete && (
        <div className="fixed inset-0 bg-black bg-opacity-60 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 max-w-md w-full mx-4">
            <h3 className="text-xl font-bold mb-4 text-red-600">Delete Room</h3>
            <p className="text-gray-600 mb-6">
              Are you sure you want to delete "{roomToDelete.name}"? This action cannot be undone.
            </p>
            
            <div className="flex gap-2">
              <button
                onClick={() => {
                  setShowDeleteConfirm(false);
                  setRoomToDelete(null);
                }}
                className="flex-1 bg-gray-500 hover:bg-gray-600 text-white py-2 px-4 rounded-lg transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={async () => {
                  try {
                    await handleDeleteRoom(roomToDelete);
                    // Close modal and redirect after successful deletion
                    setShowDeleteConfirm(false);
                    setRoomToDelete(null);
                    // Call onDeleteSuccess callback if provided
                    if (onDeleteSuccess) {
                      onDeleteSuccess();
                    }
                  } catch (error) {
                    console.error('Failed to delete room:', error);
                  }
                }}
                className="flex-1 bg-red-500 hover:bg-red-600 text-white py-2 px-4 rounded-lg transition-colors"
              >
                Delete Room
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Friends Invitation Modal */}
      {showInviteFriendsModal && invitingToRoom && (
        <div className="fixed inset-0 bg-black bg-opacity-60 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 max-w-lg w-full mx-4 max-h-[90vh] overflow-y-auto">
            <h3 className="text-xl font-bold mb-4">Invite Friends to "{invitingToRoom.name}"</h3>
            
            {/* Room Info */}
            <div className="bg-blue-50 border border-blue-200 rounded-lg p-3 mb-4">
              <p className="text-sm text-blue-800">
                <strong>Current participants:</strong> {invitingToRoom.participantCount || 0}/{invitingToRoom.maxParticipants || 4}
              </p>
              {invitingToRoom.maxParticipants > 4 && (
                <p className="text-xs text-blue-600 mt-1">
                  XP required for {invitingToRoom.maxParticipants} participants: {(invitingToRoom.maxParticipants - 4) * 300} XP
                </p>
              )}
            </div>

            {/* Friends List */}
            {friendsList.length > 0 ? (
              <div className="space-y-2 mb-4">
                <h4 className="font-medium text-gray-700">Select friends to invite:</h4>
                {friendsList.map((friend) => (
                  <div key={friend._id} className="flex items-center justify-between p-3 border border-gray-200 rounded-lg">
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 rounded-full bg-blue-500 flex items-center justify-center">
                        <span className="text-white font-bold text-sm">
                          {friend.username?.[0]?.toUpperCase() || 'F'}
                        </span>
                      </div>
                      <div>
                        <p className="font-medium text-gray-900">{friend.username || 'Unknown User'}</p>
                        <p className="text-sm text-gray-500">ID: {friend.uniqueId || friend._id}</p>
                      </div>
                    </div>
                    <button
                      onClick={() => inviteFriendToRoom(friend.id || friend._id, invitingToRoom.id || invitingToRoom._id)}
                      className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg text-sm transition-colors"
                    >
                      Invite
                    </button>
                  </div>
                ))}
              </div>
            ) : (
              <div className="text-center py-8">
                <Icon icon="mdi:account-group-outline" className="w-16 h-16 mx-auto mb-4 text-gray-400" />
                <p className="text-gray-500 mb-2">No friends found</p>
                <p className="text-sm text-gray-400">Add friends first to invite them to your room</p>
              </div>
            )}

            {/* Action buttons */}
            <div className="flex gap-2 mt-6">
              <button
                onClick={() => {
                  setShowInviteFriendsModal(false);
                  setInvitingToRoom(null);
                }}
                className="flex-1 bg-gray-500 hover:bg-gray-600 text-white py-2 px-4 rounded-lg transition-colors"
              >
                Close
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Notifications Modal */}
      {showNotifications && (
        <div className="fixed inset-0 bg-black bg-opacity-60 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 max-w-lg w-full mx-4 max-h-[90vh] overflow-y-auto">
            <h3 className="text-xl font-bold mb-4">Notifications</h3>
            
            {notifications.length > 0 ? (
              <div className="space-y-3 mb-4">
                {notifications.map((notification) => (
                  <div key={notification.id} className="border border-gray-200 rounded-lg p-3">
                    <div className="flex justify-between items-start">
                      <div className="flex-1">
                        <h4 className="font-medium text-gray-900">{notification.title}</h4>
                        <p className="text-sm text-gray-600">{notification.message}</p>
                        <p className="text-xs text-gray-400 mt-1">
                          {new Date(notification.createdAt).toLocaleDateString()}
                        </p>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="text-center py-8">
                <Icon icon="mdi:bell-off" className="w-16 h-16 mx-auto mb-4 text-gray-400" />
                <p className="text-gray-500">No notifications</p>
              </div>
            )}

            {/* Action buttons */}
            <div className="flex gap-2 mt-6">
              <button
                onClick={() => setShowNotifications(false)}
                className="flex-1 bg-gray-500 hover:bg-gray-600 text-white py-2 px-4 rounded-lg transition-colors"
              >
                Close
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Room Invitations Modal */}
      {showRoomInvitations && (
        <div className="fixed inset-0 bg-black bg-opacity-60 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 max-w-lg w-full mx-4 max-h-[90vh] overflow-y-auto">
            <h3 className="text-xl font-bold mb-4">Room Invitations</h3>
            
            {roomInvitations.length > 0 ? (
              <div className="space-y-3 mb-4">
                {roomInvitations.map((invitation) => (
                  <div key={invitation.id} className="border border-gray-200 rounded-lg p-3">
                    <div className="flex justify-between items-start">
                      <div className="flex-1">
                        <h4 className="font-medium text-gray-900">{invitation.title}</h4>
                        <p className="text-sm text-gray-600">{invitation.message}</p>
                        <p className="text-xs text-gray-400 mt-1">
                          {new Date(invitation.createdAt).toLocaleDateString()}
                        </p>
                      </div>
                      <div className="flex gap-2 ml-3">
                        <button
                          onClick={() => handleAcceptRoomInvitation(invitation.targetId)}
                          className="bg-green-500 hover:bg-green-600 text-white px-3 py-1 rounded text-sm"
                        >
                          Accept
                        </button>
                        <button
                          onClick={() => {
                            // Delete notification (reject invitation)
                            fetch(`/api/v1/notifications/${invitation.id}`, {
                              method: 'DELETE',
                              headers: {
                                'Authorization': `Bearer ${localStorage.getItem('accessToken')}`,
                                'Content-Type': 'application/json'
                              }
                            }).then(() => {
                              // Refresh invitations
                              setShowRoomInvitations(false);
                            });
                          }}
                          className="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded text-sm"
                        >
                          Decline
                        </button>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="text-center py-8">
                <Icon icon="mdi:account-group-outline" className="w-16 h-16 mx-auto mb-4 text-gray-400" />
                <p className="text-gray-500">No room invitations</p>
              </div>
            )}

            {/* Action buttons */}
            <div className="flex gap-2 mt-6">
              <button
                onClick={() => setShowRoomInvitations(false)}
                className="flex-1 bg-gray-500 hover:bg-gray-600 text-white py-2 px-4 rounded-lg transition-colors"
              >
                Close
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Share Modal */}
      {showShareModal && sharingItem && (
        <div className="fixed inset-0 bg-black bg-opacity-60 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 max-w-lg w-full mx-4 max-h-[90vh] overflow-y-auto">
            <h3 className="text-xl font-bold mb-4">Share {sharingItem.type === 'todo' ? 'Task' : sharingItem.type === 'material' ? 'File' : 'Note'}</h3>
            
            <div className="mb-4">
              <p className="text-gray-600 mb-2">
                <strong>Item:</strong> {sharingItem.type === 'todo' ? sharingItem.title : sharingItem.type === 'material' ? sharingItem.name : sharingItem.content}
              </p>
              <p className="text-sm text-gray-500">Select the rooms where you want to share this item:</p>
            </div>

            {sharedRooms.length > 0 ? (
              <div className="space-y-2 mb-4">
                {sharedRooms.map((room) => (
                  <label key={room.id || room._id} className="flex items-center space-x-3 cursor-pointer">
                    <input
                      type="checkbox"
                      checked={selectedRoomsForSharing.includes(room.id || room._id)}
                      onChange={() => toggleRoomSelection(room.id || room._id)}
                      className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                    />
                    <div className="flex-1">
                      <span className="text-gray-900 font-medium">{room.name}</span>
                      <p className="text-sm text-gray-500">{room.description || 'No description'}</p>
                    </div>
                  </label>
                ))}
              </div>
            ) : (
              <div className="text-center py-4">
                <Icon icon="mdi:account-group-outline" className="w-16 h-16 mx-auto mb-4 text-gray-400" />
                <p className="text-gray-500">No shared rooms available</p>
                <p className="text-sm text-gray-400">Create a shared room first to share items</p>
              </div>
            )}

            {/* Action buttons */}
            <div className="flex gap-2 mt-6">
              <button
                onClick={() => setShowShareModal(false)}
                className="flex-1 bg-gray-500 hover:bg-gray-600 text-white py-2 px-4 rounded-lg transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={shareItemToRooms}
                disabled={selectedRoomsForSharing.length === 0}
                className="flex-1 bg-blue-500 hover:bg-blue-600 disabled:bg-gray-300 text-white py-2 px-4 rounded-lg transition-colors"
              >
                Share to {selectedRoomsForSharing.length} Room{selectedRoomsForSharing.length !== 1 ? 's' : ''}
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
};

RoomModals.propTypes = {
  // Todo Modal
  showTodoModal: PropTypes.bool.isRequired,
  setShowTodoModal: PropTypes.func.isRequired,
  newTodo: PropTypes.string.isRequired,
  setNewTodo: PropTypes.func.isRequired,
  addTodo: PropTypes.func.isRequired,
  
  // Material Modal
  showMaterialModal: PropTypes.bool.isRequired,
  setShowMaterialModal: PropTypes.func.isRequired,
  newMaterialType: PropTypes.string.isRequired,
  setNewMaterialType: PropTypes.func.isRequired,
  newMaterialUrl: PropTypes.string.isRequired,
  setNewMaterialUrl: PropTypes.func.isRequired,
  uploadedFile: PropTypes.object,
  setUploadedFile: PropTypes.func.isRequired,
  addMaterial: PropTypes.func.isRequired,
  
  // Session XP Modal
  showSessionXPModal: PropTypes.bool.isRequired,
  setShowSessionXPModal: PropTypes.func.isRequired,
  sessionXP: PropTypes.number.isRequired,
  
  // Video Call Modal
  showVideoCall: PropTypes.bool.isRequired,
  setShowVideoCall: PropTypes.func.isRequired,
  currentRoom: PropTypes.object,
  currentUser: PropTypes.object,
  VideoCallRoom: PropTypes.elementType.isRequired,
  
  // Join Room Modal
  showJoinModal: PropTypes.bool.isRequired,
  setShowJoinModal: PropTypes.func.isRequired,
  joinRoomCode: PropTypes.string.isRequired,
  setJoinRoomCode: PropTypes.func.isRequired,
  joinRoom: PropTypes.func.isRequired,
  
  // Create Room Modal
  showCreateRoomModal: PropTypes.bool.isRequired,
  setShowCreateRoomModal: PropTypes.func.isRequired,
  newRoomData: PropTypes.object.isRequired,
  setNewRoomData: PropTypes.func.isRequired,
  createRoomError: PropTypes.string,
  setCreateRoomError: PropTypes.func.isRequired,
  handleCreateRoom: PropTypes.func.isRequired,
  maxRoomParticipants: PropTypes.number.isRequired,
  userXP: PropTypes.number.isRequired,
  
  // Room Settings Modal
  showRoomSettings: PropTypes.bool.isRequired,
  setShowRoomSettings: PropTypes.func.isRequired,
  editingRoom: PropTypes.object,
  setEditingRoom: PropTypes.func.isRequired,
  handleEditRoom: PropTypes.func.isRequired,
  generateInvitationCode: PropTypes.func.isRequired,
  errorDisplay: PropTypes.node,
  
  // Delete Room Confirmation Modal
  showDeleteConfirm: PropTypes.bool.isRequired,
  setShowDeleteConfirm: PropTypes.func.isRequired,
  roomToDelete: PropTypes.object,
  setRoomToDelete: PropTypes.func.isRequired,
  handleDeleteRoom: PropTypes.func.isRequired,
  onDeleteSuccess: PropTypes.func,
  
  // Friends Invitation Modal
  showInviteFriendsModal: PropTypes.bool.isRequired,
  setShowInviteFriendsModal: PropTypes.func.isRequired,
  invitingToRoom: PropTypes.object,
  setInvitingToRoom: PropTypes.func.isRequired,
  friendsList: PropTypes.array.isRequired,
  inviteFriendToRoom: PropTypes.func.isRequired,
  userRooms: PropTypes.array.isRequired,
  // userXP is already defined above
  
  // Notifications Modal
  showNotifications: PropTypes.bool.isRequired,
  setShowNotifications: PropTypes.func.isRequired,
  notifications: PropTypes.array.isRequired,
  
  // Room Invitations Modal
  showRoomInvitations: PropTypes.bool.isRequired,
  setShowRoomInvitations: PropTypes.func.isRequired,
  roomInvitations: PropTypes.array.isRequired,
  handleAcceptRoomInvitation: PropTypes.func.isRequired,
  
  // Share Modal
  showShareModal: PropTypes.bool.isRequired,
  setShowShareModal: PropTypes.func.isRequired,
  sharingItem: PropTypes.object,
  setSharingItem: PropTypes.func.isRequired,
  selectedRoomsForSharing: PropTypes.array.isRequired,
  toggleRoomSelection: PropTypes.func.isRequired,
  shareItemToRooms: PropTypes.func.isRequired,
  sharedRooms: PropTypes.array.isRequired
};

export default RoomModals; 