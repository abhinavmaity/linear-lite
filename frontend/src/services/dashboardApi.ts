import { apiClient } from './apiClient';
import { DashboardStats } from 'types/domain';
import { SingleResponse } from 'types/api';

export const dashboardApi = {
  getStats() {
    return apiClient.get<SingleResponse<DashboardStats>>('/dashboard/stats');
  },
};
