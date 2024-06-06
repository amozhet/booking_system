package data

import (
	"encoding/json"
	"testing"
)

func TestRuntime_MarshalJSON(t *testing.T) {
	r := Runtime(105)
	want := `"105 mins"`
	got, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("Runtime.MarshalJSON() error = %v", err)
	}
	if string(got) != want {
		t.Errorf("Runtime.MarshalJSON() = %v, want %v", string(got), want)
	}
}

func TestRuntime_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		json    []byte
		want    Runtime
		wantErr bool
	}{
		{
			name:    "valid",
			json:    []byte(`"105 mins"`),
			want:    105,
			wantErr: false,
		},
		{
			name:    "invalid format",
			json:    []byte(`"105"`),
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid integer",
			json:    []byte(`"abc mins"`),
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r Runtime
			if err := json.Unmarshal(tt.json, &r); (err != nil) != tt.wantErr {
				t.Errorf("Runtime.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if r != tt.want {
				t.Errorf("Runtime.UnmarshalJSON() = %v, want %v", r, tt.want)
			}
		})
	}
}
