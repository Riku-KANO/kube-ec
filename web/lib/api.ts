import type { Product, ProductListResponse, AuthResponse, Order } from './types';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

// Fetch wrapper with error handling
async function fetchAPI<T>(endpoint: string, options?: RequestInit): Promise<T> {
  const token = typeof window !== 'undefined' ? localStorage.getItem('access_token') : null;

  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options?.headers as Record<string, string>),
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers,
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ message: 'Request failed' }));
    throw new Error(error.message || `HTTP ${response.status}`);
  }

  return response.json();
}

// Product API
export const productApi = {
  list: async (
    page = 1,
    pageSize = 20,
    category?: string,
    search?: string
  ): Promise<ProductListResponse> => {
    const params = new URLSearchParams({
      page: page.toString(),
      page_size: pageSize.toString(),
    });
    if (category) params.append('category', category);
    if (search) params.append('search', search);

    return fetchAPI(`/products?${params}`);
  },

  get: async (id: string): Promise<Product> => {
    return fetchAPI(`/products/${id}`);
  },
};

// Auth API
export const authApi = {
  register: async (
    email: string,
    password: string,
    name: string,
    phoneNumber?: string
  ): Promise<AuthResponse> => {
    return fetchAPI('/auth/register', {
      method: 'POST',
      body: JSON.stringify({
        email,
        password,
        name,
        phone_number: phoneNumber,
      }),
    });
  },

  login: async (email: string, password: string): Promise<AuthResponse> => {
    return fetchAPI('/auth/login', {
      method: 'POST',
      body: JSON.stringify({
        email,
        password,
      }),
    });
  },
};

// Order API
export const orderApi = {
  create: async (userId: string, items: any[], shippingAddress: any): Promise<Order> => {
    return fetchAPI('/orders', {
      method: 'POST',
      body: JSON.stringify({
        user_id: userId,
        items,
        shipping_address: shippingAddress,
      }),
    });
  },

  get: async (id: string): Promise<Order> => {
    return fetchAPI(`/orders/${id}`);
  },

  list: async (userId: string, page = 1, pageSize = 20): Promise<any> => {
    const params = new URLSearchParams({
      user_id: userId,
      page: page.toString(),
      page_size: pageSize.toString(),
    });
    return fetchAPI(`/orders?${params}`);
  },
};

// Payment API
export const paymentApi = {
  create: async (orderId: string, userId: string, amount: any, method: string): Promise<any> => {
    return fetchAPI('/payments', {
      method: 'POST',
      body: JSON.stringify({
        order_id: orderId,
        user_id: userId,
        amount,
        method,
      }),
    });
  },

  process: async (paymentId: string, paymentToken: string): Promise<any> => {
    return fetchAPI('/payments/process', {
      method: 'POST',
      body: JSON.stringify({
        payment_id: paymentId,
        payment_token: paymentToken,
      }),
    });
  },
};
