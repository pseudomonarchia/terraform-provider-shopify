package provider

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDiscountAutomaticResource(t *testing.T) {
	startTime := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	endTime := time.Now().UTC().Add(24 * time.Hour).Format("2006-01-02T15:04:05Z")
	updatedStartTime := time.Now().UTC().Add(1 * time.Hour).Format("2006-01-02T15:04:05Z")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// 創建並讀取測試
			{
				Config: testAccDiscountAutomaticResourceConfig(startTime, endTime, true, false, true, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sfr_discount.test", "title", "test_discount"),
					resource.TestCheckResourceAttr("sfr_discount.test", "starts_at", startTime),
					resource.TestCheckResourceAttr("sfr_discount.test", "ends_at", endTime),
					resource.TestCheckResourceAttr("sfr_discount.test", "combines_with.order_discounts", "true"),
					resource.TestCheckResourceAttr("sfr_discount.test", "combines_with.product_discounts", "false"),
					resource.TestCheckResourceAttr("sfr_discount.test", "combines_with.shipping_discounts", "true"),
					resource.TestCheckResourceAttrSet("sfr_discount.test", "id"),
					resource.TestCheckResourceAttrSet("sfr_discount.test", "function_id"),
				),
			},
			// 更新測試
			{
				Config: testAccDiscountAutomaticResourceConfig(updatedStartTime, "", false, true, false, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sfr_discount.test", "title", "updated_discount"),
					resource.TestCheckResourceAttr("sfr_discount.test", "starts_at", updatedStartTime),
					resource.TestCheckResourceAttr("sfr_discount.test", "ends_at", ""),
					resource.TestCheckResourceAttr("sfr_discount.test", "combines_with.order_discounts", "false"),
					resource.TestCheckResourceAttr("sfr_discount.test", "combines_with.product_discounts", "true"),
					resource.TestCheckResourceAttr("sfr_discount.test", "combines_with.shipping_discounts", "false"),
				),
			},
			// 導入測試
			{
				ResourceName:      "sfr_discount.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccDiscountAutomaticResourceImportStateIdFunc,
			},
		},
	})
}

func testAccDiscountAutomaticResourceConfig(startsAt, endsAt string, orderDiscounts, productDiscounts, shippingDiscounts bool, isUpdated bool) string {
	endsAtConfig := ""
	if endsAt != "" {
		endsAtConfig = fmt.Sprintf(`ends_at = %q`, endsAt)
	}

	title := "test_discount"
	if isUpdated {
		title = "updated_discount"
	}

	return fmt.Sprintf(
		`
			resource "sfr_discount" "test" {
				function_id = "00000000-0000-0000-0000-000000000000"
				title       = %q
				starts_at   = %q
				%s
				combines_with {
					order_discounts    = %t
					product_discounts  = %t
					shipping_discounts = %t
				}
			}
		`,
		title,
		startsAt,
		endsAtConfig,
		orderDiscounts,
		productDiscounts,
		shippingDiscounts,
	)
}

func testAccDiscountAutomaticResourceImportStateIdFunc(s *terraform.State) (string, error) {
	rs, ok := s.RootModule().Resources["sfr_discount.test"]
	if !ok {
		return "", fmt.Errorf("Resource not found: sfr_discount.test")
	}

	return fmt.Sprintf("%s,%s", rs.Primary.ID, rs.Primary.Attributes["function_id"]), nil
}
