package com.prism.repository;

import com.prism.domain.Task;
import com.prism.domain.enums.GitPhase;
import com.prism.domain.enums.TaskStatus;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import java.util.List;

@Repository
public interface TaskRepository extends JpaRepository<Task, String> {

    // ========== 레거시 메서드 (기존 호환성 유지) ==========
    
    @Query("SELECT t FROM Task t WHERE " +
           "(:status IS NULL OR t.status = :status) AND " +
           "(:priority IS NULL OR t.priority = :priority) AND " +
           "(:agentId IS NULL OR t.assignedTo.id = :agentId)")
    List<Task> findByFilters(
        @Param("status") String status,
        @Param("priority") String priority,
        @Param("agentId") String agentId
    );

    @Query("SELECT t.status, COUNT(t) FROM Task t GROUP BY t.status")
    List<Object[]> countTasksByStatus();
    
    // ========== Phase 2: 3-Phase 모델 쿼리 메서드 ==========
    
    /**
     * TaskStatus와 GitPhase를 사용한 고급 필터링
     * 
     * @param taskStatus Task 상태 (CREATED, GENERATING, GENERATED, COMMITTED, PUSHED, COMPLETED, FAILED)
     * @param gitPhase Git 단계 (NONE, COMMITTED, PUSHED, PR_CREATED)
     * @param priority 우선순위 (LOW, MEDIUM, HIGH)
     * @param agentId Agent ID
     * @return 필터링된 Task 목록
     */
    @Query("SELECT t FROM Task t WHERE " +
           "(:taskStatus IS NULL OR t.taskStatus = :taskStatus) AND " +
           "(:gitPhase IS NULL OR t.gitPhase = :gitPhase) AND " +
           "(:priority IS NULL OR t.priority = :priority) AND " +
           "(:agentId IS NULL OR t.assignedTo.id = :agentId)")
    List<Task> findByAdvancedFilters(
        @Param("taskStatus") TaskStatus taskStatus,
        @Param("gitPhase") GitPhase gitPhase,
        @Param("priority") String priority,
        @Param("agentId") String agentId
    );
    
    /**
     * TaskStatus별 Task 수 집계
     */
    @Query("SELECT t.taskStatus, COUNT(t) FROM Task t GROUP BY t.taskStatus")
    List<Object[]> countTasksByTaskStatus();
    
    /**
     * GitPhase별 Task 수 집계
     */
    @Query("SELECT t.gitPhase, COUNT(t) FROM Task t GROUP BY t.gitPhase")
    List<Object[]> countTasksByGitPhase();
    
    /**
     * 특정 TaskStatus에 있는 모든 Task 조회
     */
    List<Task> findByTaskStatus(TaskStatus taskStatus);
    
    /**
     * 특정 GitPhase에 있는 모든 Task 조회
     */
    List<Task> findByGitPhase(GitPhase gitPhase);
    
    /**
     * Agent별 + TaskStatus 조회
     */
    @Query("SELECT t FROM Task t WHERE t.assignedTo.id = :agentId AND t.taskStatus = :taskStatus")
    List<Task> findByAgentIdAndTaskStatus(
        @Param("agentId") String agentId,
        @Param("taskStatus") TaskStatus taskStatus
    );
    
    /**
     * 사용자 승인 대기 중인 Task 조회 (GENERATED, COMMIT_PENDING, PUSH_PENDING)
     */
    @Query("SELECT t FROM Task t WHERE t.taskStatus IN ('GENERATED', 'COMMIT_PENDING', 'PUSH_PENDING') " +
           "AND (:agentId IS NULL OR t.assignedTo.id = :agentId) " +
           "ORDER BY t.createdAt DESC")
    List<Task> findPendingApprovalTasks(@Param("agentId") String agentId);
    
    /**
     * 자동화 설정된 Task 조회 (autoCommit 또는 autoPush가 true)
     */
    @Query("SELECT t FROM Task t WHERE (t.autoCommit = true OR t.autoPush = true) " +
           "AND t.taskStatus NOT IN ('COMPLETED', 'FAILED')")
    List<Task> findAutomatedTasks();
    
    /**
     * targetRepo별 Task 조회 (특정 GitHub 저장소 관련 작업)
     */
    List<Task> findByTargetRepo(String targetRepo);
    
    /**
     * PR이 생성된 Task 조회
     */
    @Query("SELECT t FROM Task t WHERE t.gitPhase = 'PR_CREATED' AND t.gitPrUrl IS NOT NULL")
    List<Task> findTasksWithPR();
}
