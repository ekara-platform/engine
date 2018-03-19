package engine

type tasks struct {
	values namedMap
}

func CreateTasks(p map[string]taskDef) tasks {
	ret := tasks{namedMap{}}
	for k, v := range p {
		v.name = k
		ret.values[k] = v
	}
	return ret
}

func (l tasks) GetTask(candidate string) (TaskDescription, bool) {
	if v, ok := l.values[candidate]; ok {
		return v.(TaskDescription), ok
	}
	return nil, false
}
