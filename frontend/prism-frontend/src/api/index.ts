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
}

export const getTasks = async (): Promise<Task[]> => {
  const response = await api.get('/tasks');
  return response.data;
};

export const createTask = async (task: any): Promise<Task> => {
  const response = await api.post('/tasks', task);
  return response.data;
};
