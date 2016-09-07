package cluster

import (
	host_utils "github.com/cheyang/fog/host"
	"github.com/cheyang/fog/types"
)

func provisionVMs(spec types.Spec) (hosts []types.Host, err error) {
	bus := make(chan types.Host)
	defer close(bus)
	vmSpecs, err := host_utils.BuildHostConfigs(spec)
	if err != nil {
		return hosts, err
	}

	hostCount := len(vmSpecs)
	err = host_utils.CreateInBatch(vmSpecs, bus)
	if err != nil {
		return err
	}

	hosts = make([]types.Host, hostCount)
	for i := 0; i < hostCount; i++ {
		hosts[i] = <-bus
	}

	for _, host := range hosts {
		if host.Err != nil {
			return hosts, host.Err
		}
	}

	return hosts, nil
}
