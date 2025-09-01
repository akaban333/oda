import React from 'react';
import { Icon } from '@iconify/react';
import PropTypes from 'prop-types';

const RoomSidebar = ({ 
  menuExpanded, 
  toggleMenu, 
  activeSection, 
  setActiveSection, 
  currentRoom, 
  sharedRooms, 
  handleRoomSelection, 
  refreshRooms 
}) => {
  return (
    <div 
      className={`bg-[#7681B3] h-screen relative z-10 border-r border-[#edc7cf] border-opacity-30 transition-all duration-300 ease-in-out ${menuExpanded ? 'w-64' : 'w-14'} flex flex-col justify-between`}
      role="navigation"
      aria-label="Main Navigation"
    >
      {/* Top Section - Header and Navigation */}
      <div className="flex-1 overflow-hidden">
        {/* Header */}
        <div className="pt-16 px-2">
          <h3 className={`text-white text-xl mb-4 ${menuExpanded ? 'block' : 'hidden'}`}>Room Menu</h3>
        </div>
        
        {/* Navigation Items */}
        <ul className="space-y-2 px-2">
          {/* User's Room */}
          <li>
            <button 
              className={`w-full p-2 rounded cursor-pointer flex items-center ${menuExpanded ? 'gap-2 mx-2' : 'justify-center'} transition-all duration-200 
                ${activeSection === 'userRoom' 
                  ? 'bg-white bg-opacity-20 text-white shadow-[0_0_10px_rgba(255,255,255,0.5)]' 
                  : 'text-white hover:bg-white hover:bg-opacity-10 hover:shadow-[0_0_8px_rgba(255,255,255,0.3)]'}`}
              onClick={() => setActiveSection('userRoom')}
              aria-label="Your Room"
              aria-current={activeSection === 'userRoom' ? 'page' : undefined}
            >
              <Icon icon="mdi:home" className="w-5 h-5 text-white" />
              {menuExpanded && <span className="text-sm">Your Room</span>}
            </button>
          </li>
          
          {/* Current Room - only show if in a room */}
          {currentRoom && (
            <li>
              <button 
                className={`w-full p-2 rounded cursor-pointer flex items-center ${menuExpanded ? 'gap-2 mx-2' : 'justify-center'} transition-all duration-200 
                  ${activeSection === 'sharedSpace' 
                    ? 'bg-white bg-opacity-20 text-white shadow-[0_0_10px_rgba(255,255,255,0.5)]' 
                    : 'text-white hover:bg-white hover:bg-opacity-10 hover:shadow-[0_0_8px_rgba(255,255,255,0.3)]'}`}
                onClick={() => setActiveSection('sharedSpace')}
                aria-label={`Current Room: ${currentRoom.name}`}
                aria-current={activeSection === 'sharedSpace' ? 'page' : undefined}
              >
                <Icon icon="mdi:door-open" className="w-5 h-5 text-white" />
                {menuExpanded && <span className="text-sm truncate">{currentRoom.name}</span>}
              </button>
            </li>
          )}
          
          {/* Shared Rooms */}
          <li>
            <button 
              className={`w-full p-2 rounded cursor-pointer flex items-center ${menuExpanded ? 'gap-2 mx-2' : 'justify-center'} transition-all duration-200 
                ${activeSection === 'sharedRooms' 
                  ? 'bg-white bg-opacity-20 text-white shadow-[0_0_10px_rgba(255,255,255,0.5)]' 
                  : 'text-white hover:bg-white hover:bg-opacity-10 hover:shadow-[0_0_8px_rgba(255,255,255,0.3)]'}`}
              onClick={() => {
                setActiveSection('sharedRooms');
                // Refresh rooms when switching to shared rooms view
                refreshRooms();
              }}
              aria-label="Shared Rooms"
              aria-current={activeSection === 'sharedRooms' ? 'page' : undefined}
            >
              <Icon icon="mdi:account-group" className="w-4" />
              {menuExpanded && <span className="text-sm">Shared Rooms</span>}
            </button>
          </li>
          
          {/* Individual Shared Rooms List */}
          {menuExpanded && sharedRooms.map((room) => (
            <li key={room.id || room._id} className="ml-4">
              <button 
                className={`w-full p-2 rounded cursor-pointer flex items-center gap-2 transition-all duration-200 
                  ${activeSection === 'sharedSpace' && currentRoom && (currentRoom.id || currentRoom._id) === (room.id || room._id)
                    ? 'bg-white bg-opacity-20 text-white shadow-[0_0_10px_rgba(255,255,255,0.5)]' 
                    : 'text-white hover:bg-white hover:bg-opacity-10 hover:shadow-[0_0_8px_rgba(255,255,255,0.3)]'}`}
                onClick={() => handleRoomSelection(room)}
                aria-label={`Join Room: ${room.name}`}
              >
                <Icon icon="mdi:door" className="w-4 h-4 text-white" />
                <span className="text-xs truncate">{room.name}</span>
                <span className="text-xs text-white text-opacity-70 ml-auto">
                  {room.participantCount || 0}/{room.maxParticipants || 4}
                </span>
              </button>
            </li>
          ))}
        </ul>
      </div>
      
      {/* Bottom Section - Toggle Button */}
      <div className="p-2 border-t border-white border-opacity-20 bg-[#7681B3]">
        <div className="flex justify-center">
          <button
            onClick={toggleMenu}
            className="bg-white bg-opacity-20 rounded-full p-1.5 hover:bg-opacity-30 transition-all duration-200 hover:scale-105"
            aria-label={menuExpanded ? "Collapse menu" : "Expand menu"}
            aria-expanded={menuExpanded}
          >
            <Icon 
              icon={menuExpanded ? "mdi:chevron-left" : "mdi:chevron-right"} 
              className="w-5 h-5 text-white" 
            />
          </button>
        </div>
      </div>
    </div>
  );
};

RoomSidebar.propTypes = {
  menuExpanded: PropTypes.bool.isRequired,
  toggleMenu: PropTypes.func.isRequired,
  activeSection: PropTypes.string.isRequired,
  setActiveSection: PropTypes.func.isRequired,
  currentRoom: PropTypes.object,
  sharedRooms: PropTypes.array.isRequired,
  handleRoomSelection: PropTypes.func.isRequired,
  refreshRooms: PropTypes.func.isRequired
};

export default RoomSidebar; 