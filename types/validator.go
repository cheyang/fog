package types

import (
	"fmt"
	"reflect"
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

	for _, v := range specs.Properties {
		switch v.(type) {
		case string:
		case []string:
		case int:
		case bool:
		default:
			return fmt.Errorf("The type %s of value %s", reflect.TypeOf(v), reflect.ValueOf(v))
		}
	}

	return nil
}
