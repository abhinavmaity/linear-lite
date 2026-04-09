import { apiClient } from './apiClient';
import { ProjectDetail, ProjectSummary } from 'types/domain';
import { CollectionResponse, SingleResponse } from 'types/api';

export interface ProjectListParams {
  page?: number;
  limit?: number;
  search?: string;
  sort_by?: 'name' | 'created_at' | 'updated_at';
  sort_order?: 'asc' | 'desc';
}

export interface ProjectCreateInput {
  name: string;
  key: string;
  description?: string | null;
}

export interface ProjectUpdateInput {
  name?: string;
  key?: string;
  description?: string | null;
}

export const projectsApi = {
  list(params?: ProjectListParams) {
    return apiClient.get<CollectionResponse<ProjectSummary>>('/projects', params as Record<string, unknown> | undefined);
  },
  getById(id: string) {
    return apiClient.get<SingleResponse<ProjectDetail>>(`/projects/${id}`);
  },
  create(payload: ProjectCreateInput) {
    return apiClient.post<SingleResponse<ProjectDetail>>('/projects', payload);
  },
  update(id: string, payload: ProjectUpdateInput) {
    return apiClient.put<SingleResponse<ProjectDetail>>(`/projects/${id}`, payload);
  },
  delete(id: string) {
    return apiClient.delete<void>(`/projects/${id}`);
  },
};
