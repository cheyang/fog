package cluster

import (
	"fmt"
	"testing"

	"github.com/cheyang/fog/types"
	"github.com/stretchr/testify/assert"
)

func TestBuildScaleList(t *testing.T) {
	hostList := buildMockHostList()
	desiredMap := buildDesireMap()
	spec := buildSpec()

	runningMap := make(map[string]types.VMSpec)
	for _, vmSpec := range spec.VMSpecs {
		runningMap[vmSpec.Name] = vmSpec
	}
	currentHostMap := buildcurrentHostMap(hostList, runningMap)
	desiredHostMap := make(map[string]desiredConfig)
	for name, _ := range desiredMap {
		if _, found := runningMap[name]; !found {
			assert.Equal(t, false, found, "Not found the spec")
		}

		desiredHostMap[name] = desiredConfig{
			vmSpec:    runningMap[name],
			instances: desiredMap[name],
		}
	}
	toRemoveHosts, toCreateHostSpecs, err := buildScaleList(currentHostMap, desiredHostMap)

	assert.NoError(t, err)

	fmt.Println(toRemoveHosts)
	fmt.Println(toCreateHostSpecs)
}

func buildMockHostList() []*types.Host {
	return []*types.Host{
		&types.Host{
			Name: "master-2",
		},
		&types.Host{
			Name: "master-1",
		},
		&types.Host{
			Name: "master-0",
		},
		&types.Host{
			Name: "slave-1",
		},
		&types.Host{
			Name: "slave-0",
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
		"registry": 0,
	}
}
