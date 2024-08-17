package shopify

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeliveryServiceImpl_Get(t *testing.T) {
	mockClient := new(mockShopifyAdminClient)
	service := &deliveryServiceImpl{client: mockClient}

	mockResponse := map[string]interface{}{
		"deliveryCustomization": map[string]interface{}{
			"id":      "gid://shopify/DeliveryCustomization/1",
			"title":   "Test Delivery",
			"enabled": true,
		},
	}

	mockClient.On("exec", mock.Anything, mock.AnythingOfType("string")).Return(mockResponse, nil)

	result, err := service.Get(context.Background(), "gid://shopify/DeliveryCustomization/1")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "gid://shopify/DeliveryCustomization/1", result.ID)
	assert.Equal(t, "Test Delivery", result.Title)
	assert.True(t, result.Enabled)

	mockClient.AssertExpectations(t)
}

func TestDeliveryServiceImpl_Create(t *testing.T) {
	mockClient := new(mockShopifyAdminClient)
	service := &deliveryServiceImpl{client: mockClient}

	mockResponse := map[string]interface{}{
		"deliveryCustomizationCreate": map[string]interface{}{
			"deliveryCustomization": map[string]interface{}{
				"id":      "gid://shopify/DeliveryCustomization/2",
				"title":   "New Delivery",
				"enabled": false,
			},
		},
	}

	mockClient.On("exec", mock.Anything, mock.AnythingOfType("string")).Return(mockResponse, nil)

	newDelivery := &DeliveryNode{
		Title:   "New Delivery",
		Enabled: false,
	}
	result, err := service.Create(context.Background(), "gid://shopify/Function/1", newDelivery)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "gid://shopify/DeliveryCustomization/2", result.ID)
	assert.Equal(t, "New Delivery", result.Title)
	assert.False(t, result.Enabled)

	mockClient.AssertExpectations(t)
}

func TestDeliveryServiceImpl_Update(t *testing.T) {
	mockClient := new(mockShopifyAdminClient)
	service := &deliveryServiceImpl{client: mockClient}

	mockResponse := map[string]interface{}{
		"deliveryCustomizationUpdate": map[string]interface{}{
			"deliveryCustomization": map[string]interface{}{
				"id":      "gid://shopify/DeliveryCustomization/3",
				"title":   "Updated Delivery",
				"enabled": true,
			},
		},
	}

	mockClient.On("exec", mock.Anything, mock.AnythingOfType("string")).Return(mockResponse, nil)

	updatedDelivery := &DeliveryNode{
		ID:      "gid://shopify/DeliveryCustomization/3",
		Title:   "Updated Delivery",
		Enabled: true,
	}
	result, err := service.Update(context.Background(), updatedDelivery)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "gid://shopify/DeliveryCustomization/3", result.ID)
	assert.Equal(t, "Updated Delivery", result.Title)
	assert.True(t, result.Enabled)

	mockClient.AssertExpectations(t)
}

func TestDeliveryServiceImpl_Delete(t *testing.T) {
	mockClient := new(mockShopifyAdminClient)
	service := &deliveryServiceImpl{client: mockClient}

	mockResponse := map[string]interface{}{
		"deliveryCustomizationDelete": map[string]interface{}{
			"deletedId": "gid://shopify/DeliveryCustomization/4",
		},
	}

	mockClient.On("exec", mock.Anything, mock.AnythingOfType("string")).Return(mockResponse, nil)

	result, err := service.Delete(context.Background(), "gid://shopify/DeliveryCustomization/4")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "gid://shopify/DeliveryCustomization/4", result.ID)

	mockClient.AssertExpectations(t)
}

func TestDeliveryServiceImpl_ErrorHandling(t *testing.T) {
	mockClient := new(mockShopifyAdminClient)
	service := &deliveryServiceImpl{client: mockClient}

	mockError := assert.AnError

	mockClient.On("exec", mock.Anything, mock.AnythingOfType("string")).Return(nil, mockError)

	_, err := service.Get(context.Background(), "gid://shopify/DeliveryCustomization/1")
	assert.Error(t, err)

	_, err = service.Create(context.Background(), "gid://shopify/Function/1", &DeliveryNode{})
	assert.Error(t, err)

	_, err = service.Update(context.Background(), &DeliveryNode{})
	assert.Error(t, err)

	_, err = service.Delete(context.Background(), "gid://shopify/DeliveryCustomization/1")
	assert.Error(t, err)

	mockClient.AssertExpectations(t)
}
