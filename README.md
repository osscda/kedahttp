# Autoscaling Containers

This project implements a prototype of auto-scaling containers on ACI. As HTTP requests come into the system, the container(s) that are equipped to handle that request may or may not be running and ready to accept it. If there are sufficient containers available, the request is routed to one of them.  If there are not, a container is started and the request is routed to it when it's ready.

## Architecture

This system has three components:

- Proxy
- [KEDA](https://keda.sh)
- [Redis](https://redis.io)

The **proxy** receives incoming HTTP traffic, emits events to NATS streaming, and forwards to a backend container.

KEDA is responsible for consuming events from the proxy and scaling the backend containers appropriately.

## More Information

See [this document](./docs/COMPONENTS.md) for details on the components of this system.

## TODOs (notes from @asw101 and @arschles discussion)

- [ ] Add admin "control plane" API to this, and a CLI for it Issue #18
- [ ] Add a sidecar with a NATS server Issue #17
- [ ] Express container network policy API Issue #16
- [ ] Figure out the Front Door ingress situation Issue #14
- [ ] Once ^^ is done, modify the scaling controller Issue #15
- [ ] Make the proxy / controller multi-tenant Issue #12
- [ ] Similar to ^^, add some strategies for general scaling. Issue #13
- [ ] Logging mechanism to aggregate all of the requests from the "proxies" (e.g. edge locations) Issue #11
- [ ] Dump traces from containers in from edge to your container, and from your container to other Azure services / generic URLs to App Insights Issue #10
- [ ] Support versions and rollout events for the customer's app. Issue #9
- [ ] Preview your code in a PR using a GitHub action Issue #19

## Build

### cli

Just simply run ```make cli``` command

You can then install it into your ```PATH``` or add the ```./bin``` to your ```PATH```
