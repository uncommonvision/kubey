package models

import "time"

// ResourceStatus represents the status of any Kubernetes resource
type ResourceStatus struct {
	Phase       string    `json:"phase"` // Running, Pending, Failed, Succeeded, Unknown
	Ready       bool      `json:"ready"`
	Reason      string    `json:"reason,omitempty"`
	Message     string    `json:"message,omitempty"`
	LastUpdated time.Time `json:"lastUpdated"`
}

// ResourceMetrics represents resource usage metrics
type ResourceMetrics struct {
	CPUUsage      float64 `json:"cpuUsage"`      // In millicores
	MemoryUsage   int64   `json:"memoryUsage"`   // In bytes
	CPURequest    float64 `json:"cpuRequest"`    // In millicores
	MemoryRequest int64   `json:"memoryRequest"` // In bytes
	CPULimit      float64 `json:"cpuLimit"`      // In millicores
	MemoryLimit   int64   `json:"memoryLimit"`   // In bytes
}

// NodeCapacity represents node resource capacity
type NodeCapacity struct {
	CPU              string `json:"cpu"`
	Memory           string `json:"memory"`
	Pods             string `json:"pods"`
	EphemeralStorage string `json:"ephemeralStorage,omitempty"`
	Storage          string `json:"storage,omitempty"`
}

// NodeCondition represents a node condition
type NodeCondition struct {
	Type               string    `json:"type"`
	Status             string    `json:"status"`
	Reason             string    `json:"reason,omitempty"`
	Message            string    `json:"message,omitempty"`
	LastTransitionTime time.Time `json:"lastTransitionTime"`
}

// KubeContainer represents a container within a pod
type KubeContainer struct {
	Name    string           `json:"name"`
	Image   string           `json:"image"`
	Ready   bool             `json:"ready"`
	Metrics *ResourceMetrics `json:"metrics,omitempty"`
	Status  ResourceStatus   `json:"status"`
}

// KubeVolume represents a volume in a pod
type KubeVolume struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Size string `json:"size,omitempty"`
}

// KubePod represents a Kubernetes pod
type KubePod struct {
	Name         string            `json:"name"`
	Namespace    string            `json:"namespace"`
	Containers   []KubeContainer   `json:"containers"`
	Volumes      []KubeVolume      `json:"volumes"`
	Role         string            `json:"role"`
	IP           string            `json:"ip"`
	Labels       map[string]string `json:"labels"`
	Status       ResourceStatus    `json:"status"`
	Metrics      *ResourceMetrics  `json:"metrics,omitempty"`
	NodeName     string            `json:"nodeName"`
	CreatedAt    time.Time         `json:"createdAt"`
	RestartCount int32             `json:"restartCount"`
}

// KubeNode represents a Kubernetes node
type KubeNode struct {
	Name        string            `json:"name"`
	Kubelet     string            `json:"kubelet"`
	Runtime     string            `json:"runtime"`
	Role        string            `json:"role"`
	Pods        []KubePod         `json:"pods"`
	Status      ResourceStatus    `json:"status"`
	Metrics     *ResourceMetrics  `json:"metrics,omitempty"`
	Capacity    NodeCapacity      `json:"capacity"`
	Allocatable NodeCapacity      `json:"allocatable"`
	Conditions  []NodeCondition   `json:"conditions"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations,omitempty"`
	CreatedAt   time.Time         `json:"createdAt"`
}

// KubeDeployment represents a Kubernetes deployment
type KubeDeployment struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	Replicas          int32             `json:"replicas"`
	ReadyReplicas     int32             `json:"readyReplicas"`
	AvailableReplicas int32             `json:"availableReplicas"`
	Selector          string            `json:"selector"`
	Template          string            `json:"template"`
	Strategy          string            `json:"strategy"`
	Role              string            `json:"role"`
	Status            ResourceStatus    `json:"status"`
	Labels            map[string]string `json:"labels"`
	CreatedAt         time.Time         `json:"createdAt"`
}

// KubeService represents a Kubernetes service
type KubeService struct {
	Name           string            `json:"name"`
	Namespace      string            `json:"namespace"`
	Type           string            `json:"type"`
	Pods           []KubePod         `json:"pods"`
	ClusterIP      string            `json:"clusterIP"`
	NodePort       *int32            `json:"nodePort,omitempty"`
	LoadBalancerIP string            `json:"loadBalancerIP,omitempty"`
	ExternalIPs    []string          `json:"externalIPs,omitempty"`
	Selector       string            `json:"selector"`
	Ports          []ServicePort     `json:"ports"`
	Role           string            `json:"role"`
	Status         ResourceStatus    `json:"status"`
	Labels         map[string]string `json:"labels"`
	CreatedAt      time.Time         `json:"createdAt"`
}

// ServicePort represents a service port
type ServicePort struct {
	Name       string `json:"name,omitempty"`
	Port       int32  `json:"port"`
	TargetPort string `json:"targetPort"`
	Protocol   string `json:"protocol"`
	NodePort   int32  `json:"nodePort,omitempty"`
}

// KubeControlPlane represents the control plane
type KubeControlPlane struct {
	Nodes  []KubeNode     `json:"nodes"`
	Status ResourceStatus `json:"status"`
}

// KubeNamespace represents a Kubernetes namespace
type KubeNamespace struct {
	Name        string            `json:"name"`
	Deployments []KubeDeployment  `json:"deployments"`
	Services    []KubeService     `json:"services"`
	PodCount    int               `json:"podCount"`
	Status      ResourceStatus    `json:"status"`
	Labels      map[string]string `json:"labels"`
	CreatedAt   time.Time         `json:"createdAt"`
}

// KubeCluster represents a complete Kubernetes cluster
type KubeCluster struct {
	ID           string           `json:"id"`
	Name         string           `json:"name"`
	Version      string           `json:"version"`
	Environment  string           `json:"environment,omitempty"`
	ControlPlane KubeControlPlane `json:"controlPlane"`
	Nodes        []KubeNode       `json:"nodes"`
	Namespaces   []KubeNamespace  `json:"namespaces"`
	Status       ResourceStatus   `json:"status"`
	Summary      ClusterSummary   `json:"summary"`
	CreatedAt    time.Time        `json:"createdAt,omitempty"`
}

// ClusterSummary provides high-level cluster statistics
type ClusterSummary struct {
	TotalNodes        int     `json:"totalNodes"`
	ReadyNodes        int     `json:"readyNodes"`
	TotalPods         int     `json:"totalPods"`
	RunningPods       int     `json:"runningPods"`
	PendingPods       int     `json:"pendingPods"`
	FailedPods        int     `json:"failedPods"`
	TotalNamespaces   int     `json:"totalNamespaces"`
	TotalDeployments  int     `json:"totalDeployments"`
	TotalServices     int     `json:"totalServices"`
	CPUUtilization    float64 `json:"cpuUtilization"`    // percentage
	MemoryUtilization float64 `json:"memoryUtilization"` // percentage
}
