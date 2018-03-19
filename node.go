package engine

type nodes struct {
	values namedMap
}

func CreateNodes(p map[string]nodeSetDef) nodes {
	ret := nodes{namedMap{}}
	for k, v := range p {
		v.name = k
		ret.values[k] = v
	}
	return ret
}

func (l nodes) GetNode(candidate string) (NodeDescription, bool) {
	if v, ok := l.values[candidate]; ok {
		return v.(NodeDescription), ok
	}
	return nil, false
}
