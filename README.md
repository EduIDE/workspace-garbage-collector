# Theia Workspace Garbage Collector
This operator is responsible for cleaning up the workspace of theia-ide pods.
It is a Kubernetes operator that cleans unused workspaces periodically.

## Configuration
The following values can be configured via environment variables:
- `WORKSPACE_TTL`: The time to live of a workspace. Default is 2 weeks (1209600s).
- `CHECK_INTERVAL`: The interval at which the operator checks for unused workspaces. Default is 1 hour (1800s).

## Deployment
`helm upgrade --install theia-workspace-garbage-collector ./helm -f ./helm/values.yaml`
## Releasing
Releases are manual: a release pull request (version in `helm/Chart.yaml`, both `version` and `appVersion`, plus the `CHANGELOG.md` section), then
`gh workflow run release.yml -R EduIDE/workspace-garbage-collector -f version=X.Y.Z` after it merges.
The run tags `theia-workspace-garbage-collector-X.Y.Z`, creates the GitHub Release from the CHANGELOG, pushes the image `ghcr.io/eduide/garbage-collector:X.Y.Z` and the chart `oci://ghcr.io/eduide/charts/theia-workspace-garbage-collector:X.Y.Z`, and verifies both.
Nothing releases on push to `main`. The chart's `image.tag` defaults to its `appVersion`, so a chart version pins its own image.
See `.claude/skills/release/SKILL.md`.
