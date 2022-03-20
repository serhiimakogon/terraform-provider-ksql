terraform {
  required_providers {
    ksql = {
      source = "gabriel-aranha/ksql"
      version = "1.0.0-pre"
    }
  }
}

provider "ksql" {}

data "ksql_stream" "main" {
    name = "test-stream"
}
