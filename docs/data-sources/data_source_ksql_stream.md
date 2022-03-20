---
page_title: "KSQL: ksql_stream"
subcategory: "KSQL Streams"
---

# Data Source: ksql_stream

Use this data source to get information about a KSQL Stream for use in other resources.

## Example Usage

```terraform
data "ksql_stream" "stream" {
  name = "stream_name"
}
```

## Argument Reference

* `name` - (Required) The name of the KSQL Stream. Case insensitive.

## Attributes Reference

`id` is set to the name of the KSQL Stream. In addition, the following attributes
are exported:

* `name` - The name of the KSQL Stream.
* `topic` - The topic backing the stream.
