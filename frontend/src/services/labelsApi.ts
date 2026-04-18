import { apiClient } from './apiClient';
import { Label, LabelDetail } from 'types/domain';
import { CollectionResponse, SingleResponse } from 'types/api';
import { safeArray } from 'utils/safeArray';

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

function normalizeLabel(label: Label): Label {
  return label;
}

function normalizeLabelDetail(label: LabelDetail): LabelDetail {
  return {
    ...normalizeLabel(label),
    usage_count: label.usage_count ?? 0,
  };
}

export const labelsApi = {
  list(params?: LabelListParams) {
    return apiClient
      .get<CollectionResponse<Label>>('/labels', params as Record<string, unknown> | undefined)
      .then((response) => ({
        ...response,
        items: safeArray(response.items).map(normalizeLabel),
      }));
  },
  getById(id: string) {
    return apiClient.get<SingleResponse<LabelDetail>>(`/labels/${id}`).then((response) => ({
      ...response,
      data: normalizeLabelDetail(response.data),
    }));
  },
  create(payload: LabelCreateInput) {
    return apiClient.post<SingleResponse<Label>>('/labels', payload).then((response) => ({
      ...response,
      data: normalizeLabel(response.data),
    }));
  },
  update(id: string, payload: LabelUpdateInput) {
    return apiClient.put<SingleResponse<Label>>(`/labels/${id}`, payload).then((response) => ({
      ...response,
      data: normalizeLabel(response.data),
    }));
  },
  delete(id: string) {
    return apiClient.delete<void>(`/labels/${id}`);
  },
};
