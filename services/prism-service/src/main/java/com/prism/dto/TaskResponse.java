package com.prism.dto;

import com.prism.domain.Task;
import lombok.Builder;
import lombok.Data;
import java.time.LocalDateTime;

@Data
@Builder
public class TaskResponse {
    private String id;
    private String title;
    private String description;
    private String status;
    private String priority;
    private String assignedToName;
    private String projectName;
    private LocalDateTime createdAt;
    private LocalDateTime startedAt;
    private LocalDateTime completedAt;

    public static TaskResponse from(Task task) {
        return TaskResponse.builder()
            .id(task.getId())
            .title(task.getTitle())
            .description(task.getDescription())
            .status(task.getStatus())
            .priority(task.getPriority())
            .assignedToName(task.getAssignedTo() != null ? task.getAssignedTo().getName() : null)
            .projectName(task.getProject() != null ? task.getProject().getName() : null)
            .createdAt(task.getCreatedAt())
            .startedAt(task.getStartedAt())
            .completedAt(task.getCompletedAt())
            .build();
    }
}
