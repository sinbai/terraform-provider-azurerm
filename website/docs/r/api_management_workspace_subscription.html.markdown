---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_workspace_subscription"
description: |-
  Manages a Subscription in an API Management Workspace.
---

# azurerm_api_management_workspace_subscription

Manages a Subscription in an API Management Workspace.

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
  publisher_email     = "company@example.com"
  sku_name            = "Developer_1"
}

resource "azurerm_api_management_workspace" "example" {
  name              = "example-workspace"
  api_management_id = azurerm_api_management.example.id
  display_name      = "Example Workspace"
}

resource "azurerm_api_management_workspace_subscription" "example" {
  subscription_id = "example-subscription"
  workspace_id    = azurerm_api_management_workspace.example.id
  display_name    = "Example Workspace Subscription"
  state           = "active"
}

```

## Arguments Reference

The following arguments are supported:

* `api_management_workspace_id` - (Required) The ID of the API Management Workspace where the Subscription should exist. Changing this forces a new API Management Workspace Subscription to be created.

* `display_name` - (Required) The display name of the API Management Workspace Subscription.

* `subscription_name` - (Optional) The identifier of the API Management Workspace Subscription. If not specified, it will be generated automatically. Changing this forces a new resource to be created.

* `api_id` - (Optional) The ID of the API which should be assigned to the API Management Workspace Subscription. Changing this forces a new resource to be created.

~> **Note:** Only one of `product_id` or `api_id` can be set. If neither is specified, the Subscription provides access to all APIs within the workspace.

* `product_id` - (Optional) The ID of the Product which should be assigned to the API Management Workspace Subscription. Changing this forces a new resource to be created.

~> **Note:** Only one of `product_id` or `api_id` can be set. If neither is specified, the Subscription provides access to all APIs within the workspace.

* `owner_id` - (Optional) The ID of the owner which should be assigned to the API Management Workspace Subscription. Changing this forces a new resource to be created.

* `tracing_enabled` - (Optional) Determines whether tracing is enabled for the API Management Workspace Subscription. Defaults to `true`.

* `primary_key` - (Optional) The primary subscription key to use for the API Management Workspace Subscription.

* `secondary_key` - (Optional) The secondary subscription key to use for the API Management Workspace Subscription. 

~> **Note:** If `primary_key` or `secondary_key` not specified, a key will be generated automatically.

* `state` - (Optional) The state of the API Management Workspace Subscription. Possible values are `active`, `cancelled`, `expired`, `rejected`, `submitted` and `suspended`. Defaults to `active`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management Workspace Subscription.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Workspace Subscription.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Workspace Subscription.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Workspace Subscription.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Workspace Subscription.

## Import

API Management Workspace Subscriptions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_workspace_subscription.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/service1/workspaces/workspace1/subscriptions/subscription1
```
