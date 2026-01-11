package com.prism.dto;

import com.prism.domain.ActivityLog;
import com.prism.domain.Task;
import lombok.Builder;
import lombok.Data;
import java.util.List;
import java.util.stream.Collectors;

@Data
@Builder
public class TaskDetailResponse {
    private TaskResponse task;
    private List<ActivityLogDTO> timeline;

    @Data
    @Builder
    public static class ActivityLogDTO {
        private String action;
        private String timestamp;
        private Object details;

        public static ActivityLogDTO from(ActivityLog log) {
            return ActivityLogDTO.builder()
                .action(log.getAction())
                .timestamp(log.getTimestamp().toString())
                .details(log.getDetails())
                .build();
        }
    }
}
