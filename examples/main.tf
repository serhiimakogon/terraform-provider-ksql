terraform {
  required_providers {
    ksql = {
      source = "gabriel-aranha/ksql"
      version = "1.0.9-pre"
    }
  }
}

provider "ksql" {
  # Configuration options
}

data "ksql_stream" "main" {
  name = "test_stream"
}

output "stream_name" {
  value = data.ksql_stream.main.name
}
