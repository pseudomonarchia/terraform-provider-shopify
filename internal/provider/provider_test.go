package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"shopify": providerserver.NewProtocol6WithError(New("test")()),
}

func TestAccProvider(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.shopify_function.test", "title", "product-discount"),
					resource.TestCheckResourceAttr("data.shopify_function.test", "app_title", "tf-testing"),
					resource.TestCheckResourceAttr("data.shopify_function.test", "api_type", "product_discounts"),
				),
			},
		},
	})
}

func testAccProviderConfig() string {
	return fmt.Sprintf(`
		provider "shopify" {
			store_domain       = "%s"
			store_access_token = "%s"
			store_api_version  = "%s"
		}

		data "shopify_function" "test" {
			title     = "product-discount"
			app_title = "tf-testing"
		}
	`,
		os.Getenv("SHOPIFY_STORE_DOMAIN"),
		os.Getenv("SHOPIFY_STORE_ACCESS_TOKEN"),
		os.Getenv("SHOPIFY_STORE_API_VERSION"),
	)
}
