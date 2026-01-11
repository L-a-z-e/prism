package com.prism.controller;

import com.prism.service.DashboardService;
import lombok.RequiredArgsConstructor;
import org.springframework.web.bind.annotation.*;
import io.swagger.v3.oas.annotations.tags.Tag;
import io.swagger.v3.oas.annotations.Operation;
import java.util.Map;

@RestController
@RequestMapping("/api/v1/dashboard")
@RequiredArgsConstructor
@CrossOrigin(origins = "*")
@Tag(name = "Dashboard", description = "Dashboard & Statistics API")
public class DashboardController {
    private final DashboardService dashboardService;

    @GetMapping("/stats")
    @Operation(summary = "Get dashboard statistics")
    public Map<String, Object> getStats() {
        return dashboardService.getStats();
    }

    @GetMapping("/charts/{type}")
    @Operation(summary = "Get chart data by type (activity/cost)")
    public Map<String, Object> getChartData(@PathVariable String type) {
        return dashboardService.getChartData(type);
    }
}
