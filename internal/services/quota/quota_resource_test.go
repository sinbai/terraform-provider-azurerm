package quota_test

import (
	"context"
	"fmt"

	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/quota/2025-03-01/quotainformation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type QuotaResource struct{}

func TestAccQuota_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_quota", "test")
	r := QuotaResource{}
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

func TestAccQuota_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_quota", "test")
	r := QuotaResource{}
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

func TestAccQuota_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_quota", "test")
	r := QuotaResource{}
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

func TestAccQuota_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_quota", "test")
	r := QuotaResource{}
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
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r QuotaResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := quotainformation.ParseScopedQuotaID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Quota.QuotaClient.QuotaGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r QuotaResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_quota" "test" {
  name               = "standardFSv2Family"
  location           = "%s"
  limit_object_type  = "LimitValue"
  limit_value        = 22
  resource_provider = "Microsoft.Compute"
}
`, data.Locations.Ternary)
}

func (r QuotaResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_quota" "import" {
  name                = azurerm_quota.test.name
  location            = azurerm_quota.test.location
  limit_object_type   = azurerm_quota.test.limit_object_type
  limit_value         = azurerm_quota.test.limit_value
  resource_provider   = azurerm_quota.test.resource_provider
}
`, config)
}

func (r QuotaResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_quota" "test" {
  name               = "TotalLowPriorityCores"
  location           = "%s"
  limit_object_type  = "LimitValue"
  limit_value        = 26
  resource_provider = "Microsoft.MachineLearningServices"
  resource_type  = "lowPriority"
  additional_properties = jsonencode({
      "region"= "eastus",
      "sku"="Standard_D2_v2"
  })

 lifecycle {
    ignore_changes = [resource_type, additional_properties]
 }
}
`, "eastus") //data.Locations.Ternary)
}

func (r QuotaResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_quota" "test" {
  name               = "TotalLowPriorityCores"
  location           = "%s"
  limit_object_type  = "LimitValue"
  limit_value        = 28
  resource_provider = "Microsoft.MachineLearningServices"
  resource_type  = "lowPriority"
  additional_properties = jsonencode({
      "region"= "eastus",
      "sku"="Standard_D2_v2"
  })

 lifecycle {
    ignore_changes = [resource_type, additional_properties]
 }
}
`, "eastus") //data.Locations.Ternary)
}
