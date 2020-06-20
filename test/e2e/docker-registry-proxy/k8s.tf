data "aws_eks_cluster_auth" "cluster" {
  name = module.eks_cluster.eks_cluster_id
}

provider "kubernetes" {
  host                   = module.eks_cluster.eks_cluster_endpoint
  cluster_ca_certificate = base64decode(module.eks_cluster.eks_cluster_certificate_authority_data)
  token                  = data.aws_eks_cluster_auth.cluster.token
  load_config_file       = false
}

resource "kubernetes_namespace" "dockerregistryproxy" {
  metadata {
    name = "docker-registry-proxy"
  }
}

resource "kubernetes_deployment" "dockerregistryproxy" {
  metadata {
    name = "docker-registry-proxy"
    namespace = kubernetes_namespace.dockerregistryproxy.metadata.0.name
  }

  spec {
    replicas = 3

    selector {
      match_labels = {
        app = "docker-registry-proxy"
      }
    }

    template {
      metadata {
        labels = {
          app = "docker-registry-proxy"
        }
      }

      spec {
        container {
          name = "docker-registry-proxy"

          image             = "vjftw/docker-registry-proxy:latest"
          args              = ["--network_address=tcp://:8888"]
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
          # env {
          #   name  = "DRP_AUTHENTICATION_VERIFIER"
          #   value = "_gcpidd:gcpverifier _awsidd:awsverifier"
          # }
          # env {
          #   name  = "DRP_GCPVERIFIER_PROJECT_IDS"
          #   value = ""
          # }
          # env {
          #   name  = "DRP_AWSVERIFIER_ACCOUNT_IDS"
          #   value = ""
          # }

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

resource "kubernetes_config_map" "dockerregistryproxy" {
  metadata {
    name = "docker-registry-proxy"
    namespace = kubernetes_namespace.dockerregistryproxy.metadata.0.name
  }

  data = {

  }
}

resource "kubernetes_service" "dockerregistryproxy" {
  metadata {
    name = "docker-registry-proxy"
    namespace = kubernetes_namespace.dockerregistryproxy.metadata.0.name
    annotations = {
      "service.beta.kubernetes.io/aws-load-balancer-backend-protocol" = "http"
      "service.beta.kubernetes.io/aws-load-balancer-ssl-cert"         = module.acm_request_certificate.arn
      "service.beta.kubernetes.io/aws-load-balancer-ssl-ports"        = "https"
    }
  }
  spec {
    selector = {
      app = kubernetes_deployment.dockerregistryproxy.spec.0.template.0.metadata.0.labels.app
    }
    # session_affinity = "ClientIP"

    port {
      name        = "https"
      port        = 443
      target_port = 8888
    }

    type = "LoadBalancer"
  }
}

module "acm_request_certificate" {
  source                            = "git::https://github.com/cloudposse/terraform-aws-acm-request-certificate.git?ref=0.4.0"
  domain_name                       = "e2e.test.dockerregistryproxy.vjpatel.me"
  zone_name = "vjpatel.me"
  process_domain_validation_options = true
  ttl                               = "300"
}

resource "aws_route53_record" "dockerregistryproxy" {
  zone_id = data.aws_route53_zone.main.zone_id
  name    = "e2e.test.dockerregistryproxy.vjpatel.me"
  type    = "CNAME"
  ttl     = "300"
  records = [kubernetes_service.dockerregistryproxy.load_balancer_ingress.0.hostname]
}
