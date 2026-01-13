package com.prism.dto;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * Phase 2: Git 커밋 승인 요청 DTO
 * 
 * 사용 시나리오:
 * 1. AI가 코드를 생성하고 사용자가 확인
 * 2. 사용자가 "커밋하기" 버튼 클릭
 * 3. POST /api/v1/tasks/{taskId}/commit 호출
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class CommitRequest {
    
    /**
     * Git 커밋 메시지
     * 예: "feat: Add login components"
     * 
     * 지정하지 않으면 기본값 사용:
     * "feat: Auto-generated code for task {taskTitle}"
     */
    private String commitMessage;
    
    /**
     * 커밋 후 자동으로 Push 여부 (기본: false)
     * 
     * true일 경우:
     * - 커밋 완료 후 즉시 Push
     * - 사용자 승인 없이 자동 진행
     * 
     * false일 경우:
     * - 커밋만 수행
     * - 로컬에만 저장
     * - 나중에 별도로 Push 요청 가능
     */
    @Builder.Default
    private Boolean autoCommit = false;
}
