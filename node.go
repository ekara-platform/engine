package engine

type nodes struct {
	values map[string]nodeSetDef
}

func CreateNodes(p map[string]nodeSetDef) nodes {
	ret := nodes{map[string]nodeSetDef{}}
	for k, v := range p {
		ret.values[k] = v
	}
	return ret
}

func (l nodes) Contains(candidates ...string) bool {
	for _, l1 := range candidates {
		contains := false
		for k, _ := range l.values {
			if l1 == k {
				contains = true
				break
			}
		}
		if !contains {
			return false
		}
	}
	return true
}

func (l nodes) GetNode(candidate string) (NodeDescription, bool) {
	for k, _ := range l.values {
		if candidate == k {
			return l.values[candidate], true
		}
	}
	return nil, false
}
