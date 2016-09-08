package scale

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cheyang/fog/cluster"
	"github.com/cheyang/fog/types"
	"github.com/cheyang/fog/util"
	"github.com/spf13/cobra"
)

var (
	scaleoutCmd = &cobra.Command{
		Use:   "scale out",
		Short: "scale out a cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			start := time.Now()
			defer func() {
				end := time.Now()
				fmt.Printf("scaling out a cluster takes about %v minutes.\n", end.Sub(start).Minutes())
			}()

			if len(args) == 0 {
				return errors.New("scale out command takes no arguments")
			}
			name := args[len(args)-1]
			storage, err := util.GetStorage(name)

			//load spec
			flags := cmd.Flags()
			if !flags.Changed("config-file") {
				return errors.New("--config-file are mandatory")
			}
			configFile, err := flags.GetString("config-file")
			if err != nil {
				return err
			}
			spec, err := types.LoadSpec(configFile)

			// build required role map
			roleString, err := flags.GetString("with-roles")
			if err != nil {
				return err
			}
			roles := strings.Split(roleString, ",")
			roleMap := make(map[string]bool)
			for _, role := range roles {
				roleMap[role] = true
			}

			return cluster.Scaleout(storage, spec, roleMap)
		},
	}
)

func init() {
	flags := Cmd.Flags()
	flags.StringP("config-file", "f", "", "The config file.")
	flags.StringP("with-roles", "w", "", "If you need the inventory file also includes role")
}
