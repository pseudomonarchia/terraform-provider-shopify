package shopify

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
)

var _ deliveryService = (*deliveryServiceImpl)(nil)

type deliveryService interface {
	Get(ctx context.Context, deliveryID string) (*DeliveryNode, error)
	Create(ctx context.Context, functionID string, delivery *DeliveryNode) (*DeliveryNode, error)
	Update(ctx context.Context, delivery *DeliveryNode) (*DeliveryNode, error)
	Delete(ctx context.Context, deliveryID string) (*DeliveryNode, error)
}

type deliveryServiceImpl struct {
	client shopifyAdminClient
}

type DeliveryNode struct {
	ID      string
	Title   string
	Enabled bool
}

func (d *deliveryServiceImpl) Get(ctx context.Context, deliveryID string) (*DeliveryNode, error) {
	gql := `
		query {
			deliveryCustomization(id: "%s") {
				id
				title
				enabled
			}
		}
	`

	gql = fmt.Sprintf(gql, deliveryID)
	r, err := d.client.exec(ctx, gql)
	if err != nil {
		return nil, err
	}

	jsonb, _ := json.Marshal(r)
	json := gjson.Parse(string(jsonb)).Get("deliveryCustomization")

	n := &DeliveryNode{
		ID:      json.Get("id").String(),
		Title:   json.Get("title").String(),
		Enabled: json.Get("enabled").Bool(),
	}

	return n, nil
}

func (d *deliveryServiceImpl) Create(ctx context.Context, functionID string, delivery *DeliveryNode) (*DeliveryNode, error) {
	gql := `
		mutation {
			deliveryCustomizationCreate(
				deliveryCustomization: {
					functionId: "%s"
					title: "%s"
					enabled: %t
				}
			) {
				deliveryCustomization {
					id
					title
					enabled
				}
				userErrors {
          field
          message
          code
        }
			}
		}
	`

	gql = fmt.Sprintf(gql, functionID, delivery.Title, delivery.Enabled)
	r, err := d.client.exec(ctx, gql)
	if err != nil {
		return nil, err
	}

	jsonb, _ := json.Marshal(r)
	json := gjson.Parse(string(jsonb)).Get("deliveryCustomizationCreate.deliveryCustomization")

	n := &DeliveryNode{
		ID:      json.Get("id").String(),
		Title:   json.Get("title").String(),
		Enabled: json.Get("enabled").Bool(),
	}

	return n, nil
}

func (d *deliveryServiceImpl) Update(ctx context.Context, delivery *DeliveryNode) (*DeliveryNode, error) {
	gql := `
		mutation {
			deliveryCustomizationUpdate(
				id: "%s"
				deliveryCustomization: {
					title: "%s"
					enabled: %t
				}
			) {
				deliveryCustomization {
					id
					title
					enabled
				}
			}
		}
	`

	gql = fmt.Sprintf(gql, delivery.ID, delivery.Title, delivery.Enabled)
	r, err := d.client.exec(ctx, gql)
	if err != nil {
		return nil, err
	}

	jsonb, _ := json.Marshal(r)
	json := gjson.Parse(string(jsonb)).Get("deliveryCustomizationUpdate.deliveryCustomization")

	n := &DeliveryNode{
		ID:      json.Get("id").String(),
		Title:   json.Get("title").String(),
		Enabled: json.Get("enabled").Bool(),
	}

	return n, nil
}

func (d *deliveryServiceImpl) Delete(ctx context.Context, deliveryID string) (*DeliveryNode, error) {
	gql := `
		mutation {
			deliveryCustomizationDelete(id: "%s") {
				deletedId
			}
		}
	`

	gql = fmt.Sprintf(gql, deliveryID)
	r, err := d.client.exec(ctx, gql)
	if err != nil {
		return nil, err
	}

	jsonb, _ := json.Marshal(r)
	id := gjson.Parse(string(jsonb)).Get("deliveryCustomizationDelete.deletedId").String()

	n := &DeliveryNode{
		ID: id,
	}

	return n, nil
}
