package main

import (
	_ "embed"
	"fmt"

	"github.com/a-h/gpu"
)

//go:embed add.metal
var source string

func main() {
	// Compilation has to be done once.
	gpu.Compile(source)

	input := gpu.NewMatrix[float32](2, 10, 1)
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
	output := gpu.NewMatrix[float32](1, 10, 1)

	// Run code on GPU, includes copying the matrix to the GPU.
	gpu.Run(input, output)

	// The GPU code adds the numbers in column A and B together, so the results are:
	// 0.0
	// 2.0 (1+1)
	// 4.0 (2+2)
	// 6.0 (3+3)
	// ...
	for y := 0; y < output.H; y++ {
		fmt.Printf("Summed: %v\n", output.Get(0, y, 0))
	}
}
