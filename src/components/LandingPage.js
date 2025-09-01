import React, { useEffect, useState } from 'react';
import Navbar from './Navbar';
import Door from './Door';

const LandingPage = ({ onLoginClick, onSignupClick, onLogin }) => {
  const [mousePosition, setMousePosition] = useState({ x: 0, y: 0 });

  useEffect(() => {
    const updateMousePosition = (e) => {
      setMousePosition({ x: e.clientX, y: e.clientY });
    };

    window.addEventListener('mousemove', updateMousePosition);

    return () => {
      window.removeEventListener('mousemove', updateMousePosition);
    };
  }, []);

  // Function to handle signup success from the door component
  const handleSignupSuccess = () => {
    onLogin();
  };

  return (
    <div className="min-h-screen bg-[#9494f2] relative overflow-hidden">
      {/* Custom cursor */}
      <div
        className="cursor-glow"
        style={{
          left: `${mousePosition.x}px`,
          top: `${mousePosition.y}px`,
          transform: 'translate(-50%, -50%)'
        }}
      />

      {/* Grid background */}
      <div 
        className="absolute inset-0" 
        style={{
          backgroundImage: 'linear-gradient(rgba(255, 255, 255, 0.1) 1.5px, transparent 1.5px), linear-gradient(90deg, rgba(255, 255, 255, 0.1) 1.5px, transparent 1.5px)',
          backgroundSize: '60px 60px'
        }}
      />

      {/* Content */}
      <Navbar onLoginClick={onLoginClick} onSignupClick={onSignupClick} />
      
      {/* Main content - Adjusted to be lower on the page */}
      <main className="flex-grow flex items-center justify-center pt-40 pb-20">
        <div className="text-center relative w-full max-w-4xl mx-auto">
          <div className="flex items-center justify-center gap-16">
            {/* Left text - ODA - with cute font */}
            <span className="text-[#ffe993] text-5xl font-bold text-glow" 
                  style={{ 
                    fontFamily: "'Comic Sans MS', 'Bubblegum Sans', cursive",
                    textShadow: "0 0 10px rgba(255, 233, 147, 0.7)"
                  }}>
              ODA
            </span>

            {/* Door component */}
            <div className="relative z-10">
              <Door onSignupSuccess={handleSignupSuccess} />
            </div>

            {/* Right text - SPC - with cute font */}
            <span className="text-[#ffe993] text-5xl font-bold text-glow"
                  style={{ 
                    fontFamily: "'Comic Sans MS', 'Bubblegum Sans', cursive",
                    textShadow: "0 0 10px rgba(255, 233, 147, 0.7)"
                  }}>
              SPC
            </span>
          </div>

          <p className="mt-12 text-[#ffe993] text-xl font-bold text-glow">
            where productive sessions bring us together
          </p>

          <p className="mt-6 text-white text-sm max-w-md mx-auto">
            Click on the door to sign up and join our community of focused learners and achievers.
          </p>
        </div>
      </main>
    </div>
  );
};

export default LandingPage; 