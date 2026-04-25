terraform {
  backend "gcs" {
    bucket = "englander-suite-tfstate"
    prefix = "uptime-monitor"
  }
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "6.8.0"
    }
  }
}

provider "google" {
  project = "englander-suite"
  region  = "us-east4"
  zone    = "us-east4-a"
}

resource "google_compute_network" "vpc_network" {
  name = "terraform-network"
}

resource "google_cloud_run_service" "default" {
  name = "uptime-monitor-gcr"

  template {
    spec {
      containers {
        image = "us-docker.pkg.dev/cloudrun/container/uptime-monitor"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

resource "google_cloudbuild_trigger" "tf-trigger" {
  github {
    owner = "karstenenglander"
    name  = "uptime-monitor"
    push { branch = "^main$" }
  }
  filename = "cloudbuild.yaml"
}

