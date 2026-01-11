package com.prism.service;

import com.prism.domain.Task;
import com.prism.domain.Project;
import com.prism.domain.Agent;
import com.prism.domain.ActivityLog;
import com.prism.domain.User;
import com.prism.dto.CreateTaskRequest;
import com.prism.dto.TaskResponse;
import com.prism.repository.TaskRepository;
import com.prism.repository.ProjectRepository;
import com.prism.repository.AgentRepository;
import com.prism.repository.ActivityLogRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.redis.core.StringRedisTemplate;
import java.util.Map;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.stream.Collectors;

@Slf4j
@Service
@RequiredArgsConstructor
public class TaskService {
    private final TaskRepository taskRepository;
    private final ProjectRepository projectRepository;
    private final AgentRepository agentRepository;
    private final ActivityLogRepository activityLogRepository;
    private final MockUserService mockUserService;
    private final StringRedisTemplate redisTemplate;

    @Transactional
    public TaskResponse createTask(CreateTaskRequest request) {
        User currentUser = mockUserService.getCurrentUser();

        Project project;
        if (request.getProjectId() == null || request.getProjectId().isEmpty()) {
            // MVP: Use the default project created by MockUserService
            // Assuming it's the one linked to the user's organization or just the first one found
            project = projectRepository.findAll().stream().findFirst()
                .orElseThrow(() -> new IllegalStateException("No projects found. Init failed?"));
        } else {
            project = projectRepository.findById(request.getProjectId())
                .orElseThrow(() -> new IllegalArgumentException("Invalid Project ID"));
        }

        Agent assignedAgent = null;
        if (request.getAssignedTo() != null) {
            assignedAgent = agentRepository.findById(request.getAssignedTo())
                .orElseThrow(() -> new IllegalArgumentException("Invalid Agent ID"));
        }

        Task task = Task.builder()
            .title(request.getTitle())
            .description(request.getDescription())
            .priority(request.getPriority())
            .project(project)
            .assignedTo(assignedAgent)
            .createdBy(currentUser)
            .build();

        task = taskRepository.save(task);

        if (assignedAgent != null) {
            publishTaskToAgent(task, assignedAgent);
        }

        activityLogRepository.save(ActivityLog.builder()
            .taskId(task.getId())
            .userId(currentUser.getId())
            .action("TASK_CREATED")
            .details(Map.of(
                "title", task.getTitle(),
                "assigned_to", assignedAgent != null ? assignedAgent.getName() : "Unassigned"
            ))
            .build());

        return TaskResponse.from(task);
    }

    private void publishTaskToAgent(Task task, Agent agent) {
        // Format: "taskId:title" - simplified for MVP
        String message = task.getId() + ":" + task.getTitle();
        String channel = "tasks:" + agent.getId();
        log.info("Publishing task {} to channel {}", task.getId(), channel);
        redisTemplate.convertAndSend(channel, message);
    }

    @Transactional(readOnly = true)
    public List<TaskResponse> getAllTasks() {
        return taskRepository.findAll().stream()
            .map(TaskResponse::from)
            .collect(Collectors.toList());
    }
}
