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

resource "google_service_account" "service_account" {
  account_id   = "uptime-monitor-runtime"
  display_name = "Uptime Monitor Cloud Run Runtime"
}

resource "google_project_iam_member" "runtime_cloudsql_client" {
  project = "englander-suite"
  role    = "roles/cloudsql.client"
  member  = "serviceAccount:${google_service_account.service_account.email}"
}

resource "google_project_iam_member" "runtime_cloudsql_instance_user" {
  project = "englander-suite"
  role    = "roles/cloudsql.instanceUser"
  member  = "serviceAccount:${google_service_account.service_account.email}"
}

resource "google_artifact_registry_repository_iam_member" "member" {
  project    = google_artifact_registry_repository.docker-artifact-repository.project
  location   = google_artifact_registry_repository.docker-artifact-repository.location
  repository = google_artifact_registry_repository.docker-artifact-repository.id
  role       = "roles/artifactregistry.writer"
  member     = "serviceAccount:755712906263-compute@developer.gserviceaccount.com"
}

resource "google_artifact_registry_repository" "docker-artifact-repository" {
  location      = "us-east4"
  repository_id = "docker-artifact-repository"
  description   = "Docker Artifact Repository"
  format        = "DOCKER"

}

variable "cloud_build_image" {
  type        = string
  description = "Docker image passed from cloudbuild"
}

resource "google_cloud_run_service" "default" {
  name     = "uptime-monitor-gcr"
  location = "us-east4"

  template {
    metadata {
      annotations = {
        "run.googleapis.com/cloudsql-instances" = google_sql_database_instance.instance.connection_name
      }
    }
    spec {
      service_account_name = google_service_account.service_account.email
      containers {
        image = var.cloud_build_image
        env {
          name  = "ICN_STRING"
          value = "englander-suite:us-east4:uptime-database-instance"
        }
        env {
          name  = "DATABASE_SERVICE_ACCOUNT"
          value = "user=uptime-monitor-runtime@englander-suite.iam dbname=uptime-database sslmode=disable"
        }
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

resource "google_cloudbuild_trigger" "tf-trigger" {
  name            = "uptime-monitor-main"
  location        = "global"
  service_account = "projects/englander-suite/serviceAccounts/755712906263-compute@developer.gserviceaccount.com"
  github {
    owner = "karstenenglander"
    name  = "uptime-monitor"
    push {
      branch = "^main$"
    }
  }
  filename = "cloudbuild.yaml"
}

resource "google_cloud_scheduler_job" "job" {
  paused           = true
  name             = "uptime-monitor-1m"
  description      = "Runs uptime-monitor once every 1 minute"
  schedule         = "* * * * *"
  time_zone        = "America/New_York"
  attempt_deadline = "320s"

  retry_config {
    retry_count = 1
  }

  http_target {
    http_method = "POST"
    uri         = "https://uptime-monitor-gcr-755712906263.us-east4.run.app"
    body        = base64encode("{\"foo\":\"bar\"}")
    headers = {
      "Content-Type" = "application/json"
    }
  }
}

resource "google_sql_database" "database" {
  name     = "uptime-database"
  instance = google_sql_database_instance.instance.name
}

# See versions at https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/sql_database_instance#database_version

resource "google_sql_database_instance" "instance" {
  name             = "uptime-database-instance"
  region           = "us-east4"
  database_version = "POSTGRES_18"
  settings {
    tier = "db-f1-micro"

    backup_configuration {
      enabled                        = true
      point_in_time_recovery_enabled = true
    }

    database_flags {
      name  = "cloudsql.iam_authentication"
      value = "on"
    }
  }

  deletion_protection = true
}

resource "google_sql_user" "iam_user" {
  name     = "karstenenglander@gmail.com"
  instance = google_sql_database_instance.instance.name
  type     = "CLOUD_IAM_USER"
}

resource "google_sql_user" "iam_service_account_user" {
  # Note: for Postgres only, GCP requires omitting the ".gserviceaccount.com" suffix
  # from the service account email due to length limits on database usernames.

  name     = trimsuffix(google_service_account.service_account.email, ".gserviceaccount.com")
  instance = google_sql_database_instance.instance.name
  type     = "CLOUD_IAM_SERVICE_ACCOUNT"
}

resource "google_cloud_tasks_queue" "uptime_queue" {
  name     = "uptime-queue"
  location = "us-east4"

  rate_limits {
    max_concurrent_dispatches = 3
    max_dispatches_per_second = 2
  }

  retry_config {
    max_attempts       = 5
    max_retry_duration = "300s"
    max_backoff        = "60s"
    min_backoff        = "5s"
    max_doublings      = 3
  }

  stackdriver_logging_config {
    sampling_ratio = 1.0
  }
}
