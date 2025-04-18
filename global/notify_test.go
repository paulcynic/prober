
package global

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNotify(t *testing.T) {
	n := NotifySettings{
		Timeout: 0,
		Retry: Retry{
			Times:    0,
			Interval: 0,
		},
	}

	r := n.NormalizeTimeOut(0)
	assert.Equal(t, DefaultTimeOut, r)

	r = n.NormalizeTimeOut(10)
	assert.Equal(t, time.Duration(10), r)

	n.Timeout = 20
	r = n.NormalizeTimeOut(0)
	assert.Equal(t, time.Duration(20), r)

	retry := n.NormalizeRetry(Retry{Times: 10, Interval: 0})
	assert.Equal(t, Retry{Times: 10, Interval: DefaultRetryInterval}, retry)

	retry = n.NormalizeRetry(Retry{Times: 0, Interval: 10})
	assert.Equal(t, Retry{Times: DefaultRetryTimes, Interval: 10}, retry)

	retry = n.NormalizeRetry(Retry{Times: 10, Interval: 10})
	assert.Equal(t, Retry{Times: 10, Interval: 10}, retry)

	n.Retry.Times = 20
	retry = n.NormalizeRetry(Retry{Times: 0, Interval: 0})
	assert.Equal(t, Retry{Times: 20, Interval: DefaultRetryInterval}, retry)

	n.Retry.Interval = 20
	retry = n.NormalizeRetry(Retry{Times: 0, Interval: 0})
	assert.Equal(t, Retry{Times: 20, Interval: 20}, retry)
}
