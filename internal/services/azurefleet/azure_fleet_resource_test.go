package azurefleet_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/azurefleet/2024-11-01/fleets"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AzureFleetResource struct{}

func TestAccAzureFleetFleet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureFleetFleet_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccAzureFleetFleet_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureFleetFleet_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r AzureFleetResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := fleets.ParseFleetID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.AzureFleet.FleetsClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r AzureFleetResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_public_ip" "test" {
  name                = "test-ip-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_lb" "test" {
  name                = "acctest-loadbalancer-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "internal-%[1]d"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  name            = "internal"
  loadbalancer_id = azurerm_lb.test.id
}

`, data.RandomInteger, data.Locations.Primary)
}

func (r AzureFleetResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`

%[1]s

resource "azurerm_azure_fleet" "test" {
  name                = "acctest-aff-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  spot_priority_profile {
    allocation_strategy = "PriceCapacityOptimized"
    capacity            = 2
    eviction_policy     = "Delete"
    maintain            = true
    min_capacity        = 1
  }

   regular_priority_profile {
    allocation_strategy = "LowestPrice"
    capacity            = 2
    min_capacity        = 1
  }

  vm_sizes_profile {
    name = "Standard_D2s_v3"
  }

   vm_sizes_profile {
    name = "Standard_D4s_v3"
  }

   vm_sizes_profile {
    name = "Standard_E2s_v3"
  }

  compute_profile {
    virtual_machine_profile {
      storage_profile {
         image_reference {
          offer                      = "0001-com-ubuntu-server-focal"
          publisher                  = "canonical"
          sku                        = "20_04-lts-gen2"
          version                    = "latest"
         }

         os_disk {
          caching                   = "ReadWrite"
          create_option             = "FromImage"
          os_type                   = "Linux"
          managed_disk {
            storage_account_type = "Standard_LRS"
          }
         }
      }

      os_profile {
        computer_name_prefix           = "prefix"
        admin_username                 = "azureuser"
        admin_password                 = "TestPassword$0"
        linux_configuration {
          disable_password_authentication  = false
        }
      }

      network_profile {
        network_interfaces {
          name = "networkProTest"
          properties {
            primary                       = true
            enable_accelerated_networking = false
            enable_ip_forwarding          = true
            ip_configurations {
              name = "ipConfigTest"
              properties {
                primary                    = true
                load_balancer_backend_address_pools {
                  id = azurerm_lb_backend_address_pool.test.id
                }
                subnet {
                  id = azurerm_subnet.test.id
                }
              }
            }
          }
        }
        network_api_version = "2022-07-01"
      }
    }
    compute_api_version         = "2023-09-01"
    platform_fault_domain_count = 1
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r AzureFleetResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_azure_fleet" "import" {
  name                = azurerm_azure_fleet.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  compute_profile {
    compute_api_version         = ""
    platform_fault_domain_count = 0
    additional_capabilities {
      hibernation_enabled = false
      ultra_ssd_enabled   = false
    }
    virtual_machine_profile {
      license_type = ""
      user_data_base64    = ""
      application_profile {
        gallery_application {
          configuration_reference             = ""
          enable_automatic_upgrade            = false
          order                               = 0
          package_reference_id                = ""
          tags                                = ""
          treat_failure_as_deployment_failure_enabled = false
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
            protected_settings_json = jsonencode({
              "key" : "value"
            })
            settings_json = jsonencode({
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
        network_health_probe_id {
          id = ""
        }
        network_interfaces {
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
        custom_data_base64                    = ""
        require_guest_provision_signal = false
        linux_configuration {
          disable_password_authentication  = false
          enable_vm_agent_platform_updates = false
          provision_vm_agent_enabled               = false
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
          provision_vm_agent_enabled               = false
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
      scheduled_event_os_image_enabled = true
        scheduled_event_os_image_timeout = "PT5M"
        scheduled_event_termination_enabled = true
        scheduled_event_termination_timeout = "PT15M"

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
  vm_sizes {
    name = ""
    rank = 0
  }
}
`, config, data.Locations.Primary)
}

func (r AzureFleetResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_azure_fleet" "test" {
  name                = "acctest-aff-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = []
  }
  additional_locations_profile {
    location_profiles {
      location = "%s"
      virtual_machine_profile_override {
        license_type = ""
        user_data_base64    = ""
        application_profile {
          gallery_application {
            configuration_reference             = ""
            enable_automatic_upgrade            = false
            order                               = 0
            package_reference_id                = ""
            tags                                = ""
            treat_failure_as_deployment_failure_enabled = false
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
              protected_settings_json = jsonencode({
                "key" : "value"
              })
              settings_json = jsonencode({
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
          network_health_probe_id {
            id = ""
          }
          network_interfaces {
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
          custom_data_base64                    = ""
          require_guest_provision_signal = false
          linux_configuration {
            disable_password_authentication  = false
            enable_vm_agent_platform_updates = false
            provision_vm_agent_enabled               = false
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
            provision_vm_agent_enabled               = false
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
        scheduled_event_os_image_enabled = true
        scheduled_event_os_image_timeout = "PT5M"
        scheduled_event_termination_enabled = true
        scheduled_event_termination_timeout = "PT15M"

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
  }
  compute_profile {
    compute_api_version         = ""
    platform_fault_domain_count = 0
    additional_capabilities {
      hibernation_enabled = false
      ultra_ssd_enabled   = false
    }
    virtual_machine_profile {
      license_type = ""
      user_data_base64    = ""
      application_profile {
        gallery_application {
          configuration_reference             = ""
          enable_automatic_upgrade            = false
          order                               = 0
          package_reference_id                = ""
          tags                                = ""
          treat_failure_as_deployment_failure_enabled = false
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
            protected_settings_json = jsonencode({
              "key" : "value"
            })
            settings_json = jsonencode({
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
        network_health_probe_id {
          id = ""
        }
        network_interfaces {
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
        custom_data_base64                    = ""
        require_guest_provision_signal = false
        linux_configuration {
          disable_password_authentication  = false
          enable_vm_agent_platform_updates = false
          provision_vm_agent_enabled               = false
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
          provision_vm_agent_enabled               = false
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

      scheduled_event_os_image_enabled = true
        scheduled_event_os_image_timeout = "PT5M"
        scheduled_event_termination_enabled = true
        scheduled_event_termination_timeout = "PT15M"

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
  vm_attributes {
    accelerator_support   = ""
    burstable_support     = ""
    local_storage_support = ""
    rdma_support          = ""
    excluded_vm_sizes     = []
    accelerator_count {
      max = 0
      min = 0
    }
    accelerator_manufacturers {

    }
    accelerator_types {

    }
    architecture_types {

    }
    cpu_manufacturers {

    }
    data_disk_count {
      max = 0
      min = 0
    }
    local_storage_disk_types {

    }
    local_storage_in_gi_b {
      max = 0.0
      min = 0.0
    }
    memory_in_gi_b {
      max = 0.0
      min = 0.0
    }
    memory_in_gi_b_per_vcpu {
      max = 0.0
      min = 0.0
    }
    network_bandwidth_in_mbps {
      max = 0.0
      min = 0.0
    }
    network_interface_count {
      max = 0
      min = 0
    }
    rdma_network_interface_count {
      max = 0
      min = 0
    }
    vcpu_count {
      max = 0
      min = 0
    }
    vm_categories {

    }
  }
  vm_sizes_profile {
    name = "Standard_d1_v2"
    rank = 19225
  }
  tags = {
    key = "value"
  }

  zones = []
}
`, template, data.RandomInteger, data.Locations.Primary, data.Locations.Primary)
}

func (r AzureFleetResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_azure_fleet" "test" {
  name                = "acctest-aff-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = []
  }
  additional_locations_profile {
    location_profiles {
      location = "%s"
      virtual_machine_profile_override {
        license_type = ""
        user_data_base64    = ""
        application_profile {
          gallery_application {
            configuration_reference             = ""
            enable_automatic_upgrade            = false
            order                               = 0
            package_reference_id                = ""
            tags                                = ""
            treat_failure_as_deployment_failure_enabled = false
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
              protected_settings_json = jsonencode({
                "key" : "value"
              })
              settings_json = jsonencode({
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
          network_health_probe_id {
            id = ""
          }
          network_interfaces {
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
          custom_data_base64                    = ""
          require_guest_provision_signal = false
          linux_configuration {
            disable_password_authentication  = false
            enable_vm_agent_platform_updates = false
            provision_vm_agent_enabled               = false
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
            provision_vm_agent_enabled               = false
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

        scheduled_event_os_image_enabled = true
        scheduled_event_os_image_timeout = "PT5M"
        scheduled_event_termination_enabled = true
        scheduled_event_termination_timeout = "PT15M"

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
  }
  compute_profile {
    compute_api_version         = ""
    platform_fault_domain_count = 0
    additional_capabilities {
      hibernation_enabled = false
      ultra_ssd_enabled   = false
    }
    virtual_machine_profile {
      license_type = ""
      user_data_base64    = ""
      application_profile {
        gallery_application {
          configuration_reference             = ""
          enable_automatic_upgrade            = false
          order                               = 0
          package_reference_id                = ""
          tags                                = ""
          treat_failure_as_deployment_failure_enabled = false
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
            protected_settings_json = jsonencode({
              "key" : "value"
            })
            settings_json = jsonencode({
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
        network_health_probe_id {
          id = ""
        }
        network_interfaces {
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
        custom_data_base64                    = ""
        require_guest_provision_signal = false
        linux_configuration {
          disable_password_authentication  = false
          enable_vm_agent_platform_updates = false
          provision_vm_agent_enabled               = false
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
          provision_vm_agent_enabled               = false
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

      scheduled_event_os_image_enabled = true
        scheduled_event_os_image_timeout = "PT5M"
        scheduled_event_termination_enabled = true
        scheduled_event_termination_timeout = "PT15M"

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
  vm_attributes {
    accelerator_support   = ""
    burstable_support     = ""
    local_storage_support = ""
    rdma_support          = ""
    excluded_vm_sizes     = []
    accelerator_count {
      max = 0
      min = 0
    }
    accelerator_manufacturers {

    }
    accelerator_types {

    }
    architecture_types {

    }
    cpu_manufacturers {

    }
    data_disk_count {
      max = 0
      min = 0
    }
    local_storage_disk_types {

    }
    local_storage_in_gi_b {
      max = 0.0
      min = 0.0
    }
    memory_in_gi_b {
      max = 0.0
      min = 0.0
    }
    memory_in_gi_b_per_vcpu {
      max = 0.0
      min = 0.0
    }
    network_bandwidth_in_mbps {
      max = 0.0
      min = 0.0
    }
    network_interface_count {
      max = 0
      min = 0
    }
    rdma_network_interface_count {
      max = 0
      min = 0
    }
    vcpu_count {
      max = 0
      min = 0
    }
    vm_categories {

    }
  }
  vm_sizes {
    name = ""
    rank = 0
  }
  tags = {
    key = "value"
  }

  zones = []
}
`, template, data.RandomInteger, data.Locations.Primary, data.Locations.Primary)
}
