// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azurefleet_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccAzureFleetVirtualMachineProfileOsDisk_caching(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.ephemeral(data, data.Locations.Primary, "CacheDisk"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleetVirtualMachineProfileOsDisk_storageAccountTypePremiumLRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageAccountType(data, data.Locations.Primary, "Premium_LRS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleetVirtualMachineProfileOsDisk_storageAccountTypePremiumZRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageAccountType(data, "westeurope", "Premium_ZRS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleetVirtualMachineProfileOsDisk_storageAccountTypeStandardLRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageAccountType(data, data.Locations.Primary, "Standard_LRS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleetVirtualMachineProfileOsDisk_storageAccountTypeStandardSSDLRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageAccountType(data, data.Locations.Primary, "StandardSSD_LRS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleetVirtualMachineProfileOsDisk_storageAccountTypeStandardSSDZRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageAccountType(data, "westeurope", "StandardSSD_ZRS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func (r AzureFleetResource) ephemeral(data acceptance.TestData, location string, placement string) string {
	return fmt.Sprintf(`
%[1]s


resource "azurerm_azure_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"
  platform_fault_domain_count = 2

  regular_priority_profile {
    capacity     = 1
    min_capacity = 1
  }

  vm_sizes_profile {
    name = "Standard_F4s_v2"
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
      caching              = "ReadOnly"
      diff_disk_option     = "Local"
      diff_disk_placement  = "%[4]s"
    }

    source_image_reference {
      publisher = "Canonical"
      offer     = "0001-com-ubuntu-server-jammy"
      sku       = "22_04-lts"
      version   = "latest"
    }
  }
}
`, r.template(data, location), data.RandomInteger, location, placement)
}

func (r AzureFleetResource) storageAccountType(data acceptance.TestData, location string, storageAccountType string) string {
	return fmt.Sprintf(`
%[1]s


resource "azurerm_azure_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  platform_fault_domain_count = 2
  zones = ["1", "2", "3"]

  regular_priority_profile {
    capacity     = 1
    min_capacity = 1
  }

  vm_sizes_profile {
    name = "Standard_F2s_v2"
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
      storage_account_type = "%[4]s"
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
`, r.template(data, location), data.RandomInteger, location, storageAccountType)
}
