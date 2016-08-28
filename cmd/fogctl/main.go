package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/cheyang/fog/cmd/fogctl/create"
	"github.com/cheyang/fog/cmd/fogctl/list"
	"github.com/cheyang/fog/cmd/fogctl/remove"
	"github.com/cheyang/fog/cmd/fogctl/scale"
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
	Short:        "control a cluster!",
	SilenceUsage: true,
}

func init() {
	mainCmd.PersistentFlags().StringP("config-file", "f", "", "The config file")

	mainCmd.AddCommand(
		create.Cmd,
		scale.Cmd,
		remove.Cmd,
		list.Cmd,
	)
}
