package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func newCacheCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cache",
		Short: "Manage cached screenshots",
		Long:  "Retrieve, delete, or purge cached screenshots.",
	}

	cmd.AddCommand(
		newCacheGetCmd(),
		newCacheDeleteCmd(),
		newCachePurgeCmd(),
	)

	return cmd
}

func newCacheGetCmd() *cobra.Command {
	var output string

	cmd := &cobra.Command{
		Use:   "get <cache_key>",
		Short: "Retrieve a cached screenshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			p := getPrinter()

			client, err := newClient()
			if err != nil {
				p.Error("%s", err)
				os.Exit(2)
			}

			data, err := client.Cache().Get(cmd.Context(), args[0])
			if err != nil {
				return fmt.Errorf("cache get failed: %w", err)
			}

			if data == nil {
				p.Println("Cache entry not found")
				os.Exit(1)
			}

			if output != "" {
				if err := os.WriteFile(output, data, 0600); err != nil {
					return fmt.Errorf("writing file: %w", err)
				}
				p.Print("Saved to %s (%d bytes)\n", output, len(data))
				return nil
			}

			_, err = os.Stdout.Write(data)
			return err
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "", "save to file")
	return cmd
}

func newCacheDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <cache_key>",
		Short: "Delete a single cache entry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			p := getPrinter()

			client, err := newClient()
			if err != nil {
				p.Error("%s", err)
				os.Exit(2)
			}

			deleted, err := client.Cache().Delete(cmd.Context(), args[0])
			if err != nil {
				return fmt.Errorf("cache delete failed: %w", err)
			}

			if flagJSON {
				return p.PrintJSON(map[string]interface{}{
					"key":     args[0],
					"deleted": deleted,
				})
			}

			if deleted {
				p.Println("Cache entry deleted")
			} else {
				p.Println("Cache entry not found")
			}
			return nil
		},
	}
}

func newCachePurgeCmd() *cobra.Command {
	var (
		keys    string
		urlPat  string
		before  string
		pattern string
	)

	cmd := &cobra.Command{
		Use:   "purge",
		Short: "Purge multiple cache entries",
		RunE: func(cmd *cobra.Command, args []string) error {
			p := getPrinter()

			client, err := newClient()
			if err != nil {
				p.Error("%s", err)
				os.Exit(2)
			}

			cache := client.Cache()
			ctx := cmd.Context()

			if keys != "" {
				keyList := strings.Split(keys, ",")
				for i := range keyList {
					keyList[i] = strings.TrimSpace(keyList[i])
				}
				result, err := cache.Purge(ctx, keyList)
				if err != nil {
					return fmt.Errorf("purge by keys failed: %w", err)
				}
				if flagJSON {
					return p.PrintJSON(result)
				}
				p.Print("Purged %d entries\n", result.Purged)
				return nil
			}

			if urlPat != "" {
				result, err := cache.PurgeURL(ctx, urlPat)
				if err != nil {
					return fmt.Errorf("purge by URL failed: %w", err)
				}
				if flagJSON {
					return p.PrintJSON(result)
				}
				p.Print("Purged %d entries\n", result.Purged)
				return nil
			}

			if before != "" {
				t, err := time.Parse(time.RFC3339, before)
				if err != nil {
					t, err = time.Parse("2006-01-02", before)
					if err != nil {
						return fmt.Errorf("invalid --before date, use YYYY-MM-DD or RFC3339 format")
					}
				}
				result, err := cache.PurgeBefore(ctx, t)
				if err != nil {
					return fmt.Errorf("purge by date failed: %w", err)
				}
				if flagJSON {
					return p.PrintJSON(result)
				}
				p.Print("Purged %d entries\n", result.Purged)
				return nil
			}

			if pattern != "" {
				result, err := cache.PurgePattern(ctx, pattern)
				if err != nil {
					return fmt.Errorf("purge by pattern failed: %w", err)
				}
				if flagJSON {
					return p.PrintJSON(result)
				}
				p.Print("Purged %d entries\n", result.Purged)
				return nil
			}

			return fmt.Errorf("specify at least one: --keys, --url, --before, or --pattern")
		},
	}

	cmd.Flags().StringVar(&keys, "keys", "", "purge by keys (comma-separated)")
	cmd.Flags().StringVar(&urlPat, "url", "", "purge by URL pattern")
	cmd.Flags().StringVar(&before, "before", "", "purge entries older than date")
	cmd.Flags().StringVar(&pattern, "pattern", "", "purge by storage path pattern")

	return cmd
}
