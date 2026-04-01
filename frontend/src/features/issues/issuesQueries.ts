import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { issuesApi, IssueListParams, IssueUpsertInput } from 'services/issuesApi';

export function useIssuesList(params: IssueListParams) {
  return useQuery({
    queryKey: ['issues', params],
    queryFn: () => issuesApi.list(params),
  });
}

export function useIssueDetail(id?: string) {
  return useQuery({
    queryKey: ['issue', id],
    enabled: Boolean(id),
    queryFn: () => issuesApi.getById(id!).then((response) => response.data),
  });
}

export function useUpdateIssue(id: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (payload: IssueUpsertInput) => issuesApi.update(id, payload).then((response) => response.data),
    onSuccess: (issue) => {
      queryClient.setQueryData(['issue', id], issue);
      queryClient.invalidateQueries({ queryKey: ['issues'] });
      queryClient.invalidateQueries({ queryKey: ['dashboard'] });
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
