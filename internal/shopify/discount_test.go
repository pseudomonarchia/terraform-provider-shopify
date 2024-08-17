package shopify

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDiscountService_Get(t *testing.T) {
	mockClient := new(mockShopifyAdminClient)
	service := &discountServiceImpl{client: mockClient}

	ctx := context.Background()
	discountID := "gid://shopify/DiscountAutomaticNode/12345"

	t.Run("Successful Get", func(t *testing.T) {
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

		mockClient.On("exec", ctx, mock.AnythingOfType("string")).Return(expectedResponse, nil).Once()

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
	})

	t.Run("Error in Get", func(t *testing.T) {
		mockClient.On("exec", ctx, mock.AnythingOfType("string")).Return(nil, errors.New("API error")).Once()

		discount, err := service.Get(ctx, discountID)

		assert.Error(t, err)
		assert.Nil(t, discount)
		assert.EqualError(t, err, "API error")

		mockClient.AssertExpectations(t)
	})
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

	t.Run("Successful Create", func(t *testing.T) {
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

		mockClient.On("exec", ctx, mock.AnythingOfType("string")).Return(expectedResponse, nil).Once()

		createdDiscount, err := service.Create(ctx, functionID, newDiscount)

		assert.NoError(t, err)
		assert.NotNil(t, createdDiscount)
		assert.Equal(t, "gid://shopify/DiscountAutomaticNode/12345", createdDiscount.ID)
		assert.Equal(t, newDiscount.Title, createdDiscount.Title)
		assert.Equal(t, newDiscount.StartsAt, createdDiscount.StartsAt)
		assert.Equal(t, newDiscount.EndsAt, createdDiscount.EndsAt)
		assert.Equal(t, newDiscount.CombinesWith, createdDiscount.CombinesWith)

		mockClient.AssertExpectations(t)
	})

	t.Run("Create without EndsAt", func(t *testing.T) {
		newDiscountNoEnd := &DiscountNode{
			Title:    "New Discount No End",
			StartsAt: "2023-02-01T00:00:00Z",
			CombinesWith: &DiscountCombinesWith{
				OrderDiscounts:    true,
				ProductDiscounts:  true,
				ShippingDiscounts: false,
			},
		}

		expectedResponse := map[string]interface{}{
			"discountAutomaticAppCreate": map[string]interface{}{
				"automaticAppDiscount": map[string]interface{}{
					"discountId": "gid://shopify/DiscountAutomaticNode/12346",
					"title":      "New Discount No End",
					"startsAt":   "2023-02-01T00:00:00Z",
					"endsAt":     nil,
					"combinesWith": map[string]interface{}{
						"orderDiscounts":    true,
						"productDiscounts":  true,
						"shippingDiscounts": false,
					},
				},
			},
		}

		mockClient.On("exec", ctx, mock.AnythingOfType("string")).Return(expectedResponse, nil).Once()

		createdDiscount, err := service.Create(ctx, functionID, newDiscountNoEnd)

		assert.NoError(t, err)
		assert.NotNil(t, createdDiscount)
		assert.Equal(t, "gid://shopify/DiscountAutomaticNode/12346", createdDiscount.ID)
		assert.Equal(t, newDiscountNoEnd.Title, createdDiscount.Title)
		assert.Equal(t, newDiscountNoEnd.StartsAt, createdDiscount.StartsAt)
		assert.Empty(t, createdDiscount.EndsAt)
		assert.Equal(t, newDiscountNoEnd.CombinesWith, createdDiscount.CombinesWith)

		mockClient.AssertExpectations(t)
	})

	t.Run("Error in Create", func(t *testing.T) {
		mockClient.On("exec", ctx, mock.AnythingOfType("string")).Return(nil, errors.New("API error")).Once()

		createdDiscount, err := service.Create(ctx, functionID, newDiscount)

		assert.Error(t, err)
		assert.Nil(t, createdDiscount)
		assert.EqualError(t, err, "API error")

		mockClient.AssertExpectations(t)
	})
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

	t.Run("Successful Update", func(t *testing.T) {
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

		mockClient.On("exec", ctx, mock.AnythingOfType("string")).Return(expectedResponse, nil).Once()

		result, err := service.Update(ctx, updatedDiscount)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, updatedDiscount.ID, result.ID)
		assert.Equal(t, updatedDiscount.Title, result.Title)
		assert.Equal(t, updatedDiscount.StartsAt, result.StartsAt)
		assert.Equal(t, updatedDiscount.EndsAt, result.EndsAt)
		assert.Equal(t, updatedDiscount.CombinesWith, result.CombinesWith)

		mockClient.AssertExpectations(t)
	})

	t.Run("Update without EndsAt", func(t *testing.T) {
		updatedDiscountNoEnd := &DiscountNode{
			ID:       "gid://shopify/DiscountAutomaticNode/12346",
			Title:    "Updated Discount No End",
			StartsAt: "2023-03-01T00:00:00Z",
			CombinesWith: &DiscountCombinesWith{
				OrderDiscounts:    false,
				ProductDiscounts:  true,
				ShippingDiscounts: true,
			},
		}

		expectedResponse := map[string]interface{}{
			"discountAutomaticAppUpdate": map[string]interface{}{
				"automaticAppDiscount": map[string]interface{}{
					"discountId": updatedDiscountNoEnd.ID,
					"title":      updatedDiscountNoEnd.Title,
					"startsAt":   updatedDiscountNoEnd.StartsAt,
					"endsAt":     nil,
					"combinesWith": map[string]interface{}{
						"orderDiscounts":    false,
						"productDiscounts":  true,
						"shippingDiscounts": true,
					},
				},
			},
		}

		mockClient.On("exec", ctx, mock.AnythingOfType("string")).Return(expectedResponse, nil).Once()

		result, err := service.Update(ctx, updatedDiscountNoEnd)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, updatedDiscountNoEnd.ID, result.ID)
		assert.Equal(t, updatedDiscountNoEnd.Title, result.Title)
		assert.Equal(t, updatedDiscountNoEnd.StartsAt, result.StartsAt)
		assert.Empty(t, result.EndsAt)
		assert.Equal(t, updatedDiscountNoEnd.CombinesWith, result.CombinesWith)

		mockClient.AssertExpectations(t)
	})

	t.Run("Error in Update", func(t *testing.T) {
		mockClient.On("exec", ctx, mock.AnythingOfType("string")).Return(nil, errors.New("API error")).Once()

		result, err := service.Update(ctx, updatedDiscount)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "API error")

		mockClient.AssertExpectations(t)
	})
}

func TestDiscountService_Delete(t *testing.T) {
	mockClient := new(mockShopifyAdminClient)
	service := &discountServiceImpl{client: mockClient}

	ctx := context.Background()
	discountID := "gid://shopify/DiscountAutomaticNode/12345"

	t.Run("Successful Delete", func(t *testing.T) {
		expectedResponse := map[string]interface{}{
			"discountAutomaticDelete": map[string]interface{}{
				"deletedAutomaticDiscountId": discountID,
			},
		}

		mockClient.On("exec", ctx, mock.AnythingOfType("string")).Return(expectedResponse, nil).Once()

		deletedDiscount, err := service.Delete(ctx, discountID)

		assert.NoError(t, err)
		assert.NotNil(t, deletedDiscount)
		assert.Equal(t, discountID, deletedDiscount.ID)

		mockClient.AssertExpectations(t)
	})

	t.Run("Error in Delete", func(t *testing.T) {
		mockClient.On("exec", ctx, mock.AnythingOfType("string")).Return(nil, errors.New("API error")).Once()

		deletedDiscount, err := service.Delete(ctx, discountID)

		assert.Error(t, err)
		assert.Nil(t, deletedDiscount)
		assert.EqualError(t, err, "API error")

		mockClient.AssertExpectations(t)
	})
}
