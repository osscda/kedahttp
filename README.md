# Container Apps

Container Apps

This project implements a prototype of auto-scaling containers on ACI. As HTTP requests come into the system, the container(s) that are equipped to handle that request may or may not be running and ready to accept it. If there are sufficient containers available, the request is routed to one of them.  If there are not, a container is started and the request is routed to it when it's ready.

## Architecture

This system has three components:

- [CScaler Proxy](./cmd/proxy)
- [KEDA](https://keda.sh)

The **proxy** receives incoming HTTP traffic, emits events to NATS streaming, and forwards to a backend container.

KEDA is responsible for consuming events from the proxy and scaling the backend containers appropriately.

## Installation

Run `curl -L https://raw.githubusercontent.com/arschles/containerscaler/main/scripts/install.sh?token=AAYNMMCKX2PJV4T74AKS7XS7SB4RQ | bash`

## Manual Installation

To install the application you'll need a __Kubernetes Cluster__ up and running.

### Install KEDA

You need to install KEDA first. Do so with these commands:

```shell
helm repo add kedacore https://kedacore.github.io/charts
helm repo update
helm install keda kedacore/keda --namespace cscaler --create-namespace
```

>These commands are similar to those on the [official install page](https://keda.sh/docs/1.5/deploy/#helm), but we're installing in a different namespace.

### Install the Proxy

The proxy is responsible for receiving the requests, so you'll need to install it.

```shell
helm install cscaler ./charts/cscaler-proxy -n cscaler --create-namespace
```

To upgrade:

```shell
helm upgrade cscaler ./charts/cscaler-proxy -n cscaler
```

After the install, run the following command to fetch the public IP of the proxy service:

```shell
kubectl get svc cscaler-proxy -n cscaler -o=jsonpath='{.status.loadBalancer.ingress[*].ip}'
```

### Build the app

Just simply run ```make cli``` command within the root directory. This will create a new binary file called `capps` in the `bin` directory.

You can then install it into your ```PATH``` or add the ```./bin``` to your ```PATH```, or you can just run it by typing `./bin/capps` (assuming you're on the root).

## CLI API

```shell
./bin/capps
```

Running with no parameters will give you the general help for the commands

__Root Commands__:

- `help`: General help, use it as `./bin/capps help <cmd>` to get help on any command
- `rm`: Removes a created app, has its own set of flags
- `run`: Creates a new app, has its own set of flags
- `version`: Provides the version name and number

### Create an App

```shell
./bin/capps run <app-name> --image <repository>/<image>:<tag> --port <number> --server-url <url>
```

Runs a new application based on parameters.

__Flags__

- `-i`, `--image`: The image to be downloaded from the repository.
    > Since this command will create a new set of workloads, all the logged Docker repositories within the current cluster will work

- `-s`, `--server-url`: (__Required__) The URL for the admin server. To get this, run `kubectl get svc cscaler-proxy -n cscaler -o=jsonpath="{.status.loadBalancer.ingress[*].ip}"`
    > Without the correct admin url the scaler __will not__ work

- `-p`, `--port`: Port number to be exposed, should be the port where the app listens to incoming connections.

- `--use-http`: When set, the server URL will use the `HTTP` protocol instead of `HTTPS` (_default: false_)

### Remove an App

```shell
./bin/capps rm <app-name> --server-url <url>
```

Removes a previously created app

- `-s`, `--server-url`: (__Required__) The URL for the admin server. To get this, run `kubectl get svc cscaler-proxy -n cscaler -o=jsonpath="{.status.loadBalancer.ingress[*].ip}"`
    > Without the correct admin url the scaler __will not__ work

- `--use-http`: When set, the server URL will use the `HTTP` protocol instead of `HTTPS` (_default: false_)

## Access the app

Once deployed with `capps run` you'll be able to access the application through the __[proxy IP](#install-the-proxy)__.

However, the proxy only understands DNS hostnames, which means that, if your service is called `foo`, you'll have to access it through a DNS name like `foo.domain.com` and this DNS Zone needs to have an `A` record with the name `foo` pointing to the proxy IP. This is an implementation of an automatic ingress rule.

You can either use your own domain or an Azure provided one.

> __Important__: If you could not access the endpoint and it didn't work, probably you didn't have HTTPS enabled, try to use the `--use-http` flag to test again

### Access through your domain

1. Go to your DNS zone settings in your domain registrar
2. Add a new `A` record with the __same name as your service__ – it should point to `<service>.yourdomain.com`
3. Point this DNS record to the proxy IP
4. Give it some minutes or check [dnschecker.org](https://dnschecker.org) for the propagation
5. Access the domain

You can check the logs on `kubectl logs deploy/cscaler-proxy -f -n cscaler` to check for incoming requests

### Access through Azure-provided domains

1. If your cluster is not created in Azure, check the "Enable HTTP Application Routing Addon" box when creating it
2. If your cluster already exists, run `az aks enable-addons -n <cluster-name> -g <resource-group-name> --addons http_application_routing`
3. Once the execution is complete, open the Azure Portal, navigate to the resource group named `MC_<group-name>_<cluster-name>_<location>`
4. Find the Azure DNS zone the Addon created for you
5. Follow steps 2 to 4 from the [previous section](#access-through-your-domain)
6. Access the service using `<service-name>.<dns-zone-name>`

You can check the logs on `kubectl logs deploy/cscaler-proxy -f -n cscaler` to check for incoming requests

## Debugging

If you are using [vscode ](https://code.visualstudio.com/) you can open up the dev environment in container using the [remote container extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) the nats server will be avalible in the dev container via `nats-server:4222`

If you need to do any DNS work from inside a container that's running Alpine linux, use this command:

```shell
curl -L https://github.com/sequenceiq/docker-alpine-dig/releases/download/v9.10.2/dig.tgz|tar -xzv -C /usr/local/bin/
```

Courtesy https://github.com/sequenceiq/docker-alpine-dig

## More Information

See [this document](./docs/COMPONENTS.md) for details on the components of this system.
