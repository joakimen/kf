package slice

import (
	"reflect"
	"slices"
	"testing"
)

func TestExists(t *testing.T) {
	type args struct {
		needle   string
		haystack []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "needle exists in haystack",
			args: args{
				needle:   "a",
				haystack: []string{"a", "b", "c"},
			},
			want: true,
		},

		{
			name: "needle does not exist in haystack",
			args: args{
				needle:   "d",
				haystack: []string{"a", "b", "c"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Exists(tt.args.needle, tt.args.haystack); got != tt.want {
				t.Errorf("Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnique(t *testing.T) {
	type args struct {
		elements []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "removes duplicates preserving insertion order",
			args: args{
				elements: []string{"a", "b", "a", "c", "b"},
			},
			want: []string{"a", "b", "c"},
		},
		{
			name: "does not mutate input slice",
			args: args{
				elements: []string{"c", "a", "b", "a"},
			},
			want: []string{"c", "a", "b"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := slices.Clone(tt.args.elements)
			got := Unique(tt.args.elements)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unique() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(tt.args.elements, input) {
				t.Errorf("Unique() mutated input: got %v, original %v", tt.args.elements, input)
			}
		})
	}
}

func TestTrimWhitespace(t *testing.T) {
	type args struct {
		elements []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "leading and trailing whitespace is trimmed",
			args: args{elements: []string{" a   ", "\n\rb\t"}},
			want: []string{"a", "b"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TrimWhitespace(tt.args.elements); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TrimWhitespace() = %v, want %v", got, tt.want)
			}
		})
	}
}
