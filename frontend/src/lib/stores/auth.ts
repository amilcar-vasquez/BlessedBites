import { browser } from '$app/environment';
import { writable, derived } from 'svelte/store';

export type Role = 'guest' | 'user' | 'admin';

export type AuthUser = {
  id: number;
  email: string;
  full_name: string;
  phone_no?: string;
  role: string;
};

type AuthState = {
  token: string;
  user: AuthUser | null;
};

const TOKEN_KEY = 'bb_access_token';
const USER_KEY = 'bb_user';

function readState(): AuthState {
  if (!browser) return { token: '', user: null };
  const token = localStorage.getItem(TOKEN_KEY) || '';
  const raw = localStorage.getItem(USER_KEY);
  let user: AuthUser | null = null;
  if (raw) {
    try {
      user = JSON.parse(raw) as AuthUser;
    } catch {
      user = null;
    }
  }
  return { token, user };
}

export const auth = writable<AuthState>(readState());

export const role = derived(auth, ($auth): Role => {
  if (!$auth.user) return 'guest';
  return ($auth.user.role || '').toLowerCase() === 'admin' ? 'admin' : 'user';
});

export function setSession(token: string, user: AuthUser) {
  if (browser) {
    localStorage.setItem(TOKEN_KEY, token);
    localStorage.setItem(USER_KEY, JSON.stringify(user));
  }
  auth.set({ token, user });
}

export function setAccessToken(token: string) {
  if (browser) {
    localStorage.setItem(TOKEN_KEY, token);
  }
  auth.update((s) => ({ ...s, token }));
}

export function clearSession() {
  if (browser) {
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(USER_KEY);
  }
  auth.set({ token: '', user: null });
}

export function getAccessToken(): string {
  if (!browser) return '';
  return localStorage.getItem(TOKEN_KEY) || '';
}

if (browser) {
  // Keep multiple tabs in sync.
  window.addEventListener('storage', () => auth.set(readState()));
}
