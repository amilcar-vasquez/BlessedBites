import { apiPost } from './client';

export type AuthUser = {
  id: number;
  email: string;
  full_name: string;
  role: string;
};

export type LoginResponse = {
  token: string;
  user: AuthUser;
};

export async function login(email: string, password: string): Promise<LoginResponse> {
  return apiPost<LoginResponse>('/auth/login', { email, password });
}

export async function signup(payload: {
  email: string;
  full_name: string;
  phone_no: string;
  password: string;
}): Promise<{ id: number }> {
  return apiPost<{ id: number }>('/auth/signup', payload);
}

export async function logout(): Promise<{ message: string }> {
  return apiPost<{ message: string }>('/auth/logout', {});
}

export async function requestPasswordReset(email: string): Promise<{ message: string }> {
  return apiPost<{ message: string }>('/auth/reset-password-request', { email });
}

export async function resetPassword(token: string, newPassword: string): Promise<{ message: string }> {
  return apiPost<{ message: string }>('/auth/reset-password', {
    token,
    new_password: newPassword
  });
}
