export interface UserSummary {
  id: string;
  email: string;
  name: string;
  avatar_url: string | null;
  created_at: string;
  updated_at: string;
}

export interface UserDetail extends UserSummary {
  stats: {
    total_created: number;
    total_assigned: number;
    in_progress_assigned: number;
    done_assigned: number;
  };
}

export interface IssueCounts {
  total: number;
  backlog: number;
  todo: number;
  in_progress: number;
  in_review: number;
  done: number;
  cancelled: number;
}

export interface SprintSummary {
  id: string;
  name: string;
  description: string | null;
  project_id: string;
  start_date: string;
  end_date: string;
  status: 'planned' | 'active' | 'completed';
  created_at: string;
  updated_at: string;
  issue_counts: IssueCounts;
}

export interface ProjectSummary {
  id: string;
  name: string;
  description: string | null;
  key: string;
  created_by: string;
  created_at: string;
  updated_at: string;
  issue_counts: IssueCounts;
  active_sprint: SprintSummary | null;
}

export interface ProjectDetail extends ProjectSummary {
  creator: UserSummary;
  sprints: SprintSummary[];
}

export interface SprintDetail extends SprintSummary {
  project: ProjectSummary;
}

export interface Label {
  id: string;
  name: string;
  color: string;
  description: string | null;
  created_at: string;
  updated_at: string;
}

export interface LabelDetail extends Label {
  usage_count: number;
}

export interface IssueActivity {
  id: string;
  issue_id: string;
  user_id: string;
  action: string;
  field_name: string | null;
  old_value: string | null;
  new_value: string | null;
  created_at: string;
  user: UserSummary;
}

export type IssueStatus = 'backlog' | 'todo' | 'in_progress' | 'in_review' | 'done' | 'cancelled';
export type IssuePriority = 'low' | 'medium' | 'high' | 'urgent';

export interface IssueSummary {
  id: string;
  identifier: string;
  title: string;
  description: string | null;
  status: IssueStatus;
  priority: IssuePriority;
  project_id: string;
  sprint_id: string | null;
  assignee_id: string | null;
  created_by: string;
  archived_at: string | null;
  archived_by: string | null;
  created_at: string;
  updated_at: string;
  project: ProjectSummary;
  sprint: SprintSummary | null;
  assignee: UserSummary | null;
  creator: UserSummary;
  labels: Label[];
}

export interface IssueDetail extends IssueSummary {
  activities: IssueActivity[];
}

export interface DashboardStats {
  total_issues: number;
  my_issues: number;
  in_progress: number;
  done_this_week: number;
  active_sprint: SprintSummary | null;
  recent_activity: IssueActivity[];
}

export interface AuthResponse {
  token: string;
  expires_at: string;
  user: UserSummary;
}
