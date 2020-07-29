#!/usr/bin/env bash
set -euo pipefail

# variables
RESOURCE_GROUP='200200-actions'
LOCATION='eastus'
RANDOM_STR='09c205'
# [[ -z "$RANDOM_STR" ]] && RANDOM_STR=$(openssl rand -hex 3)
REPOSITORY_NAME="hello-golang"

# create container registry
REGISTRY_NAME="acr${RANDOM_STR}"
az acr create -g $RESOURCE_GROUP -l $LOCATION --name $REGISTRY_NAME --sku Basic --admin-enabled true
# build image
CONTAINER_IMAGE=$REPOSITORY_NAME:$(date +%y%m%d)-${GITHUB_SHA}
az acr build -r $REGISTRY_NAME -t $CONTAINER_IMAGE --file Dockerfile .
# create container instance
REGISTRY_PASSWORD=$(az acr credential show -n $REGISTRY_NAME | jq -r .passwords[0].value)
CONTAINER_NAME="aci-${REPOSITORY_NAME}-${RANDOM_STR}"
az container create --resource-group $RESOURCE_GROUP --location $LOCATION \
    --name $CONTAINER_NAME \
    --image "${REGISTRY_NAME}.azurecr.io/${CONTAINER_IMAGE}" \
    --registry-login-server "${REGISTRY_NAME}.azurecr.io" \
    --registry-username $REGISTRY_NAME \
    --registry-password $REGISTRY_PASSWORD \
    --cpu 1 \
    --memory 1 \
    --ports 80 \
    --environment-variables LISTEN_PORT=80 \
    --dns-name-label ${REPOSITORY_NAME}-${RANDOM_STR}
FQDN=$(az container show -g $RESOURCE_GROUP --name $CONTAINER_NAME | jq -r .ipAddress.fqdn)
echo "http://${FQDN}"
