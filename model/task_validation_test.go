package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test loading an task without any playbook .
//
// The validation must complain only about the missing playbook
//
//- Error: empty playbook path @tasks.task1.playbook
//
func TestValidationTasksNoPlayBook(t *testing.T) {
	vErrs, _ := testEmptyContent(t, "task_playbook", false)
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.True(t, vErrs.contains(Error, "no playbook specified", "tasks.task1.playbook"))
}
