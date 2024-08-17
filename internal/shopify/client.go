package shopify

import (
	"context"
	"fmt"

	"github.com/machinebox/graphql"
)

type ShopifyAdminClinetImpl struct {
	storeDomain      string
	storeAccessToken string
	storeApiVersion  string
	local            bool

	Discount discountService
	Payment  paymentService
	Function FunctionService
	Delivery deliveryService
}

type shopifyAdminClient interface {
	exec(ctx context.Context, query string) (any, error)
}

func New(
	storeDomain string,
	storeAccessToken string,
	storeApiVersion string,
) *ShopifyAdminClinetImpl {
	c := &ShopifyAdminClinetImpl{
		storeDomain:      storeDomain,
		storeAccessToken: storeAccessToken,
		storeApiVersion:  storeApiVersion,
	}

	c.Discount = &discountServiceImpl{c}
	c.Function = &FunctionServiceImpl{c}
	c.Payment = &paymentServiceImpl{c}
	c.Delivery = &deliveryServiceImpl{c}

	return c
}

func (s *ShopifyAdminClinetImpl) exec(ctx context.Context, query string) (any, error) {
	scheme := "https"
	if s.local {
		scheme = "http"
	}

	endpoint := fmt.Sprintf("%s://%s/admin/api/%s/graphql.json", scheme, s.storeDomain, s.storeApiVersion)
	client := graphql.NewClient(endpoint)
	req := graphql.NewRequest(query)

	req.Header.Set("X-Shopify-Access-Token", s.storeAccessToken)
	req.Header.Set("Cache-Control", "no-cache")

	var res any

	if err := client.Run(ctx, req, &res); err != nil {
		return "", err
	}

	return res, nil
}
