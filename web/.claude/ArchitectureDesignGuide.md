# Architecture Design Guide: Kubey

## High-level Architecture

Kubey follows a modern React-based frontend architecture designed for scalability, maintainability, and performance. The application is built with a component-driven approach using TypeScript for type safety and Tailwind CSS for styling.

### Technology Stack

**Frontend Framework:**
- React 19.1.1 with TypeScript
- Vite for build tooling and development server
- React Router for client-side routing

**UI Framework:**
- Tailwind CSS for utility-first styling
- shadcn/ui for pre-built accessible components
- Lucide React for consistent iconography
- next-themes for dark/light mode support

**Development Tools:**
- ESLint for code linting
- TypeScript for static type checking
- Vite for fast development and optimized builds

### Architecture Layers

```
┌─────────────────┐
│   Pages         │  Route-level components
├─────────────────┤
│   Containers    │  Business logic and data fetching
├─────────────────┤
│   Components    │  Reusable UI components
│   ├── Layout    │  Layout and navigation
│   ├── UI        │  Atomic UI primitives
│   └── Theme     │  Theme management
├─────────────────┤
│   Hooks         │  Custom React hooks
├─────────────────┤
│   Services      │  API and external service integration
├─────────────────┤
│   Types         │  TypeScript type definitions
├─────────────────┤
│   Lib           │  Utility functions and configurations
└─────────────────┘
```

## Component Architecture

### Design Principles

1. **Atomic Design Pattern**: Components are organized hierarchically from atoms to pages
2. **Single Responsibility**: Each component has one clear purpose
3. **Composition over Inheritance**: Complex components built from simpler ones
4. **Type Safety**: Comprehensive TypeScript typing throughout

### Component Organization

#### UI Components (`src/components/ui/`)
Atomic, reusable components that form the foundation of the UI:

- **CardItem**: Basic card component with selection state
- **KubeClusterItem**: Specialized card for cluster display
- **SearchBar**: Search input with keyboard shortcuts
- **ThemeToggle**: Dark/light mode switcher
- **UserMenu**: User account and settings menu

#### Layout Components (`src/components/layout/`)
Page structure and navigation components:

- **Header**: Top navigation with logo, search, and user controls
- **DefaultLayout**: Standard page layout wrapper
- **CenteredFullScreenLayout**: Full-screen centered content layout

#### Container Components (`src/containers/`)
Business logic components that manage data and state:

- **CardList**: Generic list container with selection logic
- **KubeClusterList**: Cluster-specific list with filtering and sorting

#### Page Components (`src/pages/`)
Route-level components that compose containers and layouts:

- **HomePage**: Main dashboard view
- **NotFoundPage**: 404 error page

### Component Patterns

#### 1. Compound Components
```typescript
// Example: ClusterList with selectable items
<KubeClusterList clusters={clusters}>
  <KubeClusterList.Header>Production Clusters</KubeClusterList.Header>
  <KubeClusterList.Items />
  <KubeClusterList.Footer>Showing {count} clusters</KubeClusterList.Footer>
</KubeClusterList>
```

#### 2. Render Props for Flexibility
```typescript
// Example: Customizable data display
<DataTable data={clusters}>
  {(item) => (
    <ClusterRow cluster={item} onSelect={handleSelect} />
  )}
</DataTable>
```

#### 3. Custom Hooks for Logic Reuse
```typescript
// Example: Cluster data management
const useClusterData = () => {
  const [clusters, setClusters] = useState<KubeCluster[]>([])
  const [loading, setLoading] = useState(false)

  const fetchClusters = useCallback(async () => {
    setLoading(true)
    try {
      const data = await clusterApi.getAll()
      setClusters(data)
    } finally {
      setLoading(false)
    }
  }, [])

  return { clusters, loading, fetchClusters }
}
```

## Data Flow

### Current Implementation
```
User Interaction → Component Event → Container Logic → State Update → Re-render
```

### Future Data Flow (with Backend)
```
User Interaction → Component Event → Action Creator → API Call → State Update → Re-render
     ↓
  Real-time Updates (WebSocket) → State Update → Re-render
```

### Data Types Hierarchy

```typescript
interface KubeCluster {
  id: string
  name: string
  controlPlane: KubeControlPlane
  nodes: KubeNode[]
  namespaces: KubeNamespace[]
}

interface KubeNode {
  kubelet: string
  runtime: string
  role: string
  pods: KubePod[]
}

interface KubePod {
  containers: KubeContainer[]
  volumes: KubeVolume[]
  role: string
  ip: string
  labels: Record<string, string>
}
```

### State Management Strategy

#### Local Component State
- UI interaction state (selections, form inputs)
- Component-specific loading and error states
- Temporary data for optimistic updates

#### Global Application State
- User preferences (theme, language)
- Authentication state
- Global notifications and alerts

#### Cluster Data State
- Cluster list and metadata
- Cached cluster details
- Real-time metrics and health status

### State Management Implementation

**Current Approach:**
- React `useState` and `useReducer` for local state
- React Context for theme and user preferences
- URL state for filters and selections

**Future Approach:**
- Zustand or Redux Toolkit for complex state management
- React Query for server state management
- Local storage for user preferences persistence

## API Integration

### Current State
- No backend integration (static sample data)
- Component-level data simulation

### Future API Design

#### RESTful Endpoints
```
GET    /api/v1/clusters           # List all clusters
GET    /api/v1/clusters/:id       # Get cluster details
GET    /api/v1/clusters/:id/metrics # Get cluster metrics
POST   /api/v1/clusters           # Register new cluster
PUT    /api/v1/clusters/:id       # Update cluster config
DELETE /api/v1/clusters/:id       # Remove cluster
```

#### Real-time Updates
- WebSocket connection for live metrics
- Server-sent events for alerts and notifications
- Polling fallback for environments without WebSocket support

#### Authentication
- JWT-based authentication
- Cluster-specific API keys
- OAuth integration for enterprise SSO

### Service Layer Architecture

```typescript
// Service classes for API abstraction
class ClusterService {
  async getAll(): Promise<KubeCluster[]> {
    // API call logic
  }

  async getById(id: string): Promise<KubeCluster> {
    // API call logic
  }

  async getMetrics(id: string): Promise<ClusterMetrics> {
    // API call logic
  }
}

// React Query integration
const useClusters = () => {
  return useQuery({
    queryKey: ['clusters'],
    queryFn: () => clusterService.getAll(),
    staleTime: 5 * 60 * 1000, // 5 minutes
  })
}
```

## Component Design Patterns

### 1. Props Interface Design
```typescript
interface ClusterListProps {
  clusters: KubeCluster[]
  selectedClusters?: string[]
  onSelectionChange?: (selectedIds: string[]) => void
  loading?: boolean
  error?: string
  emptyMessage?: string
  gridCols?: ResponsiveGridCols
}
```

### 2. Component Composition
```typescript
// Base component with common functionality
function BaseCard({ children, className, ...props }: CardProps) {
  return (
    <div className={cn("card-styles", className)} {...props}>
      {children}
    </div>
  )
}

// Specialized component extending base
function ClusterCard({ cluster, ...props }: ClusterCardProps) {
  return (
    <BaseCard {...props}>
      <ClusterHeader cluster={cluster} />
      <ClusterMetrics cluster={cluster} />
    </BaseCard>
  )
}
```

### 3. Custom Hooks for Complex Logic
```typescript
function useClusterSearch(clusters: KubeCluster[]) {
  const [query, setQuery] = useState('')
  const [filteredClusters, setFilteredClusters] = useState(clusters)

  useEffect(() => {
    const filtered = clusters.filter(cluster =>
      cluster.name.toLowerCase().includes(query.toLowerCase())
    )
    setFilteredClusters(filtered)
  }, [clusters, query])

  return { query, setQuery, filteredClusters }
}
```

### 4. Error Boundaries
```typescript
class ClusterErrorBoundary extends Component {
  state = { hasError: false }

  static getDerivedStateFromError(error: Error) {
    return { hasError: true }
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    logError(error, errorInfo)
  }

  render() {
    if (this.state.hasError) {
      return <ClusterErrorFallback />
    }
    return this.props.children
  }
}
```

## Best Practices

### Code Organization
- **File Structure**: Group related components in feature folders
- **Naming Conventions**: PascalCase for components, camelCase for files
- **Import Order**: External libraries → internal modules → relative imports

### Performance Optimization
- **Memoization**: Use `React.memo` for expensive components
- **Virtualization**: Implement virtual scrolling for large lists
- **Code Splitting**: Lazy load route components
- **Bundle Analysis**: Monitor bundle size and optimize imports

### Accessibility
- **Semantic HTML**: Use appropriate ARIA labels and roles
- **Keyboard Navigation**: Ensure all interactive elements are keyboard accessible
- **Screen Reader Support**: Provide descriptive text for icons and images
- **Color Contrast**: Maintain WCAG AA compliance for text and UI elements

### Testing Strategy
- **Unit Tests**: Test individual components and hooks
- **Integration Tests**: Test component interactions
- **E2E Tests**: Test complete user workflows
- **Visual Regression**: Ensure UI consistency across changes

### Development Workflow
- **TypeScript Strict**: Enable strict type checking
- **ESLint Rules**: Enforce consistent code style
- **Pre-commit Hooks**: Run linting and tests before commits
- **Code Reviews**: Require review for all changes

## Future Considerations

### Scalability
- **Micro-frontend Architecture**: Split into independently deployable modules
- **Service Worker**: Implement offline capabilities and caching
- **Progressive Web App**: Add PWA features for mobile access

### Advanced Features
- **Plugin System**: Allow third-party extensions
- **Multi-tenancy**: Support multiple organizations
- **Advanced Analytics**: Machine learning insights for cluster optimization

### Infrastructure
- **CI/CD Pipeline**: Automated testing and deployment
- **Monitoring**: Application performance monitoring
- **Logging**: Centralized logging and error tracking

### Security Enhancements
- **CSP Headers**: Content Security Policy implementation
- **Input Validation**: Comprehensive input sanitization
- **Audit Trail**: Track all user actions and system changes

---

*Last Updated: October 11, 2025*