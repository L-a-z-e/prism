package com.prism.controller;

import com.prism.dto.CreateTaskRequest;
import com.prism.dto.TaskResponse;
import com.prism.service.TaskService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.*;
import io.swagger.v3.oas.annotations.tags.Tag;
import io.swagger.v3.oas.annotations.Operation;

import java.util.List;

@RestController
@RequestMapping("/api/v1/tasks")
@RequiredArgsConstructor
@CrossOrigin(origins = "*")
@Tag(name = "Tasks", description = "Task Management API")
public class TaskController {
    private final TaskService taskService;

    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    @Operation(summary = "Create a new task")
    public TaskResponse createTask(@RequestBody CreateTaskRequest request) {
        return taskService.createTask(request);
    }

    @GetMapping
    @Operation(summary = "Get all tasks with optional filters")
    public List<TaskResponse> getAllTasks(
        @RequestParam(required = false) String status,
        @RequestParam(required = false) String priority,
        @RequestParam(required = false) String agentId
    ) {
        return taskService.getAllTasks(status, priority, agentId);
    }

    @GetMapping("/{taskId}")
    @Operation(summary = "Get task details including activity log")
    public com.prism.dto.TaskDetailResponse getTask(@PathVariable String taskId) {
        return taskService.getTask(taskId);
    }
}
