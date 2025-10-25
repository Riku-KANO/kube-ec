import { apiClient, setAuthToken } from './client';
import type { components } from '@/types/api/user';

// Type aliases for better readability
export type User = components['schemas']['User'];
export type RegisterRequest = components['schemas']['RegisterRequest'];
export type RegisterResponse = components['schemas']['RegisterResponse'];
export type LoginRequest = components['schemas']['LoginRequest'];
export type LoginResponse = components['schemas']['LoginResponse'];
export type UpdateUserRequest = components['schemas']['UpdateUserRequest'];

/**
 * Register a new user
 */
export async function register(data: RegisterRequest) {
  const { data: response, error } = await apiClient.POST('/auth/register', {
    body: data,
  });

  if (error) {
    throw new Error(error.error || 'Registration failed');
  }

  // Automatically set auth token after successful registration
  if (response?.access_token) {
    setAuthToken(response.access_token);
  }

  return response;
}

/**
 * Login user
 */
export async function login(data: LoginRequest) {
  const { data: response, error } = await apiClient.POST('/auth/login', {
    body: data,
  });

  if (error) {
    throw new Error(error.error || 'Login failed');
  }

  // Automatically set auth token after successful login
  if (response?.access_token) {
    setAuthToken(response.access_token);
  }

  return response;
}

/**
 * Get user by ID
 */
export async function getUser(userId: string) {
  const { data, error } = await apiClient.GET('/users/{id}', {
    params: {
      path: { id: userId },
    },
  });

  if (error) {
    throw new Error(error.error || 'Failed to get user');
  }

  return data;
}

/**
 * Update user
 */
export async function updateUser(userId: string, data: UpdateUserRequest) {
  const { data: response, error } = await apiClient.PUT('/users/{id}', {
    params: {
      path: { id: userId },
    },
    body: data,
  });

  if (error) {
    throw new Error(error.error || 'Failed to update user');
  }

  return response;
}

/**
 * Delete user
 */
export async function deleteUser(userId: string) {
  const { data, error } = await apiClient.DELETE('/users/{id}', {
    params: {
      path: { id: userId },
    },
  });

  if (error) {
    throw new Error(error.error || 'Failed to delete user');
  }

  return data;
}
