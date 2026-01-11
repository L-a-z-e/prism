package com.prism.service;

import com.prism.repository.TaskRepository;
import com.prism.repository.ActivityLogRepository;
import com.prism.domain.Task;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import java.util.List;
import java.util.Map;

import static org.mockito.Mockito.when;
import static org.junit.jupiter.api.Assertions.*;

@ExtendWith(MockitoExtension.class)
public class DashboardServiceTest {

    @Mock
    private TaskRepository taskRepository;

    @Mock
    private ActivityLogRepository activityLogRepository;

    @InjectMocks
    private DashboardService dashboardService;

    @Test
    public void testGetStats() {
        when(taskRepository.count()).thenReturn(10L);
        when(taskRepository.countTasksByStatus()).thenReturn(List.of(
            new Object[]{"TODO", 1L},
            new Object[]{"DONE", 1L}
        ));

        Map<String, Object> stats = dashboardService.getStats();

        assertEquals(10L, stats.get("totalTasks"));
        Map<String, Long> byStatus = (Map<String, Long>) stats.get("tasksByStatus");
        assertEquals(1L, byStatus.get("TODO"));
        assertEquals(1L, byStatus.get("DONE"));
    }

    @Test
    public void testGetChartData() {
        Map<String, Object> data = dashboardService.getChartData("activity");
        assertNotNull(data.get("labels"));
        assertNotNull(data.get("datasets"));
    }
}
