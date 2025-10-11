import { Check } from 'lucide-react'
import type { KubeCluster } from '@/types/kube'

interface KubeClusterListItemProps {
  cluster: KubeCluster
  isSelected?: boolean
  onSelect?: (id: string, selected: boolean) => void
}

export default function KubeClusterListItem({
  cluster,
  isSelected = false,
  onSelect
}: KubeClusterListItemProps) {
  const totalPods = cluster.nodes.flatMap(n => n.pods).length
  const totalNodes = cluster.nodes.length
  const totalNamespaces = cluster.namespaces.length
  const totalDeployments = cluster.namespaces.flatMap(ns => ns.deployments).length
  const healthyDeployments = cluster.namespaces
    .flatMap(ns => ns.deployments)
    .filter(d => d.replicas > 0)
    .length

  const handleSelect = () => {
    onSelect?.(cluster.id, !isSelected)
  }

  return (
    <div
      className={`flex items-center px-4 py-3 border rounded-lg transition-colors cursor-pointer hover:bg-muted/50 ${
        isSelected ? 'border-primary bg-primary/5' : 'border-border'
      }`}
      onClick={handleSelect}
    >
      {/* Checkbox */}
      <div className="w-8 flex items-center">
        <div className={`w-4 h-4 border rounded flex items-center justify-center transition-colors ${
          isSelected
            ? 'border-primary bg-primary text-primary-foreground'
            : 'border-muted-foreground/30 bg-background'
        }`}>
          {isSelected && <Check className="w-3 h-3" />}
        </div>
      </div>

      {/* Cluster Name */}
      <div className="flex-1 font-medium text-foreground truncate min-w-0">
        {cluster.name}
      </div>

      {/* Metrics - Hidden on mobile, shown on larger screens */}
      <div className="hidden md:flex items-center gap-6 text-sm text-muted-foreground">
        <div className="w-16 text-center">
          <span className="font-medium text-foreground">{totalPods}</span>
          <span className="ml-1">pods</span>
        </div>
        <div className="w-16 text-center">
          <span className="font-medium text-foreground">{totalNodes}</span>
          <span className="ml-1">nodes</span>
        </div>
        <div className="w-20 text-center">
          <span className="font-medium text-foreground">{totalNamespaces}</span>
          <span className="ml-1">ns</span>
        </div>
        <div className="w-24 text-center">
          <span className="font-medium text-foreground">{totalDeployments}</span>
          <span className={`ml-1 ${
            healthyDeployments === totalDeployments
              ? 'text-green-600 dark:text-green-400'
              : 'text-yellow-600 dark:text-yellow-400'
          }`}>
            ({healthyDeployments} healthy)
          </span>
        </div>
      </div>

      {/* Status Indicator */}
      <div className="w-8 flex justify-center ml-4">
        <div
          className={`w-2 h-2 rounded-full ${
            healthyDeployments === totalDeployments
              ? 'bg-green-500'
              : 'bg-yellow-500'
          }`}
          title={healthyDeployments === totalDeployments ? 'Healthy' : 'Issues detected'}
        />
      </div>

      {/* Mobile Summary - Shown only on mobile */}
      <div className="md:hidden flex items-center gap-2 text-xs text-muted-foreground ml-2">
        <span>{totalPods}p</span>
        <span>•</span>
        <span>{totalNodes}n</span>
        <span>•</span>
        <span>{totalNamespaces}ns</span>
      </div>
    </div>
  )
}