---
page_title: "KSQL Provider"
subcategory: ""
---

# KSQL Provider

Use the KSQL provider to interact with the resources supported. You must configure the provider with the proper credentials before you can use it.

Use the navigation to the left to read about the available resources.

## Example Usage

```terraform
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
```

## Authentication

Authentication can be provided by environment variables or provider configuration block variables.

### Environment Variables

Credentials can be provided by using the `KSQLDB_URL`, `KSQLDB_USERNAME`, and `KSQLDB_PASSWORD` environment variables.

For example:

```terraform
provider "ksql" {}
```

```
$ export KSQLDB_URL="yourksqldburl"
$ export KSQLDB_USERNAME="yourksqldbusername"
$ export KSQLDB_PASSWORD="yourksqldbpassword"
```

### Provider Block Configuration

Although not recommended, it is possible to set the credential variables in the provider configuration block as follows:

```terraform
provider "ksql" {
  url      = "yourksqldburl"
  username = "yourksqldbusername"
  password = "yourksqldbpassword"
}
```
