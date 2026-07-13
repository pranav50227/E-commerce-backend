provider "google" {
  project = "ecommerce-project"
  region  = "us-central1"
}

# VPC Network
resource "google_compute_network" "vpc_network" {
  name = "ecommerce-network"
}

# VM instance for App deployment
resource "google_compute_instance" "app_instance" {
  name         = "ecommerce-app-instance"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }

  network_interface {
    network = google_compute_network.vpc_network.name
    access_config {
      // Ephemeral public IP
    }
  }
}
