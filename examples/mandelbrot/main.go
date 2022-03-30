package main

import (
	_ "embed"
	"image"
	"image/png"
	"log"
	"os"

	"github.com/a-h/gpu"
)

//go:embed mandelbrot.metal
var source string

func main() {
	gpu.Compile(source)

	size := image.Rect(0, 0, 1920, 1080)
	// Input and output are both unpopulated.
	stride := 4 // r, g, b, and A.
	input := gpu.NewMatrix[uint8](size.Dx()*stride, size.Dy(), 1)
	output := gpu.NewMatrix[uint8](size.Dx()*stride, size.Dy(), 1)

	// Run the processing.
	gpu.Run(input, output)

	// Write the output.
	fo, err := os.Create("mandelbrot.png")
	if err != nil {
		log.Fatalf("failed to create mandelbrot PNG: %v", err)
	}
	img := image.NewRGBA(size)
	img.Pix = output.Data
	err = png.Encode(fo, img)
	if err != nil {
		log.Fatalf("failed to write PNG to disk: %v", err)
	}
}
