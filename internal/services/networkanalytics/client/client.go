// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/networkanalytics/2023-11-15/dataproducts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	DataProductsClient dataproducts.DataProductsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	dataProductsClient, err := dataproducts.NewDataProductsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Network Analytics Data Product Client: %+v", err)
	}

	return &Client{
		DataProductsClient: *dataProductsClient,
	}, nil
}
