package shopify

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockShopifyAdminClient struct {
	mock.Mock
}

func (m *mockShopifyAdminClient) exec(ctx context.Context, query string) (interface{}, error) {
	args := m.Called(ctx, query)
	return args.Get(0), args.Error(1)
}

func TestNew(t *testing.T) {
	client := New("example.myshopify.com", "access_token", "2023-04")

	assert.NotNil(t, client)
	assert.Equal(t, "example.myshopify.com", client.storeDomain)
	assert.Equal(t, "access_token", client.storeAccessToken)
	assert.Equal(t, "2023-04", client.storeApiVersion)
	assert.NotNil(t, client.Discount)
	assert.NotNil(t, client.Function)
	assert.NotNil(t, client.Payment)
	assert.NotNil(t, client.Delivery)
}

func TestExec(t *testing.T) {
	os.Setenv("SHOPIFY_TEST_MODE", "true")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/admin/api/2023-04/graphql.json", r.URL.Path)
		assert.Equal(t, "access_token", r.Header.Get("X-Shopify-Access-Token"))
		assert.Equal(t, "no-cache", r.Header.Get("Cache-Control"))

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data": {"test": "success"}}`))
	}))

	defer server.Close()
	client := &ShopifyAdminClinetImpl{
		storeDomain:      server.URL[7:], // 移除 "http://"
		storeAccessToken: "access_token",
		storeApiVersion:  "2023-04",
		local:            true,
	}

	ctx := context.Background()
	query := `query { test }`

	result, err := client.exec(ctx, query)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	expectedResult := map[string]interface{}{
		"test": "success",
	}

	assert.Equal(t, expectedResult, result)
	os.Unsetenv("SHOPIFY_TEST_MODE")
}
