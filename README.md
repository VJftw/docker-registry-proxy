# Docker Registry Proxy

This project provides an authenticating (optional) proxy that can be used as a proxy hosted Docker Registries. Some use-cases are:

 - Simply provide users your own domain to pull Docker images from e.g. `docker pull docker.example.com/my-image`, is much nicer to use than `docker pull gcr.io/my-gcp-project/my-image`.
 - Provide custom authentication for pulling Docker images from your registries. At the moment, Google Container Registry (GCR) and Amazon Elastic Container Registry (ECR) only allow access to private registries via Service Accounts or Roles; this presents a key-management issue and ties client authentication to your cloud provider. In GCR, we also cannot restrict access per repository. The Docker Registry Proxy has a pluggable authentication mechanism, allowing you to implement arbitrary authentication flows. e.g. LDAP, OAuth2, static credentials etc.
- Provide access to a single source of private Docker images from multiple cloud providers via Instance Identity Documents. The **Kubelet Image Service** makes this possible with Kubernetes clusters.


## Kubelet Image Service

The Kubelet Image Service is designed to be the endpoint for the `--image-service-endpoint` flag in `kubelet` to directly intercept and transparently add arbitrary authentication when pulling Docker images in Kubernetes. This flag was merged in Aug 2016, so expected it to be available from K8s 1.5+. Using this, you no longer need to rely on `imagePullSecrets` which requires static credentials.
