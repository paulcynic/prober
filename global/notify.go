
package global

import "time"

// NotifySettings is the global notification setting
type NotifySettings struct {
	TimeFormat string
	Timeout    time.Duration
	Retry      Retry
}

// NormalizeTimeOut return a normalized timeout value
func (n *NotifySettings) NormalizeTimeOut(t time.Duration) time.Duration {
	return normalize(n.Timeout, t, 0, DefaultTimeOut)
}

// NormalizeRetry return a normalized retry value
func (n *NotifySettings) NormalizeRetry(retry Retry) Retry {
	retry.Interval = normalize(n.Retry.Interval, retry.Interval, 0, DefaultRetryInterval)
	retry.Times = normalize(n.Retry.Times, retry.Times, 0, DefaultRetryTimes)
	return retry
}
