package countrycodes

import (
	"testing"
)

func TestAssignment_Valid(t *testing.T) {
	tests := []struct {
		name string
		a    Assignment
		want bool
	}{
		{
			name: "simple passing test",
			a:    0,
			want: true,
		},
		{
			name: "simple failing test",
			a:    10,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Valid(); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}
