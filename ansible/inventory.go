package ansible

import (
	"encoding/json"
	"fmt"
)

type (
	Inventory struct {
		Hosts  map[string]Host
		Groups map[string]Group
	}

	Host struct {
		Name string
		Vars InventoryVars
	}

	Group struct {
		Children map[string]Group
		Hosts    []string
		Vars     InventoryVars
	}

	InventoryVars map[string]interface{}
)

func (i *Inventory) UnmarshalJSON(data []byte) error {
	const allGroup = "all"
	const meta = "_meta"
	const hostVars = "hostvars"

	raw := make(map[string]interface{})
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	i.Hosts = make(map[string]Host)
	i.Groups = make(map[string]Group)

	for k, v := range raw {
		if k != meta && k != allGroup {
			hostSet := make(map[string]bool)
			groups := readGroup(raw, k, v.(map[string]interface{}), hostSet)
			for hostName := range hostSet {
				host := Host{Name: hostName, Vars: InventoryVars{}}
				if _meta, ok := raw[meta].(map[string]interface{}); ok {
					if allHostVars, ok := _meta[hostVars].(map[string]interface{}); ok {
						if hostVars, ok := allHostVars[hostName].(map[string]interface{}); ok {
							host.Vars = buildInventoryVariables(hostVars)
						}
					}
				}
				i.Hosts[hostName] = host
			}
			for name, grp := range groups {
				i.Groups[name] = grp
			}
		}
	}

	return nil
}

func readGroup(root map[string]interface{}, name string, current map[string]interface{}, hostSet map[string]bool) map[string]Group {
	res := make(map[string]Group)
	group := Group{}
	hasChildren := false
	if grpChildren, ok := current["children"].([]interface{}); ok {
		for _, grpName := range grpChildren {
			name := fmt.Sprintf("%s", grpName)
			if grp, ok := root[name].(map[string]interface{}); ok {
				group.Children = readGroup(root, name, grp, hostSet)
				hasChildren = true
			}
		}
	}
	if !hasChildren {
		group.Children = make(map[string]Group)
	}

	if grpVars, ok := current["vars"].(map[string]interface{}); ok {
		group.Vars = buildInventoryVariables(grpVars)
	} else {
		group.Vars = InventoryVars{}
	}

	if grpHosts, ok := current["hosts"].([]interface{}); ok {
		var hosts []string
		for _, rHost := range grpHosts {
			host := fmt.Sprintf("%s", rHost)
			hosts = append(hosts, host)
			hostSet[host] = true
		}
		group.Hosts = hosts
	} else {
		group.Hosts = []string{}
	}
	res[fmt.Sprintf("%s", name)] = group
	return res
}

func buildInventoryVariables(src map[string]interface{}) InventoryVars {
	vars := InventoryVars{}
	if ekara, ok := src["ekara"].(map[string]interface{}); ok {
		for k, v := range ekara {
			vars[fmt.Sprintf("%s", k)] = v
		}
	}
	return vars
}
