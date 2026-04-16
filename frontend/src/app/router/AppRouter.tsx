import { Navigate, Outlet, Route, Routes } from 'react-router-dom';
import { AppShell } from 'components/common/AppShell';
import { CreateIssueModal } from 'components/issues/CreateIssueModal';
import { AuthGate } from 'features/auth/AuthGate';
import { DashboardPage } from 'pages/DashboardPage';
import { BuildStoryPage } from 'pages/BuildStoryPage';
import { IssueDetailPage } from 'pages/IssueDetailPage';
import { IssuesBoardPage } from 'pages/IssuesBoardPage';
import { IssuesListPage } from 'pages/IssuesListPage';
import { LabelsPage } from 'pages/LabelsPage';
import { LoginPage } from 'pages/LoginPage';
import { ProjectsPage } from 'pages/ProjectsPage';
import { RegisterPage } from 'pages/RegisterPage';
import { SprintsPage } from 'pages/SprintsPage';
import { TeamPage } from 'pages/TeamPage';

function ProtectedLayout() {
  return (
    <AuthGate requireAuth>
      <AppShell>
        <Outlet />
        <CreateIssueModal />
      </AppShell>
    </AuthGate>
  );
}

function PublicLayout() {
  return (
    <AuthGate requireAuth={false}>
      <Outlet />
    </AuthGate>
  );
}

export function AppRouter() {
  return (
    <Routes>
      <Route path="/build-story" element={<BuildStoryPage />} />
      <Route element={<PublicLayout />}>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
      </Route>
      <Route element={<ProtectedLayout />}>
        <Route path="/dashboard" element={<DashboardPage />} />
        <Route path="/issues" element={<IssuesListPage />} />
        <Route path="/board" element={<IssuesBoardPage />} />
        <Route path="/issues/:id" element={<IssueDetailPage />} />
        <Route path="/projects" element={<ProjectsPage />} />
        <Route path="/sprints" element={<SprintsPage />} />
        <Route path="/labels" element={<LabelsPage />} />
        <Route path="/team" element={<TeamPage />} />
      </Route>
      <Route path="/" element={<Navigate to="/dashboard" replace />} />
      <Route path="*" element={<Navigate to="/dashboard" replace />} />
    </Routes>
  );
}
