// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azurefleet_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccAzureFleetVirtualMachineProfileExtensions_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicExtensions(data, data.Locations.Primary),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleetVirtualMachineProfileExtensions_automaticUpgrade(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.extensionsAutomaticUpgrade(data, data.Locations.Primary),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleetVirtualMachineProfileExtensions_extensions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.extensions(data, data.Locations.Primary),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password", "virtual_machine_profile.0.extension.0.protected_settings_json", "virtual_machine_profile.0.extension.0.settings_json"),
		{
			Config: r.extensionsUpdate(data, data.Locations.Primary),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.extensions(data, data.Locations.Primary),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleetVirtualMachineProfileExtensions_extensionsMultiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleExtensions(data, data.Locations.Primary),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password", "virtual_machine_profile.0.extension.0.protected_settings_json", "virtual_machine_profile.0.extension.0.settings_json"),
	})
}

func TestAccAzureFleetVirtualMachineProfileExtensions_extensionsMultipleOnExistingOVMSS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleExtensionsProvisionMultipleExtensionOnExistingVMSS(data, data.Locations.Primary),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password", "virtual_machine_profile.0.extension.0.protected_settings_json", "virtual_machine_profile.0.extension.0.settings_json"),
	})
}

func TestAccAzureFleetVirtualMachineProfileExtensions_operations(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.extensionsOperationsEnabled(data, data.Locations.Primary),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password", "virtual_machine_profile.0.extension.0.protected_settings_json", "virtual_machine_profile.0.extension.0.settings_json"),
	})
}

func TestAccAzureFleetVirtualMachineProfileExtensions_protectedSettingsFromKeyVault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.extensionsProtectedSettingsFromKeyVault(data, data.Locations.Primary),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password", "virtual_machine_profile.0.extension.0.protected_settings_json"),
	})
}

func (r AzureFleetResource) extensionsAutomaticUpgrade(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_azure_fleet" "test" {
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

  virtual_machine_profile {
    os_profile {
      linux_configuration {
        computer_name_prefix = "prefix"
        admin_username       = "azureuser"
        password_authentication_enabled = true

        admin_ssh_key {
          username   = "azureuser"
          public_key = local.first_public_key
        }
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
			name                               = "testOmsAgentForLinux"
			publisher                          = "Microsoft.EnterpriseCloud.Monitoring"
			type                               = "OmsAgentForLinux"
			type_handler_version               = "1.12"
			auto_upgrade_minor_version_enabled 	= true
			automatic_upgrade_enabled 	= true
  	}
  }
}
`, r.linuxTemplate(data, location), data.RandomInteger, location)
}

func (r AzureFleetResource) extensionsOperationsEnabled(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
%[1]s


resource "azurerm_azure_fleet" "test" {
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

  virtual_machine_profile {
    os_profile {
      linux_configuration {
        computer_name_prefix = "prefix"
        admin_username       = "azureuser"

        admin_ssh_key {
          username   = "azureuser"
          public_key = local.first_public_key
        }
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
      name                               = "CustomScript"
      publisher                          = "Microsoft.Azure.Extensions"
      type                               = "CustomScript"
      type_handler_version               = "2.0"
      auto_upgrade_minor_version_enabled = true

      settings_json = jsonencode({
      "commandToExecute" = "echo $HOSTNAME"
    })

      protected_settings_json = jsonencode({
        "managedIdentity" = {}
     })
    }
  }
}
`, r.linuxTemplate(data, location), data.RandomInteger, location)
}

func (r AzureFleetResource) extensionsOperationsDisabled(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
%[1]s


resource "azurerm_azure_fleet" "test" {
  name                         = "acctest-fleet-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = "%[3]s"
  platform_fault_domain_count  = 2
  extension_operations_enabled = false

  regular_priority_profile {
    capacity     = 0
    min_capacity = 0
  }

  vm_sizes_profile {
    name = "Standard_D1_v2"
  }

  virtual_machine_profile {

    os_profile {
      linux_configuration {
        computer_name_prefix = "prefix"
        admin_username       = "azureuser"

        admin_ssh_key {
          username   = "azureuser"
          public_key = local.first_public_key
        }
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
  }
}
`, r.linuxTemplate(data, location), data.RandomInteger, location)
}

func (r AzureFleetResource) extensionsProtectedSettingsFromKeyVault(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      recover_soft_deleted_key_vaults       = false
      purge_soft_delete_on_destroy          = false
      purge_soft_deleted_keys_on_destroy    = false
      purge_soft_deleted_secrets_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {}

%[1]s

%[2]s

resource "azurerm_key_vault" "test" {
  name                   = "acctestkv%[3]s"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  tenant_id              = data.azurerm_client_config.current.tenant_id
  sku_name               = "standard"
  enabled_for_deployment = true

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    secret_permissions = [
      "Delete",
      "Get",
      "Set",
    ]
  }
}

resource "azurerm_key_vault_secret" "test" {
  name         = "secret"
  value        = "{\"commandToExecute\":\"echo $HOSTNAME\"}"
  key_vault_id = azurerm_key_vault.test.id
}

resource "azurerm_azure_fleet" "test" {
  name                        = "acctest-fleet-%[4]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = "%[5]s"
  platform_fault_domain_count = 2

  regular_priority_profile {
    capacity     = 0
    min_capacity = 0
  }

  vm_sizes_profile {
    name = "Standard_D1_v2"
  }

  virtual_machine_profile {
    os_profile {
      linux_configuration {
        computer_name_prefix = "prefix"
        admin_username       = "azureuser"

        admin_ssh_key {
          username   = "azureuser"
          public_key = local.first_public_key
        }
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
      name                               = "CustomScript"
      publisher                          = "Microsoft.Azure.Extensions"
      type                               = "CustomScript"
      type_handler_version               = "2.1"
      
			protected_settings_from_key_vault {
				secret_url      = azurerm_key_vault_secret.test.id
				source_vault_id = azurerm_key_vault.test.id
    	}
    }
  }
}
`, r.linuxPublicKeyTemplate(), r.templateWithOutProvider(data, location), data.RandomString, data.RandomInteger, location)
}

func (r AzureFleetResource) basicExtensions(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_azure_fleet" "test" {
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

  virtual_machine_profile {
    os_profile {
      linux_configuration {
        computer_name_prefix = "prefix"
        admin_username       = "azureuser"

        admin_ssh_key {
          username   = "azureuser"
          public_key = local.first_public_key
        }
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
      name                               = "CustomScript"
      publisher                          = "Microsoft.Azure.Extensions"
      type                               = "CustomScript"
      type_handler_version               = "2.0"
    }
  }
}
`, r.linuxTemplate(data, location), data.RandomInteger, location)
}

func (r AzureFleetResource) extensions(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_azure_fleet" "test" {
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

  virtual_machine_profile {
    os_profile {
      linux_configuration {
        computer_name_prefix = "prefix"
        admin_username       = "azureuser"

        admin_ssh_key {
          username   = "azureuser"
          public_key = local.first_public_key
        }
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

    extensions_time_budget = "PT30M"
    extension {
      name                               = "CustomScript"
      publisher                          = "Microsoft.Azure.Extensions"
      type                               = "CustomScript"
      type_handler_version               = "2.0"
      auto_upgrade_minor_version_enabled = true

      settings_json = jsonencode({
        "commandToExecute" = "echo $HOSTNAME"
        "timestamp"        = "1234567890"
      })

      protected_settings_json = jsonencode({
        "managedIdentity" = {}
      })
    }
  }
}
`, r.linuxTemplate(data, location), data.RandomInteger, location)
}

func (r AzureFleetResource) extensionsUpdate(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_azure_fleet" "test" {
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

  virtual_machine_profile {
    os_profile {
      linux_configuration {
        computer_name_prefix = "prefix"
        admin_username       = "azureuser"

        admin_ssh_key {
          username   = "azureuser"
          public_key = local.first_public_key
        }
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

    extensions_time_budget = "PT1H"
    extension {
      name                               = "CustomScript"
      publisher                          = "Microsoft.Azure.Extensions"
      type                               = "CustomScript"
      type_handler_version               = "2.0"
      auto_upgrade_minor_version_enabled = true
      force_extension_execution_on_change = "test"

      settings_json = jsonencode({
        "commandToExecute" = "echo $HOSTNAME"
      })

      protected_settings_json = jsonencode({
        "managedIdentity" = {}
      })
    }
  }
}
`, r.linuxTemplate(data, location), data.RandomInteger, location)
}

func (r AzureFleetResource) multipleExtensions(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
%[1]s


resource "azurerm_azure_fleet" "test" {
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

  virtual_machine_profile {
    os_profile {
      linux_configuration {
        computer_name_prefix = "prefix"
        admin_username       = "azureuser"

        admin_ssh_key {
          username   = "azureuser"
          public_key = local.first_public_key
        }
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
      name                               = "CustomScript"
      publisher                          = "Microsoft.Azure.Extensions"
      type                               = "CustomScript"
      type_handler_version               = "2.0"
      auto_upgrade_minor_version_enabled = true

      settings_json = jsonencode({
        "commandToExecute" = "echo $HOSTNAME"
      })

      protected_settings_json = jsonencode({
        "managedIdentity" = {}
      })
    }

    extension {
      name                               = "Docker"
      publisher                          = "Microsoft.Azure.Extensions"
      type                               = "DockerExtension"
      type_handler_version               = "1.0"
      auto_upgrade_minor_version_enabled = true
      extensions_to_provision_after_vm_creation = ["CustomScript"]
    }
  }
}
`, r.linuxTemplate(data, location), data.RandomInteger, location)
}

func (r AzureFleetResource) multipleExtensionsProvisionMultipleExtensionOnExistingVMSS(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_azure_fleet" "test" {
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

  virtual_machine_profile {
    os_profile {
      linux_configuration {
        computer_name_prefix            = "prefix"
        admin_username                  = "azureuser"
        admin_password                  = "P@ssw0rd1234!"
        password_authentication_enabled = true
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
      name                               = "CustomScript"
      publisher                          = "Microsoft.Azure.Extensions"
      type                               = "CustomScript"
      type_handler_version               = "2.0"
      auto_upgrade_minor_version_enabled = true

      settings_json = jsonencode({
        "commandToExecute" = "echo $HOSTNAME"
      })

      protected_settings_json = jsonencode({
        "managedIdentity" = {}
      })
    }

    extension {
      name                                      = "Docker"
      publisher                                 = "Microsoft.Azure.Extensions"
      type                                      = "DockerExtension"
      type_handler_version                      = "1.0"
      auto_upgrade_minor_version_enabled        = true
      extensions_to_provision_after_vm_creation = ["CustomScript"]
    }
  }
}
`, r.template(data, location), data.RandomInteger, location)
}
