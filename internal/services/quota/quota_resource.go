package quota

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/quota/2025-03-01/quotainformation"
	"github.com/hashicorp/go-azure-sdk/resource-manager/quota/2025-03-01/usagesinformation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type QuotaModel struct {
	Name                 string `tfschema:"name"`
	ProviderNamespace    string `tfschema:"provider_namespace"`
	Location             string `tfschema:"location"`
	LimitValue           int64  `tfschema:"limit_value"`
	LimitType            string `tfschema:"limit_type"`
	LimitObjectType      string `tfschema:"limit_object_type"`
	ResourceType         string `tfschema:"resource_type"`
	AdditionalProperties string `tfschema:"additional_properties"`
}

type QuotaResource struct{}

var _ sdk.ResourceWithUpdate = QuotaResource{}

func (r QuotaResource) ResourceType() string {
	return "azurerm_quota"
}

func (r QuotaResource) ModelObject() interface{} {
	return &QuotaModel{}
}

func (r QuotaResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return quotainformation.ValidateScopedQuotaID
}

func (r QuotaResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"limit_object_type": {
			Type:         pluginsdk.TypeString,
			ForceNew:     true,
			Required:     true,
			ValidateFunc: validation.StringInSlice(quotainformation.PossibleValuesForLimitType(), false),
		},

		"limit_value": {
			Type:     pluginsdk.TypeInt,
			Required: true,
		},

		"provider_namespace": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"additional_properties": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsJSON,
		},

		"limit_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(quotainformation.QuotaLimitTypesIndependent),
			ValidateFunc: validation.StringInSlice(quotainformation.PossibleValuesForQuotaLimitTypes(), false),
		},

		"resource_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r QuotaResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r QuotaResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model QuotaModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Quota.QuotaClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			scope := fmt.Sprintf("subscriptions/%s/providers/%s/locations/%s", subscriptionId, model.ProviderNamespace, location.Normalize(model.Location))
			id := quotainformation.NewScopedQuotaID(scope, model.Name)
			existing, err := client.QuotaGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}

			properties := &quotainformation.CurrentQuotaLimitBase{
				Properties: &quotainformation.QuotaProperties{
					Limit: quotainformation.LimitObject{
						LimitType:       pointer.To(quotainformation.QuotaLimitTypes(model.LimitType)),
						Value:           model.LimitValue,
						LimitObjectType: quotainformation.LimitType(model.LimitObjectType),
					},
					Name: &quotainformation.ResourceName{
						Value: pointer.To(model.Name),
					},
				},
			}

			if model.ResourceType != "" {
				properties.Properties.ResourceType = &model.ResourceType
			}

			if model.AdditionalProperties != "" {
				var result interface{}
				err := json.Unmarshal([]byte(model.AdditionalProperties), &result)
				if err != nil {
					return fmt.Errorf("unmarshaling `additional_properties`: %+v", err)
				}
				properties.Properties.Properties = pointer.To(result)
			}

			if err := client.QuotaCreateOrUpdateThenPoll(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r QuotaResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Quota.QuotaClient

			id, err := quotainformation.ParseScopedQuotaID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model QuotaModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.QuotaGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			payload := existing.Model
			if metadata.ResourceData.HasChange("limit_value") {
				limitObjectPtr, ok := payload.Properties.Limit.(*quotainformation.LimitObject)
				if ok {
					limitObjectPtr.Value = model.LimitValue
				}
			}

			if err := client.QuotaCreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r QuotaResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Quota.QuotaClient

			id, err := quotainformation.ParseScopedQuotaID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.QuotaGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			privider := ""
			re := regexp.MustCompile(`/providers/([^/]+)`)
			match := re.FindStringSubmatch(id.Scope)
			if len(match) > 1 {
				privider = match[1]
			}

			region := ""
			re = regexp.MustCompile(`/locations/([^/]+)`)
			match = re.FindStringSubmatch(id.Scope)
			if len(match) > 1 {
				region = match[1]
			}

			state := QuotaModel{
				Name:              id.QuotaName,
				ProviderNamespace: privider,
				Location:          location.Normalize(region),
			}

			if model := resp.Model; model != nil {
				if properties := model.Properties; properties != nil {
					limitObject := properties.Limit.(quotainformation.LimitObject)
					state.LimitValue = limitObject.Value
					state.LimitType = string(pointer.From(limitObject.LimitType))
					state.LimitObjectType = string(limitObject.LimitObjectType)
					state.ResourceType = pointer.From(properties.ResourceType)
					raw, err := json.Marshal(properties.Properties)
					if err != nil {
						return fmt.Errorf("could not marshal `additional_properties`: %+v", err)
					}
					state.AdditionalProperties = string(raw)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r QuotaResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Quota.QuotaClient
			usagesClient := metadata.Client.Quota.UsagesClient

			id, err := quotainformation.ParseScopedQuotaID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			quotaResp, err := client.QuotaGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(quotaResp.HttpResponse) {
					return fmt.Errorf("retrieving %s: %+v", *id, err)
				}
			}

			if quotaResp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			payload := quotaResp.Model
			if payload.Properties != nil {
				scopedUsageId := usagesinformation.NewScopedUsageID(id.Scope, id.QuotaName)
				usageResp, err := usagesClient.UsagesGet(ctx, scopedUsageId)
				if err != nil {
					if response.WasNotFound(usageResp.HttpResponse) {
						return fmt.Errorf("retrieving %s: %+v", scopedUsageId, err)
					}
				}

				used := int64(0)
				if model := usageResp.Model; model != nil {
					if properties := model.Properties; properties != nil {
						if usages := properties.Usages; usages != nil {
							used = usages.Value
						}
					}
				}

				// Deleting the quota means that the value of limit is equal to the used value
				limitObjectPtr, ok := payload.Properties.Limit.(*quotainformation.LimitObject)
				if ok {
					limitObjectPtr.Value = used
				}

				if err := client.QuotaCreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
					return fmt.Errorf("deleting %s: %+v", id, err)
				}
			}
			return nil
		},
	}
}
