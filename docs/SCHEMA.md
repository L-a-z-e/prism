# Prism Database Schema

## 1. Relational Database (MySQL)

### 1.1 Users & Organization
Basic hierarchy for multi-tenancy and ownership.

```sql
CREATE TABLE organizations (
    id VARCHAR(36) PRIMARY KEY, -- UUID
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY, -- UUID
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    role VARCHAR(50) DEFAULT 'USER', -- ADMIN, USER, ETC
    organization_id VARCHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (organization_id) REFERENCES organizations(id)
);

CREATE TABLE projects (
    id VARCHAR(36) PRIMARY KEY, -- UUID
    name VARCHAR(255) NOT NULL,
    key_name VARCHAR(50) NOT NULL, -- e.g., "PRISM" for ticket prefixes
    description TEXT,
    organization_id VARCHAR(36) NOT NULL,
    created_by VARCHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (organization_id) REFERENCES organizations(id),
    FOREIGN KEY (created_by) REFERENCES users(id)
);
```

### 1.2 AI Infrastructure

```sql
CREATE TABLE ai_providers (
    id VARCHAR(36) PRIMARY KEY, -- UUID
    name VARCHAR(100) NOT NULL, -- e.g., "OpenAI", "Anthropic"
    base_url VARCHAR(255),
    api_key_env_var VARCHAR(255), -- Name of env var holding the key (security best practice)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### 1.3 Core Domain: Agents & Tasks

```sql
CREATE TABLE agents (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL, -- 'PM', 'BACKEND', 'FRONTEND', 'QA', 'DEVOPS', 'CUSTOM'
    description TEXT,

    -- AI Configuration
    provider_id VARCHAR(36),
    model_name VARCHAR(100), -- e.g., "gpt-4", "claude-3-5-sonnet"
    system_prompt LONGTEXT,
    temperature DECIMAL(3, 2) DEFAULT 0.7,
    max_tokens INT DEFAULT 4096,

    -- Capabilities
    can_write_code BOOLEAN DEFAULT TRUE,
    can_run_tests BOOLEAN DEFAULT TRUE,
    can_deploy BOOLEAN DEFAULT FALSE,
    can_create_documents BOOLEAN DEFAULT TRUE,
    can_merge_pr BOOLEAN DEFAULT FALSE,

    -- Ownership
    organization_id VARCHAR(36),
    created_by VARCHAR(36),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (provider_id) REFERENCES ai_providers(id),
    FOREIGN KEY (organization_id) REFERENCES organizations(id),
    FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE TABLE tasks (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description LONGTEXT,

    -- Assignment
    assigned_to VARCHAR(36), -- Agent ID
    project_id VARCHAR(36) NOT NULL,
    created_by VARCHAR(36), -- User ID
    parent_task_id VARCHAR(36),

    -- State
    status VARCHAR(50) DEFAULT 'TODO', -- 'TODO', 'IN_PROGRESS', 'CODE_REVIEW', 'IN_TESTING', 'DEPLOYED', 'DONE'
    priority VARCHAR(50) DEFAULT 'MEDIUM', -- 'CRITICAL', 'HIGH', 'MEDIUM', 'LOW'

    -- Metrics
    estimated_hours INT,
    actual_hours INT,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,

    -- Technical Context
    git_branch VARCHAR(255),
    git_commit_hash VARCHAR(40),
    git_pr_url VARCHAR(500),
    git_pr_status VARCHAR(50), -- 'DRAFT', 'OPEN', 'APPROVED', 'MERGED'

    -- Outputs
    test_result JSON,
    build_log LONGTEXT,
    deployment_log LONGTEXT,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (assigned_to) REFERENCES agents(id),
    FOREIGN KEY (project_id) REFERENCES projects(id),
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (parent_task_id) REFERENCES tasks(id)
);
```

## 2. NoSQL Database (MongoDB)

### 2.1 Activity Logs
Stores high-volume, unstructured event data for every action taken by an agent or user.

**Collection:** `activities`

```json
{
  "_id": "ObjectId",
  "task_id": "String (UUID)",
  "agent_id": "String (UUID)",
  "user_id": "String (UUID, optional)",
  "action": "String",
  // Enum: CODE_WRITTEN, TEST_RUN, PR_CREATED, BUILD_STARTED, BUILD_FAILED, COMMENT_ADDED

  "timestamp": "ISODate",

  "details": {
    // Flexible schema based on action
    "files_modified": ["src/main.go", "README.md"],
    "lines_added": 150,
    "lines_deleted": 30,
    "git_commit": "abc1234",
    "duration_seconds": 225,
    "tokens_used": 5000,
    "cost": 0.05,
    "error_message": "..."
  }
}
```
