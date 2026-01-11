package com.prism.service.integration;

import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import java.util.UUID;

@Slf4j
@Service
public class MockNotionClient implements NotionClient {

    @Override
    public String createPage(String title, String markdownContent) {
        log.info("Mock Notion API: Creating page '{}'", title);
        log.debug("Content:\n{}", markdownContent);

        // Simulate API latency
        try {
            Thread.sleep(500);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }

        return "notion-page-" + UUID.randomUUID().toString();
    }
}
