package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/azurefleet/2024-11-01/fleets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	FleetsClient *fleets.FleetsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	fleetsClient, err := fleets.NewFleetsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building fleets  client: %+v", err)
	}
	o.Configure(fleetsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		FleetsClient: fleetsClient,
	}, nil
}
