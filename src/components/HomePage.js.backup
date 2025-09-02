import React, { useState } from 'react';
import ProfileSection from './ProfileSection';
import { Icon } from '@iconify/react';
import SocialFeed from './SocialFeed';
import RoomsInterface from './RoomsInterface';

const HomePage = ({ user, refreshUserProfile, onLogout }) => {
  console.log("HomePage rendered with user:", user);
  const [activeTab, setActiveTab] = useState('profile');

  return (
    <div className="min-h-screen bg-[#000000]">
      {/* Permanent logo bar - never moves or hides */}
      <div className="fixed top-0 left-0 right-0 h-8 z-[60] bg-[#D8C2D8]">
        {/* Logo container */}
        <div className="absolute top-0 left-1/2 -translate-x-1/2 h-8 flex items-center">
          <div className="text-white text-lg font-bold">
            LOGO
          </div>
        </div>
        
        {/* Logout button */}
        <div className="absolute top-0 right-4 h-8 flex items-center">
          <button
            onClick={onLogout}
            className="text-white text-sm font-bold hover:text-red-300 transition-colors px-2 py-1 rounded"
          >
            LOGOUT
          </button>
        </div>
      </div>

      {/* Permanent navigation tabs */}
      <div className="fixed top-8 left-0 right-0 z-[60] h-8 flex items-center">
        <div className="flex justify-between gap-1 w-full px-4">
          <button
            className={`w-[32%] h-8 bg-[#D8C2D8] rounded-b-2xl flex items-center justify-center transition-all duration-300 ${
              activeTab === 'social' 
                ? 'bg-[#D8C2D8] scale-105 shadow-lg' 
                : 'bg-[#D8C2D8]/80 hover:bg-[#D8C2D8]'
            }`}
            onClick={() => setActiveTab('social')}
          >
            <Icon icon="material-symbols:group" className="w-5 h-5 text-white" />
          </button>
          <button
            className={`w-[32%] h-8 bg-[#D8C2D8] rounded-b-2xl flex items-center justify-center transition-all duration-300 ${
              activeTab === 'rooms' 
                ? 'bg-[#D8C2D8] scale-105 shadow-lg' 
                : 'bg-[#D8C2D8]/80 hover:bg-[#D8C2D8]'
            }`}
            onClick={() => setActiveTab('rooms')}
          >
            <Icon icon="material-symbols:door-open" className="w-5 h-5 text-white" />
          </button>
          <button
            className={`w-[32%] h-8 bg-[#D8C2D8] rounded-b-2xl flex items-center justify-center transition-all duration-300 ${
              activeTab === 'profile' 
                ? 'bg-[#D8C2D8] scale-105 shadow-lg' 
                : 'bg-[#D8C2D8]/80 hover:bg-[#D8C2D8]'
            }`}
            onClick={() => setActiveTab('profile')}
          >
            <Icon icon="material-symbols:person" className="w-5 h-5 text-white" />
          </button>
        </div>
      </div>

      {/* Content sections - adjusted top padding to account for permanent navigation */}
      <div className="pt-16 transition-all duration-300 ease-in-out">
        {activeTab === 'social' && (
          <div className="h-full w-full">
            <SocialFeed />
          </div>
        )}

        {activeTab === 'rooms' && (
          <div className="h-full w-full">
            <RoomsInterface />
          </div>
        )}

        {activeTab === 'profile' && <ProfileSection user={user} refreshUserProfile={refreshUserProfile} onLogout={onLogout} />}
      </div>
    </div>
  );
};

export default HomePage;