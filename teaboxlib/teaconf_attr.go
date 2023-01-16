package teaboxlib

import (
	"fmt"
	"strconv"
	"strings"
)

/*
Argument's context attributes. They are just an array of attributes
that can be anything specific to a widget.

Example for a tabular widget:
	attributes:
	  # Adds filtering box as you type.
	  - search

	  # Adds first column with [*] for selected
	  - multiselect

	  # Sets default height
	  - height:5

	  # Expand 1st column
	  - 1:expand

	  # 3rd column is actually hidden and is serving as a value (i.e. its contents)
	  # will be returned in case the whole row is selected
	  - 3:value,hidden
*/

type TeaConfArgAttributes struct {
	kwa map[string][]string
	opt []string
}

// NewTeaConfArgAttributes constructor
func NewTeaConfArgAttributes(attr []interface{}) *TeaConfArgAttributes {
	if attr == nil {
		attr = []interface{}{}
	}

	return (&TeaConfArgAttributes{
		kwa: map[string][]string{},
		opt: []string{},
	}).parse(attr)
}

// Parse the attributes
func (tca *TeaConfArgAttributes) parse(attrs []interface{}) *TeaConfArgAttributes {
	for _, rAttr := range attrs {
		attr, ok := rAttr.(string)
		if !ok {
			continue
		}

		if strings.Contains(attr, "=") {
			tca.addKwArg(attr)
		} else {
			tca.opt = append(tca.opt, attr)
		}
	}

	return tca
}

func (tca *TeaConfArgAttributes) addKwArg(kwa string) {
	args := strings.SplitN(kwa, "=", 2)
	if len(args) != 2 {
		panic(fmt.Sprintf("Wrong keyword argument: %s", kwa))
	}

	// Parse lists
	for _, key := range strings.Split(strings.ReplaceAll(strings.TrimSpace(args[0]), " ", ""), ",") {
		for _, attr := range strings.Split(args[1], ",") {
			_, exist := tca.kwa[key]
			if !exist {
				tca.kwa[key] = []string{}
			}
			tca.kwa[key] = append(tca.kwa[key], strings.TrimSpace(attr))
		}
	}
}

// KeywordHasAllAttrs checks if a keyword has all the following attributes
func (tca *TeaConfArgAttributes) KeywordHasAllAttrs(key string, attrs ...string) bool {
	opts, exist := tca.kwa[key]
	if !exist {
		return false
	}

	// if all
	amt := len(attrs)
	offset := 0
	for _, attr := range attrs {
		for _, opt := range opts {
			if opt == attr {
				offset++
			}
		}
	}

	return amt == offset
}

// KeywordHasAnyAttrs checks if a keyword has all the following attributes
func (tca *TeaConfArgAttributes) KeywordHasAnyAttrs(key string, attrs ...string) bool {
	opts, exist := tca.kwa[key]
	if !exist {
		return false
	}

	// If at least one
	for _, attr := range attrs {
		for _, opt := range opts {
			if opt == attr {
				return true
			}
		}
	}

	return false
}

// KeywordValueAsString returns a value as string. If no value, an empty string returned.
// If there is more than one values, only the first one is returned
func (tca *TeaConfArgAttributes) KeywordValueAsString(key string) string {
	v := tca.getValue(key)
	if len(v) == 0 {
		return ""
	}

	return v[0]
}

// KeywordValueAsStrings returns a value as an array of strings. If no value, an empty array of strings returned.
func (tca *TeaConfArgAttributes) KeywordValueAsStrings(key string) []string {
	v := tca.getValue(key)
	if v == nil {
		v = []string{}
	}

	return v
}

func (tca *TeaConfArgAttributes) KeywordValueAsInt(key string) int {
	v, err := strconv.Atoi(tca.KeywordValueAsString(key))
	if err != nil {
		return -1
	}

	return v
}

// KeywordValueAsInts returns a value as an array of integers. If no value or an error, an empty array of integers returned.
func (tca *TeaConfArgAttributes) KeywordValueAsInts(key string) []int {
	ret := []int{}
	for _, rw := range tca.KeywordValueAsStrings(key) {
		v, err := strconv.Atoi(rw)
		if err != nil {
			return []int{}
		}
		ret = append(ret, v)
	}
	return ret
}

func (tca *TeaConfArgAttributes) getValue(key string) []string {
	opts, exist := tca.kwa[key]
	if !exist {
		return nil
	}

	return opts
}

// HasOption returns true of there is an option like that. :)
func (tca *TeaConfArgAttributes) HasOption(opt string) bool {
	for _, flag := range tca.opt {
		if flag == opt {
			return true
		}
	}

	return false
}
