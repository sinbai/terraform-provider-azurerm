package client

import (
	"fmt"
	"github.com/hashicorp/go-azure-sdk/resource-manager/quota/2025-03-01/quotainformation"
	"github.com/hashicorp/go-azure-sdk/resource-manager/quota/2025-03-01/usagesinformation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	QuotaClient  *quotainformation.QuotaInformationClient
	UsagesClient *usagesinformation.UsagesInformationClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {

	quotaClient, err := quotainformation.NewQuotaInformationClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building quota client: %+v", err)
	}
	o.Configure(quotaClient.Client, o.Authorizers.ResourceManager)

	usagesClient, err := usagesinformation.NewUsagesInformationClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building usages client: %+v", err)
	}
	o.Configure(usagesClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		QuotaClient:  quotaClient,
		UsagesClient: usagesClient,
	}, nil
}
