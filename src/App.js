import React, { useState, useEffect } from 'react';
import LandingPage from './components/LandingPage';
import HomePage from './components/HomePage';
import LoginPanel from './components/LoginPanel';
import SignupModal from './components/SignupModal';
import Navbar from './components/Navbar';
import { authAPI } from './utils/api';
import './styles.css';

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [user, setUser] = useState(null);
  const [isLoginOpen, setIsLoginOpen] = useState(false);
  const [isSignupOpen, setIsSignupOpen] = useState(false);
  const [loading, setLoading] = useState(true);

  // Fetch user profile from backend using JWT
  async function fetchUserProfile() {
    try {
      const data = await authAPI.getProfile();
      return data.user;
    } catch (error) {
      console.error('Error fetching user profile:', error);
      throw error;
    }
  }

  // Check if token is valid and refresh if needed
  const checkAndRefreshToken = async () => {
    const token = localStorage.getItem("accessToken");
    const refreshToken = localStorage.getItem("refreshToken");
    
    if (!token) return false;
    
    try {
      // Try to get profile with current token
      const userProfile = await fetchUserProfile();
      setUser(userProfile);
      setIsAuthenticated(true);
      return true;
    } catch (error) {
      // If token is expired, try to refresh
      if (refreshToken && error.message.includes('401')) {
        try {
          const response = await fetch('http://localhost:8080/api/v1/auth/refresh', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ refreshToken })
          });
          
          if (response.ok) {
            const data = await response.json();
            localStorage.setItem("accessToken", data.accessToken);
            if (data.refreshToken) {
              localStorage.setItem("refreshToken", data.refreshToken);
            }
            
            const userProfile = await fetchUserProfile();
            setUser(userProfile);
            setIsAuthenticated(true);
            return true;
          }
        } catch (refreshError) {
          console.error('Token refresh failed:', refreshError);
        }
      }
      
      // Clear invalid tokens
      localStorage.removeItem("accessToken");
      localStorage.removeItem("refreshToken");
      localStorage.removeItem("user");
      return false;
    }
  };

  // On initial load, check for token and fetch user profile
  useEffect(() => {
    const initializeAuth = async () => {
      const isTokenValid = await checkAndRefreshToken();
      if (!isTokenValid) {
        setIsAuthenticated(false);
        setUser(null);
      }
      setLoading(false);
    };

    initializeAuth();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // Handler for successful login
  const handleLogin = async (userObj) => {
    try {
      const userProfile = await fetchUserProfile();
      setUser(userProfile);
    } catch (error) {
      console.error('Error fetching user profile after login:', error);
      setUser(userObj); // fallback to login response
    }
    setIsAuthenticated(true);
    setIsLoginOpen(false);
  };

  // Handler for successful registration
  const handleSignupSuccess = () => {
    setIsSignupOpen(false);
    setIsLoginOpen(true); // Open login after successful registration
  };

  // Handler to open login modal
  const openLoginModal = () => setIsLoginOpen(true);
  // Handler to close login modal
  const closeLoginModal = () => setIsLoginOpen(false);
  
  // Handler to open signup modal
  const openSignupModal = () => setIsSignupOpen(true);
  // Handler to close signup modal
  const closeSignupModal = () => setIsSignupOpen(false);

  // Refresh user profile (for profile edit)
  const refreshUserProfile = async () => {
    try {
      const userProfile = await fetchUserProfile();
      setUser(userProfile);
    } catch (error) {
      console.error('Error refreshing user profile:', error);
      setIsAuthenticated(false);
      setUser(null);
      localStorage.removeItem("accessToken");
      localStorage.removeItem("refreshToken");
      localStorage.removeItem("user");
      setIsLoginOpen(true);
    }
  };

  // Logout handler
  const handleLogout = async () => {
    try {
      await authAPI.logout();
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      localStorage.removeItem("accessToken");
      localStorage.removeItem("refreshToken");
      localStorage.removeItem("user");
      setIsAuthenticated(false);
      setUser(null);
    }
  };

  // Show loading screen while checking authentication
  if (loading) {
    return (
      <div className="App pixel-font flex items-center justify-center min-h-screen bg-black">
        <div className="text-white text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-white mx-auto mb-4"></div>
          <p>Loading...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="App pixel-font">
      {!isAuthenticated && <Navbar onLoginClick={openLoginModal} onSignupClick={openSignupModal} />}
      <LoginPanel isOpen={isLoginOpen} onClose={closeLoginModal} onLogin={handleLogin} />
      <SignupModal isOpen={isSignupOpen} onClose={closeSignupModal} onSuccess={handleSignupSuccess} />
      {!isAuthenticated ? (
        <LandingPage onLoginClick={openLoginModal} onSignupClick={openSignupModal} onLogin={handleLogin} />
      ) : (
        <HomePage user={user} refreshUserProfile={refreshUserProfile} onLogout={handleLogout} />
      )}
    </div>
  );
}

export default App;
