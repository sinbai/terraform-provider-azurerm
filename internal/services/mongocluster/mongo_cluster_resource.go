// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mongocluster

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mongocluster/2024-07-01/mongoclusters"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type MongoClusterResource struct{}

var _ sdk.ResourceWithUpdate = MongoClusterResource{}

var _ sdk.ResourceWithCustomizeDiff = MongoClusterResource{}

type HighAvailabilityModel struct {
	Mode string `tfschema:"mode"`
}

func (r MongoClusterResource) ModelObject() interface{} {
	return &MongoClusterResourceModel{}
}

type MongoClusterResourceModel struct {
	Name                        string                  `tfschema:"name"`
	ResourceGroupName           string                  `tfschema:"resource_group_name"`
	Location                    string                  `tfschema:"location"`
	AdministratorLogin          string                  `tfschema:"administrator_login"`
	AdministratorLoginPassword  string                  `tfschema:"administrator_login_password"`
	CreateMode                  string                  `tfschema:"create_mode"`
	PointInTimeRestoreTimeInUtc string                  `tfschema:"point_in_time_restore_time_in_utc"`
	ShardCount                  int64                   `tfschema:"shard_count"`
	SourceServerId              string                  `tfschema:"source_server_id"`
	SourceLocation              string                  `tfschema:"source_location"`
	ComputeTier                 string                  `tfschema:"compute_tier"`
	HighAvailability            []HighAvailabilityModel `tfschema:"high_availability"`
	PublicNetworkAccessEnabled  bool                    `tfschema:"public_network_access_enabled"`
	StorageSizeInGb             int64                   `tfschema:"storage_size_in_gb"`
	Tags                        map[string]string       `tfschema:"tags"`
	Version                     string                  `tfschema:"version"`
}

func (r MongoClusterResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return mongoclusters.ValidateMongoClusterID
}

func (r MongoClusterResource) ResourceType() string {
	return "azurerm_mongo_cluster"
}

func (r MongoClusterResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 40),
				validation.StringMatch(
					regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`),
					"The name contains only lowercase letters, numbers and hyphens.",
				),
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"administrator_login": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"administrator_login_password": {
			Type:      pluginsdk.TypeString,
			Optional:  true,
			Sensitive: true,
		},

		"create_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(mongoclusters.PossibleValuesForCreateMode(), false),
		},

		"point_in_time_restore_time_in_utc": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},

		"shard_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(1),
			//ForceNew: true,
		},

		"source_server_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: mongoclusters.ValidateMongoClusterID,
		},

		"source_location": commonschema.LocationOptional(),

		"compute_tier": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				"Free",
				"M25",
				"M30",
				"M40",
				"M50",
				"M60",
				"M80",
			}, false),
		},

		"high_availability": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"mode": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(mongoclusters.PossibleValuesForHighAvailabilityMode(), false),
					},
				},
			},
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"storage_size_in_gb": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(32, 16384),
		},

		"tags": commonschema.Tags(),

		"version": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r MongoClusterResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r MongoClusterResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MongoCluster.MongoClustersClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state MongoClusterResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := mongoclusters.NewMongoClusterID(subscriptionId, state.ResourceGroupName, state.Name)
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameter := mongoclusters.MongoCluster{
				Location:   azure.NormalizeLocation(state.Location),
				Name:       pointer.To(state.Name),
				Properties: &mongoclusters.MongoClusterProperties{},
			}

			if _, ok := metadata.ResourceData.GetOk("administrator_login"); ok {
				parameter.Properties.Administrator = &mongoclusters.AdministratorProperties{
					UserName: pointer.To(state.AdministratorLogin),
					Password: pointer.To(state.AdministratorLoginPassword),
				}
			}

			if _, ok := metadata.ResourceData.GetOk("compute_tier"); ok {
				parameter.Properties.Compute = &mongoclusters.ComputeProperties{
					Tier: pointer.To(state.ComputeTier),
				}
			}

			if _, ok := metadata.ResourceData.GetOk("version"); ok {
				parameter.Properties.ServerVersion = pointer.To(state.Version)
			}

			if _, ok := metadata.ResourceData.GetOk("shard_count"); ok {
				parameter.Properties.Sharding = &mongoclusters.ShardingProperties{
					ShardCount: pointer.To(state.ShardCount),
				}
			}

			if _, ok := metadata.ResourceData.GetOk("storage_size_in_gb"); ok {
				parameter.Properties.Storage = &mongoclusters.StorageProperties{
					SizeGb: pointer.To(state.StorageSizeInGb),
				}
			}

			if _, ok := metadata.ResourceData.GetOk("create_mode"); ok {
				parameter.Properties.CreateMode = pointer.To(mongoclusters.CreateMode(state.CreateMode))
			}

			if _, ok := metadata.ResourceData.GetOk("high_availability"); ok {
				parameter.Properties.HighAvailability = &mongoclusters.HighAvailabilityProperties{
					TargetMode: pointer.To(mongoclusters.HighAvailabilityMode(state.HighAvailability[0].Mode)),
				}
			}

			if state.PublicNetworkAccessEnabled {
				parameter.Properties.PublicNetworkAccess = pointer.To(mongoclusters.PublicNetworkAccessEnabled)
			} else {
				parameter.Properties.PublicNetworkAccess = pointer.To(mongoclusters.PublicNetworkAccessDisabled)
			}

			if _, ok := metadata.ResourceData.GetOk("source_server_id"); ok {
				switch state.CreateMode {
				case string(mongoclusters.CreateModeGeoReplica):
					parameter.Properties.ReplicaParameters = &mongoclusters.MongoClusterReplicaParameters{
						SourceResourceId: state.SourceServerId,
						SourceLocation:   state.SourceLocation,
					}
				case string(mongoclusters.CreateModeReplica):
					parameter.Properties.ReplicaParameters = &mongoclusters.MongoClusterReplicaParameters{
						SourceResourceId: state.SourceServerId,
					}
				case string(mongoclusters.CreateModePointInTimeRestore):
					parameter.Properties.RestoreParameters = &mongoclusters.MongoClusterRestoreParameters{
						SourceResourceId: pointer.To(state.SourceServerId),
					}
					if state.PointInTimeRestoreTimeInUtc != "" {
						v, err := time.Parse(time.RFC3339, state.PointInTimeRestoreTimeInUtc)
						if err != nil {
							return fmt.Errorf("unable to parse `point_in_time_restore_time_in_utc` value")
						}
						parameter.Properties.RestoreParameters.SetPointInTimeUTCAsTime(v)
					}
				}
			}

			if _, ok := metadata.ResourceData.GetOk("tags"); ok {
				parameter.Tags = pointer.To(state.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, parameter); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r MongoClusterResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MongoCluster.MongoClustersClient

			id, err := mongoclusters.ParseMongoClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Info("Decoding state...")
			var state MongoClusterResourceModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := existing.Model
			if model == nil {
				return fmt.Errorf("reading %s for update: model was nil", *id)
			}

			metadata.Logger.Infof("updating %s", id)

			model.SystemData = nil

			if metadata.ResourceData.HasChange("compute_tier") {
				model.Properties.Compute = &mongoclusters.ComputeProperties{
					Tier: pointer.To(state.ComputeTier),
				}
				oldComputeTier, newComputeTier := metadata.ResourceData.GetChange("compute_tier")
				if (oldComputeTier == "Free" || oldComputeTier == "M25") && newComputeTier != "Free" && newComputeTier != "M25" {
					metadata.Logger.Infof("updating cluster tier for %s", id)
					// upgrades involving Free or M25(Burstable) cluster tier require first upgrading the cluster tier, after which other configurations can be updated.
					if err := client.CreateOrUpdateThenPoll(ctx, *id, *model); err != nil {
						return fmt.Errorf("updating %s: %+v", id, err)
					}
				}
			}

			metadata.Logger.Infof("updating other configurations for %s", id)
			if metadata.ResourceData.HasChange("administrator_login_password") {
				model.Properties.Administrator = &mongoclusters.AdministratorProperties{
					UserName: pointer.To(state.AdministratorLogin),
					Password: pointer.To(state.AdministratorLoginPassword),
				}
			}

			if metadata.ResourceData.HasChange("high_availability") {
				model.Properties.HighAvailability = &mongoclusters.HighAvailabilityProperties{
					TargetMode: pointer.To(mongoclusters.HighAvailabilityMode(state.HighAvailability[0].Mode)),
				}
			}

			if metadata.ResourceData.HasChange("public_network_access_enabled") {
				if state.PublicNetworkAccessEnabled {
					model.Properties.PublicNetworkAccess = pointer.To(mongoclusters.PublicNetworkAccessEnabled)
				} else {
					model.Properties.PublicNetworkAccess = pointer.To(mongoclusters.PublicNetworkAccessDisabled)
				}
			}

			if metadata.ResourceData.HasChange("storage_size_in_gb") {
				model.Properties.Storage = &mongoclusters.StorageProperties{
					SizeGb: pointer.To(state.StorageSizeInGb),
				}
			}

			if metadata.ResourceData.HasChange("shard_count") {
				model.Properties.Sharding = &mongoclusters.ShardingProperties{
					ShardCount: pointer.To(state.ShardCount),
				}
			}

			if metadata.ResourceData.HasChange("version") {
				model.Properties.ServerVersion = pointer.To(state.Version)
			}

			if metadata.ResourceData.HasChange("tags") {
				model.Tags = pointer.To(state.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *model); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			//metadata.Logger.Infof("updating %s", id)
			//parameter := mongoclusters.MongoClusterUpdate{
			//	Properties: &mongoclusters.MongoClusterUpdateProperties{},
			//}
			//
			//if metadata.ResourceData.HasChange("compute_tier") {
			//	parameter.Properties.Compute = &mongoclusters.ComputeProperties{
			//		Tier: pointer.To(state.ComputeTier),
			//	}
			//	oldComputeTier, newComputeTier := metadata.ResourceData.GetChange("compute_tier")
			//	if (oldComputeTier == "Free" || oldComputeTier == "M25") && newComputeTier != "Free" && newComputeTier != "M25" {
			//		metadata.Logger.Infof("updating cluster tier for %s", id)
			//		// upgrades involving Free or M25(Burstable) cluster tier require first upgrading the cluster tier, after which other configurations can be updated.
			//		if err := client.UpdateThenPoll(ctx, *id, parameter); err != nil {
			//			return fmt.Errorf("updating %s: %+v", id, err)
			//		}
			//	}
			//}
			//
			//metadata.Logger.Infof("updating other configurations for %s", id)
			//if metadata.ResourceData.HasChange("administrator_login_password") {
			//	parameter.Properties.Administrator = &mongoclusters.AdministratorProperties{
			//		UserName: pointer.To(state.AdministratorLogin),
			//		Password: pointer.To(state.AdministratorLoginPassword),
			//	}
			//}
			//
			//if metadata.ResourceData.HasChange("high_availability") {
			//	parameter.Properties.HighAvailability = &mongoclusters.HighAvailabilityProperties{
			//		TargetMode: pointer.To(mongoclusters.HighAvailabilityMode(state.HighAvailability[0].Mode)),
			//	}
			//}
			//
			//if metadata.ResourceData.HasChange("public_network_access_enabled") {
			//	if state.PublicNetworkAccessEnabled {
			//		parameter.Properties.PublicNetworkAccess = pointer.To(mongoclusters.PublicNetworkAccessEnabled)
			//	} else {
			//		parameter.Properties.PublicNetworkAccess = pointer.To(mongoclusters.PublicNetworkAccessDisabled)
			//	}
			//}
			//
			//if metadata.ResourceData.HasChange("storage_size_in_gb") {
			//	parameter.Properties.Storage = &mongoclusters.StorageProperties{
			//		SizeGb: pointer.To(state.StorageSizeInGb),
			//	}
			//}
			//
			//if metadata.ResourceData.HasChange("shard_count") {
			//	parameter.Properties.Sharding = &mongoclusters.ShardingProperties{
			//		ShardCount: pointer.To(state.ShardCount),
			//	}
			//}
			//
			//if metadata.ResourceData.HasChange("version") {
			//	parameter.Properties.ServerVersion = pointer.To(state.Version)
			//}
			//
			//if metadata.ResourceData.HasChange("tags") {
			//	parameter.Tags = pointer.To(state.Tags)
			//}
			//
			//if err := client.UpdateThenPoll(ctx, *id, parameter); err != nil {
			//	return fmt.Errorf("updating %s: %+v", id, err)
			//}

			return nil
		},
	}
}

func (r MongoClusterResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MongoCluster.MongoClustersClient

			id, err := mongoclusters.ParseMongoClusterID(metadata.ResourceData.Id())
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

			state := MongoClusterResourceModel{
				Name:              id.MongoClusterName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.NormalizeNilable(&model.Location)
				state.Tags = pointer.From(model.Tags)

				if props := model.Properties; props != nil {
					state.AdministratorLoginPassword = metadata.ResourceData.Get("administrator_login_password").(string)
					// API doesn't return the value of create_mode
					state.CreateMode = metadata.ResourceData.Get("create_mode").(string)

					if v := props.Administrator; v != nil {
						state.AdministratorLogin = pointer.From(v.UserName)
					}

					if v := props.RestoreParameters; v != nil {
						state.PointInTimeRestoreTimeInUtc = pointer.From(v.PointInTimeUTC)
						state.SourceServerId = pointer.From(v.SourceResourceId)
					}
					if v := props.ReplicaParameters; v != nil {
						state.SourceLocation = v.SourceLocation
						state.SourceServerId = v.SourceResourceId
					}
					if v := props.Sharding; v != nil {
						state.ShardCount = pointer.From(v.ShardCount)
					}
					if v := props.Compute; v != nil {
						state.ComputeTier = pointer.From(v.Tier)
					}
					state.HighAvailability = flattenMongoClusterHighAvailability(props.HighAvailability)

					state.PublicNetworkAccessEnabled = pointer.From(props.PublicNetworkAccess) == mongoclusters.PublicNetworkAccessEnabled

					if v := props.Storage; v != nil {
						state.StorageSizeInGb = pointer.From(v.SizeGb)
					}
					state.Version = pointer.From(props.ServerVersion)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r MongoClusterResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MongoCluster.MongoClustersClient

			id, err := mongoclusters.ParseMongoClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r MongoClusterResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var state MongoClusterResourceModel
			if err := metadata.DecodeDiff(&state); err != nil {
				return fmt.Errorf("DecodeDiff: %+v", err)
			}

			switch mongoclusters.CreateMode(state.CreateMode) {
			case mongoclusters.CreateModeDefault, "":
				if _, ok := metadata.ResourceDiff.GetOk("administrator_login"); !ok {
					return fmt.Errorf("`administrator_login` is required when `create_mode` is not specified or is specified as %s", string(mongoclusters.CreateModeDefault))
				}

				if _, ok := metadata.ResourceDiff.GetOk("administrator_login_password"); !ok {
					return fmt.Errorf("`administrator_login_password` is required when `create_mode` is not specified or is specified as  %s", string(mongoclusters.CreateModeDefault))
				}

				if _, ok := metadata.ResourceDiff.GetOk("compute_tier"); !ok {
					return fmt.Errorf("`compute_tier` is required when `create_mode` is not specified or is specified as %s", string(mongoclusters.CreateModeDefault))
				}

				if _, ok := metadata.ResourceDiff.GetOk("storage_size_in_gb"); !ok {
					return fmt.Errorf("`storage_size_in_gb` is required when `create_mode` is not specified or is specified as  %s", string(mongoclusters.CreateModeDefault))
				}

				if _, ok := metadata.ResourceDiff.GetOk("high_availability"); !ok {
					return fmt.Errorf("`high_availability` is required when `create_mode` is not specified or is specified as  %s", string(mongoclusters.CreateModeDefault))
				}

				if _, ok := metadata.ResourceDiff.GetOk("shard_count"); !ok {
					return fmt.Errorf("`shard_count` is required when `create_mode` is not specified or is specified as  %s", string(mongoclusters.CreateModeDefault))
				}

				if _, ok := metadata.ResourceDiff.GetOk("version"); !ok {
					return fmt.Errorf("`version` is required when `create_mode` is not specified or is specified as  %s", string(mongoclusters.CreateModeDefault))
				}
			case mongoclusters.CreateModeReplica:
				if state.SourceServerId == "" {
					return fmt.Errorf("`source_server_id` is required when create_mode is %s", string(mongoclusters.CreateModeReplica))
				}
			case mongoclusters.CreateModeGeoReplica:
				if state.SourceLocation == "" || state.SourceServerId == "" {
					return fmt.Errorf("`source_location` and `source_server_id` are required when create_mode is %s", string(mongoclusters.CreateModeGeoReplica))
				}

			case mongoclusters.CreateModePointInTimeRestore:
				if state.PointInTimeRestoreTimeInUtc == "" || state.SourceServerId == "" {
					return fmt.Errorf("`source_server_id` and `point_in_time_restore_time_in_utc` are required when create_mode is %s", string(mongoclusters.CreateModePointInTimeRestore))
				}
			}

			if state.ComputeTier == "Free" || state.ComputeTier == "M25" {
				if len(state.HighAvailability) > 0 && (state.HighAvailability[0].Mode == string(mongoclusters.HighAvailabilityModeSameZone) || state.HighAvailability[0].Mode == string(mongoclusters.HighAvailabilityModeZoneRedundantPreferred)) {
					return fmt.Errorf("high Availability is not available with the 'Free' or 'M25' Cluster Tier")
				}

				if state.ShardCount > 1 {
					return fmt.Errorf("shard count > 1 is not supported for the 'Free' or 'M25' Cluster Tier")
				}
			}

			return nil
		},
	}
}

func flattenMongoClusterHighAvailability(ha *mongoclusters.HighAvailabilityProperties) []HighAvailabilityModel {
	if ha == nil {
		return []HighAvailabilityModel{}
	}

	return []HighAvailabilityModel{
		{
			Mode: string(pointer.From(ha.TargetMode)),
		},
	}
}
