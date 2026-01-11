package com.prism.dto;

import lombok.Data;
import java.math.BigDecimal;

@Data
public class CreateAgentRequest {
    private String name;
    private String role;
    private String description;
    private String providerId;
    private String modelName;
    private String systemPrompt;
    private BigDecimal temperature;
    private Integer maxTokens;
    private Boolean canWriteCode;
    private Boolean canRunTests;
    private Boolean canDeploy;
    private Boolean canCreateDocuments;
    private Boolean canMergePr;
}
