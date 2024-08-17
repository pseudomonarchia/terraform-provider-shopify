package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"sfr": providerserver.NewProtocol6WithError(New("test")()),
}

func TestAccProvider(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					provider "sfr" {
						store_domain       = "example.myshopify.com"
						store_access_token = "test_token"
						store_api_version  = "2024-01"
					}
				`,
			},
		},
	})
}

func TestAccProviderInvalidDomain(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					provider "sfr" {
						store_domain       = "invalid-domain"
						store_access_token = "test_token"
						store_api_version  = "2024-01"
					}
				`,
				ExpectError: regexp.MustCompile(`must be a valid Shopify store domain`),
			},
		},
	})
}

func TestAccProviderInvalidAPIVersion(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					provider "sfr" {
						store_domain       = "example.myshopify.com"
						store_access_token = "test_token"
						store_api_version  = "invalid-version"
					}
				`,
				ExpectError: regexp.MustCompile(`must be a valid Shopify API version`),
			},
		},
	})
}
