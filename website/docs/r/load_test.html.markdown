---
subcategory: "Load Test"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_load_test"
description: |-
  Manages a Load Test.
---

<!-- Note: This documentation is generated. Any manual changes will be overwritten -->

# azurerm_load_test

Manages a Load Test Service.

## Example Usage

```hcl
resource "azurerm_load_example" "example" {
  location            = azurerm_resource_group.example.location
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Load Test should exist. Changing this forces a new Load Test to be created.

* `name` - (Required) Specifies the name of this Load Test. Changing this forces a new Load Test to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group within which this Load Test should exist. Changing this forces a new Load Test to be created.

* `customer_managed_key` - (Optional) A `customer_managed_key` block as defined below.

* `description` - (Optional) Description of the resource. Changing this forces a new Load Test to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Load Test.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Load Test. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Load Test.

~> **NOTE:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

A `customer_managed_key` block supports the following:

* `key_vault_key_id` - (Required) The ID of the Key Vault Key which should be used to Encrypt the data in this Load Test.

* `user_assigned_identity_id` - (Required) The ID of the User Assigned Identity that has access to the key.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Load Test.

* `data_plane_uri` - Resource data plane URI.

---



## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating this Load Test.
* `delete` - (Defaults to 30 minutes) Used when deleting this Load Test.
* `read` - (Defaults to 5 minutes) Used when retrieving this Load Test.
* `update` - (Defaults to 30 minutes) Used when updating this Load Test.

## Import

An existing Load Test can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_load_test.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.LoadTestService/loadTests/{loadTestName}
```

* Where `{subscriptionId}` is the ID of the Azure Subscription where the Load Test exists. For example `12345678-1234-9876-4563-123456789012`.
* Where `{resourceGroupName}` is the name of Resource Group where this Load Test exists. For example `example-resource-group`.
* Where `{loadTestName}` is the name of the Load Test. For example `loadTestValue`.
