import { apiClient } from './apiClient';
import { IssueDetail, IssuePriority, IssueStatus, IssueSummary } from 'types/domain';
import { CollectionResponse, SingleResponse } from 'types/api';
import { safeArray } from 'utils/safeArray';

export interface IssueListParams {
  page?: number;
  limit?: number;
  sort_by?: 'identifier' | 'title' | 'status' | 'priority' | 'created_at' | 'updated_at';
  sort_order?: 'asc' | 'desc';
  search?: string;
  status?: IssueStatus[];
  priority?: IssuePriority[];
  assignee_id?: string;
  project_id?: string;
  sprint_id?: string;
  label_id?: string[];
  label_mode?: 'any' | 'all';
  include_archived?: boolean;
}

export interface IssueUpsertInput {
  title?: string;
  description?: string | null;
  status?: IssueStatus;
  priority?: IssuePriority;
  project_id?: string;
  sprint_id?: string | null;
  assignee_id?: string | null;
  label_ids?: string[];
  archived?: boolean;
}

function normalizeIssueSummary(issue: IssueSummary): IssueSummary {
  return {
    ...issue,
    labels: safeArray(issue.labels),
  };
}

function normalizeIssueDetail(issue: IssueDetail): IssueDetail {
  return {
    ...normalizeIssueSummary(issue),
    activities: safeArray(issue.activities),
  };
}

export const issuesApi = {
  list(params: IssueListParams) {
    return apiClient
      .get<CollectionResponse<IssueSummary>>('/issues', params as Record<string, unknown>)
      .then((response) => ({
        ...response,
        items: response.items.map(normalizeIssueSummary),
      }));
  },
  getById(id: string, includeArchived = false) {
    return apiClient
      .get<SingleResponse<IssueDetail>>(`/issues/${id}`, { include_archived: includeArchived })
      .then((response) => ({
        ...response,
        data: normalizeIssueDetail(response.data),
      }));
  },
  create(payload: Required<Pick<IssueUpsertInput, 'title' | 'project_id'>> & IssueUpsertInput) {
    return apiClient.post<SingleResponse<IssueDetail>>('/issues', payload).then((response) => ({
      ...response,
      data: normalizeIssueDetail(response.data),
    }));
  },
  update(id: string, payload: IssueUpsertInput) {
    return apiClient.put<SingleResponse<IssueDetail>>(`/issues/${id}`, payload).then((response) => ({
      ...response,
      data: normalizeIssueDetail(response.data),
    }));
  },
  archive(id: string) {
    return apiClient.delete<void>(`/issues/${id}`);
  },
};
