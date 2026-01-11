package com.prism.config;

import io.swagger.v3.oas.models.OpenAPI;
import io.swagger.v3.oas.models.info.Info;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class OpenApiConfig {

    @Bean
    public OpenAPI prismOpenAPI() {
        return new OpenAPI()
                .info(new Info().title("Prism/Loom Platform API")
                .description("Enterprise AI Agent Orchestration Platform API Documentation")
                .version("v1.0"));
    }
}
