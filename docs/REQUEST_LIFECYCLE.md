# The Lifecycle of a Request

## The Proxy

The proxy has three purposes:

- Create Kubernetes [`Deployment`](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/)s when a new application is created
    - This is not yet done
- Accept and forward incoming requests to that `Deployment`
- Notify the backend scaling component when a new request comes in

A new request into the system will reach the proxy first. When it does, the proxy first publishes an event over [NATS streaming](https://github.com/nats-io/stan.go). After this point, it forwards the request to a [Kubernetes `Service`](https://kubernetes.io/docs/concepts/services-networking/service/).

Behind the scenes, the events that the proxy publishes are consumed by [KEDA](https://keda.sh), which is responsible for scaling up the `Deployment`s.

 