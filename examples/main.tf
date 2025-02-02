terraform {
  required_providers {
    ksql = {
      source = "serhiimakogon/ksql"
      version = "1.2.0-rc"
    }
  }
}

provider "ksql" {}

resource "ksql_query" "products_hot_table_table" {
  name  = "products_hot_table"
  type  = "table"
  query = "CREATE TABLE PRODUCTS_HOT_TABLE (ITEM_KEY STRING PRIMARY KEY) WITH (KAFKA_TOPIC = 'PRODUCTS_HOT',VALUE_FORMAT = 'AVRO');"
}
