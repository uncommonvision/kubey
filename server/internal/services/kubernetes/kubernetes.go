package kubernetes

import (
	"context"
	"fmt"
	"log"
	"time"

	"kubey/server/internal/models"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var clientset *kubernetes.Clientset
var kubeconfigPath string

// InitClient initializes the Kubernetes client
func InitClient(kubeconfig string) error {
	// Store kubeconfig path for multi-context support
	if kubeconfig == "" {
		kubeconfigPath = clientcmd.RecommendedHomeFile
	} else {
		kubeconfigPath = kubeconfig
	}

	// Try to build a config to verify kubeconfig is valid
	var config *rest.Config
	var err error

	if kubeconfig == "" {
		// Use in-cluster config if no kubeconfig specified
		config, err = rest.InClusterConfig()
		if err != nil {
			// Fall back to default kubeconfig
			config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
			if err != nil {
				return fmt.Errorf("failed to build kubeconfig: %v", err)
			}
		}
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return fmt.Errorf("failed to build kubeconfig from %s: %v", kubeconfig, err)
		}
	}

	// Create a clientset for the current context to verify connectivity
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create clientset: %v", err)
	}

	log.Println("Successfully connected to Kubernetes cluster")
	return nil
}

// GetClusters returns all clusters from all contexts in the kubeconfig (loaded in parallel)
func GetClusters() ([]models.KubeCluster, error) {
	if kubeconfigPath == "" {
		return nil, fmt.Errorf("kubernetes client not initialized")
	}

	// Load the kubeconfig file
	config, err := clientcmd.LoadFromFile(kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %v", err)
	}

	// Collect context names
	var contextNames []string
	for contextName := range config.Contexts {
		contextNames = append(contextNames, contextName)
	}

	if len(contextNames) == 0 {
		return nil, fmt.Errorf("no contexts found in kubeconfig")
	}

	// Process all contexts in parallel using goroutines
	type result struct {
		cluster *models.KubeCluster
		err     error
		name    string
	}

	results := make(chan result, len(contextNames))

	for _, contextName := range contextNames {
		go func(name string) {
			log.Printf("Processing context: %s", name)
			cluster, err := getClusterDataForContext(name)
			results <- result{cluster: cluster, err: err, name: name}
		}(contextName)
	}

	// Collect results from all goroutines
	var clusters []models.KubeCluster
	for i := 0; i < len(contextNames); i++ {
		res := <-results
		if res.err != nil {
			// If connection fails, create an offline cluster entry
			log.Printf("Failed to connect to context %s: %v", res.name, res.err)
			clusters = append(clusters, models.KubeCluster{
				ID:      fmt.Sprintf("context-%s", res.name),
				Name:    res.name,
				Version: "unknown",
				Status: models.ResourceStatus{
					Phase:       "Offline",
					Ready:       false,
					Reason:      "Connection failed",
					Message:     res.err.Error(),
					LastUpdated: time.Now(),
				},
			})
		} else {
			clusters = append(clusters, *res.cluster)
		}
	}

	return clusters, nil
}

// GetCluster returns a specific cluster by ID
func GetCluster(clusterID string) (*models.KubeCluster, error) {
	clusters, err := GetClusters()
	if err != nil {
		return nil, err
	}

	for _, cluster := range clusters {
		if cluster.ID == clusterID {
			return &cluster, nil
		}
	}

	return nil, nil
}

// GetClusterNodes returns nodes for a specific cluster
func GetClusterNodes(clusterID string) ([]models.KubeNode, error) {
	if clientset == nil {
		return nil, fmt.Errorf("kubernetes client not initialized")
	}

	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %v", err)
	}

	var kubeNodes []models.KubeNode
	for _, node := range nodes.Items {
		kubeNode := models.KubeNode{
			Name:        node.Name,
			Kubelet:     node.Status.NodeInfo.KubeletVersion,
			Runtime:     node.Status.NodeInfo.ContainerRuntimeVersion,
			Role:        getNodeRole(&node),
			Pods:        []models.KubePod{}, // Will be populated separately
			Labels:      node.Labels,
			Annotations: node.Annotations,
			CreatedAt:   node.CreationTimestamp.Time,
			Status:      getNodeStatus(&node),
			Capacity:    getNodeCapacity(node.Status.Capacity),
			Allocatable: getNodeCapacity(node.Status.Allocatable),
			Conditions:  getNodeConditions(&node),
		}
		kubeNodes = append(kubeNodes, kubeNode)
	}

	return kubeNodes, nil
}

// GetClusterPods returns pods for a specific cluster
func GetClusterPods(clusterID string) ([]models.KubePod, error) {
	if clientset == nil {
		return nil, fmt.Errorf("kubernetes client not initialized")
	}

	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %v", err)
	}

	var kubePods []models.KubePod
	for _, pod := range pods.Items {
		kubePod := models.KubePod{
			Name:         pod.Name,
			Namespace:    pod.Namespace,
			Role:         getPodRole(&pod),
			IP:           pod.Status.PodIP,
			Labels:       pod.Labels,
			NodeName:     pod.Spec.NodeName,
			CreatedAt:    pod.CreationTimestamp.Time,
			RestartCount: getTotalRestartCount(&pod),
			Status:       getPodStatus(&pod),
		}

		// Add containers with status
		for _, container := range pod.Spec.Containers {
			containerStatus := getContainerStatus(&pod, container.Name)
			kubeContainer := models.KubeContainer{
				Name:   container.Name,
				Image:  container.Image,
				Ready:  containerStatus.Ready,
				Status: containerStatus.Status,
			}
			kubePod.Containers = append(kubePod.Containers, kubeContainer)
		}

		// Add volumes
		for _, volume := range pod.Spec.Volumes {
			volumeType := getVolumeType(&volume)
			kubePod.Volumes = append(kubePod.Volumes, models.KubeVolume{
				Name: volume.Name,
				Type: volumeType,
			})
		}

		kubePods = append(kubePods, kubePod)
	}

	return kubePods, nil
}

// GetClusterServices returns services for a specific cluster
func GetClusterServices(clusterID string) ([]models.KubeService, error) {
	if clientset == nil {
		return nil, fmt.Errorf("kubernetes client not initialized")
	}

	services, err := clientset.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %v", err)
	}

	var kubeServices []models.KubeService
	for _, svc := range services.Items {
		kubeService := models.KubeService{
			Name:        svc.Name,
			Namespace:   svc.Namespace,
			Type:        string(svc.Spec.Type),
			ClusterIP:   svc.Spec.ClusterIP,
			Selector:    fmt.Sprintf("%v", svc.Spec.Selector),
			Role:        getServiceRole(&svc),
			Labels:      svc.Labels,
			CreatedAt:   svc.CreationTimestamp.Time,
			Status:      getServiceStatus(&svc),
			Pods:        []models.KubePod{}, // Can be populated with related pods
			ExternalIPs: svc.Spec.ExternalIPs,
		}

		// Handle LoadBalancer IP
		if len(svc.Status.LoadBalancer.Ingress) > 0 {
			kubeService.LoadBalancerIP = svc.Status.LoadBalancer.Ingress[0].IP
		}

		// Convert ports
		for _, port := range svc.Spec.Ports {
			servicePort := models.ServicePort{
				Name:       port.Name,
				Port:       port.Port,
				TargetPort: port.TargetPort.String(),
				Protocol:   string(port.Protocol),
			}
			if port.NodePort != 0 {
				servicePort.NodePort = port.NodePort
				kubeService.NodePort = &port.NodePort
			}
			kubeService.Ports = append(kubeService.Ports, servicePort)
		}

		kubeServices = append(kubeServices, kubeService)
	}

	return kubeServices, nil
}

// GetClusterDeployments returns deployments for a specific cluster
func GetClusterDeployments(clusterID string) ([]models.KubeDeployment, error) {
	if clientset == nil {
		return nil, fmt.Errorf("kubernetes client not initialized")
	}

	deployments, err := clientset.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %v", err)
	}

	var kubeDeployments []models.KubeDeployment
	for _, deployment := range deployments.Items {
		kubeDeployment := models.KubeDeployment{
			Name:              deployment.Name,
			Namespace:         deployment.Namespace,
			Replicas:          *deployment.Spec.Replicas,
			ReadyReplicas:     deployment.Status.ReadyReplicas,
			AvailableReplicas: deployment.Status.AvailableReplicas,
			Selector:          fmt.Sprintf("%v", deployment.Spec.Selector.MatchLabels),
			Template:          deployment.Spec.Template.Name,
			Strategy:          string(deployment.Spec.Strategy.Type),
			Role:              getDeploymentRole(&deployment),
			Labels:            deployment.Labels,
			CreatedAt:         deployment.CreationTimestamp.Time,
			Status:            getDeploymentStatus(&deployment),
		}

		kubeDeployments = append(kubeDeployments, kubeDeployment)
	}

	return kubeDeployments, nil
}

// GetClusterNamespaces returns namespaces for a specific cluster (uses default clientset)
func GetClusterNamespaces(clusterID string) ([]models.KubeNamespace, error) {
	if clientset == nil {
		return nil, fmt.Errorf("kubernetes client not initialized")
	}
	return getClusterNamespaces(clientset)
}

// getClusterNamespaces returns namespaces using the provided clientset
func getClusterNamespaces(cs *kubernetes.Clientset) ([]models.KubeNamespace, error) {
	namespaces, err := cs.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %v", err)
	}

	var kubeNamespaces []models.KubeNamespace
	for _, ns := range namespaces.Items {
		// Get deployments for this namespace
		deploymentList, err := cs.AppsV1().Deployments(ns.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Printf("Failed to get deployments for namespace %s: %v", ns.Name, err)
			continue
		}

		var deployments []models.KubeDeployment
		for _, deployment := range deploymentList.Items {
			deployments = append(deployments, models.KubeDeployment{
				Name:              deployment.Name,
				Namespace:         deployment.Namespace,
				Replicas:          *deployment.Spec.Replicas,
				ReadyReplicas:     deployment.Status.ReadyReplicas,
				AvailableReplicas: deployment.Status.AvailableReplicas,
				Selector:          fmt.Sprintf("%v", deployment.Spec.Selector.MatchLabels),
				Template:          deployment.Spec.Template.Name,
				Strategy:          string(deployment.Spec.Strategy.Type),
				Role:              getDeploymentRole(&deployment),
				Labels:            deployment.Labels,
				CreatedAt:         deployment.CreationTimestamp.Time,
				Status:            getDeploymentStatus(&deployment),
			})
		}

		// Get services for this namespace
		serviceList, err := cs.CoreV1().Services(ns.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Printf("Failed to get services for namespace %s: %v", ns.Name, err)
			continue
		}

		var services []models.KubeService
		for _, svc := range serviceList.Items {
			kubeService := models.KubeService{
				Name:        svc.Name,
				Namespace:   svc.Namespace,
				Type:        string(svc.Spec.Type),
				ClusterIP:   svc.Spec.ClusterIP,
				Selector:    fmt.Sprintf("%v", svc.Spec.Selector),
				Role:        getServiceRole(&svc),
				Labels:      svc.Labels,
				CreatedAt:   svc.CreationTimestamp.Time,
				Status:      getServiceStatus(&svc),
				Pods:        []models.KubePod{},
				ExternalIPs: svc.Spec.ExternalIPs,
			}

			// Convert ports
			for _, port := range svc.Spec.Ports {
				servicePort := models.ServicePort{
					Name:       port.Name,
					Port:       port.Port,
					TargetPort: port.TargetPort.String(),
					Protocol:   string(port.Protocol),
				}
				if port.NodePort != 0 {
					servicePort.NodePort = port.NodePort
					kubeService.NodePort = &port.NodePort
				}
				kubeService.Ports = append(kubeService.Ports, servicePort)
			}

			services = append(services, kubeService)
		}

		// Get pod count for this namespace
		podList, err := cs.CoreV1().Pods(ns.Name).List(context.TODO(), metav1.ListOptions{})
		podCount := 0
		if err == nil {
			podCount = len(podList.Items)
		}

		kubeNamespace := models.KubeNamespace{
			Name:        ns.Name,
			Deployments: deployments,
			Services:    services,
			PodCount:    podCount,
			Labels:      ns.Labels,
			CreatedAt:   ns.CreationTimestamp.Time,
			Status:      getNamespaceStatus(&ns),
		}

		kubeNamespaces = append(kubeNamespaces, kubeNamespace)
	}

	return kubeNamespaces, nil
}

// getClusterDataForContext creates a clientset for the given context and retrieves lightweight cluster data
func getClusterDataForContext(contextName string) (*models.KubeCluster, error) {
	// Build config for this specific context
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{CurrentContext: contextName},
	).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to build config for context %s: %v", contextName, err)
	}

	// Create a clientset for this context with a short timeout
	config.Timeout = 5 * time.Second
	contextClientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset for context %s: %v", contextName, err)
	}

	// Get lightweight cluster data (no detailed resources)
	return getClusterDataLightweight(contextName, contextClientset)
}

// getClusterDataLightweight returns minimal cluster info without loading all resources
func getClusterDataLightweight(contextName string, cs *kubernetes.Clientset) (*models.KubeCluster, error) {
	// Get basic cluster info
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	version, err := cs.Discovery().ServerVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get server version: %v", err)
	}

	cluster := &models.KubeCluster{
		ID:      fmt.Sprintf("context-%s", contextName),
		Name:    contextName,
		Version: version.GitVersion,
		Status: models.ResourceStatus{
			Phase:       "Running",
			Ready:       true,
			LastUpdated: time.Now(),
		},
	}

	// Get quick counts without loading full data
	// Node count
	nodes, err := cs.CoreV1().Nodes().List(ctx, metav1.ListOptions{Limit: 1000})
	if err == nil {
		cluster.Summary.TotalNodes = len(nodes.Items)
		for _, node := range nodes.Items {
			if isNodeReady(&node) {
				cluster.Summary.ReadyNodes++
			}
		}
	}

	// Namespace count
	namespaces, err := cs.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err == nil {
		cluster.Summary.TotalNamespaces = len(namespaces.Items)
	}

	// Pod count (just count, don't load full data)
	pods, err := cs.CoreV1().Pods("").List(ctx, metav1.ListOptions{Limit: 10000})
	if err == nil {
		cluster.Summary.TotalPods = len(pods.Items)
		for _, pod := range pods.Items {
			if pod.Status.Phase == v1.PodRunning {
				cluster.Summary.RunningPods++
			} else if pod.Status.Phase == v1.PodPending {
				cluster.Summary.PendingPods++
			}
		}
	}

	// Deployment count
	deployments, err := cs.AppsV1().Deployments("").List(ctx, metav1.ListOptions{Limit: 1000})
	if err == nil {
		cluster.Summary.TotalDeployments = len(deployments.Items)
	}

	// Service count
	services, err := cs.CoreV1().Services("").List(ctx, metav1.ListOptions{Limit: 1000})
	if err == nil {
		cluster.Summary.TotalServices = len(services.Items)
	}

	return cluster, nil
}

// isNodeReady checks if a node is ready
func isNodeReady(node *v1.Node) bool {
	for _, condition := range node.Status.Conditions {
		if condition.Type == v1.NodeReady {
			return condition.Status == v1.ConditionTrue
		}
	}
	return false
}

// getClusterData creates a cluster representation from the given clientset
func getClusterData(contextName string, cs *kubernetes.Clientset) (*models.KubeCluster, error) {
	// Get basic cluster info
	version, err := cs.Discovery().ServerVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get server version: %v", err)
	}

	cluster := &models.KubeCluster{
		ID:      fmt.Sprintf("context-%s", contextName),
		Name:    contextName,
		Version: version.GitVersion,
		Status:  getClusterStatus(),
	}

	// Get control plane nodes
	controlPlaneNodes, err := getControlPlaneNodes(cs)
	if err != nil {
		log.Printf("Failed to get control plane nodes for %s: %v", contextName, err)
	} else {
		cluster.ControlPlane.Nodes = controlPlaneNodes
	}

	// Get worker nodes
	workerNodes, err := getWorkerNodes(cs)
	if err != nil {
		log.Printf("Failed to get worker nodes for %s: %v", contextName, err)
	} else {
		cluster.Nodes = workerNodes
	}

	// Get namespaces with deployments and services
	namespaces, err := getClusterNamespaces(cs)
	if err != nil {
		log.Printf("Failed to get namespaces for %s: %v", contextName, err)
	} else {
		cluster.Namespaces = namespaces
	}

	// Calculate cluster summary
	cluster.Summary = calculateClusterSummary(cluster)

	return cluster, nil
}

// getControlPlaneNodes returns control plane nodes
func getControlPlaneNodes(cs *kubernetes.Clientset) ([]models.KubeNode, error) {
	nodes, err := cs.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{
		LabelSelector: "node-role.kubernetes.io/control-plane",
	})
	if err != nil {
		return nil, err
	}

	var controlPlaneNodes []models.KubeNode
	for _, node := range nodes.Items {
		controlPlaneNodes = append(controlPlaneNodes, models.KubeNode{
			Kubelet: node.Status.NodeInfo.KubeletVersion,
			Runtime: node.Status.NodeInfo.ContainerRuntimeVersion,
			Role:    "control-plane",
		})
	}

	return controlPlaneNodes, nil
}

// getWorkerNodes returns worker nodes
func getWorkerNodes(cs *kubernetes.Clientset) ([]models.KubeNode, error) {
	nodes, err := cs.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var workerNodes []models.KubeNode
	for _, node := range nodes.Items {
		// Skip control plane nodes
		if _, isControlPlane := node.Labels["node-role.kubernetes.io/control-plane"]; isControlPlane {
			continue
		}

		workerNodes = append(workerNodes, models.KubeNode{
			Kubelet: node.Status.NodeInfo.KubeletVersion,
			Runtime: node.Status.NodeInfo.ContainerRuntimeVersion,
			Role:    "worker",
		})
	}

	return workerNodes, nil
}

// Helper functions for determining roles
func getNodeRole(node *v1.Node) string {
	if _, ok := node.Labels["node-role.kubernetes.io/control-plane"]; ok {
		return "control-plane"
	}
	return "worker"
}

func getPodRole(pod *v1.Pod) string {
	if app, ok := pod.Labels["app"]; ok {
		return app
	}
	if k8sApp, ok := pod.Labels["k8s-app"]; ok {
		return k8sApp
	}
	return "unknown"
}

func getDeploymentRole(deployment *appsv1.Deployment) string {
	if app, ok := deployment.Labels["app"]; ok {
		return app
	}
	return "unknown"
}

func getServiceRole(service *v1.Service) string {
	if app, ok := service.Labels["app"]; ok {
		return app
	}
	return "unknown"
}

// Helper functions for status and metrics

func getPodStatus(pod *v1.Pod) models.ResourceStatus {
	phase := string(pod.Status.Phase)
	ready := isPodReady(pod)
	reason := pod.Status.Reason
	message := pod.Status.Message

	return models.ResourceStatus{
		Phase:       phase,
		Ready:       ready,
		Reason:      reason,
		Message:     message,
		LastUpdated: time.Now(), // In a real implementation, this would be the last transition time
	}
}

func isPodReady(pod *v1.Pod) bool {
	for _, condition := range pod.Status.Conditions {
		if condition.Type == v1.PodReady {
			return condition.Status == v1.ConditionTrue
		}
	}
	return false
}

func getContainerStatus(pod *v1.Pod, containerName string) struct {
	Ready  bool
	Status models.ResourceStatus
} {
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if containerStatus.Name == containerName {
			status := models.ResourceStatus{
				Ready:       containerStatus.Ready,
				LastUpdated: time.Now(),
			}

			if containerStatus.State.Running != nil {
				status.Phase = "Running"
			} else if containerStatus.State.Waiting != nil {
				status.Phase = "Waiting"
				status.Reason = containerStatus.State.Waiting.Reason
				status.Message = containerStatus.State.Waiting.Message
			} else if containerStatus.State.Terminated != nil {
				status.Phase = "Terminated"
				status.Reason = containerStatus.State.Terminated.Reason
				status.Message = containerStatus.State.Terminated.Message
			} else {
				status.Phase = "Unknown"
			}

			return struct {
				Ready  bool
				Status models.ResourceStatus
			}{
				Ready:  containerStatus.Ready,
				Status: status,
			}
		}
	}

	return struct {
		Ready  bool
		Status models.ResourceStatus
	}{
		Ready: false,
		Status: models.ResourceStatus{
			Phase:       "Unknown",
			Ready:       false,
			LastUpdated: time.Now(),
		},
	}
}

func getTotalRestartCount(pod *v1.Pod) int32 {
	var total int32
	for _, containerStatus := range pod.Status.ContainerStatuses {
		total += containerStatus.RestartCount
	}
	return total
}

func getVolumeType(volume *v1.Volume) string {
	if volume.HostPath != nil {
		return "hostPath"
	} else if volume.ConfigMap != nil {
		return "configMap"
	} else if volume.Secret != nil {
		return "secret"
	} else if volume.PersistentVolumeClaim != nil {
		return "persistentVolumeClaim"
	} else if volume.EmptyDir != nil {
		return "emptyDir"
	} else if volume.Projected != nil {
		return "projected"
	} else if volume.DownwardAPI != nil {
		return "downwardAPI"
	}
	return "unknown"
}

func getNodeStatus(node *v1.Node) models.ResourceStatus {
	ready := false
	reason := ""
	message := ""

	for _, condition := range node.Status.Conditions {
		if condition.Type == v1.NodeReady {
			ready = condition.Status == v1.ConditionTrue
			reason = condition.Reason
			message = condition.Message
			break
		}
	}

	phase := "Ready"
	if !ready {
		phase = "NotReady"
	}

	return models.ResourceStatus{
		Phase:       phase,
		Ready:       ready,
		Reason:      reason,
		Message:     message,
		LastUpdated: time.Now(),
	}
}

func getNodeCapacity(capacity v1.ResourceList) models.NodeCapacity {
	return models.NodeCapacity{
		CPU:              capacity.Cpu().String(),
		Memory:           capacity.Memory().String(),
		Pods:             capacity.Pods().String(),
		EphemeralStorage: capacity.StorageEphemeral().String(),
		Storage:          capacity.Storage().String(),
	}
}

func getNodeConditions(node *v1.Node) []models.NodeCondition {
	var conditions []models.NodeCondition
	for _, condition := range node.Status.Conditions {
		conditions = append(conditions, models.NodeCondition{
			Type:               string(condition.Type),
			Status:             string(condition.Status),
			Reason:             condition.Reason,
			Message:            condition.Message,
			LastTransitionTime: condition.LastTransitionTime.Time,
		})
	}
	return conditions
}

func getServiceStatus(service *v1.Service) models.ResourceStatus {
	// Services are generally always ready if they exist
	return models.ResourceStatus{
		Phase:       "Active",
		Ready:       true,
		LastUpdated: time.Now(),
	}
}

func getDeploymentStatus(deployment *appsv1.Deployment) models.ResourceStatus {
	ready := deployment.Status.ReadyReplicas == *deployment.Spec.Replicas
	phase := "Available"
	reason := ""
	message := ""

	if !ready {
		phase = "Progressing"
		if deployment.Status.ReadyReplicas == 0 {
			phase = "Unavailable"
		}
	}

	// Check deployment conditions for more details
	for _, condition := range deployment.Status.Conditions {
		if condition.Type == appsv1.DeploymentProgressing {
			reason = condition.Reason
			message = condition.Message
			break
		}
	}

	return models.ResourceStatus{
		Phase:       phase,
		Ready:       ready,
		Reason:      reason,
		Message:     message,
		LastUpdated: time.Now(),
	}
}

func getNamespaceStatus(namespace *v1.Namespace) models.ResourceStatus {
	phase := string(namespace.Status.Phase)
	ready := phase == "Active"

	return models.ResourceStatus{
		Phase:       phase,
		Ready:       ready,
		LastUpdated: time.Now(),
	}
}

func getClusterStatus() models.ResourceStatus {
	// Simple cluster status - in a real implementation this would check various cluster components
	return models.ResourceStatus{
		Phase:       "Running",
		Ready:       true,
		LastUpdated: time.Now(),
	}
}

func calculateClusterSummary(cluster *models.KubeCluster) models.ClusterSummary {
	summary := models.ClusterSummary{}

	// Count nodes
	summary.TotalNodes = len(cluster.Nodes) + len(cluster.ControlPlane.Nodes)
	for _, node := range cluster.Nodes {
		if node.Status.Ready {
			summary.ReadyNodes++
		}
	}
	for _, node := range cluster.ControlPlane.Nodes {
		if node.Status.Ready {
			summary.ReadyNodes++
		}
	}

	// Count resources from namespaces
	summary.TotalNamespaces = len(cluster.Namespaces)
	for _, namespace := range cluster.Namespaces {
		summary.TotalDeployments += len(namespace.Deployments)
		summary.TotalServices += len(namespace.Services)
		summary.TotalPods += namespace.PodCount

		// Count pod states (this is simplified - in reality we'd get actual pod status)
		for _, deployment := range namespace.Deployments {
			if deployment.Status.Ready {
				summary.RunningPods += int(deployment.ReadyReplicas)
			} else {
				summary.PendingPods += int(deployment.Replicas - deployment.ReadyReplicas)
			}
		}
	}

	// Calculate utilization (placeholder values - real implementation would use metrics)
	summary.CPUUtilization = 45.5    // percentage
	summary.MemoryUtilization = 62.3 // percentage

	return summary
}
