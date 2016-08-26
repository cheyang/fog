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

func (this *HostHandler) create() {

	log.Infof("Host info %s: %+v ", this.Name, this.VMSpec)

	myHost := types.Host{}
	myHost.Name = this.Name
	myHost.Roles = this.VMSpec.Roles
	myHost.Driver = this.Driver

	// store to path
	storage := persist.NewFilestore(storePath)
	storage.CreateStorePath(this.Name)

	defer func() {
		err := storage.Save(&myHost)
		if err != nil {
			log.Warnf("Error in saving to file store %s: %s ", this.Name, err)
		}
	}()

	// pre-check
	log.Infof("Running pre-create checks for  %s...\n", this.Name)

	if err := this.Driver.PreCreateCheck(); err != nil {
		myHost.Err = mcnerror.ErrDuringPreCreate{
			Cause: err,
		}
	}

	// create
	if myHost.Err == nil {
		myHost.Err = this.Driver.Create()
		if myHost.Err != nil {
			log.Warnf("Err %s in creating machine %s\n", myHost.Err.Error(), this.Name)
		} else {
			log.Infof("Creating machine for %s...\n", this.Name)
		}
	}

	// wait for
	if myHost.Err == nil {
		log.Infof("Waiting for machine to be running, this may take a few minutes %s...\n", this.Name)
		myHost.Err = mcnutils.WaitFor(drivers.MachineInState(this.Driver, state.Running))
		if myHost.Err != nil {
			log.Warnf("Err %s in waiting machine %s\n", myHost.Err.Error(), this.Name)
		}
	}

	if myHost.Err == nil {
		log.Infof("Detecting operating system of created instance %s...\n", this.Name)
		_, err := provision.DetectProvisioner(this.Driver)
		if err != nil {
			myHost.Err = fmt.Errorf("Error detecting OS: %s", err)
			log.Warnf("Error detecting OS: %s\n", err)
		}
	}

	if myHost.Err == nil {
		myHost.SSHUserName = this.Driver.GetSSHUsername()
		myHost.SSHKeyPath = this.Driver.GetSSHKeyPath()
		myHost.SSHHostname, myHost.Err = this.Driver.GetSSHHostname()

		if myHost.Err == nil {
			myHost.State, myHost.Err = this.Driver.GetState()
		} else {
			myHost.Err = myHost.Err
			log.Warnf("Failed to create host %s: %s\n", this.Name, myHost.Err)
		}

		if myHost.Err == nil {
			myHost.SSHPort, myHost.Err = this.Driver.GetSSHPort()
		} else {
			myHost.Err = myHost.Err
			log.Warnf("Failed to create host %s: %s\n", this.Name, myHost.Err)
		}

		if myHost.Err != nil {
			log.Warnf("Failed to create host %s: %s\n", this.Name, myHost.Err)
		}

	} else {

		log.Warnf("Failed to create host %s: %s\n", this.Name, myHost.Err)
	}

	// put myHost the createBus
	this.createBus <- myHost

	log.Infof("Finished creating host %s\n", this.Name)
}
