import { apiClient } from './apiClient';
import { SprintDetail, SprintSummary } from 'types/domain';
import { CollectionResponse, SingleResponse } from 'types/api';
import { safeArray } from 'utils/safeArray';

export interface SprintListParams {
  page?: number;
  limit?: number;
  project_id?: string;
  status?: 'planned' | 'active' | 'completed';
  search?: string;
  sort_by?: 'name' | 'start_date' | 'end_date' | 'created_at';
  sort_order?: 'asc' | 'desc';
}

export interface SprintCreateInput {
  name: string;
  project_id: string;
  start_date: string;
  end_date: string;
  description?: string | null;
  status?: 'planned' | 'active' | 'completed';
}

export interface SprintUpdateInput {
  name?: string;
  start_date?: string;
  end_date?: string;
  description?: string | null;
  status?: 'planned' | 'active' | 'completed';
}

function normalizeSprintSummary(sprint: SprintSummary): SprintSummary {
  return sprint;
}

function normalizeSprintDetail(sprint: SprintDetail): SprintDetail {
  return {
    ...normalizeSprintSummary(sprint),
    project: sprint.project,
  };
}

export const sprintsApi = {
  list(params?: SprintListParams) {
    return apiClient
      .get<CollectionResponse<SprintSummary>>('/sprints', params as Record<string, unknown> | undefined)
      .then((response) => ({
        ...response,
        items: safeArray(response.items).map(normalizeSprintSummary),
      }));
  },
  getById(id: string) {
    return apiClient.get<SingleResponse<SprintDetail>>(`/sprints/${id}`).then((response) => ({
      ...response,
      data: normalizeSprintDetail(response.data),
    }));
  },
  create(payload: SprintCreateInput) {
    return apiClient.post<SingleResponse<SprintDetail>>('/sprints', payload).then((response) => ({
      ...response,
      data: normalizeSprintDetail(response.data),
    }));
  },
  update(id: string, payload: SprintUpdateInput) {
    return apiClient.put<SingleResponse<SprintDetail>>(`/sprints/${id}`, payload).then((response) => ({
      ...response,
      data: normalizeSprintDetail(response.data),
    }));
  },
  delete(id: string) {
    return apiClient.delete<void>(`/sprints/${id}`);
  },
};
