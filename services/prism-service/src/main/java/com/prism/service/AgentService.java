package com.prism.service;

import com.prism.domain.Agent;
import com.prism.domain.AiProvider;
import com.prism.domain.User;
import com.prism.dto.CreateAgentRequest;
import com.prism.dto.AgentResponse;
import com.prism.repository.AgentRepository;
import com.prism.repository.AiProviderRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.stream.Collectors;

@Service
@RequiredArgsConstructor
public class AgentService {
    private final AgentRepository agentRepository;
    private final AiProviderRepository aiProviderRepository;
    private final MockUserService mockUserService;

    @Transactional
    public AgentResponse createAgent(CreateAgentRequest request) {
        User currentUser = mockUserService.getCurrentUser();

        AiProvider provider = null;
        if (request.getProviderId() != null) {
            provider = aiProviderRepository.findById(request.getProviderId())
                .orElseThrow(() -> new IllegalArgumentException("Invalid Provider ID"));
        }

        Agent agent = Agent.builder()
            .name(request.getName())
            .role(request.getRole())
            .description(request.getDescription())
            .provider(provider)
            .modelName(request.getModelName())
            .systemPrompt(request.getSystemPrompt())
            .temperature(request.getTemperature())
            .maxTokens(request.getMaxTokens())
            .canWriteCode(request.getCanWriteCode() != null ? request.getCanWriteCode() : true)
            .canRunTests(request.getCanRunTests() != null ? request.getCanRunTests() : true)
            .canDeploy(request.getCanDeploy() != null ? request.getCanDeploy() : false)
            .canCreateDocuments(request.getCanCreateDocuments() != null ? request.getCanCreateDocuments() : true)
            .canMergePr(request.getCanMergePr() != null ? request.getCanMergePr() : false)
            .organization(currentUser.getOrganization())
            .createdBy(currentUser)
            .build();

        return AgentResponse.from(agentRepository.save(agent));
    }

    @Transactional(readOnly = true)
    public List<AgentResponse> getAllAgents() {
        return agentRepository.findAll().stream()
            .map(AgentResponse::from)
            .collect(Collectors.toList());
    }
}
