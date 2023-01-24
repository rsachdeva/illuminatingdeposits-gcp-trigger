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
        # gcf-v2-sources-923961404233-us-central1	bucket with a file function-source.zip is created from
        # illuminating_gcp_trigger bucket with the uploaded file from our terraform block above illuminating-gosource.zip
        bucket = google_storage_bucket.illuminating_bucket.name
        object = google_storage_bucket_object.illuminating_src_code.name
      }
    }
  }

  lifecycle {
    replace_triggered_by  = [
      google_storage_bucket_object.illuminating_src_code
    ]
  }

  service_config {
    max_instance_count  = 1
    available_memory    = "256M"
    timeout_seconds     = 60
    # ingress_settings - (Optional) Available ingress settings. Defaults to "ALLOW_ALL" if unspecified.
    ingress_settings = "ALLOW_ALL"
    # Whether 100% of traffic is routed to the latest revision. Defaults to true.
    all_traffic_on_latest_revision = true
  }
}