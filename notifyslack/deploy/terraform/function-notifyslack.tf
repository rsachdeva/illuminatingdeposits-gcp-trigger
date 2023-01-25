resource "google_storage_bucket" "notifyslack_bucket" {
  name = "notifyslack_trigger"
  location = "US"
}

resource "google_storage_bucket_object" "notifyslack_src_code" {
  name = "notifyslack-gosource.zip"
  bucket = google_storage_bucket.notifyslack_bucket.name
  source = "notifyslack-gosource.zip"
}


resource "google_cloudfunctions2_function" "function" {
  name = "notifyslack-vzeropoint1"
  location = "us-central1"
  description = "notifyslack function from terraform script using go 1.19"

  build_config {
    runtime = "go119"
    entry_point = "NotifySlack"  # Set the entry point for exported function
    source {
      storage_source {
        bucket = google_storage_bucket.notifyslack_bucket.name
        object = google_storage_bucket_object.notifyslack_src_code.name
      }
    }
  }

  lifecycle {
    replace_triggered_by  = [
      google_storage_bucket_object.notifyslack_src_code
    ]
  }

  service_config {
    max_instance_count  = 1
    available_memory    = "256M"
    timeout_seconds     = 60
    ingress_settings = "ALLOW_ALL"
    all_traffic_on_latest_revision = true
  }
}