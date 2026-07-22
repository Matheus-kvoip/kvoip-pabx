import type {
  AuthUser,
  CallRecord,
  CreateExtensionInput,
  DashboardStats,
  Extension,
  HealthStatus,
  LoginInput,
  LoginResponse,
  Trunk,
} from '@kvoip/shared';
import { getServerToken } from './auth-server';
import { apiRequest } from './http';

function request<T>(path: string, init?: RequestInit & { token?: string }) {
  return apiRequest<T>(path, init, getServerToken);
}

export const api = {
  health: () => request<HealthStatus>('/health'),
  auth: {
    login: (body: LoginInput) =>
      request<LoginResponse>('/auth/login', {
        method: 'POST',
        body: JSON.stringify(body),
      }),
    me: () => request<AuthUser>('/auth/me'),
    logout: () =>
      request<{ ok: true }>('/auth/logout', { method: 'POST' }),
  },
  dashboard: () => request<DashboardStats>('/dashboard'),
  extensions: {
    list: () => request<Extension[]>('/extensions'),
    create: (body: CreateExtensionInput) =>
      request<Extension>('/extensions', {
        method: 'POST',
        body: JSON.stringify(body),
      }),
    remove: (id: string) =>
      request<{ ok: true }>(`/extensions/${id}`, { method: 'DELETE' }),
  },
  trunks: {
    list: () => request<Trunk[]>('/trunks'),
  },
  calls: {
    list: (activeOnly = false) =>
      request<CallRecord[]>(
        activeOnly ? '/calls?active=true' : '/calls',
      ),
  },
};
