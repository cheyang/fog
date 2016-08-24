package types

import (
	"github.com/docker/engine-api/types/container"
	"github.com/docker/machine/libmachine/state"
)

type Spec struct {
	VMSpecs         []VMSpec               `json:"Vmspecs"`
	Properties      map[string]interface{} `json:"Global,omitempty"`
	Run             container.Config       `json:"Run"`
	ClusterType     string                 `json:"ClusterType"`
	CloudDriverName string                 `json:"Driver"`
}

type VMSpec struct {
	Name            string                 `json:"Name"`
	Roles           []string               `json:"Roles"`
	Properties      map[string]interface{} `json:"Properties,omitempty"`
	CloudDriverName string                 `json:"Driver"`
	Instances       int                    `json:"Instances,omitempty"`
}

type Host struct {
	Err         error
	MachineName string
	SSHUserName string
	SSHPort     int
	SSHHostname string
	SSHKeyPath  string
	Roles       []string
	State       state.State
	VMSpec
}
