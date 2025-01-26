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
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password"),
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
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password"),
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
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			// Limited regional availability for some storage account type
			Config: r.osDiskCompleteUpdate(data, "westeurope"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.osDiskBasic(data, "westeurope"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.osDiskComplete(data, "westeurope"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password"),
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

  compute_api_version = "2024-03-01"
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
  additional_location_profile {
    location = "%[4]s"
    virtual_machine_profile_override {
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
          subnet_id = azurerm_subnet.linux_test.id

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
}
`, r.osDiskTemplate(data, location), data.RandomInteger, location, data.Locations.Secondary)
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

  compute_api_version = "2024-03-01"
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

  additional_location_profile {
    location = "%[4]s"
    virtual_machine_profile_override {
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
          subnet_id = azurerm_subnet.linux_test.id

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
        disk_encryption_set_id    = azurerm_disk_encryption_set.linux_test.id
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
}
`, r.diskEncryptionSetResourceDependencies(data), r.osDiskTemplateWithOutProvider(data, location), data.RandomInteger, location, data.Locations.Secondary)
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

  compute_api_version = "2024-03-01"
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

  additional_location_profile {
    location = "%[4]s"
    virtual_machine_profile_override {
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
          subnet_id = azurerm_subnet.linux_test.id

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
        disk_encryption_set_id    = azurerm_disk_encryption_set.linux_test.id
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
}
`, r.diskEncryptionSetResourceDependencies(data), r.osDiskTemplateWithOutProvider(data, location), data.RandomInteger, location, data.Locations.Secondary)
}

func (r ComputeFleetTestResource) osDiskTemplate(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

`, r.osDiskTemplateWithOutProvider(data, location), data.RandomInteger, location, data.Locations.Secondary)
}

func (r ComputeFleetTestResource) osDiskTemplateWithOutProvider(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
locals {
  first_public_key          = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+wWK73dCr+jgQOAxNsHAnNNNMEMWOHYEccp6wJm2gotpr9katuF/ZAdou5AaW1C61slRkHRkpRRX9FA9CYBiitZgvCCz+3nWNN7l/Up54Zps/pHWGZLHNJZRYyAB6j5yVLMVHIHriY49d/GZTZVNB8GoJv9Gakwc/fuEZYYl4YDFiGMBP///TzlI4jhiJzjKnEvqPFki5p2ZRJqcbCiF4pJrxUQR/RXqVFQdbRLZgYfJ8xGB878RENq3yQ39d8dVOkq4edbkzwcUmwwwkYVPIoDGsYLaRHnG+To7FvMeyO7xDVQkMKzopTQV8AuKpyvpqu0a9pWOMaiCyDytO7GGN you@me.com"
  second_public_key         = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0/NDMj2wG6bSa6jbn6E3LYlUsYiWMp1CQ2sGAijPALW6OrSu30lz7nKpoh8Qdw7/A4nAJgweI5Oiiw5/BOaGENM70Go+VM8LQMSxJ4S7/8MIJEZQp5HcJZ7XDTcEwruknrd8mllEfGyFzPvJOx6QAQocFhXBW6+AlhM3gn/dvV5vdrO8ihjET2GoDUqXPYC57ZuY+/Fz6W3KV8V97BvNUhpY5yQrP5VpnyvvXNFQtzDfClTvZFPuoHQi3/KYPi6O0FSD74vo8JOBZZY09boInPejkm9fvHQqfh0bnN7B6XJoUwC1Qprrx+XIy7ust5AEn5XL7d4lOvcR14MxDDKEp you@me.com"
  first_ed25519_public_key  = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIDqzSi9IHoYnbE3YQ+B2fQEVT8iGFemyPovpEtPziIVB you@me.com"
  second_ed25519_public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIDqzSi9IHoYnbE3YQ+B2fQEVT8iGFemyPovpEtPziIVB hello@world.com"
  admin_username            = "testadmin1234"
  admin_password            = "Password1234!"
  admin_password_update     = "Password1234!Update"
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-fleet-%[1]d"
  location = "%[2]s"
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

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicIP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  zones               = ["1"]
}

resource "azurerm_lb" "test" {
  name                = "acctest-loadbalancer-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "internal"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  name            = "internal"
  loadbalancer_id = azurerm_lb.test.id
}


resource "azurerm_resource_group" "linux_test" {
  name     = "acctest-rg-fleet-al-%[1]d"
  location = "%[3]s"
}

resource "azurerm_virtual_network" "linux_test" {
  name                = "acctvn-al-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.linux_test.location
  resource_group_name = azurerm_resource_group.linux_test.name
}

resource "azurerm_subnet" "linux_test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.linux_test.name
  virtual_network_name = azurerm_virtual_network.linux_test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_public_ip" "linux_test" {
  name                = "acctestpublicIP%[1]d"
  location            = azurerm_resource_group.linux_test.location
  resource_group_name = azurerm_resource_group.linux_test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  zones               = ["1"]
}

resource "azurerm_lb" "linux_test" {
  name                = "acctest-loadbalancer-%[1]d"
  location            = azurerm_resource_group.linux_test.location
  resource_group_name = azurerm_resource_group.linux_test.name
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "internal"
    public_ip_address_id = azurerm_public_ip.linux_test.id
  }
}

resource "azurerm_lb_backend_address_pool" "linux_test" {
  name            = "internal"
  loadbalancer_id = azurerm_lb.linux_test.id
}

`, data.RandomInteger, location, data.Locations.Secondary)
}
