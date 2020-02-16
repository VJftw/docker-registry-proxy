
resource "google_project_service" "gcr_apis" {
  project = google_project.my_project.project_id
  service = "containerregistry.googleapis.com"

  provisioner "local-exec" {
    command = "docker pull ghost:latest && docker tag ghost:latest gcr.io/${google_project.my_project.project_id}/ghost:latest && docker push gcr.io/${google_project.my_project.project_id}/ghost:latest"
  }
}

resource "google_service_account" "service_account" {
  project = google_project.my_project.project_id

  account_id   = "docker-registry-proxy"
  display_name = "Docker Registry Proxy"
}

resource "google_service_account_key" "mykey" {
  service_account_id = google_service_account.service_account.name
  public_key_type    = "TYPE_X509_PEM_FILE"
}

resource "google_storage_bucket_iam_member" "puller" {

  bucket = "artifacts.${google_project.my_project.project_id}.appspot.com"
  role   = "roles/storage.objectViewer"
  member = "serviceAccount:${google_service_account.service_account.email}"

  depends_on = ["google_project_service.gcr_apis"]
}
