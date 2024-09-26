// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mongocluster_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type FabricCapacityTestResource struct{}

func TestAccMongoCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_fabric_capacity", "test")
	r := FabricCapacityTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password", "create_mode"),
	})
}

func TestAccMongoCluster_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_fabric_capacity", "test")
	r := FabricCapacityTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password", "create_mode"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password", "create_mode"),
	})
}

func TestAccMongoCluster_previewFeature(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_fabric_capacity", "test")
	r := FabricCapacityTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.enablePreviewFeature(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password", "create_mode"),
	})
}

func TestAccMongoCluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_fabric_capacity", "test")
	r := FabricCapacityTestResource{}

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

func (r FabricCapacityTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := fabriccapacities.ParseMongoClusterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.MongoCluster.MongoClustersClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r FabricCapacityTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_fabric_capacity" "test" {
  name                         = "acctest-mc%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  administrator_login          = "adminTerraform"
  administrator_login_password = "QAZwsx123"
  shard_count                  = "1"
  compute_tier                 = "Free"
  high_availability_mode       = "Disabled"
  storage_size_in_gb           = "32"
  version                      = "6.0"
}
`, r.template(data), data.RandomInteger)
}

func (r FabricCapacityTestResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_fabric_capacity" "test" {
  name                         = "acctest-mc%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  administrator_login          = "adminTerraform"
  administrator_login_password = "testQAZwsx123update"
  shard_count                  = "1"
  compute_tier                 = "M30"
  high_availability_mode       = "ZoneRedundantPreferred"
  storage_size_in_gb           = "64"
  version                      = "7.0"
}
`, r.template(data), data.RandomInteger)
}

func (r FabricCapacityTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_fabric_capacity" "import" {
  name                         = azurerm_fabric_capacity.test.name
  resource_group_name          = azurerm_fabric_capacity.test.resource_group_name
  location                     = azurerm_fabric_capacity.test.location
  administrator_login          = azurerm_fabric_capacity.test.administrator_login
  administrator_login_password = azurerm_fabric_capacity.test.administrator_login_password
  shard_count                  = azurerm_fabric_capacity.test.shard_count
  compute_tier                 = azurerm_fabric_capacity.test.compute_tier
  high_availability_mode       = azurerm_fabric_capacity.test.high_availability_mode
  storage_size_in_gb           = azurerm_fabric_capacity.test.storage_size_in_gb
  version                      = azurerm_fabric_capacity.test.version
}
`, r.basic(data))
}

func (r FabricCapacityTestResource) enablePreviewFeature(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_fabric_capacity" "test" {
  name                         = "acctest-mc%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  administrator_login          = "adminTerraform"
  administrator_login_password = "testQAZwsx123"
  shard_count                  = "1"
  compute_tier                 = "M30"
  high_availability_mode       = "ZoneRedundantPreferred"
  storage_size_in_gb           = "64"
  preview_features             = ["GeoReplica"]
  version                      = "7.0"

}

resource "azurerm_fabric_capacity" "geo_replica" {
  name                         = "acctest-mc-replica%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  source_server_id             = azurerm_fabric_capacity.test.id
  source_location              = "%s"
  create_mode                  = "GeoReplica"
}

`, r.template(data), data.RandomInteger, data.RandomInteger, data.Locations.Secondary)
}

func (r FabricCapacityTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}
