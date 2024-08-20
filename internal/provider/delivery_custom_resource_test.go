package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDeliveryCustomResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeliveryCustomResourceConfig("create-test", true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("shopify_delivery.test", "title", "create-test"),
					resource.TestCheckResourceAttr("shopify_delivery.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("shopify_delivery.test", "id"),
					resource.TestCheckResourceAttrSet("shopify_delivery.test", "function_id"),
				),
			},
			{
				Config: testAccDeliveryCustomResourceConfig("update-test", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("shopify_delivery.test", "title", "update-test"),
					resource.TestCheckResourceAttr("shopify_delivery.test", "enabled", "false"),
				),
			},
			{
				ResourceName:      "shopify_delivery.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccDeliveryCustomResourceImportStateIdFunc,
			},
		},
	})
}

func testAccDeliveryCustomResourceConfig(title string, enabled bool) string {
	return fmt.Sprintf(
		`
			resource "shopify_delivery" "test" {
				function_id = "3a2c6a43-6ac1-4d4d-bbd9-59286cc33740"
				title       = "%s"
				enabled     = %t
			}
		`,
		title,
		enabled,
	)
}

func testAccDeliveryCustomResourceImportStateIdFunc(s *terraform.State) (string, error) {
	rs, ok := s.RootModule().Resources["shopify_delivery.test"]
	if !ok {
		return "", fmt.Errorf("Resource not found: shopify_delivery.test")
	}

	return fmt.Sprintf("%s,%s", rs.Primary.ID, rs.Primary.Attributes["function_id"]), nil
}
