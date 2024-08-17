package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFunctionDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// 讀取數據源測試
			{
				Config: testAccFunctionDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.shopify_function.test", "id"),
					resource.TestCheckResourceAttr("data.shopify_function.test", "title", "test_function"),
					resource.TestCheckResourceAttr("data.shopify_function.test", "app_title", "test_app"),
					resource.TestCheckResourceAttrSet("data.shopify_function.test", "api_type"),
				),
			},
		},
	})
}

func testAccFunctionDataSourceConfig() string {
	return `
		data "shopify_function" "test" {
			title     = "test_function"
			app_title = "test_app"
		}
	`
}

func TestAccFunctionDataSource_NotFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccFunctionDataSourceConfig_NotFound(),
				ExpectError: regexp.MustCompile(`No matching function found`),
			},
		},
	})
}

func testAccFunctionDataSourceConfig_NotFound() string {
	return `
		data "shopify_function" "test" {
			title     = "non_existent_function"
			app_title = "non_existent_app"
		}
	`
}
