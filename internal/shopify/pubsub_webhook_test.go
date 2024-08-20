package shopify

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPubsubWebhookService_Create(t *testing.T) {
	mockClient := new(mockShopifyAdminClient)
	service := &pubsubWebhookServiceImpl{client: mockClient}

	ctx := context.Background()
	webhook := &PubsubWebhook{
		Topic:         "ORDERS_CREATE",
		Format:        "JSON",
		PubSubProject: "test-project",
		PubSubTopic:   "test-topic",
	}

	expectedResponse := map[string]interface{}{
		"pubSubWebhookSubscriptionCreate": map[string]interface{}{
			"webhookSubscription": map[string]interface{}{
				"id":     "gid://shopify/WebhookSubscription/1",
				"topic":  "ORDERS_CREATE",
				"format": "JSON",
				"endpoint": map[string]interface{}{
					"pubSubProject": "test-project",
					"pubSubTopic":   "test-topic",
				},
			},
		},
	}

	mockClient.On("exec", ctx, mock.AnythingOfType("string")).Return(expectedResponse, nil)

	createdWebhook, err := service.Create(ctx, webhook)

	assert.NoError(t, err)
	assert.NotNil(t, createdWebhook)
	assert.Equal(t, "gid://shopify/WebhookSubscription/1", createdWebhook.ID)
	assert.Equal(t, webhook.Topic, createdWebhook.Topic)
	assert.Equal(t, webhook.Format, createdWebhook.Format)
	assert.Equal(t, webhook.PubSubProject, createdWebhook.PubSubProject)
	assert.Equal(t, webhook.PubSubTopic, createdWebhook.PubSubTopic)

	mockClient.AssertExpectations(t)
}

func TestPubsubWebhookService_Get(t *testing.T) {
	mockClient := new(mockShopifyAdminClient)
	service := &pubsubWebhookServiceImpl{client: mockClient}

	ctx := context.Background()
	webhookID := "gid://shopify/WebhookSubscription/1"

	expectedResponse := map[string]interface{}{
		"webhookSubscription": map[string]interface{}{
			"id":     webhookID,
			"topic":  "ORDERS_CREATE",
			"format": "JSON",
			"endpoint": map[string]interface{}{
				"pubSubProject": "test-project",
				"pubSubTopic":   "test-topic",
			},
		},
	}

	mockClient.On("exec", ctx, mock.AnythingOfType("string")).Return(expectedResponse, nil)

	webhook, err := service.Get(ctx, webhookID)

	assert.NoError(t, err)
	assert.NotNil(t, webhook)
	assert.Equal(t, webhookID, webhook.ID)
	assert.Equal(t, "ORDERS_CREATE", webhook.Topic)
	assert.Equal(t, "JSON", webhook.Format)
	assert.Equal(t, "test-project", webhook.PubSubProject)
	assert.Equal(t, "test-topic", webhook.PubSubTopic)

	mockClient.AssertExpectations(t)
}

func TestPubsubWebhookService_Update(t *testing.T) {
	mockClient := new(mockShopifyAdminClient)
	service := &pubsubWebhookServiceImpl{client: mockClient}

	ctx := context.Background()
	webhook := &PubsubWebhook{
		ID:            "gid://shopify/WebhookSubscription/1",
		Topic:         "ORDERS_CREATE",
		Format:        "JSON",
		PubSubProject: "updated-project",
		PubSubTopic:   "updated-topic",
	}

	expectedResponse := map[string]interface{}{
		"pubSubWebhookSubscriptionUpdate": map[string]interface{}{
			"webhookSubscription": map[string]interface{}{
				"id":     webhook.ID,
				"topic":  webhook.Topic,
				"format": webhook.Format,
				"endpoint": map[string]interface{}{
					"pubSubProject": webhook.PubSubProject,
					"pubSubTopic":   webhook.PubSubTopic,
				},
			},
		},
	}

	mockClient.On("exec", ctx, mock.AnythingOfType("string")).Return(expectedResponse, nil)

	updatedWebhook, err := service.Update(ctx, webhook)

	assert.NoError(t, err)
	assert.NotNil(t, updatedWebhook)
	assert.Equal(t, webhook.ID, updatedWebhook.ID)
	assert.Equal(t, webhook.Topic, updatedWebhook.Topic)
	assert.Equal(t, webhook.Format, updatedWebhook.Format)
	assert.Equal(t, webhook.PubSubProject, updatedWebhook.PubSubProject)
	assert.Equal(t, webhook.PubSubTopic, updatedWebhook.PubSubTopic)

	mockClient.AssertExpectations(t)
}

func TestPubsubWebhookService_Delete(t *testing.T) {
	mockClient := new(mockShopifyAdminClient)
	service := &pubsubWebhookServiceImpl{client: mockClient}

	ctx := context.Background()
	webhookID := "gid://shopify/WebhookSubscription/1"

	expectedResponse := map[string]interface{}{
		"webhookSubscriptionDelete": map[string]interface{}{
			"deletedWebhookSubscriptionId": webhookID,
		},
	}

	mockClient.On("exec", ctx, mock.AnythingOfType("string")).Return(expectedResponse, nil)

	err := service.Delete(ctx, webhookID)

	assert.NoError(t, err)

	mockClient.AssertExpectations(t)
}

func TestPubsubWebhookService_CreateError(t *testing.T) {
	mockClient := new(mockShopifyAdminClient)
	service := &pubsubWebhookServiceImpl{client: mockClient}

	ctx := context.Background()
	webhook := &PubsubWebhook{
		Topic:         "ORDERS_CREATE",
		Format:        "JSON",
		PubSubProject: "test-project",
		PubSubTopic:   "test-topic",
	}

	mockClient.On("exec", ctx, mock.AnythingOfType("string")).Return(nil, assert.AnError)

	createdWebhook, err := service.Create(ctx, webhook)

	assert.Error(t, err)
	assert.Nil(t, createdWebhook)
	assert.Equal(t, assert.AnError, err)

	mockClient.AssertExpectations(t)
}

func TestPubsubWebhookService_GetError(t *testing.T) {
	mockClient := new(mockShopifyAdminClient)
	service := &pubsubWebhookServiceImpl{client: mockClient}

	ctx := context.Background()
	webhookID := "gid://shopify/WebhookSubscription/1"

	mockClient.On("exec", ctx, mock.AnythingOfType("string")).Return(nil, assert.AnError)

	webhook, err := service.Get(ctx, webhookID)

	assert.Error(t, err)
	assert.Nil(t, webhook)
	assert.Equal(t, assert.AnError, err)

	mockClient.AssertExpectations(t)
}

func TestPubsubWebhookService_UpdateError(t *testing.T) {
	mockClient := new(mockShopifyAdminClient)
	service := &pubsubWebhookServiceImpl{client: mockClient}

	ctx := context.Background()
	webhook := &PubsubWebhook{
		ID:            "gid://shopify/WebhookSubscription/1",
		Topic:         "ORDERS_CREATE",
		Format:        "JSON",
		PubSubProject: "updated-project",
		PubSubTopic:   "updated-topic",
	}

	mockClient.On("exec", ctx, mock.AnythingOfType("string")).Return(nil, assert.AnError)

	updatedWebhook, err := service.Update(ctx, webhook)

	assert.Error(t, err)
	assert.Nil(t, updatedWebhook)
	assert.Equal(t, assert.AnError, err)

	mockClient.AssertExpectations(t)
}

func TestPubsubWebhookService_DeleteError(t *testing.T) {
	mockClient := new(mockShopifyAdminClient)
	service := &pubsubWebhookServiceImpl{client: mockClient}

	ctx := context.Background()
	webhookID := "gid://shopify/WebhookSubscription/1"

	mockClient.On("exec", ctx, mock.AnythingOfType("string")).Return(nil, assert.AnError)

	err := service.Delete(ctx, webhookID)

	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)

	mockClient.AssertExpectations(t)
}
