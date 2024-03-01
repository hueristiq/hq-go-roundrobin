package hqgoroundrobin

import (
	"errors"
	"sync"
	"sync/atomic"
)

// Item represents a single unit within the round-robin collection. It holds a value and associated statistics
// to track how many times it has been served.
type Item struct {
	// value is the content or identifier of the item.
	value string
	// Statistics holds metrics related to the item, such as its serve count.
	Statistics Statistics
}

// Value returns the underlying value of the item. This method allows accessing the item's content.
func (i Item) Value() (value string) {
	return i.value
}

// ItemInterface defines the interface that an Item must implement. This ensures that all items
// can return their underlying value.
type ItemInterface interface {
	// Value method returns the value of the item.
	Value() (value string)
}

// Statistics holds metrics related to an item, particularly how many times it has been served.
// This allows for tracking and potentially balancing the distribution of items.
type Statistics struct {
	// ServesCount is a counter for the number of times an item has been served.
	ServesCount int32
}

// IncrementServesCount atomically increases the ServesCount by a given value. This method is used
// to update the serve count in a concurrent-safe manner.
func (s *Statistics) IncrementServesCount(value int32) {
	atomic.AddInt32(&s.ServesCount, value)
}

// ResetServesCount atomically resets the ServesCount to zero. This can be used to restart
// the serve count statistics for an item
func (s *Statistics) ResetServesCount() {
	atomic.StoreInt32(&s.ServesCount, 0)
}

// StatisticsInterface defines the interface for manipulating item statistics. This abstraction
// allows for flexibility in how statistics are implemented and modified.
type StatisticsInterface interface {
	// IncrementServesCount method increases the serve count by a specified value.
	IncrementServesCount(value int32)
	// ResetServesCount method resets the serve count to zero.
	ResetServesCount()
}

// RoundRobin manages a collection of items, allowing for thread-safe addition and retrieval in a round-robin fashion.
// It supports concurrent access and ensures that items are served in a balanced order.
type RoundRobin struct {
	// items is a slice of the managed items in the round-robin.
	items []Item
	// itemsMap is used in conjunction with the slice to ensure uniqueness of items.
	itemsMap sync.Map
	// nextItemIndex is the index of the next item to serve, managed atomically to support concurrent access.
	nextItemIndex uint32
	// currentItemServesCount tracks the serve count of the currently serving item, allowing for rotation based on serve count.
	currentItemServesCount uint32
	// mutex ensures thread-safe access to the round-robin, particularly for operations that modify its state.
	mutex sync.Mutex
	// Options hold configuration settings for the round-robin, like rotation behavior.
	Options Options
}

// Items returns a copy of the items slice, allowing external access to the current state of the round-robin
// without compromising thread safety.
func (r *RoundRobin) Items() (items []Item) {
	return r.items
}

// Add inserts one or more new values into the round-robin collection. It ensures that each item is unique
// and updates the collection in a thread-safe manner.
func (r *RoundRobin) Add(values ...string) {
	for _, value := range values {
		item := Item{
			value: value,
		}

		// Attempt to store the item in the map. If it's a new item, also append it to the slice.
		if _, loaded := r.itemsMap.LoadOrStore(value, struct{}{}); !loaded {
			r.items = append(r.items, item)
		}
	}
}

// Next retrieves the next item in the round-robin order. It manages the serve count and rotates to the next item
// as necessary, ensuring thread-safe access and modification of the round-robin state.
func (r *RoundRobin) Next() (item Item) {
	r.mutex.Lock()

	defer r.mutex.Unlock()

	currentAmount := atomic.LoadUint32(&r.currentItemServesCount)

	// Rotate to the next item if the current item has reached its serve limit.
	if currentAmount >= uint32(r.Options.RotateAmount) {
		atomic.StoreUint32(&r.currentItemServesCount, 1)
		atomic.AddUint32(&r.nextItemIndex, 1)
	} else {
		atomic.AddUint32(&r.currentItemServesCount, 1)
	}

	nextItemIndex := (int(r.nextItemIndex) - 1) % len(r.items)

	// Safeguard against index out-of-bounds, defaulting to the first item if necessary.
	if nextItemIndex < 0 || nextItemIndex > len(r.items) {
		r.items[0].Statistics.IncrementServesCount(1) // Increment stats by 1 everytime item is retrieved

		return r.items[0]
	}

	r.items[nextItemIndex].Statistics.IncrementServesCount(1)

	return r.items[nextItemIndex]
}

// RoundRobinInterface defines the interface for a round-robin mechanism, abstracting the functionality
// to add items and retrieve the next item in sequence. This facilitates testing and alternative implementations.
type RoundRobinInterface interface {
	// Items method retrieves a copy of the items  in the round-robin sequence.
	Items() (items []Item)
	// Add method allows adding one or more items to the round-robin.
	Add(values ...string)
	// Next method retrieves the next item in the round-robin sequence.
	Next() (item Item)
}

// Options holds configuration settings for the round-robin, such as rotation amount.
// This allows customization of the round-robin behavior.
type Options struct {
	// RotateAmount specifies the number of serves before rotating to the next item.
	RotateAmount int32
}

var (
	// ErrNoItems indicates that no items are available for operation, typically used when initializing
	// a new RoundRobin instance without any items.
	ErrNoItems = errors.New("no items")

	// Interface assertions verify at compile time that the types implement the specified interfaces.
	_ ItemInterface       = (*Item)(nil)
	_ StatisticsInterface = (*Statistics)(nil)
	_ RoundRobinInterface = (*RoundRobin)(nil)

	// DefaultOptions provides a set of default configuration options for new round-robin instances,
	// simplifying the initialization process.
	DefaultOptions = Options{
		RotateAmount: 1,
	}
)

// New creates a new RoundRobin instance with default options. It initializes the round-robin with a set of initial items,
// returning an error if no items are provided.
func New(items ...string) (rr *RoundRobin, err error) {
	return NewWithOptions(DefaultOptions, items...)
}

// NewWithOptions creates a new RoundRobin instance with custom options. It allows for greater flexibility
// in configuring the round-robin behavior and initializes the instance with a set of initial items.
func NewWithOptions(options Options, items ...string) (rr *RoundRobin, err error) {
	if len(items) == 0 {
		err = ErrNoItems

		return
	}

	rr = &RoundRobin{
		Options: options,
	}

	rr.Add(items...)
	// Ensure the next item index starts from the first item.
	rr.nextItemIndex = 1

	return
}
