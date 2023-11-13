// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

type Client struct {
	LoadBalancersClient *loadbalancers.LoadBalancersClient
	//LoadBalancerBackendAddressPoolsClient *network.LoadBalancerBackendAddressPoolsClient
	LoadBalancingRulesClient *network.LoadBalancerLoadBalancingRulesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	loadBalancersClient, err := loadbalancers.NewLoadBalancersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building load balancers client: %+v", err)
	}
	o.Configure(loadBalancersClient.Client, o.Authorizers.ResourceManager)

	//loadBalancerBackendAddressPoolsClient := network.NewLoadBalancerBackendAddressPoolsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	//o.ConfigureClient(&loadBalancerBackendAddressPoolsClient.Client, o.ResourceManagerAuthorizer)

	loadBalancingRulesClient := network.NewLoadBalancerLoadBalancingRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&loadBalancingRulesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		LoadBalancersClient: loadBalancersClient,
		//LoadBalancerBackendAddressPoolsClient: &loadBalancerBackendAddressPoolsClient,
		LoadBalancingRulesClient: &loadBalancingRulesClient,
	}, nil
}
