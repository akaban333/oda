/* eslint-disable no-unused-vars */
import { useState, useCallback, useEffect } from 'react';
import { sessionsAPI, authAPI } from '../utils/api';

export const useSession = (currentUser) => {
  const [sessionActive, setSessionActive] = useState(false);
  const [sessionXP, setSessionXP] = useState(0);
  const [inactiveTime, setInactiveTime] = useState(0);
  const [showSessionXPModal, setShowSessionXPModal] = useState(false);
  const [lastUserActivity, setLastUserActivity] = useState(Date.now());

  // Functions for session management - memoized with useCallback
  const startSession = useCallback(async () => {
    try {
      // Start a new study session
      const sessionData = await sessionsAPI.startSession({
        startTime: new Date().toISOString()
      });
      
      setSessionActive(true);
      setSessionXP(0);
      setInactiveTime(0);
      setLastUserActivity(Date.now());
    } catch (error) {
      console.error('Error starting session:', error);
      throw new Error('Failed to start session. Please try again.');
    }
  }, []);

  const endSession = useCallback(async () => {
    try {
      // Get current session XP value at the time of calling
      const currentSessionXP = sessionXP;
      
      // End the current study session and update XP
      await sessionsAPI.endSession(null, {
        endTime: new Date().toISOString(),
        earnedXP: currentSessionXP
      });
      
      // Update user's total XP
      if (currentSessionXP > 0) {
        await authAPI.updateXP({
          xp: currentSessionXP,
          source: 'session',
          sessionId: null
        });
      }
      
      setSessionActive(false);
      setShowSessionXPModal(true);
    } catch (error) {
      console.error('Error ending session:', error);
      throw new Error('Failed to end session. Please try again.');
    }
  }, [sessionXP]);

  // Activity handler for session tracking - memoized with useCallback and debounced
  const activityHandler = useCallback((event) => {
    // Don't track activity for input fields, buttons, or form elements
    const target = event.target;
    const tagName = target.tagName;
    const className = target.className || '';
    const id = target.id || '';
    
    // Check if the target is an input, textarea, button, select, or any form element
    if (tagName === 'INPUT' || 
        tagName === 'TEXTAREA' || 
        tagName === 'BUTTON' || 
        tagName === 'SELECT' ||
        className.includes('input') ||
        className.includes('textarea') ||
        className.includes('button') ||
        className.includes('select') ||
        id.includes('input') ||
        id.includes('textarea') ||
        id.includes('button') ||
        id.includes('select') ||
        target.closest('input') ||
        target.closest('textarea') ||
        target.closest('button') ||
        target.closest('select') ||
        target.closest('form')) {
      return;
    }
    
    // Debounce activity updates to reduce state changes
    if (activityHandler.timeout) {
      clearTimeout(activityHandler.timeout);
    }
    
    activityHandler.timeout = setTimeout(() => {
      setLastUserActivity(Date.now());
      setInactiveTime(0);
    }, 100); // 100ms debounce
  }, []);

  // Handle user activity with useEffect
  useEffect(() => {
    if (!sessionActive) return;

    // Update XP counter every 5 minutes (10 XP per 5 minutes) to reduce API calls
    const xpInterval = setInterval(async () => {
      try {
        // Update session XP
        setSessionXP(prev => prev + 10);
        
        // Update user's total XP in the backend
        await authAPI.updateXP({
          xp: 10,
          source: 'session',
          sessionId: null
        });
      } catch (error) {
        console.error('Error updating XP:', error);
      }
    }, 300000); // every 5 minutes instead of every minute

    // Check for inactivity every minute
    const inactivityChecker = setInterval(() => {
      const now = Date.now();
      const timeSinceLastActivity = now - lastUserActivity;
      const inactiveMinutes = Math.floor(timeSinceLastActivity / 60000);
      
      if (inactiveMinutes >= 10) {
        // User has been inactive for 10 minutes, end session
        endSession();
      } else {
        setInactiveTime(inactiveMinutes);
      }
    }, 60000);

    // Add event listeners for user activity
    window.addEventListener('mousemove', activityHandler);
    window.addEventListener('keypress', activityHandler);
    window.addEventListener('scroll', activityHandler);
    window.addEventListener('click', activityHandler);

    return () => {
      clearInterval(xpInterval);
      clearInterval(inactivityChecker);
      window.removeEventListener('mousemove', activityHandler);
      window.removeEventListener('keypress', activityHandler);
      window.removeEventListener('scroll', activityHandler);
      window.removeEventListener('click', activityHandler);
    };
  }, [sessionActive, lastUserActivity, endSession, activityHandler]);

  return {
    sessionActive,
    setSessionActive,
    sessionXP,
    setSessionXP,
    inactiveTime,
    setInactiveTime,
    showSessionXPModal,
    setShowSessionXPModal,
    lastUserActivity,
    setLastUserActivity,
    startSession,
    endSession,
    activityHandler
  };
}; 