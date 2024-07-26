# sightseer.go

Tour your infrastructure using Go.

This is a work in progress and is not recommended for production use.

## Examples

These example aren't guaranteed to compile. They are composed of different
snippets of code copied from an existing test suite, so might need a bit of work
before they run, but all the important bits are there.

### Testing the configuration of a VM on Azure

```go
package main

import (
	"testing"

	"github.com/living-etc/sightser.go/azure"
)

const (
	subscriptionId    = "3498jf39-98h4-89u3444r-983u4r89"
	resourceGroupName = "MyResourceGroup"
)

var azureclient = azure.NewAzureClient(subscriptionId, resourceGroupName)

var tests = []struct {
	vmName    string
	privateIP string
}{
	{
		vmName:    "worker-1",
		privateIP: "10.240.0.21",
	},
	{
		vmName:    "worker-2",
		privateIP: "10.240.0.22",
	},
	{
		vmName:    "worker-3",
		privateIP: "10.240.0.23",
	},
}

func TestCompute(t *testing.T) {
	for _, tt := range tests {
		vm, _ := azureclient.VmFromName(tt.vmName)

		t.Run(tt.vmName+" correct private IP", func(t *testing.T) {
			if vm.PrivateIPAddress != tt.privateIP {
				t.Errorf("Want %v, got %v", tt.privateIP, vm.PrivateIPAddress)
			}
		})

		reachable, _ := vm.ReachableOnPort(22)
		t.Run(tt.vmName+" reachable on port 22", func(t *testing.T) {
			if !reachable {
				t.Errorf("%v: not reachable on port 22", tt.vmName)
			}
		})

		connectable, _ := vm.ConnectableOverSSH("../0-keys/id_rsa.pub")
		t.Run(tt.vmName+" connectable over SSH", func(t *testing.T) {
			if !connectable {
				t.Errorf("%v: not connectable over ssh", tt.vmName)
			}
		})

		t.Run(tt.vmName+" correct DNS name", func(t *testing.T) {
			want := tt.vmName + "-kthw-cw.uksouth.cloudapp.azure.com"
			got := vm.DnsName
			if got != want {
				t.Errorf("want %v, got %v", want, got)
			}
		})
	}
}
```

### Verify files have been upload to server

```go
package main

import (
	"testing"

	"github.com/living-etc/sightseer.go/azure"
)

const (
	subscriptionId    = "3498jf39-98h4-89u3444r-983u4r89"
	resourceGroupName = "MyResourceGroup"
)

var azureclient = azure.NewAzureClient(subscriptionId, resourceGroupName)

var tests = []struct {
	vmName    string
	privateIP string
}{
	{
		vmName:    "controller-1",
		privateIP: "10.240.0.11",
	},
	{
		vmName:    "controller-2",
		privateIP: "10.240.0.12",
	},
	{
		vmName:    "controller-3",
		privateIP: "10.240.0.13",
	},
}

func TestControllerConfig(t *testing.T) {
	for _, tt := range tests {
		vm, _ := azureclient.VmFromName(tt.vmName)

		hostname := vm.Hostname()
		if hostname != tt.vmName {
			t.Errorf("wanted %v, got %v", tt.vmName, hostname)
		}

		etcdFiles := []string{
			"/etc/etcd/ca.pem",
			"/etc/etcd/kubernetes-key.pem",
			"/etc/etcd/kubernetes.pem",
		}

		for _, file := range etcdFiles {
			t.Run(tt.vmName+" has "+file, func(t *testing.T) {
				hasFile := vm.HasFile(file)
				if !hasFile {
					t.Errorf("%v does not have %v", tt.vmName, file)
				}
			})
		}
	}
}
```
