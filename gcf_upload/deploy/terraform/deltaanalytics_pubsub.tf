resource "google_pubsub_schema" "delta_analytics_schema" {
  name = "deltaanalyticsschema"
  type = "PROTOCOL_BUFFER"
  definition = "${file("delta_analytics_schema.proto")}"
}

// used as a notification for the Simple delta analytics and also notify if any error occurs in the process
resource "google_pubsub_topic" "delta_analytics_topic" {
  name = "deltaanalyticstopic"
  depends_on = [google_pubsub_schema.delta_analytics_schema]
  schema_settings {
    schema = "projects/illuminatingdeposits-gcp/schemas/deltaanalyticsschema"
    encoding = "JSON"
  }
  message_retention_duration = "600s"
}

resource "google_pubsub_subscription" "delta_analytics_test_subscription" {
  name = "deltaanalyticstestsubscription"
  topic = google_pubsub_topic.delta_analytics_topic.name
}