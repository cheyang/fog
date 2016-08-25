package types

import (
	"github.com/docker/engine-api/types/container"
	"github.com/docker/machine/libmachine/state"
)

type Spec struct {
	Name            string                 `json:"ClusterName"` // The name of the cluster
	ClusterType     string                 `json:"ClusterType"` // The type of cluster
	VMSpecs         []VMSpec               `json:"Vmspecs"`
	Properties      map[string]interface{} `json:"Global,omitempty"`
	Run             container.Config       `json:"Run"`
	CloudDriverName string                 `json:"Driver"`
	Update          bool                   `json:"Update"` // Update an exist cluster?
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
