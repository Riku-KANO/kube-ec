'use client';

import { useState } from 'react';
import { register, login, getUser, updateUser } from '@/lib/api/user';
import type { RegisterRequest, LoginRequest } from '@/lib/api/user';

export default function ExamplePage() {
  const [userId, setUserId] = useState<string>('');

  // Register example
  const handleRegister = async () => {
    try {
      const data: RegisterRequest = {
        email: 'test@example.com',
        password: 'password123',
        name: 'Test User',
        phone_number: '+1234567890',
      };

      const response = await register(data);
      console.log('Registered:', response);

      if (response?.user?.id) {
        setUserId(response.user.id);
      }
    } catch (error) {
      console.error('Registration error:', error);
    }
  };

  // Login example
  const handleLogin = async () => {
    try {
      const data: LoginRequest = {
        email: 'test@example.com',
        password: 'password123',
      };

      const response = await login(data);
      console.log('Logged in:', response);

      if (response?.user?.id) {
        setUserId(response.user.id);
      }
    } catch (error) {
      console.error('Login error:', error);
    }
  };

  // Get user example
  const handleGetUser = async () => {
    if (!userId) {
      console.error('No user ID');
      return;
    }

    try {
      const user = await getUser(userId);
      console.log('User:', user);
    } catch (error) {
      console.error('Get user error:', error);
    }
  };

  // Update user example
  const handleUpdateUser = async () => {
    if (!userId) {
      console.error('No user ID');
      return;
    }

    try {
      const user = await updateUser(userId, {
        name: 'Updated Name',
        phone_number: '+9876543210',
      });
      console.log('Updated user:', user);
    } catch (error) {
      console.error('Update user error:', error);
    }
  };

  return (
    <div className="p-8 space-y-4">
      <h1 className="text-2xl font-bold">API Client Example</h1>

      <div className="space-y-2">
        <button
          onClick={handleRegister}
          className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
        >
          Register
        </button>

        <button
          onClick={handleLogin}
          className="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600 ml-2"
        >
          Login
        </button>

        <button
          onClick={handleGetUser}
          className="px-4 py-2 bg-purple-500 text-white rounded hover:bg-purple-600 ml-2"
          disabled={!userId}
        >
          Get User
        </button>

        <button
          onClick={handleUpdateUser}
          className="px-4 py-2 bg-orange-500 text-white rounded hover:bg-orange-600 ml-2"
          disabled={!userId}
        >
          Update User
        </button>
      </div>

      {userId && (
        <div className="mt-4 p-4 bg-gray-100 rounded">
          <p className="text-sm">Current User ID: {userId}</p>
        </div>
      )}

      <div className="mt-4 p-4 bg-gray-100 rounded">
        <p className="text-sm text-gray-600">
          Open the browser console to see API responses
        </p>
      </div>
    </div>
  );
}
