package provider

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pseudomonarchia/terraform-provider-shopify/internal/shopify"
)

var _ provider.Provider = (*funcProvider)(nil)

type funcProvider struct {
	version string
}

type funcProviderModel struct {
	StoreDomain      types.String `tfsdk:"store_domain"`
	StoreAccessToken types.String `tfsdk:"store_access_token"`
	StoreApiVersion  types.String `tfsdk:"store_api_version"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &funcProvider{
			version: version,
		}
	}
}

func (p *funcProvider) Metadata(
	_ context.Context,
	_ provider.MetadataRequest,
	resp *provider.MetadataResponse,
) {
	resp.TypeName = "shopify"
	resp.Version = p.version
}

func (p *funcProvider) Schema(
	_ context.Context,
	_ provider.SchemaRequest,
	resp *provider.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Shopify Function Registry",
		Attributes: map[string]schema.Attribute{
			"store_domain": schema.StringAttribute{
				Description: "The store's URL, formatted as <storename>.myshopify.com",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9-]+\.myshopify\.com$`),
						"must be a valid Shopify store domain (e.g., example.myshopify.com)",
					),
				},
			},
			"store_access_token": schema.StringAttribute{
				Description: "The store's access token",
				Required:    true,
				Sensitive:   true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"store_api_version": schema.StringAttribute{
				Description: "The store's API version",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^\d{4}-\d{2}$`),
						"must be a valid Shopify API version (e.g., 2024-07)",
					),
				},
			},
		},
	}
}

func (p *funcProvider) Configure(
	ctx context.Context,
	req provider.ConfigureRequest,
	resp *provider.ConfigureResponse,
) {
	var conf funcProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &conf)...)
	if resp.Diagnostics.HasError() {
		return
	}

	c := shopify.New(
		conf.StoreDomain.ValueString(),
		conf.StoreAccessToken.ValueString(),
		conf.StoreApiVersion.ValueString(),
	)

	resp.ResourceData = c
	resp.DataSourceData = c
}

func (p *funcProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewFunctionDataSource,
	}
}

func (p *funcProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDiscountAutomaticResource,
		NewPaymentCustomResource,
		NewDeliveryCustomResource,
	}
}
