import { createRouter, createWebHistory } from 'vue-router';
import AgentList from './views/AgentList.vue';
import AgentCreate from './views/AgentCreate.vue';
import TaskList from './views/TaskList.vue';
import TaskCreate from './views/TaskCreate.vue';
import TaskDetail from './views/TaskDetail.vue';

const routes = [
  { path: '/', component: AgentList },
  { path: '/agents/new', component: AgentCreate },
  { path: '/tasks', component: TaskList },
  { path: '/tasks/new', component: TaskCreate },
  { path: '/tasks/:id', component: TaskDetail },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
