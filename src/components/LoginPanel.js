/* eslint-disable no-unused-vars */
import React, { useState, useEffect, useRef } from 'react';
import { authAPI } from '../utils/api';

const LoginPanel = ({ isOpen, onClose, onLogin }) => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const panelRef = useRef(null);

  useEffect(() => {
    const handleClickOutside = (event) => {
      if (panelRef.current && !panelRef.current.contains(event.target) && isOpen) {
        onClose();
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [isOpen, onClose]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setLoading(true);
    
    try {
      const response = await fetch("http://localhost:8080/api/v1/auth/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          email,
          password,
        }),
      });
      
      if (response.ok) {
        const data = await response.json();
        
        // Store tokens
        localStorage.setItem("accessToken", data.accessToken);
        if (data.refreshToken) {
          localStorage.setItem("refreshToken", data.refreshToken);
        }
        
        // Store user data
        localStorage.setItem("user", JSON.stringify(data.user));
        
        if (onLogin) {
          console.log("onLogin called with user:", data.user);
          onLogin(data.user);
        }
      } else {
        const data = await response.json();
        setError(data.error || "Login failed");
      }
    } catch (err) {
      console.error('Login error:', err);
      setError("Network error. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  if (!isOpen) return null;

  return (
    <>
      {/* Backdrop */}
      <div className="fixed inset-0 bg-black/50 backdrop-blur-sm z-[90]" onClick={onClose}></div>

      {/* Login panel */}
      <div className="fixed inset-0 flex items-center justify-center z-[100]">
        <div 
          ref={panelRef}
          className="w-[400px] bg-[#1a1b3b] rounded-lg relative overflow-hidden"
        >
          {/* Glowing border effect */}
          <div className="absolute inset-0 bg-gradient-to-r from-[#ffe993] via-[#e8c56b] to-[#ffe993] opacity-20"></div>
          <div className="absolute inset-[1px] bg-[#1a1b3b] rounded-lg"></div>

          {/* Close button */}
          <button 
            onClick={onClose}
            className="absolute top-4 right-4 text-[#ffe993] hover:text-[#e8c56b] transition-colors z-10"
          >
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>

          {/* Login form */}
          <div className="relative p-8 font-['Montserrat']">
            <h2 className="text-[#ffe993] text-2xl mb-6 text-center pt-4">Welcome back to ODA</h2>
            
            <form onSubmit={handleSubmit}>
              {error && (
                <div className="mb-4 text-center text-red-500 text-sm font-semibold">{error}</div>
              )}
              <div className="mb-8">
                <label htmlFor="email" className="block text-[#ffe993] mb-3 text-sm">
                  Email:
                </label>
                <input
                  type="email"
                  id="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  className="w-full h-10 px-4 bg-[#ffe993] border-2 border-[#ffe993] rounded-lg text-[#1a1b3b] text-sm focus:outline-none focus:border-[#e8c56b] transition-colors placeholder-[#1a1b3b]/50 flex items-center font-['Montserrat']"
                  placeholder="Enter your email"
                  required
                  disabled={loading}
                />
              </div>

              <div className="mb-8">
                <label htmlFor="password" className="block text-[#ffe993] mb-3 text-sm">
                  Password:
                </label>
                <input
                  type="password"
                  id="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="w-full h-10 px-4 bg-[#ffe993] border-2 border-[#ffe993] rounded-lg text-[#1a1b3b] text-sm focus:outline-none focus:border-[#e8c56b] transition-colors placeholder-[#1a1b3b]/50 flex items-center font-['Montserrat']"
                  placeholder="Enter your password"
                  required
                  disabled={loading}
                />
                <div className="flex justify-end mt-2">
                  <button
                    type="button"
                    className="text-[#bbce9d] hover:text-[#747a5d] text-xs transition-colors font-light tracking-wide"
                  >
                    Forgot password?
                  </button>
                </div>
              </div>

              <div className="flex justify-center mt-8">
                <button
                  type="submit"
                  disabled={loading}
                  className={`w-48 py-3 bg-[#ffe993] text-[#1a1b3b] rounded-lg transition-colors font-bold ${
                    loading 
                      ? 'opacity-50 cursor-not-allowed' 
                      : 'hover:bg-[#e8c56b]'
                  }`}
                >
                  {loading ? 'Logging In...' : 'Log In'}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </>
  );
};

export default LoginPanel; 