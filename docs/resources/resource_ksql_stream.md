---
page_title: "KSQL: ksql_stream"
subcategory: "KSQL Streams"
---

# Resource: ksql_stream

Provides a KSQL Stream resource.

## Example Usage

```terraform
resource "ksql_stream" "stream" {
  name  = "STREAM_01"
  query = "AS SELECT * FROM STREAM_00;"
}
```

## Argument Reference

* `name` - (Required) The name of the stream. Case insensitive. Any changes to the name forces the creation of a new resource.
* `query` - (Required) The statement to create the stream. Any changes to the query forces the creation of a new resource.

## Attributes Reference

`id` is set to the name of the KSQL Stream. In addition, the following attributes
are exported:

* `name` - The name of the stream.
* `query` - The statement that created the stream.
