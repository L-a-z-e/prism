package com.prism.dto;

import com.prism.domain.Task;
import com.prism.domain.enums.GitPhase;
import com.prism.domain.enums.TaskStatus;
import lombok.Builder;
import lombok.Data;
import java.time.LocalDateTime;

/**
 * Task 조회 응답 DTO (3-Phase Git Workflow 지원)
 */
@Data
@Builder
public class TaskResponse {
    // ========== 기본 정보 ==========
    
    private String id;
    private String title;
    private String description;
    private String priority;
    private String assignedToName;
    private String projectName;
    
    // ========== 3-Phase 상태 정보 ==========
    
    /**
     * Task 전체 상태 (TaskStatus enum)
     * CREATED, GENERATING, GENERATED, COMMIT_PENDING, COMMITTED, 
     * PUSH_PENDING, PUSHED, PR_CREATED, COMPLETED, FAILED
     */
    private TaskStatus taskStatus;
    
    /**
     * Git 작업 단계 (GitPhase enum)
     * NONE, COMMITTED, PUSHED, PR_CREATED
     */
    private GitPhase gitPhase;
    
    /**
     * 다음 사용자 액션
     * 예: "COMMIT_APPROVAL", "PUSH_APPROVAL", "COMPLETE"
     * 
     * Frontend에서 어떤 버튼을 보여줄지 결정
     */
    private String nextAction;
    
    // ========== Local 작업 정보 ==========
    
    private String projectPath;
    private String outputFilePath;
    
    // ========== Git 정보 ==========
    
    private String targetRepo;
    private String gitBranch;
    private String gitCommitHash;
    private String gitPrUrl;
    private String gitPrStatus;
    
    // ========== 자동화 옵션 ==========
    
    private Boolean autoCommit;
    private Boolean autoPush;
    
    // ========== 타임스탬프 ==========
    
    private LocalDateTime createdAt;
    private LocalDateTime generatedAt;     // AI 생성 완료
    private LocalDateTime committedAt;     // Git 커밋 완료
    private LocalDateTime pushedAt;        // Push 완료
    private LocalDateTime startedAt;       // 레거시 필드
    private LocalDateTime completedAt;     // Task 전체 완료
    
    // ========== 레거시 필드 (하위 호환성) ==========
    
    @Deprecated
    private String status;  // 기존 status 필드 (향후 제거 예정)
    
    /**
     * Task 엔티티로부터 TaskResponse 생성
     */
    public static TaskResponse from(Task task) {
        // nextAction 결정 로직
        String nextAction = determineNextAction(task);
        
        return TaskResponse.builder()
            .id(task.getId())
            .title(task.getTitle())
            .description(task.getDescription())
            .priority(task.getPriority())
            .assignedToName(task.getAssignedTo() != null ? task.getAssignedTo().getName() : null)
            .projectName(task.getProject() != null ? task.getProject().getName() : null)
            
            // 3-Phase 상태
            .taskStatus(task.getTaskStatus())
            .gitPhase(task.getGitPhase())
            .nextAction(nextAction)
            
            // Local 작업
            .projectPath(task.getProjectPath())
            .outputFilePath(task.getOutputFilePath())
            
            // Git 정보
            .targetRepo(task.getTargetRepo())
            .gitBranch(task.getGitBranch())
            .gitCommitHash(task.getGitCommitHash())
            .gitPrUrl(task.getGitPrUrl())
            .gitPrStatus(task.getGitPrStatus())
            
            // 자동화
            .autoCommit(task.getAutoCommit())
            .autoPush(task.getAutoPush())
            
            // 타임스탬프
            .createdAt(task.getCreatedAt())
            .generatedAt(task.getGeneratedAt())
            .committedAt(task.getCommittedAt())
            .pushedAt(task.getPushedAt())
            .startedAt(task.getStartedAt())
            .completedAt(task.getCompletedAt())
            
            // 레거시
            .status(task.getStatus())
            .build();
    }
    
    /**
     * 현재 Task 상태에 따른 다음 액션 결정
     */
    private static String determineNextAction(Task task) {
        if (task.getTaskStatus() == null) {
            return "UNKNOWN";
        }
        
        switch (task.getTaskStatus()) {
            case CREATED:
            case GENERATING:
                return "WAIT_GENERATION";
                
            case GENERATED:
                // AI 생성 완료 -> 커밋 승인 필요
                return task.getAutoCommit() ? "AUTO_COMMIT" : "COMMIT_APPROVAL";
                
            case COMMIT_PENDING:
                return "COMMIT_APPROVAL";
                
            case COMMITTED:
                // 커밋 완료 -> Push 승인 필요
                return task.getAutoPush() ? "AUTO_PUSH" : "PUSH_APPROVAL";
                
            case PUSH_PENDING:
                return "PUSH_APPROVAL";
                
            case PUSHED:
            case PR_CREATED:
                // Push 또는 PR 생성 완료 -> Task 완료
                return "COMPLETE";
                
            case COMPLETED:
                return "COMPLETED";
                
            case FAILED:
                return "RETRY_OR_CANCEL";
                
            default:
                return "UNKNOWN";
        }
    }
}
