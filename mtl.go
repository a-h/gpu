//go:build darwin
// +build darwin

package gpu

/*
#cgo LDFLAGS: -framework Metal -framework CoreGraphics -framework Foundation
#include <stdlib.h>
#include <stdbool.h>
#include "mtl.h"
*/
import "C"
import (
	"unsafe"
)

// Load the source code.
func Setup(source string) {
	src := C.CString(source)
	defer C.free(unsafe.Pointer(src))
	C.setup(src)
}

func NewMatrix(w, h, d int) *Matrix {
	// Store matrix like a display buffer, i.e.
	// get a whole y row, and index it with x.
	// so, d, y, x
	m := &Matrix{
		W:    w,
		H:    h,
		D:    d,
		Data: make([]float32, w*h*d),
	}
	return m
}

type Matrix struct {
	W, H, D int
	Data    []float32
}

func (m Matrix) Index(x, y, z int) (i int) {
	i += z * m.W * m.H
	i += y * m.W
	i += x
	return i
}

func (m *Matrix) Set(x, y, z int, v float32) {
	m.Data[m.Index(x, y, z)] = v
}

func (m Matrix) Get(x, y, z int) float32 {
	return m.Data[m.Index(x, y, z)]
}

func (m Matrix) Size() int {
	return len(m.Data)
}

// Create the buffers to store output data.
// The output data is a 2D matrix.
func CreateBuffers(input *Matrix, output *Matrix) {
	in := (*C.float)(unsafe.Pointer(&input.Data[0]))
	out := (*C.float)(unsafe.Pointer(&output.Data[0]))
	C.createBuffers(in, C.int(input.Size()), out, C.int(output.Size()))
}

// Params matches the definitions in mtl.m etc.
type params struct {
	// Size of input matrix.
	WIn, HIn, DIn int32
	// Size of output matrix.
	WOut, HOut, DOut int32
}

func Run(input, output *Matrix) []float32 {
	// Convert the Go param struct to its C version.
	p := params{
		WIn:  int32(input.W),
		HIn:  int32(input.H),
		DIn:  int32(input.D),
		WOut: int32(output.W),
		HOut: int32(output.H),
		DOut: int32(output.D),
	}
	cp := (*C.Params)(unsafe.Pointer(&p))
	// Run.
	ptr := C.run(cp)
	return unsafe.Slice((*float32)(ptr), len(output.Data))
}
