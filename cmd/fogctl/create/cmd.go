package create

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"

	"github.com/cheyang/fog/cluster"
	"github.com/cheyang/fog/types"
	"github.com/cheyang/fog/util/yaml"
	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:   "create",
		Short: "Create a cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return errors.New("create command takes no arguments")
			}

			if !cmd.Flags().Changed("config-file") {
				return errors.New("--config-file are mandatory")
			}

			configFile, err := cmd.Flags().GetString("config-file")
			if err != nil {
				return err
			}

			// read and parse the config file
			spec := types.Spec{}
			if _, err := os.Stat(configFile); os.IsNotExist(err) {
				return err
			}
			data, err := ioutil.ReadFile(configFile)
			if err != nil {
				return err
			}
			decoder := yaml.NewYAMLToJSONDecoder(bytes.NewReader(data))
			err = decoder.Decode(&spec)
			if err != nil {
				return err
			}

			retry, err := cmd.Flags().GetBool("retry")
			if err != nil {
				return err
			}
			spec.Update = retry

			return cluster.Bootstrap(spec)
		},
	}
)

func init() {
	Cmd.Flags().StringP("config-file", "f", "", "The config file.")
	Cmd.Flags().BoolP("retry", "r", false, "retry to create cluster.")
}
