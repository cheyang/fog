package host

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cheyang/fog/types"
)

var storePath string

// Create hosts in batch
func CreateInBatch(vmSpecs []types.VMSpec, hostBus chan<- types.Host) (err error) {

	// make working directory

	for _, vm := range vmSpecs {
		// go create(vm, driverName, hostBus)

		driverName := vm.CloudDriverName
		if driverName == "" {
			return fmt.Errorf("driver name is not specified.")
		}
		driver, err := initDrivers(driverName, vm, storePath)
		if err != nil {
			return err
		}

		h := &HostHandler{
			Name:      vm.Name,
			Driver:    driver,
			VMSpec:    vm,
			createBus: hostBus,
			err:       nil,
		}

		go h.create()
	}

	return nil

}

func createStorePath(specs types.Spec) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	storePath = filepath.Join(pwd, ".fog", specs.ClusterType, specs.Name)

	// if the dir exists and not update mode
	if _, err := os.Stat(storePath); !os.IsNotExist(err) {
		if !specs.Update {
			return fmt.Errorf("working dir %s is not clean, can't work in create mode", storePath)
		}
	}

	os.MkdirAll(path, perm)

}

// step 1
func BuildHostConfigs(specs types.Spec) (vmSpecs []types.VMSpec, err error) {

	if err := createStorePath(specs); err != nil {
		return vmSpecs, err
	}

	dup := make(map[string]bool)
	for _, spec := range specs.VMSpecs {

		if spec.Name == "" {
			return vmSpecs, fmt.Errorf("Please specify the name")
		}

		if _, found := dup[spec.Name]; found {
			return nil, fmt.Errorf("duplicate name %s in configuration file.", spec.Name)
		} else {
			dup[spec.Name] = true
		}

		// if the attribute 'instances' is not specified, set it as 1
		if spec.Instances == 0 {
			spec.Instances = 1
		}

		if spec.Instances > 1 {
			for i := 0; i < spec.Instances; i++ {
				vm := spec
				vm.Name = fmt.Sprintf("%s-%d", vm.Name, i)
				vm.Properties = mergeProperties(spec.Properties, spec.Properties)
				if len(vm.Roles) == 0 {
					return vmSpecs, fmt.Errorf("please specify the role of %s", spec.Name)
				}
				// Set common cloud driver name if not specified
				if vm.CloudDriverName == "" {
					vm.CloudDriverName = specs.CloudDriverName
				}
				vmSpecs = append(vmSpecs, vm)
			}
		} else {
			vm := spec
			vm.Properties = mergeProperties(spec.Properties, spec.Properties)
			if vm.CloudDriverName == "" {
				vm.CloudDriverName = specs.CloudDriverName
			}
			vmSpecs = append(vmSpecs, vm)
		}
	}

	return vmSpecs, nil
}

func mergeProperties(global, current map[string]interface{}) map[string]interface{} {

	for k, v := range current {
		global[k] = v
	}

	return global
}
