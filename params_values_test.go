package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseParamValues(t *testing.T) {
	path := "./testdata/params_value.yaml"
	m, err := ParseParamValues(path)
	assert.NotNil(t, m)
	assert.Nil(t, err)
	assert.Equal(t, 15, len(m))

	val, ok := m["aws.params_pack1.params_pack1_key1"]
	assert.True(t, ok)
	assert.Equal(t, "params_pack1_value1", val)

	val, ok = m["aws.params_pack1.params_pack1_key2"]
	assert.True(t, ok)
	assert.Equal(t, "params_pack1_value2", val)

	val, ok = m["aws.params_pack1.params_pack1_sub_pack1.params_pack1_sub_pack1_key1"]
	assert.True(t, ok)
	assert.Equal(t, "params_pack1_sub_pack1_value1", val)

	val, ok = m["aws.params_pack1.params_pack1_sub_pack1.params_pack1_sub_pack1_key2"]
	assert.True(t, ok)
	assert.Equal(t, "params_pack1_sub_pack1_value2", val)

	val, ok = m["aws.params_pack2.params_pack2_key1"]
	assert.True(t, ok)
	assert.Equal(t, "params_pack2_value1", val)

	val, ok = m["aws.params_pack2.params_pack2_key2"]
	assert.True(t, ok)
	assert.Equal(t, "params_pack2_value2", val)

	val, ok = m["aws.params_pack2.params_pack2_sub_pack1.params_pack2_sub_pack1_key1"]
	assert.True(t, ok)
	assert.Equal(t, "params_pack2_sub_pack1_value1", val)

	val, ok = m["aws.params_pack2.params_pack2_sub_pack1.params_pack2_sub_pack1_key2"]
	assert.True(t, ok)
	assert.Equal(t, "params_pack2_sub_pack1_value2", val)

	val, ok = m["aws.params_pack2.params_pack2_sub_pack2.params_pack2_sub_pack2_key1"]
	assert.True(t, ok)
	assert.Equal(t, "params_pack2_sub_pack2_value1", val)

	val, ok = m["aws.params_pack2.params_pack2_sub_pack2.params_pack2_sub_pack2_key2"]
	assert.True(t, ok)
	assert.Equal(t, "params_pack2_sub_pack2_value2", val)

	val, ok = m["root_key1"]
	assert.True(t, ok)
	assert.Equal(t, "root_val1", val)

	val, ok = m["root_key2"]
	assert.True(t, ok)
	assert.Equal(t, "root_val2", val)

	val, ok = m["root_key3"]
	assert.True(t, ok)
	assert.Equal(t, "root_val3", val)

	_, ok = m["root_key_nil_value"]
	assert.False(t, ok)

	_, ok = m["aws.params_pack1.params_pack1_key_nil_value;"]
	assert.False(t, ok)

	val, ok = m["111"]
	assert.True(t, ok)
	assert.Equal(t, "root_111_val", val)

	val, ok = m["222"]
	assert.True(t, ok)
	assert.Equal(t, "333", val)
}
