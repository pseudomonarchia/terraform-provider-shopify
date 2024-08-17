# Shopify Provider Project

This project contains the Shopify Provider configuration for Terraform.

## Configuration Description

This project uses Terraform to manage Shopify Function resource registration for Shopify stores.
The main configuration file is `provider.tf`, which defines the settings for the Shopify Provider.

### Provider Configuration

```hcl
terraform {
  required_providers {
    shopify = {
      source  = "registry.terraform.io/pseudomonarchia/terraform-provider-shopify"
      version = "0.0.2"
    }
  }
}

provider "shopify" {
  store_domain       = "<store>.myshopify.com"
  store_access_token = "<access_token>"
  store_api_version  = "2024-07"
}
```

- `store_domain`: Your Shopify store domain
- `store_access_token`: Access token for accessing the Shopify API
- `store_api_version`: The Shopify API version being used

## Usage Instructions

1. Ensure Terraform is installed.
2. Replace `<store>` with your actual Shopify store name.
3. Replace `<access_token>` with your Shopify API access token.
4. Run `terraform init` to initialize the project.
5. Use `terraform plan` and `terraform apply` to manage your Shopify resources.

## Notes

- Make sure to securely store your access token and do not expose it or commit it to version control systems.
- Regularly check and update the API version to ensure compatibility with the latest Shopify API.

## Contributions

Issue reports and improvement suggestions are welcome.
