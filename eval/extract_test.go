
package eval

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func assertExtractor(t *testing.T, extractor Extractor, query string, vt VarType, expected interface{}, success bool) {
	extractor.SetQuery(query)
	extractor.SetVarType(vt)
	result, err := extractor.Extract()
	if success {
		assert.Nil(t, err)
	} else {
		assert.NotNil(t, err)
	}
	assert.Equal(t, expected, result)
}
func assertExtractorSucc(t *testing.T, extractor Extractor, query string, vt VarType, expected interface{}) {
	assertExtractor(t, extractor, query, vt, expected, true)
}
func assertExtractorFail(t *testing.T, extractor Extractor, query string, vt VarType, expected interface{}) {
	assertExtractor(t, extractor, query, vt, expected, false)
}


func TestJSONExtractor(t *testing.T) {
	jsonDoc := `
	{
		"company": {
			"name": "MegaEase",
			"person": [{
					"name": "Bob",
					"email": "bob@example.com",
					"age": 35,
					"salary": 35000.12,
					"birth": "1984-10-12",
					"work": "40h",
					"fulltime": true
				},
				{
					"name": "Alice",
					"email": "alice@example.com",
					"age": 25,
					"salary": 25000.12,
					"birth": "1985-10-12",
					"work": "30h",
					"fulltime": false
				}
			]
		}
	}`
	extractor := NewJSONExtractor(jsonDoc)

	assertExtractorSucc(t, extractor, "//name", String, "MegaEase")
	assertExtractorSucc(t, extractor, "//company/name", String, "MegaEase")
	assertExtractorSucc(t, extractor, "//email", String, "bob@example.com")
	assertExtractorSucc(t, extractor, "//company/person/*[1]/name", String, "Bob")
	assertExtractorSucc(t, extractor, "//company/person/*[2]/email", String, "alice@example.com")
	assertExtractorSucc(t, extractor, "//company/person/*[last()]/name", String, "Alice")
	assertExtractorSucc(t, extractor, "//company/person/*[last()]/age", Int, 25)
	assertExtractorSucc(t, extractor, "//company/person/*[salary=25000.12]/salary", Float, 25000.12)
	expected, _ := tryParseTime("1984-10-12")
	assertExtractorSucc(t, extractor, "//company/person/*[name='Bob']/birth", Time, expected)
	assertExtractorSucc(t, extractor, "//company/person/*[name='Alice']/work", Duration, 30*time.Hour)
	assertExtractorSucc(t, extractor, "//*/email[contains(.,'bob')]", String, "bob@example.com")
	assertExtractorSucc(t, extractor, "//work", Duration, 40*time.Hour)
	assertExtractorSucc(t, extractor, "//person/*[2]/fulltime", Bool, false)
}


func TestRegexExtractor(t *testing.T) {
	regexDoc := `name: Bob, email: bob@example.com, age: 35, salary: 35000.12, birth: 1984-10-12, work: 40h, fulltime: true`

	extractor := NewRegexExtractor(regexDoc)

	assertExtractorSucc(t, extractor, "name: (?P<name>[a-zA-Z0-9 ]*)", String, "Bob")
	assertExtractorSucc(t, extractor, "email: (?P<email>[a-zA-Z0-9@.]*)", String, "bob@example.com")
	assertExtractorSucc(t, extractor, "age: (?P<age>[0-9]*)", Int, 35)
	assertExtractorSucc(t, extractor, "age: (?P<age>\\d+)", Int, 35)
	assertExtractorSucc(t, extractor, "salary: (?P<salary>[0-9.]*)", Float, 35000.12)
	assertExtractorSucc(t, extractor, "salary: (?P<salary>\\d+\\.\\d+)", Float, 35000.12)
	expected, _ := tryParseTime("1984-10-12")
	assertExtractorSucc(t, extractor, "birth: (?P<birth>[0-9-]*)", Time, expected)
	assertExtractorSucc(t, extractor, "birth: (?P<birth>\\d{4}-\\d{2}-\\d{2})", Time, expected)
	assertExtractorSucc(t, extractor, "work: (?P<work>\\d+[hms])", Duration, 40*time.Hour)
	assertExtractorSucc(t, extractor, "fulltime: (?P<fulltime>true|false)", Bool, true)
	// no Submatch
	assertExtractorSucc(t, extractor, "name: ", String, "name: ")
	// no match
	assertExtractorFail(t, extractor, "mismatch", String, "")
}
