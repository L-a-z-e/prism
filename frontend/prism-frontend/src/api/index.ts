import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8085/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
});

export interface Agent {
  id: string;
  name: string;
  role: string;
  description: string;
  providerName: string;
  modelName: string;
  createdAt: string;
}

export const getAgents = async (): Promise<Agent[]> => {
  const response = await api.get('/agents');
  return response.data;
};

export const createAgent = async (agent: any): Promise<Agent> => {
  const response = await api.post('/agents', agent);
  return response.data;
};

export interface Task {
  id: string;
  title: string;
  description: string;
  status: string;
  priority: string;
  assignedToName: string;
  projectName: string;
  createdAt: string;
  gitBranch?: string;
  gitCommitHash?: string;
  gitPrUrl?: string;
}

export const getTasks = async (filters?: any): Promise<Task[]> => {
  const response = await api.get('/tasks', { params: filters });
  return response.data;
};

export const createTask = async (task: any): Promise<Task> => {
  const response = await api.post('/tasks', task);
  return response.data;
};

export interface ActivityLog {
  action: string;
  timestamp: string;
  details: any;
}

export interface TaskDetail {
  task: Task;
  timeline: ActivityLog[];
}

export const getTask = async (taskId: string): Promise<TaskDetail> => {
  const response = await api.get(`/tasks/${taskId}`);
  return response.data;
};

export const exportToNotion = async (taskId: string): Promise<{ pageId: string }> => {
  const response = await api.post(`/tasks/${taskId}/documents/notion`);
  return response.data;
};

export const downloadMarkdown = async (taskId: string) => {
  const response = await api.get(`/tasks/${taskId}/documents/raw`, { responseType: 'blob' });
  const url = window.URL.createObjectURL(new Blob([response.data]));
  const link = document.createElement('a');
  link.href = url;
  link.setAttribute('download', `task-${taskId}.md`);
  document.body.appendChild(link);
  link.click();
  link.remove();
};
