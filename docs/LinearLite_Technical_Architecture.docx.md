# **Linear-lite**

Technical Architecture Document

## **Technology Stack**

| Layer | Technology |
| :---- | :---- |
| **Frontend** | React 18 \+ TypeScript 5 \+ Vite 5 |
| **UI Components** | Tailwind CSS \+ shadcn/ui (optional) |
| **State Management** | React Query (TanStack Query) \+ Zustand |
| **Backend** | Go 1.21+ \+ Gin Web Framework |
| **Database** | PostgreSQL 15+ \+ sqlx (or GORM) |
| **Caching** | Redis 7+ (sessions, query cache) |
| **Authentication** | JWT (golang-jwt/jwt) \+ bcrypt |
| **Migrations** | golang-migrate/migrate |
| **Deployment** | Docker \+ Docker Compose |

## **System Architecture**

Linear-lite follows a clean three-tier architecture with clear separation of concerns:

* **Presentation Layer:** React SPA with TypeScript  
* **API Layer:** RESTful JSON API built with Go \+ Gin  
* **Data Layer:** PostgreSQL with Redis caching

## **Project Structure**

**Recommended monorepo structure:**

linear-lite/

├── frontend/              \# React application

│   ├── src/

│   │   ├── components/    \# Reusable UI components

│   │   ├── pages/         \# Page components

│   │   ├── hooks/         \# Custom React hooks

│   │   ├── services/      \# API client

│   │   ├── store/         \# Zustand stores

│   │   ├── types/         \# TypeScript types

│   │   └── utils/         \# Helper functions

│   ├── package.json

│   └── vite.config.ts

│

├── backend/               \# Go API server

│   ├── cmd/

│   │   └── api/           \# Main application

│   │       └── main.go

│   ├── internal/

│   │   ├── handlers/      \# HTTP handlers (controllers)

│   │   ├── services/      \# Business logic

│   │   ├── repository/    \# Database layer

│   │   ├── models/        \# Domain models

│   │   ├── middleware/    \# Auth, CORS, logging

│   │   ├── config/        \# Configuration

│   │   └── cache/         \# Redis cache layer

│   ├── migrations/        \# Database migrations

│   ├── go.mod

│   └── go.sum

│

├── docker-compose.yml     \# Local development setup

├── Dockerfile.frontend    \# Frontend container

├── Dockerfile.backend     \# Backend multi-stage build

└── README.md              \# Documentation

## **Database Schema**

### **Core Tables**

**users**

| Column | Type | Constraints | Notes |
| :---- | :---- | :---- | :---- |
| id | UUID | PRIMARY KEY |  |
| email | VARCHAR(255) | UNIQUE NOT NULL |  |
| password\_hash | VARCHAR(255) | NOT NULL | bcrypt |
| name | VARCHAR(255) | NOT NULL |  |
| avatar\_url | TEXT | NULL | Optional |
| created\_at | TIMESTAMP | DEFAULT NOW() |  |
| updated\_at | TIMESTAMP | DEFAULT NOW() |  |

Note: Complete schema for all tables (projects, sprints, issues, labels, issue\_labels, issue\_activities) is documented in the database migration files.

## **API Endpoints**

**Base URL:** http://localhost:8080/api/v1

### **Authentication**

| Method | Endpoint | Description |
| :---- | :---- | :---- |
| **POST** | /auth/register | Register new user |
| **POST** | /auth/login | Login and get JWT token |
| **GET** | /auth/me | Get current user profile |

### **Issues (Protected Routes)**

| Method | Endpoint | Description |
| :---- | :---- | :---- |
| **GET** | /issues | List all issues (with filters) |
| **GET** | /issues/:id | Get single issue by ID |
| **POST** | /issues | Create new issue |
| **PUT** | /issues/:id | Update existing issue |
| **DELETE** | /issues/:id | Delete issue |

Additional endpoints for: projects, sprints, labels, users \- follow the same RESTful pattern.

## **4-Week Implementation Plan**

### **Week 1: Foundation & Auth**

* Day 1-2: Project setup (Go modules, Vite, Docker Compose, PostgreSQL, Redis)  
* Day 3-4: Database schema design and migrations  
* Day 5-7: User authentication (registration, login, JWT middleware)

### **Week 2: Core CRUD & Filtering**

* Day 8-10: Issues CRUD (backend \+ frontend)  
* Day 11-12: Projects and Sprints CRUD  
* Day 13-14: Labels system and filtering logic

### **Week 3: Views & Advanced Features**

* Day 15-17: List view with sorting and search  
* Day 18-20: Kanban board view with drag-and-drop  
* Day 21: Activity timeline and issue details

### **Week 4: Polish & Deployment**

* Day 22-23: UI/UX polish, responsive design  
* Day 24-25: Performance optimization, Redis caching  
* Day 26-27: Docker multi-stage builds, docker-compose.yml  
* Day 28: Documentation (README, API docs, architecture diagrams)

## **Key Technical Decisions**

### **Why Gin over Echo?**

* Gin is faster (router benchmark: \~10x vs Echo)  
* More mature ecosystem and middleware  
* Better documentation and community support  
* Built-in JSON validation with struct tags

### **Why sqlx over GORM?**

* sqlx: More explicit, closer to raw SQL, better performance  
* GORM: More developer-friendly, faster prototyping, migrations built-in  
* **Recommendation: Start with GORM for speed, migrate to sqlx if performance matters**

### **Why React Query?**

* Server state management out of the box  
* Automatic caching, refetching, and background updates  
* Reduces boilerplate for API calls  
* Perfect for CRUD-heavy applications

## **Next Steps**

* **Set up development environment:** Install Go 1.21+, Node.js 20+, Docker  
* **Initialize repositories:** Create monorepo structure, initialize Go modules and npm  
* **Start with backend foundation:** Database schema, migrations, user model  
* **Build authentication first:** Registration, login, JWT middleware  
* **Iterate feature by feature:** Complete backend \+ frontend for each feature before moving on