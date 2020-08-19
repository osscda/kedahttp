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

var counter int

type Impl struct{}

func (e *Impl) IsActive(ctx context.Context, scaledObject *ScaledObjectRef) (*IsActiveResponse, error) {
	log.Printf("external.IsActive: %+v", *scaledObject)
	return &IsActiveResponse{
		Result: true,
	}, nil
}

func (e *Impl) GetMetricSpec(_ context.Context, sor *ScaledObjectRef) (*GetMetricSpecResponse, error) {
	log.Printf("external.GetMetricSpec: %+v", *sor)
	return &GetMetricSpecResponse{
		MetricSpecs: []*MetricSpec{{
			MetricName: "earthquakeThreshold",
			TargetSize: 10,
		}},
	}, nil
}

func (e *Impl) GetMetrics(_ context.Context, metricRequest *GetMetricsRequest) (*GetMetricsResponse, error) {
	log.Printf("external.GetMetrics: %+v", *metricRequest)
	counter++
	log.Printf("counter: %d", counter)
	return &GetMetricsResponse{
		MetricValues: []*MetricValue{{
			MetricName:  "earthquakeThreshold",
			MetricValue: int64(counter),
		}},
	}, nil
}

func (e *Impl) New(_ context.Context, nr *NewRequest) (*empty.Empty, error) {
	log.Printf("external.New: %+v", *nr)
	return &empty.Empty{}, nil
}

func (e *Impl) Close(_ context.Context, sor *ScaledObjectRef) (*emptypb.Empty, error) {
	log.Printf("external.Close: %+v", *sor)
	return &empty.Empty{}, nil
}
