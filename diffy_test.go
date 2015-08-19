package diffy

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	c := NewClient()

	if c == nil {
		t.Errorf("NewClient return is nil, but expected a diffy.Client struct")
	}
}
