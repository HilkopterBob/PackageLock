import React from 'react';
import { Navigate } from 'react-router-dom';
import { AuthContext } from './AuthContext';

const PrivateRoute = ({ element }) => {
  const { authToken } = React.useContext(AuthContext);

  return authToken ? element : <Navigate to="/ui/login" />;
};

export default PrivateRoute;
