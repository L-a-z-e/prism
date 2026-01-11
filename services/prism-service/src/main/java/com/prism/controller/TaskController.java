package com.prism.controller;

import com.prism.dto.CreateTaskRequest;
import com.prism.dto.TaskResponse;
import com.prism.service.TaskService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/v1/tasks")
@RequiredArgsConstructor
@CrossOrigin(origins = "*")
public class TaskController {
    private final TaskService taskService;

    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    public TaskResponse createTask(@RequestBody CreateTaskRequest request) {
        return taskService.createTask(request);
    }

    @GetMapping
    public List<TaskResponse> getAllTasks(
        @RequestParam(required = false) String status,
        @RequestParam(required = false) String priority,
        @RequestParam(required = false) String agentId
    ) {
        return taskService.getAllTasks(status, priority, agentId);
    }

    @GetMapping("/{taskId}")
    public com.prism.dto.TaskDetailResponse getTask(@PathVariable String taskId) {
        return taskService.getTask(taskId);
    }
}
