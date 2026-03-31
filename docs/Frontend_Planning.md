## **Screen Inventory**

**Total Screens: 11**

| \# | Screen Name | Purpose |
| :---- | :---- | :---- |
| **1** | **Login** | User authentication entry point |
| **2** | **Register** | New user account creation |
| **3** | **Dashboard** | Main landing page after login with overview stats |
| **4** | **Issues List View** | Table view of all issues with filtering and search |
| **5** | **Issues Board View** | Kanban board with drag-and-drop functionality |
| **6** | **Issue Detail Page** | Full issue details with edit capabilities and activity timeline |
| **7** | **Create Issue Modal** | Quick create form overlay |
| **8** | **Projects Page** | List and manage all projects |
| **9** | **Sprints Page** | View and manage sprints/milestones |
| **10** | **Labels Management** | Create, edit, and delete labels |
| **11** | **Team Page** | View all team members and their details |

## **Main User Journeys**

### **Journey 1: First-Time User Onboarding**

1. **Landing on app →** Redirected to Login screen  
2. **Click 'Sign Up' →** Register screen  
3. **Fill form →** Submit registration  
4. **Auto-login →** Dashboard with empty state  
5. **Click 'Create Issue' →** Create Issue Modal  
6. **Create first issue →** Issue List View with new issue

### **Journey 2: Daily Developer Workflow**

7. **Login →** Dashboard shows 'My Issues' summary  
8. **Navigate to Board View →** See issues organized by status  
9. **Drag issue from 'Todo' to 'In Progress' →** Status updates  
10. **Click issue card →** Issue Detail Page opens  
11. **Edit description, add labels →** Changes saved, activity logged  
12. **Go back to board →** Mark issue done by dragging to 'Done'

### **Journey 3: Sprint Planning**

13. **Navigate to Sprints Page →** View all sprints  
14. **Click 'Create Sprint' →** Modal opens  
15. **Fill sprint details →** Sprint created  
16. **Go to Issue List →** Filter by 'No Sprint'  
17. **Open issue →** Assign to new sprint  
18. **Return to Sprints →** View sprint progress

## **Detailed Screen Breakdown**

### **Screen 1: Login**

| Route | /login |
| :---- | :---- |
| **Components** | Email input, Password input, Login button, 'Sign Up' link, Error message display |
| **API Calls** | POST /api/v1/auth/login |
| **Data Needed** | None (user inputs email and password) |
| **Success Action** | Store JWT token, redirect to /dashboard |
| **Error States** | • Invalid credentials (401) \- Display: 'Invalid email or password' • Network error \- Display: 'Connection failed. Please try again.' • Empty fields \- Client-side validation before submit |
| **Edge Cases** | • User already logged in \- redirect to dashboard • Token still valid in localStorage \- skip login • Multiple rapid login attempts \- show loading state, disable button |

### **Screen 2: Register**

| Route | /register |
| :---- | :---- |
| **Components** | Name input, Email input, Password input, Confirm password input, Sign Up button, 'Login' link |
| **API Calls** | POST /api/v1/auth/register |
| **Validation** | • Email format validation • Password min 8 characters • Password and confirm password match • Name not empty |
| **Error States** | • Email already exists (409) \- 'This email is already registered' • Validation errors \- Show inline error messages • Server error (500) \- 'Registration failed. Please try again.' |

### 

### **Screen 3: Dashboard**

| Route | /dashboard (default landing after login) |
| :---- | :---- |
| **Components** | • Welcome header with user name • Stats cards: Total Issues, My Issues, In Progress, Done This Week • Recent activity feed (last 10 updates) • Quick action: 'Create Issue' button • Active sprint indicator |
| **API Calls** | • GET /api/v1/auth/me (user info) • GET /api/v1/dashboard/stats • GET /api/v1/issues?assignee=me\&limit=10 (recent activity) |
| **Data Structure** | Stats: { totalIssues, myIssues, inProgress, doneThisWeek } Recent: Array of issue objects with: id, title, status, updated\_at |
| **Empty States** | • No issues: Show onboarding message 'Get started by creating your first issue' • No activity: Show 'No recent activity' |

### **Screen 4: Issues List View**

| Route | /issues |
| :---- | :---- |
| **Components** | • Search bar (debounced) • Filter dropdowns: Status, Priority, Assignee, Labels, Sprint • 'Create Issue' button • View toggle: List / Board • Issues table: ID, Title, Status, Priority, Assignee, Labels, Sprint • Sortable column headers • Pagination controls |
| **API Calls** | • GET /api/v1/issues?page=1\&limit=50\&search=...\&status=...\&priority=...\&assignee=...\&labels=...\&sprint=... • GET /api/v1/users (for assignee dropdown) • GET /api/v1/labels (for labels dropdown) • GET /api/v1/sprints (for sprint dropdown) |
| **User Actions** | • Click row → Navigate to Issue Detail • Click status badge → Quick update status (inline) • Sort by column → Re-fetch with sort param • Apply filters → Re-fetch with filter params • Search input → Debounced search re-fetch |
| **Error States** | • No results found \- Show empty state with 'No issues match your filters' • API error \- Show error banner with retry button • Loading state \- Show skeleton table rows |
| **Performance** | • Debounce search input (500ms delay) • Cache filter dropdown data (labels, users, sprints) • Virtual scrolling if \>100 issues per page • URL persistence \- filters/search in query params |

### **Screen 5: Issues Board View (Kanban)**

| Route | /board |
| :---- | :---- |
| **Components** | • Same filters as List View • 5 columns: Backlog, Todo, In Progress, In Review, Done • Issue cards: ID, Title, Priority badge, Assignee avatar, Labels • Drag-and-drop enabled (react-beautiful-dnd or @dnd-kit) • Issue count per column |
| **API Calls** | • GET /api/v1/issues (grouped by status on frontend) • PUT /api/v1/issues/:id (when dragging to update status) |
| **User Actions** | • Drag card between columns → Optimistic update \+ API call • Click card → Navigate to Issue Detail • Apply filters → Re-fetch and regroup |
| **Error States** | • Drag fails (network error) → Revert card position, show toast error • Empty columns → Show 'No issues in this status' |
| **Critical Logic** | • Optimistic update: Move card immediately in UI before API response • On success: Update local cache • On failure: Revert card position, show error • Disable dragging while request is in-flight |

### **Screen 6: Issue Detail Page**

| Route | /issues/:id |
| :---- | :---- |
| **Layout** | Left Panel (70%): Title, Description editor, Activity timeline Right Sidebar (30%): Status, Priority, Assignee, Labels, Sprint, Project, Dates |
| **Components** | • Editable title (inline) • Markdown editor for description • Dropdowns for all attributes • Activity timeline (chronological log) • Delete button (with confirmation) • Back button |
| **API Calls** | • GET /api/v1/issues/:id (full issue details \+ activity) • PUT /api/v1/issues/:id (on any field update) • DELETE /api/v1/issues/:id (on delete) • GET /api/v1/labels, /users, /sprints, /projects (for dropdowns) |
| **Data Needed** | Full Issue: { id, title, description, status, priority, assignee, labels, sprint, project, created\_at, updated\_at, creator } Activity: Array of { timestamp, user, action, field\_changed, old\_value, new\_value } |
| **Error States** | • Issue not found (404) → Show 'Issue not found' page with back button • Update fails → Revert field, show error toast • Delete fails → Keep page open, show error |
| **Auto-save** | • Debounce title/description edits (1 second) • Show 'Saving...' indicator • Dropdown changes save immediately • Show 'Saved' confirmation briefly |

### **Screen 7: Create Issue Modal**

| Trigger | 'Create Issue' button anywhere, or keyboard shortcut 'C' |
| :---- | :---- |
| **Fields** | • Title (required) • Description (optional) • Status (defaults to 'Backlog') • Priority (defaults to 'Medium') • Assignee (optional, defaults to current user) • Labels (optional) • Sprint (optional) • Project (optional) |
| **API Call** | POST /api/v1/issues |
| **Success Action** | • Close modal • Show success toast 'Issue created' • Navigate to Issue Detail page OR stay on current page and refresh list |
| **Validation** | • Title required \- Show inline error if empty • Disable submit button until title is filled |

### **Screens 8-11: Summary**

The remaining screens (Projects, Sprints, Labels Management, Team) follow similar patterns to the Issue List View:

* **Projects Page:** List of projects with CRUD operations  
* **Sprints Page:** List of sprints with dates, status, and issue count  
* **Labels Management:** Grid of color-coded labels with edit/delete  
* **Team Page:** List of team members (read-only for MVP)

## 

## 

## **Global Patterns & Components**

### **Navigation**

* Sidebar: Dashboard, Issues (List/Board toggle), Projects, Sprints, Labels, Team  
* Top bar: Search (global), Create Issue button, User menu (Profile, Logout)  
* Keyboard shortcuts: C (create issue), / (focus search), Esc (close modal)

### **Shared Components**

* **IssueCard:** Used in Board View and Dashboard  
* **StatusBadge:** Color-coded status indicator  
* **PriorityIcon:** Icon for priority levels  
* **UserAvatar:** Profile picture or initials  
* **LabelPill:** Colored label tag  
* **ConfirmDialog:** Reusable delete confirmation  
* **Toast:** Success/error notifications

## **Global Error Handling Strategy**

### **Network Errors**

* **Connection timeout:** Show toast 'Connection timeout. Please try again.'  
* **Server unavailable:** Show banner 'Service temporarily unavailable'  
* **Rate limiting:** Show 'Too many requests. Please wait.'

### **Authentication Errors**

* **Token expired (401):** Clear token, redirect to login  
* **Unauthorized (403):** Show 'You don't have permission for this action'

### **Data Errors**

* **Resource not found (404):** Show appropriate empty state or 'Not found' page  
* **Validation errors (400):** Show inline field errors from API response  
* **Conflict (409):** e.g., 'This label already exists'

## **Critical Edge Cases to Handle**

### **Concurrent Updates**

* User A and User B edit same issue simultaneously  
* **Solution:** Last write wins (for MVP). Show warning: 'Issue was updated by another user'

### **Deleted Resources**

* User views issue, another user deletes it  
* **Solution:** Handle 404 gracefully, show 'This issue has been deleted'

### 

### **Network Instability**

* User makes changes offline  
* **Solution:** Disable all actions when offline, show 'No connection' banner

### **Browser Refresh During Action**

* User clicks save, then refreshes before response  
* **Solution:** Show 'Unsaved changes' warning before unload (only if data is dirty)

## **Loading State Guidelines**

* **Initial page load:** Full-page spinner or skeleton UI  
* **Filtering/search:** Keep previous results visible, show subtle loading indicator  
* **Button actions:** Replace button text with spinner, disable button  
* **Drag-and-drop:** Show loading state on card being dragged  
* **Auto-save:** Small 'Saving...' indicator near edited field

## **Key Takeaways**

* **Start simple:** Build core screens first, add polish later  
* **Consistent patterns:** Reuse components and patterns across screens  
* **Fail gracefully:** Every API call needs error handling  
* **Loading feedback:** Always show what's happening  
* **Optimistic updates:** Make UI feel instant, revert on error