package main

import (
	_ "embed"
	"fmt"

	"github.com/a-h/gpu"
)

//go:embed add.metal
var source string

func main() {
	gpu.Setup(source)

	input := gpu.NewMatrix(2, 10, 1)
	// Initialize like:
	// 0.0 0.0
	// 1.0 1.0
	// 2.0 2.0
	z := input.D - 1
	for y := 0; y < input.H; y++ {
		for x := 0; x < input.W; x++ {
			input.Set(x, y, z, float32(y))
		}
	}
	// 1 across, 10 down, 1 deep.
	output := gpu.NewMatrix(1, 10, 1)
	gpu.CreateBuffers(input, output)

	data := gpu.Run(input, output)
	for _, summed := range data {
		fmt.Printf("Summed: %v\n", summed)
	}
}
