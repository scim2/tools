package gen

import (
	"bytes"
	"errors"
	"github.com/elimity-com/scim/schema"
)

// GenerateStruct creates a buffer with a go representation of the resource described in the given schema.
func GenerateStruct(s schema.Schema) (*bytes.Buffer, error) {
	b := &bytes.Buffer{}
	return b, newGenWriter(b).generateStruct(s)
}

func (w *genWriter) generateStruct(s schema.Schema) error {
	if !s.Name.Present() {
		return errors.New("schema does not have a name")
	}

	if len(s.Attributes) == 0 {
		w.lnf("type %s struct {}", s.Name.Value())
		return nil
	}

	w.lnf("type %s struct {", s.Name.Value())
	w.ln("}")
	return nil
}
