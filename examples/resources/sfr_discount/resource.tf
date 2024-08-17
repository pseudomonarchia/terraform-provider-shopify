resource "sfr_discount" "example" {
  function_id = "<UUID>"
  title       = "Automatic Discount"
  starts_at   = "2024-01-09T00:00:00Z"
  combines_with = {
    order_discounts    = true
    product_discounts  = true
    shipping_discounts = true
  }
}
