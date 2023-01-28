package service_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/synapse-service/node-vs/service"
)

func TestSeparate(t *testing.T) {
	a, b := []string{"a", "b", "c"}, []string{"b", "c", "d"}
	a1, intersection, b1 := service.Separate(a, b)
	assert.Equal(t, "a", strings.Join(a1, ","), "a1 outer")
	assert.Equal(t, "b,c", strings.Join(intersection, ","), "intersection")
	assert.Equal(t, "d", strings.Join(b1, ","), "b1 outer")
}

func BenchmarkSeparate(bench *testing.B) {
	a, b := []string{"a", "b", "c", "d"}, []string{"b", "c", "d", "e"}
	for i := 0; i < bench.N; i++ {
		service.Separate(a, b)
	}
}
