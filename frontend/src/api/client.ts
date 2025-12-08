import axios from "axios";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://localhost:8080/api";

// debug: show which base URL the app built with
console.debug("API_BASE_URL:", API_BASE_URL);

const api = axios.create({
  baseURL: API_BASE_URL,
});

// If you already had interceptors for JWT, keep them:
api.interceptors.request.use((config) => {
  // AuthContext stores token under "jwt"
  const token = localStorage.getItem("jwt");
  if (token) {
    config.headers = config.headers || {};
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response interceptor to log errors for easier debugging on deployed builds
api.interceptors.response.use(
  (res) => res,
  (err) => {
    console.error("API error:", err?.response || err.message || err);
    return Promise.reject(err);
  }
);

export default api;
