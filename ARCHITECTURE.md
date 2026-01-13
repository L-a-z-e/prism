# Prism Architecture Document

## 시스템 아키텍처 개요

```
┌──────────────────────────────────────────────────────────────┐
│                      Frontend (Vue 3)                        │
│         - Kanban Board UI                                    │
│         - Project Management                                 │
│         - Agent Configuration                                │
│         - Real-time WebSocket                                │
└──────────────────────────────────────────────────────────────┘
                              │
                    WebSocket │ HTTP/REST
                              │
┌──────────────────────────────────────────────────────────────┐
│                 Backend API (Spring Boot 3)                  │
│         - Project API                                        │
│         - Agent API                                          │
│         - Task/Kanban API                                    │
│         - File Management API                                │
│         - Context Sharing API                                │
└──────────────────────────────────────────────────────────────┘
         │                              │
         │ gRPC                         │ Redis Pub/Sub
         │                              │
    ┌────▼────────┐             ┌──────▼──────────┐
    │   AI Daemon │             │  Message Queue  │
    │   (Go)      │             │  (Redis)        │
    │             │             │                 │
    │ - File Sync │             │ - Task Queue    │
    │ - LLM Calls │             │ - Event Stream  │
    │ - Code Gen  │             └─────────────────┘
    └────────┬────┘                    │
             │                         │
             │ fsnotify               │
             │                         │
    ┌────────▼────────────┐    ┌──────▼──────────┐
    │  Local File System  │    │   MySQL (DB)    │
    │  /project/...       │    │                 │
    │  - Source Code      │    │ - Projects      │
    │  - Generated Files  │    │ - Agents        │
    └─────────────────────┘    │ - Tasks         │
                               │ - File Changes  │
                               │ - Contexts      │
                               └─────────────────┘

                    ┌───────────────────────┐
                    │   AI Providers        │
                    │                       │
                    │ - OpenAI (GPT-4)      │
                    │ - Claude (Sonnet)     │
                    │ - Ollama (Local)      │
                    └───────────────────────┘
```

---

## 계층 구조

### 1. Presentation Layer (Vue 3)

```
Frontend/
├── src/
│   ├── components/
│   │   ├── KanbanBoard.vue
│   │   ├── ProjectCard.vue
│   │   ├── AgentConfig.vue
│   │   └── TaskCard.vue
│   ├── pages/
│   │   ├── Dashboard.vue
│   │   ├── ProjectDetail.vue
│   │   └── Settings.vue
│   ├── stores/
│   │   ├── projectStore.ts
│   │   ├── agentStore.ts
│   │   └── taskStore.ts
│   └── services/
│       ├── api.ts
│       └── websocket.ts
```

**책임**:
- UI 렌더링
- 사용자 입력 처리
- 실시간 업데이트 (WebSocket)
- 상태 관리 (Pinia)

---

### 2. API Layer (Spring Boot)

```
Services/
├── api-server/
│   ├── src/main/java/
│   │   └── com/prism/
│   │       ├── controller/
│   │       │   ├── ProjectController.java
│   │       │   ├── AgentController.java
│   │       │   ├── TaskController.java
│   │       │   └── FileController.java
│   │       ├── service/
│   │       │   ├── ProjectService.java
│   │       │   ├── AgentService.java
│   │       │   ├── TaskService.java
│   │       │   └── ContextService.java
│   │       ├── repository/
│   │       │   ├── ProjectRepository.java
│   │       │   ├── AgentRepository.java
│   │       │   ├── TaskRepository.java
│   │       │   └── FileChangeRepository.java
│   │       ├── entity/
│   │       │   ├── Project.java
│   │       │   ├── Agent.java
│   │       │   ├── Task.java
│   │       │   └── FileChange.java
│   │       ├── dto/
│   │       ├── config/
│   │       └── exception/
│   └── build.gradle
```

**책임**:
- REST API 제공
- 비즈니스 로직 구현
- 데이터베이스 접근
- Redis 메시지 큐 통신
- gRPC 클라이언트 (Daemon 호출)

**주요 엔드포인트**:
```
POST   /api/projects              - 프로젝트 생성
GET    /api/projects              - 프로젝트 목록
GET    /api/projects/{id}         - 프로젝트 상세
PUT    /api/projects/{id}         - 프로젝트 수정
DELETE /api/projects/{id}         - 프로젝트 삭제

POST   /api/agents                - 에이전트 생성
GET    /api/agents                - 에이전트 목록
PUT    /api/agents/{id}           - 에이전트 수정
DELETE /api/agents/{id}           - 에이전트 삭제
POST   /api/agents/{id}/start     - 에이전트 시작
POST   /api/agents/{id}/stop      - 에이전트 중지

GET    /api/tasks/kanban          - 칸반보드 조회
POST   /api/tasks                 - 작업 생성
PUT    /api/tasks/{id}            - 작업 수정
PUT    /api/tasks/{id}/move       - 작업 이동
GET    /api/tasks/{id}/files      - 변경 파일 조회
```

---

### 3. Business Logic Layer (Services)

#### ProjectService
```java
public class ProjectService {
    // 프로젝트 CRUD
    public Project create(ProjectRequest req);
    public Project findById(String id);
    public List<Project> findAll();
    public Project update(String id, ProjectRequest req);
    public void delete(String id);
    
    // 경로 검증
    public void validateLocalPath(String path);
    public boolean gitDetected(String localPath);
}
```

#### AgentService
```java
public class AgentService {
    // 에이전트 CRUD
    public Agent create(AgentRequest req);
    public Agent findById(String id);
    public List<Agent> findByProject(String projectId);
    public Agent update(String id, AgentRequest req);
    public void delete(String id);
    
    // 에이전트 제어
    public void startAgent(String id);
    public void stopAgent(String id);
    
    // API Key 암호화
    public void encryptApiKey(Agent agent);
    public String decryptApiKey(Agent agent);
}
```

#### TaskService
```java
public class TaskService {
    // 작업 CRUD
    public Task create(TaskRequest req);
    public Task findById(String id);
    public List<Task> findByProject(String projectId);
    public Task update(String id, TaskRequest req);
    public void delete(String id);
    
    // 칸반보드
    public KanbanBoard getKanban(String projectId);
    public void moveTask(String taskId, String column, int order);
    
    // 상태 관리
    public void updateStatus(String taskId, TaskStatus status);
}
```

#### ContextService
```java
public class ContextService {
    // Context 공유
    public SharedContext save(SharedContext context);
    public List<SharedContext> findByTask(String taskId);
    public List<SharedContext> findByAgent(String agentId);
}
```

---

### 4. Data Access Layer

#### Entity 정의

```java
@Entity
@Table(name = "projects")
public class Project {
    @Id
    private String id;
    private String name;
    private String localPath;
    private String description;
    private boolean gitDetected;
    private LocalDateTime createdAt;
    private LocalDateTime updatedAt;
}

@Entity
@Table(name = "agents")
public class Agent {
    @Id
    private String id;
    private String name;
    @Enumerated(EnumType.STRING)
    private AgentRole role;
    private String projectId;
    @Enumerated(EnumType.STRING)
    private AiProvider provider;
    private String modelName;
    @Convert(converter = AesEncryptionConverter.class)
    private String apiKey;
    private String systemPrompt;
    private Double temperature;
    private Integer maxTokens;
    // capabilities...
}

@Entity
@Table(name = "tasks")
public class Task {
    @Id
    private String id;
    private String title;
    private String description;
    private String projectId;
    private String assignedAgentId;
    @Enumerated(EnumType.STRING)
    private TaskStatus status; // TODO, IN_PROGRESS, REVIEW, DONE
    @Enumerated(EnumType.STRING)
    private TaskPriority priority; // LOW, MEDIUM, HIGH
    private Integer columnOrder;
    private Integer estimatedHours;
    private Integer actualHours;
    private LocalDateTime createdAt;
    private LocalDateTime updatedAt;
}

@Entity
@Table(name = "file_changes")
public class FileChange {
    @Id
    private String id;
    private String taskId;
    private String filePath;
    @Enumerated(EnumType.STRING)
    private ChangeType changeType; // CREATE, UPDATE, DELETE
    private String diff;
    private Integer sizeBytes;
    private LocalDateTime createdAt;
}

@Entity
@Table(name = "shared_contexts")
public class SharedContext {
    @Id
    private String id;
    private String fromAgentId;
    private String taskId;
    @Enumerated(EnumType.STRING)
    private ContextType contextType; // API_SPEC, CODE_OUTPUT, DOCUMENTATION
    @Lob
    private String content; // JSON
    private LocalDateTime createdAt;
}
```

---

### 5. Infrastructure Layer

#### Redis (Message Queue & Cache)
```
Redis Pub/Sub Channel: "prism:tasks:{projectId}"
  Message Format:
  {
    "type": "TASK_ASSIGNED",
    "taskId": "task-123",
    "agentId": "agent-456",
    "timestamp": "2026-01-13T10:00:00Z"
  }

Redis Key: "agent:{agentId}:status"
  Value: "ONLINE" | "OFFLINE" | "BUSY"

Redis Key: "task:{taskId}:progress"
  Value: JSON (progress percentage, current file, etc.)
```

#### MySQL (Persistent Storage)
```sql
-- 데이터베이스 선택
USE prism_db;

-- 프로젝트 테이블
CREATE TABLE projects (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    local_path VARCHAR(500) NOT NULL UNIQUE,
    description TEXT,
    git_detected BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_local_path (local_path)
);

-- 에이전트 테이블
CREATE TABLE agents (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    role VARCHAR(50) NOT NULL, -- BACKEND, FRONTEND, DOCS, TEST, DEVOPS
    project_id VARCHAR(36) NOT NULL,
    provider VARCHAR(50) NOT NULL, -- OPENAI, CLAUDE, OLLAMA
    model_name VARCHAR(100),
    api_key_encrypted TEXT,
    system_prompt LONGTEXT,
    temperature DECIMAL(3,2) DEFAULT 0.70,
    max_tokens INT DEFAULT 4096,
    can_write_code BOOLEAN DEFAULT TRUE,
    can_run_tests BOOLEAN DEFAULT TRUE,
    can_deploy BOOLEAN DEFAULT FALSE,
    can_create_documents BOOLEAN DEFAULT TRUE,
    can_merge_pr BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    INDEX idx_project_role (project_id, role)
);

-- 작업 테이블
CREATE TABLE tasks (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description LONGTEXT,
    project_id VARCHAR(36) NOT NULL,
    assigned_to VARCHAR(36),
    status VARCHAR(20) DEFAULT 'TODO', -- TODO, IN_PROGRESS, REVIEW, DONE
    priority VARCHAR(20) DEFAULT 'MEDIUM', -- LOW, MEDIUM, HIGH
    column_order INT DEFAULT 0,
    estimated_hours INT,
    actual_hours INT,
    started_at TIMESTAMP NULL,
    completed_at TIMESTAMP NULL,
    git_branch VARCHAR(200),
    git_commit_hash VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (assigned_to) REFERENCES agents(id) ON DELETE SET NULL,
    INDEX idx_project_status (project_id, status),
    INDEX idx_column_order (project_id, status, column_order)
);

-- 파일 변경 테이블
CREATE TABLE file_changes (
    id VARCHAR(36) PRIMARY KEY,
    task_id VARCHAR(36) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    change_type VARCHAR(20) NOT NULL, -- CREATE, UPDATE, DELETE
    diff LONGTEXT,
    size_bytes INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    INDEX idx_task (task_id)
);

-- Context 공유 테이블
CREATE TABLE shared_contexts (
    id VARCHAR(36) PRIMARY KEY,
    from_agent_id VARCHAR(36) NOT NULL,
    task_id VARCHAR(36),
    context_type VARCHAR(50) NOT NULL, -- API_SPEC, CODE_OUTPUT, DOCUMENTATION
    content LONGTEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (from_agent_id) REFERENCES agents(id) ON DELETE CASCADE,
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    INDEX idx_task (task_id)
);
```

---

### 6. Daemon Layer (Go)

```
daemon/
├── main.go
├── config/
│   └── config.go
├── ai/
│   ├── openai.go
│   ├── claude.go
│   ├── ollama.go
│   └── provider.go (interface)
├── task/
│   ├── listener.go    (Redis 구독)
│   ├── executor.go    (작업 실행)
│   └── uploader.go    (파일 업로드)
├── file/
│   ├── watcher.go     (fsnotify)
│   ├── reader.go      (파일 읽기)
│   └── writer.go      (파일 쓰기)
├── rpc/
│   ├── client.go      (gRPC 클라이언트)
│   └── service.pb.go
└── Dockerfile
```

**주요 기능**:

1. **Task Listener**
   - Redis Pub/Sub 구독
   - Task 메시지 수신
   - Task 실행 트리거

2. **AI Integration**
   - Provider별 클라이언트
   - 요청/응답 처리
   - 에러 핸들링
   - 재시도 로직

3. **File Watcher**
   - fsnotify로 파일 변경 감지
   - 생성/수정/삭제 추적
   - gRPC로 백엔드에 보고

4. **Code Generator**
   - 템플릿 기반 코드 생성
   - Context 기반 생성
   - 파일 시스템에 저장

**워크플로우**:
```
1. Redis 메시지 수신
   {"taskId": "task-123", "agentId": "agent-456"}

2. 에이전트 정보 조회 (gRPC)
   GetAgent("agent-456") -> Provider: OPENAI, Model: GPT-4, ...

3. AI 호출
   CreateCompletion(model, systemPrompt, userPrompt) -> response

4. 파일 생성
   WriteFile("/project/src/UserController.java", code)

5. 파일 변경 보고 (gRPC)
   ReportFileChange(taskId, filePath, "CREATE", diff)

6. Task 상태 업데이트
   UpdateTask(taskId, "IN_PROGRESS")

7. WebSocket 알림
   broadcast("file.changed", {taskId, filePath})
```

---

### 7. Communication Protocols

#### REST API (Frontend ↔ Backend)
```
HTTP/1.1
Content-Type: application/json
Authorization: Bearer {jwt_token}

Example:
GET /api/tasks/kanban?projectId=proj-123 HTTP/1.1
Host: localhost:8080
Accept: application/json
Authorization: Bearer eyJhbGc...

Response:
{
  "columns": {
    "TODO": [...],
    "IN_PROGRESS": [...],
    "REVIEW": [...],
    "DONE": [...]
  }
}
```

#### gRPC (Backend ↔ Daemon)
```protobuf
service AgentService {
  rpc GetAgent(GetAgentRequest) returns (Agent) {}
  rpc UpdateTaskStatus(UpdateTaskStatusRequest) returns (Empty) {}
  rpc ReportFileChange(FileChangeReport) returns (Empty) {}
  rpc GetProjectConfig(GetProjectConfigRequest) returns (ProjectConfig) {}
}
```

#### Redis Pub/Sub (Backend ↔ Daemon)
```
Channel: "prism:tasks:{projectId}"
Message:
{
  "eventType": "TASK_ASSIGNED",
  "taskId": "task-123",
  "agentId": "agent-456",
  "projectId": "proj-123",
  "priority": "HIGH",
  "timestamp": "2026-01-13T10:00:00Z"
}
```

#### WebSocket (Backend → Frontend)
```javascript
// 클라이언트 연결
const socket = io('http://localhost:8080', {
  auth: { token: jwt_token }
});

// 이벤트 수신
socket.on('task.status.changed', (data) => {
  // {taskId, oldStatus, newStatus}
});

socket.on('file.changed', (data) => {
  // {taskId, filePath, changeType}
});

socket.on('agent.status.changed', (data) => {
  // {agentId, status}
});
```

---

## 데이터 흐름 (Data Flow)

### 시나리오: 사용자가 Task 생성 후 AI 코드 생성

```
1. 사용자 (Frontend)
   POST /api/tasks
   {
     "projectId": "proj-123",
     "title": "User CRUD API",
     "description": "Create REST endpoints...",
     "assignedAgentId": "agent-456",
     "priority": "HIGH"
   }
   ↓

2. Backend (Spring Boot)
   - TaskService.create() 호출
   - Task 저장 (status: TODO)
   - Redis 메시지 발행
   {
     "eventType": "TASK_ASSIGNED",
     "taskId": "task-123",
     "agentId": "agent-456",
     ...
   }
   - 응답 반환
   ↓

3. Daemon (Go)
   - Redis 메시지 수신
   - gRPC로 Agent 정보 조회
   - AI Provider 결정 (OpenAI GPT-4)
   - 시스템 프롬프트 + 사용자 프롬프트 생성
   - OpenAI API 호출
   ↓

4. OpenAI
   - 요청 처리
   - Java 코드 생성
   - 응답 반환
   ↓

5. Daemon (Go)
   - 코드 파일 생성
     /project/src/main/java/UserController.java
     /project/src/main/java/UserService.java
     /project/src/main/java/UserRepository.java
   - gRPC로 FileChange 보고
   - Task 상태 업데이트 (IN_PROGRESS)
   ↓

6. Backend (Spring Boot)
   - FileChange 엔티티 저장
   - Task 상태 업데이트
   - Redis 메시지 발행
   {
     "eventType": "FILE_CHANGED",
     "taskId": "task-123",
     "filePath": "src/main/java/UserController.java",
     "changeType": "CREATE"
   }
   ↓

7. Frontend (Vue 3)
   - WebSocket 메시지 수신
   - UI 업데이트
   - Task 카드 상태 변경
   - 생성된 파일 목록 표시
```

---

## 배포 아키텍처 (Docker Compose)

```yaml
version: '3.8'

services:
  # MySQL 데이터베이스
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: prism_db
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql

  # Redis (Message Queue & Cache)
  redis:
    image: redis:7.0
    ports:
      - "6379:6379"

  # Spring Boot Backend
  api-server:
    build: ./services/api-server
    environment:
      SPRING_DATASOURCE_URL: jdbc:mysql://mysql:3306/prism_db
      SPRING_REDIS_HOST: redis
      GRPC_SERVER_PORT: 50051
    ports:
      - "8080:8080"
      - "50051:50051"
    depends_on:
      - mysql
      - redis

  # Go Daemon
  daemon:
    build: ./daemon
    environment:
      REDIS_URL: redis://redis:6379
      GRPC_SERVER_URL: api-server:50051
      OLLAMA_URL: http://ollama:11434
    depends_on:
      - redis
      - api-server

  # Ollama (Optional Local LLM)
  ollama:
    image: ollama/ollama:latest
    ports:
      - "11434:11434"
    volumes:
      - ollama_data:/root/.ollama

  # Vue Frontend
  frontend:
    build: ./frontend
    environment:
      VUE_APP_API_URL: http://localhost:8080
    ports:
      - "3000:3000"
    depends_on:
      - api-server

volumes:
  mysql_data:
  ollama_data:
```

---

## 보안 고려사항

### 1. API Key 관리
- AES-256 암호화로 저장
- 조회 시 복호화
- 환경 변수로 master key 관리

### 2. 인증 & 인가
- JWT 토큰 기반
- 만료 시간 설정
- Role-based Access Control (RBAC)

### 3. 파일 접근 제어
- 프로젝트 경로 검증
- 심볼릭 링크 차단
- 권한 확인 (읽기/쓰기)

### 4. 통신 보안
- gRPC TLS (선택적)
- HTTPS 필수
- CORS 설정

---

## 성능 최적화

### 1. 캐싱
- Agent 정보 캐싱 (Redis)
- Project 정보 캐싱
- 쿼리 결과 캐싱

### 2. 배치 처리
- 대량 Task 생성
- FileChange 배치 저장

### 3. 비동기 처리
- Task 실행 비동기화
- 파일 업로드 비동기화
- 이메일/알림 비동기화

---

**마지막 업데이트**: 2026-01-13