terraform {
  required_providers {
    ksql = {
      source = "serhiimakogon/ksql"
      version = "1.1.0"
    }
  }
}

provider "ksql" {

}

resource "ksql_query" "stream" {
  query = "CREATE TABLE PRODUCTS_HOT_TABLE (ITEM_KEY STRING PRIMARY KEY) WITH (KAFKA_TOPIC = 'PRODUCTS_HOT',VALUE_FORMAT = 'AVRO');"

  credentials {
    url      = ""
    username = ""
    password = ""
  }
}
