package engine

import (
	"github.com/ekara-platform/engine/ansible"
	"github.com/ekara-platform/engine/util"
)

type (
	runtimeContext struct {
		ekaraError error
		report     ReportFileContent
		buffer     map[string]ansible.Buffer
	}
)

func (c *runtimeContext) getBuffer(p *util.FolderPath) ansible.Buffer {
	// We check if we have a buffer corresponding to the provided folder path
	if val, ok := c.buffer[p.Path()]; ok {
		return val
	}
	return ansible.CreateBuffer()
}
