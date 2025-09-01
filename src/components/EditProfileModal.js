import React, { useState } from 'react';
import { authAPI } from '../utils/api';

const EditProfileModal = ({ isOpen, onClose, user, onProfileUpdated }) => {
  const [username, setUsername] = useState(user?.username || '');
  const [bio, setBio] = useState(user?.bio || '');
  const [avatarUrl, setAvatarUrl] = useState(user?.avatarUrl || '');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleSave = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const body = { username }; // Always send username
      if (bio.trim() !== '') body.bio = bio;
      if (avatarUrl.trim() !== '') body.avatarUrl = avatarUrl;
      
      await authAPI.updateProfile(body);
      
      // Call the callback to refresh user profile
      if (onProfileUpdated) {
        onProfileUpdated();
      }
      
      onClose();
    } catch (error) {
      console.error('Error updating profile:', error);
      setError('Failed to update profile. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black/70 backdrop-blur-sm z-50 flex items-center justify-center">
      <div className="bg-[#1f1f31] rounded-xl p-8 max-w-2xl w-full mx-4 relative">
        <button
          onClick={onClose}
          className="absolute top-4 right-4 text-[#C5B178] hover:text-white transition-colors"
        >
          X
        </button>
        <h2 className="text-2xl font-bold text-white mb-4">Edit Profile</h2>
        <div className="flex flex-col gap-4">
          <input
            type="text"
            id="username"
            name="username"
            placeholder="Username"
            className="bg-[#2a2a42] border border-[#C5B178] rounded-lg px-4 py-2 text-white font-mono"
            value={username}
            onChange={e => setUsername(e.target.value)}
          />
          <input
            type="text"
            id="avatarUrl"
            name="avatarUrl"
            placeholder="Profile Picture URL"
            className="bg-[#2a2a42] border border-[#C5B178] rounded-lg px-4 py-2 text-white font-mono"
            value={avatarUrl}
            onChange={e => setAvatarUrl(e.target.value)}
          />
          <input
            type="text"
            id="bio"
            name="bio"
            placeholder="Bio"
            className="bg-[#2a2a42] border border-[#C5B178] rounded-lg px-4 py-2 text-white font-mono"
            value={bio}
            onChange={e => setBio(e.target.value)}
          />
        </div>
        {error && <div className="text-red-400 mt-2">{error}</div>}
        <button
          className="mt-6 px-6 py-2 bg-[#C5B178] text-black rounded-lg font-mono text-sm hover:bg-[#d4b980] transition-colors"
          onClick={handleSave}
          disabled={loading}
        >
          {loading ? 'Saving...' : 'Save'}
        </button>
      </div>
    </div>
  );
};

export default EditProfileModal;
