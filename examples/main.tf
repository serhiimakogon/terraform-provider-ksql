terraform {
  required_providers {
    ksql = {
      source = "gabriel-aranha/ksql"
      version = "1.0.4-pre"
    }
  }
}

provider "ksql" {
  # Configuration options
}
