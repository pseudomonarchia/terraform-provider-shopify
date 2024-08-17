terraform {
  required_providers {
    shopify = {
      source  = "pseudomonarchia/terraform-provider-shopify"
      version = "0.0.3"
    }
  }
}

provider "shopify" {
  store_domain       = "<store>.myshopify.com"
  store_access_token = "<access_token>"
  store_api_version  = "2024-07"
}
