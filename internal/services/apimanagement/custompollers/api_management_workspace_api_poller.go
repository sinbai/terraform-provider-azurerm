// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/api"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

var _ pollers.PollerType = &apiManagementWorkspaceAPIPoller{}

var _ pollers.PollerType = &apiManagementWorkspaceAPIDeletePoller{}

var (
	workspaceApiPollingSuccess = pollers.PollResult{
		Status:       pollers.PollingStatusSucceeded,
		PollInterval: 10 * time.Second,
	}
	workspaceApiPollingInProgress = pollers.PollResult{
		Status:       pollers.PollingStatusInProgress,
		PollInterval: 10 * time.Second,
	}
)

type apiManagementWorkspaceAPIPoller struct {
	client  *api.ApiClient
	id      api.WorkspaceApiId
	asyncID string
}

type apiManagementWorkspaceAPIDeletePoller struct {
	client     *api.ApiClient
	pollingUrl *url.URL
}

type workspaceApiOptions struct {
	asyncId string
}

// NewAPIManagementWorkspaceAPIPoller - creates a new poller for API Management Workspace API operations to handle the case there is a query string
// parameter "asyncId" in the Location header of the response. This is used to poll the status of the operation.
func NewAPIManagementWorkspaceAPIPoller(cli *api.ApiClient, id api.WorkspaceApiId, response *http.Response) *apiManagementWorkspaceAPIPoller {
	urlStr := response.Header.Get("location")
	var asyncId string
	if u, err := url.Parse(urlStr); err == nil {
		asyncId = u.Query().Get("asyncId")
	}

	// sometimes the poller is not required as the API directly return 200
	if asyncId == "" {
		return nil
	}

	return &apiManagementWorkspaceAPIPoller{
		client:  cli,
		id:      id,
		asyncID: asyncId,
	}
}

func (p workspaceApiOptions) ToHeaders() *client.Headers {
	return &client.Headers{}
}

func (p workspaceApiOptions) ToOData() *odata.Query {
	return &odata.Query{}
}

func (p workspaceApiOptions) ToQuery() *client.QueryParams {
	q := client.QueryParams{}
	q.Append("asyncId", p.asyncId)
	return &q
}

func (p apiManagementWorkspaceAPIPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	if p.asyncID == "" {
		return &workspaceApiPollingSuccess, nil
	}

	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
			http.StatusAccepted,
			http.StatusCreated,
		},
		HttpMethod: http.MethodGet,
		Path:       p.id.ID(),
		OptionsObject: workspaceApiOptions{
			asyncId: p.asyncID,
		},
	}
	req, err := p.client.Client.NewRequest(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("building request: %+v", err)
	}

	resp, err := p.client.Client.Execute(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	// the response actually doesn't include a provisioningState property, so we only check the http status code
	switch resp.StatusCode {
	case http.StatusOK:
		return &workspaceApiPollingSuccess, nil
	case http.StatusAccepted, http.StatusCreated:
		return &workspaceApiPollingInProgress, nil
	}

	return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
}

// NewAPIManagementWorkspaceApiDeletePoller - creates a new poller for API Management Workspace API long-running delete operation to handle the case where the delete operation is asynchronous.
func NewAPIManagementWorkspaceApiDeletePoller(cli *api.ApiClient, response *http.Response) (*apiManagementWorkspaceAPIDeletePoller, error) {
	pollingUrl := response.Header.Get("Azure-AsyncOperation")
	if pollingUrl == "" {
		pollingUrl = response.Header.Get("Location")
	}

	if pollingUrl == "" {
		return nil, fmt.Errorf("no polling URL found in response")
	}

	url, err := url.Parse(pollingUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid polling URL %q in response: %v", pollingUrl, err)
	}
	if !url.IsAbs() {
		return nil, fmt.Errorf("invalid polling URL %q in response: URL was not absolute", pollingUrl)
	}

	return &apiManagementWorkspaceAPIDeletePoller{
		client:     cli,
		pollingUrl: url,
	}, nil
}

func (p apiManagementWorkspaceAPIDeletePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	if p.pollingUrl == nil {
		return nil, fmt.Errorf("internal error: cannot poll without a pollingUrl")
	}

	reqOpts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
			http.StatusCreated,
			http.StatusAccepted,
			http.StatusNoContent,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: nil,
		Path:          p.pollingUrl.Path,
	}

	req, err := p.client.Client.NewRequest(ctx, reqOpts)
	if err != nil {
		return nil, fmt.Errorf("building request: %+v", err)
	}

	resp, err := p.client.Client.Execute(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", p.pollingUrl.String(), err)
	}

	if resp.Response != nil {
		var respBody []byte
		respBody, err = io.ReadAll(resp.Response.Body)
		if err != nil {
			return nil, fmt.Errorf("parsing response body: %+v", err)
		}
		resp.Response.Body.Close()

		resp.Response.Body = io.NopCloser(bytes.NewReader(respBody))

		if s, ok := resp.Response.Header["Retry-After"]; ok {
			if sleep, err := strconv.ParseInt(s[0], 10, 64); err == nil {
				pollingDeleteInProgress.PollInterval = time.Second * time.Duration(sleep)
			}
		}

		// 202's don't necessarily return a body, so there's nothing to deserialize
		if resp.StatusCode == http.StatusAccepted && resp.ContentLength == 0 {
			return &workspaceApiPollingInProgress, nil
		}

		// returns a 200 OK with no Body
		if resp.StatusCode == http.StatusOK && resp.ContentLength == 0 {
			return &workspaceApiPollingSuccess, nil
		}

		if resp.Response.StatusCode == http.StatusOK {
			contentType := resp.Response.Header.Get("Content-Type")
			var op operationResult
			if strings.Contains(strings.ToLower(contentType), "application/json") {
				if err = json.Unmarshal(respBody, &op); err != nil {
					return nil, fmt.Errorf("unmarshalling response body: %+v", err)
				}
			} else {
				return nil, fmt.Errorf("internal-error: polling support for the Content-Type %q was not implemented: %+v", contentType, err)
			}

			switch string(op.Status) {
			case string(pollers.PollingStatusInProgress):
				return &workspaceApiPollingInProgress, nil
			case string(pollers.PollingStatusSucceeded):
				return &workspaceApiPollingSuccess, nil
			}
		}
	}

	return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
}
