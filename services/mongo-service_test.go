package services

import (
	"reflect"
	"testing"
)

func TestGetMongoService(t *testing.T) {
	tests := []struct {
		name    string
		wantNil bool
	}{
		{name: "test01", wantNil: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMongoService(); !reflect.DeepEqual(got, false) {
				t.Errorf("GetMongoService() = %v, want %v", got, false)
			}
		})
	}
}
