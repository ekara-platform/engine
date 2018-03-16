package engine

type stacks struct {
	values map[string]stackDef
}

func CreateStacks(p map[string]stackDef) stacks {
	ret := stacks{map[string]stackDef{}}
	for k, v := range p {
		ret.values[k] = v
	}
	return ret
}

func (l stacks) Contains(candidates ...string) bool {
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

func (l stacks) GetStack(candidate string) (StackDescription, bool) {
	for k, _ := range l.values {
		if candidate == k {
			return l.values[candidate], true
		}
	}
	return nil, false
}
