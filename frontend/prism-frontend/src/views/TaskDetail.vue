<template>
  <div class="p-6">
    <div v-if="loading" class="text-center">Loading...</div>

    <div v-else-if="taskDetail" class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Left: Task Info -->
      <div class="lg:col-span-2 space-y-6">
        <div class="bg-white p-6 rounded shadow-sm border">
          <div class="flex justify-between items-start mb-4">
            <h1 class="text-2xl font-bold">{{ taskDetail.task.title }}</h1>
            <span :class="statusClass(taskDetail.task.status)" class="px-3 py-1 rounded text-sm font-bold uppercase">
              {{ taskDetail.task.status }}
            </span>
          </div>
          <p class="text-gray-700 mb-4">{{ taskDetail.task.description }}</p>

          <div class="grid grid-cols-2 gap-4 text-sm text-gray-500 border-t pt-4">
            <div>
              <span class="block text-xs font-semibold uppercase text-gray-400">Project</span>
              {{ taskDetail.task.projectName }}
            </div>
            <div>
              <span class="block text-xs font-semibold uppercase text-gray-400">Assigned To</span>
              {{ taskDetail.task.assignedToName || 'Unassigned' }}
            </div>
            <div>
              <span class="block text-xs font-semibold uppercase text-gray-400">Created At</span>
              {{ new Date(taskDetail.task.createdAt).toLocaleString() }}
            </div>
          </div>

          <!-- Git Metadata -->
          <div v-if="taskDetail.task.gitBranch || taskDetail.task.gitPrUrl" class="mt-6 border-t pt-4">
             <h3 class="text-sm font-bold mb-2">Git Integration</h3>
             <div class="grid grid-cols-2 gap-4 text-sm">
                <div v-if="taskDetail.task.gitBranch">
                    <span class="block text-xs font-semibold uppercase text-gray-400">Branch</span>
                    <span class="font-mono bg-gray-100 px-1 rounded">{{ taskDetail.task.gitBranch }}</span>
                </div>
                <div v-if="taskDetail.task.gitCommitHash">
                    <span class="block text-xs font-semibold uppercase text-gray-400">Commit</span>
                    <span class="font-mono bg-gray-100 px-1 rounded">{{ taskDetail.task.gitCommitHash.substring(0, 7) }}</span>
                </div>
                <div v-if="taskDetail.task.gitPrUrl" class="col-span-2">
                    <span class="block text-xs font-semibold uppercase text-gray-400">Pull Request</span>
                    <a :href="taskDetail.task.gitPrUrl" target="_blank" class="text-blue-600 hover:underline">{{ taskDetail.task.gitPrUrl }}</a>
                </div>
             </div>
          </div>
        </div>

        <!-- Build/Deployment Logs (Mock for now, can be real streaming logs) -->
        <div v-if="logs.length > 0" class="bg-gray-900 text-gray-100 p-4 rounded shadow-sm font-mono text-sm h-64 overflow-y-auto">
            <div v-for="(log, i) in logs" :key="i" class="mb-1">
                <span class="text-gray-500">[{{ log.timestamp }}]</span> {{ log.message }}
            </div>
        </div>
      </div>

      <!-- Right: Timeline -->
      <div class="space-y-4">
        <h2 class="text-lg font-semibold">Activity Timeline</h2>
        <div class="relative pl-6 border-l-2 border-gray-200 space-y-8">
            <div v-for="(activity, index) in taskDetail.timeline" :key="index" class="relative">
                <div class="absolute -left-2.5 mt-1.5 h-5 w-5 rounded-full border-2 border-white" :class="activityColor(activity.action)"></div>
                <div>
                    <h3 class="text-sm font-semibold">{{ formatAction(activity.action) }}</h3>
                    <span class="text-xs text-gray-500">{{ new Date(activity.timestamp).toLocaleString() }}</span>
                    <div v-if="activity.details" class="mt-1 text-xs text-gray-600 bg-gray-50 p-2 rounded">
                        <pre class="whitespace-pre-wrap">{{ JSON.stringify(activity.details, null, 2) }}</pre>
                    </div>
                </div>
            </div>
            <div v-if="taskDetail.timeline.length === 0" class="text-gray-500 text-sm">No activity yet.</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { getTask, type TaskDetail } from '../api';
import { useWebSocket } from '../useWebSocket';

const route = useRoute();
const taskId = route.params.id as string;

const taskDetail = ref<TaskDetail | null>(null);
const loading = ref(true);
const logs = ref<any[]>([]);

const { connect, subscribe } = useWebSocket(() => {
    subscribe(`/topic/tasks/${taskId}`, (message) => {
        console.log('Received update:', message);

        // Update Status
        if (taskDetail.value && message.type === 'STATUS_UPDATE') {
            taskDetail.value.task.status = message.status;

            // Update Git info if present
            if (message.gitBranch) taskDetail.value.task.gitBranch = message.gitBranch;
            if (message.gitCommitHash) taskDetail.value.task.gitCommitHash = message.gitCommitHash;
            if (message.gitPrUrl) taskDetail.value.task.gitPrUrl = message.gitPrUrl;

            // Add to timeline
            taskDetail.value.timeline.unshift({
                action: 'TASK_STATUS_UPDATE',
                timestamp: message.timestamp,
                details: { status: message.status, details: message.details }
            });

            // Add to streaming logs (mock visualization)
            if (message.details) {
                logs.value.push({ timestamp: new Date().toLocaleTimeString(), message: message.details });
            }
        }
    });
});

onMounted(async () => {
    try {
        connect();
        taskDetail.value = await getTask(taskId);
    } catch (e) {
        console.error(e);
    } finally {
        loading.value = false;
    }
});

const statusClass = (status: string) => {
  switch (status) {
    case 'DONE': return 'bg-green-100 text-green-800';
    case 'IN_PROGRESS': return 'bg-blue-100 text-blue-800';
    case 'FAILED': return 'bg-red-100 text-red-800';
    default: return 'bg-gray-100 text-gray-800';
  }
};

const activityColor = (action: string) => {
    switch(action) {
        case 'TASK_CREATED': return 'bg-blue-500';
        case 'TASK_STATUS_UPDATE': return 'bg-purple-500';
        default: return 'bg-gray-400';
    }
};

const formatAction = (action: string) => {
    return action.replace(/_/g, ' ');
};
</script>
