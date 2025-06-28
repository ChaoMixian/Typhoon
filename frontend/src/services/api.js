import axios from "axios";
import apiConfig from "./apiConfig"; // Import the apiConfig store
import router from '@/router'; // To redirect if API call fails due to auth/config

// Create a function to get a new Axios instance, configured with current API settings
const createApiClient = () => {
  if (!apiConfig.isConfigured || !apiConfig.apiUrl) {
    // This case should ideally be caught by router guards before an API call is made
    // But as a safeguard, if somehow an API call is attempted without config:
    console.error("API client created without API URL configured.");
    // Return a dummy client that will always fail or throw an error.
    // Or, better, ensure this state leads to a redirect via router.
    // For now, let it proceed and fail, relying on router guards.
  }

  const client = axios.create({
    baseURL: apiConfig.apiUrl, // Dynamically set baseURL
    timeout: 10000, // Increased timeout for potentially slower connections
  });

  // Request interceptor
  client.interceptors.request.use(
    (config) => {
      if (apiConfig.apiToken) {
        config.headers.Authorization = `Bearer ${apiConfig.apiToken}`;
      }
      // Ensure headers are defined
      config.headers = config.headers || {};
      // Set Content-Type for relevant methods if not already set
      if (['POST', 'PUT', 'PATCH'].includes(config.method?.toUpperCase())) {
        config.headers['Content-Type'] = config.headers['Content-Type'] || 'application/json';
      }
      return config;
    },
    (error) => {
      return Promise.reject(error);
    }
  );

  // Response interceptor
  client.interceptors.response.use(
    (response) => {
      return response.data; // Return only data part of the response
    },
    (error) => {
      console.error("API Request Error:", error.response?.data || error.message);
      if (error.response) {
        // Handle specific HTTP error codes
        if (error.response.status === 401 || error.response.status === 403) {
          // Unauthorized or Forbidden
          // This could mean token is invalid or API URL is wrong.
          // Clear potentially bad config and redirect to setup.
          // apiConfig.apiUrl = ''; // Let user re-verify on setup page
          // apiConfig.apiToken = ''; // Clear token
          // apiConfig.isConfigured = false; // Trigger reactivity
          if (router.currentRoute.value.path !== '/api-setup') {
            router.push('/api-setup');
          }
          return Promise.reject(new Error("Authentication failed. Please re-configure API access."));
        }
      } else if (error.request) {
        // Network error (e.g., server down, CORS if not configured on backend)
        return Promise.reject(new Error("Network error or server unreachable."));
      }
      return Promise.reject(error); // For other errors
    }
  );
  return client;
};

// Export a function that returns a new instance of the configured client.
// This ensures that if apiConfig changes (e.g., user updates it),
// new API calls will use the new configuration.
// However, for simplicity in usage, many apps export a single instance
// and manage its reconfiguration. Let's try a single, reconfigurable instance approach first,
// but be mindful of its limitations if baseURL needs to change without app reload.

// The challenge with a single instance is that baseURL is set at creation.
// If apiConfig.apiUrl changes, the existing apiClient instance won't pick it up.
// So, we must either:
// 1. Re-create the apiClient instance when apiUrl changes (complex to manage globally).
// 2. Use a function that always returns a fresh client (as above with createApiClient).
// 3. Don't set baseURL on create, but prepend apiConfig.apiUrl to every request url (in interceptor or in each call).

// Option 3 is often the most straightforward for dynamic base URLs with a shared interceptor logic.
// Let's try a single instance that doesn't have baseURL initially, and the request interceptor adds it.

const apiClient = axios.create({
  timeout: 10000,
});

apiClient.interceptors.request.use(
  (config) => {
    if (!apiConfig.isConfigured || !apiConfig.apiUrl) {
      // This should ideally be caught by router guards.
      // If not, reject the request to prevent it from going out with a bad URL.
      const err = new Error("API URL not configured. Please set up API access.");
      console.error(err.message);
      // Redirect to setup page if not already there
      if (router.currentRoute.value.path !== '/api-setup') {
         router.push('/api-setup');
      }
      return Promise.reject(err);
    }

    // Prepend baseURL to the request URL if it's relative
    if (config.url && !config.url.startsWith('http://') && !config.url.startsWith('https://')) {
        config.url = `${apiConfig.apiUrl.replace(/\/+$/, '')}/${config.url.replace(/^\/+/, '')}`;
    }


    if (apiConfig.apiToken) {
      config.headers.Authorization = `Bearer ${apiConfig.apiToken}`;
    }
    config.headers = config.headers || {};
    if (['POST', 'PUT', 'PATCH'].includes(config.method?.toUpperCase())) {
      config.headers['Content-Type'] = config.headers['Content-Type'] || 'application/json';
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

apiClient.interceptors.response.use(
  (response) => {
    return response.data;
  },
  (error) => {
    console.error("API Request Error:", error.response?.status, error.response?.data || error.message);
    if (error.response && (error.response.status === 401 || error.response.status === 403)) {
      if (router.currentRoute.value.path !== '/api-setup') {
        // Optionally, preserve the current path to redirect back after setup
        // const currentPath = router.currentRoute.value.fullPath;
        // router.push({ path: '/api-setup', query: { redirect: currentPath } });
        router.push('/api-setup');
      }
      return Promise.reject(new Error("Authentication or Permission Error. Please re-configure API access."));
    } else if (!error.response && error.request) {
       // Network error
        if (router.currentRoute.value.path !== '/api-setup') {
            // Could be that the initial URL saved was wrong.
            // router.push('/api-setup'); // Or just let the error propagate to the caller to display.
        }
        return Promise.reject(new Error(`Network error or server at ${apiConfig.apiUrl} unreachable.`));
    }
    return Promise.reject(error);
  }
);


export async function checkBackendConnection() {
  if (!apiConfig.isConfigured || !apiConfig.apiUrl) {
    return Promise.reject(new Error("API URL not configured."));
  }
  try {
    // Assuming Mihomo's /version endpoint is proxied by Typhoon at /mihomo/version
    // And it doesn't require specific Mihomo auth if accessed via Typhoon proxy
    // Adjust endpoint if necessary. This is just a lightweight check.
    // If Typhoon itself has a health/version endpoint, that would be even better.
    // For now, let's assume such an endpoint exists on Typhoon backend (e.g. /system/health or /version)
    // As per Mihomo docs, `/version` is a simple GET.
    // If Typhoon proxies Mihomo API under a path, e.g. /mihomoapi/, then it would be /mihomoapi/version
    // I need to make an assumption here or this function cannot be fully implemented.
    // Let's assume Typhoon has a /ping or /health endpoint. If not, this will fail.
    // For the purpose of this task, I'll assume a '/system/ping' endpoint on Typhoon backend.
    // If the backend API is directly the Mihomo API, then '/version' is fine.
    // Given the context, Typhoon is a management tool, so it likely has its own API,
    // and Mihomo API is either proxied or called by Typhoon backend.
    // The API URL the user enters is for Typhoon.

    // Let's try to get the Typhoon config as a health check.
    // This means the user-provided token (if any) should grant access to this.
    await apiClient.get('config'); // Assuming 'config' is a valid GET endpoint on Typhoon API
    return Promise.resolve();
  } catch (error) {
    console.error("Backend connection check failed:", error);
    // The error interceptor in apiClient might have already processed this.
    // Rethrow a more specific message if needed, or rely on the interceptor's error.
    if (error.message && error.message.includes("Authentication or Permission Error")) {
        throw error; // Propagate the auth error from interceptor
    }
    throw new Error(`Connection to ${apiConfig.apiUrl} failed. ${error.message || ''}`);
  }
}

export default apiClient;
