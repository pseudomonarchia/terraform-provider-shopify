package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccPaymentCustomResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPaymentCustomResourceConfig("test_payment", true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("shopify_payment.test", "title", "test_payment"),
					resource.TestCheckResourceAttr("shopify_payment.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("shopify_payment.test", "id"),
					resource.TestCheckResourceAttrSet("shopify_payment.test", "function_id"),
				),
			},
			{
				Config: testAccPaymentCustomResourceConfig("updated_payment", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("shopify_payment.test", "title", "updated_payment"),
					resource.TestCheckResourceAttr("shopify_payment.test", "enabled", "false"),
				),
			},
			{
				ResourceName:      "shopify_payment.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccPaymentCustomResourceImportStateIdFunc,
			},
		},
	})
}

func testAccPaymentCustomResourceConfig(title string, enabled bool) string {
	return fmt.Sprintf(
		`
			resource "shopify_payment" "test" {
				function_id = "f2e906be-a93a-48c6-a2cc-99c64e5ab816"
				title       = %q
				enabled     = %t
			}
		`,
		title,
		enabled,
	)
}

func testAccPaymentCustomResourceImportStateIdFunc(s *terraform.State) (string, error) {
	rs, ok := s.RootModule().Resources["shopify_payment.test"]
	if !ok {
		return "", fmt.Errorf("Resource not found: shopify_payment.test")
	}

	return fmt.Sprintf("%s,%s", rs.Primary.ID, rs.Primary.Attributes["function_id"]), nil
}
