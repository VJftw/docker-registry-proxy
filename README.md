# Docker Registry Proxy

This project provides a proxy (optionally authenticating) to a Docker Registry. Usually we desire managed private Docker Registries such as Google Container Registry (GCR), Amazon Elastic Container Registry (ECR), Private Docker Hub, Quay.io, etc. as they are much cheaper to set-up, manage and reliably scale.

Some use-cases of this Docker Registry Proxy, with this in mind, are:

 - Provide users your own domain to pull Docker images from e.g. `docker pull docker.example.com/my-image`, is much nicer to use than `docker pull gcr.io/my-gcp-project/my-image`.
 - Provide arbitrary authentication for pulling Docker images from your registries. e.g. Hosted registries such as GCR and ECR only allow access to private registries via Service Accounts, Roles, or static credentials; this presents a key-management issue and/or ties client authentication to your cloud provider. In GCR, we also cannot restrict access per repository. The Docker Registry Proxy has a pluggable authentication mechanism, allowing you to implement arbitrary authentication flows. e.g. LDAP, OAuth2, static credentials etc.
- Provide access to a single source of private Docker images from multiple cloud providers via Instance Identity Documents. The **Kubelet Image Service** makes this possible with Kubernetes clusters.


## Kubelet Image Service

The Kubelet Image Service is designed to be the endpoint for the `--image-service-endpoint` flag in `kubelet` to directly intercept and transparently add arbitrary authentication when pulling Docker images in Kubernetes. This flag was merged in Aug 2016, so expected it to be available from K8s 1.5+. Using this, you no longer need to rely on `imagePullSecrets` which requires static credentials.

This is deployed in a container via a `DaemonSet`, thus requires access to the `hostPath: /var/run/docker.sock` to interact with the Docker Engine on the node.

### Managed Kubernetes
In most managed kubernetes offerings, it is difficult to modify the `kubelet` flags. The GKE documentation [recommends making host image modifications via a DaemonSet](https://cloud.google.com/kubernetes-engine/docs/concepts/node-images#modifications). This section lists the `DaemonSet` workarounds used for each Cloud Service Provider (CSP).

#### Google Kubernetes Engine (GKE)
The DaemonSet updates `/etc/default/kubelet` on the host (using `hostPath` mounts) and then kills the `kubelet` process on the host (via `hostPID: true`) which is then restarted by the host SystemD with the new kubelet flags.
