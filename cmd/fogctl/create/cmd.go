package create

import (
	"errors"
	"fmt"
	"time"

	"github.com/cheyang/fog/cluster"
	"github.com/cheyang/fog/types"
	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:   "create",
		Short: "Create a cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			start := time.Now()
			defer func() {
				end := time.Now()
				fmt.Printf("Creating a cluster takes about %v minutes.\n", end.Sub(start).Minutes())
			}()
			//load spec
			if !cmd.Flags().Changed("config-file") {
				return errors.New("--config-file are mandatory")
			}
			configFile, err := cmd.Flags().GetString("config-file")
			if err != nil {
				return err
			}
			spec, err := types.LoadSpec(configFile)

			//set retry
			retry, err := cmd.Flags().GetBool("retry")
			if err != nil {
				return err
			}
			spec.Update = retry

			hosts, err := cluster.Bootstrap(spec)

			if err != nil {
				return err
			}

			for _, host := range hosts {
				fmt.Printf("Name: %s, Ipaddress %s, SSH key path %s with roles [%v]\n",
					host.Name,
					host.SSHHostname,
					host.SSHKeyPath,
					host.Roles)
			}
		},
	}
)

func init() {
	Cmd.Flags().StringP("config-file", "f", "", "The config file.")
	Cmd.Flags().BoolP("retry", "r", false, "retry to create cluster.")
}
