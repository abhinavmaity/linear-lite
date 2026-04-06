import { getStoredToken, useAuthStore } from 'store/authStore';
import { ApiError, CollectionResponse, ErrorEnvelope, SingleResponse } from 'types/api';
import { cleanParams } from 'utils/query';
import { shouldUseMockForPath } from './env';
import { mockRequest } from './mockBackend';

const API_BASE = import.meta.env.VITE_API_BASE_URL ?? '/api/v1';

async function request<T>(path: string, init?: RequestInit) {
  const token = getStoredToken();

  if (shouldUseMockForPath(path)) {
    return mockRequest<T>(path, init, token);
  }

  const response = await fetch(`${API_BASE}${path}`, {
    ...init,
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...init?.headers,
    },
  });

  if (!response.ok) {
    const body = (await response.json().catch(() => null)) as ErrorEnvelope | null;
    if (response.status === 401) {
      useAuthStore.getState().clearSession();
    }
    throw new ApiError(
      response.status,
      body?.error ?? {
        code: 'internal_error',
        message: 'Request failed.',
      },
    );
  }

  if (response.status === 204) {
    return undefined as T;
  }

  return (await response.json()) as T;
}

export const apiClient = {
  get<T>(path: string, params?: Record<string, unknown> | URLSearchParams) {
    const query = params ? (params instanceof URLSearchParams ? params.toString() : cleanParams(params).toString()) : '';
    return request<T>(query ? `${path}?${query}` : path);
  },
  post<T>(path: string, body?: unknown) {
    return request<T>(path, {
      method: 'POST',
      body: body ? JSON.stringify(body) : undefined,
    });
  },
  put<T>(path: string, body?: unknown) {
    return request<T>(path, {
      method: 'PUT',
      body: body ? JSON.stringify(body) : undefined,
    });
  },
  delete<T>(path: string) {
    return request<T>(path, { method: 'DELETE' });
  },
};

export type Collection<T> = CollectionResponse<T>;
export type Resource<T> = SingleResponse<T>;
