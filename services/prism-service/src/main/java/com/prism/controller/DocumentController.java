package com.prism.controller;

import com.prism.service.document.DocumentService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpHeaders;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
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

    @GetMapping("/{taskId}/documents/raw")
    public ResponseEntity<String> downloadMarkdown(@PathVariable String taskId) {
        String content = documentService.getRawMarkdown(taskId);
        String filename = "task-" + taskId + ".md";

        return ResponseEntity.ok()
            .header(HttpHeaders.CONTENT_DISPOSITION, "attachment; filename=\"" + filename + "\"")
            .contentType(MediaType.TEXT_MARKDOWN)
            .body(content);
    }
}
