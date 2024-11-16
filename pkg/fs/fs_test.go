package fs

import (
	"path/filepath"
	"testing"
)

const (
	homeDirAbs = "/Users/kevin"
	pwdAbs     = "/Users/kevin/fake/dir"
)

func getenv(env string) string {
	switch env {
	case "HOME":
		return homeDirAbs
	case "PWD":
		return pwdAbs
	default:
		return ""
	}
}

func TestExpandFilePath(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{
			name: "abspaths outside user dir should be left alone",
			got:  "/etc/passwd",
			want: "/etc/passwd",
		},
		{
			name: "abspath should be cleaned",
			got:  "/relative/path/../path/file.txt",
			want: "/relative/path/file.txt",
		},
		{
			name: "user dir paths should be shrunk with tilde",
			got:  filepath.Join(homeDirAbs, "file.txt"),
			want: "~/file.txt",
		},
		{
			name: "tilde should be left alone",
			got:  "~/mydir/file.txt",
			want: "~/mydir/file.txt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SanitizeFilePath(tt.got, getenv)
			if err != nil {
				t.Errorf("SanitizeFilePath() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("SanitizeFilePath() got = %v, want %v", got, tt.want)
			}
		})
	}
}
