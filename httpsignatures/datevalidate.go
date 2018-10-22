package httpsignatures

import (
	"net/http"
	"time"
)

const maxTimeGapMillis = 30 * time.Second // 30 secs

// DateValidator checking validate by time range
type DateValidator struct {
	// TimeGap is max time different between client submit timestamp
	// and server time that considered valid. The time precision is millisecond.
	TimeGap time.Duration
}

// NewDateValidator return DateValidator with default value (30 second)
func NewDateValidator() *DateValidator {
	return &DateValidator{
		TimeGap: maxTimeGapMillis,
	}
}

// IsValid return nonce is valid or not by time range
func (v *DateValidator) IsValid(r *http.Request) bool {
	t, err := http.ParseTime(r.Header.Get("Date"))
	if err != nil {
		return false
	}
	serverTime := time.Now()
	start := serverTime.Add(-v.TimeGap)
	stop := serverTime.Add(v.TimeGap)
	if t.Before(start) || t.After(stop) {
		return false
	}
	return true
}
