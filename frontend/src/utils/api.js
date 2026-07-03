const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api';

export const apiCall = async (endpoint, options = {}) => {
  const token = localStorage.getItem('token');
  
  const headers = {
    'Content-Type': 'application/json',
    ...(token && { Authorization: `Bearer ${token}` }),
    ...options.headers,
  };

  const config = {
    ...options,
    headers,
  };

  const response = await fetch(`${API_BASE_URL}${endpoint}`, config);
  
  const data = await response.json().catch(() => null);

  if (!response.ok) {
    throw new Error(data?.error || `API request failed: ${response.status}`);
  }

  return data;
};

// Expose the raw fetch url if needed (e.g., for short URL redirects)
export const SHORT_URL_BASE = import.meta.env.VITE_APP_BASE_URL || 'http://localhost:8080/';
