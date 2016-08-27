package remove

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cheyang/fog/cluster"
	"github.com/cheyang/fog/persist"
	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:     "remove <name>",
		Short:   "Remove a cluster",
		Aliases: []string{"rm"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("cluster name missing")
			}

			if len(args) > 2 {
				return errors.New("remove command takes exactly 1 argument")
			}

			for _, name := range args {
				pwd, err := os.Getwd()
				if err != nil {
					return err
				}
				storePath := filepath.Join(pwd, ".fog", name)

				if _, err := os.Stat(storePath); os.IsNotExist(err) {
					return fmt.Errorf("Failed to find the storage of cluster %s in %s",
						name,
						filepath.Join(pwd, ".fog"))
				}
				storage := persist.NewFilestore(storePath)

				err = cluster.Remove(storePath, storage)

				if err != nil {
					return err
				}
			}

			return nil
		},
	}
)
