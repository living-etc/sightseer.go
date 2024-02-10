package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
)

type AzureClient struct {
	subscriptionId    string
	resourceGroupName string
	credential        *azidentity.DefaultAzureCredential
}

func NewAzureClient(subscriptionId string, resourceGroupName string) AzureClient {
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	check(err, "Failed to get OAuth config")

	return AzureClient{
		subscriptionId:    subscriptionId,
		resourceGroupName: resourceGroupName,
		credential:        credential,
	}
}

func (client AzureClient) VmFromName(name string) (VM, error) {
	var result VM

	subscriptionId := client.subscriptionId
	resourceGroupName := client.resourceGroupName
	credential := client.credential
	ctx := context.Background()

	vmClient, err := armcompute.NewVirtualMachinesClient(subscriptionId, credential, nil)
	check(err, "Failed to get compute client")

	vm, err := vmClient.Get(ctx, resourceGroupName, name, nil)
	check(err, "Could not retrieve instance view")

	nicRef := vm.Properties.NetworkProfile.NetworkInterfaces[0]
	nicID, err := arm.ParseResourceID(*nicRef.ID)
	check(err, "Unable to parse nic resource id")

	nicClient, err := armnetwork.NewInterfacesClient(subscriptionId, credential, nil)
	check(err, "Unable to get network interfaces client")

	nic, err := nicClient.Get(ctx, resourceGroupName, nicID.Name, nil)
	check(err, "Unable to get nic")

	for _, ipConfig := range nic.Properties.IPConfigurations {
		if ipAddress := ipConfig.Properties.PrivateIPAddress; ipAddress != nil {
			result.PrivateIPAddress = *ipAddress
		}
		if ipAddress := ipConfig.Properties.PublicIPAddress; ipAddress != nil {
			publicIPClient, err := armnetwork.NewPublicIPAddressesClient(
				subscriptionId,
				credential,
				nil,
			)
			check(err, "Could not initialise the public ip client")

			publicIPID, err := arm.ParseResourceID(*ipAddress.ID)
			check(err, "Could not parse the public ID resource ID")

			publicIP, err := publicIPClient.Get(ctx, resourceGroupName, publicIPID.Name, nil)
			check(err, "Could not get public IP address")

			result.PublicIPAddress = *publicIP.PublicIPAddress.Properties.IPAddress
			result.DnsName = *publicIP.PublicIPAddress.Properties.DNSSettings.Fqdn
		}
	}

	return result, nil
}
