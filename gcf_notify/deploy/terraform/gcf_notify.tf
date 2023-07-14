// buckets for source code illuminating_gcf_notify_src.zip
resource "google_storage_bucket" "illuminating_gcf_notify_bucket" {
  name = "illuminating_gcf_notify_bucket"
  location = "us-central1"
}

resource "google_storage_bucket_object" "illuminating_gcf_notify_src_code" {
  name = "illuminating_gcf_src_code"
  bucket = google_storage_bucket.illuminating_gcf_notify_bucket.name
  source = "illuminating_gcf_notify_src.zip"
}

resource "google_cloudfunctions2_function" "illuminating_gcf_notify" {
  name = "illuminating-gcf-notify"
  location = "us-central1"
  description = "gcf that that gets triggered by file in cloud storage illuminating_upload_json_bucket_input trigger bucket and make interest calculation for that data and stores in another illuminating_upload_json_bucket_output bucket"

  build_config {
    runtime = "go120"
    entry_point = "NotifyInvestorOfDelta"  # Set the entry point for exported function
    source {
      storage_source {
        # gcf-v2-sources-923961404233-us-central1 created bucket with a file function-source.zip. This is automatically created from
        # illuminating_gcp_trigger bucket with the uploaded file from our terraform block above illuminating-gosource.zip
        # manually clean this resource if needed to be sure when doing terraform destroy reference:
        # https://stackoverflow.com/questions/72148179/after-delete-a-cloud-function-it-still-in-gcf-sources
        bucket = google_storage_bucket.illuminating_gcf_notify_bucket.name
        object = google_storage_bucket_object.illuminating_gcf_notify_src_code.name
      }
    }
  }

  lifecycle {
    replace_triggered_by  = [
      google_storage_bucket_object.illuminating_gcf_notify_src_code
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
    environment_variables = {
      SENDER_EMAIL = file("sender_email.txt")
      RECEIVER_EMAIL = file("receiver_email.txt")
    }
  }

  event_trigger  {
    trigger_region = "us-central1"
    event_type = "google.cloud.pubsub.topic.v1.messagePublished"
    pubsub_topic = data.google_pubsub_topic.delta_analytics_topic.id
    retry_policy = "RETRY_POLICY_RETRY"
  }
}