import { apiClient } from './apiClient';
import { UserSummary } from 'types/domain';
import { CollectionResponse } from 'types/api';
import { safeArray } from 'utils/safeArray';

export const usersApi = {
  list(params?: { page?: number; limit?: number; search?: string; sort_by?: string; sort_order?: string }) {
    return apiClient.get<CollectionResponse<UserSummary>>('/users', params).then((response) => ({
      ...response,
      items: safeArray(response.items),
    }));
  },
};
