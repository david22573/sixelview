// main.go
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	width  int
	height int
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "sixelview <image-path>",
		Short: "Display images in terminals that support SIXEL graphics",
		Long:  "Render an image as SIXEL in your terminal",
		Args:  cobra.ExactArgs(1),
		RunE:  runView,
	}

	// Add flags
	rootCmd.Flags().IntVarP(&width, "width", "w", 0, "Max width in terminal cells (0 = auto-detect)")
	rootCmd.Flags().IntVarP(&height, "height", "t", 0, "Max height in terminal rows (0 = auto-detect)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func runView(cmd *cobra.Command, args []string) error {
	path := args[0]

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("image file does not exist: %s", path)
	}

	// Check SIXEL support
	if !detectSixelSupport() {
		return fmt.Errorf("terminal does not appear to support SIXEL (TERM=%s)", os.Getenv("TERM"))
	}

	// Load the image
	img, err := loadImage(path)
	if err != nil {
		return fmt.Errorf("failed to load image: %w", err)
	}

	// Determine terminal dimensions
	cols, rows := determineTerminalSize(width, height)

	// Calculate pixel dimensions
	pixelWidth := cols * 5
	pixelHeight := rows * 10

	// Resize image to fit terminal
	resized := resizeImage(img, pixelWidth, pixelHeight)

	// Encode and output as SIXEL
	if err := encodeToSixel(resized); err != nil {
		return fmt.Errorf("failed to encode SIXEL: %w", err)
	}

	return nil
}

// determineTerminalSize calculates the terminal size based on flags or auto-detection
func determineTerminalSize(width, height int) (cols, rows int) {
	cols = width
	rows = height

	// If either dimension is not specified, try to detect terminal size
	if cols == 0 || rows == 0 {
		detectedCols, detectedRows, err := getTerminalSize()
		if err == nil {
			// Use detected values for unspecified dimensions
			if cols == 0 {
				cols = detectedCols
			}
			if rows == 0 {
				rows = detectedRows
			}
		}
	}

	// Apply fallback defaults if still zero
	if cols == 0 {
		cols = 80
	}
	if rows == 0 {
		rows = 24
	}

	return cols, rows
}
