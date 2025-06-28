import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import apiConfig, { saveApiConfig, clearApiConfig, isApiConfigured } from './apiConfig';

// Mock localStorage
const localStorageMock = (() => {
  let store = {};
  return {
    getItem: vi.fn((key) => store[key] || null),
    setItem: vi.fn((key, value) => {
      store[key] = value.toString();
    }),
    removeItem: vi.fn((key) => {
      delete store[key];
    }),
    clear: vi.fn(() => {
      store = {};
    }),
  };
})();

beforeEach(() => {
  // Reset the store and mocks before each test
  localStorageMock.clear();
  // Assign the mock to global localStorage
  Object.defineProperty(window, 'localStorage', {
    value: localStorageMock,
    writable: true,
  });
  // Reset the reactive apiConfig state manually for a clean slate,
  // as it initializes from localStorage on import.
  // This requires careful handling if the module is imported once and cached.
  // A better way might be to have a reset function in apiConfig.js itself for testing.
  // For now, let's re-initialize its internal state by calling clearApiConfig
  // and then ensuring it reads from the (now mocked and empty) localStorage.
  clearApiConfig(); // Clears reactive state and localStorage via its watchers
  apiConfig.apiUrl = localStorage.getItem('typhoon_api_url') || '';
  apiConfig.apiToken = localStorage.getItem('typhoon_api_token') || '';
  apiConfig.isConfigured = !!apiConfig.apiUrl;

});

afterEach(() => {
  // Restore original localStorage if necessary, or ensure it's clean for next test module
  // For Vitest, environment is usually isolated per file, but good practice.
});

describe('apiConfig service', () => {
  it('should initialize with empty values if localStorage is empty', () => {
    expect(apiConfig.apiUrl).toBe('');
    expect(apiConfig.apiToken).toBe('');
    expect(apiConfig.isConfigured).toBe(false);
    expect(isApiConfigured()).toBe(false);
  });

  it('saveApiConfig should store URL and token in localStorage and update reactive state', () => {
    const testUrl = 'http://test.com/api';
    const testToken = 'testtoken';
    saveApiConfig(testUrl, testToken);

    expect(apiConfig.apiUrl).toBe(testUrl);
    expect(apiConfig.apiToken).toBe(testToken);
    expect(apiConfig.isConfigured).toBe(true);
    expect(isApiConfigured()).toBe(true);

    expect(localStorageMock.setItem).toHaveBeenCalledWith('typhoon_api_url', testUrl);
    expect(localStorageMock.setItem).toHaveBeenCalledWith('typhoon_api_token', testToken);
  });

  it('saveApiConfig with empty URL should update isConfigured to false', () => {
    saveApiConfig('http://test.com', 'token'); // Configure first
    expect(apiConfig.isConfigured).toBe(true);

    saveApiConfig('', 'token'); // Then save with empty URL
    expect(apiConfig.apiUrl).toBe('');
    expect(apiConfig.isConfigured).toBe(false);
    expect(isApiConfigured()).toBe(false);
    expect(localStorageMock.removeItem).toHaveBeenCalledWith('typhoon_api_url');
  });

  it('saveApiConfig with empty token should store empty token', () => {
    const testUrl = 'http://test.com/api';
    saveApiConfig(testUrl, '');

    expect(apiConfig.apiToken).toBe('');
    expect(localStorageMock.removeItem).toHaveBeenCalledWith('typhoon_api_token');
  });

  it('clearApiConfig should remove URL and token from localStorage and update state', () => {
    saveApiConfig('http://test.com/api', 'testtoken'); // Set some values first

    clearApiConfig();

    expect(apiConfig.apiUrl).toBe('');
    expect(apiConfig.apiToken).toBe('');
    expect(apiConfig.isConfigured).toBe(false);
    expect(isApiConfigured()).toBe(false);

    expect(localStorageMock.removeItem).toHaveBeenCalledWith('typhoon_api_url');
    expect(localStorageMock.removeItem).toHaveBeenCalledWith('typhoon_api_token');
  });

  it('isApiConfigured should return true only if apiUrl is set', () => {
    saveApiConfig('', 'sometoken');
    expect(isApiConfigured()).toBe(false);

    saveApiConfig('http://example.com', '');
    expect(isApiConfigured()).toBe(true);
  });
});
