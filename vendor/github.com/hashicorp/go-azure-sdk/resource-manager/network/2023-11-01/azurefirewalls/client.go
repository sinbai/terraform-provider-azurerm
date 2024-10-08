package azurefirewalls

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureFirewallsClient struct {
	Client *resourcemanager.Client
}

func NewAzureFirewallsClientWithBaseURI(sdkApi sdkEnv.Api) (*AzureFirewallsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "azurefirewalls", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AzureFirewallsClient: %+v", err)
	}

	return &AzureFirewallsClient{
		Client: client,
	}, nil
}
