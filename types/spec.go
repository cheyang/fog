package types

import (
	"github.com/docker/machine/libmachine/state"
	docker "github.com/fsouza/go-dockerclient"
)

type Spec struct {
	VMSpecs    []VMSpec               `json:"vmspec"`
	Properties map[string]interface{} `json:"Global,omitempty"`
	Run        docker.Container       `json:"run"`
}

type VMSpec struct {
	Name       string                 `json:"name"`
	Roles      []string               `json:"roles"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	Driver     string                 `json:"driver"`
	Instances  int                    `json:"instances,omitempty"`
}

type Host struct {
	Err         error
	MachineName string
	SSHUserName string
	SSHPort     string
	SSHHostname string
	SSHKeyPath  string
	Roles       []string
	State       state.State
}
