# How Autoscaling Works

The controller is responsible for autoscaling. We're using approximately the Knative [autoscaling algorithms](https://knative.dev/docs/serving/autoscaling/). Specifically, we're implementing 

```
desiredReplicas = ceil[currentReplicas * ( currentMetricValue / desiredMetricValue )]
```

## What is a Metric?

The HPA algorithm was designed assuming the system could get specific metrics from pods in real time, and use those measurements as the metrics to decide in real time whether to scale.

The controller only gets a notification of requests from the controller in real time, so we have to work with a simple counter. Here's how we're defining our "metric" using the counter:

```
$REQUESTS_SINCE_LAST_TICK / $NUM_SERVERS_AT_TICK
```

The metric essentially describes how much request load was pointed at the entire fleet of backend servers

There are a few terms to understand here. The controller sleeps and wakes up periodically to calculate the metric and decide whether it needs to autoscale. When it wakes up, we call this a **tick** (we named this after Go's [`time.Ticker`](https://pkg.go.dev/time?tab=doc#Ticker)). We're also calling the time between ticks a **duration**.




Here's the above autoscale algorithm expanded to include the metric:

```
desiredReplicas = ceil[currentReplicas * ( $REQUESTS_SINCE_LAST_TICK / $TICK_DURATION / $NUM_SERVERS_AT_TICK) / desiredMetricValue )]
```