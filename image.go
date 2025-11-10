package main

import (
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/mattn/go-sixel"
	"golang.org/x/image/draw"
)

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
		return src // don't upscale
	}

	nw := int(float64(w) * scale)
	nh := int(float64(h) * scale)

	dst := image.NewRGBA(image.Rect(0, 0, nw, nh))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, sb, draw.Over, nil)
	return dst
}

func encodeToSixel(img image.Image) error {
	if img == nil {
		return errors.New("nil image passed to encodeToSixel")
	}

	enc := sixel.NewEncoder(os.Stdout)
	if err := enc.Encode(img); err != nil {
		return fmt.Errorf("sixel encode failed: %w", err)
	}
	return nil
}
