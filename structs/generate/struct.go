package gen

import (
	"bytes"
	"errors"
	"github.com/elimity-com/scim/schema"
	"strings"
)

type StructGenerator struct {
	w   *genWriter
	s   schema.Schema
	ptr bool
}

func NewStructGenerator(s schema.Schema) (StructGenerator, error) {
	if !s.Name.Present() {
		return StructGenerator{}, errors.New("schema does not have a name")
	}

	return StructGenerator{
		w: newGenWriter(&bytes.Buffer{}),
		s: s,
	}, nil
}

// UsePtr indicates whether the generator will use pointers if the attribute is not required.
func (g *StructGenerator) UsePtr(t bool) *StructGenerator{
	g.ptr = t
	return g
}

// Generate creates a buffer with a go representation of the resource described in the given schema.
func (g *StructGenerator) Generate() *bytes.Buffer {
	g.generateStruct(g.s.Name.Value(), g.s.Description.Value(), g.s.Attributes)
	return g.w.writer.(*bytes.Buffer)
}

func (g *StructGenerator) generateStruct(name, desc string, attrs schema.Attributes) {
	w := g.w

	if desc != "" {
		w.ln(comment(wrap(desc)))
	}

	if len(attrs) == 0 {
		w.lnf("type %s struct {}", cap(name))
		return
	}

	w.lnf("type %s struct {", name)
	g.generateStructFields(name, attrs)
	w.ln("}")

	for _, attr := range attrs {
		if attr.AttributeType() == "complex" {
			typ := cap(attr.Name())
			if attr.MultiValued() && strings.HasSuffix(typ, "s") {
				typ = strings.TrimSuffix(typ, "s")
			}

			w.n()
			g.generateStruct(name+cap(attr.Name()), attr.Description(), attr.SubAttributes())
		}
	}
}

func (g *StructGenerator) generateStructFields(name string, attrs schema.Attributes) {
	w := g.w

	// get longest name to indent fields.
	var indent int
	for _, attr := range attrs {
		if l := len(cap(attr.Name())); l > indent {
			indent = l
		}
	}

	for _, attr := range attrs {
		var typ string
		switch t := attr.AttributeType(); t {
		case "decimal":
			typ = "float64"
		case "integer":
			typ = "int"
		case "boolean":
			typ = "bool"
		case "complex":
			typ = cap(name + cap(attr.Name()))
		default:
			typ = "string"
		}

		// field name
		name := cap(attr.Name())
		w.in(4).w(name)
		w.sp(indent - len(name) + 1)

		if attr.MultiValued() {
			w.w("[]")
			if strings.HasSuffix(typ, "s") {
				typ = strings.TrimSuffix(typ, "s")
			}
		} else if !attr.Required() && g.ptr {
			w.w("*")
		}

		w.ln(typ)
	}
}
