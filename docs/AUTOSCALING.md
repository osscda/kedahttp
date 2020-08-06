# How Autoscaling Works

The controller is responsible for autoscaling. We're using a portion of the Knative [autoscaling algorithms](https://knative.dev/docs/serving/autoscaling/). Specifically, we're using a subset of the configuration parameters available in the [Knative pod autoscaler](https://knative.dev/v0.13-docs/serving/configuring-autoscaling/). Instead of autoscaling pods in Kubernetes, though, we're autoscaling backend containers. Below is a list of _all_ of the Knative autoscaler configuration values complete with a description of how or whether we support it.

>See [the `serving-core` sample configuration](https://github.com/knative/serving/releases/download/v0.14.0/serving-core.yaml) (under the `config-autoscaler` `ConfigMap`) for exhaustive documentation on what each of these configuration values mean.

- `scale-to-zero-grace-period`
- `enable-scale-to-zero` - set to `true` by default
- `tick-interval` - the interval between autoscaling calculations
- `scale-to-zero-grace-period` - time that all the containers are left inactive before scaling to zero

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
