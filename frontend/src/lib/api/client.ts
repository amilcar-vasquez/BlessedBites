import { getAccessToken, setAccessToken, clearSession } from '$lib/stores/auth';

export const API_BASE = import.meta.env.PUBLIC_API_BASE_URL || '/api/v1';

export class ApiError extends Error {
  status: number;
  body: string;

  constructor(status: number, body: string) {
    super(`API request failed: ${status}`);
    this.status = status;
    this.body = body;
  }
}

let refreshPromise: Promise<string | null> | null = null;

async function refreshAccessToken(): Promise<string | null> {
  if (!refreshPromise) {
    refreshPromise = (async () => {
      try {
        const res = await fetch(`${API_BASE}/auth/refresh`, {
          method: 'POST',
          credentials: 'include'
        });
        if (!res.ok) {
          clearSession();
          return null;
        }
        const payload = (await res.json()) as { token: string };
        setAccessToken(payload.token);
        return payload.token;
      } catch {
        return null;
      } finally {
        refreshPromise = null;
      }
    })();
  }
  return refreshPromise;
}

async function request<T>(
  method: string,
  path: string,
  body?: unknown,
  explicitToken?: string
): Promise<T> {
  const doFetch = (token: string) =>
    fetch(`${API_BASE}${path}`, {
      method,
      credentials: 'include',
      headers: {
        ...(body !== undefined ? { 'Content-Type': 'application/json' } : {}),
        ...(token ? { Authorization: `Bearer ${token}` } : {})
      },
      ...(body !== undefined ? { body: JSON.stringify(body) } : {})
    });

  const token = explicitToken ?? getAccessToken();
  let res = await doFetch(token);

  // One transparent refresh + retry on expired/invalid access tokens.
  if (res.status === 401 && token && !path.startsWith('/auth/')) {
    const renewed = await refreshAccessToken();
    if (renewed) {
      res = await doFetch(renewed);
    }
  }

  if (!res.ok) {
    throw new ApiError(res.status, await res.text().catch(() => ''));
  }
  return res.json() as Promise<T>;
}

export async function apiGet<T>(path: string, token?: string): Promise<T> {
  return request<T>('GET', path, undefined, token);
}

export async function apiPost<T>(path: string, body: unknown, token?: string): Promise<T> {
  return request<T>('POST', path, body, token);
}

export async function apiPut<T>(path: string, body: unknown, token?: string): Promise<T> {
  return request<T>('PUT', path, body, token);
}

export async function apiDelete<T>(path: string, token?: string): Promise<T> {
  return request<T>('DELETE', path, undefined, token);
}
