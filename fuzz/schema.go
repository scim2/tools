package fuzz

import (
	"strings"
)

// Attribute represents an attribute of a ReferenceSchema.
type Attribute struct {
	Name            string       `json:"name"`
	Type            Type         `json:"type"`
	SubAttributes   []*Attribute `json:"subAttributes,omitempty"`
	MultiValued     bool         `json:"multiValued"`
	Description     string       `json:"description,omitempty"`
	Required        bool         `json:"required"`
	CanonicalValues []string     `json:"canonicalValues,omitempty"`
	CaseExact       bool         `json:"caseExact"`
	Mutability      Mutability   `json:"mutability"`
	Returned        Returned     `json:"returned"`
	Uniqueness      Uniqueness   `json:"uniqueness"`
	ReferenceTypes  []string     `json:"referenceTypes"`

	required bool // manually set for fuzzer (schema.NeverEmpty)
}

// ForEachAttribute calls given function on itself all sub attributes recursively.
func (attribute *Attribute) ForEachAttribute(f func(attribute *Attribute)) {
	f(attribute)
	if attribute.isComplex() {
		for _, subAttribute := range attribute.SubAttributes {
			subAttribute.ForEachAttribute(f)
		}
	}
}

func (attribute *Attribute) isComplex() bool {
	return attribute.Type == ComplexType
}

func (attribute *Attribute) neverEmpty(name string) {
	n := strings.SplitN(name, ".", 2)
	if strings.EqualFold(n[0], attribute.Name) {
		if len(n) > 1 && attribute.isComplex() {
			for _, subAttribute := range attribute.SubAttributes {
				subAttribute.neverEmpty(n[1])
			}
		} else {
			attribute.required = true
			if attribute.isComplex() {
				attribute.ForEachAttribute(func(attribute *Attribute) {
					attribute.required = true
				})
			}
		}
	}
}

func (attribute *Attribute) shouldFill() bool {
	return attribute.Required || attribute.required || attribute.isComplex()
}

type Mutability string

const (
	ReadOnly  Mutability = "readOnly"
	ReadWrite Mutability = "readWrite"
	Immutable Mutability = "immutable"
	WriteOnly Mutability = "writeOnly"
)

// ReferenceSchema represents a resource schema that is used to fuzz resources that are defined by this schema.
type ReferenceSchema struct {
	ID          string       `json:"id"`
	Name        string       `json:"name,omitempty"`
	Description string       `json:"description,omitempty"`
	Attributes  []*Attribute `json:"attributes"`
}

// ForEachAttribute calls given function on all attributes recursively.
func (schema ReferenceSchema) ForEachAttribute(f func(attribute *Attribute)) {
	for _, attribute := range schema.Attributes {
		attribute.ForEachAttribute(f)
	}
}

type Returned string

const (
	Always  Returned = "always"
	Never   Returned = "never"
	Default Returned = "default"
	Request Returned = "request"
)

type Type string

const (
	StringType    Type = "string"
	BooleanType   Type = "boolean"
	BinaryType    Type = "binary"
	DecimalType   Type = "decimal"
	IntegerType   Type = "integer"
	DateTimeType  Type = "dateTime"
	ReferenceType Type = "reference"
	ComplexType   Type = "complex"
)

type Uniqueness string

const (
	None   Uniqueness = "none"
	Server Uniqueness = "server"
	Global Uniqueness = "global"
)
