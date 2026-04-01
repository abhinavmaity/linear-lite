import { useQuery } from '@tanstack/react-query';
import { labelsApi } from 'services/labelsApi';
import { projectsApi } from 'services/projectsApi';
import { sprintsApi } from 'services/sprintsApi';
import { usersApi } from 'services/usersApi';

export function useUsersSelector() {
  return useQuery({
    queryKey: ['users', 'selector'],
    queryFn: () => usersApi.list({ limit: 100 }).then((response) => response.items),
  });
}

export function useLabelsSelector() {
  return useQuery({
    queryKey: ['labels', 'selector'],
    queryFn: () => labelsApi.list({ limit: 100 }).then((response) => response.items),
  });
}

export function useProjectsSelector() {
  return useQuery({
    queryKey: ['projects', 'selector'],
    queryFn: () => projectsApi.list({ limit: 100 }).then((response) => response.items),
  });
}

export function useSprintsSelector(projectId?: string) {
  return useQuery({
    queryKey: ['sprints', 'selector', projectId],
    queryFn: () => sprintsApi.list({ limit: 100, project_id: projectId }).then((response) => response.items),
  });
}
