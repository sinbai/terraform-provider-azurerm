// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ApiManagementWorkspaceGlobalSchemaModel struct {
	Name                     string `tfschema:"name"`
	ApiManagementWorkspaceId string `tfschema:"api_management_workspace_id"`
	Type                     string `tfschema:"type"`
	Value                    string `tfschema:"value"`
	Description              string `tfschema:"description"`
}

type ApiManagementWorkspaceGlobalSchemaResource struct{}

var _ sdk.ResourceWithUpdate = ApiManagementWorkspaceGlobalSchemaResource{}

func (r ApiManagementWorkspaceGlobalSchemaResource) ResourceType() string {
	return "azurerm_api_management_workspace_global_schema"
}

func (r ApiManagementWorkspaceGlobalSchemaResource) ModelObject() interface{} {
	return &ApiManagementWorkspaceGlobalSchemaModel{}
}

func (r ApiManagementWorkspaceGlobalSchemaResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return schema.ValidateWorkspaceSchemaID
}

func (r ApiManagementWorkspaceGlobalSchemaResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": schemaz.SchemaApiManagementChildName(),

		"api_management_workspace_id": commonschema.ResourceIDReferenceRequiredForceNew(&schema.WorkspaceId{}),

		"type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(schema.PossibleValuesForSchemaType(), false),
		},

		"value": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			ValidateFunc:     validation.StringIsNotEmpty,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r ApiManagementWorkspaceGlobalSchemaResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ApiManagementWorkspaceGlobalSchemaResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.GlobalSchemaClient_v2024_05_01

			var model ApiManagementWorkspaceGlobalSchemaModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			workspaceId, err := schema.ParseWorkspaceID(model.ApiManagementWorkspaceId)
			if err != nil {
				return err
			}

			id := schema.NewWorkspaceSchemaID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.ServiceName, workspaceId.WorkspaceId, model.Name)

			existing, err := client.WorkspaceGlobalSchemaGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := schema.GlobalSchemaContract{
				Properties: &schema.GlobalSchemaContractProperties{
					SchemaType: schema.SchemaType(model.Type),
				},
			}

			if model.Type == string(schema.SchemaTypeJson) {
				var jsonValue interface{}
				if err := json.Unmarshal([]byte(model.Value), &jsonValue); err != nil {
					return fmt.Errorf("parsing JSON value: %+v", err)
				}
				payload.Properties.Document = pointer.To(jsonValue)
			} else {
				payload.Properties.Value = pointer.To(interface{}(model.Value))
			}

			if model.Description != "" {
				payload.Properties.Description = pointer.To(model.Description)
			}

			if err := client.WorkspaceGlobalSchemaCreateOrUpdateThenPoll(ctx, id, payload, schema.DefaultWorkspaceGlobalSchemaCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ApiManagementWorkspaceGlobalSchemaResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.GlobalSchemaClient_v2024_05_01

			id, err := schema.ParseWorkspaceSchemaID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.WorkspaceGlobalSchemaGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := ApiManagementWorkspaceGlobalSchemaModel{
				Name:                     id.SchemaId,
				ApiManagementWorkspaceId: schema.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceId).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Description = pointer.From(props.Description)
					state.Type = string(props.SchemaType)

					if props.SchemaType == schema.SchemaTypeJson && props.Document != nil {
						var document []byte
						if document, err = json.Marshal(props.Document); err != nil {
							return fmt.Errorf("reading the schema document %s: %s", *id, err)
						}
						state.Value = string(document)
					} else {
						state.Value = pointer.From(props.Value).(string)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ApiManagementWorkspaceGlobalSchemaResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.GlobalSchemaClient_v2024_05_01

			var model ApiManagementWorkspaceGlobalSchemaModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := schema.ParseWorkspaceSchemaID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.WorkspaceGlobalSchemaGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			if resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			payload := resp.Model
			if metadata.ResourceData.HasChange("description") {
				payload.Properties.Description = pointer.To(model.Description)
			}

			if metadata.ResourceData.HasChange("value") {
				if model.Type == string(schema.SchemaTypeJson) {
					var jsonValue interface{}
					if err := json.Unmarshal([]byte(model.Value), &jsonValue); err != nil {
						return fmt.Errorf("parsing JSON value: %+v", err)
					}
					payload.Properties.Document = pointer.To(jsonValue)
				} else {
					payload.Properties.Value = pointer.To(interface{}(model.Value))
				}
			}

			if err := client.WorkspaceGlobalSchemaCreateOrUpdateThenPoll(ctx, *id, *payload, schema.DefaultWorkspaceGlobalSchemaCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ApiManagementWorkspaceGlobalSchemaResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.GlobalSchemaClient_v2024_05_01

			id, err := schema.ParseWorkspaceSchemaID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.WorkspaceGlobalSchemaDelete(ctx, *id, schema.DefaultWorkspaceGlobalSchemaDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ApiManagementWorkspaceGlobalSchemaResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// API behavior is that once configured , it cannot be removed(consistent with the portal behavior).
			// It can only be updated, or omitted at the time of creation.
			if oldVal, newVal := metadata.ResourceDiff.GetChange("description"); oldVal.(string) != "" && newVal.(string) == "" {
				if err := metadata.ResourceDiff.ForceNew("description"); err != nil {
					return err
				}
			}

			return nil
		},
	}
}
