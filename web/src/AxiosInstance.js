import axios from 'axios';

const axiosInstance = axios.create({
  baseURL: '/', // Adjust the base URL as needed
});

axiosInstance.interceptors.request.use(
  config => {
    const token = localStorage.getItem('authToken'); // Or use context
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`;
    }
    return config;
  },
  error => Promise.reject(error)
);

export default axiosInstance;

