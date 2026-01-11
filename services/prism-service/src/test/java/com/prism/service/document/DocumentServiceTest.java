package com.prism.service.document;

import com.prism.domain.ActivityLog;
import com.prism.domain.Task;
import com.prism.domain.Agent;
import com.prism.domain.Project;
import com.prism.repository.ActivityLogRepository;
import com.prism.repository.TaskRepository;
import com.prism.service.integration.NotionClient;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.time.LocalDateTime;
import java.util.List;
import java.util.Map;
import java.util.Optional;

import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.ArgumentMatchers.contains;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class DocumentServiceTest {

    @Mock
    private TaskRepository taskRepository;
    @Mock
    private ActivityLogRepository activityLogRepository;
    @Mock
    private NotionClient notionClient;

    @InjectMocks
    private DocumentService documentService;

    @Test
    void publishTaskToNotion_ShouldGenerateMarkdown() {
        // Arrange
        String taskId = "task-123";
        Task task = Task.builder()
            .id(taskId)
            .title("Test Task")
            .status("DONE")
            .priority("HIGH")
            .description("A sample description")
            .gitBranch("feat/test-task")
            .gitCommitHash("abc1234")
            .project(Project.builder().name("Prism").build())
            .assignedTo(Agent.builder().name("TestBot").build())
            .build();

        ActivityLog log = ActivityLog.builder()
            .taskId(taskId)
            .action("TASK_CREATED")
            .timestamp(LocalDateTime.now())
            .details(Map.of("info", "created"))
            .build();

        when(taskRepository.findById(taskId)).thenReturn(Optional.of(task));
        when(activityLogRepository.findAll()).thenReturn(List.of(log)); // In real code it filters stream
        when(notionClient.createPage(anyString(), anyString())).thenReturn("page-id-123");

        // Act
        documentService.publishTaskToNotion(taskId);

        // Assert
        verify(notionClient).createPage(eq("Task: Test Task"), contains("## Git Integration"));
        verify(notionClient).createPage(eq("Task: Test Task"), contains("- **Branch:** `feat/test-task`"));
    }
}
