// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/subscription"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementWorkspaceSubscriptionTestResource struct{}

func TestAccApiManagementWorkspaceSubscription_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_subscription", "test")
	r := ApiManagementWorkspaceSubscriptionTestResource{}

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

func TestAccApiManagementWorkspaceSubscription_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_subscription", "test")
	r := ApiManagementWorkspaceSubscriptionTestResource{}

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

func TestAccApiManagementWorkspaceSubscription_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_subscription", "test")
	r := ApiManagementWorkspaceSubscriptionTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("state").HasValue("active"),
				check.That(data.ResourceName).Key("allow_tracing").HasValue("true"),
			),
		},
		data.ImportStep("primary_key", "secondary_key"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("state").HasValue("suspended"),
				check.That(data.ResourceName).Key("allow_tracing").HasValue("false"),
			),
		},
		data.ImportStep("primary_key", "secondary_key"),
	})
}

func TestAccApiManagementWorkspaceSubscription_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_subscription", "test")
	r := ApiManagementWorkspaceSubscriptionTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("primary_key", "secondary_key"),
	})
}

func TestAccApiManagementWorkspaceSubscription_withProduct(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_subscription", "test")
	r := ApiManagementWorkspaceSubscriptionTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withProduct(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("primary_key", "secondary_key"),
	})
}

func TestAccApiManagementWorkspaceSubscription_withApi(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_subscription", "test")
	r := ApiManagementWorkspaceSubscriptionTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withApi(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("primary_key", "secondary_key"),
	})
}

func (ApiManagementWorkspaceSubscriptionTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := subscription.ParseWorkspaceSubscriptions2ID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.SubscriptionClient_v2024_05_01.WorkspaceSubscriptionGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ApiManagementWorkspaceSubscriptionTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_api_management_product" "test" {
  product_id            = "test-product"
  api_management_name   = azurerm_api_management.test.name
  resource_group_name   = azurerm_resource_group.test.name
  display_name          = "Test Product"
  subscription_required = true
  approval_required     = false
  published             = true
}

resource "azurerm_api_management_user" "test" {
  user_id             = "acctestuser%[2]d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  first_name          = "Acceptance"
  last_name           = "Test"
  email               = "azure-acctest%[2]d@example.com"
}

%[1]s

resource "azurerm_api_management_workspace_subscription" "test" {
  subscription_name    = "test-subscription-%[2]d"
  api_management_workspace_id       = azurerm_api_management_workspace.test.id
  display_name        = "Test Workspace Subscription"
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceSubscriptionTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_workspace_subscription" "import" {
  subscription_name    = azurerm_api_management_workspace_subscription.test.subscription_name
  workspace_id       = azurerm_api_management_workspace_subscription.test.workspace_id
  display_name       = azurerm_api_management_workspace_subscription.test.display_name
}
`, r.basic(data))
}

func (r ApiManagementWorkspaceSubscriptionTestResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_subscription" "test" {
  subscription_name    = "test-subscription-%d"
  workspace_id       = azurerm_api_management_workspace.test.id
  display_name       = "Updated Test Workspace Subscription"
  state              = "suspended"
  allow_tracing      = false
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceSubscriptionTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_user" "test" {
  user_id             = "test-user-%d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  first_name          = "Acceptance"
  last_name           = "Test"
  email               = "azure-acctest%d@example.com"
}

resource "azurerm_api_management_workspace_subscription" "test" {
  subscription_name    = "test-subscription-%d"
  workspace_id       = azurerm_api_management_workspace.test.id
  display_name       = "Complete Test Workspace Subscription"
  owner_id            = azurerm_api_management_user.test.id
  state              = "active"
  allow_tracing      = true
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementWorkspaceSubscriptionTestResource) withProduct(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_product" "test" {
  product_id            = "test-product-%d"
  api_management_name   = azurerm_api_management.test.name
  resource_group_name   = azurerm_resource_group.test.name
  display_name          = "Test Product"
  subscription_required = true
  approval_required     = false
  published             = true
}

resource "azurerm_api_management_workspace_subscription" "test" {
  subscription_name    = "test-subscription-%d"
  workspace_id       = azurerm_api_management_workspace.test.id
  display_name       = "Test Workspace Subscription with Product"
  product_id         = azurerm_api_management_product.test.id
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementWorkspaceSubscriptionTestResource) withApi(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_api" "test" {
  name                = "test-api-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  revision            = "1"
  display_name        = "Test API"
  path                = "test"
  protocols           = ["https"]
}

resource "azurerm_api_management_workspace_subscription" "test" {
  subscription_name    = "test-subscription-%d"
  workspace_id       = azurerm_api_management_workspace.test.id
  display_name       = "Test Workspace Subscription with API"
  api_id             = azurerm_api_management_api.test.id
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (ApiManagementWorkspaceSubscriptionTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-apim-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Premium_1"
}

resource "azurerm_api_management_workspace" "test" {
  name              = "acctest-amws-%d"
  api_management_id = azurerm_api_management.test.id
  display_name      = "Test Workspace"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
