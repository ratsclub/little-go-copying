/*
Package assert provides a set of useful assert functions to use within the
standard testing system.

Code from [Alex Edward's blog post](https://www.alexedwards.net/blog/the-9-go-test-assertions-i-use).
*/
package assert

import (
	"errors"
	"reflect"
	"regexp"
	"testing"
)

// Equal asserts that got is equals to want
func Equal[T any](t testing.TB, got, want T) {
	t.Helper()
	if !isEqual(got, want) {
		t.Errorf("got: %v; want: %v", got, want)
	}
}

// NotEqual asserts that got is different than want
func NotEqual[T any](t testing.TB, got, want T) {
	t.Helper()
	if isEqual(got, want) {
		t.Errorf("got: %v; expected values to be different", got)
	}
}

// True asserts that got is true
func True(t testing.TB, got bool) {
	t.Helper()
	if !got {
		t.Errorf("got: false; want: true")
	}
}

// False asserts that got is false
func False(t testing.TB, got bool) {
	t.Helper()
	if got {
		t.Errorf("got: true; want: false")
	}
}

// Nil asserts that got is nil
func Nil(t testing.TB, got any) {
	t.Helper()
	if !isNil(got) {
		t.Errorf("got: %v; want: nil", got)
	}
}

// NotNil asserts that got is not nil
func NotNil(t testing.TB, got any) {
	t.Helper()
	if isNil(got) {
		t.Errorf("got: nil; want: non-nil")
	}
}

// ErrorIs asserts that got error is equal to the want error
func ErrorIs(t testing.TB, got, want error) {
	t.Helper()
	if !errors.Is(got, want) {
		t.Errorf("got: %v; want: %v", got, want)
	}
}

// ErrorAs asserts that got is an error that can be assigned to target via errors.As
func ErrorAs(t testing.TB, got error, target any) {
	t.Helper()
	if got == nil {
		t.Errorf("got: nil; want assignable to: %T", target)
		return
	}
	if !errors.As(got, target) {
		t.Errorf("got: %v; want assignable to: %T", got, target)
	}
}

// MatchesRegexp asserts that got matches the pattern
func MatchesRegexp(t testing.TB, got, pattern string) {
	t.Helper()
	matched, err := regexp.MatchString(pattern, got)
	if err != nil {
		t.Fatalf("unable to parse regexp pattern %s: %s", pattern, err.Error())
		return
	}
	if !matched {
		t.Errorf("got: %q; want to match %q", got, pattern)
	}
}

func isEqual[T any](got, want T) bool {
	if isNil(got) && isNil(want) {
		return true
	}

	equalable, ok := any(got).(interface{ Equal(T) bool })
	if ok {
		return equalable.Equal(want)
	}
	return reflect.DeepEqual(got, want)
}

func isNil(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice, reflect.UnsafePointer:
		return rv.IsNil()
	}
	return false
}
