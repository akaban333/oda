import React, { useState, useEffect } from 'react';
import { Icon } from '@iconify/react';
import { friendsAPI, usersAPI, postsAPI, notificationsAPI } from '../utils/api';

const SocialFeed = () => {
  const [expandedPosts, setExpandedPosts] = useState(new Set());
  const [likedPosts, setLikedPosts] = useState(new Set());
  const [likedComments, setLikedComments] = useState(new Set());
  const [commentingOn, setCommentingOn] = useState(null);
  const [newComment, setNewComment] = useState('');
  const [searchFriends, setSearchFriends] = useState('');
  const [searchResults, setSearchResults] = useState([]);
  const [searchError, setSearchError] = useState(null);
  const [addFriendStatus, setAddFriendStatus] = useState({});
  const [posts, setPosts] = useState([]);
  const [notifications, setNotifications] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showPostForm, setShowPostForm] = useState(false);
  const [newPostContent, setNewPostContent] = useState('');
  const [mousePosition, setMousePosition] = useState({ x: 0, y: 0 });
  const [showFriendRequestModal, setShowFriendRequestModal] = useState(false);
  const [selectedFriendRequest, setSelectedFriendRequest] = useState(null);

  // Track mouse position for cursor effect
  useEffect(() => {
    const updateMousePosition = (e) => {
      setMousePosition({ x: e.clientX, y: e.clientY });
    };

    window.addEventListener('mousemove', updateMousePosition);

    return () => {
      window.removeEventListener('mousemove', updateMousePosition);
    };
  }, []);

  // Fetch posts and notifications
  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        setError(null);
        
        console.log('Fetching posts and notifications...');
        const [postsData, notificationsData] = await Promise.all([
          postsAPI.getPosts().catch(err => {
            console.error('Error fetching posts:', err);
            return { posts: [] }; // Return empty posts on error
          }),
          notificationsAPI.getNotifications().catch(err => {
            console.error('Error fetching notifications:', err);
            return { notifications: [] }; // Return empty notifications on error
          })
        ]);
        
        console.log('Posts data:', postsData);
        console.log('Notifications data:', notificationsData);
        
        // Debug: Log individual posts to see their structure
        if (postsData.posts && postsData.posts.length > 0) {
          console.log('First post structure:', postsData.posts[0]);
          console.log('First post authorName:', postsData.posts[0].authorName);
          console.log('First post authorId:', postsData.posts[0].authorId);
        }
        
        setPosts(postsData.posts || []);
        setNotifications(notificationsData.notifications || []);
        
        // Initialize liked posts and comments state
        if (postsData.posts) {
          const likedPostsSet = new Set();
          const likedCommentsSet = new Set();
          
          postsData.posts.forEach(post => {
            if (post.isLiked) {
              likedPostsSet.add(post._id || post.id);
            }
            if (post.comments) {
              post.comments.forEach(comment => {
                if (comment.isLiked) {
                  likedCommentsSet.add(comment._id || comment.id);
                }
              });
            }
          });
          
          setLikedPosts(likedPostsSet);
          setLikedComments(likedCommentsSet);
        }
        
      } catch (error) {
        console.error('Error fetching data:', error);
        console.error('Error details:', {
          message: error.message,
          stack: error.stack
        });
        // Don't set error, just show empty state
        setPosts([]);
        setNotifications([]);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  // Friend search API integration
  useEffect(() => {
    if (searchFriends.trim()) {
      const searchUsers = async () => {
        try {
          const data = await usersAPI.searchUsers(searchFriends);
          setSearchResults(data.results || []);
        } catch (error) {
          console.error('Error searching users:', error);
          setSearchResults([]);
        }
      };
      
      searchUsers();
    } else {
      setSearchResults([]);
    }
  }, [searchFriends]);

  const fetchPosts = async () => {
    try {
      const postsData = await postsAPI.getPosts();
      setPosts(postsData.posts || []);
      
      // Update liked posts and comments state
      if (postsData.posts) {
        const likedPostsSet = new Set();
        const likedCommentsSet = new Set();
        
        postsData.posts.forEach(post => {
          if (post.isLiked) {
            likedPostsSet.add(post._id || post.id);
          }
          if (post.comments) {
            post.comments.forEach(comment => {
              if (comment.isLiked) {
                likedCommentsSet.add(comment._id || comment.id);
              }
            });
          }
        });
        
        setLikedPosts(likedPostsSet);
        setLikedComments(likedCommentsSet);
      }
    } catch (error) {
      console.error('Error fetching posts:', error);
    }
  };

  const handleAddFriend = async (uniqueId) => {
    try {
      const data = await friendsAPI.sendFriendRequest({ uniqueId });
      setAddFriendStatus(s => ({ ...s, [uniqueId]: data.message || 'Friend request sent!' }));
    } catch (error) {
      console.error('Error sending friend request:', error);
      setAddFriendStatus(s => ({ ...s, [uniqueId]: error.message || 'Failed to send friend request' }));
    }
  };

  const togglePostExpansion = (postId) => {
    const newExpandedPosts = new Set(expandedPosts);
    if (newExpandedPosts.has(postId)) {
      newExpandedPosts.delete(postId);
    } else {
      newExpandedPosts.add(postId);
    }
    setExpandedPosts(newExpandedPosts);
  };

  const toggleLike = async (postId) => {
    try {
      await postsAPI.likePost(postId);
      const newLikedPosts = new Set(likedPosts);
      if (newLikedPosts.has(postId)) {
        newLikedPosts.delete(postId);
      } else {
        newLikedPosts.add(postId);
      }
      setLikedPosts(newLikedPosts);
    } catch (error) {
      console.error('Error toggling like:', error);
      setError('Failed to update like. Please try again.');
    }
  };

  const toggleCommentLike = async (commentId) => {
    try {
      console.log('Toggling comment like for comment ID:', commentId);
      const response = await postsAPI.likeComment(commentId);
      console.log('Comment like response:', response);
      
      // Update local state for comment likes
      const newLikedComments = new Set(likedComments);
      if (newLikedComments.has(commentId)) {
        newLikedComments.delete(commentId);
      } else {
        newLikedComments.add(commentId);
      }
      setLikedComments(newLikedComments);
      
      // Refresh posts to get updated like counts
      fetchPosts();
    } catch (error) {
      console.error('Error toggling comment like:', error);
      setError('Failed to update comment like. Please try again.');
    }
  };

  // Handle notification click
  const handleNotificationClick = async (notification) => {
    try {
      // Delete notification after viewing
      await notificationsAPI.deleteNotification(notification.id);
      
      // Remove from local state
      setNotifications(prev => prev.filter(n => n.id !== notification.id));

      // Handle different notification types
      switch (notification.type) {
        case 'friend_request':
          // Show friend request modal
          setSelectedFriendRequest(notification);
          setShowFriendRequestModal(true);
          break;
        case 'friend_accepted':
          // Show friend accepted message
          console.log('Friend request accepted by:', notification.data?.accepterUsername);
          break;
        case 'post_like':
          // Navigate to post or show post details
          console.log('Post liked:', notification.targetId);
          break;
        case 'comment_like':
          // Navigate to comment or show comment details
          console.log('Comment liked:', notification.targetId);
          break;
        case 'post_comment':
          // Navigate to post or show comment
          console.log('Comment on post:', notification.targetId);
          break;
        case 'room_invitation':
          // Show room invitation actions
          console.log('Room invitation to:', notification.data?.roomName);
          break;
        default:
          console.log('Notification clicked:', notification);
      }
    } catch (error) {
      console.error('Error handling notification click:', error);
    }
  };

  const handleAcceptFriendRequest = async (notification) => {
    try {
      const requesterId = notification.data?.requesterId || notification.targetId;
      await friendsAPI.acceptFriendRequest(requesterId);
      
      // AGGRESSIVE CLEANUP: Remove ALL notifications related to this user
      setNotifications(prev => prev.filter(n => 
        !(n.type === 'friend_request' && 
          (n.data?.requesterId === requesterId || 
           n.targetId === requesterId ||
           n.data?.requesterID === requesterId))
      ));
      
      setShowFriendRequestModal(false);
      setSelectedFriendRequest(null);
      
      // Force refresh notifications from backend
      try {
        const notificationsData = await notificationsAPI.getNotifications();
        setNotifications(notificationsData.notifications || []);
      } catch (refreshError) {
        console.log('Could not refresh notifications, using local cleanup');
      }
      
      // Show success message
      setError('Friend request accepted!');
      setTimeout(() => setError(''), 3000);
    } catch (error) {
      console.error('Error accepting friend request:', error);
      setError('Failed to accept friend request. Please try again.');
    }
  };

  const handleRejectFriendRequest = async (notification) => {
    try {
      const requesterId = notification.data?.requesterId || notification.targetId;
      await friendsAPI.rejectFriendRequest(requesterId);
      
      // AGGRESSIVE CLEANUP: Remove ALL notifications related to this user
      setNotifications(prev => prev.filter(n => 
        !(n.type === 'friend_request' && 
          (n.data?.requesterId === requesterId || 
           n.targetId === requesterId ||
           n.data?.requesterID === requesterId))
      ));
      
      setShowFriendRequestModal(false);
      setSelectedFriendRequest(null);
      
      // Force refresh notifications from backend
      try {
        const notificationsData = await notificationsAPI.getNotifications();
        setNotifications(notificationsData.notifications || []);
      } catch (refreshError) {
        console.log('Could not refresh notifications, using local cleanup');
      }
      
      // Show success message
      setError('Friend request rejected.');
      setTimeout(() => setError(''), 3000);
    } catch (error) {
      console.error('Error rejecting friend request:', error);
      setError('Failed to reject friend request. Please try again.');
    }
  };

  const handleCommentSubmit = async (postId) => {
    if (newComment.trim()) {
      try {
        const newCommentData = await postsAPI.createComment(postId, {
          content: newComment
        });
        
        console.log('Comment created successfully:', newCommentData);
        
        // Update the posts state to include the new comment
        setPosts(prev => prev.map(post => 
          (post._id || post.id) === postId 
            ? { ...post, comments: [...(post.comments || []), newCommentData.comment] }
            : post
        ));
        
        setNewComment('');
        setCommentingOn(null);
      } catch (error) {
        console.error('Error submitting comment:', error);
        setError('Failed to submit comment. Please try again.');
      }
    }
  };

  const handlePostSubmit = async () => {
    if (newPostContent.trim()) {
      try {
        console.log('Attempting to create post with data:', {
          content: newPostContent
        });
        
        const newPostData = await postsAPI.createPost({
          content: newPostContent
        });
        
        console.log('Post created successfully:', newPostData);
        console.log('New post structure:', newPostData.post);
        console.log('New post authorName:', newPostData.post?.authorName);
        console.log('New post authorId:', newPostData.post?.authorId);
        
        setPosts(prev => [newPostData.post, ...prev]);
        setNewPostContent('');
        setShowPostForm(false);
        setError(null); // Clear any previous errors
      } catch (error) {
        console.error('Error submitting post:', error);
        console.error('Error details:', {
          message: error.message,
          stack: error.stack,
          response: error.response
        });
        setError('Failed to submit post. Please try again.');
      }
    }
  };

  const handleDeletePost = async (postId) => {
    try {
      await postsAPI.deletePost(postId);
      setPosts(prev => prev.filter(post => (post._id || post.id) !== postId));
    } catch (error) {
      console.error('Error deleting post:', error);
      setError('Failed to delete post. Please try again.');
    }
  };

  const handleFriendSearch = (e) => {
    setSearchFriends(e.target.value);
  };

  const handleSearchKeyPress = (e) => {
    if (e.key === 'Enter' && searchFriends.trim()) {
      handleSearchUsers();
    }
  };

  const handleSearchUsers = async () => {
    if (!searchFriends.trim()) return;
    
    // Validate that the search query looks like a UUID
    const uuidRegex = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;
    if (!uuidRegex.test(searchFriends.trim())) {
      setSearchError('Please enter a valid user ID (UUID format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)');
      setSearchResults([]);
      return;
    }
    
    try {
      setSearchError(null);
      setSearchResults([]);
      
      console.log('Searching for user ID:', searchFriends.trim());
      const data = await usersAPI.searchUsers(searchFriends.trim());
      console.log('Search response:', data);
      
      if (data && data.results && Array.isArray(data.results)) {
        if (data.results.length > 0) {
          setSearchResults(data.results);
          console.log('Search results:', data.results);
        } else {
          setSearchResults([]);
          setSearchError('No users found with that ID');
        }
      } else {
        console.error('Unexpected response format:', data);
        setSearchError('Unexpected response format from server');
        setSearchResults([]);
      }
    } catch (error) {
      console.error('Error searching users:', error);
      if (error.response) {
        console.error('Error response:', error.response);
        if (error.response.status === 404) {
          setSearchError('No users found with that ID');
        } else if (error.response.status === 400) {
          setSearchError('Invalid search query format');
        } else {
          setSearchError(`Search failed: ${error.response.status}`);
        }
      } else {
        setSearchError('Failed to search users. Please try again.');
      }
      setSearchResults([]);
    }
  };

  return (
    <div className="w-full h-full bg-black" style={{ fontFamily: 'Montserrat', minHeight: '100vh' }}>
      {/* Error/Success Message Display */}
      {error && (
        <div className={`fixed top-4 left-1/2 transform -translate-x-1/2 z-50 px-6 py-3 rounded-lg text-white font-semibold ${
          error.includes('accepted') || error.includes('rejected') 
            ? 'bg-green-600' 
            : 'bg-red-600'
        }`}>
          {error}
        </div>
      )}

      {/* Main Content Area */}
      <div className="flex h-full">
        {/* Posts Section */}
        <div className="flex-1 p-6 pt-10 overflow-y-auto relative">
          <h2 className="text-white text-2xl mb-6 font-['Silkscreen']">RECENT POSTS</h2>

          <div className="pb-20"> {/* Added padding at the bottom to prevent overlap with the add post button */}
            {loading ? (
              <div className="flex items-center justify-center mt-8 mb-12">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-white mr-3"></div>
                <p className="text-white text-base">Loading posts...</p>
              </div>
            ) : error ? (
              <div className="flex items-center justify-center mt-8 mb-12">
                <Icon icon="twemoji:cat" className="w-10 h-10 mr-3" />
                <p className="text-white text-base">{error}</p>
              </div>
            ) : posts.length === 0 ? (
              <div className="flex flex-col items-center justify-center mt-8 mb-12">
                <Icon icon="twemoji:cat" className="w-10 h-10 mb-4" />
                <p className="text-white text-base mb-4">Oops no posts to show yet! Want to be the first?</p>
                <button
                  onClick={() => setShowPostForm(true)}
                  className="bg-[#7681B3] text-white px-6 py-2 rounded-lg hover:bg-[#8b96d1] transition-colors"
                >
                  Add Post
                </button>
              </div>
            ) : (
              posts.map(post => (
                <div 
                  key={post._id || post.id} 
                  className="bg-[#1B1B1D] rounded-xl mb-5 overflow-hidden border border-[#7681B3]"
                >
                  <div className="p-5">
                    <div className="flex justify-between items-start">
                      <div className="flex gap-3">
                        <div className="w-10 h-10 rounded-full bg-gray-500"></div>
                        <div>
                          <p className="text-white text-base mb-1">
                            {post.authorName || 'Unknown User'} <span className="text-[#7681B3]">@{post.authorUniqueId || post.authorId}</span>
                          </p>
                          <p className="text-white mt-2 mb-2">
                            {post.content}
                          </p>
                        </div>
                      </div>
                      <div className="flex gap-2">
                        <button 
                          onClick={() => handleDeletePost(post._id || post.id)}
                          className="text-red-400 hover:text-red-300 transition-colors"
                          title="Delete post"
                        >
                          <Icon icon="mdi:delete" className="w-5 h-5" />
                        </button>
                        <button 
                          onClick={() => togglePostExpansion(post._id || post.id)}
                          className="text-white hover:text-gray-300 transition-colors"
                        >
                          <Icon 
                            icon={expandedPosts.has(post._id || post.id) ? "material-symbols:keyboard-arrow-up" : "material-symbols:keyboard-arrow-down"} 
                            className="w-6 h-6"
                          />
                        </button>
                      </div>
                    </div>
                    
                    <div className="flex items-center justify-between mt-4">
                      <div className="flex gap-5">
                        <button 
                          onClick={() => toggleLike(post._id || post.id)}
                          className="text-[#7681B3] hover:text-[#8b96d1] transition-colors"
                        >
                          <Icon 
                            icon={likedPosts.has(post._id || post.id) ? "mdi:heart" : "mdi:heart-outline"} 
                            className="w-5 h-5"
                          />
                        </button>
                        <button 
                          className="text-[#7681B3] hover:text-[#8b96d1] transition-colors"
                          onClick={() => setCommentingOn(post._id || post.id)}
                        >
                          Comment
                        </button>
                      </div>
                      <span className="text-gray-400 text-xs">{new Date(post.createdAt).toLocaleDateString()}</span>
                    </div>
                  </div>

                  {/* Comments Section */}
                  {expandedPosts.has(post._id || post.id) && (
                    <div className="border-t border-[#7681B3]">
                      {post.comments && post.comments.map((comment, index) => (
                        <div key={comment._id || comment.id}>
                          <div className="p-4 flex gap-3">
                            <div className="w-8 h-8 rounded-full bg-gray-500"></div>
                            <div className="flex-1 ml-2">
                              <p className="text-white mb-1">
                                {comment.authorName || 'Unknown User'} <span className="text-[#7681B3]">@{comment.authorId}</span>
                              </p>
                              <p className="text-white text-sm">
                                {comment.authorName || 'Unknown User'} <span className="text-[#7681B3]">@{comment.authorUniqueId || comment.authorId}</span>
                              </p>
                              <p className="text-white text-sm">
                                {comment.content}
                              </p>
                            </div>
                            <button 
                              onClick={() => toggleCommentLike(comment._id || comment.id)}
                              className="text-[#7681B3] hover:text-[#8b96d1] transition-colors ml-2"
                            >
                              <Icon 
                                icon={likedComments.has(comment._id || comment.id) ? "mdi:heart" : "mdi:heart-outline"} 
                                className="w-5 h-5"
                              />
                            </button>
                          </div>
                          {/* Thin separator between comments */}
                          {index < post.comments.length - 1 && (
                            <div className="mx-16 border-t border-[#4a4a5e]"></div>
                          )}
                        </div>
                      ))}

                      {/* New Comment Form */}
                      {commentingOn === (post._id || post.id) && (
                        <div className="p-4 border-t border-[#7681B3]">
                          <div className="flex gap-3">
                            <div className="w-8 h-8 rounded-full bg-gray-500"></div>
                            <div className="flex-1 ml-2">
                              <textarea
                                value={newComment}
                                onChange={(e) => setNewComment(e.target.value)}
                                placeholder="Write a comment..."
                                className="w-full bg-[#2a2a42] text-white p-2 rounded-lg border border-[#7681B3] focus:outline-none"
                                rows={2}
                              />
                              <div className="flex justify-end mt-2">
                                <button
                                  onClick={() => setCommentingOn(null)}
                                  className="px-3 py-1 mr-2 bg-transparent border border-[#7681B3] text-[#7681B3] rounded-lg"
                                >
                                  Cancel
                                </button>
                                <button
                                  onClick={() => handleCommentSubmit(post._id || post.id)}
                                  className="px-3 py-1 bg-[#7681B3] text-white rounded-lg"
                                >
                                  Post Comment
                                </button>
                              </div>
                            </div>
                          </div>
                        </div>
                      )}
                    </div>
                  )}
                </div>
              ))
            )}
          </div>

          {/* Add Post Button - Fixed at the bottom of the posts section */}
          <div className="fixed bottom-6 left-0 right-0 flex justify-center z-10" style={{ width: 'calc(100% - 384px)' }}>
            <button 
              className="w-3/4 h-9 border-2 border-[#7681B3] rounded-full flex items-center justify-center text-[#7681B3] hover:bg-[#7681B3] hover:text-white transition-colors bg-black"
              onClick={() => setShowPostForm(true)}
            >
              <Icon icon="mdi:plus" className="w-5 h-5" />
            </button>
          </div>
        </div>

        {/* Right Panel */}
        <div className="w-96 bg-[#C6AE82] p-5 pt-11 flex flex-col min-h-screen sticky top-0 right-0">
          {/* User Search Container - Black rounded rectangle matching notifications style */}
          <div className="bg-black rounded-2xl flex-1 flex flex-col overflow-hidden mb-4">
            {/* Search Header */}
            <div className="p-3 border-b border-[#333] flex justify-between items-center">
              <span className="text-white ml-2">Search Users</span>
              <Icon icon="material-symbols:search" className="text-white w-5 h-5 mr-1" />
            </div>

            {/* Search Input and Results */}
            <div className="p-3 space-y-3 flex-1">
              <div className="flex flex-col gap-2">
                <input
                  type="text"
                  value={searchFriends}
                  onChange={handleFriendSearch}
                  onKeyPress={handleSearchKeyPress}
                  placeholder="Enter user ID..."
                  className="w-full bg-[#222] text-white px-3 py-2 rounded-lg border border-[#333] focus:outline-none focus:border-[#7681B3]"
                  autoFocus
                />
                <button
                  onClick={handleSearchUsers}
                  disabled={!searchFriends.trim()}
                  className="w-full bg-[#7681B3] text-white px-4 py-2 rounded-lg hover:bg-[#8b96d1] disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                >
                  Search
                </button>
              </div>
              
              <p className="text-gray-400 text-xs">Format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx</p>
              
              {searchError && (
                <div className="bg-red-900/20 border border-red-500/30 rounded-lg p-2">
                  <p className="text-red-400 text-xs">{searchError}</p>
                </div>
              )}
              
              {searchResults.length === 0 && searchFriends.trim() && !searchError && (
                <div className="bg-[#222] rounded-lg p-3 border border-[#333]">
                  <p className="text-gray-400 text-sm text-center">No users found with that ID</p>
                </div>
              )}
              
              {searchResults.length > 0 && (
                <div className="bg-[#222] rounded-lg border border-[#333] max-h-48 overflow-y-auto">
                  {searchResults.map(user => (
                    <div key={user.uniqueId} className="flex items-center justify-between py-3 px-3 border-b border-[#333] last:border-b-0">
                      <div className="flex items-center gap-3">
                        <div className="w-10 h-10 rounded-full bg-gray-500 flex items-center justify-center">
                          {user.avatarUrl ? (
                            <img src={user.avatarUrl} alt="avatar" className="w-full h-full object-cover rounded-full" />
                          ) : (
                            <div className="text-sm font-bold text-gray-300">
                              {user.username?.[0]?.toUpperCase() || 'U'}
                            </div>
                          )}
                        </div>
                        <div>
                          <div className="text-white text-sm font-medium">{user.username}</div>
                          <div className="text-[#7681B3] text-xs">@{user.uniqueId}</div>
                          {user.totalXP && (
                            <div className="text-[#C6AE82] text-xs">XP: {user.totalXP}</div>
                          )}
                        </div>
                      </div>
                      <button
                        className="bg-[#7681B3] text-white px-3 py-1 rounded text-xs hover:bg-[#8b96d1] transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                        onClick={() => handleAddFriend(user.uniqueId)}
                        disabled={addFriendStatus[user.uniqueId] === 'Friend request sent'}
                      >
                        {addFriendStatus[user.uniqueId] || 'Add Friend'}
                      </button>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>

          {/* Notifications Container - Black rounded rectangle */}
          <div className="bg-black rounded-2xl flex-1 flex flex-col overflow-hidden">
            {/* Notifications Header */}
            <div className="p-3 border-b border-[#333] flex justify-between items-center">
              <span className="text-white ml-2">Notifications</span>
              <Icon icon="material-symbols:notifications-outline" className="text-white w-5 h-5 mr-1" />
            </div>

            {/* Notification Items */}
            <div className="p-3 space-y-3 flex-1 overflow-y-auto">
              {notifications.length === 0 ? (
                <div className="text-[#7681B3] text-sm text-center py-4">
                  No notifications yet
                </div>
              ) : (
                notifications.map(notification => (
                  <div 
                    key={notification.id}
                    className={`bg-[#989572] rounded-lg p-3 cursor-pointer hover:bg-[#8a8768] transition-colors border-l-4 ${
                      notification.isRead ? 'border-[#666]' : 'border-[#C6AE82]'
                    }`}
                    onClick={() => handleNotificationClick(notification)}
                  >
                    <div className="flex items-start justify-between">
                      <div className="flex-1">
                        <h4 className="text-[#353844] text-sm font-semibold mb-1">
                          {notification.title}
                        </h4>
                        <p className="text-[#353844] text-sm">
                          {notification.message}
                          {notification.type === 'post_like' && notification.data?.postContent && (
                            <span className="block text-[#666] text-xs mt-1 italic">
                              Post: "{notification.data.postContent.substring(0, 50)}..."
                            </span>
                          )}
                          {notification.type === 'comment_like' && notification.data?.commentContent && (
                            <span className="block text-[#666] text-xs mt-1 italic">
                              Comment: "{notification.data.commentContent.substring(0, 50)}..."
                            </span>
                          )}
                        </p>
                        <p className="text-[#666] text-xs mt-1">
                          {new Date(notification.createdAt).toLocaleDateString()}
                        </p>
                      </div>
                      {!notification.isRead && (
                        <div className="w-2 h-2 bg-[#C6AE82] rounded-full ml-2"></div>
                      )}
                    </div>
                  </div>
                ))
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Post Creation Modal */}
      {showPostForm && (
        <div 
          className="fixed inset-0 bg-black/70 backdrop-blur-sm z-50 flex items-center justify-center"
          onClick={() => setShowPostForm(false)}
        >
          <div 
            className="bg-[#1B1B1D] rounded-xl p-5 w-full max-w-md mx-4 relative border border-[#7681B3]"
            onClick={(e) => e.stopPropagation()}
          >
            <div className="flex items-center mb-4">
              <div className="w-10 h-10 rounded-full bg-gray-500 mr-3"></div>
              <p className="text-white">
                Username <span className="text-[#7681B3]">@userid</span>
              </p>
            </div>
            
            <textarea
              value={newPostContent}
              onChange={(e) => setNewPostContent(e.target.value)}
              placeholder="Share something with friends :)"
              className="w-full bg-[#C6AE82] text-[#353844] p-4 rounded-lg border-none focus:outline-none mb-4"
              rows={4}
            />
            
            <div className="flex justify-end">
              <button
                onClick={handlePostSubmit}
                className="px-5 py-2 bg-black text-white rounded-lg hover:bg-gray-900 transition-colors"
              >
                Post
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Friend Request Modal */}
      {showFriendRequestModal && selectedFriendRequest && (
        <div 
          className="fixed inset-0 bg-black/70 backdrop-blur-sm z-50 flex items-center justify-center"
          onClick={() => setShowFriendRequestModal(false)}
        >
          <div 
            className="bg-[#1B1B1D] rounded-xl p-5 w-full max-w-md mx-4 relative border border-[#7681B3]"
            onClick={(e) => e.stopPropagation()}
          >
            <div className="text-center mb-6">
              <h3 className="text-white text-xl font-bold mb-2">Friend Request</h3>
              <p className="text-[#7681B3]">
                {selectedFriendRequest.data?.requesterUsername || 'Someone'} wants to be your friend!
              </p>
            </div>
            
            <div className="flex gap-3 justify-center">
              <button
                onClick={() => handleAcceptFriendRequest(selectedFriendRequest)}
                className="px-6 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors font-semibold"
              >
                Accept
              </button>
              <button
                onClick={() => handleRejectFriendRequest(selectedFriendRequest)}
                className="px-6 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors font-semibold"
              >
                Reject
              </button>
              <button
                onClick={() => setShowFriendRequestModal(false)}
                className="px-6 py-2 bg-gray-600 text-white rounded-lg hover:bg-gray-700 transition-colors font-semibold"
              >
                Cancel
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Custom Cursor */}
      <div
        className="cursor-glow"
        style={{
          left: `${mousePosition.x}px`,
          top: `${mousePosition.y}px`,
          transform: 'translate(-50%, -50%)'
        }}
      />
    </div>
  );
};

export default SocialFeed;