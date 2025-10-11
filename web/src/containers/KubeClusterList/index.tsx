import { KubeClusterItem, KubeClusterListItem } from '@/components/ui'
import type { KubeCluster } from '@/types/kube'
import { useState, useEffect } from 'react'

interface KubeClusterListProps {
  clusters: KubeCluster[]
  selectedClusters?: string[]
  onSelectionChange?: (selectedIds: string[]) => void
  view?: 'card' | 'list'
  emptyMessage?: string
  gridCols?: {
    default: number
    sm?: number
    lg?: number
    xl?: number
  }
}

export default function KubeClusterList({
  clusters,
  selectedClusters = [],
  onSelectionChange,
  view = 'card',
  emptyMessage = "No clusters available",
}: KubeClusterListProps) {
  const [internalSelection, setInternalSelection] = useState<string[]>(selectedClusters)

  useEffect(() => {
    setInternalSelection(selectedClusters)
  }, [selectedClusters])

  const handleClusterSelect = (id: string, selected: boolean) => {
    const newSelection = selected
      ? [...internalSelection, id]
      : internalSelection.filter(clusterId => clusterId !== id)

    setInternalSelection(newSelection)
    onSelectionChange?.(newSelection)
  }

  if (clusters.length === 0) {
    return (
      <div className="flex items-center justify-center py-12">
        <p className="text-muted-foreground">{emptyMessage}</p>
      </div>
    )
  }

  if (view === 'list') {
    return (
      <div className="space-y-2">
        {/* Optional header for desktop */}
        <div className="hidden md:flex px-4 py-2 text-sm font-medium text-muted-foreground border-b">
          <div className="w-8"></div>
          <div className="flex-1">Cluster</div>
          <div className="flex items-center gap-6">
            <div className="w-16 text-center">Pods</div>
            <div className="w-16 text-center">Nodes</div>
            <div className="w-20 text-center">Namespaces</div>
            <div className="w-24 text-center">Deployments</div>
          </div>
          <div className="w-8 text-center">Status</div>
        </div>
        {clusters.map((cluster) => (
          <KubeClusterListItem
            key={cluster.id}
            cluster={cluster}
            isSelected={internalSelection.includes(cluster.id)}
            onSelect={handleClusterSelect}
          />
        ))}
      </div>
    )
  }

  return (
    <div className="space-y-4">
      <div className="grid gap-6 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
        {clusters.map((cluster) => (
          <KubeClusterItem
            key={cluster.id}
            cluster={cluster}
            isSelected={internalSelection.includes(cluster.id)}
            onSelect={handleClusterSelect}
          />
        ))}
      </div>
    </div>
  )
}
