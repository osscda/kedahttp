# Components of the System

The entire system runs inside of Kubernetes and takes advantage of Kubernetes features. It cannot run outside of Kubernetes, but this functionality might be built in the future.

## The Proxy

The proxy is primarily responsible for accepting requests from the internet and forwarding them to the right _backend_. A backend is a set of containers that can scale up and down, including to 0 containers.

### Processing a Request

An incoming request to the system will reach the proxy first. When it does, the proxy first publishes an event to [Redis](https://redis.io). After this point, it forwards the request to a [Kubernetes `Service`](https://kubernetes.io/docs/concepts/services-networking/service/) for the backend intended to serve the request. This backend might not currently have any running containers.

KEDA (detailed below) is responsible for scaling (up and down) the pods that the `Service` load balances over.

### Creating a New Backend

The proxy also has an "admin" API that is not intended to be exposed to the internet without authentication. Generally speaking, a command line tool or cloud portal would do operations on this API.

The API supports two major operations: (a) create new backend and (b) delete backend.

#### Create new Backend

When the proxy gets a request to create a new backend, it does the following:

- Create a new Kubernetes [`Deployment`](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/)s and [`Service`](https://kubernetes.io/docs/concepts/services-networking/service/)s when a new application is created
- Create a Kubernetes [`Service`](https://kubernetes.io/docs/concepts/services-networking/service/) that forwards to and load-balances over the pods in the `Deployment`
- Create a new KEDA [`ScaledObject`](https://keda.sh/docs/1.5/concepts/scaling-deployments/#scaledobject-spec) to indicate that the pods in the deployment should be scaled based on NATS traffic (more on all of this next)

#### Delete a Backend

Deleting a backend reverses the steps in the "create" step. That means it will delete all of the resources that it created there.

## [KEDA](https://keda.sh)

Behind the scenes, the prometheus the proxy publishes are consumed by [KEDA](https://keda.sh). KEDA is responsible for responding to these events and scale up / down the number of pods in the `Deployment` mentioned above. The `Service` that the proxy forwards to provides a stable endpoint that load balances over all of the pods, regardless of how many there are.

