package roundrobin

import "sync/atomic"

// Statistics holds metrics for an Item, tracking how many times it has been served.
// Fields:
//   - ServesCount (int32): The number of times the item has been selected in the round-robin cycle.
type Statistics struct {
	ServesCount int32
}

// IncrementServesCount atomically increments the ServesCount by the specified value.
// This ensures thread-safe updates to the serve counter.
// Parameters:
//   - value (int32): The amount to increment ServesCount by, typically 1.
func (s *Statistics) IncrementServesCount(value int32) {
	atomic.AddInt32(&s.ServesCount, value)
}

// ResetServesCount atomically resets the ServesCount to zero.
// This is used to clear usage metrics when necessary.
func (s *Statistics) ResetServesCount() {
	atomic.StoreInt32(&s.ServesCount, 0)
}

type StatisticsInterface interface {
	IncrementServesCount(value int32)
	ResetServesCount()
}

var _ StatisticsInterface = (*Statistics)(nil)
