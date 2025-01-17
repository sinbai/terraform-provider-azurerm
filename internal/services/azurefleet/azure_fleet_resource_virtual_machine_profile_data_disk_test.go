// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azurefleet_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccAzureFleet_virtualMachineProfileDataDisk_basicManagedDisk(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicManagedDisk(data, data.Locations.Primary),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleet_virtualMachineProfileDataDisk_basicManagedDiskWithZones(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicManagedDiskWithZones(data, data.Locations.Primary),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleet_virtualMachineProfileDataDisk_loadBalancerManagedDataDisks(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.loadBalancerManagedDataDisks(data, data.Locations.Primary),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleet_virtualMachineProfileDataDisk_diskEncryptionSet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.diskEncryptionSet(data, data.Locations.Primary),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleet_virtualMachineProfileDataDisk_storageAccountTypePremiumLRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dataDiskStorageAccountType(data, data.Locations.Primary, "Premium_LRS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleet_virtualMachineProfileDataDisk_storageAccountTypePremiumV2LRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// Limited regional availability for `PremiumV2_LRS`
			// `PremiumV2_LRS` disks can only be can only be attached to zonal VMs currently
			Config: r.dataDiskStorageAccountTypePremiumV2LRS(data, "westeurope"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleet_virtualMachineProfileDataDisk_storageAccountTypePremiumZRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dataDiskStorageAccountType(data, "westeurope", "Premium_ZRS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleet_virtualMachineProfileDataDisk_storageAccountTypeStandardLRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dataDiskStorageAccountType(data, data.Locations.Primary, "Standard_LRS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleet_virtualMachineProfileDataDisk_storageAccountTypeStandardSSDLRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dataDiskStorageAccountType(data, data.Locations.Primary, "StandardSSD_LRS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleet_virtualMachineProfileDataDisk_storageAccountTypeStandardSSDZRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dataDiskStorageAccountType(data, "westeurope", "StandardSSD_ZRS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleet_virtualMachineProfileDataDisk_storageAccountTypeUltraSSDLRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// Limited regional availability for `UltraSSD_LRS`
			// `UltraSSD_LRS` disks needs to be used with `additional_capabilities_ultra_ssd_enabled` set to `true`
			Config: r.dataDiskStorageAccountTypeUltraSSDLRS(data, "eastus2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleet_virtualMachineProfileDataDisk_MarketPlaceImage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dataDiskMarketPlaceImage(data, data.Locations.Primary),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func (r AzureFleetTestResource) basicManagedDisk(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
%[1]s


resource "azurerm_azure_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  regular_priority_profile {
    capacity     = 1
    min_capacity = 1
  }

  vm_sizes_profile {
    name = "Standard_D1_v2"
  }

  virtual_machine_profile {
    source_image_reference {
      publisher = "Canonical"
      offer     = "0001-com-ubuntu-server-jammy"
      sku       = "22_04-lts"
      version   = "latest"
    }

    os_disk {
      caching              = "ReadWrite"
      storage_account_type = "Standard_LRS"
    }

    data_disk {
      lun                  = 0
      caching              = "ReadWrite"
      create_option        = "Empty"
      disk_size_in_gb      = 10
      storage_account_type = "Standard_LRS"
    }

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
        subnet_id = azurerm_subnet.test.id
        primary   = true
        public_ip_address {
          name                    = "TestPublicIPConfiguration"
          domain_name_label       = "test-domain-label"
          idle_timeout_in_minutes = 4
        }
      }
    }
  }
}
`, r.template(data, location), data.RandomInteger, location)
}

func (r AzureFleetTestResource) diskEncryptionSet(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "azurerm_azure_fleet" "test" {
  name                = "acctest-fleet-%[3]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[4]s"

  regular_priority_profile {
    capacity     = 1
    min_capacity = 1
  }

  vm_sizes_profile {
    name = "Standard_D1_v2"
  }

  virtual_machine_profile {
    source_image_reference {
      publisher = "Canonical"
      offer     = "0001-com-ubuntu-server-jammy"
      sku       = "22_04-lts"
      version   = "latest"
    }

    os_disk {
      caching              = "ReadWrite"
      storage_account_type = "Standard_LRS"
    }

    data_disk {
      lun                    = 10
      caching                = "ReadWrite"
      create_option          = "Empty"
      disk_size_in_gb        = 10
      storage_account_type   = "Standard_LRS"
      disk_encryption_set_id = azurerm_disk_encryption_set.test.id
    }

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
        subnet_id = azurerm_subnet.test.id
        primary   = true
        public_ip_address {
          name                    = "TestPublicIPConfiguration"
          domain_name_label       = "test-domain-label"
          idle_timeout_in_minutes = 4
        }
      }
	  }
  }

  depends_on = [
      "azurerm_role_assignment.disk-encryption-read-keyvault",
      "azurerm_key_vault_access_policy.disk-encryption",
  ]
}
`, r.diskEncryptionSetResourceDependencies(data), r.templateWithOutProvider(data, location), data.RandomInteger, location)
}

func (r AzureFleetTestResource) basicManagedDiskWithZones(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
%[1]s


resource "azurerm_azure_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  zones = ["1", "2", "3"]

  regular_priority_profile {
    capacity     = 1
    min_capacity = 1
  }

  vm_sizes_profile {
    name = "Standard_D1_v2"
  }

  virtual_machine_profile {
    source_image_reference {
      publisher = "Canonical"
      offer     = "0001-com-ubuntu-server-jammy"
      sku       = "22_04-lts"
      version   = "latest"
    }

    os_disk {
      caching              = "ReadWrite"
      storage_account_type = "Standard_LRS"
    }

    data_disk {
      lun                  = 0
      caching              = "ReadWrite"
      create_option        = "Empty"
      disk_size_in_gb      = 10
      storage_account_type = "Standard_LRS"
    }

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
        subnet_id = azurerm_subnet.test.id
        primary   = true
        public_ip_address {
          name                    = "TestPublicIPConfiguration"
          domain_name_label       = "test-domain-label"
          idle_timeout_in_minutes = 4
        }
      }
    }
  }
}
`, r.template(data, location), data.RandomInteger, location)
}

func (r AzureFleetTestResource) loadBalancerManagedDataDisks(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
%[1]s


resource "azurerm_azure_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  zones = ["1", "2", "3"]

  regular_priority_profile {
    capacity     = 1
    min_capacity = 1
  }

  vm_sizes_profile {
    name = "Standard_D1_v2"
  }

  virtual_machine_profile {
    source_image_reference {
      publisher = "Canonical"
      offer     = "0001-com-ubuntu-server-jammy"
      sku       = "22_04-lts"
      version   = "latest"
    }

    os_disk {
      caching              = "ReadWrite"
      storage_account_type = "Standard_LRS"
    }

    data_disk {
      storage_account_type = "Standard_LRS"
      caching              = "ReadOnly"
      disk_size_in_gb      = 10
      lun                  = 10
    }
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
        name                                   = "TestIPConfiguration"
        subnet_id                              = azurerm_subnet.test.id
        primary                                = true
        load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
      }
    }
  }
}
`, r.template(data, location), data.RandomInteger, location)
}

func (r AzureFleetTestResource) dataDiskStorageAccountType(data acceptance.TestData, location string, storageAccountType string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_azure_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  zones = ["1", "2", "3"]

  regular_priority_profile {
    capacity     = 1
    min_capacity = 1
  }

  vm_sizes_profile {
    name = "Standard_F2s_v2"
  }

  virtual_machine_profile {
    source_image_reference {
      publisher = "Canonical"
      offer     = "0001-com-ubuntu-server-jammy"
      sku       = "22_04-lts"
      version   = "latest"
    }

    os_disk {
      caching              = "ReadWrite"
      storage_account_type = "Standard_LRS"
    }

    data_disk {
      lun                  = 0
      caching              = "ReadWrite"
      create_option        = "Empty"
      disk_size_in_gb      = 10
      storage_account_type = "%[4]s"
    }
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
        name                                   = "TestIPConfiguration"
        subnet_id                              = azurerm_subnet.test.id
        primary                                = true
        load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
      }
    }
  }
}
`, r.template(data, location), data.RandomInteger, location, storageAccountType)
}

func (r AzureFleetTestResource) dataDiskStorageAccountTypePremiumV2LRS(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
%[1]s


resource "azurerm_azure_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  zones = ["1"]

  regular_priority_profile {
    capacity     = 1
    min_capacity = 1
  }

  vm_sizes_profile {
    name = "Standard_F2s_v2"
  }

  virtual_machine_profile {
    source_image_reference {
      publisher = "Canonical"
      offer     = "0001-com-ubuntu-server-jammy"
      sku       = "22_04-lts"
      version   = "latest"
    }

    os_disk {
      caching              = "ReadWrite"
      storage_account_type = "Standard_LRS"
    }

    data_disk {
      lun                  = 0
      create_option        = "Empty"
      disk_size_in_gb      = 10
      storage_account_type = "PremiumV2_LRS"
    }

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
        name                                   = "TestIPConfiguration"
        subnet_id                              = azurerm_subnet.test.id
        primary                                = true
        load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
      }
    }
  }
}
`, r.template(data, location), data.RandomInteger, location)
}

func (r AzureFleetTestResource) dataDiskMarketPlaceImage(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
%[1]s


resource "azurerm_azure_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  regular_priority_profile {
    capacity     = 1
    min_capacity = 1
  }

  vm_sizes_profile {
    name = "Standard_F2"
  }

  plan {
    name      = "arcsight_logger_72_byol"
    product   = "arcsight-logger"
    publisher = "micro-focus"
  }

  virtual_machine_profile {
    source_image_reference {
      publisher = "micro-focus"
      offer     = "arcsight-logger"
      sku       = "arcsight_logger_72_byol"
      version   = "7.2.0"
    }

    os_disk {
      caching              = "ReadWrite"
      storage_account_type = "Standard_LRS"
    }

    data_disk {
      caching              = "ReadWrite"
      disk_size_in_gb      = 900
      create_option        = "FromImage"
      storage_account_type = "Standard_LRS"
    }

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
        name                                   = "TestIPConfiguration"
        subnet_id                              = azurerm_subnet.test.id
        primary                                = true
        load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
      }
    }
  }
}
`, r.template(data, location), data.RandomInteger, location)
}

func (r AzureFleetTestResource) dataDiskStorageAccountTypeUltraSSDLRS(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_azure_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  regular_priority_profile {
    capacity     = 1
    min_capacity = 1
  }

  vm_sizes_profile {
    name = "Standard_F2s_v2"
  }
  additional_capabilities_ultra_ssd_enabled = true

  zones = ["1"]

  virtual_machine_profile {
    source_image_reference {
      publisher = "Canonical"
      offer     = "0001-com-ubuntu-server-jammy"
      sku       = "22_04-lts"
      version   = "latest"
    }

    os_disk {
      caching              = "ReadWrite"
      storage_account_type = "Standard_LRS"
    }

    data_disk {
      storage_account_type = "UltraSSD_LRS"
      disk_size_in_gb      = 10
      create_option        = "Empty"
      lun                  = 0
    }

    os_profile {
      linux_configuration {
        computer_name_prefix            = "prefix"
        admin_username                  = local.admin_username
        admin_password                  = local.admin_password
        password_authentication_enabled = true
        admin_ssh_key {
          username   = local.admin_username
          public_key = local.first_public_key
        }
      }
    }

    network_interface {
      name = "networkProTest"
      ip_configuration {
        name                                   = "ipConfigTest"
        load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
        primary                                = true
        subnet_id                              = azurerm_subnet.test.id
      }
      primary = true
    }
  }
}
`, r.template(data, location), data.RandomInteger, location)
}

func (r AzureFleetTestResource) diskEncryptionSetResourceDependencies(data acceptance.TestData) string {
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
`, data.RandomString, data.RandomInteger)
}
