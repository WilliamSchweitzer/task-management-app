package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		// ----- Happy path -----
		{name: "simple", email: "a@b.c", wantErr: false},
		{name: "realistic", email: "john.doe@example.com", wantErr: false},
		{name: "unicode local (allowed)", email: "café@domain.com", wantErr: false},
		// {name: "long but valid", email: string(make([]byte, 64)) + "@" + string(make([]byte, 189)), wantErr: false},

		// ----- Length errors -----
		{name: "too short", email: "a@b", wantErr: true},
		{name: "too long", email: string(make([]byte, 255)) + "@a.b", wantErr: true},

		// ----- @ errors -----
		{name: "no at", email: "abcd", wantErr: true},
		{name: "multiple at", email: "a@b@c.com", wantErr: true},
		{name: "at start", email: "@b.c", wantErr: true},
		{name: "at end", email: "a@b.", wantErr: true},

		// ----- Local-part errors -----
		{name: "local control char", email: "a\x00b@c.d", wantErr: true},
		{name: "local space", email: "a b@c.d", wantErr: true},
		{name: "local special", email: "a(b)@c.d", wantErr: true},
		{name: "local starts dot", email: ".a@c.d", wantErr: true},
		{name: "local ends dot", email: "a.@c.d", wantErr: true},

		// ----- Domain errors -----
		{name: "domain no dot", email: "a@bcdef", wantErr: true},
		{name: "domain leading dot", email: "a@.bc", wantErr: true},
		{name: "domain trailing dot", email: "a@bc.", wantErr: true},
		{name: "domain invalid char", email: "a@b_c.d", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateName(t *testing.T) {
	validNames := []string{
		"John Doe",
		"Jane",
		"Al",
	}
	invalidNames := []string{
		"A", // too short
		"",  // empty
	}

	for _, name := range validNames {
		err := ValidateName(name)
		assert.NoError(t, err, "Expected valid name: %s", name)
	}

	for _, name := range invalidNames {
		err := ValidateName(name)
		assert.Error(t, err, "Expected invalid name: %s", name)
	}
}
