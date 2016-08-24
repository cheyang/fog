package k8s

import (
	"github.com/cheyang/fog/types"
	"github.com/docker/machine/libmachine/drivers"
)

type Aliyun struct {
	hosts []types.Host
}

func New() *Aliyun {
	return &Aliyun{}
}

func (this *Aliyun) setConfigFromFlags(opts drivers.DriverOptions) error {
	return nil
}
func (this *Aliyun) SetHosts(hosts []types.Host) {

}
func (this *Aliyun) Configure() error {
	return nil
}
