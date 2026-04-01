import { apiClient } from './apiClient';
import { ProjectSummary } from 'types/domain';
import { CollectionResponse } from 'types/api';

export const projectsApi = {
  list(params?: { page?: number; limit?: number; search?: string; sort_by?: string; sort_order?: string }) {
    return apiClient.get<CollectionResponse<ProjectSummary>>('/projects', params);
  },
};
