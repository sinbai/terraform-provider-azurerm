package cognitive_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2022-10-01/deployments"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type cognitiveDeploymentResource struct{}

func TestAcccognitiveDeployment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_deployment", "test")
	r := cognitiveDeploymentResource{}
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

func TestAcccognitiveDeployment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_deployment", "test")
	r := cognitiveDeploymentResource{}
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

func TestAcccognitiveDeployment_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_deployment", "test")
	r := cognitiveDeploymentResource{}
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

func TestAcccognitiveDeployment_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_deployment", "test")
	r := cognitiveDeploymentResource{}
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

func (r cognitiveDeploymentResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := deployments.ParseDeploymentID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Cognitive.DeploymentsClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r cognitiveDeploymentResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}
resource "azurerm_cognitive_account" "test" {
  name                = "acctest-ca-%d"
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r cognitiveDeploymentResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_cognitive_deployment" "test" {
  name                 = "acctest-cd-%d"
  cognitive_account_id = azurerm_cognitive_account.test.id
}
`, template, data.RandomInteger)
}

func (r cognitiveDeploymentResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_cognitive_deployment" "import" {
  name                 = azurerm_cognitive_deployment.test.name
  cognitive_account_id = azurerm_cognitive_account.test.id
}
`, config)
}

func (r cognitiveDeploymentResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_cognitive_deployment" "test" {
  name                 = "acctest-cd-%d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  rai_policy_name      = ""
  model {
    format  = ""
    name    = ""
    version = ""
  }
  scale_settings {
    capacity   = 0
    scale_type = ""
  }

}
`, template, data.RandomInteger)
}

func (r cognitiveDeploymentResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_cognitive_deployment" "test" {
  name                 = "acctest-cd-%d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  rai_policy_name      = ""
  model {
    format  = ""
    name    = ""
    version = ""
  }
  scale_settings {
    capacity   = 0
    scale_type = ""
  }

}
`, template, data.RandomInteger)
}
