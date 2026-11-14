# Theia Workspace Garbage Collector
This operator is responsible for cleaning up the workspace of theia-ide pods.
It is a Kubernetes operator that cleans unused workspaces periodically.

## Configuration
The following values can be configured via environment variables:
- `WORKSPACE_TTL`: The time to live of a workspace. Default is 2 weeks (1209600s).
- `CHECK_INTERVAL`: The interval at which the operator checks for unused workspaces. Default is 1 hour (1800s).

## Deployment
Set `global.nodeSelector` in the Helm values to pin the pod to a node pool. The parent `eduide` chart sets it for you.

`helm upgrade --install theia-workspace-garbage-collector ./helm -f ./helm/values.yaml`