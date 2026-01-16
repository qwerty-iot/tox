package tox

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToFloat64Array(t *testing.T) {
	// nil
	assert.Nil(t, ToFloat64Array(nil))

	// float64
	assert.Equal(t, []float64{1.23}, ToFloat64Array(1.23))

	// []float64
	assert.Equal(t, []float64{1.23, 4.56}, ToFloat64Array([]float64{1.23, 4.56}))

	// []byte
	assert.Equal(t, []float64{1, 2, 3}, ToFloat64Array([]byte{1, 2, 3}))

	// []any
	res := ToFloat64Array([]any{1.23, 4, "not a float"})
	assert.Len(t, res, 3)
	assert.Equal(t, 1.23, res[0])
	assert.Equal(t, float64(4), res[1])
	assert.True(t, math.IsNaN(res[2]))

	// array/slice via reflection
	assert.Equal(t, []float64{1, 2}, ToFloat64Array([]int{1, 2}))
	assert.Equal(t, []float64{1, 2}, ToFloat64Array([2]int{1, 2}))

	// single value
	assert.Equal(t, []float64{123}, ToFloat64Array(123))
	assert.Equal(t, []float64{456.78}, ToFloat64Array("456.78"))

	// NaN case
	res = ToFloat64Array("abc")
	assert.Len(t, res, 1)
	assert.True(t, math.IsNaN(res[0]))
}
