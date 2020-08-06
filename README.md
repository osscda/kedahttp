# Autoscaling Containers

This project implements a prototype of auto-scaling containers on ACI. As HTTP requests come into the system, the container(s) that are equipped to handle that request may or may not be running and ready to accept it. If there are sufficient containers available, the request is routed to one of them.  If there are not, a container is started and the request is routed to it when it's ready.

>Although the featureset is comparatively basic, this project is similar in concept to [Knative serving](https://knative.dev/docs/serving/) or [Keda](https://keda.sh). 

There are some major differences, though:

- Simpler to install
    - There are two components: HTTP proxy and a scaling controller. No service mesh required
- No routes or versions (yet)
- Single tenant (at the moment)

## Architecture

This system has three components:

- Proxy
- [KEDA](https://keda.sh)
- [NATS streaming](https://docs.nats.io/nats-streaming-concepts/intro)

The **proxy** receives incoming HTTP traffic, emits events to NATS streaming, and forwards to a backend container.

KEDA is responsible for consuming events from the proxy and scaling the backend containers appropriately.

## More Information

See [this document](./docs/COMPONENTS.md) for details on the components of this system.

## TODOs (notes from @asw101 and @arschles discussion)

- [ ] Add admin "control plane" API to this, and a CLI for it
    - `csclr deploy hello-world:latest --platform=VMSS or --platform=ACI ...`
    - Also provide a standards-compliant YAML deployment (i.e. KEDA YAML or KNative `Service` YAML)
- [ ] Add a sidecar with a NATS server
- [ ] Express container network policy API - translate it to underlying service mesh API
    - Use LinkerD for service mesh?
- [ ] Figure out the Front Door ingress situation
- [ ] Once ^^ is done, modify the scaling controller to know (a) what ideal region to create new containers and (b) the "backup" regions (i.e. priority list) to spin containers up in
- [ ] Make the proxy / controller multi-tenant
- [ ] Similar to ^^, add some strategies for general scaling. For example:
    - Primary/secondary: scale first in regionA, then in regionB. i.e. a "backup" scenario
    - Geographic redundancy: on any scale event, scale up in regionA, regionB and regionC
    - Geographic load balancing: scale up in the region closest to me, where there is capacity available. have a "backup" list of regions to scale up in, if the preferred region on that scale event is unavailable
- [ ] Logging mechanism to aggregate all of the requests from the "proxies" (e.g. edge locations) and dump reports into blob store, kusto(?), Azure Arc(?), big data analytics of other kind
    - Probably in structured JSON
    - Also implement prometheus API (statsd?)
    - Also dump scale events
- [ ] Dump traces from containers in from edge to your container, and from your container to other Azure services / generic URLs to App Insights
- [ ] Support versions and rollout events for the customer's app. A/B, Green/Blue and canary. Integrate with tracing and logging
    - Tied to "native" scaling events?
    - Lean on Front Door for the routing? Just need to clean up the versioned endpoints. See [this ARM template from AaronW](https://github.com/aaronmsft/aaronmsft-com/blob/master/azure-front-door-container-instances-arm/azuredeploy.json) for example on how to do this
- [ ] Preview your code in a PR using a GitHub action
    - Have it roll out a new version

## Build

### cli

Just simply run ```make cli``` command

You can then install it into your ```PATH``` or add the ```./bin``` to your ```PATH```
