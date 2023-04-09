package util

import (
	"reflect"
	"testing"
)

func TestMerge(t *testing.T) {
	type args[T any] struct {
		a []T
		b []T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "nil",
			args: args[int]{
				a: nil,
				b: nil,
			},
			want: nil,
		},
		{
			name: "a nil",
			args: args[int]{
				a: nil,
				b: []int{4, 5, 6},
			},
			want: []int{4, 5, 6},
		},
		{
			name: "b nil",
			args: args[int]{
				a: []int{1, 2, 3},
				b: nil,
			},
			want: []int{1, 2, 3},
		},
		{
			name: "both not nil",
			args: args[int]{
				a: []int{1, 2, 3},
				b: []int{4, 5, 6},
			},
			want: []int{1, 2, 3, 4, 5, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Merge(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Merge() = %v, want %v", got, tt.want)
			}
		})
	}
}
