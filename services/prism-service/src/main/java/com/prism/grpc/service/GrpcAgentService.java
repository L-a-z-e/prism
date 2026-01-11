package com.prism.grpc.service;

import com.prism.domain.Agent;
import com.prism.domain.Task;
import com.prism.domain.ActivityLog;
import com.prism.grpc.AgentServiceGrpc;
import com.prism.grpc.HeartbeatRequest;
import com.prism.grpc.HeartbeatResponse;
import com.prism.grpc.RegisterAgentRequest;
import com.prism.grpc.RegisterAgentResponse;
import com.prism.grpc.UpdateTaskStatusRequest;
import com.prism.grpc.UpdateTaskStatusResponse;
import com.prism.repository.AgentRepository;
import com.prism.repository.TaskRepository;
import com.prism.repository.ActivityLogRepository;
import io.grpc.stub.StreamObserver;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import java.util.Map;
import net.devh.boot.grpc.server.service.GrpcService;
import org.springframework.messaging.simp.SimpMessagingTemplate;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;

@Slf4j
@GrpcService
@RequiredArgsConstructor
public class GrpcAgentService extends AgentServiceGrpc.AgentServiceImplBase {

    private final AgentRepository agentRepository;
    private final TaskRepository taskRepository;
    private final ActivityLogRepository activityLogRepository;
    private final SimpMessagingTemplate messagingTemplate;

    @Override
    public void registerAgent(RegisterAgentRequest request, StreamObserver<RegisterAgentResponse> responseObserver) {
        log.info("Registering agent: {}", request.getAgentId());
        // Simple logic: Check if exists, update heartbeat or just acknowledge
        // In real app, might update status to ONLINE

        responseObserver.onNext(RegisterAgentResponse.newBuilder()
            .setSuccess(true)
            .setMessage("Agent Registered Successfully")
            .build());
        responseObserver.onCompleted();
    }

    @Override
    public void heartbeat(HeartbeatRequest request, StreamObserver<HeartbeatResponse> responseObserver) {
        // Update last seen timestamp in DB
        responseObserver.onNext(HeartbeatResponse.newBuilder()
            .setAcknowledged(true)
            .build());
        responseObserver.onCompleted();
    }

    @Override
    @Transactional
    public void updateTaskStatus(UpdateTaskStatusRequest request, StreamObserver<UpdateTaskStatusResponse> responseObserver) {
        log.info("Received status update for Task {}: {}", request.getTaskId(), request.getStatus());

        taskRepository.findById(request.getTaskId()).ifPresentOrElse(task -> {
            task.setStatus(request.getStatus());
            if ("DONE".equals(request.getStatus())) {
                task.setCompletedAt(LocalDateTime.now());
            } else if ("IN_PROGRESS".equals(request.getStatus()) && task.getStartedAt() == null) {
                task.setStartedAt(LocalDateTime.now());
            }

            // Append log if present
            if (request.getDetails() != null && !request.getDetails().isEmpty()) {
                String existing = task.getDeploymentLog() == null ? "" : task.getDeploymentLog() + "\n";
                task.setDeploymentLog(existing + request.getDetails());
            }

            taskRepository.save(task);

            activityLogRepository.save(ActivityLog.builder()
                .taskId(task.getId())
                .agentId(request.getAgentId())
                .action("TASK_STATUS_UPDATE")
                .details(Map.of(
                    "status", request.getStatus(),
                    "details", request.getDetails()
                ))
                .build());

            // Broadcast to WebSocket
            messagingTemplate.convertAndSend("/topic/tasks/" + task.getId(), Map.of(
                "type", "STATUS_UPDATE",
                "status", request.getStatus(),
                "details", request.getDetails(),
                "timestamp", java.time.LocalDateTime.now().toString()
            ));

            responseObserver.onNext(UpdateTaskStatusResponse.newBuilder().setSuccess(true).build());
        }, () -> {
            log.warn("Task not found: {}", request.getTaskId());
            responseObserver.onNext(UpdateTaskStatusResponse.newBuilder().setSuccess(false).build());
        });

        responseObserver.onCompleted();
    }
}
