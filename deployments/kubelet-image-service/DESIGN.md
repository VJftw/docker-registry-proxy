# Kubelet Image Service Deployment design

As the Kubelete Image Service is installed onto hosts outside of Kubernetes, but via Kubernetes we propose the following deployment mechanism:

Privileged K8s (initContainer) Daemonset
- Go binary (unikernel) to:
  - copy `kubelet-image-service` binary to host from the container.
  - install `systemd` to start `kubelet-image-service`.

**Note**: `kubelet-image-service` should now configure kubelet for us (there are various configurations (GKE/EKS/Minikube etc.)).
  - reconfigure `kubelet` to use `kubelet-image-service` (via `--image-service-endpoint`).

Parameters:
 - `install_path`: where to copy the `kubelet-image-service` binary to.
 - ``