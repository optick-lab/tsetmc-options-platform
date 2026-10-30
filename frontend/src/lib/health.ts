import { apiFetch } from './api';

export interface HealthStatus {
  status: string;
}

export async function getHealth(): Promise<HealthStatus> {
  return apiFetch<HealthStatus>('/health');
}
