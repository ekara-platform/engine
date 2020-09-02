package model

//Labels represents used defined labels which will be placed on the created environment
//machines and also on the nodes for Docker
type Labels map[string]string

func (r Labels) override(parent Labels) Labels {
	dst := make(map[string]string)
	for k, v := range parent {
		dst[k] = v
	}
	for k, v := range r {
		dst[k] = v
	}
	return dst
}
