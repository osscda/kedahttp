# TODOs

- [ ] Build an [external scaler](https://keda.sh/docs/1.5/concepts/external-scalers/)
    KEDA is built for the to-be-scaled container to consume events from the same place that KEDA looks for them, so it makes certain assumptions that tie it more tightly to the applications than we'd like. Building out an external scaler breaks this tight binding
- [ ] Run the proxy behind an ingress controller
- [ ] Make the proxy aware of multiple hosts
    - This is a start to making it multi-tenant
    - It should use the `Ingress` object as the lookup table for where to forward requests
- [ ] Build an releases/versioning API. It should do this:
    - Create a _new_ app: create a new `Deployment`, `Service`, `ScaledObject`, and rule for the new hostname in the `Ingress` rule for the proxy's ingress controller
    - Upgrade an existing app: update the image name on the app's `Deployment` and triggering a rollout
        - Future: figure out how to host multiple versions, do traffic splitting, etc... This probably involves service mesh support, see below
    - Delete an app: delete the app's `Deployment`, `Service`, and `ScaledObject`
- [ ] Hook the CLI up to the releases/versioning API
- [ ] Hook the CLI up to ACR tasks. It should be able to build a container in ACR, then trigger a deploy (either create a new app or upgrade one)

## Future Ideas (notes from @asw101 and @arschles discussion)

- [ ] Add admin "control plane" API to this, and a CLI for it
    - `csclr deploy hello-world:latest --platform=VMSS or --platform=ACI ...`
    - Also provide a standards-compliant YAML deployment (i.e. KEDA YAML or KNative `Service` YAML)
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

