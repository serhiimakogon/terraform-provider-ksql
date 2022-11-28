terraform {
  required_providers {
    ksql = {
      source = "serhiimakogon/ksql"
      version = "1.0.0"
    }
  }
}

data "ksql_stream" "main" {
  name = "test_stream"
}

output "stream_name" {
  value = data.ksql_stream.main.name
}
