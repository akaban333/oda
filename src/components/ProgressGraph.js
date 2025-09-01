import React, { useState, useEffect } from 'react';
import { Icon } from '@iconify/react';
import { sessionsAPI } from '../utils/api';

const ProgressGraph = ({ sessions, xp }) => {
  const [graphData, setGraphData] = useState({
    sessions: [0, 0, 0, 0, 0, 0, 0],
    xp: [0, 0, 0, 0, 0, 0, 0],
    labels: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
  });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchSessionStats = async () => {
      try {
        setLoading(true);
        const data = await sessionsAPI.getStats();
        
        // Process the data to show last 7 days
        if (data && data.dailyStats) {
          const last7Days = data.dailyStats.slice(-7);
          setGraphData({
            sessions: last7Days.map(day => day.sessionsCount || 0),
            xp: last7Days.map(day => day.totalXP || 0),
            labels: last7Days.map(day => {
              const date = new Date(day.date);
              return date.toLocaleDateString('en', { weekday: 'short' });
            })
          });
        }
      } catch (error) {
        console.error('Error fetching session stats:', error);
        // Keep default data on error
      } finally {
        setLoading(false);
      }
    };

    fetchSessionStats();
  }, []);

  const maxSessions = Math.max(...graphData.sessions);
  const maxXP = Math.max(...graphData.xp);

  return (
    <div>
      <div className="flex items-start gap-4">
        <Icon icon="mdi:chart-line" className="w-8 h-8 text-[#edc7cf] mt-2" />
        <div className="w-[3px] bg-[#edc7cf] self-stretch"></div>
        <div className="flex flex-col flex-grow">
          <h3 className="font-['Silkscreen'] text-xl text-[#edc7cf]">PROGRESS:</h3>
          <div className="mt-2 bg-[#2a2a42] rounded-xl p-3">
            <div className="flex justify-between items-end h-24">
              {graphData.sessions.map((value, index) => (
                <div key={index} className="flex flex-col items-center">
                  <div className="relative">
                    <div 
                      className="w-3 bg-[#C5B178] rounded-t"
                      style={{ height: `${(value / maxSessions) * 100}%` }}
                    />
                    <div 
                      className="w-3 bg-[#edc7cf] rounded-t absolute bottom-0"
                      style={{ height: `${(graphData.xp[index] / maxXP) * 100}%` }}
                    />
                  </div>
                  <span className="font-mono text-xs text-gray-400 mt-1">{graphData.labels[index]}</span>
                </div>
              ))}
            </div>
            <div className="flex justify-center gap-3 mt-1">
              <div className="flex items-center gap-1">
                <div className="w-2 h-2 bg-[#C5B178] rounded"></div>
                <span className="font-mono text-xs text-gray-400">Sessions</span>
              </div>
              <div className="flex items-center gap-1">
                <div className="w-2 h-2 bg-[#edc7cf] rounded"></div>
                <span className="font-mono text-xs text-gray-400">XP</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ProgressGraph; 