package azurefleet

import (
	"context"
	"fmt"
	"regexp"
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
	AdditionalUnattendContent              []AdditionalUnattendContentModel `tfschema:"additional_unattend_content"`
	AutomaticUpdatesEnabled                bool                             `tfschema:"automatic_updates_enabled"`
	VMAgentPlatformUpdatesEnabled          bool                             `tfschema:"vm_agent_platform_updates_enabled"`
	PatchAssessmentMode                    string                           `tfschema:"patch_assessment_mode"`
	PatchBypassPlatformSafetyChecksEnabled bool                             `tfschema:"patch_bypass_platform_safety_checks_enabled"`
	PatchMode                              string                           `tfschema:"patch_mode"`
	PatchRebootSetting                     string                           `tfschema:"patch_reboot_setting"`
	HotPatchingEnabled                     bool                             `tfschema:"hot_patching_enabled"`
	ProvisionVMAgentEnabled                bool                             `tfschema:"provision_vm_agent_enabled"`
	TimeZone                               string                           `tfschema:"time_zone"`
	WinRM                                  []WinRMModel                     `tfschema:"winrm"`
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
	mode                string `tfschema:"mode"`
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
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[^_\\W][\\w\\-._]{0,62}\\w$"),
				"Azure resource names cannot",
				// Azure resource names cannot contain special characters /""[]:|<>+=;,?*@&, whitespace, or begin with '_' or end with '.' or '-'
				// The name cannot begin or end with the hyphen '-' character.
				// This field must be between 1 and 64 characters.
				//Linux VM names may only contain letters, numbers, '.', and '-'.
			),
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
						Type: pluginsdk.TypeInt,
						//Required: true,
						Optional: true,
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
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},

				"local_storage_disk_types": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
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

			properties.Properties.VMSizesProfile = pointer.From(expandVMSizeProfileModel(model.VMSizesProfile))

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
				properties.Properties.VMSizesProfile = pointer.From(expandVMSizeProfileModel(model.VMSizesProfile))
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
					additionalLocationsProfileValue, err := flattenAdditionalLocationProfileModel(props.AdditionalLocationsProfile, metadata)
					if err != nil {
						return err
					}
					state.AdditionalLocationProfile = additionalLocationsProfileValue

					computeProfileValue, err := flattenComputeProfileModel(&props.ComputeProfile, metadata)
					if err != nil {
						return err
					}
					state.ComputeProfile = computeProfileValue
					state.RegularPriorityProfile = flattenRegularPriorityProfileModel(props.RegularPriorityProfile)
					state.SpotPriorityProfile = flattenSpotPriorityProfileModel(props.SpotPriorityProfile)
					state.UniqueId = pointer.From(props.UniqueId)
					state.VMAttributes = flattenVMAttributesModel(props.VMAttributes)
					state.VMSizesProfile = flattenVMSizeProfileModel(&props.VMSizesProfile)
				}
				state.Tags = pointer.From(model.Tags)
				state.Zones = pointer.From(model.Zones)
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
		AllocationStrategy: pointer.To(fleets.RegularPriorityAllocationStrategy(input.AllocationStrategy)),
		Capacity:           pointer.To(input.Capacity),
		MinCapacity:        pointer.To(input.MinCapacity),
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
		AcceleratorManufacturers:  expandAcceleratorManufacturers(input.AcceleratorManufacturers),
		AcceleratorSupport:        pointer.To(fleets.VMAttributeSupport(input.AcceleratorSupport)),
		AcceleratorTypes:          expandAcceleratorTypes(input.AcceleratorTypes),
		ArchitectureTypes:         expandArchitectureTypes(input.ArchitectureTypes),
		BurstableSupport:          pointer.To(fleets.VMAttributeSupport(input.BurstableSupport)),
		CpuManufacturers:          expandCPUManufacturers(input.CpuManufacturers),
		DataDiskCount:             expandVMAttributeMinMaxIntegerModel(input.DataDiskCount),
		ExcludedVMSizes:           pointer.To(input.ExcludedVMSizes),
		LocalStorageDiskTypes:     expandLocalStorageDiskTypes(input.LocalStorageDiskTypes),
		LocalStorageInGiB:         expandVMAttributeMinMaxDoubleModel(input.LocalStorageInGib),
		LocalStorageSupport:       pointer.To(fleets.VMAttributeSupport(input.LocalStorageSupport)),
		MemoryInGiBPerVCPU:        expandVMAttributeMinMaxDoubleModel(input.MemoryInGibPerVCPU),
		NetworkBandwidthInMbps:    expandVMAttributeMinMaxDoubleModel(input.NetworkBandwidthInMbps),
		NetworkInterfaceCount:     expandVMAttributeMinMaxIntegerModel(input.NetworkInterfaceCount),
		RdmaNetworkInterfaceCount: expandVMAttributeMinMaxIntegerModel(input.RdmaNetworkInterfaceCount),
		RdmaSupport:               pointer.To(fleets.VMAttributeSupport(input.RdmaSupport)),
		VMCategories:              expandVMCategorys(input.VMCategories),
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

func expandVMCategorys(inputList []string) *[]fleets.VMCategory {
	if len(inputList) == 0 {
		return nil
	}

	var outputList []fleets.VMCategory
	for _, v := range inputList {
		if v != "" {
			outputList = append(outputList, fleets.VMCategory(v))
		}
	}
	return &outputList
}

func expandLocalStorageDiskTypes(inputList []string) *[]fleets.LocalStorageDiskType {
	if len(inputList) == 0 {
		return nil
	}

	var outputList []fleets.LocalStorageDiskType
	for _, v := range inputList {
		if v != "" {
			outputList = append(outputList, fleets.LocalStorageDiskType(v))
		}
	}
	return &outputList
}

func expandAcceleratorManufacturers(inputList []string) *[]fleets.AcceleratorManufacturer {
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

func expandAcceleratorTypes(inputList []string) *[]fleets.AcceleratorType {
	if len(inputList) == 0 {
		return nil
	}

	var outputList []fleets.AcceleratorType
	for _, v := range inputList {
		if v != "" {
			outputList = append(outputList, fleets.AcceleratorType(v))
		}
	}
	return &outputList
}

func expandArchitectureTypes(inputList []string) *[]fleets.ArchitectureType {
	if len(inputList) == 0 {
		return nil
	}

	var outputList []fleets.ArchitectureType
	for _, v := range inputList {
		if v != "" {
			outputList = append(outputList, fleets.ArchitectureType(v))
		}
	}
	return &outputList
}

func expandCPUManufacturers(inputList []string) *[]fleets.CPUManufacturer {
	if len(inputList) == 0 {
		return nil
	}

	var outputList []fleets.CPUManufacturer
	for _, v := range inputList {
		if v != "" {
			outputList = append(outputList, fleets.CPUManufacturer(v))
		}
	}
	return &outputList
}

func expandVMSizeProfileModel(inputList []VMSizeProfileModel) *[]fleets.VMSizeProfile {
	if len(inputList) == 0 {
		return nil
	}

	var outputList []fleets.VMSizeProfile
	for _, v := range inputList {
		input := v
		output := fleets.VMSizeProfile{
			Name: input.Name,
		}
		if input.Rank > 0 {
			output.Rank = pointer.To(input.Rank)
		}
		outputList = append(outputList, output)
	}
	return &outputList
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

func flattenAdditionalLocationProfileModel(input *fleets.AdditionalLocationsProfile, metadata sdk.ResourceMetaData) ([]AdditionalLocationProfileModel, error) {
	var outputList []AdditionalLocationProfileModel
	if input == nil {
		return outputList, nil
	}

	for _, input := range input.LocationProfiles {
		output := AdditionalLocationProfileModel{
			Location: input.Location,
		}
		virtualMachineProfileOverrideValue, err := flattenVirtualMachineProfileModel(input.VirtualMachineProfileOverride, metadata)
		if err != nil {
			return nil, err
		}

		output.VirtualMachineProfileOverride = virtualMachineProfileOverrideValue
		outputList = append(outputList, output)
	}

	output := AdditionalLocationProfileModel{}

	return append(outputList, output), nil
}

func flattenComputeProfileModel(input *fleets.ComputeProfile, metadata sdk.ResourceMetaData) ([]ComputeProfileModel, error) {
	var outputList []ComputeProfileModel
	if input == nil {
		return outputList, nil
	}

	output := ComputeProfileModel{}
	if v := input.AdditionalVirtualMachineCapabilities; v != nil {
		output.AdditionalCapabilitiesHibernationEnabled = pointer.From(v.HibernationEnabled)
		output.AdditionalCapabilitiesUltraSSDEnabled = pointer.From(v.UltraSSDEnabled)
	}

	baseVirtualMachineProfileValue, err := flattenVirtualMachineProfileModel(&input.BaseVirtualMachineProfile, metadata)
	if err != nil {
		return nil, err
	}
	output.VirtualMachineProfile = baseVirtualMachineProfileValue
	output.ComputeApiVersion = pointer.From(input.ComputeApiVersion)
	output.PlatformFaultDomainCount = pointer.From(input.PlatformFaultDomainCount)

	return append(outputList, output), nil
}

func flattenExtensionModel(input *fleets.VirtualMachineScaleSetExtensionProfile, metadata sdk.ResourceMetaData) ([]ExtensionsModel, error) {
	var outputList []ExtensionsModel
	if input == nil || input.Extensions == nil {
		return outputList, nil
	}

	output := ExtensionsModel{}
	for _, input := range *input.Extensions {
		output := ExtensionsModel{}
		if input.Name != nil {
			output.Name = pointer.From(input.Name)
		}

		if props := input.Properties; props != nil {
			output.Publisher = pointer.From(props.Publisher)
			output.Type = pointer.From(props.Type)
			output.TypeHandlerVersion = pointer.From(props.TypeHandlerVersion)
			output.AutoUpgradeMinorVersionEnabled = pointer.From(props.EnableAutomaticUpgrade)
			output.AutomaticUpgradeEnabled = pointer.From(props.EnableAutomaticUpgrade)
			output.ForceUpdateTag = pointer.From(props.ForceUpdateTag)
			// protected_settings_json isn't returned, so we get it from state otherwise set to empty string
			var model ExtensionsModel
			err := metadata.Decode(&model)
			if err != nil {
				return nil, err
			}
			if model.ProtectedSettingsJson != "" {
				output.ProtectedSettingsJson = model.ProtectedSettingsJson
			}
			output.ProtectedSettingsFromKeyVault = flattenProtectedSettingsFromKeyVaultModel(props.ProtectedSettingsFromKeyVault)
			output.ProvisionAfterExtensions = pointer.From(props.ProvisionAfterExtensions)
			var setting string
			if props.Settings != nil {
				setting, err = pluginsdk.FlattenJsonToString(*props.Settings)
				if err != nil {
					return nil, fmt.Errorf("flatenning `settings`: %+v", err)
				}
			}
			output.SettingsJson = setting
			output.SuppressFailuresEnabled = pointer.From(props.SuppressFailures)
		}

		outputList = append(outputList, output)
	}

	return append(outputList, output), nil
}

func flattenProtectedSettingsFromKeyVaultModel(input *fleets.KeyVaultSecretReference) []ProtectedSettingsFromKeyVaultModel {
	var outputList []ProtectedSettingsFromKeyVaultModel
	if input == nil {
		return outputList
	}

	output := ProtectedSettingsFromKeyVaultModel{
		SecretUrl:     input.SecretURL,
		SourceVaultId: pointer.From(input.SourceVault.Id),
	}

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
			if v := props.ApplicationGatewayBackendAddressPools; v != nil {
				output.ApplicationGatewayBackendAddressPoolIds = flattenSubResourceId(*v)
			}
			if v := props.ApplicationSecurityGroups; v != nil {
				output.ApplicationSecurityGroupIds = flattenSubResourceId(*v)
			}
			if v := props.ApplicationGatewayBackendAddressPools; v != nil {
				output.LoadBalancerBackendAddressPoolIds = flattenSubResourceId(*v)
			}
			if v := props.LoadBalancerInboundNatPools; v != nil {
				output.LoadBalancerInboundNatPoolIds = flattenSubResourceId(*v)
			}
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

func flattenSubResourceId(inputList []fleets.SubResource) []string {
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

func flattenPublicIPAddressModel(input *fleets.VirtualMachineScaleSetPublicIPAddressConfiguration) []PublicIPAddressModel {
	var outputList []PublicIPAddressModel
	if input == nil {
		return outputList
	}
	output := PublicIPAddressModel{
		Name: input.Name,
		Sku:  flattenPublicIPAddressSkuModel(input.Sku),
	}

	if props := input.Properties; props != nil {
		output.DeleteOption = string(pointer.From(props.DeleteOption))
		if v := props.DnsSettings; v != nil {
			output.DomainNameLabel = v.DomainNameLabel
			output.DomainNameLabelScope = string(pointer.From(v.DomainNameLabelScope))
		}
		output.IdleTimeoutInMinutes = pointer.From(props.IdleTimeoutInMinutes)
		output.Version = string(pointer.From(props.PublicIPAddressVersion))

		if v := props.IPTags; v != nil {
			output.IPTags = flattenIPTagModel(v)
		}
	}
	return append(outputList, output)
}

func flattenPublicIPAddressSkuModel(input *fleets.PublicIPAddressSku) []SkuModel {
	var outputList []SkuModel
	if input == nil {
		return outputList
	}
	output := SkuModel{}
	output.Name = string(pointer.From(input.Name))
	output.Tier = string(pointer.From(input.Tier))

	return append(outputList, output)
}

func flattenIPTagModel(inputList *[]fleets.VirtualMachineScaleSetIPTag) []IPTagModel {
	var outputList []IPTagModel
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := IPTagModel{}

		output.IPTagType = pointer.From(input.IPTagType)
		output.Tag = pointer.From(input.Tag)
		outputList = append(outputList, output)
	}
	return outputList
}

func flattenRegularPriorityProfileModel(input *fleets.RegularPriorityProfile) []RegularPriorityProfileModel {
	var outputList []RegularPriorityProfileModel
	if input == nil {
		return outputList
	}

	output := RegularPriorityProfileModel{}
	output.AllocationStrategy = string(pointer.From(input.AllocationStrategy))
	output.Capacity = pointer.From(input.Capacity)
	output.MinCapacity = pointer.From(input.MinCapacity)

	return append(outputList, output)
}

func flattenSpotPriorityProfileModel(input *fleets.SpotPriorityProfile) []SpotPriorityProfileModel {
	var outputList []SpotPriorityProfileModel
	if input == nil {
		return outputList
	}

	output := SpotPriorityProfileModel{}
	output.AllocationStrategy = string(pointer.From(input.AllocationStrategy))
	output.Capacity = pointer.From(input.Capacity)
	output.EvictionPolicy = string(pointer.From(input.EvictionPolicy))
	output.Maintain = pointer.From(input.Maintain)
	output.MaxPricePerVM = pointer.From(input.MaxPricePerVM)
	output.MinCapacity = pointer.From(input.MinCapacity)

	return append(outputList, output)
}

func flattenVMAttributesModel(input *fleets.VMAttributes) []VMAttributesModel {
	var outputList []VMAttributesModel
	if input == nil {
		return outputList
	}
	output := VMAttributesModel{
		AcceleratorCount:          flattenVMAttributeMinMaxIntegerModel(input.AcceleratorCount),
		AcceleratorManufacturers:  flattenToStringSlice(input.AcceleratorManufacturers),
		AcceleratorTypes:          flattenToStringSlice(input.AcceleratorTypes),
		ArchitectureTypes:         flattenToStringSlice(input.ArchitectureTypes),
		CpuManufacturers:          flattenToStringSlice(input.CpuManufacturers),
		DataDiskCount:             flattenVMAttributeMinMaxIntegerModel(input.DataDiskCount),
		LocalStorageDiskTypes:     flattenToStringSlice(input.LocalStorageDiskTypes),
		LocalStorageInGib:         flattenVMAttributeMinMaxDoubleModel(input.LocalStorageInGiB),
		MemoryInGib:               flattenVMAttributeMinMaxDoubleModel(&input.MemoryInGiB),
		MemoryInGibPerVCPU:        flattenVMAttributeMinMaxDoubleModel(input.MemoryInGiBPerVCPU),
		NetworkBandwidthInMbps:    flattenVMAttributeMinMaxDoubleModel(input.NetworkBandwidthInMbps),
		NetworkInterfaceCount:     flattenVMAttributeMinMaxIntegerModel(input.NetworkInterfaceCount),
		RdmaNetworkInterfaceCount: flattenVMAttributeMinMaxIntegerModel(input.RdmaNetworkInterfaceCount),
		VCPUCount:                 flattenVMAttributeMinMaxIntegerModel(&input.VCPUCount),
		VMCategories:              flattenToStringSlice(input.VMCategories),
	}

	output.AcceleratorSupport = string(pointer.From(input.AcceleratorSupport))
	output.BurstableSupport = string(pointer.From(input.BurstableSupport))
	output.ExcludedVMSizes = pointer.From(input.ExcludedVMSizes)
	output.LocalStorageSupport = string(pointer.From(input.LocalStorageSupport))
	output.RdmaSupport = string(pointer.From(input.RdmaSupport))

	return append(outputList, output)
}

func flattenVMAttributeMinMaxIntegerModel(input *fleets.VMAttributeMinMaxInteger) []VMAttributeMinMaxIntegerModel {
	var outputList []VMAttributeMinMaxIntegerModel
	if input == nil {
		return outputList
	}
	output := VMAttributeMinMaxIntegerModel{}
	output.Max = pointer.From(input.Max)
	output.Min = pointer.From(input.Min)

	return append(outputList, output)
}

func flattenVMAttributeMinMaxDoubleModel(input *fleets.VMAttributeMinMaxDouble) []VMAttributeMinMaxDoubleModel {
	var outputList []VMAttributeMinMaxDoubleModel
	if input == nil {
		return outputList
	}
	output := VMAttributeMinMaxDoubleModel{}
	output.Max = pointer.From(input.Max)
	output.Min = pointer.From(input.Min)

	return append(outputList, output)
}

func flattenToStringSlice[T any](inputList *[]T) []string {
	var outputList []string
	if inputList == nil {
		return outputList
	}

	result := make([]string, len(*inputList))
	for i, v := range *inputList {
		result[i] = fmt.Sprintf("%v", v)
	}

	return result
}

func flattenVMSizeProfileModel(inputList *[]fleets.VMSizeProfile) []VMSizeProfileModel {
	var outputList []VMSizeProfileModel
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := VMSizeProfileModel{
			Name: input.Name,
		}
		output.Rank = pointer.From(input.Rank)
		outputList = append(outputList, output)
	}
	return outputList
}
