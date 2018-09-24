package ansible

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildEquals(t *testing.T) {
	m := make(map[string]string)
	m["key1"] = "val1"
	m["key2"] = "val2"
	m["key3"] = "val3"

	r := buildEquals(m)

	assert.Equal(t, true, strings.Contains(r, "key1=val1"))
	assert.Equal(t, true, strings.Contains(r, "key2=val2"))
	assert.Equal(t, true, strings.Contains(r, "key3=val3"))
}


// BuildEquals converts a map[string]string into a succession
// of equalities of type "map key=map value" separated by a space
//
//	Example of returned value :
//		"key1=val1 key2=val2 key3=val3"
func buildEquals(m map[string]string) string {
	var r string
	for k, v := range m {
		r = r + k + "=" + v + " "
	}
	return r
}
