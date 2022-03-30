package main

import (
	_ "embed"
	"image"
	"image/jpeg"
	"log"
	"os"

	"github.com/a-h/gpu"
)

//go:embed grayscale.metal
var source string

func main() {
	gpu.Compile(source)

	f, err := os.Open("puppy-g7b38fec9b_1920.jpg")
	if err != nil {
		log.Fatalf("failed to read puppy JPEG: %v", err)
	}
	defer f.Close()
	jpg, err := jpeg.Decode(f)
	if err != nil {
		log.Fatalf("failed to decode JPEG: %v", err)
	}

	// Create a matrix to copy the data into.
	// Unfortunately, there's no backing byte array that's easy to access.
	// So load into the matrix.

	bounds := jpg.Bounds()
	stride := 4
	input := gpu.NewMatrix[uint8](bounds.Dx()*stride, bounds.Dy(), 1)
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			r, g, b, a := jpg.At(x, y).RGBA()
			input.Set((x*stride)+0, y, 0, uint8(r/257))
			input.Set((x*stride)+1, y, 0, uint8(g/257))
			input.Set((x*stride)+2, y, 0, uint8(b/257))
			input.Set((x*stride)+3, y, 0, uint8(a/257))
		}
	}

	// Configure the output.
	output := gpu.NewMatrix[uint8](bounds.Dx()*stride, bounds.Dy(), 1)

	// Run the processing.
	gpu.Run(input, output)

	// Write the output.
	fo, err := os.Create("gray-puppy.jpg")
	if err != nil {
		log.Fatalf("failed to create grayscale puppy: %v", err)
	}
	img := image.NewRGBA(jpg.Bounds())
	img.Pix = output.Data
	err = jpeg.Encode(fo, img, &jpeg.Options{
		Quality: 100,
	})
	if err != nil {
		log.Fatalf("failed to write JPEG to disk: %v", err)
	}
}
