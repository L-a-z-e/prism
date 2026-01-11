package com.prism.controller;

import com.prism.service.DashboardService;
import lombok.RequiredArgsConstructor;
import org.springframework.web.bind.annotation.*;
import java.util.Map;

@RestController
@RequestMapping("/api/v1/dashboard")
@RequiredArgsConstructor
@CrossOrigin(origins = "*")
public class DashboardController {
    private final DashboardService dashboardService;

    @GetMapping("/stats")
    public Map<String, Object> getStats() {
        return dashboardService.getStats();
    }

    @GetMapping("/charts/{type}")
    public Map<String, Object> getChartData(@PathVariable String type) {
        return dashboardService.getChartData(type);
    }
}
