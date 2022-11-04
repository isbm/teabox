package teaboxlib

import (
	"encoding/json"
	"fmt"
	"strings"
)

/*
TeaboxAPICall is a parsed API call, which accepts the following string format:

	<CLASS>:[TYPE]:[{KEY}]<VALUE>

The class, type and the payload should be separated by a colon ":" without a whitespace.

Types:

	string (default, can be just omitted)
	bool
	int
	json

Most of the time the value is just a string. If type is not recognised, then it is a string. :-)
Example sending a typical string to the logger status widget (two are equivalent):

	LOGGER-STATUS:string:Hello world!
	LOGGER-STATUS::Hello world!

Example increment progress bar by one allocated position. This does not need no type or value at all.
Just a notification that a progess-bar needs to be incremented by one of the allocated parts, e.g.
"jump to next of 10 allocated parts", which will cause increase progress bar by 10%:

	PROGRESS-INCREMENT::

But it can be very specific:

	PROGRESS-UPDATE:int:42

This will cause setting progress bar to 42%.

Constructions:

APIs also support simple key/values. Key should be always a string with no whitespace,
enclosed with "{" and "}". The <TYPE> of the call corresponds to the <VALUE>. Example:

	CHECKLIST-DONE:bool:{disk-setup}true
	CHECKLIST-DONE:bool:{windows-installation}false

JSON is used for more complex data exchange, like tables. In this case, json is
string-encoded with proper escaping, e.g.:

	PKG-SELECT:json:{"key": "properly \\"escaped\\" value", "2": 3}

This data structure is supported for exceptional cases and should be used sparingly. Because there
is no need to overuse JSON everywhere, since sometimes things needs to be just as simple as possible. :)
*/
type TeaboxAPICall struct {
	class    string
	datatype string
	key      string
	payload  interface{}
}

func NewTeaboxAPICall(data []byte) *TeaboxAPICall {
	ac := new(TeaboxAPICall)
	ac.parse(data)
	return ac
}

func (ac *TeaboxAPICall) parse(data []byte) {
	tokens := strings.SplitN(strings.TrimSpace(string(data)), ":", 3)
	if len(tokens) != 3 {
		return
	}

	// Set API class
	ac.class = strings.ToUpper(tokens[0])

	// Set supported types
	switch tokens[1] {
	case "bool", "int", "json":
		ac.datatype = tokens[1]
	default:
		ac.datatype = "string"
	}

	// Parse payload
	if strings.HasPrefix(tokens[2], "{") {
		keyval := strings.SplitN(tokens[2], "}", 2)
		if len(keyval) == 2 {
			ac.key = keyval[0][1:] // set a key, trim leading "{"
			ac.payload = keyval[1]
		} else {
			ac.payload = fmt.Sprintf("%v", keyval)
		}
	} else {
		ac.payload = tokens[2]
	}
}

// GetClass of the API call (address to what section)
func (ac *TeaboxAPICall) GetClass() string {
	return ac.class
}

// GetType returns a type of the payload to cast to
func (ac *TeaboxAPICall) GetType() string {
	return ac.datatype
}

// GetKey if any. Sometimes API payload call corresponds to a special key
func (ac *TeaboxAPICall) GetKey() string {
	return ac.key
}

// GetValue of the API payload call
func (ac *TeaboxAPICall) GetValue() interface{} {
	return ac.payload
}

// GetString returns value as a string type, but only if datatype is "string".
// Otherwise it will be always an empty string.
func (ac *TeaboxAPICall) GetString() string {
	if ac.GetType() == "string" {
		return ac.GetValue().(string)
	}
	return ""
}

// GetBool returns value as a boolean type, but only if datatype is "bool".
// Otherwise it will be always "false".
func (ac *TeaboxAPICall) GetBool() bool {
	if ac.GetType() == "bool" {
		v := strings.ToLower(fmt.Sprintf("%v", ac.payload))
		return v == "true" || v == "yes"
	}
	return false
}

// GetInt returns value as a boolean type, but only if datatype is "int".
// Otherwise it will be always negative value.
func (ac *TeaboxAPICall) GetInt() int {
	if ac.GetType() == "int" {
		return ac.GetValue().(int)
	}
	return -1
}

// GetJSON returns a value as an interface{} of an arbitrary JSON type.
func (ac *TeaboxAPICall) GetJSON() interface{} {
	var data interface{}
	if ac.GetType() == "json" {
		// Essentially, JSON is a string that needs to be unmarshalled
		// But we are not sure if it is there at first place.
		v, ok := ac.GetValue().(string)
		if ok {
			// We don't care here about error handling at the moment,
			// as there are no real handler to cry about this. Send your PR implementing one!
			_ = json.Unmarshal([]byte(v), data)
		}
	}
	return data
}
