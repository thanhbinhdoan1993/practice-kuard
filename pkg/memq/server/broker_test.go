package memqserver

import (
	"testing"
)

func TestUUID(t *testing.T) {
	id, err := uuid()
	if err != nil {
		t.Errorf("error when generating id: %v", err)
	}
	t.Log(id)
}
