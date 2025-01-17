// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azurefleet_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccAzureFleet_virtualMachineProfileAuth_authPassword(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authPassword(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleet_virtualMachineProfileAuth_authSSHKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authSSHKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureFleet_virtualMachineProfileAuth_authMultipleSSHKeys(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authMultipleSSHKeys(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureFleet_virtualMachineProfileAuth_authSSHKeyAndPassword(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authSSHKeyAndPassword(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleet_virtualMachineProfileAuth_authEd25519SSHKeys(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authEd25519SSHKeys(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r AzureFleetTestResource) authPassword(data acceptance.TestData) string {
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
    name = "Standard_DS1_v2"
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
    os_profile {
      linux_configuration {
        computer_name_prefix            = "prefix"
        admin_username                  = local.admin_username
        admin_password                  = local.admin_password
        password_authentication_enabled = true
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
  }
}
`, r.template(data, data.Locations.Primary), data.RandomInteger, data.Locations.Primary)
}

func (r AzureFleetTestResource) authSSHKey(data acceptance.TestData) string {
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
    name = "Standard_DS1_v2"
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
    os_profile {
      linux_configuration {
        computer_name_prefix = "prefix"
        admin_username       = local.admin_username
        admin_ssh_key {
          username   = local.admin_username
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
  }
}
`, r.template(data, data.Locations.Primary), data.RandomInteger, data.Locations.Primary)
}

func (r AzureFleetTestResource) authMultipleSSHKeys(data acceptance.TestData) string {
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
    name = "Standard_DS1_v2"
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
    os_profile {
      linux_configuration {
        computer_name_prefix = "prefix"
        admin_username       = local.admin_username
        admin_ssh_key {
          username   = local.admin_username
          public_key = local.first_public_key
        }
        admin_ssh_key {
          username   = local.admin_username
          public_key = local.second_public_key
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
  }
}
`, r.template(data, data.Locations.Primary), data.RandomInteger, data.Locations.Primary)
}

func (r AzureFleetTestResource) authSSHKeyAndPassword(data acceptance.TestData) string {
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
    name = "Standard_DS1_v2"
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
    os_profile {
      linux_configuration {
        computer_name_prefix = "prefix"
        admin_username       = local.admin_username
        admin_password       = local.admin_password
        admin_ssh_key {
          username   = local.admin_username
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
  }
}
`, r.template(data, data.Locations.Primary), data.RandomInteger, data.Locations.Primary)
}

func (r AzureFleetTestResource) authEd25519SSHKeys(data acceptance.TestData) string {
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
    name = "Standard_DS1_v2"
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

    os_profile {
      linux_configuration {
        computer_name_prefix = "prefix"
        admin_username       = local.admin_username
        admin_ssh_key {
          username   = local.admin_username
          public_key = local.ed25519_public_key
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
  }
}
`, r.template(data, data.Locations.Primary), data.RandomInteger, data.Locations.Primary)
}
