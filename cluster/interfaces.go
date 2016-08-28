package cluster

import "github.com/cheyang/fog/types"

type Deployer interface {
	SetHosts(hosts []types.Host)
	SetCommander(run interface{}) error
	Run() error
}
