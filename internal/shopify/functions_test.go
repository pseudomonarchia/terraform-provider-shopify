package shopify

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFunctionService_List(t *testing.T) {
	mockClient := new(mockShopifyAdminClient)
	service := &FunctionServiceImpl{client: mockClient}

	ctx := context.Background()

	expectedResponse := map[string]interface{}{
		"shopifyFunctions": map[string]interface{}{
			"nodes": []interface{}{
				map[string]interface{}{
					"id":      "gid://shopify/ShopifyFunction/1",
					"title":   "Function 1",
					"apiType": "DISCOUNTS",
					"app": map[string]interface{}{
						"title": "App 1",
					},
				},
				map[string]interface{}{
					"id":      "gid://shopify/ShopifyFunction/2",
					"title":   "Function 2",
					"apiType": "SHIPPING",
					"app": map[string]interface{}{
						"title": "App 2",
					},
				},
			},
		},
	}

	mockClient.On("exec", ctx, mock.AnythingOfType("string")).Return(expectedResponse, nil)

	functionNodes, err := service.List(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, functionNodes)
	assert.Len(t, functionNodes.Nodes, 2)

	// Verify the first function node
	assert.Equal(t, "gid://shopify/ShopifyFunction/1", functionNodes.Nodes[0].ID)
	assert.Equal(t, "Function 1", functionNodes.Nodes[0].Title)
	assert.Equal(t, "DISCOUNTS", functionNodes.Nodes[0].APIType)
	assert.Equal(t, "App 1", functionNodes.Nodes[0].APPName)

	// Verify the second function node
	assert.Equal(t, "gid://shopify/ShopifyFunction/2", functionNodes.Nodes[1].ID)
	assert.Equal(t, "Function 2", functionNodes.Nodes[1].Title)
	assert.Equal(t, "SHIPPING", functionNodes.Nodes[1].APIType)
	assert.Equal(t, "App 2", functionNodes.Nodes[1].APPName)

	mockClient.AssertExpectations(t)
}

func TestFunctionService_ListEmpty(t *testing.T) {
	mockClient := new(mockShopifyAdminClient)
	service := &FunctionServiceImpl{client: mockClient}

	ctx := context.Background()

	expectedResponse := map[string]interface{}{
		"shopifyFunctions": map[string]interface{}{
			"nodes": []interface{}{},
		},
	}

	mockClient.On("exec", ctx, mock.AnythingOfType("string")).Return(expectedResponse, nil)

	functionNodes, err := service.List(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, functionNodes)
	assert.Len(t, functionNodes.Nodes, 0)

	mockClient.AssertExpectations(t)
}

func TestFunctionService_ListError(t *testing.T) {
	mockClient := new(mockShopifyAdminClient)
	service := &FunctionServiceImpl{client: mockClient}

	ctx := context.Background()

	mockClient.On("exec", ctx, mock.AnythingOfType("string")).Return(nil, assert.AnError)

	functionNodes, err := service.List(ctx)

	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
	assert.Empty(t, functionNodes.Nodes)

	mockClient.AssertExpectations(t)
}
