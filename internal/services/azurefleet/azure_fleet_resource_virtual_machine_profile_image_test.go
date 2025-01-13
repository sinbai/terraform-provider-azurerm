// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azurefleet_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccAzureFleetVirtualMachineProfileImage_imageFromSourceImageReference(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{

			Config: r.imageFromSourceImageReference(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleetVirtualMachineProfileImage_imageFromImage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.imageFromImage(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleetVirtualMachineProfileImage_imageFromCommunitySharedImageGallery(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.imageFromCommunitySharedImageGallery(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccAzureFleetVirtualMachineProfileImage_imageFromSharedImageGallery(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.imageFromSharedImageGallery(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func (r AzureFleetResource) imageFromSourceImageReference(data acceptance.TestData) string {
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
`, r.template(data, data.Locations.Primary), data.RandomInteger, data.Locations.Primary, r.baseLinuxVirtualMachineProfile())
}

func (r AzureFleetResource) imageFromImage(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_public_ip" "source" {
  name                = "acctpip-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  domain_name_label   = "acctestsourcevm-%[2]d"
  sku                 = "Basic"
}

resource "azurerm_network_interface" "test" {
  name                = "acctestnic-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

   ip_configuration {
    name                          = "testconfigurationsource"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.source.id
  }
}

resource "azurerm_linux_virtual_machine" "source" {
  name                            = "acctestsourceVM-%[2]d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  size                            = "Standard_DS1_v2"
  admin_username                  = "adminuser"
  disable_password_authentication = false
  admin_password                  = "P@ssw0rd1234!"

  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.first_public_key
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}

resource "azurerm_image" "test" {
  name                      = "capture"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  source_virtual_machine_id = azurerm_linux_virtual_machine.source.id
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
    name = "Standard_DS1_v2"
  }

  virtual_machine_profile {

	source_image_id                 = azurerm_image.test.id
	
	os_disk {
		caching              = "ReadWrite"
		storage_account_type = "Standard_LRS"
	}
	
	os_profile {
		linux_configuration {
			computer_name_prefix = "prefix"
			admin_username       = "azureuser"
			admin_password       = "TestPassword$0"
			password_authentication_enabled = true
		}
	}

	network_interface {
		name                            = "networkProTest"
   	primary 												= true
		accelerated_networking_enabled  = false
		ip_forwarding_enabled           = true
		ip_configuration {
			name                                   = "ipConfigTest"
			load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
			primary                                = true
			subnet_id                              = azurerm_subnet.test.id
		}
	}
}
}
`, r.linuxTemplate(data, data.Locations.Primary), data.RandomInteger, data.Locations.Primary)
}

func (r AzureFleetResource) imageFromSharedImageGalleryVersion(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_public_ip" "source" {
  name                = "acctpip-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  domain_name_label   = "acctestsourcevm-%[2]d"
  sku                 = "Basic"
}


resource "azurerm_network_interface" "test" {
  name                = "acctestnic-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfigurationsource"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.source.id
  }
}

resource "azurerm_linux_virtual_machine" "source" {
  name                            = "acctestsourceVM-%[2]d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  size                            = "Standard_DS1_v2"
  admin_username                  = "adminuser"
  disable_password_authentication = false
  admin_password                  = "P@ssw0rd1234!"

  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.first_public_key
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}

resource "azurerm_image" "test" {
  name                      = "capture"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  source_virtual_machine_id = azurerm_linux_virtual_machine.source.id
}


resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%[2]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_shared_image" "test" {
  name                = "acctest-gallery-image"
  gallery_name        = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Linux"

  identifier {
    publisher = "AcceptanceTest-Publisher"
    offer     = "AcceptanceTest-Offer"
    sku       = "AcceptanceTest-Sku"
  }
}

resource "azurerm_shared_image_version" "test" {
  name                = "0.0.1"
  gallery_name        = azurerm_shared_image.test.gallery_name
  image_name          = azurerm_shared_image.test.name
  resource_group_name = azurerm_shared_image.test.resource_group_name
  location            = azurerm_shared_image.test.location
  managed_image_id    = azurerm_image.test.id

  target_region {
    name                   = azurerm_shared_image.test.location
    regional_replica_count = "5"
    storage_account_type   = "Standard_LRS"
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
    name = "Standard_DS1_v2"
  }

  virtual_machine_profile {

	source_image_id                 = azurerm_shared_image_version.test.id
	
	os_disk {
		caching              = "ReadWrite"
		storage_account_type = "Standard_LRS"
	}
	
	os_profile {
		linux_configuration {
			computer_name_prefix = "prefix"
			admin_username       = "azureuser"
			admin_password       = "TestPassword$0"
			password_authentication_enabled = true
		}
	}

	network_interface {
		name                            = "networkProTest"
   	primary 												= true
		accelerated_networking_enabled  = false
		ip_forwarding_enabled           = true
		ip_configuration {
			name                                   = "ipConfigTest"
			load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
			primary                                = true
			subnet_id                              = azurerm_subnet.test.id
		}
	}
}
}
`, r.linuxTemplate(data, data.Locations.Primary), data.RandomInteger, data.Locations.Primary)
}

func (r AzureFleetResource) imageFromCommunitySharedImageGallery(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_public_ip" "source" {
  name                = "acctpip-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  domain_name_label   = "acctestsourcevm-%[2]d"
  sku                 = "Basic"
}


resource "azurerm_network_interface" "test" {
  name                = "acctestnic-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfigurationsource"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.source.id
  }
}

resource "azurerm_linux_virtual_machine" "source" {
  name                            = "acctestsourceVM-%[2]d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  size                            = "Standard_DS1_v2"
  admin_username                  = "adminuser"
  disable_password_authentication = false
  admin_password                  = "P@ssw0rd1234!"

  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.first_public_key
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}

resource "azurerm_image" "test" {
  name                      = "capture"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  source_virtual_machine_id = azurerm_linux_virtual_machine.source.id
}


resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%[2]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sharing {
    permission = "Community"
    community_gallery {
      eula            = "https://eula.net"
      prefix          = "prefix"
      publisher_email = "publisher@test.net"
      publisher_uri   = "https://publisher.net"
    }
  }
}

resource "azurerm_shared_image" "test" {
  name                = "acctest-gallery-image"
  gallery_name        = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Linux"

  identifier {
    publisher = "AcceptanceTest-Publisher"
    offer     = "AcceptanceTest-Offer"
    sku       = "AcceptanceTest-Sku"
  }
}

resource "azurerm_shared_image_version" "test" {
  name                = "0.0.1"
  gallery_name        = azurerm_shared_image.test.gallery_name
  image_name          = azurerm_shared_image.test.name
  resource_group_name = azurerm_shared_image.test.resource_group_name
  location            = azurerm_shared_image.test.location
  managed_image_id    = azurerm_image.test.id

  target_region {
    name                   = azurerm_shared_image.test.location
    regional_replica_count = "5"
    storage_account_type   = "Standard_LRS"
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
    name = "Standard_DS1_v2"
  }

  virtual_machine_profile {

	source_image_id                 = "/communityGalleries/${azurerm_shared_image_gallery.test.sharing.0.community_gallery.0.name}/images/${azurerm_shared_image_version.test.image_name}/versions/${azurerm_shared_image_version.test.name}"
	
	os_disk {
		caching              = "ReadWrite"
		storage_account_type = "Standard_LRS"
	}
	
	os_profile {
		linux_configuration {
			computer_name_prefix = "prefix"
			admin_username       = "azureuser"
			admin_password       = "TestPassword$0"
			password_authentication_enabled = true
		}
	}

	network_interface {
		name                            = "networkProTest"
   	primary 												= true
		accelerated_networking_enabled  = false
		ip_forwarding_enabled           = true
		ip_configuration {
			name                                   = "ipConfigTest"
			load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
			primary                                = true
			subnet_id                              = azurerm_subnet.test.id
		}
	}
}
}
`, r.linuxTemplate(data, data.Locations.Primary), data.RandomInteger, data.Locations.Primary)
}

func (r AzureFleetResource) imageFromSharedImageGallery(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_public_ip" "source" {
  name                = "acctpip-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  domain_name_label   = "acctestsourcevm-%[2]d"
  sku                 = "Basic"
}


resource "azurerm_network_interface" "test" {
  name                = "acctestnic-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfigurationsource"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.source.id
  }
}

resource "azurerm_linux_virtual_machine" "source" {
  name                            = "acctestsourceVM-%[2]d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  size                            = "Standard_DS1_v2"
  admin_username                  = "adminuser"
  disable_password_authentication = false
  admin_password                  = "P@ssw0rd1234!"

  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.first_public_key
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}

resource "azurerm_image" "test" {
  name                      = "capture"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  source_virtual_machine_id = azurerm_linux_virtual_machine.source.id
}


resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%[2]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_shared_image" "test" {
  name                = "acctest-gallery-image"
  gallery_name        = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Linux"

  identifier {
    publisher = "AcceptanceTest-Publisher"
    offer     = "AcceptanceTest-Offer"
    sku       = "AcceptanceTest-Sku"
  }
}

resource "azurerm_shared_image_version" "test" {
  name                = "0.0.1"
  gallery_name        = azurerm_shared_image.test.gallery_name
  image_name          = azurerm_shared_image.test.name
  resource_group_name = azurerm_shared_image.test.resource_group_name
  location            = azurerm_shared_image.test.location
  managed_image_id    = azurerm_image.test.id

  target_region {
    name                   = azurerm_shared_image.test.location
    regional_replica_count = "5"
    storage_account_type   = "Standard_LRS"
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
    name = "Standard_DS1_v2"
  }

  virtual_machine_profile {

	source_image_id                 = azurerm_shared_image_version.test.id
	
	os_disk {
		caching              = "ReadWrite"
		storage_account_type = "Standard_LRS"
	}
	
	os_profile {
		linux_configuration {
			computer_name_prefix = "prefix"
			admin_username       = "azureuser"
			admin_password       = "TestPassword$0"
			password_authentication_enabled = true
		}
	}

	network_interface {
		name                            = "networkProTest"
   	primary 												= true
		accelerated_networking_enabled  = false
		ip_forwarding_enabled           = true
		ip_configuration {
			name                                   = "ipConfigTest"
			load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
			primary                                = true
			subnet_id                              = azurerm_subnet.test.id
		}
	}
}
}
`, r.linuxTemplate(data, data.Locations.Primary), data.RandomInteger, data.Locations.Primary)
}
