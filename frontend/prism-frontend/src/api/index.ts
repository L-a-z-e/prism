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
