package azurefleet

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	tagsHelper "github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurefleet/2024-05-01-preview/fleets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AzureFleetModel struct {
	Name                   string                        `tfschema:"name"`
	ResourceGroupName      string                        `tfschema:"resource_group_name"`
	ComputeProfile         []ComputeProfileModel         `tfschema:"compute_profile"`
	Location               string                        `tfschema:"location"`
	Plan                   []PlanModel                   `tfschema:"plan"`
	RegularPriorityProfile []RegularPriorityProfileModel `tfschema:"regular_priority_profile"`
	SpotPriorityProfile    []SpotPriorityProfileModel    `tfschema:"spot_priority_profile"`
	Tags                   map[string]string             `tfschema:"tags"`
	VMSizesProfile         []VMSizeProfileModel          `tfschema:"vm_sizes_profile"`
	Zones                  []string                      `tfschema:"zones"`
	TimeCreated            string                        `tfschema:"time_created"`
	UniqueId               string                        `tfschema:"unique_id"`
}

type ComputeProfileModel struct {
	BaseVirtualMachineProfile []BaseVirtualMachineProfileModel `tfschema:"base_virtual_machine_profile"`
	ComputeApiVersion         string                           `tfschema:"compute_api_version"`
	PlatformFaultDomainCount  int64                            `tfschema:"platform_fault_domain_count"`
}

type BaseVirtualMachineProfileModel struct {
	ApplicationProfile       []ApplicationProfileModel                     `tfschema:"application_profile"`
	CapacityReservation      []CapacityReservationProfileModel             `tfschema:"capacity_reservation"`
	DiagnosticsProfile       []DiagnosticsProfileModel                     `tfschema:"diagnostics_profile"`
	ExtensionProfile         []VirtualMachineScaleSetExtensionProfileModel `tfschema:"extension_profile"`
	HardwareProfile          []VirtualMachineScaleSetHardwareProfileModel  `tfschema:"hardware_profile"`
	LicenseType              string                                        `tfschema:"license_type"`
	NetworkProfile           []VirtualMachineScaleSetNetworkProfileModel   `tfschema:"network_profile"`
	OsProfile                []VirtualMachineScaleSetOSProfileModel        `tfschema:"os_profile"`
	ScheduledEventsProfile   []ScheduledEventsProfileModel                 `tfschema:"scheduled_events_profile"`
	SecurityPostureReference []SecurityPostureReferenceModel               `tfschema:"security_posture_reference"`
	SecurityProfile          []SecurityProfileModel                        `tfschema:"security_profile"`
	ServiceArtifactReference []ServiceArtifactReferenceModel               `tfschema:"service_artifact_reference"`
	StorageProfile           []VirtualMachineScaleSetStorageProfileModel   `tfschema:"storage_profile"`
	TimeCreated              string                                        `tfschema:"time_created"`
	UserData                 string                                        `tfschema:"user_data"`
}

type ApplicationProfileModel struct {
	GalleryApplications []VMGalleryApplicationModel `tfschema:"gallery_applications"`
}

type VMGalleryApplicationModel struct {
	ConfigurationReference          string `tfschema:"configuration_reference"`
	EnableAutomaticUpgrade          bool   `tfschema:"enable_automatic_upgrade"`
	Order                           int64  `tfschema:"order"`
	PackageReferenceId              string `tfschema:"package_reference_id"`
	Tags                            string `tfschema:"tags"`
	TreatFailureAsDeploymentFailure bool   `tfschema:"treat_failure_as_deployment_failure"`
}

type CapacityReservationProfileModel struct {
	CapacityReservationGroup []SubResourceModel `tfschema:"capacity_reservation_group"`
}

type SubResourceModel struct {
	Id string `tfschema:"id"`
}

type DiagnosticsProfileModel struct {
	BootDiagnostics []BootDiagnosticsModel `tfschema:"boot_diagnostics"`
}

type BootDiagnosticsModel struct {
	Enabled    bool   `tfschema:"enabled"`
	StorageUri string `tfschema:"storage_uri"`
}

type VirtualMachineScaleSetExtensionProfileModel struct {
	Extensions           []VirtualMachineScaleSetExtensionModel `tfschema:"extensions"`
	ExtensionsTimeBudget string                                 `tfschema:"extensions_time_budget"`
}

type VirtualMachineScaleSetExtensionModel struct {
	Id         string                                           `tfschema:"id"`
	Name       string                                           `tfschema:"name"`
	Properties []VirtualMachineScaleSetExtensionPropertiesModel `tfschema:"properties"`
	Type       string                                           `tfschema:"type"`
}

type VirtualMachineScaleSetExtensionPropertiesModel struct {
	AutoUpgradeMinorVersion       bool                           `tfschema:"auto_upgrade_minor_version"`
	EnableAutomaticUpgrade        bool                           `tfschema:"enable_automatic_upgrade"`
	ForceUpdateTag                string                         `tfschema:"force_update_tag"`
	ProtectedSettings             string                         `tfschema:"protected_settings"`
	ProtectedSettingsFromKeyVault []KeyVaultSecretReferenceModel `tfschema:"protected_settings_from_key_vault"`
	ProvisionAfterExtensions      []string                       `tfschema:"provision_after_extensions"`
	ProvisioningState             string                         `tfschema:"provisioning_state"`
	Publisher                     string                         `tfschema:"publisher"`
	Settings                      string                         `tfschema:"settings"`
	SuppressFailures              bool                           `tfschema:"suppress_failures"`
	Type                          string                         `tfschema:"type"`
	TypeHandlerVersion            string                         `tfschema:"type_handler_version"`
}

type KeyVaultSecretReferenceModel struct {
	SecretUrl   string             `tfschema:"secret_url"`
	SourceVault []SubResourceModel `tfschema:"source_vault"`
}

type VirtualMachineScaleSetHardwareProfileModel struct {
	VMSizeProperties []VMSizePropertiesModel `tfschema:"vm_size_properties"`
}

type VMSizePropertiesModel struct {
	VCPUsAvailable int64 `tfschema:"vcp_us_available"`
	VCPUsPerCore   int64 `tfschema:"vcp_us_per_core"`
}

type VirtualMachineScaleSetNetworkProfileModel struct {
	HealthProbe                    []ApiEntityReferenceModel                         `tfschema:"health_probe"`
	NetworkApiVersion              fleets.NetworkApiVersion                          `tfschema:"network_api_version"`
	NetworkInterfaceConfigurations []VirtualMachineScaleSetNetworkConfigurationModel `tfschema:"network_interface_configurations"`
}

type ApiEntityReferenceModel struct {
	Id string `tfschema:"id"`
}

type VirtualMachineScaleSetNetworkConfigurationModel struct {
	Name       string                                                      `tfschema:"name"`
	Properties []VirtualMachineScaleSetNetworkConfigurationPropertiesModel `tfschema:"properties"`
}

type VirtualMachineScaleSetNetworkConfigurationPropertiesModel struct {
	AuxiliaryMode               fleets.NetworkInterfaceAuxiliaryMode                         `tfschema:"auxiliary_mode"`
	AuxiliarySku                fleets.NetworkInterfaceAuxiliarySku                          `tfschema:"auxiliary_sku"`
	DeleteOption                fleets.DeleteOptions                                         `tfschema:"delete_option"`
	DisableTcpStateTracking     bool                                                         `tfschema:"disable_tcp_state_tracking"`
	DnsSettings                 []VirtualMachineScaleSetNetworkConfigurationDnsSettingsModel `tfschema:"dns_settings"`
	EnableAcceleratedNetworking bool                                                         `tfschema:"enable_accelerated_networking"`
	EnableFpga                  bool                                                         `tfschema:"enable_fpga"`
	EnableIPForwarding          bool                                                         `tfschema:"enable_ip_forwarding"`
	IPConfigurations            []VirtualMachineScaleSetIPConfigurationModel                 `tfschema:"ip_configurations"`
	NetworkSecurityGroup        []SubResourceModel                                           `tfschema:"network_security_group"`
	Primary                     bool                                                         `tfschema:"primary"`
}

type VirtualMachineScaleSetNetworkConfigurationDnsSettingsModel struct {
	DnsServers []string `tfschema:"dns_servers"`
}

type VirtualMachineScaleSetIPConfigurationModel struct {
	Name       string                                                 `tfschema:"name"`
	Properties []VirtualMachineScaleSetIPConfigurationPropertiesModel `tfschema:"properties"`
}

type VirtualMachineScaleSetIPConfigurationPropertiesModel struct {
	ApplicationGatewayBackendAddressPools []SubResourceModel                                        `tfschema:"application_gateway_backend_address_pools"`
	ApplicationSecurityGroups             []SubResourceModel                                        `tfschema:"application_security_groups"`
	LoadBalancerBackendAddressPools       []SubResourceModel                                        `tfschema:"load_balancer_backend_address_pools"`
	LoadBalancerInboundNatPools           []SubResourceModel                                        `tfschema:"load_balancer_inbound_nat_pools"`
	Primary                               bool                                                      `tfschema:"primary"`
	PrivateIPAddressVersion               fleets.IPVersion                                          `tfschema:"private_ip_address_version"`
	PublicIPAddressConfiguration          []VirtualMachineScaleSetPublicIPAddressConfigurationModel `tfschema:"public_ip_address_configuration"`
	Subnet                                []ApiEntityReferenceModel                                 `tfschema:"subnet"`
}

type VirtualMachineScaleSetPublicIPAddressConfigurationModel struct {
	Name       string                                                              `tfschema:"name"`
	Properties []VirtualMachineScaleSetPublicIPAddressConfigurationPropertiesModel `tfschema:"properties"`
	Sku        []PublicIPAddressSkuModel                                           `tfschema:"sku"`
}

type VirtualMachineScaleSetPublicIPAddressConfigurationPropertiesModel struct {
	DeleteOption           fleets.DeleteOptions                                                 `tfschema:"delete_option"`
	DnsSettings            []VirtualMachineScaleSetPublicIPAddressConfigurationDnsSettingsModel `tfschema:"dns_settings"`
	IPTags                 []VirtualMachineScaleSetIPTagModel                                   `tfschema:"ip_tags"`
	IdleTimeoutInMinutes   int64                                                                `tfschema:"idle_timeout_in_minutes"`
	PublicIPAddressVersion fleets.IPVersion                                                     `tfschema:"public_ip_address_version"`
	PublicIPPrefix         []SubResourceModel                                                   `tfschema:"public_ip_prefix"`
}

type VirtualMachineScaleSetPublicIPAddressConfigurationDnsSettingsModel struct {
	DomainNameLabel      string                           `tfschema:"domain_name_label"`
	DomainNameLabelScope fleets.DomainNameLabelScopeTypes `tfschema:"domain_name_label_scope"`
}

type VirtualMachineScaleSetIPTagModel struct {
	IPTagType string `tfschema:"ip_tag_type"`
	Tag       string `tfschema:"tag"`
}

type PublicIPAddressSkuModel struct {
	Name fleets.PublicIPAddressSkuName `tfschema:"name"`
	Tier fleets.PublicIPAddressSkuTier `tfschema:"tier"`
}

type VirtualMachineScaleSetOSProfileModel struct {
	AdminPassword               string                      `tfschema:"admin_password"`
	AdminUsername               string                      `tfschema:"admin_username"`
	AllowExtensionOperations    bool                        `tfschema:"allow_extension_operations"`
	ComputerNamePrefix          string                      `tfschema:"computer_name_prefix"`
	CustomData                  string                      `tfschema:"custom_data"`
	LinuxConfiguration          []LinuxConfigurationModel   `tfschema:"linux_configuration"`
	RequireGuestProvisionSignal bool                        `tfschema:"require_guest_provision_signal"`
	Secrets                     []VaultSecretGroupModel     `tfschema:"secrets"`
	WindowsConfiguration        []WindowsConfigurationModel `tfschema:"windows_configuration"`
}

type LinuxConfigurationModel struct {
	DisablePasswordAuthentication bool                      `tfschema:"disable_password_authentication"`
	EnableVMAgentPlatformUpdates  bool                      `tfschema:"enable_vm_agent_platform_updates"`
	PatchSettings                 []LinuxPatchSettingsModel `tfschema:"patch_settings"`
	ProvisionVMAgent              bool                      `tfschema:"provision_vm_agent"`
	Ssh                           []SshConfigurationModel   `tfschema:"ssh"`
}

type LinuxPatchSettingsModel struct {
	AssessmentMode              fleets.LinuxPatchAssessmentMode                     `tfschema:"assessment_mode"`
	AutomaticByPlatformSettings []LinuxVMGuestPatchAutomaticByPlatformSettingsModel `tfschema:"automatic_by_platform_settings"`
	PatchMode                   fleets.LinuxVMGuestPatchMode                        `tfschema:"patch_mode"`
}

type LinuxVMGuestPatchAutomaticByPlatformSettingsModel struct {
	BypassPlatformSafetyChecksOnUserSchedule bool                                                     `tfschema:"bypass_platform_safety_checks_on_user_schedule"`
	RebootSetting                            fleets.LinuxVMGuestPatchAutomaticByPlatformRebootSetting `tfschema:"reboot_setting"`
}

type SshConfigurationModel struct {
	PublicKeys []SshPublicKeyModel `tfschema:"public_keys"`
}

type SshPublicKeyModel struct {
	KeyData string `tfschema:"key_data"`
	Path    string `tfschema:"path"`
}

type VaultSecretGroupModel struct {
	SourceVault       []SubResourceModel      `tfschema:"source_vault"`
	VaultCertificates []VaultCertificateModel `tfschema:"vault_certificates"`
}

type VaultCertificateModel struct {
	CertificateStore string `tfschema:"certificate_store"`
	CertificateUrl   string `tfschema:"certificate_url"`
}

type WindowsConfigurationModel struct {
	AdditionalUnattendContent    []AdditionalUnattendContentModel `tfschema:"additional_unattend_content"`
	EnableAutomaticUpdates       bool                             `tfschema:"enable_automatic_updates"`
	EnableVMAgentPlatformUpdates bool                             `tfschema:"enable_vm_agent_platform_updates"`
	PatchSettings                []PatchSettingsModel             `tfschema:"patch_settings"`
	ProvisionVMAgent             bool                             `tfschema:"provision_vm_agent"`
	TimeZone                     string                           `tfschema:"time_zone"`
	WinRM                        []WinRMConfigurationModel        `tfschema:"win_rm"`
}

type AdditionalUnattendContentModel struct {
	ComponentName fleets.ComponentName `tfschema:"component_name"`
	Content       string               `tfschema:"content"`
	PassName      fleets.PassName      `tfschema:"pass_name"`
	SettingName   fleets.SettingNames  `tfschema:"setting_name"`
}

type PatchSettingsModel struct {
	AssessmentMode              fleets.WindowsPatchAssessmentMode                     `tfschema:"assessment_mode"`
	AutomaticByPlatformSettings []WindowsVMGuestPatchAutomaticByPlatformSettingsModel `tfschema:"automatic_by_platform_settings"`
	EnableHotpatching           bool                                                  `tfschema:"enable_hotpatching"`
	PatchMode                   fleets.WindowsVMGuestPatchMode                        `tfschema:"patch_mode"`
}

type WindowsVMGuestPatchAutomaticByPlatformSettingsModel struct {
	BypassPlatformSafetyChecksOnUserSchedule bool                                                       `tfschema:"bypass_platform_safety_checks_on_user_schedule"`
	RebootSetting                            fleets.WindowsVMGuestPatchAutomaticByPlatformRebootSetting `tfschema:"reboot_setting"`
}

type WinRMConfigurationModel struct {
	Listeners []WinRMListenerModel `tfschema:"listeners"`
}

type WinRMListenerModel struct {
	CertificateUrl string               `tfschema:"certificate_url"`
	Protocol       fleets.ProtocolTypes `tfschema:"protocol"`
}

type ScheduledEventsProfileModel struct {
	OsImageNotificationProfile   []OSImageNotificationProfileModel   `tfschema:"os_image_notification_profile"`
	TerminateNotificationProfile []TerminateNotificationProfileModel `tfschema:"terminate_notification_profile"`
}

type OSImageNotificationProfileModel struct {
	Enable           bool   `tfschema:"enable"`
	NotBeforeTimeout string `tfschema:"not_before_timeout"`
}

type TerminateNotificationProfileModel struct {
	Enable           bool   `tfschema:"enable"`
	NotBeforeTimeout string `tfschema:"not_before_timeout"`
}

type SecurityPostureReferenceModel struct {
	ExcludeExtensions []string `tfschema:"exclude_extensions"`
	Id                string   `tfschema:"id"`
	IsOverridable     bool     `tfschema:"is_overridable"`
}

type SecurityProfileModel struct {
	EncryptionAtHost   bool                      `tfschema:"encryption_at_host"`
	EncryptionIdentity []EncryptionIdentityModel `tfschema:"encryption_identity"`
	ProxyAgentSettings []ProxyAgentSettingsModel `tfschema:"proxy_agent_settings"`
	SecurityType       fleets.SecurityTypes      `tfschema:"security_type"`
	UefiSettings       []UefiSettingsModel       `tfschema:"uefi_settings"`
}

type EncryptionIdentityModel struct {
	UserAssignedIdentityResourceId string `tfschema:"user_assigned_identity_resource_id"`
}

type ProxyAgentSettingsModel struct {
	Enabled          bool        `tfschema:"enabled"`
	KeyIncarnationId int64       `tfschema:"key_incarnation_id"`
	Mode             fleets.Mode `tfschema:"mode"`
}

type UefiSettingsModel struct {
	SecureBootEnabled bool `tfschema:"secure_boot_enabled"`
	VTpmEnabled       bool `tfschema:"v_tpm_enabled"`
}

type ServiceArtifactReferenceModel struct {
	Id string `tfschema:"id"`
}

type VirtualMachineScaleSetStorageProfileModel struct {
	DataDisks          []VirtualMachineScaleSetDataDiskModel `tfschema:"data_disks"`
	DiskControllerType fleets.DiskControllerTypes            `tfschema:"disk_controller_type"`
	ImageReference     []ImageReferenceModel                 `tfschema:"image_reference"`
	OsDisk             []VirtualMachineScaleSetOSDiskModel   `tfschema:"os_disk"`
}

type VirtualMachineScaleSetDataDiskModel struct {
	Caching                 fleets.CachingTypes                                `tfschema:"caching"`
	CreateOption            fleets.DiskCreateOptionTypes                       `tfschema:"create_option"`
	DeleteOption            fleets.DiskDeleteOptionTypes                       `tfschema:"delete_option"`
	DiskIOPSReadWrite       int64                                              `tfschema:"disk_iops_read_write"`
	DiskMBpsReadWrite       int64                                              `tfschema:"disk_m_bps_read_write"`
	DiskSizeGB              int64                                              `tfschema:"disk_size_gb"`
	Lun                     int64                                              `tfschema:"lun"`
	ManagedDisk             []VirtualMachineScaleSetManagedDiskParametersModel `tfschema:"managed_disk"`
	Name                    string                                             `tfschema:"name"`
	WriteAcceleratorEnabled bool                                               `tfschema:"write_accelerator_enabled"`
}

type VirtualMachineScaleSetManagedDiskParametersModel struct {
	DiskEncryptionSet  []DiskEncryptionSetParametersModel `tfschema:"disk_encryption_set"`
	SecurityProfile    []VMDiskSecurityProfileModel       `tfschema:"security_profile"`
	StorageAccountType fleets.StorageAccountTypes         `tfschema:"storage_account_type"`
}

type DiskEncryptionSetParametersModel struct {
	Id string `tfschema:"id"`
}

type VMDiskSecurityProfileModel struct {
	DiskEncryptionSet      []DiskEncryptionSetParametersModel `tfschema:"disk_encryption_set"`
	SecurityEncryptionType fleets.SecurityEncryptionTypes     `tfschema:"security_encryption_type"`
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
	Caching                 fleets.CachingTypes                                `tfschema:"caching"`
	CreateOption            fleets.DiskCreateOptionTypes                       `tfschema:"create_option"`
	DeleteOption            fleets.DiskDeleteOptionTypes                       `tfschema:"delete_option"`
	DiffDiskSettings        []DiffDiskSettingsModel                            `tfschema:"diff_disk_settings"`
	DiskSizeGB              int64                                              `tfschema:"disk_size_gb"`
	Image                   []VirtualHardDiskModel                             `tfschema:"image"`
	ManagedDisk             []VirtualMachineScaleSetManagedDiskParametersModel `tfschema:"managed_disk"`
	Name                    string                                             `tfschema:"name"`
	OsType                  fleets.OperatingSystemTypes                        `tfschema:"os_type"`
	VhdContainers           []string                                           `tfschema:"vhd_containers"`
	WriteAcceleratorEnabled bool                                               `tfschema:"write_accelerator_enabled"`
}

type DiffDiskSettingsModel struct {
	Option    fleets.DiffDiskOptions   `tfschema:"option"`
	Placement fleets.DiffDiskPlacement `tfschema:"placement"`
}

type VirtualHardDiskModel struct {
	Uri string `tfschema:"uri"`
}

type PlanModel struct {
	Name          string `tfschema:"name"`
	Product       string `tfschema:"product"`
	PromotionCode string `tfschema:"promotion_code"`
	Publisher     string `tfschema:"publisher"`
	Version       string `tfschema:"version"`
}

type RegularPriorityProfileModel struct {
	AllocationStrategy fleets.RegularPriorityAllocationStrategy `tfschema:"allocation_strategy"`
	Capacity           int64                                    `tfschema:"capacity"`
	MinCapacity        int64                                    `tfschema:"min_capacity"`
}

type SpotPriorityProfileModel struct {
	AllocationStrategy fleets.SpotAllocationStrategy `tfschema:"allocation_strategy"`
	Capacity           int64                         `tfschema:"capacity"`
	EvictionPolicy     fleets.EvictionPolicy         `tfschema:"eviction_policy"`
	Maintain           bool                          `tfschema:"maintain"`
	MaxPricePerVM      float64                       `tfschema:"max_price_per_vm"`
	MinCapacity        int64                         `tfschema:"min_capacity"`
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
	return &AzureFleetFleetModel{}
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

		"compute_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
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

								"capacity_reservation": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"capacity_reservation_group": {
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

								"extension_profile": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
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

																				"source_vault": {
																					Type:     pluginsdk.TypeList,
																					Required: true,
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
										},
									},
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
																			string(fleets.NetworkInterfaceAuxiliaryModeFloating),
																			string(fleets.NetworkInterfaceAuxiliaryModeNone),
																			string(fleets.NetworkInterfaceAuxiliaryModeAcceleratedConnections),
																		}, false),
																	},

																	"auxiliary_sku": {
																		Type:     pluginsdk.TypeString,
																		Optional: true,
																		ValidateFunc: validation.StringInSlice([]string{
																			string(fleets.NetworkInterfaceAuxiliarySkuAFour),
																			string(fleets.NetworkInterfaceAuxiliarySkuAEight),
																			string(fleets.NetworkInterfaceAuxiliarySkuNone),
																			string(fleets.NetworkInterfaceAuxiliarySkuAOne),
																			string(fleets.NetworkInterfaceAuxiliarySkuATwo),
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
																																		string(fleets.DomainNameLabelScopeTypesResourceGroupReuse),
																																		string(fleets.DomainNameLabelScopeTypesNoReuse),
																																		string(fleets.DomainNameLabelScopeTypesTenantReuse),
																																		string(fleets.DomainNameLabelScopeTypesSubscriptionReuse),
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
																															string(fleets.IPVersionIPvFour),
																															string(fleets.IPVersionIPvSix),
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
														"source_vault": {
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
																						string(fleets.WindowsVMGuestPatchAutomaticByPlatformRebootSettingUnknown),
																						string(fleets.WindowsVMGuestPatchAutomaticByPlatformRebootSettingIfRequired),
																						string(fleets.WindowsVMGuestPatchAutomaticByPlatformRebootSettingNever),
																						string(fleets.WindowsVMGuestPatchAutomaticByPlatformRebootSettingAlways),
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
																			string(fleets.WindowsVMGuestPatchModeAutomaticByOS),
																			string(fleets.WindowsVMGuestPatchModeAutomaticByPlatform),
																			string(fleets.WindowsVMGuestPatchModeManual),
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
																string(fleets.CachingTypesReadWrite),
																string(fleets.CachingTypesNone),
																string(fleets.CachingTypesReadOnly),
															}, false),
														},

														"create_option": {
															Type:     pluginsdk.TypeString,
															Required: true,
															ValidateFunc: validation.StringInSlice([]string{
																string(fleets.DiskCreateOptionTypesAttach),
																string(fleets.DiskCreateOptionTypesCopy),
																string(fleets.DiskCreateOptionTypesRestore),
																string(fleets.DiskCreateOptionTypesFromImage),
																string(fleets.DiskCreateOptionTypesEmpty),
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
													string(fleets.DiskControllerTypesSCSI),
													string(fleets.DiskControllerTypesNVMe),
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

		"vm_sizes_profile": {
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
			var model AzureFleetFleetModel
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
				},
				Tags:  &model.Tags,
				Zones: &model.Zones,
			}

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

			var model AzureFleetFleetModel
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

			state := AzureFleetFleetModel{
				Name:              id.FleetName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(model.Location),
			}

			if model := resp.Model; model != nil {
				identityValue, err := identity.FlattenLegacySystemAndUserAssignedMap(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}

				if err := metadata.ResourceData.Set("identity", identityValue); err != nil {
					return fmt.Errorf("setting `identity`: %+v", err)
				}

				state.Plan = flattenPlanModel(model.Plan)
				if properties := model.Properties; properties != nil {
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

func expandComputeProfileModel(inputList []ComputeProfileModel) (*fleets.ComputeProfile, error) {
	if len(inputList) == 0 {
		return nil, nil
	}
	input := &inputList[0]
	output := fleets.ComputeProfile{
		PlatformFaultDomainCount: &input.PlatformFaultDomainCount,
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

func expandBaseVirtualMachineProfileModel(inputList []BaseVirtualMachineProfileModel) (*fleets.BaseVirtualMachineProfile, error) {
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

func expandVirtualMachineScaleSetExtensionModelArray(inputList []VirtualMachineScaleSetExtensionModel) (*[]fleets.VirtualMachineScaleSetExtension, error) {
	var outputList []fleets.VirtualMachineScaleSetExtension
	for _, v := range inputList {
		input := v
		output := fleets.VirtualMachineScaleSetExtension{}

		if input.Name != "" {
			output.Name = &input.Name
		}

		propertiesValue, err := expandVirtualMachineScaleSetExtensionPropertiesModel(input.Properties)
		if err != nil {
			return nil, err
		}

		output.Properties = propertiesValue

		outputList = append(outputList, output)
	}
	return &outputList, nil
}

func expandVirtualMachineScaleSetExtensionPropertiesModel(inputList []VirtualMachineScaleSetExtensionPropertiesModel) (*fleets.VirtualMachineScaleSetExtensionProperties, error) {
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
	err := json.Unmarshal([]byte(input.Settings), &settingsValue)
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
		SecretUrl: input.SecretUrl,
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

func expandVirtualMachineScaleSetNetworkProfileModel(inputList []VirtualMachineScaleSetNetworkProfileModel) *fleets.VirtualMachineScaleSetNetworkProfile {
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

func expandVirtualMachineScaleSetNetworkConfigurationModelArray(inputList []VirtualMachineScaleSetNetworkConfigurationModel) *[]fleets.VirtualMachineScaleSetNetworkConfiguration {
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
			Properties: expandVirtualMachineScaleSetIPConfigurationPropertiesModel(input.Properties),
		}

		outputList = append(outputList, output)
	}
	return &outputList
}

func expandVirtualMachineScaleSetIPConfigurationPropertiesModel(inputList []VirtualMachineScaleSetIPConfigurationPropertiesModel) *fleets.VirtualMachineScaleSetIPConfigurationProperties {
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
		Properties: expandVirtualMachineScaleSetPublicIPAddressConfigurationPropertiesModel(input.Properties),
		Sku:        expandPublicIPAddressSkuModel(input.Sku),
	}

	return &output
}

func expandVirtualMachineScaleSetPublicIPAddressConfigurationPropertiesModel(inputList []VirtualMachineScaleSetPublicIPAddressConfigurationPropertiesModel) *fleets.VirtualMachineScaleSetPublicIPAddressConfigurationProperties {
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

func expandVirtualMachineScaleSetOSProfileModel(inputList []VirtualMachineScaleSetOSProfileModel) *fleets.VirtualMachineScaleSetOSProfile {
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

func expandVirtualMachineScaleSetStorageProfileModel(inputList []VirtualMachineScaleSetStorageProfileModel) *fleets.VirtualMachineScaleSetStorageProfile {
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

func flattenComputeProfileModel(input *fleets.ComputeProfile) ([]ComputeProfileModel, error) {
	var outputList []ComputeProfileModel
	if input == nil {
		return outputList, nil
	}
	output := ComputeProfileModel{}
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

func flattenBaseVirtualMachineProfileModel(input *fleets.BaseVirtualMachineProfile) ([]BaseVirtualMachineProfileModel, error) {
	var outputList []BaseVirtualMachineProfileModel
	if input == nil {
		return outputList, nil
	}
	output := BaseVirtualMachineProfileModel{
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

		propertiesValue, err := flattenVirtualMachineScaleSetExtensionPropertiesModel(input.Properties)
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

func flattenVirtualMachineScaleSetExtensionPropertiesModel(input *fleets.VirtualMachineScaleSetExtensionProperties) ([]VirtualMachineScaleSetExtensionPropertiesModel, error) {
	var outputList []VirtualMachineScaleSetExtensionPropertiesModel
	if input == nil {
		return outputList, nil
	}
	output := VirtualMachineScaleSetExtensionPropertiesModel{
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

func flattenVirtualMachineScaleSetNetworkProfileModel(input *fleets.VirtualMachineScaleSetNetworkProfile) []VirtualMachineScaleSetNetworkProfileModel {
	var outputList []VirtualMachineScaleSetNetworkProfileModel
	if input == nil {
		return outputList
	}
	output := VirtualMachineScaleSetNetworkProfileModel{
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

func flattenVirtualMachineScaleSetNetworkConfigurationModelArray(inputList *[]fleets.VirtualMachineScaleSetNetworkConfiguration) []VirtualMachineScaleSetNetworkConfigurationModel {
	var outputList []VirtualMachineScaleSetNetworkConfigurationModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := VirtualMachineScaleSetNetworkConfigurationModel{
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
			Properties: flattenVirtualMachineScaleSetIPConfigurationPropertiesModel(input.Properties),
		}

		outputList = append(outputList, output)
	}
	return outputList
}

func flattenVirtualMachineScaleSetIPConfigurationPropertiesModel(input *fleets.VirtualMachineScaleSetIPConfigurationProperties) []VirtualMachineScaleSetIPConfigurationPropertiesModel {
	var outputList []VirtualMachineScaleSetIPConfigurationPropertiesModel
	if input == nil {
		return outputList
	}
	output := VirtualMachineScaleSetIPConfigurationPropertiesModel{
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
		Properties: flattenVirtualMachineScaleSetPublicIPAddressConfigurationPropertiesModel(input.Properties),
		Sku:        flattenPublicIPAddressSkuModel(input.Sku),
	}

	return append(outputList, output)
}

func flattenVirtualMachineScaleSetPublicIPAddressConfigurationPropertiesModel(input *fleets.VirtualMachineScaleSetPublicIPAddressConfigurationProperties) []VirtualMachineScaleSetPublicIPAddressConfigurationPropertiesModel {
	var outputList []VirtualMachineScaleSetPublicIPAddressConfigurationPropertiesModel
	if input == nil {
		return outputList
	}
	output := VirtualMachineScaleSetPublicIPAddressConfigurationPropertiesModel{
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

func flattenVirtualMachineScaleSetOSProfileModel(input *fleets.VirtualMachineScaleSetOSProfile) []VirtualMachineScaleSetOSProfileModel {
	var outputList []VirtualMachineScaleSetOSProfileModel
	if input == nil {
		return outputList
	}
	output := VirtualMachineScaleSetOSProfileModel{
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

func flattenVirtualMachineScaleSetStorageProfileModel(input *fleets.VirtualMachineScaleSetStorageProfile) []VirtualMachineScaleSetStorageProfileModel {
	var outputList []VirtualMachineScaleSetStorageProfileModel
	if input == nil {
		return outputList
	}
	output := VirtualMachineScaleSetStorageProfileModel{
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
