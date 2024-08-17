package provider

import (
	"context"
	"fmt"
	"regexp"
	"shopify-function-registry/internal/shopify"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = (*paymentCustomResource)(nil)

type paymentCustomResource struct {
	client *shopify.ShopifyAdminClinetImpl
}

type paymentCustomResourceModel struct {
	FunctionID types.String `tfsdk:"function_id"`
	ID         types.String `tfsdk:"id"`
	Title      types.String `tfsdk:"title"`
	Enabled    types.Bool   `tfsdk:"enabled"`
}

func NewPaymentCustomResource() resource.Resource {
	return &paymentCustomResource{}
}

func (r *paymentCustomResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_payment"
}

func (r *paymentCustomResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Shopify Function Payment Customization Resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"function_id": schema.StringAttribute{
				Required: true,
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

func (r *paymentCustomResource) Configure(
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

func (r *paymentCustomResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var data paymentCustomResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pn := &shopify.PaymentNode{
		Title:   data.Title.ValueString(),
		Enabled: data.Enabled.ValueBool(),
	}

	q, err := r.client.Payment.Create(ctx, data.FunctionID.ValueString(), pn)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create shopify payment customization", err.Error())
		return
	}

	data.FunctionID = types.StringValue(data.FunctionID.ValueString())
	data.ID = types.StringValue(q.ID)
	data.Title = types.StringValue(q.Title)
	data.Enabled = types.BoolValue(q.Enabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *paymentCustomResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data paymentCustomResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	q, err := r.client.Payment.Get(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get shopify payment customization", err.Error())
	}

	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = types.StringValue(q.ID)
	data.Title = types.StringValue(q.Title)
	data.Enabled = types.BoolValue(q.Enabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *paymentCustomResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var data paymentCustomResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pn := &shopify.PaymentNode{
		ID:      data.ID.ValueString(),
		Title:   data.Title.ValueString(),
		Enabled: data.Enabled.ValueBool(),
	}

	q, err := r.client.Payment.Update(ctx, pn)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update shopify payment customization", err.Error())
		return
	}

	data.ID = types.StringValue(q.ID)
	data.Title = types.StringValue(q.Title)
	data.Enabled = types.BoolValue(q.Enabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *paymentCustomResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data paymentCustomResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Payment.Delete(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete shopify payment customization", err.Error())
		return
	}
}

func (r *paymentCustomResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import Format",
			"Please use the format 'id,function_id' to import the resource",
		)

		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("function_id"), idParts[1])...)
}
