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
	err       error
}

func (this *HostHandler) create() {

	log.Infof("Host info %s: %+v ", this.Name, this.VMSpec)

	myHost := types.Host{}

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
		this.err = mcnerror.ErrDuringPreCreate{
			Cause: err,
		}
	}

	// create
	if this.err == nil {
		this.err = this.Driver.Create()
		if this.err != nil {
			log.Warnf("Err %s in creating machine %s\n", this.err.Error(), this.Name)
		} else {
			log.Infof("Creating machine for %s...\n", this.Name)
		}
	}

	// wait for
	if this.err == nil {
		log.Infof("Waiting for machine to be running, this may take a few minutes %s...\n", this.Name)
		this.err = mcnutils.WaitFor(drivers.MachineInState(this.Driver, state.Running))
		if this.err != nil {
			log.Warnf("Err %s in waiting machine %s\n", this.err.Error(), this.Name)
		}
	}

	if this.err == nil {
		log.Infof("Detecting operating system of created instance %s...\n", this.Name)
		_, err := provision.DetectProvisioner(this.Driver)
		if err != nil {
			this.err = fmt.Errorf("Error detecting OS: %s", err)
			log.Warnf("Error detecting OS: %s\n", err)
		}
	}

	if this.err == nil {
		myHost.SSHUserName = this.Driver.GetSSHUsername()

		myHost.Roles = this.VMSpec.Roles
		myHost.Name = this.Name
		myHost.SSHKeyPath = this.Driver.GetSSHKeyPath()
		myHost.SSHHostname, this.err = this.Driver.GetSSHHostname()
		myHost.Driver = this.Driver

		if this.err == nil {
			myHost.State, this.err = this.Driver.GetState()
		} else {
			myHost.Err = this.err
			log.Warnf("Failed to create host %s: %s\n", this.Name, myHost.Err)
		}

		if this.err == nil {
			myHost.SSHPort, this.err = this.Driver.GetSSHPort()
		} else {
			myHost.Err = this.err
			log.Warnf("Failed to create host %s: %s\n", this.Name, myHost.Err)
		}

		if this.err != nil {
			myHost.Err = this.err
			log.Warnf("Failed to create host %s: %s\n", this.Name, myHost.Err)
		}

	} else {
		myHost.Err = this.err
		log.Warnf("Failed to create host %s: %s\n", this.Name, myHost.Err)
	}

	// put myHost the createBus
	this.createBus <- myHost

	log.Infof("Finished creating host %s\n", this.Name)
}
