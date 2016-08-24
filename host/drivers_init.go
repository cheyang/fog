package host

import (
	"fmt"

	"github.com/cheyang/fog/types"
	"github.com/denverdino/docker-machine-driver-aliyunecs/aliyunecs"
	"github.com/docker/machine/drivers/amazonec2"
	"github.com/docker/machine/drivers/openstack"
	"github.com/docker/machine/drivers/softlayer"
	"github.com/docker/machine/libmachine/drivers"
)

var initFuncMaps = map[string]func(hostname, storePath string) drivers.Driver{
	"aliyun":    aliyunecs.NewDriver,
	"softlayer": softlayer.NewDriver,
	"aws":       amazonec2.NewDriver,
	"openstack": openstack.NewDriver,
}

func initDrivers(driverName, hostConfig types.VMSpec, storePath string) (drivers.Driver, error) {

	if driverFunc, found := initFuncMaps[driverName]; !found {
		return nil, fmt.Errorf("Driver %s is not found.", driverName)
	}
	d := driverFunc(hostConfig.Name, storePath)

	props := hostConfig.Properties
	opts := NewConfigFlagger(props)
	d.SetConfigFromFlags(opts)

	return d, nil
}
