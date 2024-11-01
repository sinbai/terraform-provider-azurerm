package azurefleet

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurefleet/2024-11-01/fleets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservationgroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type AzureFleetResourceModel struct {
	Name                       string                                     `tfschema:"name"`
	ResourceGroupName          string                                     `tfschema:"resource_group_name"`
	Identity                   []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Location                   string                                     `tfschema:"location"`
	Plan                       []PlanModel                                `tfschema:"plan"`
	AdditionalLocationsProfile []AdditionalLocationsProfileModel          `tfschema:"additional_locations_profile"`
	ComputeProfile             []ComputeProfileModel                      `tfschema:"compute_profile"`
	RegularPriorityProfile     []RegularPriorityProfileModel              `tfschema:"regular_priority_profile"`
	SpotPriorityProfile        []SpotPriorityProfileModel                 `tfschema:"spot_priority_profile"`
	UniqueId                   string                                     `tfschema:"unique_id"`
	VMAttributes               VMAttributesModel                          `tfschema:"vm_attributes"`
	VMSizesProfile             []VMSizeProfileModel                       `tfschema:"vm_sizes_profile"`
	Tags                       map[string]string                          `tfschema:"tags"`
	Zones                      []string                                   `tfschema:"zones"`
}

type AdditionalLocationsProfileModel struct {
	Location              string                       `tfschema:"location"`
	VirtualMachineProfile []VirtualMachineProfileModel `tfschema:"virtual_machine_profile"`
}

type VirtualMachineProfileModel struct {
	VMGalleryApplicationProfile    []VMGalleryApplicationsModel          `tfschema:"gallery_applications"`
	CapacityReservationGroupId     string                                `tfschema:"capacity_reservation_group_id"`
	BootDiagnostics                []BootDiagnosticsModel                `tfschema:"boot_diagnostics"`
	Extensions                     []ExtensionsModel                     `tfschema:"extensions"`
	ExtensionsTimeBudget           string                                `tfschema:"extensions_time_budget"`
	VMSize                         []VMSizeModel                         `tfschema:"vm_size"`
	LicenseType                    string                                `tfschema:"license_type"`
	NetworkProfile                 []NetworkProfileModel                 `tfschema:"network_profile"`
	OsProfile                      []OSProfileModel                      `tfschema:"os_profile"`
	OsImageNotificationProfile     []OSImageNotificationProfileModel     `tfschema:"os_image_notification_profile"`
	TerminationNotificationProfile []TerminationNotificationProfileModel `tfschema:"termination_notification_profile"`
	SecurityPostureReference       []SecurityPostureReferenceModel       `tfschema:"security_posture_reference"`
	SecurityProfile                []SecurityProfileModel                `tfschema:"security_profile"`
	ServiceArtifactId              string                                `tfschema:"service_artifact_id"`
	StorageProfile                 []StorageProfileModel                 `tfschema:"storage_profile"`
	UserData                       string                                `tfschema:"user_data"`
}

type VMGalleryApplicationsModel struct {
	ConfigurationReference          string `tfschema:"configuration_blob_uri"`
	EnableAutomaticUpgrade          bool   `tfschema:"automatic_upgrade_enabled"`
	Order                           int64  `tfschema:"order"`
	PackageReferenceId              string `tfschema:"version_id"`
	Tags                            string `tfschema:"tags"`
	TreatFailureAsDeploymentFailure bool   `tfschema:"treat_failure_as_deployment_failure"`
}

type BootDiagnosticsModel struct {
	Enabled                bool   `tfschema:"enabled"`
	StorageAccountEndpoint string `tfschema:"storage_account_endpoint"`
}

type ExtensionsModel struct {
	Name                          string                               `tfschema:"name"`
	Publisher                     string                               `tfschema:"publisher"`
	Type                          string                               `tfschema:"type"`
	TypeHandlerVersion            string                               `tfschema:"type_handler_version"`
	AutoUpgradeMinorVersion       bool                                 `tfschema:"auto_upgrade_minor_version"`
	AutomaticUpgradeEnabled       bool                                 `tfschema:"automatic_upgrade_enabled"`
	ForceUpdateTag                string                               `tfschema:"force_update_tag"`
	ProtectedSettings             map[string]interface{}               `tfschema:"protected_settings"`
	ProtectedSettingsFromKeyVault []ProtectedSettingsFromKeyVaultModel `tfschema:"protected_settings_from_key_vault"`
	ProvisionAfterExtensions      []string                             `tfschema:"provision_after_extensions"`
	Settings                      string                               `tfschema:"settings"`
	SuppressFailures              bool                                 `tfschema:"suppress_failures"`
}

type ProtectedSettingsFromKeyVaultModel struct {
	SecretUrl     string `tfschema:"secret_url"`
	SourceVaultId string `tfschema:"source_vault_id"`
}

type VMSizeModel struct {
	VCPUsAvailable int64 `tfschema:"vcp_us_available"`
	VCPUsPerCore   int64 `tfschema:"vcp_us_per_core"`
}

type NetworkProfileModel struct {
	HealthProbe                    []ApiEntityReferenceModel   `tfschema:"health_probe"`
	NetworkApiVersion              string                      `tfschema:"network_api_version"`
	NetworkInterfaceConfigurations []NetworkConfigurationModel `tfschema:"network_interface_configurations"`
}

type ApiEntityReferenceModel struct {
	Id string `tfschema:"id"`
}

type NetworkConfigurationModel struct {
	Name       string                                                      `tfschema:"name"`
	Properties []VirtualMachineScaleSetNetworkConfigurationPropertiesModel `tfschema:"properties"`
}

type VirtualMachineScaleSetNetworkConfigurationPropertiesModel struct {
	AuxiliaryMode               string                  `tfschema:"auxiliary_mode"`
	AuxiliarySku                string                  `tfschema:"auxiliary_sku"`
	DeleteOption                string                  `tfschema:"delete_option"`
	DisableTcpStateTracking     bool                    `tfschema:"disable_tcp_state_tracking"`
	DnsSettings                 []DnsSettingsModel      `tfschema:"dns_settings"`
	EnableAcceleratedNetworking bool                    `tfschema:"enable_accelerated_networking"`
	EnableFpga                  bool                    `tfschema:"enable_fpga"`
	EnableIPForwarding          bool                    `tfschema:"enable_ip_forwarding"`
	IPConfigurations            []IPConfigurationsModel `tfschema:"ip_configurations"`
	NetworkSecurityGroup        []SubResourceModel      `tfschema:"network_security_group"`
	Primary                     bool                    `tfschema:"primary"`
}

type DnsSettingsModel struct {
	DnsServers []string `tfschema:"dns_servers"`
}

type IPConfigurationsModel struct {
	Name       string                                       `tfschema:"name"`
	Properties []VirtualMachineScaleSetIPConfigurationModel `tfschema:"properties"`
}

type VirtualMachineScaleSetIPConfigurationModel struct {
	ApplicationGatewayBackendAddressPools []SubResourceModel                                        `tfschema:"application_gateway_backend_address_pools"`
	ApplicationSecurityGroups             []SubResourceModel                                        `tfschema:"application_security_groups"`
	LoadBalancerBackendAddressPools       []SubResourceModel                                        `tfschema:"load_balancer_backend_address_pools"`
	LoadBalancerInboundNatPools           []SubResourceModel                                        `tfschema:"load_balancer_inbound_nat_pools"`
	Primary                               bool                                                      `tfschema:"primary"`
	PrivateIPAddressVersion               string                                                    `tfschema:"private_ip_address_version"`
	PublicIPAddressConfiguration          []VirtualMachineScaleSetPublicIPAddressConfigurationModel `tfschema:"public_ip_address_configuration"`
	Subnet                                []ApiEntityReferenceModel                                 `tfschema:"subnet"`
}

type VirtualMachineScaleSetPublicIPAddressConfigurationModel struct {
	Name       string                                                    `tfschema:"name"`
	Properties []VirtualMachineScaleSetPublicIPAddressConfigurationModel `tfschema:"properties"`
	Sku        []PublicIPAddressSkuModel                                 `tfschema:"sku"`
}

type VirtualMachineScaleSetPublicIPAddressConfigurationModel struct {
	DeleteOption           string                                                               `tfschema:"delete_option"`
	DnsSettings            []VirtualMachineScaleSetPublicIPAddressConfigurationDnsSettingsModel `tfschema:"dns_settings"`
	IPTags                 []VirtualMachineScaleSetIPTagModel                                   `tfschema:"ip_tags"`
	IdleTimeoutInMinutes   int64                                                                `tfschema:"idle_timeout_in_minutes"`
	PublicIPAddressVersion string                                                               `tfschema:"public_ip_address_version"`
	PublicIPPrefix         []SubResourceModel                                                   `tfschema:"public_ip_prefix"`
}

type VirtualMachineScaleSetPublicIPAddressConfigurationDnsSettingsModel struct {
	DomainNameLabel      string `tfschema:"domain_name_label"`
	DomainNameLabelScope string `tfschema:"domain_name_label_scope"`
}

type VirtualMachineScaleSetIPTagModel struct {
	IPTagType string `tfschema:"ip_tag_type"`
	Tag       string `tfschema:"tag"`
}

type PublicIPAddressSkuModel struct {
	Name string `tfschema:"name"`
	Tier string `tfschema:"tier"`
}

type OSProfileModel struct {
	AdminPassword               string                      `tfschema:"admin_password"`
	AdminUsername               string                      `tfschema:"admin_username"`
	ExtensionOperationsEnabled  bool                        `tfschema:"extension_operations_enabled"`
	ComputerNamePrefix          string                      `tfschema:"computer_name_prefix"`
	CustomData                  string                      `tfschema:"custom_data"`
	LinuxConfiguration          []LinuxConfigurationModel   `tfschema:"linux_configuration"`
	RequireGuestProvisionSignal bool                        `tfschema:"require_guest_provision_signal"`
	Secrets                     []VaultSecretGroupModel     `tfschema:"secrets"`
	WindowsConfiguration        []WindowsConfigurationModel `tfschema:"windows_configuration"`
}

type LinuxConfigurationModel struct {
	PasswordAuthenticationEnabled            bool                `tfschema:"password_authentication_enabled"`
	VMAgentPlatformUpdatesEnabled            bool                `tfschema:"vm_agent_platform_updates_enabled"`
	AssessmentMode                           string              `tfschema:"assessment_mode"`
	BypassPlatformSafetyChecksOnUserSchedule bool                `tfschema:"bypass_platform_safety_checks_on_user_schedule"`
	RebootSetting                            string              `tfschema:"reboot_setting"`
	PatchMode                                string              `tfschema:"patch_mode"`
	ProvisionVMAgent                         bool                `tfschema:"provision_vm_agent"`
	SshPublicKeys                            []SshPublicKeyModel `tfschema:"ssh_public_keys"`
}

type SshPublicKeyModel struct {
	KeyData string `tfschema:"key_data"`
	Path    string `tfschema:"path"`
}

type VaultSecretGroupModel struct {
	SourceVaultIds []string           `tfschema:"source_vault_ids"`
	Certificates   []CertificateModel `tfschema:"certificates"`
}

type CertificateModel struct {
	Store string `tfschema:"store"`
	Url   string `tfschema:"url"`
}

type WindowsConfigurationModel struct {
	AdditionalUnattendContent                []AdditionalUnattendContentModel `tfschema:"additional_unattend_content"`
	AutomaticUpdatesEnabled                  bool                             `tfschema:"automatic_updates_enabled"`
	VMAgentPlatformUpdatesEnabled            bool                             `tfschema:"vm_agent_platform_updates_enabled"`
	AssessmentMode                           string                           `tfschema:"assessment_mode"`
	BypassPlatformSafetyChecksOnUserSchedule bool                             `tfschema:"bypass_platform_safety_checks_on_user_schedule"`
	RebootSetting                            string                           `tfschema:"reboot_setting"`
	HotPatchingEnabled                       bool                             `tfschema:"hot_patching_enabled"`
	PatchMode                                string                           `tfschema:"patch_mode"`
	ProvisionVMAgent                         bool                             `tfschema:"provision_vm_agent"`
	TimeZone                                 string                           `tfschema:"time_zone"`
	WinRMListeners                           []WinRMListenerModel             `tfschema:"win_rm_listeners"`
}

type AdditionalUnattendContentModel struct {
	ComponentName string `tfschema:"component_name"`
	Content       string `tfschema:"content"`
	PassName      string `tfschema:"pass_name"`
	SettingName   string `tfschema:"setting_name"`
}

type WinRMListenerModel struct {
	CertificateUrl string `tfschema:"certificate_url"`
	Protocol       string `tfschema:"protocol"`
}

type OSImageNotificationProfileModel struct {
	Enable           bool   `tfschema:"enable"`
	NotBeforeTimeout string `tfschema:"not_before_timeout"`
}

type TerminationNotificationProfileModel struct {
	Enable           bool   `tfschema:"enable"`
	NotBeforeTimeout string `tfschema:"not_before_timeout"`
}

type SecurityPostureReferenceModel struct {
	ExcludeExtensions []string `tfschema:"exclude_extensions"`
	Id                string   `tfschema:"id"`
	IsOverridable     bool     `tfschema:"is_overridable"`
}

type SecurityProfileModel struct {
	EncryptionAtHost               bool                      `tfschema:"encryption_at_host"`
	UserAssignedIdentityResourceId string                    `tfschema:"user_assigned_identity_resource_id"`
	ProxyAgentSettings             []ProxyAgentSettingsModel `tfschema:"proxy_agent_settings"`
	SecurityType                   string                    `tfschema:"security_type"`
	SecureBootEnabled              bool                      `tfschema:"secure_boot_enabled"`
	VTpmEnabled                    bool                      `tfschema:"vtpm_enabled"`
}

type ProxyAgentSettingsModel struct {
	Enabled          bool   `tfschema:"enabled"`
	KeyIncarnationId int64  `tfschema:"key_incarnation_id"`
	Mode             string `tfschema:"mode"`
}

type ServiceArtifactReferenceModel struct {
	Id string `tfschema:"id"`
}

type StorageProfileModel struct {
	DataDisks          []VirtualMachineScaleSetDataDiskModel `tfschema:"data_disks"`
	DiskControllerType string                                `tfschema:"disk_controller_type"`
	ImageReference     []ImageReferenceModel                 `tfschema:"image_reference"`
	OsDisk             []VirtualMachineScaleSetOSDiskModel   `tfschema:"os_disk"`
}

type VirtualMachineScaleSetDataDiskModel struct {
	Caching                 string                                             `tfschema:"caching"`
	CreateOption            string                                             `tfschema:"create_option"`
	DeleteOption            string                                             `tfschema:"delete_option"`
	DiskIOPSReadWrite       int64                                              `tfschema:"disk_iops_read_write"`
	DiskMBpsReadWrite       int64                                              `tfschema:"disk_m_bps_read_write"`
	DiskSizeInGB            int64                                              `tfschema:"disk_size_in_gb"`
	Lun                     int64                                              `tfschema:"lun"`
	ManagedDisk             []VirtualMachineScaleSetManagedDiskParametersModel `tfschema:"managed_disk"`
	Name                    string                                             `tfschema:"name"`
	WriteAcceleratorEnabled bool                                               `tfschema:"write_accelerator_enabled"`
}

type VirtualMachineScaleSetManagedDiskParametersModel struct {
	DiskEncryptionSetId string                       `tfschema:"disk_encryption_set_id"`
	SecurityProfile     []VMDiskSecurityProfileModel `tfschema:"security_profile"`
	StorageAccountType  string                       `tfschema:"storage_account_type"`
}

type VMDiskSecurityProfileModel struct {
	DiskEncryptionSetId    string `tfschema:"disk_encryption_set_id"`
	SecurityEncryptionType string `tfschema:"security_encryption_type"`
}

type ImageReferenceModel struct {
	CommunityGalleryImageId string `tfschema:"community_gallery_image_id"`
	ExactVersion            string `tfschema:"exact_version"`
	Id                      string `tfschema:"id"`
	Offer                   string `tfschema:"offer"`
	Publisher               string `tfschema:"publisher"`
	SharedGalleryImageId    string `tfschema:"shared_gallery_image_id"`
	Sku                     string `tfschema:"sku"`
	Version                 string `tfschema:"version"`
}

type VirtualMachineScaleSetOSDiskModel struct {
	Caching                 string                                             `tfschema:"caching"`
	CreateOption            string                                             `tfschema:"create_option"`
	DeleteOption            string                                             `tfschema:"delete_option"`
	DiffDiskSettings        []DiffDiskSettingsModel                            `tfschema:"diff_disk_settings"`
	DiskSizeInGB            int64                                              `tfschema:"disk_size_in_gb"`
	ImageUri                string                                             `tfschema:"image_uri"`
	ManagedDisk             []VirtualMachineScaleSetManagedDiskParametersModel `tfschema:"managed_disk"`
	Name                    string                                             `tfschema:"name"`
	OsType                  string                                             `tfschema:"os_type"`
	VhdContainers           []string                                           `tfschema:"vhd_containers"`
	WriteAcceleratorEnabled bool                                               `tfschema:"write_accelerator_enabled"`
}

type DiffDiskSettingsModel struct {
	Option    string `tfschema:"option"`
	Placement string `tfschema:"placement"`
}

type ComputeProfileModel struct {
	AdditionalVirtualMachineCapabilities []AdditionalCapabilitiesModel `tfschema:"additional_virtual_machine_capabilities"`
	BaseVirtualMachineProfile            []VirtualMachineProfileModel  `tfschema:"base_virtual_machine_profile"`
	ComputeApiVersion                    string                        `tfschema:"compute_api_version"`
	PlatformFaultDomainCount             int64                         `tfschema:"platform_fault_domain_count"`
}

type AdditionalCapabilitiesModel struct {
	HibernationEnabled bool `tfschema:"hibernation_enabled"`
	UltraSSDEnabled    bool `tfschema:"ultra_ssd_enabled"`
}

type PlanModel struct {
	Name          string `tfschema:"name"`
	Product       string `tfschema:"product"`
	PromotionCode string `tfschema:"promotion_code"`
	Publisher     string `tfschema:"publisher"`
	Version       string `tfschema:"version"`
}

type RegularPriorityProfileModel struct {
	AllocationStrategy string `tfschema:"allocation_strategy"`
	Capacity           int64  `tfschema:"capacity"`
	MinCapacity        int64  `tfschema:"min_capacity"`
}

type SpotPriorityProfileModel struct {
	AllocationStrategy string  `tfschema:"allocation_strategy"`
	Capacity           int64   `tfschema:"capacity"`
	EvictionPolicy     string  `tfschema:"eviction_policy"`
	Maintain           bool    `tfschema:"maintain"`
	MaxPricePerVM      float64 `tfschema:"max_price_per_vm"`
	MinCapacity        int64   `tfschema:"min_capacity"`
}

type VMAttributesModel struct {
	AcceleratorCount          []VMAttributeMinMaxIntegerModel `tfschema:"accelerator_count"`
	AcceleratorManufacturers  []string                        `tfschema:"accelerator_manufacturers"`
	AcceleratorSupport        string                          `tfschema:"accelerator_support"`
	AcceleratorTypes          []string                        `tfschema:"accelerator_types"`
	ArchitectureTypes         []string                        `tfschema:"architecture_types"`
	BurstableSupport          string                          `tfschema:"burstable_support"`
	CpuManufacturers          []string                        `tfschema:"cpu_manufacturers"`
	DataDiskCount             []VMAttributeMinMaxIntegerModel `tfschema:"data_disk_count"`
	ExcludedVMSizes           []string                        `tfschema:"excluded_vm_sizes_profile"`
	LocalStorageDiskTypes     []string                        `tfschema:"local_storage_disk_types"`
	LocalStorageInGiB         []VMAttributeMinMaxDoubleModel  `tfschema:"local_storage_in_gi_b"`
	LocalStorageSupport       string                          `tfschema:"local_storage_support"`
	MemoryInGiB               []VMAttributeMinMaxDoubleModel  `tfschema:"memory_in_gi_b"`
	MemoryInGiBPerVCPU        []VMAttributeMinMaxDoubleModel  `tfschema:"memory_in_gi_b_per_vcpu"`
	NetworkBandwidthInMbps    []VMAttributeMinMaxDoubleModel  `tfschema:"network_bandwidth_in_mbps"`
	NetworkInterfaceCount     []VMAttributeMinMaxIntegerModel `tfschema:"network_interface_count"`
	RdmaNetworkInterfaceCount []VMAttributeMinMaxIntegerModel `tfschema:"rdma_network_interface_count"`
	RdmaSupport               string                          `tfschema:"rdma_support"`
	VCPUCount                 []VMAttributeMinMaxIntegerModel `tfschema:"vcpu_count"`
	VMCategories              []string                        `tfschema:"vm_categories"`
}

type VMAttributeMinMaxIntegerModel struct {
	Max int64 `tfschema:"max"`
	Min int64 `tfschema:"min"`
}

type VMAttributeMinMaxDoubleModel struct {
	Max float64 `tfschema:"max"`
	Min float64 `tfschema:"min"`
}

type VMSizeProfileModel struct {
	Name string `tfschema:"name"`
	Rank int64  `tfschema:"rank"`
}

type AzureFleetResource struct{}

var _ sdk.ResourceWithUpdate = AzureFleetResource{}

func (r AzureFleetResource) ResourceType() string {
	return "azurerm_azure_fleet"
}

func (r AzureFleetResource) ModelObject() interface{} {
	return &AzureFleetResourceModel{}
}

func (r AzureFleetResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return fleets.ValidateFleetID
}

func (r AzureFleetResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		// "location": commonschema.LocationWithoutForceNew(),
		// "additional_locations_profile": azurefleet.additionalLocationsProfileSchema(),

		"compute_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"additional_virtual_machine_capabilities": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"hibernation_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
								},

								"ultra_ssd_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
								},
							},
						},
					},

					"base_virtual_machine_profile": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"application_profile": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"gallery_applications": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"configuration_reference": {
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},

														"enable_automatic_upgrade": {
															Type:     pluginsdk.TypeBool,
															Optional: true,
														},

														"order": {
															Type:     pluginsdk.TypeInt,
															Optional: true,
														},

														"package_reference_id": {
															Type:         pluginsdk.TypeString,
															Required:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},

														"tags": {
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},

														"treat_failure_as_deployment_failure": {
															Type:     pluginsdk.TypeBool,
															Optional: true,
														},
													},
												},
											},
										},
									},
								},

								"capacity_reservation_group_id": {
									Elem: &schema.Schema{
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: capacityreservationgroups.ValidateCapacityReservationGroupID,
									},
								},

								"diagnostics_profile": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
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

														"storage_uri": {
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

								"extensions": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"id": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"name": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"properties": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"auto_upgrade_minor_version": {
															Type:     pluginsdk.TypeBool,
															Optional: true,
														},

														"enable_automatic_upgrade": {
															Type:     pluginsdk.TypeBool,
															Optional: true,
														},

														"force_update_tag": {
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},

														"protected_settings": {
															Type:             pluginsdk.TypeString,
															Optional:         true,
															ValidateFunc:     validation.StringIsJSON,
															DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
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
																		ValidateFunc: validation.StringIsNotEmpty,
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

														"provisioning_state": {
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},

														"publisher": {
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ValidateFunc: validation.StringIsNotEmpty,
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

														"type": {
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},

														"type_handler_version": {
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},
													},
												},
											},

											"type": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
										},
									},
								},

								"extensions_time_budget": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"hardware_profile": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"vm_size_properties": {
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
										},
									},
								},

								"license_type": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"network_profile": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"health_probe": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"id": {
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},
													},
												},
											},

											"network_api_version": {
												Type:     pluginsdk.TypeString,
												Optional: true,
												ValidateFunc: validation.StringInSlice([]string{
													string(fleets.NetworkApiVersionTwoZeroTwoZeroNegativeOneOneNegativeZeroOne),
												}, false),
											},

											"network_interface_configurations": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"name": {
															Type:         pluginsdk.TypeString,
															Required:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},

														"properties": {
															Type:     pluginsdk.TypeList,
															Optional: true,
															MaxItems: 1,
															Elem: &pluginsdk.Resource{
																Schema: map[string]*pluginsdk.Schema{
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
																			string(fleets.NetworkInterfaceAuxiliarySkuNone),
																			string(fleets.NetworkInterfaceAuxiliarySkuAOne),
																			string(fleets.NetworkInterfaceAuxiliarySkuATwo),
																			string(fleets.NetworkInterfaceAuxiliarySkuAFour),
																			string(fleets.NetworkInterfaceAuxiliarySkuAEight),
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

																	"disable_tcp_state_tracking": {
																		Type:     pluginsdk.TypeBool,
																		Optional: true,
																	},

																	"dns_settings": {
																		Type:     pluginsdk.TypeList,
																		Optional: true,
																		MaxItems: 1,
																		Elem: &pluginsdk.Resource{
																			Schema: map[string]*pluginsdk.Schema{
																				"dns_servers": {
																					Type:     pluginsdk.TypeList,
																					Optional: true,
																					Elem: &pluginsdk.Schema{
																						Type: pluginsdk.TypeString,
																					},
																				},
																			},
																		},
																	},

																	"enable_accelerated_networking": {
																		Type:     pluginsdk.TypeBool,
																		Optional: true,
																	},

																	"enable_fpga": {
																		Type:     pluginsdk.TypeBool,
																		Optional: true,
																	},

																	"enable_ip_forwarding": {
																		Type:     pluginsdk.TypeBool,
																		Optional: true,
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

																				"properties": {
																					Type:     pluginsdk.TypeList,
																					Optional: true,
																					MaxItems: 1,
																					Elem: &pluginsdk.Resource{
																						Schema: map[string]*pluginsdk.Schema{
																							"application_gateway_backend_address_pools": {
																								Type:     pluginsdk.TypeList,
																								Optional: true,
																								Elem: &pluginsdk.Resource{
																									Schema: map[string]*pluginsdk.Schema{
																										"id": {
																											Type:         pluginsdk.TypeString,
																											Optional:     true,
																											ValidateFunc: validation.StringIsNotEmpty,
																										},
																									},
																								},
																							},

																							"application_security_groups": {
																								Type:     pluginsdk.TypeList,
																								Optional: true,
																								Elem: &pluginsdk.Resource{
																									Schema: map[string]*pluginsdk.Schema{
																										"id": {
																											Type:         pluginsdk.TypeString,
																											Optional:     true,
																											ValidateFunc: validation.StringIsNotEmpty,
																										},
																									},
																								},
																							},

																							"load_balancer_backend_address_pools": {
																								Type:     pluginsdk.TypeList,
																								Optional: true,
																								Elem: &pluginsdk.Resource{
																									Schema: map[string]*pluginsdk.Schema{
																										"id": {
																											Type:         pluginsdk.TypeString,
																											Optional:     true,
																											ValidateFunc: validation.StringIsNotEmpty,
																										},
																									},
																								},
																							},

																							"load_balancer_inbound_nat_pools": {
																								Type:     pluginsdk.TypeList,
																								Optional: true,
																								Elem: &pluginsdk.Resource{
																									Schema: map[string]*pluginsdk.Schema{
																										"id": {
																											Type:         pluginsdk.TypeString,
																											Optional:     true,
																											ValidateFunc: validation.StringIsNotEmpty,
																										},
																									},
																								},
																							},

																							"primary": {
																								Type:     pluginsdk.TypeBool,
																								Optional: true,
																							},

																							"private_ip_address_version": {
																								Type:     pluginsdk.TypeString,
																								Optional: true,
																								ValidateFunc: validation.StringInSlice([]string{
																									string(fleets.IPVersionIPvFour),
																									string(fleets.IPVersionIPvSix),
																								}, false),
																							},

																							"public_ip_address_configuration": {
																								Type:     pluginsdk.TypeList,
																								Optional: true,
																								MaxItems: 1,
																								Elem: &pluginsdk.Resource{
																									Schema: map[string]*pluginsdk.Schema{
																										"name": {
																											Type:         pluginsdk.TypeString,
																											Required:     true,
																											ValidateFunc: validation.StringIsNotEmpty,
																										},

																										"properties": {
																											Type:     pluginsdk.TypeList,
																											Optional: true,
																											MaxItems: 1,
																											Elem: &pluginsdk.Resource{
																												Schema: map[string]*pluginsdk.Schema{
																													"delete_option": {
																														Type:     pluginsdk.TypeString,
																														Optional: true,
																														ValidateFunc: validation.StringInSlice([]string{
																															string(fleets.DeleteOptionsDelete),
																															string(fleets.DeleteOptionsDetach),
																														}, false),
																													},

																													"dns_settings": {
																														Type:     pluginsdk.TypeList,
																														Optional: true,
																														MaxItems: 1,
																														Elem: &pluginsdk.Resource{
																															Schema: map[string]*pluginsdk.Schema{
																																"domain_name_label": {
																																	Type:         pluginsdk.TypeString,
																																	Required:     true,
																																	ValidateFunc: validation.StringIsNotEmpty,
																																},

																																"domain_name_label_scope": {
																																	Type:     pluginsdk.TypeString,
																																	Optional: true,
																																	ValidateFunc: validation.StringInSlice([]string{
																																		string(fleets.DomainNameLabelScopeTypesTenantReuse),
																																		string(fleets.DomainNameLabelScopeTypesSubscriptionReuse),
																																		string(fleets.DomainNameLabelScopeTypesResourceGroupReuse),
																																		string(fleets.DomainNameLabelScopeTypesNoReuse),
																																	}, false),
																																},
																															},
																														},
																													},

																													"ip_tags": {
																														Type:     pluginsdk.TypeList,
																														Optional: true,
																														Elem: &pluginsdk.Resource{
																															Schema: map[string]*pluginsdk.Schema{
																																"ip_tag_type": {
																																	Type:         pluginsdk.TypeString,
																																	Optional:     true,
																																	ValidateFunc: validation.StringIsNotEmpty,
																																},

																																"tag": {
																																	Type:         pluginsdk.TypeString,
																																	Optional:     true,
																																	ValidateFunc: validation.StringIsNotEmpty,
																																},
																															},
																														},
																													},

																													"idle_timeout_in_minutes": {
																														Type:     pluginsdk.TypeInt,
																														Optional: true,
																													},

																													"public_ip_address_version": {
																														Type:     pluginsdk.TypeString,
																														Optional: true,
																														ValidateFunc: validation.StringInSlice([]string{
																															string(fleets.IPVersionIPvSix),
																															string(fleets.IPVersionIPvFour),
																														}, false),
																													},

																													"public_ip_prefix": {
																														Type:     pluginsdk.TypeList,
																														Optional: true,
																														MaxItems: 1,
																														Elem: &pluginsdk.Resource{
																															Schema: map[string]*pluginsdk.Schema{
																																"id": {
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

																							"subnet": {
																								Type:     pluginsdk.TypeList,
																								Optional: true,
																								MaxItems: 1,
																								Elem: &pluginsdk.Resource{
																									Schema: map[string]*pluginsdk.Schema{
																										"id": {
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
																			},
																		},
																	},

																	"network_security_group": {
																		Type:     pluginsdk.TypeList,
																		Optional: true,
																		MaxItems: 1,
																		Elem: &pluginsdk.Resource{
																			Schema: map[string]*pluginsdk.Schema{
																				"id": {
																					Type:         pluginsdk.TypeString,
																					Optional:     true,
																					ValidateFunc: validation.StringIsNotEmpty,
																				},
																			},
																		},
																	},

																	"primary": {
																		Type:     pluginsdk.TypeBool,
																		Optional: true,
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},

								"os_profile": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"admin_password": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"admin_username": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"allow_extension_operations": {
												Type:     pluginsdk.TypeBool,
												Optional: true,
											},

											"computer_name_prefix": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"custom_data": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"linux_configuration": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"disable_password_authentication": {
															Type:     pluginsdk.TypeBool,
															Optional: true,
														},

														"enable_vm_agent_platform_updates": {
															Type:     pluginsdk.TypeBool,
															Optional: true,
														},

														"patch_settings": {
															Type:     pluginsdk.TypeList,
															Optional: true,
															MaxItems: 1,
															Elem: &pluginsdk.Resource{
																Schema: map[string]*pluginsdk.Schema{
																	"assessment_mode": {
																		Type:     pluginsdk.TypeString,
																		Optional: true,
																		ValidateFunc: validation.StringInSlice([]string{
																			string(fleets.LinuxPatchAssessmentModeImageDefault),
																			string(fleets.LinuxPatchAssessmentModeAutomaticByPlatform),
																		}, false),
																	},

																	"automatic_by_platform_settings": {
																		Type:     pluginsdk.TypeList,
																		Optional: true,
																		MaxItems: 1,
																		Elem: &pluginsdk.Resource{
																			Schema: map[string]*pluginsdk.Schema{
																				"bypass_platform_safety_checks_on_user_schedule": {
																					Type:     pluginsdk.TypeBool,
																					Optional: true,
																				},

																				"reboot_setting": {
																					Type:     pluginsdk.TypeString,
																					Optional: true,
																					ValidateFunc: validation.StringInSlice([]string{
																						string(fleets.LinuxVMGuestPatchAutomaticByPlatformRebootSettingNever),
																						string(fleets.LinuxVMGuestPatchAutomaticByPlatformRebootSettingAlways),
																						string(fleets.LinuxVMGuestPatchAutomaticByPlatformRebootSettingUnknown),
																						string(fleets.LinuxVMGuestPatchAutomaticByPlatformRebootSettingIfRequired),
																					}, false),
																				},
																			},
																		},
																	},

																	"patch_mode": {
																		Type:     pluginsdk.TypeString,
																		Optional: true,
																		ValidateFunc: validation.StringInSlice([]string{
																			string(fleets.LinuxVMGuestPatchModeImageDefault),
																			string(fleets.LinuxVMGuestPatchModeAutomaticByPlatform),
																		}, false),
																	},
																},
															},
														},

														"provision_vm_agent": {
															Type:     pluginsdk.TypeBool,
															Optional: true,
														},

														"ssh": {
															Type:     pluginsdk.TypeList,
															Optional: true,
															MaxItems: 1,
															Elem: &pluginsdk.Resource{
																Schema: map[string]*pluginsdk.Schema{
																	"public_keys": {
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
														"source_vault_id": commonschema.ResourceIDReferenceOptional(&commonids.KeyVaultId{}),

														"vault_certificates": {
															Type:     pluginsdk.TypeList,
															Optional: true,
															Elem: &pluginsdk.Resource{
																Schema: map[string]*pluginsdk.Schema{
																	"certificate_store": {
																		Type:         pluginsdk.TypeString,
																		Optional:     true,
																		ValidateFunc: validation.StringIsNotEmpty,
																	},

																	"certificate_url": {
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

														"enable_automatic_updates": {
															Type:     pluginsdk.TypeBool,
															Optional: true,
														},

														"enable_vm_agent_platform_updates": {
															Type:     pluginsdk.TypeBool,
															Optional: true,
														},

														"patch_settings": {
															Type:     pluginsdk.TypeList,
															Optional: true,
															MaxItems: 1,
															Elem: &pluginsdk.Resource{
																Schema: map[string]*pluginsdk.Schema{
																	"assessment_mode": {
																		Type:     pluginsdk.TypeString,
																		Optional: true,
																		ValidateFunc: validation.StringInSlice([]string{
																			string(fleets.WindowsPatchAssessmentModeImageDefault),
																			string(fleets.WindowsPatchAssessmentModeAutomaticByPlatform),
																		}, false),
																	},

																	"automatic_by_platform_settings": {
																		Type:     pluginsdk.TypeList,
																		Optional: true,
																		MaxItems: 1,
																		Elem: &pluginsdk.Resource{
																			Schema: map[string]*pluginsdk.Schema{
																				"bypass_platform_safety_checks_on_user_schedule": {
																					Type:     pluginsdk.TypeBool,
																					Optional: true,
																				},

																				"reboot_setting": {
																					Type:     pluginsdk.TypeString,
																					Optional: true,
																					ValidateFunc: validation.StringInSlice([]string{
																						string(fleets.WindowsVMGuestPatchAutomaticByPlatformRebootSettingIfRequired),
																						string(fleets.WindowsVMGuestPatchAutomaticByPlatformRebootSettingNever),
																						string(fleets.WindowsVMGuestPatchAutomaticByPlatformRebootSettingAlways),
																						string(fleets.WindowsVMGuestPatchAutomaticByPlatformRebootSettingUnknown),
																					}, false),
																				},
																			},
																		},
																	},

																	"enable_hotpatching": {
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
																},
															},
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

														"win_rm": {
															Type:     pluginsdk.TypeList,
															Optional: true,
															MaxItems: 1,
															Elem: &pluginsdk.Resource{
																Schema: map[string]*pluginsdk.Schema{
																	"listeners": {
																		Type:     pluginsdk.TypeList,
																		Optional: true,
																		Elem: &pluginsdk.Resource{
																			Schema: map[string]*pluginsdk.Schema{
																				"certificate_url": {
																					Type:         pluginsdk.TypeString,
																					Optional:     true,
																					ValidateFunc: validation.StringIsNotEmpty,
																				},

																				"protocol": {
																					Type:     pluginsdk.TypeString,
																					Optional: true,
																					ValidateFunc: validation.StringInSlice([]string{
																						string(fleets.ProtocolTypesHTTP),
																						string(fleets.ProtocolTypesHTTPS),
																					}, false),
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},

								"scheduled_events_profile": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"os_image_notification_profile": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"enable": {
															Type:     pluginsdk.TypeBool,
															Optional: true,
														},

														"not_before_timeout": {
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},
													},
												},
											},

											"terminate_notification_profile": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"enable": {
															Type:     pluginsdk.TypeBool,
															Optional: true,
														},

														"not_before_timeout": {
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

								"security_profile": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"encryption_at_host": {
												Type:     pluginsdk.TypeBool,
												Optional: true,
											},

											"encryption_identity": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"user_assigned_identity_resource_id": {
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},
													},
												},
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

											"uefi_settings": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"secure_boot_enabled": {
															Type:     pluginsdk.TypeBool,
															Optional: true,
														},

														"v_tpm_enabled": {
															Type:     pluginsdk.TypeBool,
															Optional: true,
														},
													},
												},
											},
										},
									},
								},

								"service_artifact_reference": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"id": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
										},
									},
								},

								"storage_profile": {
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
																string(fleets.DiskCreateOptionTypesFromImage),
																string(fleets.DiskCreateOptionTypesEmpty),
																string(fleets.DiskCreateOptionTypesAttach),
																string(fleets.DiskCreateOptionTypesCopy),
																string(fleets.DiskCreateOptionTypesRestore),
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

														"disk_iops_read_write": {
															Type:     pluginsdk.TypeInt,
															Optional: true,
														},

														"disk_m_bps_read_write": {
															Type:     pluginsdk.TypeInt,
															Optional: true,
														},

														"disk_size_gb": {
															Type:     pluginsdk.TypeInt,
															Optional: true,
														},

														"lun": {
															Type:     pluginsdk.TypeInt,
															Required: true,
														},

														"managed_disk": {
															Type:     pluginsdk.TypeList,
															Optional: true,
															MaxItems: 1,
															Elem: &pluginsdk.Resource{
																Schema: map[string]*pluginsdk.Schema{
																	"disk_encryption_set": {
																		Type:     pluginsdk.TypeList,
																		Optional: true,
																		MaxItems: 1,
																		Elem: &pluginsdk.Resource{
																			Schema: map[string]*pluginsdk.Schema{
																				"id": {
																					Type:         pluginsdk.TypeString,
																					Optional:     true,
																					ValidateFunc: validation.StringIsNotEmpty,
																				},
																			},
																		},
																	},

																	"security_profile": {
																		Type:     pluginsdk.TypeList,
																		Optional: true,
																		MaxItems: 1,
																		Elem: &pluginsdk.Resource{
																			Schema: map[string]*pluginsdk.Schema{
																				"disk_encryption_set": {
																					Type:     pluginsdk.TypeList,
																					Optional: true,
																					MaxItems: 1,
																					Elem: &pluginsdk.Resource{
																						Schema: map[string]*pluginsdk.Schema{
																							"id": {
																								Type:         pluginsdk.TypeString,
																								Optional:     true,
																								ValidateFunc: validation.StringIsNotEmpty,
																							},
																						},
																					},
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
																			string(fleets.StorageAccountTypesStandardSSDLRS),
																			string(fleets.StorageAccountTypesUltraSSDLRS),
																			string(fleets.StorageAccountTypesPremiumZRS),
																			string(fleets.StorageAccountTypesStandardSSDZRS),
																			string(fleets.StorageAccountTypesPremiumVTwoLRS),
																			string(fleets.StorageAccountTypesStandardLRS),
																			string(fleets.StorageAccountTypesPremiumLRS),
																		}, false),
																	},
																},
															},
														},

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

														"offer": {
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},

														"publisher": {
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},

														"shared_gallery_image_id": {
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},

														"sku": {
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},

														"version": {
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
																string(fleets.DiskCreateOptionTypesEmpty),
																string(fleets.DiskCreateOptionTypesAttach),
																string(fleets.DiskCreateOptionTypesCopy),
																string(fleets.DiskCreateOptionTypesRestore),
																string(fleets.DiskCreateOptionTypesFromImage),
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
															MaxItems: 1,
															Elem: &pluginsdk.Resource{
																Schema: map[string]*pluginsdk.Schema{
																	"option": {
																		Type:     pluginsdk.TypeString,
																		Optional: true,
																		ValidateFunc: validation.StringInSlice([]string{
																			string(fleets.DiffDiskOptionsLocal),
																		}, false),
																	},

																	"placement": {
																		Type:     pluginsdk.TypeString,
																		Optional: true,
																		ValidateFunc: validation.StringInSlice([]string{
																			string(fleets.DiffDiskPlacementCacheDisk),
																			string(fleets.DiffDiskPlacementResourceDisk),
																			string(fleets.DiffDiskPlacementNVMeDisk),
																		}, false),
																	},
																},
															},
														},

														"disk_size_gb": {
															Type:     pluginsdk.TypeInt,
															Optional: true,
														},

														"image": {
															Type:     pluginsdk.TypeList,
															Optional: true,
															MaxItems: 1,
															Elem: &pluginsdk.Resource{
																Schema: map[string]*pluginsdk.Schema{
																	"uri": {
																		Type:         pluginsdk.TypeString,
																		Optional:     true,
																		ValidateFunc: validation.StringIsNotEmpty,
																	},
																},
															},
														},

														"managed_disk": {
															Type:     pluginsdk.TypeList,
															Optional: true,
															MaxItems: 1,
															Elem: &pluginsdk.Resource{
																Schema: map[string]*pluginsdk.Schema{
																	"disk_encryption_set": {
																		Type:     pluginsdk.TypeList,
																		Optional: true,
																		MaxItems: 1,
																		Elem: &pluginsdk.Resource{
																			Schema: map[string]*pluginsdk.Schema{
																				"id": {
																					Type:         pluginsdk.TypeString,
																					Optional:     true,
																					ValidateFunc: validation.StringIsNotEmpty,
																				},
																			},
																		},
																	},

																	"security_profile": {
																		Type:     pluginsdk.TypeList,
																		Optional: true,
																		MaxItems: 1,
																		Elem: &pluginsdk.Resource{
																			Schema: map[string]*pluginsdk.Schema{
																				"disk_encryption_set": {
																					Type:     pluginsdk.TypeList,
																					Optional: true,
																					MaxItems: 1,
																					Elem: &pluginsdk.Resource{
																						Schema: map[string]*pluginsdk.Schema{
																							"id": {
																								Type:         pluginsdk.TypeString,
																								Optional:     true,
																								ValidateFunc: validation.StringIsNotEmpty,
																							},
																						},
																					},
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
																			string(fleets.StorageAccountTypesPremiumLRS),
																			string(fleets.StorageAccountTypesStandardSSDLRS),
																			string(fleets.StorageAccountTypesUltraSSDLRS),
																			string(fleets.StorageAccountTypesPremiumZRS),
																			string(fleets.StorageAccountTypesStandardSSDZRS),
																			string(fleets.StorageAccountTypesPremiumVTwoLRS),
																			string(fleets.StorageAccountTypesStandardLRS),
																		}, false),
																	},
																},
															},
														},

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
														},
													},
												},
											},
										},
									},
								},

								"time_created": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"user_data": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"compute_api_version": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"platform_fault_domain_count": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
				},
			},
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"location": commonschema.Location(),

		"plan": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"product": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"promotion_code": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"publisher": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"version": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"regular_priority_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"allocation_strategy": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(fleets.RegularPriorityAllocationStrategyLowestPrice),
							string(fleets.RegularPriorityAllocationStrategyPrioritized),
						}, false),
					},

					"capacity": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"min_capacity": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
				},
			},
		},

		"spot_priority_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"allocation_strategy": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(fleets.SpotAllocationStrategyPriceCapacityOptimized),
							string(fleets.SpotAllocationStrategyLowestPrice),
							string(fleets.SpotAllocationStrategyCapacityOptimized),
						}, false),
					},

					"capacity": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"eviction_policy": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(fleets.EvictionPolicyDelete),
							string(fleets.EvictionPolicyDeallocate),
						}, false),
					},

					"maintain": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"max_price_per_vm": {
						Type:     pluginsdk.TypeFloat,
						Optional: true,
					},

					"min_capacity": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
				},
			},
		},

		"tags": commonschema.Tags(),

		"vm_attributes": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"accelerator_count": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"max": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},

								"min": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},
							},
						},
					},

					"accelerator_manufacturers": {
						Type:     pluginsdk.TypeList,
						Optional: true,
					},

					"accelerator_support": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(fleets.VMAttributeSupportIncluded),
							string(fleets.VMAttributeSupportRequired),
							string(fleets.VMAttributeSupportExcluded),
						}, false),
					},

					"accelerator_types": {
						Type:     pluginsdk.TypeList,
						Optional: true,
					},

					"architecture_types": {
						Type:     pluginsdk.TypeList,
						Optional: true,
					},

					"burstable_support": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(fleets.VMAttributeSupportRequired),
							string(fleets.VMAttributeSupportExcluded),
							string(fleets.VMAttributeSupportIncluded),
						}, false),
					},

					"cpu_manufacturers": {
						Type:     pluginsdk.TypeList,
						Optional: true,
					},

					"data_disk_count": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"max": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},

								"min": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},
							},
						},
					},

					"excluded_vm_sizes_profile": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"local_storage_disk_types": {
						Type:     pluginsdk.TypeList,
						Optional: true,
					},

					"local_storage_in_gi_b": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"max": {
									Type:     pluginsdk.TypeFloat,
									Optional: true,
								},

								"min": {
									Type:     pluginsdk.TypeFloat,
									Optional: true,
								},
							},
						},
					},

					"local_storage_support": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(fleets.VMAttributeSupportRequired),
							string(fleets.VMAttributeSupportExcluded),
							string(fleets.VMAttributeSupportIncluded),
						}, false),
					},

					"memory_in_gi_b": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"max": {
									Type:     pluginsdk.TypeFloat,
									Optional: true,
								},

								"min": {
									Type:     pluginsdk.TypeFloat,
									Optional: true,
								},
							},
						},
					},

					"memory_in_gi_b_per_vcpu": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"max": {
									Type:     pluginsdk.TypeFloat,
									Optional: true,
								},

								"min": {
									Type:     pluginsdk.TypeFloat,
									Optional: true,
								},
							},
						},
					},

					"network_bandwidth_in_mbps": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"max": {
									Type:     pluginsdk.TypeFloat,
									Optional: true,
								},

								"min": {
									Type:     pluginsdk.TypeFloat,
									Optional: true,
								},
							},
						},
					},

					"network_interface_count": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"max": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},

								"min": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},
							},
						},
					},

					"rdma_network_interface_count": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"max": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},

								"min": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},
							},
						},
					},

					"rdma_support": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(fleets.VMAttributeSupportExcluded),
							string(fleets.VMAttributeSupportIncluded),
							string(fleets.VMAttributeSupportRequired),
						}, false),
					},

					"vcpu_count": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"max": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},

								"min": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},
							},
						},
					},

					"vm_categories": {
						Type:     pluginsdk.TypeList,
						Optional: true,
					},
				},
			},
		},

		"vm_sizes_profile_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"rank": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
				},
			},
		},

		"zones": commonschema.ZonesMultipleOptional(),
	}
}

func (r AzureFleetResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"time_created": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"unique_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r AzureFleetResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model AzureFleetResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.AzureFleet.FleetsClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := fleets.NewFleetID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			identityValue, err := identity.ExpandLegacySystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}
			properties := &fleets.Fleet{
				Identity: identityValue,
				Location: location.Normalize(model.Location),
				Plan:     expandPlanModel(model.Plan),
				Properties: &fleets.FleetProperties{
					RegularPriorityProfile: expandRegularPriorityProfileModel(model.RegularPriorityProfile),
					SpotPriorityProfile:    expandSpotPriorityProfileModel(model.SpotPriorityProfile),
					VMAttributes:           expandVMAttributesModel(model.VMAttributes),
				},
				Tags:  &model.Tags,
				Zones: &model.Zones,
			}

			additionalLocationsProfileValue, err := expandAdditionalLocationsProfileModel(model.AdditionalLocationsProfile)
			if err != nil {
				return err
			}

			properties.Properties.AdditionalLocationsProfile = additionalLocationsProfileValue

			computeProfileValue, err := expandComputeProfileModel(model.ComputeProfile)
			if err != nil {
				return err
			}

			properties.Properties.ComputeProfile = pointer.From(computeProfileValue)

			properties.Properties.VMSizesProfile = pointer.From(expandVMSizeProfileModelArray(model.VMSizesProfile))

			if err := client.CreateOrUpdateThenPoll(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r AzureFleetResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureFleet.FleetsClient

			id, err := fleets.ParseFleetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model AzureFleetResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}
			if resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			properties := resp.Model

			if metadata.ResourceData.HasChange("identity") {
				identityValue, err := identity.ExpandLegacySystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				properties.Identity = identityValue
			}

			if metadata.ResourceData.HasChange("plan") {
				properties.Plan = expandPlanModel(model.Plan)
			}

			if metadata.ResourceData.HasChange("additional_locations_profile") {
				additionalLocationsProfileValue, err := expandAdditionalLocationsProfileModel(model.AdditionalLocationsProfile)
				if err != nil {
					return err
				}

				properties.Properties.AdditionalLocationsProfile = additionalLocationsProfileValue
			}

			if metadata.ResourceData.HasChange("compute_profile") {
				computeProfileValue, err := expandComputeProfileModel(model.ComputeProfile)
				if err != nil {
					return err
				}

				properties.Properties.ComputeProfile = pointer.From(computeProfileValue)
			}

			if metadata.ResourceData.HasChange("regular_priority_profile") {
				properties.Properties.RegularPriorityProfile = expandRegularPriorityProfileModel(model.RegularPriorityProfile)
			}

			if metadata.ResourceData.HasChange("spot_priority_profile") {
				properties.Properties.SpotPriorityProfile = expandSpotPriorityProfileModel(model.SpotPriorityProfile)
			}

			if metadata.ResourceData.HasChange("vm_attributes") {
				properties.Properties.VMAttributes = expandVMAttributesModel(model.VMAttributes)
			}

			if metadata.ResourceData.HasChange("vm_sizes_profile") {
				properties.Properties.VMSizesProfile = pointer.From(expandVMSizeProfileModelArray(model.VMSizesProfile))
			}

			properties.SystemData = nil

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = &model.Tags
			}

			if metadata.ResourceData.HasChange("zones") {
				properties.Zones = &model.Zones
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r AzureFleetResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureFleet.FleetsClient

			id, err := fleets.ParseFleetID(metadata.ResourceData.Id())
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

			state := AzureFleetResourceModel{
				Name:              id.FleetName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)

				identityValue, err := identity.FlattenLegacySystemAndUserAssignedMap(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}

				if err := metadata.ResourceData.Set("identity", identityValue); err != nil {
					return fmt.Errorf("setting `identity`: %+v", err)
				}

				state.Plan = flattenPlanModel(model.Plan)
				if properties := model.Properties; properties != nil {
					additionalLocationsProfileValue, err := flattenAdditionalLocationsProfileModel(properties.AdditionalLocationsProfile)
					if err != nil {
						return err
					}
					state.AdditionalLocationsProfile = additionalLocationsProfileValue

					computeProfileValue, err := flattenComputeProfileModel(&properties.ComputeProfile)
					if err != nil {
						return err
					}
					state.ComputeProfile = computeProfileValue

					state.RegularPriorityProfile = flattenRegularPriorityProfileModel(properties.RegularPriorityProfile)

					state.SpotPriorityProfile = flattenSpotPriorityProfileModel(properties.SpotPriorityProfile)

					if properties.TimeCreated != nil {
						state.TimeCreated = *properties.TimeCreated
					}

					if properties.UniqueId != nil {
						state.UniqueId = *properties.UniqueId
					}

					state.VMAttributes = flattenVMAttributesModel(properties.VMAttributes)

					state.VMSizesProfile = flattenVMSizeProfileModelArray(&properties.VMSizesProfile)
				}
				if model.Tags != nil {
					state.Tags = *model.Tags
				}
				if model.Zones != nil {
					state.Zones = *model.Zones
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r AzureFleetResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureFleet.FleetsClient

			id, err := fleets.ParseFleetID(metadata.ResourceData.Id())
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

func expandAdditionalLocationsProfileModel(inputList []AdditionalLocationsProfileModel) (*fleets.AdditionalLocationsProfile, error) {
	if len(inputList) == 0 {
		return nil, nil
	}
	input := &inputList[0]
	output := fleets.AdditionalLocationsProfile{}

	locationProfilesValue, err := expandLocationProfileModelArray(input.LocationProfiles)
	if err != nil {
		return nil, err
	}

	output.LocationProfiles = pointer.From(locationProfilesValue)

	return &output, nil
}

func expandLocationProfileModelArray(inputList []LocationProfileModel) (*[]fleets.LocationProfile, error) {
	var outputList []fleets.LocationProfile
	for _, v := range inputList {
		input := v
		output := fleets.LocationProfile{
			Location: input.Location,
		}

		virtualMachineProfileOverrideValue, err := expandBaseVirtualMachineProfileModel(input.VirtualMachineProfile)
		if err != nil {
			return nil, err
		}

		output.VirtualMachineProfile = virtualMachineProfileOverrideValue

		outputList = append(outputList, output)
	}
	return &outputList, nil
}

func expandBaseVirtualMachineProfileModel(inputList []VirtualMachineProfileModel) (*fleets.BaseVirtualMachineProfile, error) {
	if len(inputList) == 0 {
		return nil, nil
	}
	input := &inputList[0]
	output := fleets.BaseVirtualMachineProfile{
		ApplicationProfile:       expandApplicationProfileModel(input.ApplicationProfile),
		CapacityReservation:      expandCapacityReservationProfileModel(input.CapacityReservation),
		DiagnosticsProfile:       expandDiagnosticsProfileModel(input.DiagnosticsProfile),
		HardwareProfile:          expandVirtualMachineScaleSetHardwareProfileModel(input.HardwareProfile),
		NetworkProfile:           expandVirtualMachineScaleSetNetworkProfileModel(input.NetworkProfile),
		OsProfile:                expandVirtualMachineScaleSetOSProfileModel(input.OsProfile),
		ScheduledEventsProfile:   expandScheduledEventsProfileModel(input.ScheduledEventsProfile),
		SecurityPostureReference: expandSecurityPostureReferenceModel(input.SecurityPostureReference),
		SecurityProfile:          expandSecurityProfileModel(input.SecurityProfile),
		ServiceArtifactReference: expandServiceArtifactReferenceModel(input.ServiceArtifactReference),
		StorageProfile:           expandVirtualMachineScaleSetStorageProfileModel(input.StorageProfile),
	}

	extensionProfileValue, err := expandVirtualMachineScaleSetExtensionProfileModel(input.ExtensionProfile)
	if err != nil {
		return nil, err
	}

	output.ExtensionProfile = extensionProfileValue

	if input.LicenseType != "" {
		output.LicenseType = &input.LicenseType
	}

	if input.UserData != "" {
		output.UserData = &input.UserData
	}

	return &output, nil
}

func expandApplicationProfileModel(inputList []ApplicationProfileModel) *fleets.ApplicationProfile {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.ApplicationProfile{
		GalleryApplications: expandVMGalleryApplicationModelArray(input.GalleryApplications),
	}

	return &output
}

func expandVMGalleryApplicationModelArray(inputList []VMGalleryApplicationModel) *[]fleets.VMGalleryApplication {
	var outputList []fleets.VMGalleryApplication
	for _, v := range inputList {
		input := v
		output := fleets.VMGalleryApplication{
			EnableAutomaticUpgrade:          &input.EnableAutomaticUpgrade,
			Order:                           &input.Order,
			PackageReferenceId:              input.PackageReferenceId,
			TreatFailureAsDeploymentFailure: &input.TreatFailureAsDeploymentFailure,
		}

		if input.ConfigurationReference != "" {
			output.ConfigurationReference = &input.ConfigurationReference
		}

		if input.Tags != "" {
			output.Tags = &input.Tags
		}
		outputList = append(outputList, output)
	}
	return &outputList
}

func expandCapacityReservationProfileModel(inputList []CapacityReservationProfileModel) *fleets.CapacityReservationProfile {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.CapacityReservationProfile{
		CapacityReservationGroup: expandSubResourceModel(input.CapacityReservationGroup),
	}

	return &output
}

func expandSubResourceModel(inputList []SubResourceModel) *fleets.SubResource {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.SubResource{}
	if input.Id != "" {
		output.Id = &input.Id
	}

	return &output
}

func expandDiagnosticsProfileModel(inputList []DiagnosticsProfileModel) *fleets.DiagnosticsProfile {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.DiagnosticsProfile{
		BootDiagnostics: expandBootDiagnosticsModel(input.BootDiagnostics),
	}

	return &output
}

func expandBootDiagnosticsModel(inputList []BootDiagnosticsModel) *fleets.BootDiagnostics {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.BootDiagnostics{
		Enabled: &input.Enabled,
	}
	if input.StorageUri != "" {
		output.StorageUri = &input.StorageUri
	}

	return &output
}

func expandVirtualMachineScaleSetExtensionProfileModel(inputList []VirtualMachineScaleSetExtensionProfileModel) (*fleets.VirtualMachineScaleSetExtensionProfile, error) {
	if len(inputList) == 0 {
		return nil, nil
	}
	input := &inputList[0]
	output := fleets.VirtualMachineScaleSetExtensionProfile{}

	extensionsValue, err := expandVirtualMachineScaleSetExtensionModelArray(input.Extensions)
	if err != nil {
		return nil, err
	}

	output.Extensions = extensionsValue

	if input.ExtensionsTimeBudget != "" {
		output.ExtensionsTimeBudget = &input.ExtensionsTimeBudget
	}

	return &output, nil
}

func expandVirtualMachineScaleSetExtensionModelArray(inputList []ExtensionsModel) (*[]fleets.VirtualMachineScaleSetExtension, error) {
	var outputList []fleets.VirtualMachineScaleSetExtension
	for _, v := range inputList {
		input := v
		output := fleets.VirtualMachineScaleSetExtension{}

		if input.Name != "" {
			output.Name = &input.Name
		}

		propertiesValue, err := expandVirtualMachineScaleSetExtensionModel(input.Properties)
		if err != nil {
			return nil, err
		}

		output.Properties = propertiesValue

		outputList = append(outputList, output)
	}
	return &outputList, nil
}

func expandVirtualMachineScaleSetExtensionModel(inputList []VirtualMachineScaleSetExtensionModel) (*fleets.VirtualMachineScaleSetExtensionProperties, error) {
	if len(inputList) == 0 {
		return nil, nil
	}
	input := &inputList[0]
	output := fleets.VirtualMachineScaleSetExtensionProperties{
		AutoUpgradeMinorVersion:       &input.AutoUpgradeMinorVersion,
		EnableAutomaticUpgrade:        &input.EnableAutomaticUpgrade,
		ProtectedSettingsFromKeyVault: expandKeyVaultSecretReferenceModel(input.ProtectedSettingsFromKeyVault),
		ProvisionAfterExtensions:      &input.ProvisionAfterExtensions,
		SuppressFailures:              &input.SuppressFailures,
	}
	if input.ForceUpdateTag != "" {
		output.ForceUpdateTag = &input.ForceUpdateTag
	}

	var protectedSettingsValue interface{}
	err := json.Unmarshal([]byte(input.ProtectedSettings), &protectedSettingsValue)
	if err != nil {
		return nil, err
	}

	output.ProtectedSettings = &protectedSettingsValue

	if input.Publisher != "" {
		output.Publisher = &input.Publisher
	}

	var settingsValue interface{}
	err = json.Unmarshal([]byte(input.Settings), &settingsValue)
	if err != nil {
		return nil, err
	}

	output.Settings = &settingsValue

	if input.Type != "" {
		output.Type = &input.Type
	}

	if input.TypeHandlerVersion != "" {
		output.TypeHandlerVersion = &input.TypeHandlerVersion
	}

	return &output, nil
}

func expandKeyVaultSecretReferenceModel(inputList []KeyVaultSecretReferenceModel) *fleets.KeyVaultSecretReference {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.KeyVaultSecretReference{
		SecretURL: input.SecretUrl,
	}

	output.SourceVault = pointer.From(expandSubResourceModel(input.SourceVault))

	return &output
}

func expandVirtualMachineScaleSetHardwareProfileModel(inputList []VirtualMachineScaleSetHardwareProfileModel) *fleets.VirtualMachineScaleSetHardwareProfile {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.VirtualMachineScaleSetHardwareProfile{
		VMSizeProperties: expandVMSizePropertiesModel(input.VMSizeProperties),
	}

	return &output
}

func expandVMSizePropertiesModel(inputList []VMSizePropertiesModel) *fleets.VMSizeProperties {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.VMSizeProperties{
		VCPUsAvailable: &input.VCPUsAvailable,
		VCPUsPerCore:   &input.VCPUsPerCore,
	}

	return &output
}

func expandVirtualMachineScaleSetNetworkProfileModel(inputList []NetworkProfileModel) *fleets.VirtualMachineScaleSetNetworkProfile {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.VirtualMachineScaleSetNetworkProfile{
		HealthProbe:                    expandApiEntityReferenceModel(input.HealthProbe),
		NetworkApiVersion:              &input.NetworkApiVersion,
		NetworkInterfaceConfigurations: expandVirtualMachineScaleSetNetworkConfigurationModelArray(input.NetworkInterfaceConfigurations),
	}

	return &output
}

func expandApiEntityReferenceModel(inputList []ApiEntityReferenceModel) *fleets.ApiEntityReference {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.ApiEntityReference{}
	if input.Id != "" {
		output.Id = &input.Id
	}

	return &output
}

func expandVirtualMachineScaleSetNetworkConfigurationModelArray(inputList []NetworkConfigurationModel) *[]fleets.VirtualMachineScaleSetNetworkConfiguration {
	var outputList []fleets.VirtualMachineScaleSetNetworkConfiguration
	for _, v := range inputList {
		input := v
		output := fleets.VirtualMachineScaleSetNetworkConfiguration{
			Name:       input.Name,
			Properties: expandVirtualMachineScaleSetNetworkConfigurationPropertiesModel(input.Properties),
		}

		outputList = append(outputList, output)
	}
	return &outputList
}

func expandVirtualMachineScaleSetNetworkConfigurationPropertiesModel(inputList []VirtualMachineScaleSetNetworkConfigurationPropertiesModel) *fleets.VirtualMachineScaleSetNetworkConfigurationProperties {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.VirtualMachineScaleSetNetworkConfigurationProperties{
		AuxiliaryMode:               &input.AuxiliaryMode,
		AuxiliarySku:                &input.AuxiliarySku,
		DeleteOption:                &input.DeleteOption,
		DisableTcpStateTracking:     &input.DisableTcpStateTracking,
		DnsSettings:                 expandVirtualMachineScaleSetNetworkConfigurationDnsSettingsModel(input.DnsSettings),
		EnableAcceleratedNetworking: &input.EnableAcceleratedNetworking,
		EnableFpga:                  &input.EnableFpga,
		EnableIPForwarding:          &input.EnableIPForwarding,
		NetworkSecurityGroup:        expandSubResourceModel(input.NetworkSecurityGroup),
		Primary:                     &input.Primary,
	}

	output.IPConfigurations = pointer.From(expandVirtualMachineScaleSetIPConfigurationModelArray(input.IPConfigurations))

	return &output
}

func expandVirtualMachineScaleSetNetworkConfigurationDnsSettingsModel(inputList []VirtualMachineScaleSetNetworkConfigurationDnsSettingsModel) *fleets.VirtualMachineScaleSetNetworkConfigurationDnsSettings {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.VirtualMachineScaleSetNetworkConfigurationDnsSettings{
		DnsServers: &input.DnsServers,
	}

	return &output
}

func expandVirtualMachineScaleSetIPConfigurationModelArray(inputList []VirtualMachineScaleSetIPConfigurationModel) *[]fleets.VirtualMachineScaleSetIPConfiguration {
	var outputList []fleets.VirtualMachineScaleSetIPConfiguration
	for _, v := range inputList {
		input := v
		output := fleets.VirtualMachineScaleSetIPConfiguration{
			Name:       input.Name,
			Properties: expandVirtualMachineScaleSetIPConfigurationModel(input.Properties),
		}

		outputList = append(outputList, output)
	}
	return &outputList
}

func expandVirtualMachineScaleSetIPConfigurationModel(inputList []VirtualMachineScaleSetIPConfigurationModel) *fleets.VirtualMachineScaleSetIPConfigurationProperties {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.VirtualMachineScaleSetIPConfigurationProperties{
		ApplicationGatewayBackendAddressPools: expandSubResourceModelArray(input.ApplicationGatewayBackendAddressPools),
		ApplicationSecurityGroups:             expandSubResourceModelArray(input.ApplicationSecurityGroups),
		LoadBalancerBackendAddressPools:       expandSubResourceModelArray(input.LoadBalancerBackendAddressPools),
		LoadBalancerInboundNatPools:           expandSubResourceModelArray(input.LoadBalancerInboundNatPools),
		Primary:                               &input.Primary,
		PrivateIPAddressVersion:               &input.PrivateIPAddressVersion,
		PublicIPAddressConfiguration:          expandVirtualMachineScaleSetPublicIPAddressConfigurationModel(input.PublicIPAddressConfiguration),
		Subnet:                                expandApiEntityReferenceModel(input.Subnet),
	}

	return &output
}

func expandVirtualMachineScaleSetPublicIPAddressConfigurationModel(inputList []VirtualMachineScaleSetPublicIPAddressConfigurationModel) *fleets.VirtualMachineScaleSetPublicIPAddressConfiguration {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.VirtualMachineScaleSetPublicIPAddressConfiguration{
		Name:       input.Name,
		Properties: expandVirtualMachineScaleSetPublicIPAddressConfigurationModel(input.Properties),
		Sku:        expandPublicIPAddressSkuModel(input.Sku),
	}

	return &output
}

func expandVirtualMachineScaleSetPublicIPAddressConfigurationModel(inputList []VirtualMachineScaleSetPublicIPAddressConfigurationModel) *fleets.VirtualMachineScaleSetPublicIPAddressConfigurationProperties {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.VirtualMachineScaleSetPublicIPAddressConfigurationProperties{
		DeleteOption:           &input.DeleteOption,
		DnsSettings:            expandVirtualMachineScaleSetPublicIPAddressConfigurationDnsSettingsModel(input.DnsSettings),
		IPTags:                 expandVirtualMachineScaleSetIPTagModelArray(input.IPTags),
		IdleTimeoutInMinutes:   &input.IdleTimeoutInMinutes,
		PublicIPAddressVersion: &input.PublicIPAddressVersion,
		PublicIPPrefix:         expandSubResourceModel(input.PublicIPPrefix),
	}

	return &output
}

func expandVirtualMachineScaleSetPublicIPAddressConfigurationDnsSettingsModel(inputList []VirtualMachineScaleSetPublicIPAddressConfigurationDnsSettingsModel) *fleets.VirtualMachineScaleSetPublicIPAddressConfigurationDnsSettings {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.VirtualMachineScaleSetPublicIPAddressConfigurationDnsSettings{
		DomainNameLabel:      input.DomainNameLabel,
		DomainNameLabelScope: &input.DomainNameLabelScope,
	}

	return &output
}

func expandVirtualMachineScaleSetIPTagModelArray(inputList []VirtualMachineScaleSetIPTagModel) *[]fleets.VirtualMachineScaleSetIPTag {
	var outputList []fleets.VirtualMachineScaleSetIPTag
	for _, v := range inputList {
		input := v
		output := fleets.VirtualMachineScaleSetIPTag{}

		if input.IPTagType != "" {
			output.IPTagType = &input.IPTagType
		}

		if input.Tag != "" {
			output.Tag = &input.Tag
		}
		outputList = append(outputList, output)
	}
	return &outputList
}

func expandPublicIPAddressSkuModel(inputList []PublicIPAddressSkuModel) *fleets.PublicIPAddressSku {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.PublicIPAddressSku{
		Name: &input.Name,
		Tier: &input.Tier,
	}

	return &output
}

func expandVirtualMachineScaleSetOSProfileModel(inputList []OSProfileModel) *fleets.VirtualMachineScaleSetOSProfile {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.VirtualMachineScaleSetOSProfile{
		AllowExtensionOperations:    &input.AllowExtensionOperations,
		LinuxConfiguration:          expandLinuxConfigurationModel(input.LinuxConfiguration),
		RequireGuestProvisionSignal: &input.RequireGuestProvisionSignal,
		Secrets:                     expandVaultSecretGroupModelArray(input.Secrets),
		WindowsConfiguration:        expandWindowsConfigurationModel(input.WindowsConfiguration),
	}
	if input.AdminPassword != "" {
		output.AdminPassword = &input.AdminPassword
	}

	if input.AdminUsername != "" {
		output.AdminUsername = &input.AdminUsername
	}

	if input.ComputerNamePrefix != "" {
		output.ComputerNamePrefix = &input.ComputerNamePrefix
	}

	if input.CustomData != "" {
		output.CustomData = &input.CustomData
	}

	return &output
}

func expandLinuxConfigurationModel(inputList []LinuxConfigurationModel) *fleets.LinuxConfiguration {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.LinuxConfiguration{
		DisablePasswordAuthentication: &input.DisablePasswordAuthentication,
		EnableVMAgentPlatformUpdates:  &input.EnableVMAgentPlatformUpdates,
		PatchSettings:                 expandLinuxPatchSettingsModel(input.PatchSettings),
		ProvisionVMAgent:              &input.ProvisionVMAgent,
		Ssh:                           expandSshConfigurationModel(input.Ssh),
	}

	return &output
}

func expandLinuxPatchSettingsModel(inputList []LinuxPatchSettingsModel) *fleets.LinuxPatchSettings {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.LinuxPatchSettings{
		AssessmentMode:              &input.AssessmentMode,
		AutomaticByPlatformSettings: expandLinuxVMGuestPatchAutomaticByPlatformSettingsModel(input.AutomaticByPlatformSettings),
		PatchMode:                   &input.PatchMode,
	}

	return &output
}

func expandLinuxVMGuestPatchAutomaticByPlatformSettingsModel(inputList []LinuxVMGuestPatchAutomaticByPlatformSettingsModel) *fleets.LinuxVMGuestPatchAutomaticByPlatformSettings {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.LinuxVMGuestPatchAutomaticByPlatformSettings{
		BypassPlatformSafetyChecksOnUserSchedule: &input.BypassPlatformSafetyChecksOnUserSchedule,
		RebootSetting:                            &input.RebootSetting,
	}

	return &output
}

func expandSshConfigurationModel(inputList []SshConfigurationModel) *fleets.SshConfiguration {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.SshConfiguration{
		PublicKeys: expandSshPublicKeyModelArray(input.PublicKeys),
	}

	return &output
}

func expandSshPublicKeyModelArray(inputList []SshPublicKeyModel) *[]fleets.SshPublicKey {
	var outputList []fleets.SshPublicKey
	for _, v := range inputList {
		input := v
		output := fleets.SshPublicKey{}

		if input.KeyData != "" {
			output.KeyData = &input.KeyData
		}

		if input.Path != "" {
			output.Path = &input.Path
		}
		outputList = append(outputList, output)
	}
	return &outputList
}

func expandVaultSecretGroupModelArray(inputList []VaultSecretGroupModel) *[]fleets.VaultSecretGroup {
	var outputList []fleets.VaultSecretGroup
	for _, v := range inputList {
		input := v
		output := fleets.VaultSecretGroup{
			SourceVault:       expandSubResourceModel(input.SourceVault),
			VaultCertificates: expandVaultCertificateModelArray(input.VaultCertificates),
		}

		outputList = append(outputList, output)
	}
	return &outputList
}

func expandVaultCertificateModelArray(inputList []VaultCertificateModel) *[]fleets.VaultCertificate {
	var outputList []fleets.VaultCertificate
	for _, v := range inputList {
		input := v
		output := fleets.VaultCertificate{}

		if input.CertificateStore != "" {
			output.CertificateStore = &input.CertificateStore
		}

		if input.CertificateUrl != "" {
			output.CertificateUrl = &input.CertificateUrl
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
		AdditionalUnattendContent:    expandAdditionalUnattendContentModelArray(input.AdditionalUnattendContent),
		EnableAutomaticUpdates:       &input.EnableAutomaticUpdates,
		EnableVMAgentPlatformUpdates: &input.EnableVMAgentPlatformUpdates,
		PatchSettings:                expandPatchSettingsModel(input.PatchSettings),
		ProvisionVMAgent:             &input.ProvisionVMAgent,
		WinRM:                        expandWinRMConfigurationModel(input.WinRM),
	}
	if input.TimeZone != "" {
		output.TimeZone = &input.TimeZone
	}

	return &output
}

func expandAdditionalUnattendContentModelArray(inputList []AdditionalUnattendContentModel) *[]fleets.AdditionalUnattendContent {
	var outputList []fleets.AdditionalUnattendContent
	for _, v := range inputList {
		input := v
		output := fleets.AdditionalUnattendContent{
			ComponentName: &input.ComponentName,
			PassName:      &input.PassName,
			SettingName:   &input.SettingName,
		}

		if input.Content != "" {
			output.Content = &input.Content
		}
		outputList = append(outputList, output)
	}
	return &outputList
}

func expandPatchSettingsModel(inputList []PatchSettingsModel) *fleets.PatchSettings {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.PatchSettings{
		AssessmentMode:              &input.AssessmentMode,
		AutomaticByPlatformSettings: expandWindowsVMGuestPatchAutomaticByPlatformSettingsModel(input.AutomaticByPlatformSettings),
		EnableHotpatching:           &input.EnableHotpatching,
		PatchMode:                   &input.PatchMode,
	}

	return &output
}

func expandWindowsVMGuestPatchAutomaticByPlatformSettingsModel(inputList []WindowsVMGuestPatchAutomaticByPlatformSettingsModel) *fleets.WindowsVMGuestPatchAutomaticByPlatformSettings {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.WindowsVMGuestPatchAutomaticByPlatformSettings{
		BypassPlatformSafetyChecksOnUserSchedule: &input.BypassPlatformSafetyChecksOnUserSchedule,
		RebootSetting:                            &input.RebootSetting,
	}

	return &output
}

func expandWinRMConfigurationModel(inputList []WinRMConfigurationModel) *fleets.WinRMConfiguration {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.WinRMConfiguration{
		Listeners: expandWinRMListenerModelArray(input.Listeners),
	}

	return &output
}

func expandWinRMListenerModelArray(inputList []WinRMListenerModel) *[]fleets.WinRMListener {
	var outputList []fleets.WinRMListener
	for _, v := range inputList {
		input := v
		output := fleets.WinRMListener{
			Protocol: &input.Protocol,
		}

		if input.CertificateUrl != "" {
			output.CertificateUrl = &input.CertificateUrl
		}
		outputList = append(outputList, output)
	}
	return &outputList
}

func expandScheduledEventsProfileModel(inputList []ScheduledEventsProfileModel) *fleets.ScheduledEventsProfile {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.ScheduledEventsProfile{
		OsImageNotificationProfile:   expandOSImageNotificationProfileModel(input.OsImageNotificationProfile),
		TerminateNotificationProfile: expandTerminateNotificationProfileModel(input.TerminateNotificationProfile),
	}

	return &output
}

func expandOSImageNotificationProfileModel(inputList []OSImageNotificationProfileModel) *fleets.OSImageNotificationProfile {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.OSImageNotificationProfile{
		Enable: &input.Enable,
	}
	if input.NotBeforeTimeout != "" {
		output.NotBeforeTimeout = &input.NotBeforeTimeout
	}

	return &output
}

func expandTerminateNotificationProfileModel(inputList []TerminateNotificationProfileModel) *fleets.TerminateNotificationProfile {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.TerminateNotificationProfile{
		Enable: &input.Enable,
	}
	if input.NotBeforeTimeout != "" {
		output.NotBeforeTimeout = &input.NotBeforeTimeout
	}

	return &output
}

func expandSecurityPostureReferenceModel(inputList []SecurityPostureReferenceModel) *fleets.SecurityPostureReference {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.SecurityPostureReference{
		ExcludeExtensions: &input.ExcludeExtensions,
		IsOverridable:     &input.IsOverridable,
	}
	if input.Id != "" {
		output.Id = &input.Id
	}

	return &output
}

func expandSecurityProfileModel(inputList []SecurityProfileModel) *fleets.SecurityProfile {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.SecurityProfile{
		EncryptionAtHost:   &input.EncryptionAtHost,
		EncryptionIdentity: expandEncryptionIdentityModel(input.EncryptionIdentity),
		ProxyAgentSettings: expandProxyAgentSettingsModel(input.ProxyAgentSettings),
		SecurityType:       &input.SecurityType,
		UefiSettings:       expandUefiSettingsModel(input.UefiSettings),
	}

	return &output
}

func expandEncryptionIdentityModel(inputList []EncryptionIdentityModel) *fleets.EncryptionIdentity {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.EncryptionIdentity{}
	if input.UserAssignedIdentityResourceId != "" {
		output.UserAssignedIdentityResourceId = &input.UserAssignedIdentityResourceId
	}

	return &output
}

func expandProxyAgentSettingsModel(inputList []ProxyAgentSettingsModel) *fleets.ProxyAgentSettings {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.ProxyAgentSettings{
		Enabled:          &input.Enabled,
		KeyIncarnationId: &input.KeyIncarnationId,
		Mode:             &input.Mode,
	}

	return &output
}

func expandUefiSettingsModel(inputList []UefiSettingsModel) *fleets.UefiSettings {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.UefiSettings{
		SecureBootEnabled: &input.SecureBootEnabled,
		VTpmEnabled:       &input.VTpmEnabled,
	}

	return &output
}

func expandServiceArtifactReferenceModel(inputList []ServiceArtifactReferenceModel) *fleets.ServiceArtifactReference {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.ServiceArtifactReference{}
	if input.Id != "" {
		output.Id = &input.Id
	}

	return &output
}

func expandVirtualMachineScaleSetStorageProfileModel(inputList []StorageProfileModel) *fleets.VirtualMachineScaleSetStorageProfile {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.VirtualMachineScaleSetStorageProfile{
		DataDisks:          expandVirtualMachineScaleSetDataDiskModelArray(input.DataDisks),
		DiskControllerType: &input.DiskControllerType,
		ImageReference:     expandImageReferenceModel(input.ImageReference),
		OsDisk:             expandVirtualMachineScaleSetOSDiskModel(input.OsDisk),
	}

	return &output
}

func expandVirtualMachineScaleSetDataDiskModelArray(inputList []VirtualMachineScaleSetDataDiskModel) *[]fleets.VirtualMachineScaleSetDataDisk {
	var outputList []fleets.VirtualMachineScaleSetDataDisk
	for _, v := range inputList {
		input := v
		output := fleets.VirtualMachineScaleSetDataDisk{
			Caching:                 &input.Caching,
			CreateOption:            input.CreateOption,
			DeleteOption:            &input.DeleteOption,
			DiskIOPSReadWrite:       &input.DiskIOPSReadWrite,
			DiskMBpsReadWrite:       &input.DiskMBpsReadWrite,
			DiskSizeGB:              &input.DiskSizeGB,
			Lun:                     input.Lun,
			ManagedDisk:             expandVirtualMachineScaleSetManagedDiskParametersModel(input.ManagedDisk),
			WriteAcceleratorEnabled: &input.WriteAcceleratorEnabled,
		}

		if input.Name != "" {
			output.Name = &input.Name
		}
		outputList = append(outputList, output)
	}
	return &outputList
}

func expandVirtualMachineScaleSetManagedDiskParametersModel(inputList []VirtualMachineScaleSetManagedDiskParametersModel) *fleets.VirtualMachineScaleSetManagedDiskParameters {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.VirtualMachineScaleSetManagedDiskParameters{
		DiskEncryptionSet:  expandDiskEncryptionSetParametersModel(input.DiskEncryptionSet),
		SecurityProfile:    expandVMDiskSecurityProfileModel(input.SecurityProfile),
		StorageAccountType: &input.StorageAccountType,
	}

	return &output
}

func expandDiskEncryptionSetParametersModel(inputList []DiskEncryptionSetParametersModel) *fleets.DiskEncryptionSetParameters {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.DiskEncryptionSetParameters{}
	if input.Id != "" {
		output.Id = &input.Id
	}

	return &output
}

func expandVMDiskSecurityProfileModel(inputList []VMDiskSecurityProfileModel) *fleets.VMDiskSecurityProfile {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.VMDiskSecurityProfile{
		DiskEncryptionSet:      expandDiskEncryptionSetParametersModel(input.DiskEncryptionSet),
		SecurityEncryptionType: &input.SecurityEncryptionType,
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
		output.CommunityGalleryImageId = &input.CommunityGalleryImageId
	}

	if input.Id != "" {
		output.Id = &input.Id
	}

	if input.Offer != "" {
		output.Offer = &input.Offer
	}

	if input.Publisher != "" {
		output.Publisher = &input.Publisher
	}

	if input.SharedGalleryImageId != "" {
		output.SharedGalleryImageId = &input.SharedGalleryImageId
	}

	if input.Sku != "" {
		output.Sku = &input.Sku
	}

	if input.Version != "" {
		output.Version = &input.Version
	}

	return &output
}

func expandVirtualMachineScaleSetOSDiskModel(inputList []VirtualMachineScaleSetOSDiskModel) *fleets.VirtualMachineScaleSetOSDisk {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.VirtualMachineScaleSetOSDisk{
		Caching:                 &input.Caching,
		CreateOption:            input.CreateOption,
		DeleteOption:            &input.DeleteOption,
		DiffDiskSettings:        expandDiffDiskSettingsModel(input.DiffDiskSettings),
		DiskSizeGB:              &input.DiskSizeGB,
		Image:                   expandVirtualHardDiskModel(input.Image),
		ManagedDisk:             expandVirtualMachineScaleSetManagedDiskParametersModel(input.ManagedDisk),
		OsType:                  &input.OsType,
		VhdContainers:           &input.VhdContainers,
		WriteAcceleratorEnabled: &input.WriteAcceleratorEnabled,
	}
	if input.Name != "" {
		output.Name = &input.Name
	}

	return &output
}

func expandDiffDiskSettingsModel(inputList []DiffDiskSettingsModel) *fleets.DiffDiskSettings {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.DiffDiskSettings{
		Option:    &input.Option,
		Placement: &input.Placement,
	}

	return &output
}

func expandVirtualHardDiskModel(inputList []VirtualHardDiskModel) *fleets.VirtualHardDisk {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.VirtualHardDisk{}
	if input.Uri != "" {
		output.Uri = &input.Uri
	}

	return &output
}

func expandComputeProfileModel(inputList []ComputeProfileModel) (*fleets.ComputeProfile, error) {
	if len(inputList) == 0 {
		return nil, nil
	}
	input := &inputList[0]
	output := fleets.ComputeProfile{
		AdditionalVirtualMachineCapabilities: expandAdditionalCapabilitiesModel(input.AdditionalVirtualMachineCapabilities),
		PlatformFaultDomainCount:             &input.PlatformFaultDomainCount,
	}

	baseVirtualMachineProfileValue, err := expandBaseVirtualMachineProfileModel(input.BaseVirtualMachineProfile)
	if err != nil {
		return nil, err
	}

	output.BaseVirtualMachineProfile = pointer.From(baseVirtualMachineProfileValue)

	if input.ComputeApiVersion != "" {
		output.ComputeApiVersion = &input.ComputeApiVersion
	}

	return &output, nil
}

func expandAdditionalCapabilitiesModel(inputList []AdditionalCapabilitiesModel) *fleets.AdditionalCapabilities {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.AdditionalCapabilities{
		HibernationEnabled: &input.HibernationEnabled,
		UltraSSDEnabled:    &input.UltraSSDEnabled,
	}

	return &output
}

func expandPlanModel(inputList []PlanModel) *fleets.Plan {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.Plan{
		Name:      input.Name,
		Product:   input.Product,
		Publisher: input.Publisher,
	}
	if input.PromotionCode != "" {
		output.PromotionCode = &input.PromotionCode
	}

	if input.Version != "" {
		output.Version = &input.Version
	}

	return &output
}

func expandRegularPriorityProfileModel(inputList []RegularPriorityProfileModel) *fleets.RegularPriorityProfile {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.RegularPriorityProfile{
		AllocationStrategy: &input.AllocationStrategy,
		Capacity:           &input.Capacity,
		MinCapacity:        &input.MinCapacity,
	}

	return &output
}

func expandSpotPriorityProfileModel(inputList []SpotPriorityProfileModel) *fleets.SpotPriorityProfile {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.SpotPriorityProfile{
		AllocationStrategy: &input.AllocationStrategy,
		Capacity:           &input.Capacity,
		EvictionPolicy:     &input.EvictionPolicy,
		Maintain:           &input.Maintain,
		MaxPricePerVM:      &input.MaxPricePerVM,
		MinCapacity:        &input.MinCapacity,
	}

	return &output
}

func expandVMAttributesModel(inputList []VMAttributesModel) *fleets.VMAttributes {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.VMAttributes{
		AcceleratorCount:          expandVMAttributeMinMaxIntegerModel(input.AcceleratorCount),
		AcceleratorManufacturers:  expandAcceleratorManufacturerModelArray(input.AcceleratorManufacturers),
		AcceleratorSupport:        pointer.To(fleets.VMAttributeSupport(input.AcceleratorSupport)),
		AcceleratorTypes:          expandAcceleratorTypeModelArray(input.AcceleratorTypes),
		ArchitectureTypes:         expandArchitectureTypeModelArray(input.ArchitectureTypes),
		BurstableSupport:          &input.BurstableSupport,
		CpuManufacturers:          expandCPUManufacturerModelArray(input.CpuManufacturers),
		DataDiskCount:             expandVMAttributeMinMaxIntegerModel(input.DataDiskCount),
		ExcludedVMSizes:           &input.ExcludedVMSizes,
		LocalStorageDiskTypes:     expandLocalStorageDiskTypeModelArray(input.LocalStorageDiskTypes),
		LocalStorageInGiB:         expandVMAttributeMinMaxDoubleModel(input.LocalStorageInGiB),
		LocalStorageSupport:       &input.LocalStorageSupport,
		MemoryInGiBPerVCPU:        expandVMAttributeMinMaxDoubleModel(input.MemoryInGiBPerVCPU),
		NetworkBandwidthInMbps:    expandVMAttributeMinMaxDoubleModel(input.NetworkBandwidthInMbps),
		NetworkInterfaceCount:     expandVMAttributeMinMaxIntegerModel(input.NetworkInterfaceCount),
		RdmaNetworkInterfaceCount: expandVMAttributeMinMaxIntegerModel(input.RdmaNetworkInterfaceCount),
		RdmaSupport:               &input.RdmaSupport,
		VMCategories:              expandVMCategoryModelArray(input.VMCategories),
	}

	output.MemoryInGiB = pointer.From(expandVMAttributeMinMaxDoubleModel(input.MemoryInGiB))

	output.VCPUCount = pointer.From(expandVMAttributeMinMaxIntegerModel(input.VCPUCount))

	return &output
}

func expandVMAttributeMinMaxIntegerModel(inputList []VMAttributeMinMaxIntegerModel) *fleets.VMAttributeMinMaxInteger {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.VMAttributeMinMaxInteger{
		Max: &input.Max,
		Min: &input.Min,
	}

	return &output
}

func expandAcceleratorManufacturerModelArray(inputList []AcceleratorManufacturerModel) *[]fleets.AcceleratorManufacturer {
	var outputList []fleets.AcceleratorManufacturer
	for _, v := range inputList {
		input := v
		output := fleets.AcceleratorManufacturer{}

		outputList = append(outputList, output)
	}
	return &outputList
}

func expandAcceleratorTypeModelArray(inputList []AcceleratorTypeModel) *[]fleets.AcceleratorType {
	var outputList []fleets.AcceleratorType
	for _, v := range inputList {
		input := v
		output := fleets.AcceleratorType{}

		outputList = append(outputList, output)
	}
	return &outputList
}

func expandArchitectureTypeModelArray(inputList []ArchitectureTypeModel) *[]fleets.ArchitectureType {
	var outputList []fleets.ArchitectureType
	for _, v := range inputList {
		input := v
		output := fleets.ArchitectureType{}

		outputList = append(outputList, output)
	}
	return &outputList
}

func expandCPUManufacturerModelArray(inputList []CPUManufacturerModel) *[]fleets.CPUManufacturer {
	var outputList []fleets.CPUManufacturer
	for _, v := range inputList {
		input := v
		output := fleets.CPUManufacturer{}

		outputList = append(outputList, output)
	}
	return &outputList
}

func expandLocalStorageDiskTypeModelArray(inputList []LocalStorageDiskTypeModel) *[]fleets.LocalStorageDiskType {
	var outputList []fleets.LocalStorageDiskType
	for _, v := range inputList {
		input := v
		output := fleets.LocalStorageDiskType{}

		outputList = append(outputList, output)
	}
	return &outputList
}

func expandVMAttributeMinMaxDoubleModel(inputList []VMAttributeMinMaxDoubleModel) *fleets.VMAttributeMinMaxDouble {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.VMAttributeMinMaxDouble{
		Max: &input.Max,
		Min: &input.Min,
	}

	return &output
}

func expandVMCategoryModelArray(inputList []VMCategoryModel) *[]fleets.VMCategory {
	var outputList []fleets.VMCategory
	for _, v := range inputList {
		input := v
		output := fleets.VMCategory{}

		outputList = append(outputList, output)
	}
	return &outputList
}

func expandVMSizeProfileModelArray(inputList []VMSizeProfileModel) *[]fleets.VMSizeProfile {
	var outputList []fleets.VMSizeProfile
	for _, v := range inputList {
		input := v
		output := fleets.VMSizeProfile{
			Name: input.Name,
			Rank: &input.Rank,
		}

		outputList = append(outputList, output)
	}
	return &outputList
}

func flattenAdditionalLocationsProfileModel(input *fleets.AdditionalLocationsProfile) ([]AdditionalLocationsProfileModel, error) {
	var outputList []AdditionalLocationsProfileModel
	if input == nil {
		return outputList, nil
	}
	output := AdditionalLocationsProfileModel{}
	locationProfilesValue, err := flattenLocationProfileModelArray(&input.LocationProfiles)
	if err != nil {
		return nil, err
	}

	output.LocationProfiles = locationProfilesValue

	return append(outputList, output), nil
}

func flattenLocationProfileModelArray(inputList *[]fleets.LocationProfile) ([]LocationProfileModel, error) {
	var outputList []LocationProfileModel
	if inputList == nil {
		return outputList, nil
	}
	for _, input := range *inputList {
		output := LocationProfileModel{
			Location: input.Location,
		}

		virtualMachineProfileOverrideValue, err := flattenBaseVirtualMachineProfileModel(input.VirtualMachineProfile)
		if err != nil {
			return nil, err
		}

		output.VirtualMachineProfile = virtualMachineProfileOverrideValue
		outputList = append(outputList, output)
	}
	return outputList, nil
}

func flattenBaseVirtualMachineProfileModel(input *fleets.BaseVirtualMachineProfile) ([]VirtualMachineProfileModel, error) {
	var outputList []VirtualMachineProfileModel
	if input == nil {
		return outputList, nil
	}
	output := VirtualMachineProfileModel{
		ApplicationProfile:       flattenApplicationProfileModel(input.ApplicationProfile),
		CapacityReservation:      flattenCapacityReservationProfileModel(input.CapacityReservation),
		DiagnosticsProfile:       flattenDiagnosticsProfileModel(input.DiagnosticsProfile),
		HardwareProfile:          flattenVirtualMachineScaleSetHardwareProfileModel(input.HardwareProfile),
		NetworkProfile:           flattenVirtualMachineScaleSetNetworkProfileModel(input.NetworkProfile),
		OsProfile:                flattenVirtualMachineScaleSetOSProfileModel(input.OsProfile),
		ScheduledEventsProfile:   flattenScheduledEventsProfileModel(input.ScheduledEventsProfile),
		SecurityPostureReference: flattenSecurityPostureReferenceModel(input.SecurityPostureReference),
		SecurityProfile:          flattenSecurityProfileModel(input.SecurityProfile),
		ServiceArtifactReference: flattenServiceArtifactReferenceModel(input.ServiceArtifactReference),
		StorageProfile:           flattenVirtualMachineScaleSetStorageProfileModel(input.StorageProfile),
	}
	extensionProfileValue, err := flattenVirtualMachineScaleSetExtensionProfileModel(input.ExtensionProfile)
	if err != nil {
		return nil, err
	}

	output.ExtensionProfile = extensionProfileValue

	if input.LicenseType != nil {
		output.LicenseType = *input.LicenseType
	}

	if input.TimeCreated != nil {
		output.TimeCreated = *input.TimeCreated
	}

	if input.UserData != nil {
		output.UserData = *input.UserData
	}

	return append(outputList, output), nil
}

func flattenApplicationProfileModel(input *fleets.ApplicationProfile) []ApplicationProfileModel {
	var outputList []ApplicationProfileModel
	if input == nil {
		return outputList
	}
	output := ApplicationProfileModel{
		GalleryApplications: flattenVMGalleryApplicationModelArray(input.GalleryApplications),
	}

	return append(outputList, output)
}

func flattenVMGalleryApplicationModelArray(inputList *[]fleets.VMGalleryApplication) []VMGalleryApplicationModel {
	var outputList []VMGalleryApplicationModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := VMGalleryApplicationModel{
			PackageReferenceId: input.PackageReferenceId,
		}

		if input.ConfigurationReference != nil {
			output.ConfigurationReference = *input.ConfigurationReference
		}

		if input.EnableAutomaticUpgrade != nil {
			output.EnableAutomaticUpgrade = *input.EnableAutomaticUpgrade
		}

		if input.Order != nil {
			output.Order = *input.Order
		}

		if input.Tags != nil {
			output.Tags = *input.Tags
		}

		if input.TreatFailureAsDeploymentFailure != nil {
			output.TreatFailureAsDeploymentFailure = *input.TreatFailureAsDeploymentFailure
		}
		outputList = append(outputList, output)
	}
	return outputList
}

func flattenCapacityReservationProfileModel(input *fleets.CapacityReservationProfile) []CapacityReservationProfileModel {
	var outputList []CapacityReservationProfileModel
	if input == nil {
		return outputList
	}
	output := CapacityReservationProfileModel{
		CapacityReservationGroup: flattenSubResourceModel(input.CapacityReservationGroup),
	}

	return append(outputList, output)
}

func flattenSubResourceModel(input *fleets.SubResource) []SubResourceModel {
	var outputList []SubResourceModel
	if input == nil {
		return outputList
	}
	output := SubResourceModel{}
	if input.Id != nil {
		output.Id = *input.Id
	}

	return append(outputList, output)
}

func flattenDiagnosticsProfileModel(input *fleets.DiagnosticsProfile) []DiagnosticsProfileModel {
	var outputList []DiagnosticsProfileModel
	if input == nil {
		return outputList
	}
	output := DiagnosticsProfileModel{
		BootDiagnostics: flattenBootDiagnosticsModel(input.BootDiagnostics),
	}

	return append(outputList, output)
}

func flattenBootDiagnosticsModel(input *fleets.BootDiagnostics) []BootDiagnosticsModel {
	var outputList []BootDiagnosticsModel
	if input == nil {
		return outputList
	}
	output := BootDiagnosticsModel{}
	if input.Enabled != nil {
		output.Enabled = *input.Enabled
	}

	if input.StorageUri != nil {
		output.StorageUri = *input.StorageUri
	}

	return append(outputList, output)
}

func flattenVirtualMachineScaleSetExtensionProfileModel(input *fleets.VirtualMachineScaleSetExtensionProfile) ([]VirtualMachineScaleSetExtensionProfileModel, error) {
	var outputList []VirtualMachineScaleSetExtensionProfileModel
	if input == nil {
		return outputList, nil
	}
	output := VirtualMachineScaleSetExtensionProfileModel{}
	extensionsValue, err := flattenVirtualMachineScaleSetExtensionModelArray(input.Extensions)
	if err != nil {
		return nil, err
	}

	output.Extensions = extensionsValue

	if input.ExtensionsTimeBudget != nil {
		output.ExtensionsTimeBudget = *input.ExtensionsTimeBudget
	}

	return append(outputList, output), nil
}

func flattenVirtualMachineScaleSetExtensionModelArray(inputList *[]fleets.VirtualMachineScaleSetExtension) ([]VirtualMachineScaleSetExtensionModel, error) {
	var outputList []VirtualMachineScaleSetExtensionModel
	if inputList == nil {
		return outputList, nil
	}
	for _, input := range *inputList {
		output := VirtualMachineScaleSetExtensionModel{}

		if input.Id != nil {
			output.Id = *input.Id
		}

		if input.Name != nil {
			output.Name = *input.Name
		}

		propertiesValue, err := flattenVirtualMachineScaleSetExtensionModel(input.Properties)
		if err != nil {
			return nil, err
		}

		output.Properties = propertiesValue

		if input.Type != nil {
			output.Type = *input.Type
		}
		outputList = append(outputList, output)
	}
	return outputList, nil
}

func flattenVirtualMachineScaleSetExtensionModel(input *fleets.VirtualMachineScaleSetExtensionProperties) ([]VirtualMachineScaleSetExtensionModel, error) {
	var outputList []VirtualMachineScaleSetExtensionModel
	if input == nil {
		return outputList, nil
	}
	output := VirtualMachineScaleSetExtensionModel{
		ProtectedSettingsFromKeyVault: flattenKeyVaultSecretReferenceModel(input.ProtectedSettingsFromKeyVault),
	}
	if input.AutoUpgradeMinorVersion != nil {
		output.AutoUpgradeMinorVersion = *input.AutoUpgradeMinorVersion
	}

	if input.EnableAutomaticUpgrade != nil {
		output.EnableAutomaticUpgrade = *input.EnableAutomaticUpgrade
	}

	if input.ForceUpdateTag != nil {
		output.ForceUpdateTag = *input.ForceUpdateTag
	}

	if input.ProtectedSettings != nil && *input.ProtectedSettings != nil {

		protectedSettingsValue, err := json.Marshal(*input.ProtectedSettings)
		if err != nil {
			return outputList, err
		}

		output.ProtectedSettings = string(protectedSettingsValue)
	}

	if input.ProvisionAfterExtensions != nil {
		output.ProvisionAfterExtensions = *input.ProvisionAfterExtensions
	}

	if input.ProvisioningState != nil {
		output.ProvisioningState = *input.ProvisioningState
	}

	if input.Publisher != nil {
		output.Publisher = *input.Publisher
	}

	if input.Settings != nil && *input.Settings != nil {

		settingsValue, err := json.Marshal(*input.Settings)
		if err != nil {
			return outputList, err
		}

		output.Settings = string(settingsValue)
	}

	if input.SuppressFailures != nil {
		output.SuppressFailures = *input.SuppressFailures
	}

	if input.Type != nil {
		output.Type = *input.Type
	}

	if input.TypeHandlerVersion != nil {
		output.TypeHandlerVersion = *input.TypeHandlerVersion
	}

	return append(outputList, output), nil
}

func flattenKeyVaultSecretReferenceModel(input *fleets.KeyVaultSecretReference) []KeyVaultSecretReferenceModel {
	var outputList []KeyVaultSecretReferenceModel
	if input == nil {
		return outputList
	}
	output := KeyVaultSecretReferenceModel{
		SecretUrl:   input.SecretUrl,
		SourceVault: flattenSubResourceModel(&input.SourceVault),
	}

	return append(outputList, output)
}

func flattenVirtualMachineScaleSetHardwareProfileModel(input *fleets.VirtualMachineScaleSetHardwareProfile) []VirtualMachineScaleSetHardwareProfileModel {
	var outputList []VirtualMachineScaleSetHardwareProfileModel
	if input == nil {
		return outputList
	}
	output := VirtualMachineScaleSetHardwareProfileModel{
		VMSizeProperties: flattenVMSizePropertiesModel(input.VMSizeProperties),
	}

	return append(outputList, output)
}

func flattenVMSizePropertiesModel(input *fleets.VMSizeProperties) []VMSizePropertiesModel {
	var outputList []VMSizePropertiesModel
	if input == nil {
		return outputList
	}
	output := VMSizePropertiesModel{}
	if input.VCPUsAvailable != nil {
		output.VCPUsAvailable = *input.VCPUsAvailable
	}

	if input.VCPUsPerCore != nil {
		output.VCPUsPerCore = *input.VCPUsPerCore
	}

	return append(outputList, output)
}

func flattenVirtualMachineScaleSetNetworkProfileModel(input *fleets.VirtualMachineScaleSetNetworkProfile) []NetworkProfileModel {
	var outputList []NetworkProfileModel
	if input == nil {
		return outputList
	}
	output := NetworkProfileModel{
		HealthProbe:                    flattenApiEntityReferenceModel(input.HealthProbe),
		NetworkInterfaceConfigurations: flattenVirtualMachineScaleSetNetworkConfigurationModelArray(input.NetworkInterfaceConfigurations),
	}
	if input.NetworkApiVersion != nil {
		output.NetworkApiVersion = *input.NetworkApiVersion
	}

	return append(outputList, output)
}

func flattenApiEntityReferenceModel(input *fleets.ApiEntityReference) []ApiEntityReferenceModel {
	var outputList []ApiEntityReferenceModel
	if input == nil {
		return outputList
	}
	output := ApiEntityReferenceModel{}
	if input.Id != nil {
		output.Id = *input.Id
	}

	return append(outputList, output)
}

func flattenVirtualMachineScaleSetNetworkConfigurationModelArray(inputList *[]fleets.VirtualMachineScaleSetNetworkConfiguration) []NetworkConfigurationModel {
	var outputList []NetworkConfigurationModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := NetworkConfigurationModel{
			Name:       input.Name,
			Properties: flattenVirtualMachineScaleSetNetworkConfigurationPropertiesModel(input.Properties),
		}

		outputList = append(outputList, output)
	}
	return outputList
}

func flattenVirtualMachineScaleSetNetworkConfigurationPropertiesModel(input *fleets.VirtualMachineScaleSetNetworkConfigurationProperties) []VirtualMachineScaleSetNetworkConfigurationPropertiesModel {
	var outputList []VirtualMachineScaleSetNetworkConfigurationPropertiesModel
	if input == nil {
		return outputList
	}
	output := VirtualMachineScaleSetNetworkConfigurationPropertiesModel{
		DnsSettings:          flattenVirtualMachineScaleSetNetworkConfigurationDnsSettingsModel(input.DnsSettings),
		IPConfigurations:     flattenVirtualMachineScaleSetIPConfigurationModelArray(&input.IPConfigurations),
		NetworkSecurityGroup: flattenSubResourceModel(input.NetworkSecurityGroup),
	}
	if input.AuxiliaryMode != nil {
		output.AuxiliaryMode = *input.AuxiliaryMode
	}

	if input.AuxiliarySku != nil {
		output.AuxiliarySku = *input.AuxiliarySku
	}

	if input.DeleteOption != nil {
		output.DeleteOption = *input.DeleteOption
	}

	if input.DisableTcpStateTracking != nil {
		output.DisableTcpStateTracking = *input.DisableTcpStateTracking
	}

	if input.EnableAcceleratedNetworking != nil {
		output.EnableAcceleratedNetworking = *input.EnableAcceleratedNetworking
	}

	if input.EnableFpga != nil {
		output.EnableFpga = *input.EnableFpga
	}

	if input.EnableIPForwarding != nil {
		output.EnableIPForwarding = *input.EnableIPForwarding
	}

	if input.Primary != nil {
		output.Primary = *input.Primary
	}

	return append(outputList, output)
}

func flattenVirtualMachineScaleSetNetworkConfigurationDnsSettingsModel(input *fleets.VirtualMachineScaleSetNetworkConfigurationDnsSettings) []VirtualMachineScaleSetNetworkConfigurationDnsSettingsModel {
	var outputList []VirtualMachineScaleSetNetworkConfigurationDnsSettingsModel
	if input == nil {
		return outputList
	}
	output := VirtualMachineScaleSetNetworkConfigurationDnsSettingsModel{}
	if input.DnsServers != nil {
		output.DnsServers = *input.DnsServers
	}

	return append(outputList, output)
}

func flattenVirtualMachineScaleSetIPConfigurationModelArray(inputList *[]fleets.VirtualMachineScaleSetIPConfiguration) []VirtualMachineScaleSetIPConfigurationModel {
	var outputList []VirtualMachineScaleSetIPConfigurationModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := VirtualMachineScaleSetIPConfigurationModel{
			Name:       input.Name,
			Properties: flattenVirtualMachineScaleSetIPConfigurationModel(input.Properties),
		}

		outputList = append(outputList, output)
	}
	return outputList
}

func flattenVirtualMachineScaleSetIPConfigurationModel(input *fleets.VirtualMachineScaleSetIPConfigurationProperties) []VirtualMachineScaleSetIPConfigurationModel {
	var outputList []VirtualMachineScaleSetIPConfigurationModel
	if input == nil {
		return outputList
	}
	output := VirtualMachineScaleSetIPConfigurationModel{
		ApplicationGatewayBackendAddressPools: flattenSubResourceModelArray(input.ApplicationGatewayBackendAddressPools),
		ApplicationSecurityGroups:             flattenSubResourceModelArray(input.ApplicationSecurityGroups),
		LoadBalancerBackendAddressPools:       flattenSubResourceModelArray(input.LoadBalancerBackendAddressPools),
		LoadBalancerInboundNatPools:           flattenSubResourceModelArray(input.LoadBalancerInboundNatPools),
		PublicIPAddressConfiguration:          flattenVirtualMachineScaleSetPublicIPAddressConfigurationModel(input.PublicIPAddressConfiguration),
		Subnet:                                flattenApiEntityReferenceModel(input.Subnet),
	}
	if input.Primary != nil {
		output.Primary = *input.Primary
	}

	if input.PrivateIPAddressVersion != nil {
		output.PrivateIPAddressVersion = *input.PrivateIPAddressVersion
	}

	return append(outputList, output)
}

func flattenVirtualMachineScaleSetPublicIPAddressConfigurationModel(input *fleets.VirtualMachineScaleSetPublicIPAddressConfiguration) []VirtualMachineScaleSetPublicIPAddressConfigurationModel {
	var outputList []VirtualMachineScaleSetPublicIPAddressConfigurationModel
	if input == nil {
		return outputList
	}
	output := VirtualMachineScaleSetPublicIPAddressConfigurationModel{
		Name:       input.Name,
		Properties: flattenVirtualMachineScaleSetPublicIPAddressConfigurationModel(input.Properties),
		Sku:        flattenPublicIPAddressSkuModel(input.Sku),
	}

	return append(outputList, output)
}

func flattenVirtualMachineScaleSetPublicIPAddressConfigurationModel(input *fleets.VirtualMachineScaleSetPublicIPAddressConfigurationProperties) []VirtualMachineScaleSetPublicIPAddressConfigurationModel {
	var outputList []VirtualMachineScaleSetPublicIPAddressConfigurationModel
	if input == nil {
		return outputList
	}
	output := VirtualMachineScaleSetPublicIPAddressConfigurationModel{
		DnsSettings:    flattenVirtualMachineScaleSetPublicIPAddressConfigurationDnsSettingsModel(input.DnsSettings),
		IPTags:         flattenVirtualMachineScaleSetIPTagModelArray(input.IPTags),
		PublicIPPrefix: flattenSubResourceModel(input.PublicIPPrefix),
	}
	if input.DeleteOption != nil {
		output.DeleteOption = *input.DeleteOption
	}

	if input.IdleTimeoutInMinutes != nil {
		output.IdleTimeoutInMinutes = *input.IdleTimeoutInMinutes
	}

	if input.PublicIPAddressVersion != nil {
		output.PublicIPAddressVersion = *input.PublicIPAddressVersion
	}

	return append(outputList, output)
}

func flattenVirtualMachineScaleSetPublicIPAddressConfigurationDnsSettingsModel(input *fleets.VirtualMachineScaleSetPublicIPAddressConfigurationDnsSettings) []VirtualMachineScaleSetPublicIPAddressConfigurationDnsSettingsModel {
	var outputList []VirtualMachineScaleSetPublicIPAddressConfigurationDnsSettingsModel
	if input == nil {
		return outputList
	}
	output := VirtualMachineScaleSetPublicIPAddressConfigurationDnsSettingsModel{
		DomainNameLabel: input.DomainNameLabel,
	}
	if input.DomainNameLabelScope != nil {
		output.DomainNameLabelScope = *input.DomainNameLabelScope
	}

	return append(outputList, output)
}

func flattenVirtualMachineScaleSetIPTagModelArray(inputList *[]fleets.VirtualMachineScaleSetIPTag) []VirtualMachineScaleSetIPTagModel {
	var outputList []VirtualMachineScaleSetIPTagModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := VirtualMachineScaleSetIPTagModel{}

		if input.IPTagType != nil {
			output.IPTagType = *input.IPTagType
		}

		if input.Tag != nil {
			output.Tag = *input.Tag
		}
		outputList = append(outputList, output)
	}
	return outputList
}

func flattenPublicIPAddressSkuModel(input *fleets.PublicIPAddressSku) []PublicIPAddressSkuModel {
	var outputList []PublicIPAddressSkuModel
	if input == nil {
		return outputList
	}
	output := PublicIPAddressSkuModel{}
	if input.Name != nil {
		output.Name = *input.Name
	}

	if input.Tier != nil {
		output.Tier = *input.Tier
	}

	return append(outputList, output)
}

func flattenVirtualMachineScaleSetOSProfileModel(input *fleets.VirtualMachineScaleSetOSProfile) []OSProfileModel {
	var outputList []OSProfileModel
	if input == nil {
		return outputList
	}
	output := OSProfileModel{
		LinuxConfiguration:   flattenLinuxConfigurationModel(input.LinuxConfiguration),
		Secrets:              flattenVaultSecretGroupModelArray(input.Secrets),
		WindowsConfiguration: flattenWindowsConfigurationModel(input.WindowsConfiguration),
	}
	if input.AdminPassword != nil {
		output.AdminPassword = *input.AdminPassword
	}

	if input.AdminUsername != nil {
		output.AdminUsername = *input.AdminUsername
	}

	if input.AllowExtensionOperations != nil {
		output.AllowExtensionOperations = *input.AllowExtensionOperations
	}

	if input.ComputerNamePrefix != nil {
		output.ComputerNamePrefix = *input.ComputerNamePrefix
	}

	if input.CustomData != nil {
		output.CustomData = *input.CustomData
	}

	if input.RequireGuestProvisionSignal != nil {
		output.RequireGuestProvisionSignal = *input.RequireGuestProvisionSignal
	}

	return append(outputList, output)
}

func flattenLinuxConfigurationModel(input *fleets.LinuxConfiguration) []LinuxConfigurationModel {
	var outputList []LinuxConfigurationModel
	if input == nil {
		return outputList
	}
	output := LinuxConfigurationModel{
		PatchSettings: flattenLinuxPatchSettingsModel(input.PatchSettings),
		Ssh:           flattenSshConfigurationModel(input.Ssh),
	}
	if input.DisablePasswordAuthentication != nil {
		output.DisablePasswordAuthentication = *input.DisablePasswordAuthentication
	}

	if input.EnableVMAgentPlatformUpdates != nil {
		output.EnableVMAgentPlatformUpdates = *input.EnableVMAgentPlatformUpdates
	}

	if input.ProvisionVMAgent != nil {
		output.ProvisionVMAgent = *input.ProvisionVMAgent
	}

	return append(outputList, output)
}

func flattenLinuxPatchSettingsModel(input *fleets.LinuxPatchSettings) []LinuxPatchSettingsModel {
	var outputList []LinuxPatchSettingsModel
	if input == nil {
		return outputList
	}
	output := LinuxPatchSettingsModel{
		AutomaticByPlatformSettings: flattenLinuxVMGuestPatchAutomaticByPlatformSettingsModel(input.AutomaticByPlatformSettings),
	}
	if input.AssessmentMode != nil {
		output.AssessmentMode = *input.AssessmentMode
	}

	if input.PatchMode != nil {
		output.PatchMode = *input.PatchMode
	}

	return append(outputList, output)
}

func flattenLinuxVMGuestPatchAutomaticByPlatformSettingsModel(input *fleets.LinuxVMGuestPatchAutomaticByPlatformSettings) []LinuxVMGuestPatchAutomaticByPlatformSettingsModel {
	var outputList []LinuxVMGuestPatchAutomaticByPlatformSettingsModel
	if input == nil {
		return outputList
	}
	output := LinuxVMGuestPatchAutomaticByPlatformSettingsModel{}
	if input.BypassPlatformSafetyChecksOnUserSchedule != nil {
		output.BypassPlatformSafetyChecksOnUserSchedule = *input.BypassPlatformSafetyChecksOnUserSchedule
	}

	if input.RebootSetting != nil {
		output.RebootSetting = *input.RebootSetting
	}

	return append(outputList, output)
}

func flattenSshConfigurationModel(input *fleets.SshConfiguration) []SshConfigurationModel {
	var outputList []SshConfigurationModel
	if input == nil {
		return outputList
	}
	output := SshConfigurationModel{
		PublicKeys: flattenSshPublicKeyModelArray(input.PublicKeys),
	}

	return append(outputList, output)
}

func flattenSshPublicKeyModelArray(inputList *[]fleets.SshPublicKey) []SshPublicKeyModel {
	var outputList []SshPublicKeyModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := SshPublicKeyModel{}

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

func flattenVaultSecretGroupModelArray(inputList *[]fleets.VaultSecretGroup) []VaultSecretGroupModel {
	var outputList []VaultSecretGroupModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := VaultSecretGroupModel{
			SourceVault:       flattenSubResourceModel(input.SourceVault),
			VaultCertificates: flattenVaultCertificateModelArray(input.VaultCertificates),
		}

		outputList = append(outputList, output)
	}
	return outputList
}

func flattenVaultCertificateModelArray(inputList *[]fleets.VaultCertificate) []VaultCertificateModel {
	var outputList []VaultCertificateModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := VaultCertificateModel{}

		if input.CertificateStore != nil {
			output.CertificateStore = *input.CertificateStore
		}

		if input.CertificateUrl != nil {
			output.CertificateUrl = *input.CertificateUrl
		}
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
		AdditionalUnattendContent: flattenAdditionalUnattendContentModelArray(input.AdditionalUnattendContent),
		PatchSettings:             flattenPatchSettingsModel(input.PatchSettings),
		WinRM:                     flattenWinRMConfigurationModel(input.WinRM),
	}
	if input.EnableAutomaticUpdates != nil {
		output.EnableAutomaticUpdates = *input.EnableAutomaticUpdates
	}

	if input.EnableVMAgentPlatformUpdates != nil {
		output.EnableVMAgentPlatformUpdates = *input.EnableVMAgentPlatformUpdates
	}

	if input.ProvisionVMAgent != nil {
		output.ProvisionVMAgent = *input.ProvisionVMAgent
	}

	if input.TimeZone != nil {
		output.TimeZone = *input.TimeZone
	}

	return append(outputList, output)
}

func flattenAdditionalUnattendContentModelArray(inputList *[]fleets.AdditionalUnattendContent) []AdditionalUnattendContentModel {
	var outputList []AdditionalUnattendContentModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := AdditionalUnattendContentModel{}

		if input.ComponentName != nil {
			output.ComponentName = *input.ComponentName
		}

		if input.Content != nil {
			output.Content = *input.Content
		}

		if input.PassName != nil {
			output.PassName = *input.PassName
		}

		if input.SettingName != nil {
			output.SettingName = *input.SettingName
		}
		outputList = append(outputList, output)
	}
	return outputList
}

func flattenPatchSettingsModel(input *fleets.PatchSettings) []PatchSettingsModel {
	var outputList []PatchSettingsModel
	if input == nil {
		return outputList
	}
	output := PatchSettingsModel{
		AutomaticByPlatformSettings: flattenWindowsVMGuestPatchAutomaticByPlatformSettingsModel(input.AutomaticByPlatformSettings),
	}
	if input.AssessmentMode != nil {
		output.AssessmentMode = *input.AssessmentMode
	}

	if input.EnableHotpatching != nil {
		output.EnableHotpatching = *input.EnableHotpatching
	}

	if input.PatchMode != nil {
		output.PatchMode = *input.PatchMode
	}

	return append(outputList, output)
}

func flattenWindowsVMGuestPatchAutomaticByPlatformSettingsModel(input *fleets.WindowsVMGuestPatchAutomaticByPlatformSettings) []WindowsVMGuestPatchAutomaticByPlatformSettingsModel {
	var outputList []WindowsVMGuestPatchAutomaticByPlatformSettingsModel
	if input == nil {
		return outputList
	}
	output := WindowsVMGuestPatchAutomaticByPlatformSettingsModel{}
	if input.BypassPlatformSafetyChecksOnUserSchedule != nil {
		output.BypassPlatformSafetyChecksOnUserSchedule = *input.BypassPlatformSafetyChecksOnUserSchedule
	}

	if input.RebootSetting != nil {
		output.RebootSetting = *input.RebootSetting
	}

	return append(outputList, output)
}

func flattenWinRMConfigurationModel(input *fleets.WinRMConfiguration) []WinRMConfigurationModel {
	var outputList []WinRMConfigurationModel
	if input == nil {
		return outputList
	}
	output := WinRMConfigurationModel{
		Listeners: flattenWinRMListenerModelArray(input.Listeners),
	}

	return append(outputList, output)
}

func flattenWinRMListenerModelArray(inputList *[]fleets.WinRMListener) []WinRMListenerModel {
	var outputList []WinRMListenerModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := WinRMListenerModel{}

		if input.CertificateUrl != nil {
			output.CertificateUrl = *input.CertificateUrl
		}

		if input.Protocol != nil {
			output.Protocol = *input.Protocol
		}
		outputList = append(outputList, output)
	}
	return outputList
}

func flattenScheduledEventsProfileModel(input *fleets.ScheduledEventsProfile) []ScheduledEventsProfileModel {
	var outputList []ScheduledEventsProfileModel
	if input == nil {
		return outputList
	}
	output := ScheduledEventsProfileModel{
		OsImageNotificationProfile:   flattenOSImageNotificationProfileModel(input.OsImageNotificationProfile),
		TerminateNotificationProfile: flattenTerminateNotificationProfileModel(input.TerminateNotificationProfile),
	}

	return append(outputList, output)
}

func flattenOSImageNotificationProfileModel(input *fleets.OSImageNotificationProfile) []OSImageNotificationProfileModel {
	var outputList []OSImageNotificationProfileModel
	if input == nil {
		return outputList
	}
	output := OSImageNotificationProfileModel{}
	if input.Enable != nil {
		output.Enable = *input.Enable
	}

	if input.NotBeforeTimeout != nil {
		output.NotBeforeTimeout = *input.NotBeforeTimeout
	}

	return append(outputList, output)
}

func flattenTerminateNotificationProfileModel(input *fleets.TerminateNotificationProfile) []TerminateNotificationProfileModel {
	var outputList []TerminateNotificationProfileModel
	if input == nil {
		return outputList
	}
	output := TerminateNotificationProfileModel{}
	if input.Enable != nil {
		output.Enable = *input.Enable
	}

	if input.NotBeforeTimeout != nil {
		output.NotBeforeTimeout = *input.NotBeforeTimeout
	}

	return append(outputList, output)
}

func flattenSecurityPostureReferenceModel(input *fleets.SecurityPostureReference) []SecurityPostureReferenceModel {
	var outputList []SecurityPostureReferenceModel
	if input == nil {
		return outputList
	}
	output := SecurityPostureReferenceModel{}
	if input.ExcludeExtensions != nil {
		output.ExcludeExtensions = *input.ExcludeExtensions
	}

	if input.Id != nil {
		output.Id = *input.Id
	}

	if input.IsOverridable != nil {
		output.IsOverridable = *input.IsOverridable
	}

	return append(outputList, output)
}

func flattenSecurityProfileModel(input *fleets.SecurityProfile) []SecurityProfileModel {
	var outputList []SecurityProfileModel
	if input == nil {
		return outputList
	}
	output := SecurityProfileModel{
		EncryptionIdentity: flattenEncryptionIdentityModel(input.EncryptionIdentity),
		ProxyAgentSettings: flattenProxyAgentSettingsModel(input.ProxyAgentSettings),
		UefiSettings:       flattenUefiSettingsModel(input.UefiSettings),
	}
	if input.EncryptionAtHost != nil {
		output.EncryptionAtHost = *input.EncryptionAtHost
	}

	if input.SecurityType != nil {
		output.SecurityType = *input.SecurityType
	}

	return append(outputList, output)
}

func flattenEncryptionIdentityModel(input *fleets.EncryptionIdentity) []EncryptionIdentityModel {
	var outputList []EncryptionIdentityModel
	if input == nil {
		return outputList
	}
	output := EncryptionIdentityModel{}
	if input.UserAssignedIdentityResourceId != nil {
		output.UserAssignedIdentityResourceId = *input.UserAssignedIdentityResourceId
	}

	return append(outputList, output)
}

func flattenProxyAgentSettingsModel(input *fleets.ProxyAgentSettings) []ProxyAgentSettingsModel {
	var outputList []ProxyAgentSettingsModel
	if input == nil {
		return outputList
	}
	output := ProxyAgentSettingsModel{}
	if input.Enabled != nil {
		output.Enabled = *input.Enabled
	}

	if input.KeyIncarnationId != nil {
		output.KeyIncarnationId = *input.KeyIncarnationId
	}

	if input.Mode != nil {
		output.Mode = *input.Mode
	}

	return append(outputList, output)
}

func flattenUefiSettingsModel(input *fleets.UefiSettings) []UefiSettingsModel {
	var outputList []UefiSettingsModel
	if input == nil {
		return outputList
	}
	output := UefiSettingsModel{}
	if input.SecureBootEnabled != nil {
		output.SecureBootEnabled = *input.SecureBootEnabled
	}

	if input.VTpmEnabled != nil {
		output.VTpmEnabled = *input.VTpmEnabled
	}

	return append(outputList, output)
}

func flattenServiceArtifactReferenceModel(input *fleets.ServiceArtifactReference) []ServiceArtifactReferenceModel {
	var outputList []ServiceArtifactReferenceModel
	if input == nil {
		return outputList
	}
	output := ServiceArtifactReferenceModel{}
	if input.Id != nil {
		output.Id = *input.Id
	}

	return append(outputList, output)
}

func flattenVirtualMachineScaleSetStorageProfileModel(input *fleets.VirtualMachineScaleSetStorageProfile) []StorageProfileModel {
	var outputList []StorageProfileModel
	if input == nil {
		return outputList
	}
	output := StorageProfileModel{
		DataDisks:      flattenVirtualMachineScaleSetDataDiskModelArray(input.DataDisks),
		ImageReference: flattenImageReferenceModel(input.ImageReference),
		OsDisk:         flattenVirtualMachineScaleSetOSDiskModel(input.OsDisk),
	}
	if input.DiskControllerType != nil {
		output.DiskControllerType = *input.DiskControllerType
	}

	return append(outputList, output)
}

func flattenVirtualMachineScaleSetDataDiskModelArray(inputList *[]fleets.VirtualMachineScaleSetDataDisk) []VirtualMachineScaleSetDataDiskModel {
	var outputList []VirtualMachineScaleSetDataDiskModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := VirtualMachineScaleSetDataDiskModel{
			CreateOption: input.CreateOption,
			Lun:          input.Lun,
			ManagedDisk:  flattenVirtualMachineScaleSetManagedDiskParametersModel(input.ManagedDisk),
		}

		if input.Caching != nil {
			output.Caching = *input.Caching
		}

		if input.DeleteOption != nil {
			output.DeleteOption = *input.DeleteOption
		}

		if input.DiskIOPSReadWrite != nil {
			output.DiskIOPSReadWrite = *input.DiskIOPSReadWrite
		}

		if input.DiskMBpsReadWrite != nil {
			output.DiskMBpsReadWrite = *input.DiskMBpsReadWrite
		}

		if input.DiskSizeGB != nil {
			output.DiskSizeGB = *input.DiskSizeGB
		}

		if input.Name != nil {
			output.Name = *input.Name
		}

		if input.WriteAcceleratorEnabled != nil {
			output.WriteAcceleratorEnabled = *input.WriteAcceleratorEnabled
		}
		outputList = append(outputList, output)
	}
	return outputList
}

func flattenVirtualMachineScaleSetManagedDiskParametersModel(input *fleets.VirtualMachineScaleSetManagedDiskParameters) []VirtualMachineScaleSetManagedDiskParametersModel {
	var outputList []VirtualMachineScaleSetManagedDiskParametersModel
	if input == nil {
		return outputList
	}
	output := VirtualMachineScaleSetManagedDiskParametersModel{
		DiskEncryptionSet: flattenDiskEncryptionSetParametersModel(input.DiskEncryptionSet),
		SecurityProfile:   flattenVMDiskSecurityProfileModel(input.SecurityProfile),
	}
	if input.StorageAccountType != nil {
		output.StorageAccountType = *input.StorageAccountType
	}

	return append(outputList, output)
}

func flattenDiskEncryptionSetParametersModel(input *fleets.DiskEncryptionSetParameters) []DiskEncryptionSetParametersModel {
	var outputList []DiskEncryptionSetParametersModel
	if input == nil {
		return outputList
	}
	output := DiskEncryptionSetParametersModel{}
	if input.Id != nil {
		output.Id = *input.Id
	}

	return append(outputList, output)
}

func flattenVMDiskSecurityProfileModel(input *fleets.VMDiskSecurityProfile) []VMDiskSecurityProfileModel {
	var outputList []VMDiskSecurityProfileModel
	if input == nil {
		return outputList
	}
	output := VMDiskSecurityProfileModel{
		DiskEncryptionSet: flattenDiskEncryptionSetParametersModel(input.DiskEncryptionSet),
	}
	if input.SecurityEncryptionType != nil {
		output.SecurityEncryptionType = *input.SecurityEncryptionType
	}

	return append(outputList, output)
}

func flattenImageReferenceModel(input *fleets.ImageReference) []ImageReferenceModel {
	var outputList []ImageReferenceModel
	if input == nil {
		return outputList
	}
	output := ImageReferenceModel{}
	if input.CommunityGalleryImageId != nil {
		output.CommunityGalleryImageId = *input.CommunityGalleryImageId
	}

	if input.ExactVersion != nil {
		output.ExactVersion = *input.ExactVersion
	}

	if input.Id != nil {
		output.Id = *input.Id
	}

	if input.Offer != nil {
		output.Offer = *input.Offer
	}

	if input.Publisher != nil {
		output.Publisher = *input.Publisher
	}

	if input.SharedGalleryImageId != nil {
		output.SharedGalleryImageId = *input.SharedGalleryImageId
	}

	if input.Sku != nil {
		output.Sku = *input.Sku
	}

	if input.Version != nil {
		output.Version = *input.Version
	}

	return append(outputList, output)
}

func flattenVirtualMachineScaleSetOSDiskModel(input *fleets.VirtualMachineScaleSetOSDisk) []VirtualMachineScaleSetOSDiskModel {
	var outputList []VirtualMachineScaleSetOSDiskModel
	if input == nil {
		return outputList
	}
	output := VirtualMachineScaleSetOSDiskModel{
		CreateOption:     input.CreateOption,
		DiffDiskSettings: flattenDiffDiskSettingsModel(input.DiffDiskSettings),
		Image:            flattenVirtualHardDiskModel(input.Image),
		ManagedDisk:      flattenVirtualMachineScaleSetManagedDiskParametersModel(input.ManagedDisk),
	}
	if input.Caching != nil {
		output.Caching = *input.Caching
	}

	if input.DeleteOption != nil {
		output.DeleteOption = *input.DeleteOption
	}

	if input.DiskSizeGB != nil {
		output.DiskSizeGB = *input.DiskSizeGB
	}

	if input.Name != nil {
		output.Name = *input.Name
	}

	if input.OsType != nil {
		output.OsType = *input.OsType
	}

	if input.VhdContainers != nil {
		output.VhdContainers = *input.VhdContainers
	}

	if input.WriteAcceleratorEnabled != nil {
		output.WriteAcceleratorEnabled = *input.WriteAcceleratorEnabled
	}

	return append(outputList, output)
}

func flattenDiffDiskSettingsModel(input *fleets.DiffDiskSettings) []DiffDiskSettingsModel {
	var outputList []DiffDiskSettingsModel
	if input == nil {
		return outputList
	}
	output := DiffDiskSettingsModel{}
	if input.Option != nil {
		output.Option = *input.Option
	}

	if input.Placement != nil {
		output.Placement = *input.Placement
	}

	return append(outputList, output)
}

func flattenVirtualHardDiskModel(input *fleets.VirtualHardDisk) []VirtualHardDiskModel {
	var outputList []VirtualHardDiskModel
	if input == nil {
		return outputList
	}
	output := VirtualHardDiskModel{}
	if input.Uri != nil {
		output.Uri = *input.Uri
	}

	return append(outputList, output)
}

func flattenComputeProfileModel(input *fleets.ComputeProfile) ([]ComputeProfileModel, error) {
	var outputList []ComputeProfileModel
	if input == nil {
		return outputList, nil
	}
	output := ComputeProfileModel{
		AdditionalVirtualMachineCapabilities: flattenAdditionalCapabilitiesModel(input.AdditionalVirtualMachineCapabilities),
	}
	baseVirtualMachineProfileValue, err := flattenBaseVirtualMachineProfileModel(&input.BaseVirtualMachineProfile)
	if err != nil {
		return nil, err
	}

	output.BaseVirtualMachineProfile = baseVirtualMachineProfileValue

	if input.ComputeApiVersion != nil {
		output.ComputeApiVersion = *input.ComputeApiVersion
	}

	if input.PlatformFaultDomainCount != nil {
		output.PlatformFaultDomainCount = *input.PlatformFaultDomainCount
	}

	return append(outputList, output), nil
}

func flattenAdditionalCapabilitiesModel(input *fleets.AdditionalCapabilities) []AdditionalCapabilitiesModel {
	var outputList []AdditionalCapabilitiesModel
	if input == nil {
		return outputList
	}
	output := AdditionalCapabilitiesModel{}
	if input.HibernationEnabled != nil {
		output.HibernationEnabled = *input.HibernationEnabled
	}

	if input.UltraSSDEnabled != nil {
		output.UltraSSDEnabled = *input.UltraSSDEnabled
	}

	return append(outputList, output)
}

func flattenBaseVirtualMachineProfileModel(input *fleets.BaseVirtualMachineProfile) ([]VirtualMachineProfileModel, error) {
	var outputList []VirtualMachineProfileModel
	if input == nil {
		return outputList, nil
	}
	output := VirtualMachineProfileModel{
		ApplicationProfile:       flattenApplicationProfileModel(input.ApplicationProfile),
		CapacityReservation:      flattenCapacityReservationProfileModel(input.CapacityReservation),
		DiagnosticsProfile:       flattenDiagnosticsProfileModel(input.DiagnosticsProfile),
		HardwareProfile:          flattenVirtualMachineScaleSetHardwareProfileModel(input.HardwareProfile),
		NetworkProfile:           flattenVirtualMachineScaleSetNetworkProfileModel(input.NetworkProfile),
		OsProfile:                flattenVirtualMachineScaleSetOSProfileModel(input.OsProfile),
		ScheduledEventsProfile:   flattenScheduledEventsProfileModel(input.ScheduledEventsProfile),
		SecurityPostureReference: flattenSecurityPostureReferenceModel(input.SecurityPostureReference),
		SecurityProfile:          flattenSecurityProfileModel(input.SecurityProfile),
		ServiceArtifactReference: flattenServiceArtifactReferenceModel(input.ServiceArtifactReference),
		StorageProfile:           flattenVirtualMachineScaleSetStorageProfileModel(input.StorageProfile),
	}
	extensionProfileValue, err := flattenVirtualMachineScaleSetExtensionProfileModel(input.ExtensionProfile)
	if err != nil {
		return nil, err
	}

	output.ExtensionProfile = extensionProfileValue

	if input.LicenseType != nil {
		output.LicenseType = *input.LicenseType
	}

	if input.TimeCreated != nil {
		output.TimeCreated = *input.TimeCreated
	}

	if input.UserData != nil {
		output.UserData = *input.UserData
	}

	return append(outputList, output), nil
}

func flattenApplicationProfileModel(input *fleets.ApplicationProfile) []ApplicationProfileModel {
	var outputList []ApplicationProfileModel
	if input == nil {
		return outputList
	}
	output := ApplicationProfileModel{
		GalleryApplications: flattenVMGalleryApplicationModelArray(input.GalleryApplications),
	}

	return append(outputList, output)
}

func flattenVMGalleryApplicationModelArray(inputList *[]fleets.VMGalleryApplication) []VMGalleryApplicationModel {
	var outputList []VMGalleryApplicationModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := VMGalleryApplicationModel{
			PackageReferenceId: input.PackageReferenceId,
		}

		if input.ConfigurationReference != nil {
			output.ConfigurationReference = *input.ConfigurationReference
		}

		if input.EnableAutomaticUpgrade != nil {
			output.EnableAutomaticUpgrade = *input.EnableAutomaticUpgrade
		}

		if input.Order != nil {
			output.Order = *input.Order
		}

		if input.Tags != nil {
			output.Tags = *input.Tags
		}

		if input.TreatFailureAsDeploymentFailure != nil {
			output.TreatFailureAsDeploymentFailure = *input.TreatFailureAsDeploymentFailure
		}
		outputList = append(outputList, output)
	}
	return outputList
}

func flattenCapacityReservationProfileModel(input *fleets.CapacityReservationProfile) []CapacityReservationProfileModel {
	var outputList []CapacityReservationProfileModel
	if input == nil {
		return outputList
	}
	output := CapacityReservationProfileModel{
		CapacityReservationGroup: flattenSubResourceModel(input.CapacityReservationGroup),
	}

	return append(outputList, output)
}

func flattenSubResourceModel(input *fleets.SubResource) []SubResourceModel {
	var outputList []SubResourceModel
	if input == nil {
		return outputList
	}
	output := SubResourceModel{}
	if input.Id != nil {
		output.Id = *input.Id
	}

	return append(outputList, output)
}

func flattenDiagnosticsProfileModel(input *fleets.DiagnosticsProfile) []DiagnosticsProfileModel {
	var outputList []DiagnosticsProfileModel
	if input == nil {
		return outputList
	}
	output := DiagnosticsProfileModel{
		BootDiagnostics: flattenBootDiagnosticsModel(input.BootDiagnostics),
	}

	return append(outputList, output)
}

func flattenBootDiagnosticsModel(input *fleets.BootDiagnostics) []BootDiagnosticsModel {
	var outputList []BootDiagnosticsModel
	if input == nil {
		return outputList
	}
	output := BootDiagnosticsModel{}
	if input.Enabled != nil {
		output.Enabled = *input.Enabled
	}

	if input.StorageUri != nil {
		output.StorageUri = *input.StorageUri
	}

	return append(outputList, output)
}

func flattenVirtualMachineScaleSetExtensionProfileModel(input *fleets.VirtualMachineScaleSetExtensionProfile) ([]VirtualMachineScaleSetExtensionProfileModel, error) {
	var outputList []VirtualMachineScaleSetExtensionProfileModel
	if input == nil {
		return outputList, nil
	}
	output := VirtualMachineScaleSetExtensionProfileModel{}
	extensionsValue, err := flattenVirtualMachineScaleSetExtensionModelArray(input.Extensions)
	if err != nil {
		return nil, err
	}

	output.Extensions = extensionsValue

	if input.ExtensionsTimeBudget != nil {
		output.ExtensionsTimeBudget = *input.ExtensionsTimeBudget
	}

	return append(outputList, output), nil
}

func flattenVirtualMachineScaleSetExtensionModelArray(inputList *[]fleets.VirtualMachineScaleSetExtension) ([]VirtualMachineScaleSetExtensionModel, error) {
	var outputList []VirtualMachineScaleSetExtensionModel
	if inputList == nil {
		return outputList, nil
	}
	for _, input := range *inputList {
		output := VirtualMachineScaleSetExtensionModel{}

		if input.Id != nil {
			output.Id = *input.Id
		}

		if input.Name != nil {
			output.Name = *input.Name
		}

		propertiesValue, err := flattenVirtualMachineScaleSetExtensionModel(input.Properties)
		if err != nil {
			return nil, err
		}

		output.Properties = propertiesValue

		if input.Type != nil {
			output.Type = *input.Type
		}
		outputList = append(outputList, output)
	}
	return outputList, nil
}

func flattenVirtualMachineScaleSetExtensionModel(input *fleets.VirtualMachineScaleSetExtensionProperties) ([]VirtualMachineScaleSetExtensionModel, error) {
	var outputList []VirtualMachineScaleSetExtensionModel
	if input == nil {
		return outputList, nil
	}
	output := VirtualMachineScaleSetExtensionModel{
		ProtectedSettingsFromKeyVault: flattenKeyVaultSecretReferenceModel(input.ProtectedSettingsFromKeyVault),
	}
	if input.AutoUpgradeMinorVersion != nil {
		output.AutoUpgradeMinorVersion = *input.AutoUpgradeMinorVersion
	}

	if input.EnableAutomaticUpgrade != nil {
		output.EnableAutomaticUpgrade = *input.EnableAutomaticUpgrade
	}

	if input.ForceUpdateTag != nil {
		output.ForceUpdateTag = *input.ForceUpdateTag
	}

	if input.ProtectedSettings != nil && *input.ProtectedSettings != nil {

		protectedSettingsValue, err := json.Marshal(*input.ProtectedSettings)
		if err != nil {
			return outputList, err
		}

		output.ProtectedSettings = string(protectedSettingsValue)
	}

	if input.ProvisionAfterExtensions != nil {
		output.ProvisionAfterExtensions = *input.ProvisionAfterExtensions
	}

	if input.ProvisioningState != nil {
		output.ProvisioningState = *input.ProvisioningState
	}

	if input.Publisher != nil {
		output.Publisher = *input.Publisher
	}

	if input.Settings != nil && *input.Settings != nil {

		settingsValue, err := json.Marshal(*input.Settings)
		if err != nil {
			return outputList, err
		}

		output.Settings = string(settingsValue)
	}

	if input.SuppressFailures != nil {
		output.SuppressFailures = *input.SuppressFailures
	}

	if input.Type != nil {
		output.Type = *input.Type
	}

	if input.TypeHandlerVersion != nil {
		output.TypeHandlerVersion = *input.TypeHandlerVersion
	}

	return append(outputList, output), nil
}

func flattenKeyVaultSecretReferenceModel(input *fleets.KeyVaultSecretReference) []KeyVaultSecretReferenceModel {
	var outputList []KeyVaultSecretReferenceModel
	if input == nil {
		return outputList
	}
	output := KeyVaultSecretReferenceModel{
		SecretUrl:   input.SecretUrl,
		SourceVault: flattenSubResourceModel(&input.SourceVault),
	}

	return append(outputList, output)
}

func flattenVirtualMachineScaleSetHardwareProfileModel(input *fleets.VirtualMachineScaleSetHardwareProfile) []VirtualMachineScaleSetHardwareProfileModel {
	var outputList []VirtualMachineScaleSetHardwareProfileModel
	if input == nil {
		return outputList
	}
	output := VirtualMachineScaleSetHardwareProfileModel{
		VMSizeProperties: flattenVMSizePropertiesModel(input.VMSizeProperties),
	}

	return append(outputList, output)
}

func flattenVMSizePropertiesModel(input *fleets.VMSizeProperties) []VMSizePropertiesModel {
	var outputList []VMSizePropertiesModel
	if input == nil {
		return outputList
	}
	output := VMSizePropertiesModel{}
	if input.VCPUsAvailable != nil {
		output.VCPUsAvailable = *input.VCPUsAvailable
	}

	if input.VCPUsPerCore != nil {
		output.VCPUsPerCore = *input.VCPUsPerCore
	}

	return append(outputList, output)
}

func flattenVirtualMachineScaleSetNetworkProfileModel(input *fleets.VirtualMachineScaleSetNetworkProfile) []NetworkProfileModel {
	var outputList []NetworkProfileModel
	if input == nil {
		return outputList
	}
	output := NetworkProfileModel{
		HealthProbe:                    flattenApiEntityReferenceModel(input.HealthProbe),
		NetworkInterfaceConfigurations: flattenVirtualMachineScaleSetNetworkConfigurationModelArray(input.NetworkInterfaceConfigurations),
	}
	if input.NetworkApiVersion != nil {
		output.NetworkApiVersion = *input.NetworkApiVersion
	}

	return append(outputList, output)
}

func flattenApiEntityReferenceModel(input *fleets.ApiEntityReference) []ApiEntityReferenceModel {
	var outputList []ApiEntityReferenceModel
	if input == nil {
		return outputList
	}
	output := ApiEntityReferenceModel{}
	if input.Id != nil {
		output.Id = *input.Id
	}

	return append(outputList, output)
}

func flattenVirtualMachineScaleSetNetworkConfigurationModelArray(inputList *[]fleets.VirtualMachineScaleSetNetworkConfiguration) []NetworkConfigurationModel {
	var outputList []NetworkConfigurationModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := NetworkConfigurationModel{
			Name:       input.Name,
			Properties: flattenVirtualMachineScaleSetNetworkConfigurationPropertiesModel(input.Properties),
		}

		outputList = append(outputList, output)
	}
	return outputList
}

func flattenVirtualMachineScaleSetNetworkConfigurationPropertiesModel(input *fleets.VirtualMachineScaleSetNetworkConfigurationProperties) []VirtualMachineScaleSetNetworkConfigurationPropertiesModel {
	var outputList []VirtualMachineScaleSetNetworkConfigurationPropertiesModel
	if input == nil {
		return outputList
	}
	output := VirtualMachineScaleSetNetworkConfigurationPropertiesModel{
		DnsSettings:          flattenVirtualMachineScaleSetNetworkConfigurationDnsSettingsModel(input.DnsSettings),
		IPConfigurations:     flattenVirtualMachineScaleSetIPConfigurationModelArray(&input.IPConfigurations),
		NetworkSecurityGroup: flattenSubResourceModel(input.NetworkSecurityGroup),
	}
	if input.AuxiliaryMode != nil {
		output.AuxiliaryMode = *input.AuxiliaryMode
	}

	if input.AuxiliarySku != nil {
		output.AuxiliarySku = *input.AuxiliarySku
	}

	if input.DeleteOption != nil {
		output.DeleteOption = *input.DeleteOption
	}

	if input.DisableTcpStateTracking != nil {
		output.DisableTcpStateTracking = *input.DisableTcpStateTracking
	}

	if input.EnableAcceleratedNetworking != nil {
		output.EnableAcceleratedNetworking = *input.EnableAcceleratedNetworking
	}

	if input.EnableFpga != nil {
		output.EnableFpga = *input.EnableFpga
	}

	if input.EnableIPForwarding != nil {
		output.EnableIPForwarding = *input.EnableIPForwarding
	}

	if input.Primary != nil {
		output.Primary = *input.Primary
	}

	return append(outputList, output)
}

func flattenVirtualMachineScaleSetNetworkConfigurationDnsSettingsModel(input *fleets.VirtualMachineScaleSetNetworkConfigurationDnsSettings) []VirtualMachineScaleSetNetworkConfigurationDnsSettingsModel {
	var outputList []VirtualMachineScaleSetNetworkConfigurationDnsSettingsModel
	if input == nil {
		return outputList
	}
	output := VirtualMachineScaleSetNetworkConfigurationDnsSettingsModel{}
	if input.DnsServers != nil {
		output.DnsServers = *input.DnsServers
	}

	return append(outputList, output)
}

func flattenVirtualMachineScaleSetIPConfigurationModelArray(inputList *[]fleets.VirtualMachineScaleSetIPConfiguration) []VirtualMachineScaleSetIPConfigurationModel {
	var outputList []VirtualMachineScaleSetIPConfigurationModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := VirtualMachineScaleSetIPConfigurationModel{
			Name:       input.Name,
			Properties: flattenVirtualMachineScaleSetIPConfigurationModel(input.Properties),
		}

		outputList = append(outputList, output)
	}
	return outputList
}

func flattenVirtualMachineScaleSetIPConfigurationModel(input *fleets.VirtualMachineScaleSetIPConfigurationProperties) []VirtualMachineScaleSetIPConfigurationModel {
	var outputList []VirtualMachineScaleSetIPConfigurationModel
	if input == nil {
		return outputList
	}
	output := VirtualMachineScaleSetIPConfigurationModel{
		ApplicationGatewayBackendAddressPools: flattenSubResourceModelArray(input.ApplicationGatewayBackendAddressPools),
		ApplicationSecurityGroups:             flattenSubResourceModelArray(input.ApplicationSecurityGroups),
		LoadBalancerBackendAddressPools:       flattenSubResourceModelArray(input.LoadBalancerBackendAddressPools),
		LoadBalancerInboundNatPools:           flattenSubResourceModelArray(input.LoadBalancerInboundNatPools),
		PublicIPAddressConfiguration:          flattenVirtualMachineScaleSetPublicIPAddressConfigurationModel(input.PublicIPAddressConfiguration),
		Subnet:                                flattenApiEntityReferenceModel(input.Subnet),
	}
	if input.Primary != nil {
		output.Primary = *input.Primary
	}

	if input.PrivateIPAddressVersion != nil {
		output.PrivateIPAddressVersion = *input.PrivateIPAddressVersion
	}

	return append(outputList, output)
}

func flattenVirtualMachineScaleSetPublicIPAddressConfigurationModel(input *fleets.VirtualMachineScaleSetPublicIPAddressConfiguration) []VirtualMachineScaleSetPublicIPAddressConfigurationModel {
	var outputList []VirtualMachineScaleSetPublicIPAddressConfigurationModel
	if input == nil {
		return outputList
	}
	output := VirtualMachineScaleSetPublicIPAddressConfigurationModel{
		Name:       input.Name,
		Properties: flattenVirtualMachineScaleSetPublicIPAddressConfigurationModel(input.Properties),
		Sku:        flattenPublicIPAddressSkuModel(input.Sku),
	}

	return append(outputList, output)
}

func flattenVirtualMachineScaleSetPublicIPAddressConfigurationModel(input *fleets.VirtualMachineScaleSetPublicIPAddressConfigurationProperties) []VirtualMachineScaleSetPublicIPAddressConfigurationModel {
	var outputList []VirtualMachineScaleSetPublicIPAddressConfigurationModel
	if input == nil {
		return outputList
	}
	output := VirtualMachineScaleSetPublicIPAddressConfigurationModel{
		DnsSettings:    flattenVirtualMachineScaleSetPublicIPAddressConfigurationDnsSettingsModel(input.DnsSettings),
		IPTags:         flattenVirtualMachineScaleSetIPTagModelArray(input.IPTags),
		PublicIPPrefix: flattenSubResourceModel(input.PublicIPPrefix),
	}
	if input.DeleteOption != nil {
		output.DeleteOption = *input.DeleteOption
	}

	if input.IdleTimeoutInMinutes != nil {
		output.IdleTimeoutInMinutes = *input.IdleTimeoutInMinutes
	}

	if input.PublicIPAddressVersion != nil {
		output.PublicIPAddressVersion = *input.PublicIPAddressVersion
	}

	return append(outputList, output)
}

func flattenVirtualMachineScaleSetPublicIPAddressConfigurationDnsSettingsModel(input *fleets.VirtualMachineScaleSetPublicIPAddressConfigurationDnsSettings) []VirtualMachineScaleSetPublicIPAddressConfigurationDnsSettingsModel {
	var outputList []VirtualMachineScaleSetPublicIPAddressConfigurationDnsSettingsModel
	if input == nil {
		return outputList
	}
	output := VirtualMachineScaleSetPublicIPAddressConfigurationDnsSettingsModel{
		DomainNameLabel: input.DomainNameLabel,
	}
	if input.DomainNameLabelScope != nil {
		output.DomainNameLabelScope = *input.DomainNameLabelScope
	}

	return append(outputList, output)
}

func flattenVirtualMachineScaleSetIPTagModelArray(inputList *[]fleets.VirtualMachineScaleSetIPTag) []VirtualMachineScaleSetIPTagModel {
	var outputList []VirtualMachineScaleSetIPTagModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := VirtualMachineScaleSetIPTagModel{}

		if input.IPTagType != nil {
			output.IPTagType = *input.IPTagType
		}

		if input.Tag != nil {
			output.Tag = *input.Tag
		}
		outputList = append(outputList, output)
	}
	return outputList
}

func flattenPublicIPAddressSkuModel(input *fleets.PublicIPAddressSku) []PublicIPAddressSkuModel {
	var outputList []PublicIPAddressSkuModel
	if input == nil {
		return outputList
	}
	output := PublicIPAddressSkuModel{}
	if input.Name != nil {
		output.Name = *input.Name
	}

	if input.Tier != nil {
		output.Tier = *input.Tier
	}

	return append(outputList, output)
}

func flattenVirtualMachineScaleSetOSProfileModel(input *fleets.VirtualMachineScaleSetOSProfile) []OSProfileModel {
	var outputList []OSProfileModel
	if input == nil {
		return outputList
	}
	output := OSProfileModel{
		LinuxConfiguration:   flattenLinuxConfigurationModel(input.LinuxConfiguration),
		Secrets:              flattenVaultSecretGroupModelArray(input.Secrets),
		WindowsConfiguration: flattenWindowsConfigurationModel(input.WindowsConfiguration),
	}
	if input.AdminPassword != nil {
		output.AdminPassword = *input.AdminPassword
	}

	if input.AdminUsername != nil {
		output.AdminUsername = *input.AdminUsername
	}

	if input.AllowExtensionOperations != nil {
		output.AllowExtensionOperations = *input.AllowExtensionOperations
	}

	if input.ComputerNamePrefix != nil {
		output.ComputerNamePrefix = *input.ComputerNamePrefix
	}

	if input.CustomData != nil {
		output.CustomData = *input.CustomData
	}

	if input.RequireGuestProvisionSignal != nil {
		output.RequireGuestProvisionSignal = *input.RequireGuestProvisionSignal
	}

	return append(outputList, output)
}

func flattenLinuxConfigurationModel(input *fleets.LinuxConfiguration) []LinuxConfigurationModel {
	var outputList []LinuxConfigurationModel
	if input == nil {
		return outputList
	}
	output := LinuxConfigurationModel{
		PatchSettings: flattenLinuxPatchSettingsModel(input.PatchSettings),
		Ssh:           flattenSshConfigurationModel(input.Ssh),
	}
	if input.DisablePasswordAuthentication != nil {
		output.DisablePasswordAuthentication = *input.DisablePasswordAuthentication
	}

	if input.EnableVMAgentPlatformUpdates != nil {
		output.EnableVMAgentPlatformUpdates = *input.EnableVMAgentPlatformUpdates
	}

	if input.ProvisionVMAgent != nil {
		output.ProvisionVMAgent = *input.ProvisionVMAgent
	}

	return append(outputList, output)
}

func flattenLinuxPatchSettingsModel(input *fleets.LinuxPatchSettings) []LinuxPatchSettingsModel {
	var outputList []LinuxPatchSettingsModel
	if input == nil {
		return outputList
	}
	output := LinuxPatchSettingsModel{
		AutomaticByPlatformSettings: flattenLinuxVMGuestPatchAutomaticByPlatformSettingsModel(input.AutomaticByPlatformSettings),
	}
	if input.AssessmentMode != nil {
		output.AssessmentMode = *input.AssessmentMode
	}

	if input.PatchMode != nil {
		output.PatchMode = *input.PatchMode
	}

	return append(outputList, output)
}

func flattenLinuxVMGuestPatchAutomaticByPlatformSettingsModel(input *fleets.LinuxVMGuestPatchAutomaticByPlatformSettings) []LinuxVMGuestPatchAutomaticByPlatformSettingsModel {
	var outputList []LinuxVMGuestPatchAutomaticByPlatformSettingsModel
	if input == nil {
		return outputList
	}
	output := LinuxVMGuestPatchAutomaticByPlatformSettingsModel{}
	if input.BypassPlatformSafetyChecksOnUserSchedule != nil {
		output.BypassPlatformSafetyChecksOnUserSchedule = *input.BypassPlatformSafetyChecksOnUserSchedule
	}

	if input.RebootSetting != nil {
		output.RebootSetting = *input.RebootSetting
	}

	return append(outputList, output)
}

func flattenSshConfigurationModel(input *fleets.SshConfiguration) []SshConfigurationModel {
	var outputList []SshConfigurationModel
	if input == nil {
		return outputList
	}
	output := SshConfigurationModel{
		PublicKeys: flattenSshPublicKeyModelArray(input.PublicKeys),
	}

	return append(outputList, output)
}

func flattenSshPublicKeyModelArray(inputList *[]fleets.SshPublicKey) []SshPublicKeyModel {
	var outputList []SshPublicKeyModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := SshPublicKeyModel{}

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

func flattenVaultSecretGroupModelArray(inputList *[]fleets.VaultSecretGroup) []VaultSecretGroupModel {
	var outputList []VaultSecretGroupModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := VaultSecretGroupModel{
			SourceVault:       flattenSubResourceModel(input.SourceVault),
			VaultCertificates: flattenVaultCertificateModelArray(input.VaultCertificates),
		}

		outputList = append(outputList, output)
	}
	return outputList
}

func flattenVaultCertificateModelArray(inputList *[]fleets.VaultCertificate) []VaultCertificateModel {
	var outputList []VaultCertificateModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := VaultCertificateModel{}

		if input.CertificateStore != nil {
			output.CertificateStore = *input.CertificateStore
		}

		if input.CertificateUrl != nil {
			output.CertificateUrl = *input.CertificateUrl
		}
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
		AdditionalUnattendContent: flattenAdditionalUnattendContentModelArray(input.AdditionalUnattendContent),
		PatchSettings:             flattenPatchSettingsModel(input.PatchSettings),
		WinRM:                     flattenWinRMConfigurationModel(input.WinRM),
	}
	if input.EnableAutomaticUpdates != nil {
		output.EnableAutomaticUpdates = *input.EnableAutomaticUpdates
	}

	if input.EnableVMAgentPlatformUpdates != nil {
		output.EnableVMAgentPlatformUpdates = *input.EnableVMAgentPlatformUpdates
	}

	if input.ProvisionVMAgent != nil {
		output.ProvisionVMAgent = *input.ProvisionVMAgent
	}

	if input.TimeZone != nil {
		output.TimeZone = *input.TimeZone
	}

	return append(outputList, output)
}

func flattenAdditionalUnattendContentModelArray(inputList *[]fleets.AdditionalUnattendContent) []AdditionalUnattendContentModel {
	var outputList []AdditionalUnattendContentModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := AdditionalUnattendContentModel{}

		if input.ComponentName != nil {
			output.ComponentName = *input.ComponentName
		}

		if input.Content != nil {
			output.Content = *input.Content
		}

		if input.PassName != nil {
			output.PassName = *input.PassName
		}

		if input.SettingName != nil {
			output.SettingName = *input.SettingName
		}
		outputList = append(outputList, output)
	}
	return outputList
}

func flattenPatchSettingsModel(input *fleets.PatchSettings) []PatchSettingsModel {
	var outputList []PatchSettingsModel
	if input == nil {
		return outputList
	}
	output := PatchSettingsModel{
		AutomaticByPlatformSettings: flattenWindowsVMGuestPatchAutomaticByPlatformSettingsModel(input.AutomaticByPlatformSettings),
	}
	if input.AssessmentMode != nil {
		output.AssessmentMode = *input.AssessmentMode
	}

	if input.EnableHotpatching != nil {
		output.EnableHotpatching = *input.EnableHotpatching
	}

	if input.PatchMode != nil {
		output.PatchMode = *input.PatchMode
	}

	return append(outputList, output)
}

func flattenWindowsVMGuestPatchAutomaticByPlatformSettingsModel(input *fleets.WindowsVMGuestPatchAutomaticByPlatformSettings) []WindowsVMGuestPatchAutomaticByPlatformSettingsModel {
	var outputList []WindowsVMGuestPatchAutomaticByPlatformSettingsModel
	if input == nil {
		return outputList
	}
	output := WindowsVMGuestPatchAutomaticByPlatformSettingsModel{}
	if input.BypassPlatformSafetyChecksOnUserSchedule != nil {
		output.BypassPlatformSafetyChecksOnUserSchedule = *input.BypassPlatformSafetyChecksOnUserSchedule
	}

	if input.RebootSetting != nil {
		output.RebootSetting = *input.RebootSetting
	}

	return append(outputList, output)
}

func flattenWinRMConfigurationModel(input *fleets.WinRMConfiguration) []WinRMConfigurationModel {
	var outputList []WinRMConfigurationModel
	if input == nil {
		return outputList
	}
	output := WinRMConfigurationModel{
		Listeners: flattenWinRMListenerModelArray(input.Listeners),
	}

	return append(outputList, output)
}

func flattenWinRMListenerModelArray(inputList *[]fleets.WinRMListener) []WinRMListenerModel {
	var outputList []WinRMListenerModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := WinRMListenerModel{}

		if input.CertificateUrl != nil {
			output.CertificateUrl = *input.CertificateUrl
		}

		if input.Protocol != nil {
			output.Protocol = *input.Protocol
		}
		outputList = append(outputList, output)
	}
	return outputList
}

func flattenScheduledEventsProfileModel(input *fleets.ScheduledEventsProfile) []ScheduledEventsProfileModel {
	var outputList []ScheduledEventsProfileModel
	if input == nil {
		return outputList
	}
	output := ScheduledEventsProfileModel{
		OsImageNotificationProfile:   flattenOSImageNotificationProfileModel(input.OsImageNotificationProfile),
		TerminateNotificationProfile: flattenTerminateNotificationProfileModel(input.TerminateNotificationProfile),
	}

	return append(outputList, output)
}

func flattenOSImageNotificationProfileModel(input *fleets.OSImageNotificationProfile) []OSImageNotificationProfileModel {
	var outputList []OSImageNotificationProfileModel
	if input == nil {
		return outputList
	}
	output := OSImageNotificationProfileModel{}
	if input.Enable != nil {
		output.Enable = *input.Enable
	}

	if input.NotBeforeTimeout != nil {
		output.NotBeforeTimeout = *input.NotBeforeTimeout
	}

	return append(outputList, output)
}

func flattenTerminateNotificationProfileModel(input *fleets.TerminateNotificationProfile) []TerminateNotificationProfileModel {
	var outputList []TerminateNotificationProfileModel
	if input == nil {
		return outputList
	}
	output := TerminateNotificationProfileModel{}
	if input.Enable != nil {
		output.Enable = *input.Enable
	}

	if input.NotBeforeTimeout != nil {
		output.NotBeforeTimeout = *input.NotBeforeTimeout
	}

	return append(outputList, output)
}

func flattenSecurityPostureReferenceModel(input *fleets.SecurityPostureReference) []SecurityPostureReferenceModel {
	var outputList []SecurityPostureReferenceModel
	if input == nil {
		return outputList
	}
	output := SecurityPostureReferenceModel{}
	if input.ExcludeExtensions != nil {
		output.ExcludeExtensions = *input.ExcludeExtensions
	}

	if input.Id != nil {
		output.Id = *input.Id
	}

	if input.IsOverridable != nil {
		output.IsOverridable = *input.IsOverridable
	}

	return append(outputList, output)
}

func flattenSecurityProfileModel(input *fleets.SecurityProfile) []SecurityProfileModel {
	var outputList []SecurityProfileModel
	if input == nil {
		return outputList
	}
	output := SecurityProfileModel{
		EncryptionIdentity: flattenEncryptionIdentityModel(input.EncryptionIdentity),
		ProxyAgentSettings: flattenProxyAgentSettingsModel(input.ProxyAgentSettings),
		UefiSettings:       flattenUefiSettingsModel(input.UefiSettings),
	}
	if input.EncryptionAtHost != nil {
		output.EncryptionAtHost = *input.EncryptionAtHost
	}

	if input.SecurityType != nil {
		output.SecurityType = *input.SecurityType
	}

	return append(outputList, output)
}

func flattenEncryptionIdentityModel(input *fleets.EncryptionIdentity) []EncryptionIdentityModel {
	var outputList []EncryptionIdentityModel
	if input == nil {
		return outputList
	}
	output := EncryptionIdentityModel{}
	if input.UserAssignedIdentityResourceId != nil {
		output.UserAssignedIdentityResourceId = *input.UserAssignedIdentityResourceId
	}

	return append(outputList, output)
}

func flattenProxyAgentSettingsModel(input *fleets.ProxyAgentSettings) []ProxyAgentSettingsModel {
	var outputList []ProxyAgentSettingsModel
	if input == nil {
		return outputList
	}
	output := ProxyAgentSettingsModel{}
	if input.Enabled != nil {
		output.Enabled = *input.Enabled
	}

	if input.KeyIncarnationId != nil {
		output.KeyIncarnationId = *input.KeyIncarnationId
	}

	if input.Mode != nil {
		output.Mode = *input.Mode
	}

	return append(outputList, output)
}

func flattenUefiSettingsModel(input *fleets.UefiSettings) []UefiSettingsModel {
	var outputList []UefiSettingsModel
	if input == nil {
		return outputList
	}
	output := UefiSettingsModel{}
	if input.SecureBootEnabled != nil {
		output.SecureBootEnabled = *input.SecureBootEnabled
	}

	if input.VTpmEnabled != nil {
		output.VTpmEnabled = *input.VTpmEnabled
	}

	return append(outputList, output)
}

func flattenServiceArtifactReferenceModel(input *fleets.ServiceArtifactReference) []ServiceArtifactReferenceModel {
	var outputList []ServiceArtifactReferenceModel
	if input == nil {
		return outputList
	}
	output := ServiceArtifactReferenceModel{}
	if input.Id != nil {
		output.Id = *input.Id
	}

	return append(outputList, output)
}

func flattenVirtualMachineScaleSetStorageProfileModel(input *fleets.VirtualMachineScaleSetStorageProfile) []StorageProfileModel {
	var outputList []StorageProfileModel
	if input == nil {
		return outputList
	}
	output := StorageProfileModel{
		DataDisks:      flattenVirtualMachineScaleSetDataDiskModelArray(input.DataDisks),
		ImageReference: flattenImageReferenceModel(input.ImageReference),
		OsDisk:         flattenVirtualMachineScaleSetOSDiskModel(input.OsDisk),
	}
	if input.DiskControllerType != nil {
		output.DiskControllerType = *input.DiskControllerType
	}

	return append(outputList, output)
}

func flattenVirtualMachineScaleSetDataDiskModelArray(inputList *[]fleets.VirtualMachineScaleSetDataDisk) []VirtualMachineScaleSetDataDiskModel {
	var outputList []VirtualMachineScaleSetDataDiskModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := VirtualMachineScaleSetDataDiskModel{
			CreateOption: input.CreateOption,
			Lun:          input.Lun,
			ManagedDisk:  flattenVirtualMachineScaleSetManagedDiskParametersModel(input.ManagedDisk),
		}

		if input.Caching != nil {
			output.Caching = *input.Caching
		}

		if input.DeleteOption != nil {
			output.DeleteOption = *input.DeleteOption
		}

		if input.DiskIOPSReadWrite != nil {
			output.DiskIOPSReadWrite = *input.DiskIOPSReadWrite
		}

		if input.DiskMBpsReadWrite != nil {
			output.DiskMBpsReadWrite = *input.DiskMBpsReadWrite
		}

		if input.DiskSizeGB != nil {
			output.DiskSizeGB = *input.DiskSizeGB
		}

		if input.Name != nil {
			output.Name = *input.Name
		}

		if input.WriteAcceleratorEnabled != nil {
			output.WriteAcceleratorEnabled = *input.WriteAcceleratorEnabled
		}
		outputList = append(outputList, output)
	}
	return outputList
}

func flattenVirtualMachineScaleSetManagedDiskParametersModel(input *fleets.VirtualMachineScaleSetManagedDiskParameters) []VirtualMachineScaleSetManagedDiskParametersModel {
	var outputList []VirtualMachineScaleSetManagedDiskParametersModel
	if input == nil {
		return outputList
	}
	output := VirtualMachineScaleSetManagedDiskParametersModel{
		DiskEncryptionSet: flattenDiskEncryptionSetParametersModel(input.DiskEncryptionSet),
		SecurityProfile:   flattenVMDiskSecurityProfileModel(input.SecurityProfile),
	}
	if input.StorageAccountType != nil {
		output.StorageAccountType = *input.StorageAccountType
	}

	return append(outputList, output)
}

func flattenDiskEncryptionSetParametersModel(input *fleets.DiskEncryptionSetParameters) []DiskEncryptionSetParametersModel {
	var outputList []DiskEncryptionSetParametersModel
	if input == nil {
		return outputList
	}
	output := DiskEncryptionSetParametersModel{}
	if input.Id != nil {
		output.Id = *input.Id
	}

	return append(outputList, output)
}

func flattenVMDiskSecurityProfileModel(input *fleets.VMDiskSecurityProfile) []VMDiskSecurityProfileModel {
	var outputList []VMDiskSecurityProfileModel
	if input == nil {
		return outputList
	}
	output := VMDiskSecurityProfileModel{
		DiskEncryptionSet: flattenDiskEncryptionSetParametersModel(input.DiskEncryptionSet),
	}
	if input.SecurityEncryptionType != nil {
		output.SecurityEncryptionType = *input.SecurityEncryptionType
	}

	return append(outputList, output)
}

func flattenImageReferenceModel(input *fleets.ImageReference) []ImageReferenceModel {
	var outputList []ImageReferenceModel
	if input == nil {
		return outputList
	}
	output := ImageReferenceModel{}
	if input.CommunityGalleryImageId != nil {
		output.CommunityGalleryImageId = *input.CommunityGalleryImageId
	}

	if input.ExactVersion != nil {
		output.ExactVersion = *input.ExactVersion
	}

	if input.Id != nil {
		output.Id = *input.Id
	}

	if input.Offer != nil {
		output.Offer = *input.Offer
	}

	if input.Publisher != nil {
		output.Publisher = *input.Publisher
	}

	if input.SharedGalleryImageId != nil {
		output.SharedGalleryImageId = *input.SharedGalleryImageId
	}

	if input.Sku != nil {
		output.Sku = *input.Sku
	}

	if input.Version != nil {
		output.Version = *input.Version
	}

	return append(outputList, output)
}

func flattenVirtualMachineScaleSetOSDiskModel(input *fleets.VirtualMachineScaleSetOSDisk) []VirtualMachineScaleSetOSDiskModel {
	var outputList []VirtualMachineScaleSetOSDiskModel
	if input == nil {
		return outputList
	}
	output := VirtualMachineScaleSetOSDiskModel{
		CreateOption:     input.CreateOption,
		DiffDiskSettings: flattenDiffDiskSettingsModel(input.DiffDiskSettings),
		Image:            flattenVirtualHardDiskModel(input.Image),
		ManagedDisk:      flattenVirtualMachineScaleSetManagedDiskParametersModel(input.ManagedDisk),
	}
	if input.Caching != nil {
		output.Caching = *input.Caching
	}

	if input.DeleteOption != nil {
		output.DeleteOption = *input.DeleteOption
	}

	if input.DiskSizeGB != nil {
		output.DiskSizeGB = *input.DiskSizeGB
	}

	if input.Name != nil {
		output.Name = *input.Name
	}

	if input.OsType != nil {
		output.OsType = *input.OsType
	}

	if input.VhdContainers != nil {
		output.VhdContainers = *input.VhdContainers
	}

	if input.WriteAcceleratorEnabled != nil {
		output.WriteAcceleratorEnabled = *input.WriteAcceleratorEnabled
	}

	return append(outputList, output)
}

func flattenDiffDiskSettingsModel(input *fleets.DiffDiskSettings) []DiffDiskSettingsModel {
	var outputList []DiffDiskSettingsModel
	if input == nil {
		return outputList
	}
	output := DiffDiskSettingsModel{}
	if input.Option != nil {
		output.Option = *input.Option
	}

	if input.Placement != nil {
		output.Placement = *input.Placement
	}

	return append(outputList, output)
}

func flattenVirtualHardDiskModel(input *fleets.VirtualHardDisk) []VirtualHardDiskModel {
	var outputList []VirtualHardDiskModel
	if input == nil {
		return outputList
	}
	output := VirtualHardDiskModel{}
	if input.Uri != nil {
		output.Uri = *input.Uri
	}

	return append(outputList, output)
}

func flattenPlanModel(input *fleets.Plan) []PlanModel {
	var outputList []PlanModel
	if input == nil {
		return outputList
	}
	output := PlanModel{
		Name:      input.Name,
		Product:   input.Product,
		Publisher: input.Publisher,
	}
	if input.PromotionCode != nil {
		output.PromotionCode = *input.PromotionCode
	}

	if input.Version != nil {
		output.Version = *input.Version
	}

	return append(outputList, output)
}

func flattenRegularPriorityProfileModel(input *fleets.RegularPriorityProfile) []RegularPriorityProfileModel {
	var outputList []RegularPriorityProfileModel
	if input == nil {
		return outputList
	}
	output := RegularPriorityProfileModel{}
	if input.AllocationStrategy != nil {
		output.AllocationStrategy = *input.AllocationStrategy
	}

	if input.Capacity != nil {
		output.Capacity = *input.Capacity
	}

	if input.MinCapacity != nil {
		output.MinCapacity = *input.MinCapacity
	}

	return append(outputList, output)
}

func flattenSpotPriorityProfileModel(input *fleets.SpotPriorityProfile) []SpotPriorityProfileModel {
	var outputList []SpotPriorityProfileModel
	if input == nil {
		return outputList
	}
	output := SpotPriorityProfileModel{}
	if input.AllocationStrategy != nil {
		output.AllocationStrategy = *input.AllocationStrategy
	}

	if input.Capacity != nil {
		output.Capacity = *input.Capacity
	}

	if input.EvictionPolicy != nil {
		output.EvictionPolicy = *input.EvictionPolicy
	}

	if input.Maintain != nil {
		output.Maintain = *input.Maintain
	}

	if input.MaxPricePerVM != nil {
		output.MaxPricePerVM = *input.MaxPricePerVM
	}

	if input.MinCapacity != nil {
		output.MinCapacity = *input.MinCapacity
	}

	return append(outputList, output)
}

func flattenVMAttributesModel(input *fleets.VMAttributes) []VMAttributesModel {
	var outputList []VMAttributesModel
	if input == nil {
		return outputList
	}
	output := VMAttributesModel{
		AcceleratorCount:          flattenVMAttributeMinMaxIntegerModel(input.AcceleratorCount),
		AcceleratorManufacturers:  flattenAcceleratorManufacturerModelArray(input.AcceleratorManufacturers),
		AcceleratorTypes:          flattenAcceleratorTypeModelArray(input.AcceleratorTypes),
		ArchitectureTypes:         flattenArchitectureTypeModelArray(input.ArchitectureTypes),
		CpuManufacturers:          flattenCPUManufacturerModelArray(input.CpuManufacturers),
		DataDiskCount:             flattenVMAttributeMinMaxIntegerModel(input.DataDiskCount),
		LocalStorageDiskTypes:     flattenLocalStorageDiskTypeModelArray(input.LocalStorageDiskTypes),
		LocalStorageInGiB:         flattenVMAttributeMinMaxDoubleModel(input.LocalStorageInGiB),
		MemoryInGiB:               flattenVMAttributeMinMaxDoubleModel(&input.MemoryInGiB),
		MemoryInGiBPerVCPU:        flattenVMAttributeMinMaxDoubleModel(input.MemoryInGiBPerVCPU),
		NetworkBandwidthInMbps:    flattenVMAttributeMinMaxDoubleModel(input.NetworkBandwidthInMbps),
		NetworkInterfaceCount:     flattenVMAttributeMinMaxIntegerModel(input.NetworkInterfaceCount),
		RdmaNetworkInterfaceCount: flattenVMAttributeMinMaxIntegerModel(input.RdmaNetworkInterfaceCount),
		VCPUCount:                 flattenVMAttributeMinMaxIntegerModel(&input.VCPUCount),
		VMCategories:              flattenVMCategoryModelArray(input.VMCategories),
	}
	if input.AcceleratorSupport != nil {
		output.AcceleratorSupport = *input.AcceleratorSupport
	}

	if input.BurstableSupport != nil {
		output.BurstableSupport = *input.BurstableSupport
	}

	if input.ExcludedVMSizes != nil {
		output.ExcludedVMSizes = *input.ExcludedVMSizes
	}

	if input.LocalStorageSupport != nil {
		output.LocalStorageSupport = *input.LocalStorageSupport
	}

	if input.RdmaSupport != nil {
		output.RdmaSupport = *input.RdmaSupport
	}

	return append(outputList, output)
}

func flattenVMAttributeMinMaxIntegerModel(input *fleets.VMAttributeMinMaxInteger) []VMAttributeMinMaxIntegerModel {
	var outputList []VMAttributeMinMaxIntegerModel
	if input == nil {
		return outputList
	}
	output := VMAttributeMinMaxIntegerModel{}
	if input.Max != nil {
		output.Max = *input.Max
	}

	if input.Min != nil {
		output.Min = *input.Min
	}

	return append(outputList, output)
}

func flattenAcceleratorManufacturerModelArray(inputList *[]fleets.AcceleratorManufacturer) []AcceleratorManufacturerModel {
	var outputList []AcceleratorManufacturerModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := AcceleratorManufacturerModel{}

		outputList = append(outputList, output)
	}
	return outputList
}

func flattenAcceleratorTypeModelArray(inputList *[]fleets.AcceleratorType) []AcceleratorTypeModel {
	var outputList []AcceleratorTypeModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := AcceleratorTypeModel{}

		outputList = append(outputList, output)
	}
	return outputList
}

func flattenArchitectureTypeModelArray(inputList *[]fleets.ArchitectureType) []ArchitectureTypeModel {
	var outputList []ArchitectureTypeModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := ArchitectureTypeModel{}

		outputList = append(outputList, output)
	}
	return outputList
}

func flattenCPUManufacturerModelArray(inputList *[]fleets.CPUManufacturer) []CPUManufacturerModel {
	var outputList []CPUManufacturerModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := CPUManufacturerModel{}

		outputList = append(outputList, output)
	}
	return outputList
}

func flattenLocalStorageDiskTypeModelArray(inputList *[]fleets.LocalStorageDiskType) []LocalStorageDiskTypeModel {
	var outputList []LocalStorageDiskTypeModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := LocalStorageDiskTypeModel{}

		outputList = append(outputList, output)
	}
	return outputList
}

func flattenVMAttributeMinMaxDoubleModel(input *fleets.VMAttributeMinMaxDouble) []VMAttributeMinMaxDoubleModel {
	var outputList []VMAttributeMinMaxDoubleModel
	if input == nil {
		return outputList
	}
	output := VMAttributeMinMaxDoubleModel{}
	if input.Max != nil {
		output.Max = *input.Max
	}

	if input.Min != nil {
		output.Min = *input.Min
	}

	return append(outputList, output)
}

func flattenVMCategoryModelArray(inputList *[]fleets.VMCategory) []VMCategoryModel {
	var outputList []VMCategoryModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := VMCategoryModel{}

		outputList = append(outputList, output)
	}
	return outputList
}

func flattenVMSizeProfileModelArray(inputList *[]fleets.VMSizeProfile) []VMSizeProfileModel {
	var outputList []VMSizeProfileModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := VMSizeProfileModel{
			Name: input.Name,
		}

		if input.Rank != nil {
			output.Rank = *input.Rank
		}
		outputList = append(outputList, output)
	}
	return outputList
}
