# Keda improvements

## Docs for external scaling

- Example with real metrics in the [external scaling docs overview](https://keda.sh/docs/1.5/concepts/external-scalers/#overview) - plug it in to the [HPA formula](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/#algorithm-details)
- Mention how `maxReplicaCount`, `minReplicaCount`, `cooldownPeriod`, and `pollingInterval` factor in to external scaling
- Include how to start the Go gRPC server on the [scaling docs](https://keda.sh/docs/1.5/concepts/external-scalers/#overview)