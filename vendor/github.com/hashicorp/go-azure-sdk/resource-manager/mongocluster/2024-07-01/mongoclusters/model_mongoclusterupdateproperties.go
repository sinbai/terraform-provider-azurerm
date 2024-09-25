package mongoclusters

<<<<<<< HEAD
import (
	
)

=======
>>>>>>> 7a921d7afc5b9cf5038ddcdec068d7c1c5160c66
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoClusterUpdateProperties struct {
	Administrator       *AdministratorProperties    `json:"administrator,omitempty"`
	Backup              *BackupProperties           `json:"backup,omitempty"`
	Compute             *ComputeProperties          `json:"compute,omitempty"`
	HighAvailability    *HighAvailabilityProperties `json:"highAvailability,omitempty"`
	PreviewFeatures     *[]PreviewFeature           `json:"previewFeatures,omitempty"`
	PublicNetworkAccess *PublicNetworkAccess        `json:"publicNetworkAccess,omitempty"`
	ServerVersion       *string                     `json:"serverVersion,omitempty"`
	Sharding            *ShardingProperties         `json:"sharding,omitempty"`
	Storage             *StorageProperties          `json:"storage,omitempty"`
}
