package assert_test

import (
	"testing"

	"github.com/ratsclub/little-go-copying/assert"
)

func ExampleEqual() {
	t := &testing.T{} // provided by the test

	tests := []struct {
		name      string
		got, want any
	}{
		{
			name: "should be equal",
			got:  "john doe",
			want: "john doe",
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.got, tt.want)
	}
}
