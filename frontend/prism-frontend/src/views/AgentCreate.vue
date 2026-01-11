<template>
  <div class="p-6 max-w-2xl mx-auto">
    <h1 class="text-2xl font-bold mb-6">Create New Agent</h1>

    <form @submit.prevent="submit" class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-gray-700">Name</label>
        <input v-model="form.name" type="text" required class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2" />
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700">Role</label>
        <select v-model="form.role" class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2">
          <option value="BACKEND">Backend Developer</option>
          <option value="FRONTEND">Frontend Developer</option>
          <option value="QA">QA Engineer</option>
          <option value="PM">Project Manager</option>
          <option value="DEVOPS">DevOps Engineer</option>
        </select>
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700">Description</label>
        <textarea v-model="form.description" rows="3" class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2"></textarea>
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700">Model Name</label>
        <input v-model="form.modelName" type="text" placeholder="e.g. gpt-4" class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2" />
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700">System Prompt</label>
        <textarea v-model="form.systemPrompt" rows="5" class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 font-mono text-sm"></textarea>
      </div>

      <div class="flex justify-end gap-2 mt-6">
        <router-link to="/" class="px-4 py-2 text-gray-700 border rounded hover:bg-gray-50">Cancel</router-link>
        <button type="submit" class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700">Create Agent</button>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { createAgent } from '../api';

const router = useRouter();
const form = ref({
  name: '',
  role: 'BACKEND',
  description: '',
  modelName: '',
  systemPrompt: '',
  temperature: 0.7,
  maxTokens: 4096
});

const submit = async () => {
  try {
    await createAgent(form.value);
    router.push('/');
  } catch (e) {
    console.error(e);
    alert('Failed to create agent');
  }
};
</script>
