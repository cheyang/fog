package host

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/cheyang/fog/types"
	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/mcnerror"
	"github.com/docker/machine/libmachine/mcnutils"
	"github.com/docker/machine/libmachine/provision"
	"github.com/docker/machine/libmachine/state"
)

type HostCreator struct {
	d   drivers.Driver
	h   types.VMSpec
	bus chan<- types.Host
	err error
}

func (this *HostCreator) create() {

	myHost := types.Host{}

	// pre-check
	log.Infof("Running pre-create checks for  %s...\n", this.h.Name)

	if err := this.d.PreCreateCheck(); err != nil {
		this.err = mcnerror.ErrDuringPreCreate{
			Cause: err,
		}
	}

	// create
	if this.err == nil {
		this.err = this.d.Create()
		if this.err != nil {
			log.Warnf("Err %s in creating machine %s\n", this.err.Error(), this.h.Name)
		} else {
			log.Infof("Creating machine for %s...\n", this.h.Name)
		}
	}

	// wait for
	if this.err == nil {
		log.Infof("Waiting for machine to be running, this may take a few minutes %s...\n", this.h.Name)
		this.err = mcnutils.WaitFor(drivers.MachineInState(this.d, state.Running))
		if this.err != nil {
			log.Warnf("Err %s in waiting machine %s\n", this.err.Error(), this.h.Name)
		}
	}

	if this.err == nil {
		log.Infof("Detecting operating system of created instance %s...\n", this.h.Name)
		_, err := provision.DetectProvisioner(this.d)
		if err != nil {
			this.err = fmt.Errorf("Error detecting OS: %s", err)
			log.Warnf("Error detecting OS: %s\n", err)
		}
	}

	if this.err == nil {
		myHost.SSHUserName = this.d.GetSSHUsername()

		myHost.Roles = this.h.Roles
		myHost.MachineName = this.d.GetMachineName()
		myHost.SSHKeyPath = this.d.GetSSHKeyPath()
		myHost.SSHHostname, this.err = this.d.GetSSHHostname()

		if this.err == nil {
			myHost.State, this.err = this.d.GetState()
		} else {
			myHost.Err = this.err
			log.Warnf("Failed to create host %s: %s\n", this.h.Name, myHost.Err)
		}

		if this.err == nil {
			myHost.SSHPort, this.err = this.d.GetSSHPort()
		} else {
			myHost.Err = this.err
			log.Warnf("Failed to create host %s: %s\n", this.h.Name, myHost.Err)
		}

		if this.err != nil {
			myHost.Err = this.err
			log.Warnf("Failed to create host %s: %s\n", this.h.Name, myHost.Err)
		}

	} else {
		myHost.Err = this.err
		log.Warnf("Failed to create host %s: %s\n", this.h.Name, myHost.Err)
	}

	// put myHost the bus
	this.bus <- myHost

	log.Infof("Finished creating host %s\n", this.h.Name)
}
