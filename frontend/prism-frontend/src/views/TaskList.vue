<template>
  <div class="p-6">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">Tasks</h1>
      <router-link to="/tasks/new" class="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600">
        + New Task
      </router-link>
    </div>

    <div v-if="loading" class="text-center">Loading...</div>

    <div v-else class="space-y-4">
      <div v-for="task in tasks" :key="task.id" class="border p-4 rounded shadow-sm bg-white hover:shadow-md transition">
        <div class="flex justify-between">
          <h2 class="text-xl font-semibold">{{ task.title }}</h2>
          <span :class="statusClass(task.status)" class="text-xs px-2 py-1 rounded font-bold uppercase">{{ task.status }}</span>
        </div>
        <p class="text-gray-600 mt-1">{{ task.description }}</p>
        <div class="mt-3 text-sm text-gray-500 flex gap-4">
          <span>Assigned to: <strong>{{ task.assignedToName || 'Unassigned' }}</strong></span>
          <span>Project: {{ task.projectName }}</span>
          <span>Priority: {{ task.priority }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { getTasks, type Task } from '../api';

const tasks = ref<Task[]>([]);
const loading = ref(true);

onMounted(async () => {
  try {
    tasks.value = await getTasks();
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
</script>
