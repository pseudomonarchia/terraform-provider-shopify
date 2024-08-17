package shopify

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPaymentService_Get(t *testing.T) {
	mockClient := new(mockShopifyAdminClient)
	service := &paymentServiceImpl{client: mockClient}

	ctx := context.Background()
	paymentID := "gid://shopify/PaymentCustomization/1"

	expectedResponse := map[string]interface{}{
		"paymentCustomization": map[string]interface{}{
			"id":      paymentID,
			"title":   "Test Payment",
			"enabled": true,
		},
	}

	mockClient.On("exec", ctx, mock.AnythingOfType("string")).Return(expectedResponse, nil)

	payment, err := service.Get(ctx, paymentID)

	assert.NoError(t, err)
	assert.NotNil(t, payment)
	assert.Equal(t, paymentID, payment.ID)
	assert.Equal(t, "Test Payment", payment.Title)
	assert.True(t, payment.Enabled)

	mockClient.AssertExpectations(t)
}

func TestPaymentService_Create(t *testing.T) {
	mockClient := new(mockShopifyAdminClient)
	service := &paymentServiceImpl{client: mockClient}

	ctx := context.Background()
	functionID := "gid://shopify/ShopifyFunction/1"

	newPayment := &PaymentNode{
		Title:   "New Payment",
		Enabled: true,
	}

	expectedResponse := map[string]interface{}{
		"paymentCustomizationCreate": map[string]interface{}{
			"paymentCustomization": map[string]interface{}{
				"id":      "gid://shopify/PaymentCustomization/1",
				"title":   "New Payment",
				"enabled": true,
			},
		},
	}

	mockClient.On("exec", ctx, mock.AnythingOfType("string")).Return(expectedResponse, nil)

	createdPayment, err := service.Create(ctx, functionID, newPayment)

	assert.NoError(t, err)
	assert.NotNil(t, createdPayment)
	assert.Equal(t, "gid://shopify/PaymentCustomization/1", createdPayment.ID)
	assert.Equal(t, newPayment.Title, createdPayment.Title)
	assert.Equal(t, newPayment.Enabled, createdPayment.Enabled)

	mockClient.AssertExpectations(t)
}

func TestPaymentService_Update(t *testing.T) {
	mockClient := new(mockShopifyAdminClient)
	service := &paymentServiceImpl{client: mockClient}

	ctx := context.Background()
	updatedPayment := &PaymentNode{
		ID:      "gid://shopify/PaymentCustomization/1",
		Title:   "Updated Payment",
		Enabled: false,
	}

	expectedResponse := map[string]interface{}{
		"paymentCustomizationUpdate": map[string]interface{}{
			"paymentCustomization": map[string]interface{}{
				"id":      updatedPayment.ID,
				"title":   updatedPayment.Title,
				"enabled": updatedPayment.Enabled,
			},
		},
	}

	mockClient.On("exec", ctx, mock.AnythingOfType("string")).Return(expectedResponse, nil)

	result, err := service.Update(ctx, updatedPayment)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, updatedPayment.ID, result.ID)
	assert.Equal(t, updatedPayment.Title, result.Title)
	assert.Equal(t, updatedPayment.Enabled, result.Enabled)

	mockClient.AssertExpectations(t)
}

func TestPaymentService_Delete(t *testing.T) {
	mockClient := new(mockShopifyAdminClient)
	service := &paymentServiceImpl{client: mockClient}

	ctx := context.Background()
	paymentID := "gid://shopify/PaymentCustomization/1"

	expectedResponse := map[string]interface{}{
		"paymentCustomizationDelete": map[string]interface{}{
			"deletedId": paymentID,
		},
	}

	mockClient.On("exec", ctx, mock.AnythingOfType("string")).Return(expectedResponse, nil)

	deletedPayment, err := service.Delete(ctx, paymentID)

	assert.NoError(t, err)
	assert.NotNil(t, deletedPayment)
	assert.Equal(t, paymentID, deletedPayment.ID)

	mockClient.AssertExpectations(t)
}

func TestPaymentService_GetError(t *testing.T) {
	mockClient := new(mockShopifyAdminClient)
	service := &paymentServiceImpl{client: mockClient}

	ctx := context.Background()
	paymentID := "gid://shopify/PaymentCustomization/1"

	mockClient.On("exec", ctx, mock.AnythingOfType("string")).Return(nil, assert.AnError)

	payment, err := service.Get(ctx, paymentID)

	assert.Error(t, err)
	assert.Nil(t, payment)
	assert.Equal(t, assert.AnError, err)

	mockClient.AssertExpectations(t)
}

func TestPaymentService_CreateError(t *testing.T) {
	mockClient := new(mockShopifyAdminClient)
	service := &paymentServiceImpl{client: mockClient}

	ctx := context.Background()
	functionID := "gid://shopify/ShopifyFunction/1"
	newPayment := &PaymentNode{
		Title:   "New Payment",
		Enabled: true,
	}

	mockClient.On("exec", ctx, mock.AnythingOfType("string")).Return(nil, assert.AnError)

	createdPayment, err := service.Create(ctx, functionID, newPayment)

	assert.Error(t, err)
	assert.Nil(t, createdPayment)
	assert.Equal(t, assert.AnError, err)

	mockClient.AssertExpectations(t)
}

func TestPaymentService_UpdateError(t *testing.T) {
	mockClient := new(mockShopifyAdminClient)
	service := &paymentServiceImpl{client: mockClient}

	ctx := context.Background()
	updatedPayment := &PaymentNode{
		ID:      "gid://shopify/PaymentCustomization/1",
		Title:   "Updated Payment",
		Enabled: false,
	}

	mockClient.On("exec", ctx, mock.AnythingOfType("string")).Return(nil, assert.AnError)

	result, err := service.Update(ctx, updatedPayment)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, assert.AnError, err)

	mockClient.AssertExpectations(t)
}

func TestPaymentService_DeleteError(t *testing.T) {
	mockClient := new(mockShopifyAdminClient)
	service := &paymentServiceImpl{client: mockClient}

	ctx := context.Background()
	paymentID := "gid://shopify/PaymentCustomization/1"

	mockClient.On("exec", ctx, mock.AnythingOfType("string")).Return(nil, assert.AnError)

	deletedPayment, err := service.Delete(ctx, paymentID)

	assert.Error(t, err)
	assert.Nil(t, deletedPayment)
	assert.Equal(t, assert.AnError, err)

	mockClient.AssertExpectations(t)
}
