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
	})
}

func TestAccAzureFleetVirtualMachineProfileOthers_licenseType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.licenseTypeWindows(data, data.Locations.Primary, "Windows_Client"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.windows_configuration.0.admin_password"),
		{
			Config: r.licenseTypeWindows(data, data.Locations.Primary, "Windows_Server"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.windows_configuration.0.admin_password"),
	})
}

func TestAccAzureFleetVirtualMachineProfileOthers_UserData(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.userData(data, data.Locations.Primary, "Hello World"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.userData(data, data.Locations.Primary, "Goodbye World"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleetVirtualMachineProfileOthers_galleryApplication(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.galleryApplication(data, data.Locations.Primary, "test"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.galleryApplication(data, data.Locations.Primary, "testUpdate"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleetVirtualMachineProfileOthers_serviceArtifactReference(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.serviceArtifactReference(data, data.Locations.Primary),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.serviceArtifactReferenceUpdate(data, data.Locations.Primary),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleetVirtualMachineProfileOthers_scheduledEvent(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.scheduledEvent(data, data.Locations.Primary),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.scheduledEventUpdate(data, data.Locations.Primary),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleetVirtualMachineProfileOthers_encryptionAtHost(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.encryptionAtHost(data, data.Locations.Primary),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
		//{
		//	Config: r.encryptionAtHostUpdate(data, data.Locations.Primary),
		//	Check: acceptance.ComposeTestCheckFunc(
		//		check.That(data.ResourceName).ExistsInAzure(r),
		//	),
		//},
		//data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func (r AzureFleetResource) linuxVMGuestPatching(data acceptance.TestData, location string, patchMode string) string {
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
        admin_username                  = "azureuser"
        admin_password                  = "P@ssw0rd1234!"
        password_authentication_enabled = true

        patch_mode = "%[4]s"
        patch_assessment_mode = "%[4]s"
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

func (r AzureFleetResource) hotPatchingWindows(data acceptance.TestData, location string, enabled bool) string {
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

  virtual_machine_profile {
    source_image_reference {
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
				hot_patching_enabled = %[4]t
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
`, r.template(data, location), data.RandomInteger, location, enabled)
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
    source_image_reference {
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

resource "azurerm_azure_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"
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
`, r.template(data, location), data.RandomInteger, location)
}

func (r AzureFleetResource) capacityReservationGroupUpdate(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
%[1]s


resource "azurerm_capacity_reservation_group" "test" {
  name                = "acctest-ccrg-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_capacity_reservation" "test" {
  name                          = "acctest-ccr-%[2]d"
  capacity_reservation_group_id = azurerm_capacity_reservation_group.test.id
  sku {
    name     = "Standard_F2"
    capacity = 1
  }
}

resource "azurerm_capacity_reservation" "test2" {
  name                          = "acctest-ccr2-%[2]d"
  capacity_reservation_group_id = azurerm_capacity_reservation_group.test.id
  sku {
    name     = "Standard_F2"
    capacity = 1
  }
}

resource "azurerm_azure_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  capacity_reservation_group_id = azurerm_capacity_reservation_group.test2.id

  regular_priority_profile {
    capacity     = 1
    min_capacity = 1
  }

  vm_sizes_profile {
    name = "Standard_D2s_v3"
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
`, r.template(data, location), data.RandomInteger, location)
}

func (r AzureFleetResource) capacityReservationGroup(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_capacity_reservation_group" "test" {
  name                = "acctest-ccrg-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_capacity_reservation" "test" {
  name                          = "acctest-ccr-%[2]d"
  capacity_reservation_group_id = azurerm_capacity_reservation_group.test.id
  sku {
    name     = "Standard_F2"
    capacity = 1
  }
}

resource "azurerm_azure_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  capacity_reservation_group_id = azurerm_capacity_reservation_group.test.id

  regular_priority_profile {
    capacity     = 1
    min_capacity = 1
  }

  vm_sizes_profile {
    name = "Standard_D2s_v3"
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
`, r.template(data, location), data.RandomInteger, location)
}

func (r AzureFleetResource) galleryApplication(data acceptance.TestData, location string, tag string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "accteststr%[4]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "test"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "blob"
}

resource "azurerm_storage_blob" "test" {
  name                   = "script"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Page"
  size                   = 512
}

resource "azurerm_storage_blob" "test2" {
  name                   = "script2"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Page"
  size                   = 512
}

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_gallery_application" "test" {
  name              = "acctest-app-%[2]d"
  gallery_id        = azurerm_shared_image_gallery.test.id
  location          = azurerm_shared_image_gallery.test.location
  supported_os_type = "Linux"
}

resource "azurerm_gallery_application_version" "test" {
  name                   = "0.0.1"
  gallery_application_id = azurerm_gallery_application.test.id
  location               = azurerm_gallery_application.test.location

  source {
    media_link                 = azurerm_storage_blob.test.id
    default_configuration_link = azurerm_storage_blob.test.id
  }

  manage_action {
    install = "[install command]"
    remove  = "[remove command]"
  }

  target_region {
    name                   = azurerm_gallery_application.test.location
    regional_replica_count = 1
    storage_account_type   = "Premium_LRS"
  }
}

resource "azurerm_azure_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  regular_priority_profile {
    capacity     = 1
    min_capacity = 1
  }

  vm_sizes_profile {
    name = "Standard_D2s_v3"
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

    gallery_application {
			version_id             = azurerm_gallery_application_version.test.id
			configuration_blob_uri = azurerm_storage_blob.test2.id
			order                  = 1
			tag                    = "%[5]s"
      automatic_upgrade_enabled = false
			treat_failure_as_deployment_failure_enabled = false
		}
  }
}
`, r.template(data, location), data.RandomInteger, location, data.RandomString, tag)
}

func (r AzureFleetResource) serviceArtifactReference(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "accteststr%[4]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "test"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "blob"
}

resource "azurerm_storage_blob" "test" {
  name                   = "script"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Page"
  size                   = 512
}

resource "azurerm_storage_blob" "test2" {
  name                   = "script2"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Page"
  size                   = 512
}

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_gallery_application" "test" {
  name              = "acctest-app-%[2]d"
  gallery_id        = azurerm_shared_image_gallery.test.id
  location          = azurerm_shared_image_gallery.test.location
  supported_os_type = "Linux"
}

resource "azurerm_gallery_application_version" "test" {
  name                   = "0.0.1"
  gallery_application_id = azurerm_gallery_application.test.id
  location               = azurerm_gallery_application.test.location

  source {
    media_link                 = azurerm_storage_blob.test.id
    default_configuration_link = azurerm_storage_blob.test.id
  }

  manage_action {
    install = "[install command]"
    remove  = "[remove command]"
  }

  target_region {
    name                   = azurerm_gallery_application.test.location
    regional_replica_count = 1
    storage_account_type   = "Premium_LRS"
  }
}

resource "azurerm_azure_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  regular_priority_profile {
    capacity     = 1
    min_capacity = 1
  }

  vm_sizes_profile {
    name = "Standard_D2s_v3"
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

    service_artifact_id = azurerm_gallery_application_version.test.id
  }
}
`, r.template(data, location), data.RandomInteger, location, data.RandomString)
}

func (r AzureFleetResource) serviceArtifactReferenceUpdate(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "accteststr%[4]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "test"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "blob"
}

resource "azurerm_storage_blob" "test" {
  name                   = "script"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Page"
  size                   = 512
}

resource "azurerm_storage_blob" "test2" {
  name                   = "script2"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Page"
  size                   = 512
}

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_gallery_application" "test" {
  name              = "acctest-app-%[2]d"
  gallery_id        = azurerm_shared_image_gallery.test.id
  location          = azurerm_shared_image_gallery.test.location
  supported_os_type = "Linux"
}

resource "azurerm_gallery_application_version" "test" {
  name                   = "0.0.1"
  gallery_application_id = azurerm_gallery_application.test.id
  location               = azurerm_gallery_application.test.location

  source {
    media_link                 = azurerm_storage_blob.test.id
    default_configuration_link = azurerm_storage_blob.test.id
  }

  manage_action {
    install = "[install command]"
    remove  = "[remove command]"
  }

  target_region {
    name                   = azurerm_gallery_application.test.location
    regional_replica_count = 1
    storage_account_type   = "Premium_LRS"
  }
}

resource "azurerm_gallery_application_version" "test2" {
  name                   = "0.0.2"
  gallery_application_id = azurerm_gallery_application.test.id
  location               = azurerm_gallery_application.test.location

  source {
    media_link                 = azurerm_storage_blob.test.id
    default_configuration_link = azurerm_storage_blob.test.id
  }

  manage_action {
    install = "[install command]"
    remove  = "[remove command]"
  }

  target_region {
    name                   = azurerm_gallery_application.test.location
    regional_replica_count = 1
    storage_account_type   = "Premium_LRS"
  }
}

resource "azurerm_azure_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  regular_priority_profile {
    capacity     = 1
    min_capacity = 1
  }

  vm_sizes_profile {
    name = "Standard_D2s_v3"
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

    service_artifact_id = azurerm_gallery_application_version.test2.id
  }
}
`, r.template(data, location), data.RandomInteger, location, data.RandomString)
}

func (r AzureFleetResource) licenseTypeWindows(data acceptance.TestData, location string, lType string) string {
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
    name = "Standard_D2s_v3"
  }

  virtual_machine_profile {
    source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter-Server-Core"
    version   = "latest"
    }

    os_disk {
      caching              = "ReadWrite"
      storage_account_type = "Standard_LRS"
    }

    os_profile {
     windows_configuration {
			computer_name_prefix = "testvm"
      admin_username       = "myadmin"
      admin_password       = "Passwword1234"

      automatic_updates_enabled  = true
      provision_vm_agent_enabled = true
      time_zone                   = "W. Europe Standard Time"

 			winrm_listener {
        protocol = "Http"
      }
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
    license_type = "%[4]s"
  }
}
`, r.template(data, location), data.RandomInteger, location, lType)
}

func (r AzureFleetResource) userData(data acceptance.TestData, location string, userDta string) string {
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

    os_profile {
    	linux_configuration {
				computer_name_prefix = "testvm-%[2]d"
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
     user_data_base64 = base64encode("%[4]s")
  }
}
`, r.template(data, location), data.RandomInteger, location, userDta)
}
func (r AzureFleetResource) scheduledEvent(data acceptance.TestData, location string) string {
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

    os_profile {
    	linux_configuration {
				computer_name_prefix = "testvm-%[2]d"
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

     scheduled_event_termination_timeout = "PT5M"
 		 scheduled_event_os_image_timeout =  "PT15M"
  }
}
`, r.template(data, location), data.RandomInteger, location)
}

func (r AzureFleetResource) scheduledEventUpdate(data acceptance.TestData, location string) string {
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

    os_profile {
    	linux_configuration {
				computer_name_prefix = "testvm-%[2]d"
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

     ScheduledEventTerminationTimeout = "PT15M"
  }
}
`, r.template(data, location), data.RandomInteger, location)
}
func (r AzureFleetResource) encryptionAtHost(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
%[1]s

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
  name                        = "acctestkv%[4]s"
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
  name                = "acctestdes-%[2]d"
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

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_azure_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  regular_priority_profile {
    capacity     = 1
    min_capacity = 1
  }

  vm_sizes_profile {
    name = "Standard_D2s_v3"
  }

   virtual_machine_profile {
    
     security_profile {
      encryption_at_host_enabled = true
			proxy_agent{
					mode = "Audit"
			}
			security_type = "ConfidentialVM"
			uefi_secure_boot_enabled = true
			uefi_vtpm_enabled = true
      user_assigned_identity_id = azurerm_user_assigned_identity.test.id
    }

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
      disk_encryption_set_id = azurerm_disk_encryption_set.test.id
      storage_account_type   = "Standard_LRS"
      security_encryption_type         = "DiskWithVMGuestState"
      security_disk_encryption_set_id = azurerm_disk_encryption_set.test.id
    }

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
`, r.templateWithOutProvider(data, location), data.RandomInteger, location, data.RandomString)
}
