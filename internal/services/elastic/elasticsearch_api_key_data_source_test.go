package elastic_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ElasticsearchApiKeyDataSource struct{}

func TestAccElasticsearchApiKeyDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_elastic_cloud_elasticsearch_api_key", "test")
	r := ElasticsearchApiKeyDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("api_key").Exists(),
			),
		},
	})
}

func (ElasticsearchApiKeyDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  client_id               = ""
  client_certificate_path = ""
  client_secret           = ""
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-elastic-%[1]d"
  location = "%[2]s"
  lifecycle {
    ignore_changes = ["tags"]
  }
}

resource "azurerm_elastic_cloud_elasticsearch" "test" {
  name                        = "acctest-estc%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sku_name                    = "ess-monthly-consumption_Monthly"
  elastic_cloud_email_address = "v-elenaxin@microsoft.com"
  lifecycle {
    ignore_changes = ["tags"]
  }
}

data "azurerm_elastic_cloud_elasticsearch_api_key" "test" {
  email_address = "v-elenaxin@microsoft.com"
  depends_on    = [azurerm_elastic_cloud_elasticsearch.test]
}
`, data.RandomInteger, data.Locations.Primary)
}
