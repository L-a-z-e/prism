package com.prism.domain.enums;

/**
 * Git 샛력 단계 전넸
 * 
 * NONE: Git 작업 없음 (로컬 단계)
 * COMMITTED: 로컬 커밋 만 (Push X)
 * PUSHED: 원격 저장소에 Push된 상태
 * PR_CREATED: PR 생성된 상태
 */
public enum GitPhase {
    NONE,              // Git 작업 없음
    COMMITTED,         // 로컬 커밋만
    PUSHED,            // Push 완료
    PR_CREATED         // PR 생성됨
}
