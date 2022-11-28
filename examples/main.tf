terraform {
  required_providers {
    ksql = {
      source = "serhiimakogon/ksql"
      version = "1.0.0"
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
