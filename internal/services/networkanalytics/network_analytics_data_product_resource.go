// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package networkanalytics

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2023-05-01/cognitiveservicesaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/volumegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/networkanalytics/2023-11-15/dataproducts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/purview/2021-07-01/account"
	commonValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type NetworkAnalyticsDataProductModel struct {
	Name                              string            `tfschema:"name"`
	Location                          string            `tfschema:"location"`
	ResourceGroupName                 string            `tfschema:"resource_group_name"`
	MajorVersion                      string            `tfschema:"major_version"`
	ProductName                       string            `tfschema:"product_name"`
	Publisher                         string            `tfschema:"publisher"`
	KeyVaultKeyId                     string            `tfschema:"key_vault_key_id"`
	ManagedResourceGroupConfiguration ManagedResourceGroupConfigurationModel            `tfschema:"managed_resource_group_configuration"`
	NetworkAcls NetworkAclsModel            `tfschema:"network_acls"`
	PublicNetworkAccessEnabled        string            `tfschema:"public_network_access_enabled"`
	RedundancyEnabled                 string            `tfschema:"redundancy_enabled"`
	CurrentMinorVersion               string            `tfschema:"current_minor_version"`
	Identity                          string            `tfschema:"identity"`
	Owners                            string            `tfschema:"owners"`
	PrivateLinkEnabled                string            `tfschema:"private_link_enabled"`
	PurviewId                         string            `tfschema:"purview_id"`
	PurviewCollection                 string            `tfschema:"purview_collection"`
	Tags                              map[string]string `tfschema:"tags"`
}
type ManagedResourceGroupConfigurationModel struct {
	Name   []string `tfschema:"name"`
	Location string `tfschema:"location"`
}

type NetworkAclsModel struct {
	AllowedQueryIpRanges   []string `tfschema:"allowed_query_ip_ranges"`
	DefaultAction string `tfschema:"default_action"`
	IpRules IpRuleModel `tfschema:"ip_rules"`
	VirtualNetworkRules VirtualNetworkRuleModel `tfschema:"virtual_network_rules"`
}


type IpRuleModel struct {
	Action   []string `tfschema:"action"`
	Value string `tfschema:"value"`

}

type VirtualNetworkRuleModel struct {
	SubnetId   []string `tfschema:"subnet_id"`
	Action string `tfschema:"action"`

}

type NetworkAnalyticsDataProductResource struct{}

var _ sdk.ResourceWithUpdate = NetworkAnalyticsDataProductResource{}

func (r NetworkAnalyticsDataProductResource) ResourceType() string {
	return "azurerm_network_analytics_data_product"
}

func (r NetworkAnalyticsDataProductResource) ModelObject() interface{} {
	return &NetworkAnalyticsDataProductModel{}
}

func (r NetworkAnalyticsDataProductResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return dataproducts.ValidateDataProductID
}

func (r NetworkAnalyticsDataProductResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-z][a-z0-9]{3, 63}$`), "must be between 3 and 63 characters in length and contains only lowercase letters or numbers"),
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"major_version": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"product_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"publisher": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"key_vault_key_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: keyVaultValidate.NestedItemIdWithOptionalVersion,
		},

		"managed_resource_group_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},
					"location": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},
				},
			},
		},

		"network_acls": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			ForceNew:     true,
			MaxItems:     1,
			RequiredWith: []string{"custom_subdomain_name"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"allowed_query_ip_ranges": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.Any(
								commonValidate.IPv4Address,
							),
						},
					},
					"default_action": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(dataproducts.PossibleValuesForDefaultAction(), false),
					},
					"ip_rules": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"action": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(dataproducts.PossibleValuesForDefaultAction(), false),
								},

								"value": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ValidateFunc: validation.Any(
										commonValidate.IPv4Address,
										commonValidate.CIDR,
									),
									Set: set.HashIPv4AddressOrCIDR,
								},
							},
						},
					},

					"virtual_network_rules": {
						Type:       pluginsdk.TypeSet,
						Required:   true,
						ConfigMode: pluginsdk.SchemaConfigModeAuto,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"subnet_id": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"action": {
									Type:         pluginsdk.TypeBool,
									Optional:     true,
									ValidateFunc: validation.StringInSlice(dataproducts.PossibleValuesForDefaultAction(), false),
								},
							},
						},
					},
				},
			},
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  true,
		},

		"redundancy_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"current_minor_version": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

		"owners": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"private_link_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"purview_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: account.ValidateAccountID,
		},

		"purview_collection": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": commonschema.Tags(),
	}
}

func (r NetworkAnalyticsDataProductResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r NetworkAnalyticsDataProductResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model NetworkAnalyticsDataProductModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.NetworkAnalytics.DataProductsClient

			subscriptionId := metadata.Client.Account.SubscriptionId
			id := dataproducts.NewDataProductID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := dataproducts.DataProduct{
				Identity: nil,
				Location: "",
				Properties: &dataproducts.DataProductProperties{
					CurrentMinorVersion: nil,
					CustomerEncryptionKey: &dataproducts.EncryptionKeyDetails{
						KeyName:     "",
						KeyVaultUri: "",
						KeyVersion:  "",
					},
					CustomerManagedKeyEncryptionEnabled: nil,
					MajorVersion:                        "",
					ManagedResourceGroupConfiguration: &dataproducts.ManagedResourceGroupConfiguration{
						Location: "",
						Name:     "",
					},
					Networkacls:         expandNetworkAnalyticsDataProductNetworkAcls(),
					Owners:              nil,
					PrivateLinksEnabled: nil,
					Product:             "",
					PublicNetworkAccess: nil,
					Publisher:           "",
					PurviewAccount:      nil,
					PurviewCollection:   nil,
					Redundancy:          nil,
				},
				Tags: &model.Tags,
			}

			if err := client.CreateThenPoll(ctx, id, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r NetworkAnalyticsDataProductResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkAnalytics.DataProductsClient

			id, err := dataproducts.ParseDataProductID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model NetworkAnalyticsDataProductModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = &model.Tags
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r NetworkAnalyticsDataProductResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkAnalytics.DataProductsClient

			id, err := dataproducts.ParseDataProductID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := NetworkAnalyticsDataProductModel{
				Name:              id.AzureTrafficCollectorName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				if properties := model.Properties; properties != nil {
					state.CollectorPolicies = flattenCollectorPolicyModelArray(properties.CollectorPolicies)
					state.VirtualHub = flattenVirtualHubModel(properties.VirtualHub)
				}

				if model.Tags != nil {
					state.Tags = *model.Tags
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r NetworkAnalyticsDataProductResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkAnalytics.DataProductsClient

			id, err := dataproducts.ParseAzureTrafficCollectorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func flattenCollectorPolicyModelArray(inputList *[]dataproducts.ResourceReference) []string {
	var outputList []string
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		if input.Id != nil {
			outputList = append(outputList, *input.Id)
		}
	}

	return outputList
}

func flattenVirtualHubModel(input *dataproducts.ResourceReference) []string {
	var outputList []string
	if input != nil && input.Id != nil {
		outputList = append(outputList, *input.Id)
	}

	return outputList
}

func expandNetworkAnalyticsDataProductNetworkAcls(input []NetworkAnalyticsDataProductNetworkAclsModel) *dataproducts.DataProductNetworkAcls {
	if input == nil || len(input) == 0 {
		return &dataproducts.DataProductNetworkAcls{}
	}



	v := input.(map[string]interface{})

	defaultAction := cognitiveservicesaccounts.NetworkRuleAction(v["default_action"].(string))


	ipRulesRaw := v["ip_rules"].(*pluginsdk.Set)
	ipRules := make([]cognitiveservicesaccounts.IPRule, 0)

	for _, v := range ipRulesRaw.List() {
		rule := cognitiveservicesaccounts.IPRule{
			Value: v.(string),
		}
		ipRules = append(ipRules, rule)
	}

	var networkRule []dataproducts.VirtualNetworkRule
	for _, rule := range input {
		networkRules = append(networkRules, dataproducts.VirtualNetworkRule{
			Id:     rule.,
			Action: pointer.To(dataproducts.Action(rule.Action)),
		})
	}

	return &dataproducts.DataProductNetworkAcls{
		AllowedQueryIPRangeList: nil,
		DefaultAction:           "",
		IPRules:                 nil,
		VirtualNetworkRule:      networkRule,
	}
}
