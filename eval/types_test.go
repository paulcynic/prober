
package eval

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func testVarType(t *testing.T, expect VarType, str string) {
	var result VarType
	result.Type(str)
	assert.Equal(t, expect, result)
}

func TestType(t *testing.T) {
	testVarType(t, Unknown, "unknown")
	testVarType(t, Int, "int")
	testVarType(t, Float, "float")
	testVarType(t, String, "string")
	testVarType(t, Bool, "bool")
	testVarType(t, Time, "time")
	testVarType(t, Duration, "duration")
	testVarType(t, Unknown, "unrecognized")
}

func testVarTypeYAML(t *testing.T, expect VarType, str string) {
	var result VarType
	buf, err := yaml.Marshal(expect)
	assert.Nil(t, err)
	assert.Equal(t, str, string(buf))

	assert.Nil(t, yaml.Unmarshal(buf, &result))
	assert.Equal(t, expect, result)
}

func TestVarTypeYAML(t *testing.T) {
	testVarTypeYAML(t, Unknown, "unknown\n")
	testVarTypeYAML(t, Int, "int\n")
	testVarTypeYAML(t, Float, "float\n")
	testVarTypeYAML(t, String, "string\n")
	testVarTypeYAML(t, Bool, "bool\n")
	testVarTypeYAML(t, Time, "time\n")
	testVarTypeYAML(t, Duration, "duration\n")

	str := "-name:: value\n"
	var result VarType
	assert.NotNil(t, yaml.Unmarshal([]byte(str), &result))
}

//------------------------------------------------------------------------------

func testDocType(t *testing.T, expect DocType, str string) {
	var result DocType
	result.Type(str)
	assert.Equal(t, expect, result)
}

func TestDocType(t *testing.T) {
	testDocType(t, Unsupported, "unsupported")
	testDocType(t, HTML, "html")
	testDocType(t, XML, "xml")
	testDocType(t, JSON, "json")
	testDocType(t, TEXT, "text")
}

func testDocTypeYAML(t *testing.T, expect DocType, str string) {
	var result DocType
	buf, err := yaml.Marshal(expect)
	assert.Nil(t, err)
	assert.Equal(t, str, string(buf))

	assert.Nil(t, yaml.Unmarshal(buf, &result))
	assert.Equal(t, expect, result)
}

func TestDocTypeYAML(t *testing.T) {
	testDocTypeYAML(t, Unsupported, "unsupported\n")
	testDocTypeYAML(t, HTML, "html\n")
	testDocTypeYAML(t, XML, "xml\n")
	testDocTypeYAML(t, JSON, "json\n")
	testDocTypeYAML(t, TEXT, "text\n")

	str := "-name:: value\n"
	var result DocType
	assert.NotNil(t, yaml.Unmarshal([]byte(str), &result))
}
