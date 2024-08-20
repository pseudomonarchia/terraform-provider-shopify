package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPubsubWebhookResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubWebhookResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("shopify_pubsub_webhook.test", "topic", "DISCOUNTS_CREATE"),
					resource.TestCheckResourceAttr("shopify_pubsub_webhook.test", "format", "JSON"),
					resource.TestCheckResourceAttr("shopify_pubsub_webhook.test", "pubsub_project", "test-project"),
					resource.TestCheckResourceAttr("shopify_pubsub_webhook.test", "pubsub_topic", "test-topic"),
				),
			},
			{
				ResourceName:      "shopify_pubsub_webhook.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPubsubWebhookResourceConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("shopify_pubsub_webhook.test", "topic", "DISCOUNTS_CREATE"),
					resource.TestCheckResourceAttr("shopify_pubsub_webhook.test", "format", "JSON"),
					resource.TestCheckResourceAttr("shopify_pubsub_webhook.test", "pubsub_project", "updated-project"),
					resource.TestCheckResourceAttr("shopify_pubsub_webhook.test", "pubsub_topic", "updated-topic"),
				),
			},
		},
	})
}

func testAccPubsubWebhookResourceConfig() string {
	return `
		resource "shopify_pubsub_webhook" "test" {
			topic          = "DISCOUNTS_CREATE"
			format         = "JSON"
			pubsub_project = "test-project"
			pubsub_topic   = "test-topic"
		}
	`
}

func testAccPubsubWebhookResourceConfigUpdate() string {
	return `
		resource "shopify_pubsub_webhook" "test" {
			topic          = "DISCOUNTS_CREATE"
			format         = "JSON"
			pubsub_project = "updated-project"
			pubsub_topic   = "updated-topic"
		}
	`
}
