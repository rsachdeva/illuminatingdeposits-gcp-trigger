// buckets for source code illuminating_gcf_interestcal_src.zip
resource "google_storage_bucket" "illuminating_gcf_interestcal_bucket" {
  name = "illuminating_gcf_interestcal_bucket"
  location = "us-central1"
}

resource "google_storage_bucket_object" "illuminating_gcf_interestcal_src_code" {
  name = "illuminating_gcf_src_code"
  bucket = google_storage_bucket.illuminating_gcf_interestcal_bucket.name
  source = "illuminating_gcf_interestcal_src.zip"
}

resource "google_cloudfunctions2_function" "illuminating_gcf_interestcal" {
  name = "illuminating-gcf-interestcal"
  location = "us-central1"
  description = "gcf that that gets triggered by file in cloud storage illuminating_upload_json_bucket_input trigger bucket and make interest calculation for that data and stores in another illuminating_upload_json_bucket_output bucket"

  build_config {
    runtime = "go120"
    entry_point = "InterestCalStorage"  # Set the entry point for exported function
    source {
      storage_source {
        # gcf-v2-sources-923961404233-us-central1 created bucket with a file function-source.zip. This is automatically created from
        # illuminating_gcp_trigger bucket with the uploaded file from our terraform block above illuminating-gosource.zip
        # manually clean this resource if needed to be sure when doing terraform destroy reference:
        # https://stackoverflow.com/questions/72148179/after-delete-a-cloud-function-it-still-in-gcf-sources
        bucket = google_storage_bucket.illuminating_gcf_interestcal_bucket.name
        object = google_storage_bucket_object.illuminating_gcf_interestcal_src_code.name
      }
    }
  }

  lifecycle {
    replace_triggered_by  = [
      google_storage_bucket_object.illuminating_gcf_interestcal_src_code
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

  event_trigger {
    trigger_region        = "us-central1" # The trigger must be in the same location as the bucket
    event_type            = "google.cloud.storage.object.v1.finalized"
    event_filters {
      attribute = "bucket"
      value     = data.google_storage_bucket.trigger_bucket.name
    }
  }
}
