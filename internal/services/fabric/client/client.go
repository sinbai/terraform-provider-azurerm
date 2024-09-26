// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	FabricCapacityClient *fabriccapacities.MongoClustersClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {

	fabricCapacityClient, err := fabriccapacities.NewMongoClustersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building fabric capacity client: %+v", err)
	}
	o.Configure(fabricCapacityClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		FabricCapacityClient: fabricCapacityClient,
	}, nil
}
