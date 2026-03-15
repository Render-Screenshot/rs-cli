package cli

import (
	"fmt"
	"os"

	"github.com/Render-Screenshot/rs-cli/internal/config"
	"github.com/Render-Screenshot/rs-cli/internal/output"
	"github.com/spf13/cobra"
)

func newWhoamiCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "whoami",
		Short: "Show current authentication status",
		RunE: func(cmd *cobra.Command, args []string) error {
			p := getPrinter()

			client, err := newClient()
			if err != nil {
				p.Error("%s", err)
				os.Exit(2)
			}

			usage, err := client.Usage(cmd.Context())
			if err != nil {
				return fmt.Errorf("fetching account info: %w", err)
			}

			apiKey := config.ResolveAPIKey(flagAPIKey)

			if flagJSON {
				return p.PrintJSON(map[string]interface{}{
					"api_key":      output.MaskKey(apiKey),
					"credits":      usage.Credits,
					"used":         usage.Used,
					"remaining":    usage.Remaining,
					"period_start": usage.PeriodStart,
					"period_end":   usage.PeriodEnd,
					"config_path":  config.Path(),
				})
			}

			p.Print("API Key:      %s\n", output.MaskKey(apiKey))
			p.Print("Credits:      %d\n", usage.Credits)
			p.Print("Used:         %d\n", usage.Used)
			p.Print("Remaining:    %d\n", usage.Remaining)
			p.Print("Period:       %s to %s\n", usage.PeriodStart, usage.PeriodEnd)
			p.Print("Config:       %s\n", config.Path())

			return nil
		},
	}
}
