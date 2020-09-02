package main

import (
	context "context"
	"math/rand"
	"time"

	"github.com/arschles/containerscaler/externalscaler"
	empty "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/emptypb"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type impl struct {
	reqCounter *reqCounter
}

func newImpl(reqCounter *reqCounter) *impl {
	return &impl{reqCounter: reqCounter}
}

func (e *impl) Ping(context.Context, *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func (e *impl) IsActive(ctx context.Context, scaledObject *externalscaler.ScaledObjectRef) (*externalscaler.IsActiveResponse, error) {
	return &externalscaler.IsActiveResponse{
		Result: true,
	}, nil
}

func (e *impl) GetMetricSpec(_ context.Context, sor *externalscaler.ScaledObjectRef) (*externalscaler.GetMetricSpecResponse, error) {
	return &externalscaler.GetMetricSpecResponse{
		MetricSpecs: []*externalscaler.MetricSpec{{
			MetricName: "earthquakeThreshold",
			TargetSize: 100,
		}},
	}, nil
}

func (e *impl) GetMetrics(_ context.Context, metricRequest *externalscaler.GetMetricsRequest) (*externalscaler.GetMetricsResponse, error) {
	return &externalscaler.GetMetricsResponse{
		MetricValues: []*externalscaler.MetricValue{{
			MetricName:  "earthquakeThreshold",
			MetricValue: int64(e.reqCounter.get()),
		}},
	}, nil
}

func (e *impl) New(_ context.Context, nr *externalscaler.NewRequest) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func (e *impl) Close(_ context.Context, sor *externalscaler.ScaledObjectRef) (*emptypb.Empty, error) {
	return &empty.Empty{}, nil
}
