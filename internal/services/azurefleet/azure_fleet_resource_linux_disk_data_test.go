// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azurefleet_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccFleetLinux_disksDataDiskBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.linuxTemplate(data, data.Locations.Primary, r.disksDataDiskBasicVirtualMachineProfile()),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("compute_profile.0.virtual_machine_profile.0.os_profile.0.admin_password"),
	})
}

func TestAccFleetLinux_disksDataDiskCaching(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.linuxTemplate(data, data.Locations.Primary, r.disksDataDiskCachingVirtualMachineProfile()),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("compute_profile.0.virtual_machine_profile.0.os_profile.0.admin_password"),
	})
}

func TestAccFleetLinux_disksDataDiskDiskEncryptionSet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.linuxTemplateWithoutProvider(data, data.Locations.Primary, r.disksDataDisk_diskEncryptionSetVirtualMachineProfile(), r.disksDataDisk_diskEncryptionSetResource(data)),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("compute_profile.0.virtual_machine_profile.0.os_profile.0.admin_password"),
	})
}

func TestAccFleetLinux_disksDataDiskMultiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.linuxTemplate(data, data.Locations.Primary, r.disksDataDiskMultipleVirtualMachineProfile()),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("compute_profile.0.virtual_machine_profile.0.os_profile.0.admin_password"),
	})
}

func TestAccFleetLinux_disksDataDiskStorageAccountTypeUltraSSDLRSWithIOPSAndMBPS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}
	// Are supported in East US 2, SouthEast Asia, and North Europe, in two availability zones per region

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// Are supported in East US 2, SouthEast Asia, and North Europe, in two availability zones per region
			Config: r.linuxTemplate(data, "eastus2", r.disksDataDiskStorageAccountTypeUltraSSDLRSWithIOPSAndMBPSVirtualMachineProfile()),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("compute_profile.0.virtual_machine_profile.0.os_profile.0.admin_password"),
	})
}

func TestAccFleetLinux_disksDataDiskWriteAcceleratorEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.linuxTemplate(data, data.Locations.Primary, r.disksDataDiskStorageAccountTypeUltraSSDLRSWithIOPSAndMBPSVirtualMachineProfile()),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("compute_profile.0.virtual_machine_profile.0.os_profile.0.admin_password"),
	})
}

func (r AzureFleetResource) disksDataDiskBasicVirtualMachineProfile() string {
	return `
virtual_machine_profile {
    storage_profile {
        image_reference {
          publisher = "Canonical"
          offer     = "0001-com-ubuntu-server-jammy"
          sku       = "22_04-lts"
          version   = "latest"
        }

        os_disk {
          caching = "ReadWrite"
          create_option = "FromImage"
          os_type       = "Linux"
          managed_disk {
            storage_account_type = "Standard_LRS"
          }
        }

		data_disk {
			managed_disk {
            storage_account_type = "Premium_LRS"
            }
			disk_size_in_gb         = 10
			lun                  = 10
		  }
      }

      os_profile {
        computer_name_prefix = "prefix"
        admin_username       = "azureuser"
        admin_password       = "P@ssw0rd1234!"
        linux_configuration {
          password_authentication_enabled = true
          ssh_keys {
			username   = "azureuser"
			public_key = local.first_public_key
		  }
        }
      }

      network_interface {
        name                           = "networkProTest"
        accelerated_networking_enabled = false
        ip_forwarding_enabled          = true
        ip_configuration {
          name                                   = "ipConfigTest"
          load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
          primary                                = true
          subnet_id                              = azurerm_subnet.test.id
        }
        primary = true
      }

      network_api_version = "2020-11-01"
    }
`
}

func (r AzureFleetResource) disksDataDiskCachingVirtualMachineProfile() string {
	return `
virtual_machine_profile {
    storage_profile {
        image_reference {
          publisher = "Canonical"
          offer     = "0001-com-ubuntu-server-jammy"
          sku       = "22_04-lts"
          version   = "latest"
        }

        os_disk {
          caching = "ReadWrite"
          create_option = "FromImage"
          os_type       = "Linux"
          managed_disk {
            storage_account_type = "Standard_LRS"
          }
        }

		data_disk {
			managed_disk {
            storage_account_type = "Standard_LRS"
            }
            caching              = "ReadOnly"
			disk_size_in_gb         = 10
			lun                  = 10
		  }
      }

      os_profile {
        computer_name_prefix = "prefix"
        admin_username       = "azureuser"
        admin_password       = "P@ssw0rd1234!"
        linux_configuration {
          password_authentication_enabled = true
          ssh_keys {
			username   = "azureuser"
			public_key = local.first_public_key
		  }
        }
      }

      network_interface {
        name                           = "networkProTest"
        accelerated_networking_enabled = false
        ip_forwarding_enabled          = true
        ip_configuration {
          name                                   = "ipConfigTest"
          load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
          primary                                = true
          subnet_id                              = azurerm_subnet.test.id
        }
        primary = true
      }

      network_api_version = "2020-11-01"
    }
`
}

func (r AzureFleetResource) disksDataDisk_diskEncryptionSetVirtualMachineProfile() string {
	return `
virtual_machine_profile {
    storage_profile {
        image_reference {
          publisher = "Canonical"
          offer     = "0001-com-ubuntu-server-jammy"
          sku       = "22_04-lts"
          version   = "latest"
        }

        os_disk {
          caching = "ReadWrite"
          create_option = "FromImage"
          os_type       = "Linux"
          managed_disk {
            storage_account_type = "Standard_LRS"
          }
        }

		data_disk {
			managed_disk {
            storage_account_type = "Standard_LRS"
            }
        caching                = "ReadWrite"
        disk_encryption_set_id = azurerm_disk_encryption_set.test.id
			disk_size_in_gb         = 10
			lun                  = 10
		  }
      }

      os_profile {
        computer_name_prefix = "prefix"
        admin_username       = "azureuser"
        admin_password       = "P@ssw0rd1234!"
        linux_configuration {
          password_authentication_enabled = true
          ssh_keys {
			username   = "azureuser"
			public_key = local.first_public_key
		  }
        }
      }

      network_interface {
        name                           = "networkProTest"
        accelerated_networking_enabled = false
        ip_forwarding_enabled          = true
        ip_configuration {
          name                                   = "ipConfigTest"
          load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
          primary                                = true
          subnet_id                              = azurerm_subnet.test.id
        }
        primary = true
      }

      network_api_version = "2020-11-01"
    }
`
}

func (r AzureFleetResource) disksDataDiskMultipleVirtualMachineProfile() string {
	return `
virtual_machine_profile {
    storage_profile {
        image_reference {
          publisher = "Canonical"
          offer     = "0001-com-ubuntu-server-jammy"
          sku       = "22_04-lts"
          version   = "latest"
        }

        os_disk {
          caching = "ReadWrite"
          create_option = "FromImage"
          os_type       = "Linux"
          managed_disk {
            storage_account_type = "Standard_LRS"
          }
        }

		data_disk {
			managed_disk {
            storage_account_type = "Standard_LRS"
            }
        caching                = "ReadWrite"
       
			disk_size_in_gb         = 10
			lun                  = 10
		  }

         data_disk {
			managed_disk {
            storage_account_type = "Standard_LRS"
            }
        caching                = "ReadWrite"
       
			disk_size_in_gb         = 10
			lun                  = 20
		  }
      }

      os_profile {
        computer_name_prefix = "prefix"
        admin_username       = "azureuser"
        admin_password       = "P@ssw0rd1234!"
        linux_configuration {
          password_authentication_enabled = true
          ssh_keys {
			username   = "azureuser"
			public_key = local.first_public_key
		  }
        }
      }

      network_interface {
        name                           = "networkProTest"
        accelerated_networking_enabled = false
        ip_forwarding_enabled          = true
        ip_configuration {
          name                                   = "ipConfigTest"
          load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
          primary                                = true
          subnet_id                              = azurerm_subnet.test.id
        }
        primary = true
      }

      network_api_version = "2020-11-01"
    }
`
}

func (r AzureFleetResource) disksDataDiskStorageAccountTypeUltraSSDLRSWithIOPSAndMBPSVirtualMachineProfile() string {
	return `
additional_capabilities_ultra_ssd_enabled = true

virtual_machine_profile {
    storage_profile {
        image_reference {
          publisher = "Canonical"
          offer     = "0001-com-ubuntu-server-jammy"
          sku       = "22_04-lts"
          version   = "latest"
        }

        os_disk {
          caching = "ReadWrite"
          create_option = "FromImage"
          os_type       = "Linux"
          managed_disk {
            storage_account_type = "Standard_LRS"
          }
        }

		data_disk {
			managed_disk {
            storage_account_type = "UltraSSD_LRS"
            }
			disk_size_in_gb         = 10
			lun                  = 10
            disk_iops_read_write = 101
            disk_mbps_read_write = 11
		  }
      }

      os_profile {
        computer_name_prefix = "prefix"
        admin_username       = "azureuser"
        admin_password       = "P@ssw0rd1234!"
        linux_configuration {
          password_authentication_enabled = true
          ssh_keys {
			username   = "azureuser"
			public_key = local.first_public_key
		  }
        }
      }

      network_interface {
        name                           = "networkProTest"
        accelerated_networking_enabled = false
        ip_forwarding_enabled          = true
        ip_configuration {
          name                                   = "ipConfigTest"
          load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
          primary                                = true
          subnet_id                              = azurerm_subnet.test.id
        }
        primary = true
      }

      network_api_version = "2020-11-01"
    }
`
}

func (r AzureFleetResource) disksDataDiskWriteAcceleratorEnabledVirtualMachineProfile() string {
	return `
virtual_machine_profile {
    storage_profile {
        image_reference {
          publisher = "Canonical"
          offer     = "0001-com-ubuntu-server-jammy"
          sku       = "22_04-lts"
          version   = "latest"
        }

        os_disk {
          caching = "ReadWrite"
          create_option = "FromImage"
          os_type       = "Linux"
          managed_disk {
            storage_account_type = "Standard_LRS"
          }
        }

		data_disk {
			managed_disk {
            storage_account_type = "Premium_LRS"
            }
			disk_size_in_gb         = 10
			lun                  = 10
write_accelerator_enabled = true
		  }
      }

      os_profile {
        computer_name_prefix = "prefix"
        admin_username       = "azureuser"
        admin_password       = "P@ssw0rd1234!"
        linux_configuration {
          password_authentication_enabled = true
          ssh_keys {
			username   = "azureuser"
			public_key = local.first_public_key
		  }
        }
      }

      network_interface {
        name                           = "networkProTest"
        accelerated_networking_enabled = false
        ip_forwarding_enabled          = true
        ip_configuration {
          name                                   = "ipConfigTest"
          load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
          primary                                = true
          subnet_id                              = azurerm_subnet.test.id
        }
        primary = true
      }

      network_api_version = "2020-11-01"
    }
`
}

func (AzureFleetResource) disksDataDisk_diskEncryptionSetDependencies(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      recover_soft_deleted_key_vaults    = false
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                        = "acctestkv%s"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  sku_name                    = "standard"
  purge_protection_enabled    = true
  enabled_for_disk_encryption = true
}

resource "azurerm_key_vault_access_policy" "service-principal" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "Purge",
    "Update",
    "GetRotationPolicy",
  ]

  secret_permissions = [
    "Get",
    "Delete",
    "Set",
  ]
}

resource "azurerm_key_vault_key" "test" {
  name         = "examplekey"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]

  depends_on = ["azurerm_key_vault_access_policy.service-principal"]
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (r AzureFleetResource) disksDataDisk_diskEncryptionSetResource(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_disk_encryption_set" "test" {
  name                = "acctestdes-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  key_vault_key_id    = azurerm_key_vault_key.test.id

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_key_vault_access_policy" "disk-encryption" {
  key_vault_id = azurerm_key_vault.test.id

  key_permissions = [
    "Get",
    "WrapKey",
    "UnwrapKey",
    "GetRotationPolicy",
  ]

  tenant_id = azurerm_disk_encryption_set.test.identity.0.tenant_id
  object_id = azurerm_disk_encryption_set.test.identity.0.principal_id
}

resource "azurerm_role_assignment" "disk-encryption-read-keyvault" {
  scope                = azurerm_key_vault.test.id
  role_definition_name = "Reader"
  principal_id         = azurerm_disk_encryption_set.test.identity.0.principal_id
}
`, r.disksDataDisk_diskEncryptionSetDependencies(data), data.RandomInteger)
}
