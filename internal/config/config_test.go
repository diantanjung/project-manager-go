package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    time.Duration
		wantErr bool
	}{
		{name: "days", value: "7d", want: 7 * 24 * time.Hour},
		{name: "hours", value: "2h", want: 2 * time.Hour},
		{name: "minutes", value: "15m", want: 15 * time.Minute},
		{name: "seconds", value: "30s", want: 30 * time.Second},
		{name: "bad unit", value: "1w", wantErr: true},
		{name: "bad number", value: "xm", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDuration(tt.value)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadEnvFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	err := os.WriteFile(path, []byte(`
EXISTING=from_file
NEW_VALUE="from env"
# ignored
`), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	t.Setenv("EXISTING", "from_process")
	if err := loadEnvFile(path); err != nil {
		t.Fatalf("loadEnvFile() error = %v", err)
	}
	if got := os.Getenv("EXISTING"); got != "from_process" {
		t.Fatalf("EXISTING = %q, want process env preserved", got)
	}
	if got := os.Getenv("NEW_VALUE"); got != "from env" {
		t.Fatalf("NEW_VALUE = %q, want file value", got)
	}
}
