// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package networkanalytics_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/networkfunction/2022-11-01/azuretrafficcollectors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NetworkAnalyticsDataProductTestResource struct{}

func TestAccNetworkFunctionAzureTrafficCollector_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_analytics_data_product", "test")
	r := NetworkAnalyticsDataProductTestResource{}
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

func TestAccNetworkFunctionAzureTrafficCollector_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_analytics_data_product", "test")
	r := NetworkAnalyticsDataProductTestResource{}
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

func TestAccNetworkFunctionAzureTrafficCollector_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_analytics_data_product", "test")
	r := NetworkAnalyticsDataProductTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkFunctionAzureTrafficCollector_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_analytics_data_product", "test")
	r := NetworkAnalyticsDataProductTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
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

func (r NetworkAnalyticsDataProductTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := azuretrafficcollectors.ParseAzureTrafficCollectorID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.NetworkFunction.AzureTrafficCollectorsClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r NetworkAnalyticsDataProductTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r NetworkAnalyticsDataProductTestResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_network_analytics_data_product" "test" {
  name                = "acctest-nadp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r NetworkAnalyticsDataProductTestResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_analytics_data_product" "import" {
  name                = azurerm_network_analytics_data_product.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
}
`, config, data.Locations.Primary)
}

func (r NetworkAnalyticsDataProductTestResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_analytics_data_product" "test" {
  name                = "acctest-nfatc-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  tags = {
    key = "value"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r NetworkAnalyticsDataProductTestResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_analytics_data_product" "test" {
  name                = "acctest-nfatc-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  tags = {
    key = "value2"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}
