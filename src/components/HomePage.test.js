import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import HomePage from './HomePage';

// Mock the child components
jest.mock('./ProfileSection', () => {
  return function MockProfileSection({ user, refreshUserProfile, onLogout }) {
    return (
      <div data-testid="profile-section">
        <h2>Profile Section</h2>
        <p>User: {user?.username || 'No user'}</p>
        <button onClick={refreshUserProfile}>Refresh Profile</button>
        <button onClick={onLogout}>Logout</button>
      </div>
    );
  };
});

jest.mock('./SocialFeed', () => {
  return function MockSocialFeed() {
    return (
      <div data-testid="social-feed">
        <h2>Social Feed</h2>
        <p>Social feed content goes here</p>
      </div>
    );
  };
});

jest.mock('./RoomsInterface', () => {
  return function MockRoomsInterface() {
    return (
      <div data-testid="rooms-interface">
        <h2>Rooms Interface</h2>
        <p>Rooms interface content goes here</p>
      </div>
    );
  };
});

// Mock the Icon component
jest.mock('@iconify/react', () => ({
  Icon: ({ icon, className }) => (
    <span data-testid={`icon-${icon}`} className={className}>
      {icon}
    </span>
  ),
}));

describe('HomePage', () => {
  const mockUser = {
    username: 'testuser',
    email: 'test@example.com',
    uniqueID: '12345',
    xp: 100,
  };

  const mockRefreshUserProfile = jest.fn();
  const mockOnLogout = jest.fn();

  beforeEach(() => {
    jest.clearAllMocks();
  });

  test('renders without crashing', () => {
    render(
      <HomePage
        user={mockUser}
        refreshUserProfile={mockRefreshUserProfile}
        onLogout={mockOnLogout}
      />
    );
    
    expect(screen.getByText('LOGO')).toBeInTheDocument();
    expect(screen.getByText('LOGOUT')).toBeInTheDocument();
  });

  test('displays user information in profile section', () => {
    render(
      <HomePage
        user={mockUser}
        refreshUserProfile={mockRefreshUserProfile}
        onLogout={mockOnLogout}
      />
    );
    
    expect(screen.getByText('User: testuser')).toBeInTheDocument();
  });

  test('shows profile section by default', () => {
    render(
      <HomePage
        user={mockUser}
        refreshUserProfile={mockRefreshUserProfile}
        onLogout={mockOnLogout}
      />
    );
    
    expect(screen.getByTestId('profile-section')).toBeInTheDocument();
    expect(screen.queryByTestId('social-feed')).not.toBeInTheDocument();
    expect(screen.queryByTestId('rooms-interface')).not.toBeInTheDocument();
  });

  test('switches to social feed tab when clicked', async () => {
    render(
      <HomePage
        user={mockUser}
        refreshUserProfile={mockRefreshUserProfile}
        onLogout={mockOnLogout}
      />
    );
    
    const socialTab = screen.getByTestId('icon-material-symbols:group').closest('button');
    fireEvent.click(socialTab);
    
    await waitFor(() => {
      expect(screen.getByTestId('social-feed')).toBeInTheDocument();
      expect(screen.queryByTestId('profile-section')).not.toBeInTheDocument();
    });
  });

  test('switches to rooms tab when clicked', async () => {
    render(
      <HomePage
        user={mockUser}
        refreshUserProfile={mockRefreshUserProfile}
        onLogout={mockOnLogout}
      />
    );
    
    const roomsTab = screen.getByTestId('icon-material-symbols:door-open').closest('button');
    fireEvent.click(roomsTab);
    
    await waitFor(() => {
      expect(screen.getByTestId('rooms-interface')).toBeInTheDocument();
      expect(screen.queryByTestId('profile-section')).not.toBeInTheDocument();
    });
  });

  test('switches back to profile tab when clicked', async () => {
    render(
      <HomePage
        user={mockUser}
        refreshUserProfile={mockRefreshUserProfile}
        onLogout={mockOnLogout}
      />
    );
    
    // First switch to social tab
    const socialTab = screen.getByTestId('icon-material-symbols:group').closest('button');
    fireEvent.click(socialTab);
    
    await waitFor(() => {
      expect(screen.getByTestId('social-feed')).toBeInTheDocument();
    });
    
    // Then switch back to profile tab
    const profileTab = screen.getByTestId('icon-material-symbols:person').closest('button');
    fireEvent.click(profileTab);
    
    await waitFor(() => {
      expect(screen.getByTestId('profile-section')).toBeInTheDocument();
      expect(screen.queryByTestId('social-feed')).not.toBeInTheDocument();
    });
  });

  test('calls onLogout when logout button is clicked', () => {
    render(
      <HomePage
        user={mockUser}
        refreshUserProfile={mockRefreshUserProfile}
        onLogout={mockOnLogout}
      />
    );
    
    const logoutButton = screen.getByText('LOGOUT');
    fireEvent.click(logoutButton);
    
    expect(mockOnLogout).toHaveBeenCalledTimes(1);
  });

  test('calls refreshUserProfile when refresh button is clicked in profile section', () => {
    render(
      <HomePage
        user={mockUser}
        refreshUserProfile={mockRefreshUserProfile}
        onLogout={mockOnLogout}
      />
    );
    
    const refreshButton = screen.getByText('Refresh Profile');
    fireEvent.click(refreshButton);
    
    expect(mockRefreshUserProfile).toHaveBeenCalledTimes(1);
  });

  test('applies active tab styling correctly', async () => {
    render(
      <HomePage
        user={mockUser}
        refreshUserProfile={mockRefreshUserProfile}
        onLogout={mockOnLogout}
      />
    );
    
    // Profile tab should be active by default
    const profileTab = screen.getByTestId('icon-material-symbols:person').closest('button');
    expect(profileTab).toHaveClass('scale-105', 'shadow-lg');
    
    // Social tab should not be active
    const socialTab = screen.getByTestId('icon-material-symbols:group').closest('button');
    expect(socialTab).not.toHaveClass('scale-105', 'shadow-lg');
  });

  test('handles missing user gracefully', () => {
    render(
      <HomePage
        user={null}
        refreshUserProfile={mockRefreshUserProfile}
        onLogout={mockOnLogout}
      />
    );
    
    expect(screen.getByText('User: No user')).toBeInTheDocument();
  });

  test('maintains navigation structure', () => {
    render(
      <HomePage
        user={mockUser}
        refreshUserProfile={mockRefreshUserProfile}
        onLogout={mockOnLogout}
      />
    );
    
    // Check that all navigation tabs are present
    expect(screen.getByTestId('icon-material-symbols:group')).toBeInTheDocument();
    expect(screen.getByTestId('icon-material-symbols:door-open')).toBeInTheDocument();
    expect(screen.getByTestId('icon-material-symbols:person')).toBeInTheDocument();
    
    // Check that logo and logout are always visible
    expect(screen.getByText('LOGO')).toBeInTheDocument();
    expect(screen.getByText('LOGOUT')).toBeInTheDocument();
  });
}); 