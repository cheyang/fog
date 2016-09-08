package scale

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/cheyang/fog/cluster"
	"github.com/cheyang/fog/util"
	"github.com/spf13/cobra"
)

var (
	scaleinCmd = &cobra.Command{
		Use:   "scale in",
		Short: "scale in a cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("scale in command takes no arguments")
			}

			scaleInMap := make(map[string]int)
			var name string
			for i, arg := range args {
				if i == len(args)-1 {
					name = arg
					break
				}

				kv := strings.Split(arg, "=")

				if len(kv) == 2 {
					// scaleInMap[kv[0]]
					value, err := strconv.Atoi(kv[1])
					if err != nil {
						return err
					}
					key := kv[0]
					scaleInMap[key] = value
				} else {
					return fmt.Errorf("the format of %s is not correct!", arg)
				}
			}
			storage, err := util.GetStorage(name)
			if err != nil {
				return err
			}

			return cluster.Scalein(storage, scaleInMap)
		},
	}
)
