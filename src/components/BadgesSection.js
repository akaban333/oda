import React from 'react';
import { Icon } from '@iconify/react';

const BadgesSection = ({ sessions }) => {
  // Define badge thresholds and their icons
  const badges = [
    { threshold: 10, icon: 'mdi:medal-bronze', label: 'Beginner' },
    { threshold: 25, icon: 'mdi:medal-silver', label: 'Regular' },
    { threshold: 50, icon: 'mdi:medal-gold', label: 'Expert' },
    { threshold: 100, icon: 'mdi:medal-diamond', label: 'Master' }
  ];

  return (
    <div className="mt-8">
      <div className="flex items-start gap-4">
        <Icon icon="mdi:medal" className="w-8 h-8 text-[#C5B178] mt-2" />
        <div className="w-[3px] h-24 bg-[#C5B178]"></div>
        <div className="flex flex-col flex-grow">
          <h3 className="font-['Silkscreen'] text-2xl text-[#C5B178]">BADGES:</h3>
          <div className="mt-4 flex gap-4">
            {badges.map((badge, index) => (
              <div 
                key={index}
                className={`flex flex-col items-center ${sessions >= badge.threshold ? 'opacity-100' : 'opacity-30'}`}
              >
                <Icon icon={badge.icon} className="w-12 h-12 text-[#C5B178]" />
                <span className="font-mono text-sm text-white mt-1">{badge.label}</span>
              </div>
            ))}
          </div>
        </div>
        <div className="w-[3px] h-24 bg-[#C5B178]"></div>
      </div>
    </div>
  );
};

export default BadgesSection; 