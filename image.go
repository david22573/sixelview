package main

import (
	"image"
	"os"

	"github.com/mattn/go-sixel"
	"golang.org/x/image/draw"
)

// loadImage opens and decodes an image file into image.Image.
func loadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// resizeImage rescales src to fit within maxW x maxH pixels while preserving aspect ratio.
func resizeImage(src image.Image, maxW, maxH int) image.Image {
	if maxW <= 0 || maxH <= 0 {
		return src
	}

	sb := src.Bounds()
	w := sb.Dx()
	h := sb.Dy()

	scaleW := float64(maxW) / float64(w)
	scaleH := float64(maxH) / float64(h)
	scale := scaleW
	if scaleH < scaleW {
		scale = scaleH
	}
	if scale >= 1.0 {
		return src // no upscaling
	}

	nw := int(float64(w) * scale)
	nh := int(float64(h) * scale)

	dst := image.NewRGBA(image.Rect(0, 0, nw, nh))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, sb, draw.Over, nil)
	return dst
}

// encodeToSixel writes SIXEL bytes for the provided image to stdout.
func encodeToSixel(img image.Image) error {
	if img == nil {
		return errors.New("nil image")
	}

	enc := sixel.NewEncoder()
	enc.SetWriter(os.Stdout)
	// Optionally set palette or dithering here via enc.SetDither(true) etc.
	if err := enc.Encode(img); err != nil {
		return fmt.Errorf("sixel encode failed: %w", err)
	}
	return nil
}
