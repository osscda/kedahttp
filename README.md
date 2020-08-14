# Autoscaling Containers

This project implements a prototype of auto-scaling containers on ACI. As HTTP requests come into the system, the container(s) that are equipped to handle that request may or may not be running and ready to accept it. If there are sufficient containers available, the request is routed to one of them.  If there are not, a container is started and the request is routed to it when it's ready.

## Architecture

This system has three components:

- Proxy
- [KEDA](https://keda.sh)
- [Redis](https://redis.io)

The **proxy** receives incoming HTTP traffic, emits events to NATS streaming, and forwards to a backend container.

KEDA is responsible for consuming events from the proxy and scaling the backend containers appropriately.

## Installation

You need to install KEDA first. Do so with these commands:

```shell
helm repo add kedacore https://kedacore.github.io/charts
helm repo update
helm install keda kedacore/keda --namespace cscaler --create-namespace
```

>These commands are similar to those on the [official install page](https://keda.sh/docs/1.5/deploy/#helm), but we're installing in a different namespace.

You'll also need to install Redis:

```shell
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
helm install redis bitnami/redis -n cscaler --create-namespace \
    --set cluster.enabled=false \
    --set usePassword=true \
    --set password=abcd \
    --set cluster.enabled=false
```

## Install the Proxy & Dummy App (for now)


```shell
helm install -n cscaler cscaler ./charts/proxy --create-namespace \
    --set redisAddr=redis-master.cscaler.svc.cluster.local:6379 \
    --set redisPass=abcd
```

## More Information

See [this document](./docs/COMPONENTS.md) for details on the components of this system.

## Build

### cli

Just simply run ```make cli``` command

You can then install it into your ```PATH``` or add the ```./bin``` to your ```PATH```
