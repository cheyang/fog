package cluster

import (
	"fmt"
	"testing"

	"github.com/cheyang/fog/types"
	"github.com/stretchr/testify/assert"
)

func TestbuildScaleList(t *testing.T) {
	hostList := buildMockHostList()
	desireMap := buildDesireMap()
	spec := buildSpec()

	runningMap := make(map[string]types.VMSpec)
	for _, vmSpec := range spec.VMSpecs {
		runningMap[vmSpec.Name] = vmSpec
	}
	currentHostMap := buildMockHostList(hostList, runningMap)
	desiredHostMap := make(map[string]desiredConfig)
	for name, _ := range desiredMap {
		if _, found := runningMap[name]; !found {
			assert.Error(t, err)
		}

		desiredHostMap[name] = desiredConfig{
			vmSpec:    runningMap[name],
			instances: desiredMap[name],
		}
	}
	buildScaleList(currentHostMap, desiredHostMap)

	fmt.Println(toRemoveHosts)
	fmt.Println(toCreateHostSpecs)
}

func buildMockHostList() []*types.Host {
	return []*types.Host{
		&types.Host{
			Name: "master-0",
		},
		&types.Host{
			Name: "master-1",
		},
		&types.Host{
			Name: "master-2",
		},
		&types.Host{
			Name: "slave-0",
		},
		&types.Host{
			Name: "slave-1",
		},
		&types.Host{
			Name: "slave-2",
		},
		&types.Host{
			Name: "registry-0",
		},
		&types.Host{
			Name: "registry-1",
		},
		&types.Host{
			Name: "registry-2",
		},
	}
}

func buildSpec() *types.Spec {
	return &types.Spec{
		VMSpecs: []types.VMSpec{
			types.VMSpec{
				Name: "master",
			},
			types.VMSpec{
				Name: "slave",
			},
			types.VMSpec{
				Name: "registry",
			},
		},
	}
}

func buildDesireMap() map[string]int {
	return map[string]int{
		"master":   1,
		"slave":    5,
		"registry": 3,
	}
}
