package ansible

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const jsonInventory = `{
		"_meta": {
			"hostvars": {
				"host001": {
					"toto": "titi"
				}
			}
		},
		"all": {
			"children": ["group001"]
		},
		"group001": {
			"hosts": ["host001", "host002", "host003"],
			"vars": {
				"ekara": {
					"var1": true
				}
			},
			"children": ["group002"]
		},
		"group002": {
			"hosts": ["host003", "host004"],
			"vars": {
				"ekara": {
					"var2": 500
				}
			},
			"children":[]
		}
	}`

func TestParse(t *testing.T) {
	inv := Inventory{}
	err := inv.UnmarshalJSON([]byte(jsonInventory))
	assert.Nil(t, err)
	assert.Contains(t, inv.Hosts, "host001")
	assert.Contains(t, inv.Hosts, "host002")
	assert.Contains(t, inv.Hosts, "host003")
	assert.Contains(t, inv.Hosts, "host004")
}

func TestParseAndMarshal(t *testing.T) {
	inv := Inventory{}
	err := inv.UnmarshalJSON([]byte(jsonInventory))
	assert.Nil(t, err)
	b, err := json.Marshal(inv)
	assert.Nil(t, err)
	fmt.Println(string(b))
}
