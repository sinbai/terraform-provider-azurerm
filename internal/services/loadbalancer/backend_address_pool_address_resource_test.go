package loadbalancer_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-02-01/network"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var _ types.TestResourceVerifyingRemoved = BackendAddressPoolAddressResourceTests{}

type BackendAddressPoolAddressResourceTests struct{}

func TestAccBackendAddressPoolAddressBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool_address", "test")
	r := BackendAddressPoolAddressResourceTests{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackendAddressPoolAddressWithFrontendIPConfigAndSubnet(t *testing.T) {
	skip(t)

	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool_address", "test")
	r := BackendAddressPoolAddressResourceTests{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withFEIPConfigAndSubnet(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackendAddressPoolAddressWithFrontendIPConfigAndSubnetUpdate(t *testing.T) {
	skip(t)

	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool_address", "test")
	r := BackendAddressPoolAddressResourceTests{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withFEIPConfigAndSubnet(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateFEIPConfigAndSubnet(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackendAddressPoolAddressRequiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool_address", "test")
	r := BackendAddressPoolAddressResourceTests{}
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

func TestAccBackendAddressPoolAddressDisappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool_address", "test")
	r := BackendAddressPoolAddressResourceTests{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccBackendAddressPoolAddressUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool_address", "test")
	r := BackendAddressPoolAddressResourceTests{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (BackendAddressPoolAddressResourceTests) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.BackendAddressPoolAddressID(state.ID)
	if err != nil {
		return nil, err
	}

	pool, err := client.LoadBalancers.LoadBalancerBackendAddressPoolsClient.Get(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if pool.BackendAddressPoolPropertiesFormat == nil {
		return nil, fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	if pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses != nil {
		for _, address := range *pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses {
			if address.Name == nil {
				continue
			}

			if *address.Name == id.AddressName {
				return utils.Bool(true), nil
			}
		}
	}
	return utils.Bool(false), nil
}

func (BackendAddressPoolAddressResourceTests) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.BackendAddressPoolAddressID(state.ID)
	if err != nil {
		return nil, err
	}

	pool, err := client.LoadBalancers.LoadBalancerBackendAddressPoolsClient.Get(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if pool.BackendAddressPoolPropertiesFormat == nil {
		return nil, fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	addresses := make([]network.LoadBalancerBackendAddress, 0)
	if pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses != nil {
		addresses = *pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses
	}
	newAddresses := make([]network.LoadBalancerBackendAddress, 0)
	for _, address := range addresses {
		if address.Name == nil {
			continue
		}

		if *address.Name != id.AddressName {
			newAddresses = append(newAddresses, address)
		}
	}
	pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses = &newAddresses

	future, err := client.LoadBalancers.LoadBalancerBackendAddressPoolsClient.CreateOrUpdate(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName, pool)
	if err != nil {
		return nil, fmt.Errorf("updating %s: %+v", *id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.LoadBalancers.LoadBalancerBackendAddressPoolsClient.Client); err != nil {
		return nil, fmt.Errorf("waiting for update of %s: %+v", *id, err)
	}
	return utils.Bool(true), nil
}

// nolint unused - for future use
func (BackendAddressPoolAddressResourceTests) backendAddressPoolHasAddresses(expected int) acceptance.ClientCheckFunc {
	return func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
		id, err := parse.LoadBalancerBackendAddressPoolID(state.ID)
		if err != nil {
			return err
		}

		client := clients.LoadBalancers.LoadBalancerBackendAddressPoolsClient
		pool, err := client.Get(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)
		if err != nil {
			return err
		}
		if pool.BackendAddressPoolPropertiesFormat == nil {
			return fmt.Errorf("`properties` is nil")
		}
		if pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses == nil {
			return fmt.Errorf("`properties.loadBalancerBackendAddresses` is nil")
		}

		actual := len(*pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses)
		if actual != expected {
			return fmt.Errorf("expected %d but got %d addresses", expected, actual)
		}

		return nil
	}
}

func (t BackendAddressPoolAddressResourceTests) basic(data acceptance.TestData) string {
	template := t.template(data, string(network.LoadBalancerSkuTierRegional))
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_lb_backend_address_pool_address" "test" {
  name                    = "address"
  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id
  virtual_network_id      = azurerm_virtual_network.test.id
  ip_address              = "191.168.0.1"
}
`, template)
}

func (t BackendAddressPoolAddressResourceTests) withFEIPConfigAndSubnet(data acceptance.TestData) string {
	template := t.template(data, string(network.LoadBalancerSkuTierGlobal))
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_subnet" "test" {
  name                 = "acctest-subnet-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes       = ["192.168.1.0/24"]
}

resource "azurerm_public_ip" "test1" {
  name                = "acctest-pip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  sku_tier            = "Regional"
}


resource "azurerm_public_ip" "test2" {
  name                = "acctest-another-pip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  sku_tier            = "Regional"
}

resource "azurerm_lb" "test1" {
  name                = "acctest-lb-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  sku_tier            = "Regional"

  frontend_ip_configuration {
    name                 = "feip1"
    public_ip_address_id = azurerm_public_ip.test1.id
  }

  frontend_ip_configuration {
    name                 = "feip2"
    public_ip_address_id = azurerm_public_ip.test2.id
  }
}

resource "azurerm_lb_backend_address_pool_address" "test" {
  name                         = "acctest-Address"
  backend_address_pool_id      = azurerm_lb_backend_address_pool.test.id
  virtual_network_id           = azurerm_virtual_network.test.id
  frontend_ip_configuration_id = azurerm_lb.test1.frontend_ip_configuration[0].id
  subnet_id                    = azurerm_subnet.test.id
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (t BackendAddressPoolAddressResourceTests) updateFEIPConfigAndSubnet(data acceptance.TestData) string {
	template := t.template(data, string(network.LoadBalancerSkuTierGlobal))
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_subnet" "test1" {
  name                 = "acctest-another-subnet-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes       = ["192.168.1.0/24"]
}

resource "azurerm_public_ip" "test1" {
  name                = "acctest-pip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  sku_tier            = "Regional"
}


resource "azurerm_public_ip" "test2" {
  name                = "acctest-another-pip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  sku_tier            = "Regional"
}

resource "azurerm_lb" "test1" {
  name                = "acctest-lb-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  sku_tier            = "Regional"

  frontend_ip_configuration {
    name                 = "feip1"
    public_ip_address_id = azurerm_public_ip.test1.id
  }

  frontend_ip_configuration {
    name                 = "feip2"
    public_ip_address_id = azurerm_public_ip.test2.id
  }
}

resource "azurerm_lb_backend_address_pool_address" "test" {
  name                         = "acctest-Address"
  backend_address_pool_id      = azurerm_lb_backend_address_pool.test.id
  virtual_network_id           = azurerm_virtual_network.test.id
  frontend_ip_configuration_id = azurerm_lb.test1.frontend_ip_configuration[1].id
  subnet_id                    = azurerm_subnet.test1.id
}

`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (t BackendAddressPoolAddressResourceTests) requiresImport(data acceptance.TestData) string {
	template := t.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_lb_backend_address_pool_address" "import" {
  name                    = azurerm_lb_backend_address_pool_address.test.name
  backend_address_pool_id = azurerm_lb_backend_address_pool_address.test.backend_address_pool_id
  virtual_network_id      = azurerm_lb_backend_address_pool_address.test.virtual_network_id
  ip_address              = azurerm_lb_backend_address_pool_address.test.ip_address
}
`, template)
}

func (t BackendAddressPoolAddressResourceTests) update(data acceptance.TestData) string {
	template := t.template(data, string(network.LoadBalancerSkuTierRegional))
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_lb_backend_address_pool_address" "test" {
  name                    = "address"
  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id
  virtual_network_id      = azurerm_virtual_network.test.id
  ip_address              = "191.168.0.2"
}
`, template)
}

func (BackendAddressPoolAddressResourceTests) template(data acceptance.TestData, lbSkuTier string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["192.168.0.0/16"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  sku_tier            = "%s"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  sku_tier            = "%s"

  frontend_ip_configuration {
    name                 = "feip"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  name            = "internal"
  loadbalancer_id = azurerm_lb.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, lbSkuTier, data.RandomInteger, lbSkuTier)
}

func skip(t *testing.T) {
	supportedLocs := map[string]int{
		azure.NormalizeLocation("East US 2"):      1,
		azure.NormalizeLocation("West US"):        2,
		azure.NormalizeLocation("West Europe"):    3,
		azure.NormalizeLocation("Southeast Asia"): 4,
		azure.NormalizeLocation("Central US"):     5,
		azure.NormalizeLocation("North Europe"):   6,
		azure.NormalizeLocation("East Asia"):      7,
	}

	location := os.Getenv("ARM_TEST_LOCATION")

	if _, ok := supportedLocs[azure.NormalizeLocation(location)]; !ok {
		t.Skip(fmt.Sprintf("Skipping as the cross-region load balancer or Public IP in Global tier can only be deployed to %q regions", "East US 2,West US,West Europe,Southeast Asia,Central US,North Europe,East Asia"))
	}
}
