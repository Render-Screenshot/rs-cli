package cli

import (
	"fmt"

	"github.com/Render-Screenshot/rs-cli/internal/config"
	"github.com/Render-Screenshot/rs-cli/internal/output"
	"github.com/spf13/cobra"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "View and manage CLI configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Default: show all config
			return showConfig()
		},
	}

	cmd.AddCommand(
		newConfigGetCmd(),
		newConfigSetCmd(),
		newConfigPathCmd(),
		newConfigShowCmd(),
	)

	return cmd
}

func newConfigShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show all configuration values",
		RunE: func(cmd *cobra.Command, args []string) error {
			return showConfig()
		},
	}
}

func showConfig() error {
	p := getPrinter()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	if flagJSON {
		return p.PrintJSON(map[string]interface{}{
			"api_key":       output.MaskKey(cfg.APIKey),
			"public_key_id": cfg.PublicKeyID,
			"secret_key":    output.MaskKey(cfg.SecretKey),
			"path":          config.Path(),
		})
	}

	p.Print("api_key:       %s\n", output.MaskKey(cfg.APIKey))
	p.Print("public_key_id: %s\n", cfg.PublicKeyID)
	p.Print("secret_key:    %s\n", output.MaskKey(cfg.SecretKey))
	p.Print("path:          %s\n", config.Path())

	return nil
}

func newConfigGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <key>",
		Short: "Get a configuration value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			p := getPrinter()

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			val, err := cfg.Get(args[0])
			if err != nil {
				return err
			}

			// Mask sensitive values
			switch args[0] {
			case "api_key", "secret_key":
				val = output.MaskKey(val)
			}

			p.Println(val)
			return nil
		},
	}
}

func newConfigSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			p := getPrinter()

			cfg, err := config.Load()
			if err != nil {
				cfg = &config.Config{}
			}

			if err := cfg.Set(args[0], args[1]); err != nil {
				return err
			}

			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("saving config: %w", err)
			}

			p.Print("Set %s\n", args[0])
			return nil
		},
	}
}

func newConfigPathCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "Show config file path",
		Run: func(cmd *cobra.Command, args []string) {
			p := getPrinter()
			p.Println(config.Path())
		},
	}
}
