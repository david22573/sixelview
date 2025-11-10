package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func cmdView() *cli.Command {
	return &cli.Command{
		Name:      "view",
		Usage:     "Render an image as SIXEL in your terminal",
		ArgsUsage: "<image-path>",
		Flags: []cli.Flag{
			&cli.IntFlag{Name: "width", Aliases: []string{"w"}, Usage: "Max width in terminal cells"},
			&cli.IntFlag{Name: "height", Aliases: []string{"ht"}, Usage: "Max height in terminal rows"},
		},
		Action: func(c *cli.Context) error {
			if c.Args().Len() == 0 {
				return fmt.Errorf("expected an image file path")
			}

			path := c.Args().First()

			if !detectSixelSupport() {
				return fmt.Errorf("terminal does not appear to support SIXEL (TERM=%s)", os.Getenv("TERM"))
			}

			img, err := loadImage(path)
			if err != nil {
				return fmt.Errorf("load image: %w", err)
			}

			cols, rows, err := getTerminalSize()
			if err != nil {
				// fallback to provided flags or default
				cols = c.Int("width")
				rows = c.Int("height")
			}

			if cols == 0 {
				cols = 80
			}
			if rows == 0 {
				rows = 24
			}

			resized := resizeImage(img, cols*2, rows*4) // approximate pixels per cell

			if err := encodeToSixel(resized); err != nil {
				return fmt.Errorf("encode sixel: %w", err)
			}

			return nil
		},
	}
}
