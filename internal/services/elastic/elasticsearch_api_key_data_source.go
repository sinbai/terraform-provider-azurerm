// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package elastic

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2023-06-01/apikey"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/elastic/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.DataSource = ElasticsearchApiKeyDataSource{}

type ElasticsearchApiKeyDataSource struct{}

type ElasticsearchApiKeyDataSourceModel struct {
	EmailAddress string `tfschema:"email_address"`
	ApiKey       string `tfschema:"api_key"`
}

func (e ElasticsearchApiKeyDataSource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"email_address": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ElasticEmailAddress,
		},
	}
}

func (e ElasticsearchApiKeyDataSource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"api_key": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (e ElasticsearchApiKeyDataSource) ModelObject() interface{} {
	return &ElasticsearchApiKeyDataSourceModel{}
}

func (e ElasticsearchApiKeyDataSource) ResourceType() string {
	return "azurerm_elastic_cloud_elasticsearch_api_key"
}

func (e ElasticsearchApiKeyDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			apiKeyClient := metadata.Client.Elastic.ApiKeyClient

			var model ElasticsearchApiKeyDataSourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			userEmail := apikey.UserEmailId{
				EmailId: pointer.To(model.EmailAddress),
			}

			subscriptionId := commonids.SubscriptionId{
				SubscriptionId: metadata.Client.Account.SubscriptionId,
			}

			resp, err := apiKeyClient.OrganizationsGetApiKey(ctx, subscriptionId, userEmail)
			if err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("retrieving api key for %s: %+v", model.EmailAddress, err)
				}
			}

			if m := resp.Model; m != nil && m.Properties != nil {
				model.ApiKey = pointer.From(m.Properties.ApiKey)
			}

			return metadata.Encode(&model)
		},
		Timeout: 5 * time.Minute,
	}
}
