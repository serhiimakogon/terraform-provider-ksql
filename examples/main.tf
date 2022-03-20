terraform {
  required_providers {
    ksql = {
      source = "gabriel-aranha/ksql"
      version = "1.0.3-pre"
    }
  }
}

provider "ksql" {
  # Configuration options
}
