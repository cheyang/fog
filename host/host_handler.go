package host

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/cheyang/fog/persist"
	"github.com/cheyang/fog/types"
	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/mcnerror"
	"github.com/docker/machine/libmachine/mcnutils"
	"github.com/docker/machine/libmachine/provision"
	"github.com/docker/machine/libmachine/state"
)

type HostHandler struct {
	Name      string
	Driver    drivers.Driver
	VMSpec    types.VMSpec
	createBus chan<- types.Host
}

func (this *HostHandler) createOrGet() {

	log.Infof("Host info %s: %+v ", this.Name, this.VMSpec)
	// store to path
	storage := persist.NewFilestore(storePath)
	exist, err := storage.Exists(this.Name)
	var host *types.Host

	if err == nil {
		if exist {
			host = this.get(storage)
		} else {
			host = this.create(storage)
		}
	} else {
		host := &types.Host{
			Name: this.Name,
			Err:  err,
		}
	}

	// put host the createBus
	this.createBus <- *host

	log.Infof("Finished creating host %s\n", this.Name)
}

func (this *HostHandler) create(s persist.Store) *types.Host {
	log.Infof("start creating host %s from storage", this.Name)
	host := &types.Host{
		Name:       this.Name,
		Roles:      this.VMSpec.Roles,
		Driver:     this.Driver,
		DriverName: this.VMSpec.CloudDriverName,
		VMSpec:     this.VMSpec,
	}
	s.CreateStorePath(this.Name)

	defer func() {
		err := s.Save(host)
		if err != nil {
			log.Warnf("Error in saving to file store %s: %s ", this.Name, err)
		}
	}()

	// pre-check
	log.Infof("Running pre-create checks for  %s...\n", this.Name)
	if err := this.Driver.PreCreateCheck(); err != nil {
		host.Err = mcnerror.ErrDuringPreCreate{
			Cause: err,
		}
	}

	// create
	if host.Err == nil {
		host.Err = this.Driver.Create()
		if host.Err != nil {
			log.Warnf("Err %s in creating machine %s\n", host.Err.Error(), this.Name)
		} else {
			log.Infof("Creating machine for %s...\n", this.Name)
		}
	}

	// wait for
	if host.Err == nil {
		log.Infof("Waiting for machine to be running, this may take a few minutes %s...\n", this.Name)
		host.Err = mcnutils.WaitFor(drivers.MachineInState(this.Driver, state.Running))
		if host.Err != nil {
			log.Warnf("Err %s in waiting machine %s\n", host.Err.Error(), this.Name)
		}
	}

	if host.Err == nil {
		log.Infof("Detecting operating system of created instance %s...\n", this.Name)
		_, err := provision.DetectProvisioner(this.Driver)
		if err != nil {
			host.Err = fmt.Errorf("Error detecting OS: %s", err)
			log.Warnf("Error detecting OS: %s\n", err)
		}
	}

	if host.Err == nil {
		host.SSHUserName = this.Driver.GetSSHUsername()
		host.SSHKeyPath = this.Driver.GetSSHKeyPath()
		host.SSHHostname, host.Err = this.Driver.GetSSHHostname()

		if host.Err == nil {
			host.State, host.Err = this.Driver.GetState()
		} else {
			host.Err = host.Err
			log.Warnf("Failed to create host %s: %s\n", this.Name, host.Err)
		}

		if host.Err == nil {
			host.SSHPort, host.Err = this.Driver.GetSSHPort()
		} else {
			host.Err = host.Err
			log.Warnf("Failed to create host %s: %s\n", this.Name, host.Err)
		}

		if host.Err != nil {
			log.Warnf("Failed to create host %s: %s\n", this.Name, host.Err)
		}

	} else {

		log.Warnf("Failed to create host %s: %s\n", this.Name, host.Err)
	}

	return host
}

func (this *HostHandler) get(s persist.Store) *types.Host {
	log.Infof("start getting host %s from storage", this.Name)
	host, err := s.Load(this.Name)
	if err != nil {
		host.Err = err
	}

	if host.SSHHostname == "" {
		host.Err = fmt.Errorf("SSHHostName is empty")
	}
	if host.SSHKeyPath == "" {
		host.Err = fmt.Errorf("SSHKeyPath is empty")
	}
	if host.SSHUserName == "" {
		host.Err = fmt.Errorf("SSHUserName is empty")
	}

	log.Infof("Waiting for machine to be running, this may take a few minutes %s...\n", host.Name)
	host.Err = mcnutils.WaitFor(drivers.MachineInState(host.Driver, state.Running))
	if host.Err != nil {
		log.Warnf("Err %s in waiting machine %s\n", host.Err.Error(), host.Name)
	}

	if host.Err == nil {
		log.Infof("Detecting operating system of created instance %s...\n", host.Name)
		_, err := provision.DetectProvisioner(host.Driver)
		if err != nil {
			host.Err = fmt.Errorf("Error detecting OS: %s", err)
			log.Warnf("Error detecting OS: %s\n", err)
		}
	}

	return host
}
