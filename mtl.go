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

type GPUType interface {
	uint8 | uint32 | int32 | float32
}

func NewMatrix[T GPUType](w, h, d int) *Matrix[T] {
	// Store matrix like a display buffer, i.e.
	// get a whole y row, and index it with x.
	// so, d, y, x
	m := &Matrix[T]{
		W:    w,
		H:    h,
		D:    d,
		Data: make([]T, w*h*d),
	}
	return m
}

func NewMatrixFromData[T GPUType](w, h, d int, data []T) *Matrix[T] {
	// Store matrix like a display buffer, i.e.
	// get a whole y row, and index it with x.
	// so, d, y, x
	m := &Matrix[T]{
		W:    w,
		H:    h,
		D:    d,
		Data: data,
	}
	return m
}

type Matrix[T GPUType] struct {
	W, H, D int
	Data    []T
}

func (m Matrix[T]) Index(x, y, z int) (i int) {
	i += z * m.W * m.H
	i += y * m.W
	i += x
	return i
}

func (m *Matrix[T]) Set(x, y, z int, v T) {
	m.Data[m.Index(x, y, z)] = v
}

func (m Matrix[T]) Get(x, y, z int) T {
	return m.Data[m.Index(x, y, z)]
}

func (m Matrix[T]) Size() int {
	return len(m.Data)
}

// Setup the GPU, passing the metal source and the
// input / output matrices.
func Setup[TIn, TOut GPUType](source string, input []TIn, output []TOut) {
	src := C.CString(source)
	defer C.free(unsafe.Pointer(src))
	C.setup(src)
	// Create the buffers to store output data.
	var in unsafe.Pointer
	if len(input) > 0 {
		in = unsafe.Pointer(&input[0])
	}
	var out unsafe.Pointer
	if len(output) > 0 {
		out = unsafe.Pointer(&output[0])
	}
	C.createBuffers(in, C.int(dataSizeBytes[TIn]()), C.int(len(input)),
		out, C.int(dataSizeBytes[TOut]()), C.int(len(output)))
}

func dataSizeBytes[T GPUType]() int32 {
	var v T
	switch any(v).(type) {
	case uint8:
		return 1
	case int32:
		return 4
	case uint32:
		return 4
	case float32:
		return 4
	}
	panic("unknown data size for GPU type")
}

// Params matches the definitions in mtl.m etc.
type params struct {
	// Size of input matrix.
	WIn, HIn, DIn int32
	// Size of output matrix.
	WOut, HOut, DOut int32
}

func Run[TIn GPUType, TOut GPUType](input *Matrix[TIn], output *Matrix[TOut]) {
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
	output.Data = unsafe.Slice((*TOut)(ptr), len(output.Data))
	return
}
