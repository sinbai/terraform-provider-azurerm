---
subcategory: "cognitive"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cognitive_commitment_plan"
description: |-
  Manages a Cognitive Commitment Plans.
---

# azurerm_cognitive_commitment_plan

Manages a Cognitive Commitment Plans.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_cognitive_account" "example" {
  name                = "example-ca"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_cognitive_commitment_plan" "example" {
  name                 = "example-ccp"
  cognitive_account_id = azurerm_cognitive_account.test.id
  auto_renew           = false
  hosting_model        = ""
  plan_type            = ""
  current {
    count = 0
    tier  = ""
  }
  next {
    count = 0
    tier  = ""
  }

}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Cognitive Commitment Plans. Changing this forces a new Cognitive Commitment Plans to be created.

* `cognitive_account_id` - (Required) Specifies the ID of the Cognitive Commitment Plans. Changing this forces a new Cognitive Commitment Plans to be created.

* `auto_renew` - (Optional) AutoRenew commitment plan.

* `current` - (Optional) A `current` block as defined below.

* `hosting_model` - (Optional) Account hosting model.

* `next` - (Optional) A `next` block as defined below.

* `plan_type` - (Optional) Commitment plan type.

---

A `current` block supports the following:

* `count` - (Optional) Commitment period commitment count.

* `tier` - (Optional) Commitment period commitment tier.

---

A `next` block supports the following:

* `count` - (Optional) Commitment period commitment count.

* `tier` - (Optional) Commitment period commitment tier.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Cognitive Commitment Plans.

* `current` - A `current` block as defined below.

* `last` - A `last` block as defined below.

* `next` - A `next` block as defined below.

---

A `current` block exports the following:

* `end_date` - Commitment period end date.

* `quota` - A `quota` block as defined below.

* `start_date` - Commitment period start date.

---

A `quota` block exports the following:

* `quantity` - Commitment quota quantity.

* `unit` - Commitment quota unit.

---

A `last` block exports the following:

* `count` - Commitment period commitment count.

* `end_date` - Commitment period end date.

* `quota` - A `quota` block as defined below.

* `start_date` - Commitment period start date.

* `tier` - Commitment period commitment tier.

---

A `quota` block exports the following:

* `quantity` - Commitment quota quantity.

* `unit` - Commitment quota unit.

---

A `next` block exports the following:

* `end_date` - Commitment period end date.

* `quota` - A `quota` block as defined below.

* `start_date` - Commitment period start date.

---

A `quota` block exports the following:

* `quantity` - Commitment quota quantity.

* `unit` - Commitment quota unit.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Cognitive Commitment Plans.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cognitive Commitment Plans.
* `update` - (Defaults to 30 minutes) Used when updating the Cognitive Commitment Plans.
* `delete` - (Defaults to 30 minutes) Used when deleting the Cognitive Commitment Plans.

## Import

Cognitive Commitment Plans can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cognitive_commitment_plan.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.CognitiveServices/accounts/account1/commitmentPlans/commitmentPlan1
```
