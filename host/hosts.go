package host

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cheyang/fog/types"
)

// step 2
func CreateInBatch(vmConfigs []types.VMSpec, driverName string, hostBus chan<- Host) (count int, err error) {

	if driverName == "" {
		return count, fmt.Errorf("driver name is not specified.")
	}

	// make working directory
	pwd, err := os.Getwd()
	if err != nil {
		return count, err
	}
	t := time.Now()
	timestamp := fmt.Sprint(t.Format("20060102150405"))
	fogDir := filepath.Join(pwd, ".fog", timestamp)
	storePath := filepath.Join(pwd, ".fog", "latest")
	os.Remove(storePath)
	err = os.Symlink(fogDir, storePath)
	if err != nil {
		return count, err
	}

	// vmConfigs, err := BuildHostConfigs(spec)
	// if err == nil {
	// 	return count, err
	// }

	for _, vm := range vmConfigs {
		// go create(vm, driverName, hostBus)
		driver, err := initDrivers(driverName, vm, storePath)
		if err != nil {
			return count, err
		}

		h := &HostCreator{
			d:   driver,
			h:   vm,
			bus: hostBus,
			err: nil,
		}

		go h.create()
	}

	return len(vmConfigs), nil

}

// step 1
func BuildHostConfigs(specs types.Spec) (vmConfigs []types.VMSpec, err error) {

	dup := make(map[string]bool)
	for _, spec := range specs.VMSpecs {

		if spec.Name == "" {
			return vmConfigs, fmt.Errorf("Please specify the name")
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
				vmConfigs = append(vmConfigs, vm)
			}
		} else {
			vm := spec
			vm.Properties = mergeProperties(spec.Properties, spec.Properties)
			vmConfigs = append(vmConfigs, vm)
		}
	}

	return vmConfigs, nil
}

func mergeProperties(global, current map[string]interface{}) map[string]interface{} {

	for k, v := range current {
		global[k] = v
	}

	return global
}
