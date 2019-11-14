package flason

import (
	"testing"
)

func eq(a, b []JsonPair) bool {
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for index := range a {
		if a[index].Path != b[index].Path || a[index].Value != b[index].Value {
			return false
		}
	}

	return true
}

func TestBasicTypes(t *testing.T) {
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
			name: "Single value",
			args: args{
				str: `"a"`,
			},
			want: []JsonPair{
				JsonPair{
					Path:  "",
					Value: "a",
				},
			},
			wantErr: false,
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
		{
			name: "Basic null",
			args: args{
				str: `{ "key": null }`,
			},
			want: []JsonPair{
				JsonPair{
					Path:  ".key",
					Value: "null",
				},
			},
			wantErr: false,
		},
		{
			name: "Basic boolean",
			args: args{
				str: `{ "key": true, "otherKey": false }`,
			},
			want: []JsonPair{
				JsonPair{
					Path:  ".key",
					Value: "true",
				},
				JsonPair{
					Path:  ".otherKey",
					Value: "false",
				},
			},
			wantErr: false,
		},
		{
			name: "Basic array",
			args: args{
				str: `[ "a", 1, true ]`,
			},
			want: []JsonPair{
				JsonPair{
					Path:  "[0]",
					Value: "a",
				},
				JsonPair{
					Path:  "[1]",
					Value: "1",
				},
				JsonPair{
					Path:  "[2]",
					Value: "true",
				},
			},
			wantErr: false,
		},
		{
			name: "Basic object",
			args: args{
				str: `{ "inv": { "true": true, "false": false } }`,
			},
			want: []JsonPair{
				JsonPair{
					Path:  ".inv.false",
					Value: "false",
				},
				JsonPair{
					Path:  ".inv.true",
					Value: "true",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FlattenJson(tt.args.str, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("FlattenJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !eq(got, tt.want) {
				t.Errorf("FlattenJson() = %v, want %v", got, tt.want)
			}
		})
	}
}
