#!/bin/bash

function shoutln {
  echo
  printf '%s\n' "$1" | awk '{ print toupper($0) }'
  echo
}

function shout {
  printf '%s' "$1" | awk '{ print toupper($0) }'
}

function log {
  echo -n "-> $1"
}

function logln () {
  echo "-> $1"
}

commandExists () {
    type "$1" &> /dev/null ;
}

function pause () {
  read -s -n 1 -p "$*"
}

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
REPOSITORY_LOCATION=/tmp/containerscaler
rm -rf $REPOSITORY_LOCATION

shoutln "=== container scaler init script ==="
echo
shoutln "Checking for prerequisites"

log "Checking for Helm... "
if ! commandExists helm ; then
  echo
  echo "  Helm is not installed. Installing it now..."
  curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash
fi
echo "done"

log "Checking for the Azure CLI... "
if ! commandExists az ; then
  echo
  echo "  Azure CLI is not installed. Installing it now..."
  curl -L https://aka.ms/InstallAzureCli | bash
  echo "  I'll prompt you to log in to your Azure Account, please follow the next steps"
  az login
fi
echo "done"

logln "This script will create all resources under your current default subscription in Azure CLI, if that's not what you want, please stop it and use \"az account set --subscription {your subscription id}\" to set a new subscription and then run this command again."

pause "  If that's the correct subscription press [Enter] to continue"

log "Checking for the Kubectl... "
if ! commandExists kubectl ; then
  echo
  echo "  Kubectl is not installed. Installing it now..."
  az aks install-cli
fi
echo "done"

log "Checking for Git... "
if ! commandExists git ; then
  echo
  echo "  Git is not installed. Please refer to https://git-scm.com/book/en/v2/Getting-Started-Installing-Git to install it and run this command again"
  exit 1
fi
echo "done"

shoutln "=== creating kubernetes cluster ==="
while read -p "==> What is the RESOURCE GROUP name I should use to create this cluster? " RESOURCE_GROUP_NAME && [[ -z "$RESOURCE_GROUP_NAME" ]]
do
  echo "  Resource group name is required"
done

while read -p "==> What is the LOCATION I should use to create this cluster? " LOCATION && [[ -z "$LOCATION" ]]
do
  echo "  Location is required"
done

resource_group_exists=$(az group list --query "[?name =='$RESOURCE_GROUP_NAME'].name" -o tsv)
if [ -n "$resource_group_exists" ]
then
  logln "Resource group with this name already exists, using it..."
else
  az group create -n "$RESOURCE_GROUP_NAME" --location "$LOCATION"
fi

while read -p "==> What will be the name of the CLUSTER? " AKS_NAME && [[ -z "$AKS_NAME" ]]
do
  echo "  Name is required"
done

logln "Creating Azure Resource"
echo "  Be patient, this can take a while..."
az aks create -g "$RESOURCE_GROUP_NAME" -n "$AKS_NAME" --node-count=2 --generate-ssh-keys --node-vm-size=Standard_B2s --location "$LOCATION" --enable-addons=http_application_routing
[ $? -ne 0 ] && echo "ERROR: Error creating Azure Resources, stopping script" && exit 2

echo "  Cluster install finished!"

logln "Getting credentials"
az aks get-credentials -g $RESOURCE_GROUP_NAME -n $AKS_NAME
kubectl config set-context $AKS_NAME

shoutln "=== adding helm repos ==="
helm repo add kedacore https://kedacore.github.io/charts
helm repo update

shoutln "=== cloning repository ==="
git clone https://github.com/arschles/containerscaler "$REPOSITORY_LOCATION"

shoutln "=== installing cscaler ==="

read -p "==> What will be the name of the NAMESPACE where I should install everything? [cscaler]" NAMESPACE
NAMESPACE=${NAMESPACE:-cscaler}

logln "Installing KEDA"
helm install keda kedacore/keda --namespace "$NAMESPACE" --create-namespace

logln "Installing Proxy"

AKS_HAR_ZONE_NAME=$(az aks show -g "$RESOURCE_GROUP_NAME" -n "$AKS_NAME" -o tsv --query addonProfiles.httpApplicationRouting.config.HTTPApplicationRoutingZoneName)
CONFIG_LOCATION=$HOME/.capps

helm install cscaler-proxy $REPOSITORY_LOCATION/charts/cscaler-proxy -n "$NAMESPACE" --create-namespace --set cscalerProxyDNSZoneName="$AKS_HAR_ZONE_NAME"

logln "Creating config file..."
mkdir -p $CONFIG_LOCATION
echo "server_url: cscaler-admin.$AKS_HAR_ZONE_NAME" > $CONFIG_LOCATION/cappsconfig
echo " Config file written to '$CONFIG_LOCATION'"

shoutln "=== compiling binary ==="

make -C $REPOSITORY_LOCATION cli

if [[ ":$PATH:" == *":/usr/local/bin:"* ]]; then
  [ ! -d "/usr/local/bin" ] && mkdir -p /usr/local/bin
  mv $REPOSITORY_LOCATION/bin/capps /usr/local/bin
  rm -rf $REPOSITORY_LOCATION
else
  logln "You don't seem to have \"/usr/local/bin\" in your PATH variable, you might want to add that."
  logln "The CLI is built, please copy it from \"$REPOSITORY_LOCATION/bin/capps\" to where you'd like to use it"
fi

shout "===! finished !==="
logln "Use 'capps' to access the CLI"
exit 0
