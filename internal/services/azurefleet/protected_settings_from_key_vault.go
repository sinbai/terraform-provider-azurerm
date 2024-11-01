// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azurefleet

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	`github.com/hashicorp/go-azure-sdk/resource-manager/azurefleet/2024-11-01/fleets`
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func protectedSettingsFromKeyVaultSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:          pluginsdk.TypeList,
		Optional:      true,
		MaxItems:      1,
		ConflictsWith: []string{"protected_settings"},
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"secret_url": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: keyVaultValidate.NestedItemId,
				},

				"source_vault_id": commonschema.ResourceIDReferenceRequired(&commonids.KeyVaultId{}),
			},
		},
	}
}

func expandProtectedSettingsFromKeyVault(input []interface{}) *fleets.KeyVaultSecretReference {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &fleets.KeyVaultSecretReference{
		SecretURL: v["secret_url"].(string),
		SourceVault: fleets.SubResource{
			Id: pointer.To(v["source_vault_id"].(string)),
		},
	}
}

func flattenProtectedSettingsFromKeyVault(input *fleets.KeyVaultSecretReference) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	sourceVaultId := ""
	if input.SourceVault.Id != nil {
		sourceVaultId = *input.SourceVault.Id
	}

	return []interface{}{
		map[string]interface{}{
			"secret_url":      input.SecretURL,
			"source_vault_id": sourceVaultId,
		},
	}
}
