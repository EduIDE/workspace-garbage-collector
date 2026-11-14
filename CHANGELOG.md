# Changelog

## [Unreleased]

### Features

- Every release now publishes the image as `ghcr.io/eduide/garbage-collector:X.Y.Z`. The chart's `image.tag` defaults to its `appVersion`, so a chart version pins its own image instead of `latest` ([#25](https://github.com/EduIDE/workspace-garbage-collector/pull/25)).
