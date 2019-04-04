package engine

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeFormat(t *testing.T) {
	assert.Equal(t, fmtDuration(1*time.Millisecond), "00h00m00s001")
	assert.Equal(t, fmtDuration(100*time.Millisecond), "00h00m00s100")
	assert.Equal(t, fmtDuration(1000*time.Millisecond), "00h00m01s000")
	assert.Equal(t, fmtDuration(60000*time.Millisecond), "00h01m00s000")
	assert.Equal(t, fmtDuration(3600000*time.Millisecond), "01h00m00s000")
	assert.Equal(t, fmtDuration(3600001*time.Millisecond), "01h00m00s001")
	assert.Equal(t, fmtDuration(3600100*time.Millisecond), "01h00m00s100")
	assert.Equal(t, fmtDuration(3601001*time.Millisecond), "01h00m01s001")
	assert.Equal(t, fmtDuration(3661001*time.Millisecond), "01h01m01s001")

}
