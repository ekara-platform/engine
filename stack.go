package engine

type stacks struct {
	values namedMap
}

func CreateStacks(desc *environmentDef) stacks {
	ret := stacks{namedMap{}}
	for k, v := range desc.Stacks {
		v.name = k
		v.desc = desc
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
