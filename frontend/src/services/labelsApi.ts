import { apiClient } from './apiClient';
import { Label, LabelDetail } from 'types/domain';
import { CollectionResponse, SingleResponse } from 'types/api';

export interface LabelListParams {
  page?: number;
  limit?: number;
  search?: string;
  sort_by?: 'name' | 'created_at';
  sort_order?: 'asc' | 'desc';
}

export interface LabelCreateInput {
  name: string;
  color: string;
  description?: string | null;
}

export interface LabelUpdateInput {
  name?: string;
  color?: string;
  description?: string | null;
}

export const labelsApi = {
  list(params?: LabelListParams) {
    return apiClient.get<CollectionResponse<Label>>('/labels', params as Record<string, unknown> | undefined);
  },
  getById(id: string) {
    return apiClient.get<SingleResponse<LabelDetail>>(`/labels/${id}`);
  },
  create(payload: LabelCreateInput) {
    return apiClient.post<SingleResponse<Label>>('/labels', payload);
  },
  update(id: string, payload: LabelUpdateInput) {
    return apiClient.put<SingleResponse<Label>>(`/labels/${id}`, payload);
  },
  delete(id: string) {
    return apiClient.delete<void>(`/labels/${id}`);
  },
};
