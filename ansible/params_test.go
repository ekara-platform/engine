package ansible

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {

	bp := BuildBaseParam("client_val", "uid_val", "provider_val", "pubK_val", "privK_val")
	body := bp.Body

	val, ok := body["connectionConfig"]
	assert.True(t, ok)

	v, okType := val.(map[string]interface{})
	assert.True(t, okType)
	for k, v := range v {
		switch k {
		case "provider":
			assert.Equal(t, "provider_val", v)
		case "machine_public_key":
			assert.Equal(t, "pubK_val", v)
		case "machine_private_key":
			assert.Equal(t, "privK_val", v)
		default:
			assert.Fail(t, "unknown key")
		}
	}

	val, ok = body["environment"]
	assert.True(t, ok)

	v, okType = val.(map[string]interface{})
	assert.True(t, okType)
	for k, v := range v {
		switch k {
		case "name":
			assert.Equal(t, "client_val", v)
		case "uid":
			assert.Equal(t, "uid_val", v)
		default:
			assert.Fail(t, "unknown key")
		}
	}
}

func TestAddString(t *testing.T) {

	bp := BuildBaseParam("client_val", "uid_val", "provider_val", "pubK_val", "privK_val")

	bp.AddString("string_key1", "string_val1")
	bp.AddString("string_key2", "string_val2")

	body := bp.Body

	val, ok := body["string_key1"]
	assert.True(t, ok)

	v, okType := val.(string)
	assert.True(t, okType)
	vString := string(v)
	assert.Equal(t, "string_val1", vString)

	val, ok = body["string_key2"]
	assert.True(t, ok)

	v, okType = val.(string)
	assert.True(t, okType)
	vString = string(v)
	assert.Equal(t, "string_val2", vString)
}

func TestAddInt(t *testing.T) {

	bp := BuildBaseParam("client_val", "uid_val", "provider_val", "pubK_val", "privK_val")

	bp.AddInt("string_key1", 11)
	bp.AddInt("string_key2", 22)

	body := bp.Body

	val, ok := body["string_key1"]
	assert.True(t, ok)

	v, okType := val.(int)
	assert.True(t, okType)
	vInt := int(v)
	assert.Equal(t, vInt, 11)

	val, ok = body["string_key2"]
	assert.True(t, ok)

	v, okType = val.(int)
	assert.True(t, okType)
	vInt = int(v)
	assert.Equal(t, vInt, 22)
}

func TestAddMapString(t *testing.T) {

	bp := BuildBaseParam("client_val", "uid_val", "provider_val", "pubK_val", "privK_val")

	m := make(map[string]interface{})
	m["string_key1"] = "string_val1"
	m["string_key2"] = "string_val2"
	bp.AddMap(m)

	body := bp.Body

	val, ok := body["string_key1"]
	assert.True(t, ok)

	v, okType := val.(string)
	assert.True(t, okType)
	vString := string(v)
	assert.Equal(t, "string_val1", vString)

	val, ok = body["string_key2"]
	assert.True(t, ok)

	v, okType = val.(string)
	assert.True(t, okType)
	vString = string(v)
	assert.Equal(t, "string_val2", vString)
}

func TestAddMapInt(t *testing.T) {

	bp := BuildBaseParam("client_val", "uid_val", "provider_val", "pubK_val", "privK_val")

	m := make(map[string]interface{})
	m["string_key1"] = 11
	m["string_key2"] = 22
	bp.AddMap(m)

	body := bp.Body

	val, ok := body["string_key1"]
	assert.True(t, ok)

	v, okType := val.(int)
	assert.True(t, okType)
	vInt := int(v)
	assert.Equal(t, vInt, 11)

	val, ok = body["string_key2"]
	assert.True(t, ok)

	v, okType = val.(int)
	assert.True(t, okType)
	vInt = int(v)
	assert.Equal(t, vInt, 22)
}

func TestAddInterface(t *testing.T) {

	bp := BuildBaseParam("client_val", "uid_val", "provider_val", "pubK_val", "privK_val")

	m := make(map[string]interface{})
	m["string_key1"] = "String"
	m["string_key2"] = 22
	m["string_key3"] = true
	bp.AddMap(m)

	body := bp.Body

	// string
	val, ok := body["string_key1"]
	assert.True(t, ok)

	vs, okType := val.(string)
	assert.True(t, okType)
	vString := string(vs)
	assert.Equal(t, vString, "String")

	// int
	val, ok = body["string_key2"]
	assert.True(t, ok)

	vi, okType := val.(int)
	assert.True(t, okType)
	vInt := int(vi)
	assert.Equal(t, vInt, 22)

	// boolean
	val, ok = body["string_key3"]
	assert.True(t, ok)

	vb, okType := val.(bool)
	assert.True(t, okType)
	vBool := bool(vb)
	assert.Equal(t, vBool, true)
}

func TestAddNamedMapString(t *testing.T) {

	bp := BuildBaseParam("client_val", "uid_val", "provider_val", "pubK_val", "privK_val")

	m := make(map[string]interface{})
	m["string_key1"] = "string_val1"
	m["string_key2"] = "string_val2"
	bp.AddNamedMap("master_key", m)

	body := bp.Body

	val, ok := body["master_key"]
	assert.True(t, ok)

	v, okType := val.(map[string]interface{})
	assert.True(t, okType)
	for k, v := range v {
		switch k {
		case "string_key1":
			assert.Equal(t, "string_val1", v)
		case "string_key2":
			assert.Equal(t, "string_val2", v)
		default:
			assert.Fail(t, "unknown key")
		}
	}
}

func TestAddBuffer(t *testing.T) {

	bp := BuildBaseParam("client_val", "uid_val", "provider_val", "pubK_val", "privK_val")

	buf := CreateBuffer()
	buf.Param["string_key1"] = "string_val1"
	buf.Param["string_key2"] = "string_val2"

	bp.AddBuffer(buf)

	body := bp.Body

	val, ok := body["string_key1"]
	assert.True(t, ok)

	v, okType := val.(string)
	assert.True(t, okType)
	vString := string(v)
	assert.Equal(t, "string_val1", vString)

	val, ok = body["string_key2"]
	assert.True(t, ok)

	v, okType = val.(string)
	assert.True(t, okType)
	vString = string(v)
	assert.Equal(t, "string_val2", vString)
}
