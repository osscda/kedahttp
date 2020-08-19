package externalscaler

import (
	context "context"
	"log"
	"math/rand"
	"time"

	empty "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/emptypb"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Impl struct{}

func (e *Impl) IsActive(ctx context.Context, scaledObject *ScaledObjectRef) (*IsActiveResponse, error) {
	log.Printf("external.IsActive")
	return &IsActiveResponse{
		Result: true,
	}, nil
}

func (e *Impl) GetMetricSpec(context.Context, *ScaledObjectRef) (*GetMetricSpecResponse, error) {
	log.Printf("external.GetMetricSpec")
	return &GetMetricSpecResponse{
		MetricSpecs: []*MetricSpec{{
			MetricName: "earthquakeThreshold",
			TargetSize: 10,
		}},
	}, nil
}

func (e *Impl) GetMetrics(_ context.Context, metricRequest *GetMetricsRequest) (*GetMetricsResponse, error) {
	log.Printf("external.GetMetrics")
	return &GetMetricsResponse{
		MetricValues: []*MetricValue{{
			MetricName:  "earthquakeThreshold",
			MetricValue: int64(rand.Intn(10)),
		}},
	}, nil
}

func (e *Impl) New(context.Context, *NewRequest) (*empty.Empty, error) {
	log.Printf("external.New")
	return nil, nil
}

func (e *Impl) Close(_ context.Context, _ *ScaledObjectRef) (*emptypb.Empty, error) {
	log.Printf("external.Close")
	return nil, nil
}
