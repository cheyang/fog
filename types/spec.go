package types

import (
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
