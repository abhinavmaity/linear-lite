import { useMemo, useState } from 'react';
import { ErrorBanner } from 'components/common/ErrorBanner';
import { PageHeader } from 'components/common/PageHeader';
import { Spinner } from 'components/common/Spinner';
import { IssueCard } from 'components/issues/IssueCard';
import { useIssuesList } from 'features/issues/issuesQueries';
import { issuesApi } from 'services/issuesApi';
import { useUIStore } from 'store/uiStore';
import { IssueSummary } from 'types/domain';

type BoardColumn = 'backlog' | 'todo' | 'in_progress' | 'in_review' | 'done';

const columns: BoardColumn[] = ['backlog', 'todo', 'in_progress', 'in_review', 'done'];

export function IssuesBoardPage() {
  const pushToast = useUIStore((state) => state.pushToast);
  const [items, setItems] = useState<ReturnType<typeof groupByStatus> | null>(null);
  const issues = useIssuesList({ page: 1, limit: 100, sort_by: 'updated_at', sort_order: 'desc' });

  const grouped = useMemo(() => {
    const base = groupByStatus(issues.data?.items ?? []);
    return items ?? base;
  }, [issues.data?.items, items]);

  async function moveIssue(issueId: string, nextStatus: BoardColumn) {
    if (!issues.data) return;
    const previous = groupByStatus(issues.data.items);
    const optimistic = groupByStatus(
      issues.data.items.map((issue) => (issue.id === issueId ? { ...issue, status: nextStatus } : issue)),
    );
    setItems(optimistic);

    try {
      await issuesApi.update(issueId, { status: nextStatus });
      pushToast({ tone: 'success', message: 'Issue updated.' });
      setItems(null);
    } catch (error) {
      setItems(previous);
      pushToast({ tone: 'error', message: error instanceof Error ? error.message : 'Failed to update issue.' });
    }
  }

  return (
    <div>
      <PageHeader title="Board" subtitle="Shared issue query layer with optimistic drag-and-drop status updates." />
      {issues.isLoading ? <Spinner label="Loading board" /> : null}
      {issues.isError ? <ErrorBanner message={(issues.error as Error).message} /> : null}
      {issues.data ? (
        <div
          style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(5, minmax(220px, 1fr))',
            gap: 16,
            alignItems: 'start',
            overflowX: 'auto',
          }}
        >
          {columns.map((column) => (
            <section
              key={column}
              className="panel"
              onDragOver={(event) => event.preventDefault()}
              onDrop={(event) => {
                const issueId = event.dataTransfer.getData('text/plain');
                if (issueId) {
                  moveIssue(issueId, column);
                }
              }}
              style={{ padding: 14, minHeight: 420 }}
            >
              <div className="label" style={{ fontSize: 18, marginBottom: 12 }}>
                {column.replace('_', ' ')} ({grouped[column].length})
              </div>
              <div style={{ display: 'grid', gap: 12 }}>
                {grouped[column].map((issue) => (
                  <div key={issue.id} draggable onDragStart={(event) => event.dataTransfer.setData('text/plain', issue.id)}>
                    <IssueCard issue={issue} compact />
                  </div>
                ))}
              </div>
            </section>
          ))}
        </div>
      ) : null}
    </div>
  );
}

function groupByStatus(items: IssueSummary[]): Record<BoardColumn, IssueSummary[]> {
  return {
    backlog: items.filter((issue) => issue.status === 'backlog'),
    todo: items.filter((issue) => issue.status === 'todo'),
    in_progress: items.filter((issue) => issue.status === 'in_progress'),
    in_review: items.filter((issue) => issue.status === 'in_review'),
    done: items.filter((issue) => issue.status === 'done'),
  };
}
