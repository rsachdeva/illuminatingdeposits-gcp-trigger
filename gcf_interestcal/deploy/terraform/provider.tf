terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "4.47.0"
    }
  }
}

provider "google" {
  # Configuration options
  project     = "illuminatingdeposits-gcp"
  region      = "us-central1"
  zone        = "us-central1-a"
  credentials = "keys.json"
}











