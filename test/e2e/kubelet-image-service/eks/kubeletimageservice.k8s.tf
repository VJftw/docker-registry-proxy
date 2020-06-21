data "aws_eks_cluster_auth" "cluster" {
  name = module.eks_cluster.eks_cluster_id
}

provider "kubernetes" {
  host                   = module.eks_cluster.eks_cluster_endpoint
  cluster_ca_certificate = base64decode(module.eks_cluster.eks_cluster_certificate_authority_data)
  token                  = data.aws_eks_cluster_auth.cluster.token
  load_config_file       = false
}

resource "kubernetes_namespace" "kubeletimageservice" {
  metadata {
    name = "kubelet-image-service"
  }
}

resource "kubernetes_daemonset" "dockerregistryproxy" {
  metadata {
    name = "kubelet-image-service"
    namespace = kubernetes_namespace.dockerregistryproxy.metadata.0.name
  }

  spec {
    replicas = 3

    selector {
      match_labels = {
        app = "kubelet-image-service"
      }
    }

    template {
      metadata {
        labels = {
          app = "kubelet-image-service"
        }
      }

      spec {
        container {
          name = "kubelet-image-service"

          image             = "vjftw/kubelet-image-service:latest"
          args              = ["--network_address=unix:///var/run/kubelet/kubelet-image-service.sock"]
          image_pull_policy = "Always"

          port {
            container_port = 8888
          }

          env {
            name  = "DRP_UPSTREAM_REPOSITORY"
            value = data.aws_ecr_authorization_token.test.proxy_endpoint
          }
          env {
            name  = "DRP_PLUGINS"
            value = "auth-provider_static:upstream_static auth-verifier_gcp-instanceidentitydocument:gcpverifier auth-verifier_aws-instanceidentitydocument:awsverifier"
          }
          env {
            name  = "DRP_UPSTREAM_AUTHENTICATION"
            value = "upstream_static"
          }
          env {
            name  = "DRP_UPSTREAM_STATIC_USERNAME"
            value = data.aws_ecr_authorization_token.test.user_name
          }
          env {
            name  = "DRP_UPSTREAM_STATIC_PASSWORD"
            value = data.aws_ecr_authorization_token.test.password
          }
          env {
            name  = "DRP_AUTHENTICATION_VERIFIER"
            value = var.docker_registry_proxy_authentication_verifier
          }
          env {
            name  = "DRP_GCPVERIFIER_PROJECT_IDS"
            value = var.docker_registry_proxy_gcpverifier_project_ids
          }
          env {
            name  = "DRP_AWSVERIFIER_ACCOUNT_IDS"
            value = var.docker_registry_proxy_awsverifier_account_ids
          }

          resources {
            limits {
              cpu    = "0.5"
              memory = "512Mi"
            }
            requests {
              cpu    = "250m"
              memory = "50Mi"
            }
          }
        }
      }
    }
  }
}