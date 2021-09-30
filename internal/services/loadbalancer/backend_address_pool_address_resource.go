package loadbalancer

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-02-01/network"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/validate"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var _ sdk.Resource = BackendAddressPoolAddressResource{}
var _ sdk.ResourceWithUpdate = BackendAddressPoolAddressResource{}

type BackendAddressPoolAddressResource struct{}

type BackendAddressPoolAddressModel struct {
	Name                      string `tfschema:"name"`
	BackendAddressPoolId      string `tfschema:"backend_address_pool_id"`
	VirtualNetworkId          string `tfschema:"virtual_network_id"`
	IPAddress                 string `tfschema:"ip_address"`
	FrontendIpConfigurationId string `tfschema:"frontend_ip_configuration_id"`
	SubnetId                  string `tfschema:"subnet_id"`
}

func (r BackendAddressPoolAddressResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"backend_address_pool_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.LoadBalancerBackendAddressPoolID,
		},

		"virtual_network_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: networkValidate.VirtualNetworkID,
		},

		"ip_address": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.IsIPAddress,
		},

		"frontend_ip_configuration_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.LoadBalancerFrontendIpConfigurationID,
		},

		"subnet_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: networkValidate.SubnetID,
		},
	}
}

func (r BackendAddressPoolAddressResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r BackendAddressPoolAddressResource) ModelObject() interface{} {
	return &BackendAddressPoolAddressModel{}
}

func (r BackendAddressPoolAddressResource) ResourceType() string {
	return "azurerm_lb_backend_address_pool_address"
}

func (r BackendAddressPoolAddressResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LoadBalancers.LoadBalancerBackendAddressPoolsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model BackendAddressPoolAddressModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			poolId, err := parse.LoadBalancerBackendAddressPoolID(model.BackendAddressPoolId)
			if err != nil {
				return err
			}

			locks.ByName(poolId.BackendAddressPoolName, backendAddressPoolResourceName)
			defer locks.UnlockByName(poolId.BackendAddressPoolName, backendAddressPoolResourceName)

			// Backend Addresses can only be created for Standard LB's - not Basic, so we have to check
			lb, err := metadata.Client.LoadBalancers.LoadBalancersClient.Get(ctx, poolId.ResourceGroup, poolId.LoadBalancerName, "")
			if err != nil {
				return fmt.Errorf("retrieving Load Balancer %q (Resource Group %q): %+v", poolId.LoadBalancerName, poolId.ResourceGroup, err)
			}
			isStandardSku := false
			if lb.Sku != nil && lb.Sku.Name == network.LoadBalancerSkuNameStandard {
				isStandardSku = true
			}
			if !isStandardSku {
				return fmt.Errorf("Backend Addresses are only supported on Standard SKU Load Balancers")
			}

			if lb.Sku.Tier == network.LoadBalancerSkuTierRegional && model.IPAddress == "" {
				return fmt.Errorf("IP Address is required on Regional SKU Tier Load Balancers")
			}

			id := parse.NewBackendAddressPoolAddressID(subscriptionId, poolId.ResourceGroup, poolId.LoadBalancerName, poolId.BackendAddressPoolName, model.Name)
			pool, err := client.Get(ctx, poolId.ResourceGroup, poolId.LoadBalancerName, poolId.BackendAddressPoolName)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *poolId, err)
			}
			if pool.BackendAddressPoolPropertiesFormat == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *poolId)
			}

			addresses := make([]network.LoadBalancerBackendAddress, 0)
			if pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses != nil {
				addresses = *pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses
			}

			metadata.Logger.Infof("checking for existing %s..", id)
			for _, address := range addresses {
				if address.Name == nil {
					continue
				}

				if *address.Name == id.AddressName {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			backendAddressParameters := &network.LoadBalancerBackendAddressPropertiesFormat{
				VirtualNetwork: &network.SubResource{
					ID: utils.String(model.VirtualNetworkId),
				},
			}

			if model.IPAddress != "" {
				backendAddressParameters.IPAddress = utils.String(model.IPAddress)
			}

			if model.FrontendIpConfigurationId != "" {
				backendAddressParameters.LoadBalancerFrontendIPConfiguration = &network.SubResource{
					ID: utils.String(model.FrontendIpConfigurationId),
				}
			}

			if model.SubnetId != "" {
				backendAddressParameters.Subnet = &network.SubResource{
					ID: utils.String(model.SubnetId),
				}
			}

			addresses = append(addresses, network.LoadBalancerBackendAddress{
				LoadBalancerBackendAddressPropertiesFormat: backendAddressParameters,
				Name: utils.String(id.AddressName),
			})
			pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses = &addresses

			metadata.Logger.Infof("adding %s..", id)
			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName, pool)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			metadata.Logger.Infof("waiting for update %s..", id)
			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for update of %s: %+v", id, err)
			}
			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r BackendAddressPoolAddressResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LoadBalancers.LoadBalancerBackendAddressPoolsClient
			id, err := parse.BackendAddressPoolAddressID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			pool, err := client.Get(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if pool.BackendAddressPoolPropertiesFormat == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			var backendAddress *network.LoadBalancerBackendAddress
			if pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses != nil {
				for _, address := range *pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses {
					if address.Name == nil {
						continue
					}

					if *address.Name == id.AddressName {
						backendAddress = &address
						break
					}
				}
			}
			if backendAddress == nil {
				return metadata.MarkAsGone(id)
			}

			backendAddressPoolId := parse.NewLoadBalancerBackendAddressPoolID(id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)
			model := BackendAddressPoolAddressModel{
				Name:                 id.AddressName,
				BackendAddressPoolId: backendAddressPoolId.ID(),
			}

			if props := backendAddress.LoadBalancerBackendAddressPropertiesFormat; props != nil {
				if props.IPAddress != nil {
					model.IPAddress = *props.IPAddress
				}

				if props.VirtualNetwork != nil && props.VirtualNetwork.ID != nil {
					model.VirtualNetworkId = *props.VirtualNetwork.ID
				}

				if props.LoadBalancerFrontendIPConfiguration != nil && props.LoadBalancerFrontendIPConfiguration.ID != nil {
					model.FrontendIpConfigurationId = *props.LoadBalancerFrontendIPConfiguration.ID
				}

				if props.Subnet != nil && props.Subnet.ID != nil {
					model.SubnetId = *props.Subnet.ID
				}
			}

			return metadata.Encode(&model)
		},
		Timeout: 5 * time.Minute,
	}
}

func (r BackendAddressPoolAddressResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LoadBalancers.LoadBalancerBackendAddressPoolsClient
			id, err := parse.BackendAddressPoolAddressID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.BackendAddressPoolName, backendAddressPoolResourceName)
			defer locks.UnlockByName(id.BackendAddressPoolName, backendAddressPoolResourceName)

			pool, err := client.Get(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if pool.BackendAddressPoolPropertiesFormat == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			addresses := make([]network.LoadBalancerBackendAddress, 0)
			if pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses != nil {
				addresses = *pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses
			}

			newAddresses := make([]network.LoadBalancerBackendAddress, 0)
			for _, address := range addresses {
				if address.Name == nil {
					continue
				}

				if *address.Name != id.AddressName {
					newAddresses = append(newAddresses, address)
				}
			}
			pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses = &newAddresses

			metadata.Logger.Infof("removing %s..", *id)
			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName, pool)
			if err != nil {
				return fmt.Errorf("removing %s: %+v", *id, err)
			}
			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for removal of %s: %+v", *id, err)
			}

			//time.Sleep(2*time.Minute)
			//// there appears to be an eventual consistency issue here
			//timeout, _ := ctx.Deadline()
			//log.Printf("[DEBUG] Waiting for %s to be eventually deleted", *id)
			//stateConf := &pluginsdk.StateChangeConf{
			//	Pending:                   []string{"Exists"},
			//	Target:                    []string{"NotFound"},
			//	Refresh:                   loadbalancerBackendAddressDeleteStateRefreshFunc(ctx, client, *id),
			//	MinTimeout:                10 * time.Second,
			//	ContinuousTargetOccurence: 10,
			//	Timeout:                   time.Until(timeout),
			//}
			//
			//if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			//	return fmt.Errorf("waiting for %s to be deleted: %+v", *id, err)
			//}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

//func loadbalancerBackendAddressDeleteStateRefreshFunc(ctx context.Context, client *network.LoadBalancerBackendAddressPoolsClient, id parse.BackendAddressPoolAddressId) pluginsdk.StateRefreshFunc {
//	// Whilst the load balancer backend address is deleted quickly, it appears it's not actually finished replicating at this time
//	// so the deletion of the parent Shared Image fails with "can not delete until nested resources are deleted"
//	// ergo we need to poll on this for a bit
//	return func() (interface{}, string, error) {
//		pool, err := client.Get(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)
//		if err != nil {
//			return nil, "", fmt.Errorf("failed to poll to check if the load balancer backend address has been deleted: %+v", err)
//		}
//		if pool.BackendAddressPoolPropertiesFormat == nil {
//			return "NotFound", "NotFound", nil
//		}
//
//		exist := false
//		if pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses != nil {
//			for _, address := range *pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses {
//				if address.Name == nil {
//					continue
//				}
//
//				if *address.Name == id.AddressName {
//					exist = true
//					break
//				}
//			}
//
//			if !exist {
//				return "NotFound", "NotFound", nil
//			}
//
//		}else{
//			return "NotFound", "NotFound", nil
//		}
//
//		return pool, "Exists", nil
//	}
//}

func (r BackendAddressPoolAddressResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.BackendAddressPoolAddressID
}

func (r BackendAddressPoolAddressResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LoadBalancers.LoadBalancerBackendAddressPoolsClient
			id, err := parse.BackendAddressPoolAddressID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.BackendAddressPoolName, backendAddressPoolResourceName)
			defer locks.UnlockByName(id.BackendAddressPoolName, backendAddressPoolResourceName)

			var model BackendAddressPoolAddressModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			pool, err := client.Get(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if pool.BackendAddressPoolPropertiesFormat == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			addresses := make([]network.LoadBalancerBackendAddress, 0)
			if pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses != nil {
				addresses = *pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses
			}
			index := -1
			for i, address := range addresses {
				if address.Name == nil {
					continue
				}

				if *address.Name == id.AddressName {
					index = i
					break
				}
			}
			if index == -1 {
				return fmt.Errorf("%s was not found", *id)
			}

			backendAddressParameters := &network.LoadBalancerBackendAddressPropertiesFormat{
				VirtualNetwork: &network.SubResource{
					ID: utils.String(model.VirtualNetworkId),
				},
			}

			if model.IPAddress != "" {
				backendAddressParameters.IPAddress = utils.String(model.IPAddress)
			}

			if model.FrontendIpConfigurationId != "" {
				backendAddressParameters.LoadBalancerFrontendIPConfiguration = &network.SubResource{
					ID: utils.String(model.FrontendIpConfigurationId),
				}
			}

			if model.SubnetId != "" {
				backendAddressParameters.Subnet = &network.SubResource{
					ID: utils.String(model.SubnetId),
				}
			}

			addresses[index] = network.LoadBalancerBackendAddress{
				LoadBalancerBackendAddressPropertiesFormat: backendAddressParameters,
				Name: utils.String(id.AddressName),
			}
			pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses = &addresses

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName, pool)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}
			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for update of %s: %+v", *id, err)
			}
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
