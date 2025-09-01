import React, { useState, useEffect } from 'react';
import SignupModal from './SignupModal';

const Door = ({ onSignupSuccess }) => {
  const [isSignupOpen, setIsSignupOpen] = useState(false);
  const [isHovered, setIsHovered] = useState(false);
  const [isClicked, setIsClicked] = useState(false);

  const handleDoorClick = () => {
    // Small delay to show the door opening animation before the modal appears
    setTimeout(() => {
      setIsSignupOpen(true);
      // Reset the clicked state after the modal appears
      setTimeout(() => setIsClicked(false), 300);
    }, 400);
  };

  const handleSignupSuccess = () => {
    setIsSignupOpen(false);
    if (onSignupSuccess) {
      onSignupSuccess();
    }
  };

  return (
    <div className="relative w-64 h-96 mx-auto">
      {/* Glow effect behind the door when hovered */}
      <div 
        className={`absolute inset-0 bg-[#ffe993]/30 rounded-full blur-2xl transform scale-75 transition-opacity duration-500 ${isHovered ? 'opacity-60' : 'opacity-0'}`}
      ></div>
      
      {/* Door frame - stays static */}
      <div 
        className="absolute inset-0 bg-gradient-to-b from-[#ffe993] to-[#e8c56b] rounded-lg p-3 group border border-[#1a1b3b] shadow-lg shadow-[#ffe993]/20 door-glow"
        onMouseEnter={() => setIsHovered(true)}
        onMouseLeave={() => setIsHovered(false)}
      >
        {/* Decorative elements for the door frame */}
        <div className="absolute top-1 left-1 right-1 bottom-1 pointer-events-none">
          <div className="absolute -top-4 -right-4 w-16 h-16 bg-[#ffe993]/10 rounded-full blur-lg"></div>
          <div className="absolute -bottom-4 -left-4 w-12 h-12 bg-[#ffe993]/10 rounded-full blur-lg"></div>
        </div>
        
        {/* Door interior - stays in frame */}
        <div className="relative w-full h-full bg-gradient-to-br from-[#9494f2] to-[#7e7ee0] rounded-lg border border-black/10 overflow-hidden">
          {/* Portal container */}
          <div className="absolute inset-0 flex items-center justify-center overflow-hidden z-10">
            {/* Spiral portal */}
            <div className="opacity-0 group-hover:opacity-100 transition-opacity duration-500 w-full h-full flex items-center justify-center">
              <div className="w-[300%] h-[300%] relative">
                <div className="absolute inset-0 spiral-pattern animate-spin-slow"></div>
                {/* Portal glow effect */}
                <div className="absolute inset-0 bg-gradient-radial from-[#ffe993] via-transparent to-transparent opacity-10 animate-pulse"></div>
              </div>
            </div>
            
            {/* Join Us text */}
            <div className="absolute inset-0 flex items-center justify-center">
              <span className="text-[#38395c] text-xl font-bold text-glow opacity-0 group-hover:opacity-100 transition-opacity duration-500 delay-300 z-20">
                JOIN US
              </span>
            </div>
          </div>
        </div>

        {/* The actual door that swings - now with click handler */}
        <div 
          className={`absolute inset-[3px] door-container cursor-pointer transition-transform ${
            isClicked ? 'transform scale-105' : ''
          }`}
          onClick={handleDoorClick}
        >
          {/* Door face */}
          <div 
            className={`absolute inset-0 bg-gradient-to-r from-[#ffe993] to-[#e8c56b] rounded-lg door-panel ${
              isHovered ? 'shadow-lg shadow-[#ffe993]/30' : ''
            }`}
          >
            {/* Decorative elements on door surface */}
            <div className="absolute inset-0 pointer-events-none overflow-hidden rounded-lg">
              <div className="absolute top-[10%] left-[20%] w-1.5 h-1.5 bg-[#ffe993]/60 rounded-full sparkle"></div>
              <div className="absolute bottom-[15%] right-[30%] w-2 h-2 bg-[#ffe993]/60 rounded-full sparkle delay-100"></div>
            </div>
            
            {/* Door handle */}
            <div className="absolute right-4 top-1/2 -translate-y-1/2 w-8 h-8 bg-[#ff9522] rounded-full shadow-[-4px_0_0_#21333e] hover:scale-110 transition-transform">
              {/* Handle sparkle effect */}
              <div className={`absolute w-2 h-2 rounded-full bg-white top-1 right-1 opacity-0 ${isHovered ? 'animate-pulse opacity-80' : ''}`}></div>
            </div>
            
            {/* Click indicator */}
            {isHovered && (
              <div className="absolute inset-0 flex items-center justify-center">
                <span className="bg-[#1a1b3b] text-[#ffe993] px-4 py-2 rounded-full text-sm font-medium animate-pulse transform hover:scale-105 transition-transform shadow-md shadow-black/20">
                  Click to Sign Up
                </span>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Signup Modal */}
      <SignupModal 
        isOpen={isSignupOpen} 
        onClose={() => setIsSignupOpen(false)} 
        onSuccess={handleSignupSuccess}
      />
    </div>
  );
};

export default Door; 