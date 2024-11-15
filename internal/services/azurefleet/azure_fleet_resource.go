package azurefleet

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurefleet/2024-11-01/fleets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type AzureFleetResourceModel struct {
	Name                      string                                     `tfschema:"name"`
	ResourceGroupName         string                                     `tfschema:"resource_group_name"`
	Identity                  []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Location                  string                                     `tfschema:"location"`
	Plan                      []PlanModel                                `tfschema:"plan"`
	AdditionalLocationProfile []AdditionalLocationProfileModel           `tfschema:"additional_location_profile"`
	ComputeProfile            []ComputeProfileModel                      `tfschema:"compute_profile"`
	RegularPriorityProfile    []RegularPriorityProfileModel              `tfschema:"regular_priority_profile"`
	SpotPriorityProfile       []SpotPriorityProfileModel                 `tfschema:"spot_priority_profile"`
	UniqueId                  string                                     `tfschema:"unique_id"`
	VMAttributes              []VMAttributesModel                        `tfschema:"vm_attributes"`
	VMSizesProfile            []VMSizeProfileModel                       `tfschema:"vm_sizes_profile"`
	Tags                      map[string]string                          `tfschema:"tags"`
	Zones                     []string                                   `tfschema:"zones"`
}

type AdditionalLocationProfileModel struct {
	Location                      string                       `tfschema:"location"`
	VirtualMachineProfileOverride []VirtualMachineProfileModel `tfschema:"virtual_machine_profile_override"`
}

type VirtualMachineProfileModel struct {
	GalleryApplicationProfile            []GalleryApplicationModel       `tfschema:"gallery_application"`
	CapacityReservationGroupId           string                          `tfschema:"capacity_reservation_group_id"`
	BootDiagnosticEnabled                bool                            `tfschema:"boot_diagnostic_enabled"`
	BootDiagnosticStorageAccountEndpoint string                          `tfschema:"boot_diagnostic_storage_account_endpoint"`
	Extensions                           []ExtensionsModel               `tfschema:"extensions"`
	ExtensionsTimeBudget                 string                          `tfschema:"extensions_time_budget"`
	VMSize                               []VMSizeModel                   `tfschema:"vm_size"`
	LicenseType                          string                          `tfschema:"license_type"`
	NetworkHealthProbeId                 string                          `tfschema:"network_health_probe_id"`
	NetworkApiVersion                    string                          `tfschema:"network_api_version"`
	NetworkInterface                     []NetworkInterfaceModel         `tfschema:"network_interface"`
	OsProfile                            []OSProfileModel                `tfschema:"os_profile"`
	ScheduledEventTerminationEnabled     bool                            `tfschema:"scheduled_event_termination_enabled"`
	ScheduledEventTerminationTimeout     string                          `tfschema:"scheduled_event_termination_timeout"`
	ScheduledEventOsImageEnabled         bool                            `tfschema:"scheduled_event_os_image_enabled"`
	ScheduledEventOsImageTimeout         string                          `tfschema:"scheduled_event_os_image_timeout"`
	SecurityPostureReference             []SecurityPostureReferenceModel `tfschema:"security_posture_reference"`
	SecurityProfile                      []SecurityProfileModel          `tfschema:"security_profile"`
	ServiceArtifactId                    string                          `tfschema:"service_artifact_id"`
	StorageProfile                       []StorageProfileModel           `tfschema:"storage_profile"`
	UserDataBase64                       string                          `tfschema:"user_data_base64"`
}

type GalleryApplicationModel struct {
	ConfigurationBlobUri                   string `tfschema:"configuration_blob_uri"`
	AutomaticUpgradeEnabled                bool   `tfschema:"automatic_upgrade_enabled"`
	Order                                  int64  `tfschema:"order"`
	VersionId                              string `tfschema:"version_id"`
	Tags                                   string `tfschema:"tags"`
	TreatFailureAsDeploymentFailureEnabled bool   `tfschema:"treat_failure_as_deployment_failure_enabled"`
}

type ExtensionsModel struct {
	Name                           string                               `tfschema:"name"`
	Publisher                      string                               `tfschema:"publisher"`
	Type                           string                               `tfschema:"type"`
	TypeHandlerVersion             string                               `tfschema:"type_handler_version"`
	AutoUpgradeMinorVersionEnabled bool                                 `tfschema:"auto_upgrade_minor_version_enabled"`
	AutomaticUpgradeEnabled        bool                                 `tfschema:"automatic_upgrade_enabled"`
	ForceUpdateTag                 string                               `tfschema:"force_update_tag"`
	ProtectedSettingsJson          string                               `tfschema:"protected_settings_json"`
	ProtectedSettingsFromKeyVault  []ProtectedSettingsFromKeyVaultModel `tfschema:"protected_settings_from_key_vault"`
	ProvisionAfterExtensions       []string                             `tfschema:"provision_after_extensions"`
	SettingsJson                   string                               `tfschema:"settings_json"`
	SuppressFailuresEnabled        bool                                 `tfschema:"suppress_failures_enabled"`
}

type ProtectedSettingsFromKeyVaultModel struct {
	SecretUrl     string `tfschema:"secret_url"`
	SourceVaultId string `tfschema:"source_vault_id"`
}

type VMSizeModel struct {
	VCPUAvailableCount int64 `tfschema:"vcpu_available_count"`
	VCPUPerCoreCount   int64 `tfschema:"vcpu_per_core_count"`
}

type NetworkInterfaceModel struct {
	Name                         string                 `tfschema:"name"`
	AuxiliaryMode                string                 `tfschema:"auxiliary_mode"`
	AuxiliarySku                 string                 `tfschema:"auxiliary_sku"`
	DeleteOption                 string                 `tfschema:"delete_option"`
	TcpStateTrackingEnabled      bool                   `tfschema:"tcp_state_tracking_enabled"`
	DnsServers                   []string               `tfschema:"dns_servers"`
	AcceleratedNetworkingEnabled bool                   `tfschema:"accelerated_networking_enabled"`
	FpgaEnabled                  bool                   `tfschema:"fpga_enabled"`
	IPForwardingEnabled          bool                   `tfschema:"ip_forwarding_enabled"`
	IPConfiguration              []IPConfigurationModel `tfschema:"ip_configuration"`
	NetworkSecurityGroupId       string                 `tfschema:"network_security_group_id"`
	Primary                      bool                   `tfschema:"primary"`
}

type IPConfigurationModel struct {
	Name                                    string                 `tfschema:"name"`
	ApplicationGatewayBackendAddressPoolIds []string               `tfschema:"application_gateway_backend_address_pool_ids"`
	ApplicationSecurityGroupIds             []string               `tfschema:"application_security_group_ids"`
	LoadBalancerBackendAddressPoolIds       []string               `tfschema:"load_balancer_backend_address_pool_ids"`
	LoadBalancerInboundNatPoolIds           []string               `tfschema:"load_balancer_inbound_nat_rules_ids"`
	Primary                                 bool                   `tfschema:"primary"`
	Version                                 string                 `tfschema:"version"`
	PublicIPAddress                         []PublicIPAddressModel `tfschema:"public_ip_address"`
	SubnetId                                string                 `tfschema:"subnet_id"`
}

type PublicIPAddressModel struct {
	Name                 string       `tfschema:"name"`
	DeleteOption         string       `tfschema:"delete_option"`
	DomainNameLabel      string       `tfschema:"domain_name_label"`
	DomainNameLabelScope string       `tfschema:"domain_name_label_scope"`
	IdleTimeoutInMinutes int64        `tfschema:"idle_timeout_in_minutes"`
	IPTags               []IPTagModel `tfschema:"ip_tags"`
	Version              string       `tfschema:"version"`
	PublicIPPrefix       string       `tfschema:"public_ip_prefix_id"`
	Sku                  []SkuModel   `tfschema:"sku"`
}

type IPTagModel struct {
	IPTagType string `tfschema:"ip_tag_type"`
	Tag       string `tfschema:"tag"`
}

type SkuModel struct {
	Name string `tfschema:"name"`
	Tier string `tfschema:"tier"`
}

type OSProfileModel struct {
	AdminPassword               string                      `tfschema:"admin_password"`
	AdminUsername               string                      `tfschema:"admin_username"`
	ExtensionOperationsEnabled  bool                        `tfschema:"extension_operations_enabled"`
	ComputerNamePrefix          string                      `tfschema:"computer_name_prefix"`
	CustomDataBase64            string                      `tfschema:"custom_data_base64"`
	LinuxConfiguration          []LinuxConfigurationModel   `tfschema:"linux_configuration"`
	RequireGuestProvisionSignal bool                        `tfschema:"require_guest_provision_signal"`
	OsProfileSecrets            []OsProfileSecretsModel     `tfschema:"os_profile_secrets"`
	WindowsConfiguration        []WindowsConfigurationModel `tfschema:"windows_configuration"`
}

type LinuxConfigurationModel struct {
	PasswordAuthenticationEnabled          bool          `tfschema:"password_authentication_enabled"`
	VMAgentPlatformUpdatesEnabled          bool          `tfschema:"vm_agent_platform_updates_enabled"`
	PatchAssessmentMode                    string        `tfschema:"patch_assessment_mode"`
	PatchBypassPlatformSafetyChecksEnabled bool          `tfschema:"patch_bypass_platform_safety_checks_enabled"`
	PatchRebootSetting                     string        `tfschema:"patch_reboot_setting"`
	PatchMode                              string        `tfschema:"patch_mode"`
	ProvisionVMAgentEnabled                bool          `tfschema:"provision_vm_agent_enabled"`
	SshKeys                                []SshKeyModel `tfschema:"ssh_keys"`
}

type SshKeyModel struct {
	KeyData string `tfschema:"key_data"`
	Path    string `tfschema:"path"`
}

type OsProfileSecretsModel struct {
	SourceVaultId     string                  `tfschema:"source_vault_id"`
	VaultCertificates []VaultCertificateModel `tfschema:"vault_certificates"`
}

type VaultCertificateModel struct {
	CertificateStore string `tfschema:"certificate_store"`
	CertificateUrl   string `tfschema:"certificate_url"`
}

type WindowsConfigurationModel struct {
	AdditionalUnattendContent                       []AdditionalUnattendContentModel `tfschema:"additional_unattend_content"`
	AutomaticUpdatesEnabled                         bool                             `tfschema:"automatic_updates_enabled"`
	VMAgentPlatformUpdatesEnabled                   bool                             `tfschema:"vm_agent_platform_updates_enabled"`
	PatchAssessmentMode                             string                           `tfschema:"patch_assessment_mode"`
	BypassPlatformSafetyChecksOnUserScheduleEnabled bool                             `tfschema:"patch_bypass_platform_safety_checks_enabled"`
	PatchRebootSetting                              string                           `tfschema:"patch_reboot_setting"`
	HotPatchingEnabled                              bool                             `tfschema:"hot_patching_enabled"`
	PatchMode                                       string                           `tfschema:"patch_mode"`
	ProvisionVMAgentEnabled                         bool                             `tfschema:"provision_vm_agent_enabled"`
	TimeZone                                        string                           `tfschema:"time_zone"`
	WinRM                                           []WinRMModel                     `tfschema:"winrm"`
}

type AdditionalUnattendContentModel struct {
	ComponentName string `tfschema:"component_name"`
	Content       string `tfschema:"content"`
	PassName      string `tfschema:"pass_name"`
	SettingName   string `tfschema:"setting_name"`
}

type WinRMModel struct {
	CertificateUrl string `tfschema:"certificate_url"`
	Protocol       string `tfschema:"protocol"`
}

type SecurityPostureReferenceModel struct {
	ExcludeExtensions []string `tfschema:"exclude_extensions"`
	Id                string   `tfschema:"id"`
	OverrideEnabled   bool     `tfschema:"override_enabled"`
}

type SecurityProfileModel struct {
	EncryptionAtHostEnabled bool              `tfschema:"encryption_at_host_enabled"`
	UserAssignedIdentityId  string            `tfschema:"user_assigned_identity_id"`
	ProxyAgent              []ProxyAgentModel `tfschema:"proxy_agent"`
	SecurityType            string            `tfschema:"security_type"`
	UefiSecureBootEnabled   bool              `tfschema:"uefi_secure_boot_enabled"`
	UefiVTpmEnabled         bool              `tfschema:"uefi_vtpm_enabled"`
}

type ProxyAgentModel struct {
	KeyIncarnationValue int64  `tfschema:"key_incarnation_value"`
	model               string `tfschema:"mode"`
}

type StorageProfileModel struct {
	StorageProfileDataDisks      []StorageProfileDataDiskModel       `tfschema:"storage_profile_data_disk"`
	DiskControllerType           string                              `tfschema:"disk_controller_type"`
	StorageProfileImageReference []StorageProfileImageReferenceModel `tfschema:"storage_profile_image_reference"`
	StorageProfileOsDisk         []StorageProfileOSDiskModel         `tfschema:"storage_profile_os_disk"`
}

type StorageProfileDataDiskModel struct {
	Caching                 string             `tfschema:"caching"`
	CreateOption            string             `tfschema:"create_option"`
	DeleteOption            string             `tfschema:"delete_option"`
	DiskIOPSReadWrite       int64              `tfschema:"disk_iops_read_write"`
	DiskMBpsReadWrite       int64              `tfschema:"disk_m_bps_read_write"`
	DiskSizeInGB            int64              `tfschema:"disk_size_in_gb"`
	Lun                     int64              `tfschema:"lun"`
	ManagedDisk             []ManagedDiskModel `tfschema:"managed_disk"`
	Name                    string             `tfschema:"name"`
	WriteAcceleratorEnabled bool               `tfschema:"write_accelerator_enabled"`
}

type ManagedDiskModel struct {
	DiskEncryptionSetId         string `tfschema:"disk_encryption_set_id"`
	SecurityDiskEncryptionSetId string `tfschema:"security_disk_encryption_set_id"`
	SecurityEncryptionType      string `tfschema:"security_encryption_type"`
	StorageAccountType          string `tfschema:"storage_account_type"`
}

type StorageProfileImageReferenceModel struct {
	CommunityGalleryImageId string `tfschema:"community_gallery_image_id"`
	Id                      string `tfschema:"id"`
	Offer                   string `tfschema:"offer"`
	Publisher               string `tfschema:"publisher"`
	SharedGalleryImageId    string `tfschema:"shared_gallery_image_id"`
	Sku                     string `tfschema:"sku"`
	Version                 string `tfschema:"version"`
}

type StorageProfileOSDiskModel struct {
	Caching                 string             `tfschema:"caching"`
	CreateOption            string             `tfschema:"create_option"`
	DeleteOption            string             `tfschema:"delete_option"`
	DiffDiskOption          string             `tfschema:"diff_disk_option"`
	DiffDiskPlacement       string             `tfschema:"diff_disk_placement"`
	DiskSizeInGB            int64              `tfschema:"disk_size_in_gb"`
	ImageUri                string             `tfschema:"image_uri"`
	ManagedDisk             []ManagedDiskModel `tfschema:"managed_disk"`
	Name                    string             `tfschema:"name"`
	OsType                  string             `tfschema:"os_type"`
	VhdContainers           []string           `tfschema:"vhd_containers"`
	WriteAcceleratorEnabled bool               `tfschema:"write_accelerator_enabled"`
}

type ComputeProfileModel struct {
	AdditionalCapabilitiesHibernationEnabled bool                         `tfschema:"additional_capabilities_ultra_ssd_enabled"`
	AdditionalCapabilitiesUltraSSDEnabled    bool                         `tfschema:"additional_capabilities_hibernation_enabled"`
	VirtualMachineProfile                    []VirtualMachineProfileModel `tfschema:"virtual_machine_profile"`
	ComputeApiVersion                        string                       `tfschema:"compute_api_version"`
	PlatformFaultDomainCount                 int64                        `tfschema:"platform_fault_domain_count"`
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
	LocalStorageInGib         []VMAttributeMinMaxDoubleModel  `tfschema:"local_storage_in_gib"`
	LocalStorageSupport       string                          `tfschema:"local_storage_support"`
	MemoryInGib               []VMAttributeMinMaxDoubleModel  `tfschema:"memory_in_gib"`
	MemoryInGibPerVCPU        []VMAttributeMinMaxDoubleModel  `tfschema:"memory_in_gib_per_vcpu"`
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

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"additional_location_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"location": commonschema.LocationWithoutForceNew(),

					"virtual_machine_profile_override": virtualMachineProfileSchema(),
				},
			},
		},

		"compute_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"compute_api_version": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"platform_fault_domain_count": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},

					"virtual_machine_profile": virtualMachineProfileSchema(),

					"additional_capabilities_hibernation_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"additional_capabilities_ultra_ssd_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
				},
			},
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

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

					"publisher": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"promotion_code": {
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

		"regular_priority_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"allocation_strategy": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(fleets.RegularPriorityAllocationStrategyLowestPrice),
							string(fleets.RegularPriorityAllocationStrategyPrioritized),
						}, false),
					},

					"capacity": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},

					"min_capacity": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},
				},
			},
		},

		"spot_priority_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"allocation_strategy": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(fleets.SpotAllocationStrategyPriceCapacityOptimized),
							string(fleets.SpotAllocationStrategyLowestPrice),
							string(fleets.SpotAllocationStrategyCapacityOptimized),
						}, false),
					},

					"capacity": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},

					"eviction_policy": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(fleets.EvictionPolicyDelete),
							string(fleets.EvictionPolicyDeallocate),
						}, false),
					},

					"maintain": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},

					"min_capacity": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},

					"max_price_per_vm": {
						Type:     pluginsdk.TypeFloat,
						Optional: true,
					},
				},
			},
		},

		// need to confirm is this a required property?
		"tags": commonschema.Tags(),

		"vm_attributes": vmAttributesSchema(),

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

func vmAttributesSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
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
						Schema: vmAttributesMaxMinSchema(pluginsdk.TypeInt),
					},
				},

				"accelerator_manufacturers": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
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
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},

				"architecture_types": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
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
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},

				"data_disk_count": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: vmAttributesMaxMinSchema(pluginsdk.TypeInt),
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

				"local_storage_in_gib": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: vmAttributesMaxMinSchema(pluginsdk.TypeFloat),
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

				"memory_in_gib": {
					Type:     pluginsdk.TypeList,
					Required: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: vmAttributesMaxMinSchema(pluginsdk.TypeFloat),
					},
				},

				"memory_in_gib_per_vcpu": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: vmAttributesMaxMinSchema(pluginsdk.TypeFloat),
					},
				},

				"network_bandwidth_in_mbps": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: vmAttributesMaxMinSchema(pluginsdk.TypeFloat),
					},
				},

				"network_interface_count": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: vmAttributesMaxMinSchema(pluginsdk.TypeInt),
					},
				},

				"rdma_network_interface_count": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: vmAttributesMaxMinSchema(pluginsdk.TypeInt),
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
						Schema: vmAttributesMaxMinSchema(pluginsdk.TypeInt),
					},
				},

				"vm_categories": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},
	}
}

func vmAttributesMaxMinSchema(inputType pluginsdk.ValueType) map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"max": {
			Type:     inputType,
			Optional: true,
		},

		"min": {
			Type:     inputType,
			Optional: true,
		},
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
			client := metadata.Client.AzureFleet.FleetsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model AzureFleetResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := fleets.NewFleetID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := fleets.Fleet{
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

			expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}
			properties.Identity = expandedIdentity

			additionalLocationsProfileValue, err := expandAdditionalLocationProfileModel(model.AdditionalLocationProfile)
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

			if err := client.CreateOrUpdateThenPoll(ctx, id, properties); err != nil {
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

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			properties := existing.Model
			if metadata.ResourceData.HasChange("identity") {
				identityValue, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(model.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				properties.Identity = identityValue
			}

			if metadata.ResourceData.HasChange("plan") {
				properties.Plan = expandPlanModel(model.Plan)
			}

			if metadata.ResourceData.HasChange("additional_location_profile") {
				additionalLocationsProfileValue, err := expandAdditionalLocationProfileModel(model.AdditionalLocationProfile)
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

			//properties.SystemData = nil

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

				identityValue, err := identity.FlattenSystemAndUserAssignedMapToModel(pointer.To((identity.SystemAndUserAssignedMap)(*model.Identity)))
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}
				state.Identity = pointer.From(identityValue)

				state.Plan = flattenPlanModel(model.Plan)

				if props := model.Properties; props != nil {
					additionalLocationsProfileValue, err := flattenAdditionalLocationProfileModel(props.AdditionalLocationsProfile)
					if err != nil {
						return err
					}
					state.AdditionalLocationProfile = additionalLocationsProfileValue

					computeProfileValue, err := flattenComputeProfileModel(&props.ComputeProfile)
					if err != nil {
						return err
					}
					state.ComputeProfile = computeProfileValue

					state.RegularPriorityProfile = flattenRegularPriorityProfileModel(props.RegularPriorityProfile)

					state.SpotPriorityProfile = flattenSpotPriorityProfileModel(props.SpotPriorityProfile)

					if props.UniqueId != nil {
						state.UniqueId = *props.UniqueId
					}

					state.VMAttributes = flattenVMAttributesModel(props.VMAttributes)

					state.VMSizesProfile = flattenVMSizeProfileModelArray(&props.VMSizesProfile)
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

func expandAdditionalLocationProfileModel(inputList []AdditionalLocationProfileModel) (*fleets.AdditionalLocationsProfile, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	output := fleets.AdditionalLocationsProfile{}
	var outputList []fleets.LocationProfile
	for _, v := range inputList {
		input := v
		output := fleets.LocationProfile{
			Location: input.Location,
		}

		virtualMachineProfileOverrideValue, err := expandBaseVirtualMachineProfileModel(input.VirtualMachineProfileOverride)
		if err != nil {
			return nil, err
		}

		output.VirtualMachineProfileOverride = virtualMachineProfileOverrideValue

		outputList = append(outputList, output)
	}

	output.LocationProfiles = outputList

	return &output, nil
}

func expandBaseVirtualMachineProfileModel(inputList []VirtualMachineProfileModel) (*fleets.BaseVirtualMachineProfile, error) {
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
		OsProfile:                expandOSProfileModel(input.OsProfile),
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

func expandSubResourceArray(inputList []string) *[]fleets.SubResource {
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
	extensionsValue, err := expandVirtualMachineScaleSetExtensionModelArray(inputList)
	if err != nil {
		return nil, err
	}
	output.Extensions = extensionsValue

	if timeBudget != "" {
		output.ExtensionsTimeBudget = pointer.To(timeBudget)
	}
	return &output, nil
}

func expandVirtualMachineScaleSetExtensionModelArray(inputList []ExtensionsModel) (*[]fleets.VirtualMachineScaleSetExtension, error) {
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
		NetworkInterfaceConfigurations: expandNetworkConfigurationArray(inputList),
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

func expandNetworkConfigurationArray(inputList []NetworkInterfaceModel) *[]fleets.VirtualMachineScaleSetNetworkConfiguration {
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
		AuxiliaryMode:               pointer.To(fleets.NetworkInterfaceAuxiliaryMode(input.AuxiliaryMode)),
		AuxiliarySku:                pointer.To(fleets.NetworkInterfaceAuxiliarySku(input.AuxiliarySku)),
		DeleteOption:                pointer.To(fleets.DeleteOptions(input.DeleteOption)),
		DisableTcpStateTracking:     pointer.To(!input.TcpStateTrackingEnabled),
		DnsSettings:                 expandNetworkConfigurationDnsSettings(input.DnsServers),
		EnableAcceleratedNetworking: pointer.To(input.AcceleratedNetworkingEnabled),
		EnableFpga:                  pointer.To(input.FpgaEnabled),
		EnableIPForwarding:          pointer.To(input.IPForwardingEnabled),
		NetworkSecurityGroup:        expandSubResource(input.NetworkSecurityGroupId),
		Primary:                     pointer.To(input.Primary),
	}

	output.IPConfigurations = pointer.From(expandIPConfigurationModelArray(input.IPConfiguration))

	return &output
}

func expandNetworkConfigurationDnsSettings(input []string) *fleets.VirtualMachineScaleSetNetworkConfigurationDnsSettings {
	if len(input) == 0 {
		return nil
	}

	return &fleets.VirtualMachineScaleSetNetworkConfigurationDnsSettings{
		DnsServers: pointer.To(input),
	}
}

func expandIPConfigurationModelArray(inputList []IPConfigurationModel) *[]fleets.VirtualMachineScaleSetIPConfiguration {
	var outputList []fleets.VirtualMachineScaleSetIPConfiguration
	for _, v := range inputList {
		input := v
		output := fleets.VirtualMachineScaleSetIPConfiguration{
			Name: input.Name,
			Properties: &fleets.VirtualMachineScaleSetIPConfigurationProperties{
				ApplicationGatewayBackendAddressPools: expandSubResourceArray(input.ApplicationGatewayBackendAddressPoolIds),
				ApplicationSecurityGroups:             expandSubResourceArray(input.ApplicationSecurityGroupIds),
				LoadBalancerBackendAddressPools:       expandSubResourceArray(input.LoadBalancerBackendAddressPoolIds),
				LoadBalancerInboundNatPools:           expandSubResourceArray(input.LoadBalancerInboundNatPoolIds),
				Primary:                               pointer.To(input.Primary),
				PrivateIPAddressVersion:               pointer.To(fleets.IPVersion(input.Version)),
				PublicIPAddressConfiguration:          expandPublicIPAddressModel(input.PublicIPAddress),
				Subnet:                                expandApiEntityReferenceModel(input.SubnetId),
			},
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
			DeleteOption:           pointer.To(fleets.DeleteOptions(input.DeleteOption)),
			DnsSettings:            expandPublicIPAddressDnsSettings(input.DomainNameLabel, input.DomainNameLabelScope),
			IPTags:                 expandIPTagModelArray(input.IPTags),
			IdleTimeoutInMinutes:   pointer.To(input.IdleTimeoutInMinutes),
			PublicIPAddressVersion: pointer.To(fleets.IPVersion(input.Version)),
			PublicIPPrefix:         expandSubResource(input.PublicIPPrefix),
		},
		Sku: expandPublicIPAddressSkuModel(input.Sku),
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

func expandIPTagModelArray(inputList []IPTagModel) *[]fleets.VirtualMachineScaleSetIPTag {
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

func expandOSProfileModel(inputList []OSProfileModel) *fleets.VirtualMachineScaleSetOSProfile {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.VirtualMachineScaleSetOSProfile{
		AdminUsername:               pointer.To(input.AdminUsername),
		AdminPassword:               pointer.To(input.AdminPassword),
		ComputerNamePrefix:          pointer.To(input.ComputerNamePrefix),
		AllowExtensionOperations:    pointer.To(input.ExtensionOperationsEnabled),
		LinuxConfiguration:          expandLinuxConfigurationModel(input.LinuxConfiguration),
		RequireGuestProvisionSignal: pointer.To(input.RequireGuestProvisionSignal),
		Secrets:                     expandVaultSecretGroupModelArray(input.OsProfileSecrets),
		WindowsConfiguration:        expandWindowsConfigurationModel(input.WindowsConfiguration),
	}

	if input.CustomDataBase64 != "" {
		output.CustomData = pointer.To(input.CustomDataBase64)
	}

	return &output
}

func expandLinuxConfigurationModel(inputList []LinuxConfigurationModel) *fleets.LinuxConfiguration {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.LinuxConfiguration{
		DisablePasswordAuthentication: pointer.To(input.PasswordAuthenticationEnabled),
		EnableVMAgentPlatformUpdates:  pointer.To(input.VMAgentPlatformUpdatesEnabled),
		PatchSettings:                 expandLinuxPatchSettingsModel(input),
		ProvisionVMAgent:              pointer.To(input.ProvisionVMAgentEnabled),
		Ssh:                           expandSshConfigurationModel(input.SshKeys),
	}

	return &output
}

func expandLinuxPatchSettingsModel(input *LinuxConfigurationModel) *fleets.LinuxPatchSettings {
	output := fleets.LinuxPatchSettings{
		AssessmentMode:              pointer.To(fleets.LinuxPatchAssessmentMode(input.PatchAssessmentMode)),
		AutomaticByPlatformSettings: expandLinuxAutomaticByPlatformSettingsModel(input.PatchBypassPlatformSafetyChecksEnabled, input.PatchRebootSetting),
		PatchMode:                   pointer.To(fleets.LinuxVMGuestPatchMode(input.PatchMode)),
	}

	return &output
}

func expandLinuxAutomaticByPlatformSettingsModel(checkEnabled bool, setting string) *fleets.LinuxVMGuestPatchAutomaticByPlatformSettings {
	output := fleets.LinuxVMGuestPatchAutomaticByPlatformSettings{
		BypassPlatformSafetyChecksOnUserSchedule: pointer.To(checkEnabled),
	}

	if setting != "" {
		output.RebootSetting = pointer.To(fleets.LinuxVMGuestPatchAutomaticByPlatformRebootSetting(setting))
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

func expandVaultSecretGroupModelArray(inputList []OsProfileSecretsModel) *[]fleets.VaultSecretGroup {
	if len(inputList) == 0 {
		return nil
	}

	var outputList []fleets.VaultSecretGroup
	for _, v := range inputList {
		input := v
		output := fleets.VaultSecretGroup{
			SourceVault:       expandSubResource(input.SourceVaultId),
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
			output.CertificateStore = pointer.To(input.CertificateStore)
		}

		if input.CertificateUrl != "" {
			output.CertificateURL = pointer.To(input.CertificateUrl)
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
		EnableAutomaticUpdates:       pointer.To(input.AutomaticUpdatesEnabled),
		EnableVMAgentPlatformUpdates: pointer.To(input.VMAgentPlatformUpdatesEnabled),
		PatchSettings:                expandWindowsPatchSettingsModel(input),
		ProvisionVMAgent:             pointer.To(input.ProvisionVMAgentEnabled),
		WinRM:                        expandWinRM(input.WinRM),
	}
	if input.TimeZone != "" {
		output.TimeZone = pointer.To(input.TimeZone)
	}

	return &output
}

func expandWindowsPatchSettingsModel(input *WindowsConfigurationModel) *fleets.PatchSettings {
	output := fleets.PatchSettings{
		AssessmentMode:              pointer.To(fleets.WindowsPatchAssessmentMode(input.PatchAssessmentMode)),
		AutomaticByPlatformSettings: expandWindowsAutomaticByPlatformSettingsModel(input.BypassPlatformSafetyChecksOnUserScheduleEnabled, input.PatchRebootSetting),
		PatchMode:                   pointer.To(fleets.WindowsVMGuestPatchMode(input.PatchMode)),
	}

	return &output
}

func expandWindowsAutomaticByPlatformSettingsModel(checkEnabled bool, setting string) *fleets.WindowsVMGuestPatchAutomaticByPlatformSettings {
	output := fleets.WindowsVMGuestPatchAutomaticByPlatformSettings{
		BypassPlatformSafetyChecksOnUserSchedule: pointer.To(checkEnabled),
	}

	if setting != "" {
		output.RebootSetting = pointer.To(fleets.WindowsVMGuestPatchAutomaticByPlatformRebootSetting(setting))
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
		Mode:             pointer.To(fleets.Mode(input.model)),
	}

	return &output
}

func expandStorageProfileModel(inputList []StorageProfileModel) *fleets.VirtualMachineScaleSetStorageProfile {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := fleets.VirtualMachineScaleSetStorageProfile{
		DataDisks:          expandDataDiskModelArray(input.StorageProfileDataDisks),
		DiskControllerType: pointer.To(fleets.DiskControllerTypes(input.DiskControllerType)),
		ImageReference:     expandImageReferenceModel(input.StorageProfileImageReference),
		OsDisk:             expandOSDiskModel(input.StorageProfileOsDisk),
	}

	return &output
}

func expandDataDiskModelArray(inputList []StorageProfileDataDiskModel) *[]fleets.VirtualMachineScaleSetDataDisk {
	var outputList []fleets.VirtualMachineScaleSetDataDisk
	for _, v := range inputList {
		input := v
		output := fleets.VirtualMachineScaleSetDataDisk{
			Caching:                 pointer.To(fleets.CachingTypes(input.Caching)),
			CreateOption:            fleets.DiskCreateOptionTypes(input.CreateOption),
			DeleteOption:            pointer.To(fleets.DiskDeleteOptionTypes(input.DeleteOption)),
			DiskIOPSReadWrite:       pointer.To(input.DiskIOPSReadWrite),
			DiskMBpsReadWrite:       pointer.To(input.DiskMBpsReadWrite),
			DiskSizeGB:              pointer.To(input.DiskSizeInGB),
			Lun:                     input.Lun,
			ManagedDisk:             expandManagedDiskModel(input.ManagedDisk),
			WriteAcceleratorEnabled: pointer.To(input.WriteAcceleratorEnabled),
		}

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
		DiskEncryptionSet:  expandDiskEncryptionSetModel(input),
		SecurityProfile:    expandVMDiskSecurityProfileModel(input),
		StorageAccountType: pointer.To(fleets.StorageAccountTypes(input.StorageAccountType)),
	}

	return &output
}

func expandDiskEncryptionSetModel(input *ManagedDiskModel) *fleets.DiskEncryptionSetParameters {
	if input == nil || input.DiskEncryptionSetId == "" {
		return nil
	}

	output := fleets.DiskEncryptionSetParameters{}
	if input.DiskEncryptionSetId != "" {
		output.Id = pointer.To(input.DiskEncryptionSetId)
	}

	return &output
}

func expandVMDiskSecurityProfileModel(input *ManagedDiskModel) *fleets.VMDiskSecurityProfile {
	if input == nil {
		return nil
	}

	output := fleets.VMDiskSecurityProfile{
		DiskEncryptionSet: expandDiskEncryptionSetModel(input),
	}

	if input.SecurityEncryptionType != "" {
		output.SecurityEncryptionType = pointer.To(fleets.SecurityEncryptionTypes(input.SecurityEncryptionType))
	}
	return &output
}

func expandImageReferenceModel(inputList []StorageProfileImageReferenceModel) *fleets.ImageReference {
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

func expandOSDiskModel(inputList []StorageProfileOSDiskModel) *fleets.VirtualMachineScaleSetOSDisk {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.VirtualMachineScaleSetOSDisk{
		Caching:                 pointer.To(fleets.CachingTypes(input.Caching)),
		CreateOption:            fleets.DiskCreateOptionTypes(input.CreateOption),
		DeleteOption:            pointer.To(fleets.DiskDeleteOptionTypes(input.DeleteOption)),
		DiffDiskSettings:        expandDiffDiskSettingsModel(input),
		DiskSizeGB:              pointer.To(input.DiskSizeInGB),
		Image:                   expandImage(input.ImageUri),
		ManagedDisk:             expandManagedDiskModel(input.ManagedDisk),
		OsType:                  pointer.To(fleets.OperatingSystemTypes(input.OsType)),
		VhdContainers:           pointer.To(input.VhdContainers),
		WriteAcceleratorEnabled: pointer.To(input.WriteAcceleratorEnabled),
	}

	if input.Name != "" {
		output.Name = pointer.To(input.Name)
	}

	return &output
}

func expandDiffDiskSettingsModel(input *StorageProfileOSDiskModel) *fleets.DiffDiskSettings {
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

func expandComputeProfileModel(inputList []ComputeProfileModel) (*fleets.ComputeProfile, error) {
	if len(inputList) == 0 {
		return nil, nil
	}
	input := &inputList[0]
	output := fleets.ComputeProfile{
		AdditionalVirtualMachineCapabilities: expandAdditionalCapabilitiesModel(input),
		PlatformFaultDomainCount:             pointer.To(input.PlatformFaultDomainCount),
	}

	baseVirtualMachineProfileValue, err := expandBaseVirtualMachineProfileModel(input.VirtualMachineProfile)
	if err != nil {
		return nil, err
	}
	output.BaseVirtualMachineProfile = pointer.From(baseVirtualMachineProfileValue)

	if input.ComputeApiVersion != "" {
		output.ComputeApiVersion = pointer.To(input.ComputeApiVersion)
	}

	return &output, nil
}

func expandAdditionalCapabilitiesModel(input *ComputeProfileModel) *fleets.AdditionalCapabilities {
	if input == nil {
		return nil
	}

	output := fleets.AdditionalCapabilities{
		HibernationEnabled: pointer.To(input.AdditionalCapabilitiesHibernationEnabled),
		UltraSSDEnabled:    pointer.To(input.AdditionalCapabilitiesUltraSSDEnabled),
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
		output.PromotionCode = pointer.To(input.PromotionCode)
	}

	if input.Version != "" {
		output.Version = pointer.To(input.Version)
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
		AllocationStrategy: pointer.To(fleets.SpotAllocationStrategy(input.AllocationStrategy)),
		Capacity:           pointer.To(input.Capacity),
		EvictionPolicy:     pointer.To(fleets.EvictionPolicy(input.EvictionPolicy)),
		Maintain:           pointer.To(input.Maintain),
		MaxPricePerVM:      pointer.To(input.MaxPricePerVM),
		MinCapacity:        pointer.To(input.MinCapacity),
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
		BurstableSupport:          pointer.To(fleets.VMAttributeSupport(input.BurstableSupport)),
		CpuManufacturers:          expandCPUManufacturerModelArray(input.CpuManufacturers),
		DataDiskCount:             expandVMAttributeMinMaxIntegerModel(input.DataDiskCount),
		ExcludedVMSizes:           pointer.To(input.ExcludedVMSizes),
		LocalStorageDiskTypes:     expandLocalStorageDiskTypeModelArray(input.LocalStorageDiskTypes),
		LocalStorageInGiB:         expandVMAttributeMinMaxDoubleModel(input.LocalStorageInGib),
		LocalStorageSupport:       pointer.To(fleets.VMAttributeSupport(input.LocalStorageSupport)),
		MemoryInGiBPerVCPU:        expandVMAttributeMinMaxDoubleModel(input.MemoryInGibPerVCPU),
		NetworkBandwidthInMbps:    expandVMAttributeMinMaxDoubleModel(input.NetworkBandwidthInMbps),
		NetworkInterfaceCount:     expandVMAttributeMinMaxIntegerModel(input.NetworkInterfaceCount),
		RdmaNetworkInterfaceCount: expandVMAttributeMinMaxIntegerModel(input.RdmaNetworkInterfaceCount),
		RdmaSupport:               pointer.To(fleets.VMAttributeSupport(input.RdmaSupport)),
		VMCategories:              expandVMCategoryModelArray(input.VMCategories),
	}

	output.MemoryInGiB = pointer.From(expandVMAttributeMinMaxDoubleModel(input.MemoryInGib))

	output.VCPUCount = pointer.From(expandVMAttributeMinMaxIntegerModel(input.VCPUCount))

	return &output
}

func expandVMAttributeMinMaxIntegerModel(inputList []VMAttributeMinMaxIntegerModel) *fleets.VMAttributeMinMaxInteger {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.VMAttributeMinMaxInteger{
		Max: pointer.To(input.Max),
		Min: pointer.To(input.Min),
	}

	return &output
}

func expandAcceleratorManufacturerModelArray(inputList []string) *[]fleets.AcceleratorManufacturer {
	if len(inputList) == 0 {
		return nil
	}

	result := make([]fleets.AcceleratorManufacturer, 0)

	for _, v := range inputList {
		if v != "" {
			result = append(result, fleets.AcceleratorManufacturer(v))
		}
	}

	return &result
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
	if len(inputList) == 0 {
		return nil
	}

	var outputList []fleets.VMSizeProfile
	for _, v := range inputList {
		input := v
		output := fleets.VMSizeProfile{
			Name: input.Name,
			Rank: pointer.To(input.Rank),
		}

		outputList = append(outputList, output)
	}
	return &outputList
}

func flattenAdditionalLocationProfileModel(input *fleets.AdditionalLocationsProfile) ([]AdditionalLocationProfileModel, error) {
	var outputList []AdditionalLocationProfileModel
	if input == nil {
		return outputList, nil
	}

	for _, input := range input.LocationProfiles {
		output := AdditionalLocationProfileModel{
			Location: input.Location,
		}
		virtualMachineProfileOverrideValue, err := flattenVirtualMachineProfileModel(input.VirtualMachineProfileOverride)
		if err != nil {
			return nil, err
		}

		output.VirtualMachineProfileOverride = virtualMachineProfileOverrideValue
		outputList = append(outputList, output)
	}

	output := AdditionalLocationProfileModel{}

	return append(outputList, output), nil
}

func flattenVirtualMachineProfileModel(input *fleets.BaseVirtualMachineProfile) ([]VirtualMachineProfileModel, error) {
	var outputList []VirtualMachineProfileModel
	if input == nil {
		return outputList, nil
	}
	output := VirtualMachineProfileModel{
		GalleryApplicationProfile:        flattenApplicationProfileModel(input.ApplicationProfile),
		VMSize:                           flattenVirtualMachineScaleSetHardwareProfileModel(input.HardwareProfile),
		NetworkInterface:                 flattenNetworkInterfaceModel(input.NetworkProfile),
		OsProfile:                        flattenVirtualMachineScaleSetOSProfileModel(input.OsProfile),
		ScheduledEventTerminationEnabled: flattenScheduledEventsProfileModel(input.ScheduledEventsProfile),
		ScheduledEventTerminationTimeout: flattenScheduledEventsProfileModel(input.ScheduledEventsProfile),
		ScheduledEventOsImageEnabled:     flattenScheduledEventsProfileModel(input.ScheduledEventsProfile),
		ScheduledEventOsImageTimeout:     flattenScheduledEventsProfileModel(input.ScheduledEventsProfile),
		SecurityPostureReference:         flattenSecurityPostureReferenceModel(input.SecurityPostureReference),
		SecurityProfile:                  flattenSecurityProfileModel(input.SecurityProfile),
		ServiceArtifactId:                flattenServiceArtifactReferenceModel(input.ServiceArtifactReference),
		StorageProfile:                   flattenVirtualMachineScaleSetStorageProfileModel(input.StorageProfile),
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

	extensionProfileValue, err := flattenVirtualMachineScaleSetExtensionProfileModel(input.ExtensionProfile)
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

func flattenCapacityReservationProfileModel(input *fleets.CapacityReservationProfile) string {
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
		output.AutoUpgradeMinorVersionEnabled = *input.AutoUpgradeMinorVersion
	}

	if input.EnableAutomaticUpgrade != nil {
		output.EnableAutomaticUpgrade = *input.EnableAutomaticUpgrade
	}

	if input.ForceUpdateTag != nil {
		output.ForceUpdateTag = *input.ForceUpdateTag
	}

	if input.ProtectedSettingsJson != nil && *input.ProtectedSettingsJson != nil {

		protectedSettingsValue, err := json.Marshal(*input.ProtectedSettingsJson)
		if err != nil {
			return outputList, err
		}

		output.ProtectedSettingsJson = string(protectedSettingsValue)
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

	if input.SettingsJson != nil && *input.SettingsJson != nil {

		settingsValue, err := json.Marshal(*input.SettingsJson)
		if err != nil {
			return outputList, err
		}

		output.SettingsJson = string(settingsValue)
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
		output.VCPUAvailableCount = *input.VCPUsAvailable
	}

	if input.VCPUsPerCore != nil {
		output.VCPUPerCoreCount = *input.VCPUsPerCore
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
			output.AuxiliaryMode = string(pointer.From(props.AuxiliaryMode))
			output.AuxiliarySku = string(pointer.From(props.AuxiliarySku))

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

func flattenIPConfigurationModel(inputList []fleets.VirtualMachineScaleSetIPConfiguration) []IPConfigurationModel {
	var outputList []IPConfigurationModel
	if len(inputList) == 0 {
		return outputList
	}
	for _, input := range inputList {
		output := IPConfigurationModel{
			Name:            input.Name,
			PublicIPAddress: flattenVirtualMachineScaleSetPublicIPAddressConfigurationModel(input.PublicIPAddressConfiguration),
		}

		if props := input.Properties; props != nil {
			output.Primary = pointer.From(props.Primary)
			output.Version = string(pointer.From(props.PrivateIPAddressVersion))
			if v := props.Subnet; v != nil {
				output.SubnetId = pointer.From(v.Id)
			}

			if v := props.ApplicationGatewayBackendAddressPools; v != nil {
				output.ApplicationGatewayBackendAddressPoolIds = flattenSubResource(*v)
			}
			if v := props.ApplicationSecurityGroups; v != nil {
				output.ApplicationSecurityGroupIds = flattenSubResource(*v)
			}
			if v := props.ApplicationGatewayBackendAddressPools; v != nil {
				output.LoadBalancerBackendAddressPoolIds = flattenSubResource(*v)
			}
			if v := props.LoadBalancerInboundNatPools; v != nil {
				output.LoadBalancerInboundNatPoolIds = flattenSubResource(*v)
			}
		}

		outputList = append(outputList, output)
	}
	return outputList
}

func flattenSubResource(inputList []fleets.SubResource) []string {
	var outputList []string
	if len(inputList) == 0 {
		return outputList
	}
	for _, input := range inputList {
		output := pointer.From(input.Id)
		outputList = append(outputList, output)
	}
	return outputList
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

func flattenVirtualMachineScaleSetNetworkConfigurationPropertiesModel(input *fleets.VirtualMachineScaleSetNetworkConfigurationProperties) []NetworkConfigurationModel {
	var outputList []NetworkConfigurationModel
	if input == nil {
		return outputList
	}
	output := NetworkConfigurationModel{
		DnsSettings:          flattenVirtualMachineScaleSetNetworkConfigurationDnsSettingsModel(input.DnsSettings),
		IPConfiguration:      flattenIPConfigurationModelArray(&input.IPConfigurations),
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

func flattenIPConfigurationModelArray(inputList *[]fleets.VirtualMachineScaleSetIPConfiguration) []VirtualMachineScaleSetIPConfigurationModel {
	var outputList []VirtualMachineScaleSetIPConfigurationModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := VirtualMachineScaleSetIPConfigurationModel{
			Name:       input.Name,
			Properties: flattenIPConfigurationModel(input.Properties),
		}

		outputList = append(outputList, output)
	}
	return outputList
}

func flattenIPConfigurationModel(input *fleets.VirtualMachineScaleSetIPConfigurationProperties) []VirtualMachineScaleSetIPConfigurationModel {
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

	if input.CustomDataBase64 != nil {
		output.CustomDataBase64 = *input.CustomDataBase64
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

	if input.ProvisionVMAgentEnabled != nil {
		output.ProvisionVMAgentEnabled = *input.ProvisionVMAgentEnabled
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

	if input.ProvisionVMAgentEnabled != nil {
		output.ProvisionVMAgentEnabled = *input.ProvisionVMAgentEnabled
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
		AdditionalCapabilities: flattenAdditionalCapabilitiesModel(input.AdditionalVirtualMachineCapabilities),
	}
	baseVirtualMachineProfileValue, err := flattenVirtualMachineProfileModel(&input.BaseVirtualMachineProfile)
	if err != nil {
		return nil, err
	}

	output.VirtualMachineProfile = baseVirtualMachineProfileValue

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

func flattenVirtualMachineScaleSetExtensionModel(input *fleets.VirtualMachineScaleSetExtensionProperties) ([]ExtensionsModel, error) {
	var outputList []ExtensionsModel
	if input == nil {
		return outputList, nil
	}
	output := ExtensionsModel{
		ProtectedSettingsFromKeyVault: flattenKeyVaultSecretReferenceModel(input.ProtectedSettingsFromKeyVault),
	}
	if input.AutoUpgradeMinorVersion != nil {
		output.AutoUpgradeMinorVersionEnabled = *input.AutoUpgradeMinorVersion
	}

	if input.EnableAutomaticUpgrade != nil {
		output.AutomaticUpgradeEnabled = *input.EnableAutomaticUpgrade
	}

	if input.ForceUpdateTag != nil {
		output.ForceUpdateTag = *input.ForceUpdateTag
	}

	if input.ProtectedSettings != nil && *input.ProtectedSettings != nil {

		protectedSettingsValue, err := json.Marshal(*input.ProtectedSettings)
		if err != nil {
			return outputList, err
		}

		output.ProtectedSettingsJson = string(protectedSettingsValue)
	}

	if input.ProvisionAfterExtensions != nil {
		output.ProvisionAfterExtensions = *input.ProvisionAfterExtensions
	}

	if input.Publisher != nil {
		output.Publisher = *input.Publisher
	}

	if input.Settings != nil && *input.Settings != nil {

		settingsValue, err := json.Marshal(*input.Settings)
		if err != nil {
			return outputList, err
		}

		output.SettingsJson = string(settingsValue)
	}

	if input.SuppressFailures != nil {
		output.SuppressFailuresEnabled = *input.SuppressFailures
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
		output.VCPUAvailableCount = *input.VCPUsAvailable
	}

	if input.VCPUsPerCore != nil {
		output.VCPUPerCoreCount = *input.VCPUsPerCore
	}

	return append(outputList, output)
}

func flattenNetworkInterfaceModel(input *fleets.VirtualMachineScaleSetNetworkProfile) []NetworkProfileModel {
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

func flattenVirtualMachineScaleSetNetworkConfigurationPropertiesModel(input *fleets.VirtualMachineScaleSetNetworkConfigurationProperties) []NetworkConfigurationModel {
	var outputList []NetworkConfigurationModel
	if input == nil {
		return outputList
	}
	output := NetworkConfigurationModel{
		DnsSettings:          flattenVirtualMachineScaleSetNetworkConfigurationDnsSettingsModel(input.DnsSettings),
		IPConfiguration:      flattenIPConfigurationModelArray(&input.IPConfigurations),
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

func flattenIPConfigurationModelArray(inputList *[]fleets.VirtualMachineScaleSetIPConfiguration) []VirtualMachineScaleSetIPConfigurationModel {
	var outputList []VirtualMachineScaleSetIPConfigurationModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := VirtualMachineScaleSetIPConfigurationModel{
			Name:       input.Name,
			Properties: flattenIPConfigurationModel(input.Properties),
		}

		outputList = append(outputList, output)
	}
	return outputList
}

func flattenIPConfigurationModel(input *fleets.VirtualMachineScaleSetIPConfigurationProperties) []VirtualMachineScaleSetIPConfigurationModel {
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

	if input.CustomDataBase64 != nil {
		output.CustomDataBase64 = *input.CustomDataBase64
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

	if input.ProvisionVMAgentEnabled != nil {
		output.ProvisionVMAgentEnabled = *input.ProvisionVMAgentEnabled
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

	if input.ProvisionVMAgentEnabled != nil {
		output.ProvisionVMAgentEnabled = *input.ProvisionVMAgentEnabled
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

	output.PromotionCode = pointer.From(input.PromotionCode)
	output.Version = pointer.From(input.Version)

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
