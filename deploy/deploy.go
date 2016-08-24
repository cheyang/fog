package deploy

import (
	"github.com/cheyang/fog/host"
	"github.com/cheyang/fog/types"
	"github.com/cheyang/fog/util/dump"
)

func Run(spec types.Spec) error {

	//register dump tool
	dump.InstallCoreDumpGenerator()

	bus := make(chan types.Host)
	defer close(bus)
	vmSpecs, err := host.BuildHostConfigs(spec)
	if err != nil {
		return err
	}

	hostCount := len(vmSpecs)
	err = host.CreateInBatch(vmSpecs, bus)
	if err != nil {
		return err
	}

	hosts := make([]types.Host, hostCount)
	for i := 0; i < hostCount; i++ {
		hosts[i] = <-bus
	}

}
