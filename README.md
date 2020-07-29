# Autoscaling Containers

This project implements a prototype of auto-scaling containers on either Kubernetes or ACI. Although the featureset is very basic, it's similar in concept to [Knative serving](https://knative.dev/docs/serving/). There are some major differences, though:

- Simpler to install
    - There are two components: HTTP proxy and a scaling controller. No service mesh required
- No routes or versions (yet)
- Single tenant (at the moment)

## Architecture

This system has two components:

The **proxy** receives incoming HTTP traffic, looks up where to send that traffic in its database, and forwards it on to a URL. This URL can be any DNS name or IP. On Kubernetes, this URL should be a `Service` DNS name. For forwarding to ACI, this URL should be either a public IP or DNS name of one or more ACI containers. If a request comes in for a container that's not yet available, the proxy will wait for one to become available.

The **scaling controller** is responsible for fetching traffic metrics from the proxy and scaling containers up and down based on request volume.

## How to Run This

The proxy and controller both depend on NATS, so that's the first thing to run. Do so with Docker and this command:

```shell
docker run -p 4222:4222 -ti nats:latest
```

Or, if you don't want to use Docker, you can install NATS as a binary. Follow the directions in the [installation page](https://docs.nats.io/nats-server/installation) for how to do it. Note that the Mac Homebrew installation instructions work for Linux and Linuxbrew. If you use Linuxbrew, you'll see a warning that it's a Mac-specific installer. That's fine and won't affect you. Simply run `nats-server` on the command line to get running.

## FAQ

_Why don't you use Horizontal Pod Autoscaling, ingress controllers, or service meshes to do these things?_

Because those systems don't work for independently running processes that need to forward to ACI containers. They are generally for Kubernetes or other container orchestrators only.

## TODOs

- [ ] Add a sidecar with a NATS server 
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
