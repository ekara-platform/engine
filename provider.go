package engine

type providers struct {
	values namedMap
}

func CreateProviders(desc *environmentDef) providers {
	ret := providers{namedMap{}}
	for k, v := range desc.Providers {
		v.name = k
		v.desc = desc
		ret.values[k] = v
	}
	return ret
}

func (l providers) GetProvider(candidate string) (ProviderDescription, bool) {
	if v, ok := l.values[candidate]; ok {
		return v.(ProviderDescription), ok
	}
	return nil, false
}
