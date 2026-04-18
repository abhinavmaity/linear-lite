import { apiClient } from './apiClient';
import { ProjectDetail, ProjectSummary } from 'types/domain';
import { CollectionResponse, SingleResponse } from 'types/api';
import { safeArray } from 'utils/safeArray';

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

function normalizeProjectDetail(project: ProjectDetail): ProjectDetail {
  return {
    ...project,
    sprints: safeArray(project.sprints),
  };
}

export const projectsApi = {
  list(params?: ProjectListParams) {
    return apiClient
      .get<CollectionResponse<ProjectSummary>>('/projects', params as Record<string, unknown> | undefined)
      .then((response) => ({
        ...response,
        items: safeArray(response.items),
      }));
  },
  getById(id: string) {
    return apiClient.get<SingleResponse<ProjectDetail>>(`/projects/${id}`).then((response) => ({
      ...response,
      data: normalizeProjectDetail(response.data),
    }));
  },
  create(payload: ProjectCreateInput) {
    return apiClient.post<SingleResponse<ProjectDetail>>('/projects', payload).then((response) => ({
      ...response,
      data: normalizeProjectDetail(response.data),
    }));
  },
  update(id: string, payload: ProjectUpdateInput) {
    return apiClient.put<SingleResponse<ProjectDetail>>(`/projects/${id}`, payload).then((response) => ({
      ...response,
      data: normalizeProjectDetail(response.data),
    }));
  },
  delete(id: string) {
    return apiClient.delete<void>(`/projects/${id}`);
  },
};
