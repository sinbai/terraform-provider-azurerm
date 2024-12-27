// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azurefleet_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccFleetLinux_authPassword(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authenticationTest(data, r.authPasswordVirtualMachineProfile()),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("compute_profile.0.virtual_machine_profile.0.os_profile.0.admin_password"),
	})
}

func TestAccFleetLinux_authSSHKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authenticationTest(data, r.authSSHKeyVirtualMachineProfile()),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFleetLinux_authSSHKeyAndPassword(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authenticationTest(data, r.authSSHKeyAndPasswordVirtualMachineProfile()),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("compute_profile.0.virtual_machine_profile.0.os_profile.0.admin_password"),
	})
}

func TestAccFleetLinux_authMultipleSSHKeys(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authenticationTest(data, r.authMultipleSSHKeysVirtualMachineProfile()),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFleetLinux_authUpdatingSSHKeys(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		//{
		//	Config: r.authSSHKey(data),
		//	Check: acceptance.ComposeTestCheckFunc(
		//		check.That(data.ResourceName).ExistsInAzure(r),
		//	),
		//},
		//data.ImportStep(),
		{
			Config: r.authSSHKeyUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFleetLinux_authEd25519SSHKeys(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authEd25519SSHKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		//{
		//	Config: r.authSSHKey(data),
		//	Check: acceptance.ComposeTestCheckFunc(
		//		check.That(data.ResourceName).ExistsInAzure(r),
		//	),
		//},
		//data.ImportStep(),
		{
			Config: r.authEd25519SSHKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFleetLinux_authDisablePasswordAuthUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		//{
		//	// disable it
		//	Config: r.authSSHKey(data),
		//	Check: acceptance.ComposeTestCheckFunc(
		//		check.That(data.ResourceName).ExistsInAzure(r),
		//	),
		//},
		//data.ImportStep("admin_password"),
		{
			// enable it
			Config: r.authenticationTest(data, r.authPasswordVirtualMachineProfile()),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
		//{
		//	// disable it
		//	Config: r.authSSHKey(data),
		//	Check: acceptance.ComposeTestCheckFunc(
		//		check.That(data.ResourceName).ExistsInAzure(r),
		//	),
		//},
		//data.ImportStep("admin_password"),
	})
}

func (r AzureFleetResource) authenticationTest(data acceptance.TestData, vmProfile string) string {
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

  compute_profile {
    %[4]s
  }
}
`, r.templateLinux(data), data.RandomInteger, data.Locations.Primary, vmProfile)
}

func (r AzureFleetResource) authPasswordVirtualMachineProfile() string {
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
      }

      os_profile {
        computer_name_prefix = "prefix"
        admin_username       = "azureuser"
        admin_password       = "TestPassword$0"
        linux_configuration {
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

      network_api_version = "2020-11-01"
    }
`
}

func (r AzureFleetResource) authSSHKeyVirtualMachineProfile() string {
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
      }

      os_profile {
        computer_name_prefix = "prefix"
        admin_username       = "azureuser"
        linux_configuration {
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

func (r AzureFleetResource) authMultipleSSHKeysVirtualMachineProfile() string {
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
      }

      os_profile {
        computer_name_prefix = "prefix"
        admin_username       = "azureuser"
        linux_configuration {
          ssh_keys {
			username   = "azureuser"
			public_key = local.first_public_key
		  }
          ssh_keys {
			username   = "azureuser"
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

      network_api_version = "2020-11-01"
    }
`
}

func (r AzureFleetResource) authSSHKeyAndPasswordVirtualMachineProfile() string {
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
      }

      os_profile {
        computer_name_prefix = "prefix"
        admin_username       = "azureuser"
        linux_configuration {
          ssh_keys {
			username   = "azureuser"
            admin_password      = "P@ssw0rd1234!"
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

func (r AzureFleetResource) authEd25519SSHKey(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_azure_fleet" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.ed25519_public_key
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, r.templateLinux(data), data.RandomInteger)
}

func (r AzureFleetResource) authSSHKeyUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_azure_fleet" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.second_public_key
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, r.templateLinux(data), data.RandomInteger)
}

func (r AzureFleetResource) tmp() string {
	return `
`
}
func (r AzureFleetResource) authMultipleSSHKeys(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_azure_fleet" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.first_public_key
  }

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.second_public_key
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, r.templateLinux(data), data.RandomInteger)
}
func (r AzureFleetResource) fleetTemplate(data acceptance.TestData, computeProfile string) string {
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

	%[4]s
}
`, r.templateLinux(data), data.RandomInteger, data.Locations.Primary, computeProfile)
}
