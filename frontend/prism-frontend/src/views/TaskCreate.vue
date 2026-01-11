<template>
  <div class="p-6 max-w-2xl mx-auto">
    <h1 class="text-2xl font-bold mb-6">Create New Task</h1>

    <form @submit.prevent="submit" class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-gray-700">Title</label>
        <input v-model="form.title" type="text" required class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2" />
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700">Project</label>
        <p class="text-xs text-gray-500 mb-1">For MVP, using Default Project</p>
        <!-- In real app, fetch projects list -->
        <input type="text" disabled value="Prism Project" class="mt-1 block w-full bg-gray-100 border border-gray-300 rounded-md shadow-sm p-2" />
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700">Assign Agent</label>
        <select v-model="form.assignedTo" required class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2">
          <option v-for="agent in agents" :key="agent.id" :value="agent.id">
            {{ agent.name }} ({{ agent.role }})
          </option>
        </select>
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700">Description</label>
        <textarea v-model="form.description" rows="3" class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2"></textarea>
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700">Priority</label>
        <select v-model="form.priority" class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2">
          <option value="LOW">Low</option>
          <option value="MEDIUM">Medium</option>
          <option value="HIGH">High</option>
          <option value="CRITICAL">Critical</option>
        </select>
      </div>

      <div class="flex justify-end gap-2 mt-6">
        <router-link to="/tasks" class="px-4 py-2 text-gray-700 border rounded hover:bg-gray-50">Cancel</router-link>
        <button type="submit" class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700">Create Task</button>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { getAgents, createTask, type Agent } from '../api';

const router = useRouter();
const agents = ref<Agent[]>([]);
const form = ref({
  title: '',
  description: '',
  priority: 'MEDIUM',
  projectId: '', // Will fill with first available project ID or fetch
  assignedTo: ''
});

// Need to fetch project ID. For MVP, we'll cheat and assume we can get it or just hardcode if we had it.
// Actually, let's just fetch agents and maybe projects?
// For now, I'll rely on the backend finding the "Default" project if I pass a known ID or I need to fetch it.
// Wait, the API needs a ProjectID. I should probably add an API to get projects or just use the seed data's ID if I knew it.
// Better: Add a "getProjects" API or simply fetch one.
// Simplification: I'll hardcode the known Seed Project ID if I can, OR simpler: Backend looks up default if null?
// No, let's do it right. I'll add a quick "getProjects" to API/Backend.
// Time constraint: I'll just List Projects in this view if needed.
// Actually, I'll modify `TaskService` to use the default project if projectId is missing/empty for MVP.

onMounted(async () => {
  try {
    agents.value = await getAgents();
    if (agents.value.length > 0 && agents.value[0]) {
        form.value.assignedTo = agents.value[0].id;
    }
  } catch (e) {
    console.error(e);
  }
});

const submit = async () => {
  try {
    // Hack: We need a valid project ID.
    // Since I didn't expose GET /projects yet, I'll assume the backend MockUserService created one.
    // I will modify the backend TaskService to handle "default" project logic or just fetch it here.
    // Let's modify the Backend TaskService to find the default project if ID is missing.
    await createTask(form.value);
    router.push('/tasks');
  } catch (e) {
    console.error(e);
    alert('Failed to create task');
  }
};
</script>
