package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	renderscreenshot "github.com/Render-Screenshot/rs-go"

	"github.com/Render-Screenshot/rs-cli/internal/config"
	"github.com/Render-Screenshot/rs-cli/internal/output"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// stdinReader is shared across the login flow to avoid losing buffered bytes
// when multiple prompts read from os.Stdin in sequence.
var stdinReader *bufio.Reader

func getStdinReader() *bufio.Reader {
	if stdinReader == nil {
		stdinReader = bufio.NewReader(os.Stdin)
	}
	return stdinReader
}

func newLoginCmd() *cobra.Command {
	var signedURLs bool

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with your API key",
		Long:  "Save your API key to the config file. The key is validated against the API before saving.",
		RunE: func(cmd *cobra.Command, args []string) error {
			p := getPrinter()

			// Prompt for API key
			fmt.Fprint(os.Stderr, "Enter your API key: ")
			apiKey, err := readSecret()
			if err != nil {
				return fmt.Errorf("reading API key: %w", err)
			}
			fmt.Fprintln(os.Stderr) // newline after masked input

			apiKey = strings.TrimSpace(apiKey)
			if apiKey == "" {
				return fmt.Errorf("API key cannot be empty")
			}

			// Validate key
			p.Debug("validating API key")
			client, err := renderscreenshot.New(apiKey)
			if err != nil {
				return fmt.Errorf("invalid API key format: %w", err)
			}

			usage, err := client.Usage(cmd.Context())
			if err != nil {
				return fmt.Errorf("authentication failed: %w", err)
			}

			// Save config
			cfg, err := config.Load()
			if err != nil {
				cfg = &config.Config{}
			}
			cfg.APIKey = apiKey

			// Optionally prompt for signing keys
			if signedURLs {
				reader := getStdinReader()

				fmt.Fprint(os.Stderr, "Enter your public key ID (rs_pub_...): ")
				pubKey, _ := reader.ReadString('\n')
				pubKey = strings.TrimSpace(pubKey)
				if pubKey != "" {
					cfg.PublicKeyID = pubKey
				}

				fmt.Fprint(os.Stderr, "Enter your secret key: ")
				secKey, err := readSecret()
				if err == nil {
					secKey = strings.TrimSpace(secKey)
					if secKey != "" {
						cfg.SecretKey = secKey
					}
				}
				fmt.Fprintln(os.Stderr)
			}

			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("saving config: %w", err)
			}

			p.Print("Authenticated successfully!\n")
			p.Print("  API Key:    %s\n", output.MaskKey(apiKey))
			p.Print("  Credits:    %d remaining\n", usage.Remaining)
			p.Print("  Config:     %s\n", config.Path())

			return nil
		},
	}

	cmd.Flags().BoolVar(&signedURLs, "signed-urls", false, "also configure keys for signed URLs")
	return cmd
}

// readSecret reads input without echoing (for passwords/keys).
func readSecret() (string, error) {
	fd := int(os.Stdin.Fd())
	if term.IsTerminal(fd) {
		b, err := term.ReadPassword(fd)
		return string(b), err
	}
	// Non-terminal (piped input) — use shared reader
	reader := getStdinReader()
	line, err := reader.ReadString('\n')
	return strings.TrimSpace(line), err
}
