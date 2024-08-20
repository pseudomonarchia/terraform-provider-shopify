package shopify

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
)

var _ pubsubWebhookService = (*pubsubWebhookServiceImpl)(nil)

type pubsubWebhookService interface {
	Create(ctx context.Context, webhook *PubsubWebhook) (*PubsubWebhook, error)
	Get(ctx context.Context, id string) (*PubsubWebhook, error)
	Update(ctx context.Context, webhook *PubsubWebhook) (*PubsubWebhook, error)
	Delete(ctx context.Context, id string) error
}

type pubsubWebhookServiceImpl struct {
	client shopifyAdminClient
}

type PubsubWebhook struct {
	ID            string
	Topic         string
	Format        string
	PubSubProject string
	PubSubTopic   string
}

func (p *pubsubWebhookServiceImpl) Create(
	ctx context.Context,
	webhook *PubsubWebhook,
) (*PubsubWebhook, error) {
	gql := `
		mutation pubSubWebhookSubscriptionCreate {
			pubSubWebhookSubscriptionCreate(
				topic: %s
				webhookSubscription: {
					pubSubProject: "%s"
					pubSubTopic: "%s"
					format: %s
				}
		) {
				webhookSubscription {
					id
					topic
					format
					endpoint {
						... on WebhookPubSubEndpoint {
							pubSubProject
							pubSubTopic
						}
					}
				}
			}
		}
	`

	gql = fmt.Sprintf(
		gql,
		webhook.Topic,
		webhook.PubSubProject,
		webhook.PubSubTopic,
		webhook.Format,
	)

	r, err := p.client.exec(ctx, gql)
	if err != nil {
		return nil, err
	}

	jsonb, _ := json.Marshal(r)
	json := gjson.Parse(string(jsonb)).Get("pubSubWebhookSubscriptionCreate.webhookSubscription")

	hook := &PubsubWebhook{
		ID:            json.Get("id").String(),
		Topic:         json.Get("topic").String(),
		Format:        json.Get("format").String(),
		PubSubProject: json.Get("endpoint.pubSubProject").String(),
		PubSubTopic:   json.Get("endpoint.pubSubTopic").String(),
	}

	return hook, nil
}

func (p *pubsubWebhookServiceImpl) Get(
	ctx context.Context,
	id string,
) (*PubsubWebhook, error) {
	gql := `
		query {
			webhookSubscription(id: "%s") {
				id
				topic
				format
				endpoint {
					... on WebhookHttpEndpoint {
						callbackUrl
					}
					... on WebhookEventBridgeEndpoint {
						arn
					}
					... on WebhookPubSubEndpoint {
						pubSubProject
						pubSubTopic
					}
				}
			}
		}
	`

	gql = fmt.Sprintf(gql, id)
	r, err := p.client.exec(ctx, gql)
	if err != nil {
		return nil, err
	}

	jsonb, _ := json.Marshal(r)
	json := gjson.Parse(string(jsonb)).Get("webhookSubscription")

	hook := &PubsubWebhook{
		ID:            json.Get("id").String(),
		Topic:         json.Get("topic").String(),
		Format:        json.Get("format").String(),
		PubSubProject: json.Get("endpoint.pubSubProject").String(),
		PubSubTopic:   json.Get("endpoint.pubSubTopic").String(),
	}

	return hook, nil
}

func (p *pubsubWebhookServiceImpl) Update(
	ctx context.Context,
	webhook *PubsubWebhook,
) (*PubsubWebhook, error) {
	gql := `
		mutation {
		  pubSubWebhookSubscriptionUpdate(
		    id: "%s",
		    webhookSubscription: {
		      pubSubProject: "%s"
		      pubSubTopic: "%s"
		      format: %s
		    }
		) {
		    webhookSubscription {
		      id
		      topic
		      format
		      endpoint {
		        ... on WebhookPubSubEndpoint {
		          pubSubProject
		          pubSubTopic
		        }
		      }
		    }
		  }
		}
	`

	gql = fmt.Sprintf(
		gql,
		webhook.ID,
		webhook.PubSubProject,
		webhook.PubSubTopic,
		webhook.Format,
	)

	r, err := p.client.exec(ctx, gql)
	if err != nil {
		return nil, err
	}

	jsonb, _ := json.Marshal(r)
	json := gjson.Parse(string(jsonb)).Get("pubSubWebhookSubscriptionUpdate.webhookSubscription")

	hook := &PubsubWebhook{
		ID:            json.Get("id").String(),
		Topic:         json.Get("topic").String(),
		Format:        json.Get("format").String(),
		PubSubProject: json.Get("endpoint.pubSubProject").String(),
		PubSubTopic:   json.Get("endpoint.pubSubTopic").String(),
	}

	return hook, nil
}

func (p *pubsubWebhookServiceImpl) Delete(
	ctx context.Context,
	id string,
) error {
	gql := `
		mutation {
			webhookSubscriptionDelete(id: "%s") {
				deletedWebhookSubscriptionId
			}
		}
	`

	gql = fmt.Sprintf(gql, id)
	_, err := p.client.exec(ctx, gql)
	return err
}
