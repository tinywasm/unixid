package unixid_test

import (
	"testing"

	"github.com/cdvelop/unixid"
)

func TestValidate(t *testing.T) {
	uid, err := unixid.NewUnixID()
	if err != nil {
		t.Fatal("Error creating unixid:", err)
	}

	testCases := []struct {
		name  string
		input string
		valid bool
	}{
		{"valid server ID", "1624397134562544800", true},
		{"valid client ID", "1624397134562544800.42", true},
		{"invalid multiple dots", "1624397134562544800.42.42", false},
		{"invalid letter", "1624397134562544800a", false},
		{"invalid starts with dot", ".1624397134562544800", false},
		{"invalid ends with dot", "1624397134562544800.", false},
		{"invalid empty", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := uid.Validate(tc.input)
			if tc.valid && err != nil {
				t.Errorf("expected valid, got error: %v", err)
			}
			if !tc.valid && err == nil {
				t.Errorf("expected error, got nil")
			}
		})
	}
}
