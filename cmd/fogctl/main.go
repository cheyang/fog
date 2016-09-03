package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/cheyang/fog/cmd/fogctl/create"
	"github.com/cheyang/fog/cmd/fogctl/list"
	"github.com/cheyang/fog/cmd/fogctl/remove"
	"github.com/cheyang/fog/cmd/fogctl/scale"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/swarmkit/log"
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
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {

		flag, err := cmd.PersistentFlags().GetString("log-level")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		level, err := logrus.ParseLevel(flag)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		logrus.SetLevel(level)

		debugFlag, err := cmd.PersistentFlags().GetBool("debug-docker-machine")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		if debugFlag {
			log.SetDebug(true)
			fmt.Printf("Enable docker machine debug %b\n", debugFlag)
		}

	},
}

func init() {
	mainCmd.PersistentFlags().StringP("config-file", "f", "", "The config file")
	mainCmd.PersistentFlags().StringP("log-level", "l", "info", "Log level (options \"debug\", \"info\", \"warn\", \"error\", \"fatal\", \"panic\")")
	mainCmd.PersistentFlags().BoolP("debug-docker-machine", "D", false, "Debug the docker machine libraray")
	mainCmd.AddCommand(
		create.Cmd,
		scale.Cmd,
		remove.Cmd,
		list.Cmd,
	)
}
