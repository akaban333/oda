import React, { useState, useEffect } from 'react';
import { Icon } from '@iconify/react';
import { postsAPI } from '../utils/api';

const PostsPanel = ({ isOpen, onClose, refreshTrigger }) => {
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  // Fetch all posts and filter for current user when panel opens or refreshTrigger changes
  useEffect(() => {
    if (isOpen) {
      const fetchPosts = async () => {
        try {
          setLoading(true);
          setError(null);
          // Get all posts like SocialFeed does
          const data = await postsAPI.getPosts();
          // Filter posts for current user
          const userPosts = (data.posts || []).filter(post => post.authorId === refreshTrigger);
          setPosts(userPosts);
        } catch (error) {
          console.error('Error fetching posts:', error);
          setError('Failed to load posts. Please try again.');
        } finally {
          setLoading(false);
        }
      };

      fetchPosts();
    }
  }, [isOpen, refreshTrigger]);

  const handleDelete = async (postId) => {
    try {
      await postsAPI.deletePost(postId);
      // Remove deleted post from state
      setPosts(posts.filter(post => (post.id || post._id) !== postId));
    } catch (error) {
      console.error('Error deleting post:', error);
      setError('Failed to delete post. Please try again.');
    }
  };

  return (
    <>
      {/* Backdrop */}
      {isOpen && (
        <div 
          className="fixed inset-0 bg-black/50 backdrop-blur-sm z-40"
          onClick={onClose}
        />
      )}

      {/* Posts Panel */}
      <div className={`fixed right-0 top-0 w-1/2 h-full bg-[#1f1f31] z-50 transform transition-transform duration-300 ease-in-out ${isOpen ? 'translate-x-0' : 'translate-x-full'} shadow-2xl`}>
        <div className="p-8 pt-20">
          {/* Header */}
          <div className="flex items-center justify-between mb-8">
            <h2 className="font-['Silkscreen'] text-3xl text-[#C5B178]">MY POSTS</h2>
            <button onClick={onClose}>
              <Icon icon="material-symbols:close" className="w-8 h-8 text-[#C5B178]" />
            </button>
          </div>



          {/* Posts List */}
          <div className="space-y-6">
            {loading && (
              <div className="text-center text-white/60">Loading posts...</div>
            )}
            {error && (
              <div className="text-center text-red-400">{error}</div>
            )}
            {!loading && !error && posts.length === 0 && (
              <div className="text-center text-white/60">
                No posts yet. Share your study progress!
              </div>
            )}
            {posts.map(post => (
              <div 
                key={post.id || post._id} 
                className="bg-[#2a2a42] rounded-xl p-6 relative border border-[#edc7cf]"
              >
                <div className="flex justify-between items-start mb-4">
                  <div className="flex-1">
                    <h3 className="font-['Silkscreen'] text-xl text-white mb-2">
                      {post.authorName || 'You'}
                    </h3>
                    <p className="font-mono text-gray-300 text-sm leading-relaxed">
                      {post.content}
                    </p>
                  </div>
                  <button 
                    onClick={() => handleDelete(post.id || post._id)}
                    className="text-[#edc7cf] hover:text-red-500 transition-colors ml-4"
                    title="Delete post"
                  >
                    <Icon icon="material-symbols:delete-outline" className="w-6 h-6" />
                  </button>
                </div>
                
                {/* Post Stats */}
                <div className="flex items-center gap-4 text-sm text-gray-400 mb-3">
                  <span className="flex items-center gap-1">
                    <Icon icon="mdi:heart" className="w-4 h-4" />
                    {post.likesCount || 0} likes
                  </span>
                  <span className="flex items-center gap-1">
                    <Icon icon="mdi:comment" className="w-4 h-4" />
                    {post.commentsCount || 0} comments
                  </span>
                  <span className="text-[#C5B178]">
                    {new Date(post.createdAt).toLocaleDateString()} at {new Date(post.createdAt).toLocaleTimeString()}
                  </span>
                </div>
                
                {/* Comments */}
                {post.comments && post.comments.length > 0 && (
                  <div className="border-t border-gray-600 pt-3">
                    <h4 className="text-[#C5B178] font-semibold mb-2">Comments:</h4>
                    <div className="space-y-2">
                      {post.comments.map(comment => (
                        <div key={comment._id} className="bg-[#1a1a2e] rounded-lg p-3">
                          <div className="flex justify-between items-start">
                            <div className="flex-1">
                              <p className="text-white text-sm font-semibold mb-1">
                                {comment.authorName || 'Unknown'}
                              </p>
                              <p className="text-gray-300 text-sm">{comment.content}</p>
                            </div>
                            <div className="flex items-center gap-2 ml-3">
                              <span className="text-xs text-gray-500">
                                {new Date(comment.createdAt).toLocaleDateString()}
                              </span>
                              <span className="text-xs text-gray-500">
                                {comment.likesCount || 0} likes
                              </span>
                            </div>
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>
      </div>
    </>
  );
};

export default PostsPanel; 