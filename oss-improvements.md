# Keda improvements

## Docs for external scaling

- Example with real metrics in the [external scaling docs overview](https://keda.sh/docs/1.5/concepts/external-scalers/#overview) - plug it in to the [HPA formula](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/#algorithm-details)
- Mention how `maxReplicaCount`, `minReplicaCount`, `cooldownPeriod`, and `pollingInterval` factor in to external scaling
- Include how to start the Go gRPC server on the [scaling docs](https://keda.sh/docs/1.5/concepts/external-scalers/#overview)
- KEDA redis issue (https://github.com/kedacore/keda/issues/905)


## Helm

- How to configure dependency configs (same as subchart configs)
- You can install helm via linuxbrew

This on the getting started page:

```shell
DESKTOP-DQP07VM :: ~/src/containerscaler ‹controller*› » helm install stable/mysql --generate-name
Error: failed to download "stable/mysql" (hint: running `helm repo update` may help)
```
