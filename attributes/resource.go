package attributes

import (
	"fmt"
	"reflect"
	"strings"
)

// Add stores the given key and value, returns an error if the key in use.
func Add(resource map[string]interface{}, key string, value interface{}) error {
	if err := validKey(resource, key); err != nil {
		return err
	}

	if _, ok := resource[key]; ok {
		return fmt.Errorf("duplicate key: %s", key)
	}
	resource[key] = value
	return nil
}

// AddEmptyComplexAttribute creates a new map with given key, returns an error if the key in use.
func AddEmptyComplexAttribute(resource map[string]interface{}, key string) (map[string]interface{}, error) {
	if err := validKey(resource, key); err != nil {
		return nil, err
	}

	if err := Add(resource, key, make(map[string]interface{})); err != nil {
		return nil, err
	}
	return resource[key].(map[string]interface{}), nil
}

// AppendComplexMultiValuedAttribute adds given complex value to the complex multi valued attribute with the given key.
// Tries to fill up nil values before appending.
// Returns an error IF:
// - the slice is not present.
// - it is not a slice of maps.
func AppendComplexMultiValuedAttribute(resource map[string]interface{}, key string, value map[string]interface{}) error {
	if err := validKey(resource, key); err != nil {
		return err
	}

	if resourceValue, found := resource[key]; found {
		if sliceValue, isSlice := resourceValue.([]map[string]interface{}); isSlice {
			for k, e := range value {
				var filled bool
				for _, r := range sliceValue {
					if Add(r, k, e) == nil {
						filled = true
						break
					}
				}
				if !filled {
					sliceValue = append(sliceValue, map[string]interface{}{
						k: e,
					})
				}
			}
			resource[key] = sliceValue
			return nil
		}
		return fmt.Errorf("key value was not complex and multi valued: %s %s", key, value)
	}
	return fmt.Errorf("key not found: %s", key)
}

// AppendMultiValuedAttribute adds given value to the multi valued attribute with the given key.
// Tries to fill up nil values before appending.
// Returns an error IF:
// - the slice is not present.
// - it is not a slice.
// - the value type does not match the type of the content.
func AppendMultiValuedAttribute(resource map[string]interface{}, key string, value interface{}) error {
	if err := validKey(resource, key); err != nil {
		return err
	}

	if resourceValue, found := resource[key]; found {
		if sliceValue, isSlice := resourceValue.([]interface{}); isSlice {
			if len(sliceValue) != 0 {
				elementType := reflect.TypeOf(sliceValue[0])
				if t := reflect.TypeOf(value); t != elementType {
					return fmt.Errorf("type does not match %s slice type: %s", elementType, value)
				}
			}

			// Try to fill nil values.
			for i := 0; i < len(sliceValue); i++ {
				if sliceValue[i] == nil {
					sliceValue[i] = value
					resource[key] = sliceValue
					return nil
				}
			}

			// Otherwise append value to slice.
			resource[key] = append(sliceValue, value)
			return nil
		}
		return fmt.Errorf("key value was not multi valued: %s %s", key, value)
	}
	return fmt.Errorf("key not found: %s", key)
}

// Depth returns the amount of nested maps.
func Depth(resource map[string]interface{}) int {
	var depth int
	for _, attribute := range resource {
		switch attribute := attribute.(type) {
		case map[string]interface{}:
			if d := Depth(attribute); depth < d {
				depth = d
			}
		case []interface{}:
			for _, v := range attribute {
				if r, ok := v.(map[string]interface{}); ok {
					if d := Depth(r); depth < d {
						depth = d
					}
				}
			}
		}
	}
	return 1 + depth
}

// EnsureComplexAttribute gets the complex attribute based on the given key.
// IF not present -> it creates an empty map.
// IF not a map   -> it overwrites it with an empty map.
func EnsureComplexAttribute(resource map[string]interface{}, key string) map[string]interface{} {
	if value, found := resource[key]; found {
		if mapValue, isMap := value.(map[string]interface{}); isMap {
			return mapValue
		}
	}
	resource[key] = make(map[string]interface{})
	return resource[key].(map[string]interface{})
}

// EnsureComplexMultiValuedAttribute gets the complex multi valued attribute based on the given key.
// Makes sure the size is at least the given length (by appending empty maps if smaller).
// IF not present -> it creates an empty slice of empty maps of given length.
// IF not a slice -> it overwrites it with an empty slice of empty maps of given length.
func EnsureComplexMultiValuedAttribute(resource map[string]interface{}, key string, length int) []map[string]interface{} {
	if value, found := resource[key]; found {
		if sliceValue, isSlice := value.([]map[string]interface{}); isSlice {
			for len(sliceValue) <= length {
				sliceValue = append(sliceValue, nil)
			}

			resource[key] = sliceValue
			return sliceValue
		}
	}
	resources := make([]map[string]interface{}, 0)
	for i := 0; i < length; i++ {
		resources = append(resources, make(map[string]interface{}))
	}
	resource[key] = resources
	return resources
}

// EnsureMultiValuedAttribute gets the multi valued attribute based on the given key.
// Makes sure the size is at least the given length (by appending nil if smaller).
// IF not present -> it creates an empty slice of given length.
// IF not a map   -> it overwrites it with an empty slice of given length.
func EnsureMultiValuedAttribute(resource map[string]interface{}, key string, length int) []interface{} {
	if value, found := resource[key]; found {
		if sliceValue, isSlice := value.([]interface{}); isSlice {
			for len(sliceValue) <= length {
				sliceValue = append(sliceValue, nil)
			}
			return sliceValue
		}
	}
	resource[key] = make([]interface{}, length)
	return resource[key].([]interface{})
}

// Exists checks whether the key exists in the map.
func Exists(resource map[string]interface{}, key string) bool {
	if validKey(resource, key) != nil {
		return true
	}

	_, ok := resource[key]
	return ok
}

// validKey checks whether there is another case insensitive key with the same value.
// i.e. ("x", "X") -> false
//		("x", "x") -> true
//		("x", "y") -> true
func validKey(resource map[string]interface{}, key string) error {
	for k := range resource {
		if strings.EqualFold(k, key) && k != key {
			return fmt.Errorf("duplicate keys: %s and %s", k, key)
		}
	}
	return nil
}
