package goroutine_pool

import "testing"

func TestNewError(t *testing.T) {
	pool := NewPool(-1, 1, 1, 1, 1, 1)
}
