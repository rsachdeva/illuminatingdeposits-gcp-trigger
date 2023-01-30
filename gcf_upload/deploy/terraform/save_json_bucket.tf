# storage bucket to save the uploaded http json body
resource "google_storage_bucket" "illuminating_save_json_bucket" {
  name = "illuminating_save_json_bucket"
  location = "US"
}