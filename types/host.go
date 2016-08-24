package types

import "github.com/docker/machine/libmachine/state"

type Host struct {
	Err         error
	SSHUserName string
	SSHPort     string
	SSHHostname string
	State       state.State
}
