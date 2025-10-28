import type { KubeCluster } from '@/types/kube'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

export class ApiError extends Error {
  constructor(
    message: string,
    public status?: number,
    public statusText?: string
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

async function fetchApi<T>(endpoint: string): Promise<T> {
  try {
    const response = await fetch(`${API_BASE_URL}${endpoint}`)

    if (!response.ok) {
      throw new ApiError(
        `API request failed: ${response.statusText}`,
        response.status,
        response.statusText
      )
    }

    const data = await response.json()
    return data as T
  } catch (error) {
    if (error instanceof ApiError) {
      throw error
    }
    throw new ApiError(
      `Network error: ${error instanceof Error ? error.message : 'Unknown error'}`
    )
  }
}

export const api = {
  // Cluster endpoints
  getClusters: () => fetchApi<KubeCluster[]>('/api/clusters'),

  getCluster: (id: string) => fetchApi<KubeCluster>(`/api/clusters/${id}`),

  getClusterNodes: (id: string) => fetchApi(`/api/clusters/${id}/nodes`),

  getClusterPods: (id: string) => fetchApi(`/api/clusters/${id}/pods`),

  getClusterServices: (id: string) => fetchApi(`/api/clusters/${id}/services`),

  getClusterDeployments: (id: string) => fetchApi(`/api/clusters/${id}/deployments`),

  getClusterNamespaces: (id: string) => fetchApi(`/api/clusters/${id}/namespaces`),

  // Health check
  healthCheck: () => fetchApi<{ status: string; time: number }>('/health'),

  // WebSocket status
  getWebSocketStatus: () => fetchApi<{ connected_clients: number; status: string }>('/ws/status'),
}
