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

func TestAccAzureFleet_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_fleet", "test")
	r := AzureFleetResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.fleet(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("compute_profile.0.virtual_machine_profile.0.os_profile.0.admin_password"),
		{
			Config: r.fleetUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("compute_profile.0.virtual_machine_profile.0.os_profile.0.admin_password"),
		{
			Config: r.fleet(data),
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

func (r AzureFleetResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

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

`, data.RandomInteger, data.Locations.Primary)
}

func (r AzureFleetResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_azure_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  spot_priority_profile {
    maintain_enabled = false
  }

  vm_sizes_profile {
    name = "Standard_DS1_v2"
  }

  compute_profile {
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
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r AzureFleetResource) spotCapacity(data acceptance.TestData) string {
	template := r.template(data)
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
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r AzureFleetResource) spotCapacityUpdate(data acceptance.TestData) string {
	template := r.template(data)
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
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r AzureFleetResource) tempTest(data acceptance.TestData) string {
	template := r.template(data)
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
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r AzureFleetResource) tempTestUpdate(data acceptance.TestData) string {
	template := r.template(data)
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
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r AzureFleetResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_azure_fleet" "import" {
  name                = azurerm_azure_fleet.test.name
  resource_group_name = azurerm_azure_fleet.test.resource_group_name
  location            = azurerm_azure_fleet.test.location

  spot_priority_profile {
    maintain_enabled = false
  }

  vm_sizes_profile {
    name = "Standard_DS1_v2"
  }

  compute_profile {
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
}
`, config)
}

func (r AzureFleetResource) fleet(data acceptance.TestData) string {
	template := r.template(data)
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
    compute_api_version = "2023-09-01"
  }

  additional_location_profile {
    location = "%[4]s"
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
          load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
          primary                                = true
          subnet_id                              = azurerm_subnet.test.id
        }
        primary = true
      }
      network_api_version = "2020-11-01"
    }
  }

   identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }


 # plan {
 #   name      = "os"
   # product   = "rancheros"
  #  publisher = "rancher"
  #  promotion_code = "test"
  #  version        = "1.0"
  #}

  regular_priority_profile {
    allocation_strategy = "LowestPrice"
    min_capacity        = 1
    capacity            = 1
  }

  spot_priority_profile {
    allocation_strategy     = "LowestPrice"
    eviction_policy         = "Delete"
    maintain_enabled        = false
    max_hourly_price_per_vm = 1
    min_capacity            = 0
    capacity                = 1
  }

  tags = {
    Hello = "There"
    World = "Example"
  }

  vm_attributes {
    memory_in_gib {
      max = 2.0
      min = 1.0
    }

    vcpu_count {
      max = 2
      min = 1
    }

    accelerator_count {
      max = 2
      min = 1
    }

    accelerator_manufacturers = ["AMD"]
    accelerator_support       = "Included"
    accelerator_types         = ["GPU"]
    architecture_types        = ["X64"]
    burstable_support         = "Included"
    cpu_manufacturers         = ["Microsoft"]
    data_disk_count {
      max = 2
      min = 1
    }
    excluded_vm_sizes_profile = ["Standard_E2s_v3"]
    local_storage_disk_types  = ["HDD"]
    local_storage_in_gib {
      max = 2
      min = 1
    }
    local_storage_support = "Included"
    memory_in_gib_per_vcpu {
      max = 2
      min = 1.0
    }
    network_bandwidth_in_mbps {
      max = 2.0
      min = 1.0
    }
    network_interface_count {
      max = 2
      min = 1
    }
    rdma_network_interface_count {
      max = 2
      min = 1
    }
    rdma_support  = "Included"
    vm_categories = ["ComputeOptimized"]
  }

  zones = ["1"]
}
`, template, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (r AzureFleetResource) fleetUpdate(data acceptance.TestData) string {
	template := r.template(data)
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

  vm_sizes_profile {
    name = "Standard_D2s_v3"
    rank = "10001"
  }

  additional_location_profile {
    location = "%[4]s"
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
          load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
          primary                                = true
          subnet_id                              = azurerm_subnet.test.id
        }
        primary = true
      }
      network_api_version = "2024-01-01"
    }
    compute_api_version = "2024-03-01"
  }

  #identity {
   # type         = "UserAssigned"
   # identity_ids = [azurerm_user_assigned_identity.test.id]
  #}

 # plan {
  #  name           = "nrTest2"
   # product        = "NewRelic2"
   # publisher      = "NewRelic2"
   # promotion_code = "test2"
   # version        = "2.0"
  #}

  regular_priority_profile {
    allocation_strategy = "LowestPrice"
    min_capacity        = 1
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
`, template, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}
