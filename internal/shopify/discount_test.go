package shopify

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDiscountService_Get(t *testing.T) {
	mockClient := new(mockShopifyAdminClient)
	service := &discountServiceImpl{client: mockClient}

	ctx := context.Background()
	discountID := "gid://shopify/DiscountAutomaticNode/12345"

	expectedResponse := map[string]interface{}{
		"discountNode": map[string]interface{}{
			"discount": map[string]interface{}{
				"discountId": discountID,
				"title":      "Test Discount",
				"startsAt":   "2023-01-01T00:00:00Z",
				"endsAt":     "2023-12-31T23:59:59Z",
				"combinesWith": map[string]interface{}{
					"orderDiscounts":    true,
					"productDiscounts":  false,
					"shippingDiscounts": true,
				},
			},
		},
	}

	mockClient.On("exec", ctx, mock.AnythingOfType("string")).Return(expectedResponse, nil)

	discount, err := service.Get(ctx, discountID)

	assert.NoError(t, err)
	assert.NotNil(t, discount)
	assert.Equal(t, discountID, discount.ID)
	assert.Equal(t, "Test Discount", discount.Title)
	assert.Equal(t, "2023-01-01T00:00:00Z", discount.StartsAt)
	assert.Equal(t, "2023-12-31T23:59:59Z", discount.EndsAt)
	assert.True(t, discount.CombinesWith.OrderDiscounts)
	assert.False(t, discount.CombinesWith.ProductDiscounts)
	assert.True(t, discount.CombinesWith.ShippingDiscounts)

	mockClient.AssertExpectations(t)
}

func TestDiscountService_Create(t *testing.T) {
	mockClient := new(mockShopifyAdminClient)
	service := &discountServiceImpl{client: mockClient}

	ctx := context.Background()
	functionID := "gid://shopify/DiscountAutomaticNode/67890"

	newDiscount := &DiscountNode{
		Title:    "New Discount",
		StartsAt: "2023-02-01T00:00:00Z",
		EndsAt:   "2023-11-30T23:59:59Z",
		CombinesWith: &DiscountCombinesWith{
			OrderDiscounts:    true,
			ProductDiscounts:  true,
			ShippingDiscounts: false,
		},
	}

	expectedResponse := map[string]interface{}{
		"discountAutomaticAppCreate": map[string]interface{}{
			"automaticAppDiscount": map[string]interface{}{
				"discountId": "gid://shopify/DiscountAutomaticNode/12345",
				"title":      "New Discount",
				"startsAt":   "2023-02-01T00:00:00Z",
				"endsAt":     "2023-11-30T23:59:59Z",
				"combinesWith": map[string]interface{}{
					"orderDiscounts":    true,
					"productDiscounts":  true,
					"shippingDiscounts": false,
				},
			},
		},
	}

	mockClient.On("exec", ctx, mock.AnythingOfType("string")).Return(expectedResponse, nil)

	createdDiscount, err := service.Create(ctx, functionID, newDiscount)

	assert.NoError(t, err)
	assert.NotNil(t, createdDiscount)
	assert.Equal(t, "gid://shopify/DiscountAutomaticNode/12345", createdDiscount.ID)
	assert.Equal(t, newDiscount.Title, createdDiscount.Title)
	assert.Equal(t, newDiscount.StartsAt, createdDiscount.StartsAt)
	assert.Equal(t, newDiscount.EndsAt, createdDiscount.EndsAt)
	assert.Equal(t, newDiscount.CombinesWith, createdDiscount.CombinesWith)

	mockClient.AssertExpectations(t)
}

func TestDiscountService_Update(t *testing.T) {
	mockClient := new(mockShopifyAdminClient)
	service := &discountServiceImpl{client: mockClient}

	ctx := context.Background()
	updatedDiscount := &DiscountNode{
		ID:       "gid://shopify/DiscountAutomaticNode/12345",
		Title:    "Updated Discount",
		StartsAt: "2023-03-01T00:00:00Z",
		EndsAt:   "2023-10-31T23:59:59Z",
		CombinesWith: &DiscountCombinesWith{
			OrderDiscounts:    false,
			ProductDiscounts:  true,
			ShippingDiscounts: true,
		},
	}

	expectedResponse := map[string]interface{}{
		"discountAutomaticAppUpdate": map[string]interface{}{
			"automaticAppDiscount": map[string]interface{}{
				"discountId": updatedDiscount.ID,
				"title":      updatedDiscount.Title,
				"startsAt":   updatedDiscount.StartsAt,
				"endsAt":     updatedDiscount.EndsAt,
				"combinesWith": map[string]interface{}{
					"orderDiscounts":    false,
					"productDiscounts":  true,
					"shippingDiscounts": true,
				},
			},
		},
	}

	mockClient.On("exec", ctx, mock.AnythingOfType("string")).Return(expectedResponse, nil)

	result, err := service.Update(ctx, updatedDiscount)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, updatedDiscount.ID, result.ID)
	assert.Equal(t, updatedDiscount.Title, result.Title)
	assert.Equal(t, updatedDiscount.StartsAt, result.StartsAt)
	assert.Equal(t, updatedDiscount.EndsAt, result.EndsAt)
	assert.Equal(t, updatedDiscount.CombinesWith, result.CombinesWith)

	mockClient.AssertExpectations(t)
}

func TestDiscountService_Delete(t *testing.T) {
	mockClient := new(mockShopifyAdminClient)
	service := &discountServiceImpl{client: mockClient}

	ctx := context.Background()
	discountID := "gid://shopify/DiscountAutomaticNode/12345"

	expectedResponse := map[string]interface{}{
		"discountAutomaticDelete": map[string]interface{}{
			"deletedAutomaticDiscountId": discountID,
		},
	}

	mockClient.On("exec", ctx, mock.AnythingOfType("string")).Return(expectedResponse, nil)

	deletedDiscount, err := service.Delete(ctx, discountID)

	assert.NoError(t, err)
	assert.NotNil(t, deletedDiscount)
	assert.Equal(t, discountID, deletedDiscount.ID)

	mockClient.AssertExpectations(t)
}
