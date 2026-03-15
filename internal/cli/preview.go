package cli

import (
	"fmt"

	"github.com/Render-Screenshot/rs-cli/internal/preview"
	"github.com/spf13/cobra"
)

func newPreviewCmd() *cobra.Command {
	var timeout int

	cmd := &cobra.Command{
		Use:   "preview <url>",
		Short: "Fetch page metadata without taking a screenshot",
		Long:  "Fetch title, description, Open Graph, and Twitter Card metadata from a URL. No API credits consumed.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			p := getPrinter()

			p.Debug("fetching metadata from %s", args[0])

			meta, err := preview.Fetch(args[0], timeout)
			if err != nil {
				return fmt.Errorf("fetching metadata: %w", err)
			}

			if flagJSON {
				return p.PrintJSON(meta)
			}

			// Human-readable output
			p.Print("Title:       %s\n", meta.Title)
			p.Print("Description: %s\n", meta.Description)
			p.Print("Favicon:     %s\n", meta.Favicon)
			p.Print("\n")

			if meta.OG.Title != "" || meta.OG.Description != "" || meta.OG.Image != "" {
				p.Println("Open Graph:")
				if meta.OG.Title != "" {
					p.Print("  og:title:       %s\n", meta.OG.Title)
				}
				if meta.OG.Description != "" {
					p.Print("  og:description: %s\n", meta.OG.Description)
				}
				if meta.OG.Image != "" {
					p.Print("  og:image:       %s\n", meta.OG.Image)
				}
				if meta.OG.URL != "" {
					p.Print("  og:url:         %s\n", meta.OG.URL)
				}
				if meta.OG.Type != "" {
					p.Print("  og:type:        %s\n", meta.OG.Type)
				}
				if meta.OG.SiteName != "" {
					p.Print("  og:site_name:   %s\n", meta.OG.SiteName)
				}
				p.Print("\n")
			}

			if meta.Twitter.Card != "" || meta.Twitter.Title != "" {
				p.Println("Twitter Card:")
				if meta.Twitter.Card != "" {
					p.Print("  twitter:card:        %s\n", meta.Twitter.Card)
				}
				if meta.Twitter.Site != "" {
					p.Print("  twitter:site:        %s\n", meta.Twitter.Site)
				}
				if meta.Twitter.Title != "" {
					p.Print("  twitter:title:       %s\n", meta.Twitter.Title)
				}
				if meta.Twitter.Description != "" {
					p.Print("  twitter:description: %s\n", meta.Twitter.Description)
				}
				if meta.Twitter.Image != "" {
					p.Print("  twitter:image:       %s\n", meta.Twitter.Image)
				}
			}

			return nil
		},
	}

	cmd.Flags().IntVar(&timeout, "timeout", 10, "maximum wait time in seconds")

	return cmd
}
