import React, { useState } from 'react';
import AboutPanel from './AboutPanel';
import PurposePanel from './PurposePanel';
import FAQPanel from './FAQPanel';
import ContactPanel from './ContactPanel';

const Navbar = ({ onLoginClick, onSignupClick }) => {
  const [isAboutOpen, setIsAboutOpen] = useState(false);
  const [isPurposeOpen, setIsPurposeOpen] = useState(false);
  const [isFAQOpen, setIsFAQOpen] = useState(false);
  const [isContactOpen, setIsContactOpen] = useState(false);

  return (
    <>
      <nav className="fixed top-0 left-0 right-0 bg-[#1a1b3b]/80 backdrop-blur-sm z-50">
        <div className="max-w-7xl mx-auto px-8">
          <div className="flex justify-between items-center h-14 border-b border-[#ffe993]/20">
            {/* Left side menu items */}
            <div className="flex space-x-12">
              <button 
                onClick={() => setIsAboutOpen(true)}
                className="text-[#ffe993] hover:text-[#fff] px-3 py-1 rounded-md text-sm font-bold"
              >
                ABOUT US
              </button>
              <button 
                onClick={() => setIsPurposeOpen(true)}
                className="text-[#ffe993] hover:text-[#fff] px-3 py-1 rounded-md text-sm font-bold"
              >
                PURPOSE
              </button>
              <button 
                onClick={() => setIsFAQOpen(true)}
                className="text-[#ffe993] hover:text-[#fff] px-3 py-1 rounded-md text-sm font-bold"
              >
                FAQ
              </button>
              <button 
                onClick={() => setIsContactOpen(true)}
                className="text-[#ffe993] hover:text-[#fff] px-3 py-1 rounded-md text-sm font-bold"
              >
                CONTACT US
              </button>
            </div>
            
            {/* Auth buttons */}
            <div className="flex space-x-4">
              <button 
                onClick={onSignupClick}
                className="bg-[#ffe993] text-[#1a1b3b] px-6 py-1.5 rounded-md text-sm font-bold hover:bg-[#e8c56b] transition-colors duration-200"
              >
                SIGN UP
              </button>
              <button 
                onClick={onLoginClick}
                className="bg-[#1a1b3b] text-[#ffe993] px-6 py-1.5 rounded-md text-sm font-bold border-2 border-[#ffe993] hover:bg-[#ffe993] hover:text-[#1a1b3b] transition-colors duration-200"
              >
                LOG IN
              </button>
            </div>
          </div>
        </div>
      </nav>

      {/* Panels */}
      <AboutPanel isOpen={isAboutOpen} onClose={() => setIsAboutOpen(false)} />
      <PurposePanel isOpen={isPurposeOpen} onClose={() => setIsPurposeOpen(false)} />
      <FAQPanel isOpen={isFAQOpen} onClose={() => setIsFAQOpen(false)} />
      <ContactPanel isOpen={isContactOpen} onClose={() => setIsContactOpen(false)} />
    </>
  );
};

export default Navbar; 