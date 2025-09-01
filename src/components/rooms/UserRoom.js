import React from 'react';
import { Icon } from '@iconify/react';
import PropTypes from 'prop-types';

const UserRoom = ({ 
  sessionActive, 
  sessionXP, 
  startSession, 
  endSession,
  todos,
  markTodoAsDone,
  deleteTodo,
  userMaterials,
  deleteMaterial,
  userNotes,
  deleteNote,
  handleAddNote,
  newNote,
  setNewNote,
  // Timer state from context
  isTimerRunning,
  workTime,
  setWorkTime,
  breakTime,
  setBreakTime,
  currentTime,
  isBreak,
  showTimerNotification,
  timerNotificationMessage,
  toggleTimer,
  formatTime,
  closeTimerNotification,
  // Modal state setters
  setShowTodoModal,
  setShowMaterialModal,
  // Room information and delete functionality
  currentRoom,
  currentUser,
  onDeleteRoom
}) => {
  return (
    <div className="h-full p-2" style={{ background: 'linear-gradient(135deg, #7681b3 0%, #8b96d1 50%, #a5b0e8 100%)', backgroundSize: 'cover' }}>

      {/* Room Header with Delete Button */}
      {currentRoom && (
        <div className="flex justify-between items-center mb-3 p-2 bg-white bg-opacity-20 backdrop-blur-md rounded-lg border border-white border-opacity-30 shadow-lg">
          <div className="flex items-center gap-3">
            <h2 className="text-white font-bold text-lg drop-shadow-lg">
              {currentRoom.name}
            </h2>
            <span className="text-white text-sm bg-blue-500 bg-opacity-80 px-2 py-1 rounded-full">
                              {currentRoom.participantCount || 0}/{currentRoom.maxParticipants || 4}
            </span>
          </div>
          
          {/* Delete Room Button - Only show for room owner */}
          {currentRoom.creatorId === currentUser?.uniqueId && (
            <button
              onClick={() => onDeleteRoom && onDeleteRoom(currentRoom.id || currentRoom._id)}
              className="bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded-lg shadow-lg hover:shadow-xl hover:shadow-red-500/30 transition-all duration-200 hover:scale-105 flex items-center gap-2"
              title="Delete this room"
            >
              <Icon icon="mdi:trash" className="w-4 h-4" />
              Delete Room
            </button>
          )}
        </div>
      )}

      {/* Main content grid - Original compact sizes */}
      <div className="grid grid-cols-3 gap-2 h-[calc(100vh-130px)]">
        {/* Left column */}
        <div className="flex flex-col space-y-2 h-full">
          {/* Pomodoro Timer */}
          <div className="bg-white bg-opacity-75 p-2 rounded-lg shadow-sm flex flex-col border border-[#83A0BC] h-[40%]">
            <h3 className="text-[#1D154A] font-mono uppercase text-xs mb-1 font-bold text-left">POMODORO TIMER</h3>
            
            {!isTimerRunning ? (
              // Settings view
              <div className="bg-[#FFD169] bg-opacity-80 p-2 border border-[#E5C090] mb-1 text-xs">
                <div className="mb-1">
                  <label htmlFor="work-time" className="text-[#1D154A] text-xs block mb-0.5 font-mono">Work time:</label>
                  <select 
                    id="work-time"
                    value={workTime}
                    onChange={(e) => setWorkTime(parseInt(e.target.value))}
                    className="w-full py-0.5 px-1 bg-white border border-[#E5C090] text-xs"
                    disabled={isTimerRunning}
                  >
                    <option value={25}>25 minutes</option>
                    <option value={30}>30 minutes</option>
                    <option value={35}>35 minutes</option>
                    <option value={40}>40 minutes</option>
                    <option value={45}>45 minutes</option>
                    <option value={50}>50 minutes</option>
                    <option value={55}>55 minutes</option>
                    <option value={60}>60 minutes</option>
                  </select>
                </div>
                <div>
                  <label htmlFor="break-time" className="text-[#1D154A] text-xs block mb-0.5 font-mono">Break time:</label>
                  <select 
                    id="break-time"
                    value={breakTime}
                    onChange={(e) => setBreakTime(parseInt(e.target.value))}
                    className="w-full py-0.5 px-1 bg-white border border-[#E5C090] text-xs"
                    disabled={isTimerRunning}
                  >
                    <option value={5}>5 minutes</option>
                    <option value={10}>10 minutes</option>
                    <option value={15}>15 minutes</option>
                  </select>
                </div>
              </div>
            ) : (
              // Timer view
              <div className="bg-[#FFD169] bg-opacity-80 p-2 border border-[#E5C090] mb-1 text-xs">
                <div className="text-center">
                  <div className="text-[#1D154A] font-mono text-lg font-bold mb-1">
                    {formatTime(currentTime)}
                  </div>
                  <div className="text-[#1D154A] text-xs">
                    {isBreak ? 'Break Time' : 'Work Time'}
                  </div>
                </div>
              </div>
            )}
            
            <div className="flex justify-center mt-auto mb-1">
              <button 
                className="bg-white text-[#1D154A] px-4 py-1 rounded-full font-mono uppercase text-xs hover:bg-gray-100 transition-colors shadow-sm border border-[#D9A384]"
                onClick={toggleTimer}
                aria-label={isTimerRunning ? "Stop timer" : "Start timer"}
              >
                {isTimerRunning ? "STOP TIMER" : "START TIMER"}
              </button>
            </div>
          </div>
          
          {/* Notes Section */}
          <div className="bg-white bg-opacity-75 p-2 rounded-lg shadow-sm flex flex-col border border-[#83A0BC] h-[60%]">
            <h3 className="text-[#1D154A] font-mono uppercase text-xs mb-1 font-bold text-left">NOTES</h3>
            <div className="space-y-2 overflow-auto hide-scrollbar flex-1 pr-1">
              {userNotes.slice(0, 4).map((note) => (
                <div key={note.id} className="bg-[#83A0BC] bg-opacity-80 p-2 border border-[#6A8099]">
                  <p className="text-white text-xs font-mono line-clamp-3">
                    {note.content}
                  </p>
                  <div className="flex justify-end mt-1 gap-1">
                    <button 
                      onClick={() => deleteNote(note.id)}
                      className="text-white opacity-70 hover:opacity-100"
                      aria-label="Delete note"
                    >
                      <Icon icon="mdi:trash-can-outline" className="w-4 h-4 cursor-pointer" />
                    </button>
                  </div>
                </div>
              ))}
            </div>
            <div className="mt-1 bg-white p-1 rounded-full flex items-center shadow-sm border border-[#83A0BC]">
              <input 
                type="text" 
                placeholder="Add note..." 
                value={newNote}
                onChange={(e) => setNewNote(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' && handleAddNote()}
                className="w-full border-none outline-none text-xs font-mono pl-2"
                aria-label="Add a new note"
              />
              <button 
                onClick={handleAddNote}
                className="text-[#83a0bc] hover:text-[#6A8099]"
                aria-label="Submit note"
              >
                <Icon icon="mdi:arrow-right" className="w-4 h-4 cursor-pointer mr-1" />
              </button>
            </div>
          </div>
        </div>

        {/* To-Do List (middle column) */}
        <div className="bg-white bg-opacity-75 p-2 rounded-lg shadow-sm flex flex-col border border-[#83A0BC] h-full relative">
          <div className="flex mb-1 items-center">
            <h3 className="text-[#1D154A] font-mono uppercase text-xs font-bold text-left">TO-DO LIST</h3>
            <div className="ml-auto">
              <Icon icon="mdi:view-list" className="text-[#1D154A] w-4 h-4" />
            </div>
          </div>
          <div className="flex-1 space-y-2 overflow-auto hide-scrollbar pr-1">
            {todos.slice(0, 4).map((todo) => (
              <div key={todo.id} className="bg-[#C6AE82] bg-opacity-80 p-1 flex items-center border border-[#B09C70]">
                <div className="flex items-start gap-1 flex-1">
                  <div className="w-4 h-4 flex-shrink-0 flex items-center justify-center mt-0.5">
                    <Icon icon="mdi:checkbox-outline" className="text-white w-4 h-4" />
                  </div>
                  <div className="flex-1 flex justify-between items-center">
                    <span className="text-white text-xs font-mono line-clamp-1">{todo.title}</span>
                    <div className="ml-1 flex gap-1">
                      <button
                        onClick={() => markTodoAsDone(todo.id)}
                        className="text-[10px] font-mono text-white bg-transparent whitespace-nowrap"
                        aria-label={`Mark ${todo.title} as done`}
                      >
                        Done
                      </button>
                      <button
                        onClick={() => deleteTodo(todo.id)}
                        className="text-red-400 hover:text-red-300 transition-colors"
                        aria-label={`Delete ${todo.title}`}
                      >
                        <Icon icon="mdi:trash-can-outline" className="w-3 h-3" />
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
          {/* Add button positioned at the bottom center */}
          <div className="absolute bottom-2 left-1/2 transform -translate-x-1/2">
            <button 
              className="bg-white rounded-full p-1.5 hover:bg-gray-100 transition-colors shadow-sm border border-[#83A0BC]"
              onClick={() => setShowTodoModal(true)}
              aria-label="Add new todo"
            >
              <Icon icon="mdi:plus" className="text-[#1D154A] w-4 h-4" />
            </button>
          </div>
        </div>

        {/* Materials (right column) */}
        <div className="bg-white bg-opacity-75 p-2 rounded-lg shadow-sm flex flex-col border border-[#83A0BC] h-full relative">
          <div className="flex mb-1 items-center">
            <h3 className="text-[#1D154A] font-mono uppercase text-xs font-bold text-left">MATERIALS</h3>
            <div className="ml-auto">
              <Icon icon="mdi:bookmark" className="text-[#1D154A] w-4 h-4" />
            </div>
          </div>
          
          <div className="flex-1 space-y-2 overflow-auto hide-scrollbar pr-1">
            {userMaterials.map((material) => (
              <div key={material.id} className="bg-[#8E8E9A] bg-opacity-80 p-1 flex items-start border border-[#7A7A85] h-10">
                <div className="flex items-start w-full">
                  <div className="bg-white p-1 flex items-center justify-center border border-[#7A7A85] h-9 w-9">
                    {material.type === 'pdf' && <span className="text-[#1D154A] font-bold text-xs">PDF</span>}
                    {material.type === 'figma' && (
                      <div className="flex gap-0.5">
                        <div className="w-2 h-2 bg-[#1D154A] rounded-full"></div>
                        <div className="w-2 h-2 bg-[#1D154A] rounded-full"></div>
                        <div className="w-2 h-2 bg-[#1D154A] rounded-full flex items-center justify-center">
                          <div className="w-1 h-1 bg-white rounded-full"></div>
                        </div>
                      </div>
                    )}
                    {material.type === 'word' && <span className="text-[#1D154A] font-bold text-xs">W</span>}
                    {material.type === 'link' && <span className="text-[#1D154A] font-bold text-xs">🔗</span>}
                    {material.type === 'file' && <span className="text-[#1D154A] font-bold text-xs">📁</span>}
                    {!['pdf', 'figma', 'word', 'link', 'file'].includes(material.type) && <span className="text-[#1D154A] font-bold text-xs">F</span>}
                  </div>
                  <div className="flex-1 ml-2 text-white flex items-center">
                    <button 
                      className="text-xs font-mono truncate max-w-[90px] text-left hover:text-blue-200 transition-colors"
                      onClick={() => {
                        if (material.fileUrl && material.fileUrl.startsWith('http')) {
                          // For links, open in new tab
                          window.open(material.fileUrl, '_blank');
                        } else if (material.fileUrl && material.fileUrl.startsWith('file://')) {
                          // For uploaded files, create a download link
                          const link = document.createElement('a');
                          link.href = material.fileUrl;
                          link.download = material.name;
                          link.click();
                        } else if (material.fileUrl && material.fileUrl.startsWith('blob:')) {
                          // For blob URLs (file uploads), create download
                          const link = document.createElement('a');
                          link.href = material.fileUrl;
                          link.download = material.name;
                          link.click();
                        } else {
                          alert('No valid URL for this material.');
                        }
                      }}
                      title={`Open ${material.name}`}
                    >
                      {material.name}
                    </button>
                  </div>
                  <button 
                    className="text-white opacity-70 hover:opacity-100"
                    onClick={() => deleteMaterial(material.id)}
                    aria-label={`Delete ${material.name}`}
                  >
                    <Icon icon="mdi:trash-can-outline" className="w-3.5 h-3.5" />
                  </button>
                </div>
              </div>
            ))}
          </div>
          {/* Add button positioned at the bottom center */}
          <div className="absolute bottom-2 left-1/2 transform -translate-x-1/2">
            <button 
              className="bg-white rounded-full p-1.5 hover:bg-gray-100 transition-colors shadow-sm border border-[#83A0BC]"
              onClick={() => setShowMaterialModal(true)}
              aria-label="Add new material"
            >
              <Icon icon="mdi:plus" className="text-[#1D154A] w-4 h-4" />
            </button>
          </div>
        </div>
      </div>

      {/* Timer Notification */}
      {showTimerNotification && (
        <div 
          className="fixed bottom-16 right-4 bg-[#FFD169] bg-opacity-90 p-3 rounded-lg shadow-lg border border-[#E5C090] z-50 max-w-xs"
          role="alert"
          aria-live="assertive"
        >
          <div className="flex justify-between items-center mb-1">
            <h4 className="text-[#1D154A] font-bold text-sm">Timer Complete</h4>
            <button 
              onClick={closeTimerNotification}
              aria-label="Close notification"
            >
              <Icon icon="mdi:close" className="text-[#1D154A] w-4 h-4" />
            </button>
          </div>
          <p className="text-[#1D154A] text-xs">{timerNotificationMessage}</p>
        </div>
      )}
    </div>
  );
};

UserRoom.propTypes = {
  sessionActive: PropTypes.bool.isRequired,
  sessionXP: PropTypes.number.isRequired,
  startSession: PropTypes.func.isRequired,
  endSession: PropTypes.func.isRequired,
  todos: PropTypes.array.isRequired,
  markTodoAsDone: PropTypes.func.isRequired,
  deleteTodo: PropTypes.func.isRequired,
  userMaterials: PropTypes.array.isRequired,
  deleteMaterial: PropTypes.func.isRequired,
  userNotes: PropTypes.array.isRequired,
  deleteNote: PropTypes.func.isRequired,
  handleAddNote: PropTypes.func.isRequired,
  newNote: PropTypes.string.isRequired,
  setNewNote: PropTypes.func.isRequired,
  isTimerRunning: PropTypes.bool.isRequired,
  workTime: PropTypes.number.isRequired,
  setWorkTime: PropTypes.func.isRequired,
  breakTime: PropTypes.number.isRequired,
  setBreakTime: PropTypes.func.isRequired,
  currentTime: PropTypes.number.isRequired,
  isBreak: PropTypes.bool.isRequired,
  showTimerNotification: PropTypes.bool.isRequired,
  timerNotificationMessage: PropTypes.string.isRequired,
  toggleTimer: PropTypes.func.isRequired,
  formatTime: PropTypes.func.isRequired,
  closeTimerNotification: PropTypes.func.isRequired,
  setShowTodoModal: PropTypes.func.isRequired,
  setShowMaterialModal: PropTypes.func.isRequired,
  currentRoom: PropTypes.object,
  currentUser: PropTypes.object,
  onDeleteRoom: PropTypes.func
};

export default UserRoom; 