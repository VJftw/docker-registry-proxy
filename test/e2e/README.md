# E2E Testing of the Docker Registry Proxy + Kubelet Image Service

## Docker Registry Proxy

1. :heavy_check_mark: Start a EKS cluster
2. :heavy_check_mark: Deploy Docker Registry Proxy deployment with ACM ingress backed by ECR
3. Test pulling the test image (from `terraform.tfvars`, no-auth required)

## Kubelet Image Service (Parallel)

### AWS EKS
$0.10/hour + nodes

1. Start an AWS EKS cluster
2. Deploy Kubelet Image Service daemonset
3. Start a Deployment that uses an image via Docker Registry Proxy

### GCP GKE
$0.10/hour + nodes

1. Start a GCP GKE cluster
2. Deploy Kubelet Image Service daemonset
3. Start a Deployment that uses an image via Docker Registry Proxy

 