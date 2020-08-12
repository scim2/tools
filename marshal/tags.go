package marshal

import (
	"reflect"
	"strconv"
	"strings"
)

type tag struct {
	name        string
	multiValued bool
	indexes     []int
	allowZero   bool
	ignore      bool
	sub         *tag
}

func parseTags(field reflect.StructField) tag {
	var t tag

	scimTag := field.Tag.Get("scim")
	tags := strings.Split(scimTag, ",")
	if scimTag == "" {
		t.name = lowerFirstRune(field.Name)
	} else {
		parts := strings.Split(tags[0], ".")
		if parts[0] != "" {
			t.name = parts[0]
		} else {
			t.name = lowerFirstRune(field.Name)
		}

		if len(parts) > 1 {
			t.sub = &tag{
				name: parts[1],
			}
		}
	}

	if len(tags) > 1 {
		for _, option := range tags[1:] {
			var sub bool
			if strings.HasPrefix(option, "_") {
				if t.sub == nil {
					continue
				}
				option = strings.TrimPrefix(option, "_")
				sub = true
			}

			switch option {
			case "multiValued", "mV":
				if !sub {
					t.multiValued = true
				} else {
					t.sub.multiValued = true
				}
			case "zero", "0":
				if !sub {
					t.allowZero = true
				} else {
					t.sub.allowZero = true
				}
			case "ignore", "!":
				if !sub {
					t.ignore = true
				} else {
					t.sub.ignore = true
				}
			}

			if strings.HasPrefix(option, "index=") || strings.HasPrefix(option, "i=") {
				option = strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(option, "i"), "ndex"), "=")
				if option == "all" {
					if !sub {
						t.indexes = []int{-1}
					} else {
						t.sub.indexes = []int{-1}
					}
					break
				}
				for _, v := range strings.Split(option, ";") {
					if !strings.Contains(v, "-") {
						if i, err := strconv.Atoi(v); err == nil {
							if 0 <= i {
								if !sub {
									t.indexes = append(t.indexes, i)
								} else {
									t.sub.indexes = append(t.indexes, i)
								}
							}
						}
					} else {
						parts := strings.Split(v, "-")
						if len(parts) != 2 {
							break
						}
						from, err := strconv.Atoi(parts[0])
						if err != nil {
							break
						}
						to, err := strconv.Atoi(parts[1])
						if err != nil {
							break
						}
						for i := from; i <= to; i++ {
							if !sub {
								t.indexes = append(t.indexes, i)
							} else {
								t.sub.indexes = append(t.indexes, i)
							}
						}
					}
				}
			}
		}
	}

	return t
}

func (t tag) all() bool {
	return len(t.indexes) == 1 && t.indexes[0] == -1
}

func (t tag) max() int {
	var max int
	for _, i := range t.indexes {
		if i > max {
			max = i
		}
	}
	return max
}
