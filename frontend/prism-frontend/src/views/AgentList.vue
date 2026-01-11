<template>
  <div class="p-6">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">Agents</h1>
      <router-link to="/agents/new" class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600">
        + New Agent
      </router-link>
    </div>

    <div v-if="loading" class="text-center">Loading...</div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <div v-for="agent in agents" :key="agent.id" class="border p-4 rounded shadow-sm hover:shadow-md transition">
        <h2 class="text-xl font-semibold">{{ agent.name }}</h2>
        <span class="inline-block bg-gray-100 text-gray-800 text-xs px-2 py-1 rounded mt-1">{{ agent.role }}</span>
        <p class="text-gray-600 mt-2">{{ agent.description }}</p>
        <div class="mt-4 text-sm text-gray-500">
          <p>Model: {{ agent.modelName || 'N/A' }}</p>
          <p>Provider: {{ agent.providerName || 'N/A' }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { getAgents, type Agent } from '../api';

const agents = ref<Agent[]>([]);
const loading = ref(true);

onMounted(async () => {
  try {
    agents.value = await getAgents();
  } catch (e) {
    console.error(e);
  } finally {
    loading.value = false;
  }
});
</script>
