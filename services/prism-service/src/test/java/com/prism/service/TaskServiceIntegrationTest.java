package com.prism.service;

import com.prism.domain.ActivityLog;
import com.prism.domain.Task;
import com.prism.dto.CreateTaskRequest;
import com.prism.dto.TaskResponse;
import com.prism.repository.ActivityLogRepository;
import com.prism.repository.AgentRepository;
import com.prism.repository.ProjectRepository;
import com.prism.repository.TaskRepository;
import com.prism.repository.UserRepository;
import com.prism.repository.OrganizationRepository;
import com.prism.repository.AiProviderRepository;
import com.prism.domain.Organization;
import com.prism.domain.User;
import com.prism.domain.Project;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.mock.mockito.MockBean;
import org.springframework.data.redis.core.StringRedisTemplate;
import org.springframework.messaging.simp.SimpMessagingTemplate;
import org.springframework.test.context.ActiveProfiles;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

@SpringBootTest
@ActiveProfiles("test")
class TaskServiceIntegrationTest {

    @Autowired
    private TaskService taskService;

    @Autowired
    private TaskRepository taskRepository;

    @Autowired
    private OrganizationRepository organizationRepository;
    @Autowired
    private UserRepository userRepository;
    @Autowired
    private ProjectRepository projectRepository;

    @MockBean
    private ActivityLogRepository activityLogRepository; // Mock Mongo

    @MockBean
    private StringRedisTemplate redisTemplate; // Mock Redis

    @MockBean
    private SimpMessagingTemplate messagingTemplate; // Mock Websocket

    private String projectId;

    @BeforeEach
    void setUp() {
        taskRepository.deleteAll();
        projectRepository.deleteAll();
        userRepository.deleteAll();
        organizationRepository.deleteAll();

        // Seed basic data for foreign keys
        Organization org = organizationRepository.save(Organization.builder().name("Test Org").build());
        User user = userRepository.save(User.builder().username("dev_user").email("test@example.com").organization(org).build());
        Project project = projectRepository.save(Project.builder().name("Test Project").keyName("TEST").organization(org).createdBy(user).build());
        projectId = project.getId();
    }

    @Test
    void createTask_ShouldPersistAndPublish() {
        CreateTaskRequest request = new CreateTaskRequest();
        request.setTitle("Integration Test Task");
        request.setDescription("Testing H2 persistence");
        request.setPriority("HIGH");
        request.setProjectId(projectId);

        TaskResponse response = taskService.createTask(request);

        assertNotNull(response.getId());
        assertEquals("Integration Test Task", response.getTitle());

        // Verify DB
        Task persisted = taskRepository.findById(response.getId()).orElse(null);
        assertNotNull(persisted);
        assertEquals("HIGH", persisted.getPriority());

        // Verify WebSocket broadcast
        verify(messagingTemplate).convertAndSend(any(String.class), any(TaskResponse.class));

        // Verify Activity Log attempt
        verify(activityLogRepository).save(any(ActivityLog.class));
    }
}
