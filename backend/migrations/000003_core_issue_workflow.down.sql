DROP INDEX IF EXISTS idx_issue_activities_user_id_created_at;
DROP INDEX IF EXISTS idx_issue_activities_issue_id_created_at;
DROP TABLE IF EXISTS issue_activities;

DROP INDEX IF EXISTS idx_issue_labels_label_id;
DROP INDEX IF EXISTS idx_issue_labels_issue_id;
DROP TABLE IF EXISTS issue_labels;

DROP TRIGGER IF EXISTS trg_issues_updated_at ON issues;
DROP INDEX IF EXISTS idx_issues_search_vector;
DROP INDEX IF EXISTS idx_issues_assignee_status_updated_at;
DROP INDEX IF EXISTS idx_issues_project_status_updated_at;
DROP INDEX IF EXISTS idx_issues_archived_at;
DROP INDEX IF EXISTS idx_issues_priority;
DROP INDEX IF EXISTS idx_issues_status;
DROP INDEX IF EXISTS idx_issues_created_by;
DROP INDEX IF EXISTS idx_issues_assignee_id;
DROP INDEX IF EXISTS idx_issues_sprint_id;
DROP INDEX IF EXISTS idx_issues_project_id;
DROP INDEX IF EXISTS uq_issues_identifier;
DROP TABLE IF EXISTS issues;

DROP TRIGGER IF EXISTS trg_labels_updated_at ON labels;
DROP INDEX IF EXISTS uq_labels_lower_name;
DROP TABLE IF EXISTS labels;

DROP TRIGGER IF EXISTS trg_sprints_updated_at ON sprints;
DROP INDEX IF EXISTS uq_sprints_one_active_per_project;
DROP INDEX IF EXISTS idx_sprints_project_dates;
DROP INDEX IF EXISTS idx_sprints_status;
DROP INDEX IF EXISTS idx_sprints_project_id;
DROP TABLE IF EXISTS sprints;

DROP TRIGGER IF EXISTS trg_projects_updated_at ON projects;
DROP INDEX IF EXISTS idx_projects_created_by;
DROP INDEX IF EXISTS uq_projects_key;
DROP TABLE IF EXISTS projects;
