// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azurefleet

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"

	"github.com/hashicorp/go-azure-sdk/resource-manager/azurefleet/2024-11-01/fleets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservationgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplicationversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-07-01/virtualmachinescalesets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/applicationsecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networksecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/publicipprefixes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// "location": commonschema.LocationWithoutForceNew(),
func virtualMachineProfileSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"gallery_applications": galleryApplicationsProfileSchema(),

				"capacity_reservation_group_id": {
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: capacityreservationgroups.ValidateCapacityReservationGroupID,
					},
				},

				"boot_diagnostics": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"enabled": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
							},

							"storage_account_endpoint": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},

				"extensions": extensionProfileSchema(),

				"extensions_time_budget": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      "PT1H30M",
					ValidateFunc: azValidate.ISO8601DurationBetween("PT15M", "PT2H"),
				},

				"vm_size": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"vcp_us_available": {
								Type:     pluginsdk.TypeInt,
								Optional: true,
							},

							"vcp_us_per_core": {
								Type:     pluginsdk.TypeInt,
								Optional: true,
							},
						},
					},
				},

				"license_type": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"network_interface": networkInterfaceSchema(),

				"health_probe_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: azure.ValidateResourceID,
				},

				"network_api_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(fleets.NetworkApiVersionTwoZeroTwoZeroNegativeOneOneNegativeZeroOne),
					}, false),
				},

				"os_profile": osProfileSchema(),

				"os_image_notification_profile": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"enable": {
								Type:     pluginsdk.TypeBool,
								Required: true,
							},

							"not_before_timeout": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},

				"termination_notification_profile": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"enable": {
								Type:     pluginsdk.TypeBool,
								Required: true,
							},

							"not_before_timeout": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
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

							"is_overridable": {
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

				"storage_profile": storageProfileSchema(),

				"user_data": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsBase64,
				},
			},
		},
	}
}

func galleryApplicationsProfileSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"version_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: galleryapplicationversions.ValidateApplicationVersionID,
				},

				"configuration_blob_uri": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"order": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      0,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(0, 2147483647),
				},

				"tags": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"automatic_upgrade_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},

				"treat_failure_as_deployment_failure": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},
			},
		},
	}
}

func extensionProfileSchema() *pluginsdk.Schema {
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

				"auto_upgrade_minor_version": {
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

				"protected_settings": {
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
						Type: pluginsdk.TypeString,
					},
				},

				"settings": {
					Type:             pluginsdk.TypeString,
					Optional:         true,
					ValidateFunc:     validation.StringIsJSON,
					DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
				},

				"suppress_failures": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},
			},
		},
	}
}

func networkInterfaceSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"ip_configurations": {
					Type:     pluginsdk.TypeList,
					Required: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"name": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
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
								Set: pluginsdk.HashString,
							},

							"load_balancer_backend_address_pool_ids": {
								Type:     pluginsdk.TypeSet,
								Optional: true,
								Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
								Set:      pluginsdk.HashString,
							},

							"load_balancer_inbound_nat_rules_ids": {
								Type:     pluginsdk.TypeSet,
								Optional: true,
								Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
								Set:      pluginsdk.HashString,
							},

							"primary": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
								Default:  false,
							},

							"public_ip_address": {
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
											Type:     pluginsdk.TypeString,
											Optional: true,
											ValidateFunc: validation.StringInSlice([]string{
												string(fleets.DeleteOptionsDelete),
												string(fleets.DeleteOptionsDetach),
											}, false),
										},

										"domain_name_label": {
											Type:         pluginsdk.TypeString,
											Optional:     true,
											ValidateFunc: validation.StringIsNotEmpty,
										},

										"domain_name_label_scope": {
											Type:     pluginsdk.TypeString,
											Optional: true,
											ValidateFunc: validation.StringInSlice([]string{
												string(fleets.DomainNameLabelScopeTypesSubscriptionReuse),
												string(fleets.DomainNameLabelScopeTypesResourceGroupReuse),
												string(fleets.DomainNameLabelScopeTypesNoReuse),
												string(fleets.DomainNameLabelScopeTypesTenantReuse),
											}, false),
										},

										"ip_tags": {
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

										"idle_timeout_in_minutes": {
											Type:     pluginsdk.TypeInt,
											Optional: true,
										},

										"version": {
											Type:     pluginsdk.TypeString,
											Optional: true,
											ForceNew: true,
											Default:  string(virtualmachinescalesets.IPVersionIPvFour),
											ValidateFunc: validation.StringInSlice([]string{
												string(fleets.IPVersionIPvFour),
												string(fleets.IPVersionIPvSix),
											}, false),
										},

										"public_ip_prefix_id": {
											Type:         pluginsdk.TypeList,
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
														Type:     pluginsdk.TypeString,
														Optional: true,
														ValidateFunc: validation.StringInSlice([]string{
															string(fleets.PublicIPAddressSkuNameBasic),
															string(fleets.PublicIPAddressSkuNameStandard),
														}, false),
													},

													"tier": {
														Type:     pluginsdk.TypeString,
														Optional: true,
														ValidateFunc: validation.StringInSlice([]string{
															string(fleets.PublicIPAddressSkuTierRegional),
															string(fleets.PublicIPAddressSkuTierGlobal),
														}, false),
													},
												},
											},
										},
									},
								},
							},

							"subnet_id": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: commonids.ValidateSubnetID,
							},

							"version": {
								Type:     pluginsdk.TypeString,
								Optional: true,
								ValidateFunc: validation.StringInSlice([]string{
									string(fleets.IPVersionIPvFour),
									string(fleets.IPVersionIPvSix),
								}, false),
							},
						},
					},
				},

				"dns_servers": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},

				"accelerated_networking_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"ip_forwarding_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"network_security_group_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: networksecuritygroups.ValidateNetworkSecurityGroupID,
				},

				"primary": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"auxiliary_mode": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(fleets.NetworkInterfaceAuxiliaryModeNone),
						string(fleets.NetworkInterfaceAuxiliaryModeAcceleratedConnections),
						string(fleets.NetworkInterfaceAuxiliaryModeFloating),
					}, false),
				},

				"auxiliary_sku": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(fleets.NetworkInterfaceAuxiliarySkuATwo),
						string(fleets.NetworkInterfaceAuxiliarySkuAFour),
						string(fleets.NetworkInterfaceAuxiliarySkuAEight),
						string(fleets.NetworkInterfaceAuxiliarySkuNone),
						string(fleets.NetworkInterfaceAuxiliarySkuAOne),
					}, false),
				},

				"delete_option": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(fleets.DeleteOptionsDelete),
						string(fleets.DeleteOptionsDetach),
					}, false),
				},

				"tcp_state_tracking_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},

				"fpga_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	}
}

func osProfileSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"admin_username": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"admin_password": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"extension_operations_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  true,
					ForceNew: true,
				},

				"computer_name_prefix": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"custom_data": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsBase64,
				},

				"linux_configuration": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"password_authentication_enabled": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
							},

							"vm_agent_platform_updates_enabled": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
							},

							"assessment_mode": {
								Type:     pluginsdk.TypeString,
								Optional: true,
								ValidateFunc: validation.StringInSlice([]string{
									string(fleets.LinuxPatchAssessmentModeAutomaticByPlatform),
									string(fleets.LinuxPatchAssessmentModeImageDefault),
								}, false),
							},

							"bypass_platform_safety_checks_on_user_schedule": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
							},

							"reboot_setting": {
								Type:     pluginsdk.TypeString,
								Optional: true,
								ValidateFunc: validation.StringInSlice([]string{
									string(fleets.LinuxVMGuestPatchAutomaticByPlatformRebootSettingIfRequired),
									string(fleets.LinuxVMGuestPatchAutomaticByPlatformRebootSettingNever),
									string(fleets.LinuxVMGuestPatchAutomaticByPlatformRebootSettingAlways),
									string(fleets.LinuxVMGuestPatchAutomaticByPlatformRebootSettingUnknown),
								}, false),
							},

							"patch_mode": {
								Type:     pluginsdk.TypeString,
								Optional: true,
								ValidateFunc: validation.StringInSlice([]string{
									string(fleets.LinuxVMGuestPatchModeImageDefault),
									string(fleets.LinuxVMGuestPatchModeAutomaticByPlatform),
								}, false),
							},

							"provision_vm_agent": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
								Default:  true,
								ForceNew: true,
							},

							"ssh_public_keys": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"key_data": {
											Type:         pluginsdk.TypeString,
											Optional:     true,
											ValidateFunc: validation.StringIsNotEmpty,
										},

										"path": {
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

				"require_guest_provision_signal": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},

				"secrets": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"source_vault_ids": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								Elem: &pluginsdk.Schema{
									Type:         pluginsdk.TypeString,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},

							"certificates": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"store": {
											Type:     pluginsdk.TypeString,
											Optional: true,
										},

										"url": {
											Type:         pluginsdk.TypeString,
											Required:     true,
											ValidateFunc: keyVaultValidate.NestedItemId,
										},
									},
								},
							},
						},
					},
				},

				"windows_configuration": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"additional_unattend_content": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"component_name": {
											Type:     pluginsdk.TypeString,
											Optional: true,
											ValidateFunc: validation.StringInSlice([]string{
												string(fleets.ComponentNameMicrosoftNegativeWindowsNegativeShellNegativeSetup),
											}, false),
										},

										"content": {
											Type:         pluginsdk.TypeString,
											Optional:     true,
											ValidateFunc: validation.StringIsNotEmpty,
										},

										"pass_name": {
											Type:     pluginsdk.TypeString,
											Optional: true,
											ValidateFunc: validation.StringInSlice([]string{
												string(fleets.PassNameOobeSystem),
											}, false),
										},

										"setting_name": {
											Type:     pluginsdk.TypeString,
											Optional: true,
											ValidateFunc: validation.StringInSlice([]string{
												string(fleets.SettingNamesAutoLogon),
												string(fleets.SettingNamesFirstLogonCommands),
											}, false),
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

							"assessment_mode": {
								Type:     pluginsdk.TypeString,
								Optional: true,
								ValidateFunc: validation.StringInSlice([]string{
									string(fleets.WindowsPatchAssessmentModeImageDefault),
									string(fleets.WindowsPatchAssessmentModeAutomaticByPlatform),
								}, false),
							},

							"bypass_platform_safety_checks_on_user_schedule": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
							},

							"reboot_setting": {
								Type:     pluginsdk.TypeString,
								Optional: true,
								ValidateFunc: validation.StringInSlice([]string{
									string(fleets.WindowsVMGuestPatchAutomaticByPlatformRebootSettingUnknown),
									string(fleets.WindowsVMGuestPatchAutomaticByPlatformRebootSettingIfRequired),
									string(fleets.WindowsVMGuestPatchAutomaticByPlatformRebootSettingNever),
									string(fleets.WindowsVMGuestPatchAutomaticByPlatformRebootSettingAlways),
								}, false),
							},

							"hot_patching_enabled": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
							},

							"patch_mode": {
								Type:     pluginsdk.TypeString,
								Optional: true,
								ValidateFunc: validation.StringInSlice([]string{
									string(fleets.WindowsVMGuestPatchModeManual),
									string(fleets.WindowsVMGuestPatchModeAutomaticByOS),
									string(fleets.WindowsVMGuestPatchModeAutomaticByPlatform),
								}, false),
							},

							"provision_vm_agent": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
							},

							"time_zone": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},

							"win_rm_listeners": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"protocol": {
											Type:     pluginsdk.TypeString,
											Required: true,
											ForceNew: true,
											ValidateFunc: validation.StringInSlice([]string{
												string(fleets.ProtocolTypesHTTP),
												string(fleets.ProtocolTypesHTTPS),
											}, false),
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

func securityProfileSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"encryption_at_host": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},

				"user_assigned_identity_resource_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"proxy_agent_settings": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"enabled": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
							},

							"key_incarnation_id": {
								Type:     pluginsdk.TypeInt,
								Optional: true,
							},

							"mode": {
								Type:     pluginsdk.TypeString,
								Optional: true,
								ValidateFunc: validation.StringInSlice([]string{
									string(fleets.ModeAudit),
									string(fleets.ModeEnforce),
								}, false),
							},
						},
					},
				},

				"security_type": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(fleets.SecurityTypesTrustedLaunch),
						string(fleets.SecurityTypesConfidentialVM),
					}, false),
				},

				"secure_boot_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ForceNew: true,
				},

				"vtpm_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ForceNew: true,
				},
			},
		},
	}
}

func managedDiskSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"disk_encryption_set_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validate.DiskEncryptionSetID,
				},

				"security_profile": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"disk_encryption_set_id": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validate.DiskEncryptionSetID,
							},

							"security_encryption_type": {
								Type:     pluginsdk.TypeString,
								Optional: true,
								ValidateFunc: validation.StringInSlice([]string{
									string(fleets.SecurityEncryptionTypesVMGuestStateOnly),
									string(fleets.SecurityEncryptionTypesDiskWithVMGuestState),
									string(fleets.SecurityEncryptionTypesNonPersistedTPM),
								}, false),
							},
						},
					},
				},

				"storage_account_type": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(fleets.StorageAccountTypesPremiumVTwoLRS),
						string(fleets.StorageAccountTypesStandardLRS),
						string(fleets.StorageAccountTypesPremiumLRS),
						string(fleets.StorageAccountTypesStandardSSDLRS),
						string(fleets.StorageAccountTypesUltraSSDLRS),
						string(fleets.StorageAccountTypesPremiumZRS),
						string(fleets.StorageAccountTypesStandardSSDZRS),
					}, false),
				},
			},
		},
	}
}

func storageProfileSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"data_disks": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"caching": {
								Type:     pluginsdk.TypeString,
								Optional: true,
								ValidateFunc: validation.StringInSlice([]string{
									string(fleets.CachingTypesNone),
									string(fleets.CachingTypesReadOnly),
									string(fleets.CachingTypesReadWrite),
								}, false),
							},

							"create_option": {
								Type:     pluginsdk.TypeString,
								Required: true,
								ValidateFunc: validation.StringInSlice([]string{
									string(fleets.DiskCreateOptionTypesRestore),
									string(fleets.DiskCreateOptionTypesFromImage),
									string(fleets.DiskCreateOptionTypesEmpty),
									string(fleets.DiskCreateOptionTypesAttach),
									string(fleets.DiskCreateOptionTypesCopy),
								}, false),
							},

							"delete_option": {
								Type:     pluginsdk.TypeString,
								Optional: true,
								ValidateFunc: validation.StringInSlice([]string{
									string(fleets.DiskDeleteOptionTypesDetach),
									string(fleets.DiskDeleteOptionTypesDelete),
								}, false),
							},

							"disk_iops_read_write": {
								Type:     pluginsdk.TypeInt,
								Optional: true,
							},

							"disk_m_bps_read_write": {
								Type:     pluginsdk.TypeInt,
								Optional: true,
							},

							"disk_size_in_gb": {
								Type:     pluginsdk.TypeInt,
								Optional: true,
							},

							"lun": {
								Type:     pluginsdk.TypeInt,
								Required: true,
							},

							"managed_disk": managedDiskSchema(),

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
				},

				"disk_controller_type": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(fleets.DiskControllerTypesNVMe),
						string(fleets.DiskControllerTypesSCSI),
					}, false),
				},

				"image_reference": {
					Type:     pluginsdk.TypeList,
					Optional: true,
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
							"community_gallery_image_id": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},

							"exact_version": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},

							"id": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},

							"shared_gallery_image_id": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},

				"os_disk": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"caching": {
								Type:     pluginsdk.TypeString,
								Required: true,
								ValidateFunc: validation.StringInSlice([]string{
									string(fleets.CachingTypesReadOnly),
									string(fleets.CachingTypesReadWrite),
									string(fleets.CachingTypesNone),
								}, false),
							},

							"create_option": {
								Type:     pluginsdk.TypeString,
								Required: true,
								ValidateFunc: validation.StringInSlice([]string{
									string(fleets.DiskCreateOptionTypesCopy),
									string(fleets.DiskCreateOptionTypesRestore),
									string(fleets.DiskCreateOptionTypesFromImage),
									string(fleets.DiskCreateOptionTypesEmpty),
									string(fleets.DiskCreateOptionTypesAttach),
								}, false),
							},

							"delete_option": {
								Type:     pluginsdk.TypeString,
								Optional: true,
								ValidateFunc: validation.StringInSlice([]string{
									string(fleets.DiskDeleteOptionTypesDelete),
									string(fleets.DiskDeleteOptionTypesDetach),
								}, false),
							},

							"diff_disk_settings": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								ForceNew: true,
								MaxItems: 1,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"option": {
											Type:     pluginsdk.TypeString,
											Required: true,
											ForceNew: true,
											ValidateFunc: validation.StringInSlice([]string{
												string(fleets.DiffDiskOptionsLocal),
											}, false),
										},

										"placement": {
											Type:     pluginsdk.TypeString,
											Optional: true,
											ForceNew: true,
											ValidateFunc: validation.StringInSlice([]string{
												string(fleets.DiffDiskPlacementNVMeDisk),
												string(fleets.DiffDiskPlacementCacheDisk),
												string(fleets.DiffDiskPlacementResourceDisk),
											}, false),
										},
									},
								},
							},

							"disk_size_in_gb": {
								Type:     pluginsdk.TypeInt,
								Optional: true,
							},

							"image_uri": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},

							"managed_disk": managedDiskSchema(),

							"name": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},

							"os_type": {
								Type:     pluginsdk.TypeString,
								Optional: true,
								ValidateFunc: validation.StringInSlice([]string{
									string(fleets.OperatingSystemTypesWindows),
									string(fleets.OperatingSystemTypesLinux),
								}, false),
							},

							"vhd_containers": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								Elem: &pluginsdk.Schema{
									Type: pluginsdk.TypeString,
								},
							},

							"write_accelerator_enabled": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
								Default:  false,
							},
						},
					},
				},
			},
		},
	}
}

func applicationProfileSchema4() *pluginsdk.Schema {
	return &pluginsdk.Schema{}
}

func applicationProfileSchema5() *pluginsdk.Schema {
	return &pluginsdk.Schema{}
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
