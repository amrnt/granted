package granted

import (
	"testing"
)

func TestDefault(t *testing.T) {
	if DefaultConfig != Default.Config {
		t.Error("Error")
	}
}
