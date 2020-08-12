# Autoscaling Containers

This project implements a prototype of auto-scaling containers on ACI. As HTTP requests come into the system, the container(s) that are equipped to handle that request may or may not be running and ready to accept it. If there are sufficient containers available, the request is routed to one of them.  If there are not, a container is started and the request is routed to it when it's ready.

## Architecture

This system has three components:

- Proxy
- [KEDA](https://keda.sh)
- [Redis](https://redis.io)

The **proxy** receives incoming HTTP traffic, emits events to NATS streaming, and forwards to a backend container.

KEDA is responsible for consuming events from the proxy and scaling the backend containers appropriately.

The entire thing gets installed on Kubernetes with [Helm](https://helm.sh). Here's the command to install it (from the [./charts/cscaler-proxy](./charts/cscaler-proxy) directory):

```shell
helm install -n cscaler cscaler .
```

### Install Keda

Follow the instructions here to install Keda: 

https://keda.sh/docs/1.5/deploy/#helm

## More Information

See [this document](./docs/COMPONENTS.md) for details on the components of this system.

## Build

### cli

Just simply run ```make cli``` command

You can then install it into your ```PATH``` or add the ```./bin``` to your ```PATH```
