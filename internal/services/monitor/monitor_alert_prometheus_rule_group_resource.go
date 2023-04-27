package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2023-03-01/prometheusrulegroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type AlertPrometheusRuleGroupResourceModel struct {
	Name              string                `tfschema:"name"`
	ResourceGroupName string                `tfschema:"resource_group_name"`
	ClusterName       string                `tfschema:"cluster_name"`
	Description       string                `tfschema:"description"`
	Enabled           bool                  `tfschema:"enabled"`
	Interval          string                `tfschema:"interval"`
	Location          string                `tfschema:"location"`
	Rules             []PrometheusRuleModel `tfschema:"rules"`
	Scopes            []string              `tfschema:"scopes"`
	Tags              map[string]string     `tfschema:"tags"`
}

type PrometheusRuleModel struct {
	Actions              []PrometheusRuleGroupActionModel          `tfschema:"actions"`
	Alert                string                                    `tfschema:"alert"`
	Annotations          map[string]string                         `tfschema:"annotations"`
	Enabled              bool                                      `tfschema:"enabled"`
	Expression           string                                    `tfschema:"expression"`
	For                  string                                    `tfschema:"for"`
	Labels               map[string]string                         `tfschema:"labels"`
	Record               string                                    `tfschema:"record"`
	ResolveConfiguration []PrometheusRuleResolveConfigurationModel `tfschema:"resolve_configuration"`
	Severity             int64                                     `tfschema:"severity"`
}

type PrometheusRuleGroupActionModel struct {
	ActionGroupId    string            `tfschema:"action_group_id"`
	ActionProperties map[string]string `tfschema:"action_properties"`
}

type PrometheusRuleResolveConfigurationModel struct {
	AutoResolved  bool   `tfschema:"auto_resolved"`
	TimeToResolve string `tfschema:"time_to_resolve"`
}

type AlertPrometheusRuleGroupResource struct{}

var _ sdk.ResourceWithUpdate = AlertPrometheusRuleGroupResource{}

func (r AlertPrometheusRuleGroupResource) ResourceType() string {
	return "azurerm_monitor_alert_prometheus_rule_group"
}

func (r AlertPrometheusRuleGroupResource) ModelObject() interface{} {
	return &AlertPrometheusRuleGroupResourceModel{}
}

func (r AlertPrometheusRuleGroupResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return prometheusrulegroups.ValidatePrometheusRuleGroupID
}

func (r AlertPrometheusRuleGroupResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"rules": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"actions": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"action_group_id": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"action_properties": {
									Type:     pluginsdk.TypeMap,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"alert": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"annotations": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"expression": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"for": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"labels": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: &schema.Schema{
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},

					"record": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"resolve_configuration": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"auto_resolved": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
								},

								"time_to_resolve": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"severity": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
				},
			},
		},

		"scopes": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"cluster_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"interval": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": commonschema.Tags(),
	}
}

func (r AlertPrometheusRuleGroupResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r AlertPrometheusRuleGroupResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model AlertPrometheusRuleGroupResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Monitor.AlertPrometheusRuleGroupClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := prometheusrulegroups.NewPrometheusRuleGroupID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := prometheusrulegroups.PrometheusRuleGroupResource{
				Location: location.Normalize(model.Location),
				Properties: prometheusrulegroups.PrometheusRuleGroupProperties{
					Enabled: &model.Enabled,
					Scopes:  model.Scopes,
				},
				Tags: &model.Tags,
			}

			if model.ClusterName != "" {
				properties.Properties.ClusterName = &model.ClusterName
			}

			if model.Description != "" {
				properties.Properties.Description = &model.Description
			}

			if model.Interval != "" {
				properties.Properties.Interval = &model.Interval
			}

			rulesValue, err := expandPrometheusRuleModel(model.Rules)
			if err != nil {
				return err
			}

			if rulesValue != nil {
				properties.Properties.Rules = *rulesValue
			}

			if _, err := client.CreateOrUpdate(ctx, id, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r AlertPrometheusRuleGroupResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.AlertPrometheusRuleGroupClient

			id, err := prometheusrulegroups.ParsePrometheusRuleGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model AlertPrometheusRuleGroupResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			if metadata.ResourceData.HasChange("cluster_name") {
				properties.Properties.ClusterName = pointer.To(model.ClusterName)
			}

			if metadata.ResourceData.HasChange("description") {
				properties.Properties.Description = pointer.To(model.Description)
			}

			if metadata.ResourceData.HasChange("enabled") {
				properties.Properties.Enabled = pointer.To(model.Enabled)
			}

			if metadata.ResourceData.HasChange("interval") {
				properties.Properties.Interval = pointer.To(model.Interval)
			}

			if metadata.ResourceData.HasChange("rules") {
				rulesValue, err := expandPrometheusRuleModel(model.Rules)
				if err != nil {
					return err
				}

				if rulesValue != nil {
					properties.Properties.Rules = pointer.From(rulesValue)
				}
			}

			if metadata.ResourceData.HasChange("scopes") {
				properties.Properties.Scopes = model.Scopes
			}

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = pointer.To(model.Tags)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r AlertPrometheusRuleGroupResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.AlertPrometheusRuleGroupClient

			id, err := prometheusrulegroups.ParsePrometheusRuleGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			state := AlertPrometheusRuleGroupResourceModel{
				Name:              id.PrometheusRuleGroupName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(model.Location),
			}

			state.ClusterName = pointer.From(model.Properties.ClusterName)

			state.Description = pointer.From(model.Properties.Description)

			state.Enabled = pointer.From(model.Properties.Enabled)

			state.Interval = pointer.From(model.Properties.Interval)

			rulesValue, err := flattenPrometheusRuleModel(&model.Properties.Rules)
			if err != nil {
				return err
			}

			state.Rules = rulesValue

			state.Scopes = model.Properties.Scopes

			state.Tags = pointer.From(model.Tags)

			return metadata.Encode(&state)
		},
	}
}

func (r AlertPrometheusRuleGroupResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.AlertPrometheusRuleGroupClient

			id, err := prometheusrulegroups.ParsePrometheusRuleGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandPrometheusRuleModel(inputList []PrometheusRuleModel) (*[]prometheusrulegroups.PrometheusRule, error) {
	outputList := make([]prometheusrulegroups.PrometheusRule, 0)
	for _, v := range inputList {
		input := v
		output := prometheusrulegroups.PrometheusRule{
			Annotations: &input.Annotations,
			Enabled:     &input.Enabled,
			Expression:  input.Expression,
			Labels:      &input.Labels,
			Severity:    &input.Severity,
		}

		actionsValue, err := expandPrometheusRuleGroupActionModel(input.Actions)
		if err != nil {
			return nil, err
		}

		output.Actions = actionsValue

		if input.Alert != "" {
			output.Alert = &input.Alert
		}

		if input.For != "" {
			output.For = &input.For
		}

		if input.Record != "" {
			output.Record = &input.Record
		}

		resolveConfigurationValue, err := expandPrometheusRuleResolveConfigurationModel(input.ResolveConfiguration)
		if err != nil {
			return nil, err
		}

		output.ResolveConfiguration = resolveConfigurationValue

		outputList = append(outputList, output)
	}

	return &outputList, nil
}

func expandPrometheusRuleGroupActionModel(inputList []PrometheusRuleGroupActionModel) (*[]prometheusrulegroups.PrometheusRuleGroupAction, error) {
	outputList := make([]prometheusrulegroups.PrometheusRuleGroupAction, 0)
	for _, v := range inputList {
		input := v
		output := prometheusrulegroups.PrometheusRuleGroupAction{
			ActionProperties: &input.ActionProperties,
		}

		if input.ActionGroupId != "" {
			output.ActionGroupId = &input.ActionGroupId
		}

		outputList = append(outputList, output)
	}

	return &outputList, nil
}

func expandPrometheusRuleResolveConfigurationModel(inputList []PrometheusRuleResolveConfigurationModel) (*prometheusrulegroups.PrometheusRuleResolveConfiguration, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	input := &inputList[0]
	output := prometheusrulegroups.PrometheusRuleResolveConfiguration{
		AutoResolved: &input.AutoResolved,
	}

	if input.TimeToResolve != "" {
		output.TimeToResolve = &input.TimeToResolve
	}

	return &output, nil
}

func flattenPrometheusRuleModel(inputList *[]prometheusrulegroups.PrometheusRule) ([]PrometheusRuleModel, error) {
	outputList := make([]PrometheusRuleModel, 0)
	if inputList == nil {
		return outputList, nil
	}

	for _, input := range *inputList {
		output := PrometheusRuleModel{
			Expression: input.Expression,
		}

		actionsValue, err := flattenPrometheusRuleGroupActionModel(input.Actions)
		if err != nil {
			return nil, err
		}

		output.Actions = actionsValue

		output.Alert = pointer.From(input.Alert)

		output.Annotations = pointer.From(input.Annotations)

		output.Enabled = pointer.From(input.Enabled)

		output.For = pointer.From(input.For)

		output.Labels = pointer.From(input.Labels)

		output.Record = pointer.From(input.Record)

		resolveConfigurationValue, err := flattenPrometheusRuleResolveConfigurationModel(input.ResolveConfiguration)
		if err != nil {
			return nil, err
		}

		output.ResolveConfiguration = resolveConfigurationValue

		output.Severity = pointer.From(input.Severity)

		outputList = append(outputList, output)
	}

	return outputList, nil
}

func flattenPrometheusRuleGroupActionModel(inputList *[]prometheusrulegroups.PrometheusRuleGroupAction) ([]PrometheusRuleGroupActionModel, error) {
	outputList := make([]PrometheusRuleGroupActionModel, 0)
	if inputList == nil {
		return outputList, nil
	}

	for _, input := range *inputList {
		output := PrometheusRuleGroupActionModel{}

		output.ActionGroupId = pointer.From(input.ActionGroupId)

		output.ActionProperties = pointer.From(input.ActionProperties)

		outputList = append(outputList, output)
	}

	return outputList, nil
}

func flattenPrometheusRuleResolveConfigurationModel(input *prometheusrulegroups.PrometheusRuleResolveConfiguration) ([]PrometheusRuleResolveConfigurationModel, error) {
	outputList := make([]PrometheusRuleResolveConfigurationModel, 0)
	if input == nil {
		return outputList, nil
	}

	output := PrometheusRuleResolveConfigurationModel{}

	output.AutoResolved = pointer.From(input.AutoResolved)

	output.TimeToResolve = pointer.From(input.TimeToResolve)

	return append(outputList, output), nil
}
