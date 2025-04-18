
package eval

import (
	"strings"

	"github.com/megaease/easeprobe/global"
)

// -----------------------------------------------------------------------------

// DocType is the different type of document
type DocType int

// The Document Type
const (
	Unsupported DocType = iota
	HTML
	XML
	JSON
	TEXT
)

var docTypeToStr = map[DocType]string{
	Unsupported: "unsupported",
	HTML:        "html",
	XML:         "xml",
	JSON:        "json",
	TEXT:        "text",
}

var strToDocType = global.ReverseMap(docTypeToStr)

// String covert the DocType to string
func (t DocType) String() string {
	return docTypeToStr[t]
}

// Type covert the string to Type
func (t *DocType) Type(s string) {
	*t = strToDocType[strings.ToLower(s)]
}

// MarshalYAML is marshal the type
func (t DocType) MarshalYAML() (interface{}, error) {
	return global.EnumMarshalYaml(docTypeToStr, t, "DocType")
}

// UnmarshalYAML is unmarshal the type
func (t *DocType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return global.EnumUnmarshalYaml(unmarshal, strToDocType, t, Unsupported, "DocType")
}

// -----------------------------------------------------------------------------

// VarType is an enum for the different types of values
type VarType int

// The value types
const (
	Unknown VarType = iota
	Int
	Float
	String
	Bool
	Time
	Duration
)

var varTypeToStr = map[VarType]string{
	Unknown:  "unknown",
	Int:      "int",
	Float:    "float",
	String:   "string",
	Bool:     "bool",
	Time:     "time",
	Duration: "duration",
}

var strToVarType = global.ReverseMap(varTypeToStr)

// String covert the Type to string
func (t VarType) String() string {
	return varTypeToStr[t]
}

// Type covert the string to Type
func (t *VarType) Type(s string) {
	*t = strToVarType[strings.ToLower(s)]
}

// MarshalYAML is marshal the type
func (t VarType) MarshalYAML() (interface{}, error) {
	return global.EnumMarshalYaml(varTypeToStr, t, "Variable")
}

// UnmarshalYAML is unmarshal the type
func (t *VarType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return global.EnumUnmarshalYaml(unmarshal, strToVarType, t, Unknown, "Variable")
}
