package plugin_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/synapse-service/node-vs/pkg/plugin"
)

func TestReaderWriter(t *testing.T) {
	body := []byte{0, 0, 7}
	var buf bytes.Buffer

	w := plugin.NewWriter(&buf)
	n, err := w.Write(body)
	assert.NoError(t, err, "write error")
	assert.Equal(t, 3, n, "write bytes length")

	r := plugin.NewReader(&buf)
	b, err := r.ReadAll()
	assert.NoError(t, err, "read all error")
	assert.Equal(t, 3, len(b), "read bytes length")
}
