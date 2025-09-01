import React from 'react';

// Mock Icon component
export const Icon = ({ icon, ...props }) => {
  return React.createElement('span', {
    'data-testid': 'icon',
    'data-icon': icon,
    ...props,
  });
};

// Mock default export
export default Icon; 