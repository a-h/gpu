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
	"sync"
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
		init: &sync.Once{},
	}
	return m
}

type Matrix[T GPUType] struct {
	W, H, D int
	Data    []T
	init    *sync.Once
}

func (m *Matrix[T]) Populate() {
	m.init.Do(func() {
		if m.Data == nil {
			m.Data = make([]T, m.W*m.H*m.D)
		}
	})
}

func (m Matrix[T]) Index(x, y, z int) (i int) {
	i += z * m.W * m.H
	i += y * m.W
	i += x
	return i
}

func (m *Matrix[T]) Set(x, y, z int, v T) {
	m.Populate()
	m.Data[m.Index(x, y, z)] = v
}

func (m Matrix[T]) Get(x, y, z int) T {
	m.Populate()
	return m.Data[m.Index(x, y, z)]
}

func (m Matrix[T]) Size() int {
	return m.W * m.H * m.D
}

// Compile the shader. Only needs to be done once.
func Compile(shaderCode string) {
	src := C.CString(shaderCode)
	defer C.free(unsafe.Pointer(src))
	C.compile(src)
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
	// Setup.
	var in unsafe.Pointer
	var inputSize int
	if len(input.Data) > 0 {
		in = unsafe.Pointer(&input.Data[0])
		inputSize = input.Size()
	}
	C.createBuffers(in, C.int(dataSizeBytes[TIn]()), C.int(inputSize),
		C.int(dataSizeBytes[TOut]()), C.int(output.Size()))
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
	output.Data = unsafe.Slice((*TOut)(ptr), output.Size())
	return
}
