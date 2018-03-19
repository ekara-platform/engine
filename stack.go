package engine

type stacks struct {
	values namedMap
}

func CreateStacks(p map[string]stackDef) stacks {
	ret := stacks{namedMap{}}
	for k, v := range p {
		v.name = k
		ret.values[k] = v
	}
	return ret
}

func (l stacks) GetStack(candidate string) (StackDescription, bool) {
	if v, ok := l.values[candidate]; ok {
		return v.(StackDescription), ok
	}
	return nil, false
}
