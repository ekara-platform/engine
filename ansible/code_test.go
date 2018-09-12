package ansible

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError_Fetch(t *testing.T) {
	e := ReturnedError(ansible_ok)
	assert.Nil(t, e)

	e = ReturnedError(ansible_error)
	assert.NotNil(t, e)

	e = ReturnedError(ansible_failed)
	assert.NotNil(t, e)

	e = ReturnedError(ansible_unreachable)
	assert.NotNil(t, e)

	e = ReturnedError(ansible_parser)
	assert.NotNil(t, e)

	e = ReturnedError(ansible_bar_or_incomplete)
	assert.NotNil(t, e)

	e = ReturnedError(ansible_user_insterrupted)
	assert.NotNil(t, e)

	e = ReturnedError(ansible_unexpected_error)
	assert.NotNil(t, e)

	e = ReturnedError(999999)
	assert.Nil(t, e)
}
