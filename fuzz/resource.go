package fuzz

import (
	"encoding/base64"
	"github.com/google/gofuzz"
)

type (
	Resource       map[string]interface{}
	ResourceFuzzer func(resource *Resource, c fuzz.Continue)
)

func NewResourceFuzzer(reference ReferenceSchema) ResourceFuzzer {
	return func(r *Resource, c fuzz.Continue) {
		resource := make(Resource)
		var fuzzAttribute func(resource map[string]interface{}, attribute *Attribute)
		fuzzAttribute = func(resource map[string]interface{}, attribute *Attribute) {
			switch attribute.Type {
			case StringType, ReferenceType:
				randString := randAlphaString(c.Rand, 10)
				if attribute.isRequired() {
					if attribute.MultiValued {
						resource[attribute.Name] = []string{
							randString,
						}
					} else {
						resource[attribute.Name] = randString
					}
				}
			case BooleanType:
				randBool := c.RandBool()
				if attribute.isRequired() {
					if attribute.MultiValued {
						resource[attribute.Name] = []bool{
							randBool,
						}
					} else {
						resource[attribute.Name] = randBool
					}
				}
			case BinaryType:
				randBase64String := base64.StdEncoding.EncodeToString([]byte(randAlphaString(c.Rand, 10)))
				if attribute.isRequired() {
					if attribute.MultiValued {
						resource[attribute.Name] = []string{
							randBase64String,
						}
					} else {
						resource[attribute.Name] = randBase64String
					}
				}
			case DecimalType:
				var randFloat64 float64
				c.Fuzz(&randFloat64)
				if attribute.isRequired() {
					if attribute.MultiValued {
						resource[attribute.Name] = []float64{
							randFloat64,
						}
					} else {
						resource[attribute.Name] = randFloat64
					}
				}
			case IntegerType:
				var randInt int
				c.Fuzz(&randInt)
				if attribute.isRequired() {
					if attribute.MultiValued {
						resource[attribute.Name] = []int{
							randInt,
						}
					} else {
						resource[attribute.Name] = randInt
					}
				}
			case DateTimeType:
				randDateTimeString := randDateTime()
				if attribute.isRequired() {
					if attribute.MultiValued {
						resource[attribute.Name] = []string{
							randDateTimeString,
						}
					} else {
						resource[attribute.Name] = randDateTimeString
					}
				}
			case ComplexType:
				complexResource := make(map[string]interface{})
				for _, subAttribute := range attribute.SubAttributes {
					fuzzAttribute(complexResource, subAttribute)
				}
				if len(complexResource) != 0 {
					if attribute.MultiValued {
						resource[attribute.Name] = []map[string]interface{}{
							complexResource,
						}
					} else {
						resource[attribute.Name] = complexResource
					}
				}
			}
		}

		for _, attribute := range reference.Attributes {
			fuzzAttribute(resource, attribute)
		}

		*r = resource
	}
}
