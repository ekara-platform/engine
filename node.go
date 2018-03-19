package engine

type nodes struct {
	values namedMap
}

func CreateNodes(desc *environmentDef) nodes {
	ret := nodes{namedMap{}}
	for k, v := range desc.Nodes {
		v.name = k
		v.desc = desc
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
