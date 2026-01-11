import { createRouter, createWebHistory } from 'vue-router';
import AgentList from './views/AgentList.vue';
import AgentCreate from './views/AgentCreate.vue';

const routes = [
  { path: '/', component: AgentList },
  { path: '/agents/new', component: AgentCreate },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
