package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegexp(t *testing.T) {
	assert.True(t, isValidQualifiedName("aa_bb"))
	assert.True(t, isValidQualifiedName("AA_BB"))
	assert.True(t, isValidQualifiedName("11_22"))
	assert.True(t, isValidQualifiedName("aaAA11_aaBB11"))

	assert.True(t, isValidQualifiedName("aa_BB"))
	assert.True(t, isValidQualifiedName("aa_11"))
	assert.True(t, isValidQualifiedName("11_bb"))
	assert.True(t, isValidQualifiedName("11_BB"))

	assert.True(t, isValidQualifiedName("aaAA11_aaBB11"))

	assert.False(t, isValidQualifiedName("a-b"))
}

// Tests on Envionment with name and qualifier
func TestFullQualifiedName(t *testing.T) {
	env := Environment{
		QName: QualifiedName{
			Name:      "ABC",
			Qualifier: "DEF",
		},
	}
	assert.Equal(t, "ABC_DEF", env.QName.String())
}

// Tests on Envionment with qualifier only
func TestPartialQualifiedName(t *testing.T) {
	env := Environment{
		QName: QualifiedName{Name: "ABC"},
	}
	assert.Equal(t, "ABC", env.QName.String())
}

// Tests on Environment without name and qualifier
func TestEmptyQualifiedName(t *testing.T) {
	env := Environment{}
	// check the zero value
	assert.Equal(t, QualifiedName{}, env.QName)
	assert.Equal(t, "", env.QName.String())
}

// Tests on Yaml Envionment with name and qualifier
func TestFullQualifiedNameYaml(t *testing.T) {
	env := yamlEnvironment{
		Name:      "ABC",
		Qualifier: "DEF",
	}
	assert.Equal(t, "ABC_DEF", env.yamlQualifiedName().String())
}

// Tests on Yaml Envionment with qualifier only
func TestPartialQualifiedNameYaml(t *testing.T) {
	env := yamlEnvironment{
		Name: "ABC",
	}
	assert.Equal(t, "ABC", env.yamlQualifiedName().String())
}

// Tests on Yaml Envionment without name and qualifier
func TestEmptyQualifiedNameYaml(t *testing.T) {
	env := yamlEnvironment{}
	// check the zero value
	assert.Equal(t, QualifiedName{}, env.yamlQualifiedName())
	assert.Equal(t, "", env.yamlQualifiedName().String())
}

func TestValidQualifiedName(t *testing.T) {
	env := Environment{
		QName: QualifiedName{
			Name:      "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			Qualifier: "abcdefghijklmnopqrstuvwxyz",
		},
		loc: DescriptorLocation{},
	}
	assert.True(t, !env.QName.validate(env, env.loc).HasErrors())

	env.QName = QualifiedName{Name: "0123456789", Qualifier: "0123456789"}
	assert.True(t, !env.QName.validate(env, env.loc).HasErrors())

	env.QName = QualifiedName{Name: "à"}
	assert.False(t, !env.QName.validate(env, env.loc).HasErrors())

	env.QName = QualifiedName{Name: "é"}
	assert.False(t, !env.QName.validate(env, env.loc).HasErrors())

	env.QName = QualifiedName{Name: "ù"}
	assert.False(t, !env.QName.validate(env, env.loc).HasErrors())

	env.QName = QualifiedName{Name: "è"}
	assert.False(t, !env.QName.validate(env, env.loc).HasErrors())

	env.QName = QualifiedName{Name: "ç"}
	assert.False(t, !env.QName.validate(env, env.loc).HasErrors())

	env.QName = QualifiedName{Name: "!"}
	assert.False(t, !env.QName.validate(env, env.loc).HasErrors())

	env.QName = QualifiedName{Name: "-"}
	assert.False(t, !env.QName.validate(env, env.loc).HasErrors())

	env.QName = QualifiedName{Name: "&"}
	assert.False(t, !env.QName.validate(env, env.loc).HasErrors())

	env.QName = QualifiedName{Name: "#"}
	assert.False(t, !env.QName.validate(env, env.loc).HasErrors())
}

// QualifiedName returns the concatenation of the environment name and qualifier
// separated by a "_".
// If the environment qualifier is not defined it will return just the name
func (r yamlEnvironment) yamlQualifiedName() QualifiedName {
	if len(r.Qualifier) == 0 {
		return QualifiedName{Name: r.Name}
	}
	return QualifiedName{Name: r.Name + "_" + r.Qualifier}
}
