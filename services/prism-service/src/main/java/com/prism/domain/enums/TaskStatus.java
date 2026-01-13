package com.prism.domain.enums;

/**
 * Task 상태 단계 전넸 (3-Phase 모델)
 * 
 * CREATED: Task 생성됨
 * GENERATING: AI 코드 생성 중
 * GENERATED: 생성 완료, 사용자 검론 대기 중
 * COMMIT_PENDING: 커밋 승인 대기 중
 * COMMITTED: 로컬 커밋 완료, Push 승인 대기 중
 * PUSH_PENDING: Push 승인 대기 중
 * PUSHED: 원격 Push 완료
 * PR_CREATED: GitHub PR 생성됨
 * COMPLETED: Task 완료
 * FAILED: 실패
 */
public enum TaskStatus {
    CREATED,           // Task 생성됨
    GENERATING,        // AI 코드 생성 중
    GENERATED,         // 생성 완료, 사용자 검론 대기
    COMMIT_PENDING,    // 커밋 승인 대기
    COMMITTED,         // 로컬 커밋 완료, Push 승인 대기
    PUSH_PENDING,      // Push 승인 대기
    PUSHED,            // Push 완료
    PR_CREATED,        // PR 생성됨
    COMPLETED,         // Task 완료
    FAILED             // 실패
}
