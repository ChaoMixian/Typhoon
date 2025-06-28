import { reactive, watch } from 'vue';

const STORAGE_KEY_API_URL = 'typhoon_api_url';
const STORAGE_KEY_API_TOKEN = 'typhoon_api_token';

const apiConfig = reactive({
  apiUrl: localStorage.getItem(STORAGE_KEY_API_URL) || '',
  apiToken: localStorage.getItem(STORAGE_KEY_API_TOKEN) || '',
  isConfigured: !!localStorage.getItem(STORAGE_KEY_API_URL), // True if apiUrl is set
});

watch(() => apiConfig.apiUrl, (newUrl) => {
  if (newUrl) {
    localStorage.setItem(STORAGE_KEY_API_URL, newUrl);
    apiConfig.isConfigured = true;
  } else {
    localStorage.removeItem(STORAGE_KEY_API_URL);
    apiConfig.isConfigured = false;
  }
});

watch(() => apiConfig.apiToken, (newToken) => {
  if (newToken) {
    localStorage.setItem(STORAGE_KEY_API_TOKEN, newToken);
  } else {
    localStorage.removeItem(STORAGE_KEY_API_TOKEN);
  }
});

export function saveApiConfig(url, token) {
  apiConfig.apiUrl = url;
  apiConfig.apiToken = token; // Token can be empty if not used
}

export function clearApiConfig() {
  apiConfig.apiUrl = '';
  apiConfig.apiToken = '';
}

// Function to check if the API is configured, primarily for router guards
export function isApiConfigured() {
  return apiConfig.isConfigured && !!apiConfig.apiUrl;
}

export default apiConfig;
