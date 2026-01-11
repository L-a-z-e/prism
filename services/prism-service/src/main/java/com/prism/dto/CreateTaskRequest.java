package com.prism.dto;

import lombok.Data;

@Data
public class CreateTaskRequest {
    private String title;
    private String description;
    private String priority; // MEDIUM, HIGH, etc
    private String projectId;
    private String assignedTo; // Agent ID
}
