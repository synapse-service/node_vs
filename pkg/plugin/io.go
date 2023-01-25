package plugin

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

const HeaderLength = 8

func NewReader(r io.Reader) *Reader { return &Reader{r: r} }

type Reader struct{ r io.Reader }

func (r *Reader) ReadAll() ([]byte, error) {
	reader := bufio.NewReader(r.r)
	h := make([]byte, HeaderLength)
	_, err := reader.Read(h)
	if err != nil {
		return nil, errors.Wrap(err, "read header")
	}
	l := binary.BigEndian.Uint64(h)

	if l > 4096 {
		return nil, fmt.Errorf("too large message: %d", l)
	}

	b := make([]byte, 0, l)
	var total int
	for total != cap(b) {
		tmp := make([]byte, cap(b))
		n, err := reader.Read(tmp)
		if err != nil {
			return nil, errors.Wrap(err, "read body")
		}
		tmp = tmp[:n]
		b = append(b, tmp...)
		total += n
	}

	if uint64(total) != l {
		return nil, fmt.Errorf("unexpected length: received %d bytes, but %d expected", total, l)
	}

	return b, nil
}

func NewWriter(w io.Writer) *Writer { return &Writer{w: w} }

type Writer struct{ w io.Writer }

func (w *Writer) Write(b []byte) (int, error) {
	h := make([]byte, HeaderLength)
	binary.BigEndian.PutUint64(h, uint64(len(b)))
	n, err := w.w.Write(append(h, b...))
	if err != nil {
		return 0, errors.Wrap(err, "send error")
	}
	return n - HeaderLength, nil
}
