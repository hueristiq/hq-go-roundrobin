package roundrobin

import (
	"sync"
	"sync/atomic"

	hqgoerrors "github.com/hueristiq/hq-go-errors"
)

// RoundRobin manages a thread-safe collection of items for round-robin selection.
// It supports adding items, retrieving the next item based on rotation rules, and tracking statistics.
// Fields:
//   - items ([]Item): The slice of Items in the round-robin cycle.
//   - itemsMap (sync.Map): A thread-safe map ensuring uniqueness of item values.
//   - nextItemIndex (uint32): The index of the next item to be served, updated atomically.
//   - currentItemServesCount (uint32): The number of times the current item has been served, updated atomically.
//   - mutex (sync.Mutex): Synchronizes access to non-atomic operations, such as modifying the items slice.
//   - Options (Options): Configuration for the round-robin behavior, such as rotation frequency.
type RoundRobin struct {
	items                  []Item
	itemsMap               sync.Map
	nextItemIndex          uint32
	currentItemServesCount uint32
	mutex                  sync.Mutex
	Options                Options
}

// Items returns the current list of items in the round-robin cycle.
// Returns:
//   - items ([]Item): A slice containing all items in the round-robin cycle.
func (r *RoundRobin) Items() (items []Item) {
	return r.items
}

// Add appends one or more unique string values to the round-robin cycle.
// Duplicate values are ignored to maintain uniqueness, using a thread-safe map.
// Parameters:
//   - values ([]string): Variable number of string values to add as Items.
func (r *RoundRobin) Add(values ...string) {
	for _, value := range values {
		item := Item{
			value: value,
		}

		if _, loaded := r.itemsMap.LoadOrStore(value, struct{}{}); !loaded {
			r.items = append(r.items, item)
		}
	}
}

// Next returns the next Item in the round-robin cycle based on the rotation policy.
// It uses a mutex for thread-safe access and updates the serve count atomically.
// Rotation occurs after the current item has been served RotateAmount times, as specified in Options.
// Returns:
//   - item (Item): The next Item in the cycle. If the index is invalid, the first item is returned.
func (r *RoundRobin) Next() (item Item) {
	r.mutex.Lock()

	defer r.mutex.Unlock()

	currentAmount := atomic.LoadUint32(&r.currentItemServesCount)

	if currentAmount >= uint32(r.Options.RotateAmount) {
		atomic.StoreUint32(&r.currentItemServesCount, 1)
		atomic.AddUint32(&r.nextItemIndex, 1)
	} else {
		atomic.AddUint32(&r.currentItemServesCount, 1)
	}

	nextItemIndex := (int(r.nextItemIndex) - 1) % len(r.items)

	if nextItemIndex < 0 || nextItemIndex > len(r.items) {
		r.items[0].Statistics.IncrementServesCount(1)

		return r.items[0]
	}

	r.items[nextItemIndex].Statistics.IncrementServesCount(1)

	return r.items[nextItemIndex]
}

type RoundRobinInterface interface {
	Items() (items []Item)
	Add(values ...string)
	Next() (item Item)
}

// Options configures the behavior of the RoundRobin instance.
// Fields:
//   - RotateAmount (int32): The number of times an item is served before rotating to the next item.
type Options struct {
	RotateAmount int32
}

var _ RoundRobinInterface = (*RoundRobin)(nil)

var DefaultOptions = Options{
	RotateAmount: 1,
}

// New creates a new RoundRobin instance with default options and the provided items.
// It initializes the round-robin cycle with the given items, ensuring at least one item is provided.
// Parameters:
//   - items ([]string): Variable number of string values to initialize the round-robin cycle.
//
// Returns:
//   - rr (*RoundRobin): A pointer to the initialized RoundRobin instance, or nil if an error occurs.
//   - err (error): An error if no items are provided, otherwise nil.
func New(items ...string) (rr *RoundRobin, err error) {
	return NewWithOptions(DefaultOptions, items...)
}

// NewWithOptions creates a new RoundRobin instance with the specified options and items.
// It initializes the round-robin cycle with the given items and configuration, ensuring at least one item is provided.
// Parameters:
//   - options (Options): Configuration options for the round-robin cycle, such as RotateAmount.
//   - items ([]string): Variable number of string values to initialize the round-robin cycle.
//
// Returns:
//   - rr (*RoundRobin): A pointer to the initialized RoundRobin instance, or nil if an error occurs.
//   - err (error): An error if no items are provided, otherwise nil.
func NewWithOptions(options Options, items ...string) (rr *RoundRobin, err error) {
	if len(items) == 0 {
		err = hqgoerrors.New("no items")

		return
	}

	rr = &RoundRobin{
		Options: options,
	}

	rr.Add(items...)
	rr.nextItemIndex = 1

	return
}
