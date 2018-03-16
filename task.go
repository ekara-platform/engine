package engine

type tasks struct {
	values map[string]taskDef
}

func CreateTasks(p map[string]taskDef) tasks {
	ret := tasks{map[string]taskDef{}}
	for k, v := range p {
		ret.values[k] = v
	}
	return ret
}

func (l tasks) Contains(candidates ...string) bool {
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

func (l tasks) GetTask(candidate string) (TaskDescription, bool) {
	for k, _ := range l.values {
		if candidate == k {
			return l.values[candidate], true
		}
	}
	return nil, false
}
