---
subcategory: "Azure Fleet"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_azure_fleet"
description: |-
  Manages an Azure Fleet.
---

# azurerm_azure_fleet

Manages an Azure Fleet.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "East US"
}

resource "azurerm_azure_fleet" "example" {
  name                = "example-aff"
  resource_group_name = azurerm_resource_group.example.name
  location            = "West Europe"
  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = []
  }
  compute_profile {
    compute_api_version         = ""
    platform_fault_domain_count = 0
    base_virtual_machine_profile {
      license_type = ""
      user_data    = ""
      application_profile {
        gallery_applications {
          configuration_reference             = ""
          enable_automatic_upgrade            = false
          order                               = 0
          package_reference_id                = ""
          tags                                = ""
          treat_failure_as_deployment_failure = false
        }
      }
      capacity_reservation {
        capacity_reservation_group {
          id = ""
        }
      }
      diagnostics_profile {
        boot_diagnostics {
          enabled     = false
          storage_uri = ""
        }
      }
      extension_profile {
        extensions_time_budget = ""
        extensions {
          name = ""
          properties {
            auto_upgrade_minor_version = false
            enable_automatic_upgrade   = false
            force_update_tag           = ""
            publisher                  = ""
            suppress_failures          = false
            type                       = ""
            type_handler_version       = ""
            provision_after_extensions = []
            protected_settings_from_key_vault {
              secret_url = ""
              source_vault {
                id = ""
              }
            }
            protected_settings = jsonencode({
              "key" : "value"
            })
            settings = jsonencode({
              "key" : "value"
            })
          }
        }
      }
      hardware_profile {
        vm_size_properties {
          vcp_us_available = 0
          vcp_us_per_core  = 0
        }
      }
      network_profile {
        network_api_version = ""
        health_probe {
          id = ""
        }
        network_interface_configurations {
          name = ""
          properties {
            auxiliary_mode                = ""
            auxiliary_sku                 = ""
            delete_option                 = ""
            disable_tcp_state_tracking    = false
            enable_accelerated_networking = false
            enable_fpga                   = false
            enable_ip_forwarding          = false
            primary                       = false
            dns_settings {
              dns_servers = []
            }
            ip_configurations {
              name = ""
              properties {
                primary                    = false
                private_ip_address_version = ""
                application_gateway_backend_address_pools {
                  id = ""
                }
                application_security_groups {
                  id = ""
                }
                load_balancer_backend_address_pools {
                  id = ""
                }
                load_balancer_inbound_nat_pools {
                  id = ""
                }
                public_ip_address_configuration {
                  name = ""
                  properties {
                    delete_option             = ""
                    idle_timeout_in_minutes   = 0
                    public_ip_address_version = ""
                    dns_settings {
                      domain_name_label       = ""
                      domain_name_label_scope = ""
                    }
                    ip_tags {
                      ip_tag_type = ""
                      tag         = ""
                    }
                    public_ip_prefix {
                      id = ""
                    }
                  }
                  sku {
                    name = ""
                    tier = ""
                  }
                }
                subnet {
                  id = ""
                }
              }
            }
            network_security_group {
              id = ""
            }
          }
        }
      }
      os_profile {
        admin_password                 = ""
        admin_username                 = ""
        allow_extension_operations     = false
        computer_name_prefix           = ""
        custom_data                    = ""
        require_guest_provision_signal = false
        linux_configuration {
          disable_password_authentication  = false
          enable_vm_agent_platform_updates = false
          provision_vm_agent               = false
          patch_settings {
            assessment_mode = ""
            patch_mode      = ""
            automatic_by_platform_settings {
              bypass_platform_safety_checks_on_user_schedule = false
              reboot_setting                                 = ""
            }
          }
          ssh {
            public_keys {
              key_data = ""
              path     = ""
            }
          }
        }
        secrets {
          source_vault {
            id = ""
          }
          vault_certificates {
            certificate_store = ""
            certificate_url   = ""
          }
        }
        windows_configuration {
          enable_automatic_updates         = false
          enable_vm_agent_platform_updates = false
          provision_vm_agent               = false
          time_zone                        = ""
          additional_unattend_content {
            component_name = ""
            content        = ""
            pass_name      = ""
            setting_name   = ""
          }
          patch_settings {
            assessment_mode    = ""
            enable_hotpatching = false
            patch_mode         = ""
            automatic_by_platform_settings {
              bypass_platform_safety_checks_on_user_schedule = false
              reboot_setting                                 = ""
            }
          }
          win_rm {
            listeners {
              certificate_url = ""
              protocol        = ""
            }
          }
        }
      }
      scheduled_events_profile {
        os_image_notification_profile {
          enable             = false
          not_before_timeout = ""
        }
        terminate_notification_profile {
          enable             = false
          not_before_timeout = ""
        }
      }
      security_posture_reference {
        id                 = ""
        is_overridable     = false
        exclude_extensions = []
      }
      security_profile {
        encryption_at_host = false
        security_type      = ""
        encryption_identity {
          user_assigned_identity_resource_id = ""
        }
        proxy_agent_settings {
          enabled            = false
          key_incarnation_id = 0
          mode               = ""
        }
        uefi_settings {
          secure_boot_enabled = false
          v_tpm_enabled       = false
        }
      }
      service_artifact_reference {
        id = ""
      }
      storage_profile {
        disk_controller_type = ""
        data_disks {
          caching                   = ""
          create_option             = ""
          delete_option             = ""
          disk_iops_read_write      = 0
          disk_m_bps_read_write     = 0
          disk_size_gb              = 0
          lun                       = 0
          name                      = ""
          write_accelerator_enabled = false
          managed_disk {
            storage_account_type = ""
            disk_encryption_set {
              id = ""
            }
            security_profile {
              security_encryption_type = ""
              disk_encryption_set {
                id = ""
              }
            }
          }
        }
        image_reference {
          community_gallery_image_id = ""
          id                         = ""
          offer                      = ""
          publisher                  = ""
          shared_gallery_image_id    = ""
          sku                        = ""
          version                    = ""
        }
        os_disk {
          caching                   = ""
          create_option             = ""
          delete_option             = ""
          disk_size_gb              = 0
          name                      = ""
          os_type                   = ""
          write_accelerator_enabled = false
          vhd_containers            = []
          diff_disk_settings {
            option    = ""
            placement = ""
          }
          image {
            uri = ""
          }
          managed_disk {
            storage_account_type = ""
            disk_encryption_set {
              id = ""
            }
            security_profile {
              security_encryption_type = ""
              disk_encryption_set {
                id = ""
              }
            }
          }
        }
      }
    }
  }
  plan {
    name           = ""
    product        = ""
    promotion_code = ""
    publisher      = ""
    version        = ""
  }
  regular_priority_profile {
    allocation_strategy = ""
    capacity            = 0
    min_capacity        = 0
  }
  spot_priority_profile {
    allocation_strategy = ""
    capacity            = 0
    eviction_policy     = ""
    maintain            = false
    max_price_per_vm    = 0.0
    min_capacity        = 0
  }
  vm_sizes_profile {
    name = ""
    rank = 0
  }
  tags = {
    key = "value"
  }

  zones = []
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Azure Fleet. Changing this forces a new Azure Fleet to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Azure Fleet should exist. Changing this forces a new Azure Fleet to be created.

* `compute_profile` - (Required) A `compute_profile` block as defined below.

* `location` - (Required) Specifies the Azure Region where the Azure Fleet should exist. Changing this forces a new Azure Fleet to be created.

* `vm_sizes_profile` - (Required) A `vm_sizes_profile` block as defined below.

* `identity` - (Optional) An `identity` block as defined below.

* `plan` - (Optional) A `plan` block as defined below.

* `regular_priority_profile` - (Optional) A `regular_priority_profile` block as defined below.

* `spot_priority_profile` - (Optional) A `spot_priority_profile` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Azure Fleet.

* `zones` - (Optional) Zones in which the Compute Fleet is available.

---

A `compute_profile` block supports the following:

* `base_virtual_machine_profile` - (Required) A `base_virtual_machine_profile` block as defined below.

* `compute_api_version` - (Optional) Specifies the Microsoft.Compute API version to use when creating underlying Virtual Machine scale sets and Virtual Machines.
The default value will be the latest supported computeApiVersion by Compute Fleet.

* `platform_fault_domain_count` - (Optional) Specifies the number of fault domains to use when creating the underlying VMSS.
A fault domain is a logical group of hardware within an Azure datacenter.
VMs in the same fault domain share a common power source and network switch.
If not specified, defaults to 1, which represents "Max Spreading" (using as many fault domains as possible).
This property cannot be updated.

---

A `base_virtual_machine_profile` block supports the following:

* `application_profile` - (Optional) An `application_profile` block as defined below.

* `capacity_reservation` - (Optional) A `capacity_reservation` block as defined below.

* `diagnostics_profile` - (Optional) A `diagnostics_profile` block as defined below.

* `extension_profile` - (Optional) An `extension_profile` block as defined below.

* `hardware_profile` - (Optional) A `hardware_profile` block as defined below.

* `license_type` - (Optional) Specifies that the image or disk that is being used was licensed on-premises.
<br><br> Possible values for Windows Server operating system are: <br><br>
Windows_Client <br><br> Windows_Server <br><br> Possible values for Linux
Server operating system are: <br><br> RHEL_BYOS (for RHEL) <br><br> SLES_BYOS
(for SUSE) <br><br> For more information, see [Azure Hybrid Use Benefit for
Windows
Server](https://docs.microsoft.com/azure/virtual-machines/windows/hybrid-use-benefit-licensing)
<br><br> [Azure Hybrid Use Benefit for Linux
Server](https://docs.microsoft.com/azure/virtual-machines/linux/azure-hybrid-benefit-linux)
<br><br> Minimum api-version: 2015-06-15.

* `network_profile` - (Optional) A `network_profile` block as defined below.

* `os_profile` - (Optional) An `os_profile` block as defined below.

* `scheduled_events_profile` - (Optional) A `scheduled_events_profile` block as defined below.

* `security_posture_reference` - (Optional) A `security_posture_reference` block as defined below.

* `security_profile` - (Optional) A `security_profile` block as defined below.

* `service_artifact_reference` - (Optional) A `service_artifact_reference` block as defined below.

* `storage_profile` - (Optional) A `storage_profile` block as defined below.

* `user_data` - (Optional) UserData for the virtual machines in the scale set, which must be base-64
encoded. Customer should not pass any secrets in here. Minimum api-version:
2021-03-01.

---

An `application_profile` block supports the following:

* `gallery_applications` - (Optional) A `gallery_applications` block as defined below.

---

A `gallery_applications` block supports the following:

* `configuration_reference` - (Optional) Optional, Specifies the uri to an azure blob that will replace the default
configuration for the package if provided.

* `enable_automatic_upgrade` - (Optional) If set to true, when a new Gallery Application version is available in PIR/SIG,
it will be automatically updated for the VM/VMSS.

* `order` - (Optional) Optional, Specifies the order in which the packages have to be installed.

* `package_reference_id` - (Required) Specifies the GalleryApplicationVersion resource id on the form of
/subscriptions/{SubscriptionId}/resourceGroups/{ResourceGroupName}/providers/Microsoft.Compute/galleries/{galleryName}/applications/{application}/versions/{version}.

* `tags` - (Optional) Optional, Specifies a passthrough value for more generic context.

* `treat_failure_as_deployment_failure` - (Optional) Optional, If true, any failure for any operation in the VmApplication will fail
the deployment.

---

A `capacity_reservation` block supports the following:

* `capacity_reservation_group` - (Optional) A `capacity_reservation_group` block as defined below.

---

A `capacity_reservation_group` block supports the following:

* `id` - (Optional) Resource Id.

---

A `diagnostics_profile` block supports the following:

* `boot_diagnostics` - (Optional) A `boot_diagnostics` block as defined below.

---

A `boot_diagnostics` block supports the following:

* `enabled` - (Optional) Whether boot diagnostics should be enabled on the Virtual Machine.

* `storage_uri` - (Optional) Uri of the storage account to use for placing the console output and
screenshot. If storageUri is not specified while enabling boot diagnostics,
managed storage will be used.

---

An `extension_profile` block supports the following:

* `extensions` - (Optional) An `extensions` block as defined below.

* `extensions_time_budget` - (Optional) Specifies the time alloted for all extensions to start. The time duration
should be between 15 minutes and 120 minutes (inclusive) and should be
specified in ISO 8601 format. The default value is 90 minutes (PT1H30M).
Minimum api-version: 2020-06-01.

---

An `extensions` block supports the following:

* `name` - (Optional) Specifies the name of the extension.

* `properties` - (Optional) A `properties` block as defined below.

---

A `properties` block supports the following:

* `auto_upgrade_minor_version` - (Optional) Indicates whether the extension should use a newer minor version if one is
available at deployment time. Once deployed, however, the extension will not
upgrade minor versions unless redeployed, even with this property set to true.

* `enable_automatic_upgrade` - (Optional) Indicates whether the extension should be automatically upgraded by the
platform if there is a newer version of the extension available.

* `force_update_tag` - (Optional) If a value is provided and is different from the previous value, the extension
handler will be forced to update even if the extension configuration has not
changed.

* `protected_settings` - (Optional) Specifies the extension can contain either protectedSettings or
protectedSettingsFromKeyVault or no protected settings at all.

* `protected_settings_from_key_vault` - (Optional) A `protected_settings_from_key_vault` block as defined below.

* `provision_after_extensions` - (Optional) Collection of extension names after which this extension needs to be
provisioned.

* `publisher` - (Optional) Specifies the name of the extension handler publisher.

* `settings` - (Optional) Json formatted public settings for the extension.

* `suppress_failures` - (Optional) Indicates whether failures stemming from the extension will be suppressed
(Operational failures such as not connecting to the VM will not be suppressed
regardless of this value). The default is false.

* `type` - (Optional) Specifies the type of the extension; an example is "CustomScriptExtension".

* `type_handler_version` - (Optional) Specifies the version of the script handler.

---

A `protected_settings_from_key_vault` block supports the following:

* `secret_url` - (Required) Specifies the URL referencing a secret in a Key Vault.

* `source_vault` - (Required) A `source_vault` block as defined below.

---

A `source_vault` block supports the following:

* `id` - (Optional) Resource Id.

---

A `hardware_profile` block supports the following:

* `vm_size_properties` - (Optional) A `vm_size_properties` block as defined below.

---

A `vm_size_properties` block supports the following:

* `vcp_us_available` - (Optional) Specifies the number of vCPUs available for the VM. When this property is not
specified in the request body the default behavior is to set it to the value of
vCPUs available for that VM size exposed in api response of [List all available
virtual machine sizes in a
region](https://docs.microsoft.com/en-us/rest/api/compute/resource-skus/list).

* `vcp_us_per_core` - (Optional) Specifies the vCPU to physical core ratio. When this property is not specified
in the request body the default behavior is set to the value of vCPUsPerCore
for the VM Size exposed in api response of [List all available virtual machine
sizes in a
region](https://docs.microsoft.com/en-us/rest/api/compute/resource-skus/list).
**Setting this property to 1 also means that hyper-threading is disabled.**.

---

A `network_profile` block supports the following:

* `health_probe` - (Optional) A `health_probe` block as defined below.

* `network_api_version` - (Optional) Specifies the Microsoft.Network API version used when creating networking
resources in the Network Interface Configurations for Virtual Machine Scale Set
with orchestration mode 'Flexible'.

* `network_interface_configurations` - (Optional) A `network_interface_configurations` block as defined below.

---

A `health_probe` block supports the following:

* `id` - (Optional) Specifies the ARM resource id in the form of
/subscriptions/{SubscriptionId}/resourceGroups/{ResourceGroupName}/...

---

A `network_interface_configurations` block supports the following:

* `name` - (Required) Specifies the network configuration name.

* `properties` - (Optional) A `properties` block as defined below.

---

A `properties` block supports the following:

* `auxiliary_mode` - (Optional) Specifies whether the Auxiliary mode is enabled for the Network Interface
resource.

* `auxiliary_sku` - (Optional) Specifies whether the Auxiliary sku is enabled for the Network Interface
resource.

* `delete_option` - (Optional) Specify what happens to the network interface when the VM is deleted.

* `disable_tcp_state_tracking` - (Optional) Specifies whether the network interface is disabled for tcp state tracking.

* `dns_settings` - (Optional) A `dns_settings` block as defined below.

* `enable_accelerated_networking` - (Optional) Specifies whether the network interface is accelerated networking-enabled.

* `enable_fpga` - (Optional) Specifies whether the network interface is FPGA networking-enabled.

* `enable_ip_forwarding` - (Optional) Whether IP forwarding enabled on this NIC.

* `ip_configurations` - (Required) An `ip_configurations` block as defined below.

* `network_security_group` - (Optional) A `network_security_group` block as defined below.

* `primary` - (Optional) Specifies the primary network interface in case the virtual machine has more
than 1 network interface.

---

A `dns_settings` block supports the following:

* `dns_servers` - (Optional) List of DNS servers IP addresses.

---

An `ip_configurations` block supports the following:

* `name` - (Required) Specifies the IP configuration name.

* `properties` - (Optional) A `properties` block as defined below.

---

A `properties` block supports the following:

* `application_gateway_backend_address_pools` - (Optional) An `application_gateway_backend_address_pools` block as defined below.

* `application_security_groups` - (Optional) An `application_security_groups` block as defined below.

* `load_balancer_backend_address_pools` - (Optional) A `load_balancer_backend_address_pools` block as defined below.

* `load_balancer_inbound_nat_pools` - (Optional) A `load_balancer_inbound_nat_pools` block as defined below.

* `primary` - (Optional) Specifies the primary network interface in case the virtual machine has more
than 1 network interface.

* `private_ip_address_version` - (Optional) Available from Api-Version 2017-03-30 onwards, it represents whether the
specific ipconfiguration is IPv4 or IPv6. Default is taken as IPv4.  Possible
values are: 'IPv4' and 'IPv6'.

* `public_ip_address_configuration` - (Optional) A `public_ip_address_configuration` block as defined below.

* `subnet` - (Optional) A `subnet` block as defined below.

---

An `application_gateway_backend_address_pools` block supports the following:

* `id` - (Optional) Resource Id.

---

An `application_security_groups` block supports the following:

* `id` - (Optional) Resource Id.

---

A `load_balancer_backend_address_pools` block supports the following:

* `id` - (Optional) Resource Id.

---

A `load_balancer_inbound_nat_pools` block supports the following:

* `id` - (Optional) Resource Id.

---

A `public_ip_address_configuration` block supports the following:

* `name` - (Required) Specifies the publicIP address configuration name.

* `properties` - (Optional) A `properties` block as defined below.

* `sku` - (Optional) A `sku` block as defined below.

---

A `properties` block supports the following:

* `delete_option` - (Optional) Specify what happens to the public IP when the VM is deleted.

* `dns_settings` - (Optional) A `dns_settings` block as defined below.

* `ip_tags` - (Optional) An `ip_tags` block as defined below.

* `idle_timeout_in_minutes` - (Optional) Specifies the idle timeout of the public IP address.

* `public_ip_address_version` - (Optional) Available from Api-Version 2019-07-01 onwards, it represents whether the
specific ipconfiguration is IPv4 or IPv6. Default is taken as IPv4. Possible
values are: 'IPv4' and 'IPv6'.

* `public_ip_prefix` - (Optional) A `public_ip_prefix` block as defined below.

---

A `dns_settings` block supports the following:

* `domain_name_label` - (Required) Specifies the Domain name label.The concatenation of the domain name label and vm index
will be the domain name labels of the PublicIPAddress resources that will be
created.

* `domain_name_label_scope` - (Optional) Specifies the Domain name label scope.The concatenation of the hashed domain name label
that generated according to the policy from domain name label scope and vm
index will be the domain name labels of the PublicIPAddress resources that will
be created.

---

An `ip_tags` block supports the following:

* `ip_tag_type` - (Optional) IP tag type. Example: FirstPartyUsage.

* `tag` - (Optional) IP tag associated with the public IP. Example: SQL, Storage etc.

---

A `public_ip_prefix` block supports the following:

* `id` - (Optional) Resource Id.

---

A `sku` block supports the following:

* `name` - (Optional) Specify public IP sku name.

* `tier` - (Optional) Specify public IP sku tier.

---

A `subnet` block supports the following:

* `id` - (Optional) Specifies the ARM resource id in the form of
/subscriptions/{SubscriptionId}/resourceGroups/{ResourceGroupName}/...

---

A `network_security_group` block supports the following:

* `id` - (Optional) Resource Id.

---

An `os_profile` block supports the following:

* `admin_password` - (Optional) Specifies the password of the administrator account. <br><br> **Minimum-length
(Windows):** 8 characters <br><br> **Minimum-length (Linux):** 6 characters
<br><br> **Max-length (Windows):** 123 characters <br><br> **Max-length
(Linux):** 72 characters <br><br> **Complexity requirements:** 3 out of 4
conditions below need to be fulfilled <br> Has lower characters <br>Has upper
characters <br> Has a digit <br> Has a special character (Regex match [\W_])
<br><br> **Disallowed values:** "abc@123", "P@$$w0rd", "P@ssw0rd",
"P@ssword123", "Pa$$word", "pass@word1", "Password!", "Password1",
"Password22", "iloveyou!" <br><br> For resetting the password, see [How to
reset the Remote Desktop service or its login password in a Windows
VM](https://docs.microsoft.com/troubleshoot/azure/virtual-machines/reset-rdp)
<br><br> For resetting root password, see [Manage users, SSH, and check or
repair disks on Azure Linux VMs using the VMAccess
Extension](https://docs.microsoft.com/troubleshoot/azure/virtual-machines/troubleshoot-ssh-connection).

* `admin_username` - (Optional) Specifies the name of the administrator account. <br><br> **Windows-only
restriction:** Cannot end in "." <br><br> **Disallowed values:**
"administrator", "admin", "user", "user1", "test", "user2", "test1", "user3",
"admin1", "1", "123", "a", "actuser", "adm", "admin2", "aspnet", "backup",
"console", "david", "guest", "john", "owner", "root", "server", "sql",
"support", "support_388945a0", "sys", "test2", "test3", "user4", "user5".
<br><br> **Minimum-length (Linux):** 1  character <br><br> **Max-length
(Linux):** 64 characters <br><br> **Max-length (Windows):** 20 characters.

* `allow_extension_operations` - (Optional) Specifies whether extension operations should be allowed on the virtual machine
scale set. This may only be set to False when no extensions are present on the
virtual machine scale set.

* `computer_name_prefix` - (Optional) Specifies the computer name prefix for all of the virtual machines in the scale
set. Computer name prefixes must be 1 to 15 characters long.

* `custom_data` - (Optional) Specifies a base-64 encoded string of custom data. The base-64 encoded string
is decoded to a binary array that is saved as a file on the Virtual Machine.
The maximum length of the binary array is 65535 bytes. For using cloud-init for
your VM, see [Using cloud-init to customize a Linux VM during
creation](https://docs.microsoft.com/azure/virtual-machines/linux/using-cloud-init).

* `linux_configuration` - (Optional) A `linux_configuration` block as defined below.

* `require_guest_provision_signal` - (Optional) Optional property which must either be set to True or omitted.

* `secrets` - (Optional) A `secrets` block as defined below.

* `windows_configuration` - (Optional) A `windows_configuration` block as defined below.

---

A `linux_configuration` block supports the following:

* `disable_password_authentication` - (Optional) Specifies whether password authentication should be disabled.

* `enable_vm_agent_platform_updates` - (Optional) Indicates whether VMAgent Platform Updates is enabled for the Linux virtual
machine. Default value is false.

* `patch_settings` - (Optional) A `patch_settings` block as defined below.

* `provision_vm_agent` - (Optional) Indicates whether virtual machine agent should be provisioned on the virtual
machine. When this property is not specified in the request body, default
behavior is to set it to true. This will ensure that VM Agent is installed on
the VM so that extensions can be added to the VM later.

* `ssh` - (Optional) A `ssh` block as defined below.

---

A `patch_settings` block supports the following:

* `assessment_mode` - (Optional) Specifies the mode of VM Guest Patch Assessment for the IaaS virtual
machine.<br /><br /> Possible values are:<br /><br /> **ImageDefault** - You
control the timing of patch assessments on a virtual machine. <br /><br />
**AutomaticByPlatform** - The platform will trigger periodic patch assessments.
The property provisionVMAgent must be true.

* `automatic_by_platform_settings` - (Optional) An `automatic_by_platform_settings` block as defined below.

* `patch_mode` - (Optional) Specifies the mode of VM Guest Patching to IaaS virtual machine or virtual
machines associated to virtual machine scale set with OrchestrationMode as
Flexible.<br /><br /> Possible values are:<br /><br /> **ImageDefault** - The
virtual machine's default patching configuration is used. <br /><br />
**AutomaticByPlatform** - The virtual machine will be automatically updated by
the platform. The property provisionVMAgent must be true.

---

An `automatic_by_platform_settings` block supports the following:

* `bypass_platform_safety_checks_on_user_schedule` - (Optional) Enables customer to schedule patching without accidental upgrades.

* `reboot_setting` - (Optional) Specifies the reboot setting for all AutomaticByPlatform patch installation
operations.

---

A `ssh` block supports the following:

* `public_keys` - (Optional) A `public_keys` block as defined below.

---

A `public_keys` block supports the following:

* `key_data` - (Optional) SSH public key certificate used to authenticate with the VM through ssh. The
key needs to be at least 2048-bit and in ssh-rsa format. For creating ssh keys,
see [Create SSH keys on Linux and Mac for Linux VMs in
Azure]https://docs.microsoft.com/azure/virtual-machines/linux/create-ssh-keys-detailed).

* `path` - (Optional) Specifies the full path on the created VM where ssh public key is stored. If
the file already exists, the specified key is appended to the file. Example:
/home/user/.ssh/authorized_keys.

---

A `secrets` block supports the following:

* `source_vault` - (Optional) A `source_vault` block as defined below.

* `vault_certificates` - (Optional) A `vault_certificates` block as defined below.

---

A `source_vault` block supports the following:

* `id` - (Optional) Resource Id.

---

A `vault_certificates` block supports the following:

* `certificate_store` - (Optional) For Windows VMs, specifies the certificate store on the Virtual Machine to
which the certificate should be added. The specified certificate store is
implicitly in the LocalMachine account. For Linux VMs, the certificate file is
placed under the /var/lib/waagent directory, with the file name
&lt;UppercaseThumbprint&gt;.crt for the X509 certificate file and
&lt;UppercaseThumbprint&gt;.prv for private key. Both of these files are .pem
formatted.

* `certificate_url` - (Optional) This is the URL of a certificate that has been uploaded to Key Vault as a
secret. For adding a secret to the Key Vault, see [Add a key or secret to the
key
vault](https://docs.microsoft.com/azure/key-vault/key-vault-get-started/#add).
In this case, your certificate needs to be It is the Base64 encoding of the
following JSON Object which is encoded in UTF-8: <br><br> {<br>
"data":"<Base64-encoded-certificate>",<br>  "dataType":"pfx",<br>
"password":"<pfx-file-password>"<br>} <br> To install certificates on a virtual
machine it is recommended to use the [Azure Key Vault virtual machine extension
for
Linux](https://docs.microsoft.com/azure/virtual-machines/extensions/key-vault-linux)
or the [Azure Key Vault virtual machine extension for
Windows](https://docs.microsoft.com/azure/virtual-machines/extensions/key-vault-windows).

---

A `windows_configuration` block supports the following:

* `additional_unattend_content` - (Optional) An `additional_unattend_content` block as defined below.

* `enable_automatic_updates` - (Optional) Indicates whether Automatic Updates is enabled for the Windows virtual machine.
Default value is true. For virtual machine scale sets, this property can be
updated and updates will take effect on OS reprovisioning.

* `enable_vm_agent_platform_updates` - (Optional) Indicates whether VMAgent Platform Updates is enabled for the Windows virtual
machine. Default value is false.

* `patch_settings` - (Optional) A `patch_settings` block as defined below.

* `provision_vm_agent` - (Optional) Indicates whether virtual machine agent should be provisioned on the virtual
machine. When this property is not specified in the request body, it is set to
true by default. This will ensure that VM Agent is installed on the VM so that
extensions can be added to the VM later.

* `time_zone` - (Optional) Specifies the time zone of the virtual machine. e.g. "Pacific Standard Time".
Possible values can be
[TimeZoneInfo.Id](https://docs.microsoft.com/dotnet/api/system.timezoneinfo.id?#System_TimeZoneInfo_Id)
value from time zones returned by
[TimeZoneInfo.GetSystemTimeZones](https://docs.microsoft.com/dotnet/api/system.timezoneinfo.getsystemtimezones).

* `win_rm` - (Optional) A `win_rm` block as defined below.

---

An `additional_unattend_content` block supports the following:

* `component_name` - (Optional) Specifies the component name. Currently, the only allowable value is
Microsoft-Windows-Shell-Setup.

* `content` - (Optional) Specifies the XML formatted content that is added to the unattend.xml file for
the specified path and component. The XML must be less than 4KB and must
include the root element for the setting or feature that is being inserted.

* `pass_name` - (Optional) Specifies the pass name. Currently, the only allowable value is OobeSystem.

* `setting_name` - (Optional) Specifies the name of the setting to which the content applies. Possible values
are: FirstLogonCommands and AutoLogon.

---

A `patch_settings` block supports the following:

* `assessment_mode` - (Optional) Specifies the mode of VM Guest patch assessment for the IaaS virtual
machine.<br /><br /> Possible values are:<br /><br /> **ImageDefault** - You
control the timing of patch assessments on a virtual machine.<br /><br />
**AutomaticByPlatform** - The platform will trigger periodic patch assessments.
The property provisionVMAgent must be true.

* `automatic_by_platform_settings` - (Optional) An `automatic_by_platform_settings` block as defined below.

* `enable_hotpatching` - (Optional) Enables customers to patch their Azure VMs without requiring a reboot. For
enableHotpatching, the 'provisionVMAgent' must be set to true and 'patchMode'
must be set to 'AutomaticByPlatform'.

* `patch_mode` - (Optional) Specifies the mode of VM Guest Patching to IaaS virtual machine or virtual
machines associated to virtual machine scale set with OrchestrationMode as
Flexible.<br /><br /> Possible values are:<br /><br /> **Manual** - You
control the application of patches to a virtual machine. You do this by
applying patches manually inside the VM. In this mode, automatic updates are
disabled; the property WindowsConfiguration.enableAutomaticUpdates must be
false<br /><br /> **AutomaticByOS** - The virtual machine will automatically be
updated by the OS. The property WindowsConfiguration.enableAutomaticUpdates
must be true. <br /><br /> **AutomaticByPlatform** - the virtual machine will
automatically updated by the platform. The properties provisionVMAgent and
WindowsConfiguration.enableAutomaticUpdates must be true.

---

An `automatic_by_platform_settings` block supports the following:

* `bypass_platform_safety_checks_on_user_schedule` - (Optional) Enables customer to schedule patching without accidental upgrades.

* `reboot_setting` - (Optional) Specifies the reboot setting for all AutomaticByPlatform patch installation
operations.

---

A `win_rm` block supports the following:

* `listeners` - (Optional) A `listeners` block as defined below.

---

A `listeners` block supports the following:

* `certificate_url` - (Optional) This is the URL of a certificate that has been uploaded to Key Vault as a
secret. For adding a secret to the Key Vault, see [Add a key or secret to the
key
vault](https://docs.microsoft.com/azure/key-vault/key-vault-get-started/#add).
In this case, your certificate needs to be the Base64 encoding of the following
JSON Object which is encoded in UTF-8: <br><br> {<br>
"data":"<Base64-encoded-certificate>",<br>  "dataType":"pfx",<br>
"password":"<pfx-file-password>"<br>} <br> To install certificates on a virtual
machine it is recommended to use the [Azure Key Vault virtual machine extension
for
Linux](https://docs.microsoft.com/azure/virtual-machines/extensions/key-vault-linux)
or the [Azure Key Vault virtual machine extension for
Windows](https://docs.microsoft.com/azure/virtual-machines/extensions/key-vault-windows).

* `protocol` - (Optional) Specifies the protocol of WinRM listener. Possible values are: **http,**
**https.**.

---

A `scheduled_events_profile` block supports the following:

* `os_image_notification_profile` - (Optional) An `os_image_notification_profile` block as defined below.

* `terminate_notification_profile` - (Optional) A `terminate_notification_profile` block as defined below.

---

An `os_image_notification_profile` block supports the following:

* `enable` - (Optional) Specifies whether the OS Image Scheduled event is enabled or disabled.

* `not_before_timeout` - (Optional) Length of time a Virtual Machine being reimaged or having its OS upgraded will
have to potentially approve the OS Image Scheduled Event before the event is
auto approved (timed out). The configuration is specified in ISO 8601 format,
and the value must not exceed 15 minutes (PT15M).

---

A `terminate_notification_profile` block supports the following:

* `enable` - (Optional) Specifies whether the Terminate Scheduled event is enabled or disabled.

* `not_before_timeout` - (Optional) Configurable length of time a Virtual Machine being deleted will have to
potentially approve the Terminate Scheduled Event before the event is auto
approved (timed out). The configuration must be specified in ISO 8601 format,
the default value is 5 minutes (PT5M).

---

A `security_posture_reference` block supports the following:

* `exclude_extensions` - (Optional) List of virtual machine extension names to exclude when applying the security
posture.

* `id` - (Optional) Specifies the security posture reference id in the form of
/CommunityGalleries/{communityGalleryName}/securityPostures/{securityPostureName}/versions/{major.minor.patch}|{major.*}|latest.

* `is_overridable` - (Optional) Whether the security posture can be overridden by the user.

---

A `security_profile` block supports the following:

* `encryption_at_host` - (Optional) This property can be used by user in the request to enable or disable the Host
Encryption for the virtual machine or virtual machine scale set. This will
enable the encryption for all the disks including Resource/Temp disk at host
itself. The default behavior is: The Encryption at host will be disabled unless
this property is set to true for the resource.

* `encryption_identity` - (Optional) An `encryption_identity` block as defined below.

* `proxy_agent_settings` - (Optional) A `proxy_agent_settings` block as defined below.

* `security_type` - (Optional) Specifies the SecurityType of the virtual machine. It has to be set to any
specified value to enable UefiSettings. The default behavior is: UefiSettings
will not be enabled unless this property is set.

* `uefi_settings` - (Optional) An `uefi_settings` block as defined below.

---

An `encryption_identity` block supports the following:

* `user_assigned_identity_resource_id` - (Optional) Specifies ARM Resource ID of one of the user identities associated with the VM.

---

A `proxy_agent_settings` block supports the following:

* `enabled` - (Optional) Specifies whether ProxyAgent feature should be enabled on the virtual machine
or virtual machine scale set.

* `key_incarnation_id` - (Optional) Increase the value of this property allows user to reset the key used for
securing communication channel between guest and host.

* `mode` - (Optional) Specifies the mode that ProxyAgent will execute on if the feature is enabled.
ProxyAgent will start to audit or monitor but not enforce access control over
requests to host endpoints in Audit mode, while in Enforce mode it will enforce
access control. The default value is Enforce mode.

---

An `uefi_settings` block supports the following:

* `secure_boot_enabled` - (Optional) Specifies whether secure boot should be enabled on the virtual machine. Minimum
api-version: 2020-12-01.

* `v_tpm_enabled` - (Optional) Specifies whether vTPM should be enabled on the virtual machine. Minimum
api-version: 2020-12-01.

---

A `service_artifact_reference` block supports the following:

* `id` - (Optional) Specifies the service artifact reference id in the form of
/subscriptions/{subscriptionId}/resourceGroups/{resourceGroup}/providers/Microsoft.Compute/galleries/{galleryName}/serviceArtifacts/{serviceArtifactName}/vmArtifactsProfiles/{vmArtifactsProfilesName}.

---

A `storage_profile` block supports the following:

* `data_disks` - (Optional) A `data_disks` block as defined below.

* `disk_controller_type` - (Optional) Specifies the disk controller type configured for the virtual machines in the scale set. Minimum api-version: 2022-08-01.

* `image_reference` - (Optional) An `image_reference` block as defined below.

* `os_disk` - (Optional) An `os_disk` block as defined below.

---

A `data_disks` block supports the following:

* `caching` - (Optional) Specifies the caching requirements. Possible values are: **None,**
**ReadOnly,** **ReadWrite.** The default values are: **None for Standard
storage. ReadOnly for Premium storage.**.

* `create_option` - (Required) Specifies the create option.

* `delete_option` - (Optional) Specifies whether data disk should be deleted or detached upon VMSS Flex
deletion (This feature is available for VMSS with Flexible OrchestrationMode
only).<br><br> Possible values: <br><br> **Delete** If this value is used, the
data disk is deleted when the VMSS Flex VM is deleted.<br><br> **Detach** If
this value is used, the data disk is retained after VMSS Flex VM is
deleted.<br><br> The default value is set to **Delete**.

* `disk_iops_read_write` - (Optional) Specifies the Read-Write IOPS for the managed disk. Should be used only when
StorageAccountType is UltraSSD_LRS. If not specified, a default value would be
assigned based on diskSizeGB.

* `disk_m_bps_read_write` - (Optional) Specifies the bandwidth in MB per second for the managed disk. Should be used
only when StorageAccountType is UltraSSD_LRS. If not specified, a default value
would be assigned based on diskSizeGB.

* `disk_size_gb` - (Optional) Specifies the size of an empty data disk in gigabytes. This element can be used
to overwrite the size of the disk in a virtual machine image. The property
diskSizeGB is the number of bytes x 1024^3 for the disk and the value cannot be
larger than 1023.

* `lun` - (Required) Specifies the logical unit number of the data disk. This value is used to
identify data disks within the VM and therefore must be unique for each data
disk attached to a VM.

* `managed_disk` - (Optional) A `managed_disk` block as defined below.

* `name` - (Optional) Specifies the disk name.

* `write_accelerator_enabled` - (Optional) Specifies whether writeAccelerator should be enabled or disabled on the disk.

---

A `managed_disk` block supports the following:

* `disk_encryption_set` - (Optional) A `disk_encryption_set` block as defined below.

* `security_profile` - (Optional) A `security_profile` block as defined below.

* `storage_account_type` - (Optional) Specifies the storage account type for the managed disk. NOTE: UltraSSD_LRS can
only be used with data disks, it cannot be used with OS Disk.

---

A `disk_encryption_set` block supports the following:

* `id` - (Optional) Resource Id.

---

A `security_profile` block supports the following:

* `disk_encryption_set` - (Optional) A `disk_encryption_set` block as defined below.

* `security_encryption_type` - (Optional) Specifies the EncryptionType of the managed disk. It is set to
DiskWithVMGuestState for encryption of the managed disk along with VMGuestState
blob, VMGuestStateOnly for encryption of just the VMGuestState blob, and
NonPersistedTPM for not persisting firmware state in the VMGuestState blob..
**Note:** It can be set for only Confidential VMs.

---

A `disk_encryption_set` block supports the following:

* `id` - (Optional) Resource Id.

---

An `image_reference` block supports the following:

* `community_gallery_image_id` - (Optional) Specified the community gallery image unique id for vm deployment. This can be
fetched from community gallery image GET call.

* `id` - (Optional) Resource Id.

* `offer` - (Optional) Specifies the offer of the platform image or marketplace image used to create
the virtual machine.

* `publisher` - (Optional) Specifies the image publisher.

* `shared_gallery_image_id` - (Optional) Specified the shared gallery image unique id for vm deployment. This can be
fetched from shared gallery image GET call.

* `sku` - (Optional) Specifies the image SKU.

* `version` - (Optional) Specifies the version of the platform image or marketplace image used to create
the virtual machine. The allowed formats are Major.Minor.Build or 'latest'.
Major, Minor, and Build are decimal numbers. Specify 'latest' to use the latest
version of an image available at deploy time. Even if you use 'latest', the VM
image will not automatically update after deploy time even if a new version
becomes available. Please do not use field 'version' for gallery image
deployment, gallery image should always use 'id' field for deployment, to use 'latest'
version of gallery image, just set
'/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/galleries/{galleryName}/images/{imageName}'
in the 'id' field without version input.

---

An `os_disk` block supports the following:

* `caching` - (Optional) Specifies the caching requirements. Possible values are: **None,**
**ReadOnly,** **ReadWrite.** The default values are: **None for Standard
storage. ReadOnly for Premium storage.**.

* `create_option` - (Required) Specifies how the virtual machines in the scale set should be created. The only
allowed value is: **FromImage.** This value is used when you are using an image
to create the virtual machine. If you are using a platform image, you also use
the imageReference element described above. If you are using a marketplace
image, you  also use the plan element previously described.

* `delete_option` - (Optional) Specifies whether OS Disk should be deleted or detached upon VMSS Flex deletion
(This feature is available for VMSS with Flexible OrchestrationMode only).
<br><br> Possible values: <br><br> **Delete** If this value is used, the OS
disk is deleted when VMSS Flex VM is deleted.<br><br> **Detach** If this value
is used, the OS disk is retained after VMSS Flex VM is deleted. <br><br> The
default value is set to **Delete**. For an Ephemeral OS Disk, the default value
is set to **Delete**. User cannot change the delete option for Ephemeral OS
Disk.

* `diff_disk_settings` - (Optional) A `diff_disk_settings` block as defined below.

* `disk_size_gb` - (Optional) Specifies the size of an empty data disk in gigabytes. This element can be used
to overwrite the size of the disk in a virtual machine image. The property 'diskSizeGB'
is the number of bytes x 1024^3 for the disk and the value cannot
be larger than 1023.

* `image` - (Optional) An `image` block as defined below.

* `managed_disk` - (Optional) A `managed_disk` block as defined below.

* `name` - (Optional) Specifies the disk name.

* `os_type` - (Optional) This property allows you to specify the type of the OS that is included in the
disk if creating a VM from user-image or a specialized VHD. Possible values
are: **Windows,** **Linux.**.

* `vhd_containers` - (Optional) Specifies the container urls that are used to store operating system disks for
the scale set.

* `write_accelerator_enabled` - (Optional) Specifies whether writeAccelerator should be enabled or disabled on the disk.

---

A `diff_disk_settings` block supports the following:

* `option` - (Optional) Specifies the ephemeral disk settings for operating system disk.

* `placement` - (Optional) Specifies the ephemeral disk placement for operating system disk. Possible
values are: **CacheDisk,** **ResourceDisk.** The defaulting behavior is:
**CacheDisk** if one is configured for the VM size otherwise **ResourceDisk**
is used. Refer to the VM size documentation for Windows VM at
https://docs.microsoft.com/azure/virtual-machines/windows/sizes and Linux VM at
https://docs.microsoft.com/azure/virtual-machines/linux/sizes to check which VM
sizes exposes a cache disk.

---

An `image` block supports the following:

* `uri` - (Optional) Specifies the virtual hard disk's uri.

---

A `managed_disk` block supports the following:

* `disk_encryption_set` - (Optional) A `disk_encryption_set` block as defined below.

* `security_profile` - (Optional) A `security_profile` block as defined below.

* `storage_account_type` - (Optional) Specifies the storage account type for the managed disk. NOTE: UltraSSD_LRS can
only be used with data disks, it cannot be used with OS Disk.

---

A `disk_encryption_set` block supports the following:

* `id` - (Optional) Resource Id.

---

A `security_profile` block supports the following:

* `disk_encryption_set` - (Optional) A `disk_encryption_set` block as defined below.

* `security_encryption_type` - (Optional) Specifies the EncryptionType of the managed disk. It is set to
DiskWithVMGuestState for encryption of the managed disk along with VMGuestState
blob, VMGuestStateOnly for encryption of just the VMGuestState blob, and
NonPersistedTPM for not persisting firmware state in the VMGuestState blob..
**Note:** It can be set for only Confidential VMs.

---

A `disk_encryption_set` block supports the following:

* `id` - (Optional) Resource Id.

---

A `vm_sizes_profile` block supports the following:

* `name` - (Required) Specifies the Sku name (e.g. 'Standard_DS1_v2').

* `rank` - (Optional) Specifies the rank of the VM size. This is used with 'RegularPriorityAllocationStrategy.Prioritized'
The lower the number, the higher the priority. Starting with 0.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) A list of IDs for User Assigned Managed Identity resources to be assigned.

---

A `plan` block supports the following:

* `name` - (Required) A user defined name of the 3rd Party Artifact that is being procured.

* `product` - (Required) Specifies the 3rd Party artifact that is being procured. E.g. NewRelic. Product maps to the OfferID specified for the artifact at the time of Data Market onboarding. .

* `promotion_code` - (Optional) A publisher provided promotion code as provisioned in Data Market for the said product/artifact.

* `publisher` - (Required) Specifies the publisher of the 3rd Party Artifact that is being bought. E.g. NewRelic.

* `version` - (Optional) Specifies the version of the desired product/artifact.

---

A `regular_priority_profile` block supports the following:

* `allocation_strategy` - (Optional) Allocation strategy to follow when determining the VM sizes distribution for Regular VMs.

* `capacity` - (Optional) Total capacity to achieve. It is currently in terms of number of VMs.

* `min_capacity` - (Optional) Minimum capacity to achieve which cannot be updated. If we will not be able to "guarantee" minimum capacity, we will reject the request in the sync path itself.

---

A `spot_priority_profile` block supports the following:

* `allocation_strategy` - (Optional) Allocation strategy to follow when determining the VM sizes distribution for Spot VMs.

* `capacity` - (Optional) Total capacity to achieve. It is currently in terms of number of VMs.

* `eviction_policy` - (Optional) Eviction Policy to follow when evicting Spot VMs.

* `maintain` - (Optional) Flag to enable/disable continuous goal seeking for the desired capacity and restoration of evicted Spot VMs.
If maintain is enabled, AzureFleetRP will use all VM sizes in vmSizesProfile to create new VMs (if VMs are evicted deleted)
or update existing VMs with new VM sizes (if VMs are evicted deallocated or failed to allocate due to capacity constraint) in order to achieve the desired capacity.
Maintain is enabled by default.

* `max_price_per_vm` - (Optional) Price per hour of each Spot VM will never exceed this.

* `min_capacity` - (Optional) Minimum capacity to achieve which cannot be updated. If we will not be able to "guarantee" minimum capacity, we will reject the request in the sync path itself.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure Fleet.

* `compute_profile` - A `compute_profile` block as defined below.

* `identity` - An `identity` block as defined below.

* `time_created` - Specifies the time at which the Compute Fleet is created.

* `unique_id` - Specifies the ID which uniquely identifies a Compute Fleet.

---

A `compute_profile` block exports the following:

* `base_virtual_machine_profile` - A `base_virtual_machine_profile` block as defined below.

---

A `base_virtual_machine_profile` block exports the following:

* `extension_profile` - An `extension_profile` block as defined below.

* `storage_profile` - A `storage_profile` block as defined below.

* `time_created` - Specifies the time in which this VM profile for the Virtual Machine Scale Set
was created. Minimum API version for this property is 2023-09-01. This value
will be added to VMSS Flex VM tags when creating/updating the VMSS VM Profile
with minimum api-version 2023-09-01. Examples: "2024-07-01T00:00:01.1234567+00:00".

---

An `extension_profile` block exports the following:

* `extensions` - An `extensions` block as defined below.

---

An `extensions` block exports the following:

* `id` - Resource Id.

* `properties` - A `properties` block as defined below.

* `type` - Resource type.

---

A `properties` block exports the following:

* `provisioning_state` - The provisioning state, which only appears in the response.

---

A `storage_profile` block exports the following:

* `image_reference` - An `image_reference` block as defined below.

---

An `image_reference` block exports the following:

* `exact_version` - Specifies in decimal numbers, the version of platform image or marketplace
image used to create the virtual machine. This readonly field differs from 'version',
only if the value specified in 'version' field is 'latest'.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Azure Fleet.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Fleet.
* `update` - (Defaults to 30 minutes) Used when updating the Azure Fleet.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure Fleet.

## Import

Azure Fleet can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_azure_fleet.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.AzureFleet/fleets/fleet1
```
