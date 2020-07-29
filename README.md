# Autoscaling Containers

This project implements a prototype of auto-scaling containers on either Kubernetes or ACI. Although the featureset is very basic, this is similar to [Knative serving](https://knative.dev/docs/serving/), with the following exceptions:

- Simpler to install
    - A HTTP proxy and a scaling component, no service mesh required
- No routes or versions
- Single tenant (at the moment)

## Architecture

This system has two components:

The **proxy** receives incoming HTTP traffic, looks up where to send that traffic in its database, and forwards it on to a URL. This URL can be any DNS name or IP. On Kubernetes, this URL should be a `Service` DNS name. For forwarding to ACI, this URL should be either a public IP or DNS name of one or more ACI containers.

The **scaling controller** is responsible for fetching traffic metrics from the proxy and scaling containers up and down based on request volume.

## How to Run This

The proxy and controller both depend on NATS, so that's the first thing to run. Do so with Docker and this command:

```shell
docker run -p 4222:4222 -ti nats:latest
```

## FAQ

_Why don't you use Horizontal Pod Autoscaling, ingress controllers, or service meshes to do these things?_

Because those systems don't work for independently running processes that need to forward to ACI containers. They are generally for Kubernetes or other container orchestrators only.