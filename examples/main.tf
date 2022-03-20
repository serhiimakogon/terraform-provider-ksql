terraform {
  required_providers {
    ksql = {
      source = "gabriel-aranha/ksql"
    }
  }
}

provider "ksql" {}

data "ksql_stream" "main" {
    name = "test-stream"
}
