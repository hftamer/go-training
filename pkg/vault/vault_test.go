package vault

import (
	"os"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		want    Vault
		wantErr bool
		key string
	}{
		{"no encoding key set", Vault{}, true, "" },
		{"new file initialization", Vault{encodingKey: "124", data: make(map[string]string)}, false, "124" },
		{"correct decoding key for the existing file", Vault{encodingKey: "124", data: make(map[string]string)}, false, "124" },
		{"incorrect decoding key", Vault{encodingKey: "123", data: make(map[string]string)}, true, "123" },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(PMGR_ENCODING_KEY, tt.key)
			defer os.Unsetenv(tt.key)
			got, err := New()
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
