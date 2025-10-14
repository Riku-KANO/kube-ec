export interface Product {
  id: string;
  name: string;
  description: string;
  price: Money;
  stock_quantity: number;
  category: string;
  image_urls: string[];
  sku: string;
  is_active: boolean;
}

export interface Money {
  currency: string;
  amount: number;
}

export interface User {
  id: string;
  email: string;
  name: string;
  phone_number: string;
}

export interface Order {
  id: string;
  user_id: string;
  items: OrderItem[];
  total_amount: Money;
  status: string;
  shipping_address?: Address;
  payment_id: string;
}

export interface OrderItem {
  product_id: string;
  product_name: string;
  quantity: number;
  unit_price: Money;
  subtotal: Money;
}

export interface Address {
  postal_code: string;
  prefecture: string;
  city: string;
  address_line1: string;
  address_line2: string;
  phone_number: string;
}

export interface AuthResponse {
  user: User;
  access_token: string;
  refresh_token: string;
}

export interface PaginationResponse {
  total_count: number;
  total_pages: number;
  current_page: number;
}

export interface ProductListResponse {
  products: Product[];
  pagination: PaginationResponse;
}
