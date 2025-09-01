import React, { useState, useEffect } from 'react';
import { Icon } from '@iconify/react';

const SignupModal = ({ isOpen = true, onClose = () => {}, onSuccess }) => {
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    password: '',
    confirmPassword: '',
    major: '',
    acceptTerms: false
  });
  
  const [success, setSuccess] = useState(false);
  const [serverError, setServerError] = useState('');
  const [errors, setErrors] = useState({});
  
  // Add keyboard event listener for Escape key
  useEffect(() => {
    const handleEscapeKey = (e) => {
      if (e.key === 'Escape' && isOpen) {
        onClose();
      }
    };

    window.addEventListener('keydown', handleEscapeKey);
    
    // Clean up event listener
    return () => {
      window.removeEventListener('keydown', handleEscapeKey);
    };
  }, [isOpen, onClose]);
  
  if (!isOpen) return null;

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData({
      ...formData,
      [name]: type === 'checkbox' ? checked : value
    });
  };

  const validateForm = () => {
    const newErrors = {};
    
    if (!formData.username.trim()) newErrors.username = 'Username is required';
    if (!formData.email.trim()) newErrors.email = 'Email is required';
    else if (!/\S+@\S+\.\S+/.test(formData.email)) newErrors.email = 'Email is invalid';
    
    if (!formData.password) newErrors.password = 'Password is required';
    else if (formData.password.length < 8) newErrors.password = 'Password must be at least 8 characters';
    
    if (formData.password !== formData.confirmPassword) newErrors.confirmPassword = 'Passwords do not match';
    if (!formData.acceptTerms) newErrors.acceptTerms = 'You must accept the terms and conditions';
    
    setErrors(newErrors);
    setServerError('');
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    console.log('Form submitted!');
    setServerError('');
    if (validateForm()) {
      console.log('Validation passed, sending request...');
      try {
        const response = await fetch('http://localhost:8080/api/v1/auth/register', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            username: formData.username,
            email: formData.email,
            password: formData.password
          })
        });
        console.log('Fetch response status:', response.status);
        if (response.status === 201) {
          setSuccess(true);
          setTimeout(() => {
            if (onSuccess) onSuccess();
            else onClose();
          }, 2000);
        } else {
          const data = await response.json();
          setServerError(data.error || 'Registration failed');
          console.error('Registration error:', data.error);
        }
      } catch (err) {
        setServerError('Network error. Please try again.');
        console.error('Network error:', err);
      }
    } else {
      console.log('Validation failed:', serverError);
    }
  };

  // Click outside to close
  const handleBackdropClick = (e) => {
    if (e.target === e.currentTarget) {
      onClose();
    }
  };

  return (
    <div 
      className="fixed inset-0 bg-black/70 backdrop-blur-sm z-50 flex items-center justify-center"
      onClick={handleBackdropClick}
    >
      <div className="bg-gradient-to-br from-[#2a2a42] to-[#1f1f31] rounded-lg p-4 max-w-xs w-full mx-4 animate-fade-in relative shadow-[0_0_20px_rgba(255,233,147,0.15)] border border-[#ffe993]/20 door-glow">
        {/* Decorative elements */}
        <div className="absolute top-0 left-0 w-full h-full overflow-hidden rounded-lg pointer-events-none">
          <div className="absolute -top-8 -right-8 w-16 h-16 bg-[#ffe993]/10 rounded-full blur-lg"></div>
          <div className="absolute -bottom-4 -left-4 w-12 h-12 bg-[#ffe993]/10 rounded-full blur-lg"></div>
          
          {/* Particle effects - reduced */}
          <div className="absolute top-8 right-8 w-1.5 h-1.5 bg-[#ffe993]/60 rounded-full sparkle"></div>
          <div className="absolute bottom-12 left-6 w-1 h-1 bg-[#ffe993]/60 rounded-full sparkle delay-100"></div>
        </div>
        
        <button
          onClick={onClose}
          className="absolute top-2 right-2 text-[#ffe993] hover:text-white transition-colors bg-[#1a1b3b]/50 rounded-full p-1 hover:bg-[#1a1b3b] hover:scale-110 transform duration-200 z-10"
        >
          <Icon icon="material-symbols:close" className="w-3 h-3" />
        </button>
        
        {!success ? (
          <>
            <div className="mb-4">
              <h2 className="text-xl font-bold text-[#ffe993] text-glow">Join ODA SPC</h2>
            </div>
            
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="space-y-3">
                <div className="group">
                  <label className="text-left block text-[#ffe993] text-xs font-medium">Username:</label>
                  <div className="relative">
                    <div className="absolute left-2 top-1/2 -translate-y-1/2 text-[#ffe993]">
                      <Icon icon="mdi:account" className="w-3 h-3" />
                    </div>
                    <input
                      type="text"
                      name="username"
                      value={formData.username}
                      onChange={handleChange}
                      className={`w-full bg-[#1a1b3b]/70 border ${errors.username ? 'border-red-500' : 'border-[#ffe993]/30 group-hover:border-[#ffe993]/70'} rounded px-2 py-1.5 pl-7 text-white text-xs placeholder:text-[#ffe993]/40 focus:border-[#ffe993] focus:ring-1 focus:ring-[#ffe993]/50 outline-none transition-all duration-200 input-focus-effect`}
                      placeholder="Choose a username"
                    />
                  </div>
                  {errors.username && <p className="text-red-400 text-[10px] mt-0.5 pl-1">{errors.username}</p>}
                </div>
                
                <div className="group">
                  <label className="text-left block text-[#ffe993] text-xs font-medium">Email:</label>
                  <div className="relative">
                    <div className="absolute left-2 top-1/2 -translate-y-1/2 text-[#ffe993]">
                      <Icon icon="mdi:email-outline" className="w-3 h-3" />
                    </div>
                    <input
                      type="email"
                      name="email"
                      value={formData.email}
                      onChange={handleChange}
                      className={`w-full bg-[#1a1b3b]/70 border ${errors.email ? 'border-red-500' : 'border-[#ffe993]/30 group-hover:border-[#ffe993]/70'} rounded px-2 py-1.5 pl-7 text-white text-xs placeholder:text-[#ffe993]/40 focus:border-[#ffe993] focus:ring-1 focus:ring-[#ffe993]/50 outline-none transition-all duration-200 input-focus-effect`}
                      placeholder="your@email.com"
                    />
                  </div>
                  {errors.email && <p className="text-red-400 text-[10px] mt-0.5 pl-1">{errors.email}</p>}
                </div>
                
                <div className="group">
                  <label className="text-left block text-[#ffe993] text-xs font-medium">Password:</label>
                  <div className="relative">
                    <div className="absolute left-2 top-1/2 -translate-y-1/2 text-[#ffe993]">
                      <Icon icon="mdi:lock-outline" className="w-3 h-3" />
                    </div>
                    <input
                      type="password"
                      name="password"
                      value={formData.password}
                      onChange={handleChange}
                      className={`w-full bg-[#1a1b3b]/70 border ${errors.password ? 'border-red-500' : 'border-[#ffe993]/30 group-hover:border-[#ffe993]/70'} rounded px-2 py-1.5 pl-7 text-white text-xs placeholder:text-[#ffe993]/40 focus:border-[#ffe993] focus:ring-1 focus:ring-[#ffe993]/50 outline-none transition-all duration-200 input-focus-effect`}
                      placeholder="Create a password"
                    />
                  </div>
                  {errors.password && <p className="text-red-400 text-[10px] mt-0.5 pl-1">{errors.password}</p>}
                </div>
                
                <div className="group">
                  <label className="text-left block text-[#ffe993] text-xs font-medium">Confirm Password:</label>
                  <div className="relative">
                    <div className="absolute left-2 top-1/2 -translate-y-1/2 text-[#ffe993]">
                      <Icon icon="mdi:lock-check-outline" className="w-3 h-3" />
                    </div>
                    <input
                      type="password"
                      name="confirmPassword"
                      value={formData.confirmPassword}
                      onChange={handleChange}
                      className={`w-full bg-[#1a1b3b]/70 border ${errors.confirmPassword ? 'border-red-500' : 'border-[#ffe993]/30 group-hover:border-[#ffe993]/70'} rounded px-2 py-1.5 pl-7 text-white text-xs placeholder:text-[#ffe993]/40 focus:border-[#ffe993] focus:ring-1 focus:ring-[#ffe993]/50 outline-none transition-all duration-200 input-focus-effect`}
                      placeholder="Confirm your password"
                    />
                  </div>
                  {errors.confirmPassword && <p className="text-red-400 text-[10px] mt-0.5 pl-1">{errors.confirmPassword}</p>}
                </div>
                
                <div className="group">
                  <label className="text-left block text-[#ffe993] text-xs font-medium">Major/Field:</label>
                  <div className="relative">
                    <div className="absolute left-2 top-1/2 -translate-y-1/2 text-[#ffe993]">
                      <Icon icon="mdi:book-education-outline" className="w-3 h-3" />
                    </div>
                    <input
                      type="text"
                      name="major"
                      value={formData.major}
                      onChange={handleChange}
                      className="w-full bg-[#1a1b3b]/70 border border-[#ffe993]/30 group-hover:border-[#ffe993]/70 rounded px-2 py-1.5 pl-7 text-white text-xs placeholder:text-[#ffe993]/40 focus:border-[#ffe993] focus:ring-1 focus:ring-[#ffe993]/50 outline-none transition-all duration-200 input-focus-effect"
                      placeholder="E.g., Computer Science"
                    />
                  </div>
                </div>
              </div>
              
              <div className="flex items-start p-1.5 bg-[#1a1b3b]/40 rounded border border-[#ffe993]/10 hover:border-[#ffe993]/30 transition-colors mt-3">
                <input
                  type="checkbox"
                  name="acceptTerms"
                  checked={formData.acceptTerms}
                  onChange={handleChange}
                  className="mt-0.5 w-3 h-3 accent-[#ffe993] rounded border-[#ffe993]"
                />
                <label className="ml-1.5 text-white/90 text-[10px]">
                  I agree to the <a href="#" className="text-[#ffe993] hover:underline">Terms</a> and <a href="#" className="text-[#ffe993] hover:underline">Privacy Policy</a>
                </label>
              </div>
              {errors.acceptTerms && <p className="text-red-400 text-[10px] pl-1 mt-1">{errors.acceptTerms}</p>}
              
              <button
                type="submit"
                className="w-full bg-gradient-to-r from-[#ffe993] to-[#e8c56b] text-[#1a1b3b] py-1.5 rounded text-xs font-bold hover:from-[#e8c56b] hover:to-[#d8b967] transition-all duration-300 shadow shadow-[#ffe993]/20 hover:shadow-md hover:shadow-[#ffe993]/30 transform hover:translate-y-[-1px] active:translate-y-[1px] button-hover-effect animate-gradient mt-4"
              >
                Create Account
              </button>
              
              <div className="text-center mt-3">
                <p className="text-white/60 text-[10px]">
                  Already have an account? <a href="#" className="text-[#ffe993] hover:underline">Log In</a>
                </p>
              </div>
            </form>
          </>
        ) : (
          <div className="text-left py-2">
            <div className="relative w-12 h-12 mx-auto mb-2">
              <div className="absolute inset-0 bg-green-500/10 rounded-full animate-pulse"></div>
              <div className="relative flex items-center justify-center bg-[#1a1b3b] rounded-full p-2 border border-green-500/30">
                <Icon icon="mdi:check-circle" className="w-6 h-6 text-green-500" />
              </div>
            </div>
            
            <h2 className="text-lg font-bold text-[#ffe993] mb-1 text-glow text-center">Account Created!</h2>
            <p className="text-white/80 text-xs mb-3 text-center">
              Your account has been created.<br />Welcome to ODA SPC!
            </p>
            
            <div className="relative max-w-xs mx-auto overflow-hidden rounded">
              <div className="absolute inset-0 bg-gradient-to-r from-[#ffe993]/20 via-transparent to-[#ffe993]/20 animate-pulse"></div>
              <div className="w-full bg-gradient-to-r from-[#ffe993] to-[#e8c56b] text-[#1a1b3b] py-1.5 rounded font-bold relative animate-gradient">
                <div className="flex items-center justify-center">
                  <span className="text-xs">Redirecting</span>
                  <span className="ml-1 flex space-x-0.5">
                    <span className="animate-bounce delay-100">.</span>
                    <span className="animate-bounce delay-200">.</span>
                    <span className="animate-bounce delay-300">.</span>
                  </span>
                </div>
              </div>
            </div>
          </div>
        )}
        {serverError && <p className="text-red-400 text-xs text-center mt-2">{serverError}</p>}
      </div>
    </div>
  );
};

export default SignupModal; 