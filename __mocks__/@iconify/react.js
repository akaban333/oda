import React from 'react';

// Mock Icon component
export const Icon = ({ icon, ...props }) => {
  return React.createElement('span', {
    'data-testid': 'icon',
    'data-icon': icon,
    ...props,
  });
};

// Mock all possible exports
export const addIcon = () => {};
export const addCollection = () => {};
export const addIconify = () => {};
export const getIcon = () => null;
export const listIcons = () => [];
export const searchIcons = () => {};

// Mock default export
export default Icon; 