package gen

import (
	"bytes"
	"errors"
	"github.com/scim2/tools/schema"
	"strings"
)

type StructGenerator struct {
	w   *genWriter
	s   schema.ReferenceSchema
	ptr bool
}

func NewStructGenerator(s schema.ReferenceSchema) (StructGenerator, error) {
	if s.Name == "" {
		return StructGenerator{}, errors.New("schema does not have a name")
	}

	return StructGenerator{
		w: newGenWriter(&bytes.Buffer{}),
		s: s,
	}, nil
}

// UsePtr indicates whether the generator will use pointers if the attribute is not required.
func (g *StructGenerator) UsePtr(t bool) *StructGenerator {
	g.ptr = t
	return g
}

// Generate creates a buffer with a go representation of the resource described in the given schema.
func (g *StructGenerator) Generate() *bytes.Buffer {
	g.generateStruct(g.s.Name, g.s.Description, g.s.Attributes)
	return g.w.writer.(*bytes.Buffer)
}

func (g *StructGenerator) generateStruct(name, desc string, attrs []*schema.Attribute) {
	w := g.w

	name = keepAlpha(name) // remove all non alpha characters

	if desc != "" {
		w.ln(comment(wrap(desc, 117))) // 120 - "// "
	}

	if len(attrs) == 0 {
		w.lnf("type %s struct {}", name)
		return
	}

	w.lnf("type %s struct {", name)
	g.generateStructFields(name, attrs)
	w.ln("}")

	for _, attr := range attrs {
		if attr.Type == schema.ComplexType {
			typ := cap(attr.Name)
			if attr.MultiValued && strings.HasSuffix(typ, "s") {
				typ = strings.TrimSuffix(typ, "s")
			}
			w.n()
			g.generateStruct(name+typ, attr.Description, attr.SubAttributes)
		}
	}
}

func (g *StructGenerator) generateStructFields(name string, attrs []*schema.Attribute) {
	w := g.w

	name = keepAlpha(name) // remove all non alpha characters

	// get longest name to indent fields.
	var indent int
	for _, attr := range attrs {
		if l := len(cap(attr.Name)); l > indent {
			indent = l
		}
	}

	for _, attr := range attrs {
		var typ string
		switch t := attr.Type; t {
		case "decimal":
			typ = "float64"
		case "integer":
			typ = "int"
		case "boolean":
			typ = "bool"
		case "complex":
			typ = cap(name + cap(attr.Name))
		default:
			typ = "string"
		}

		// field name
		name := cap(keepAlpha(attr.Name))
		w.in(4).w(name)
		w.sp(indent - len(name) + 1)

		if attr.MultiValued {
			w.w("[]")
			if strings.HasSuffix(typ, "s") {
				typ = strings.TrimSuffix(typ, "s")
			}
		} else if !attr.Required && g.ptr {
			w.w("*")
		}

		w.ln(typ)
	}
}
