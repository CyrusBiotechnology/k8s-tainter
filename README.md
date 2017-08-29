# k8s-tainter

From this discussion: https://groups.google.com/forum/#!topic/kubernetes-users/KUm233PUp-I

Taint Kubernetes nodes automatically from a configuration.

# Config

    taints:
      - label: cloud.google.com/gke-preemptible
        key: preemptible
      - label: high-idle
        key: high-priority
