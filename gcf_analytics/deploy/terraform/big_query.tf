// terraform big query dataset and table setup
resource "google_bigquery_dataset" "bd_ds" {
  dataset_id          = "gcfdeltaanalytics"
}

resource "google_bigquery_table" "table_tf" {
  table_id            = "delta_calculations"
  dataset_id          = google_bigquery_dataset.bd_ds.dataset_id
  # Whether or not to allow Terraform to destroy the instance; for now allowing to destroy
  deletion_protection = false
  #  Specify nested and repeated columns in table schemas https://cloud.google.com/bigquery/docs/nested-repeated
  schema              = "${file("delta_calculations_schema.json")}"
}
