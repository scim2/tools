// Package gen provides generated go structures based on a given SCIM schema.
package gen

import (
	"fmt"
	"io"
	"strings"
)

type genWriter struct {
	writer io.Writer
	prefix string
}

func newGenWriter(w io.Writer) *genWriter {
	return &genWriter{
		writer: w,
	}
}

func (w *genWriter) Write(p []byte) (int, error) {
	return w.writer.Write(append([]byte(w.prefix), p...))
}

func (w *genWriter) w(p string) error {
	_, err := w.Write([]byte(p))
	return err
}

func (w *genWriter) n() error {
	_, err := w.Write([]byte("\n"))
	return err
}

func (w *genWriter) ln(p string) error {
	_, err := w.Write([]byte(p + "\n"))
	return err
}

func (w *genWriter) f(format string, args ...interface{}) error {
	return w.w(fmt.Sprintf(format, args...))
}

func (w *genWriter) lnf(format string, args ...interface{}) error {
	return w.ln(fmt.Sprintf(format, args...))
}

// sp add n spaces
func (w *genWriter) sp(n int) error {
	return w.w(strings.Repeat(" ", n))
}

// in adds n spaces as before the prefix.
func (w *genWriter) in(n int) *genWriter {
	return &genWriter{
		writer: w,
		prefix: strings.Repeat(" ", n) + w.prefix,
	}
}
