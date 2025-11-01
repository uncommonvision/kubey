import type { KubeCluster } from '@/types/kube'

// Use empty string since Vite proxy will handle the API calls
const API_BASE_URL = import.meta.env.VITE_API_URL || ''

// Functional ApiError factory to replace class approach
export type ApiError = Error & { status?: number; statusText?: string }

export const createApiError = (
  message: string,
  status?: number,
  statusText?: string
): ApiError => {
  const err = new Error(message) as ApiError
  err.name = 'ApiError'
  err.status = status
  err.statusText = statusText
  return err
}

// Helper for instanceof checks - maintain compatibility with existing code
export const isApiError = (error: unknown): error is ApiError => {
  return error instanceof Error && error.name === 'ApiError'
}

async function fetchApi<T>(endpoint: string): Promise<T> {
  try {
    const response = await fetch(`${API_BASE_URL}${endpoint}`)

    if (!response.ok) {
      throw createApiError(
        `API request failed: ${response.statusText}`,
        response.status,
        response.statusText
      )
    }

    const data = await response.json()
    return data as T
  } catch (error) {
    if (isApiError(error)) {
      throw error
    }
    throw createApiError(
      `Network error: ${error instanceof Error ? error.message : 'Unknown error'}`
    )
  }
}

// Cluster Service - Individual function exports for cluster-related operations
export const getClusters = () => fetchApi<KubeCluster[]>('/api/clusters')

export const getCluster = (id: string) => fetchApi<KubeCluster>(`/api/clusters/${id}`)

export const getClusterNodes = (id: string) => fetchApi(`/api/clusters/${id}/nodes`)

export const getClusterPods = (id: string) => fetchApi(`/api/clusters/${id}/pods`)

export const getClusterServices = (id: string) => fetchApi(`/api/clusters/${id}/services`)

export const getClusterDeployments = (id: string) => fetchApi(`/api/clusters/${id}/deployments`)

export const getClusterNamespaces = (id: string) => fetchApi(`/api/clusters/${id}/namespaces`)

// Export ApiError constructor for convenience
export const ApiError = createApiError
