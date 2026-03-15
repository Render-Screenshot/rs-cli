package cli

import (
	"fmt"

	"github.com/Render-Screenshot/rs-cli/internal/config"
	"github.com/spf13/cobra"
)

func newLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Remove stored credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			p := getPrinter()

			if err := config.Delete(); err != nil {
				return fmt.Errorf("removing credentials: %w", err)
			}

			p.Println("Credentials removed")
			return nil
		},
	}
}
