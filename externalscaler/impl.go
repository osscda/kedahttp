package externalscaler

import (
	context "context"
	"math/rand"
	"sync/atomic"
	"time"

	empty "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/emptypb"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var counter int

type Impl struct {
	reqCounter *int64
}

func NewImpl(reqCounter *int64) *Impl {
	return &Impl{reqCounter: reqCounter}
}

func (e *Impl) Ping(context.Context, *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func (e *Impl) IsActive(ctx context.Context, scaledObject *ScaledObjectRef) (*IsActiveResponse, error) {
	return &IsActiveResponse{
		Result: true,
	}, nil
}

func (e *Impl) GetMetricSpec(_ context.Context, sor *ScaledObjectRef) (*GetMetricSpecResponse, error) {
	return &GetMetricSpecResponse{
		MetricSpecs: []*MetricSpec{{
			MetricName: "earthquakeThreshold",
			TargetSize: 100,
		}},
	}, nil
}

func (e *Impl) GetMetrics(_ context.Context, metricRequest *GetMetricsRequest) (*GetMetricsResponse, error) {
	counter := atomic.LoadInt64(e.reqCounter)
	return &GetMetricsResponse{
		MetricValues: []*MetricValue{{
			MetricName:  "earthquakeThreshold",
			MetricValue: counter,
		}},
	}, nil
}

func (e *Impl) New(_ context.Context, nr *NewRequest) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func (e *Impl) Close(_ context.Context, sor *ScaledObjectRef) (*emptypb.Empty, error) {
	return &empty.Empty{}, nil
}
