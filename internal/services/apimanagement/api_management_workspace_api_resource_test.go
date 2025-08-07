// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/api"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementWorkspaceApiTestResource struct{}

func TestAccApiManagementWorkspaceApi_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_api", "test")
	r := ApiManagementWorkspaceApiTestResource{}

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

func TestAccApiManagementWorkspaceApi_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_api", "test")
	r := ApiManagementWorkspaceApiTestResource{}

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

func TestAccApiManagementWorkspaceApi_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_api", "test")
	r := ApiManagementWorkspaceApiTestResource{}

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

func TestAccApiManagementWorkspaceApi_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_api", "test")
	r := ApiManagementWorkspaceApiTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updated(data),
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
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementWorkspaceApi_importApi(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_api", "test")
	r := ApiManagementWorkspaceApiTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.importApi(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("import"),
	})
}

func TestAccApiManagementWorkspaceApi_sourceApiId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_api", "test")
	r := ApiManagementWorkspaceApiTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sourceApiId(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("source_api_id"),
	})
}

func TestAccApiManagementWorkspaceApi_apiVersionSet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_api", "test")
	r := ApiManagementWorkspaceApiTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.apiVersionSet(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (ApiManagementWorkspaceApiTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := api.ParseWorkspaceApiID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.V20240501ApiClient.WorkspaceApiGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ApiManagementWorkspaceApiTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_api_management_openid_connect_provider" "test" {
  name                = "acctest-%[2]d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  client_id           = "00001111-2222-3333-%[2]d"
  client_secret       = "%[2]d-cwdavsxbacsaxZX-%[2]d"
  display_name        = "Initial Name"
  metadata_endpoint   = "https://azacceptance.hashicorptest.com/example/foo"
}

resource "azurerm_api_management_authorization_server" "test" {
  name                         = "acctestauthsrv-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  api_management_name          = azurerm_api_management.test.name
  display_name                 = "Test Group"
  authorization_endpoint       = "https://azacceptance.hashicorptest.com/client/authorize"
  client_id                    = "42424242-4242-4242-4242-424242424242"
  client_registration_endpoint = "https://azacceptance.hashicorptest.com/client/register"

  grant_types = [
    "implicit",
  ]

  authorization_methods = [
    "GET",
  ]
}

resource "azurerm_api_management_workspace_api" "test" {
  name                        = "acctestapi-%[2]d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  revision                    = "1"
  display_name                = "api1"
  protocols                   = ["https"]
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceApiTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_workspace_api" "import" {
  name                        = azurerm_api_management_workspace_api.test.name
  api_management_workspace_id = azurerm_api_management_workspace_api.test.api_management_workspace_id
  display_name                = azurerm_api_management_workspace_api.test.display_name
  path                        = azurerm_api_management_workspace_api.test.path
  protocols                   = azurerm_api_management_workspace_api.test.protocols
  revision                    = azurerm_api_management_workspace_api.test.revision
}
`, r.basic(data))
}

func (r ApiManagementWorkspaceApiTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_api_management_openid_connect_provider" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  client_id           = "00001111-2222-3333-%[2]d"
  client_secret       = "%[2]d-cwdavsxbacsaxZX-%[2]d"
  display_name        = "Initial Name"
  metadata_endpoint   = "https://azacceptance.hashicorptest.com/example/foo"
}

resource "azurerm_api_management_authorization_server" "test" {
  name                         = "acctestauthsrv-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  api_management_name          = azurerm_api_management.test.name
  display_name                 = "Test Group"
  authorization_endpoint       = "https://azacceptance.hashicorptest.com/client/authorize"
  client_id                    = "42424242-4242-4242-4242-424242424242"
  client_registration_endpoint = "https://azacceptance.hashicorptest.com/client/register"

  grant_types = [
    "implicit",
  ]

  authorization_methods = [
    "GET",
  ]
}

resource "azurerm_api_management_workspace_api" "test" {
  name                        = "acctestapi-%[2]d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  api_type                    = "graphql"
  description                 = "Test API Description"
  display_name                = "Test API"
  path                        = "test"
  protocols                   = ["https"]
  revision                    = "1"
  revision_description        = "Test Revision Description"
  service_url                 = "https://api.example.com"
  subscription_enabled        = true
  terms_of_service_url        = "https://example.com/terms"

  contact {
    name  = "API Support"
    email = "support@example.com"
    url   = "https://example.com/support"
  }

  openid_authentication {
    openid_provider_name = azurerm_api_management_openid_connect_provider.test.name
    bearer_token_sending_methods = [
      "authorizationHeader",
      "query",
    ]
  }

  subscription_key_parameter_names {
    header = "X-API-KEY"
    query  = "api-key"
  }

  license {
    name = "MIT License"
    url  = "https://opensource.org/licenses/MIT"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceApiTestResource) updated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_api_management_openid_connect_provider" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  client_id           = "00001111-2222-3333-%[2]d"
  client_secret       = "%[2]d-cwdavsxbacsaxZX-%[2]d"
  display_name        = "Initial Name"
  metadata_endpoint   = "https://azacceptance.hashicorptest.com/example/foo"
}

resource "azurerm_api_management_authorization_server" "test" {
  name                         = "acctestauthsrv-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  api_management_name          = azurerm_api_management.test.name
  display_name                 = "Test Group"
  authorization_endpoint       = "https://azacceptance.hashicorptest.com/client/authorize"
  client_id                    = "42424242-4242-4242-4242-424242424242"
  client_registration_endpoint = "https://azacceptance.hashicorptest.com/client/register"

  grant_types = [
    "implicit",
  ]

  authorization_methods = [
    "GET",
  ]
}

resource "azurerm_api_management_workspace_api" "test" {
  name                        = "acctestapi-%[2]d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  api_type                    = "soap"
  description                 = "Updated API Description"
  display_name                = "Update Test API"
  path                        = "updatetest"
  protocols                   = ["https", "http"]
  revision                    = "1"
  revision_description        = "Test Revision2 Description"
  service_url                 = "https://updated-api.example.com"
  subscription_enabled        = false
  terms_of_service_url        = "https://example:8080/service"

  contact {
    email = "test@test.com"
    name  = "test"
    url   = "https://example:8080"
  }

  oauth2_authorization {
    authorization_server_name = azurerm_api_management_authorization_server.test.name
    scope                     = "acctest"
  }

  subscription_key_parameter_names {
    header = "X-Butter-Robot-API-Key"
    query  = "location"
  }

  license {
    name = "test-license"
    url  = "https://example:8080/license"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceApiTestResource) importApi(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_api" "test" {
  name                        = "acctestapi-%d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  display_name                = "api1"
  path                        = "api1"
  protocols                   = ["https"]
  revision                    = "1"

  import {
    content_value  = file("testdata/api_management_api_wsdl_multiple.xml")
    content_format = "wsdl"

    wsdl_selector {
      service_name  = "Calculator"
      endpoint_name = "CalculatorHttpsSoap11Endpoint"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceApiTestResource) apiVersionSet(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_api_version_set" "test" {
  name                        = "acctestversionset-%d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  display_name                = "Test Version Set"
  versioning_scheme           = "Segment"
}

resource "azurerm_api_management_workspace_api" "test" {
  name                        = "acctestapi-%d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  display_name                = "Test API"
  path                        = "test"
  protocols                   = ["https"]
  revision                    = "1"
  version                     = "v1"
  version_description         = "Version 1.0"
  version_set_id              = azurerm_api_management_workspace_api_version_set.test.id
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementWorkspaceApiTestResource) sourceApiId(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_api" "source" {
  name                        = "acctestsourceapi-%d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  display_name                = "Source API"
  path                        = "source"
  protocols                   = ["https"]
  revision                    = "1"
}

resource "azurerm_api_management_workspace_api" "test" {
  name                        = "acctestapi-%d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  display_name                = "Test API"
  path                        = "test"
  protocols                   = ["https"]
  revision                    = "1"
  source_api_id               = azurerm_api_management_workspace_api.source.id
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (ApiManagementWorkspaceApiTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`


resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
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
  name              = "acctestAMWS-%d"
  api_management_id = azurerm_api_management.test.id
  display_name      = "Test Workspace"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
