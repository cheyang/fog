package cluster

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	provider_registry "github.com/cheyang/fog/cloudprovider/registry"
	"github.com/cheyang/fog/cluster/ansible"
	"github.com/cheyang/fog/cluster/deploy"
	"github.com/cheyang/fog/host"
	"github.com/cheyang/fog/persist"
	"github.com/cheyang/fog/types"
	"github.com/cheyang/fog/util"
	"github.com/cheyang/fog/util/dump"
)

func Bootstrap(spec types.Spec) error {

	err := types.Validate(spec)
	if err != nil {
		return err
	}

	logrus.Infof("spec: %+v", spec)

	//register core dump tool
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

	for _, host := range hosts {
		if host.Err != nil {
			return host.Err
		}
	}
	cp := provider_registry.GetProvider(spec.CloudDriverName, spec.ClusterType)
	if cp != nil {
		cp.SetHosts(hosts)
		cp.Configure() // configure IaaS
	}

	var deployer deploy.Deployer
	deployer, err = ansible.NewDeployer(spec.Name)
	if err != nil {
		return err
	}
	deployer.SetHosts(hosts)
	if len(spec.Run) > 0 {
		deployer.SetCommander(spec.Run)
	} else {
		deployer.SetCommander(spec.DockerRun)
	}

	err = deployer.Run()
	if err != nil {
		return err
	}

	name := spec.Name
	storePath, err := util.GetStorePath(name)
	if err != nil {
		return err
	}

	if _, err := os.Stat(storePath); os.IsNotExist(err) {
		return fmt.Errorf("Failed to find the storage of cluster %s in %s",
			name,
			storePath)
	}
	storage := persist.NewFilestore(storePath)

	for _, host := range hosts {
		host.Deployed = true
		err = storage.Save(&host)
		logrus.WithError(err).Warningf("failed to save %s", host.Name)
	}

	return nil
}
