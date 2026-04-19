import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { issuesApi, IssueListParams, IssueUpsertInput } from 'services/issuesApi';

export function useIssuesList(params: IssueListParams) {
  return useQuery({
    queryKey: ['issues', params],
    queryFn: () => issuesApi.list(params),
  });
}

export function useIssueDetail(id?: string, includeArchived = false) {
  return useQuery({
    queryKey: ['issue', id, includeArchived],
    enabled: Boolean(id),
    queryFn: () => issuesApi.getById(id!, includeArchived).then((response) => response.data),
  });
}

export function useUpdateIssue(id: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (payload: IssueUpsertInput) => issuesApi.update(id, payload).then((response) => response.data),
    onSuccess: (issue, variables) => {
      queryClient.setQueryData(['issue', id, false], issue);
      queryClient.setQueryData(['issue', id, true], issue);

      if (shouldInvalidateIssuesList(variables)) {
        queryClient.invalidateQueries({ queryKey: ['issues'] });
      }
      if (shouldInvalidateDashboard(variables)) {
        queryClient.invalidateQueries({ queryKey: ['dashboard'] });
      }
    },
  });
}

export function useArchiveIssue(id: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => issuesApi.archive(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['issues'] });
      queryClient.invalidateQueries({ queryKey: ['dashboard'] });
    },
  });
}

function shouldInvalidateIssuesList(payload: IssueUpsertInput) {
  const keys = Object.keys(payload);
  if (keys.length === 0) return false;
  return keys.some((key) => key !== 'title' && key !== 'description');
}

function shouldInvalidateDashboard(payload: IssueUpsertInput) {
  return ['status', 'priority', 'project_id', 'sprint_id', 'assignee_id', 'label_ids', 'archived'].some((key) =>
    Object.prototype.hasOwnProperty.call(payload, key),
  );
}
