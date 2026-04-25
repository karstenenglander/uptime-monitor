terraform {
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

