# The Lifecycle of a Request

## The Proxy

The proxy has three purposes:

- Stores a list of _backend containers_ that can serve a given request
- Forwards a given request to a backend container
- Requests a new container to serve a given request, waits for it, and stores it in the database

A new request into the system will first hit the proxy. At this point, the proxy first emits a `reqcounter` event (the controller listens for this event).

Next, it looks into its internal database (it uses [BoltDB](https://github.com/boltdb/bolt)) for a container that can serve that request. If it finds one or more, it routes the request to a random one and finishes.

If not, it waits for a `scaledup` event from the controller indicating that a new container has been created to serve the event, forwards the request to it, saves it to its database, and finishes.

## The Controller

The controller 
