resource "google_project_service" "gke_apis" {
  project = "${google_project.my_project.project_id}"
  service = "container.googleapis.com"

  provisioner "local-exec" {
    command = "sleep 60;"
  }
}

resource "google_container_cluster" "primary" {
  project = "${google_project.my_project.project_id}"

  name     = "kubelet-image-service"
  location = "europe-west1"

  # We can't create a cluster with no node pool defined, but we want to only use
  # separately managed node pools. So we create the smallest possible default
  # node pool and immediately delete it.
  remove_default_node_pool = true
  initial_node_count       = 1

  master_auth {
    username = ""
    password = ""

    client_certificate_config {
      issue_client_certificate = false
    }
  }
}

resource "google_container_node_pool" "primary_preemptible_nodes" {
  project = "${google_project.my_project.project_id}"

  name       = "preemptible-node-pool"
  location   = "europe-west1"
  cluster    = google_container_cluster.primary.name
  node_count = 1

  node_config {
    preemptible  = true
    machine_type = "n1-standard-1"

    metadata = {
      disable-legacy-endpoints = "true"
    }

    oauth_scopes = [
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]
  }
}
