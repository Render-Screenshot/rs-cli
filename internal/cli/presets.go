package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newPresetsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "presets",
		Short: "List available screenshot presets",
		RunE: func(cmd *cobra.Command, args []string) error {
			p := getPrinter()

			client, err := newClient()
			if err != nil {
				p.Error("%s", err)
				os.Exit(2)
			}

			presets, err := client.Presets(cmd.Context())
			if err != nil {
				return fmt.Errorf("fetching presets: %w", err)
			}

			if flagJSON {
				return p.PrintJSON(presets)
			}

			headers := []string{"ID", "NAME", "WIDTH", "HEIGHT"}
			rows := make([][]string, len(presets))
			for i, preset := range presets {
				rows[i] = []string{
					preset.ID,
					preset.Name,
					fmt.Sprintf("%d", preset.Width),
					fmt.Sprintf("%d", preset.Height),
				}
			}
			p.Table(headers, rows)
			return nil
		},
	}
}
