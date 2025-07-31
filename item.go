package roundrobin

// Item represents an element in the round-robin cycle, holding a value and its usage statistics.
// Fields:
//   - value (string): The string value associated with the item.
//   - Statistics (Statistics): Metrics tracking the item's usage, such as serve count.
type Item struct {
	value      string
	Statistics Statistics
}

// Value returns the string value of the Item.
// Returns:
//   - value (string): The string value stored in the Item.
func (i Item) Value() (value string) {
	value = i.value

	return
}

type ItemInterface interface {
	Value() (value string)
}

var _ ItemInterface = (*Item)(nil)
