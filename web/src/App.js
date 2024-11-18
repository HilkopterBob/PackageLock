
import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import LoginPageComponent from './components/LoginPage';
import { AuthProvider } from './AuthContext';
// Import other components as needed

function App() {
  return (
    <AuthProvider>
      <Router>
        <Routes>
          <Route path="/ui/login" element={<LoginPageComponent />} />
          {/* Add other routes here */}
          <Route path="/" element={<PrivateRoute element={<HomePage />} />} />
          {/* Fallback route */}
          <Route path="*" element={<NotFoundPage />} />
        </Routes>
      </Router>
    </AuthProvider>
  );
}

export default App;

