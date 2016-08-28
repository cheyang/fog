package types

import (
	"fmt"
	"reflect"

	"github.com/Sirupsen/logrus"
)

func Validate(specs Spec) error {

	if specs.CloudDriverName == "" {
		return fmt.Errorf("cloud driver name is not specified")
	}

	if specs.ClusterType == "" {
		return fmt.Errorf("cluster type is not specified")
	}

	if specs.Name == "" {
		return fmt.Errorf("cluster name is not specified")
	}

	if err := validateMap(specs.Properties); err != nil {
		return err
	}
	for _, vmSpec := range specs.VMSpecs {
		if err := validateMap(vmSpec.Properties); err != nil {
			return err
		}

		if vmSpec.Instances < 1 {
			return fmt.Errorf("the instances %d of %s is not specified", vmSpec.Instances, vmSpec.Name)
		}
	}

	if len(specs.Run) > 0 && specs.DockerRun != nil {
		return fmt.Errorf("DockerRun and Run can't be specified together")
	}

	if len(specs.Run) == 0 && specs.DockerRun == nil {
		return fmt.Errorf("DockerRun and Run can't be empty either")
	}

	return nil
}

func validateMap(props map[string]interface{}) error {
	for k, v := range props {
		logrus.Debugf("The type %s of value %s, and its key is ", reflect.TypeOf(v), reflect.ValueOf(v))
		switch v.(type) {
		case string:
		case []string:
		case int:
		case bool:
		case float64:
		default:
			return fmt.Errorf("The type %s of value %s, and its key is %s", reflect.TypeOf(v), reflect.ValueOf(v), k)
		}
	}

	return nil
}
