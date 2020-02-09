resource "random_id" "project_id" {
  byte_length = 4
}

resource "google_project" "my_project" {
  name            = "test-kis"
  project_id      = "test-kis-${random_id.project_id.hex}"
  billing_account = "${var.billing_account_id}"
}
