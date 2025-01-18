// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azurefleet_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccAzureFleet_virtualMachineProfileNetwork_multiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multiple(data, data.Locations.Primary),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleet_virtualMachineProfileNetwork_networkSecurityGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.networkSecurityGroup(data, data.Locations.Primary),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleet_virtualMachineProfileNetwork_acceleratedNetworking(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.acceleratedNetworking(data, data.Locations.Primary, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.acceleratedNetworking(data, data.Locations.Primary, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),

		{
			Config: r.acceleratedNetworking(data, data.Locations.Primary, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleet_virtualMachineProfileNetwork_dnsNameLabel(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dnsNameLabel(data, data.Locations.Primary, "test-domain-label", "ResourceGroupReuse"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.dnsNameLabel(data, data.Locations.Primary, "updated-domain-label", "SubscriptionReuse"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleet_virtualMachineProfileNetwork_ipForwarding(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.ipForwarding(data, data.Locations.Primary, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.ipForwarding(data, data.Locations.Primary, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.ipForwarding(data, data.Locations.Primary, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleet_virtualMachineProfileNetwork_publicIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicPublicIP(data, data.Locations.Primary),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleet_virtualMachineProfileNetwork_publicIPSku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.publicIPSku(data, data.Locations.Primary),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleet_virtualMachineProfileNetwork_publicIPVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.publicIPVersion(data, data.Locations.Primary),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleet_virtualMachineProfileNetwork_basicDNSSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dNSSettings(data, data.Locations.Primary, "\"8.8.8.8\""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.dNSSettings(data, data.Locations.Primary, "\"8.8.8.8\", \"8.8.4.4\""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.dNSSettings(data, data.Locations.Primary, "\"8.8.8.8\""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleet_virtualMachineProfileNetwork_loadBalancer(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.loadBalancer(data, data.Locations.Primary),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

//func TestAccAzureFleet_virtualMachineProfileNetwork_fpga(t *testing.T) {
//	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
//	r := AzureFleetTestResource{}
//
//	data.ResourceTest(t, r, []acceptance.TestStep{
//		{
//			Config: r.fpga(data, data.Locations.Primary, true),
//			Check: acceptance.ComposeTestCheckFunc(
//				check.That(data.ResourceName).ExistsInAzure(r),
//			),
//		},
//		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
//		{
//			Config: r.fpga(data, data.Locations.Primary, false),
//			Check: acceptance.ComposeTestCheckFunc(
//				check.That(data.ResourceName).ExistsInAzure(r),
//			),
//		},
//		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
//	})
//}

//func (r AzureFleetTestResource) fpga(data acceptance.TestData, location string, enabled bool) string {
//	return fmt.Sprintf(`
//%[1]s
//
//resource "azurerm_subnet" "test1" {
//	name                 = "acctestSubnet1"
//	resource_group_name  = azurerm_resource_group.test.name
//	virtual_network_name = azurerm_virtual_network.test.name
//	address_prefixes     = ["10.0.1.0/24"]
//
//	delegation {
//		name = "acctestdelegation"
//		service_delegation {
//			name    = "Microsoft.Network/fpgaNetworkInterfaces"
//			actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
//		}
//	}
//}
//
//resource "azurerm_azure_fleet" "test" {
//  name                = "acctest-fleet-%[2]d"
//  resource_group_name = azurerm_resource_group.test.name
//  location            = "%[3]s"
//  platform_fault_domain_count = 2
//
//  regular_priority_profile {
//    capacity     = 2
//    min_capacity = 0
//  }
//
//  vm_sizes_profile {
//    name = "Standard_D2s_v3"
//  }
//
//  virtual_machine_profile {
//    os_profile {
//      linux_configuration {
//        computer_name_prefix            = "prefix"
//        admin_username                  = local.admin_username
//        admin_password                  = local.admin_password
//        password_authentication_enabled = true
//      }
//    }
//
//    network_interface {
//      name    = "primary-networkProTest"
//      primary = true
//      fpga_enabled = %[4]t
//
//      ip_configuration {
//        name      = "primary"
//        primary   = true
//        subnet_id = azurerm_subnet.test1.id
//
//        public_ip_address {
//          name                    = "TestPublicIPConfiguration"
//          domain_name_label       = "test-domain-label"
//          idle_timeout_in_minutes = 4
//        }
//      }
//    }
//
//    os_disk {
//      storage_account_type = "Standard_LRS"
//      caching              = "ReadWrite"
//    }
//
//    source_image_reference {
//      publisher = "Canonical"
//      offer     = "0001-com-ubuntu-server-jammy"
//      sku       = "22_04-lts"
//      version   = "latest"
//    }
//  }
//}
//`, r.template(data, location), data.RandomInteger, location, enabled)
//}

func (r AzureFleetTestResource) multiple(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_azure_fleet" "test" {
  name                        = "acctest-fleet-%[2]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = "%[3]s"
  platform_fault_domain_count = 2

  regular_priority_profile {
    capacity     = 2
    min_capacity = 0
  }

  vm_sizes_profile {
    name = "Standard_D1_v2"
  }

  virtual_machine_profile {
    os_profile {
      linux_configuration {
        computer_name_prefix            = "prefix"
        admin_username                  = local.admin_username
        admin_password                  = local.admin_password
        password_authentication_enabled = true
      }
    }

    network_interface {
      name    = "primary-networkProTest"
      primary = true

      ip_configuration {
        name      = "primary"
        primary   = true
        subnet_id = azurerm_subnet.test.id

        public_ip_address {
          name                    = "TestPublicIPConfiguration"
          domain_name_label       = "test-domain-label"
          idle_timeout_in_minutes = 4
        }
      }
    }

    network_interface {
      name    = "secondary-networkProTest"
      primary = true

      ip_configuration {
        name      = "secondary"
        primary   = true
        subnet_id = azurerm_subnet.test.id
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
`, r.template(data, location), data.RandomInteger, location)
}

func (r AzureFleetTestResource) acceleratedNetworking(data acceptance.TestData, location string, enabled bool) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_azure_fleet" "test" {
  name                        = "acctest-fleet-%[2]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = "%[3]s"
  platform_fault_domain_count = 2

  regular_priority_profile {
    capacity     = 2
    min_capacity = 0
  }

  vm_sizes_profile {
    name = "Standard_D2s_v3"
  }

  virtual_machine_profile {
    os_profile {
      linux_configuration {
        computer_name_prefix            = "prefix"
        admin_username                  = local.admin_username
        admin_password                  = local.admin_password
        password_authentication_enabled = true
      }
    }

    network_interface {
      name                          = "primary-networkProTest"
      primary                       = true
      enable_accelerated_networking = %[4]t

      ip_configuration {
        name      = "primary"
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
`, r.template(data, location), data.RandomInteger, location, enabled)
}

func (r AzureFleetTestResource) networkSecurityGroup(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
%[1]s


resource "azurerm_network_security_group" "test" {
  name                = "acceptanceTestSecurityGroup-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}


resource "azurerm_azure_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  platform_fault_domain_count = 2

  regular_priority_profile {
    capacity     = 2
    min_capacity = 0
  }

  vm_sizes_profile {
    name = "Standard_D1_v2"
  }

  virtual_machine_profile {
    os_profile {
      linux_configuration {
        computer_name_prefix            = "prefix"
        admin_username                  = local.admin_username
        admin_password                  = local.admin_password
        password_authentication_enabled = true
      }
    }

    network_interface {
      name                      = "networkProTest"
      primary                   = true
      network_security_group_id = azurerm_network_security_group.test.id

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
`, r.template(data, location), data.RandomInteger, location)
}

func (r AzureFleetTestResource) dnsNameLabel(data acceptance.TestData, location string, domainNamelabel string, scope string) string {
	return fmt.Sprintf(`
%[1]s


resource "azurerm_azure_fleet" "test" {
  name                        = "acctest-fleet-%[2]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = "%[3]s"
  platform_fault_domain_count = 2

  regular_priority_profile {
    capacity     = 2
    min_capacity = 0
  }

  vm_sizes_profile {
    name = "Standard_D1_v2"
  }

  virtual_machine_profile {
    os_profile {
      linux_configuration {
        computer_name_prefix            = "prefix"
        admin_username                  = local.admin_username
        admin_password                  = local.admin_password
        password_authentication_enabled = true
      }
    }

    network_interface {
      name    = "primary-networkProTest"
      primary = true

      ip_configuration {
        name      = "primary"
        primary   = true
        subnet_id = azurerm_subnet.test.id

        public_ip_address {
          name                    = "TestPublicIPConfiguration"
          domain_name_label       = "%[4]s"
          domain_name_label_scope = "%[5]s"
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
`, r.template(data, location), data.RandomInteger, location, domainNamelabel, scope)
}

func (r AzureFleetTestResource) ipForwarding(data acceptance.TestData, location string, enabled bool) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_azure_fleet" "test" {
  name                        = "acctest-fleet-%[2]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = "%[3]s"
  platform_fault_domain_count = 2

  regular_priority_profile {
    capacity     = 1
    min_capacity = 0
  }

  vm_sizes_profile {
    name = "Standard_D1_v2"
  }

  virtual_machine_profile {
    os_profile {
      linux_configuration {
        computer_name_prefix            = "prefix"
        admin_username                  = local.admin_username
        admin_password                  = local.admin_password
        password_authentication_enabled = true
      }
    }

    network_interface {
      name                  = "primary-networkProTest"
      primary               = true
      ip_forwarding_enabled = "%[4]t"

      ip_configuration {
        name      = "primary"
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
`, r.template(data, location), data.RandomInteger, location, enabled)
}

func (r AzureFleetTestResource) basicPublicIP(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_azure_fleet" "test" {
  name                        = "acctest-fleet-%[2]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = "%[3]s"
  platform_fault_domain_count = 2

  regular_priority_profile {
    capacity     = 1
    min_capacity = 0
  }

  vm_sizes_profile {
    name = "Standard_D1_v2"
  }

  virtual_machine_profile {
    os_profile {
      linux_configuration {
        computer_name_prefix            = "prefix"
        admin_username                  = local.admin_username
        admin_password                  = local.admin_password
        password_authentication_enabled = true
      }
    }

    network_interface {
      name    = "primary-networkProTest"
      primary = true

      ip_configuration {
        name      = "primary"
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
`, r.template(data, location), data.RandomInteger, location)
}

func (r AzureFleetTestResource) publicIPSku(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
%[1]s


resource "azurerm_azure_fleet" "test" {
  name                        = "acctest-fleet-%[2]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = "%[3]s"
  platform_fault_domain_count = 2

  regular_priority_profile {
    capacity     = 1
    min_capacity = 0
  }

  vm_sizes_profile {
    name = "Standard_D1_v2"
  }

  virtual_machine_profile {
    os_profile {
      linux_configuration {
        computer_name_prefix            = "prefix"
        admin_username                  = local.admin_username
        admin_password                  = local.admin_password
        password_authentication_enabled = true
      }
    }

    network_interface {
      name    = "primary-networkProTest"
      primary = true

      ip_configuration {
        name      = "primary"
        primary   = true
        subnet_id = azurerm_subnet.test.id

        public_ip_address {
          name                    = "TestPublicIPConfigurationUpdate"
          domain_name_label       = "update-domain-label"
          idle_timeout_in_minutes = 3
          delete_option           = "Detach"
          sku {
            name = "Standard"
            tier = "Regional"
          }
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
`, r.template(data, location), data.RandomInteger, location)
}

func (r AzureFleetTestResource) publicIPVersion(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
%[1]s


resource "azurerm_azure_fleet" "test" {
  name                        = "acctest-fleet-%[2]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = "%[3]s"
  platform_fault_domain_count = 2

  regular_priority_profile {
    capacity     = 1
    min_capacity = 0
  }

  vm_sizes_profile {
    name = "Standard_D1_v2"
  }

  virtual_machine_profile {
    os_profile {
      linux_configuration {
        computer_name_prefix            = "prefix"
        admin_username                  = local.admin_username
        admin_password                  = local.admin_password
        password_authentication_enabled = true
      }
    }

    network_interface {
      name    = "primary-networkProTest"
      primary = true

      ip_configuration {
        name      = "primary"
        primary   = true
        subnet_id = azurerm_subnet.test.id

        public_ip_address {
          name                    = "TestPublicIPConfiguration"
          domain_name_label       = "test-domain-label"
          idle_timeout_in_minutes = 4
        }
      }
      ip_configuration {
        name    = "second"
        version = "IPv6"

        public_ip_address {
          name                    = "second"
          idle_timeout_in_minutes = 4
          version                 = "IPv6"
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
`, r.template(data, location), data.RandomInteger, location)
}

func (r AzureFleetTestResource) dNSSettings(data acceptance.TestData, location string, dnsServers string) string {
	return fmt.Sprintf(`
%[1]s


resource "azurerm_azure_fleet" "test" {
  name                        = "acctest-fleet-%[2]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = "%[3]s"
  platform_fault_domain_count = 2

  regular_priority_profile {
    capacity     = 1
    min_capacity = 0
  }

  vm_sizes_profile {
    name = "Standard_D4_v2"
  }

  virtual_machine_profile {
    os_profile {
      linux_configuration {
        computer_name_prefix            = "prefix"
        admin_username                  = local.admin_username
        admin_password                  = local.admin_password
        password_authentication_enabled = true
      }
    }

    network_interface {
      name        = "primary-networkProTest"
      primary     = true
      dns_servers = [%[4]s]
      ip_configuration {
        name      = "primary"
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
`, r.template(data, location), data.RandomInteger, location, dnsServers)
}

func (r AzureFleetTestResource) loadBalancer(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
%[1]s


resource "azurerm_azure_fleet" "test" {
  name                        = "acctest-fleet-%[2]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = "%[3]s"
  platform_fault_domain_count = 2

  regular_priority_profile {
    capacity     = 1
    min_capacity = 0
  }

  vm_sizes_profile {
    name = "Standard_D1_v2"
  }

  virtual_machine_profile {
    os_profile {
      linux_configuration {
        computer_name_prefix            = "prefix"
        admin_username                  = local.admin_username
        admin_password                  = local.admin_password
        password_authentication_enabled = true
      }
    }

    network_interface {
      name    = "primary-networkProTest"
      primary = true

      ip_configuration {
        name      = "primary"
        primary   = true
        subnet_id = azurerm_subnet.test.id

        load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
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
`, r.template(data, location), data.RandomInteger, location)
}
