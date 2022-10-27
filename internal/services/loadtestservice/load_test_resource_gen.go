package loadtestservice

// NOTE: this file is generated - manual changes will be overwritten.
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.
import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/utils"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2022-12-01/loadtests"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = LoadTestResource{}
var _ sdk.ResourceWithUpdate = LoadTestResource{}

type LoadTestResource struct{}

func (r LoadTestResource) ModelObject() interface{} {
	return &LoadTestResourceSchema{}
}

type LoadTestResourceSchema struct {
	DataPlaneURI       string                 `tfschema:"data_plane_uri"`
	Description        string                 `tfschema:"description"`
	CustomerManagedKey []CustomerManagedKey   `tfschema:"customer_managed_key"`
	Location           string                 `tfschema:"location"`
	Name               string                 `tfschema:"name"`
	ResourceGroupName  string                 `tfschema:"resource_group_name"`
	Tags               map[string]interface{} `tfschema:"tags"`
}

type CustomerManagedKey struct {
	UserAssignedIdentityId string `tfschema:"user_assigned_identity_id"`
	KeyVaultKeyID          string `tfschema:"key_vault_key_id"`
}

func (r LoadTestResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return loadtests.ValidateLoadTestID
}
func (r LoadTestResource) ResourceType() string {
	return "azurerm_load_test"
}
func (r LoadTestResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),
		"name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"resource_group_name": commonschema.ResourceGroupName(),
		"description": {
			ForceNew: true,
			Optional: true,
			Type:     pluginsdk.TypeString,
		},
		"customer_managed_key": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_vault_key_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: keyVaultValidate.NestedItemIdWithOptionalVersion,
					},

					"user_assigned_identity_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
					},
				},
			},
		},
		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),
		"tags":     commonschema.Tags(),
	}
}
func (r LoadTestResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"data_plane_uri": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
	}
}
func (r LoadTestResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LoadTestService.LoadTests

			var config LoadTestResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			id := loadtests.NewLoadTestID(subscriptionId, config.ResourceGroupName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			var payload loadtests.LoadTestResource
			if err := r.mapLoadTestResourceSchemaToLoadTestResource(metadata, config, &payload); err != nil {
				return fmt.Errorf("mapping schema model to sdk model: %+v", err)
			}

			if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
func (r LoadTestResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LoadTestService.LoadTests
			schema := LoadTestResourceSchema{}

			id, err := loadtests.ParseLoadTestID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				schema.Name = id.LoadTestName
				schema.ResourceGroupName = id.ResourceGroupName
				if err := r.mapLoadTestResourceToLoadTestResourceSchema(metadata, *model, &schema); err != nil {
					return fmt.Errorf("flattening model: %+v", err)
				}
			}

			return metadata.Encode(&schema)
		},
	}
}
func (r LoadTestResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LoadTestService.LoadTests

			id, err := loadtests.ParseLoadTestID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
func (r LoadTestResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LoadTestService.LoadTests

			id, err := loadtests.ParseLoadTestID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config LoadTestResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			var payload loadtests.LoadTestResourcePatchRequestBody
			if err := r.mapLoadTestResourceSchemaToLoadTestResourcePatchRequestBody(metadata, config, &payload); err != nil {
				return fmt.Errorf("mapping schema model to sdk model: %+v", err)
			}

			if _, err := client.Update(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r LoadTestResource) mapLoadTestResourceSchemaToLoadTestProperties(input LoadTestResourceSchema, output *loadtests.LoadTestProperties) error {
	output.Description = &input.Description

	if cmk := input.CustomerManagedKey; len(cmk) > 0 {
		t := loadtests.TypeUserAssigned
		output.Encryption = &loadtests.EncryptionProperties{
			Identity: &loadtests.EncryptionPropertiesIdentity{
				ResourceId: utils.String(cmk[0].UserAssignedIdentityId),
				Type:       &t,
			},
			KeyUrl: utils.String(cmk[0].KeyVaultKeyID),
		}
	}
	return nil
}

func (r LoadTestResource) mapLoadTestPropertiesToLoadTestResourceSchema(input loadtests.LoadTestProperties, output *LoadTestResourceSchema) error {
	output.DataPlaneURI = pointer.From(input.DataPlaneURI)
	output.Description = pointer.From(input.Description)

	if e := input.Encryption; e != nil {
		userAssignedIdentityId := ""
		if v := e.Identity; v != nil {
			userAssignedIdentityId = *v.ResourceId
		}
		keyVaultKeyID := ""
		if v := e.KeyUrl; v != nil {
			keyVaultKeyID = *v
		}

		result := CustomerManagedKey{
			UserAssignedIdentityId: userAssignedIdentityId,
			KeyVaultKeyID:          keyVaultKeyID,
		}
		output.CustomerManagedKey = []CustomerManagedKey{result}
	}
	return nil
}

func (r LoadTestResource) mapLoadTestResourceSchemaToLoadTestResourcePatchRequestBodyProperties(input LoadTestResourceSchema, output *loadtests.LoadTestResourcePatchRequestBodyProperties) error {
	output.Description = &input.Description
	if cmk := input.CustomerManagedKey; len(cmk) > 0 {
		t := loadtests.TypeUserAssigned
		output.Encryption = &loadtests.EncryptionProperties{
			Identity: &loadtests.EncryptionPropertiesIdentity{
				ResourceId: utils.String(cmk[0].UserAssignedIdentityId),
				Type:       &t,
			},
			KeyUrl: utils.String(cmk[0].KeyVaultKeyID),
		}
	}
	return nil
}

func (r LoadTestResource) mapLoadTestResourcePatchRequestBodyPropertiesToLoadTestResourceSchema(input loadtests.LoadTestResourcePatchRequestBodyProperties, output *LoadTestResourceSchema) error {
	output.Description = pointer.From(input.Description)
	if e := input.Encryption; e != nil {
		userAssignedIdentityId := ""
		if v := e.Identity; v != nil {
			userAssignedIdentityId = *v.ResourceId
		}
		keyVaultKeyID := ""
		if v := e.KeyUrl; v != nil {
			keyVaultKeyID = *v
		}

		result := CustomerManagedKey{
			UserAssignedIdentityId: userAssignedIdentityId,
			KeyVaultKeyID:          keyVaultKeyID,
		}
		output.CustomerManagedKey = []CustomerManagedKey{result}
	}
	return nil
}

func (r LoadTestResource) mapLoadTestResourceSchemaToLoadTestResource(metadata sdk.ResourceMetaData, input LoadTestResourceSchema, output *loadtests.LoadTestResource) error {
	identity, err := identity.ExpandLegacySystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}
	output.Identity = identity
	output.Location = location.Normalize(input.Location)
	output.Tags = tags.Expand(input.Tags)

	if output.Properties == nil {
		output.Properties = &loadtests.LoadTestProperties{}
	}
	if err := r.mapLoadTestResourceSchemaToLoadTestProperties(input, output.Properties); err != nil {
		return fmt.Errorf("mapping Schema to SDK Field %q / Model %q: %+v", "LoadTestProperties", "Properties", err)
	}

	return nil
}

func (r LoadTestResource) mapLoadTestResourceToLoadTestResourceSchema(metadata sdk.ResourceMetaData, input loadtests.LoadTestResource, output *LoadTestResourceSchema) error {

	identityValue, err := identity.FlattenLegacySystemAndUserAssignedMap(input.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}

	if err := metadata.ResourceData.Set("identity", identityValue); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	output.Location = location.Normalize(input.Location)
	output.Tags = tags.Flatten(input.Tags)

	if input.Properties == nil {
		input.Properties = &loadtests.LoadTestProperties{}
	}
	if err := r.mapLoadTestPropertiesToLoadTestResourceSchema(*input.Properties, output); err != nil {
		return fmt.Errorf("mapping SDK Field %q / Model %q to Schema: %+v", "LoadTestProperties", "Properties", err)
	}

	return nil
}

func (r LoadTestResource) mapLoadTestResourceSchemaToLoadTestResourcePatchRequestBody(metadata sdk.ResourceMetaData, input LoadTestResourceSchema, output *loadtests.LoadTestResourcePatchRequestBody) error {

	identity, err := identity.ExpandLegacySystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding SystemAssigned Identity: %+v", err)
	}
	output.Identity = identity

	output.Tags = tags.Expand(input.Tags)

	if output.Properties == nil {
		output.Properties = &loadtests.LoadTestResourcePatchRequestBodyProperties{}
	}
	if err := r.mapLoadTestResourceSchemaToLoadTestResourcePatchRequestBodyProperties(input, output.Properties); err != nil {
		return fmt.Errorf("mapping Schema to SDK Field %q / Model %q: %+v", "LoadTestResourcePatchRequestBodyProperties", "Properties", err)
	}

	return nil
}

func (r LoadTestResource) mapLoadTestResourcePatchRequestBodyToLoadTestResourceSchema(metadata sdk.ResourceMetaData, input loadtests.LoadTestResourcePatchRequestBody, output *LoadTestResourceSchema) error {

	identityValue, err := identity.FlattenLegacySystemAndUserAssignedMap(input.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}

	if err := metadata.ResourceData.Set("identity", identityValue); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	output.Tags = tags.Flatten(input.Tags)

	if input.Properties == nil {
		input.Properties = &loadtests.LoadTestResourcePatchRequestBodyProperties{}
	}
	if err := r.mapLoadTestResourcePatchRequestBodyPropertiesToLoadTestResourceSchema(*input.Properties, output); err != nil {
		return fmt.Errorf("mapping SDK Field %q / Model %q to Schema: %+v", "LoadTestResourcePatchRequestBodyProperties", "Properties", err)
	}

	return nil
}
