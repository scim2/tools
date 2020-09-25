package marshal

import (
	"fmt"
	"github.com/scim2/tools/structs"
	"testing"
)

func TestComplex(t *testing.T) {
	type structPtrString struct {
		Ptr *string `scim:"ptr"`
	}

	type complex struct {
		Bool    bool                   `scim:"complex.bool"`
		Int     int                    `scim:"complex.int"`
		Array   [2]interface{}         `scim:"complex.array"`
		Map     map[string]interface{} `scim:"complex.map"`
		Ptr     *string                `scim:"complex.ptr"`
		Slice   []interface{}          `scim:"complex.slice"`
		String  string                 `scim:"complex.string"`
		Struct  structPtrString        `scim:"complex.struct"`
		Structs []structPtrString      `scim:"complex.structs"`
	}

	str := "_"

	t.Run("valid", func(t *testing.T) {
		resource, err := Marshal(complex{
			Bool:   true,
			Int:    1,
			Map:    map[string]interface{}{str: str},
			Ptr:    &str,
			String: str,
			Struct: structPtrString{Ptr: &str},
		})
		if err != nil {
			t.Error(err)
		}
		ref := structs.Resource{
			"complex": structs.Resource{
				"bool":   true,
				"int":    1,
				"map":    structs.Resource{str: str},
				"ptr":    str,
				"string": str,
				"struct": structs.Resource{"ptr": str},
			},
		}
		if fmt.Sprintf("%#v", resource) != fmt.Sprintf("%#v", ref) {
			t.Error(fmt.Sprintf("\n%#v", resource), fmt.Sprintf("\n%#v", ref))
		}
	})

	t.Run("invalid", func(t *testing.T) {
		t.Run("array", func(t *testing.T) {
			if _, err := Marshal(complex{
				Array: [2]interface{}{str},
			}); err == nil {
				t.Error("error expected, got none")
			}
		})

		t.Run("slice", func(t *testing.T) {
			if _, err := Marshal(complex{
				Slice: []interface{}{str},
			}); err == nil {
				t.Error("error expected, got none")
			}
		})

		t.Run("slice", func(t *testing.T) {
			if _, err := Marshal(complex{
				Structs: []structPtrString{
					{Ptr: &str},
				},
			}); err == nil {
				t.Error("error expected, got none")
			}
		})
	})
}

func TestComplexMultiValued(t *testing.T) {
	type structPtrString struct {
		Ptr *string `scim:",mV"`
	}

	type simple struct {
		Bool    bool                   `scim:"complexMV.bool,mV,_mV"`
		Int     int                    `scim:"complexMV1.int,mV"`
		Array   [2]interface{}         `scim:"complexMV2.array,mV"`
		Map     map[string]interface{} `scim:"complexMV3.map,mV"`
		Ptr     *string                `scim:"complex.ptr,_mV"`
		Slice   []interface{}          `scim:"complexMV5.slice,mV,_mV"`
		String  string                 `scim:"complex.string,_mV"`
		Struct  structPtrString        `scim:"complex.struct,_mV"`
		Structs []structPtrString      `scim:"complexMV8.structs,mV,_mV"`
	}

	str := "_"

	resource, err := Marshal(simple{
		Bool:    true,
		Int:     1,
		Array:   [2]interface{}{str, str},
		Map:     map[string]interface{}{str: str},
		Ptr:     &str,
		Slice:   []interface{}{str},
		String:  str,
		Struct:  structPtrString{Ptr: &str},
		Structs: []structPtrString{{Ptr: &str}, {Ptr: &str}},
	})
	if err != nil {
		t.Error(err)
	}
	ref := structs.Resource{
		"complex": structs.Resource{
			"ptr":    str,
			"string": str,
			"struct": structs.Resource{"ptr": []interface{}{str}},
		},
		"complexMV":  []structs.Resource{{"bool": []interface{}{true}}},
		"complexMV1": []structs.Resource{{"int": 1}},
		"complexMV2": []structs.Resource{{"array": str}, {"array": str}},
		"complexMV3": []structs.Resource{{"map": structs.Resource{str: str}}},
		"complexMV5": []structs.Resource{{"slice": []interface{}{str}}},
		"complexMV8": []structs.Resource{
			{"structs": []structs.Resource{{"ptr": []interface{}{str}}}},
			{"structs": []structs.Resource{{"ptr": []interface{}{str}}}},
		},
	}
	if fmt.Sprintf("%#v", resource) != fmt.Sprintf("%#v", ref) {
		t.Error(fmt.Sprintf("\n%#v", resource), fmt.Sprintf("\n%#v", ref))
	}
}

func TestSimple(t *testing.T) {
	type structPtrString struct {
		Ptr *string
	}

	type simple struct {
		Bool    bool
		Int     int
		Array   [2]interface{}
		Map     map[string]interface{}
		Ptr     *string
		Slice   []interface{}
		String  string
		Struct  structPtrString
		Structs []structPtrString
	}
	str := "_"

	t.Run("valid", func(t *testing.T) {
		resource, err := Marshal(simple{
			Bool:   true,
			Int:    1,
			Map:    map[string]interface{}{str: str},
			Ptr:    &str,
			String: str,
			Struct: structPtrString{
				Ptr: &str,
			},
		})
		if err != nil {
			t.Error(err)
		}
		ref := structs.Resource{
			"bool":   true,
			"int":    1,
			"map":    structs.Resource{str: str},
			"ptr":    str,
			"string": str,
			"struct": structs.Resource{"ptr": str},
		}
		if fmt.Sprintf("%#v", resource) != fmt.Sprintf("%#v", ref) {
			t.Error(fmt.Sprintf("\n%#v", resource), fmt.Sprintf("\n%#v", ref))
		}
	})

	t.Run("invalid", func(t *testing.T) {
		t.Run("array", func(t *testing.T) {
			if _, err := Marshal(simple{
				Array: [2]interface{}{str},
			}); err == nil {
				t.Error("error expected, got none")
			}
		})

		t.Run("slice", func(t *testing.T) {
			if _, err := Marshal(simple{
				Slice: []interface{}{str},
			}); err == nil {
				t.Error("error expected, got none")
			}
		})

		t.Run("slice", func(t *testing.T) {
			if _, err := Marshal(simple{
				Structs: []structPtrString{
					{Ptr: &str},
				},
			}); err == nil {
				t.Error("error expected, got none")
			}
		})

		t.Run("map", func(t *testing.T) {
			invalid := []interface{}{
				[]interface{}{str},
				[2]interface{}{str},
				map[string]interface{}{str: str},
				structPtrString{},
			}

			for _, test := range invalid {
				if _, err := Marshal(simple{
					Map: map[string]interface{}{
						"slice": test,
					},
				}); err == nil {
					t.Error("error expected, got none")
				}
			}
		})

		t.Run("nested", func(t *testing.T) {
			type nested struct {
				Name string
				N    *nested
			}

			if _, err := Marshal(nested{
				Name: str,
				N: &nested{
					Name: str,
				},
			}); err != nil {
				t.Errorf("no error expected, got %q", err)
			}

			if _, err := Marshal(nested{
				Name: str,
				N: &nested{
					Name: str,
					N: &nested{
						Name: str,
					},
				},
			}); err == nil {
				t.Error("expected error, got none")
			}
		})
	})
}

func TestSimpleMultiValued(t *testing.T) {
	type structPtrString struct {
		Ptr *string `scim:",mV"`
	}

	type simple struct {
		Bool    bool                   `scim:",mV"`
		Int     int                    `scim:",mV"`
		Array   [2]interface{}         `scim:",mV"`
		Map     map[string]interface{} `scim:",mV"`
		Ptr     *string                `scim:",mV"`
		Slice   []interface{}          `scim:",mV"`
		String  string                 `scim:",mV"`
		Struct  structPtrString        `scim:",mV"`
		Structs []structPtrString      `scim:",mV"`
	}
	str := "_"

	resource, err := Marshal(simple{
		Bool:   true,
		Int:    1,
		Array:  [2]interface{}{str},
		Map:    map[string]interface{}{str: str},
		Ptr:    &str,
		Slice:  []interface{}{str, str},
		String: str,
		Struct: structPtrString{Ptr: &str},
		Structs: []structPtrString{
			{Ptr: &str},
			{Ptr: &str},
		},
	})
	if err != nil {
		t.Error(err)
	}
	ref := structs.Resource{
		"bool":    []interface{}{true},
		"int":     []interface{}{1},
		"array":   []interface{}{str},
		"map":     []structs.Resource{{str: str}},
		"ptr":     []interface{}{str},
		"slice":   []interface{}{str, str},
		"string":  []interface{}{str},
		"struct":  []structs.Resource{{"ptr": []interface{}{str}}},
		"structs": []structs.Resource{{"ptr": []interface{}{str}}, {"ptr": []interface{}{str}}},
	}
	if fmt.Sprintf("%#v", resource) != fmt.Sprintf("%#v", ref) {
		t.Error(fmt.Sprintf("\n%#v", resource), fmt.Sprintf("\n%#v", ref))
	}
}
