# storage bucket to save the uploaded http json body
resource "google_storage_bucket" "illuminating_upload_json_bucket" {
  name = "illuminating_upload_json_bucket"
  location = "US"
  force_destroy = true
}