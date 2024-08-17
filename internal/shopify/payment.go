package shopify

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
)

var _ paymentService = (*paymentServiceImpl)(nil)

type paymentService interface {
	Get(ctx context.Context, paymentID string) (*PaymentNode, error)
	Create(ctx context.Context, functionID string, payment *PaymentNode) (*PaymentNode, error)
	Update(ctx context.Context, payment *PaymentNode) (*PaymentNode, error)
	Delete(ctx context.Context, paymentID string) (*PaymentNode, error)
}

type paymentServiceImpl struct {
	client shopifyAdminClient
}

type PaymentNode struct {
	ID      string
	Title   string
	Enabled bool
}

func (p *paymentServiceImpl) Get(ctx context.Context, paymentID string) (*PaymentNode, error) {
	gql := `
		query {
			paymentCustomization(id: "%s") {
				id
				title
				enabled
			}
		}
	`

	gql = fmt.Sprintf(gql, paymentID)
	r, err := p.client.exec(ctx, gql)
	if err != nil {
		return nil, err
	}

	jsonb, _ := json.Marshal(r)
	json := gjson.Parse(string(jsonb)).Get("paymentCustomization")

	n := &PaymentNode{
		ID:      json.Get("id").String(),
		Title:   json.Get("title").String(),
		Enabled: json.Get("enabled").Bool(),
	}

	return n, nil
}

func (p *paymentServiceImpl) Create(ctx context.Context, functionID string, payment *PaymentNode) (*PaymentNode, error) {
	gql := `
		mutation {
			paymentCustomizationCreate(
				paymentCustomization: {
					functionId: "%s"
					title: "%s"
					enabled: %t
				}
			) {
				paymentCustomization {
					id
					title
					enabled
				}
			}
		}
	`

	gql = fmt.Sprintf(gql, functionID, payment.Title, payment.Enabled)
	r, err := p.client.exec(ctx, gql)
	if err != nil {
		return nil, err
	}

	jsonb, _ := json.Marshal(r)
	json := gjson.Parse(string(jsonb)).Get("paymentCustomizationCreate.paymentCustomization")

	n := &PaymentNode{
		ID:      json.Get("id").String(),
		Title:   json.Get("title").String(),
		Enabled: json.Get("enabled").Bool(),
	}

	return n, nil
}

func (p *paymentServiceImpl) Update(ctx context.Context, payment *PaymentNode) (*PaymentNode, error) {
	gql := `
		mutation {
			paymentCustomizationUpdate(
				id: "%s"
				paymentCustomization: {
					title: "%s"
					enabled: %t
				}
			) {
				paymentCustomization {
					id
					title
					enabled
				}
			}
		}
	`

	gql = fmt.Sprintf(gql, payment.ID, payment.Title, payment.Enabled)
	r, err := p.client.exec(ctx, gql)
	if err != nil {
		return nil, err
	}

	jsonb, _ := json.Marshal(r)
	json := gjson.Parse(string(jsonb)).Get("paymentCustomizationUpdate.paymentCustomization")

	n := &PaymentNode{
		ID:      json.Get("id").String(),
		Title:   json.Get("title").String(),
		Enabled: json.Get("enabled").Bool(),
	}

	return n, nil
}

func (p *paymentServiceImpl) Delete(ctx context.Context, paymentID string) (*PaymentNode, error) {
	gql := `
		mutation {
			paymentCustomizationDelete(id: "%s") {
				deletedId
			}
		}
	`

	gql = fmt.Sprintf(gql, paymentID)
	r, err := p.client.exec(ctx, gql)
	if err != nil {
		return nil, err
	}

	jsonb, _ := json.Marshal(r)
	id := gjson.Parse(string(jsonb)).Get("paymentCustomizationDelete.deletedId").String()

	n := &PaymentNode{
		ID: id,
	}

	return n, nil
}
