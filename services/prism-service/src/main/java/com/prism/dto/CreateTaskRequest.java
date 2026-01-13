package com.prism.dto;

import lombok.Data;

/**
 * Task 생성 요청 DTO (3-Phase Git Workflow 지원)
 * 
 * Phase 1: AI 코드 생성 (필수)
 * Phase 2: Git 커밋 (선택)
 * Phase 3: Push & PR (선택)
 */
@Data
public class CreateTaskRequest {
    // ========== 기본 정보 (Phase 1 필수) ==========
    
    private String title;
    private String description;
    private String priority;       // LOW, MEDIUM, HIGH
    private String projectId;
    private String assignedTo;     // Agent ID
    
    // ========== Local 작업 관련 (Phase 1) ==========
    
    /**
     * 로컬 프로젝트 경로
     * 예: "/home/dev/projects/project-A"
     */
    private String projectPath;
    
    // ========== Git 설정 (Phase 2, 3 선택사항) ==========
    
    /**
     * 대상 GitHub 저장소
     * 예: "L-a-z-e/project-a"
     */
    private String targetRepo;
    
    /**
     * Git 브랜치명 (지정하지 않으면 자동 생성: feat/task-{taskId})
     */
    private String gitBranch;
    
    // ========== 자동화 옵션 ==========
    
    /**
     * 자동 커밋 여부 (기본: false)
     * true일 경우 AI 코드 생성 후 자동으로 커밋
     */
    private Boolean autoCommit;
    
    /**
     * 자동 Push & PR 생성 여부 (기본: false)
     * true일 경우 커밋 후 자동으로 Push + PR 생성
     * 주의: autoCommit가 false면 이 옵션은 무시됨
     */
    private Boolean autoPush;
    
    /**
     * PR 제목 (지정하지 않으면 Task title 사용)
     */
    private String prTitle;
    
    /**
     * PR 베이스 브랜치 (기본: "main" 또는 "develop")
     */
    private String prBaseBranch;
}
