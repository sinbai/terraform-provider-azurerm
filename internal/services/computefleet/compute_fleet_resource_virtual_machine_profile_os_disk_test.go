// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package computefleet_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccComputeFleet_virtualMachineProfileOsDisk_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.osDiskBasic(data, data.Locations.Primary),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_virtualMachineProfileOsDisk_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.osDiskComplete(data, data.Locations.Primary),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_virtualMachineProfileOsDisk_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.osDiskComplete(data, "westeurope"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			// Limited regional availability for some storage account type
			Config: r.osDiskCompleteUpdate(data, "westeurope"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.osDiskBasic(data, "westeurope"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.osDiskComplete(data, "westeurope"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func (r ComputeFleetTestResource) osDiskBasic(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_compute_fleet" "test" {
  name                        = "acctest-fleet-%[2]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = "%[3]s"
  platform_fault_domain_count = 1

  regular_priority_profile {
    capacity     = 1
    min_capacity = 0
  }

  vm_sizes_profile {
    name = "Standard_F4s_v2"
  }

  virtual_machine_profile {
    network_api_version = "2020-11-01"
    os_profile {
      linux_configuration {
        computer_name_prefix            = "prefix"
        admin_username                  = local.admin_username
        admin_password                  = local.admin_password
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

    os_disk {}

    source_image_reference {
      publisher = "Canonical"
      offer     = "0001-com-ubuntu-server-jammy"
      sku       = "22_04-lts"
      version   = "latest"
    }
  }
}
`, r.template(data, location), data.RandomInteger, location)
}

func (r ComputeFleetTestResource) osDiskComplete(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
%[1]s
%[2]s
resource "azurerm_compute_fleet" "test" {
  name                        = "acctest-fleet-%[3]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = "%[4]s"
  platform_fault_domain_count = 1

  regular_priority_profile {
    capacity     = 0
    min_capacity = 0
  }

  vm_sizes_profile {
    name = "Standard_DC8ads_v5"
  }

  virtual_machine_profile {
    network_api_version = "2020-11-01"
    secure_boot_enabled = true
    vtpm_enabled        = true

    os_profile {
      linux_configuration {
        computer_name_prefix = "prefix"
        admin_username       = local.admin_username
        admin_password       = local.admin_password
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
      caching                   = "ReadOnly"
      delete_option             = "Delete"
      diff_disk_option          = "Local"
      diff_disk_placement       = "ResourceDisk"
      disk_size_in_gb           = 30
      storage_account_type      = "Premium_LRS"
      security_encryption_type  = "DiskWithVMGuestState"
      disk_encryption_set_id    = azurerm_disk_encryption_set.test.id
      write_accelerator_enabled = false
    }

    source_image_reference {
      publisher = "Canonical"
      offer     = "0001-com-ubuntu-confidential-vm-jammy"
      sku       = "22_04-lts-cvm"
      version   = "latest"
    }
  }
}
`, r.diskEncryptionSetResourceDependencies(data), r.templateWithOutProvider(data, location), data.RandomInteger, location)
}

func (r ComputeFleetTestResource) osDiskCompleteUpdate(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
%[1]s
%[2]s
resource "azurerm_compute_fleet" "test" {
  name                        = "acctest-fleet-%[3]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = "%[4]s"
  platform_fault_domain_count = 1

  regular_priority_profile {
    capacity     = 0
    min_capacity = 0
  }

  vm_sizes_profile {
    name = "Standard_DC8ads_v5"
  }

  virtual_machine_profile {
    network_api_version = "2020-11-01"
    secure_boot_enabled = true
    vtpm_enabled        = true

    os_profile {
      linux_configuration {
        computer_name_prefix = "prefix"
        admin_username       = local.admin_username
        admin_password       = local.admin_password
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
      caching                   = "ReadOnly"
      delete_option             = "Delete"
      diff_disk_option          = "Local"
      diff_disk_placement       = "ResourceDisk"
      disk_size_in_gb           = 50
      storage_account_type      = "Premium_ZRS"
      security_encryption_type  = "DiskWithVMGuestState"
      disk_encryption_set_id    = azurerm_disk_encryption_set.test.id
      write_accelerator_enabled = false
    }

    source_image_reference {
      publisher = "Canonical"
      offer     = "0001-com-ubuntu-confidential-vm-jammy"
      sku       = "22_04-lts-cvm"
      version   = "latest"
    }
  }
}
`, r.diskEncryptionSetResourceDependencies(data), r.templateWithOutProvider(data, location), data.RandomInteger, location)
}
