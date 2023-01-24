resource "google_storage_bucket" "illuminating_bucket" {
  name = "illuminating_gcp_trigger"
  location = "US"
}

resource "google_storage_bucket_object" "illuminating_src_code" {
  name = "illuminating-gosource.zip"
  bucket = google_storage_bucket.illuminating_bucket.name
  source = "illuminating-gosource.zip"
}

#resource "google_cloudfunctions_function" "illuminating_deposits_func" {
#  name = "illuminating_deposits_func"
#  runtime = "go119"
#  description = "function from terraform script using go 1.19"
#
#  available_memory_mb = 128
#  source_archive_bucket = google_storage_bucket.illuminating_bucket.name
#  source_archive_object = google_storage_bucket_object.illuminating_src_code.name
#
#  trigger_http = true
#  entry_point = "HelloHTTP"
#
#}


resource "google_cloudfunctions2_function" "function" {
  name = "illuminating-deposits-vzeropoint1"
  location = "us-central1"
  description = "function from terraform script using go 1.19"

  build_config {
    runtime = "go119"
    entry_point = "HelloHTTP"  # Set the entry point for exported function
    source {
      storage_source {
        bucket = google_storage_bucket.illuminating_bucket.name
        object = google_storage_bucket_object.illuminating_src_code.name
      }
    }
  }

  service_config {
    max_instance_count  = 1
    available_memory    = "256M"
    timeout_seconds     = 60
  }
}