import { useState, useCallback, useEffect, useRef } from 'react';
import { authAPI } from '../utils/api';

export const useTimer = (sessionActive) => {
  const [isTimerRunning, setIsTimerRunning] = useState(false);
  const [workTime, setWorkTime] = useState(25);
  const [breakTime, setBreakTime] = useState(5);
  const [currentTime, setCurrentTime] = useState(25 * 60); // in seconds
  const [isBreak, setIsBreak] = useState(false);
  const [showTimerNotification, setShowTimerNotification] = useState(false);
  const [timerNotificationMessage, setTimerNotificationMessage] = useState("");
  const [timerInterval, setTimerInterval] = useState(null);
  
  const timerIntervalRef = useRef(null);

  // Toggle timer state - memoized with useCallback
  const toggleTimer = useCallback(() => {
    console.log('Toggle timer called, current state:', { isTimerRunning, currentTime, isBreak });
    
    if (!isTimerRunning) {
      // Start the timer
      console.log('Starting timer with work time:', workTime);
      setIsTimerRunning(true);
      setCurrentTime(workTime * 60);
      setIsBreak(false);
      
      const interval = setInterval(() => {
        setCurrentTime(prev => {
          console.log('Timer tick, remaining:', prev);
          if (prev <= 1) {
            // Timer finished
            console.log('Timer finished');
            clearInterval(interval);
            setIsTimerRunning(false);
            
            if (!isBreak) {
              // Work session finished, start break
              console.log('Work session finished, starting break');
              setIsBreak(true);
              setCurrentTime(breakTime * 60);
              setShowTimerNotification(true);
              setTimerNotificationMessage("Work session completed! Starting break...");
              
              // Start break timer
              const breakInterval = setInterval(() => {
                setCurrentTime(prev => {
                  if (prev <= 1) {
                    // Break finished
                    console.log('Break finished');
                    clearInterval(breakInterval);
                    setIsTimerRunning(false);
                    setIsBreak(false);
                    setShowTimerNotification(true);
                    setTimerNotificationMessage("Break completed! Ready for next session.");
                    
                    // Earn XP for completing a full Pomodoro cycle
                    if (sessionActive) {
                      authAPI.updateXP({
                        xp: 30,
                        source: 'pomodoro',
                        sessionId: null
                      }).catch(error => {
                        console.error('Error updating XP:', error);
                      });
                    }
                    return 0;
                  }
                  return prev - 1;
                });
              }, 1000);
            } else {
              // Break finished
              console.log('Break finished');
              setIsBreak(false);
              setShowTimerNotification(true);
              setTimerNotificationMessage("Break completed! Ready for next session.");
            }
            return 0;
          }
          return prev - 1;
        });
      }, 1000);
      
      setTimerInterval(interval);
      timerIntervalRef.current = interval;
    } else {
      // Stop the timer
      console.log('Stopping timer');
      if (timerIntervalRef.current) {
        clearInterval(timerIntervalRef.current);
        timerIntervalRef.current = null;
      }
      setIsTimerRunning(false);
      setCurrentTime(workTime * 60);
      setIsBreak(false);
      setShowTimerNotification(true);
      setTimerNotificationMessage("Timer stopped. No XP earned.");
    }
  }, [isTimerRunning, workTime, breakTime, sessionActive]);

  // Format time for display - memoized with useCallback
  const formatTime = useCallback((seconds) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
  }, []);

  // Close timer notification - memoized with useCallback
  const closeTimerNotification = useCallback(() => {
    setShowTimerNotification(false);
  }, []);

  // Cleanup timer on unmount
  useEffect(() => {
    return () => {
      if (timerIntervalRef.current) {
        clearInterval(timerIntervalRef.current);
      }
    };
  }, []);

  return {
    isTimerRunning,
    setIsTimerRunning,
    workTime,
    setWorkTime,
    breakTime,
    setBreakTime,
    currentTime,
    setCurrentTime,
    isBreak,
    setIsBreak,
    showTimerNotification,
    setShowTimerNotification,
    timerNotificationMessage,
    setTimerNotificationMessage,
    timerInterval,
    setTimerInterval,
    toggleTimer,
    formatTime,
    closeTimerNotification
  };
}; 