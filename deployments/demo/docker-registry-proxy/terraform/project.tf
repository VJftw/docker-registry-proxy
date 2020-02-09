resource "random_id" "project_id" {
  byte_length = 4
}

resource "google_project" "my_project" {
  name            = "test-docker-proxy"
  project_id      = "test-docker-proxy-${random_id.project_id.hex}"
  billing_account = "${var.billing_account_id}"
}
