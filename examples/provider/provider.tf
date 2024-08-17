terraform {
  required_providers {
    sfr = {
      source  = "registry.terraform.io/pseudomonarchia/shopify-function-registry"
      version = "0.0.1"
    }
  }
}

provider "sfr" {
  store_domain       = "<store>.myshopify.com"
  store_access_token = "<access_token>"
  store_api_version  = "2024-07"
}
