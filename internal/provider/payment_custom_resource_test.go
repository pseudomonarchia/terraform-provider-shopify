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
			// 創建並讀取測試
			{
				Config: testAccPaymentCustomResourceConfig("test_payment", true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sfr_payment.test", "title", "test_payment"),
					resource.TestCheckResourceAttr("sfr_payment.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("sfr_payment.test", "id"),
					resource.TestCheckResourceAttrSet("sfr_payment.test", "function_id"),
				),
			},
			// 更新測試
			{
				Config: testAccPaymentCustomResourceConfig("updated_payment", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sfr_payment.test", "title", "updated_payment"),
					resource.TestCheckResourceAttr("sfr_payment.test", "enabled", "false"),
				),
			},
			// 導入測試
			{
				ResourceName:      "sfr_payment.test",
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
			resource "sfr_payment" "test" {
				function_id = "00000000-0000-0000-0000-000000000000"
				title       = %q
				enabled     = %t
			}
		`,
		title,
		enabled,
	)
}

func testAccPaymentCustomResourceImportStateIdFunc(s *terraform.State) (string, error) {
	rs, ok := s.RootModule().Resources["sfr_payment.test"]
	if !ok {
		return "", fmt.Errorf("Resource not found: sfr_payment.test")
	}

	return fmt.Sprintf("%s,%s", rs.Primary.ID, rs.Primary.Attributes["function_id"]), nil
}
