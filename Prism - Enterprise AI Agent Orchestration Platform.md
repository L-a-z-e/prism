# 📋 Prism - Enterprise AI Agent Orchestration Platform

## PRD (Product Requirements Document)


***

## 1. 프로젝트 개요

### 1.1 프로젝트명

**Prism** - Enterprise AI Agent Orchestration Platform

### 1.2 한 줄 설명

개발팀이 AI 에이전트들을 중앙에서 관리하고, 각 개발자의 로컬 환경에서 자동으로 코드를 생성/수정/테스트하는 엔터프라이즈 협업 플랫폼

### 1.3 목표

- 🤖 **AI 에이전트 중앙 관리**: 여러 AI Provider(Claude, Gemini, OpenAI)를 하나의 플랫폼에서 제어
- 🔧 **자동화된 개발 프로세스**: 요청 → 코드 생성 → 테스트 → 빌드 → PR 생성까지 완전 자동화
- 👥 **팀 협업 강화**: PM과 개발자들이 AI 에이전트와 함께 협업
- 📊 **투명한 추적**: 모든 작업의 이력과 활동을 시각화
- 🚀 **엔터프라이즈 준비**: 확장성, 보안, 감시 기능 완비


### 1.4 사용자

- **Primary**: Backend/Frontend/DevOps 엔지니어
- **Secondary**: 프로덕트 매니저, QA 엔지니어
- **Organization**: 10-500명 규모의 개발팀


### 1.5 배포 환경

- 웹 기반 UI (Portal Universe 내)
- Daemon 기반 로컬 에이전트 (각 개발자 PC)
- Docker/Kubernetes (MSA 배포)

***

## 2. 현재 상황 분석

### 2.1 기존 솔루션의 문제점

| 방식 | 문제 | 영향 |
| :-- | :-- | :-- |
| **GitHub Copilot** | IDE 내장, 팀 관리 불가 | 개인 작업만 가능 |
| **Cursor** | 개인 개발 도구, 협업 미흡 | 팀 규모 확장 불가 |
| **Olly Molly** | 로컬 CLI 기반, 항상 실행 필요 | UX 나쁨, 관리 복잡 |
| **커스텀 솔루션** | 구현 복잡, 유지보수 어려움 | 개발 리소스 낭비 |

### 2.2 시장의 기회

- AI 기반 코드 생성 수요 급증
- 엔터프라이즈 팀 협업 도구 부족
- 자동화된 워크플로우 필요성 증가


### 2.3 Portal Universe와의 시너지

- 기존 MSA 인프라 활용
- 사용자 인증/권한 관리 재사용
- 데이터베이스/메시지 큐 공유
- Module Federation으로 UI 통합

***

## 3. 아키텍처

### 3.1 전체 시스템 구조

```
┌─────────────────────────────────────────────────────┐
│ Portal Universe (Central Platform)                  │
├─────────────────────────────────────────────────────┤
│                                                     │
│  ┌──────────────────────────────────────────────┐  │
│  │ Web UI (Portal Shell)                        │  │
│  │ - Agent Manager Dashboard                    │  │
│  │ - Task Board (Kanban)                        │  │
│  │ - Activity Log & Timeline                    │  │
│  │ - Settings & Configuration                   │  │
│  └──────────────────────────────────────────────┘  │
│                                                     │
│  ┌──────────────────────────────────────────────┐  │
│  │ Backend Services (Spring Boot MSA)           │  │
│  │                                              │  │
│  │ ├─ Agent Manager Service                    │  │
│  │ │  • 에이전트 설정 관리                     │  │
│  │ │  • AI Provider API 키 관리                │  │
│  │ │  • 프롬프트 템플릿 관리                   │  │
│  │ │                                           │  │
│  │ ├─ Task Orchestrator Service               │  │
│  │ │  • 작업 큐 관리 (Redis)                  │  │
│  │ │  • Daemon에게 작업 할당                  │  │
│  │ │  • 상태 변화 추적                        │  │
│  │ │                                           │  │
│  │ ├─ Git Integration Service                 │  │
│  │ │  • PR 생성/관리                          │  │
│  │ │  • 커밋 정보 저장                        │  │
│  │ │  • 브랜치 관리                           │  │
│  │ │                                           │  │
│  │ ├─ Activity Logger Service                 │  │
│  │ │  • 모든 작업 로깅                        │  │
│  │ │  • 타임라인 생성                         │  │
│  │ │  • 감시 및 감사                          │  │
│  │ │                                           │  │
│  │ ├─ Document Generator Service              │  │
│  │ │  • Notion 연동                           │  │
│  │ │  • Markdown 생성                         │  │
│  │ │  • GitHub Wiki 업데이트                  │  │
│  │ │                                           │  │
│  │ └─ Deployment Service                      │  │
│  │    • E2E 테스트 트리거                     │  │
│  │    • 배포 자동화                           │  │
│  │    • 롤백 관리                             │  │
│  │                                              │  │
│  └──────────────────────────────────────────────┘  │
│                                                     │
│  ┌──────────────────────────────────────────────┐  │
│  │ Data Layer                                   │  │
│  │ ├─ MySQL: 구조화 데이터                     │  │
│  │ ├─ MongoDB: 활동 로그                       │  │
│  │ ├─ Redis: 작업 큐, 캐시                    │  │
│  │ └─ Kafka: 서비스 간 통신                   │  │
│  │                                              │  │
│  └──────────────────────────────────────────────┘  │
│                                                     │
└─────────────────────────────────────────────────────┘
     ↕ gRPC / WebSocket
┌─────────────────────────────────────────────────────┐
│ Each Developer's PC                                 │
├─────────────────────────────────────────────────────┤
│                                                     │
│  🔧 Prism Daemon                                    │
│                                                     │
│  ├─ gRPC Server (:5000)                           │
│  ├─ Redis Listener (작업 큐 수신)                │
│  ├─ Project Manager (로컬 프로젝트 감시)         │
│  ├─ Claude/Gemini/OpenAI CLI 실행                │
│  ├─ Local Environment (npm/python/go 실행)      │
│  └─ Git Operations (커밋, 푸시)                  │
│                                                     │
│  ✨ 특징:                                         │
│  • 시스템 서비스로 자동 시작                      │
│  • 항상 백그라운드 실행                          │
│  • UI와 무관하게 독립 작동                       │
│  • 멀티 프로젝트 지원                            │
│                                                     │
└─────────────────────────────────────────────────────┘
```


### 3.2 통신 프로토콜

```
Portal ←→ Daemon: gRPC + WebSocket
├─ Task Assignment (gRPC)
├─ Task Progress Updates (WebSocket)
├─ Result Submission (gRPC)
└─ Log Streaming (WebSocket)

Portal ←→ External APIs:
├─ GitHub API (Git 작업)
├─ Claude/Gemini/OpenAI API (LLM)
├─ Notion API (문서)
└─ Slack/Email (알림)
```


***

## 4. 핵심 기능 요구사항

### 4.1 에이전트 관리 (Agent Management)

#### 4.1.1 다중 AI Provider 지원

**기능:**

- [x] 여러 AI Provider 지원
    - Claude (Anthropic)
    - Gemini (Google)
    - GPT-4 (OpenAI)
    - LLaMA (로컬 모델 지원)
- [x] API Key 암호화 저장
- [x] Provider별 모델 선택 가능
- [x] 비용 추적 (토큰 사용량)

**데이터 모델:**

```sql
CREATE TABLE ai_providers (
  id VARCHAR(36) PRIMARY KEY,
  name ENUM('CLAUDE', 'GEMINI', 'OPENAI', 'LLAMA'),
  api_key_encrypted VARCHAR(500),
  model_name VARCHAR(100),
  cost_per_1k_tokens DECIMAL(8, 4),
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);
```


#### 4.1.2 동적 에이전트 역할 생성

**기본 역할:**

- PM (Product Manager)
- Backend Developer
- Frontend Developer
- QA Engineer
- DevOps Engineer

**커스텀 역할:**

- 사용자 정의 역할 추가 가능
- 각 역할별 시스템 프롬프트 설정
- 권한/기능 커스터마이징

**데이터 모델:**

```sql
CREATE TABLE agents (
  id VARCHAR(36) PRIMARY KEY,
  name VARCHAR(255),
  role ENUM('PM', 'BACKEND', 'FRONTEND', 'QA', 'DEVOPS', 'CUSTOM'),
  description TEXT,
  
  -- AI Configuration
  provider_id VARCHAR(36),
  model_name VARCHAR(100),
  
  -- Capabilities
  system_prompt LONGTEXT,
  custom_instructions LONGTEXT,
  temperature DECIMAL(2, 1),
  max_tokens INT,
  
  -- Permissions
  can_write_code BOOLEAN DEFAULT TRUE,
  can_run_tests BOOLEAN DEFAULT TRUE,
  can_deploy BOOLEAN DEFAULT FALSE,
  can_create_documents BOOLEAN DEFAULT TRUE,
  can_merge_pr BOOLEAN DEFAULT FALSE,
  
  -- Organization
  organization_id VARCHAR(36),
  created_by VARCHAR(36),
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);
```


#### 4.1.3 에이전트 설정 DB 저장

**저장 항목:**

- 에이전트 프로필 (이름, 역할, 권한)
- AI Provider 설정
- 커스텀 시스템 프롬프트
- 능력 (코드 작성, 테스트, 배포)
- 사용 통계 (작업 수, 토큰 사용량, 비용)

**조회 API:**

```http
GET /api/agents
GET /api/agents/{agentId}
POST /api/agents
PUT /api/agents/{agentId}
DELETE /api/agents/{agentId}
```


***

### 4.2 작업 관리 (Task Management)

#### 4.2.1 향상된 Kanban 보드

**상태 워크플로우:**

```
TODO
  ↓
IN_PROGRESS (에이전트 작업 중)
  ↓
CODE_REVIEW (코드 리뷰 대기)
  ↓
IN_TESTING (테스트 실행 중)
  ↓
DEPLOYED (배포 완료)
  ↓
DONE (완전 종료)

추가 상태:
HOLD (일시 정지)
BLOCKED (차단됨)
FAILED (실패)
```

**우선순위:**

- CRITICAL (즉시)
- HIGH (1-2시간)
- MEDIUM (1-2일)
- LOW (1주)

**칼럼:**

```
┌──────────┬──────────────┬────────────┬────────────┬──────────┬──────────┐
│   TODO   │ IN_PROGRESS  │ CODE_REVIEW│ IN_TESTING │ DEPLOYED │   DONE   │
├──────────┼──────────────┼────────────┼────────────┼──────────┼──────────┤
│ 📋 Task  │ 🔄 Task      │ 👀 Task    │ 🧪 Task    │ 🚀 Task  │ ✅ Task  │
│ Priority │ Agent Name   │ Reviewer   │ Test %     │ Env      │ Duration │
│ Due Date │ Progress     │ Blockers   │ Fail Logs  │ Status   │ PR Link  │
└──────────┴──────────────┴────────────┴────────────┴──────────┴──────────┘
```


#### 4.2.2 작업 배정 로직

**수동 배정:**

- PM이 명시적으로 담당자 선택
- 에이전트 선택

**자동 배정:**

- 태그/카테고리 기반 (예: "backend" → Backend Agent)
- 라운드 로빈 (부하 분산)
- 능력 기반 (예: "test" → QA Agent)

**의존성 관리:**

- Parent-Child 작업 관계
- 자동 차단 상태 설정
- 완료 순서 강제

**데이터 모델:**

```sql
CREATE TABLE tasks (
  id VARCHAR(36) PRIMARY KEY,
  title VARCHAR(255),
  description LONGTEXT,
  
  -- Assignment
  assigned_to VARCHAR(36),
  created_by VARCHAR(36),
  
  -- Status & Priority
  status ENUM('TODO', 'IN_PROGRESS', 'CODE_REVIEW', 'IN_TESTING', 'DEPLOYED', 'DONE'),
  priority ENUM('CRITICAL', 'HIGH', 'MEDIUM', 'LOW'),
  
  -- Tracking
  estimated_hours INT,
  actual_hours INT,
  started_at TIMESTAMP,
  completed_at TIMESTAMP,
  
  -- Git Integration
  git_branch VARCHAR(255),
  git_commit_hash VARCHAR(40),
  git_pr_url VARCHAR(500),
  git_pr_status ENUM('DRAFT', 'OPEN', 'APPROVED', 'MERGED'),
  
  -- Documents
  notion_page_id VARCHAR(255),
  markdown_path VARCHAR(500),
  
  -- Relationships
  project_id VARCHAR(36),
  parent_task_id VARCHAR(36),
  
  -- Results
  test_result JSON,
  build_log LONGTEXT,
  deployment_log LONGTEXT,
  
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);
```


#### 4.2.3 작업 히스토리 추적

**기록 항목:**

- 상태 변화 (TODO → IN_PROGRESS)
- 담당자 변경
- 우선순위 변경
- 댓글/토론
- Git 이벤트 (커밋, PR)

**타임라인 뷰:**

- 크로노로지 순서
- 각 이벤트의 상세 정보
- 담당자/변경자 표시
- 소요 시간

***

### 4.3 Git 연동 (Git Integration)

#### 4.3.1 에이전트별 브랜치 관리

**자동 브랜치 생성:**

```
형식: feat/task-{taskId}-{taskName}
예시: feat/task-123-user-authentication

브랜치 전략:
├── main (프로덕션)
├── develop (개발)
└── feat/ (각 에이전트의 작업 브랜치)
```

**브랜치 생명주기:**

1. 작업 할당 → 브랜치 생성
2. 작업 진행 → 커밋/푸시
3. 완료 → PR 생성
4. 리뷰/테스트 → 병합 또는 반려
5. 병합 → 브랜치 삭제

#### 4.3.2 자동 커밋 + PR 생성

**커밋 메시지 형식:**

```
[Agent] 작업 제목

작업 ID: task-123
에이전트: Backend Developer Agent
소요시간: 3m 45s

수정 파일:
- src/auth/index.ts (+150 lines)
- src/auth/login.ts (+80 lines)
- tests/auth.test.ts (+50 lines)

테스트 결과: 45/45 통과
빌드 상태: ✅ 성공
```

**PR 생성:**

```
제목: [Agent] 작업 제목
본문:
- 작업 설명
- 변경 사항
- 테스트 커버리지
- 빌드 로그
- 관련 이슈
```


#### 4.3.3 머지 승인 워크플로우

**옵션:**

1. 자동 머지 (에이전트 권한 있을 때)
2. 수동 리뷰 (선택된 리뷰어)
3. 승인 대기 (PR 상태로 유지)

**병합 조건:**

- [ ] 모든 테스트 통과
- [ ] 빌드 성공
- [ ] 코드 리뷰 승인 (필요 시)
- [ ] 충돌 없음


#### 4.3.4 작업별 Git 히스토리 연결

**추적 정보:**

- 커밋 해시
- 작성자 (에이전트)
- 작성 시간
- 파일 변경 사항
- 라인 추가/삭제
- PR 링크

**조회 API:**

```http
GET /api/tasks/{taskId}/git-history
GET /api/tasks/{taskId}/commits
GET /api/tasks/{taskId}/pr
```


***

### 4.4 문서 생성 (Document Generation)

#### 4.4.1 Notion 연동

**자동 생성:**

- 작업 생성 시 Notion 페이지 자동 생성
- PRD, API Spec, 테스트 시나리오 작성
- 완료 시 자동 아카이빙

**양방향 동기화:**

- Notion 수정 → Task 업데이트
- Task 업데이트 → Notion 반영

**템플릿:**

```
📄 PRD Template
├── Overview
├── Features
├── API Specification
├── Database Schema
└── Implementation Timeline

📝 API Spec Template
├── Endpoints
├── Request/Response
├── Error Codes
└── Examples

🧪 Test Scenario Template
├── Unit Tests
├── Integration Tests
├── E2E Tests
└── Performance Tests
```


#### 4.4.2 Markdown 파일 생성

**자동 생성:**

- README.md (프로젝트)
- CHANGELOG.md (변경사항)
- API.md (API 문서)
- TEST_REPORT.md (테스트 결과)

**GitHub Wiki 업데이트:**

- 문서 자동 푸시
- 인덱스 자동 생성
- 버전 관리


#### 4.4.3 템플릿 시스템

**제공 템플릿:**

- 마이크로서비스 설계
- REST API 스펙
- 데이터베이스 스키마
- 배포 체크리스트
- 성능 테스트 시나리오

**커스텀 템플릿:**

- 조직별 커스텀 가능
- 프롬프트 템플릿 저장
- 버전 관리

***

### 4.5 작업 로깅 \& 모니터링 (Activity Log)

#### 4.5.1 상세 로그 기록

**기록 항목:**

```json
{
  "timestamp": "2025-01-15T14:30:00Z",
  "task_id": "task-123",
  "agent_id": "agent-456",
  "action": "CODE_WRITTEN",
  
  "details": {
    "files_modified": 3,
    "lines_added": 150,
    "lines_deleted": 30,
    "git_commit": "abc1234def567",
    "duration_seconds": 225,
    "tokens_used": 5000,
    "cost_usd": 0.15
  },
  
  "metadata": {
    "provider": "CLAUDE",
    "model": "claude-3-5-sonnet",
    "temperature": 0.7,
    "completion_tokens": 4200
  }
}
```


#### 4.5.2 타임라인 뷰

**표시 정보:**

```
Timeline View:
─────────────────────────────────────

14:30:00 ▶ Backend Developer
         📝 Task assigned
         → "User Authentication"

14:30:15 ▶ Prism Daemon
         🚀 Started execution
         → Claude CLI v2.0

14:32:00 ▶ Claude
         ✍️  Code written
         → 3 files modified (+150 lines)

14:33:00 ▶ Prism Daemon
         🧪 Tests running
         → npm run test

14:33:45 ▶ Test Result
         ✅ 45/45 tests passed
         → Coverage: 92%

14:34:00 ▶ Build
         🔨 Building...
         → npm run build

14:34:30 ▶ Build Complete
         ✅ Build succeeded
         → Size: 245KB

14:35:00 ▶ Git
         🔀 Commit & Push
         → feat/task-123-auth

14:35:15 ▶ GitHub
         🔗 PR Created
         → #42: User Authentication

14:35:30 ▶ Portal
         ✅ Task Completed
         → Status: IN_REVIEW
```


#### 4.5.3 작업 추적 레포트 (Markdown 내보내기)

**내보내기 형식:**

```markdown
# 작업 추적 레포트
작업 ID: task-123
작업명: User Authentication
담당자: Backend Developer Agent
기간: 2025-01-15 14:30 ~ 14:35 (5분 5초)

## 요약
- 상태: IN_REVIEW
- 파일 변경: 3개
- 라인 추가: +150
- 테스트: 45/45 ✅
- 빌드: ✅ 성공

## 타임라인
| 시간 | 작업 | 상태 |
|------|------|------|
| 14:30 | 작업 시작 | ✅ |
| 14:32 | 코드 작성 | ✅ |
| 14:33 | 테스트 | ✅ |
| 14:34 | 빌드 | ✅ |
| 14:35 | Git Push | ✅ |

## Git 정보
- Commit: abc1234def567
- Branch: feat/task-123-auth
- PR: #42

## 비용
- 토큰 사용: 5,000
- 비용: $0.15
```


#### 4.5.4 실시간 모니터링

**대시보드:**

- 활성 작업 수
- 에이전트별 작업량
- 성공/실패율
- 평균 소요시간
- 비용 추이

**알림:**

- 작업 완료
- 테스트 실패
- 빌드 실패
- PR 충돌
- 예산 초과

***

### 4.6 E2E 테스트 \& 배포 (CI/CD Integration)

#### 4.6.1 에이전트가 테스트 실행 가능

**자동 테스트 실행:**

```bash
# Daemon에서 자동 실행
npm run test
pytest tests/
go test ./...

# 결과 수집
coverage: 92%
passed: 45/45
failed: 0
duration: 15s
```

**테스트 결과 기록:**

```json
{
  "task_id": "task-123",
  "test_framework": "jest",
  "total_tests": 45,
  "passed": 45,
  "failed": 0,
  "skipped": 0,
  "coverage": 92,
  "duration_ms": 15000,
  "details": {
    "unit_tests": { "passed": 30, "failed": 0 },
    "integration_tests": { "passed": 15, "failed": 0 },
    "e2e_tests": { "passed": 0, "skipped": 0 }
  }
}
```


#### 4.6.2 테스트 결과 자동 기록

**기록 위치:**

- Task 테이블 (test_result JSON)
- Activity Log (MongoDB)
- GitHub Check (PR에 표시)

**표시:**

```
Pull Request:
  ✅ Tests: 45/45 passed (92% coverage)
  ✅ Build: Success (245KB)
  ✅ Lint: No issues
  ⏳ E2E: Running... (5분 50초 남음)
```


#### 4.6.3 자동 배포 트리거

**배포 조건:**

```
IF task.status == "IN_REVIEW" 
  AND pr.merged == true
  AND test.passed == true
  AND build.success == true
THEN
  trigger_deployment()
```

**배포 환경:**

- Staging (자동)
- Production (수동 승인)

**배포 스크립트:**

```bash
# 자동 실행
npm run deploy:staging

# 또는
docker push registry.example.com/service:task-123
kubectl apply -f deployment.yaml
```


#### 4.6.4 배포 상태 추적

**배포 단계:**

```
Deployment Started (14:35:30)
  ↓
Building Docker Image (14:35:45)
  ✅ Built: sha256:abc123...
  ↓
Pushing to Registry (14:36:00)
  ✅ Pushed: registry.example.com/service:task-123
  ↓
Deploying to Staging (14:36:15)
  ✅ Deployment: 3/3 replicas ready
  ↓
Smoke Tests (14:36:30)
  ✅ Health Check: OK
  ✅ API Test: OK
  ↓
Deployment Completed (14:36:45)
  ✅ Duration: 1m 15s
  ✅ Status: SUCCESS
```


***

### 4.7 에이전트 간 협업 (Agent Collaboration)

#### 4.7.1 에이전트 간 메시지 전달

**시나리오:**

```
Backend Dev Agent 작업 완료
  ↓
"API 준비됨. Frontend에서 호출 가능"
  ↓
Frontend Dev Agent에 자동 알림
  ↓
Frontend Dev Agent가 해당 작업 시작
```

**메시지 시스템:**

```
Task Comment:
┌──────────────────────────────┐
│ Backend Developer Agent      │
│ 2025-01-15 14:35:00          │
├──────────────────────────────┤
│ ✅ API 엔드포인트 완료:      │
│ POST /api/auth/login         │
│ POST /api/auth/register      │
│                              │
│ Frontend에서 호출 가능합니다 │
│                              │
│ @Frontend Developer Agent    │
└──────────────────────────────┘

@Mention 알림:
→ Frontend Developer Agent에 통지
→ Dashboard에 "Dependency Ready" 표시
```


#### 4.7.2 병렬 작업 관리

**작업 의존성:**

```
Task-A (Backend API)
  ↓
Task-B (Frontend UI) ─┐
  ↓                  │
Task-C (E2E Test) ←──┘

Parallel Execution:
- Task-B와 Task-D 동시 실행 가능
- Task-E는 Task-B 완료 후 시작
```

**블로킹 상태:**

- Task가 다른 Task를 기다리는 경우
- 자동으로 BLOCKED 상태
- 의존성 완료 시 자동 시작


#### 4.7.3 의존성 표시

**UI 표현:**

```
Task Board:

Task-A (Backend API)
  │
  └─── DEPENDS ON ──→ Task-B (Frontend UI)
       BLOCKS ◀────── Task-C (E2E Test)

Gantt Chart:
[====== Task-A ======]
            [== Task-B ==]
                    [==== Task-C ====]
```


***

## 5. 사용자 흐름 (User Flow)

### 5.1 초기 설정

```
Step 1: Administrator가 Portal에서 에이전트 생성
────────────────────────────────────────────

웹 UI: "New Agent"
  → Name: "Backend Developer"
  → Role: "BACKEND"
  → Provider: "Claude"
  → Model: "claude-3-5-sonnet"
  → System Prompt: "[설정]"
  → Capabilities: Code, Tests, Deploy
  → Save

Step 2: 각 개발자가 로컬 Daemon 설치
────────────────────────────────────────────

$ npx loom-daemon setup
  ? Portal Server URL: https://portal.example.com
  ? API Token: [사용자 토큰]
  ? Project Path: ~/my-project
  ? Project Name: my-project
  
✅ Daemon installed
✅ System service enabled
✅ Auto-start configured

Step 3: 확인 (선택)
────────────────────────────────────────────

$ npx loom-desktop
  → Local status 확인
  → Daemon 연결 상태 확인
  → 작업 로그 확인
  (언제든 닫아도 됨)
```


### 5.2 일상 작업 흐름

```
Step 1: PM이 웹 UI에서 작업 생성
────────────────────────────────────────────

웹 UI (Portal Shell)
  → "New Task"
  → Title: "User Authentication"
  → Description: "[요구사항]"
  → Priority: HIGH
  → Agent: "Backend Developer Agent"
  → Project: "my-project"
  → Save

Step 2: 자동 처리
────────────────────────────────────────────

Portal Server:
  1. Task 생성 (DB 저장)
  2. 상태: TODO → IN_PROGRESS
  3. Redis에 발행: "tasks:backend-agent"

Daemon (User PC - 자동 수신):
  1. Redis 메시지 수신
  2. Git 브랜치 생성: feat/task-123-auth
  3. Claude CLI 실행
  4. 코드 작성 (로컬 파일)
  5. npm run test (테스트 자동 실행)
  6. npm run build (빌드 검증)
  7. Git commit & push
  8. GitHub PR 생성
  9. 결과 전송

Step 3: Portal에서 결과 표시
────────────────────────────────────────────

웹 UI:
  Task Status: IN_REVIEW
  ✅ Tests: 45/45
  ✅ Build: Success
  ✅ Commit: abc1234
  ✅ PR: #42

Timeline:
  14:35:00 Task Started
  14:35:15 Claude Running
  14:32:00 Files Modified (3)
  14:33:00 Tests Passed
  14:34:00 Build Completed
  14:35:00 PR Created

Step 4: 리뷰 및 완료
────────────────────────────────────────────

Reviewer가 PR 확인
  → "Looks good!"
  → Merge PR

Task 자동 업데이트:
  Status: DEPLOYED
  Merged: true
  Merged At: 2025-01-15 14:40:00

완료!
```


***

## 6. 데이터 모델

### 6.1 핵심 테이블

```sql
-- 에이전트
CREATE TABLE agents (
  id VARCHAR(36) PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  role ENUM('PM', 'BACKEND', 'FRONTEND', 'QA', 'DEVOPS', 'CUSTOM'),
  description TEXT,
  provider_id VARCHAR(36),
  model_name VARCHAR(100),
  system_prompt LONGTEXT,
  temperature DECIMAL(2, 1),
  max_tokens INT,
  can_write_code BOOLEAN DEFAULT TRUE,
  can_run_tests BOOLEAN DEFAULT TRUE,
  can_deploy BOOLEAN DEFAULT FALSE,
  can_create_documents BOOLEAN DEFAULT TRUE,
  can_merge_pr BOOLEAN DEFAULT FALSE,
  organization_id VARCHAR(36),
  created_by VARCHAR(36),
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  FOREIGN KEY (provider_id) REFERENCES ai_providers(id),
  FOREIGN KEY (organization_id) REFERENCES organizations(id)
);

-- 작업
CREATE TABLE tasks (
  id VARCHAR(36) PRIMARY KEY,
  title VARCHAR(255),
  description LONGTEXT,
  assigned_to VARCHAR(36),
  created_by VARCHAR(36),
  status ENUM('TODO', 'IN_PROGRESS', 'CODE_REVIEW', 'IN_TESTING', 'DEPLOYED', 'DONE'),
  priority ENUM('CRITICAL', 'HIGH', 'MEDIUM', 'LOW'),
  estimated_hours INT,
  actual_hours INT,
  started_at TIMESTAMP,
  completed_at TIMESTAMP,
  git_branch VARCHAR(255),
  git_commit_hash VARCHAR(40),
  git_pr_url VARCHAR(500),
  git_pr_status ENUM('DRAFT', 'OPEN', 'APPROVED', 'MERGED'),
  notion_page_id VARCHAR(255),
  markdown_path VARCHAR(500),
  project_id VARCHAR(36),
  parent_task_id VARCHAR(36),
  test_result JSON,
  build_log LONGTEXT,
  deployment_log LONGTEXT,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  FOREIGN KEY (assigned_to) REFERENCES agents(id),
  FOREIGN KEY (created_by) REFERENCES users(id),
  FOREIGN KEY (project_id) REFERENCES projects(id)
);

-- 활동 로그 (MongoDB)
db.activities.insertOne({
  _id: ObjectId(),
  task_id: "task-123",
  agent_id: "agent-456",
  action: "CODE_WRITTEN",
  timestamp: ISODate("2025-01-15T14:30:00Z"),
  details: {
    files_modified: 3,
    lines_added: 150,
    lines_deleted: 30,
    git_commit: "abc1234",
    duration_seconds: 225,
    tokens_used: 5000
  }
});
```


### 6.2 관계도

```
Organization
  ├─ Projects
  │   └─ Tasks
  │       ├─ Agents (assigned_to)
  │       ├─ Comments
  │       ├─ Git Operations
  │       └─ Test Results
  ├─ Agents
  │   ├─ AI Providers
  │   ├─ System Prompts
  │   └─ Task History
  └─ Users
      ├─ Team Members
      ├─ Permissions
      └─ Activity Logs
```


***

## 7. API 스펙

### 7.1 Agent Management APIs

```http
# 에이전트 생성
POST /api/agents
Content-Type: application/json

{
  "name": "Backend Developer",
  "role": "BACKEND",
  "provider": "CLAUDE",
  "model_name": "claude-3-5-sonnet",
  "system_prompt": "You are an expert backend developer...",
  "temperature": 0.7,
  "max_tokens": 4096
}

# Response
{
  "id": "agent-789",
  "name": "Backend Developer",
  "status": "active",
  "created_at": "2025-01-15T14:00:00Z"
}

# 에이전트 목록
GET /api/agents

# 에이전트 상세
GET /api/agents/{agentId}

# 에이전트 수정
PUT /api/agents/{agentId}

# 에이전트 삭제
DELETE /api/agents/{agentId}
```


### 7.2 Task Management APIs

```http
# 작업 생성
POST /api/tasks
{
  "title": "User Authentication",
  "description": "Implement login and registration",
  "assigned_to": "agent-789",
  "priority": "HIGH",
  "project_id": "proj-123",
  "estimated_hours": 2
}

# Response
{
  "id": "task-456",
  "status": "TODO",
  "created_at": "2025-01-15T14:00:00Z"
}

# 작업 목록
GET /api/tasks?status=TODO&priority=HIGH

# 작업 상세
GET /api/tasks/{taskId}

# 작업 업데이트
PUT /api/tasks/{taskId}

# 작업 배정
POST /api/tasks/{taskId}/assign
{
  "agent_id": "agent-789"
}
```


### 7.3 Activity Log APIs

```http
# 활동 로그 조회
GET /api/tasks/{taskId}/activities

# Response
{
  "activities": [
    {
      "timestamp": "2025-01-15T14:35:00Z",
      "agent": "Backend Developer",
      "action": "CODE_WRITTEN",
      "details": {
        "files_modified": 3,
        "lines_added": 150
      }
    }
  ],
  "timeline": "..."
}

# Markdown 내보내기
GET /api/tasks/{taskId}/activities/export?format=markdown

# 타임라인 뷰
GET /api/tasks/{taskId}/timeline
```


***

## 8. 기술 스택

### 8.1 Frontend

| 기술 | 용도 | 이유 |
| :-- | :-- | :-- |
| **Vue 3** | UI 프레임워크 | Portal Universe 통합 |
| **Module Federation** | 마이크로 프론트엔드 | 기능별 독립 배포 |
| **WebSocket** | 실시간 업데이트 | 타임라인 실시간 표시 |
| **Tailwind CSS** | 스타일링 | 빠른 UI 개발 |
| **Vite** | 빌드 도구 | 고속 개발 |

### 8.2 Backend

| 기술 | 용도 | 이유 |
| :-- | :-- | :-- |
| **Spring Boot 3** | 웹 프레임워크 | MSA 패턴 |
| **gRPC** | 내부 통신 | 고성능, 타입안전 |
| **REST** | 외부 API | 표준 인터페이스 |
| **MySQL** | 구조화 데이터 | 트랜잭션 지원 |
| **MongoDB** | 활동 로그 | 비구조 데이터 |
| **Redis** | 작업 큐, 캐시 | 고속 처리 |
| **Kafka** | 메시지 큐 | 비동기 통신 |

### 8.3 로컬 Daemon

| 기술 | 용도 | 이유 |
| :-- | :-- | :-- |
| **Rust 또는 Go** | CLI 도구 | 경량, 빠름 |
| **gRPC** | 서버 통신 | 안정적 |
| **Redis Client** | 작업 큐 수신 | 간단함 |
| **Git2 Lib** | Git 작업 | 자동화 |
| **systemd** | 시스템 서비스 | 자동 시작 |

### 8.4 External APIs

| API | 용도 | 이유 |
| :-- | :-- | :-- |
| **GitHub API** | Git 작업 | PR 자동 생성 |
| **Claude API** | LLM | 코드 생성 |
| **Gemini API** | LLM | 다중 지원 |
| **OpenAI API** | LLM | GPT-4 지원 |
| **Notion API** | 문서 | 자동 문서화 |
| **Slack API** | 알림 | 팀 통지 |

### 8.5 Deployment

| 기술 | 용도 | 이유 |
| :-- | :-- | :-- |
| **Docker** | 컨테이너화 | 이식성 |
| **Kubernetes** | 오케스트레이션 | 확장성 |
| **GitHub Actions** | CI/CD | 자동 배포 |
| **Prometheus** | 모니터링 | 성능 추적 |
| **Grafana** | 대시보드 | 시각화 |


***

## 9. 구현 로드맵

### Phase 1: MVP (8주)

**주 1-2: 프로젝트 설정 \& 기본 구조**

- [ ] Repository 구조 설정
- [ ] Spring Boot 프로젝트 초기화
- [ ] Database 스키마 정의
- [ ] gRPC 정의 (.proto 파일)

**주 3-4: Agent Manager Service**

- [ ] 에이전트 CRUD API
- [ ] AI Provider 관리
- [ ] 프롬프트 템플릿 관리
- [ ] 데이터베이스 연동

**주 5-6: Task Orchestrator \& Daemon 기본**

- [ ] Task CRUD API
- [ ] Redis 큐 구현
- [ ] Prism Daemon 프로토타입 (Rust/Go)
- [ ] Daemon ↔ Server 통신 (gRPC)

**주 7-8: 기본 UI \& 통합**

- [ ] Web UI - Agent Dashboard
- [ ] Web UI - Task Board (기본)
- [ ] WebSocket 실시간 업데이트
- [ ] E2E 테스트

**Deliverables:**

- ✅ 에이전트 생성/관리 가능
- ✅ 작업 생성/배정 가능
- ✅ Daemon이 작업 수신 및 실행 가능
- ✅ 기본 UI에서 모니터링 가능

***

### Phase 2: 고도화 (6주)

**주 1-2: Git Integration Service**

- [ ] GitHub API 통합
- [ ] 자동 브랜치 생성
- [ ] PR 자동 생성
- [ ] 커밋 메시지 자동 생성

**주 3: Document Generator Service**

- [ ] Notion API 통합
- [ ] Markdown 생성
- [ ] 템플릿 시스템

**주 4-5: Activity Logger \& Advanced UI**

- [ ] MongoDB 활동 로그
- [ ] 타임라인 뷰
- [ ] 고급 필터링
- [ ] Markdown 내보내기

**주 6: E2E 테스트 \& 배포 자동화**

- [ ] E2E 테스트 통합
- [ ] 배포 자동화
- [ ] 성공/실패 추적

**Deliverables:**

- ✅ 자동 Git 작업
- ✅ 문서 자동 생성
- ✅ 활동 로그 추적
- ✅ E2E 테스트 자동화

***

### Phase 3: 팀 협업 \& 최적화 (4주)

**주 1-2: 에이전트 간 협업**

- [ ] 메시지 시스템
- [ ] Task 의존성
- [ ] 병렬 실행 관리

**주 3: 대시보드 \& 모니터링**

- [ ] 통계 대시보드
- [ ] 성능 모니터링
- [ ] 비용 추적

**주 4: 최적화 \& 문서**

- [ ] 성능 최적화
- [ ] 보안 감사
- [ ] 완전한 문서화

**Deliverables:**

- ✅ 완전한 협업 기능
- ✅ 성능 최적화
- ✅ Production-ready

***

## 10. 성공 지표 (KPI)

### 10.1 기술 지표

| 지표 | 목표 | 측정 방법 |
| :-- | :-- | :-- |
| **API 응답시간** | < 200ms | APM 도구 |
| **작업 완료율** | > 95% | Task 상태 추적 |
| **테스트 커버리지** | > 80% | Code coverage 도구 |
| **배포 빈도** | 1회/일 | CI/CD 로그 |
| **평균 소요시간** | < 10분 | Activity log |

### 10.2 비즈니스 지표

| 지표 | 목표 | 측정 방법 |
| :-- | :-- | :-- |
| **개발 생산성 증가** | +30% | 완료된 작업 수 |
| **버그율 감소** | -50% | 테스트 결과 |
| **배포 실패율** | < 5% | 배포 로그 |
| **팀 만족도** | > 4.0/5.0 | 설문조사 |

### 10.3 비용 지표

| 항목 | 예상 |
| :-- | :-- |
| **LLM API 비용** | \$500-2000/월 |
| **서버 비용** | \$1000-2000/월 |
| **ROI** | 3-6개월 |


***

## 11. 리스크 및 완화 전략

| 리스크 | 영향 | 완화 방안 |
| :-- | :-- | :-- |
| **LLM API 레이트 제한** | 작업 중단 | 큐 시스템, 재시도 로직, 배치 처리 |
| **Git 충돌** | 머지 실패 | 자동 충돌 해결, PM 알림 |
| **무단 코드 수정** | 보안 위협 | 감시 모드, 승인 워크플로우 |
| **비용 폭증** | 예산 오버 | 토큰 모니터링, 제한 설정 |
| **네트워크 끊김** | Daemon 실패 | 오프라인 큐, 자동 재시도 |
| **동시성 문제** | 데이터 손상 | 트랜잭션, 분산 락 |


***

## 12. 보안 \& 규정 준수

### 12.1 보안 조치

- [ ] API 키 암호화 저장 (AES-256)
- [ ] gRPC mTLS 인증
- [ ] 역할 기반 접근 제어 (RBAC)
- [ ] 감시 로깅 (모든 작업 기록)
- [ ] 정기 보안 감사


### 12.2 규정 준수

- [ ] GDPR 준수 (데이터 삭제 권리)
- [ ] SOC 2 Type II 준비
- [ ] 데이터 암호화 (저장/전송)
- [ ] 감사 추적 유지

***

## 13. 결론

**Prism**은 AI 에이전트를 중앙에서 관리하면서도, 각 개발자의 로컬 환경에서 실제로 작동하도록 설계된 엔터프라이즈 플랫폼입니다.

**핵심 혁신:**

1. **로컬 Daemon 기반**: 실제 로컬 파일 수정, 테스트, 빌드 검증
2. **완전 자동화**: 요청부터 배포까지 완전 자동
3. **팀 협업**: 중앙 제어 + 분산 실행
4. **투명한 추적**: 모든 작업의 타임라인 기록
5. **엔터프라이즈 준비**: 확장성, 보안, 규정 준수

**기대 효과:**

- 개발 생산성 30% 증가
- 배포 빈도 10배 증가
- 버그율 50% 감소
- 개발자 만족도 향상

</parameter>
</invoke>

