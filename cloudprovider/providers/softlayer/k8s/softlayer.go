package k8s

import (
	"github.com/cheyang/fog/cloudprovider"
	"github.com/cheyang/fog/types"
	"github.com/docker/machine/libmachine/drivers"
)

type Softlayer struct {
	hosts []*types.Host
}

func New() cloudprovider.CloudInterface {
	return &Softlayer{}
}

func (this *Softlayer) SetConfigFromFlags(opts drivers.DriverOptions) error {
	return nil
}
func (this *Softlayer) SetHosts(hosts []*types.Host) {

}
func (this *Softlayer) Configure() error {
	return nil
}
