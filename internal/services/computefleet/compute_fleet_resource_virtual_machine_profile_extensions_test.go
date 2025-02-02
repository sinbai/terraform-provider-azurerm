// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package computefleet_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccComputeFleet_virtualMachineProfileExtensions_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.extensionsBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_virtualMachineProfileExtensions_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.extensionsComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"virtual_machine_profile.0.extension.2.protected_settings_json",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.extension.2.protected_settings_json"),
		{
			Config: r.extensionsCompleteUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"virtual_machine_profile.0.extension.2.protected_settings_json",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.extension.2.protected_settings_json"),
		{
			Config: r.extensionsBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.extensionsComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"virtual_machine_profile.0.extension.2.protected_settings_json",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.extension.2.protected_settings_json"),
	})
}

func TestAccComputeFleet_virtualMachineProfileExtensions_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.extensionsComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"virtual_machine_profile.0.extension.2.protected_settings_json",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.extension.2.protected_settings_json"),
	})
}

func (r ComputeFleetTestResource) extensionsBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_compute_fleet" "test" {
  name                        = "acctest-fleet-%[2]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = "%[3]s"
  platform_fault_domain_count = 2

  regular_priority_profile {
    capacity     = 0
    min_capacity = 0
  }

  vm_sizes_profile {
    name = "Standard_D1_v2"
  }

  compute_api_version = "2024-03-01"
  virtual_machine_profile {
    network_api_version = "2020-11-01"
    os_profile {
      linux_configuration {
        computer_name_prefix            = "prefix"
        admin_username                  = local.admin_username
        password_authentication_enabled = false
        admin_ssh_keys                  = [local.first_public_key]
      }
    }
    network_interface {
      name    = "networkProTest"
      primary = true

      ip_configuration {
        name      = "TestIPConfiguration"
        primary   = true
        subnet_id = azurerm_subnet.test.id

        public_ip_address {
          name                    = "TestPublicIPConfiguration"
          domain_name_label       = "test-domain-label"
          idle_timeout_in_minutes = 4
        }
      }
    }

    os_disk {
      storage_account_type = "Standard_LRS"
      caching              = "ReadWrite"
    }

    source_image_reference {
      publisher = "Canonical"
      offer     = "0001-com-ubuntu-server-jammy"
      sku       = "22_04-lts"
      version   = "latest"
    }

    extension {
      name                 = "CustomScript"
      publisher            = "Microsoft.Azure.Extensions"
      type                 = "CustomScript"
      type_handler_version = "2.0"
    }
  }
  additional_location_profile {
    location = "%[4]s"
    virtual_machine_profile_override {
      network_api_version = "2020-11-01"
      os_profile {
        linux_configuration {
          computer_name_prefix            = "prefix"
          admin_username                  = local.admin_username
          password_authentication_enabled = false
          admin_ssh_keys                  = [local.first_public_key]
        }
      }
      network_interface {
        name    = "networkProTest"
        primary = true

        ip_configuration {
          name      = "TestIPConfiguration"
          primary   = true
          subnet_id = azurerm_subnet.linux_test.id

          public_ip_address {
            name                    = "TestPublicIPConfiguration"
            domain_name_label       = "test-domain-label"
            idle_timeout_in_minutes = 4
          }
        }
      }

      os_disk {
        storage_account_type = "Standard_LRS"
        caching              = "ReadWrite"
      }

      source_image_reference {
        publisher = "Canonical"
        offer     = "0001-com-ubuntu-server-jammy"
        sku       = "22_04-lts"
        version   = "latest"
      }

      extension {
        name                 = "CustomScript"
        publisher            = "Microsoft.Azure.Extensions"
        type                 = "CustomScript"
        type_handler_version = "2.0"
      }
    }
  }
}
`, r.baseAndAdditionalLocationLinuxTemplate(data), data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (r ComputeFleetTestResource) extensionsComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_compute_fleet" "test" {
  name                        = "acctest-fleet-%[2]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = "%[3]s"
  platform_fault_domain_count = 2

  regular_priority_profile {
    capacity     = 0
    min_capacity = 0
  }

  vm_sizes_profile {
    name = "Standard_D1_v2"
  }

  compute_api_version = "2024-03-01"
  virtual_machine_profile {
    network_api_version = "2020-11-01"
    os_profile {
      linux_configuration {
        computer_name_prefix            = "prefix"
        admin_username                  = local.admin_username
        password_authentication_enabled = false
        admin_ssh_keys                  = [local.first_public_key]
      }
    }
    network_interface {
      name    = "networkProTest"
      primary = true

      ip_configuration {
        name      = "TestIPConfiguration"
        primary   = true
        subnet_id = azurerm_subnet.test.id

        public_ip_address {
          name                    = "TestPublicIPConfiguration"
          domain_name_label       = "test-domain-label"
          idle_timeout_in_minutes = 4
        }
      }
    }

    os_disk {
      storage_account_type = "Standard_LRS"
      caching              = "ReadWrite"
    }

    source_image_reference {
      publisher = "Canonical"
      offer     = "0001-com-ubuntu-server-jammy"
      sku       = "22_04-lts"
      version   = "latest"
    }

    extension_operations_enabled = true
    extension {
      name                               = "testOmsAgentForLinux"
      publisher                          = "Microsoft.EnterpriseCloud.Monitoring"
      type                               = "OmsAgentForLinux"
      type_handler_version               = "1.12"
      auto_upgrade_minor_version_enabled = true
      automatic_upgrade_enabled          = true
      failure_suppression_enabled        = true
    }

    extension {
      name                               = "CustomScript"
      publisher                          = "Microsoft.Azure.Extensions"
      type                               = "CustomScript"
      type_handler_version               = "2.0"
      auto_upgrade_minor_version_enabled = true
    }

    extension {
      name                                      = "Docker"
      publisher                                 = "Microsoft.Azure.Extensions"
      type                                      = "DockerExtension"
      type_handler_version                      = "1.0"
      auto_upgrade_minor_version_enabled        = true
      extensions_to_provision_after_vm_creation = ["CustomScript"]
      force_extension_execution_on_change       = "test"

      settings_json = jsonencode({
        "commandToExecute" = "echo $HOSTNAME"
      })

      protected_settings_json = jsonencode({
        "managedIdentity" = {}
      })
    }
    extensions_time_budget = "PT30M"
  }

  additional_location_profile {
    location = "%[4]s"
    virtual_machine_profile_override {
      network_api_version = "2020-11-01"
      os_profile {
        linux_configuration {
          computer_name_prefix            = "prefix"
          admin_username                  = local.admin_username
          password_authentication_enabled = false
          admin_ssh_keys                  = [local.first_public_key]
        }
      }
      network_interface {
        name    = "networkProTest"
        primary = true

        ip_configuration {
          name      = "TestIPConfiguration"
          primary   = true
          subnet_id = azurerm_subnet.linux_test.id

          public_ip_address {
            name                    = "TestPublicIPConfiguration"
            domain_name_label       = "test-domain-label"
            idle_timeout_in_minutes = 4
          }
        }
      }

      os_disk {
        storage_account_type = "Standard_LRS"
        caching              = "ReadWrite"
      }

      source_image_reference {
        publisher = "Canonical"
        offer     = "0001-com-ubuntu-server-jammy"
        sku       = "22_04-lts"
        version   = "latest"
      }

      extension_operations_enabled = true
      extension {
        name                               = "testOmsAgentForLinux"
        publisher                          = "Microsoft.EnterpriseCloud.Monitoring"
        type                               = "OmsAgentForLinux"
        type_handler_version               = "1.12"
        auto_upgrade_minor_version_enabled = true
        automatic_upgrade_enabled          = true
        failure_suppression_enabled        = true
      }

      extension {
        name                               = "CustomScript"
        publisher                          = "Microsoft.Azure.Extensions"
        type                               = "CustomScript"
        type_handler_version               = "2.0"
        auto_upgrade_minor_version_enabled = true
      }

      extension {
        name                                      = "Docker"
        publisher                                 = "Microsoft.Azure.Extensions"
        type                                      = "DockerExtension"
        type_handler_version                      = "1.0"
        auto_upgrade_minor_version_enabled        = true
        extensions_to_provision_after_vm_creation = ["CustomScript"]
        force_extension_execution_on_change       = "test"

        settings_json = jsonencode({
          "commandToExecute" = "echo $HOSTNAME"
        })

        protected_settings_json = jsonencode({
          "managedIdentity" = {}
        })
      }
      extensions_time_budget = "PT30M"
    }
  }
}
`, r.baseAndAdditionalLocationLinuxTemplate(data), data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (r ComputeFleetTestResource) extensionsCompleteUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_compute_fleet" "test" {
  name                        = "acctest-fleet-%[2]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = "%[3]s"
  platform_fault_domain_count = 2

  regular_priority_profile {
    capacity     = 0
    min_capacity = 0
  }

  vm_sizes_profile {
    name = "Standard_D1_v2"
  }

  compute_api_version = "2024-03-01"
  virtual_machine_profile {
    network_api_version = "2020-11-01"
    os_profile {
      linux_configuration {
        computer_name_prefix            = "prefix"
        admin_username                  = local.admin_username
        password_authentication_enabled = false
        admin_ssh_keys                  = [local.first_public_key]
      }
    }
    network_interface {
      name    = "networkProTest"
      primary = true

      ip_configuration {
        name      = "TestIPConfiguration"
        primary   = true
        subnet_id = azurerm_subnet.test.id

        public_ip_address {
          name                    = "TestPublicIPConfiguration"
          domain_name_label       = "test-domain-label"
          idle_timeout_in_minutes = 4
        }
      }
    }

    os_disk {
      storage_account_type = "Standard_LRS"
      caching              = "ReadWrite"
    }

    source_image_reference {
      publisher = "Canonical"
      offer     = "0001-com-ubuntu-server-jammy"
      sku       = "22_04-lts"
      version   = "latest"
    }

    extension_operations_enabled = true
    extension {
      name                               = "testOmsAgentForLinux"
      publisher                          = "Microsoft.EnterpriseCloud.Monitoring"
      type                               = "OmsAgentForLinux"
      type_handler_version               = "1.12"
      auto_upgrade_minor_version_enabled = false
      automatic_upgrade_enabled          = false
      failure_suppression_enabled        = false
    }

    extension {
      name                               = "CustomScript"
      publisher                          = "Microsoft.Azure.Extensions"
      type                               = "CustomScript"
      type_handler_version               = "2.0"
      auto_upgrade_minor_version_enabled = true
    }

    extension {
      name                 = "Docker"
      publisher            = "Microsoft.Azure.Extensions"
      type                 = "DockerExtension"
      type_handler_version = "1.0"

      auto_upgrade_minor_version_enabled        = false
      extensions_to_provision_after_vm_creation = ["CustomScript"]
      force_extension_execution_on_change       = "testUpdate"

      settings_json = jsonencode({
        "commandToExecute" = "echo $(date)"
      })

      protected_settings_json = jsonencode({
        "reset_ssh" = "True"
      })
    }
    extensions_time_budget = "PT1H"
  }
  additional_location_profile {
    location = "%[4]s"
    virtual_machine_profile_override {
      network_api_version = "2020-11-01"
      os_profile {
        linux_configuration {
          computer_name_prefix            = "prefix"
          admin_username                  = local.admin_username
          password_authentication_enabled = false
          admin_ssh_keys                  = [local.first_public_key]
        }
      }
      network_interface {
        name    = "networkProTest"
        primary = true

        ip_configuration {
          name      = "TestIPConfiguration"
          primary   = true
          subnet_id = azurerm_subnet.linux_test.id

          public_ip_address {
            name                    = "TestPublicIPConfiguration"
            domain_name_label       = "test-domain-label"
            idle_timeout_in_minutes = 4
          }
        }
      }

      os_disk {
        storage_account_type = "Standard_LRS"
        caching              = "ReadWrite"
      }

      source_image_reference {
        publisher = "Canonical"
        offer     = "0001-com-ubuntu-server-jammy"
        sku       = "22_04-lts"
        version   = "latest"
      }

      extension_operations_enabled = true
      extension {
        name                               = "testOmsAgentForLinux"
        publisher                          = "Microsoft.EnterpriseCloud.Monitoring"
        type                               = "OmsAgentForLinux"
        type_handler_version               = "1.12"
        auto_upgrade_minor_version_enabled = false
        automatic_upgrade_enabled          = false
        failure_suppression_enabled        = false
      }

      extension {
        name                               = "CustomScript"
        publisher                          = "Microsoft.Azure.Extensions"
        type                               = "CustomScript"
        type_handler_version               = "2.0"
        auto_upgrade_minor_version_enabled = true
      }

      extension {
        name                 = "Docker"
        publisher            = "Microsoft.Azure.Extensions"
        type                 = "DockerExtension"
        type_handler_version = "1.0"

        auto_upgrade_minor_version_enabled        = false
        extensions_to_provision_after_vm_creation = ["CustomScript"]
        force_extension_execution_on_change       = "testUpdate"

        settings_json = jsonencode({
          "commandToExecute" = "echo $(date)"
        })

        protected_settings_json = jsonencode({
          "reset_ssh" = "True"
        })
      }
      extensions_time_budget = "PT1H"
    }
  }
}
`, r.baseAndAdditionalLocationLinuxTemplate(data), data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}
