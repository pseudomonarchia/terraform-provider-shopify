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
			// 創建並讀取測試
			{
				Config: testAccDeliveryCustomResourceConfig("test_delivery", true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("shopify_delivery.test", "title", "test_delivery"),
					resource.TestCheckResourceAttr("shopify_delivery.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("shopify_delivery.test", "id"),
					resource.TestCheckResourceAttrSet("shopify_delivery.test", "function_id"),
				),
			},
			// 更新測試
			{
				Config: testAccDeliveryCustomResourceConfig("updated_delivery", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("shopify_delivery.test", "title", "updated_delivery"),
					resource.TestCheckResourceAttr("shopify_delivery.test", "enabled", "false"),
				),
			},
			// 導入測試
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
				function_id = "00000000-0000-0000-0000-000000000000"
				title       = %q
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
