import React, { createContext, useState } from 'react';

export const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const existingToken = localStorage.getItem('authToken');
  const [authToken, setAuthToken] = useState(existingToken);

  const setToken = (token) => {
    localStorage.setItem('authToken', token);
    setAuthToken(token);
  };

  const clearToken = () => {
    localStorage.removeItem('authToken');
    setAuthToken(null);
  };

  return (
    <AuthContext.Provider value={{ authToken, setToken, clearToken }}>
      {children}
    </AuthContext.Provider>
  );
};

