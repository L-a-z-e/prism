package com.prism.controller;

import com.prism.dto.CreateAgentRequest;
import com.prism.dto.AgentResponse;
import com.prism.service.AgentService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/v1/agents")
@RequiredArgsConstructor
@CrossOrigin(origins = "*") // Allow frontend access
public class AgentController {
    private final AgentService agentService;

    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    public AgentResponse createAgent(@RequestBody CreateAgentRequest request) {
        return agentService.createAgent(request);
    }

    @GetMapping
    public List<AgentResponse> getAllAgents() {
        return agentService.getAllAgents();
    }
}
