// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/product"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/subscription"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/workspace"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ApiManagementWorkspaceSubscriptionModel struct {
	SubscriptionName         string `tfschema:"subscription_name"`
	ApiManagementWorkspaceId string `tfschema:"api_management_workspace_id"`
	DisplayName              string `tfschema:"display_name"`
	OwnerId                  string `tfschema:"owner_id"`
	ProductId                string `tfschema:"product_id"`
	ApiId                    string `tfschema:"api_id"`
	TracingEnabled           bool   `tfschema:"tracing_enabled"`
	PrimaryKey               string `tfschema:"primary_key"`
	SecondaryKey             string `tfschema:"secondary_key"`
	State                    string `tfschema:"state"`
}

type ApiManagementWorkspaceSubscriptionResource struct{}

var _ sdk.ResourceWithUpdate = ApiManagementWorkspaceSubscriptionResource{}

func (r ApiManagementWorkspaceSubscriptionResource) ResourceType() string {
	return "azurerm_api_management_workspace_subscription"
}

func (r ApiManagementWorkspaceSubscriptionResource) ModelObject() interface{} {
	return &ApiManagementWorkspaceSubscriptionModel{}
}

func (r ApiManagementWorkspaceSubscriptionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return subscription.ValidateWorkspaceSubscriptions2ID
}

func (r ApiManagementWorkspaceSubscriptionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"subscription_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validate.ApiManagementChildName,
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"api_management_workspace_id": commonschema.ResourceIDReferenceRequiredForceNew(&workspace.WorkspaceId{}),

		"api_id": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"product_id"},
			ValidateFunc:  api.ValidateApiID,
		},

		"owner_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"product_id": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"api_id"},
			ValidateFunc:  product.ValidateProductID,
		},

		"tracing_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"primary_key": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"secondary_key": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"state": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(subscription.SubscriptionStateSubmitted),
			ValidateFunc: validation.StringInSlice([]string{
				string(subscription.SubscriptionStateActive),
				string(subscription.SubscriptionStateCancelled),
				string(subscription.SubscriptionStateExpired),
				string(subscription.SubscriptionStateRejected),
				string(subscription.SubscriptionStateSubmitted),
				string(subscription.SubscriptionStateSuspended),
			}, false),
		},
	}
}

func (r ApiManagementWorkspaceSubscriptionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ApiManagementWorkspaceSubscriptionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.SubscriptionClient_v2024_05_01

			var model ApiManagementWorkspaceSubscriptionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			workspaceId, err := workspace.ParseWorkspaceID(model.ApiManagementWorkspaceId)
			if err != nil {
				return err
			}

			id := subscription.NewWorkspaceSubscriptions2ID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.ServiceName, workspaceId.WorkspaceId, model.SubscriptionName)

			existing, err := client.WorkspaceSubscriptionGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			var scope string
			switch {
			case model.ProductId != "":
				scope = model.ProductId
			case model.ApiId != "":
				scope = model.ApiId
			default:
				scope = workspaceId.ID() + "/apis"
			}

			payload := subscription.SubscriptionCreateParameters{
				Properties: &subscription.SubscriptionCreateParameterProperties{
					DisplayName:  model.DisplayName,
					Scope:        scope,
					State:        pointer.To(subscription.SubscriptionState(model.State)),
					AllowTracing: pointer.To(model.TracingEnabled),
				},
			}

			if model.OwnerId != "" {
				payload.Properties.OwnerId = pointer.To(model.OwnerId)
			}

			if model.PrimaryKey != "" {
				payload.Properties.PrimaryKey = pointer.To(model.PrimaryKey)
			}

			if model.SecondaryKey != "" {
				payload.Properties.SecondaryKey = pointer.To(model.SecondaryKey)
			}

			if _, err := client.WorkspaceSubscriptionCreateOrUpdate(ctx, id, payload, subscription.WorkspaceSubscriptionCreateOrUpdateOperationOptions{}); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ApiManagementWorkspaceSubscriptionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.SubscriptionClient_v2024_05_01

			id, err := subscription.ParseWorkspaceSubscriptions2ID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.WorkspaceSubscriptionGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := ApiManagementWorkspaceSubscriptionModel{
				SubscriptionName:         id.SubscriptionName,
				ApiManagementWorkspaceId: workspace.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceId).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.DisplayName = pointer.From(props.DisplayName)
					state.TracingEnabled = pointer.From(props.AllowTracing)
					state.State = string(props.State)
					state.OwnerId = pointer.From(props.OwnerId)

					productId := ""
					apiId := ""
					if props.Scope != "" && !strings.HasSuffix(props.Scope, "/apis") {
						// the scope is either a product or api id
						parseId, err := product.ParseProductIDInsensitively(props.Scope)
						if err == nil {
							productId = parseId.ID()
						} else {
							parsedApiId, err := api.ParseApiIDInsensitively(props.Scope)
							if err == nil {
								apiId = parsedApiId.ID()
							}
						}
						if id.SubscriptionName == "master" && productId == "" && apiId == "" {
							// Built-in "master" subscription has the API Management service ID as its scope and is unmodifiable, so it should not be managed by the AzureRM Provider.
							return fmt.Errorf("built-in subscription is system-generated and cannot be managed with AzureRM Provider")
						}
					}

					state.ProductId = productId
					state.ApiId = apiId
				}
			}

			// Primary and secondary keys must be got from this additional api
			keyResp, err := client.WorkspaceSubscriptionListSecrets(ctx, *id)
			if err != nil {
				return fmt.Errorf("listing Subscription %q Primary and Secondary Keys (API Management Service %q / Resource Group %q): %+v", id.SubscriptionId, id.ServiceName, id.ResourceGroupName, err)
			}
			if model := keyResp.Model; model != nil {
				state.PrimaryKey = pointer.From(model.PrimaryKey)
				state.SecondaryKey = pointer.From(model.SecondaryKey)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ApiManagementWorkspaceSubscriptionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.SubscriptionClient_v2024_05_01

			id, err := subscription.ParseWorkspaceSubscriptions2ID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ApiManagementWorkspaceSubscriptionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := subscription.SubscriptionUpdateParameters{
				Properties: &subscription.SubscriptionUpdateParameterProperties{},
			}

			if metadata.ResourceData.HasChange("display_name") {
				parameters.Properties.DisplayName = pointer.To(model.DisplayName)
			}

			if metadata.ResourceData.HasChange("tracing_enabled") {
				parameters.Properties.AllowTracing = pointer.To(model.TracingEnabled)
			}

			if metadata.ResourceData.HasChange("state") {
				parameters.Properties.State = pointer.To(subscription.SubscriptionState(model.State))
			}

			if metadata.ResourceData.HasChange("primary_key") {
				parameters.Properties.PrimaryKey = pointer.To(model.PrimaryKey)
			}

			if metadata.ResourceData.HasChange("secondary_key") {
				parameters.Properties.SecondaryKey = pointer.To(model.SecondaryKey)
			}

			regularSubscriptionId := subscription.NewSubscriptions2ID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.SubscriptionName)
			if _, err := client.Update(ctx, regularSubscriptionId, parameters, subscription.DefaultUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", id.ID(), err)
			}

			return nil
		},
	}
}

func (r ApiManagementWorkspaceSubscriptionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.SubscriptionClient_v2024_05_01

			id, err := subscription.ParseWorkspaceSubscriptions2ID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.WorkspaceSubscriptionDelete(ctx, *id, subscription.WorkspaceSubscriptionDeleteOperationOptions{}); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
