package twerk

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		v      interface{}
		config Config
	}
	tests := []struct {
		name    string
		args    args
		want    *twerk
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.v, tt.args.config)
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
