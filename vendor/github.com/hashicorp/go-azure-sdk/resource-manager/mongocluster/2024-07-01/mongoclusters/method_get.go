package mongoclusters

import (
	"context"
<<<<<<< HEAD

	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"

=======
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
>>>>>>> 7a921d7afc5b9cf5038ddcdec068d7c1c5160c66
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *MongoCluster
}

// Get ...
func (c MongoClustersClient) Get(ctx context.Context, id MongoClusterId) (result GetOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       id.ID(),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	var model MongoCluster
	result.Model = &model
<<<<<<< HEAD

=======
>>>>>>> 7a921d7afc5b9cf5038ddcdec068d7c1c5160c66
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
