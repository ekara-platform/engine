package engine

type params struct {
	values map[string]string
}

func CreateParameters(p map[string]string) Parameters {
	ret := params{map[string]string{}}
	for k, v := range p {
		ret.values[k] = v
	}
	return ret
}

func (p params) AsMap() map[string]string {
	ret := map[string]string{}
	for k, v := range p.values {
		ret[k] = v
	}
	return ret
}
