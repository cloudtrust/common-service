package healthcheck

import "time"

// TimeProvider is the interface we use to provide real time to ease testing
type TimeProvider interface {
	Now() time.Time
}

type RealTimeProvider struct{}

func (tp RealTimeProvider) Now() time.Time {
	return time.Now()
}
