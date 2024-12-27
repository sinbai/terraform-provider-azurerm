// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azurefleet_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/azurefleet/2024-11-01/fleets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AzureFleetResource struct{}

func TestAccAzureFleet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("compute_profile.0.virtual_machine_profile.0.os_profile.0.admin_password"),
	})
}

func TestAccAzureFleet_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccAzureFleet_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("compute_profile.0.virtual_machine_profile.0.os_profile.0.admin_password"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("compute_profile.0.virtual_machine_profile.0.os_profile.0.admin_password"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("compute_profile.0.virtual_machine_profile.0.os_profile.0.admin_password"),
	})
}

func TestAccAzureFleet_spotCapacity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.spotCapacity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("compute_profile.0.virtual_machine_profile.0.os_profile.0.admin_password"),
		{
			Config: r.spotCapacityUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("compute_profile.0.virtual_machine_profile.0.os_profile.0.admin_password"),
		{
			Config: r.spotCapacity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("compute_profile.0.virtual_machine_profile.0.os_profile.0.admin_password"),
	})
}

func TestAccAzureFleet_spotVmSizeProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.spotVmSizeProfile(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("compute_profile.0.virtual_machine_profile.0.os_profile.0.admin_password"),
		{
			Config: r.spotVmSizeProfileUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("compute_profile.0.virtual_machine_profile.0.os_profile.0.admin_password"),
		{
			Config: r.spotVmSizeProfile(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("compute_profile.0.virtual_machine_profile.0.os_profile.0.admin_password"),
	})
}

func TestAccAzureFleet_additionalLocation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.additionalLocation(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("compute_profile.0.virtual_machine_profile.0.os_profile.0.admin_password"),
	})
}

func TestAccAzureFleet_tempTest(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.tempTest(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("compute_profile.0.virtual_machine_profile.0.os_profile.0.admin_password", "compute_profile.0.compute_api_version"),
		{
			Config: r.tempTestUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("compute_profile.0.virtual_machine_profile.0.os_profile.0.admin_password", "compute_profile.0.compute_api_version"),
	})
}

func (r AzureFleetResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := fleets.ParseFleetID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.AzureFleet.FleetsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r AzureFleetResource) basic(data acceptance.TestData) string {
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
`, r.template(data, data.Locations.Primary), data.RandomInteger, data.Locations.Primary, r.basicVirtualMachineProfile())
}

func (r AzureFleetResource) spotCapacity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_azure_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  spot_priority_profile {
    maintain_enabled = true
    capacity         = 1
  }

  vm_sizes_profile {
    name = "Standard_DS1_v2"
  }
  vm_sizes_profile {
    name = "Standard_D2s_v3"
  }
  vm_sizes_profile {
    name = "Standard_D2as_v4"
  }

  compute_profile {
    %[4]s
  }
  zones = ["1", "2", "3"]
}
`, r.template(data, data.Locations.Primary), data.RandomInteger, data.Locations.Primary, r.basicVirtualMachineProfile())
}

func (r AzureFleetResource) spotCapacityUpdate(data acceptance.TestData) string {
	template := r.template(data, data.Locations.Primary)
	vmProfile := r.basicVirtualMachineProfile()
	return fmt.Sprintf(`
%[1]s

resource "azurerm_azure_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  spot_priority_profile {
    maintain_enabled = true
    capacity         = 2
  }

  vm_sizes_profile {
    name = "Standard_DS1_v2"
  }
  vm_sizes_profile {
    name = "Standard_D2s_v3"
  }
  vm_sizes_profile {
    name = "Standard_D2as_v4"
  }

  compute_profile {
    %[4]s
  }
  zones = ["1", "2", "3"]
}
`, template, data.RandomInteger, data.Locations.Primary, vmProfile)
}

func (r AzureFleetResource) spotVmSizeProfile(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_azure_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  spot_priority_profile {
    maintain_enabled = true
    capacity         = 1
  }

  vm_sizes_profile {
    name = "Standard_DS1_v2"
  }
  vm_sizes_profile {
    name = "Standard_D2s_v3"
  }
  vm_sizes_profile {
    name = "Standard_D2as_v4"
  }

  compute_profile {

    %[4]s

  }

  zones = ["1", "2", "3"]
}
`, r.template(data, data.Locations.Primary), data.RandomInteger, data.Locations.Primary, r.basicVirtualMachineProfile())
}

func (r AzureFleetResource) spotVmSizeProfileUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_azure_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  spot_priority_profile {
    maintain_enabled = true
    capacity         = 1
  }

  vm_sizes_profile {
    name = "Standard_DS1_v2"
  }
  vm_sizes_profile {
    name = "Standard_D2s_v3"
  }
  vm_sizes_profile {
    name = "Standard_DS2_v2"
  }

  compute_profile {

    %[4]s

  }
  zones = ["1", "2", "3"]
}
`, r.template(data, data.Locations.Primary), data.RandomInteger, data.Locations.Primary, r.basicVirtualMachineProfile())
}

func (r AzureFleetResource) tempTest(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_azure_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  spot_priority_profile {
    capacity         = 1
    maintain_enabled = false
  }

  regular_priority_profile {
    capacity     = 1
    min_capacity = 1
  }

  vm_sizes_profile {
    name = "Standard_DS1_v2"
  }

  compute_profile {
    compute_api_version = "2023-09-01"
    virtual_machine_profile {
      storage_profile {
        image_reference {
          offer     = "0001-com-ubuntu-server-focal"
          publisher = "canonical"
          sku       = "20_04-lts-gen2"
          version   = "latest"
        }

        os_disk {
          caching       = "ReadWrite"
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
  }

  zones = ["1", "2", "3"]
}
`, r.template(data, data.Locations.Primary), data.RandomInteger, data.Locations.Primary)
}

func (r AzureFleetResource) tempTestUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_azure_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  spot_priority_profile {
    capacity         = 1
    maintain_enabled = false
  }

  regular_priority_profile {
    capacity     = 1
    min_capacity = 1
  }

  vm_sizes_profile {
    name = "Standard_DS1_v2"
  }

  compute_profile {
    compute_api_version = "2024-03-01"
    virtual_machine_profile {
      storage_profile {
        image_reference {
          offer     = "0001-com-ubuntu-server-focal"
          publisher = "canonical"
          sku       = "20_04-lts-gen2"
          version   = "latest"
        }

        os_disk {
          caching       = "ReadWrite"
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
  }

  zones = ["1", "2", "3"]
}
`, r.template(data, data.Locations.Primary), data.RandomInteger, data.Locations.Primary)
}

func (r AzureFleetResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_azure_fleet" "import" {
  name                = azurerm_azure_fleet.test.name
  resource_group_name = azurerm_azure_fleet.test.resource_group_name
  location            = azurerm_azure_fleet.test.location

  regular_priority_profile {
    capacity     = 1
    min_capacity = 1
  }

  vm_sizes_profile {
    name = "Standard_DS1_v2"
  }

  compute_profile {
    %s
  }
}
`, r.basic(data), r.basicVirtualMachineProfile())
}

func (r AzureFleetResource) additionalLocation(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_azure_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  compute_profile {
    %[4]s
  }

  additional_location_profile {
    location = "%[5]s"

  }

  spot_priority_profile {
    allocation_strategy = "LowestPrice"
    maintain_enabled    = false
    min_capacity        = 0
    capacity            = 1
  }

  vm_sizes_profile {
    name = "Standard_DS1_v2"
  }

}
`, r.additionalLocationTemplate(data), data.RandomInteger, data.Locations.Primary, r.basicVirtualMachineProfile(), data.Locations.Secondary) //, r.virtualMachineProfileOverride()) #%[6]s
}

func (r AzureFleetResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

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
    allocation_strategy = "LowestPrice"
    min_capacity        = 0
    capacity            = 2
  }

  spot_priority_profile {
    allocation_strategy     = "LowestPrice"
    eviction_policy         = "Delete"
    maintain_enabled        = false
    max_hourly_price_per_vm = 1
    min_capacity            = 0
    capacity                = 2
  }

  vm_sizes_profile {
    name = "Standard_D2s_v3"
    rank = "10001"
  }

  compute_profile {
    %[4]s
  }
  additional_location_profile {
    location = "%[5]s"
    %[6]s
    compute_api_version = "2024-03-01"
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  # plan {
  #  name           = "nrTest2"
  # product        = "NewRelic2"
  # publisher      = "NewRelic2"
  # promotion_code = "test2"
  # version        = "2.0"
  #}

  tags = {
    Hello = "ThereUpdate"
    World = "ExampleUpdate"
  }

  vm_attributes {
    memory_in_gib {
      max = 3.0
      min = 2.0
    }

    vcpu_count {
      max = 3
      min = 2
    }

    accelerator_count {
      max = 3
      min = 2
    }

    accelerator_manufacturers = ["Nvidia"]
    accelerator_support       = "Required"
    accelerator_types         = ["FPGA"]
    architecture_types        = ["ARM64"]
    burstable_support         = "Required"
    cpu_manufacturers         = ["Intel"]
    data_disk_count {
      max = 3
      min = 2
    }
    excluded_vm_sizes_profile = ["Standard_D8s_v3"]
    local_storage_disk_types  = ["SSD"]
    local_storage_in_gib {
      max = 3
      min = 2
    }
    local_storage_support = "Required"
    memory_in_gib_per_vcpu {
      max = 3
      min = 3
    }
    network_bandwidth_in_mbps {
      max = 3.0
      min = 2.0
    }
    network_interface_count {
      max = 3
      min = 2
    }
    rdma_network_interface_count {
      max = 3
      min = 2
    }
    rdma_support  = "Required"
    vm_categories = ["StorageOptimized"]
  }

  zones = ["1"]
}
`, r.template(data, data.Locations.Primary), data.RandomInteger, data.Locations.Primary, r.basicVirtualMachineProfile(), data.Locations.Secondary, r.virtualMachineProfileOverride())
}

func (r AzureFleetResource) basicVirtualMachineProfile() string {
	return fmt.Sprintf(`
virtual_machine_profile {
      storage_profile {
        image_reference {
          offer     = "0001-com-ubuntu-server-focal"
          publisher = "canonical"
          sku       = "20_04-lts-gen2"
          version   = "latest"
        }

        os_disk {
          caching       = "ReadWrite"
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
`)
}

func (r AzureFleetResource) virtualMachineProfileOverride() string {
	return fmt.Sprintf(`
virtual_machine_profile_override {
      storage_profile {
        image_reference {
          offer     = "0001-com-ubuntu-server-focal"
          publisher = "canonical"
          sku       = "20_04-lts-gen2"
          version   = "latest"
        }

        os_disk {
          caching       = "ReadWrite"
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
          load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test2.id]
          primary                                = true
          subnet_id                              = azurerm_subnet.test2.id
        }
        primary = true
      }
      network_api_version = "2020-11-01"
}
`)
}

func (AzureFleetResource) linuxPublicKeyTemplate() string {
	return `
# note: whilst these aren't used in all tests, it saves us redefining these everywhere
locals {
  first_public_key   = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+wWK73dCr+jgQOAxNsHAnNNNMEMWOHYEccp6wJm2gotpr9katuF/ZAdou5AaW1C61slRkHRkpRRX9FA9CYBiitZgvCCz+3nWNN7l/Up54Zps/pHWGZLHNJZRYyAB6j5yVLMVHIHriY49d/GZTZVNB8GoJv9Gakwc/fuEZYYl4YDFiGMBP///TzlI4jhiJzjKnEvqPFki5p2ZRJqcbCiF4pJrxUQR/RXqVFQdbRLZgYfJ8xGB878RENq3yQ39d8dVOkq4edbkzwcUmwwwkYVPIoDGsYLaRHnG+To7FvMeyO7xDVQkMKzopTQV8AuKpyvpqu0a9pWOMaiCyDytO7GGN you@me.com"
  second_public_key  = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0/NDMj2wG6bSa6jbn6E3LYlUsYiWMp1CQ2sGAijPALW6OrSu30lz7nKpoh8Qdw7/A4nAJgweI5Oiiw5/BOaGENM70Go+VM8LQMSxJ4S7/8MIJEZQp5HcJZ7XDTcEwruknrd8mllEfGyFzPvJOx6QAQocFhXBW6+AlhM3gn/dvV5vdrO8ihjET2GoDUqXPYC57ZuY+/Fz6W3KV8V97BvNUhpY5yQrP5VpnyvvXNFQtzDfClTvZFPuoHQi3/KYPi6O0FSD74vo8JOBZZY09boInPejkm9fvHQqfh0bnN7B6XJoUwC1Qprrx+XIy7ust5AEn5XL7d4lOvcR14MxDDKEp you@me.com"
  ed25519_public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIDqzSi9IHoYnbE3YQ+B2fQEVT8iGFemyPovpEtPziIVB you@me.com"
}
`
}

func (r AzureFleetResource) template(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

`, r.templateWithOutProvider(data, location), data.RandomInteger, location)
}

func (r AzureFleetResource) templateWithOutProvider(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%[1]d"
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
    name                 = "internal-%[1]d"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  name            = "internal"
  loadbalancer_id = azurerm_lb.test.id
}

`, data.RandomInteger, location)
}

func (r AzureFleetResource) additionalLocationTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_resource_group" "test2" {
  name     = "acctest-fleet-al-%[2]d"
  location = "%[3]s"
}

resource "azurerm_virtual_network" "test2" {
  name                = "acctvn-al-%[2]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test2.location
  resource_group_name = azurerm_resource_group.test2.name
}

resource "azurerm_subnet" "test2" {
  name                 = "acctsub-al-%[2]d"
  resource_group_name  = azurerm_resource_group.test2.name
  virtual_network_name = azurerm_virtual_network.test2.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_public_ip" "test2" {
  name                = "acctestpublicIP-al-%[2]d"
  location            = azurerm_resource_group.test2.location
  resource_group_name = azurerm_resource_group.test2.name
  allocation_method   = "Static"
  sku                 = "Standard"
  zones               = ["1"]
}

resource "azurerm_lb" "test2" {
  name                = "acctest-loadbalancer-al-%[2]d"
  location            = azurerm_resource_group.test2.location
  resource_group_name = azurerm_resource_group.test2.name
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "internal-al-%[2]d"
    public_ip_address_id = azurerm_public_ip.test2.id
  }
}

resource "azurerm_lb_backend_address_pool" "test2" {
  name            = "internal-al"
  loadbalancer_id = azurerm_lb.test2.id
}

`, r.template(data, data.Locations.Primary), data.RandomInteger, data.Locations.Secondary)
}

func (r AzureFleetResource) linuxTemplate(data acceptance.TestData, location string, vmProfile string) string {
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
    name = "Standard_DS1_v2"
  }

  compute_profile {
		%[5]s
  }

  zones = [1, 2, 3]
}


`, r.linuxPublicKeyTemplate(), r.template(data, location), data.RandomInteger, location, vmProfile)
}

func (r AzureFleetResource) linuxTemplateWithoutProvider(data acceptance.TestData, location string, vmProfile string, dependencies string) string {
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
    name = "Standard_DS1_v2"
  }

  compute_profile {
		%[5]s
  }
}

%[6]s


`, r.linuxPublicKeyTemplate(), r.templateWithOutProvider(data, location), data.RandomInteger, location, vmProfile, dependencies)
}
