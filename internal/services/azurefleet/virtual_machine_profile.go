// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azurefleet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurefleet/2024-11-01/fleets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservationgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/images"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplicationversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/applicationsecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networksecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/publicipprefixes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func virtualMachineProfileSchema(required bool) *pluginsdk.Schema {
	vmProfile := &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: required,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"network_api_version": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForNetworkApiVersion(), false),
				},

				"network_interface": networkInterfaceSchema(),

				"os_profile": osProfileSchema(),

				"storage_profile": storageProfileSchema(),

				"boot_diagnostic_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},

				"boot_diagnostic_storage_account_endpoint": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"capacity_reservation_group_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: capacityreservationgroups.ValidateCapacityReservationGroupID,
				},

				"extensions": extensionSchema(),

				"extensions_time_budget": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      "PT1H30M",
					ValidateFunc: azValidate.ISO8601DurationBetween("PT15M", "PT2H"),
				},

				"gallery_application": galleryApplicationSchema(),

				"license_type": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"network_health_probe_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: azure.ValidateResourceID,
				},

				// if it is specified os_image_notification_profile enable is set to true.
				"scheduled_event_os_image_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},

				"scheduled_event_os_image_timeout": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  "PT5M",
					ValidateFunc: validation.StringInSlice([]string{
						"PT5M",
					}, false),
				},

				"scheduled_event_termination_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},

				// if it is specified terminate_notification_profile enable is set to true.
				"scheduled_event_termination_timeout": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  "PT15M",
					ValidateFunc: validation.StringInSlice([]string{
						"PT15M",
					}, false),
				},

				"security_posture_reference": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"exclude_extensions": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								Elem: &pluginsdk.Schema{
									Type: pluginsdk.TypeString,
								},
							},

							"id": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},

							"override_enabled": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
							},
						},
					},
				},

				"security_profile": securityProfileSchema(),

				"service_artifact_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"user_data_base64": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsBase64,
				},

				"vm_size": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"vcpu_available_count": {
								Type:     pluginsdk.TypeInt,
								Optional: true,
							},

							"vcpu_per_core_count": {
								Type:     pluginsdk.TypeInt,
								Optional: true,
							},
						},
					},
				},
			},
		},
	}
	if !required {
		vmProfile.Optional = true
	}
	return vmProfile
}

func galleryApplicationSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 100,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"version_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: galleryapplicationversions.ValidateApplicationVersionID,
				},

				// Example: https://mystorageaccount.blob.core.windows.net/configurations/settings_json.config
				"configuration_blob_uri": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IsURLWithHTTPorHTTPS,
				},

				"order": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      0,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(0, 2147483647),
				},

				// NOTE: Per the service team, "this is a pass through value that we just add to the model but don't depend on. It can be any string."
				"tag": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"automatic_upgrade_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},

				"treat_failure_as_deployment_failure_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},
			},
		},
	}
}

func extensionSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"publisher": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"type": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"type_handler_version": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"auto_upgrade_minor_version_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  true,
				},

				"automatic_upgrade_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"force_update_tag": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"protected_settings_json": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsJSON,
				},

				"protected_settings_from_key_vault": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
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
				},

				"provision_after_extensions": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},

				"settings_json": {
					Type:             pluginsdk.TypeString,
					Optional:         true,
					ValidateFunc:     validation.StringIsJSON,
					DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
				},

				"suppress_failures_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	}
}

func networkInterfaceSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"accelerated_networking_enabled": {
					Type:     pluginsdk.TypeBool,
					Required: true,
				},

				"ip_configuration": ipConfigurationSchema(),

				"ip_forwarding_enabled": {
					Type:     pluginsdk.TypeBool,
					Required: true,
				},

				"primary": {
					Type:     pluginsdk.TypeBool,
					Required: true,
				},

				"dns_servers": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},

				"network_security_group_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: networksecuritygroups.ValidateNetworkSecurityGroupID,
				},

				"auxiliary_mode": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						// NOTE: because there is a `None` value in the possible values, it's handled in the Create/Update and Read functions.
						string(fleets.NetworkInterfaceAuxiliaryModeAcceleratedConnections),
						string(fleets.NetworkInterfaceAuxiliaryModeFloating),
					}, false),
				},

				"auxiliary_sku": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						// NOTE: because there is a `None` value in the possible values, it's handled in the Create/Update and Read functions.
						string(fleets.NetworkInterfaceAuxiliarySkuATwo),
						string(fleets.NetworkInterfaceAuxiliarySkuAFour),
						string(fleets.NetworkInterfaceAuxiliarySkuAEight),
						string(fleets.NetworkInterfaceAuxiliarySkuAOne),
					}, false),
				},

				"delete_option": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForDiskDeleteOptionTypes(), false),
				},

				"tcp_state_tracking_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					// need to confirm if DisableTcpStateTracking default value is false?
					Default: true,
				},

				"fpga_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},
			},
		},
	}
}

func ipConfigurationSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"load_balancer_backend_address_pool_ids": {
					Type:     pluginsdk.TypeSet,
					Required: true,
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					Set:      pluginsdk.HashString,
				},

				"primary": {
					Type:     pluginsdk.TypeBool,
					Required: true,
				},

				"subnet_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: commonids.ValidateSubnetID,
				},

				"application_gateway_backend_address_pool_ids": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					Set:      pluginsdk.HashString,
				},

				"application_security_group_ids": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: applicationsecuritygroups.ValidateApplicationSecurityGroupID,
					},
					Set:      pluginsdk.HashString,
					MaxItems: 20,
				},

				"load_balancer_inbound_nat_rules_ids": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					Set:      pluginsdk.HashString,
				},

				"public_ip_address": publicIPAddressSchema(),

				"version": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					Default:      string(fleets.IPVersionIPvFour),
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForIPVersion(), false),
				},
			},
		},
	}
}

func publicIPAddressSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"delete_option": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForDeleteOptions(), false),
				},

				"domain_name_label": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"domain_name_label_scope": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForDomainNameLabelScopeTypes(), false),
				},

				"idle_timeout_in_minutes": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(4, 32),
				},
				"ip_tag": {
					// TODO: does this want to be a Set?
					Type:     pluginsdk.TypeList,
					Optional: true,
					ForceNew: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"tag": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ForceNew:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"type": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ForceNew:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
				"version": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					Default:      string(fleets.IPVersionIPvFour),
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForIPVersion(), false),
				},
				// TODO: preview feature
				// $ az feature register --namespace Microsoft.Network --name AllowBringYourOwnPublicIpAddress
				// $ az provider register -n Microsoft.Network
				"public_ip_prefix_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: publicipprefixes.ValidatePublicIPPrefixID,
				},

				"sku": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"name": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForPublicIPAddressSkuName(), false),
							},

							"tier": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForPublicIPAddressSkuTier(), false),
							},
						},
					},
				},
			},
		},
	}
}

func osProfileSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"computer_name_prefix": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"admin_username": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"admin_password": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"extension_operations_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ForceNew: true,
				},

				"custom_data_base64": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsBase64,
				},

				"linux_configuration": {
					Type:          pluginsdk.TypeList,
					Optional:      true,
					MaxItems:      1,
					ConflictsWith: []string{"compute_profile.0.virtual_machine_profile.0.os_profile.0.windows_configuration"},
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"password_authentication_enabled": {
								Type:     pluginsdk.TypeBool,
								Required: true,
								ForceNew: true,
							},

							"vm_agent_platform_updates_enabled": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
							},

							"patch_setting": linuxPatchSettingSchema(),

							"provision_vm_agent_enabled": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
								ForceNew: true,
							},

							"ssh_keys": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"path": {
											Type:         pluginsdk.TypeString,
											Required:     true,
											ValidateFunc: validation.StringIsNotEmpty,
										},
										"key_data": {
											Type:         pluginsdk.TypeString,
											Optional:     true,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},
								},
							},
						},
					},
				},

				"require_guest_provision_signal_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},

				"secret": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"key_vault_id": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: azure.ValidateResourceID,
							},

							"certificate": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"url": {
											Type:         pluginsdk.TypeString,
											Required:     true,
											ValidateFunc: validation.StringIsNotEmpty,
										},
										// linux does not contain this property?
										"store": {
											Type:         pluginsdk.TypeString,
											Optional:     true,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},
								},
							},
						},
					},
				},

				"windows_configuration": {
					Type:          pluginsdk.TypeList,
					Optional:      true,
					MaxItems:      1,
					ConflictsWith: []string{"compute_profile.0.virtual_machine_profile.0.os_profile.0.linux_configuration"},
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"additional_unattend_content": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"pass_name": {
											Type:         pluginsdk.TypeString,
											Required:     true,
											ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForPassName(), false),
										},
										"component_name": {
											Type:         pluginsdk.TypeString,
											Required:     true,
											ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForComponentName(), false),
										},
										"setting_name": {
											Type:         pluginsdk.TypeString,
											Required:     true,
											ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForSettingNames(), false),
										},
										"content": {
											Type:         pluginsdk.TypeString,
											Required:     true,
											Sensitive:    true,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},
								},
							},

							"automatic_updates_enabled": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
							},

							"vm_agent_platform_updates_enabled": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
							},

							"patch_setting": windowsPatchSettingSchema(),

							"provision_vm_agent_enabled": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
								Default:  true,
								ForceNew: true,
							},

							"time_zone": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},

							"winrm": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"protocol": {
											Type:         pluginsdk.TypeString,
											Required:     true,
											ForceNew:     true,
											ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForProtocolTypes(), false),
										},
										"certificate_url": {
											Type:         pluginsdk.TypeString,
											Optional:     true,
											ForceNew:     true,
											ValidateFunc: keyVaultValidate.NestedItemId,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func storageProfileSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"image_reference": storageProfileImageReferenceSchema(),

				"os_disk": storageProfileOsDiskSchema(),

				"disk_controller_type": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForDiskControllerTypes(), false),
				},

				"data_disk": storageProfileDataDiskSchema(),
			},
		},
	}
}

func linuxPatchSettingSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"assessment_mode": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForLinuxPatchAssessmentMode(), false),
				},

				"automatic_by_platform_setting": {
					Type:     pluginsdk.TypeList,
					Required: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"bypass_platform_safety_checks_enabled": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
							},

							"reboot_setting": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForLinuxVMGuestPatchAutomaticByPlatformRebootSetting(), false),
							},
						},
					},
				},

				"patch_mode": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForLinuxVMGuestPatchMode(), false),
				},
			},
		},
	}
}

func windowsPatchSettingSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"assessment_mode": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForLinuxPatchAssessmentMode(), false),
				},

				"automatic_by_platform_setting": {
					Type:     pluginsdk.TypeList,
					Required: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"bypass_platform_safety_checks_enabled": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
							},

							"reboot_setting": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForLinuxVMGuestPatchAutomaticByPlatformRebootSetting(), false),
							},
						},
					},
				},

				"patch_mode": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForLinuxVMGuestPatchMode(), false),
				},

				"hot_patching_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},
			},
		},
	}
}

func securityProfileSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"encryption_at_host_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},

				"user_assigned_identity_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"proxy_agent": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							// there is another property enable exists
							// need to confirm whether the following properties should be set when this feature is enable?
							// key_incarnation_value is required?
							"key_incarnation_value": {
								Type:     pluginsdk.TypeInt,
								Required: true,
							},
							"mode": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								Default:      string(fleets.ModeEnforce),
								ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForMode(), false),
							},
						},
					},
				},

				"security_type": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForSecurityTypes(), false),
				},

				"uefi_secure_boot_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ForceNew: true,
				},

				"uefi_vtpm_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ForceNew: true,
				},
			},
		},
	}
}

func osDiskManagedDiskSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"storage_account_type": {
					Type:     pluginsdk.TypeString,
					Required: true,
					// `PremiumV2_LRS` and `UltraSSD_LRS` is not supported OS Disk
					ValidateFunc: validation.StringInSlice([]string{
						string(fleets.StorageAccountTypesPremiumLRS),
						string(fleets.StorageAccountTypesPremiumZRS),
						string(fleets.StorageAccountTypesStandardLRS),
						string(fleets.StorageAccountTypesStandardSSDLRS),
						string(fleets.StorageAccountTypesStandardSSDZRS),
					}, false),
				},

				"disk_encryption_set_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validate.DiskEncryptionSetID,
				},

				"security_disk_encryption_set_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validate.DiskEncryptionSetID,
				},

				"security_encryption_type": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForSecurityTypes(), false),
				},
			},
		},
	}
}

func dataDiskManagedDiskSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"storage_account_type": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForStorageAccountTypes(), false),
				},

				"disk_encryption_set_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validate.DiskEncryptionSetID,
				},

				"security_disk_encryption_set_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validate.DiskEncryptionSetID,
				},

				"security_encryption_type": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForSecurityTypes(), false),
				},
			},
		},
	}
}

func storageProfileDataDiskSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"lun": {
					Type:     pluginsdk.TypeInt,
					Required: true,
				},

				"create_option": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForDiskCreateOptionTypes(), false),
				},

				"caching": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						// NOTE: because there is a `None` value in the possible values, it's handled in the Create/Update and Read functions.
						string(fleets.CachingTypesReadOnly),
						string(fleets.CachingTypesReadWrite),
					}, false),
				},

				"delete_option": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForDiskDeleteOptionTypes(), false),
				},

				"disk_iops_read_write": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
				},

				"disk_mbps_read_write": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
				},

				"disk_size_in_gb": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(1, 32767),
				},

				"managed_disk": dataDiskManagedDiskSchema(),

				"name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"write_accelerator_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},
			},
		},
	}
}

func storageProfileOsDiskSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"create_option": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForDiskCreateOptionTypes(), false),
				},

				"caching": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						// NOTE: because there is a `None` value in the possible values, it's handled in the Create/Update and Read functions.
						//string(fleets.CachingTypesNone),
						string(fleets.CachingTypesReadOnly),
						string(fleets.CachingTypesReadWrite),
					}, false),
				},

				"os_type": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"managed_disk": osDiskManagedDiskSchema(),

				"delete_option": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForDiskDeleteOptionTypes(), false),
				},

				"diff_disk_option": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForDiffDiskOptions(), false),
				},
				"diff_disk_placement": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForDiffDiskPlacement(), false),
				},

				"disk_size_in_gb": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(1, 32767),
				},

				"image_uri": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"vhd_containers": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					Set:      pluginsdk.HashString,
				},

				"write_accelerator_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},
			},
		},
	}
}

func storageProfileImageReferenceSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"publisher": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"offer": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"sku": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"version": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"id": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.Any(
						images.ValidateImageID,
						validate.SharedImageID,
						validate.SharedImageVersionID,
					),
				},

				"community_gallery_image_id": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.Any(
						images.ValidateImageID,
						validate.CommunityGalleryImageID,
						validate.CommunityGalleryImageVersionID,
					),
				},

				"shared_gallery_image_id": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.Any(
						images.ValidateImageID,
						validate.SharedGalleryImageID,
						validate.SharedGalleryImageVersionID,
					),
				},
			},
		},
	}
}

func resourceStorageProfileImageReferenceHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		if v, ok := m["publisher"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
		if v, ok := m["offer"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
		if v, ok := m["sku"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
		if v, ok := m["version"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
		if v, ok := m["id"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
		if v, ok := m["community_gallery_image_id"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
		if v, ok := m["shared_gallery_image_id"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
	}

	return pluginsdk.HashString(buf.String())
}

func expandBaseVirtualMachineProfileModel(inputList []VirtualMachineProfileModel, d *schema.ResourceData) (*fleets.BaseVirtualMachineProfile, error) {
	if len(inputList) == 0 {
		return nil, nil
	}
	input := &inputList[0]
	output := fleets.BaseVirtualMachineProfile{
		ApplicationProfile:       expandApplicationProfile(input.GalleryApplicationProfile),
		CapacityReservation:      expandCapacityReservation(input.CapacityReservationGroupId),
		DiagnosticsProfile:       expandDiagnosticsProfile(input.BootDiagnosticEnabled, input.BootDiagnosticStorageAccountEndpoint),
		HardwareProfile:          expandHardwareProfile(input.VMSize),
		NetworkProfile:           expandNetworkProfile(input.NetworkInterface, input.NetworkApiVersion, input.NetworkHealthProbeId),
		OsProfile:                expandOSProfileModel(input.OsProfile, d),
		ScheduledEventsProfile:   expandScheduledEventsProfile(input),
		SecurityPostureReference: expandSecurityPostureReferenceModel(input.SecurityPostureReference),
		SecurityProfile:          expandSecurityProfileModel(input.SecurityProfile),
		ServiceArtifactReference: expandServiceArtifactReference(input.ServiceArtifactId),
		StorageProfile:           expandStorageProfileModel(input.StorageProfile),
	}

	extensionProfileValue, err := expandExtensionProfileModel(input.Extensions, input.ExtensionsTimeBudget)
	if err != nil {
		return nil, err
	}
	output.ExtensionProfile = extensionProfileValue

	if input.LicenseType != "" {
		output.LicenseType = pointer.To(input.LicenseType)
	}

	if input.UserDataBase64 != "" {
		output.UserData = pointer.To(input.UserDataBase64)
	}

	return &output, nil
}

func expandServiceArtifactReference(id string) *fleets.ServiceArtifactReference {
	if id == "" {
		return nil
	}
	return &fleets.ServiceArtifactReference{Id: pointer.To(id)}
}

func expandApplicationProfile(inputList []GalleryApplicationModel) *fleets.ApplicationProfile {
	if len(inputList) == 0 {
		return nil
	}

	var outputList []fleets.VMGalleryApplication
	for _, v := range inputList {
		input := v
		output := fleets.VMGalleryApplication{
			EnableAutomaticUpgrade:          pointer.To(input.AutomaticUpgradeEnabled),
			Order:                           pointer.To(input.Order),
			PackageReferenceId:              input.VersionId,
			TreatFailureAsDeploymentFailure: pointer.To(input.TreatFailureAsDeploymentFailureEnabled),
		}

		if input.ConfigurationBlobUri != "" {
			output.ConfigurationReference = pointer.To(input.ConfigurationBlobUri)
		}

		if input.Tags != "" {
			output.Tags = pointer.To(input.Tags)
		}
		outputList = append(outputList, output)
	}

	output := fleets.ApplicationProfile{
		GalleryApplications: &outputList,
	}

	return &output
}

func expandCapacityReservation(input string) *fleets.CapacityReservationProfile {
	if input == "" {
		return nil
	}

	return &fleets.CapacityReservationProfile{
		CapacityReservationGroup: expandSubResource(input),
	}
}

func expandSubResource(input string) *fleets.SubResource {
	if input == "" {
		return nil
	}

	return &fleets.SubResource{
		Id: pointer.To(input),
	}
}

func expandSubResources(inputList []string) *[]fleets.SubResource {
	if len(inputList) == 0 {
		return nil
	}
	var outputList []fleets.SubResource

	for _, v := range inputList {
		input := v

		output := expandSubResource(input)
		if output != nil {
			outputList = append(outputList, pointer.From(output))
		}
	}

	return &outputList
}

func expandDiagnosticsProfile(enabled bool, endpoint string) *fleets.DiagnosticsProfile {
	bootDiagnostics := fleets.BootDiagnostics{
		Enabled:    pointer.To(enabled),
		StorageUri: pointer.To(endpoint),
	}

	output := fleets.DiagnosticsProfile{
		BootDiagnostics: pointer.To(bootDiagnostics),
	}

	return &output
}

func expandExtensionProfileModel(inputList []ExtensionsModel, timeBudget string) (*fleets.VirtualMachineScaleSetExtensionProfile, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	output := fleets.VirtualMachineScaleSetExtensionProfile{}
	extensionsValue, err := expandExtensionsModel(inputList)
	if err != nil {
		return nil, err
	}
	output.Extensions = extensionsValue

	if timeBudget != "" {
		output.ExtensionsTimeBudget = pointer.To(timeBudget)
	}
	return &output, nil
}

func expandExtensionsModel(inputList []ExtensionsModel) (*[]fleets.VirtualMachineScaleSetExtension, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	var outputList []fleets.VirtualMachineScaleSetExtension
	for _, v := range inputList {
		input := v
		output := fleets.VirtualMachineScaleSetExtension{}

		if input.Name != "" {
			output.Name = pointer.To(input.Name)
		}

		propertiesValue, err := expandExtensionModel(input)
		if err != nil {
			return nil, err
		}
		output.Properties = propertiesValue

		outputList = append(outputList, output)
	}
	return &outputList, nil
}

func expandExtensionModel(input ExtensionsModel) (*fleets.VirtualMachineScaleSetExtensionProperties, error) {
	output := fleets.VirtualMachineScaleSetExtensionProperties{
		AutoUpgradeMinorVersion:       pointer.To(input.AutoUpgradeMinorVersionEnabled),
		EnableAutomaticUpgrade:        pointer.To(input.AutomaticUpgradeEnabled),
		ProtectedSettingsFromKeyVault: expandKeyVaultSecretReferenceModel(input.ProtectedSettingsFromKeyVault),
		SuppressFailures:              pointer.To(input.SuppressFailuresEnabled),
	}
	if input.ForceUpdateTag != "" {
		output.ForceUpdateTag = pointer.To(input.ForceUpdateTag)
	}

	if input.ProtectedSettingsJson != "" {
		protectedSettingsValue := make(map[string]interface{})
		err := json.Unmarshal([]byte(input.ProtectedSettingsJson), &protectedSettingsValue)
		if err != nil {
			return nil, err
		}
		output.ProtectedSettings = pointer.To(protectedSettingsValue)
	}

	if len(input.ProvisionAfterExtensions) > 0 {
		output.ProvisionAfterExtensions = pointer.To(input.ProvisionAfterExtensions)
	}

	if input.Publisher != "" {
		output.Publisher = pointer.To(input.Publisher)
	}

	if input.SettingsJson != "" {
		settingsValue := make(map[string]interface{})
		err := json.Unmarshal([]byte(input.SettingsJson), &settingsValue)
		if err != nil {
			return nil, err
		}
		output.Settings = pointer.To(settingsValue)
	}

	if input.Type != "" {
		output.Type = pointer.To(input.Type)
	}

	if input.TypeHandlerVersion != "" {
		output.TypeHandlerVersion = pointer.To(input.TypeHandlerVersion)
	}

	return &output, nil
}

func expandKeyVaultSecretReferenceModel(inputList []ProtectedSettingsFromKeyVaultModel) *fleets.KeyVaultSecretReference {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.KeyVaultSecretReference{
		SecretURL: input.SecretUrl,
	}

	output.SourceVault = pointer.From(expandSubResource(input.SourceVaultId))

	return &output
}

func expandHardwareProfile(inputList []VMSizeModel) *fleets.VirtualMachineScaleSetHardwareProfile {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	vmSize := fleets.VMSizeProperties{
		VCPUsAvailable: pointer.To(input.VCPUAvailableCount),
		VCPUsPerCore:   pointer.To(input.VCPUPerCoreCount),
	}

	return &fleets.VirtualMachineScaleSetHardwareProfile{
		VMSizeProperties: pointer.To(vmSize),
	}
}

func expandNetworkProfile(inputList []NetworkInterfaceModel, version string, healthProbe string) *fleets.VirtualMachineScaleSetNetworkProfile {
	if len(inputList) == 0 {
		return nil
	}

	output := fleets.VirtualMachineScaleSetNetworkProfile{
		HealthProbe:                    expandApiEntityReferenceModel(healthProbe),
		NetworkApiVersion:              pointer.To(fleets.NetworkApiVersion(version)),
		NetworkInterfaceConfigurations: expandNetworkInterfaceModel(inputList),
	}

	return &output
}

func expandApiEntityReferenceModel(input string) *fleets.ApiEntityReference {
	if input == "" {
		return nil
	}

	return &fleets.ApiEntityReference{
		Id: pointer.To(input),
	}
}

func expandNetworkInterfaceModel(inputList []NetworkInterfaceModel) *[]fleets.VirtualMachineScaleSetNetworkConfiguration {
	var outputList []fleets.VirtualMachineScaleSetNetworkConfiguration
	for _, v := range inputList {
		input := v
		output := fleets.VirtualMachineScaleSetNetworkConfiguration{
			Name:       input.Name,
			Properties: expandNetworkConfigurationPropertiesModel(input),
		}

		outputList = append(outputList, output)
	}
	return &outputList
}

func expandNetworkConfigurationPropertiesModel(input NetworkInterfaceModel) *fleets.VirtualMachineScaleSetNetworkConfigurationProperties {
	output := fleets.VirtualMachineScaleSetNetworkConfigurationProperties{
		DisableTcpStateTracking:     pointer.To(!input.TcpStateTrackingEnabled),
		EnableAcceleratedNetworking: pointer.To(input.AcceleratedNetworkingEnabled),
		EnableFpga:                  pointer.To(input.FpgaEnabled),
		EnableIPForwarding:          pointer.To(input.IPForwardingEnabled),
		NetworkSecurityGroup:        expandSubResource(input.NetworkSecurityGroupId),
		Primary:                     pointer.To(input.Primary),
	}

	if len(input.DnsServers) > 0 {
		output.DnsSettings = &fleets.VirtualMachineScaleSetNetworkConfigurationDnsSettings{
			DnsServers: pointer.To(input.DnsServers),
		}
	}

	auxiliaryMode := fleets.NetworkInterfaceAuxiliaryModeNone
	if input.AuxiliaryMode != "" {
		auxiliaryMode = fleets.NetworkInterfaceAuxiliaryMode(input.AuxiliaryMode)

	}
	output.AuxiliaryMode = pointer.To(auxiliaryMode)

	if input.DeleteOption != "" {
		output.DeleteOption = pointer.To(fleets.DeleteOptions(input.DeleteOption))
	}

	auxiliarySku := fleets.NetworkInterfaceAuxiliarySkuNone
	if input.AuxiliaryMode != "" {
		auxiliarySku = fleets.NetworkInterfaceAuxiliarySku(input.AuxiliarySku)

	}
	output.AuxiliarySku = pointer.To(auxiliarySku)

	output.IPConfigurations = pointer.From(expandIPConfigurationModel(input.IPConfiguration))

	return &output
}

func expandIPConfigurationModel(inputList []IPConfigurationModel) *[]fleets.VirtualMachineScaleSetIPConfiguration {
	var outputList []fleets.VirtualMachineScaleSetIPConfiguration
	for _, v := range inputList {
		input := v
		output := fleets.VirtualMachineScaleSetIPConfiguration{
			Name: input.Name,
			Properties: &fleets.VirtualMachineScaleSetIPConfigurationProperties{
				ApplicationGatewayBackendAddressPools: expandSubResources(input.ApplicationGatewayBackendAddressPoolIds),
				ApplicationSecurityGroups:             expandSubResources(input.ApplicationSecurityGroupIds),
				LoadBalancerBackendAddressPools:       expandSubResources(input.LoadBalancerBackendAddressPoolIds),
				LoadBalancerInboundNatPools:           expandSubResources(input.LoadBalancerInboundNatPoolIds),
				Primary:                               pointer.To(input.Primary),
				PublicIPAddressConfiguration:          expandPublicIPAddressModel(input.PublicIPAddress),
				Subnet:                                expandApiEntityReferenceModel(input.SubnetId),
			},
		}

		if input.Version != "" {
			output.Properties.PrivateIPAddressVersion = pointer.To(fleets.IPVersion(input.Version))
		}
		outputList = append(outputList, output)
	}
	return &outputList
}

func expandPublicIPAddressModel(inputList []PublicIPAddressModel) *fleets.VirtualMachineScaleSetPublicIPAddressConfiguration {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := fleets.VirtualMachineScaleSetPublicIPAddressConfiguration{
		Name: input.Name,
		Properties: &fleets.VirtualMachineScaleSetPublicIPAddressConfigurationProperties{
			DnsSettings:          expandPublicIPAddressDnsSettings(input.DomainNameLabel, input.DomainNameLabelScope),
			IPTags:               expandIPTagModel(input.IPTags),
			IdleTimeoutInMinutes: pointer.To(input.IdleTimeoutInMinutes),
			PublicIPPrefix:       expandSubResource(input.PublicIPPrefix),
		},
		Sku: expandPublicIPAddressSkuModel(input.Sku),
	}

	if input.DeleteOption != "" {
		output.Properties.DeleteOption = pointer.To(fleets.DeleteOptions(input.DeleteOption))
	}
	if input.Version != "" {
		output.Properties.PublicIPAddressVersion = pointer.To(fleets.IPVersion(input.Version))
	}

	return &output
}

func expandPublicIPAddressDnsSettings(domainNameLabel string, domainNameLabelScope string) *fleets.VirtualMachineScaleSetPublicIPAddressConfigurationDnsSettings {
	if domainNameLabel == "" && domainNameLabelScope == "" {
		return nil
	}
	output := fleets.VirtualMachineScaleSetPublicIPAddressConfigurationDnsSettings{}
	if domainNameLabel != "" {
		output.DomainNameLabel = domainNameLabel
	}

	if domainNameLabelScope != "" {
		output.DomainNameLabelScope = pointer.To(fleets.DomainNameLabelScopeTypes(domainNameLabelScope))
	}

	return &output
}

func expandIPTagModel(inputList []IPTagModel) *[]fleets.VirtualMachineScaleSetIPTag {
	var outputList []fleets.VirtualMachineScaleSetIPTag
	for _, v := range inputList {
		input := v
		output := fleets.VirtualMachineScaleSetIPTag{}

		if input.IPTagType != "" {
			output.IPTagType = pointer.To(input.IPTagType)
		}

		if input.Tag != "" {
			output.Tag = pointer.To(input.Tag)
		}
		outputList = append(outputList, output)
	}
	return &outputList
}

func expandPublicIPAddressSkuModel(inputList []SkuModel) *fleets.PublicIPAddressSku {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.PublicIPAddressSku{
		Name: pointer.To(fleets.PublicIPAddressSkuName(input.Name)),
		Tier: pointer.To(fleets.PublicIPAddressSkuTier(input.Tier)),
	}

	return &output
}

func expandOSProfileModel(inputList []OSProfileModel, d *schema.ResourceData) *fleets.VirtualMachineScaleSetOSProfile {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.VirtualMachineScaleSetOSProfile{
		AdminUsername:            pointer.To(input.AdminUsername),
		AdminPassword:            pointer.To(input.AdminPassword),
		AllowExtensionOperations: pointer.To(input.ExtensionOperationsEnabled),
		LinuxConfiguration:       expandLinuxConfigurationModel(input.LinuxConfiguration),
		Secrets:                  expandOsProfileSecretsModel(input.Secret),
		WindowsConfiguration:     expandWindowsConfigurationModel(input.WindowsConfiguration),
	}

	if input.CustomDataBase64 != "" {
		output.CustomData = pointer.To(input.CustomDataBase64)
	}
	if input.ComputerNamePrefix != "" {
		output.ComputerNamePrefix = pointer.To(input.ComputerNamePrefix)
	}

	// The property 'osProfile.RequireGuestProvisionSignalEnabled' is not valid because the 'Microsoft.Compute/Agentless' feature is not enabled for this subscription
	// it must either be set to True or omitted.
	if v := d.GetRawConfig().AsValueMap()["compute_profile"].AsValueSlice()[0].AsValueMap()["virtual_machine_profile"].AsValueSlice()[0].AsValueMap()["os_profile"].AsValueSlice()[0].AsValueMap()["require_guest_provision_signal_enabled"]; !v.IsNull() {
		output.RequireGuestProvisionSignal = pointer.To(input.RequireGuestProvisionSignalEnabled)
	}

	return &output
}

func expandLinuxConfigurationModel(inputList []LinuxConfigurationModel) *fleets.LinuxConfiguration {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.LinuxConfiguration{
		DisablePasswordAuthentication: pointer.To(!input.PasswordAuthenticationEnabled),
		EnableVMAgentPlatformUpdates:  pointer.To(input.VMAgentPlatformUpdatesEnabled),
		PatchSettings:                 expandLinuxPatchSettingsModel(input.PatchSetting),
		ProvisionVMAgent:              pointer.To(input.ProvisionVMAgentEnabled),
		Ssh:                           expandSshConfigurationModel(input.SshKeys),
	}

	return &output
}

func expandLinuxPatchSettingsModel(inputList []LinuxPatchSettingModel) *fleets.LinuxPatchSettings {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.LinuxPatchSettings{}
	if input.PatchMode != "" {
		output.PatchMode = pointer.To(fleets.LinuxVMGuestPatchMode(input.PatchMode))
		if input.PatchMode == string(fleets.LinuxPatchAssessmentModeAutomaticByPlatform) {
			// AutomaticByPlatformSettings cannot be set if the PatchMode is not 'AutomaticByPlatform'.
			output.AutomaticByPlatformSettings = expandLinuxAutomaticByPlatformSettingModel(input.AutomaticByPlatformSetting)
		}
	}

	if input.AssessmentMode != "" {
		output.AssessmentMode = pointer.To(fleets.LinuxPatchAssessmentMode(input.AssessmentMode))
	}
	return &output
}

func expandLinuxAutomaticByPlatformSettingModel(inputList []AutomaticByPlatformSettingModel) *fleets.LinuxVMGuestPatchAutomaticByPlatformSettings {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.LinuxVMGuestPatchAutomaticByPlatformSettings{
		BypassPlatformSafetyChecksOnUserSchedule: pointer.To(input.BypassPlatformSafetyChecksEnabled),
	}

	if input.RebootSetting != "" {
		output.RebootSetting = pointer.To(fleets.LinuxVMGuestPatchAutomaticByPlatformRebootSetting(input.RebootSetting))
	}

	return &output
}

func expandSshConfigurationModel(inputList []SshKeyModel) *fleets.SshConfiguration {
	if len(inputList) == 0 {
		return nil
	}

	var publicKeys []fleets.SshPublicKey
	for _, v := range inputList {
		input := v
		output := fleets.SshPublicKey{
			Path: pointer.To(input.Path),
		}

		if input.KeyData != "" {
			output.KeyData = &input.KeyData
		}

		publicKeys = append(publicKeys, output)
	}

	output := fleets.SshConfiguration{
		PublicKeys: pointer.To(publicKeys),
	}

	return &output
}

func expandOsProfileSecretsModel(inputList []SecretModel) *[]fleets.VaultSecretGroup {
	if len(inputList) == 0 {
		return nil
	}

	var outputList []fleets.VaultSecretGroup
	for _, v := range inputList {
		input := v
		output := fleets.VaultSecretGroup{
			SourceVault:       expandSubResource(input.KeyVaultId),
			VaultCertificates: expandVaultCertificateModel(input.Certificates),
		}

		outputList = append(outputList, output)
	}
	return &outputList
}

func expandVaultCertificateModel(inputList []CertificateModel) *[]fleets.VaultCertificate {
	var outputList []fleets.VaultCertificate
	for _, v := range inputList {
		input := v
		output := fleets.VaultCertificate{}

		if input.Store != "" {
			output.CertificateStore = pointer.To(input.Store)
		}

		if input.Url != "" {
			output.CertificateURL = pointer.To(input.Url)
		}
		outputList = append(outputList, output)
	}
	return &outputList
}

func expandWindowsConfigurationModel(inputList []WindowsConfigurationModel) *fleets.WindowsConfiguration {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := fleets.WindowsConfiguration{
		AdditionalUnattendContent:    expandAdditionalUnattendContentModel(input.AdditionalUnattendContent),
		EnableAutomaticUpdates:       pointer.To(input.AutomaticUpdatesEnabled),
		EnableVMAgentPlatformUpdates: pointer.To(input.VMAgentPlatformUpdatesEnabled),
		PatchSettings:                expandWindowsPatchSettingModel(input.PatchSetting),
		ProvisionVMAgent:             pointer.To(input.ProvisionVMAgentEnabled),
		WinRM:                        expandWinRM(input.WinRM),
	}
	if input.TimeZone != "" {
		output.TimeZone = pointer.To(input.TimeZone)
	}

	return &output
}

func expandWindowsPatchSettingModel(inputList []WindowsPatchSettingModel) *fleets.PatchSettings {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := fleets.PatchSettings{
		AssessmentMode:              pointer.To(fleets.WindowsPatchAssessmentMode(input.AssessmentMode)),
		AutomaticByPlatformSettings: expandAutomaticByPlatformSettingModel(input.AutomaticByPlatformSetting),
		PatchMode:                   pointer.To(fleets.WindowsVMGuestPatchMode(input.PatchMode)),
		EnableHotpatching:           pointer.To(input.HotPatchingEnabled),
	}

	return &output
}

func expandAutomaticByPlatformSettingModel(inputList []AutomaticByPlatformSettingModel) *fleets.WindowsVMGuestPatchAutomaticByPlatformSettings {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := fleets.WindowsVMGuestPatchAutomaticByPlatformSettings{
		BypassPlatformSafetyChecksOnUserSchedule: pointer.To(input.BypassPlatformSafetyChecksEnabled),
	}

	if input.RebootSetting != "" {
		output.RebootSetting = pointer.To(fleets.WindowsVMGuestPatchAutomaticByPlatformRebootSetting(input.RebootSetting))
	}

	return &output
}

func expandAdditionalUnattendContentModel(inputList []AdditionalUnattendContentModel) *[]fleets.AdditionalUnattendContent {
	if len(inputList) == 0 {
		return nil
	}

	var outputList []fleets.AdditionalUnattendContent
	for _, v := range inputList {
		input := v
		output := fleets.AdditionalUnattendContent{
			ComponentName: pointer.To(fleets.ComponentName(input.ComponentName)),
			PassName:      pointer.To(fleets.PassName(input.PassName)),
			SettingName:   pointer.To(fleets.SettingNames(input.SettingName)),
			Content:       pointer.To(input.Content),
		}
		outputList = append(outputList, output)
	}
	return &outputList
}

func expandWinRM(inputList []WinRMModel) *fleets.WinRMConfiguration {
	if len(inputList) == 0 {
		return nil
	}

	var listenerList []fleets.WinRMListener
	for _, v := range inputList {
		input := v
		output := fleets.WinRMListener{
			Protocol: pointer.To(fleets.ProtocolTypes(input.Protocol)),
		}

		if input.CertificateUrl != "" {
			output.CertificateURL = pointer.To(input.CertificateUrl)
		}
		listenerList = append(listenerList, output)
	}

	return &fleets.WinRMConfiguration{
		Listeners: pointer.To(listenerList),
	}
}

func expandScheduledEventsProfile(input *VirtualMachineProfileModel) *fleets.ScheduledEventsProfile {
	if input == nil {
		return nil
	}
	return &fleets.ScheduledEventsProfile{
		OsImageNotificationProfile: &fleets.OSImageNotificationProfile{
			Enable:           pointer.To(input.ScheduledEventOsImageEnabled),
			NotBeforeTimeout: pointer.To(input.ScheduledEventOsImageTimeout),
		},

		TerminateNotificationProfile: &fleets.TerminateNotificationProfile{
			Enable:           pointer.To(input.ScheduledEventTerminationEnabled),
			NotBeforeTimeout: pointer.To(input.ScheduledEventTerminationTimeout),
		},
	}
}

func expandSecurityPostureReferenceModel(inputList []SecurityPostureReferenceModel) *fleets.SecurityPostureReference {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := fleets.SecurityPostureReference{
		ExcludeExtensions: pointer.To(input.ExcludeExtensions),
		IsOverridable:     pointer.To(input.OverrideEnabled),
	}
	if input.Id != "" {
		output.Id = pointer.To(input.Id)
	}

	return &output
}

func expandSecurityProfileModel(inputList []SecurityProfileModel) *fleets.SecurityProfile {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := fleets.SecurityProfile{
		EncryptionAtHost:   pointer.To(input.EncryptionAtHostEnabled),
		EncryptionIdentity: expandEncryptionIdentityModel(input.UserAssignedIdentityId),
		ProxyAgentSettings: expandProxyAgentModel(input.ProxyAgent),
		SecurityType:       pointer.To(fleets.SecurityTypes(input.SecurityType)),
		UefiSettings: &fleets.UefiSettings{
			SecureBootEnabled: pointer.To(input.UefiSecureBootEnabled),
			VTpmEnabled:       pointer.To(input.UefiVTpmEnabled),
		},
	}

	return &output
}

func expandEncryptionIdentityModel(id string) *fleets.EncryptionIdentity {
	if id == "" {
		return nil
	}

	return &fleets.EncryptionIdentity{
		UserAssignedIdentityResourceId: pointer.To(id),
	}
}

func expandProxyAgentModel(inputList []ProxyAgentModel) *fleets.ProxyAgentSettings {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := fleets.ProxyAgentSettings{
		Enabled:          pointer.To(true),
		KeyIncarnationId: pointer.To(input.KeyIncarnationValue),
		Mode:             pointer.To(fleets.Mode(input.mode)),
	}

	return &output
}

func expandStorageProfileModel(inputList []StorageProfileModel) *fleets.VirtualMachineScaleSetStorageProfile {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := fleets.VirtualMachineScaleSetStorageProfile{
		DataDisks:          expandDataDiskModel(input.DataDisks),
		DiskControllerType: pointer.To(fleets.DiskControllerTypes(input.DiskControllerType)),
		ImageReference:     expandImageReferenceModel(input.ImageReference),
		OsDisk:             expandOSDiskModel(input.OsDisk),
	}

	return &output
}

func expandDataDiskModel(inputList []DataDiskModel) *[]fleets.VirtualMachineScaleSetDataDisk {
	var outputList []fleets.VirtualMachineScaleSetDataDisk
	for _, v := range inputList {
		input := v
		output := fleets.VirtualMachineScaleSetDataDisk{
			CreateOption:            fleets.DiskCreateOptionTypes(input.CreateOption),
			DeleteOption:            pointer.To(fleets.DiskDeleteOptionTypes(input.DeleteOption)),
			DiskIOPSReadWrite:       pointer.To(input.DiskIOPSReadWrite),
			DiskMBpsReadWrite:       pointer.To(input.DiskMbpsReadWrite),
			Lun:                     input.Lun,
			ManagedDisk:             expandManagedDiskModel(input.ManagedDisk),
			WriteAcceleratorEnabled: pointer.To(input.WriteAcceleratorEnabled),
		}

		if input.DiskSizeInGB > 0 {
			output.DiskSizeGB = pointer.To(input.DiskSizeInGB)
		}

		caching := string(fleets.CachingTypesNone)
		if v := input.Caching; v != "" {
			caching = v
		}
		output.Caching = pointer.To(fleets.CachingTypes(caching))

		if input.Name != "" {
			output.Name = &input.Name
		}
		outputList = append(outputList, output)
	}
	return &outputList
}

func expandManagedDiskModel(inputList []ManagedDiskModel) *fleets.VirtualMachineScaleSetManagedDiskParameters {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.VirtualMachineScaleSetManagedDiskParameters{
		DiskEncryptionSet:  expandDiskEncryptionSetModel(input.DiskEncryptionSetId),
		SecurityProfile:    expandVMDiskSecurityProfileModel(input.SecurityEncryptionType, input.SecurityDiskEncryptionSetId),
		StorageAccountType: pointer.To(fleets.StorageAccountTypes(input.StorageAccountType)),
	}

	return &output
}

func expandDiskEncryptionSetModel(diskEncryptionSetId string) *fleets.DiskEncryptionSetParameters {
	if diskEncryptionSetId == "" {
		return nil
	}

	return &fleets.DiskEncryptionSetParameters{
		Id: pointer.To(diskEncryptionSetId),
	}
}

func expandVMDiskSecurityProfileModel(securityEncryptionType string, securityDiskEncryptionSetId string) *fleets.VMDiskSecurityProfile {
	if securityEncryptionType == "" && securityDiskEncryptionSetId == "" {
		return nil
	}

	output := fleets.VMDiskSecurityProfile{
		DiskEncryptionSet: expandDiskEncryptionSetModel(securityDiskEncryptionSetId),
	}
	if securityEncryptionType != "" {
		output.SecurityEncryptionType = pointer.To(fleets.SecurityEncryptionTypes(securityEncryptionType))
	}

	return &output
}

func expandImageReferenceModel(inputList []ImageReferenceModel) *fleets.ImageReference {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := fleets.ImageReference{}

	if input.CommunityGalleryImageId != "" {
		output.CommunityGalleryImageId = pointer.To(input.CommunityGalleryImageId)
	}

	if input.Id != "" {
		output.Id = pointer.To(input.Id)
	}

	if input.Offer != "" {
		output.Offer = pointer.To(input.Offer)
	}

	if input.Publisher != "" {
		output.Publisher = pointer.To(input.Publisher)
	}

	if input.SharedGalleryImageId != "" {
		output.SharedGalleryImageId = pointer.To(input.SharedGalleryImageId)
	}

	if input.Sku != "" {
		output.Sku = pointer.To(input.Sku)
	}

	if input.Version != "" {
		output.Version = pointer.To(input.Version)
	}

	return &output
}

func expandOSDiskModel(inputList []OSDiskModel) *fleets.VirtualMachineScaleSetOSDisk {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]

	output := fleets.VirtualMachineScaleSetOSDisk{
		CreateOption:            fleets.DiskCreateOptionTypes(input.CreateOption),
		OsType:                  pointer.To(fleets.OperatingSystemTypes(input.OsType)),
		DiffDiskSettings:        expandDiffDiskSettingsModel(input),
		ManagedDisk:             expandManagedDiskModel(input.ManagedDisk),
		WriteAcceleratorEnabled: pointer.To(input.WriteAcceleratorEnabled),
	}

	if input.DiskSizeInGB > 0 {
		output.DiskSizeGB = pointer.To(input.DiskSizeInGB)
	}

	caching := string(fleets.CachingTypesNone)
	if v := input.Caching; v != "" {
		caching = v
	}
	output.Caching = pointer.To(fleets.CachingTypes(caching))

	if input.DeleteOption != "" {
		output.DeleteOption = pointer.To(fleets.DiskDeleteOptionTypes(input.DeleteOption))
	}
	if input.ImageUri != "" {
		output.Image = expandImage(input.ImageUri)
	}
	if input.Name != "" {
		output.Name = pointer.To(input.Name)
	}
	if len(input.VhdContainers) > 0 {
		output.VhdContainers = pointer.To(input.VhdContainers)
	}

	return &output
}

func expandDiffDiskSettingsModel(input *OSDiskModel) *fleets.DiffDiskSettings {
	if input == nil || (input.DiffDiskOption == "" && input.DiffDiskPlacement == "") {
		return nil
	}

	output := fleets.DiffDiskSettings{}
	if input.DiffDiskOption != "" {
		output.Option = pointer.To(fleets.DiffDiskOptions(input.DiffDiskOption))
	}
	if input.DiffDiskPlacement != "" {
		output.Placement = pointer.To(fleets.DiffDiskPlacement(input.DiffDiskPlacement))
	}

	return &output
}

func expandImage(input string) *fleets.VirtualHardDisk {
	if input == "" {
		return nil
	}

	return &fleets.VirtualHardDisk{Uri: pointer.To(input)}
}

func flattenVirtualMachineProfileModel(input *fleets.BaseVirtualMachineProfile, metadata sdk.ResourceMetaData) ([]VirtualMachineProfileModel, error) {
	var outputList []VirtualMachineProfileModel
	if input == nil {
		return outputList, nil
	}
	output := VirtualMachineProfileModel{
		GalleryApplicationProfile: flattenApplicationProfileModel(input.ApplicationProfile),
		VMSize:                    flattenVMSizeModel(input.HardwareProfile),
		NetworkInterface:          flattenNetworkInterfaceModel(input.NetworkProfile),
		SecurityPostureReference:  flattenSecurityPostureReferenceModel(input.SecurityPostureReference),
		SecurityProfile:           flattenSecurityProfileModel(input.SecurityProfile),
		StorageProfile:            flattenStorageProfileModel(input.StorageProfile),
	}
	osProfile, err := flattenOSProfileModel(input.OsProfile, metadata.ResourceData)
	if err != nil {
		return outputList, err
	}
	output.OsProfile = osProfile

	if v := input.ServiceArtifactReference; v != nil {
		output.ServiceArtifactId = pointer.From(v.Id)
	}

	if se := input.ScheduledEventsProfile; se != nil {
		if v := se.TerminateNotificationProfile; v != nil {
			output.ScheduledEventTerminationEnabled = pointer.From(v.Enable)
			output.ScheduledEventTerminationTimeout = pointer.From(v.NotBeforeTimeout)
		}
		if v := se.OsImageNotificationProfile; v != nil {
			output.ScheduledEventOsImageEnabled = pointer.From(v.Enable)
			output.ScheduledEventOsImageTimeout = pointer.From(v.NotBeforeTimeout)
		}
	}

	if cr := input.CapacityReservation; cr != nil {
		if v := cr.CapacityReservationGroup; v != nil {
			output.CapacityReservationGroupId = pointer.From(v.Id)
		}
	}

	if dp := input.DiagnosticsProfile; dp != nil {
		if v := dp.BootDiagnostics; v != nil {
			output.BootDiagnosticEnabled = pointer.From(v.Enabled)
			output.BootDiagnosticStorageAccountEndpoint = pointer.From(v.StorageUri)
		}
	}

	if np := input.NetworkProfile; np != nil {
		if v := np.HealthProbe; v != nil {
			output.NetworkHealthProbeId = pointer.From(v.Id)
		}
		output.NetworkApiVersion = string(pointer.From(np.NetworkApiVersion))
	}

	extensionProfileValue, err := flattenExtensionModel(input.ExtensionProfile, metadata)
	if err != nil {
		return nil, err
	}
	output.Extensions = extensionProfileValue

	if input.LicenseType != nil {
		output.LicenseType = *input.LicenseType
	}

	if input.UserData != nil {
		output.UserDataBase64 = *input.UserData
	}

	return append(outputList, output), nil
}

func flattenLinuxConfigurationModel(input *fleets.LinuxConfiguration) []LinuxConfigurationModel {
	var outputList []LinuxConfigurationModel
	if input == nil {
		return outputList
	}

	output := LinuxConfigurationModel{}
	output.PasswordAuthenticationEnabled = !pointer.From(input.DisablePasswordAuthentication)
	output.VMAgentPlatformUpdatesEnabled = pointer.From(input.ProvisionVMAgent)
	output.ProvisionVMAgentEnabled = pointer.From(input.ProvisionVMAgent)
	output.SshKeys = flattenSshKeyModel(input.Ssh)
	output.PatchSetting = flattenLinuxPatchSettingModel(input.PatchSettings)

	return append(outputList, output)
}

func flattenLinuxPatchSettingModel(input *fleets.LinuxPatchSettings) []LinuxPatchSettingModel {
	var outputList []LinuxPatchSettingModel
	if input == nil {
		return outputList
	}

	output := LinuxPatchSettingModel{
		AutomaticByPlatformSetting: flattenLinuxAutomaticByPlatformSettingModel(input.AutomaticByPlatformSettings),
	}
	output.AssessmentMode = string(pointer.From(input.AssessmentMode))
	output.PatchMode = string(pointer.From(input.PatchMode))

	return append(outputList, output)
}

func flattenLinuxAutomaticByPlatformSettingModel(input *fleets.LinuxVMGuestPatchAutomaticByPlatformSettings) []AutomaticByPlatformSettingModel {
	var outputList []AutomaticByPlatformSettingModel
	if input == nil {
		return outputList
	}
	output := AutomaticByPlatformSettingModel{}
	output.BypassPlatformSafetyChecksEnabled = pointer.From(input.BypassPlatformSafetyChecksOnUserSchedule)
	output.RebootSetting = string(pointer.From(input.RebootSetting))

	return append(outputList, output)
}

func flattenSshKeyModel(input *fleets.SshConfiguration) []SshKeyModel {
	var outputList []SshKeyModel
	if input == nil || input.PublicKeys == nil {
		return outputList
	}

	for _, input := range *input.PublicKeys {
		output := SshKeyModel{}
		if input.KeyData != nil {
			output.KeyData = *input.KeyData
		}
		if input.Path != nil {
			output.Path = *input.Path
		}
		outputList = append(outputList, output)
	}

	return outputList
}

func flattenApplicationProfileModel(input *fleets.ApplicationProfile) []GalleryApplicationModel {
	var outputList []GalleryApplicationModel
	if input == nil {
		return outputList
	}

	for _, input := range *input.GalleryApplications {
		output := GalleryApplicationModel{}
		output.VersionId = input.PackageReferenceId
		output.ConfigurationBlobUri = pointer.From(input.ConfigurationReference)
		output.AutomaticUpgradeEnabled = pointer.From(input.EnableAutomaticUpgrade)
		output.Order = pointer.From(input.Order)
		output.Tags = pointer.From(input.Tags)
		output.TreatFailureAsDeploymentFailureEnabled = pointer.From(input.TreatFailureAsDeploymentFailure)

		outputList = append(outputList, output)
	}

	return outputList
}

func flattenVMSizeModel(input *fleets.VirtualMachineScaleSetHardwareProfile) []VMSizeModel {
	var outputList []VMSizeModel
	if input == nil {
		return outputList
	}

	output := VMSizeModel{}
	if props := input.VMSizeProperties; props != nil {
		output.VCPUAvailableCount = pointer.From(props.VCPUsAvailable)
		output.VCPUPerCoreCount = pointer.From(props.VCPUsPerCore)
	}

	return append(outputList, output)
}

func flattenNetworkInterfaceModel(input *fleets.VirtualMachineScaleSetNetworkProfile) []NetworkInterfaceModel {
	var outputList []NetworkInterfaceModel
	if input == nil {
		return outputList
	}

	for _, input := range *input.NetworkInterfaceConfigurations {
		output := NetworkInterfaceModel{
			Name: input.Name,
		}

		if props := input.Properties; props != nil {
			auxiliaryMode := ""
			if v := props.AuxiliaryMode; v != nil && *v != fleets.NetworkInterfaceAuxiliaryModeNone {
				auxiliaryMode = string(*v)
			}
			output.AuxiliaryMode = auxiliaryMode

			auxiliarySku := ""
			if v := props.AuxiliarySku; v != nil && *v != fleets.NetworkInterfaceAuxiliarySkuNone {
				auxiliarySku = string(*v)
			}
			output.AuxiliarySku = auxiliarySku

			output.DeleteOption = string(pointer.From(props.DeleteOption))

			output.TcpStateTrackingEnabled = !pointer.From(props.DisableTcpStateTracking)

			if v := props.DnsSettings; v != nil {
				output.DnsServers = pointer.From(v.DnsServers)
			}

			output.TcpStateTrackingEnabled = !pointer.From(props.DisableTcpStateTracking)

			output.AcceleratedNetworkingEnabled = pointer.From(props.EnableAcceleratedNetworking)

			output.FpgaEnabled = pointer.From(props.EnableFpga)

			output.IPForwardingEnabled = pointer.From(props.EnableIPForwarding)

			output.IPConfiguration = flattenIPConfigurationModel(props.IPConfigurations)

			if v := props.NetworkSecurityGroup; v != nil {
				output.NetworkSecurityGroupId = pointer.From(v.Id)
			}

			output.Primary = pointer.From(props.Primary)
		}

		outputList = append(outputList, output)
	}

	return outputList
}

func flattenOsProfileSecretsModel(inputList *[]fleets.VaultSecretGroup) []SecretModel {
	var outputList []SecretModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := SecretModel{
			Certificates: flattenVaultCertificateModel(input.VaultCertificates),
		}
		if v := input.SourceVault; v != nil {
			output.KeyVaultId = pointer.From(v.Id)
		}

		outputList = append(outputList, output)
	}
	return outputList
}

func flattenOSProfileModel(input *fleets.VirtualMachineScaleSetOSProfile, d *schema.ResourceData) ([]OSProfileModel, error) {
	var outputList []OSProfileModel
	if input == nil {
		return outputList, nil
	}
	output := OSProfileModel{
		LinuxConfiguration:   flattenLinuxConfigurationModel(input.LinuxConfiguration),
		Secret:               flattenOsProfileSecretsModel(input.Secrets),
		WindowsConfiguration: flattenWindowsConfigurationModel(input.WindowsConfiguration),
	}

	output.AdminPassword = d.Get("compute_profile.0.virtual_machine_profile.0.os_profile.0.admin_password").(string)
	output.AdminUsername = pointer.From(input.AdminUsername)
	output.ExtensionOperationsEnabled = pointer.From(input.AllowExtensionOperations)
	output.ComputerNamePrefix = pointer.From(input.ComputerNamePrefix)
	output.CustomDataBase64 = pointer.From(input.CustomData)
	output.RequireGuestProvisionSignalEnabled = pointer.From(input.RequireGuestProvisionSignal)

	return append(outputList, output), nil
}

func flattenVaultCertificateModel(inputList *[]fleets.VaultCertificate) []CertificateModel {
	var outputList []CertificateModel
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := CertificateModel{}

		output.Store = pointer.From(input.CertificateStore)
		output.Url = pointer.From(input.CertificateURL)

		outputList = append(outputList, output)
	}
	return outputList
}

func flattenWindowsConfigurationModel(input *fleets.WindowsConfiguration) []WindowsConfigurationModel {
	var outputList []WindowsConfigurationModel
	if input == nil {
		return outputList
	}
	output := WindowsConfigurationModel{
		AdditionalUnattendContent: flattenAdditionalUnattendContentModel(input.AdditionalUnattendContent),
		WinRM:                     flattenWinRMModel(input.WinRM),
	}

	output.AutomaticUpdatesEnabled = !pointer.From(input.EnableAutomaticUpdates)
	output.VMAgentPlatformUpdatesEnabled = pointer.From(input.ProvisionVMAgent)
	output.ProvisionVMAgentEnabled = pointer.From(input.ProvisionVMAgent)
	output.TimeZone = pointer.From(input.TimeZone)
	output.PatchSetting = flattenWindowsPatchSettingModel(input.PatchSettings)

	return append(outputList, output)
}

func flattenWindowsPatchSettingModel(input *fleets.PatchSettings) []WindowsPatchSettingModel {
	var outputList []WindowsPatchSettingModel
	if input == nil {
		return outputList
	}

	output := WindowsPatchSettingModel{
		AutomaticByPlatformSetting: flattenWindowsAutomaticByPlatformSettingModel(input.AutomaticByPlatformSettings),
	}
	output.AssessmentMode = string(pointer.From(input.AssessmentMode))
	output.PatchMode = string(pointer.From(input.PatchMode))
	output.HotPatchingEnabled = pointer.From(input.EnableHotpatching)

	return append(outputList, output)
}

func flattenWindowsAutomaticByPlatformSettingModel(input *fleets.WindowsVMGuestPatchAutomaticByPlatformSettings) []AutomaticByPlatformSettingModel {
	var outputList []AutomaticByPlatformSettingModel
	if input == nil {
		return outputList
	}
	output := AutomaticByPlatformSettingModel{}
	output.BypassPlatformSafetyChecksEnabled = pointer.From(input.BypassPlatformSafetyChecksOnUserSchedule)
	output.RebootSetting = string(pointer.From(input.RebootSetting))

	return append(outputList, output)
}

func flattenAdditionalUnattendContentModel(inputList *[]fleets.AdditionalUnattendContent) []AdditionalUnattendContentModel {
	var outputList []AdditionalUnattendContentModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := AdditionalUnattendContentModel{}
		output.ComponentName = string(pointer.From(input.ComponentName))
		output.Content = pointer.From(input.Content)
		output.PassName = string(pointer.From(input.PassName))
		output.SettingName = string(pointer.From(input.SettingName))
		outputList = append(outputList, output)
	}
	return outputList
}

func flattenWinRMModel(input *fleets.WinRMConfiguration) []WinRMModel {
	var outputList []WinRMModel
	if input == nil || input.Listeners == nil {
		return outputList
	}

	for _, input := range *input.Listeners {
		output := WinRMModel{}
		output.CertificateUrl = pointer.From(input.CertificateURL)
		output.Protocol = string(pointer.From(input.Protocol))
		outputList = append(outputList, output)
	}

	return outputList
}

func flattenSecurityPostureReferenceModel(input *fleets.SecurityPostureReference) []SecurityPostureReferenceModel {
	var outputList []SecurityPostureReferenceModel
	if input == nil {
		return outputList
	}
	output := SecurityPostureReferenceModel{}
	output.ExcludeExtensions = pointer.From(input.ExcludeExtensions)
	output.Id = pointer.From(input.Id)
	output.OverrideEnabled = pointer.From(input.IsOverridable)

	return append(outputList, output)
}

func flattenSecurityProfileModel(input *fleets.SecurityProfile) []SecurityProfileModel {
	var outputList []SecurityProfileModel
	if input == nil {
		return outputList
	}
	output := SecurityProfileModel{
		ProxyAgent: flattenProxyAgentModel(input.ProxyAgentSettings),
	}

	output.EncryptionAtHostEnabled = pointer.From(input.EncryptionAtHost)

	if v := input.EncryptionIdentity; v != nil {
		output.UserAssignedIdentityId = pointer.From(v.UserAssignedIdentityResourceId)
	}

	if v := input.UefiSettings; v != nil {
		output.UefiSecureBootEnabled = pointer.From(v.SecureBootEnabled)
		output.UefiVTpmEnabled = pointer.From(v.VTpmEnabled)
	}

	output.SecurityType = string(pointer.From(input.SecurityType))

	return append(outputList, output)
}

func flattenProxyAgentModel(input *fleets.ProxyAgentSettings) []ProxyAgentModel {
	var outputList []ProxyAgentModel
	if input == nil {
		return outputList
	}

	output := ProxyAgentModel{}
	output.KeyIncarnationValue = pointer.From(input.KeyIncarnationId)
	output.mode = string(pointer.From(input.Mode))

	return append(outputList, output)
}

func flattenStorageProfileModel(input *fleets.VirtualMachineScaleSetStorageProfile) []StorageProfileModel {
	var outputList []StorageProfileModel
	if input == nil {
		return outputList
	}
	output := StorageProfileModel{
		DataDisks:      flattenDataDiskModel(input.DataDisks),
		ImageReference: flattenImageReferenceModel(input.ImageReference),
		OsDisk:         flattenOSDiskModel(input.OsDisk),
	}

	output.DiskControllerType = string(pointer.From(input.DiskControllerType))

	return append(outputList, output)
}

func flattenDataDiskModel(inputList *[]fleets.VirtualMachineScaleSetDataDisk) []DataDiskModel {
	var outputList []DataDiskModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := DataDiskModel{
			CreateOption: string(input.CreateOption),
			Lun:          input.Lun,
			ManagedDisk:  flattenManagedDiskModel(input.ManagedDisk),
		}

		caching := ""
		if v := input.Caching; v != nil && *v != fleets.CachingTypesNone {
			caching = string(*v)
		}
		output.Caching = caching

		output.DeleteOption = string(pointer.From(input.DeleteOption))
		output.DiskIOPSReadWrite = pointer.From(input.DiskIOPSReadWrite)
		output.DiskMbpsReadWrite = pointer.From(input.DiskMBpsReadWrite)
		output.DiskSizeInGB = pointer.From(input.DiskSizeGB)
		output.Name = pointer.From(input.Name)
		output.WriteAcceleratorEnabled = pointer.From(input.WriteAcceleratorEnabled)

		outputList = append(outputList, output)
	}
	return outputList
}

func flattenManagedDiskModel(input *fleets.VirtualMachineScaleSetManagedDiskParameters) []ManagedDiskModel {
	var outputList []ManagedDiskModel
	if input == nil {
		return outputList
	}

	output := ManagedDiskModel{}
	if v := input.DiskEncryptionSet; v != nil {
		output.DiskEncryptionSetId = pointer.From(v.Id)
	}

	if sp := input.SecurityProfile; sp != nil {
		if v := sp.DiskEncryptionSet; v != nil {
			output.SecurityDiskEncryptionSetId = pointer.From(v.Id)
		}
		output.SecurityEncryptionType = string(pointer.From(sp.SecurityEncryptionType))
	}
	output.StorageAccountType = string(pointer.From(input.StorageAccountType))

	return append(outputList, output)
}

func flattenImageReferenceModel(input *fleets.ImageReference) []ImageReferenceModel {
	var outputList []ImageReferenceModel
	if input == nil {
		return outputList
	}

	output := ImageReferenceModel{}
	output.CommunityGalleryImageId = pointer.From(input.CommunityGalleryImageId)
	output.Version = pointer.From(input.ExactVersion)
	output.Id = pointer.From(input.Id)
	output.Offer = pointer.From(input.Offer)
	output.Publisher = pointer.From(input.Publisher)
	output.SharedGalleryImageId = pointer.From(input.SharedGalleryImageId)
	output.Sku = pointer.From(input.Sku)
	output.Version = pointer.From(input.Version)

	return append(outputList, output)
}

func flattenOSDiskModel(input *fleets.VirtualMachineScaleSetOSDisk) []OSDiskModel {
	var outputList []OSDiskModel
	if input == nil {
		return outputList
	}

	output := OSDiskModel{
		ManagedDisk: flattenManagedDiskModel(input.ManagedDisk),
	}
	if v := input.DiffDiskSettings; v != nil {
		output.DiffDiskOption = string(pointer.From(v.Option))
		output.DiffDiskPlacement = string(pointer.From(v.Placement))
	}

	if v := input.Image; v != nil {
		output.ImageUri = pointer.From(v.Uri)
	}
	output.CreateOption = string(input.CreateOption)

	caching := ""
	if v := input.Caching; v != nil && *v != fleets.CachingTypesNone {
		caching = string(*v)
	}
	output.Caching = caching
	output.DeleteOption = string(pointer.From(input.DeleteOption))
	output.DiskSizeInGB = pointer.From(input.DiskSizeGB)
	output.Name = pointer.From(input.Name)
	output.OsType = string(pointer.From(input.OsType))
	output.VhdContainers = pointer.From(input.VhdContainers)

	//vhdContainers := make([]string, 0)
	//if v := input.VhdContainers; v != nil {
	//	vhdContainers = pointer.From(input.VhdContainers)
	//}
	output.VhdContainers = pointer.From(input.VhdContainers)
	output.WriteAcceleratorEnabled = pointer.From(input.WriteAcceleratorEnabled)

	return append(outputList, output)
}

func flattenIPConfigurationModel(inputList []fleets.VirtualMachineScaleSetIPConfiguration) []IPConfigurationModel {
	var outputList []IPConfigurationModel
	if len(inputList) == 0 {
		return outputList
	}
	for _, input := range inputList {
		output := IPConfigurationModel{
			Name: input.Name,
		}
		if props := input.Properties; props != nil {
			output.Primary = pointer.From(props.Primary)
			output.Version = string(pointer.From(props.PrivateIPAddressVersion))

			addressPools := make([]string, 0)
			if v := props.ApplicationGatewayBackendAddressPools; v != nil {
				addressPools = flattenSubResourceId(*v)
			}
			output.ApplicationGatewayBackendAddressPoolIds = addressPools

			lbAddressPools := make([]string, 0)
			if v := props.LoadBalancerBackendAddressPools; v != nil {
				lbAddressPools = flattenSubResourceId(*v)
			}
			output.LoadBalancerBackendAddressPoolIds = lbAddressPools

			groupIds := make([]string, 0)
			if v := props.ApplicationSecurityGroups; v != nil {
				groupIds = flattenSubResourceId(*v)
			}
			output.ApplicationSecurityGroupIds = groupIds

			natPools := make([]string, 0)
			if v := props.LoadBalancerInboundNatPools; v != nil {
				natPools = flattenSubResourceId(*v)
			}
			output.LoadBalancerInboundNatPoolIds = natPools

			if v := props.PublicIPAddressConfiguration; v != nil {
				output.PublicIPAddress = flattenPublicIPAddressModel(v)
			}

			if v := props.Subnet; v != nil {
				output.SubnetId = pointer.From(v.Id)
			}
		}

		outputList = append(outputList, output)
	}
	return outputList
}
