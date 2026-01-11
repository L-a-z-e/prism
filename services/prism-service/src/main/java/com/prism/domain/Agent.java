package com.prism.domain;

import jakarta.persistence.Entity;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import jakarta.persistence.ManyToOne;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.Column;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import java.math.BigDecimal;
import java.time.LocalDateTime;

@Entity
@Table(name = "agents")
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class Agent {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private String id;

    private String name;
    private String role; // PM, BACKEND, FRONTEND, etc.
    private String description;

    // AI Config
    @ManyToOne
    @JoinColumn(name = "provider_id")
    private AiProvider provider;

    private String modelName;

    @Column(columnDefinition = "LONGTEXT")
    private String systemPrompt;

    @Builder.Default
    private BigDecimal temperature = new BigDecimal("0.7");

    @Builder.Default
    private Integer maxTokens = 4096;

    // Capabilities
    @Builder.Default
    private boolean canWriteCode = true;
    @Builder.Default
    private boolean canRunTests = true;
    @Builder.Default
    private boolean canDeploy = false;
    @Builder.Default
    private boolean canCreateDocuments = true;
    @Builder.Default
    private boolean canMergePr = false;

    // Ownership
    @ManyToOne
    @JoinColumn(name = "organization_id")
    private Organization organization;

    @ManyToOne
    @JoinColumn(name = "created_by")
    private User createdBy;

    @Builder.Default
    private LocalDateTime createdAt = LocalDateTime.now();
    @Builder.Default
    private LocalDateTime updatedAt = LocalDateTime.now();
}
