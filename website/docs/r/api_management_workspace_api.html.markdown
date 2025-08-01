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
  name                = "example-workspace"
  api_management_id   = azurerm_api_management.example.id
  display_name        = "Example Workspace"
}

resource "azurerm_api_management_workspace_api" "example" {
  name              = "example-api"
  api_management_id = azurerm_api_management.example.id
  workspace_name    = azurerm_api_management_workspace.example.name
  revision          = "1"
  display_name      = "Example API"
  protocols         = ["https"]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the API Management Workspace API. Changing this forces a new resource to be created.

* `api_management_workspace_id` - (Required) The ID of the API Management Workspace. Changing this forces a new resource to be created.

* `revision` - (Required) The Revision of the API Management Workspace API. Changing this forces a new resource to be created.

* `display_name` - (Optional) The display name of the API Management Workspace API.

* `path` - (Optional) The path for the API Management Workspace API.

* `protocols` - (Optional) A set of protocols on which the operations in the API Management Workspace API can be invoked. Possible values are `Http`, `Https`, `Ws` and `Wss`.

* `api_type` - (Optional) The type of the API Management Workspace API. Possible values are `graphql`, `grpc`, `http`, `odata`, `soap` and `websocket`.

* `contact` - (Optional) A `contact` block as defined below.

* `description` - (Optional) A description of the API Management Workspace API.

* `import` - (Optional) An `import` block as defined below.

* `license` - (Optional) A `license` block as defined below.

* `oauth2_authorization` - (Optional) An `oauth2_authorization` block as defined below.

* `openid_authentication` - (Optional) An `openid_authentication` block as defined below.

* `revision_description` - (Optional) A description of the API Management Workspace API revision.

* `service_url` - (Optional) Absolute URL of the backend service implementing the API Management Workspace API.

* `source_api_id` - (Optional) The ID of the source API Management Workspace API to be copied.

* `subscription_key_parameter_names` - (Optional) A `subscription_key_parameter_names` block as defined below.

* `subscription_required` - (Optional) Should the API Management Workspace API require a subscription key? Defaults to `true`.

* `terms_of_service_url` - (Optional) Absolute URL of the terms of service for the API Management Workspace API.

* `version` - (Optional) The version number of this API Management Workspace API, if the API Management Workspace API is versioned.

* `version_description` - (Optional) A description of the API Management Workspace API version.

* `version_set_id` - (Optional) The ID of the version set which the API Management Workspace API is associated with.

---

A `oauth2_authorization` block supports the following:

* `authorization_server_name` - (Required) OAuth authorization server name.

* `scope` - (Optional) Operations scope.

---

A `openid_authentication` block supports the following:

* `openid_provider_name` - (Required) OpenID Connect provider identifier. The name of an [OpenID Connect Provider](https://www.terraform.io/docs/providers/azurerm/r/api_management_openid_connect_provider.html).

* `bearer_token_sending_methods` - (Optional) How to send the token to the server. A list of zero or more methods. Possible values are `authorizationHeader` and `query`.

---

A `contact` block supports the following:

* `email` - (Optional) The email address of the contact person/organization.

* `name` - (Optional) The name of the contact person/organization.

* `url` - (Optional) The URL pointing to the contact information.

---

An `import` block supports the following:

* `content_value` - (Required) The content from which the API Management Workspace API Definition should be imported. When a `content_format` of `*-link-*` is specified, this must be a URL, otherwise this must be defined inline.

* `content_format` - (Required) The format of the content from which the API Management Workspace API Definition should be imported. Possible values are: `openapi`, `openapi+json`, `openapi+json-link`, `openapi-link`, `swagger-json`, `swagger-link-json`, `wadl-link-json`, `wadl-xml`, `wsdl`, `wsdl-link` and `graphql-link`.

* `wsdl_selector` - (Optional) A `wsdl_selector` block as defined below, which allows you to limit the import of a WSDL to only a subset of the document. This can only be specified when `content_format` is `wsdl` or `wsdl-link`.

---

A `wsdl_selector` block supports the following:

* `service_name` - (Required) The name of the service to import from WSDL.

* `endpoint_name` - (Required) The name of the endpoint (port) to import from WSDL.

---

A `license` block supports the following:

* `name` - (Optional) The name of the license.

* `url` - (Optional) The URL to the license used for the API Management Workspace API.

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