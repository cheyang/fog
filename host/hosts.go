package host

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cheyang/fog/types"
)

// Create hosts in batch
func CreateInBatch(vmSpecs []types.VMSpec, hostBus chan<- Host) (err error) {

	// make working directory
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	t := time.Now()
	timestamp := fmt.Sprint(t.Format("20060102150405"))
	fogDir := filepath.Join(pwd, ".fog", timestamp)
	storePath := filepath.Join(pwd, ".fog", "latest")
	os.Remove(storePath)
	err = os.Symlink(fogDir, storePath)
	if err != nil {
		return err
	}

	// vmSpecs, err := BuildHostConfigs(spec)
	// if err == nil {
	// 	return count, err
	// }

	for _, vm := range vmSpecs {
		// go create(vm, driverName, hostBus)

		driverName := vm.Driver
		if driverName == "" {
			return fmt.Errorf("driver name is not specified.")
		}
		driver, err := initDrivers(driverName, vm, storePath)
		if err != nil {
			return err
		}

		h := &HostCreator{
			d:   driver,
			h:   vm,
			bus: hostBus,
			err: nil,
		}

		go h.create()
	}

	return nil

}

// step 1
func BuildHostConfigs(specs types.Spec) (vmSpecs []types.VMSpec, err error) {

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
					return ni, fmt.Errorf("please specify the role of %s", spec.Name)
				}
				vmSpecs = append(vmSpecs, vm)
			}
		} else {
			vm := spec
			vm.Properties = mergeProperties(spec.Properties, spec.Properties)
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
