# Prism - AI Multi-Agent Development Platform
## Product Requirements Document (PRD)

**Version**: 1.0  
**Date**: 2026-01-13  
**Author**: Development Team  
**Status**: Active Development

---

## 1. 제품 개요

### 1.1 비전
개발자가 로컬 프로젝트에 여러 AI 에이전트를 배치하여, 각 에이전트가 전문 영역(백엔드, 프론트엔드, 문서, 테스트)을 담당하며 협업하는 AI 기반 개발 플랫폼

### 1.2 핵심 가치 제안
- **멀티 에이전트 협업**: 백엔드/프론트/문서/테스트 에이전트가 자동 협업
- **AI Provider 유연성**: OpenAI, Claude, Ollama 등 자유롭게 선택
- **로컬 작업**: 사용자의 로컬 프로젝트에서 직접 파일 생성/수정
- **시각적 관리**: 칸반보드로 작업 흐름 관리
- **실시간 피드백**: 파일 변경사항 즉시 확인

### 1.3 목표 사용자
- **Primary**: 개인 개발자, 스타트업 개발팀 (1-5명)
- **Secondary**: 중소기업 개발팀, 프리랜서
- **사용 사례**:
  - 빠른 프로토타입 개발
  - 반복 작업 자동화 (CRUD API, 문서화)
  - 레거시 코드 리팩토링
  - 테스트 코드 생성

---

## 2. 핵심 기능

### 2.1 프로젝트 관리

#### F-101: 프로젝트 등록
**설명**: 사용자가 로컬 프로젝트 경로를 등록하여 AI 에이전트가 접근 가능하도록 설정

**Acceptance Criteria**:
- [ ] 프로젝트 이름 입력 (필수)
- [ ] 로컬 경로 선택 (디렉토리 브라우저)
- [ ] 경로 유효성 검증 (존재 여부, 권한)
- [ ] Git 저장소 자동 감지
- [ ] 프로젝트 정보 저장 (DB)

**API**:
```
POST /api/projects
{
  "name": "My E-commerce",
  "localPath": "/Users/laze/projects/ecommerce",
  "description": "Online shopping platform"
}

Response 201:
{
  "id": "proj-123",
  "name": "My E-commerce",
  "localPath": "/Users/laze/projects/ecommerce",
  "gitDetected": true,
  "createdAt": "2026-01-13T10:00:00Z"
}
```

---

### 2.2 에이전트 관리

#### F-201: 에이전트 생성 및 설정
**설명**: 프로젝트에 특화된 역할의 AI 에이전트를 생성하고 AI Provider 설정

**에이전트 역할 타입**:
| 역할 | 설명 | 기본 권한 |
|------|------|----------|
| BACKEND | Java/Spring Boot 백엔드 개발 | 코드 작성, 테스트 실행 |
| FRONTEND | Vue/React 프론트엔드 개발 | 코드 작성, 컴포넌트 생성 |
| DOCS | 문서 작성 (README, API 문서) | 문서 작성, 코드 리뷰 |
| TEST | 테스트 코드 작성 및 실행 | 테스트 작성, 실행, 리포트 |
| DEVOPS | 배포 및 인프라 관리 | 배포, CI/CD 설정 |

**API**:
```
POST /api/agents
{
  "projectId": "proj-123",
  "name": "Backend Developer",
  "role": "BACKEND",
  "aiProvider": "OPENAI",
  "model": "gpt-4-turbo",
  "apiKey": "sk-xxx...",
  "systemPrompt": "You are an expert Java Spring Boot developer...",
  "capabilities": {
    "canWriteCode": true,
    "canRunTests": true,
    "canDeploy": false,
    "canCreateDocuments": true
  }
}

Response 201:
{
  "id": "agent-456",
  "name": "Backend Developer",
  "role": "BACKEND",
  "status": "ONLINE",
  "createdAt": "2026-01-13T10:05:00Z"
}
```

---

### 2.3 칸반보드 작업 관리

#### F-301: 칸반보드 뷰
**설명**: 드래그 앤 드롭 가능한 칸반보드로 작업 상태 시각화

**칸반 컬럼**:
1. **TODO** (할 일)
2. **IN_PROGRESS** (작업 중)
3. **REVIEW** (리뷰 대기)
4. **DONE** (완료)

**작업 카드 속성**:
- 제목, 설명, 상태, 우선순위
- 할당된 에이전트
- 변경된 파일 목록
- 진행도 표시

**API**:
```
GET /api/tasks/kanban?projectId=proj-123

Response 200:
{
  "columns": {
    "TODO": [...],
    "IN_PROGRESS": [...],
    "REVIEW": [...],
    "DONE": [...]
  }
}

PUT /api/tasks/{taskId}/move
{
  "toColumn": "IN_PROGRESS",
  "newOrder": 1
}
```

---

#### F-302: 작업 생성 및 할당
**설명**: 새 작업을 생성하고 에이전트에 자동/수동 할당

**Acceptance Criteria**:
- [ ] 작업 제목, 설명 입력
- [ ] 에이전트 선택 (드롭다운)
- [ ] 우선순위 설정 (LOW/MEDIUM/HIGH)
- [ ] 예상 시간 입력 (선택)
- [ ] 의존 작업 설정 (선택)

---

### 2.4 AI 코드 생성 및 파일 관리

#### F-401: AI 코드 생성
**설명**: 에이전트가 작업 설명을 기반으로 코드 자동 생성

**워크플로우**:
```
1. 사용자가 Task 생성 및 에이전트 할당
   ↓
2. Daemon이 Task 수신 (Redis Pub/Sub)
   ↓
3. Agent가 AI에게 요청 (Chat API)
   - System Prompt: "You are a backend developer..."
   - User Prompt: "Create User CRUD API..."
   ↓
4. AI 응답 수신 (20-60초)
   ↓
5. 로컬 파일 생성
   - /project/src/main/java/UserController.java
   - /project/src/main/java/UserService.java
   ↓
6. 상태 업데이트 (TODO → IN_PROGRESS → REVIEW)
   ↓
7. 웹 UI에 실시간 알림 (WebSocket)
```

**생성 가능한 파일 타입**:
- Java: `.java` (Controller, Service, Repository, Entity)
- JavaScript/TypeScript: `.js`, `.ts`, `.vue`, `.tsx`
- Markdown: `.md` (README, API 문서)
- Config: `.yml`, `.properties`, `.json`
- Test: `*Test.java`, `*.spec.ts`

---

#### F-402: 파일 변경 추적
**설명**: 에이전트가 생성/수정한 파일을 실시간 추적 및 Diff 표시

**API**:
```
GET /api/tasks/{taskId}/files

Response 200:
{
  "files": [
    {
      "path": "src/main/java/UserController.java",
      "changeType": "CREATE",
      "size": 1024,
      "diff": "@@ -0,0 +1,15 @@\n+@RestController...",
      "createdAt": "2026-01-13T10:10:00Z"
    }
  ]
}
```

---

### 2.5 에이전트 협업

#### F-501: Context Sharing (컨텍스트 공유)
**설명**: 에이전트 간 작업 결과물 및 정보 공유

**사용 사례**:
```
Scenario: Backend가 API 완성 → Frontend가 자동으로 컴포넌트 생성

1. Backend Agent가 User API 완성
   - POST /api/users
   - GET /api/users
   
2. Context 공유
   {
     "type": "API_SPEC",
     "fromAgent": "Backend Developer",
     "content": {
       "endpoints": [
         {"method": "POST", "path": "/api/users", "body": {...}},
         {"method": "GET", "path": "/api/users", "response": {...}}
       ]
     }
   }
   
3. Frontend Agent가 Context 읽음

4. 자동으로 Vue 컴포넌트 생성
   - UserList.vue (GET /api/users 호출)
   - UserForm.vue (POST /api/users 호출)
```

---

## 3. 기술 스택

### 3.1 프론트엔드
- **Framework**: Vue 3 + TypeScript
- **UI Library**: Vuetify / Element Plus
- **State Management**: Pinia
- **API Client**: Axios
- **Real-time**: WebSocket (Socket.io)
- **Code Diff**: Monaco Editor / Diff2Html

### 3.2 백엔드
- **Framework**: Spring Boot 3.x
- **Language**: Java 17+
- **Database**: MySQL 8.0
- **Cache**: Redis 7.0
- **Message Queue**: Redis Pub/Sub
- **gRPC**: Spring gRPC
- **Security**: Spring Security + JWT

### 3.3 AI Daemon
- **Language**: Go 1.24+
- **AI Clients**:
  - OpenAI (GPT-4, GPT-4-turbo)
  - Anthropic Claude (Claude 3.5 Sonnet)
  - Ollama (Local LLM - DeepSeek, Llama)
- **File System**: fsnotify (파일 감시)
- **Communication**: gRPC, Redis

### 3.4 인프라
- **Container**: Docker + Docker Compose
- **Orchestration**: Kubernetes (Optional)
- **CI/CD**: GitHub Actions
- **Monitoring**: Prometheus + Grafana

---

## 4. 데이터 모델

### 4.1 핵심 엔티티

```sql
-- 프로젝트
CREATE TABLE projects (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    local_path VARCHAR(500) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 에이전트
CREATE TABLE agents (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    role VARCHAR(50) NOT NULL,
    provider_id VARCHAR(36),
    model_name VARCHAR(100),
    api_key_encrypted TEXT,
    system_prompt LONGTEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 작업
CREATE TABLE tasks (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description LONGTEXT,
    assigned_to VARCHAR(36),
    project_id VARCHAR(36) NOT NULL,
    status VARCHAR(20) DEFAULT 'TODO',
    priority VARCHAR(20) DEFAULT 'MEDIUM',
    column_order INT DEFAULT 0,
    estimated_hours INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 파일 변경
CREATE TABLE file_changes (
    id VARCHAR(36) PRIMARY KEY,
    task_id VARCHAR(36) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    change_type VARCHAR(20) NOT NULL,
    diff LONGTEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## 5. API 명세

### 5.1 프로젝트 API

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/projects` | 프로젝트 생성 |
| GET | `/api/projects` | 프로젝트 목록 |
| GET | `/api/projects/{id}` | 프로젝트 상세 |
| PUT | `/api/projects/{id}` | 프로젝트 수정 |
| DELETE | `/api/projects/{id}` | 프로젝트 삭제 |

### 5.2 에이전트 API

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/agents` | 에이전트 생성 |
| GET | `/api/agents?projectId={id}` | 에이전트 목록 |
| GET | `/api/agents/{id}` | 에이전트 상세 |
| PUT | `/api/agents/{id}` | 에이전트 수정 |
| DELETE | `/api/agents/{id}` | 에이전트 삭제 |

### 5.3 작업 API

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/tasks/kanban?projectId={id}` | 칸반보드 조회 |
| POST | `/api/tasks` | 작업 생성 |
| GET | `/api/tasks/{id}` | 작업 상세 |
| PUT | `/api/tasks/{id}` | 작업 수정 |
| PUT | `/api/tasks/{id}/move` | 작업 이동 |
| DELETE | `/api/tasks/{id}` | 작업 삭제 |
| GET | `/api/tasks/{id}/files` | 변경된 파일 목록 |

---

## 6. 비기능적 요구사항

### 6.1 성능
- AI 응답 시간: 20-60초 (모델별 차이)
- API 응답 시간: < 200ms (P95)
- 파일 변경 감지: < 1초
- WebSocket 지연: < 100ms

### 6.2 보안
- API Key 암호화 저장 (AES-256)
- JWT 기반 인증
- 로컬 파일 접근 권한 검증
- HTTPS 필수

### 6.3 확장성
- 동시 사용자: 100명
- 프로젝트 수: 무제한
- 에이전트 수/프로젝트: 10개
- 작업 수/프로젝트: 1000개

### 6.4 가용성
- Uptime: 99%
- 데이터 백업: 일 1회
- 장애 복구: < 1시간

---

## 7. 개발 로드맵

### Phase 1: MVP (4주)
**목표**: 기본 칸반보드 + 단일 에이전트 코드 생성

**Week 1-2**: 백엔드 기반
- [ ] Project에 `localPath` 추가
- [ ] Agent에 `apiKeyEncrypted` 추가
- [ ] Task에 `columnOrder` 추가
- [ ] 칸반보드 API 구현

**Week 3**: 프론트엔드
- [ ] 프로젝트 설정 페이지
- [ ] 에이전트 설정 페이지
- [ ] 칸반보드 UI (드래그 앤 드롭)

**Week 4**: Daemon + AI 통합
- [ ] OpenAI/Claude 클라이언트 추가
- [ ] 로컬 파일 생성 로직
- [ ] 테스트

### Phase 2: 협업 기능 (3주)
- [ ] Context Sharing 구현
- [ ] 멀티 에이전트 조정
- [ ] 파일 Diff 뷰어

### Phase 3: 고도화 (2주)
- [ ] 테스트 에이전트
- [ ] DevOps 에이전트
- [ ] 성능 최적화

---

## 8. 성공 지표 (KPI)

### 사용자 지표
- **활성 사용자**: 50명 (3개월)
- **프로젝트 생성**: 100개
- **작업 완료**: 1000개

### 품질 지표
- **AI 코드 정확도**: > 80%
- **사용자 만족도**: > 4.0/5.0
- **버그 발생률**: < 5%

---

## 9. 리스크 및 대응

| 리스크 | 영향도 | 확률 | 대응 방안 |
|--------|--------|------|----------|
| AI API 비용 초과 | 높음 | 중간 | Ollama 로컬 옵션, 요청 캐싱 |
| 로컬 파일 권한 문제 | 중간 | 높음 | 권한 확인 로직, 명확한 오류 메시지 |
| 멀티 에이전트 충돌 | 중간 | 중간 | 작업 잠금, 순차 실행 옵션 |
| AI 응답 품질 낮음 | 높음 | 낮음 | 프롬프트 엔지니어링, 피드백 루프 |

---

**승인**: _________________  
**날짜**: 2026-01-13