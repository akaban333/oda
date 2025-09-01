import { render, screen } from '@testing-library/react';
import App from './App';

test('renders app with loading state', () => {
  render(<App />);
  
  // Check that the app renders with loading state
  const loadingElement = screen.getByText(/loading/i);
  expect(loadingElement).toBeInTheDocument();
  
  // Check that the app container exists
  const appContainer = screen.getByText(/loading/i).closest('.App');
  expect(appContainer).toBeInTheDocument();
});
