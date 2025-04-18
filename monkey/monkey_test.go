
package monkey

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type myStruct struct{}

func (s *myStruct) Method() string {
	return "original"
}

func TestPatch(t *testing.T) {
	assert.Equal(t, "original", strings.Clone("original"))

	Patch(strings.Clone, func(_ string) string { return "replacement" })
	assert.Equal(t, "replacement", strings.Clone("original"))

	Patch(strings.Clone, func(_ string) string { return "replacement2" })
	assert.Equal(t, "replacement2", strings.Clone("original"))

	Unpatch(strings.Clone)
	assert.Equal(t, "original", strings.Clone("original"))

	result := Unpatch(strings.Clone)
	assert.Equal(t, false, result)
}

func TestPatchInstanceMethod(t *testing.T) {
	assert.Equal(t, "original", (&myStruct{}).Method())

	PatchInstanceMethod(reflect.TypeOf(&myStruct{}), "Method", func(*myStruct) string { return "replacement" })
	assert.Equal(t, "replacement", (&myStruct{}).Method())

	PatchInstanceMethod(reflect.TypeOf(&myStruct{}), "Method", func(*myStruct) string { return "replacement2" })
	assert.Equal(t, "replacement2", (&myStruct{}).Method())

	UnpatchInstanceMethod(reflect.TypeOf(&myStruct{}), "Method")
	assert.Equal(t, "original", (&myStruct{}).Method())

	result := UnpatchInstanceMethod(reflect.TypeOf(&myStruct{}), "Method")
	assert.Equal(t, false, result)
}

func TestUnpatchAll(t *testing.T) {
	assert.Equal(t, "original", strings.Clone("original"))
	Patch(strings.Clone, func() string { return "replacement" })
	assert.Equal(t, "replacement", strings.Clone("original"))

	assert.Equal(t, "original", (&myStruct{}).Method())
	PatchInstanceMethod(reflect.TypeOf(&myStruct{}), "Method", func(*myStruct) string { return "replacement" })
	assert.Equal(t, "replacement", (&myStruct{}).Method())

	UnpatchAll()
	assert.Equal(t, "original", strings.Clone("original"))
	assert.Equal(t, "original", (&myStruct{}).Method())
}
