package mongoclusters

<<<<<<< HEAD
import (
	
)

=======
>>>>>>> 7a921d7afc5b9cf5038ddcdec068d7c1c5160c66
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CheckNameAvailabilityResponse struct {
	Message       *string                      `json:"message,omitempty"`
	NameAvailable *bool                        `json:"nameAvailable,omitempty"`
	Reason        *CheckNameAvailabilityReason `json:"reason,omitempty"`
}
