---
subcategory: "KSQL Streams"
page_title: "KSQL: ksql_streams"
description: |-
Provides a KSQL Streams data source.
---

# Data Source: ksql_streams

Use this data source to get information about KSQL Streams for use in other resources.

## Example Usage

### List all KSQL Streams
```terraform
data "ksql_streams" "stream" {}
```

### List all KSQL Streams by tag
```terraform
data "ksql_streams" "stream" {
  tag = "DEV"
}
```

### List all KSQL Streams by topic
```terraform
data "ksql_streams" "stream" {
  topic = "main_topic"
}
```

## Argument Reference
Currently this data source does not allow setting both `tag` and `topic` variables at the same time.

* `tag` - (Optional) The tag to filter the streams. Case sensitive.
* `topic` - (Optional) The topic to filter the streams. Case sensitive.

## Attributes Reference

`id` is set to a random `UUID`. In addition, the following attributes
are exported:

* `streams` - The list of streams found. Can be empty if no streams are found.
  * `name` - The name of the KSQL Stream.
  * `topic` - The topic backing the stream.
