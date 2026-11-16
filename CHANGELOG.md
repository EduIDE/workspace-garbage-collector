# Changelog

## [Unreleased]

## [1.0.0] - 2026-11-16

### Features

- Release tags follow the org-wide `vX.Y.Z` format, like every other EduIDE repository. Releases up to 0.2.0 keep their old `theia-workspace-garbage-collector-X.Y.Z` tags. Chart and image versions stay `X.Y.Z`.
- Every release now publishes the image as `ghcr.io/eduide/garbage-collector:X.Y.Z`. The chart's `image.tag` defaults to its `appVersion`, so a chart version pins its own image instead of `latest` ([#25](https://github.com/EduIDE/workspace-garbage-collector/pull/25)).

### Fixes

- The garbage collector no longer crashes on start when `CHECK_INTERVAL` is unset, and `WORKSPACE_TTL` no longer fails every run when unset: both fall back to the documented defaults (30 minutes and 2 weeks). Values without a unit are read as seconds, and zero, negative or malformed values stop the start with a clear error. ([#27](https://github.com/EduIDE/workspace-garbage-collector/pull/27))
