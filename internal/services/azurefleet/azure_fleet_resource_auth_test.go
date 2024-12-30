// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azurefleet_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccFleetLinux_authPassword(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.linuxTemplate(data, data.Locations.Primary, r.authPasswordVirtualMachineProfile()),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.admin_password"),
	})
}

func TestAccFleetLinux_authSSHKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.linuxTemplate(data, data.Locations.Primary, r.authSSHKeyVirtualMachineProfile()),
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
			Config: r.linuxTemplate(data, data.Locations.Primary, r.authSSHKeyAndPasswordVirtualMachineProfile()),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.admin_password"),
	})
}

func TestAccFleetLinux_authMultipleSSHKeys(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.linuxTemplate(data, data.Locations.Primary, r.authMultipleSSHKeysVirtualMachineProfile()),
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
			Config: r.linuxTemplate(data, data.Locations.Primary, r.authEd25519SSHKeyVirtualMachineProfile()),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
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
        admin_password       = "P@ssw0rd1234!"
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

func (r AzureFleetResource) authEd25519SSHKeyVirtualMachineProfile() string {
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

      network_api_version = "2020-11-01"
    }
`
}
