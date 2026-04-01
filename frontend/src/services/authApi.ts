import { apiClient } from './apiClient';
import { AuthResponse, UserSummary } from 'types/domain';
import { SingleResponse } from 'types/api';

export const authApi = {
  login(payload: { email: string; password: string }) {
    return apiClient.post<SingleResponse<AuthResponse>>('/auth/login', payload);
  },
  register(payload: { name: string; email: string; password: string }) {
    return apiClient.post<SingleResponse<AuthResponse>>('/auth/register', payload);
  },
  me() {
    return apiClient.get<SingleResponse<UserSummary>>('/auth/me');
  },
};
