            {/* Friends */}
            <div 
              className="w-[32%] h-32 bg-transparent border border-[#edc7cf] rounded-xl relative cursor-pointer hover:bg-[#2a2a42] transition-colors"
              onClick={() => setIsFriendsOpen(true)}
              style={{ padding: '20px' }} // Increased padding for larger clickable area
            >
              <Icon icon="fa-solid:user-friends" className="w-8 h-8 text-[#edc7cf] absolute top-4 left-5" />
              <div className="h-full flex items-center justify-center">
                <p className="text-white font-mono text-lg">{mockFriends} FRIENDS</p>
              </div>
            </div>

            {/* XP */}
            <div 
              className="w-[32%] h-32 bg-transparent border border-[#edc7cf] rounded-xl relative cursor-pointer hover:bg-[#2a2a42] transition-colors"
              onClick={() => setIsXPOpen(true)}
              style={{ padding: '20px' }} // Increased padding for larger clickable area
            >
              <Icon icon="game-icons:fluffy-flame" className="w-8 h-8 text-[#edc7cf] absolute top-4 left-5" />
              <div className="h-full flex items-center justify-center">
                <p className="text-white font-mono text-lg">2300 TOTAL XP</p>
              </div>
            </div>

            {/* Sessions */}
            <div 
              className="w-[32%] h-32 bg-transparent border border-[#edc7cf] rounded-xl relative cursor-pointer hover:bg-[#2a2a42] transition-colors"
              onClick={() => setIsProgressChartOpen(true)}
              style={{ padding: '20px' }} // Increased padding for larger clickable area
            >
              <Icon icon="material-symbols:book-outline-rounded" className="w-8 h-8 text-[#edc7cf] absolute top-4 left-5" />
              <div className="h-full flex items-center justify-center">
                <p className="text-white font-mono text-lg">36 SESSIONS</p>
              </div>
            </div>

            {/* Tabs Section */}
            <div className="flex justify-around mt-4">
              <div 
                className="cursor-pointer p-4 hover:bg-[#2a2a42] transition-colors" // Increased padding for larger clickable area
                onClick={() => setActiveTab('socialFeed')}
              >
                <p className="text-white">Social Feed</p>
              </div>
              <div 
                className="cursor-pointer p-4 hover:bg-[#2a2a42] transition-colors" // Increased padding for larger clickable area
                onClick={() => setActiveTab('rooms')}
              >
                <p className="text-white">Rooms</p>
              </div>
              <div 
                className="cursor-pointer p-4 hover:bg-[#2a2a42] transition-colors" // Increased padding for larger clickable area
                onClick={() => setActiveTab('profile')}
              >
                <p className="text-white">Profile</p>
              </div>
            </div> 