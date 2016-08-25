package main

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/cheyang/fog/cluster"
	"github.com/cheyang/fog/types"
	"github.com/cheyang/fog/util/yaml"
	"github.com/spf13/cobra"
)

/**

 */

func main() {
	if err := mainCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}

var mainCmd = &cobra.Command{
	Use:          os.Args[0],
	Short:        "Run the raw control command!",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
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

		return cluster.Bootstrap(spec)

	},
}

func init() {
	mainCmd.Flags().StringP("config-file", "c", "", "The config file")
	mainCmd.Flags().StringP("volume", "v", "", "[host-src:]container-dest[:<options>]: Bind mount a volume.")
}
