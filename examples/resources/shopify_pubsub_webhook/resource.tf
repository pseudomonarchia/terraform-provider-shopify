resource "shopify_pubsub_webhook" "example" {
  topic          = "ORDERS_CREATE"
  format         = "JSON"
  pubsub_project = "test-project"
  pubsub_topic   = "test-topic"
}
