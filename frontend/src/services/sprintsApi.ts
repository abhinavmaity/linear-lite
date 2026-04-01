import { apiClient } from './apiClient';
import { SprintSummary } from 'types/domain';
import { CollectionResponse } from 'types/api';

export const sprintsApi = {
  list(params?: {
    page?: number;
    limit?: number;
    project_id?: string;
    status?: string;
    search?: string;
    sort_by?: string;
    sort_order?: string;
  }) {
    return apiClient.get<CollectionResponse<SprintSummary>>('/sprints', params);
  },
};
