// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package computefleet_test

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

type ComputeFleetTestResource struct{}

func TestAccComputeFleet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}
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

func TestAccComputeFleet_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}
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

func TestAccComputeFleet_completeExceptVMSS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeExceptVMSS(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_additionalLocationLinux(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.additionalLocationLinux(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password",
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
		),
	})
}

func TestAccComputeFleet_additionalLocationWindows(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.additionalLocationWindows(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.windows_configuration.0.admin_password",
			"virtual_machine_profile.0.os_profile.0.windows_configuration.0.admin_password",
		),
	})
}

func TestAccComputeFleet_vmAttributes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.vmAttributesBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password", "compute_api_version"),
		{
			Config: r.vmAttributesAppend(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password", "compute_api_version"),
		{
			Config: r.vmAttributesUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password", "compute_api_version"),
	})
}

func TestAccComputeFleet_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identityNone(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.identity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.identityUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.identityNone(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_plan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.plan(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_zones(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.zones(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_spotPriorityProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.spotPriorityProfile(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.spotPriorityProfileUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.spotPriorityProfile(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_spotAndRegulaPriorityProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.spotAndRegulaPriorityProfile(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.spotAndRegulaPriorityProfile(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.spotAndRegulaPriorityProfile(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func (r ComputeFleetTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := fleets.ParseFleetID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ComputeFleet.ComputeFleetClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r ComputeFleetTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_compute_fleet" "test" {
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

func (r ComputeFleetTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_compute_fleet" "import" {
  name                = azurerm_compute_fleet.test.name
  resource_group_name = azurerm_compute_fleet.test.resource_group_name
  location            = azurerm_compute_fleet.test.location

  regular_priority_profile {
    capacity     = 1
    min_capacity = 1
  }

  vm_sizes_profile {
    name = "Standard_DS1_v2"
  }

  %s
}
`, r.basic(data), r.baseLinuxVirtualMachineProfile())
}

func (r ComputeFleetTestResource) identityNone(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_compute_fleet" "test" {
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

func (r ComputeFleetTestResource) identity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_user_assigned_identity" "test2" {
  name                = "acctest2%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_compute_fleet" "test" {
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

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  %[4]s
}
`, r.template(data, data.Locations.Primary), data.RandomInteger, data.Locations.Primary, r.baseLinuxVirtualMachineProfile()) // need to update location after the bug mentioned in email is fixed after two weeks.
}

func (r ComputeFleetTestResource) identityUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_user_assigned_identity" "test2" {
  name                = "acctest2%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_compute_fleet" "test" {
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

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test2.id]
  }

  %[4]s
}
`, r.template(data, data.Locations.Primary), data.RandomInteger, data.Locations.Primary, r.baseLinuxVirtualMachineProfile())
}

func (r ComputeFleetTestResource) additionalLocationLinux(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

resource "azurerm_compute_fleet" "test" {
  name                = "acctest-fleet-%[3]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[4]s"

  %[5]s

  additional_location_profile {
    location = "%[6]s"
    %[7]s
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
`, r.template(data, data.Locations.Primary), r.additionalLinuxTemplate(data, data.Locations.Secondary), data.RandomInteger, data.Locations.Primary, r.baseLinuxVirtualMachineProfile(), data.Locations.Secondary, r.additionalLinuxVirtualMachineProfile())
}

func (r ComputeFleetTestResource) additionalLocationWindows(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

resource "azurerm_compute_fleet" "test" {
  name                = "acctest-fleet-%[3]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[4]s"

  %[5]s

  additional_location_profile {
    location = "%[6]s"
    %[7]s
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
`, r.template(data, data.Locations.Primary), r.additionalWindowsTemplate(data, data.Locations.Secondary), data.RandomInteger, data.Locations.Primary, r.baseWindowsVirtualMachineProfile(), data.Locations.Secondary, r.additionalWindowsVirtualMachineProfile())
}

func (r ComputeFleetTestResource) vmAttributesBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_compute_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  spot_priority_profile {
    allocation_strategy = "LowestPrice"
    maintain_enabled    = false
    capacity            = 2
  }

  %[4]s

  vm_attributes {
    memory_in_gib {
      max = 10
      min = 1
    }
    vcpu_count {
      max = 10
      min = 1
    }
  }

  compute_api_version = "2024-03-01"
}
`, r.template(data, data.Locations.Primary), data.RandomInteger, data.Locations.Primary, r.baseLinuxVirtualMachineProfile())
}

func (r ComputeFleetTestResource) vmAttributesAppend(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_compute_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  spot_priority_profile {
    allocation_strategy = "LowestPrice"
    maintain_enabled    = false
    capacity            = 2
  }

  %[4]s

  vm_attributes {
    memory_in_gib {
      max = 10
      min = 1
    }
    vcpu_count {
      max = 10
      min = 1
    }
    data_disk_count {
      max = 10
      min = 1
    }
    local_storage_in_gib {
      max = 100
      min = 1
    }
    memory_in_gib_per_vcpu {
      max = 10
      min = 0
    }
    local_storage_support    = "Included"
    local_storage_disk_types = ["SSD"]
    architecture_types       = ["X64", "ARM64"]
    cpu_manufacturers        = ["Intel"]
    network_bandwidth_in_mbps {
      max = 500
      min = 0
    }
    network_interface_count {
      max = 10
      min = 0
    }
    excluded_vm_sizes = ["Standard_DS1_v2"]
    vm_categories     = ["GeneralPurpose"]
    burstable_support = "Excluded"
    rdma_support      = "Included"
    rdma_network_interface_count {
      max = 10
      min = 0
    }
    accelerator_support = "Included"
    accelerator_count {
      max = 3
      min = 0
    }
  }

  compute_api_version = "2024-03-01"
}
`, r.template(data, data.Locations.Primary), data.RandomInteger, data.Locations.Primary, r.baseLinuxVirtualMachineProfile())
}

func (r ComputeFleetTestResource) vmAttributesUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_compute_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  spot_priority_profile {
    allocation_strategy = "LowestPrice"
    maintain_enabled    = false
    capacity            = 2
  }

  %[4]s

  vm_attributes {
    memory_in_gib {
      max = 9
      min = 1
    }
    vcpu_count {
      max = 9
      min = 1
    }
    data_disk_count {
      max = 9
      min = 1
    }
    local_storage_in_gib {
      max = 99
      min = 1
    }
    memory_in_gib_per_vcpu {
      max = 9
      min = 0
    }
    local_storage_support    = "Included"
    local_storage_disk_types = ["HDD", "SSD"]
    architecture_types       = ["X64"]
    cpu_manufacturers        = ["Intel", "Microsoft"]
    network_bandwidth_in_mbps {
      max = 501
      min = 0
    }
    network_interface_count {
      max = 9
      min = 0
    }
    excluded_vm_sizes = ["Standard_D2s_v3"]
    vm_categories     = ["GeneralPurpose", "ComputeOptimized"]
    burstable_support = "Included"
    rdma_support      = "Included"
    rdma_network_interface_count {
      max = 9
      min = 0
    }
    accelerator_support = "Included"
    accelerator_count {
      max = 2
      min = 0
    }
  }

  compute_api_version = "2024-03-01"
}
`, r.template(data, data.Locations.Primary), data.RandomInteger, data.Locations.Primary, r.baseLinuxVirtualMachineProfile())
}

func (r ComputeFleetTestResource) baseLinuxVirtualMachineProfile() string {
	return fmt.Sprintf(`
virtual_machine_profile {
	network_api_version = "2020-11-01"
	source_image_reference {
		offer     = "0001-com-ubuntu-server-focal"
		publisher = "canonical"
		sku       = "20_04-lts-gen2"
		version   = "latest"
	}
	
	os_profile {
		linux_configuration {
			computer_name_prefix = "prefix"
			admin_username       = local.admin_username
			admin_password       = local.admin_password
			password_authentication_enabled = true
		}
	}

	network_interface {
		name                            = "networkProTest"
   	primary 												= true
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
# ignore_changes as when os_disk block is not specified the API return default values for caching, delete_option, disk_size_in_gb and storage_account_type
lifecycle {
	ignore_changes = [virtual_machine_profile.0.os_disk]
}
`)
}

func (r ComputeFleetTestResource) baseWindowsVirtualMachineProfile() string {
	return fmt.Sprintf(`
virtual_machine_profile {
	network_api_version = "2020-11-01"
	source_image_reference {
		publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter-Server-Core"
    version   = "latest"
	}

	os_profile {
		windows_configuration {
			computer_name_prefix = "testvm"
      admin_username       = local.admin_username
      admin_password       = local.admin_password
		}
	}

	network_interface {
		name                            = "networkProTest"
   	primary 												= true
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
}
# ignore_changes as when os_disk block is not specified the API return default values for caching, delete_option, disk_size_in_gb and storage_account_type
lifecycle {
	iignore_changes = [virtual_machine_profile.0.os_disk]
}
`)
}

func (r ComputeFleetTestResource) additionalLinuxVirtualMachineProfile() string {
	return fmt.Sprintf(`
virtual_machine_profile_override {
	source_image_reference {
		offer     = "0001-com-ubuntu-server-focal"
		publisher = "canonical"
		sku       = "20_04-lts-gen2"
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
			password_authentication_enabled = true
		}
	}

	network_interface {
      name    = "networkProTest"
      primary = true
      ip_configuration {
        name      = "TestIPConfiguration"
        subnet_id = azurerm_subnet.linux_test.id
        primary   = true
        public_ip_address {
          name                    = "TestPublicIPConfiguration"
          domain_name_label       = "test-domain-label"
          idle_timeout_in_minutes = 4
        }
      }
    }
}
`)
}

func (r ComputeFleetTestResource) additionalWindowsVirtualMachineProfile() string {
	return fmt.Sprintf(`
virtual_machine_profile_override {
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
      admin_username       = local.admin_username
      admin_password       = local.admin_password

      automatic_updates_enabled  = true
      provision_vm_agent_enabled = true

 			winrm_listener {
        protocol = "Http"
      }
		}
	}

	network_interface {
		name                            = "networkProTest"
   	primary 												= true

		ip_configuration {
        name                                   = "ipConfigTest"
        load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.windows_test.id]
        primary                                = true
        subnet_id                              = azurerm_subnet.windows_test.id
      }
	}
}
`)
}

func (r ComputeFleetTestResource) template(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

`, r.templateWithOutProvider(data, location), data.RandomInteger, location)
}

func (r ComputeFleetTestResource) templateWithOutProvider(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
locals {
  first_public_key          = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+wWK73dCr+jgQOAxNsHAnNNNMEMWOHYEccp6wJm2gotpr9katuF/ZAdou5AaW1C61slRkHRkpRRX9FA9CYBiitZgvCCz+3nWNN7l/Up54Zps/pHWGZLHNJZRYyAB6j5yVLMVHIHriY49d/GZTZVNB8GoJv9Gakwc/fuEZYYl4YDFiGMBP///TzlI4jhiJzjKnEvqPFki5p2ZRJqcbCiF4pJrxUQR/RXqVFQdbRLZgYfJ8xGB878RENq3yQ39d8dVOkq4edbkzwcUmwwwkYVPIoDGsYLaRHnG+To7FvMeyO7xDVQkMKzopTQV8AuKpyvpqu0a9pWOMaiCyDytO7GGN you@me.com"
  second_public_key         = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0/NDMj2wG6bSa6jbn6E3LYlUsYiWMp1CQ2sGAijPALW6OrSu30lz7nKpoh8Qdw7/A4nAJgweI5Oiiw5/BOaGENM70Go+VM8LQMSxJ4S7/8MIJEZQp5HcJZ7XDTcEwruknrd8mllEfGyFzPvJOx6QAQocFhXBW6+AlhM3gn/dvV5vdrO8ihjET2GoDUqXPYC57ZuY+/Fz6W3KV8V97BvNUhpY5yQrP5VpnyvvXNFQtzDfClTvZFPuoHQi3/KYPi6O0FSD74vo8JOBZZY09boInPejkm9fvHQqfh0bnN7B6XJoUwC1Qprrx+XIy7ust5AEn5XL7d4lOvcR14MxDDKEp you@me.com"
  first_ed25519_public_key  = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIDqzSi9IHoYnbE3YQ+B2fQEVT8iGFemyPovpEtPziIVB you@me.com"
  second_ed25519_public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIDqzSi9IHoYnbE3YQ+B2fQEVT8iGFemyPovpEtPziIVB hello@world.com"
  admin_username            = "testadmin1234"
  admin_password            = "Password1234!"
  admin_password_update     = "Password1234!Update"
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-fleet-%[1]d"
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
    name                 = "internal"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  name            = "internal"
  loadbalancer_id = azurerm_lb.test.id
}


`, data.RandomInteger, location)
}

func (r ComputeFleetTestResource) additionalLinuxTemplate(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`

resource "azurerm_resource_group" "linux_test" {
  name     = "acctest-rg-fleet-al-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "linux_test" {
  name                = "acctvn-al-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.linux_test.location
  resource_group_name = azurerm_resource_group.linux_test.name
}

resource "azurerm_subnet" "linux_test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.linux_test.name
  virtual_network_name = azurerm_virtual_network.linux_test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_public_ip" "linux_test" {
  name                = "acctestpublicIP%[1]d"
  location            = azurerm_resource_group.linux_test.location
  resource_group_name = azurerm_resource_group.linux_test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  zones               = ["1"]
}

resource "azurerm_lb" "linux_test" {
  name                = "acctest-loadbalancer-%[1]d"
  location            = azurerm_resource_group.linux_test.location
  resource_group_name = azurerm_resource_group.linux_test.name
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "internal"
    public_ip_address_id = azurerm_public_ip.linux_test.id
  }
}

resource "azurerm_lb_backend_address_pool" "linux_test" {
  name            = "internal"
  loadbalancer_id = azurerm_lb.linux_test.id
}
`, data.RandomInteger, location)
}

func (r ComputeFleetTestResource) additionalWindowsTemplate(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "windows_test" {
  name     = "acctest-rg-fleet-al-win-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "windows_test" {
  name                = "acctvn-al-win-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.windows_test.location
  resource_group_name = azurerm_resource_group.windows_test.name
}

resource "azurerm_subnet" "windows_test" {
  name                 = "acctsub-al-win-%[1]d"
  resource_group_name  = azurerm_resource_group.windows_test.name
  virtual_network_name = azurerm_virtual_network.windows_test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_public_ip" "windows_test" {
  name                = "acctestpublicIP-al-%[1]d"
  location            = azurerm_resource_group.windows_test.location
  resource_group_name = azurerm_resource_group.windows_test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  zones               = ["1"]
}

resource "azurerm_lb" "windows_test" {
  name                = "acctest-loadbalancer-al-%[1]d"
  location            = azurerm_resource_group.windows_test.location
  resource_group_name = azurerm_resource_group.windows_test.name
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "internal-al-%[1]d"
    public_ip_address_id = azurerm_public_ip.windows_test.id
  }
}

resource "azurerm_lb_backend_address_pool" "windows_test" {
  name            = "internal-al"
  loadbalancer_id = azurerm_lb.windows_test.id
}
`, data.RandomInteger, location)
}

func (r ComputeFleetTestResource) completeExceptVMSS(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%[1]s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_user_assigned_identity" "test2" {
  name                = "acctest2%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_compute_fleet" "test" {
  name                        = "acctest-fleet-%[2]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = "%[3]s"
  platform_fault_domain_count = 1
  compute_api_version         = "2024-03-01"

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  plan {
    name           = "arcsight_logger_72_byol"
    product        = "arcsight-logger"
    publisher      = "micro-focus"
    promotion_code = "test"
  }

  spot_priority_profile {
    min_capacity     = 0
    maintain_enabled = false
    capacity         = 0
  }

  regular_priority_profile {
    capacity     = 0
    min_capacity = 0
  }


  vm_sizes_profile {
    name = "Standard_DS1_v2"
  }

  additional_location_profile {
    location = "%[4]s"
    %[5]s
  }

  virtual_machine_profile {
    network_api_version = "2020-11-01"
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
  tags = {
    Hello = "There"
    World = "Example"
  }
  zones = ["1", "2", "3"]
}
	`, r.template(data, data.Locations.Primary), data.RandomInteger, data.Locations.Primary, data.Locations.Secondary, r.additionalLinuxVirtualMachineProfile())
}

func (r ComputeFleetTestResource) spotPriorityProfile(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_compute_fleet" "test" {
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

  %[4]s

  zones = ["1", "2", "3"]

}
`, r.template(data, data.Locations.Primary), data.RandomInteger, data.Locations.Primary, r.baseLinuxVirtualMachineProfile())
}

func (r ComputeFleetTestResource) spotPriorityProfileUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_compute_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  spot_priority_profile {
    allocation_strategy     = "PriceCapacityOptimized"
    eviction_policy         = "Delete"
    maintain_enabled        = true
    max_hourly_price_per_vm = -1
    capacity                = 2
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

  %[4]s

  zones = ["1", "2", "3"]
}
`, r.template(data, data.Locations.Primary), data.RandomInteger, data.Locations.Primary, r.baseLinuxVirtualMachineProfile())
}

func (r ComputeFleetTestResource) spotAndRegulaPriorityProfile(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_compute_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  spot_priority_profile {
    min_capacity     = 1
    maintain_enabled = false
    capacity         = 1
  }

  regular_priority_profile {
    capacity     = 1
    min_capacity = 1
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

  %[4]s

  zones = ["1", "2", "3"]

}
`, r.template(data, data.Locations.Primary), data.RandomInteger, data.Locations.Primary, r.baseLinuxVirtualMachineProfile())
}

func (r ComputeFleetTestResource) spotAndRegulaPriorityProfileUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_compute_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  spot_priority_profile {
    min_capacity     = 1
    maintain_enabled = false
    capacity         = 2
  }

  regular_priority_profile {
    allocation_strategy = "LowestPrice"
    min_capacity        = 1
    capacity            = 2
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

  %[4]s

  zones = ["1", "2", "3"]

}
`, r.template(data, data.Locations.Primary), data.RandomInteger, data.Locations.Primary, r.baseLinuxVirtualMachineProfile())
}

func (r ComputeFleetTestResource) zones(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_compute_fleet" "test" {
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

  zones = ["1", "2"]
}
`, r.template(data, data.Locations.Primary), data.RandomInteger, data.Locations.Primary, r.baseLinuxVirtualMachineProfile())
}

func (r ComputeFleetTestResource) plan(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%[1]s

resource "azurerm_compute_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  plan {
    name           = "arcsight_logger_72_byol"
    product        = "arcsight-logger"
    publisher      = "micro-focus"
    promotion_code = "test"
  }

  regular_priority_profile {
    capacity     = 1
    min_capacity = 1
  }

  vm_sizes_profile {
    name = "Standard_DS1_v2"
  }

  virtual_machine_profile {
    network_api_version = "2020-11-01"
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
	`, r.template(data, data.Locations.Primary), data.RandomInteger, data.Locations.Primary)
}

func (r ComputeFleetTestResource) vmSizeProfile(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_compute_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  spot_priority_profile {
    allocation_strategy     = "LowestPrice"
    eviction_policy         = "Delete"
    maintain_enabled        = false
    max_hourly_price_per_vm = 1
    min_capacity            = 0
    capacity                = 2
  }

  vm_sizes_profile {
    name = "Standard_DS1_v2"
  }

  %[4]s

  zones = ["1", "2", "3"]
}
`, r.template(data, data.Locations.Primary), data.RandomInteger, data.Locations.Primary, r.baseLinuxVirtualMachineProfile())
}

func (r ComputeFleetTestResource) vmSizeProfileUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_compute_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  spot_priority_profile {
    allocation_strategy     = "LowestPrice"
    eviction_policy         = "Delete"
    maintain_enabled        = false
    max_hourly_price_per_vm = 1
    min_capacity            = 0
    capacity                = 2
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

  %[4]s

  zones = ["1", "2", "3"]
}
`, r.template(data, data.Locations.Primary), data.RandomInteger, data.Locations.Primary, r.baseLinuxVirtualMachineProfile())
}
