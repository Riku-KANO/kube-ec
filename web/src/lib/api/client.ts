import createClient from 'openapi-fetch';
import type { paths } from '@/types/api/user';

// API client configuration
const baseUrl =
  process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

// Create typed API client
export const apiClient = createClient<paths>({ baseUrl });

// Helper to set authorization token
export function setAuthToken(token: string) {
  apiClient.use({
    onRequest({ request }) {
      request.headers.set('Authorization', `Bearer ${token}`);
      return request;
    },
  });
}

// Helper to clear authorization token
export function clearAuthToken() {
  apiClient.use({
    onRequest({ request }) {
      request.headers.delete('Authorization');
      return request;
    },
  });
}
