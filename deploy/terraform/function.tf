resource "google_storage_bucket" "illuminating_bucket" {
  name = "illuminating_gcp_trigger"
  location = "US"
}

resource "google_storage_bucket_object" "illuminating_src_code" {
  name = "gomain.zip"
  bucket = google_storage_bucket.illuminating_bucket.name
  source = "gomain.zip"
}

resource "google_cloudfunctions_function" "illuminating_deposits_func" {
  name = "illuminating_deposits_func"
  runtime = "go119"
  description = "function from terraform script using go 1.19"

  available_memory_mb = 128
  source_archive_bucket = google_storage_bucket.illuminating_bucket.name
  source_archive_object = google_storage_bucket_object.illuminating_src_code.name

  trigger_http = true
  entry_point = "HelloHTTP"

}