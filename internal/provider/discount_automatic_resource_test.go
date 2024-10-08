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
			{
				Config: testAccDiscountAutomaticResourceConfig(startTime, "", true, false, true, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("shopify_discount.test", "title", "test_discount"),
					resource.TestCheckResourceAttr("shopify_discount.test", "starts_at", startTime),
					resource.TestCheckResourceAttr("shopify_discount.test", "combines_with.order_discounts", "true"),
					resource.TestCheckResourceAttr("shopify_discount.test", "combines_with.product_discounts", "false"),
					resource.TestCheckResourceAttr("shopify_discount.test", "combines_with.shipping_discounts", "true"),
					resource.TestCheckResourceAttrSet("shopify_discount.test", "id"),
					resource.TestCheckResourceAttrSet("shopify_discount.test", "function_id"),
				),
			},
			{
				Config: testAccDiscountAutomaticResourceConfig(updatedStartTime, endTime, false, true, false, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("shopify_discount.test", "title", "updated_discount"),
					resource.TestCheckResourceAttr("shopify_discount.test", "starts_at", updatedStartTime),
					resource.TestCheckResourceAttr("shopify_discount.test", "ends_at", endTime),
					resource.TestCheckResourceAttr("shopify_discount.test", "combines_with.order_discounts", "false"),
					resource.TestCheckResourceAttr("shopify_discount.test", "combines_with.product_discounts", "true"),
					resource.TestCheckResourceAttr("shopify_discount.test", "combines_with.shipping_discounts", "false"),
				),
			},
			{
				ResourceName:      "shopify_discount.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccDiscountAutomaticResourceImportStateIdFunc,
			},
		},
	})
}

func testAccDiscountAutomaticResourceConfig(
	startsAt,
	endsAt string,
	orderDiscounts,
	productDiscounts,
	shippingDiscounts bool,
	isUpdated bool,
) string {
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
			resource "shopify_discount" "test" {
				function_id = "07224386-3c16-4f9e-b8ba-da049b6afc66"
				title       = %q
				starts_at   = %q
				%s
				combines_with = {
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
	rs, ok := s.RootModule().Resources["shopify_discount.test"]
	if !ok {
		return "", fmt.Errorf("Resource not found: shopify_discount.test")
	}

	return fmt.Sprintf("%s,%s", rs.Primary.ID, rs.Primary.Attributes["function_id"]), nil
}
