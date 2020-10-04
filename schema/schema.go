package schema

var (
	SchemasAttribute = &Attribute{
		MultiValued: true,
		Mutability:  Immutable,
		Name:        "schemas",
		Type:        StringType,
		Required:    true,
	}
	IDAttribute = &Attribute{
		CaseExact:   true,
		Description: "A unique identifier for a SCIM resource as defined by the service provider.",
		Mutability:  ReadOnly,
		Name:        "id",
		Type:        StringType,
		Required:    true,
		Returned:    Always,
		Uniqueness:  Server,
	}
	ExternalIDAttribute = &Attribute{
		CaseExact:   true,
		Description: "A String that is an identifier for the resource as defined by the\nprovisioning client.",
		Name:        "externalId",
		Type:        StringType,
	}
	MetaAttribute = &Attribute{
		Description: "A complex attribute containing resource metadata.",
		Mutability:  ReadOnly,
		Name:        "meta",
		Type:        ComplexType,
		SubAttributes: []*Attribute{
			{
				CaseExact:   true,
				Description: "The name of the resource type of the resource.",
				Mutability:  ReadOnly,
				Name:        "resourceType",
			},
			{
				Description: "The DateTime that the resource was added to the service provider.",
				Mutability:  ReadOnly,
				Name:        "created",
			},
			{
				Description: "The most recent DateTime that the details of this resource were updated at the service provider.",
				Mutability:  ReadOnly,
				Name:        "lastModified",
			},
			{
				Description: "The URI of the resource being returned.",
				Mutability:  ReadOnly,
				Name:        "location",
			},
			{
				CaseExact:   true,
				Description: "The version of the resource being returned.",
				Mutability:  ReadOnly,
				Name:        "version",
			},
		},
	}

	CoreAttributes = []*Attribute{
		SchemasAttribute,
		IDAttribute,
		ExternalIDAttribute,
		MetaAttribute,
	}
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
}

// ForEachAttribute calls given function on itself all sub attributes recursively.
func (attribute *Attribute) ForEachAttribute(f func(attribute *Attribute)) {
	f(attribute)
	if attribute.Type == ComplexType {
		for _, subAttribute := range attribute.SubAttributes {
			subAttribute.ForEachAttribute(f)
		}
	}
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
func (s ReferenceSchema) ForEachAttribute(f func(attribute *Attribute)) {
	for _, attribute := range s.Attributes {
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
