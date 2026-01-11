package com.prism.dto;

import com.prism.domain.Agent;
import lombok.Builder;
import lombok.Data;
import java.math.BigDecimal;
import java.time.LocalDateTime;

@Data
@Builder
public class AgentResponse {
    private String id;
    private String name;
    private String role;
    private String description;
    private String providerName;
    private String modelName;
    private BigDecimal temperature;
    private Integer maxTokens;
    private LocalDateTime createdAt;

    public static AgentResponse from(Agent agent) {
        return AgentResponse.builder()
            .id(agent.getId())
            .name(agent.getName())
            .role(agent.getRole())
            .description(agent.getDescription())
            .providerName(agent.getProvider() != null ? agent.getProvider().getName() : null)
            .modelName(agent.getModelName())
            .temperature(agent.getTemperature())
            .maxTokens(agent.getMaxTokens())
            .createdAt(agent.getCreatedAt())
            .build();
    }
}
