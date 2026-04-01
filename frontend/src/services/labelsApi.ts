import { apiClient } from './apiClient';
import { Label } from 'types/domain';
import { CollectionResponse } from 'types/api';

export const labelsApi = {
  list(params?: { page?: number; limit?: number; search?: string; sort_by?: string; sort_order?: string }) {
    return apiClient.get<CollectionResponse<Label>>('/labels', params);
  },
};
