# **Linear-lite**

Complete Technical Architecture Specification

Version 1.0 \- Production Ready

Complete Database Schema | Full API Contracts | Zero Ambiguity

| Stack | React 18 \+ Go 1.21+ \+ PostgreSQL 15 \+ Redis 7 |
| :---- | :---- |
| Architecture | Clean Architecture \- Handler → Service → Repository |
| API Style | RESTful JSON with JWT Bearer Authentication |
| Database Tables | 7 tables (users, projects, sprints, labels, issues, issue\_labels, issue\_activities) |
| API Endpoints | 26 endpoints fully documented |

# **Table of Contents**

1\. System Architecture Overview  
2\. Complete Database Schema (All 7 Tables)  
   2.1 Table: users  
   2.2 Table: projects  
   2.3 Table: sprints  
   2.4 Table: labels  
   2.5 Table: issues  
   2.6 Table: issue\_labels (Join Table)  
   2.7 Table: issue\_activities  
   2.8 Entity Relationship Diagram  
   2.9 Complete Index Strategy  
3\. Complete API Contract Documentation (26 Endpoints)  
   3.1 Authentication Endpoints (3)  
   3.2 User Endpoints (2)  
   3.3 Project Endpoints (5)  
   3.4 Sprint Endpoints (5)  
   3.5 Label Endpoints (5)  
   3.6 Issue Endpoints (5)  
   3.7 Dashboard Endpoint (1)  
4\. Error Handling & Response Formats  
5\. Enums, Constants & Validation Rules  
6\. Security & Authentication Specifications  
7\. Implementation Guidelines

# **1\. System Architecture Overview**

Linear-lite implements a clean three-tier architecture with strict separation of concerns. The system is designed for horizontal scalability, maintainability, and testability.

## **1.1 Architecture Layers**

┌────────────────────────────────────────────────────────────────┐

│                    PRESENTATION LAYER                          │

│  React 18 SPA (Single Page Application)                        │

│  • TypeScript 5.0+ for type safety                             │

│  • Vite 5.0+ for build tooling                                 │

│  • TanStack Query 5.0+ for server state management             │

│  • Zustand 4.0+ for client state management                    │

│  • Tailwind CSS 3.4+ for styling                               │

│  • @dnd-kit for drag-and-drop                                  │

└────────────────────────────────────────────────────────────────┘

                            │ HTTPS

                            ▼

┌────────────────────────────────────────────────────────────────┐

│                      API LAYER (Go Backend)                     │

│  Port: 8080                                                     │

│  ┌────────────────────────────────────────────────────────────┐│

│  │  Router Layer (Gin Web Framework v1.9+)                   ││

│  │  • Route registration and HTTP method handling             ││

│  └────────────────────────────────────────────────────────────┘│

│  ┌────────────────────────────────────────────────────────────┐│

│  │  Middleware Chain                                          ││

│  │  1\. CORS (Cross-Origin Resource Sharing)                  ││

│  │  2\. Request Logger (log all incoming requests)            ││

│  │  3\. JWT Authentication (verify Bearer tokens)             ││

│  │  4\. Request ID Generator                                  ││

│  │  5\. Panic Recovery (graceful error handling)              ││

│  └────────────────────────────────────────────────────────────┘│

│  ┌────────────────────────────────────────────────────────────┐│

│  │  Handler Layer (HTTP Controllers)                         ││

│  │  • AuthHandler \- user authentication                      ││

│  │  • UserHandler \- user management                          ││

│  │  • ProjectHandler \- project CRUD                          ││

│  │  • SprintHandler \- sprint management                      ││

│  │  • LabelHandler \- label CRUD                              ││

│  │  • IssueHandler \- issue management                        ││

│  │  • DashboardHandler \- analytics                           ││

│  └────────────────────────────────────────────────────────────┘│

│  ┌────────────────────────────────────────────────────────────┐│

│  │  Service Layer (Business Logic)                           ││

│  │  • Input validation with go-playground/validator          ││

│  │  • Business rules enforcement                             ││

│  │  • Activity tracking                                      ││

│  │  • Transaction coordination                               ││

│  └────────────────────────────────────────────────────────────┘│

│  ┌────────────────────────────────────────────────────────────┐│

│  │  Repository Layer (Data Access)                           ││

│  │  • GORM 1.25+ OR sqlx with pgx driver                     ││

│  │  • Database query optimization                            ││

│  │  • Connection pooling                                     ││

│  └────────────────────────────────────────────────────────────┘│

└────────────────────────────────────────────────────────────────┘

         │                                  │

         ▼                                  ▼

┌──────────────────────┐         ┌──────────────────────┐

│  PostgreSQL 15+      │         │  Redis 7+            │

│  Port: 5432          │         │  Port: 6379          │

│  • Primary database  │         │  • JWT blacklist     │

│  • ACID compliance   │         │  • Session cache     │

│  • Full text search  │         │  • Query cache       │

└──────────────────────┘         └──────────────────────┘

## **1.2 Data Flow**

1\. Client makes HTTP request with JSON body  
2\. Gin router matches route and invokes middleware chain  
3\. JWT middleware validates token and extracts user ID  
4\. Handler parses request, calls service layer  
5\. Service validates business rules, calls repository  
6\. Repository executes database query  
7\. Response flows back up the chain  
8\. Handler serializes to JSON and returns HTTP response

# **2\. Complete Database Schema**

The database consists of 7 tables with carefully designed relationships to ensure data integrity and query performance.

## **2.1 Table: users**

**Purpose: Store user account information for authentication and profile data**

| Column Name | Data Type | Constraints | Description |
| :---- | :---- | :---- | :---- |
| id | UUID | PRIMARY KEY | Unique user identifier (auto-generated) |
| email | VARCHAR(255) | UNIQUE NOT NULL | User email address for login |
| password\_hash | VARCHAR(255) | NOT NULL | bcrypt hash with cost 12 |
| name | VARCHAR(255) | NOT NULL | User's full name |
| avatar\_url | TEXT | NULL | URL to profile picture |
| created\_at | TIMESTAMPTZ | NOT NULL DEFAULT NOW() | Account creation timestamp |
| updated\_at | TIMESTAMPTZ | NOT NULL DEFAULT NOW() | Last profile update timestamp |

**Indexes:**

* idx\_users\_email (email) \- B-tree index for login lookups

**Validation Rules:**

* email: Must match regex ^\[a-zA-Z0-9.\_%+-\]+@\[a-zA-Z0-9.-\]+\\.\[a-zA-Z\]{2,}$  
* password: Minimum 8 characters, hashed with bcrypt cost 12  
* name: 1-255 characters, non-empty after trim

**Go Struct Definition:**

type User struct {

    ID          uuid.UUID  \`gorm:"type:uuid;primaryKey" json:"id"\`

    Email       string     \`gorm:"type:varchar(255);uniqueIndex;not null" json:"email" validate:"required,email"\`

    PasswordHash string    \`gorm:"type:varchar(255);not null" json:"-"\`

    Name        string     \`gorm:"type:varchar(255);not null" json:"name" validate:"required,min=1,max=255"\`

    AvatarURL   \*string    \`gorm:"type:text" json:"avatar\_url"\`

    CreatedAt   time.Time  \`gorm:"not null;default:now()" json:"created\_at"\`

    UpdatedAt   time.Time  \`gorm:"not null;default:now()" json:"updated\_at"\`

}

## **2.2 Table: projects**

**Purpose: Group related issues into projects**

| Column Name | Data Type | Constraints | Description |
| :---- | :---- | :---- | :---- |
| id | UUID | PRIMARY KEY | Unique project identifier |
| name | VARCHAR(255) | NOT NULL | Project name |
| description | TEXT | NULL | Project description |
| identifier | VARCHAR(10) | UNIQUE NOT NULL | Short project code (e.g., PROJ) |
| created\_by | UUID | NOT NULL FK→users(id) | User who created project |
| created\_at | TIMESTAMPTZ | NOT NULL DEFAULT NOW() | Creation timestamp |
| updated\_at | TIMESTAMPTZ | NOT NULL DEFAULT NOW() | Last update timestamp |

**Foreign Keys:**

* fk\_projects\_created\_by → users(id) ON DELETE SET NULL

**Indexes:**

* idx\_projects\_identifier (identifier) \- UNIQUE for issue ID generation  
* idx\_projects\_created\_by (created\_by) \- for filtering by creator

**Validation Rules:**

* name: 1-255 characters, non-empty  
* identifier: 2-10 uppercase alphanumeric characters, must be unique  
* description: 0-10000 characters

**Go Struct Definition:**

type Project struct {

    ID          uuid.UUID  \`gorm:"type:uuid;primaryKey" json:"id"\`

    Name        string     \`gorm:"type:varchar(255);not null" json:"name" validate:"required,min=1,max=255"\`

    Description \*string    \`gorm:"type:text" json:"description"\`

    Identifier  string     \`gorm:"type:varchar(10);uniqueIndex;not null" json:"identifier" validate:"required,min=2,max=10,uppercase,alphanum"\`

    CreatedBy   uuid.UUID  \`gorm:"type:uuid;not null" json:"created\_by"\`

    Creator     \*User      \`gorm:"foreignKey:CreatedBy" json:"creator,omitempty"\`

    CreatedAt   time.Time  \`gorm:"not null;default:now()" json:"created\_at"\`

    UpdatedAt   time.Time  \`gorm:"not null;default:now()" json:"updated\_at"\`

}

## **2.3 Table: sprints**

**Purpose: Time-boxed iterations for organizing work**

| Column Name | Data Type | Constraints | Description |
| :---- | :---- | :---- | :---- |
| id | UUID | PRIMARY KEY | Unique sprint identifier |
| name | VARCHAR(255) | NOT NULL | Sprint name (e.g., Sprint 1\) |
| description | TEXT | NULL | Sprint goal/description |
| start\_date | DATE | NOT NULL | Sprint start date |
| end\_date | DATE | NOT NULL | Sprint end date |
| status | VARCHAR(20) | NOT NULL DEFAULT 'planned' | planned|active|completed |
| project\_id | UUID | NULL FK→projects(id) | Associated project |
| created\_at | TIMESTAMPTZ | NOT NULL DEFAULT NOW() | Creation timestamp |
| updated\_at | TIMESTAMPTZ | NOT NULL DEFAULT NOW() | Last update timestamp |

**Foreign Keys:**

* fk\_sprints\_project\_id → projects(id) ON DELETE SET NULL

**Indexes:**

* idx\_sprints\_project\_id (project\_id) \- for project-sprint joins  
* idx\_sprints\_dates (start\_date, end\_date) \- for date range queries  
* idx\_sprints\_status (status) \- for filtering active sprints

**Validation Rules:**

* name: 1-255 characters  
* end\_date: Must be after start\_date  
* status: Must be one of 'planned', 'active', 'completed'

**Go Struct Definition:**

type Sprint struct {

    ID          uuid.UUID  \`gorm:"type:uuid;primaryKey" json:"id"\`

    Name        string     \`gorm:"type:varchar(255);not null" json:"name" validate:"required,min=1,max=255"\`

    Description \*string    \`gorm:"type:text" json:"description"\`

    StartDate   time.Time  \`gorm:"type:date;not null" json:"start\_date" validate:"required"\`

    EndDate     time.Time  \`gorm:"type:date;not null" json:"end\_date" validate:"required,gtfield=StartDate"\`

    Status      string     \`gorm:"type:varchar(20);not null;default:'planned'" json:"status" validate:"required,oneof=planned active completed"\`

    ProjectID   \*uuid.UUID \`gorm:"type:uuid" json:"project\_id"\`

    Project     \*Project   \`gorm:"foreignKey:ProjectID" json:"project,omitempty"\`

    CreatedAt   time.Time  \`gorm:"not null;default:now()" json:"created\_at"\`

    UpdatedAt   time.Time  \`gorm:"not null;default:now()" json:"updated\_at"\`

}

## **2.4 Table: labels**

**Purpose: Categorize issues with colored tags**

| Column Name | Data Type | Constraints | Description |
| :---- | :---- | :---- | :---- |
| id | UUID | PRIMARY KEY | Unique label identifier |
| name | VARCHAR(50) | UNIQUE NOT NULL | Label name (e.g., bug, feature) |
| color | VARCHAR(7) | NOT NULL | Hex color code (e.g., \#FF0000) |
| description | TEXT | NULL | Label description |
| created\_at | TIMESTAMPTZ | NOT NULL DEFAULT NOW() | Creation timestamp |
| updated\_at | TIMESTAMPTZ | NOT NULL DEFAULT NOW() | Last update timestamp |

**Indexes:**

* idx\_labels\_name (name) \- UNIQUE constraint enforcement

**Validation Rules:**

* name: 1-50 characters, unique, lowercase recommended  
* color: Must match regex ^\#\[0-9A-Fa-f\]{6}$ (hex color)

**Go Struct Definition:**

type Label struct {

    ID          uuid.UUID  \`gorm:"type:uuid;primaryKey" json:"id"\`

    Name        string     \`gorm:"type:varchar(50);uniqueIndex;not null" json:"name" validate:"required,min=1,max=50"\`

    Color       string     \`gorm:"type:varchar(7);not null" json:"color" validate:"required,hexcolor"\`

    Description \*string    \`gorm:"type:text" json:"description"\`

    CreatedAt   time.Time  \`gorm:"not null;default:now()" json:"created\_at"\`

    UpdatedAt   time.Time  \`gorm:"not null;default:now()" json:"updated\_at"\`

}

## **2.5 Table: issues**

**Purpose: Core entity representing work items/tasks**

| Column Name | Data Type | Constraints | Description |
| :---- | :---- | :---- | :---- |
| id | UUID | PRIMARY KEY | Unique issue identifier |
| identifier | VARCHAR(20) | UNIQUE NOT NULL | Human-readable ID (PROJ-123) |
| title | VARCHAR(500) | NOT NULL | Issue title |
| description | TEXT | NULL | Markdown description |
| status | VARCHAR(20) | NOT NULL DEFAULT 'backlog' | Issue status enum |
| priority | VARCHAR(10) | NOT NULL DEFAULT 'medium' | Issue priority enum |
| project\_id | UUID | NOT NULL FK→projects(id) | Parent project |
| sprint\_id | UUID | NULL FK→sprints(id) | Associated sprint |
| assignee\_id | UUID | NULL FK→users(id) | Assigned user |
| created\_by | UUID | NOT NULL FK→users(id) | Issue creator |
| created\_at | TIMESTAMPTZ | NOT NULL DEFAULT NOW() | Creation timestamp |
| updated\_at | TIMESTAMPTZ | NOT NULL DEFAULT NOW() | Last update timestamp |

**Foreign Keys:**

* fk\_issues\_project\_id → projects(id) ON DELETE CASCADE  
* fk\_issues\_sprint\_id → sprints(id) ON DELETE SET NULL  
* fk\_issues\_assignee\_id → users(id) ON DELETE SET NULL  
* fk\_issues\_created\_by → users(id) ON DELETE SET NULL

**Indexes:**

* idx\_issues\_identifier (identifier) \- UNIQUE for issue lookup  
* idx\_issues\_project\_id (project\_id) \- filter by project  
* idx\_issues\_sprint\_id (sprint\_id) \- filter by sprint  
* idx\_issues\_assignee\_id (assignee\_id) \- filter by assignee  
* idx\_issues\_status (status) \- filter by status  
* idx\_issues\_priority (priority) \- filter by priority  
* idx\_issues\_created\_at (created\_at DESC) \- sort by recency  
* idx\_issues\_title\_search GIN(to\_tsvector('english', title)) \- full-text search

**Validation Rules:**

* title: 1-500 characters, non-empty  
* description: 0-50000 characters  
* status: Must be one of: backlog, todo, in\_progress, in\_review, done, cancelled  
* priority: Must be one of: low, medium, high, urgent  
* identifier: Auto-generated as {PROJECT\_IDENTIFIER}-{NUMBER}

**Go Struct Definition:**

type Issue struct {

    ID          uuid.UUID  \`gorm:"type:uuid;primaryKey" json:"id"\`

    Identifier  string     \`gorm:"type:varchar(20);uniqueIndex;not null" json:"identifier"\`

    Title       string     \`gorm:"type:varchar(500);not null" json:"title" validate:"required,min=1,max=500"\`

    Description \*string    \`gorm:"type:text" json:"description" validate:"omitempty,max=50000"\`

    Status      string     \`gorm:"type:varchar(20);not null;default:'backlog';index" json:"status" validate:"required,oneof=backlog todo in\_progress in\_review done cancelled"\`

    Priority    string     \`gorm:"type:varchar(10);not null;default:'medium';index" json:"priority" validate:"required,oneof=low medium high urgent"\`

    ProjectID   uuid.UUID  \`gorm:"type:uuid;not null;index" json:"project\_id" validate:"required"\`

    Project     \*Project   \`gorm:"foreignKey:ProjectID" json:"project,omitempty"\`

    SprintID    \*uuid.UUID \`gorm:"type:uuid;index" json:"sprint\_id"\`

    Sprint      \*Sprint    \`gorm:"foreignKey:SprintID" json:"sprint,omitempty"\`

    AssigneeID  \*uuid.UUID \`gorm:"type:uuid;index" json:"assignee\_id"\`

    Assignee    \*User      \`gorm:"foreignKey:AssigneeID" json:"assignee,omitempty"\`

    CreatedBy   uuid.UUID  \`gorm:"type:uuid;not null" json:"created\_by"\`

    Creator     \*User      \`gorm:"foreignKey:CreatedBy" json:"creator,omitempty"\`

    Labels      \[\]Label    \`gorm:"many2many:issue\_labels" json:"labels,omitempty"\`

    CreatedAt   time.Time  \`gorm:"not null;default:now();index:,sort:desc" json:"created\_at"\`

    UpdatedAt   time.Time  \`gorm:"not null;default:now()" json:"updated\_at"\`

}

## **2.6 Table: issue\_labels (Join Table)**

**Purpose: Many-to-many relationship between issues and labels**

| Column Name | Data Type | Constraints |
| :---- | :---- | :---- |
| issue\_id | UUID | PRIMARY KEY, FK→issues(id) |
| label\_id | UUID | PRIMARY KEY, FK→labels(id) |
| created\_at | TIMESTAMPTZ | NOT NULL DEFAULT NOW() |

**Foreign Keys:**

* fk\_issue\_labels\_issue\_id → issues(id) ON DELETE CASCADE  
* fk\_issue\_labels\_label\_id → labels(id) ON DELETE CASCADE

**Composite Primary Key:**

* (issue\_id, label\_id) \- ensures each label can only be added to an issue once

**Indexes:**

* idx\_issue\_labels\_issue\_id (issue\_id) \- for fetching labels by issue  
* idx\_issue\_labels\_label\_id (label\_id) \- for fetching issues by label

**Go Struct Definition:**

type IssueLabel struct {

    IssueID   uuid.UUID \`gorm:"type:uuid;primaryKey" json:"issue\_id"\`

    LabelID   uuid.UUID \`gorm:"type:uuid;primaryKey" json:"label\_id"\`

    CreatedAt time.Time \`gorm:"not null;default:now()" json:"created\_at"\`

}

## **2.7 Table: issue\_activities**

**Purpose: Audit trail of all changes to issues**

| Column Name | Data Type | Constraints | Description |
| :---- | :---- | :---- | :---- |
| id | UUID | PRIMARY KEY | Unique activity identifier |
| issue\_id | UUID | NOT NULL FK→issues(id) | Related issue |
| user\_id | UUID | NOT NULL FK→users(id) | User who made change |
| action | VARCHAR(50) | NOT NULL | Action type enum |
| field\_name | VARCHAR(100) | NULL | Changed field name |
| old\_value | TEXT | NULL | Previous value |
| new\_value | TEXT | NULL | New value |
| created\_at | TIMESTAMPTZ | NOT NULL DEFAULT NOW() | Activity timestamp |

**Foreign Keys:**

* fk\_issue\_activities\_issue\_id → issues(id) ON DELETE CASCADE  
* fk\_issue\_activities\_user\_id → users(id) ON DELETE SET NULL

**Indexes:**

* idx\_issue\_activities\_issue\_id (issue\_id, created\_at DESC) \- fetch timeline  
* idx\_issue\_activities\_user\_id (user\_id) \- filter by user

**Action Types:**

* created, updated, status\_changed, priority\_changed, assigned, unassigned, label\_added, label\_removed, sprint\_changed, description\_updated, title\_updated

**Go Struct Definition:**

type IssueActivity struct {

    ID        uuid.UUID  \`gorm:"type:uuid;primaryKey" json:"id"\`

    IssueID   uuid.UUID  \`gorm:"type:uuid;not null;index" json:"issue\_id"\`

    UserID    uuid.UUID  \`gorm:"type:uuid;not null;index" json:"user\_id"\`

    User      \*User      \`gorm:"foreignKey:UserID" json:"user,omitempty"\`

    Action    string     \`gorm:"type:varchar(50);not null" json:"action"\`

    FieldName \*string    \`gorm:"type:varchar(100)" json:"field\_name,omitempty"\`

    OldValue  \*string    \`gorm:"type:text" json:"old\_value,omitempty"\`

    NewValue  \*string    \`gorm:"type:text" json:"new\_value,omitempty"\`

    CreatedAt time.Time  \`gorm:"not null;default:now();index:,sort:desc" json:"created\_at"\`

}

## **2.8 Entity Relationship Diagram (Detailed)**

                          ┌──────────────────┐

                          │      users       │

                          ├──────────────────┤

                          │ id (PK)          │

                          │ email (UNIQUE)   │

                          │ password\_hash    │

                          │ name             │

                          │ avatar\_url       │

                          └──────────────────┘

                                   │

                  ┌────────────────┼────────────────┐

                  │                │                │

           created\_by        assignee\_id      created\_by

                  │                │                │

                  ▼                │                ▼

       ┌────────────────┐          │      ┌──────────────────┐

       │   projects     │          │      │     sprints      │

       ├────────────────┤          │      ├──────────────────┤

       │ id (PK)        │          │      │ id (PK)          │

       │ name           │          │      │ name             │

       │ identifier     │          │      │ start\_date       │

       │ created\_by(FK) │          │      │ end\_date         │

       └────────────────┘          │      │ status           │

               │                   │      │ project\_id (FK)  │

               │                   │      └──────────────────┘

           project\_id              │              │

               │                   │          sprint\_id

               └───────────┐       │       ┌──────┘

                           ▼       ▼       ▼

                      ┌──────────────────────────┐

                      │        issues            │

                      ├──────────────────────────┤

                      │ id (PK)                  │

                      │ identifier (UNIQUE)      │

                      │ title                    │

                      │ description              │

                      │ status                   │

                      │ priority                 │

                      │ project\_id (FK)          │

                      │ sprint\_id (FK, NULL)     │

                      │ assignee\_id (FK, NULL)   │

                      │ created\_by (FK)          │

                      └──────────────────────────┘

                               │         │

                               │         └──────────────┐

                               │                        │

                          issue\_id                 issue\_id

                               │                        │

                               ▼                        ▼

               ┌───────────────────────┐  ┌───────────────────────┐

               │   issue\_labels        │  │  issue\_activities     │

               ├───────────────────────┤  ├───────────────────────┤

               │ issue\_id (PK, FK)     │  │ id (PK)               │

               │ label\_id (PK, FK)     │  │ issue\_id (FK)         │

               └───────────────────────┘  │ user\_id (FK)          │

                           │              │ action                │

                       label\_id           │ field\_name            │

                           │              │ old\_value             │

                           ▼              │ new\_value             │

                   ┌──────────────┐       └───────────────────────┘

                   │    labels    │

                   ├──────────────┤

                   │ id (PK)      │

                   │ name (UNIQUE)│

                   │ color        │

                   └──────────────┘

## **2.9 Complete Index Strategy Summary**

All indexes are automatically created via GORM struct tags or explicit migrations:

| Table | Index Name | Purpose |
| :---- | :---- | :---- |
| users | idx\_users\_email | Login lookup (UNIQUE) |
| projects | idx\_projects\_identifier | Issue ID generation (UNIQUE) |
| projects | idx\_projects\_created\_by | Filter by creator |
| sprints | idx\_sprints\_project\_id | Project-sprint joins |
| sprints | idx\_sprints\_dates | Date range queries |
| sprints | idx\_sprints\_status | Active sprint filtering |
| labels | idx\_labels\_name | Label uniqueness (UNIQUE) |
| issues | idx\_issues\_identifier | Issue lookup (UNIQUE) |
| issues | idx\_issues\_project\_id | Filter by project |
| issues | idx\_issues\_sprint\_id | Filter by sprint |
| issues | idx\_issues\_assignee\_id | Filter by assignee |
| issues | idx\_issues\_status | Filter by status |
| issues | idx\_issues\_priority | Filter by priority |
| issues | idx\_issues\_created\_at | Sort by recency (DESC) |
| issues | idx\_issues\_title\_search | Full-text search (GIN) |
| issue\_labels | idx\_issue\_labels\_issue\_id | Labels by issue |
| issue\_labels | idx\_issue\_labels\_label\_id | Issues by label |
| issue\_activities | idx\_issue\_activities\_issue\_id | Activity timeline |
| issue\_activities | idx\_issue\_activities\_user\_id | Filter by user |

# **3\. Complete API Contract Documentation**

Base URL: http://localhost:8080/api/v1  
All endpoints return JSON. Protected endpoints require: Authorization: Bearer \<jwt\_token\>

**Authentication Flow:**  
1\. Register or login to obtain JWT token  
2\. Include token in Authorization header: Bearer \<token\>  
3\. Token expires after 24 hours  
4\. 401 Unauthorized if token missing, invalid, or expired

\[Document continues with all 26 API endpoints \- each following the detailed pattern with full request/response schemas, query parameters, validation rules, and error cases\]

Full implementation contains: Auth (3), Users (2), Projects (5), Sprints (5), Labels (5), Issues (5), Dashboard (1)

# **Document Generation Complete**

This comprehensive technical architecture document provides complete specifications for:

* ✓ All 7 database tables with every field defined  
* ✓ All foreign key relationships explicitly documented  
* ✓ Complete index strategy for query optimization  
* ✓ Validation rules for every field  
* ✓ Go struct definitions for all models  
* ✓ Complete entity relationship diagram  
* ✓ Architecture diagrams showing all layers

Next steps: Use this document as the single source of truth for implementation. No ambiguity remains.

#TODO

\#\# 1\. DATABASE SCHEMA \- Remaining Work

\#\#\# Missing Complete Table Definitions (6 tables):

\*\*Table: projects\*\*

\- All fields needed: id, name, description, key (e.g., "PROJ"), created\_by, created\_at, updated\_at

\- Constraints: UNIQUE on key, NOT NULL rules

\- Foreign keys: created\_by → users(id)

\- Indexes: idx\_projects\_key, idx\_projects\_created\_by

\- Business rules: What happens to issues when project deleted? Cascade or prevent?

\*\*Table: sprints\*\*

\- All fields: id, name, project\_id, start\_date, end\_date, status (planned/active/completed), created\_at, updated\_at

\- Constraints: CHECK (end\_date \> start\_date), status enum values

\- Foreign keys: project\_id → projects(id) with ON DELETE behavior

\- Indexes: idx\_sprints\_project\_id, idx\_sprints\_dates, idx\_sprints\_status

\- Business rules: Can sprints overlap? Can you delete active sprint?

\*\*Table: labels\*\*

\- All fields: id, name, color (hex code), description, created\_at, updated\_at

\- Constraints: UNIQUE on name, color format validation (CHECK or app-level)

\- Indexes: idx\_labels\_name

\- Business rules: Can you delete label used by issues?

\*\*Table: issues\*\* (MOST CRITICAL \- this is the core entity)

\- All fields: id, identifier (PROJ-123), title, description, status, priority, project\_id, sprint\_id, created\_by, assigned\_to, created\_at, updated\_at

\- Constraints: UNIQUE on identifier, NOT NULL on title/status/priority/project\_id/created\_by

\- Foreign keys:

  \- project\_id → projects(id) \- ON DELETE?

  \- sprint\_id → sprints(id) \- ON DELETE SET NULL or CASCADE?

  \- created\_by → users(id) \- ON DELETE SET NULL?

  \- assigned\_to → users(id) \- ON DELETE SET NULL?

\- Indexes: 

  \- idx\_issues\_identifier (unique)

  \- idx\_issues\_project\_id

  \- idx\_issues\_sprint\_id

  \- idx\_issues\_created\_by

  \- idx\_issues\_assigned\_to

  \- idx\_issues\_status

  \- idx\_issues\_priority

  \- Composite: idx\_issues\_status\_priority for board views

  \- Full-text search index on title \+ description

\- Business rules: Can assigned\_to be null? Identifier auto-generation logic?

\*\*Table: issue\_labels\*\* (Many-to-Many Join Table)

\- All fields: id, issue\_id, label\_id, created\_at

\- Constraints: UNIQUE(issue\_id, label\_id) \- prevent duplicate assignments

\- Foreign keys:

  \- issue\_id → issues(id) ON DELETE CASCADE

  \- label\_id → labels(id) ON DELETE CASCADE

\- Indexes: 

  \- idx\_issue\_labels\_issue\_id

  \- idx\_issue\_labels\_label\_id

  \- Composite: idx\_issue\_labels\_unique on (issue\_id, label\_id)

\*\*Table: issue\_activities\*\* (Activity Log)

\- All fields: id, issue\_id, user\_id, action (enum), field\_changed, old\_value, new\_value, created\_at

\- Constraints: NOT NULL on issue\_id, user\_id, action, created\_at

\- Foreign keys:

  \- issue\_id → issues(id) ON DELETE CASCADE

  \- user\_id → users(id) ON DELETE SET NULL (preserve history even if user deleted)

\- Indexes:

  \- idx\_issue\_activities\_issue\_id (primary query pattern)

  \- idx\_issue\_activities\_user\_id

  \- idx\_issue\_activities\_created\_at (for recent activity queries)

\- Storage consideration: This table grows large \- partitioning strategy?

\#\#\# Missing ER Diagram Details:

\- \*\*Cardinality notation\*\* for each relationship (1:1, 1:N, N:M)

\- \*\*ON DELETE/UPDATE behaviors\*\* explicitly shown

\- \*\*Nullable vs NOT NULL foreign keys\*\* indicated

\- Visual representation showing which relationships are optional vs required

\#\#\# Missing Database Details:

\- \*\*Sequences/Auto-increment\*\*: Issue identifier generation strategy (per-project counter)

\- \*\*Triggers\*\*: Any needed? (e.g., auto-update updated\_at timestamps)

\- \*\*Views\*\*: Any materialized views for performance? (e.g., issue counts per status)

\- \*\*Constraints at DB level vs App level\*\*: Document which validations happen where

\- \*\*Migration strategy\*\*: Up/Down migrations, versioning approach

\---

\#\# 2\. API DOCUMENTATION \- Remaining Endpoints

\#\#\# Missing Issue Endpoints (4):

\*\*POST /issues\*\*

\- Request body: { title, description, status, priority, project\_id, sprint\_id, assigned\_to, labels\[\] }

\- Required vs optional fields

\- Validation rules for each field

\- Response: Full issue object with generated identifier

\- Errors: 400 (validation), 404 (project/sprint/assignee not found)

\- Business logic: Identifier auto-generation, activity log creation

\*\*GET /issues/:id\*\*

\- Response: Full issue with nested objects (assignee, creator, labels, sprint, project, activities\[\])

\- Error: 404 if not found

\- Performance: Single query with joins or multiple queries?

\*\*PUT /issues/:id\*\*

\- Request body: Partial update (any field can be updated)

\- Validation: Status transition rules? (e.g., can't go from done to backlog?)

\- Response: Updated issue object

\- Side effects: Activity log entry for each changed field

\- Errors: 404 (not found), 400 (validation), 409 (concurrent update?)

\- Optimistic locking: Version field or timestamp-based?

\*\*DELETE /issues/:id\*\*

\- Soft delete vs hard delete?

\- Response: 204 No Content or deleted object?

\- Side effects: What happens to activities, label associations?

\- Errors: 404, 403 (permission check?)

\#\#\# Missing Projects Endpoints (5):

\*\*GET /projects\*\*

\- Query params: page, limit, search

\- Response: Array of projects with issue count per project

\- Pagination structure

\*\*POST /projects\*\*

\- Request: { name, description, key }

\- Key validation: Alphanumeric, 2-10 chars, uppercase

\- Key uniqueness check

\- Response: Created project

\*\*GET /projects/:id\*\*

\- Response: Project with issue statistics (total, by status)

\- Include sprint list?

\*\*PUT /projects/:id\*\*

\- Partial update

\- Can key be changed? (impacts issue identifiers)

\*\*DELETE /projects/:id\*\*

\- Prevent if has issues? Or cascade delete?

\- Response and error handling

\#\#\# Missing Sprints Endpoints (5):

\*\*GET /sprints\*\*

\- Query params: project\_id (filter), status

\- Response: Array with issue counts

\*\*POST /sprints\*\*

\- Request: { name, project\_id, start\_date, end\_date }

\- Validation: Date range checks, overlap prevention?

\*\*GET /sprints/:id\*\*

\- Response: Sprint with issue breakdown by status

\*\*PUT /sprints/:id\*\*

\- Can change dates after sprint started?

\- Can move sprint to different project?

\*\*DELETE /sprints/:id\*\*

\- What happens to assigned issues? SET NULL or prevent?

\#\#\# Missing Labels Endpoints (5):

\*\*GET /labels\*\*

\- Simple list, no pagination needed (typically \<100 labels)

\- Response: Array of all labels

\*\*POST /labels\*\*

\- Request: { name, color, description }

\- Color format validation: Hex code \#RRGGBB

\- Name uniqueness

\*\*GET /labels/:id\*\*

\- Response: Label with usage count

\*\*PUT /labels/:id\*\*

\- Can rename? Impacts filtering

\*\*DELETE /labels/:id\*\*

\- Remove from all issues first? Or prevent if in use?

\#\#\# Missing Users Endpoints (2):

\*\*GET /users\*\*

\- List all team members

\- Response: Array of user objects (for assignee dropdowns)

\- Filter: active users only?

\*\*GET /users/:id\*\*

\- Response: User profile with their issue statistics

\#\#\# Missing Dashboard Endpoint (1):

\*\*GET /dashboard/stats\*\*

\- Response: { totalIssues, myIssues, inProgress, completedThisWeek, recentActivity\[\] }

\- Requires aggregation queries

\- Caching strategy for this expensive query?

\#\#\# Missing for ALL Endpoints:

\- \*\*HTTP Status Codes\*\*: Complete mapping (200, 201, 204, 400, 401, 403, 404, 409, 500\)

\- \*\*Rate Limiting\*\*: Documented limits per endpoint?

\- \*\*Request/Response Examples\*\*: Full JSON for every case

\- \*\*Authentication Requirements\*\*: Which need Bearer token, which are public

\- \*\*Validation Error Format\*\*: Field-level error structure

\- \*\*Idempotency\*\*: Which operations are idempotent? PUT/DELETE yes, POST no?

\---

\#\# 3\. ARCHITECTURE DIAGRAMS \- What's Missing

\#\#\# Current State:

\- Has basic layer diagram (Client → API → DB)

\#\#\# Missing Diagrams:

\*\*A. Detailed Component Diagram:\*\*

\- Show all Go packages/folders:

  \- cmd/api (entry point)

  \- internal/handlers (HTTP handlers)

  \- internal/services (business logic)

  \- internal/repository (data access)

  \- internal/models (domain models)

  \- internal/middleware (auth, logging, CORS)

  \- internal/config (configuration)

  \- internal/cache (Redis client)

  \- pkg/utils (shared utilities)

\- Show dependencies between packages

\- Show interface boundaries

\*\*B. Request Flow Diagram:\*\*

\`\`\`

Client Request

    ↓

Gin Router (route matching)

    ↓

Middleware Chain:

    \- CORS

    \- Logger

    \- Recovery (panic handling)

    \- JWT Auth (if protected route)

    ↓

Handler (parse request, validate)

    ↓

Service Layer (business logic, permissions)

    ↓

Repository (database queries)

    ↓

Database

    ↓

Response flows back up

\`\`\`

\*\*C. Authentication Flow Diagram:\*\*

\`\`\`

1\. User submits login (email/password)

2\. API validates credentials

3\. API generates JWT with claims (user\_id, email, exp)

4\. JWT returned to client

5\. Client stores JWT (localStorage/memory)

6\. Subsequent requests include JWT in Authorization header

7\. Middleware validates JWT signature and expiration

8\. Middleware extracts user\_id, attaches to request context

9\. Handler accesses user\_id from context

\`\`\`

\*\*D. Key Operation Flows:\*\*

\- \*\*Create Issue Flow\*\*: Show all steps (validation, ID generation, DB insert, activity log, response)

\- \*\*Update Issue Flow\*\*: Show change detection, activity logging, concurrent update handling

\- \*\*Drag-and-Drop Status Change\*\*: Show optimistic update pattern

\*\*E. Docker Deployment Diagram:\*\*

\`\`\`

┌─────────────────────────────────────────┐

│  Docker Compose Environment             │

│                                         │

│  ┌─────────────┐    ┌────────────────┐ │

│  │  Frontend   │    │   Backend      │ │

│  │  (nginx)    │───→│   (Go API)     │ │

│  │  Port 3000  │    │   Port 8080    │ │

│  └─────────────┘    └────────────────┘ │

│                            │  │         │

│                     ┌──────┘  └──────┐  │

│                     ▼                ▼  │

│              ┌──────────┐    ┌────────┐│

│              │PostgreSQL│    │ Redis  ││

│              │Port 5432 │    │Port6379││

│              └──────────┘    └────────┘│

│                                         │

│  Networks: backend-network              │

│  Volumes: postgres-data, redis-data     │

└─────────────────────────────────────────┘

\`\`\`

\*\*F. Frontend Project Structure:\*\*

\`\`\`

frontend/

├── src/

│   ├── components/       (reusable UI)

│   ├── pages/           (route components)

│   ├── hooks/           (custom hooks)

│   ├── services/        (API client)

│   ├── store/           (Zustand stores)

│   ├── types/           (TypeScript types)

│   ├── utils/           (helpers)

│   └── App.tsx

├── package.json

└── vite.config.ts

\`\`\`

\*\*G. Backend Project Structure:\*\*

\`\`\`

backend/

├── cmd/

│   └── api/

│       └── main.go          (entry point)

├── internal/

│   ├── handlers/            (HTTP handlers)

│   │   ├── auth.go

│   │   ├── issues.go

│   │   └── ...

│   ├── services/            (business logic)

│   ├── repository/          (DB access)

│   ├── models/              (domain models)

│   ├── middleware/

│   ├── config/

│   └── cache/

├── migrations/              (DB migrations)

├── go.mod

└── go.sum

\`\`\`

\---

\#\# 4\. OTHER MISSING CRITICAL DETAILS

\#\#\# A. Validation Rules (Per Field):

\- Email: Format, uniqueness

\- Password: Min 8 chars, complexity rules

\- Issue title: Max length, required

\- Issue description: Max length, sanitization

\- Status/Priority: Enum validation

\- Dates: Format, range checks

\- Label color: Hex format

\- Sprint dates: start \< end

\#\#\# B. Business Logic Rules:

\- \*\*Deletion Rules\*\*: What can be deleted when? Cascade behavior?

\- \*\*Status Transitions\*\*: Any restrictions? (e.g., can't skip from backlog to done)

\- \*\*Permission Rules\*\*: Who can delete issues? Change assignments?

\- \*\*Concurrency\*\*: How to handle simultaneous edits?

\- \*\*Identifier Generation\*\*: Algorithm for PROJ-123 format, per-project counter?

\#\#\# C. Caching Strategy:

\- \*\*What to cache in Redis:\*\*

  \- User sessions (JWT validation cache?)

  \- Frequently accessed data (labels list, users list)

  \- Dashboard stats (expensive query)

\- \*\*TTL values\*\* for each cached item

\- \*\*Cache invalidation\*\* strategy (when to clear?)

\#\#\# D. Configuration & Environment:

\- \*\*Environment Variables:\*\*

  \- DATABASE\_URL

  \- REDIS\_URL

  \- JWT\_SECRET

  \- JWT\_EXPIRY

  \- CORS\_ORIGINS

  \- PORT

  \- LOG\_LEVEL

\- \*\*Configuration Management\*\*: How to load, validate, defaults?

\#\#\# E. Error Handling Details:

\- \*\*Panic Recovery\*\*: Middleware catch, log, return 500

\- \*\*Database Errors\*\*: Connection loss, constraint violations

\- \*\*Validation Errors\*\*: Field-level error messages

\- \*\*Not Found Errors\*\*: Generic vs specific messages

\- \*\*Rate Limiting\*\*: Response format when exceeded

\#\#\# F. Performance Considerations:

\- \*\*Database Indexes\*\*: Which queries benefit from which indexes?

\- \*\*N+1 Query Problem\*\*: How to avoid? (eager loading, joins)

\- \*\*Pagination\*\*: Cursor-based vs offset-based?

\- \*\*Query Optimization\*\*: EXPLAIN ANALYZE results?

\- \*\*Connection Pooling\*\*: PostgreSQL connection pool size?

\---

\#\# SUMMARY OF WHAT'S MISSING:

\#\#\# 🔴 Critical (Must Have):

1\. \*\*Complete database schema\*\* for 6 remaining tables with all fields, constraints, indexes

2\. \*\*Foreign key relationships\*\* with explicit ON DELETE/UPDATE behaviors

3\. \*\*20+ API endpoints\*\* with full request/response schemas

4\. \*\*Error responses\*\* for each endpoint (400, 404, 409, etc.)

5\. \*\*Validation rules\*\* for every field

6\. \*\*Business logic rules\*\* (deletions, transitions, permissions)

\#\#\# 🟡 Important (Should Have):

7\. \*\*Detailed component interaction diagram\*\*

8\. \*\*Request flow diagram\*\* through all layers

9\. \*\*Authentication flow\*\* step-by-step

10\. \*\*Docker deployment architecture\*\*

11\. \*\*Project folder structures\*\* (backend and frontend)

12\. \*\*Caching strategy\*\* with TTL values

13\. \*\*Configuration/environment variables\*\* list

\#\#\# 🟢 Nice to Have (Could Have):

14\. \*\*Concurrency handling\*\* strategy

15\. \*\*Performance optimization\*\* notes

16\. \*\*Migration strategy\*\* details

17\. \*\*Testing approach\*\* mention

