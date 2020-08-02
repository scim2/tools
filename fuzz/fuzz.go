package fuzz

import (
	"encoding/base64"
	"fmt"
	"github.com/google/gofuzz"
	"math/rand"
	"time"
)

type Fuzzer struct {
	schema ReferenceSchema

	fuzzer *fuzz.Fuzzer
	r      *rand.Rand

	emptyChance float64
	minElements int
	maxElements int
}

// New returns a new Fuzzer.
func New(schema ReferenceSchema) *Fuzzer {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	return &Fuzzer{
		schema:      schema,
		fuzzer:      fuzz.New().RandSource(r),
		r:           r,
		emptyChance: .2,
		minElements: 1,
		maxElements: 10,
	}
}

// RandSource causes the Fuzzer to get values from the given source of randomness.
func (f *Fuzzer) RandSource(s rand.Source) *Fuzzer {
	f.r = rand.New(s)
	f.fuzzer.RandSource(s)
	return f
}

// EmptyChance sets the probability of creating an empty field map to the given chance.
// The chance should be between 0 (no empty fields) and 1 (all empty), inclusive.
func (f *Fuzzer) EmptyChance(p float64) *Fuzzer {
	if p < 0 || p > 1 {
		panic("p should be between 0 and 1, inclusive.")
	}
	f.emptyChance = p
	return f
}

// NumElements sets the minimum and maximum number of elements that will be added.
// If the elements are not required, it is possible to get less elements than the given parameter.
func (f *Fuzzer) NumElements(atLeast, atMost int) *Fuzzer {
	if atLeast > atMost {
		panic("atLeast must be <= atMost")
	}
	if atLeast < 0 {
		panic("atLeast must be >= 0")
	}
	f.minElements = atLeast
	f.maxElements = atMost
	return f
}

// NeverEmpty makes sure that all passed attribute names are never empty during fuzzing.
// Setting a complex attribute on never empty will also make their sub attributes never empty.
// i.e. "displayName", "name.givenName" or "emails.value"
func (f *Fuzzer) NeverEmpty(names ...string) *Fuzzer {
	for _, attribute := range f.schema.Attributes {
		for _, name := range names {
			attribute.neverEmpty(name)
		}
	}
	return f
}

func (f *Fuzzer) elementCount() int {
	if f.minElements == f.maxElements {
		return f.minElements
	}
	return f.minElements + f.r.Intn(f.maxElements-f.minElements+1)
}

func (f *Fuzzer) shouldFill() bool {
	return f.r.Float64() > f.emptyChance
}

// Fuzz recursively fills a Resource based on fields the ReferenceSchema of the Fuzzer.
func (f *Fuzzer) Fuzz() Resource {
	var resource Resource
	f.fuzzer.Funcs(f.newResourceFuzzer()).Fuzz(&resource)
	return resource
}

func (f *Fuzzer) newResourceFuzzer() func(resource *Resource, c fuzz.Continue) {
	return func(r *Resource, c fuzz.Continue) {
		resource := make(Resource)
		for _, attribute := range f.schema.Attributes {
			f.fuzzAttribute(resource, attribute, c)
		}
		*r = resource
	}
}

func (f *Fuzzer) fuzzAttribute(resource map[string]interface{}, attribute *Attribute, c fuzz.Continue) {
	if attribute.MultiValued {
		var elements []interface{}
		for i := 0; i < f.elementCount(); i++ {
			value := f.fuzzSingleAttribute(attribute, c)
			if value != nil {
				elements = append(elements, f.fuzzSingleAttribute(attribute, c))
			}
		}
		if len(elements) != 0 {
			resource[attribute.Name] = elements
		}
		return
	}

	if attribute.shouldFill() || f.shouldFill() {
		value := f.fuzzSingleAttribute(attribute, c)
		if value != nil {
			resource[attribute.Name] = f.fuzzSingleAttribute(attribute, c)
		}
	}
}

func (f *Fuzzer) fuzzSingleAttribute(attribute *Attribute, c fuzz.Continue) interface{} {
	switch attribute.Type {
	case StringType, ReferenceType:
		var randString string
		if len(attribute.CanonicalValues) == 0 {
			randString = randAlphaString(c.Rand, 10)
		} else {
			randString = randStringFromSlice(c.Rand, attribute.CanonicalValues)
		}
		return randString
	case BooleanType:
		randBool := c.RandBool()
		return randBool
	case BinaryType:
		randBase64String := base64.StdEncoding.EncodeToString([]byte(randAlphaString(c.Rand, 10)))
		return randBase64String
	case DecimalType:
		var randFloat64 float64
		c.Fuzz(&randFloat64)
		return randFloat64
	case IntegerType:
		var randInt int
		c.Fuzz(&randInt)
		return randInt
	case DateTimeType:
		randDateTimeString := randDateTime()
		return randDateTimeString
	case ComplexType:
		complexResource := make(map[string]interface{})
		for _, subAttribute := range attribute.SubAttributes {
			f.fuzzAttribute(complexResource, subAttribute, c)
		}
		if len(complexResource) == 0 {
			return nil
		}
		return complexResource
	default:
		panic(fmt.Sprintf("unknown attribute type %s", attribute.Type))
	}
}
