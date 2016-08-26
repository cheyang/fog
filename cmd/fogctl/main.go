package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/cheyang/fog/cmd/fogctl/create"
	"github.com/cheyang/fog/cmd/fogctl/list"
	"github.com/cheyang/fog/cmd/fogctl/remove"
	"github.com/cheyang/fog/cmd/fogctl/update"
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
	mainCmd.PersistentFlags().StringP("config-file", "c", "", "The config file")
	mainCmd.PersistentFlags().StringP("volume", "v", "", "[host-src:]container-dest[:<options>]: Bind mount a volume.")

	mainCmd.AddCommand(
		create.Cmd,
		update.Cmd,
		remove.Cmd,
		list.Cmd,
	)
}
