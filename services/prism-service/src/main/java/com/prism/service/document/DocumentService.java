package com.prism.service.document;

import com.prism.domain.Task;
import com.prism.domain.ActivityLog;
import com.prism.repository.TaskRepository;
import com.prism.repository.ActivityLogRepository;
import com.prism.service.integration.NotionClient;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.format.DateTimeFormatter;
import java.util.List;
import java.util.stream.Collectors;

@Service
@RequiredArgsConstructor
public class DocumentService {
    private final TaskRepository taskRepository;
    private final ActivityLogRepository activityLogRepository;
    private final NotionClient notionClient;

    @Transactional
    public String publishTaskToNotion(String taskId) {
        Task task = taskRepository.findById(taskId)
            .orElseThrow(() -> new IllegalArgumentException("Task not found"));

        String markdown = generateMarkdown(task);
        String pageId = notionClient.createPage("Task: " + task.getTitle(), markdown);

        task.setNotionPageId(pageId);
        taskRepository.save(task);

        return pageId;
    }

    private String generateMarkdown(Task task) {
        StringBuilder sb = new StringBuilder();

        sb.append("# ").append(task.getTitle()).append("\n\n");
        sb.append("**Status:** ").append(task.getStatus()).append("\n");
        sb.append("**Priority:** ").append(task.getPriority()).append("\n");
        sb.append("**Assigned To:** ").append(task.getAssignedTo() != null ? task.getAssignedTo().getName() : "Unassigned").append("\n");
        sb.append("**Project:** ").append(task.getProject() != null ? task.getProject().getName() : "N/A").append("\n\n");

        sb.append("## Description\n");
        sb.append(task.getDescription() != null ? task.getDescription() : "No description").append("\n\n");

        if (task.getGitBranch() != null) {
            sb.append("## Git Integration\n");
            sb.append("- **Branch:** `").append(task.getGitBranch()).append("`\n");
            if (task.getGitCommitHash() != null) {
                sb.append("- **Commit:** `").append(task.getGitCommitHash()).append("`\n");
            }
            if (task.getGitPrUrl() != null) {
                sb.append("- **PR:** [Link](").append(task.getGitPrUrl()).append(")\n");
            }
            sb.append("\n");
        }

        sb.append("## Activity Timeline\n");
        List<ActivityLog> logs = activityLogRepository.findAll().stream()
            .filter(l -> l.getTaskId().equals(task.getId()))
            .sorted((a, b) -> b.getTimestamp().compareTo(a.getTimestamp()))
            .collect(Collectors.toList());

        DateTimeFormatter formatter = DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss");

        for (ActivityLog log : logs) {
            sb.append("- **").append(log.getTimestamp().format(formatter)).append("** - ")
              .append(log.getAction());
             if (log.getDetails() != null) {
                 sb.append(" `").append(log.getDetails().toString()).append("`");
             }
             sb.append("\n");
        }

        return sb.toString();
    }
}
