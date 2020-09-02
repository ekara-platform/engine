package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDescribableEnvironment(t *testing.T) {
	o := Environment{
		QName: QualifiedName{
			Name:      "MyName",
			Qualifier: "MyQualifier",
		},
	}
	assert.Equal(t, "MyName_MyQualifier", o.DescName())
	assert.Equal(t, "Environment", o.DescType())
}

func TestDescribableNodeSet(t *testing.T) {
	o := NodeSet{Name: "Name1"}
	assert.Equal(t, "Name1", o.DescName())
	assert.Equal(t, "NodeSet", o.DescType())
}

func TestDescribableStack(t *testing.T) {
	o := Stack{Name: "Name1"}
	assert.Equal(t, "Name1", o.DescName())
	assert.Equal(t, "Stack", o.DescType())
}

func TestDescribableProvider(t *testing.T) {
	o := Provider{Name: "Name1"}
	assert.Equal(t, "Name1", o.DescName())
	assert.Equal(t, "Provider", o.DescType())
}

func TestChainOne(t *testing.T) {

	n1 := NodeSet{Name: "NameNode1"}
	chained := ChainDescribable(n1)

	assert.Equal(t, "NameNode1", chained.DescName())
	assert.Equal(t, "NodeSet", chained.DescType())
}

func TestChainMultiples(t *testing.T) {

	n1 := NodeSet{Name: "NameNode1"}
	n2 := NodeSet{Name: "NameNode2"}
	p1 := Provider{Name: "NameProvider1"}
	p2 := Provider{Name: "NameProvider2"}

	chained := ChainDescribable(n1, n2, p1, p2)

	assert.Equal(t, "NameNode1-NameNode2-NameProvider1-NameProvider2", chained.DescName())
	assert.Equal(t, "NodeSet-NodeSet-Provider-Provider", chained.DescType())
}

func ExampleChainDescribable() {
	p := Provider{Name: "MyProviderName"}
	n := NodeSet{Name: "MyNodesetName"}

	c := ChainDescribable(p, n)
	fmt.Printf("Chained types:%s, names:%s", c.DescType(), c.DescName())
	// Output: Chained types:Provider-NodeSet, names:MyProviderName-MyNodesetName

}
