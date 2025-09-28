---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_workspace_global_schema"
description: |-
  Manages a Global Schema within an API Management Workspace.
---

# azurerm_api_management_workspace_global_schema

Manages a Global Schema within an API Management Workspace.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_api_management" "example" {
  name                = "example-apimanagement"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Premium_1"
}

resource "azurerm_api_management_workspace" "example" {
  name              = "example-workspace"
  api_management_id = azurerm_api_management.example.id
  display_name      = "Example Workspace"
}

resource "azurerm_api_management_workspace_global_schema" "example" {
  name                         = "example-schema"
  api_management_workspace_id = azurerm_api_management_workspace.example.id
  type                        = "json"
  description                 = "Example JSON schema"
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
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for the API Management Workspace Global Schema. Changing this forces a new resource to be created.

* `api_management_workspace_id` - (Required) Specifies the ID of the API Management Workspace in which the Global Schema should be created. Changing this forces a new resource to be created.

* `type` - (Required) Specifies the content type of the API Management Workspace Global Schema. Possible values are `xml` and `json`. Changing this forces a new resource to be created.

* `value` - (Required) Specifies the string defining the document representing the API Management Workspace Global Schema.

---

* `description` - (Optional) Specifies the description of the API Management Workspace Global Schema.

~> **Note:** Once set it cannot be removed. Removing it will force a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management Workspace Global Schema.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Workspace Global Schema.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Workspace Global Schema.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Workspace Global Schema.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Workspace Global Schema.

## Import

API Management Workspace Global Schemas can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_workspace_global_schema.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/instance1/workspaces/workspace1/schemas/schema1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.ApiManagement` - 2024-05-01