package engine

// Interfaces

type Labeled interface {
	GetLabels() Labels
}

type Labels interface {
	Contains(...string) bool
	AsStrings() []string
}

// Implementation

type labels struct {
	values []string
}

func CreateLabels(values ...string) Labels {
	ret := labels{make([]string, len(values))}
	copy(ret.values, values)
	return ret
}

func (l labels) Contains(candidates ...string) bool {
	for _, l1 := range candidates {
		contains := false
		for _, l2 := range l.values {
			if l1 == l2 {
				contains = true
			}
		}
		if !contains {
			return false
		}
	}
	return true
}

func (l labels) AsStrings() []string {
	ret := make([]string, len(l.values))
	copy(ret, l.values)
	return ret
}
