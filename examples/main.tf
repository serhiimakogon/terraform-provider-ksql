terraform {
  required_providers {
    ksql = {
      source = "gabriel-aranha/ksql"
      version = "1.0.2-pre"
    }
  }
}

provider "ksql" {
  # Configuration options
}

data "ksql_stream" "main" {
  name = "test_stream"
}

resource "ksql_stream" "main" {
  name = "test_stream"
  query = "AS SELECT * FROM KSQL_PROCESSING_LOG;"
}
