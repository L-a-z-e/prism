package com.prism.service;

import com.prism.domain.Organization;
import com.prism.domain.Project;
import com.prism.domain.User;
import com.prism.domain.AiProvider;
import com.prism.repository.OrganizationRepository;
import com.prism.repository.ProjectRepository;
import com.prism.repository.UserRepository;
import com.prism.repository.AiProviderRepository;
import jakarta.annotation.PostConstruct;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
@RequiredArgsConstructor
public class MockUserService {
    private final UserRepository userRepository;
    private final OrganizationRepository organizationRepository;
    private final ProjectRepository projectRepository;
    private final AiProviderRepository aiProviderRepository;

    private static final String MOCK_USERNAME = "dev_user";

    @PostConstruct
    @Transactional
    public void initData() {
        if (organizationRepository.count() == 0) {
            Organization org = organizationRepository.save(Organization.builder()
                .name("Default Organization")
                .description("Default organization for development")
                .build());

            User user = userRepository.save(User.builder()
                .username(MOCK_USERNAME)
                .email("dev@example.com")
                .role("ADMIN")
                .organization(org)
                .build());

            projectRepository.save(Project.builder()
                .name("Prism Project")
                .keyName("PRISM")
                .description("Main development project")
                .organization(org)
                .createdBy(user)
                .build());

            aiProviderRepository.save(AiProvider.builder()
                .name("Anthropic")
                .baseUrl("https://api.anthropic.com")
                .apiKeyEnvVar("ANTHROPIC_API_KEY")
                .build());
        }
    }

    public User getCurrentUser() {
        return userRepository.findByUsername(MOCK_USERNAME)
            .orElseThrow(() -> new RuntimeException("Mock user not found. Data init failed?"));
    }
}
