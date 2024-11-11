package shopify

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
)

var _ discountService = (*discountServiceImpl)(nil)

type discountService interface {
	Get(ctx context.Context, discountID string) (*DiscountNode, error)
	Create(ctx context.Context, functionID string, discount *DiscountNode) (*DiscountNode, error)
	Update(ctx context.Context, discount *DiscountNode) (*DiscountNode, error)
	Delete(ctx context.Context, discountID string) (*DiscountNode, error)
}

type discountServiceImpl struct {
	client shopifyAdminClient
}

type DiscountNode struct {
	ID           string
	Title        string
	StartsAt     string
	EndsAt       string
	CombinesWith *DiscountCombinesWith
}

type DiscountCombinesWith struct {
	OrderDiscounts    bool
	ProductDiscounts  bool
	ShippingDiscounts bool
}

func (d *discountServiceImpl) Get(
	ctx context.Context,
	discountID string,
) (*DiscountNode, error) {
	gql := `
		query {
			discountNode(id: "%s"){
				discount {
				... on DiscountAutomaticApp {
						discountId
						title
						startsAt
						endsAt
						combinesWith {
							orderDiscounts
							productDiscounts
							shippingDiscounts
						}
					}
				}
			}
		}
	`

	gql = fmt.Sprintf(gql, discountID)
	r, err := d.client.exec(ctx, gql)
	if err != nil {
		return nil, err
	}

	jsonb, _ := json.Marshal(r)
	json := gjson.Parse(string(jsonb)).Get("discountNode.discount")

	n := &DiscountNode{
		ID:       json.Get("discountId").String(),
		Title:    json.Get("title").String(),
		StartsAt: json.Get("startsAt").String(),
		EndsAt:   json.Get("endsAt").String(),
		CombinesWith: &DiscountCombinesWith{
			OrderDiscounts:    json.Get("combinesWith.orderDiscounts").Bool(),
			ProductDiscounts:  json.Get("combinesWith.productDiscounts").Bool(),
			ShippingDiscounts: json.Get("combinesWith.shippingDiscounts").Bool(),
		},
	}

	return n, nil
}

func (d *discountServiceImpl) Create(
	ctx context.Context,
	functionID string,
	discount *DiscountNode,
) (*DiscountNode, error) {
	gql := `
		mutation {
			discountAutomaticAppCreate(
				automaticAppDiscount: {
					functionId: "%s"
					title: "%s"
					startsAt: "%s"
					%s
					combinesWith: {
						orderDiscounts: %t
						productDiscounts: %t
						shippingDiscounts: %t
					}
				}
			) {
				automaticAppDiscount {
					discountId
					title
					startsAt
					endsAt
					combinesWith {
						orderDiscounts
						productDiscounts
						shippingDiscounts
					}
				}
			}
		}
	`

	endsAtField := "endsAt: null"
	if discount.EndsAt != "" {
		endsAtField = fmt.Sprintf(`endsAt: "%s"`, discount.EndsAt)
	}

	gql = fmt.Sprintf(
		gql,
		functionID,
		discount.Title,
		discount.StartsAt,
		endsAtField,
		discount.CombinesWith.OrderDiscounts,
		discount.CombinesWith.ProductDiscounts,
		discount.CombinesWith.ShippingDiscounts,
	)

	r, err := d.client.exec(ctx, gql)
	if err != nil {
		return nil, err
	}

	jsonb, _ := json.Marshal(r)
	json := gjson.
		Parse(string(jsonb)).
		Get("discountAutomaticAppCreate.automaticAppDiscount")

	n := &DiscountNode{
		ID:       json.Get("discountId").String(),
		Title:    json.Get("title").String(),
		StartsAt: json.Get("startsAt").String(),
		EndsAt:   json.Get("endsAt").String(),
		CombinesWith: &DiscountCombinesWith{
			OrderDiscounts:    json.Get("combinesWith.orderDiscounts").Bool(),
			ProductDiscounts:  json.Get("combinesWith.productDiscounts").Bool(),
			ShippingDiscounts: json.Get("combinesWith.shippingDiscounts").Bool(),
		},
	}

	return n, nil
}

func (d *discountServiceImpl) Update(
	ctx context.Context,
	discount *DiscountNode,
) (*DiscountNode, error) {
	gql := `
		mutation {
			discountAutomaticAppUpdate(
				id: "%s"
				automaticAppDiscount: { 
					title: "%s"
					startsAt: "%s"
					%s
					combinesWith: {
						orderDiscounts: %t
						productDiscounts: %t
						shippingDiscounts: %t
					}
				}
			) 
			{
				automaticAppDiscount {
					discountId
					title
					startsAt
					endsAt
					combinesWith {
						orderDiscounts
						productDiscounts
						shippingDiscounts
					}
				}
				userErrors {
      		field
      		message
      		code
    		}
			}
		}
	`

	endsAtField := "endsAt: null"
	if discount.EndsAt != "" {
		endsAtField = fmt.Sprintf(`endsAt: "%s"`, discount.EndsAt)
	}

	gql = fmt.Sprintf(
		gql,
		discount.ID,
		discount.Title,
		discount.StartsAt,
		endsAtField,
		discount.CombinesWith.OrderDiscounts,
		discount.CombinesWith.ProductDiscounts,
		discount.CombinesWith.ShippingDiscounts,
	)

	r, err := d.client.exec(ctx, gql)
	if err != nil {
		return nil, err
	}

	jsonb, _ := json.Marshal(r)
	json := gjson.
		Parse(string(jsonb)).
		Get("discountAutomaticAppUpdate.automaticAppDiscount")

	n := &DiscountNode{
		ID:       json.Get("discountId").String(),
		Title:    json.Get("title").String(),
		StartsAt: json.Get("startsAt").String(),
		EndsAt:   json.Get("endsAt").String(),
		CombinesWith: &DiscountCombinesWith{
			OrderDiscounts:    json.Get("combinesWith.orderDiscounts").Bool(),
			ProductDiscounts:  json.Get("combinesWith.productDiscounts").Bool(),
			ShippingDiscounts: json.Get("combinesWith.shippingDiscounts").Bool(),
		},
	}

	return n, nil
}

func (d *discountServiceImpl) Delete(
	ctx context.Context,
	discountID string,
) (*DiscountNode, error) {
	gql := `
		mutation discountAutomaticDelete {
			discountAutomaticDelete(id: "%s") {
				deletedAutomaticDiscountId
			}
		}
	`

	gql = fmt.Sprintf(gql, discountID)
	r, err := d.client.exec(ctx, gql)
	if err != nil {
		return nil, err
	}

	jsonb, _ := json.Marshal(r)
	id := gjson.
		Parse(string(jsonb)).
		Get("discountAutomaticDelete.deletedAutomaticDiscountId").
		String()

	n := &DiscountNode{
		ID: id,
	}

	return n, nil
}
