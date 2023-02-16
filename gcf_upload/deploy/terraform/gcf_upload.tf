// buckets for source code illuminating_gcf_upload_src.zip
resource "google_storage_bucket" "illuminating_gcf_upload_bucket" {
  name = "illuminating_gcf_upload_bucket"
  location = "us-central1"
}

resource "google_storage_bucket_object" "illuminating_gcf_upload_src_code" {
  name = "illuminating_gcf_src_code"
  bucket = google_storage_bucket.illuminating_gcf_upload_bucket.name
  source = "illuminating_gcf_upload_src.zip"
}


resource "google_cloudfunctions2_function" "illuminating_gcf_upload" {
  name = "illuminating-gcf-upload"
  location = "us-central1"
  description = "illuminating gcf that takes the interest cal actual body of json through http post and stores as a file in cloud storage from terraform script"

  build_config {
    runtime = "go119"
    entry_point = "UploadHTTP"  # Set the entry point for exported function
    source {
      storage_source {
        # gcf-v2-sources-923961404233-us-central1 created bucket with a file function-source.zip. This is automatically created from
        # illuminating_gcp_trigger bucket with the uploaded file from our terraform block above illuminating-gosource.zip
        # manually clean this resource if needed to be sure when doing terraform destroy reference:
        # https://stackoverflow.com/questions/72148179/after-delete-a-cloud-function-it-still-in-gcf-sources
        bucket = google_storage_bucket.illuminating_gcf_upload_bucket.name
        object = google_storage_bucket_object.illuminating_gcf_upload_src_code.name
      }
    }
  }

  lifecycle {
    replace_triggered_by  = [
      google_storage_bucket_object.illuminating_gcf_upload_src_code
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