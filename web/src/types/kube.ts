export interface ResourceStatus {
  phase: string       // Running, Pending, Failed, Succeeded, Unknown
  ready: boolean
  reason?: string
  message?: string
  lastUpdated: string // ISO string
}

export interface ResourceMetrics {
  cpuUsage: number      // In millicores
  memoryUsage: number   // In bytes
  cpuRequest: number    // In millicores
  memoryRequest: number // In bytes
  cpuLimit: number      // In millicores
  memoryLimit: number   // In bytes
}

export interface NodeCapacity {
  cpu: string
  memory: string
  pods: string
  ephemeralStorage?: string
  storage?: string
}

export interface NodeCondition {
  type: string
  status: string
  reason?: string
  message?: string
  lastTransitionTime: string // ISO string
}

export interface KubeContainer {
  name: string
  image: string
  ready: boolean
  metrics?: ResourceMetrics
  status: ResourceStatus
}

export interface KubeVolume {
  name: string
  type: string
  size?: string
}

export interface KubePod {
  name: string
  namespace: string
  containers: KubeContainer[]
  volumes: KubeVolume[]
  role: string
  ip: string
  labels: Record<string, string>
  status: ResourceStatus
  metrics?: ResourceMetrics
  nodeName: string
  createdAt: string // ISO string
  restartCount: number
}

export interface KubeNode {
  name: string
  kubelet: string
  runtime: string
  role: string
  pods: KubePod[]
  status: ResourceStatus
  metrics?: ResourceMetrics
  capacity: NodeCapacity
  allocatable: NodeCapacity
  conditions: NodeCondition[]
  labels: Record<string, string>
  annotations?: Record<string, string>
  createdAt: string // ISO string
}

export interface KubeDeployment {
  name: string
  namespace: string
  replicas: number
  readyReplicas: number
  availableReplicas: number
  selector: string
  template: string
  strategy: string
  role: string
  status: ResourceStatus
  labels: Record<string, string>
  createdAt: string // ISO string
}

export interface ServicePort {
  name?: string
  port: number
  targetPort: string
  protocol: string
  nodePort?: number
}

export interface KubeService {
  name: string
  namespace: string
  type: string
  pods: KubePod[]
  clusterIP: string
  nodePort?: number
  loadBalancerIP?: string
  externalIPs?: string[]
  selector: string
  ports: ServicePort[]
  role: string
  status: ResourceStatus
  labels: Record<string, string>
  createdAt: string // ISO string
}

export interface KubeControlPlane {
  nodes: KubeNode[]
  status: ResourceStatus
}

export interface KubeNamespace {
  name: string
  deployments: KubeDeployment[]
  services: KubeService[]
  podCount: number
  status: ResourceStatus
  labels: Record<string, string>
  createdAt: string // ISO string
}

export interface ClusterSummary {
  totalNodes: number
  readyNodes: number
  totalPods: number
  runningPods: number
  pendingPods: number
  failedPods: number
  totalNamespaces: number
  totalDeployments: number
  totalServices: number
  cpuUtilization: number    // percentage
  memoryUtilization: number // percentage
}

export interface KubeCluster {
  id: string
  name: string
  version: string
  environment?: string
  controlPlane: KubeControlPlane
  nodes: KubeNode[]
  namespaces: KubeNamespace[]
  status: ResourceStatus
  summary: ClusterSummary
  createdAt?: string // ISO string
}
