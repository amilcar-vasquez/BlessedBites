import { apiDelete, apiGet, apiPost, apiPut } from './client';
import type { MenuItem } from './menu';
import type { Category } from './categories';

export type AdminOrder = {
  id: number;
  user_id: number;
  full_name: string;
  total_cost: number;
  status: string;
  created_at: string;
};

export type AdminUser = {
  id: number;
  email: string;
  full_name: string;
  phone_no: string;
  role: 'admin' | 'customer' | string;
  created_at: string;
};

type ListResponse<T> = {
  items: T[];
};

function getAccessToken(): string {
  if (typeof window === 'undefined') return '';
  return localStorage.getItem('bb_access_token') || '';
}

export async function listAdminMenu(): Promise<MenuItem[]> {
  const token = getAccessToken();
  const payload = await apiGet<ListResponse<MenuItem>>('/admin/menu', token);
  return payload.items;
}

export async function createAdminMenu(input: {
  name: string;
  description: string;
  price: number;
  category_id: number;
  image_url: string;
  is_active?: boolean;
}): Promise<MenuItem> {
  return apiPost<MenuItem>('/admin/menu', input, getAccessToken());
}

export async function updateAdminMenu(
  id: number,
  input: {
    name: string;
    description: string;
    price: number;
    category_id: number;
    image_url: string;
    is_active?: boolean;
  }
): Promise<MenuItem> {
  return apiPut<MenuItem>(`/admin/menu/${id}`, input, getAccessToken());
}

export async function deleteAdminMenu(id: number): Promise<{ deleted: number }> {
  return apiDelete<{ deleted: number }>(`/admin/menu/${id}`, getAccessToken());
}

export async function createAdminCategory(name: string): Promise<Category> {
  return apiPost<Category>('/admin/category', { name }, getAccessToken());
}

export async function deleteAdminCategory(id: number): Promise<{ deleted: number }> {
  return apiDelete<{ deleted: number }>(`/admin/category/${id}`, getAccessToken());
}

export async function listAdminOrders(): Promise<AdminOrder[]> {
  const token = getAccessToken();
  const payload = await apiGet<ListResponse<AdminOrder>>('/admin/orders', token);
  return payload.items;
}

export async function listAdminUsers(): Promise<AdminUser[]> {
  const token = getAccessToken();
  const payload = await apiGet<ListResponse<AdminUser>>('/admin/users', token);
  return payload.items;
}

export async function updateAdminUser(
  id: number,
  input: Partial<Pick<AdminUser, 'email' | 'full_name' | 'phone_no' | 'role'>>
): Promise<AdminUser> {
  return apiPut<AdminUser>(`/admin/users/${id}`, input, getAccessToken());
}

export async function deleteAdminUser(id: number): Promise<{ deleted: number }> {
  return apiDelete<{ deleted: number }>(`/admin/users/${id}`, getAccessToken());
}
