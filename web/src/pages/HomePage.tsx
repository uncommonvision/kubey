import { useState } from 'react'
import { KubeClusterList } from "@/containers"
import { ComparisonBanner, ViewToggle } from "@/components/ui"
import { sampleClusters } from "@/data/sampleClusters"
import { DefaultLayout } from "@/components/layout"
import { useKeydownShortcut } from '@/hooks/useKeydownShortcut'

export default function HomePage() {
  const [selectedClusters, setSelectedClusters] = useState<string[]>([])
  const [view, setView] = useState<'card' | 'list'>('card')

  const handleSelectionChange = (selectedIds: string[]) => {
    setSelectedClusters(selectedIds)
  }

  // Add view toggle keybind
  useKeydownShortcut(
    { key: 'v', ctrl: false, alt: false, shift: false, meta: false },
    () => setView(prev => prev === 'card' ? 'list' : 'card'),
    'Toggle View',
    'Switch between card and list view'
  )

  return (
    <DefaultLayout>
      <div className="space-y-6">
        <div className="flex flex-col sm:flex-row sm:justify-between sm:items-start gap-4">
          <div>
            <h1 className="text-3xl font-bold text-foreground mb-2">
              Kubernetes Clusters
            </h1>
            <p className="text-lg text-muted-foreground">
              Monitor and manage your Kubernetes infrastructure across environments.
            </p>
          </div>
          <ViewToggle view={view} onViewChange={setView} />
        </div>

        <ComparisonBanner selectedItems={selectedClusters} />

        <KubeClusterList
          clusters={sampleClusters}
          selectedClusters={selectedClusters}
          onSelectionChange={handleSelectionChange}
          view={view}
          emptyMessage="No clusters available"
        />
      </div>
    </DefaultLayout>
  )
}
