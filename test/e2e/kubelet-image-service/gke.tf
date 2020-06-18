locals {
    region = "europe-west1"
    node_instance_type = "n1-standard1"
    billing_account_id = "00C32F-55AB5C-78A691"
    name = "kis"
    namespace = "test"
    stage = "e2e"
    tags = {}
}

resource "random_id" "id" {
  byte_length = 2
}

module "base_label" {
  source     = "git::https://github.com/cloudposse/terraform-null-label.git?ref=0.16.0"
  namespace  = local.namespace
  stage      = local.stage
  environment = random_id.id.hex
  name = local.name
  tags = local.tags
}

resource "google_project" "test_project" {
  name            = module.base_label.id
  project_id      = module.base_label.id
  billing_account = local.billing_account_id
}

resource "google_project_service" "gke_apis" {
  project = google_project.test_project.project_id
  service = "container.googleapis.com"

#   provisioner "local-exec" {
#     command = "sleep 60;"
#   }
}

resource "google_container_cluster" "primary" {
    depends_on = [
        google_project_service.gke_apis,
    ]
  project = google_project.test_project.project_id

  name     = module.base_label.id
  location = local.region

  # We can't create a cluster with no node pool defined, but we want to only use
  # separately managed node pools. So we create the smallest possible default
  # node pool and immediately delete it.
#   remove_default_node_pool = true
  initial_node_count       = 1

  master_auth {
    username = ""
    password = ""

    client_certificate_config {
      issue_client_certificate = false
    }
  }
}

# resource "google_container_node_pool" "primary_preemptible_nodes" {
#   project = google_project.test_project.project_id

#   name       = module.base_label.id
#   location   = local.region
#   cluster    = google_container_cluster.primary.name
#   node_count = 1

#   node_config {
#     preemptible  = true
#     machine_type = local.node_instance_type

#     metadata = {
#       disable-legacy-endpoints = "true"
#     }

#     oauth_scopes = [
#       "https://www.googleapis.com/auth/logging.write",
#       "https://www.googleapis.com/auth/monitoring",
#     ]
#   }
# }
