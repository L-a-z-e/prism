package com.prism.domain;

import jakarta.persistence.Entity;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import jakarta.persistence.ManyToOne;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.Column;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
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

    @Builder.Default
    private String status = "TODO";
    @Builder.Default
    private String priority = "MEDIUM";

    private Integer estimatedHours;
    private Integer actualHours;
    private LocalDateTime startedAt;
    private LocalDateTime completedAt;

    private String gitBranch;
    private String gitCommitHash;
    private String gitPrUrl;
    private String gitPrStatus;

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
