// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementWorkspaceGlobalSchemaTestResource struct{}

func TestAccApiManagementWorkspaceGlobalSchema_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_global_schema", "test")
	r := ApiManagementWorkspaceGlobalSchemaTestResource{}
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

func TestAccApiManagementWorkspaceGlobalSchema_xmlSchema(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_global_schema", "test")
	r := ApiManagementWorkspaceGlobalSchemaTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.xmlSchema(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.xmlSchemaUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.xmlSchema(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementWorkspaceGlobalSchema_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_global_schema", "test")
	r := ApiManagementWorkspaceGlobalSchemaTestResource{}
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

func TestAccApiManagementWorkspaceGlobalSchema_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_global_schema", "test")
	r := ApiManagementWorkspaceGlobalSchemaTestResource{}
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

func TestAccApiManagementWorkspaceGlobalSchema_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_global_schema", "test")
	r := ApiManagementWorkspaceGlobalSchemaTestResource{}
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

func (r ApiManagementWorkspaceGlobalSchemaTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := schema.ParseWorkspaceSchemaID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.GlobalSchemaClient_v2024_05_01.WorkspaceGlobalSchemaGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ApiManagementWorkspaceGlobalSchemaTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_global_schema" "test" {
  name                         = "acctest-schema-%d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  type                        = "json"
  value                       = jsonencode({
    type = "object"
    properties = {
      id = {
        type = "integer"
      }
    }
  })
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceGlobalSchemaTestResource) xmlSchema(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_global_schema" "test" {
  name                         = "acctest-schema-%d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  type                        = "xml"
  value                       = <<XML
    <xsd:schema xmlns:xsd="http://www.w3.org/2001/XMLSchema"
    xmlns:tns="http://tempuri.org/PurchaseOrderSchema.xsd" targetNamespace="http://tempuri.org/PurchaseOrderSchema.xsd" elementFormDefault="qualified">
    <xsd:element name="PurchaseOrder" type="tns:PurchaseOrderType"/>
    <xsd:complexType name="PurchaseOrderType">
        <xsd:sequence>
            <xsd:element name="ShipTo" type="tns:USAddress" maxOccurs="2"/>
            <xsd:element name="BillTo" type="tns:USAddress"/>
        </xsd:sequence>
        <xsd:attribute name="OrderDate" type="xsd:date"/>
    </xsd:complexType>
    <xsd:complexType name="USAddress">
        <xsd:sequence>
            <xsd:element name="name" type="xsd:string"/>
            <xsd:element name="street" type="xsd:string"/>
            <xsd:element name="city" type="xsd:string"/>
            <xsd:element name="state" type="xsd:string"/>
            <xsd:element name="zip" type="xsd:integer"/>
        </xsd:sequence>
        <xsd:attribute name="country" type="xsd:NMTOKEN" fixed="US"/>
    </xsd:complexType>
</xsd:schema>
XML
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceGlobalSchemaTestResource) xmlSchemaUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_global_schema" "test" {
  name                         = "acctest-schema-%d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  type                        = "xml"
  value                       = <<XML
    <xsd:schema xmlns:xsd="http://www.w3.org/2001/XMLSchema"
    xmlns:tns="http://tempuri.org/PurchaseOrderSchema.xsd" targetNamespace="http://tempuri.org/PurchaseOrderSchema.xsd" elementFormDefault="qualified">
    <xsd:element name="PurchaseOrder" type="tns:PurchaseOrderType"/>
    <xsd:complexType name="PurchaseOrderType">
        <xsd:sequence>
            <xsd:element name="BillTo" type="tns:USAddress"/>
            <xsd:element name="ShipTo" type="tns:USAddress" maxOccurs="3"/>
            <xsd:element name="DiscountCode" type="xsd:string" minOccurs="0"/>
        </xsd:sequence>
        <xsd:attribute name="OrderDate" type="xsd:dateTime"/>
    </xsd:complexType>
    <xsd:complexType name="USAddress">
        <xsd:sequence>
            <xsd:element name="name" type="xsd:string"/>
            <xsd:element name="street" type="xsd:string"/>
            <xsd:element name="city" type="xsd:string"/>
            <xsd:element name="state" type="xsd:string"/>
            <xsd:element name="zip" type="xsd:string"/>
            <xsd:element name="phone" type="xsd:string" minOccurs="0"/>
        </xsd:sequence>
        <xsd:attribute name="country" type="xsd:NMTOKEN" fixed="USA"/>
    </xsd:complexType>
</xsd:schema>
XML
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceGlobalSchemaTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_global_schema" "import" {
  name                         = azurerm_api_management_workspace_global_schema.test.name
  api_management_workspace_id = azurerm_api_management_workspace_global_schema.test.api_management_workspace_id
  type                        = azurerm_api_management_workspace_global_schema.test.type
  value                       = azurerm_api_management_workspace_global_schema.test.value
}
`, r.basic(data))
}

func (r ApiManagementWorkspaceGlobalSchemaTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_global_schema" "test" {
  name                         = "acctest-schema-%d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  type                        = "json"
  description                 = "A JSON schema for testing"
  value                       = jsonencode({
    type = "object"
    properties = {
      id = {
        type = "integer"
      }
      name = {
        type = "string"
      }
    }
  })
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceGlobalSchemaTestResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_global_schema" "test" {
  name                         = "acctest-schema-%d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  description                 = "An updated JSON schema for testing"
  type                = "json"
  value                 = jsonencode({
    type = "object"
    properties = {
      id = {
        type = "integer"
      }
      full_name = {
        type = "string"
      }
      age = {
        type = "integer"
      }
    }
  })
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceGlobalSchemaTestResource) template(data acceptance.TestData) string {
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
  name               = "acctest-ws-%d"
  api_management_id = azurerm_api_management.test.id
  display_name      = "Test Workspace"
  description       = "A test workspace"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
