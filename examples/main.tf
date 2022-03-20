terraform {
  required_providers {
    ksql = {
      source = "gabriel-aranha/ksql"
      version = "1.0.5-pre"
    }
  }
}

provider "ksql" {
  # Configuration options
}

data "ksql_stream" "main" {
  name = "TEST_STREAM"
}

data "ksql_streams" "main" {
  tag = "TEST"
}

resource "ksql_stream" "main" {
  name = "ANOTHER_STREAM"
  query = "AS SELECT * FROM TEST_STREAM;"
}
