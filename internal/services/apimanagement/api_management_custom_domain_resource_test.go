// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apimanagementservice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementCustomDomainTestResource struct{}

func TestAccApiManagementCustomDomain_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_custom_domain", "test")
	r := ApiManagementCustomDomainTestResource{}

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

func TestAccApiManagementCustomDomain_basicWithUserIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_custom_domain", "test")
	r := ApiManagementCustomDomainTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithUserIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementCustomDomain_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_custom_domain", "test")
	r := ApiManagementCustomDomainTestResource{}

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

func TestAccApiManagementCustomDomain_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_custom_domain", "test")
	r := ApiManagementCustomDomainTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.proxyOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.developerPortalOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccApiManagementCustomDomain_certificateManaged(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_custom_domain", "test")
	r := ApiManagementCustomDomainTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.certificateManaged(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.certificateManagedUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("gateway.0.certificate", // not returned from API, sensitive
			"gateway.0.certificate_password", // not returned from API, sensitive
		),
		{
			Config: r.certificateManagedUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.certificateManaged(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (ApiManagementCustomDomainTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.CustomDomainID(state.ID)
	if err != nil {
		return nil, err
	}

	serviceId := apimanagementservice.NewServiceID(id.SubscriptionId, id.ResourceGroup, id.ServiceName)
	resp, err := clients.ApiManagement.ServiceClient.Get(ctx, serviceId)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil && resp.Model.Id != nil), nil
}

func (r ApiManagementCustomDomainTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_custom_domain" "test" {
  api_management_id = azurerm_api_management.test.id

  gateway {
    host_name    = "api.example.com"
    key_vault_id = azurerm_key_vault_certificate.test.secret_id
  }

  developer_portal {
    host_name    = "portal.example.com"
    key_vault_id = azurerm_key_vault_certificate.test.secret_id
  }
}
`, r.template(data, true))
}

func (r ApiManagementCustomDomainTestResource) proxyOnly(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_custom_domain" "test" {
  api_management_id = azurerm_api_management.test.id

  gateway {
    host_name    = "api.example.com"
    key_vault_id = azurerm_key_vault_certificate.test.secret_id
  }
}
`, r.template(data, true))
}

func (r ApiManagementCustomDomainTestResource) developerPortalOnly(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_custom_domain" "test" {
  api_management_id = azurerm_api_management.test.id

  developer_portal {
    host_name    = "portal.example.com"
    key_vault_id = azurerm_key_vault_certificate.test.secret_id
  }
}
`, r.template(data, true))
}

func (r ApiManagementCustomDomainTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_custom_domain" "import" {
  api_management_id = azurerm_api_management_custom_domain.test.api_management_id

  gateway {
    host_name    = "api.example.com"
    key_vault_id = azurerm_key_vault_certificate.test.secret_id
  }

  developer_portal {
    host_name    = "portal.example.com"
    key_vault_id = azurerm_key_vault_certificate.test.secret_id
  }
}
`, r.basic(data))
}

func (ApiManagementCustomDomainTestResource) template(data acceptance.TestData, systemAssignedIdentity bool) string {
	identitySnippet := `
  identity {
    type = "SystemAssigned"
  }
`
	if !systemAssignedIdentity {
		identitySnippet = `
  identity {
    type = "SystemAssigned, UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }
`
	}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Developer_1"

%[4]s

}

resource "azurerm_key_vault" "test" {
  name                = "apimkv%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    certificate_permissions = [
      "Create",
      "Delete",
      "Get",
      "Update",
      "Purge",
    ]

    key_permissions = [
      "Create",
      "Get",
    ]

    secret_permissions = [
      "Get",
    ]
  }

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = azurerm_api_management.test.identity.0.principal_id

    certificate_permissions = [
      "Get",
    ]

    secret_permissions = [
      "Get",
    ]
  }
}

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctestcert%[3]s"
  key_vault_id = azurerm_key_vault.test.id

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }

    lifetime_action {
      action {
        action_type = "AutoRenew"
      }

      trigger {
        days_before_expiry = 30
      }
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }

    x509_certificate_properties {
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]

      subject            = "CN=api.example.com"
      validity_in_months = 12

      subject_alternative_names {
        dns_names = [
          "api.example.com",
          "portal.example.com",
        ]
      }
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, identitySnippet)
}

func (r ApiManagementCustomDomainTestResource) basicWithUserIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "bp-user-example"
}

resource "azurerm_api_management_custom_domain" "test" {
  api_management_id = azurerm_api_management.test.id

  gateway {
    host_name    = "api.example.com"
    key_vault_id = azurerm_key_vault_certificate.test.secret_id
  }

  developer_portal {
    host_name    = "portal.example.com"
    key_vault_id = azurerm_key_vault_certificate.test.secret_id
  }
}
`, r.template(data, false))
}

func (r ApiManagementCustomDomainTestResource) certificateManaged(data acceptance.TestData) string {
	dnsZone := "sinbai.store"                                        //os.Getenv("ARM_TEST_DNS_ZONE")
	dnsZoneResourceGroup := "elenatest"                              //os.Getenv("ARM_TEST_DNS_ZONE_RESOURCE_GROUP")
	txtRecordValue := "7fGTFuan5vQ2VZNUc4DGEPZ9c+WMCWZ+5tfl2ejjtMg=" //os.Getenv("ARM_TEST_DNS_ZONE_RESOURCE_GROUP")
	return fmt.Sprintf(`

%[1]s

data "azurerm_dns_zone" "test" {
  name                = "%[2]s"
  resource_group_name = "%[3]s"
}

resource "azurerm_dns_cname_record" "test" {
  name                = "%[4]s"
  zone_name           = data.azurerm_dns_zone.test.name
  resource_group_name = data.azurerm_dns_zone.test.resource_group_name
  ttl                 = 3600
  record              = "${azurerm_api_management.test.name}.azure-api.net"
}

resource "azurerm_dns_txt_record" "test" {
  name                = "apimuid.%[4]s"
  resource_group_name = data.azurerm_dns_zone.test.resource_group_name
  zone_name           = data.azurerm_dns_zone.test.name
  ttl                 = 3600
  record {
    value = "%[5]s"
  }
}

resource "azurerm_api_management_custom_domain" "test" {
  api_management_id = azurerm_api_management.test.id
  
  gateway {
    host_name           = "${azurerm_dns_cname_record.test.name}.${data.azurerm_dns_zone.test.name}"
    certificate_source  = "Managed"
  }

  developer_portal {
    host_name    = "portal.example.com"
    key_vault_id = azurerm_key_vault_certificate.test.secret_id
  }
}
`, r.template(data, true), dnsZone, dnsZoneResourceGroup, data.RandomString, txtRecordValue)
}

func (r ApiManagementCustomDomainTestResource) certificateManagedUpdate(data acceptance.TestData) string {
	dnsZone := "sinbai.store"                                        //os.Getenv("ARM_TEST_DNS_ZONE")
	dnsZoneResourceGroup := "elenatest"                              //os.Getenv("ARM_TEST_DNS_ZONE_RESOURCE_GROUP")
	txtRecordValue := "7fGTFuan5vQ2VZNUc4DGEPZ9c+WMCWZ+5tfl2ejjtMg=" //os.Getenv("ARM_TEST_DNS_ZONE_RESOURCE_GROUP")
	return fmt.Sprintf(`

%[1]s

data "azurerm_dns_zone" "test" {
  name                = "%[2]s"
  resource_group_name = "%[3]s"
}

resource "azurerm_dns_cname_record" "test" {
  name                = "%[4]s"
  zone_name           = data.azurerm_dns_zone.test.name
  resource_group_name = data.azurerm_dns_zone.test.resource_group_name
  ttl                 = 3600
  record              = "${azurerm_api_management.test.name}.azure-api.net"
}

resource "azurerm_dns_txt_record" "test" {
  name                = "apimuid.%[4]s"
  resource_group_name = data.azurerm_dns_zone.test.resource_group_name
  zone_name           = data.azurerm_dns_zone.test.name
  ttl                 = 3600
  record {
    value = "%[5]s"
  }
}

resource "azurerm_api_management_custom_domain" "test" {
  api_management_id = azurerm_api_management.test.id
  
  gateway {
    host_name            = "api.terraform.io"
    certificate          = filebase64("testdata/api_management_api_test.pfx")
    certificate_password = "terraform"
    negotiate_client_certificate = false
  }

  developer_portal {
    host_name    = "portal.example.com"
    key_vault_id = azurerm_key_vault_certificate.test.secret_id
  }
}
`, r.template(data, true), dnsZone, dnsZoneResourceGroup, data.RandomString, txtRecordValue)
}
