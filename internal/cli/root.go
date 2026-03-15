package cli

import (
	"context"
	"fmt"
	"os"

	renderscreenshot "github.com/Render-Screenshot/rs-go"

	"github.com/Render-Screenshot/rs-cli/internal/config"
	"github.com/Render-Screenshot/rs-cli/internal/output"
	"github.com/spf13/cobra"
)

// Version is set at build time via ldflags.
var Version = "dev"

var (
	// Global flags
	flagAPIKey  string
	flagJSON    bool
	flagQuiet   bool
	flagVerbose bool

	printer *output.Printer
)

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rs",
		Short: "RenderScreenshot CLI — capture web screenshots from the command line",
		Long: `rs is the command-line interface for the RenderScreenshot API.

Capture web screenshots, generate signed URLs, manage cache,
and preview page metadata — all from your terminal.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			printer = output.New(flagJSON, flagQuiet, flagVerbose)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		Version: Version,
	}

	cmd.SetVersionTemplate("rs version {{.Version}}\n")

	// Global persistent flags
	pf := cmd.PersistentFlags()
	pf.StringVar(&flagAPIKey, "api-key", "", "API key (overrides RS_API_KEY env and config)")
	pf.BoolVar(&flagJSON, "json", false, "output as JSON")
	pf.BoolVar(&flagQuiet, "quiet", false, "suppress progress output")
	pf.BoolVar(&flagVerbose, "verbose", false, "show detailed request/response info")

	// Register subcommands
	cmd.AddCommand(
		newTakeCmd(),
		newBatchCmd(),
		newSignedURLCmd(),
		newPreviewCmd(),
		newCacheCmd(),
		newPresetsCmd(),
		newDevicesCmd(),
		newWhoamiCmd(),
		newConfigCmd(),
		newLoginCmd(),
		newLogoutCmd(),
		newVersionCmd(),
	)

	return cmd
}

// Execute runs the root command.
func Execute() error {
	return newRootCmd().ExecuteContext(context.Background())
}

// newClient creates an SDK client from resolved credentials.
// Returns an error with exit code 2 if no API key is available.
func newClient() (*renderscreenshot.Client, error) {
	apiKey := config.ResolveAPIKey(flagAPIKey)
	if apiKey == "" {
		return nil, fmt.Errorf("no API key configured. Run 'rs login' or set RS_API_KEY")
	}

	opts := []renderscreenshot.Option{}

	if baseURL := os.Getenv("RS_BASE_URL"); baseURL != "" {
		opts = append(opts, renderscreenshot.WithBaseURL(baseURL))
	}

	pubKey, secKey := config.ResolveSigningKeys()
	if pubKey != "" {
		opts = append(opts, renderscreenshot.WithPublicKeyID(pubKey))
	}
	if secKey != "" {
		opts = append(opts, renderscreenshot.WithSigningKey(secKey))
	}

	client, err := renderscreenshot.New(apiKey, opts...)
	if err != nil {
		return nil, fmt.Errorf("creating client: %w", err)
	}
	return client, nil
}

// getPrinter returns the global printer, initializing if needed.
func getPrinter() *output.Printer {
	if printer == nil {
		printer = output.New(flagJSON, flagQuiet, flagVerbose)
	}
	return printer
}
