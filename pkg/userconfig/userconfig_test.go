package userconfig

import (
	"reflect"
	"testing"
)

func TestSanitizeUserConfig(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "sorts entries alphabetically",
			args: []string{"~/z/file.txt", "~/a/file.txt", "~/m/file.txt"},
			want: []string{"~/a/file.txt", "~/m/file.txt", "~/z/file.txt"},
		},
		{
			name: "removes duplicates and sorts",
			args: []string{"~/b.txt", "~/a.txt", "~/b.txt"},
			want: []string{"~/a.txt", "~/b.txt"},
		},
		{
			name: "trims whitespace and sorts",
			args: []string{" ~/b.txt ", " ~/a.txt "},
			want: []string{"~/a.txt", "~/b.txt"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sanitizeUserConfig(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sanitizeUserConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
