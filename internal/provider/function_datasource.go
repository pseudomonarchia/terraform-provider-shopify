package provider

import (
	"context"
	"fmt"
	"shopify-function-registry/internal/shopify"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &FunctinDataSource{}

type FunctinDataSource struct {
	client *shopify.ShopifyAdminClinetImpl
}

type FunctionDataSourceModel struct {
	ID      types.String `tfsdk:"id"`
	Title   types.String `tfsdk:"title"`
	APIType types.String `tfsdk:"api_type"`
	APPName types.String `tfsdk:"app_title"`
}

func NewFunctionDataSource() datasource.DataSource {
	return &FunctinDataSource{}
}

func (d *FunctinDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_function"
}

func (d *FunctinDataSource) Schema(
	_ context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Shopify Function Data Source",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"title": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"api_type": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"app_title": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	}
}

func (d *FunctinDataSource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*shopify.ShopifyAdminClinetImpl)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf(
				"Expected *shopify.ShopifyAdminClinetImpl, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)

		return
	}

	d.client = c
}

func (d *FunctinDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data FunctionDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	q, err := d.client.Function.List(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get shopify function", err.Error())
		return
	}

	for _, node := range q.Nodes {
		if node.APPName == data.APPName.ValueString() && node.Title == data.Title.ValueString() {
			data.ID = types.StringValue(node.ID)
			break
		}
	}

	if data.ID.IsNull() {
		resp.Diagnostics.AddError(
			"No matching function found",
			fmt.Sprintf("Unable to find function with APP name '%s' and title '%s'",
				data.APPName.ValueString(),
				data.Title.ValueString(),
			),
		)

		return
	}

	tflog.Trace(ctx, "read a shopify function data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
