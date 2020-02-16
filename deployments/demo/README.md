# Demo

## Architecture

This directory contains a demonstration with 2 components:
 - docker-registry-proxy
 - kubelet-image-service

### Docker Registry Proxy

This is where you pull your Docker images from.

This component creates a GCP Project `demo-docker-registry-proxy-*` to demonstrate a simple deployment of the Docker Registry Proxy onto a Kubernetes Cluster (GKE) backed by Google Container Registry (GCR). It creates an ingress so that the Docker Registry Proxy can be exposed.

See `up.sh`.

### Kubelet Image Service

This is how to configure your Kubernetes workers to authenticate and pull images from the Docker Registry Proxy.

This component creates GKE and EKS clusters to demonstrate a simple deployment of the Kubelet Image Service and updating of Kubelet flags on Google Kubernetes Engine (GKE). There is also a `ghost-k8s` directory to deploy a Ghost CMS that references the Docker Registry Proxy.

See `up.sh`.
