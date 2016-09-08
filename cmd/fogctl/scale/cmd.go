package scale

import "github.com/spf13/cobra"

var (
	Cmd = &cobra.Command{
		Use:   "scale",
		Short: "scale out/in a cluster",
	}
)

func init() {
	Cmd.AddCommand(
		scaleoutCmd,
		scaleinCmd,
	)
}
