package cognitive

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2022-10-01/cognitiveservicesaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2022-10-01/commitmentplans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type cognitiveCommitmentPlanModel struct {
	Name               string                       `tfschema:"name"`
	CognitiveAccountId string                       `tfschema:"cognitive_account_id"`
	AutoRenew          bool                         `tfschema:"auto_renew"`
	Current            []CommitmentPeriodModel      `tfschema:"current"`
	HostingModel       commitmentplans.HostingModel `tfschema:"hosting_model"`
	Next               []CommitmentPeriodModel      `tfschema:"next"`
	PlanType           string                       `tfschema:"plan_type"`
}

type CommitmentPeriodModel struct {
	Count int64  `tfschema:"count"`
	Tier  string `tfschema:"tier"`
}

type cognitiveCommitmentPlanResource struct{}

var _ sdk.ResourceWithUpdate = cognitiveCommitmentPlanResource{}

func (r cognitiveCommitmentPlanResource) ResourceType() string {
	return "azurerm_cognitive_commitment_plan"
}

func (r cognitiveCommitmentPlanResource) ModelObject() interface{} {
	return &cognitiveCommitmentPlanModel{}
}

func (r cognitiveCommitmentPlanResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commitmentplans.ValidateCommitmentPlanID
}

func (r cognitiveCommitmentPlanResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"cognitive_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: cognitiveservicesaccounts.ValidateAccountID,
		},

		"auto_renew": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"current": CommitmentPeriod(),

		"hosting_model": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(commitmentplans.HostingModelWeb),
				string(commitmentplans.HostingModelConnectedContainer),
				string(commitmentplans.HostingModelDisconnectedContainer),
			}, false),
		},

		"next": CommitmentPeriod(),

		"plan_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r cognitiveCommitmentPlanResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func CommitmentPeriod() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"count": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
				},

				"tier": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func (r cognitiveCommitmentPlanResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model cognitiveCommitmentPlanModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Cognitive.CommitmentPlansClient
			accountId, err := cognitiveservicesaccounts.ParseAccountID(model.CognitiveAccountId)
			if err != nil {
				return err
			}

			id := commitmentplans.NewCommitmentPlanID(accountId.SubscriptionId, accountId.ResourceGroupName, accountId.AccountName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &commitmentplans.CommitmentPlan{
				Properties: &commitmentplans.CommitmentPlanProperties{
					AutoRenew:    &model.AutoRenew,
					HostingModel: &model.HostingModel,
				},
			}

			currentValue, err := expandCommitmentPeriodModel(model.Current)
			if err != nil {
				return err
			}

			properties.Properties.Current = currentValue

			nextValue, err := expandCommitmentPeriodModel(model.Next)
			if err != nil {
				return err
			}

			properties.Properties.Next = nextValue

			if model.PlanType != "" {
				properties.Properties.PlanType = &model.PlanType
			}

			if _, err := client.CreateOrUpdate(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r cognitiveCommitmentPlanResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.CommitmentPlansClient

			id, err := commitmentplans.ParseCommitmentPlanID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model cognitiveCommitmentPlanModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("auto_renew") {
				properties.Properties.AutoRenew = &model.AutoRenew
			}

			if metadata.ResourceData.HasChange("current") {
				currentValue, err := expandCommitmentPeriodModel(model.Current)
				if err != nil {
					return err
				}

				properties.Properties.Current = currentValue
			}

			if metadata.ResourceData.HasChange("hosting_model") {
				properties.Properties.HostingModel = &model.HostingModel
			}

			if metadata.ResourceData.HasChange("next") {
				nextValue, err := expandCommitmentPeriodModel(model.Next)
				if err != nil {
					return err
				}

				properties.Properties.Next = nextValue
			}

			if metadata.ResourceData.HasChange("plan_type") {
				if model.PlanType != "" {
					properties.Properties.PlanType = &model.PlanType
				} else {
					properties.Properties.PlanType = nil
				}
			}

			properties.SystemData = nil

			if _, err := client.CreateOrUpdate(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r cognitiveCommitmentPlanResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.CommitmentPlansClient

			id, err := commitmentplans.ParseCommitmentPlanID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := cognitiveCommitmentPlanModel{
				Name:               id.CommitmentPlanName,
				CognitiveAccountId: cognitiveservicesaccounts.NewAccountID(id.SubscriptionId, id.ResourceGroupName, id.AccountName).ID(),
			}

			if properties := model.Properties; properties != nil {
				if properties.AutoRenew != nil {
					state.AutoRenew = *properties.AutoRenew
				}

				currentValue, err := flattenCommitmentPeriodModel(properties.Current)
				if err != nil {
					return err
				}

				state.Current = currentValue

				if properties.HostingModel != nil {
					state.HostingModel = *properties.HostingModel
				}

				nextValue, err := flattenCommitmentPeriodModel(properties.Next)
				if err != nil {
					return err
				}

				state.Next = nextValue

				if properties.PlanType != nil {
					state.PlanType = *properties.PlanType
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r cognitiveCommitmentPlanResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.CommitmentPlansClient

			id, err := commitmentplans.ParseCommitmentPlanID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandCommitmentPeriodModel(inputList []CommitmentPeriodModel) (*commitmentplans.CommitmentPeriod, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	input := &inputList[0]
	output := commitmentplans.CommitmentPeriod{
		Count: &input.Count,
	}

	if input.Tier != "" {
		output.Tier = &input.Tier
	}

	return &output, nil
}

func flattenCommitmentPeriodModel(input *commitmentplans.CommitmentPeriod) ([]CommitmentPeriodModel, error) {
	var outputList []CommitmentPeriodModel
	if input == nil {
		return outputList, nil
	}

	output := CommitmentPeriodModel{}

	if input.Count != nil {
		output.Count = *input.Count
	}

	if input.Tier != nil {
		output.Tier = *input.Tier
	}

	return append(outputList, output), nil
}
