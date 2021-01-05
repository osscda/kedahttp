package main

import (
	context "context"
	"math/rand"
	"time"

	"github.com/arschles/containerscaler/externalscaler"
	empty "github.com/golang/protobuf/ptypes/empty"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type impl struct {
	reqCounter *reqCounter
	externalscaler.UnimplementedExternalScalerServer
}

func newImpl(reqCounter *reqCounter) *impl {
	return &impl{reqCounter: reqCounter}
}

func (e *impl) Ping(context.Context, *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func (e *impl) IsActive(
	ctx context.Context,
	in *externalscaler.ScaledObjectRef,
) (*externalscaler.IsActiveResponse, error) {
	return &externalscaler.IsActiveResponse{
		Result: true,
	}, nil
}

func (e *impl) StreamIsActive(
	in *externalscaler.ScaledObjectRef,
	server externalscaler.ExternalScaler_StreamIsActiveServer,
) error {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-server.Context().Done():
			return nil
		case <-ticker.C:
			server.Send(&externalscaler.IsActiveResponse{
				Result: true,
			})
		}
	}
	return nil
}

func (e *impl) GetMetricSpec(
	ctx context.Context,
	in *externalscaler.ScaledObjectRef,
) (*externalscaler.GetMetricSpecResponse, error) {
	return &externalscaler.GetMetricSpecResponse{
		MetricSpecs: []*externalscaler.MetricSpec{
			{
				MetricName: "proxyCounter",
				TargetSize: 100,
			},
		},
	}, nil
}
func (e *impl) GetMetrics(
	ctx context.Context,
	in *externalscaler.GetMetricsRequest,
) (*externalscaler.GetMetricsResponse, error) {
	return &externalscaler.GetMetricsResponse{
		MetricValues: []*externalscaler.MetricValue{
			{
				MetricName:  "proxyCounter",
				MetricValue: int64(e.reqCounter.get()),
			},
		},
	}, nil
}
