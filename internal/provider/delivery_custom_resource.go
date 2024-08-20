package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pseudomonarchia/terraform-provider-shopify/internal/shopify"
)

var _ resource.Resource = (*deliveryCustomResource)(nil)

type deliveryCustomResource struct {
	client *shopify.ShopifyAdminClinetImpl
}

type deliveryCustomResourceModel struct {
	FunctionID types.String `tfsdk:"function_id"`
	ID         types.String `tfsdk:"id"`
	Title      types.String `tfsdk:"title"`
	Enabled    types.Bool   `tfsdk:"enabled"`
}

func NewDeliveryCustomResource() resource.Resource {
	return &deliveryCustomResource{}
}

func (r *deliveryCustomResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_delivery"
}

func (r *deliveryCustomResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Shopify Function Delivery Customization Resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"function_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`),
						"Must be a valid UUID",
					),
				},
			},
			"title": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"enabled": schema.BoolAttribute{
				Required: true,
			},
		},
	}
}

func (r *deliveryCustomResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
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

	r.client = c
}

func (r *deliveryCustomResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var data deliveryCustomResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dn := &shopify.DeliveryNode{
		Title:   data.Title.ValueString(),
		Enabled: data.Enabled.ValueBool(),
	}

	q, err := r.client.Delivery.Create(ctx, data.FunctionID.ValueString(), dn)
	if err != nil {
		resp.Diagnostics.AddError("Create Shopify Delivery Custom Failed", err.Error())
		return
	}

	data.FunctionID = types.StringValue(data.FunctionID.ValueString())
	data.ID = types.StringValue(q.ID)
	data.Title = types.StringValue(q.Title)
	data.Enabled = types.BoolValue(q.Enabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *deliveryCustomResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data deliveryCustomResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	q, err := r.client.Delivery.Get(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Get Shopify Delivery Custom Failed", err.Error())
	}

	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = types.StringValue(q.ID)
	data.Title = types.StringValue(q.Title)
	data.Enabled = types.BoolValue(q.Enabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *deliveryCustomResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var data deliveryCustomResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dn := &shopify.DeliveryNode{
		ID:      data.ID.ValueString(),
		Title:   data.Title.ValueString(),
		Enabled: data.Enabled.ValueBool(),
	}

	q, err := r.client.Delivery.Update(ctx, dn)
	if err != nil {
		resp.Diagnostics.AddError("Update Shopify Delivery Custom Failed", err.Error())
		return
	}

	data.ID = types.StringValue(q.ID)
	data.Title = types.StringValue(q.Title)
	data.Enabled = types.BoolValue(q.Enabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *deliveryCustomResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data deliveryCustomResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Delivery.Delete(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete Shopify Delivery Custom Failed", err.Error())
		return
	}
}

func (r *deliveryCustomResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import Format",
			"Please use the format 'id,function_id' to import resources",
		)

		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("function_id"), idParts[1])...)
}
