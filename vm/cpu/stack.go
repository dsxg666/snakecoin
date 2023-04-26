package cpu

import "errors"

// Stack holds return-addresses when the `call` operation is being
// completed.  It can also be used for storing ints.
type Stack struct {
	// The entries on our stack
	entries []int
}

//
// Stack functions
//

// NewStack creates a new stack object.
func NewStack() *Stack {
	return &Stack{}
}

// Empty returns true if the stack is empty.
func (s *Stack) Empty() bool {
	return (len(s.entries) <= 0)
}

// Size retrieves the length of the stack.
func (s *Stack) Size() int {
	return (len(s.entries))
}

// Push adds a value to the stack.
func (s *Stack) Push(value int) {
	s.entries = append(s.entries, value)
}

// Pop removes a value from the stack.
func (s *Stack) Pop() (int, error) {
	if s.Empty() {
		return 0, errors.New("Pop from an empty stack")
	}

	// get top
	l := len(s.entries)
	top := s.entries[l-1]

	// truncate
	s.entries = s.entries[:l-1]

	return top, nil
}
