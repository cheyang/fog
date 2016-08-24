package cluster

import (
	"github.com/cheyang/fog/cloudprovider"
	"github.com/cheyang/fog/host"
	"github.com/cheyang/fog/types"
	"github.com/cheyang/fog/util/dump"
)

func Bootstrap(spec types.Spec) error {

	err := types.Validate(spec)
	if err != nil {
		return err
	}

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

	cp := cloudprovider.InitProivder(spec)
	if cp != nil {
		cp.SetHosts(hosts)
		cp.Configure() // configure infrastructure
	}

	var deployer Deployer = &ansibleDeployer{}
	deployer.SetHosts(hosts)
	deployer.Run()

	return nil
}
