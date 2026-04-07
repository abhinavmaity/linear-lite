-- Milestone 3 foundation: core issue-workflow schema.
-- Source of truth: docs/Technical_Architecture.md

CREATE TABLE projects (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(255) NOT NULL,
  description TEXT NULL,
  key VARCHAR(10) NOT NULL,
  next_issue_number INTEGER NOT NULL DEFAULT 1,
  created_by UUID NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT chk_projects_name_not_blank CHECK (btrim(name) <> ''),
  CONSTRAINT chk_projects_name_len CHECK (char_length(name) BETWEEN 1 AND 255),
  CONSTRAINT chk_projects_key_format CHECK (key ~ '^[A-Z0-9]{2,10}$'),
  CONSTRAINT chk_projects_next_issue_number_positive CHECK (next_issue_number >= 1),
  CONSTRAINT fk_projects_created_by
    FOREIGN KEY (created_by)
    REFERENCES users(id)
    ON DELETE RESTRICT
    ON UPDATE RESTRICT
);

CREATE UNIQUE INDEX uq_projects_key ON projects (key);
CREATE INDEX idx_projects_created_by ON projects (created_by);

CREATE TRIGGER trg_projects_updated_at
BEFORE UPDATE ON projects
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TABLE sprints (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(255) NOT NULL,
  description TEXT NULL,
  project_id UUID NOT NULL,
  start_date DATE NOT NULL,
  end_date DATE NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'planned',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT chk_sprints_name_not_blank CHECK (btrim(name) <> ''),
  CONSTRAINT chk_sprints_name_len CHECK (char_length(name) BETWEEN 1 AND 255),
  CONSTRAINT chk_sprints_status CHECK (status IN ('planned', 'active', 'completed')),
  CONSTRAINT chk_sprints_date_range CHECK (end_date >= start_date),
  CONSTRAINT fk_sprints_project_id
    FOREIGN KEY (project_id)
    REFERENCES projects(id)
    ON DELETE RESTRICT
    ON UPDATE RESTRICT
);

CREATE INDEX idx_sprints_project_id ON sprints (project_id);
CREATE INDEX idx_sprints_status ON sprints (status);
CREATE INDEX idx_sprints_project_dates ON sprints (project_id, start_date, end_date);
CREATE UNIQUE INDEX uq_sprints_one_active_per_project
  ON sprints (project_id)
  WHERE status = 'active';

CREATE TRIGGER trg_sprints_updated_at
BEFORE UPDATE ON sprints
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TABLE labels (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(50) NOT NULL,
  color VARCHAR(7) NOT NULL,
  description TEXT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT chk_labels_name_not_blank CHECK (btrim(name) <> ''),
  CONSTRAINT chk_labels_name_len CHECK (char_length(name) BETWEEN 1 AND 50),
  CONSTRAINT chk_labels_color_hex CHECK (color ~ '^#[0-9A-Fa-f]{6}$')
);

CREATE UNIQUE INDEX uq_labels_lower_name ON labels (LOWER(name));

CREATE TRIGGER trg_labels_updated_at
BEFORE UPDATE ON labels
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TABLE issues (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  identifier VARCHAR(32) NOT NULL,
  title VARCHAR(500) NOT NULL,
  description TEXT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'backlog',
  priority VARCHAR(10) NOT NULL DEFAULT 'medium',
  project_id UUID NOT NULL,
  sprint_id UUID NULL,
  assignee_id UUID NULL,
  created_by UUID NOT NULL,
  archived_at TIMESTAMPTZ NULL,
  archived_by UUID NULL,
  search_vector tsvector GENERATED ALWAYS AS (
    setweight(to_tsvector('english', coalesce(title, '')), 'A') ||
    setweight(to_tsvector('english', coalesce(description, '')), 'B')
  ) STORED,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT chk_issues_title_not_blank CHECK (btrim(title) <> ''),
  CONSTRAINT chk_issues_title_len CHECK (char_length(title) BETWEEN 1 AND 500),
  CONSTRAINT chk_issues_status CHECK (
    status IN ('backlog', 'todo', 'in_progress', 'in_review', 'done', 'cancelled')
  ),
  CONSTRAINT chk_issues_priority CHECK (
    priority IN ('low', 'medium', 'high', 'urgent')
  ),
  CONSTRAINT chk_issues_identifier_not_blank CHECK (btrim(identifier) <> ''),
  CONSTRAINT chk_issues_archive_pair CHECK (
    (archived_at IS NULL AND archived_by IS NULL) OR
    (archived_at IS NOT NULL AND archived_by IS NOT NULL)
  ),
  CONSTRAINT fk_issues_project_id
    FOREIGN KEY (project_id)
    REFERENCES projects(id)
    ON DELETE RESTRICT
    ON UPDATE RESTRICT,
  CONSTRAINT fk_issues_sprint_id
    FOREIGN KEY (sprint_id)
    REFERENCES sprints(id)
    ON DELETE RESTRICT
    ON UPDATE RESTRICT,
  CONSTRAINT fk_issues_assignee_id
    FOREIGN KEY (assignee_id)
    REFERENCES users(id)
    ON DELETE SET NULL
    ON UPDATE RESTRICT,
  CONSTRAINT fk_issues_created_by
    FOREIGN KEY (created_by)
    REFERENCES users(id)
    ON DELETE RESTRICT
    ON UPDATE RESTRICT,
  CONSTRAINT fk_issues_archived_by
    FOREIGN KEY (archived_by)
    REFERENCES users(id)
    ON DELETE SET NULL
    ON UPDATE RESTRICT
);

CREATE UNIQUE INDEX uq_issues_identifier ON issues (identifier);
CREATE INDEX idx_issues_project_id ON issues (project_id);
CREATE INDEX idx_issues_sprint_id ON issues (sprint_id);
CREATE INDEX idx_issues_assignee_id ON issues (assignee_id);
CREATE INDEX idx_issues_created_by ON issues (created_by);
CREATE INDEX idx_issues_status ON issues (status);
CREATE INDEX idx_issues_priority ON issues (priority);
CREATE INDEX idx_issues_archived_at ON issues (archived_at);
CREATE INDEX idx_issues_project_status_updated_at
  ON issues (project_id, status, updated_at DESC);
CREATE INDEX idx_issues_assignee_status_updated_at
  ON issues (assignee_id, status, updated_at DESC);
CREATE INDEX idx_issues_search_vector ON issues USING GIN (search_vector);

CREATE TRIGGER trg_issues_updated_at
BEFORE UPDATE ON issues
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TABLE issue_labels (
  issue_id UUID NOT NULL,
  label_id UUID NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (issue_id, label_id),
  CONSTRAINT fk_issue_labels_issue_id
    FOREIGN KEY (issue_id)
    REFERENCES issues(id)
    ON DELETE CASCADE
    ON UPDATE RESTRICT,
  CONSTRAINT fk_issue_labels_label_id
    FOREIGN KEY (label_id)
    REFERENCES labels(id)
    ON DELETE RESTRICT
    ON UPDATE RESTRICT
);

CREATE INDEX idx_issue_labels_issue_id ON issue_labels (issue_id);
CREATE INDEX idx_issue_labels_label_id ON issue_labels (label_id);

CREATE TABLE issue_activities (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  issue_id UUID NOT NULL,
  user_id UUID NOT NULL,
  action VARCHAR(50) NOT NULL,
  field_name VARCHAR(100) NULL,
  old_value TEXT NULL,
  new_value TEXT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT chk_issue_activities_action CHECK (
    action IN (
      'created',
      'updated',
      'title_changed',
      'description_changed',
      'status_changed',
      'priority_changed',
      'assignee_changed',
      'sprint_changed',
      'project_changed',
      'label_added',
      'label_removed',
      'archived',
      'restored'
    )
  ),
  CONSTRAINT fk_issue_activities_issue_id
    FOREIGN KEY (issue_id)
    REFERENCES issues(id)
    ON DELETE CASCADE
    ON UPDATE RESTRICT,
  CONSTRAINT fk_issue_activities_user_id
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE RESTRICT
    ON UPDATE RESTRICT
);

CREATE INDEX idx_issue_activities_issue_id_created_at
  ON issue_activities (issue_id, created_at DESC);
CREATE INDEX idx_issue_activities_user_id_created_at
  ON issue_activities (user_id, created_at DESC);
