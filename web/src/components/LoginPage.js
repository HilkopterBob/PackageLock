import React from 'react';
import {
  LoginFooterItem,
  LoginForm,
  LoginMainFooterBandItem,
  LoginMainFooterLinksItem,
  LoginPage,
  ListItem,
  ListVariant,
  Button
} from '@patternfly/react-core';
import ExclamationCircleIcon from '@patternfly/react-icons/dist/esm/icons/exclamation-circle-icon';
import brandImg2 from '../../assets/brandImgColor2.svg'; // Adjust the import path as needed
import axios from 'axios';
import { AuthContext } from '../AuthContext';
import axiosInstance from '../axiosInstance';

export const LoginPageComponent = () => {

  const { setToken } = useContext(AuthContext);
  const [showHelperText, setShowHelperText] = React.useState(false);
  const [helperText, setHelperText] = React.useState('');
  const [username, setUsername] = React.useState('');
  const [isValidUsername, setIsValidUsername] = React.useState(true);
  const [password, setPassword] = React.useState('');
  const [isValidPassword, setIsValidPassword] = React.useState(true);
  const [isRememberMeChecked, setIsRememberMeChecked] = React.useState(false);

  const handleUsernameChange = (value) => {
    setUsername(value);
    setIsValidUsername(true);
  };

  const handlePasswordChange = (value) => {
    setPassword(value);
    setIsValidPassword(true);
  };

  const onRememberMeClick = () => {
    setIsRememberMeChecked(!isRememberMeChecked);
  };

  const onLoginButtonClick = (event) => {
    event.preventDefault();

    // Basic validation
    const isValidForm = username && password;
    setIsValidUsername(!!username);
    setIsValidPassword(!!password);
    if (!isValidForm) {
      setShowHelperText(true);
      setHelperText('Please fill out all fields.');
      return;
    }

    // Prepare the data to send
    const loginData = {
      username: username,
      password: password,
    };

    // Make the API call
    axios.post('/auth/login', loginData)
      .then(response => {
        // Handle success
        const token = response.data.token; // Adjust according to your API's response
        // Store the token in localStorage or a context
        localStorage.setItem('authToken', token);

        // Redirect to the desired page or update the app state
        window.location.href = '/'; // Adjust the redirect path as needed
      })
      .catch(error => {
        // Handle error
        console.error('Error logging in:', error);
        setShowHelperText(true);
        if (error.response && error.response.status === 401) {
          setHelperText('Invalid username or password.');
        } else {
          setHelperText('An error occurred. Please try again.');
        }
      });
    setToken(token);
    navigate('/');
  };

  const loginForm = (
    <LoginForm
      showHelperText={showHelperText}
      helperText={helperText}
      helperTextIcon={< ExclamationCircleIcon />}
      usernameLabel="Username"
      usernameValue={username}
      onChangeUsername={handleUsernameChange}
      isValidUsername={isValidUsername}
      passwordLabel="Password"
      passwordValue={password}
      onChangePassword={handlePasswordChange}
      isValidPassword={isValidPassword}
      rememberMeLabel="Keep me logged in for 30 days."
      isRememberMeChecked={isRememberMeChecked}
      onChangeRememberMe={onRememberMeClick}
      onLoginButtonClick={onLoginButtonClick}
      loginButtonLabel="Log in"
    />
  );


  // You can customize or remove the social media login content and footer items as needed.

  return (
    <LoginPage
      // You can customize these props as needed
      brandImgSrc={brandImg2}
      brandImgAlt="Your Brand Logo"
      backgroundImgSrc="/assets/images/pfbg-icon.svg"
      loginTitle="Log in to your account"
      loginSubtitle="Enter your credentials."
    >
      {loginForm}
    </LoginPage>
  );
};

export default LoginPageComponent;

