package action

type graph struct {
	nodes   []string
	outputs map[string]map[string]int
	inputs  map[string]int
}

func newGraph(cap int) *graph {
	return &graph{
		nodes:   make([]string, 0, cap),
		inputs:  make(map[string]int),
		outputs: make(map[string]map[string]int),
	}
}

func (g *graph) addNode(name string) bool {
	g.nodes = append(g.nodes, name)

	if _, ok := g.outputs[name]; ok {
		return false
	}
	g.outputs[name] = make(map[string]int)
	g.inputs[name] = 0
	return true
}

func (g *graph) addEdge(dependency, to string) bool {
	m, ok := g.outputs[dependency]
	if !ok {
		return false
	}

	m[to] = len(m) + 1
	g.inputs[to]++

	return true
}

func (g *graph) sort() ([]string, bool) {
	L := make([]string, 0, len(g.nodes))
	S := make([]string, 0, len(g.nodes))

	for _, n := range g.nodes {
		if g.inputs[n] == 0 {
			S = append(S, n)
		}
	}

	for len(S) > 0 {
		var n string
		n, S = S[0], S[1:]
		L = append(L, n)

		ms := make([]string, len(g.outputs[n]))
		for m, i := range g.outputs[n] {
			ms[i-1] = m
		}

		for _, m := range ms {
			delete(g.outputs[n], m)
			g.inputs[m]--

			if g.inputs[m] == 0 {
				S = append(S, m)
			}
		}
	}

	N := 0
	for _, v := range g.inputs {
		N += v
	}

	if N > 0 {
		return L, false
	}

	return L, true
}
