<template>
  <div class="p-6 max-w-7xl mx-auto">
    <h1 class="text-3xl font-bold mb-8 text-gray-800">Dashboard</h1>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      <div v-for="(value, key) in summaryStats" :key="key" class="bg-white p-6 rounded-lg shadow-md border border-gray-100">
        <h3 class="text-sm font-medium text-gray-500 uppercase tracking-wider">{{ key }}</h3>
        <p class="mt-2 text-3xl font-bold text-gray-900">{{ value }}</p>
      </div>
    </div>

    <!-- Charts -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
      <div class="bg-white p-6 rounded-lg shadow-md border border-gray-100">
        <h3 class="text-lg font-bold mb-4">Activity Volume (Last 7 Days)</h3>
        <div class="h-64">
           <Bar v-if="activityData" :data="activityData" :options="chartOptions" />
           <p v-else class="text-gray-400 text-center py-10">Loading Chart...</p>
        </div>
      </div>

      <div class="bg-white p-6 rounded-lg shadow-md border border-gray-100">
        <h3 class="text-lg font-bold mb-4">Cost Trend (Last 7 Days)</h3>
        <div class="h-64">
           <Bar v-if="costData" :data="costData" :options="chartOptions" />
           <p v-else class="text-gray-400 text-center py-10">Loading Chart...</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { getDashboardStats, getChartData } from '../api';
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  BarElement,
  CategoryScale,
  LinearScale
} from 'chart.js'
import { Bar } from 'vue-chartjs'

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend)

const summaryStats = ref<Record<string, any>>({});
const activityData = ref(null);
const costData = ref(null);

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false
};

onMounted(async () => {
  try {
    const stats = await getDashboardStats();
    summaryStats.value = {
        'Total Tasks': stats.totalTasks,
        'Active Agents': stats.activeAgents,
        'Total Cost ($)': stats.totalCost,
        'Avg Completion (Hrs)': stats.avgCompletionTimeHours
    };

    activityData.value = await getChartData('activity');
    costData.value = await getChartData('cost');
  } catch (e) {
    console.error(e);
  }
});
</script>
