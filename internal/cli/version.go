package cli

import (
	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Run: func(cmd *cobra.Command, args []string) {
			p := getPrinter()
			if flagJSON {
				p.PrintJSON(map[string]string{"version": Version})
				return
			}
			p.Print("rs version %s\n", Version)
		},
	}
}
