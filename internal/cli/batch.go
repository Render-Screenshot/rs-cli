package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	renderscreenshot "github.com/Render-Screenshot/rs-go"

	"github.com/Render-Screenshot/rs-cli/internal/flags"
	"github.com/spf13/cobra"
)

func newBatchCmd() *cobra.Command {
	var tf *flags.TakeFlags
	var (
		file     string
		manifest string
	)

	cmd := &cobra.Command{
		Use:   "batch [urls...]",
		Short: "Take screenshots of multiple URLs",
		Long:  "Take screenshots of multiple URLs with shared options. URLs can be positional args, from a file, or from stdin.",
		RunE: func(cmd *cobra.Command, args []string) error {
			p := getPrinter()

			urls, err := collectURLs(args, file)
			if err != nil {
				return err
			}

			// Manifest mode: per-URL options
			if manifest != "" {
				return runManifestBatch(cmd, manifest)
			}

			if len(urls) == 0 {
				return fmt.Errorf("no URLs provided. Pass URLs as arguments, --file, or pipe to stdin")
			}

			client, err := newClient()
			if err != nil {
				p.Error("%s", err)
				os.Exit(2)
			}

			opts, err := tf.BuildTakeOptions("")
			if err != nil {
				return err
			}

			p.Debug("batch processing %d URLs", len(urls))

			resp, err := client.Batch(cmd.Context(), urls, opts)
			if err != nil {
				return fmt.Errorf("batch failed: %w", err)
			}

			if flagJSON {
				return p.PrintJSON(resp)
			}

			// Table output
			headers := []string{"URL", "STATUS", "IMAGE URL"}
			rows := make([][]string, 0, len(resp.Results))
			for _, r := range resp.Results {
				status := r.Status
				imgURL := r.ImageURL
				if r.Error != "" {
					status = "error"
					imgURL = r.Error
				}
				rows = append(rows, []string{r.URL, status, imgURL})
			}
			p.Table(headers, rows)

			return nil
		},
	}

	cmd.Flags().StringVar(&file, "file", "", "read URLs from file (one per line)")
	cmd.Flags().StringVar(&manifest, "manifest", "", "JSON manifest with per-URL options")

	cmd.AddCommand(newBatchStatusCmd())

	tf = flags.Register(cmd)
	return cmd
}

func newBatchStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status <batch_id>",
		Short: "Check the status of a batch job",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			p := getPrinter()

			client, err := newClient()
			if err != nil {
				p.Error("%s", err)
				os.Exit(2)
			}

			resp, err := client.GetBatch(cmd.Context(), args[0])
			if err != nil {
				return fmt.Errorf("fetching batch status: %w", err)
			}

			if flagJSON {
				return p.PrintJSON(resp)
			}

			p.Print("Batch ID:   %s\n", resp.ID)
			p.Print("Status:     %s\n", resp.Status)
			p.Print("Total:      %d\n", resp.Total)
			p.Print("Completed:  %d\n", resp.Completed)
			p.Print("Failed:     %d\n", resp.Failed)

			if len(resp.Results) > 0 {
				p.Print("\n")
				headers := []string{"URL", "STATUS", "IMAGE URL"}
				rows := make([][]string, 0, len(resp.Results))
				for _, r := range resp.Results {
					status := r.Status
					imgURL := r.ImageURL
					if r.Error != "" {
						status = "error"
						imgURL = r.Error
					}
					rows = append(rows, []string{r.URL, status, imgURL})
				}
				p.Table(headers, rows)
			}

			return nil
		},
	}
}

func collectURLs(args []string, file string) ([]string, error) {
	var urls []string

	// From positional args
	urls = append(urls, args...)

	// From file
	if file != "" {
		f, err := os.Open(file)
		if err != nil {
			return nil, fmt.Errorf("opening file: %w", err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" && !strings.HasPrefix(line, "#") {
				urls = append(urls, line)
			}
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("reading file: %w", err)
		}
	}

	// From stdin if "-" was passed or no other input
	if len(urls) == 1 && urls[0] == "-" {
		urls = urls[:0] // clear the "-"
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" && !strings.HasPrefix(line, "#") {
				urls = append(urls, line)
			}
		}
	}

	return urls, nil
}

func runManifestBatch(cmd *cobra.Command, manifestPath string) error {
	p := getPrinter()

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("reading manifest: %w", err)
	}

	var entries []map[string]interface{}
	if err := json.Unmarshal(data, &entries); err != nil {
		return fmt.Errorf("parsing manifest: %w", err)
	}

	client, err := newClient()
	if err != nil {
		p.Error("%s", err)
		os.Exit(2)
	}

	requests := make([]renderscreenshot.BatchRequest, len(entries))
	for i, e := range entries {
		urlStr, _ := e["url"].(string)
		if urlStr == "" {
			return fmt.Errorf("manifest entry %d: missing \"url\" field", i)
		}
		opts := applyManifestOptions(urlStr, e)
		requests[i] = renderscreenshot.BatchRequest{
			URL:     urlStr,
			Options: opts,
		}
	}

	resp, err := client.BatchAdvanced(cmd.Context(), requests)
	if err != nil {
		return fmt.Errorf("batch failed: %w", err)
	}

	if flagJSON {
		return p.PrintJSON(resp)
	}

	headers := []string{"URL", "STATUS", "IMAGE URL"}
	rows := make([][]string, 0, len(resp.Results))
	for _, r := range resp.Results {
		status := r.Status
		imgURL := r.ImageURL
		if r.Error != "" {
			status = "error"
			imgURL = r.Error
		}
		rows = append(rows, []string{r.URL, status, imgURL})
	}
	p.Table(headers, rows)

	return nil
}

// applyManifestOptions builds TakeOptions from a manifest entry map,
// calling public setter methods for each recognized key.
func applyManifestOptions(url string, m map[string]interface{}) *renderscreenshot.TakeOptions {
	opts := renderscreenshot.URL(url)

	if v, ok := m["preset"].(string); ok {
		opts.Preset(v)
	}
	if v, ok := m["device"].(string); ok {
		opts.Device(v)
	}

	// Viewport
	if v, ok := toInt(m["width"]); ok {
		opts.Width(v)
	}
	if v, ok := toInt(m["height"]); ok {
		opts.Height(v)
	}
	if v, ok := toFloat(m["scale"]); ok {
		opts.Scale(v)
	}
	if v, ok := m["mobile"].(bool); ok && v {
		opts.Mobile()
	}

	// Capture
	if v, ok := m["full_page"].(bool); ok && v {
		opts.FullPage()
	}
	if v, ok := m["element"].(string); ok {
		opts.Element(v)
	}
	if v, ok := m["format"].(string); ok {
		opts.Format(renderscreenshot.ImageFormat(v))
	}
	if v, ok := toInt(m["quality"]); ok {
		opts.Quality(v)
	}

	// Wait
	if v, ok := toInt(m["delay"]); ok {
		opts.Delay(v)
	}
	if v, ok := m["wait_for"].(string); ok {
		opts.WaitFor(renderscreenshot.WaitCondition(v))
	}
	if v, ok := m["wait_selector"].(string); ok {
		opts.WaitForSelector(v)
	}

	// Page manipulation
	if v, ok := m["click"].(string); ok {
		opts.Click(v)
	}
	if v, ok := m["inject_script"].(string); ok {
		opts.InjectScript(v)
	}
	if v, ok := m["inject_style"].(string); ok {
		opts.InjectStyle(v)
	}

	// Content blocking
	if v, ok := m["block_ads"].(bool); ok && v {
		opts.BlockAds()
	}
	if v, ok := m["block_trackers"].(bool); ok && v {
		opts.BlockTrackers()
	}
	if v, ok := m["block_cookie_banners"].(bool); ok && v {
		opts.BlockCookieBanners()
	}
	if v, ok := m["block_chat_widgets"].(bool); ok && v {
		opts.BlockChatWidgets()
	}

	// Browser emulation
	if v, ok := m["dark_mode"].(bool); ok && v {
		opts.DarkMode()
	}
	if v, ok := m["reduced_motion"].(bool); ok && v {
		opts.ReducedMotion()
	}
	if v, ok := m["user_agent"].(string); ok {
		opts.UserAgent(v)
	}
	if v, ok := m["timezone"].(string); ok {
		opts.Timezone(v)
	}
	if v, ok := m["locale"].(string); ok {
		opts.Locale(v)
	}

	// Cache
	if v, ok := toInt(m["cache_ttl"]); ok {
		opts.CacheTTL(v)
	}
	if v, ok := m["cache_refresh"].(bool); ok && v {
		opts.CacheRefresh()
	}

	return opts
}

// toInt extracts an int from a JSON-decoded value (which may be float64).
func toInt(v interface{}) (int, bool) {
	switch n := v.(type) {
	case float64:
		return int(n), true
	case int:
		return n, true
	}
	return 0, false
}

// toFloat extracts a float64 from a JSON-decoded value.
func toFloat(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	}
	return 0, false
}
