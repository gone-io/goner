package g

import (
	"reflect"
	"testing"
)

func TestNewService(t *testing.T) {
	type args struct {
		name    string
		ip      string
		port    int
		meta    Metadata
		healthy bool
		weight  float64
	}
	tests := []struct {
		name string
		args args
		want Service
	}{
		{
			name: "test",
			args: args{
				name:    "test",
				ip:      "127.0.0.1",
				port:    8080,
				meta:    Metadata{},
				healthy: true,
				weight:  1.0,
			},
			want: &service{
				name:    "test",
				ip:      "127.0.0.1",
				port:    8080,
				meta:    Metadata{},
				healthy: true,
				weight:  1.0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewService(tt.args.name, tt.args.ip, tt.args.port, tt.args.meta, tt.args.healthy, tt.args.weight); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService(t *testing.T) {
	metadata := Metadata{
		"test": "test",
	}

	newService := NewService("test", "127.0.0.1", 8080, metadata, true, 1.0)
	t.Run("test", func(t *testing.T) {
		if newService.GetName() != "test" {
			t.Errorf("GetName() = %v, want %v", newService.GetName(), "test")
		}
		if newService.GetIP() != "127.0.0.1" {
			t.Errorf("GetIP() = %v, want %v", newService.GetIP(), "127.0.0.1")
		}
		if newService.GetPort() != 8080 {
			t.Errorf("GetPort() = %v, want %v", newService.GetPort(), 8080)
		}
		if newService.IsHealthy() != true {
			t.Errorf("IsHealthy() = %v, want %v", newService.IsHealthy(), true)
		}
		if newService.GetMetadata()["test"] != metadata["test"] {
			t.Errorf("GetMetadata() = %v, want %v", newService.GetMetadata(), Metadata{})
		}
		if newService.GetWeight() != 1.0 {
			t.Errorf("GetWeight() = %v, want %v", newService.GetWeight(), 1.0)
		}
	})
}
