package k8s

import (
	"github.com/cheyang/fog/types"
	"github.com/docker/machine/libmachine/drivers"
)

type Softlayer struct {
	hosts []types.Host
}

func New() *Softlayer {
	return &Softlayer{}
}

func (this *Softlayer) setConfigFromFlags(opts drivers.DriverOptions) error {
	return nil
}
func (this *Softlayer) SetHosts(hosts []types.Host) {

}
func (this *Softlayer) Configure() error {
	return nil
}
