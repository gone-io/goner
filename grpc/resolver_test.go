package grpc

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	mock "github.com/gone-io/gone"
	"github.com/gone-io/goner/g"
	gMock "github.com/gone-io/goner/g/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/resolver"
)

func TestResolverBuilder_Build(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	service := gMock.NewMockService(controller)
	service.EXPECT().GetIP().Return("127.0.0.1").AnyTimes()
	service.EXPECT().GetPort().Return(8080).AnyTimes()
	service.EXPECT().GetName().Return("svc1").AnyTimes()
	service.EXPECT().GetWeight().Return(100.0).AnyTimes()

	tests := []struct {
		name           string
		serviceName    string
		instances      []g.Service
		watchCh        chan []g.Service
		getInstanceErr error
		watchErr       error
		expectErr      bool
		scheme         string
	}{
		{
			name:        "success case",
			serviceName: "test-service",
			instances: []g.Service{
				service,
			},
			watchCh:   make(chan []g.Service),
			expectErr: false,
			scheme:    "dns",
		},
		{
			name:           "get instances error",
			serviceName:    "test-service",
			getInstanceErr: fmt.Errorf("get instances error"),
			watchCh:        make(chan []g.Service),
			expectErr:      true,
			scheme:         "dns",
		},
		{
			name:        "watch error",
			serviceName: "test-service",
			watchErr:    fmt.Errorf("watch error"),
			expectErr:   true,
			scheme:      "dns",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			discovery := gMock.NewMockServiceDiscovery(controller)
			logger := mock.NewMockLogger(controller)
			cc := NewMockClientConn(controller)

			discovery.EXPECT().Watch(tt.serviceName).Return(tt.watchCh, func() error { return nil }, tt.watchErr).AnyTimes()
			discovery.EXPECT().GetInstances(tt.serviceName).Return(tt.instances, tt.getInstanceErr).AnyTimes()

			cc.EXPECT().UpdateState(gomock.Any()).Return(nil).AnyTimes()

			// Create resolver builder and build resolver
			builder := NewResolverBuilder(discovery, logger)
			r, err := builder.Build(resolver.Target{
				URL: url.URL{
					Scheme: tt.scheme,
					Path:   tt.serviceName,
				},
			}, cc, resolver.BuildOptions{})

			// Verify expectations
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, r)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, r)

				// Test Close method
				logger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
				r.Close()
			}
		})
	}
}

func TestDiscoveryResolver_ResolveNow(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	service := gMock.NewMockService(controller)
	service.EXPECT().GetIP().Return("127.0.0.1").AnyTimes()
	service.EXPECT().GetPort().Return(8080).AnyTimes()
	service.EXPECT().GetName().Return("svc1").AnyTimes()
	service.EXPECT().GetWeight().Return(100.0).AnyTimes()

	tests := []struct {
		name           string
		serviceName    string
		instances      []g.Service
		getInstanceErr error
		expectLogError bool
		updateStateErr error
	}{
		{
			name:        "success case",
			serviceName: "test-service",
			instances: []g.Service{
				service,
			},
			expectLogError: false,
		},
		{
			name:           "get instances error",
			serviceName:    "test-service",
			getInstanceErr: fmt.Errorf("get instances error"),
			expectLogError: true,
		},
		{
			name:           "update state error",
			serviceName:    "test-service",
			instances:      []g.Service{service},
			updateStateErr: fmt.Errorf("update state error"),
			expectLogError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			discovery := gMock.NewMockServiceDiscovery(controller)
			logger := mock.NewMockLogger(controller)
			cc := NewMockClientConn(controller)

			discovery.EXPECT().GetInstances(tt.serviceName).Return(tt.instances, tt.getInstanceErr)
			if tt.expectLogError {
				logger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
			}
			if tt.instances != nil && tt.getInstanceErr == nil {
				cc.EXPECT().UpdateState(gomock.Any()).Return(tt.updateStateErr)
				if tt.updateStateErr != nil {
					logger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
				}
			}

			// Create resolver
			r := &discoveryResolver{
				discovery:   discovery,
				logger:      logger,
				cc:          cc,
				serviceName: tt.serviceName,
			}

			// Call ResolveNow
			r.ResolveNow(resolver.ResolveNowOptions{})
		})
	}
}

func TestDiscoveryResolver_Watch(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	service := gMock.NewMockService(controller)
	service.EXPECT().GetIP().Return("127.0.0.1").AnyTimes()
	service.EXPECT().GetPort().Return(8080).AnyTimes()
	service.EXPECT().GetName().Return("svc1").AnyTimes()
	service.EXPECT().GetWeight().Return(100.0).AnyTimes()

	tests := []struct {
		name        string
		services    []g.Service
		updateError error
	}{
		{
			name:     "success case",
			services: []g.Service{service},
		},
		{
			name:        "update error",
			services:    []g.Service{service},
			updateError: fmt.Errorf("update error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := mock.NewMockLogger(controller)
			cc := NewMockClientConn(controller)
			updateCh := make(chan []g.Service, 1)

			if tt.updateError != nil {
				logger.EXPECT().Errorf(gomock.Any(), gomock.Any())
			}

			cc.EXPECT().UpdateState(gomock.Any()).Return(tt.updateError)

			// Create resolver
			r := &discoveryResolver{
				logger:   logger,
				cc:       cc,
				updateCh: updateCh,
			}

			// Start watch in goroutine
			go r.watch()

			// Send update
			updateCh <- tt.services

			// Wait for update to be processed
			time.Sleep(100 * time.Millisecond)
		})
	}
}
