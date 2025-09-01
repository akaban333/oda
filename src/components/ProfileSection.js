import React, { useState, useEffect } from 'react';
import { Icon } from '@iconify/react';
import EditProfileModal from './EditProfileModal';
import PostsPanel from './PostsPanel';
import { friendsAPI, authAPI, sessionsAPI } from '../utils/api';

const ProfileSection = ({ user, refreshUserProfile }) => {
  const [isEditProfileModalOpen, setIsEditProfileModalOpen] = useState(false);
  const [showFriendsModal, setShowFriendsModal] = useState(false);
  const [showPostsPanel, setShowPostsPanel] = useState(false);
  const [friendsList, setFriendsList] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [userStats, setUserStats] = useState({
    friendsCount: 0,
    totalXP: 0,
    sessionsCount: 0
  });

  // Fetch user statistics
  useEffect(() => {
    const fetchUserStats = async () => {
      try {
        const [friendsData, sessionStats] = await Promise.all([
          friendsAPI.getFriends().catch(() => ({ friends: [] })),
          sessionsAPI.getStats().catch(() => ({ totalSessions: 0, totalXP: 0 }))
        ]);

        // Count all friend relationships (pending, requested, accepted)
        const totalFriendsCount = friendsData.friends?.length || 0;
        
        setUserStats({
          friendsCount: totalFriendsCount,
          totalXP: sessionStats.totalXP || user?.xp || 0,
          sessionsCount: sessionStats.totalSessions || 0
        });
      } catch (error) {
        console.error('Error fetching user stats:', error);
        // Use fallback values from user object
        setUserStats({
          friendsCount: user?.friendsCount || 0,
          totalXP: user?.xp || 0,
          sessionsCount: user?.sessionsCount || 0
        });
      }
    };

    if (user) {
      fetchUserStats();
    }
  }, [user]);

  const handleShowFriends = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const data = await friendsAPI.getFriends();
      
      setFriendsList(data.friends || []);
      setShowFriendsModal(true);
    } catch (error) {
      console.error('Error fetching friends:', error);
      setError('Failed to load friends. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = async () => {
    try {
      await authAPI.logout();
      localStorage.clear();
      window.location.reload();
    } catch (error) {
      console.error('Error logging out:', error);
      // Still clear localStorage and reload even if logout API fails
      localStorage.clear();
      window.location.reload();
    }
  };

  const handleAcceptFriend = async (userId) => {
    try {
      await friendsAPI.acceptFriendRequest(userId);
      // Refresh friends list and user profile
      const data = await friendsAPI.getFriends();
      setFriendsList(data.friends || []);
      refreshUserProfile();
    } catch (error) {
      console.error('Error accepting friend request:', error);
      setError('Failed to accept friend request.');
    }
  };

  const handleRejectFriend = async (userId) => {
    try {
      await friendsAPI.rejectFriendRequest(userId);
      // Refresh friends list and user profile
      const data = await friendsAPI.getFriends();
      setFriendsList(data.friends || []);
      refreshUserProfile();
    } catch (error) {
      console.error('Error rejecting friend request:', error);
      setError('Failed to reject friend request.');
    }
  };

  const handleRemoveFriend = async (userId) => {
    try {
      await friendsAPI.removeFriend(userId);
      // Refresh friends list and user profile
      const data = await friendsAPI.getFriends();
      setFriendsList(data.friends || []);
      refreshUserProfile();
    } catch (error) {
      console.error('Error removing friend:', error);
      setError('Failed to remove friend.');
    }
  };

  return (
    <div className="min-h-screen bg-[#282435] -mt-20 pt-20">
      <div className="max-w-4xl mx-auto p-8">
        {/* Main Profile Information Section */}
        <div className="bg-[#282435] rounded-2xl p-8 mb-8">
          <div className="flex items-center justify-between">
            {/* Left side: Avatar and User Info */}
            <div className="flex items-center gap-8">
              {/* Avatar */}
              <div className="w-24 h-24 bg-gray-300 flex items-center justify-center overflow-hidden rounded-full">
                {user.avatarUrl ? (
                  <img 
                    src={user.avatarUrl} 
                    alt="avatar" 
                    className="w-full h-full object-cover rounded-full" 
                  />
                ) : (
                  <div className="text-4xl font-bold text-gray-600">
                    {user.username?.[0]?.toUpperCase() || 'U'}
                  </div>
                )}
              </div>
              
              {/* User Info */}
              <div className="flex items-center gap-8">
                <div>
                  <div className="text-white text-2xl font-bold mb-1">
                    {user.username || 'USERNAME'}
                  </div>
                  <div className="text-gray-400 text-lg flex items-center gap-2">
                    <span>@{user.uniqueId || user.username || 'USERID'}</span>
                    {user.uniqueId && (
                      <button
                        onClick={() => {
                          navigator.clipboard.writeText(user.uniqueId);
                          alert('User ID copied to clipboard!');
                        }}
                        className="text-[#B8A17D] hover:text-[#A08F6B] transition-colors"
                        title="Copy User ID"
                      >
                        <Icon icon="mdi:content-copy" className="w-4 h-4" />
                      </button>
                    )}
                  </div>
                </div>
                
                {/* Vertical separator */}
                <div className="w-px h-12 bg-gray-600"></div>
                
                <div>
                  <div className="text-white text-lg font-bold mb-1">FIELD</div>
                  <div className="text-gray-400 text-lg">SOFTWARE ENGINEERING</div>
                </div>
              </div>
            </div>
            
            {/* Right side: Edit Profile Button */}
            <div className="flex gap-4">
              <button 
                onClick={() => setIsEditProfileModalOpen(true)}
                className="bg-[#B8A17D] text-white px-6 py-3 rounded-lg font-bold text-lg hover:bg-[#A08F6B] transition-colors"
              >
                EDIT PROFILE
              </button>
              <button 
                onClick={() => setShowPostsPanel(true)}
                className="bg-[#7681B3] text-white px-6 py-3 rounded-lg font-bold text-lg hover:bg-[#8b96d1] transition-colors"
              >
                VIEW POSTS
              </button>
            </div>
          </div>
        </div>

        {/* Error Display */}
        {error && (
          <div className="bg-red-600 text-white p-4 rounded-lg mb-6">
            <div className="flex items-center justify-between">
              <span>{error}</span>
              <button 
                onClick={() => setError(null)}
                className="text-white hover:text-red-200"
              >
                ×
              </button>
            </div>
          </div>
        )}

        {/* Statistics Section */}
        <div className="flex gap-6 mb-8">
          {/* Friends */}
          <div 
            className="flex-1 bg-[#282435] border border-gray-600 rounded-xl p-6 cursor-pointer hover:bg-[#2F2A3F] transition-colors flex flex-col items-center justify-center"
            onClick={handleShowFriends}
          >
            <div className="flex items-center gap-3 mb-2">
              <Icon icon="fa-solid:user-friends" className="w-6 h-6 text-[#D8C2D8]" />
            </div>
            <div className="text-white text-center font-bold text-lg">
              {loading ? 'Loading...' : `${userStats.friendsCount} FRIENDS`}
            </div>
          </div>
          
          {/* XP */}
          <div className="flex-1 bg-[#282435] border border-gray-600 rounded-xl p-6 flex flex-col items-center justify-center">
            <div className="flex items-center gap-3 mb-2">
              <Icon icon="game-icons:fluffy-flame" className="w-6 h-6 text-[#D8C2D8]" />
            </div>
            <div className="text-white text-center font-bold text-lg">
              {userStats.totalXP} TOTAL XP
            </div>
          </div>
          
          {/* Sessions */}
          <div className="flex-1 bg-[#282435] border border-gray-600 rounded-xl p-6 flex flex-col items-center justify-center">
            <div className="flex items-center gap-3 mb-2">
              <Icon icon="material-symbols:book-outline-rounded" className="w-6 h-6 text-[#D8C2D8]" />
            </div>
            <div className="text-white text-center font-bold text-lg">
              {userStats.sessionsCount} SESSIONS
            </div>
          </div>
        </div>

        {/* Bio Section */}
        <div className="bg-[#282435] border border-gray-600 rounded-xl p-6 mb-8">
          <div className="flex items-start gap-4">
            <div className="w-px h-16 bg-gray-600 mt-2"></div>
            <div className="flex-1">
              <div className="text-[#B8A17D] font-bold text-lg mb-2">BIO:</div>
              <div className="text-white text-lg leading-relaxed">
                {user.bio || 'HI! I AM ERSA, I AM 20 YEARS OLD AND I AM IN MY THIRD YEAR OF BACHELOR STUDIES FOR SOFTWARE ENGINEERING.'}
              </div>
            </div>
          </div>
        </div>

        {/* Action Buttons */}
        <div className="flex justify-between items-center">
          <button 
            onClick={handleLogout}
            className="bg-[#B8A17D] text-white px-8 py-3 rounded-lg font-bold text-lg hover:bg-[#A08F6B] transition-colors"
          >
            LOG OUT
          </button>
        </div>
      </div>

      {/* Edit Profile Modal */}
      {isEditProfileModalOpen && (
        <EditProfileModal
          user={user}
          onClose={() => setIsEditProfileModalOpen(false)}
          onProfileUpdated={refreshUserProfile}
        />
      )}
      
      {/* Friends Modal */}
      {showFriendsModal && (
        <div className="fixed inset-0 bg-black bg-opacity-60 flex items-center justify-center z-50">
          <div className="bg-[#282435] rounded-lg p-6 max-w-2xl w-full border border-gray-600 max-h-[80vh] overflow-y-auto">
            <h3 className="text-xl text-white mb-4">Friends & Requests</h3>
            
            {/* Friend Requests Section */}
            {friendsList.filter(friend => friend.status === 'requested').length > 0 && (
              <div className="mb-6">
                <h4 className="text-lg text-[#B8A17D] font-semibold mb-3">Friend Requests</h4>
                <div className="space-y-3">
                  {friendsList.filter(friend => friend.status === 'requested').map(friend => (
                    <div key={friend.id} className="bg-[#1F1D2B] rounded-lg p-4 border border-gray-600">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-3">
                          <div className="w-10 h-10 bg-gray-500 flex items-center justify-center text-lg font-bold overflow-hidden rounded-full">
                            {friend.avatarUrl ? (
                              <img 
                                src={friend.avatarUrl} 
                                alt="avatar" 
                                className="w-full h-full object-cover rounded-full" 
                              />
                            ) : (
                              (friend.displayName?.[0]?.toUpperCase() || friend.username?.[0]?.toUpperCase() || '?')
                            )}
                          </div>
                          <div>
                            <div className="text-white font-semibold">{friend.displayName || friend.username || 'Unknown User'}</div>
                            <div className="text-[#7681B3] text-sm">@{friend.uniqueId}</div>
                          </div>
                        </div>
                        <div className="flex gap-2">
                          <button 
                            onClick={() => handleAcceptFriend(friend.userId)}
                            disabled={loading}
                            className="bg-green-600 text-white px-4 py-2 rounded-lg text-sm font-semibold hover:bg-green-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                          >
                            {loading ? 'Accepting...' : 'Accept'}
                          </button>
                          <button 
                            onClick={() => handleRejectFriend(friend.userId)}
                            disabled={loading}
                            className="bg-red-600 text-white px-4 py-2 rounded-lg text-sm font-semibold hover:bg-red-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                          >
                            {loading ? 'Rejecting...' : 'Reject'}
                          </button>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Pending Requests Section */}
            {friendsList.filter(friend => friend.status === 'pending').length > 0 && (
              <div className="mb-6">
                <h4 className="text-lg text-[#B8A17D] font-semibold mb-3">Pending Requests</h4>
                <div className="space-y-3">
                  {friendsList.filter(friend => friend.status === 'pending').map(friend => (
                    <div key={friend.id} className="bg-[#1F1D2B] rounded-lg p-4 border border-gray-600">
                      <div className="flex items-center gap-3">
                        <div className="w-10 h-10 bg-gray-500 flex items-center justify-center text-lg font-bold overflow-hidden rounded-full">
                          {friend.avatarUrl ? (
                            <img 
                              src={friend.avatarUrl} 
                              alt="avatar" 
                              className="w-full h-full object-cover rounded-full" 
                            />
                          ) : (
                            (friend.displayName?.[0]?.toUpperCase() || friend.username?.[0]?.toUpperCase() || '?')
                          )}
                        </div>
                        <div>
                          <div className="text-white font-semibold">{friend.displayName || friend.username || 'Unknown User'}</div>
                          <div className="text-[#7681B3] text-sm">@{friend.uniqueId}</div>
                          <div className="text-yellow-400 text-xs">Request sent - waiting for response</div>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Accepted Friends Section */}
            {friendsList.filter(friend => friend.status === 'accepted').length > 0 && (
              <div className="mb-6">
                <h4 className="text-lg text-[#B8A17D] font-semibold mb-3">Friends</h4>
                <div className="space-y-3">
                  {friendsList.filter(friend => friend.status === 'accepted').map(friend => (
                    <div key={friend.id} className="bg-[#1F1D2B] rounded-lg p-4 border border-gray-600">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-3">
                          <div className="w-10 h-10 bg-gray-500 flex items-center justify-center text-lg font-bold overflow-hidden rounded-full">
                            {friend.avatarUrl ? (
                              <img 
                                src={friend.avatarUrl} 
                                alt="avatar" 
                                className="w-full h-full object-cover rounded-full" 
                              />
                            ) : (
                              (friend.displayName?.[0]?.toUpperCase() || friend.username?.[0]?.toUpperCase() || '?')
                            )}
                          </div>
                          <div>
                            <div className="text-white font-semibold">{friend.displayName || friend.username || 'Unknown User'}</div>
                            <div className="text-[#7681B3] text-sm">@{friend.uniqueId}</div>
                            <div className="text-green-400 text-xs">Friends since {new Date(friend.since).toLocaleDateString()}</div>
                          </div>
                        </div>
                        <button 
                          onClick={() => handleRemoveFriend(friend.userId)}
                          disabled={loading}
                          className="bg-gray-600 text-white px-3 py-1 rounded text-sm hover:bg-gray-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                        >
                          {loading ? 'Removing...' : 'Remove'}
                        </button>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* No Friends Message */}
            {friendsList.length === 0 && (
              <div className="text-center py-8">
                <div className="text-gray-400 text-lg">No friends or requests yet.</div>
                <div className="text-gray-500 text-sm mt-2">Start connecting with other students!</div>
              </div>
            )}

            <button onClick={() => setShowFriendsModal(false)} className="mt-6 bg-[#B8A17D] text-white px-6 py-2 rounded-lg font-bold hover:bg-[#A08F6B] transition-colors">Close</button>
          </div>
        </div>
      )}
      
      {/* PostsPanel */}
      <PostsPanel 
        isOpen={showPostsPanel} 
        onClose={() => setShowPostsPanel(false)} 
        refreshTrigger={user?.id}
        key={`posts-${user?.id}`}
      />
      
      {!user && <div className="text-red-500">User not loaded!</div>}
    </div>
  );
};

export default ProfileSection;
