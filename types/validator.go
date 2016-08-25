package types

import "fmt"

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

	return nil
}
