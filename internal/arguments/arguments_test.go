package arguments

import "testing"

func TestCheck(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"no arguments", args{}, true},
		{"printall", args{args: []string{"printall"}}, false},
		{"add with correct args", args{args: []string{"add", "key", "value"}}, false},
		{"add with wrong args", args{args: []string{"add", "key"}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Check(tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("Check() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
