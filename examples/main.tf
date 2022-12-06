terraform {
  required_providers {
    ksql = {
      source = "serhiimakogon/ksql"
      version = "0.1.0"
    }
  }
}

provider "ksql" {
  url      = ""
  username = ""
  password = ""
}

resource "ksql_query" "products_hot_table_table" {
  name  = "products_hot_table"
  type  = "table"
  query = "CREATE TABLE PRODUCTS_HOT_TABLE (ITEM_KEY STRING PRIMARY KEY) WITH (KAFKA_TOPIC = 'PRODUCTS_HOT',VALUE_FORMAT = 'AVRO');"
}
