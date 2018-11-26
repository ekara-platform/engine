package ansible

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError_Fetch(t *testing.T) {
	e := ReturnedError(ansibleOk)
	assert.Nil(t, e)

	e = ReturnedError(ansibleError)
	assert.NotNil(t, e)

	e = ReturnedError(ansibleFailed)
	assert.NotNil(t, e)

	e = ReturnedError(ansibleUnreachable)
	assert.NotNil(t, e)

	e = ReturnedError(ansibleParser)
	assert.NotNil(t, e)

	e = ReturnedError(ansibleBarOrIncomplete)
	assert.NotNil(t, e)

	e = ReturnedError(ansibleUserInsterrupted)
	assert.NotNil(t, e)

	e = ReturnedError(ansibleUnexpectedError)
	assert.NotNil(t, e)

	e = ReturnedError(999999)
	assert.Nil(t, e)
}
