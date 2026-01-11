package com.prism.service.integration;

public interface NotionClient {
    String createPage(String title, String markdownContent);
}
