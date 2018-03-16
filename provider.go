package engine

type providers struct {
	values map[string]providerDef
}

func CreateProviders(p map[string]providerDef) providers {
	ret := providers{map[string]providerDef{}}
	for k, v := range p {
		ret.values[k] = v
	}
	return ret
}

func (l providers) Contains(candidates ...string) bool {
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

func (l providers) GetProvider(candidate string) (ProviderDescription, bool) {
	for k, _ := range l.values {
		if candidate == k {
			return l.values[candidate], true
		}
	}
	return nil, false
}
