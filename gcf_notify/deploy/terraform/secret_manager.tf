resource "google_secret_manager_secret" "secret-basic" {
  secret_id = "sendgrid-api-key"

  labels = {
    label = "sendgrid"
  }

  replication {
    automatic = true
  }
}


resource "google_secret_manager_secret_version" "secret-version-basic" {
  secret = google_secret_manager_secret.secret-basic.id
   // read the secret value from a file that is ignored by git
  secret_data = file("sendgrid_api_value.txt")
}