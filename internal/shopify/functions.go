package shopify

import (
	"context"
	"encoding/json"

	"github.com/tidwall/gjson"
)

var _ FunctionService = (*FunctionServiceImpl)(nil)

type FunctionService interface {
	List(ctx context.Context) (FunctionNodes, error)
}

type FunctionServiceImpl struct {
	client shopifyAdminClient
}

type FunctionNode struct {
	ID      string
	Title   string
	APIType string
	APPName string
}

type FunctionNodes struct {
	Nodes []FunctionNode
}

func (f *FunctionServiceImpl) List(ctx context.Context) (FunctionNodes, error) {
	gql := `
		query {
			shopifyFunctions(first: 250) {
				nodes {
					id
					title
					apiType
					app {
						title
					}
				}
			}
		}
	`

	r, err := f.client.exec(ctx, gql)
	var functionNodes FunctionNodes
	if err != nil {
		return functionNodes, err
	}

	jsonb, _ := json.Marshal(r)
	gjson.
		Parse(string(jsonb)).
		Get("shopifyFunctions.nodes").
		ForEach(func(_, value gjson.Result) bool {
			functionNode := FunctionNode{
				ID:      value.Get("id").String(),
				Title:   value.Get("title").String(),
				APIType: value.Get("apiType").String(),
				APPName: value.Get("app.title").String(),
			}

			functionNodes.Nodes = append(functionNodes.Nodes, functionNode)
			return true
		})

	return functionNodes, nil
}
