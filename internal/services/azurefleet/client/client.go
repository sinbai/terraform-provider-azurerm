package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurefleet/2024-05-01-preview/fleets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	FleetsClient *fleets.FleetsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	fleetsClient := fleets.NewFleetsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&fleetsClient.Client, o.ResourceManagerAuthorizer)

	fleetsClient, err := fleets.NewFleetsClientWithBaseURI(o.ResourceManagerEndpoint)
	if err != nil {
		return nil, fmt.Errorf("building fleets  client: %+v", err)
	}
	o.Configure(fleetsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		FleetsClient: fleetsClient,
	}, nil
}
