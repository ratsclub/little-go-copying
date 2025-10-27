package assert

import (
	"errors"
	"fmt"
	"testing"
)

type mockT struct {
	testing.TB
	errorfCalled bool
	fatalfCalled bool
	message      string
}

func (m *mockT) Errorf(format string, args ...any) {
	m.errorfCalled = true
	m.message = format
}

func (m *mockT) Fatalf(format string, args ...any) {
	m.fatalfCalled = true
	m.message = format
}

func TestEqual(t *testing.T) {
	tests := []struct {
		name      string
		got, want any
		isEqual   bool
	}{
		{"success on nil values", nil, nil, true},
		{"success on equal ints", 5, 5, true},
		{"error on different ints", 5, 10, false},
		{"success on equal strings", "hello", "hello", true},
		{"error different strings", "hello", "world", false},
		{"success on identical structs", struct{ name string }{name: "john"}, struct{ name string }{name: "john"}, true},
		{"error on different structs", struct{ name string }{name: "john"}, struct{ name string }{name: "mary"}, false},
		{"error on different struct paddings", struct {
			name string
			age  int
		}{name: "john", age: 18}, struct {
			age  int
			name string
		}{name: "john", age: 18}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt := &mockT{t, false, false, ""}

			Equal(mt, tt.got, tt.want)

			if tt.isEqual && mt.errorfCalled {
				t.Errorf("Equal() = called=%v, want %v", mt.errorfCalled, tt.isEqual)
			}
		})
	}
}

func TestNotEqual(t *testing.T) {
	tests := []struct {
		name      string
		got, want any
		isEqual   bool
	}{
		{"success on nil values", nil, nil, true},
		{"error on equal ints", 5, 5, true},
		{"success on different ints", 5, 10, false},
		{"error on equal strings", "hello", "hello", true},
		{"success different strings", "hello", "world", false},
		{"error on identical structs", struct{ name string }{name: "john"}, struct{ name string }{name: "john"}, true},
		{"success on different structs", struct{ name string }{name: "john"}, struct{ name string }{name: "mary"}, false},
		{"success on different struct paddings", struct {
			name string
			age  int
		}{name: "john", age: 18}, struct {
			age  int
			name string
		}{name: "john", age: 18}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt := &mockT{t, false, false, ""}

			NotEqual(mt, tt.got, tt.want)
			if !tt.isEqual && mt.errorfCalled {
				t.Errorf("NotEqual() = called=%v, want %v", mt.errorfCalled, tt.isEqual)
			}
		})
	}
}

func TestTrue(t *testing.T) {
	tests := []struct {
		name      string
		got       bool
		expectErr bool
	}{
		{"got true", true, false},
		{"got false", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt := &mockT{t, false, false, ""}
			True(mt, tt.got)

			if mt.errorfCalled != tt.expectErr {
				t.Errorf("%s: True() called=%v, want %v", tt.name, mt.errorfCalled, tt.expectErr)
			}
		})
	}
}

func TestFalse(t *testing.T) {
	tests := []struct {
		name      string
		got       bool
		expectErr bool
	}{
		{"got false", false, false},
		{"got true", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt := &mockT{t, false, false, ""}
			False(mt, tt.got)

			if mt.errorfCalled != tt.expectErr {
				t.Errorf("%s: False() called=%v, want %v", tt.name, mt.errorfCalled, tt.expectErr)
			}
		})
	}
}

func TestNil(t *testing.T) {
	var p *int = nil
	var s []string = nil

	tests := []struct {
		name      string
		got       any
		expectErr bool
	}{
		{"got nil pointer", p, false},
		{"got nil slice", s, false},
		{"got non-nil int", 42, true},
		{"got non-nil string", "hello", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt := &mockT{t, false, false, ""}
			Nil(mt, tt.got)

			if mt.errorfCalled != tt.expectErr {
				t.Errorf("%s: Nil() called=%v, want %v", tt.name, mt.errorfCalled, tt.expectErr)
			}
		})
	}
}

func TestNotNil(t *testing.T) {
	var p *int = nil
	var s []string = nil

	tests := []struct {
		name      string
		got       any
		expectErr bool
	}{
		{"got non-nil int", 42, false},
		{"got non-nil string", "hello", false},
		{"got nil pointer", p, true},
		{"got nil slice", s, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt := &mockT{t, false, false, ""}
			NotNil(mt, tt.got)

			if mt.errorfCalled != tt.expectErr {
				t.Errorf("%s: NotNil() called=%v, want %v", tt.name, mt.errorfCalled, tt.expectErr)
			}
		})
	}
}

func TestErrorIs(t *testing.T) {
	baseErr := errors.New("base error")
	wrappedErr := fmt.Errorf("wrapped: %w", baseErr)
	otherErr := errors.New("different error")

	tests := []struct {
		name      string
		got, want error
		expectErr bool
	}{
		{"same error", baseErr, baseErr, false},
		{"wrapped error matches base", wrappedErr, baseErr, false},
		{"different errors", otherErr, baseErr, true},
		{"got nil error", nil, baseErr, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt := &mockT{t, false, false, ""}
			ErrorIs(mt, tt.got, tt.want)

			if mt.errorfCalled != tt.expectErr {
				t.Errorf("%s: ErrorIs() called=%v, want %v", tt.name, mt.errorfCalled, tt.expectErr)
			}
		})
	}
}

type customError struct {
	msg string
}

func (e *customError) Error() string { return e.msg }

func TestErrorAs(t *testing.T) {
	var target *customError

	baseErr := errors.New("base")
	custom := &customError{"custom"}
	wrapped := fmt.Errorf("wrapped: %w", custom)

	tests := []struct {
		name      string
		got       error
		target    any
		expectErr bool
	}{
		{"got matches target type", custom, &target, false},
		{"wrapped matches target type", wrapped, &target, false},
		{"base (non-matching) error", baseErr, &target, true},
		{"got is nil", nil, &target, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt := &mockT{t, false, false, ""}
			ErrorAs(mt, tt.got, tt.target)

			if mt.errorfCalled != tt.expectErr {
				t.Errorf("%s: ErrorAs() called=%v, want %v", tt.name, mt.errorfCalled, tt.expectErr)
			}
		})
	}
}

func TestMatchesRegexp(t *testing.T) {
	tests := []struct {
		name        string
		got         string
		pattern     string
		expectError bool
		expectFatal bool
	}{
		{
			name:        "string matches pattern",
			got:         "hello123",
			pattern:     `^hello\d+$`,
			expectError: false,
			expectFatal: false,
		},
		{
			name:        "string does not match pattern",
			got:         "hello",
			pattern:     `^\d+$`,
			expectError: true,
			expectFatal: false,
		},
		{
			name:        "invalid regexp pattern",
			got:         "anything",
			pattern:     `*invalid[`,
			expectError: false,
			expectFatal: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt := &mockT{t, false, false, ""}
			MatchesRegexp(mt, tt.got, tt.pattern)

			if mt.errorfCalled != tt.expectError {
				t.Errorf("%s: expected Errorf called=%v, got %v", tt.name, tt.expectError, mt.errorfCalled)
			}

			if mt.fatalfCalled != tt.expectFatal {
				t.Errorf("%s: expected Fatalf called=%v, got %v", tt.name, tt.expectFatal, mt.fatalfCalled)
			}
		})
	}
}
