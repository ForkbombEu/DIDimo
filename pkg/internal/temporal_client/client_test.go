package temporalclient

import "testing"

func TestTemporalClientGet(t *testing.T) {
	have, err := GetTemporalClient()
	if err != nil {
		t.Error(err)
	}

	if have == nil {
		t.Error("Expected a non-nil temporal client, got nil")
	}
}
