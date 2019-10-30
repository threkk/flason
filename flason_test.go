package flason

import (
	"reflect"
	"testing"
)

func TestFlattenJson(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    []JsonPair
		wantErr bool
	}{
		{
			name: "Invalid JSON",
			args: args{
				str: "Not a JSON",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Basic string",
			args: args{
				str: `{ "key": "value" }`,
			},
			want: []JsonPair{
				JsonPair{
					Path:  ".key",
					Value: "value",
				},
			},
			wantErr: false,
		},
		{
			name: "Basic number",
			args: args{
				str: `{ "key": 1.5 }`,
			},
			want: []JsonPair{
				JsonPair{
					Path:  ".key",
					Value: "1.5",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FlattenJson(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("FlattenJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FlattenJson() = %v, want %v", got, tt.want)
			}
		})
	}
}
