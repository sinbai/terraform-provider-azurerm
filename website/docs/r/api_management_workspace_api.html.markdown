---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_workspace_api"
description: |-
  Manages an API Management Workspace API.
---

# azurerm_api_management_workspace_api

Manages an API Management Workspace API.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_api_management" "example" {
  name                = "example-apim"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  publisher_name      = "My Company"
  publisher_email     = "company@terraform.io"
  sku_name            = "Premium_1"
}

resource "azurerm_api_management_workspace" "example" {
  name              = "example-workspace"
  api_management_id = azurerm_api_management.example.id
  display_name      = "Example Workspace"
}

resource "azurerm_api_management_workspace_api" "example" {
  name                          = "example-api"
  api_management_workspace_id   = azurerm_api_management_workspace.example.id
  api_type                      = "http"
  description                   = "Test API Description"
  display_name                  = "Test API"
  path                          = "test"
  protocols                     = ["https"]
  revision                      = "1"
  revision_description          = "Test Revision Description"
  service_url                   = "https://api.example.com"
  subscription_required_enabled = true
  terms_of_service_url          = "https://example.com/terms"

  contact {
    name  = "API Support"
    email = "support@example.com"
    url   = "https://example.com/support"
  }

  openid_authentication {
    openid_provider_name = "openidProviderExample"
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
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the API Management Workspace API. Changing this forces a new resource to be created.

* `api_management_workspace_id` - (Required) Specifies the ID of the API Management Workspace. Changing this forces a new resource to be created.

* `revision` - (Required) Specifies the revision of the API Management Workspace API. Changing this forces a new resource to be created.

* `api_type` - (Optional) Specifies the type of the API Management Workspace API. Possible values are `graphql`, `grpc`, `http`, `odata`, `soap` and `websocket`.

* `contact` - (Optional) A `contact` block as defined below.

* `description` - (Optional) Specifies the description of the API Management Workspace API.

* `display_name` - (Optional) Specifies the display name of the API Management Workspace API.

* `import` - (Optional) An `import` block as defined below.

* `license` - (Optional) A `license` block as defined below.

* `oauth2_authorization` - (Optional) An `oauth2_authorization` block as defined below.

* `openid_authentication` - (Optional) An `openid_authentication` block as defined below.

* `path` - (Optional) Specifies the path for the API Management Workspace API.

* `protocols` - (Optional) Specifies the set of protocols on which the operations in the API Management Workspace API can be invoked. Possible values are `Http`, `Https`, `Ws` and `Wss`.

* `revision_description` - (Optional) Specifies the description of the API Management Workspace API revision.

* `service_url` - (Optional) Specifies the absolute URL of the backend service implementing the API Management Workspace API.

* `source_api_id` - (Optional) Specifies the ID of the source API Management Workspace API to be copied.

* `subscription_key_parameter_names` - (Optional) A `subscription_key_parameter_names` block as defined below.

* `subscription_enabled` - (Optional) Whether the API Management Workspace API require a subscription key is enabled. Defaults to `true`.

* `terms_of_service_url` - (Optional) Specifies the absolute URL of the terms of service for the API Management Workspace API.

* `version` - (Optional) Specifies the version number of this API Management Workspace API, if versioning is used.

* `version_description` - (Optional) Specifies the description of the API Management Workspace API version.

* `version_set_id` - (Optional) Specifies the ID of the version set associated with the API Management Workspace API.

---

A `contact` block supports the following:

* `email` - (Optional) The email address of the contact person/organization.

* `name` - (Optional) The name of the contact person/organization.

* `url` - (Optional) The URL pointing to the contact information.

---

An `import` block supports the following:

* `content_value` - (Required) The content used to import the API Management Workspace API Definition. 

~> **Note:** If the `content_format` of `*-link-*` is specified, `content_value` must be a URL; otherwise, it must be provided inline.

* `content_format` - (Required) The format of the content to be imported. Possible values are `openapi`, `openapi+json`, `openapi+json-link`, `openapi-link`, `swagger-json`, `swagger-link-json`, `wadl-link-json`, `wadl-xml`, `wsdl`, `wsdl-link`, `grpc`, `grpc-link`, `odata`, `odata-link` and `graphql-link`.

* `wsdl_selector` - (Optional) A `wsdl_selector` block as defined below.

~> **Note:** The `wsdl_selector` allows you to limit the import of a WSDL to a specific service and endpoint. It can only be specified when `content_format` is set to `wsdl` or `wsdl-link`.

---

A `wsdl_selector` block supports the following:

* `service_name` - (Required) The Name of the WSDL service to import.

* `endpoint_name` - (Required) The name of the endpoint (port) to import.

---

A `license` block supports the following:

* `name` - (Optional) The name of the license.

* `url` - (Optional) The URL of the license used by the API Management Workspace API.

---

A `oauth2_authorization` block supports the following:

* `authorization_server_name` - (Required) The name of the OAuth 2.0 authorization server to be used.

* `scope` - (Optional) The scope that defines the operations or access level permitted.

---

A `openid_authentication` block supports the following:

* `openid_provider_name` - (Required) The name of the OpenID Connect Provider.

* `bearer_token_sending_methods` - (Optional) Specifies a list of supported methods for sending the token to the server. Possible values are `authorizationHeader` and `query`.

---

A `subscription_key_parameter_names` block supports the following:

* `header` - (Optional) The name of the HTTP Header which should be used for the subscription Key.

* `query` - (Optional) The name of the querystring parameter which should be used for the subscription Key.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management Workspace API.

* `is_current` - Is this the current API Revision?

* `is_online` - Is this API Revision online/accessible via the Gateway?

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Workspace API.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Workspace API.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Workspace API.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Workspace API.

## Import

API Management Workspace APIs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_workspace_api.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/service1/workspaces/workspace1/apis/api1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.ApiManagement` - 2024-05-01
