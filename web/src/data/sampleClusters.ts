import type { KubeCluster } from '@/types/kube'

export const sampleClusters: KubeCluster[] = [
  {
    id: 'prod-east-01',
    name: 'Production East',
    controlPlane: {
      nodes: [
        {
          kubelet: 'v1.28.2',
          runtime: 'containerd://1.7.1',
          role: 'control-plane',
          pods: [
            {
              containers: [
                { name: 'kube-apiserver', image: 'registry.k8s.io/kube-apiserver:v1.28.2' },
                { name: 'kube-controller-manager', image: 'registry.k8s.io/kube-controller-manager:v1.28.2' },
                { name: 'kube-scheduler', image: 'registry.k8s.io/kube-scheduler:v1.28.2' },
                { name: 'etcd', image: 'registry.k8s.io/etcd:3.5.9-0' }
              ],
              volumes: [
                { name: 'etcd-data', type: 'hostPath' },
                { name: 'kubelet-config', type: 'configMap' }
              ],
              role: 'control-plane',
              ip: '10.0.1.10',
              labels: { 'node-role.kubernetes.io/control-plane': '' }
            }
          ]
        }
      ]
    },
    nodes: [
      {
        kubelet: 'v1.28.2',
        runtime: 'containerd://1.7.1',
        role: 'worker',
        pods: [
          {
            containers: [
              { name: 'nginx', image: 'nginx:1.24-alpine' },
              { name: 'sidecar', image: 'envoyproxy/envoy:v1.26-latest' }
            ],
            volumes: [{ name: 'nginx-config', type: 'configMap' }],
            role: 'web',
            ip: '10.0.1.11',
            labels: { 'app': 'web-frontend', 'tier': 'frontend' }
          },
          {
            containers: [
              { name: 'api-server', image: 'node:18-alpine' }
            ],
            volumes: [{ name: 'app-config', type: 'secret' }],
            role: 'api',
            ip: '10.0.1.12',
            labels: { 'app': 'api-server', 'tier': 'backend' }
          },
          {
            containers: [
              { name: 'postgres', image: 'postgres:15-alpine' }
            ],
            volumes: [
              { name: 'postgres-data', type: 'persistentVolumeClaim' },
              { name: 'postgres-config', type: 'configMap' }
            ],
            role: 'database',
            ip: '10.0.1.13',
            labels: { 'app': 'database', 'tier': 'data' }
          }
        ]
      },
      {
        kubelet: 'v1.28.2',
        runtime: 'containerd://1.7.1',
        role: 'worker',
        pods: [
          {
            containers: [
              { name: 'redis', image: 'redis:7-alpine' }
            ],
            volumes: [{ name: 'redis-data', type: 'persistentVolumeClaim' }],
            role: 'cache',
            ip: '10.0.1.14',
            labels: { 'app': 'redis', 'tier': 'cache' }
          },
          {
            containers: [
              { name: 'worker', image: 'python:3.11-slim' }
            ],
            volumes: [{ name: 'worker-config', type: 'configMap' }],
            role: 'worker',
            ip: '10.0.1.15',
            labels: { 'app': 'job-worker', 'tier': 'processing' }
          }
        ]
      },
      {
        kubelet: 'v1.28.2',
        runtime: 'containerd://1.7.1',
        role: 'worker',
        pods: [
          {
            containers: [
              { name: 'prometheus', image: 'prom/prometheus:v2.45.0' },
              { name: 'config-reloader', image: 'jimmidyson/configmap-reload:v0.8.0' }
            ],
            volumes: [
              { name: 'prometheus-config', type: 'configMap' },
              { name: 'prometheus-data', type: 'persistentVolumeClaim' }
            ],
            role: 'monitoring',
            ip: '10.0.1.16',
            labels: { 'app': 'prometheus', 'tier': 'monitoring' }
          },
          {
            containers: [
              { name: 'grafana', image: 'grafana/grafana:9.5.2' }
            ],
            volumes: [{ name: 'grafana-data', type: 'persistentVolumeClaim' }],
            role: 'monitoring',
            ip: '10.0.1.17',
            labels: { 'app': 'grafana', 'tier': 'monitoring' }
          }
        ]
      }
    ],
    namespaces: [
      {
        deployments: [
          {
            replicas: 3,
            selector: 'app=web-frontend',
            template: 'web-frontend-template',
            strategy: 'RollingUpdate',
            role: 'frontend'
          },
          {
            replicas: 2,
            selector: 'app=api-server',
            template: 'api-server-template',
            strategy: 'RollingUpdate',
            role: 'backend'
          },
          {
            replicas: 1,
            selector: 'app=database',
            template: 'postgres-template',
            strategy: 'Recreate',
            role: 'data'
          }
        ],
        services: [
          {
            pods: [],
            clusterIP: '10.96.1.1',
            nodePort: 30001,
            selector: 'app=web-frontend',
            ports: {
              80: {
                containers: [{ name: 'nginx', image: 'nginx:1.24-alpine' }],
                volumes: [],
                role: 'web',
                ip: '10.0.1.11',
                labels: { 'app': 'web-frontend' }
              },
              443: {
                containers: [{ name: 'nginx-ssl', image: 'nginx:1.24-alpine' }],
                volumes: [],
                role: 'web',
                ip: '10.0.1.11',
                labels: { 'app': 'web-frontend' }
              }
            },
            role: 'frontend'
          },
          {
            pods: [],
            clusterIP: '10.96.1.2',
            nodePort: 30002,
            selector: 'app=api-server',
            ports: {
              8080: {
                containers: [{ name: 'api-server', image: 'node:18-alpine' }],
                volumes: [],
                role: 'api',
                ip: '10.0.1.12',
                labels: { 'app': 'api-server' }
              }
            },
            role: 'backend'
          }
        ]
      },
      {
        deployments: [
          {
            replicas: 1,
            selector: 'app=redis',
            template: 'redis-template',
            strategy: 'Recreate',
            role: 'cache'
          },
          {
            replicas: 2,
            selector: 'app=job-worker',
            template: 'worker-template',
            strategy: 'RollingUpdate',
            role: 'processing'
          }
        ],
        services: [
          {
            pods: [],
            clusterIP: '10.96.2.1',
            nodePort: 30003,
            selector: 'app=redis',
            ports: {
              6379: {
                containers: [{ name: 'redis', image: 'redis:7-alpine' }],
                volumes: [],
                role: 'cache',
                ip: '10.0.1.14',
                labels: { 'app': 'redis' }
              }
            },
            role: 'cache'
          }
        ]
      },
      {
        deployments: [
          {
            replicas: 1,
            selector: 'app=prometheus',
            template: 'prometheus-template',
            strategy: 'Recreate',
            role: 'monitoring'
          },
          {
            replicas: 1,
            selector: 'app=grafana',
            template: 'grafana-template',
            strategy: 'Recreate',
            role: 'monitoring'
          }
        ],
        services: [
          {
            pods: [],
            clusterIP: '10.96.3.1',
            nodePort: 30004,
            selector: 'app=prometheus',
            ports: {
              9090: {
                containers: [{ name: 'prometheus', image: 'prom/prometheus:v2.45.0' }],
                volumes: [],
                role: 'monitoring',
                ip: '10.0.1.16',
                labels: { 'app': 'prometheus' }
              }
            },
            role: 'monitoring'
          },
          {
            pods: [],
            clusterIP: '10.96.3.2',
            nodePort: 30005,
            selector: 'app=grafana',
            ports: {
              3000: {
                containers: [{ name: 'grafana', image: 'grafana/grafana:9.5.2' }],
                volumes: [],
                role: 'monitoring',
                ip: '10.0.1.17',
                labels: { 'app': 'grafana' }
              }
            },
            role: 'monitoring'
          }
        ]
      }
    ]
  },
  {
    id: 'staging-west-01',
    name: 'Staging West',
    controlPlane: {
      nodes: [
        {
          kubelet: 'v1.28.1',
          runtime: 'containerd://1.7.0',
          role: 'control-plane',
          pods: [
            {
              containers: [
                { name: 'kube-apiserver', image: 'registry.k8s.io/kube-apiserver:v1.28.1' },
                { name: 'kube-controller-manager', image: 'registry.k8s.io/kube-controller-manager:v1.28.1' },
                { name: 'kube-scheduler', image: 'registry.k8s.io/kube-scheduler:v1.28.1' },
                { name: 'etcd', image: 'registry.k8s.io/etcd:3.5.9-0' }
              ],
              volumes: [
                { name: 'etcd-data', type: 'hostPath' }
              ],
              role: 'control-plane',
              ip: '10.1.1.10',
              labels: { 'node-role.kubernetes.io/control-plane': '' }
            }
          ]
        }
      ]
    },
    nodes: [
      {
        kubelet: 'v1.28.1',
        runtime: 'containerd://1.7.0',
        role: 'worker',
        pods: [
          {
            containers: [
              { name: 'nginx-staging', image: 'nginx:1.24-alpine' }
            ],
            volumes: [{ name: 'nginx-config', type: 'configMap' }],
            role: 'web',
            ip: '10.1.1.11',
            labels: { 'app': 'web-frontend', 'environment': 'staging' }
          },
          {
            containers: [
              { name: 'api-staging', image: 'node:18-alpine' }
            ],
            volumes: [{ name: 'app-config', type: 'secret' }],
            role: 'api',
            ip: '10.1.1.12',
            labels: { 'app': 'api-server', 'environment': 'staging' }
          }
        ]
      },
      {
        kubelet: 'v1.28.1',
        runtime: 'containerd://1.7.0',
        role: 'worker',
        pods: [
          {
            containers: [
              { name: 'postgres-staging', image: 'postgres:15-alpine' }
            ],
            volumes: [{ name: 'postgres-data', type: 'persistentVolumeClaim' }],
            role: 'database',
            ip: '10.1.1.13',
            labels: { 'app': 'database', 'environment': 'staging' }
          }
        ]
      }
    ],
    namespaces: [
      {
        deployments: [
          {
            replicas: 2,
            selector: 'app=web-frontend,environment=staging',
            template: 'web-staging-template',
            strategy: 'RollingUpdate',
            role: 'frontend'
          },
          {
            replicas: 1,
            selector: 'app=api-server,environment=staging',
            template: 'api-staging-template',
            strategy: 'RollingUpdate',
            role: 'backend'
          },
          {
            replicas: 1,
            selector: 'app=database,environment=staging',
            template: 'postgres-staging-template',
            strategy: 'Recreate',
            role: 'data'
          }
        ],
        services: [
          {
            pods: [],
            clusterIP: '10.97.1.1',
            nodePort: 31001,
            selector: 'app=web-frontend,environment=staging',
            ports: {
              80: {
                containers: [{ name: 'nginx-staging', image: 'nginx:1.24-alpine' }],
                volumes: [],
                role: 'web',
                ip: '10.1.1.11',
                labels: { 'app': 'web-frontend', 'environment': 'staging' }
              }
            },
            role: 'frontend'
          }
        ]
      }
    ]
  },
  {
    id: 'dev-local-01',
    name: 'Development Local',
    controlPlane: {
      nodes: [
        {
          kubelet: 'v1.28.0',
          runtime: 'docker://24.0.2',
          role: 'control-plane',
          pods: [
            {
              containers: [
                { name: 'kube-apiserver', image: 'registry.k8s.io/kube-apiserver:v1.28.0' },
                { name: 'etcd', image: 'registry.k8s.io/etcd:3.5.8-0' }
              ],
              volumes: [{ name: 'etcd-data', type: 'emptyDir' }],
              role: 'control-plane',
              ip: '10.2.1.10',
              labels: { 'node-role.kubernetes.io/control-plane': '' }
            }
          ]
        }
      ]
    },
    nodes: [
      {
        kubelet: 'v1.28.0',
        runtime: 'docker://24.0.2',
        role: 'worker',
        pods: [
          {
            containers: [
              { name: 'dev-app', image: 'node:18-alpine' }
            ],
            volumes: [{ name: 'source-code', type: 'hostPath' }],
            role: 'development',
            ip: '10.2.1.11',
            labels: { 'app': 'dev-app', 'environment': 'development' }
          },
          {
            containers: [
              { name: 'postgres-dev', image: 'postgres:15-alpine' }
            ],
            volumes: [{ name: 'dev-data', type: 'emptyDir' }],
            role: 'database',
            ip: '10.2.1.12',
            labels: { 'app': 'database', 'environment': 'development' }
          }
        ]
      }
    ],
    namespaces: [
      {
        deployments: [
          {
            replicas: 1,
            selector: 'app=dev-app',
            template: 'dev-app-template',
            strategy: 'Recreate',
            role: 'development'
          },
          {
            replicas: 1,
            selector: 'app=database',
            template: 'postgres-dev-template',
            strategy: 'Recreate',
            role: 'data'
          }
        ],
        services: [
          {
            pods: [],
            clusterIP: '10.98.1.1',
            nodePort: 32001,
            selector: 'app=dev-app',
            ports: {
              3000: {
                containers: [{ name: 'dev-app', image: 'node:18-alpine' }],
                volumes: [],
                role: 'development',
                ip: '10.2.1.11',
                labels: { 'app': 'dev-app' }
              }
            },
            role: 'development'
          }
        ]
      }
    ]
  },
  {
    id: 'test-ci-01',
    name: 'Testing CI/CD',
    controlPlane: {
      nodes: [
        {
          kubelet: 'v1.27.5',
          runtime: 'containerd://1.6.20',
          role: 'control-plane',
          pods: [
            {
              containers: [
                { name: 'kube-apiserver', image: 'registry.k8s.io/kube-apiserver:v1.27.5' }
              ],
              volumes: [],
              role: 'control-plane',
              ip: '10.3.1.10',
              labels: { 'node-role.kubernetes.io/control-plane': '' }
            }
          ]
        }
      ]
    },
    nodes: [
      {
        kubelet: 'v1.27.5',
        runtime: 'containerd://1.6.20',
        role: 'worker',
        pods: [
          {
            containers: [
              { name: 'test-runner', image: 'node:18-alpine' }
            ],
            volumes: [{ name: 'test-results', type: 'emptyDir' }],
            role: 'testing',
            ip: '10.3.1.11',
            labels: { 'app': 'test-runner', 'type': 'ci' }
          },
          {
            containers: [
              { name: 'selenium', image: 'selenium/standalone-chrome:4.11' }
            ],
            volumes: [],
            role: 'testing',
            ip: '10.3.1.12',
            labels: { 'app': 'selenium', 'type': 'e2e' }
          }
        ]
      }
    ],
    namespaces: [
      {
        deployments: [
          {
            replicas: 1,
            selector: 'app=test-runner',
            template: 'test-runner-template',
            strategy: 'Recreate',
            role: 'testing'
          },
          {
            replicas: 0, // Scaled down for cost savings
            selector: 'app=selenium',
            template: 'selenium-template',
            strategy: 'Recreate',
            role: 'testing'
          }
        ],
        services: [
          {
            pods: [],
            clusterIP: '10.99.1.1',
            nodePort: 33001,
            selector: 'app=test-runner',
            ports: {
              8080: {
                containers: [{ name: 'test-runner', image: 'node:18-alpine' }],
                volumes: [],
                role: 'testing',
                ip: '10.3.1.11',
                labels: { 'app': 'test-runner' }
              }
            },
            role: 'testing'
          }
        ]
      }
    ]
  }
]