---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "shopify_delivery Resource - shopify"
subcategory: ""
description: |-
  Shopify Function Delivery Customization Resource
---

# shopify_delivery (Resource)

Shopify Function Delivery Customization Resource

## Example Usage

```terraform
resource "shopify_delivery" "example" {
  function_id = "<UUID>"
  title       = "Delivery Customization"
  enabled     = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `enabled` (Boolean)
- `function_id` (String)
- `title` (String)

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import shopify_delivery.example <delivery_id>,<function_id>
```
