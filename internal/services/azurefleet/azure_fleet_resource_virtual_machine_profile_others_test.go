// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azurefleet_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccAzureFleetVirtualMachineProfileOthers_additionalCapabilities(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.additionalCapabilities(data, data.Locations.Primary),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleetVirtualMachineProfileOthers_automaticVMGuestPatchingLinux(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.linuxVMGuestPatching(data, data.Locations.Primary, "ImageDefault"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.linuxVMGuestPatching(data, data.Locations.Primary, "AutomaticByPlatform"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.linuxVMGuestPatching(data, data.Locations.Primary, "ImageDefault"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleetVirtualMachineProfileOthers_automaticVMGuestPatchingWindows(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsVMGuestPatching(data, data.Locations.Primary, "AutomaticByOS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.windows_configuration.0.admin_password"),
		{
			Config: r.windowsVMGuestPatching(data, data.Locations.Primary, "AutomaticByPlatform"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.windows_configuration.0.admin_password"),
		{
			Config: r.windowsVMGuestPatching(data, data.Locations.Primary, "AutomaticByOS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.windows_configuration.0.admin_password"),
	})
}

func TestAccAzureFleetVirtualMachineProfileOthers_hotPatching(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.hotPatchingWindows(data, data.Locations.Primary, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.hotPatchingWindows(data, data.Locations.Primary, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.hotPatchingWindows(data, data.Locations.Primary, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleetVirtualMachineProfileOthers_capacityReservationGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.capacityReservationGroup(data, data.Locations.Primary),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.capacityReservationGroupUpdate(data, data.Locations.Primary),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.capacityReservationGroup(data, data.Locations.Primary),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func (r AzureFleetResource) linuxVMGuestPatching(data acceptance.TestData, location string, patchMode string) string {
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
    name = "Standard_F2"
  }

  virtual_machine_profile {
    image_reference {
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
        admin_username                  = "azureuser"
        admin_password                  = "P@ssw0rd1234!"
        password_authentication_enabled = true

        patch_mode = "%[5]s"
        patch_assessment_mode = "%[5]s"
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

 extension {
    name                               = "HealthExtension"
    publisher                          = "Microsoft.ManagedServices"
    type                               = "ApplicationHealthLinux"
    type_handler_version               = "1.0"

    settings_json = jsonencode({
      "protocol"    = "http"
      "port"        = 80
      "requestPath" = "/healthEndpoint"
    })
  }
  }
}
`, r.linuxPublicKeyTemplate(), r.template(data, location), data.RandomInteger, location, patchMode)
}

func (r AzureFleetResource) hotPatchingWindows(data acceptance.TestData, location string, enabled bool) string {
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
    name = "Standard_F2s_v2"
  }

  virtual_machine_profile {
    image_reference {
      publisher = "MicrosoftWindowsServer"
       offer     = "WindowsServer"
    sku       = "2016-Datacenter"
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
      windows_configuration {
        computer_name_prefix            = "prefix"
        admin_username                  = "azureuser"
        admin_password                  = "P@ssw0rd1234!"

        patch_mode = "AutomaticByPlatform"
				hot_patching_enabled = %[5]t
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

 extension {
    name                               = "HealthExtension"
    publisher                          = "Microsoft.ManagedServices"
    type                               = "ApplicationHealthLinux"
    type_handler_version               = "1.0"

    settings_json = jsonencode({
      "protocol"    = "http"
      "port"        = 80
      "requestPath" = "/healthEndpoint"
    })
  }
  }
}
`, r.linuxPublicKeyTemplate(), r.template(data, location), data.RandomInteger, location, enabled)
}

func (r AzureFleetResource) windowsVMGuestPatching(data acceptance.TestData, location string, patchMode string) string {
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
    image_reference {
      publisher = "MicrosoftWindowsServer"
       offer     = "WindowsServer"
    sku       = "2016-Datacenter"
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
      windows_configuration {
        computer_name_prefix            = "prefix"
        admin_username                  = "azureuser"
        admin_password                  = "P@ssw0rd1234!"

        patch_mode = "%[4]s"
        patch_assessment_mode = "ImageDefault"
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

 extension {
    name                               = "HealthExtension"
    publisher                          = "Microsoft.ManagedServices"
    type                               = "ApplicationHealthLinux"
    type_handler_version               = "1.0"

    settings_json = jsonencode({
      "protocol"    = "http"
      "port"        = 80
      "requestPath" = "/healthEndpoint"
    })
  }
  }
}
`, r.template(data, location), data.RandomInteger, location, patchMode)
}

func (r AzureFleetResource) additionalCapabilities(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "azurerm_azure_fleet" "test" {
  name                = "acctest-fleet-%[3]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[4]s"
additional_capabilities_ultra_ssd_enabled = true
additional_capabilities_hibernation_enabled = true
  zones = ["1"]

  regular_priority_profile {
    capacity     = 1
    min_capacity = 1
  }

  vm_sizes_profile {
    name = "Standard_D2s_v3"
  }

  virtual_machine_profile {
    image_reference {
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
      admin_username       = "myadmin"
      admin_password       = "Passwword1234"
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

extension {
    name                               = "HealthExtension"
    publisher                          = "Microsoft.ManagedServices"
    type                               = "ApplicationHealthLinux"
    type_handler_version               = "1.0"

    settings_json = jsonencode({
      "protocol"    = "http"
      "port"        = 80
      "requestPath" = "/healthEndpoint"
    })
  }
  }
}
`, r.linuxPublicKeyTemplate(), r.template(data, location), data.RandomInteger, location)
}

func (r AzureFleetResource) capacityReservationGroupUpdate(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "azurerm_capacity_reservation_group" "test" {
  name                = "acctest-ccrg-%[3]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_capacity_reservation" "test" {
  name                          = "acctest-ccr-%[3]d"
  capacity_reservation_group_id = azurerm_capacity_reservation_group.test.id
  sku {
    name     = "Standard_F2"
    capacity = 1
  }
}

resource "azurerm_capacity_reservation" "test2" {
  name                          = "acctest-ccr2-%[3]d"
  capacity_reservation_group_id = azurerm_capacity_reservation_group.test.id
  sku {
    name     = "Standard_F2"
    capacity = 1
  }
}

resource "azurerm_azure_fleet" "test" {
  name                = "acctest-fleet-%[3]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[4]s"

  capacity_reservation_group_id = azurerm_capacity_reservation_group.test2.id

  regular_priority_profile {
    capacity     = 1
    min_capacity = 1
  }

  vm_sizes_profile {
    name = "Standard_D2s_v3"
  }

  virtual_machine_profile {
    image_reference {
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
      admin_username       = "myadmin"
      admin_password       = "Passwword1234"
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
depends_on = [azurerm_capacity_reservation.test]
}
`, r.linuxPublicKeyTemplate(), r.template(data, location), data.RandomInteger, location)
}
func (r AzureFleetResource) capacityReservationGroup(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "azurerm_capacity_reservation_group" "test" {
  name                = "acctest-ccrg-%[3]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_capacity_reservation" "test" {
  name                          = "acctest-ccr-%[3]d"
  capacity_reservation_group_id = azurerm_capacity_reservation_group.test.id
  sku {
    name     = "Standard_F2"
    capacity = 1
  }
}

resource "azurerm_azure_fleet" "test" {
  name                = "acctest-fleet-%[3]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[4]s"

  capacity_reservation_group_id = azurerm_capacity_reservation_group.test.id

  regular_priority_profile {
    capacity     = 1
    min_capacity = 1
  }

  vm_sizes_profile {
    name = "Standard_D2s_v3"
  }

  virtual_machine_profile {
    image_reference {
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
      admin_username       = "myadmin"
      admin_password       = "Passwword1234"
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
depends_on = [azurerm_capacity_reservation.test]
}
`, r.linuxPublicKeyTemplate(), r.template(data, location), data.RandomInteger, location)
}
