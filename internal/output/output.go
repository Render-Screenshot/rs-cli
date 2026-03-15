package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

// Printer handles output formatting based on global flags.
type Printer struct {
	JSON    bool
	Quiet   bool
	Verbose bool
	Out     io.Writer
	Err     io.Writer
}

// New creates a Printer with default stdout/stderr.
func New(jsonMode, quiet, verbose bool) *Printer {
	return &Printer{
		JSON:    jsonMode,
		Quiet:   quiet,
		Verbose: verbose,
		Out:     os.Stdout,
		Err:     os.Stderr,
	}
}

// Print writes to stdout (suppressed in quiet mode).
func (p *Printer) Print(format string, args ...interface{}) {
	if p.Quiet {
		return
	}
	fmt.Fprintf(p.Out, format, args...)
}

// Println writes a line to stdout (suppressed in quiet mode).
func (p *Printer) Println(msg string) {
	if p.Quiet {
		return
	}
	fmt.Fprintln(p.Out, msg)
}

// Error writes to stderr (never suppressed).
func (p *Printer) Error(format string, args ...interface{}) {
	fmt.Fprintf(p.Err, "Error: "+format+"\n", args...)
}

// Debug writes to stderr only in verbose mode.
func (p *Printer) Debug(format string, args ...interface{}) {
	if !p.Verbose {
		return
	}
	fmt.Fprintf(p.Err, "[debug] "+format+"\n", args...)
}

// PrintJSON marshals v as indented JSON to stdout.
func (p *Printer) PrintJSON(v interface{}) error {
	enc := json.NewEncoder(p.Out)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

// Table prints a slice of rows as an aligned table.
// headers is the first row; rows contains the data.
func (p *Printer) Table(headers []string, rows [][]string) {
	w := tabwriter.NewWriter(p.Out, 0, 0, 2, ' ', 0)

	fmt.Fprintln(w, strings.Join(headers, "\t"))
	fmt.Fprintln(w, strings.Repeat("─", len(strings.Join(headers, "  "))))

	for _, row := range rows {
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}
	w.Flush()
}

// MaskKey masks an API key, showing prefix and last 4 chars.
func MaskKey(key string) string {
	if key == "" {
		return "(not set)"
	}
	if len(key) <= 12 {
		return "****"
	}
	// Find prefix (e.g., "rs_live_")
	parts := strings.SplitN(key, "_", 3)
	if len(parts) < 3 {
		return key[:4] + "****" + key[len(key)-4:]
	}
	prefix := parts[0] + "_" + parts[1] + "_"
	suffix := key[len(key)-4:]
	return prefix + "****" + suffix
}
