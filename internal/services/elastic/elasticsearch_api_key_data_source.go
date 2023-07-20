// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package elastic

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2023-06-01/apikey"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/elastic/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceElasticsearchApiKey() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceElasticsearchApiKeyRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"email_address": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ElasticEmailAddress,
			},
			"api_key": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceElasticsearchApiKeyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	apiKeyClient := meta.(*clients.Client).Elastic.ApiKeyClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	userEmail := apikey.UserEmailId{
		EmailId: pointer.To(d.Get("email_address").(string)),
	}

	subscriptionId := commonids.SubscriptionId{
		SubscriptionId: meta.(*clients.Client).Account.SubscriptionId,
	}

	respApiKey, err := apiKeyClient.OrganizationsGetApiKey(ctx, subscriptionId, userEmail)
	if err != nil {
		if !response.WasNotFound(respApiKey.HttpResponse) {
			return fmt.Errorf("retrieving apikey for %s: %+v", d.Get("email_address").(string), err)
		}
	}

	if model := respApiKey.Model; model != nil && model.Properties != nil {
		d.Set("api_key", pointer.From(model.Properties.ApiKey))
	}

	return nil
}
