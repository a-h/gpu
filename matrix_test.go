package gpu

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMatrix(t *testing.T) {
	t.Run("1D matrix", func(t *testing.T) {
		m := NewMatrix[float32](10, 1, 1)
		m.Set(0, 0, 0, 1.0)
		m.Set(9, 0, 0, 1.0)
		expected := []float32{1.0, 0, 0, 0, 0, 0, 0, 0, 0, 1.0}
		expectMatrix(t, expected, m.Data)
		expectValue(t, m, 0, 0, 0, 1.0)
		expectValue(t, m, 9, 0, 0, 1.0)
	})
	t.Run("2D matrix", func(t *testing.T) {
		m := NewMatrix[float32](3, 3, 1)
		m.Set(0, 0, 0, 1.0)
		m.Set(1, 1, 0, 1.0)
		m.Set(2, 2, 0, 1.0)
		expected := []float32{
			1.0, 0.0, 0.0,
			0.0, 1.0, 0.0,
			0.0, 0.0, 1.0,
		}
		expectMatrix(t, expected, m.Data)
		expectValue(t, m, 2, 2, 0, 1.0)
	})
	t.Run("#D matrix", func(t *testing.T) {
		m := NewMatrix[float32](3, 3, 2)
		m.Set(0, 0, 1, 1.0)
		m.Set(1, 1, 1, 1.0)
		m.Set(2, 2, 1, 1.0)
		expected := []float32{
			0.0, 0.0, 0.0,
			0.0, 0.0, 0.0,
			0.0, 0.0, 0.0,
			1.0, 0.0, 0.0,
			0.0, 1.0, 0.0,
			0.0, 0.0, 1.0,
		}
		expectMatrix(t, expected, m.Data)
		expectValue(t, m, 2, 2, 0, 0.0)
		expectValue(t, m, 2, 2, 1, 1.0)
	})
}

func expectMatrix[T GPUType](t *testing.T, want, got []T) {
	if diff := cmp.Diff(want, got); diff != "" {
		t.Error(diff)
	}
}

func expectValue[T GPUType](t *testing.T, m *Matrix[T], x, y, z int, expected T) {
	if got := m.Get(x, y, z); got != expected {
		t.Errorf("expected %v, got %v", expected, got)
	}
}
