package cli

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/Render-Screenshot/rs-cli/internal/config"
	"github.com/Render-Screenshot/rs-cli/internal/flags"
	"github.com/spf13/cobra"
)

func newSignedURLCmd() *cobra.Command {
	var tf *flags.TakeFlags
	var (
		expires string
		copy    bool
	)

	cmd := &cobra.Command{
		Use:   "signed-url <url>",
		Short: "Generate a signed URL",
		Long:  "Generate a signed URL for public embedding without exposing your API key.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			p := getPrinter()

			client, err := newClient()
			if err != nil {
				p.Error("%s", err)
				os.Exit(2)
			}

			pubKeyID, secKey := config.ResolveSigningKeys()
			if pubKeyID == "" || secKey == "" {
				return fmt.Errorf("signed URLs require public_key_id and secret_key.\nRun 'rs login --signed-urls' or set RS_PUBLIC_KEY_ID and RS_SECRET_KEY")
			}

			duration, err := parseDuration(expires)
			if err != nil {
				return fmt.Errorf("invalid --expires: %w", err)
			}

			expiresAt := time.Now().Add(duration)
			opts, err := tf.BuildTakeOptions(args[0])
			if err != nil {
				return err
			}

			signedURL, err := client.GenerateURL(opts, expiresAt, secKey, pubKeyID)
			if err != nil {
				return fmt.Errorf("generating signed URL: %w", err)
			}

			if flagJSON {
				return p.PrintJSON(map[string]interface{}{
					"url":        signedURL,
					"expires_at": expiresAt.Format(time.RFC3339),
				})
			}

			p.Println(signedURL)

			if copy {
				if err := copyToClipboard(signedURL); err != nil {
					p.Error("could not copy to clipboard: %s", err)
				} else {
					p.Print("Copied to clipboard\n")
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&expires, "expires", "24h", "expiration: 1h, 24h, 7d, 30d")
	cmd.Flags().BoolVar(&copy, "copy", false, "copy URL to clipboard")

	tf = flags.Register(cmd)
	return cmd
}

// parseDuration handles Go durations plus "d" suffix for days.
func parseDuration(s string) (time.Duration, error) {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, "d") {
		numStr := strings.TrimSuffix(s, "d")
		var days int
		if _, err := fmt.Sscanf(numStr, "%d", &days); err != nil {
			return 0, fmt.Errorf("invalid day format: %s", s)
		}
		if days > 30 {
			return 0, fmt.Errorf("maximum expiration is 30 days")
		}
		return time.Duration(days) * 24 * time.Hour, nil
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, err
	}
	if d > 30*24*time.Hour {
		return 0, fmt.Errorf("maximum expiration is 30 days")
	}
	return d, nil
}

func copyToClipboard(text string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		cmd = exec.Command("xclip", "-selection", "clipboard")
	case "windows":
		cmd = exec.Command("clip")
	default:
		return fmt.Errorf("unsupported platform")
	}
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}
