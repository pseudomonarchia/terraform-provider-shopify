package provider

import (
	"context"
	"fmt"

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

var _ resource.Resource = (*pubsubWebhookResource)(nil)

type pubsubWebhookResource struct {
	client *shopify.ShopifyAdminClinetImpl
}

type pubsubWebhookResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Topic         types.String `tfsdk:"topic"`
	Format        types.String `tfsdk:"format"`
	PubSubProject types.String `tfsdk:"pubsub_project"`
	PubSubTopic   types.String `tfsdk:"pubsub_topic"`
}

func NewPubsubWebhookResource() resource.Resource {
	return &pubsubWebhookResource{}
}

func (r *pubsubWebhookResource) Metadata(
	c context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_pubsub_webhook"
}

func (r *pubsubWebhookResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Shopify PubSub Webhook Resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"topic": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"format": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("JSON", "XML"),
				},
			},
			"pubsub_project": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"pubsub_topic": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	}
}

func (r *pubsubWebhookResource) Configure(
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

func (r *pubsubWebhookResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var data pubsubWebhookResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	webhook := &shopify.PubsubWebhook{
		Topic:         data.Topic.ValueString(),
		Format:        data.Format.ValueString(),
		PubSubProject: data.PubSubProject.ValueString(),
		PubSubTopic:   data.PubSubTopic.ValueString(),
	}

	createdWebhook, err := r.client.PubsubWebhook.Create(ctx, webhook)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create shopify pubsub webhook", err.Error())
		return
	}

	data.ID = types.StringValue(createdWebhook.ID)
	data.Topic = types.StringValue(createdWebhook.Topic)
	data.Format = types.StringValue(createdWebhook.Format)
	data.PubSubProject = types.StringValue(createdWebhook.PubSubProject)
	data.PubSubTopic = types.StringValue(createdWebhook.PubSubTopic)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *pubsubWebhookResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data pubsubWebhookResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	webhook, err := r.client.PubsubWebhook.Get(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get shopify pubsub webhook", err.Error())
		return
	}

	data.ID = types.StringValue(webhook.ID)
	data.Topic = types.StringValue(webhook.Topic)
	data.Format = types.StringValue(webhook.Format)
	data.PubSubProject = types.StringValue(webhook.PubSubProject)
	data.PubSubTopic = types.StringValue(webhook.PubSubTopic)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *pubsubWebhookResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var data pubsubWebhookResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	webhook := &shopify.PubsubWebhook{
		ID:            data.ID.ValueString(),
		Topic:         data.Topic.ValueString(),
		Format:        data.Format.ValueString(),
		PubSubProject: data.PubSubProject.ValueString(),
		PubSubTopic:   data.PubSubTopic.ValueString(),
	}

	updatedWebhook, err := r.client.PubsubWebhook.Update(ctx, webhook)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update shopify pubsub webhook", err.Error())
		return
	}

	data.ID = types.StringValue(updatedWebhook.ID)
	data.Topic = types.StringValue(updatedWebhook.Topic)
	data.Format = types.StringValue(updatedWebhook.Format)
	data.PubSubProject = types.StringValue(updatedWebhook.PubSubProject)
	data.PubSubTopic = types.StringValue(updatedWebhook.PubSubTopic)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *pubsubWebhookResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data pubsubWebhookResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.PubsubWebhook.Delete(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete shopify pubsub webhook", err.Error())
		return
	}
}

func (r *pubsubWebhookResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
