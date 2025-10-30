import { useState, useEffect } from 'react'
import { KubeClusterList } from "@/containers"
import { ComparisonBanner, ViewToggle } from "@/components/ui"
import { DefaultLayout } from "@/components/layout"
import { useKeydownShortcut } from '@/hooks/useKeydownShortcut'
import { getClusters, ApiError } from '@/services'
import type { KubeCluster } from '@/types/kube'

export default function HomePage() {
  const [selectedClusters, setSelectedClusters] = useState<string[]>([])
  const [view, setView] = useState<'card' | 'list'>('card')
  const [clusters, setClusters] = useState<KubeCluster[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const handleSelectionChange = (selectedIds: string[]) => {
    setSelectedClusters(selectedIds)
  }

  // Fetch clusters from API
  useEffect(() => {
    const fetchClusters = async () => {
      try {
        setLoading(true)
        setError(null)
        const data = await getClusters()
        setClusters(data)
      } catch (err) {
        if (err instanceof ApiError) {
          setError(err.message)
        } else if (err instanceof Error) {
          setError(err.message)
        } else {
          setError('Failed to fetch clusters')
        }
        console.error('Error fetching clusters:', err)
      } finally {
        setLoading(false)
      }
    }

    fetchClusters()
  }, [])

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

        {loading && (
          <div className="flex justify-center items-center py-12">
            <div className="text-muted-foreground">Loading clusters...</div>
          </div>
        )}

        {error && (
          <div className="bg-destructive/10 border border-destructive/20 rounded-lg p-4">
            <div className="flex items-start gap-3">
              <div className="flex-1">
                <h3 className="font-semibold text-destructive mb-1">Error loading clusters</h3>
                <p className="text-sm text-destructive/80">{error}</p>
                <p className="text-sm text-muted-foreground mt-2">
                  Make sure the backend server is running on http://localhost:8080
                </p>
              </div>
            </div>
          </div>
        )}

        {!loading && !error && (
          <KubeClusterList
            clusters={clusters}
            selectedClusters={selectedClusters}
            onSelectionChange={handleSelectionChange}
            view={view}
            emptyMessage="No clusters available"
          />
        )}
      </div>
    </DefaultLayout>
  )
}
