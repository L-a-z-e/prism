package com.prism.controller;

import com.prism.service.document.DocumentService;
import lombok.RequiredArgsConstructor;
import org.springframework.web.bind.annotation.*;

import java.util.Map;

@RestController
@RequestMapping("/api/v1/tasks")
@RequiredArgsConstructor
@CrossOrigin(origins = "*")
public class DocumentController {
    private final DocumentService documentService;

    @PostMapping("/{taskId}/documents/notion")
    public Map<String, String> exportToNotion(@PathVariable String taskId) {
        String pageId = documentService.publishTaskToNotion(taskId);
        return Map.of("pageId", pageId, "status", "PUBLISHED");
    }
}
