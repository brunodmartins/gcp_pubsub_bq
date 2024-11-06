provider "google" {
  project = "zoo-poc"
  region  = "us-central1"
}

data "google_project" "default" {}

resource "google_bigquery_dataset" "zoo_dataset" {
  dataset_id                  = "zoo"
  friendly_name               = "Zoo Dataset"
  description                 = "This dataset contains a table for zoo tables"
  location                    = "US"
  default_table_expiration_ms = 3600000

  labels = {
    env = "default"
  }
}

resource "google_bigquery_table" "animals_table" {
  dataset_id = google_bigquery_dataset.zoo_dataset.dataset_id
  table_id   = "animals"

  deletion_protection=false

  time_partitioning {
    type = "DAY"
  }

  labels = {
    env = "default"
  }
  schema = <<EOF
[
  {
    "name": "name",
    "type": "STRING",
    "mode": "REQUIRED",
    "description": "Name of the animal"
  },
  {
    "name": "age",
    "type": "INTEGER",
    "mode": "NULLABLE",
    "description": "Age of the animal in years"
  },
  {
    "name": "species",
    "type": "STRING",
    "mode": "REQUIRED",
    "description": "Species of the animal"
  },
  {
    "name": "family",
    "type": "STRING",
    "mode": "NULLABLE",
    "description": "Family to which the animal belongs"
  },
  {
      "name": "gender",
      "type": "STRING",
      "mode": "NULLABLE",
      "description": "Gender of the animal"
    },
    {
      "name": "born_date",
      "type": "DATE",
      "mode": "NULLABLE",
      "description": "Birth date of the animal"
    }
  ]
  EOF

}

resource "google_project_iam_member" "pubsub_member_iam" {
    project = data.google_project.default.project_id
    role    = "roles/bigquery.dataEditor"
    member  = "serviceAccount:service-${data.google_project.default.number}@gcp-sa-pubsub.iam.gserviceaccount.com"
  }

resource "google_pubsub_topic" "animals_topic" {
  name = "animals-topic"
}

resource "google_pubsub_subscription" "animals_subscription" {
  name  = "animals-subscription"
  topic = google_pubsub_topic.animals_topic.id

  bigquery_config {
    table = "${google_bigquery_table.animals_table.project}.${google_bigquery_table.animals_table.dataset_id}.${google_bigquery_table.animals_table.table_id}"
    use_table_schema = true
  }

  retain_acked_messages = true
  message_retention_duration = "86400s"  # Retain for 1 day
}
