package com.prism.domain;

import jakarta.persistence.Entity;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import jakarta.persistence.ManyToOne;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.Column;
import jakarta.persistence.Enumerated;
import jakarta.persistence.EnumType;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import com.prism.domain.enums.GitPhase;
import com.prism.domain.enums.TaskStatus;
import java.time.LocalDateTime;

@Entity
@Table(name = "tasks")
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class Task {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private String id;

    private String title;

    @Column(columnDefinition = "LONGTEXT")
    private String description;

    @ManyToOne
    @JoinColumn(name = "assigned_to")
    private Agent assignedTo;

    @ManyToOne
    @JoinColumn(name = "project_id")
    private Project project;

    @ManyToOne
    @JoinColumn(name = "created_by")
    private User createdBy;

    private String parentTaskId;

    // ========== Phase 1: 3-Phase 모델 추가 ==========
    
    // Status enum (기존 status 필드 대체)
    @Enumerated(EnumType.STRING)
    @Builder.Default
    private TaskStatus taskStatus = TaskStatus.CREATED;
    
    // Git Phase (선택 사항 추적)
    @Enumerated(EnumType.STRING)
    @Builder.Default
    private GitPhase gitPhase = GitPhase.NONE;
    
    // Local 작업 관련
    private String projectPath;              // "/home/dev/projects/project-A"
    private String outputFilePath;           // 생성된 파일 경로
    
    // Git 관련 필드 (선택)
    private String targetRepo;               // "L-a-z-e/project-a"
    private String gitBranch;                // "feat/task-123"
    private String gitCommitHash;            // Commit SHA
    private String gitPrUrl;                 // PR URL
    private String gitPrStatus;
    
    // 자동화 옵션
    @Builder.Default
    private Boolean autoCommit = false;      // 자동 커밋 여부 (기본: false)
    @Builder.Default
    private Boolean autoPush = false;        // 자동 Push 여부 (기본: false)
    
    // 타임스탬프 (Phase별 진행 시간)
    private LocalDateTime generatedAt;       // AI 생성 완료
    private LocalDateTime committedAt;       // Git 커밋 완료
    private LocalDateTime pushedAt;          // Push 완료
    
    // 기존 필드들 (호환성 유지)
    @Builder.Default
    private String status = "TODO";          // 레거시 필드 (향후 제거)
    
    @Builder.Default
    private String priority = "MEDIUM";

    private Integer estimatedHours;
    private Integer actualHours;
    private LocalDateTime startedAt;
    private LocalDateTime completedAt;

    private String notionPageId;
    private String markdownPath;

    @Column(columnDefinition = "JSON")
    private String testResult;

    @Column(columnDefinition = "LONGTEXT")
    private String buildLog;

    @Column(columnDefinition = "LONGTEXT")
    private String deploymentLog;

    @Builder.Default
    private LocalDateTime createdAt = LocalDateTime.now();
    @Builder.Default
    private LocalDateTime updatedAt = LocalDateTime.now();
}
