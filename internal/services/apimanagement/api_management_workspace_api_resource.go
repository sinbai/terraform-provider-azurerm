// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/workspace"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/api"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ApiManagementWorkspaceApiResource struct{}

var _ sdk.ResourceWithUpdate = ApiManagementWorkspaceApiResource{}

var _ sdk.ResourceWithCustomizeDiff = ApiManagementWorkspaceApiResource{}

type ApiManagementWorkspaceApiModel struct {
	Name                          string                               `tfschema:"name"`
	ApiManagementWorkspaceId      string                               `tfschema:"api_management_workspace_id"`
	DisplayName                   string                               `tfschema:"display_name"`
	Path                          string                               `tfschema:"path"`
	Protocols                     []string                             `tfschema:"protocols"`
	Revision                      string                               `tfschema:"revision"`
	RevisionDescription           string                               `tfschema:"revision_description"`
	ApiType                       string                               `tfschema:"api_type"`
	Contact                       []ContactModel                       `tfschema:"contact"`
	Description                   string                               `tfschema:"description"`
	Import                        []ImportModel                        `tfschema:"import"`
	License                       []LicenseModel                       `tfschema:"license"`
	ServiceUrl                    string                               `tfschema:"service_url"`
	SubscriptionKeyParameterNames []SubscriptionKeyParameterNamesModel `tfschema:"subscription_key_parameter_names"`
	SubscriptionEnabled           bool                                 `tfschema:"subscription_enabled"`
	TermsOfServiceUrl             string                               `tfschema:"terms_of_service_url"`
	SourceApiId                   string                               `tfschema:"source_api_id"`
	OAuth2Authorization           []OAuth2AuthorizationModel           `tfschema:"oauth2_authorization"`
	OpenidAuthentication          []OpenidAuthenticationModel          `tfschema:"openid_authentication"`
	Version                       string                               `tfschema:"version"`
	VersionDescription            string                               `tfschema:"version_description"`
	VersionSetId                  string                               `tfschema:"version_set_id"`
	// Computed
	IsCurrent bool `tfschema:"is_current"`
	IsOnline  bool `tfschema:"is_online"`
}

type OAuth2AuthorizationModel struct {
	AuthorizationServerName string `tfschema:"authorization_server_name"`
	Scope                   string `tfschema:"scope"`
}

type OpenidAuthenticationModel struct {
	OpenidProviderName        string   `tfschema:"openid_provider_name"`
	BearerTokenSendingMethods []string `tfschema:"bearer_token_sending_methods"`
}

type ContactModel struct {
	Email string `tfschema:"email"`
	Name  string `tfschema:"name"`
	Url   string `tfschema:"url"`
}

type ImportModel struct {
	ContentFormat string              `tfschema:"content_format"`
	ContentValue  string              `tfschema:"content_value"`
	WsdlSelector  []WsdlSelectorModel `tfschema:"wsdl_selector"`
}

type WsdlSelectorModel struct {
	ServiceName  string `tfschema:"service_name"`
	EndpointName string `tfschema:"endpoint_name"`
}

type LicenseModel struct {
	Name string `tfschema:"name"`
	Url  string `tfschema:"url"`
}

type SubscriptionKeyParameterNamesModel struct {
	Header string `tfschema:"header"`
	Query  string `tfschema:"query"`
}

func (r ApiManagementWorkspaceApiResource) ResourceType() string {
	return "azurerm_api_management_workspace_api"
}

func (r ApiManagementWorkspaceApiResource) ModelObject() interface{} {
	return &ApiManagementWorkspaceApiModel{}
}

func (r ApiManagementWorkspaceApiResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return api.ValidateWorkspaceApiID
}

func (r ApiManagementWorkspaceApiResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ApiManagementApiName,
		},

		"api_management_workspace_id": commonschema.ResourceIDReferenceRequiredForceNew(&workspace.WorkspaceId{}),

		"revision": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"api_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringInSlice(api.PossibleValuesForApiType(), false),
		},

		"contact": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"email": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validate.EmailAddress,
					},
					"name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"url": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},
				},
			},
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"import": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"content_value": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"content_format": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(api.PossibleValuesForContentFormat(), false),
					},

					"wsdl_selector": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"service_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"endpoint_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
				},
			},
		},

		"license": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"url": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},
				},
			},
		},

		"oauth2_authorization": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"openid_authentication"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"authorization_server_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.ApiManagementChildName,
					},
					"scope": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						// There is currently no validation, as any length and characters can be used in the field
					},
				},
			},
		},

		"openid_authentication": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"oauth2_authorization"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"openid_provider_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.ApiManagementChildName,
					},
					"bearer_token_sending_methods": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(api.PossibleValuesForBearerTokenSendingMethods(), false),
						},
					},
				},
			},
		},

		"path": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validate.ApiManagementApiPath,
		},

		"protocols": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice(api.PossibleValuesForProtocol(), false),
			},
		},

		"revision_description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"service_url": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"source_api_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: api.ValidateWorkspaceApiID,
		},

		"subscription_key_parameter_names": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"header": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"query": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"subscription_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"terms_of_service_url": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},

		"version": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			RequiredWith: []string{"version_set_id"},
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"version_description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"version_set_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			RequiredWith: []string{"version"},
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r ApiManagementWorkspaceApiResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"is_current": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"is_online": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
	}
}

func (r ApiManagementWorkspaceApiResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ApiManagementWorkspaceApiModel
			if err := metadata.DecodeDiff(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			values := metadata.ResourceDiff.GetRawConfig().AsValueMap()
			protocols := expandApiManagementWorkspaceApiProtocols(model.Protocols)
			if values["source_api_id"].IsNull() && (values["display_name"].IsNull() || protocols == nil) {
				return errors.New("`display_name`, `protocols` are required when `source_api_id` is not set")
			}

			if model.ApiType == string(api.ApiTypeWebsocket) && model.ServiceUrl == "" {
				return errors.New("`service_url` is required when `api_type` is `websocket`")
			}

			return nil
		},
	}
}

func (r ApiManagementWorkspaceApiResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.V20240501ApiClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ApiManagementWorkspaceApiModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			workspaceId, err := api.ParseWorkspaceID(model.ApiManagementWorkspaceId)
			if err != nil {
				return err
			}

			id := api.NewWorkspaceApiID(subscriptionId, workspaceId.ResourceGroupName, workspaceId.ServiceName, workspaceId.WorkspaceId, model.Name)

			metadata.Logger.Infof("Import check for %s", id)
			existing, err := client.WorkspaceApiGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			metadata.Logger.Infof("Creating %s", id)

			apiType := api.ApiTypeHTTP
			if model.ApiType != "" {
				apiType = api.ApiType(model.ApiType)
			}
			soapApiType := workspaceApiSoapApiTypeFromApiType(apiType)

			// If import is used, we need to send properties to Azure API in two operations.
			// First we execute import and then updated the other props.
			if len(model.Import) > 0 {
				if apiParams := expandApiManagementWorkspaceApiImport(model.Import, apiType, soapApiType, model.Path, model.ServiceUrl, model.Version, model.VersionSetId); apiParams != nil {
					result, err := client.WorkspaceApiCreateOrUpdate(ctx, id, *apiParams, api.WorkspaceApiCreateOrUpdateOperationOptions{})
					if err != nil {
						return fmt.Errorf("creating with import of %s: %+v", id, err)
					}

					if pollerType := custompollers.NewAPIManagementWorkspaceAPIPoller(client, id, result.HttpResponse); pollerType != nil {
						poller := pollers.NewPoller(pollerType, 5*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
						if err := poller.PollUntilDone(ctx); err != nil {
							return fmt.Errorf("polling import %s: %+v", id, err)
						}
					}
				}
			}

			parameters := api.ApiCreateOrUpdateParameter{
				Properties: &api.ApiCreateOrUpdateProperties{
					SubscriptionRequired: pointer.To(model.SubscriptionEnabled),
				},
			}

			if model.DisplayName != "" {
				parameters.Properties.DisplayName = pointer.To(model.DisplayName)
			}

			if model.Description != "" {
				parameters.Properties.Description = pointer.To(model.Description)
			}

			if model.RevisionDescription != "" {
				parameters.Properties.ApiRevisionDescription = pointer.To(model.RevisionDescription)
			}

			if model.VersionDescription != "" {
				parameters.Properties.ApiVersionDescription = pointer.To(model.VersionDescription)
			}

			if model.Path != "" {
				parameters.Properties.Path = model.Path
			}

			if model.ServiceUrl != "" {
				parameters.Properties.ServiceURL = pointer.To(model.ServiceUrl)
			}

			if model.TermsOfServiceUrl != "" {
				parameters.Properties.TermsOfServiceURL = pointer.To(model.TermsOfServiceUrl)
			}

			if model.SourceApiId != "" {
				parameters.Properties.SourceApiId = pointer.To(model.SourceApiId)
			}

			if len(model.Protocols) > 0 {
				parameters.Properties.Protocols = expandApiManagementWorkspaceApiProtocols(model.Protocols)
			}

			if model.ApiType != "" {
				parameters.Properties.Type = pointer.To(api.ApiType(model.ApiType))
				parameters.Properties.ApiType = pointer.To(soapApiType)
			}

			if model.Version != "" {
				parameters.Properties.ApiVersion = pointer.To(model.Version)
			}

			if model.VersionSetId != "" {
				parameters.Properties.ApiVersionSetId = pointer.To(model.VersionSetId)
			}

			if len(model.OAuth2Authorization) > 0 {
				parameters.Properties.AuthenticationSettings = &api.AuthenticationSettingsContract{
					OAuth2: expandApiManagementWorkspaceApiOAuth2AuthenticationSettingsContract(model.OAuth2Authorization),
				}
			}

			if len(model.OpenidAuthentication) > 0 {
				if parameters.Properties.AuthenticationSettings == nil {
					parameters.Properties.AuthenticationSettings = &api.AuthenticationSettingsContract{}
				}
				parameters.Properties.AuthenticationSettings.Openid = expandApiManagementWorkspaceApiOpenIDAuthenticationSettingsContract(model.OpenidAuthentication)
			}

			if len(model.Contact) > 0 {
				parameters.Properties.Contact = expandApiManagementWorkspaceApiContact(model.Contact)
			}

			if len(model.License) > 0 {
				parameters.Properties.License = expandApiManagementWorkspaceApiLicense(model.License)
			}

			if len(model.SubscriptionKeyParameterNames) > 0 {
				parameters.Properties.SubscriptionKeyParameterNames = expandApiManagementWorkspaceApiSubscriptionKeyParameterNames(model.SubscriptionKeyParameterNames)
			}

			result, err := client.WorkspaceApiCreateOrUpdate(ctx, id, parameters, api.WorkspaceApiCreateOrUpdateOperationOptions{})
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			if pollerType := custompollers.NewAPIManagementWorkspaceAPIPoller(client, id, result.HttpResponse); pollerType != nil {
				poller := pollers.NewPoller(pollerType, 5*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
				if err := poller.PollUntilDone(ctx); err != nil {
					return fmt.Errorf("polling creating %s: %+v", id, err)
				}
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ApiManagementWorkspaceApiResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.V20240501ApiClient

			id, err := api.ParseWorkspaceApiID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Decoding state for %s", *id)
			var state ApiManagementWorkspaceApiModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			metadata.Logger.Infof("Updating %s", *id)
			apiType := api.ApiTypeHTTP
			if state.ApiType != "" {
				apiType = api.ApiType(state.ApiType)
			}
			soapApiType := workspaceApiSoapApiTypeFromApiType(apiType)

			// If import is used, we need to send properties to Azure API in two operations.
			// First we execute import and then updated the other props.
			if metadata.ResourceData.HasChange("import") {
				if len(state.Import) > 0 {
					metadata.ResourceData.Partial(true)
					if apiParams := expandApiManagementWorkspaceApiImport(state.Import, apiType, soapApiType,
						state.Path, state.ServiceUrl, state.Version, state.VersionSetId); apiParams != nil {
						result, err := client.WorkspaceApiCreateOrUpdate(ctx, *id, *apiParams, api.WorkspaceApiCreateOrUpdateOperationOptions{})
						if err != nil {
							return fmt.Errorf("updating with import of %s: %+v", id, err)
						}

						if pollerType := custompollers.NewAPIManagementWorkspaceAPIPoller(client, *id, result.HttpResponse); pollerType != nil {
							poller := pollers.NewPoller(pollerType, 5*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
							if err := poller.PollUntilDone(ctx); err != nil {
								return fmt.Errorf("polling import %s: %+v", *id, err)
							}
						}
					}
					metadata.ResourceData.Partial(false)
				}
			}

			existing, err := client.WorkspaceApiGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			payload := *existing.Model
			if payload.Properties.Type != nil {
				soapApiType = workspaceApiSoapApiTypeFromApiType(pointer.From(payload.Properties.Type))
			}

			prop := &api.ApiCreateOrUpdateProperties{
				Path:                          payload.Properties.Path,
				Protocols:                     payload.Properties.Protocols,
				ServiceURL:                    payload.Properties.ServiceURL,
				Description:                   payload.Properties.Description,
				ApiVersionDescription:         payload.Properties.ApiVersionDescription,
				ApiRevisionDescription:        payload.Properties.ApiRevisionDescription,
				SubscriptionRequired:          payload.Properties.SubscriptionRequired,
				SubscriptionKeyParameterNames: payload.Properties.SubscriptionKeyParameterNames,
				Contact:                       payload.Properties.Contact,
				License:                       payload.Properties.License,
				SourceApiId:                   payload.Properties.SourceApiId,
				DisplayName:                   payload.Properties.DisplayName,
				ApiVersion:                    payload.Properties.ApiVersion,
				ApiVersionSetId:               payload.Properties.ApiVersionSetId,
				TermsOfServiceURL:             payload.Properties.TermsOfServiceURL,
				Type:                          payload.Properties.Type,
				ApiType:                       pointer.To(soapApiType),
			}

			// When the `import` property changes, values in the import file (e.g., `info.title`, `info.description`)
			// may overlap with explicitly configured API properties like `display_name` or `description`.
			// For example, if `display_name = "My first API"` remains unchanged in the config, but the import sets
			// `info.title = "My test API"`, the imported value will override the configured one.
			// To avoid this, always set explicitly configured properties during updates—even if unchanged—
			// so they take precedence over imported values.
			if metadata.ResourceData.HasChange("import") {
				if state.Path != "" {
					prop.Path = state.Path
				}

				if len(state.Protocols) > 0 {
					prop.Protocols = expandApiManagementWorkspaceApiProtocols(state.Protocols)
				}

				if state.ServiceUrl != "" {
					prop.ServiceURL = pointer.To(state.ServiceUrl)
				}

				if state.Description != "" {
					prop.Description = pointer.To(state.Description)
				}

				if state.VersionDescription != "" {
					prop.ApiVersionDescription = pointer.To(state.VersionDescription)
				}

				if state.RevisionDescription != "" {
					prop.ApiRevisionDescription = pointer.To(state.RevisionDescription)
				}

				if len(state.SubscriptionKeyParameterNames) > 0 {
					prop.SubscriptionKeyParameterNames = expandApiManagementWorkspaceApiSubscriptionKeyParameterNames(state.SubscriptionKeyParameterNames)
				}

				if len(state.Contact) > 0 {
					prop.Contact = expandApiManagementWorkspaceApiContact(state.Contact)
				}

				if len(state.License) > 0 {
					prop.License = expandApiManagementWorkspaceApiLicense(state.License)
				}

				if state.SourceApiId != "" {
					prop.SourceApiId = pointer.To(state.SourceApiId)
				}

				if state.DisplayName != "" {
					prop.DisplayName = pointer.To(state.DisplayName)
				}

				if state.Version != "" {
					prop.ApiVersion = pointer.To(state.Version)
				}

				if state.VersionSetId != "" {
					prop.ApiVersionSetId = pointer.To(state.VersionSetId)
				}

				if state.TermsOfServiceUrl != "" {
					prop.TermsOfServiceURL = pointer.To(state.TermsOfServiceUrl)
				}

				if state.ApiType != "" {
					prop.Type = pointer.To(api.ApiType(state.ApiType))
					prop.ApiType = pointer.To(workspaceApiSoapApiTypeFromApiType(api.ApiType(state.ApiType)))
				}
			}

			// For the setting of `AuthenticationSettingsContract`, the PUT payload restrictions are as follows:
			//   1. Cannot have both 'oAuth2' and 'openid' set
			//   2. Cannot use `OAuth2AuthenticationSettings` in combination with `OAuth2` nor `openid`
			//   3. Cannot use `OpenidAuthenticationSettings` in combination with `Openid` nor `OAuth2`
			// If specifying `oauth2_authorization`/`openid_authentication` when creating a resource and then updating the resource, the error #2/#3 mentioned above will occur.
			// This is because starting from the 2022-08-01 version, the Get API additionally returns a collection of `oauth2_authorization`/`openid_authentication` authentication settings, which property name is `OAuth2AuthenticationSettings`/`OpenidAuthenticationSetting`.
			// Given the API behavior, the update here should only read the specified property `oauth2_authorization`/`openid_authentication` to exclude `OAuth2AuthenticationSettings`/`OpenidAuthenticationSetting` to ensure the update works properly.
			if v := payload.Properties.AuthenticationSettings; v != nil {
				authenticationSettings := &api.AuthenticationSettingsContract{}
				if v.OAuth2 != nil {
					authenticationSettings.OAuth2 = v.OAuth2
					prop.AuthenticationSettings = authenticationSettings
				}

				if v.Openid != nil {
					authenticationSettings.Openid = v.Openid
					prop.AuthenticationSettings = authenticationSettings
				}
			}

			if metadata.ResourceData.HasChange("path") {
				prop.Path = state.Path
			}

			if metadata.ResourceData.HasChange("protocols") {
				prop.Protocols = expandApiManagementWorkspaceApiProtocols(state.Protocols)
			}

			if metadata.ResourceData.HasChange("api_type") {
				prop.Type = pointer.To(apiType)
				prop.ApiType = pointer.To(soapApiType)
			}

			if metadata.ResourceData.HasChange("service_url") {
				prop.ServiceURL = pointer.To(state.ServiceUrl)
			}

			if metadata.ResourceData.HasChange("description") {
				prop.Description = pointer.To(state.Description)
			}

			if metadata.ResourceData.HasChange("revision_description") {
				prop.ApiRevisionDescription = pointer.To(state.RevisionDescription)
			}

			if metadata.ResourceData.HasChange("version_description") {
				prop.ApiVersionDescription = pointer.To(state.VersionDescription)
			}

			if metadata.ResourceData.HasChange("subscription_enabled") {
				prop.SubscriptionRequired = pointer.To(state.SubscriptionEnabled)
			}

			if metadata.ResourceData.HasChange("subscription_key_parameter_names") {
				prop.SubscriptionKeyParameterNames = expandApiManagementWorkspaceApiSubscriptionKeyParameterNames(state.SubscriptionKeyParameterNames)
			}

			if metadata.ResourceData.HasChange("oauth2_authorization") {
				if prop.AuthenticationSettings == nil {
					prop.AuthenticationSettings = &api.AuthenticationSettingsContract{}
				}
				prop.AuthenticationSettings.OAuth2 = expandApiManagementWorkspaceApiOAuth2AuthenticationSettingsContract(state.OAuth2Authorization)
			}

			if metadata.ResourceData.HasChange("openid_authentication") {
				if prop.AuthenticationSettings == nil {
					prop.AuthenticationSettings = &api.AuthenticationSettingsContract{}
				}
				prop.AuthenticationSettings.Openid = expandApiManagementWorkspaceApiOpenIDAuthenticationSettingsContract(state.OpenidAuthentication)
			}

			if metadata.ResourceData.HasChange("contact") {
				prop.Contact = expandApiManagementWorkspaceApiContact(state.Contact)
			}

			if metadata.ResourceData.HasChange("license") {
				prop.License = expandApiManagementWorkspaceApiLicense(state.License)
			}

			if metadata.ResourceData.HasChange("source_api_id") {
				prop.SourceApiId = pointer.To(state.SourceApiId)
			}

			if metadata.ResourceData.HasChange("display_name") {
				prop.DisplayName = pointer.To(state.DisplayName)
			}

			if metadata.ResourceData.HasChange("version") {
				prop.ApiVersion = pointer.To(state.Version)
			}

			if metadata.ResourceData.HasChange("version_set_id") {
				prop.ApiVersionSetId = pointer.To(state.VersionSetId)
			}

			if metadata.ResourceData.HasChange("terms_of_service_url") {
				prop.TermsOfServiceURL = pointer.To(state.TermsOfServiceUrl)
			}

			params := api.ApiCreateOrUpdateParameter{
				Properties: prop,
			}

			result, err := client.WorkspaceApiCreateOrUpdate(ctx, *id, params, api.WorkspaceApiCreateOrUpdateOperationOptions{})
			if err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			if pollerType := custompollers.NewAPIManagementWorkspaceAPIPoller(client, *id, result.HttpResponse); pollerType != nil {
				poller := pollers.NewPoller(pollerType, 5*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
				if err := poller.PollUntilDone(ctx); err != nil {
					return fmt.Errorf("polling updating %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}

func (r ApiManagementWorkspaceApiResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.V20240501ApiClient

			id, err := api.ParseWorkspaceApiID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Reading %s", id)
			resp, err := client.WorkspaceApiGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var config ApiManagementWorkspaceApiModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			state := ApiManagementWorkspaceApiModel{
				Name:                     id.ApiId,
				ApiManagementWorkspaceId: api.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceId).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					apiType := string(pointer.From(props.Type))
					if apiType == "" {
						apiType = string(api.ApiTypeHTTP)
					}
					state.ApiType = apiType
					state.DisplayName = pointer.From(props.DisplayName)
					state.Description = pointer.From(props.Description)
					state.IsCurrent = pointer.From(props.IsCurrent)
					state.IsOnline = pointer.From(props.IsOnline)
					state.Path = props.Path
					state.ServiceUrl = pointer.From(props.ServiceURL)
					state.Revision = pointer.From(props.ApiRevision)
					state.SubscriptionEnabled = pointer.From(props.SubscriptionRequired)
					state.Version = pointer.From(props.ApiVersion)
					state.VersionSetId = pointer.From(props.ApiVersionSetId)
					state.RevisionDescription = pointer.From(props.ApiRevisionDescription)
					state.VersionDescription = pointer.From(props.ApiVersionDescription)
					state.TermsOfServiceUrl = pointer.From(props.TermsOfServiceURL)

					if props.Protocols != nil {
						protocols := make([]string, len(*props.Protocols))
						for i, protocol := range *props.Protocols {
							protocols[i] = string(protocol)
						}
						state.Protocols = flattenApiManagementWorkspaceApiProtocols(props.Protocols)
					}

					state.SubscriptionKeyParameterNames = flattenSubscriptionKeyParameterNames(props.SubscriptionKeyParameterNames)

					if authenticationSettings := props.AuthenticationSettings; authenticationSettings != nil {
						state.OAuth2Authorization = flattenApiManagementWorkspaceApiOAuth2Authorization(authenticationSettings.OAuth2)
						state.OpenidAuthentication = flattenApiManagementWorkspaceApiOpenIDAuthentication(authenticationSettings.Openid)
					}

					if props.Contact != nil {
						state.Contact = flattenContact(props.Contact)
					}

					if props.License != nil {
						state.License = flattenLicense(props.License)
					}

					state.Import = config.Import
					state.SourceApiId = config.SourceApiId
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ApiManagementWorkspaceApiResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.V20240501ApiClient

			id, err := api.ParseWorkspaceApiID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Deleting %s", id)

			if resp, err := client.WorkspaceApiDelete(ctx, *id, api.DefaultWorkspaceApiDeleteOperationOptions()); err != nil {
				if resp.HttpResponse.StatusCode == http.StatusAccepted {
					pollerType, err := custompollers.NewAPIManagementWorkspaceApiDeletePoller(client, resp.HttpResponse)
					if err != nil {
						return fmt.Errorf("polling deleting %s: %+v", id, err)
					}
					if pollerType != nil {
						poller := pollers.NewPoller(pollerType, 20*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
						if err := poller.PollUntilDone(ctx); err != nil {
							return fmt.Errorf("polling purging the deleting %s: %+v", id, err)
						}
					}
				} else {
					if !response.WasNotFound(resp.HttpResponse) {
						return fmt.Errorf("deleting %s: %+v", *id, err)
					}
				}
			}

			return nil
		},
	}
}

func expandApiManagementWorkspaceApiContact(input []ContactModel) *api.ApiContactInformation {
	if len(input) == 0 {
		return nil
	}

	return &api.ApiContactInformation{
		Email: pointer.To(input[0].Email),
		Name:  pointer.To(input[0].Name),
		Url:   pointer.To(input[0].Url),
	}
}

func expandApiManagementWorkspaceApiLicense(input []LicenseModel) *api.ApiLicenseInformation {
	if len(input) == 0 {
		return nil
	}

	return &api.ApiLicenseInformation{
		Name: pointer.To(input[0].Name),
		Url:  pointer.To(input[0].Url),
	}
}

func expandApiManagementWorkspaceApiSubscriptionKeyParameterNames(input []SubscriptionKeyParameterNamesModel) *api.SubscriptionKeyParameterNamesContract {
	if len(input) == 0 {
		return nil
	}

	return &api.SubscriptionKeyParameterNamesContract{
		Header: pointer.To(input[0].Header),
		Query:  pointer.To(input[0].Query),
	}
}

func workspaceApiSoapApiTypeFromApiType(apiType api.ApiType) api.SoapApiType {
	return map[api.ApiType]api.SoapApiType{
		api.ApiTypeGraphql:   api.SoapApiTypeGraphql,
		api.ApiTypeHTTP:      api.SoapApiTypeHTTP,
		api.ApiTypeSoap:      api.SoapApiTypeSoap,
		api.ApiTypeWebsocket: api.SoapApiTypeWebsocket,
		api.ApiTypeOdata:     api.SoapApiTypeOdata,
		api.ApiTypeGrpc:      api.SoapApiTypeGrpc,
	}[apiType]
}

func expandApiManagementWorkspaceApiImport(importVs []ImportModel, apiType api.ApiType, soapApiType api.SoapApiType, path, serviceUrl, version, versionSetId string) *api.ApiCreateOrUpdateParameter {
	if len(importVs) == 0 {
		return nil
	}

	importV := importVs[0]
	apiParams := api.ApiCreateOrUpdateParameter{
		Properties: &api.ApiCreateOrUpdateProperties{
			Format: pointer.To(api.ContentFormat(importV.ContentFormat)),
			Value:  pointer.To(importV.ContentValue),
		},
	}

	if apiType != "" {
		apiParams.Properties.Type = pointer.To(apiType)
		apiParams.Properties.ApiType = pointer.To(soapApiType)
	}

	if path != "" {
		apiParams.Properties.Path = path
	}

	if wsdlSelectorVs := importV.WsdlSelector; len(wsdlSelectorVs) > 0 {
		apiParams.Properties.WsdlSelector = &api.ApiCreateOrUpdatePropertiesWsdlSelector{
			WsdlServiceName:  pointer.To(wsdlSelectorVs[0].ServiceName),
			WsdlEndpointName: pointer.To(wsdlSelectorVs[0].EndpointName),
		}
	}

	if serviceUrl != "" {
		apiParams.Properties.ServiceURL = pointer.To(serviceUrl)
	}

	if version != "" {
		apiParams.Properties.ApiVersion = pointer.To(version)
	}

	if versionSetId != "" {
		apiParams.Properties.ApiVersionSetId = pointer.To(versionSetId)
	}

	return &apiParams
}

func expandApiManagementWorkspaceApiOAuth2AuthenticationSettingsContract(input []OAuth2AuthorizationModel) *api.OAuth2AuthenticationSettingsContract {
	if len(input) == 0 {
		return nil
	}

	return &api.OAuth2AuthenticationSettingsContract{
		AuthorizationServerId: pointer.To(input[0].AuthorizationServerName),
		Scope:                 pointer.To(input[0].Scope),
	}
}

func expandApiManagementWorkspaceApiOpenIDAuthenticationSettingsContract(input []OpenidAuthenticationModel) *api.OpenIdAuthenticationSettingsContract {
	if len(input) == 0 {
		return nil
	}

	return &api.OpenIdAuthenticationSettingsContract{
		OpenidProviderId:          pointer.To(input[0].OpenidProviderName),
		BearerTokenSendingMethods: expandApiManagementWorkspaceOpenIDAuthenticationSettingsBearerTokenSendingMethods(input[0].BearerTokenSendingMethods),
	}
}

func expandApiManagementWorkspaceOpenIDAuthenticationSettingsBearerTokenSendingMethods(input []string) *[]api.BearerTokenSendingMethods {
	if input == nil {
		return nil
	}

	results := make([]api.BearerTokenSendingMethods, 0)
	for _, v := range input {
		results = append(results, api.BearerTokenSendingMethods(v))
	}

	return &results
}

func expandApiManagementWorkspaceApiProtocols(input []string) *[]api.Protocol {
	if len(input) == 0 {
		return nil
	}
	results := make([]api.Protocol, 0)

	for _, v := range input {
		results = append(results, api.Protocol(v))
	}

	return &results
}

func flattenContact(input *api.ApiContactInformation) []ContactModel {
	output := make([]ContactModel, 0)
	if input == nil {
		return output
	}

	return append(output, ContactModel{
		Email: pointer.From(input.Email),
		Name:  pointer.From(input.Name),
		Url:   pointer.From(input.Url),
	})
}

func flattenLicense(input *api.ApiLicenseInformation) []LicenseModel {
	output := make([]LicenseModel, 0)
	if input == nil {
		return output
	}

	return append(output, LicenseModel{
		Name: pointer.From(input.Name),
		Url:  pointer.From(input.Url),
	})
}

func flattenSubscriptionKeyParameterNames(input *api.SubscriptionKeyParameterNamesContract) []SubscriptionKeyParameterNamesModel {
	outputList := make([]SubscriptionKeyParameterNamesModel, 0)
	if input == nil {
		return outputList
	}

	return append(outputList, SubscriptionKeyParameterNamesModel{
		Header: pointer.From(input.Header),
		Query:  pointer.From(input.Query),
	})
}

func flattenApiManagementWorkspaceApiProtocols(input *[]api.Protocol) []string {
	outputList := make([]string, 0)
	if input == nil {
		return outputList
	}

	for _, v := range *input {
		outputList = append(outputList, string(v))
	}

	return outputList
}

func flattenApiManagementWorkspaceApiOAuth2Authorization(input *api.OAuth2AuthenticationSettingsContract) []OAuth2AuthorizationModel {
	outputList := make([]OAuth2AuthorizationModel, 0)
	if input == nil {
		return outputList
	}

	return append(outputList, OAuth2AuthorizationModel{
		AuthorizationServerName: pointer.From(input.AuthorizationServerId),
		Scope:                   pointer.From(input.Scope),
	})
}

func flattenApiManagementWorkspaceApiOpenIDAuthentication(input *api.OpenIdAuthenticationSettingsContract) []OpenidAuthenticationModel {
	outputList := make([]OpenidAuthenticationModel, 0)
	if input == nil {
		return outputList
	}

	return append(outputList, OpenidAuthenticationModel{
		OpenidProviderName:        pointer.From(input.OpenidProviderId),
		BearerTokenSendingMethods: flattenApiManagementWorkspaceOpenIDAuthenticationSettingsBearerTokenSendingMethods(input.BearerTokenSendingMethods),
	})
}

func flattenApiManagementWorkspaceOpenIDAuthenticationSettingsBearerTokenSendingMethods(input *[]api.BearerTokenSendingMethods) []string {
	outputList := make([]string, 0)
	if input == nil {
		return outputList
	}

	bearerTokenSendingMethods := make([]string, 0)
	if s := input; s != nil {
		for _, v := range *s {
			bearerTokenSendingMethods = append(bearerTokenSendingMethods, string(v))
		}
	}

	return bearerTokenSendingMethods
}
