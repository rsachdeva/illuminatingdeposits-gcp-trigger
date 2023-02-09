# storage bucket to save the uploaded http json body
resource "google_storage_bucket" "illuminating_upload_json_bucket_output" {
  name = "illuminating_upload_json_bucket_output"
  location = "us-central1"
  force_destroy = true
}