package flags

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	renderscreenshot "github.com/Render-Screenshot/rs-go"

	"github.com/spf13/cobra"
)

// TakeFlags holds all registered take-related flag values.
type TakeFlags struct {
	// Input
	HTML string

	// Output
	Output string
	Format string
	Quality int
	Stdout  bool
	Open    bool

	// Viewport
	Width  int
	Height int
	Scale  float64
	Mobile bool
	Device string

	// Capture
	FullPage bool
	Element  string
	Preset   string

	// Wait
	Delay        int
	WaitFor      string
	WaitSelector string
	Timeout      int

	// Page manipulation
	Click        string
	Hide         string
	Remove       string
	InjectScript string
	InjectStyle  string

	// Content blocking
	BlockAds       bool
	BlockTrackers  bool
	BlockCookies   bool
	BlockChat      bool
	BlockURLs      string
	BlockResources string

	// Browser emulation
	DarkMode      bool
	ReducedMotion bool
	Media         string
	Timezone      string
	Locale        string
	UserAgent     string
	Geolocation   string

	// Network
	Headers    string
	Cookies    string
	AuthBasic  string
	AuthBearer string
	BypassCSP  bool

	// Cache
	CacheTTL     int
	CacheRefresh bool
	NoCache      bool

	// PDF
	PDFPaper           string
	PDFWidth           string
	PDFHeight          string
	PDFLandscape       bool
	PDFMargin          string
	PDFScale           float64
	PDFBackground      bool
	PDFHeader          string
	PDFFooter          string
	PDFFitOnePage      bool
	PDFPageRanges      string
	PDFPreferCSSPageSize bool

	// Storage
	StoragePath string
	StorageACL  string
}

// Register adds all take flags to the given command.
func Register(cmd *cobra.Command) *TakeFlags {
	f := &TakeFlags{}

	// Input
	cmd.Flags().StringVar(&f.HTML, "html", "", "render HTML content instead of URL")

	// Output
	cmd.Flags().StringVarP(&f.Output, "output", "o", "", "save to file path")
	cmd.Flags().StringVarP(&f.Format, "format", "f", "png", "output format: png, jpeg, webp, pdf")
	cmd.Flags().IntVarP(&f.Quality, "quality", "q", 80, "image quality 1-100 (jpeg/webp)")
	cmd.Flags().BoolVar(&f.Stdout, "stdout", false, "write binary to stdout")
	cmd.Flags().BoolVar(&f.Open, "open", false, "open result in default viewer")

	// Viewport
	cmd.Flags().IntVar(&f.Width, "width", 0, "viewport width in pixels")
	cmd.Flags().IntVar(&f.Height, "height", 0, "viewport height in pixels")
	cmd.Flags().Float64Var(&f.Scale, "scale", 0, "device scale factor 1-3")
	cmd.Flags().BoolVar(&f.Mobile, "mobile", false, "enable mobile emulation")
	cmd.Flags().StringVar(&f.Device, "device", "", "emulate device (e.g., iphone_14_pro)")

	// Capture
	cmd.Flags().BoolVar(&f.FullPage, "full-page", false, "capture full scrollable page")
	cmd.Flags().StringVar(&f.Element, "element", "", "capture specific CSS element")
	cmd.Flags().StringVar(&f.Preset, "preset", "", "use preset: og_card, twitter_card, etc.")

	// Wait
	cmd.Flags().IntVar(&f.Delay, "delay", 0, "wait additional ms before capture")
	cmd.Flags().StringVar(&f.WaitFor, "wait-for", "", "wait condition: load, networkidle, domcontentloaded")
	cmd.Flags().StringVar(&f.WaitSelector, "wait-selector", "", "wait for CSS selector to appear")
	cmd.Flags().IntVar(&f.Timeout, "timeout", 0, "maximum wait time in seconds")

	// Page manipulation
	cmd.Flags().StringVar(&f.Click, "click", "", "CSS selector to click before capture")
	cmd.Flags().StringVar(&f.Hide, "hide", "", "comma-separated selectors to hide")
	cmd.Flags().StringVar(&f.Remove, "remove", "", "comma-separated selectors to remove")
	cmd.Flags().StringVar(&f.InjectScript, "inject-script", "", "JavaScript to inject before capture")
	cmd.Flags().StringVar(&f.InjectStyle, "inject-style", "", "CSS to inject before capture")

	// Content blocking
	cmd.Flags().BoolVar(&f.BlockAds, "block-ads", false, "block ad network requests")
	cmd.Flags().BoolVar(&f.BlockTrackers, "block-trackers", false, "block analytics/tracking")
	cmd.Flags().BoolVar(&f.BlockCookies, "block-cookies", false, "dismiss cookie banners")
	cmd.Flags().BoolVar(&f.BlockChat, "block-chat", false, "block chat widgets")
	cmd.Flags().StringVar(&f.BlockURLs, "block-urls", "", "block URLs matching patterns (comma-separated)")
	cmd.Flags().StringVar(&f.BlockResources, "block-resources", "", "block resource types (comma-separated)")

	// Browser emulation
	cmd.Flags().BoolVar(&f.DarkMode, "dark-mode", false, "force dark color scheme")
	cmd.Flags().BoolVar(&f.ReducedMotion, "reduced-motion", false, "prefer reduced motion")
	cmd.Flags().StringVar(&f.Media, "media", "", "media type: screen or print")
	cmd.Flags().StringVar(&f.Timezone, "timezone", "", "set timezone (e.g., America/New_York)")
	cmd.Flags().StringVar(&f.Locale, "locale", "", "set locale (e.g., fr-FR)")
	cmd.Flags().StringVar(&f.UserAgent, "user-agent", "", "custom User-Agent string")
	cmd.Flags().StringVar(&f.Geolocation, "geolocation", "", "spoof geolocation (lat,lng)")

	// Network
	cmd.Flags().StringVar(&f.Headers, "headers", "", "custom HTTP headers as JSON")
	cmd.Flags().StringVar(&f.Cookies, "cookies", "", "cookies as JSON array")
	cmd.Flags().StringVar(&f.AuthBasic, "auth-basic", "", "HTTP Basic auth (user:pass)")
	cmd.Flags().StringVar(&f.AuthBearer, "auth-bearer", "", "bearer token for target page")
	cmd.Flags().BoolVar(&f.BypassCSP, "bypass-csp", false, "bypass Content Security Policy")

	// Cache
	cmd.Flags().IntVar(&f.CacheTTL, "cache-ttl", 0, "cache duration in seconds")
	cmd.Flags().BoolVar(&f.CacheRefresh, "cache-refresh", false, "force fresh render")
	cmd.Flags().BoolVar(&f.NoCache, "no-cache", false, "disable caching entirely")

	// PDF
	cmd.Flags().StringVar(&f.PDFPaper, "pdf-paper", "", "paper size: a4, letter, legal, tabloid")
	cmd.Flags().StringVar(&f.PDFWidth, "pdf-width", "", "custom PDF width (e.g., 210mm, 8.5in)")
	cmd.Flags().StringVar(&f.PDFHeight, "pdf-height", "", "custom PDF height (e.g., 297mm, 11in)")
	cmd.Flags().BoolVar(&f.PDFLandscape, "pdf-landscape", false, "landscape orientation")
	cmd.Flags().StringVar(&f.PDFMargin, "pdf-margin", "", "uniform margin (e.g., 1cm)")
	cmd.Flags().Float64Var(&f.PDFScale, "pdf-scale", 0, "scale factor 0.1-2.0")
	cmd.Flags().BoolVar(&f.PDFBackground, "pdf-background", false, "include background graphics")
	cmd.Flags().StringVar(&f.PDFHeader, "pdf-header", "", "header HTML template")
	cmd.Flags().StringVar(&f.PDFFooter, "pdf-footer", "", "footer HTML template")
	cmd.Flags().BoolVar(&f.PDFFitOnePage, "pdf-fit-one-page", false, "fit content to one page")
	cmd.Flags().StringVar(&f.PDFPageRanges, "pdf-page-ranges", "", "pages to include (e.g., 1-5)")
	cmd.Flags().BoolVar(&f.PDFPreferCSSPageSize, "pdf-prefer-css-page-size", false, "prefer CSS-defined page size")

	// Storage
	cmd.Flags().StringVar(&f.StoragePath, "storage-path", "", "storage path pattern")
	cmd.Flags().StringVar(&f.StorageACL, "storage-acl", "", "storage ACL: public-read or private")

	return f
}

// BuildTakeOptions converts flag values into SDK TakeOptions.
// Returns an error if flag values are malformed (e.g. invalid JSON).
func (f *TakeFlags) BuildTakeOptions(url string) (*renderscreenshot.TakeOptions, error) {
	var opts *renderscreenshot.TakeOptions

	switch {
	case f.HTML != "":
		opts = renderscreenshot.HTML(f.HTML)
	case url != "":
		opts = renderscreenshot.URL(url)
	default:
		opts = renderscreenshot.URL("")
	}

	// Preset / Device
	if f.Preset != "" {
		opts.Preset(f.Preset)
	}
	if f.Device != "" {
		opts.Device(f.Device)
	}

	// Viewport
	if f.Width > 0 {
		opts.Width(f.Width)
	}
	if f.Height > 0 {
		opts.Height(f.Height)
	}
	if f.Scale > 0 {
		opts.Scale(f.Scale)
	}
	if f.Mobile {
		opts.Mobile()
	}

	// Capture
	if f.FullPage {
		opts.FullPage()
	}
	if f.Element != "" {
		opts.Element(f.Element)
	}

	// Format
	switch strings.ToLower(f.Format) {
	case "jpeg", "jpg":
		opts.Format(renderscreenshot.FormatJPEG)
	case "webp":
		opts.Format(renderscreenshot.FormatWebP)
	case "pdf":
		opts.Format(renderscreenshot.FormatPDF)
	default:
		opts.Format(renderscreenshot.FormatPNG)
	}
	if f.Quality != 80 {
		opts.Quality(f.Quality)
	}

	// Wait
	if f.Delay > 0 {
		opts.Delay(f.Delay)
	}
	if f.WaitFor != "" {
		switch strings.ToLower(f.WaitFor) {
		case "networkidle":
			opts.WaitFor(renderscreenshot.WaitNetworkIdle)
		case "domcontentloaded":
			opts.WaitFor(renderscreenshot.WaitDOMContentLoaded)
		default:
			opts.WaitFor(renderscreenshot.WaitLoad)
		}
	}
	if f.WaitSelector != "" {
		opts.WaitForSelector(f.WaitSelector)
	}
	if f.Timeout > 0 {
		opts.WaitForTimeout(f.Timeout * 1000) // convert seconds to ms
	}

	// Page manipulation
	if f.Click != "" {
		opts.Click(f.Click)
	}
	if f.Hide != "" {
		opts.Hide(splitComma(f.Hide))
	}
	if f.Remove != "" {
		opts.Remove(splitComma(f.Remove))
	}
	if f.InjectScript != "" {
		opts.InjectScript(f.InjectScript)
	}
	if f.InjectStyle != "" {
		opts.InjectStyle(f.InjectStyle)
	}

	// Content blocking
	if f.BlockAds {
		opts.BlockAds()
	}
	if f.BlockTrackers {
		opts.BlockTrackers()
	}
	if f.BlockCookies {
		opts.BlockCookieBanners()
	}
	if f.BlockChat {
		opts.BlockChatWidgets()
	}
	if f.BlockURLs != "" {
		opts.BlockURLs(splitComma(f.BlockURLs))
	}
	if f.BlockResources != "" {
		opts.BlockResources(splitComma(f.BlockResources))
	}

	// Browser emulation
	if f.DarkMode {
		opts.DarkMode()
	}
	if f.ReducedMotion {
		opts.ReducedMotion()
	}
	if f.Media != "" {
		switch strings.ToLower(f.Media) {
		case "print":
			opts.SetMediaType(renderscreenshot.MediaPrint)
		default:
			opts.SetMediaType(renderscreenshot.MediaScreen)
		}
	}
	if f.Timezone != "" {
		opts.Timezone(f.Timezone)
	}
	if f.Locale != "" {
		opts.Locale(f.Locale)
	}
	if f.UserAgent != "" {
		opts.UserAgent(f.UserAgent)
	}
	if f.Geolocation != "" {
		lat, lon, err := parseGeolocation(f.Geolocation)
		if err != nil {
			return nil, fmt.Errorf("invalid --geolocation: %w", err)
		}
		opts.SetGeolocation(lat, lon)
	}

	// Network
	if f.Headers != "" {
		var h map[string]string
		if err := json.Unmarshal([]byte(f.Headers), &h); err != nil {
			return nil, fmt.Errorf("invalid --headers JSON: %w", err)
		}
		opts.Headers(h)
	}
	if f.Cookies != "" {
		var cookies []renderscreenshot.Cookie
		if err := json.Unmarshal([]byte(f.Cookies), &cookies); err != nil {
			return nil, fmt.Errorf("invalid --cookies JSON: %w", err)
		}
		opts.Cookies(cookies)
	}
	if f.AuthBasic != "" {
		parts := strings.SplitN(f.AuthBasic, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid --auth-basic: expected user:password format")
		}
		opts.AuthBasic(parts[0], parts[1])
	}
	if f.AuthBearer != "" {
		opts.AuthBearer(f.AuthBearer)
	}
	if f.BypassCSP {
		opts.BypassCSP()
	}

	// Cache
	if f.CacheTTL > 0 {
		opts.CacheTTL(f.CacheTTL)
	}
	if f.CacheRefresh {
		opts.CacheRefresh()
	}
	if f.NoCache {
		opts.CacheTTL(0)
	}

	// PDF options
	if f.PDFPaper != "" {
		switch strings.ToLower(f.PDFPaper) {
		case "a3":
			opts.PDFPaperSize(renderscreenshot.PaperA3)
		case "a4":
			opts.PDFPaperSize(renderscreenshot.PaperA4)
		case "a5":
			opts.PDFPaperSize(renderscreenshot.PaperA5)
		case "legal":
			opts.PDFPaperSize(renderscreenshot.PaperLegal)
		case "letter":
			opts.PDFPaperSize(renderscreenshot.PaperLetter)
		case "ledger", "tabloid":
			opts.PDFPaperSize(renderscreenshot.PaperLedger)
		}
	}
	if f.PDFWidth != "" {
		opts.PDFWidth(f.PDFWidth)
	}
	if f.PDFHeight != "" {
		opts.PDFHeight(f.PDFHeight)
	}
	if f.PDFLandscape {
		opts.PDFLandscape()
	}
	if f.PDFMargin != "" {
		opts.PDFMarginUniform(f.PDFMargin)
	}
	if f.PDFScale > 0 {
		opts.PDFScale(f.PDFScale)
	}
	if f.PDFBackground {
		opts.PDFPrintBackground()
	}
	if f.PDFHeader != "" {
		opts.PDFHeader(f.PDFHeader)
	}
	if f.PDFFooter != "" {
		opts.PDFFooter(f.PDFFooter)
	}
	if f.PDFFitOnePage {
		opts.PDFFitOnePage()
	}
	if f.PDFPageRanges != "" {
		opts.PDFPageRanges(f.PDFPageRanges)
	}
	if f.PDFPreferCSSPageSize {
		opts.PDFPreferCSSPageSize()
	}

	// Storage
	if f.StoragePath != "" {
		opts.StorageEnabled()
		opts.StoragePath(f.StoragePath)
	}
	if f.StorageACL != "" {
		switch strings.ToLower(f.StorageACL) {
		case "public-read":
			opts.StorageACL(renderscreenshot.ACLPublicRead)
		case "private":
			opts.StorageACL(renderscreenshot.ACLPrivate)
		}
	}

	return opts, nil
}

// FileExtension returns the file extension for the current format.
func (f *TakeFlags) FileExtension() string {
	switch strings.ToLower(f.Format) {
	case "jpeg", "jpg":
		return ".jpg"
	case "webp":
		return ".webp"
	case "pdf":
		return ".pdf"
	default:
		return ".png"
	}
}

func splitComma(s string) []string {
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

func parseGeolocation(s string) (float64, float64, error) {
	parts := strings.SplitN(s, ",", 2)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid geolocation format, expected lat,lng")
	}
	lat, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid latitude: %w", err)
	}
	lon, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid longitude: %w", err)
	}
	return lat, lon, nil
}
