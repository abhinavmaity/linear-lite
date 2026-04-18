import { apiClient } from './apiClient';
import { DashboardStats } from 'types/domain';
import { SingleResponse } from 'types/api';
import { safeArray } from 'utils/safeArray';

function normalizeDashboardStats(stats: DashboardStats): DashboardStats {
  return {
    ...stats,
    recent_activity: safeArray(stats.recent_activity),
  };
}

export const dashboardApi = {
  getStats() {
    return apiClient.get<SingleResponse<DashboardStats>>('/dashboard/stats').then((response) => ({
      ...response,
      data: normalizeDashboardStats(response.data),
    }));
  },
};
