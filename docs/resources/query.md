---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "ksql_query Resource - terraform-provider-ksql"
subcategory: ""
description: |-
  Provides a KSQL Query resource.
---

# ksql_query (Resource)

Provides a KSQL Query resource.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) KSQL query name.
- `query` (String) KSQL query.
- `type` (String) KSQL query type [table|stream].

### Optional

- `credentials` (Block List, Max: 1) The KSQL Cluster API Credentials. (see [below for nested schema](#nestedblock--credentials))
- `delete_topic_on_destroy` (Boolean) Delete topic on destroy.
- `ignore_already_exists` (Boolean) Ignore already exists errors.
- `query_properties` (Map of String) Map of query properties
- `terminate_persistent_query` (Boolean) Terminate persistent query if needed.

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--credentials"></a>
### Nested Schema for `credentials`

Optional:

- `password` (String, Sensitive) The KSQL password.
- `url` (String, Sensitive) The KSQL URL.
- `username` (String, Sensitive) The KSQL username.


