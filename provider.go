package engine

type providers struct {
	values namedMap
}

func CreateProviders(p map[string]providerDef) providers {
	ret := providers{namedMap{}}
	for k, v := range p {
		v.name = k
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
