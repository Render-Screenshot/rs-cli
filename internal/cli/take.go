package cli

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/Render-Screenshot/rs-cli/internal/flags"
	"github.com/spf13/cobra"
)

func newTakeCmd() *cobra.Command {
	var tf *flags.TakeFlags

	cmd := &cobra.Command{
		Use:   "take <url>",
		Short: "Capture a screenshot",
		Long:  "Capture a screenshot of a URL or HTML content.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			p := getPrinter()

			// Determine input source
			inputURL := ""
			if len(args) > 0 {
				inputURL = args[0]
			}

			if inputURL == "" && tf.HTML == "" {
				return fmt.Errorf("provide a URL or --html")
			}

			client, err := newClient()
			if err != nil {
				p.Error("%s", err)
				os.Exit(2)
			}

			opts, err := tf.BuildTakeOptions(inputURL)
			if err != nil {
				return err
			}
			ctx := cmd.Context()

			// JSON output mode
			if flagJSON {
				p.Debug("taking screenshot (JSON mode)")
				resp, err := client.TakeJSON(ctx, opts)
				if err != nil {
					return fmt.Errorf("screenshot failed: %w", err)
				}
				return p.PrintJSON(resp)
			}

			// Binary output mode
			p.Debug("taking screenshot")
			data, err := client.Take(ctx, opts)
			if err != nil {
				return fmt.Errorf("screenshot failed: %w", err)
			}

			// Route output
			switch {
			case tf.Stdout:
				_, err = os.Stdout.Write(data)
				return err

			default:
				filename := tf.Output
				if filename == "" {
					filename = autoFilename(inputURL, tf.FileExtension())
				}

				if err := os.WriteFile(filename, data, 0600); err != nil {
					return fmt.Errorf("writing file: %w", err)
				}
				p.Print("Saved to %s (%d bytes)\n", filename, len(data))

				if tf.Open {
					return openFile(filename)
				}
			}
			return nil
		},
	}

	tf = flags.Register(cmd)
	return cmd
}

// autoFilename generates a filename from a URL: example.com → example-com.png
func autoFilename(rawURL, ext string) string {
	if rawURL == "" {
		return "screenshot" + ext
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return "screenshot" + ext
	}
	host := u.Hostname()
	if host == "" {
		return "screenshot" + ext
	}
	name := strings.ReplaceAll(host, ".", "-")
	return name + ext
}

// openFile opens a file with the system default application.
func openFile(path string) error {
	var cmd string
	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
	case "linux":
		cmd = "xdg-open"
	case "windows":
		cmd = "start"
	default:
		return fmt.Errorf("unsupported platform for --open")
	}
	return exec.Command(cmd, path).Start()
}
