package com.prism.service;

import com.prism.repository.TaskRepository;
import com.prism.repository.ActivityLogRepository;
import com.prism.domain.Task;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import java.util.Map;
import java.util.HashMap;
import java.util.List;
import java.util.stream.Collectors;

@Service
@RequiredArgsConstructor
public class DashboardService {
    private final TaskRepository taskRepository;
    private final ActivityLogRepository activityLogRepository;

    public Map<String, Object> getStats() {
        Map<String, Object> stats = new HashMap<>();
        long totalTasks = taskRepository.count();

        List<Object[]> statusCounts = taskRepository.countTasksByStatus();
        Map<String, Long> tasksByStatus = statusCounts.stream()
            .collect(Collectors.toMap(
                row -> (String) row[0],
                row -> (Long) row[1]
            ));

        stats.put("totalTasks", totalTasks);
        stats.put("tasksByStatus", tasksByStatus);

        // Mock data for metrics not yet fully tracked
        stats.put("activeAgents", 3);
        stats.put("totalCost", 15.42);
        stats.put("avgCompletionTimeHours", 4.5);

        return stats;
    }

    public Map<String, Object> getChartData(String type) {
        if ("activity".equals(type)) {
             // In a real app, aggregation over ActivityLog by date
             return Map.of(
                 "labels", List.of("Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"),
                 "datasets", List.of(Map.of(
                     "label", "Activities",
                     "data", List.of(12, 19, 3, 5, 2, 8, 15),
                     "backgroundColor", "#3B82F6"
                 ))
             );
        } else if ("cost".equals(type)) {
             return Map.of(
                 "labels", List.of("Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"),
                 "datasets", List.of(Map.of(
                     "label", "Cost ($)",
                     "data", List.of(1.2, 2.5, 0.5, 1.0, 0.8, 3.2, 1.5),
                     "backgroundColor", "#EF4444"
                 ))
             );
        }
        return Map.of();
    }
}
