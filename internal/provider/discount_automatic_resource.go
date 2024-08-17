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

var _ resource.Resource = (*discountAutomaticResource)(nil)

type discountAutomaticResource struct {
	client *shopify.ShopifyAdminClinetImpl
}

type discountAutomaticResourceModel struct {
	FunctionID   types.String                                `tfsdk:"function_id"`
	ID           types.String                                `tfsdk:"id"`
	Title        types.String                                `tfsdk:"title"`
	StartsAt     types.String                                `tfsdk:"starts_at"`
	EndsAt       types.String                                `tfsdk:"ends_at"`
	CombinesWith *discountAutomaticCombinesWithResourceModel `tfsdk:"combines_with"`
}

type discountAutomaticCombinesWithResourceModel struct {
	OrderDiscounts    types.Bool `tfsdk:"order_discounts"`
	ProductDiscounts  types.Bool `tfsdk:"product_discounts"`
	ShippingDiscounts types.Bool `tfsdk:"shipping_discounts"`
}

func NewDiscountAutomaticResource() resource.Resource {
	return &discountAutomaticResource{}
}

func (r *discountAutomaticResource) Metadata(
	c context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_discount"
}

func (r *discountAutomaticResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Shopify Function Discount Automatic Resource",
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
			"starts_at": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
						"Must be a valid datetime format (example: 2024-01-01T00:00:00Z)",
					),
				},
			},
			"ends_at": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
						"Must be a valid datetime format (example: 2024-01-01T00:00:00Z)",
					),
				},
			},
			"combines_with": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"order_discounts": schema.BoolAttribute{
						Required: true,
					},
					"product_discounts": schema.BoolAttribute{
						Required: true,
					},
					"shipping_discounts": schema.BoolAttribute{
						Required: true,
					},
				},
			},
		},
	}
}

func (r *discountAutomaticResource) Configure(
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

func (r *discountAutomaticResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var data discountAutomaticResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dn := &shopify.DiscountNode{
		Title:    data.Title.ValueString(),
		StartsAt: data.StartsAt.ValueString(),
		CombinesWith: &shopify.DiscountCombinesWith{
			OrderDiscounts:    data.CombinesWith.OrderDiscounts.ValueBool(),
			ProductDiscounts:  data.CombinesWith.ProductDiscounts.ValueBool(),
			ShippingDiscounts: data.CombinesWith.ShippingDiscounts.ValueBool(),
		},
	}

	if !data.EndsAt.IsNull() {
		dn.EndsAt = data.EndsAt.ValueString()
	}

	q, err := r.client.Discount.Create(ctx, data.FunctionID.ValueString(), dn)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create shopify discount automatic", err.Error())
		return
	}

	data.FunctionID = types.StringValue(data.FunctionID.ValueString())
	data.ID = types.StringValue(q.ID)
	data.Title = types.StringValue(q.Title)
	data.StartsAt = types.StringValue(q.StartsAt)
	data.CombinesWith = &discountAutomaticCombinesWithResourceModel{
		OrderDiscounts:    types.BoolValue(q.CombinesWith.OrderDiscounts),
		ProductDiscounts:  types.BoolValue(q.CombinesWith.ProductDiscounts),
		ShippingDiscounts: types.BoolValue(q.CombinesWith.ShippingDiscounts),
	}

	if q.EndsAt != "" {
		data.EndsAt = types.StringValue(q.EndsAt)
	} else {
		data.EndsAt = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *discountAutomaticResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data discountAutomaticResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	q, err := r.client.Discount.Get(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get shopify discount automatic", err.Error())
	}

	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = types.StringValue(q.ID)
	data.Title = types.StringValue(q.Title)
	data.StartsAt = types.StringValue(q.StartsAt)
	data.CombinesWith = &discountAutomaticCombinesWithResourceModel{
		OrderDiscounts:    types.BoolValue(q.CombinesWith.OrderDiscounts),
		ProductDiscounts:  types.BoolValue(q.CombinesWith.ProductDiscounts),
		ShippingDiscounts: types.BoolValue(q.CombinesWith.ShippingDiscounts),
	}

	if q.EndsAt != "" {
		data.EndsAt = types.StringValue(q.EndsAt)
	} else {
		data.EndsAt = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *discountAutomaticResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var data discountAutomaticResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dn := &shopify.DiscountNode{
		ID:       data.ID.ValueString(),
		Title:    data.Title.ValueString(),
		StartsAt: data.StartsAt.ValueString(),
		CombinesWith: &shopify.DiscountCombinesWith{
			OrderDiscounts:    data.CombinesWith.OrderDiscounts.ValueBool(),
			ProductDiscounts:  data.CombinesWith.ProductDiscounts.ValueBool(),
			ShippingDiscounts: data.CombinesWith.ShippingDiscounts.ValueBool(),
		},
	}

	if !data.EndsAt.IsNull() {
		dn.EndsAt = data.EndsAt.ValueString()
	}

	q, err := r.client.Discount.Update(ctx, dn)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update shopify discount automatic", err.Error())
		return
	}

	data.ID = types.StringValue(q.ID)
	data.Title = types.StringValue(q.Title)
	data.StartsAt = types.StringValue(q.StartsAt)
	data.CombinesWith = &discountAutomaticCombinesWithResourceModel{
		OrderDiscounts:    types.BoolValue(q.CombinesWith.OrderDiscounts),
		ProductDiscounts:  types.BoolValue(q.CombinesWith.ProductDiscounts),
		ShippingDiscounts: types.BoolValue(q.CombinesWith.ShippingDiscounts),
	}

	if q.EndsAt != "" {
		data.EndsAt = types.StringValue(q.EndsAt)
	} else {
		data.EndsAt = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *discountAutomaticResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data discountAutomaticResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Discount.Delete(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete shopify discount automatic", err.Error())
		return
	}
}

func (r *discountAutomaticResource) ImportState(
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
