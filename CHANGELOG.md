# Changelog

## [Unreleased]

### Features

- Release tags follow the org-wide `vX.Y.Z` format, like every other EduIDE repository. Releases up to 0.2.0 keep their old `theia-workspace-garbage-collector-X.Y.Z` tags. Chart and image versions stay `X.Y.Z`.
- Every release now publishes the image as `ghcr.io/eduide/garbage-collector:X.Y.Z`. The chart's `image.tag` defaults to its `appVersion`, so a chart version pins its own image instead of `latest` ([#25](https://github.com/EduIDE/workspace-garbage-collector/pull/25)).
