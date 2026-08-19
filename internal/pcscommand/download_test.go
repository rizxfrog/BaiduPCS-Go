package pcscommand

import "testing"

func TestResolveDownloadNoCheck(t *testing.T) {
	tests := []struct {
		name        string
		configValue bool
		noCheck     bool
		forceCheck  bool
		want        bool
	}{
		{name: "enabled by config", configValue: true, want: true},
		{name: "disabled by config", configValue: false, want: false},
		{name: "no-check flag", noCheck: true, want: true},
		{name: "force check overrides config", configValue: true, forceCheck: true, want: false},
		{name: "force check overrides no-check", noCheck: true, forceCheck: true, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveDownloadNoCheck(tt.configValue, tt.noCheck, tt.forceCheck); got != tt.want {
				t.Fatalf("resolveDownloadNoCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}
