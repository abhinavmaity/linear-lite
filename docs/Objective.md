# **Summary**

Linear-lite is a streamlined issue tracking and project management application designed for small to medium-sized development teams. It provides essential functionality for managing tasks, organizing sprints, and tracking project progress without the complexity of enterprise-level tools.

The application focuses on core workflows that teams use daily: creating and managing issues, organizing work into sprints, filtering and searching tasks, and visualizing progress through different views. By limiting the scope to essential features, Linear-lite delivers a clean, fast, and intuitive experience that teams can adopt immediately.

# **Product Overview**

## **Purpose**

Linear-lite addresses the need for a lightweight, easy-to-use issue tracker that balances simplicity with functionality. Many existing tools are either too basic for team collaboration or too complex for quick adoption. Linear-lite fills this gap by providing recognizable patterns from tools like Linear and Jira while maintaining a minimal feature set.

## **Target Users**

* Small development teams (3-15 people)  
* Startup engineering teams needing quick project setup  
* Side project teams and open-source maintainers  
* Teams wanting self-hosted or locally-run project management

## **Key Principles**

* **Essential features only:** Focus on 80% of daily workflows  
* **Fast and lightweight:** Quick load times, minimal clicks  
* **Recognizable patterns:** Familiar to users of Linear, Jira, or GitHub Issues  
* **Self-hostable:** Can run locally or be deployed independently

# **Core Features**

## **1\. Issue Management**

### **Create Issues**

* Quick create with title and description  
* Automatic unique identifier generation (e.g., PROJ-123)  
* Rich text description with Markdown support  
* Set initial status, priority, and assignee during creation

### 

### **Edit Issues**

* Update title, description, and all attributes  
* Change status (Backlog, Todo, In Progress, In Review, Done, Cancelled)  
* Set priority (Low, Medium, High, Urgent)  
* Assign/reassign to team members  
* Add/remove labels

### **Delete Issues**

* Soft delete with confirmation  
* Archive functionality for completed issues

### **Issue Details**

* Activity timeline showing all changes  
* Creation and last updated timestamps  
* Creator and current assignee information

## **2\. Labels & Organization**

### **Label System**

* Create custom labels with name and color  
* Apply multiple labels to any issue  
* Common label types: bug, feature, enhancement, documentation  
* Label management: edit name/color, delete unused labels

## **3\. Sprint & Project Organization**

### **Projects**

* Create projects with name and description  
* Assign issues to projects  
* Project-level view showing all related issues

### **Sprints/Milestones**

* Create sprints with name, start date, and end date  
* Assign issues to specific sprints  
* View sprint progress (total, completed, in progress)  
* Active sprint indicator  
* Close/archive completed sprints

## **4\. Multiple Views**

### **List View**

* Table-style view with all issues  
* Columns: ID, Title, Status, Priority, Assignee, Labels, Sprint  
* Sortable columns  
* Quick inline status updates

### **Board View (Kanban)**

* Columns for each status (Backlog, Todo, In Progress, In Review, Done)  
* Drag-and-drop issues between columns to update status  
* Card view showing: ID, title, assignee, priority, labels  
* Vertical scrolling within columns

## **5\. Filtering & Search**

### **Filter Capabilities**

* Filter by status (single or multiple)  
* Filter by priority  
* Filter by assignee (including unassigned)  
* Filter by labels (AND/OR logic)  
* Filter by sprint/project  
* Combine multiple filters

### **Search**

* Full-text search across issue titles and descriptions  
* Search by issue ID (e.g., PROJ-123)  
* Real-time search results

## **6\. User & Team Management**

### **User Accounts**

* User registration with email and password  
* User login with session management  
* Basic profile: name, email, avatar (optional)

### **Team Management**

* View all team members  
* Assign issues to any team member  
* Filter issues by specific assignees

# **User Roles & Permissions**

| Role | Capabilities | MVP Priority |
| :---- | :---- | :---- |
| **Team Member** | Create, edit, and delete own issues; view all issues; comment on issues; update issue status | High |
| **Admin** | All team member capabilities plus: manage projects, sprints, labels; manage team members; delete any issue | Medium |

Note: For MVP, permissions can be simplified to allow all authenticated users to perform all actions. Formal role-based access control can be added post-MVP.

# **Key User Stories**

### **As a developer, I want to...**

* Quickly create a bug report with description and steps to reproduce  
* See all issues assigned to me in one view  
* Move issues across status columns as I work on them  
* Filter issues by labels to find related work  
* See what's in the current sprint and track progress

### **As a project manager, I want to...**

* Create and manage sprints with clear start and end dates  
* Organize issues by priority to guide the team  
* View sprint progress and identify blockers  
* Search for specific issues by ID or keywords  
* See who is working on what across the team

# **Success Metrics**

### **MVP Success Criteria**

* Complete CRUD operations for issues, projects, sprints, and labels  
* Functional board view with drag-and-drop  
* Working filters and search functionality  
* User authentication and session management  
* Deployable via Docker with single command  
* Responsive design that works on desktop and tablet

### **Performance Targets**

* Page load time under 2 seconds  
* Issue creation/update under 500ms  
* Search results return under 1 second for 1000+ issues

# **Out of Scope for MVP**

The following features are intentionally excluded from the MVP to maintain focus.

* Comments and discussions on issues  
* File attachments and image uploads  
* Email notifications  
* Real-time collaboration (WebSocket updates)  
* Time tracking and estimates  
* Issue dependencies and blocking relationships  
* Advanced reporting and analytics  
* Bulk operations (bulk edit, bulk delete)  
* Custom fields or issue templates  
* GitHub/GitLab integration  
* Mobile apps (iOS/Android)  
* API for external integrations /Multi-workspace or multi-organization support

# **Implementation Priorities**

1. ### **Foundation**

* User authentication and authorization  
* Database schema and models  
* Basic issue CRUD operations

2. ### **Core Functionality**

* Labels, priorities, and assignments  
* Projects and sprints management  
* Advanced filtering logic  
* Full-text search implementation

3. ### **Views and UX**

* List view with sorting  
* Kanban board view with drag-and-drop  
* Activity timeline  
* Responsive design

4. ### **Polish and Deployment**

* UI/UX refinements  
* Performance optimization  
* Docker configuration  
* Documentation and README  
* Bug fixes and testing