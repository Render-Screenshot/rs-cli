package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newDevicesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "devices",
		Short: "List available device presets",
		RunE: func(cmd *cobra.Command, args []string) error {
			p := getPrinter()

			client, err := newClient()
			if err != nil {
				p.Error("%s", err)
				os.Exit(2)
			}

			devices, err := client.Devices(cmd.Context())
			if err != nil {
				return fmt.Errorf("fetching devices: %w", err)
			}

			if flagJSON {
				return p.PrintJSON(devices)
			}

			headers := []string{"ID", "NAME", "WIDTH", "HEIGHT"}
			rows := make([][]string, len(devices))
			for i, d := range devices {
				rows[i] = []string{
					d.ID,
					d.Name,
					fmt.Sprintf("%d", d.Width),
					fmt.Sprintf("%d", d.Height),
				}
			}
			p.Table(headers, rows)
			return nil
		},
	}
}
