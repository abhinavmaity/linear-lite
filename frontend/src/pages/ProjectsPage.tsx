import { useQuery } from '@tanstack/react-query';
import { EmptyState } from 'components/common/EmptyState';
import { ErrorBanner } from 'components/common/ErrorBanner';
import { PageHeader } from 'components/common/PageHeader';
import { Spinner } from 'components/common/Spinner';
import { projectsApi } from 'services/projectsApi';

export function ProjectsPage() {
  const projects = useQuery({
    queryKey: ['projects', 'page'],
    queryFn: () => projectsApi.list({ limit: 50 }).then((response) => response.items),
  });

  return (
    <div>
      <PageHeader title="Projects" subtitle="Scaffolded route with live list data from the project summary endpoint." />
      {projects.isLoading ? <Spinner label="Loading projects" /> : null}
      {projects.isError ? <ErrorBanner message={(projects.error as Error).message} /> : null}
      {projects.data?.length ? (
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(260px, 1fr))', gap: 16 }}>
          {projects.data.map((project) => (
            <article key={project.id} className="panel" style={{ padding: 18 }}>
              <div className="label" style={{ color: 'var(--text-secondary)', marginBottom: 8 }}>
                {project.key}
              </div>
              <h3 style={{ margin: 0 }}>{project.name}</h3>
              <p style={{ color: 'var(--text-secondary)' }}>{project.description ?? 'No description'}</p>
            </article>
          ))}
        </div>
      ) : null}
      {projects.data && projects.data.length === 0 ? (
        <EmptyState title="No projects" description="Projects will appear here once the backend contains project data." />
      ) : null}
    </div>
  );
}
