package cognitive_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2022-10-01/commitmentplans"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type cognitiveCommitmentPlanResource struct{}

func TestAcccognitiveCommitmentPlan_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_commitment_plan", "test")
	r := cognitiveCommitmentPlanResource{}
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

func TestAcccognitiveCommitmentPlan_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_commitment_plan", "test")
	r := cognitiveCommitmentPlanResource{}
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

func TestAcccognitiveCommitmentPlan_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_commitment_plan", "test")
	r := cognitiveCommitmentPlanResource{}
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

func TestAcccognitiveCommitmentPlan_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_commitment_plan", "test")
	r := cognitiveCommitmentPlanResource{}
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

func (r cognitiveCommitmentPlanResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commitmentplans.ParseCommitmentPlanID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Cognitive.CommitmentPlansClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r cognitiveCommitmentPlanResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}
resource "azurerm_cognitive_account" "test" {
  name                = "acctestcogacc-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "SpeechServices"
  sku_name            = "S0"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r cognitiveCommitmentPlanResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_cognitive_commitment_plan" "test" {
  name                 = "acctest-ccp-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  cognitive_account_id = azurerm_cognitive_account.test.id

  hosting_model = "Web"
  plan_type     = "Speech2Text"
}
`, template, data.RandomInteger)
}

func (r cognitiveCommitmentPlanResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_cognitive_commitment_plan" "import" {
  name                 = azurerm_cognitive_commitment_plan.test.name
  cognitive_account_id = azurerm_cognitive_account.test.id
}
`, config)
}

func (r cognitiveCommitmentPlanResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_cognitive_commitment_plan" "test" {
  name                 = "acctest-ccp-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  cognitive_account_id = azurerm_cognitive_account.test.id
  auto_renew           = true
  hosting_model        = "Web"
  plan_type            = "Speech2Text"

  current {
    count = 0
    tier  = "T1"
  }

  next {
    count = 1
    tier  = "T2"
  }
}
`, template, data.RandomInteger)
}

func (r cognitiveCommitmentPlanResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_cognitive_commitment_plan" "test" {
  name                 = "acctest-ccp-%d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  auto_renew           = false
  hosting_model        = "Web"
  plan_type            = "Speech2Text"

  current {
    count = 1
    tier  = "T2"
  }

  next {
    count = 0
    tier  = "T1"
  }
}
`, template, data.RandomInteger)
}
