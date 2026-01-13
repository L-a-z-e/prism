package com.prism.dto;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * Phase 3: Push & PR 생성 요청 DTO
 * 
 * 사용 시나리오:
 * 1. 커밋이 완료된 상태
 * 2. 사용자가 "Push & PR 생성" 버튼 클릭
 * 3. POST /api/v1/tasks/{taskId}/push 호출
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class PushRequest {
    
    /**
     * PR 생성 여부 (기본: true)
     * 
     * true일 경우:
     * - Push 후 자동으로 PR 생성
     * 
     * false일 경우:
     * - Push만 수행 (PR은 생성하지 않음)
     * - 나중에 GitHub UI에서 수동으로 PR 생성 가능
     */
    @Builder.Default
    private Boolean createPR = true;
    
    /**
     * PR 제목
     * 예: "feat: Add login components"
     * 
     * 지정하지 않으면 Task title 사용
     */
    private String prTitle;
    
    /**
     * PR 설명 (Body)
     * 예: "Resolves #123\n\n- Added LoginForm component\n- Added LoginModal component"
     * 
     * 지정하지 않으면 Task description 사용
     */
    private String prBody;
    
    /**
     * PR 베이스 브랜치 (기본: "main")
     * 예: "develop", "main", "staging"
     * 
     * 병합될 대상 브랜치 지정
     */
    private String prBaseBranch;
    
    /**
     * PR 리뷰어 목록 (선택사항)
     * 예: ["reviewer1", "reviewer2"]
     */
    private String[] prReviewers;
    
    /**
     * PR 라벨 목록 (선택사항)
     * 예: ["enhancement", "ai-generated"]
     */
    private String[] prLabels;
    
    /**
     * 자동 Push 여부 (기본: false)
     * 
     * true일 경우:
     * - 사용자 승인 없이 자동으로 Push + PR 생성
     * 
     * false일 경우:
     * - 사용자가 명시적으로 승인해야 Push 진행
     */
    @Builder.Default
    private Boolean autoPush = false;
}
