# Filecoin state types
[![CircleCI](https://circleci.com/gh/filecoin-project/go-state-types.svg?style=svg)](https://circleci.com/gh/filecoin-project/go-state-types)
[![codecov](https://codecov.io/gh/filecoin-project/go-state-types/branch/master/graph/badge.svg)](https://codecov.io/gh/filecoin-project/go-state-types)

This repository contains primitive and low level types used in the Filecoin blockchain state representation.

These are primarily intended for use by [Actors](https://github.com/filecoin-project/specs-actors) and other
modules that read chain state directly.

## Versioning

We adopt a policy similar to the [Builtin-Actors repository](https://github.com/filecoin-project/builtin-actors?tab=readme-ov-file#versioning), with these key differences:
- The minor number in the version correlates with the `ActorVersion` from the Builtin-Actors repository.  (Builtin-Actors uses the major version number for this.)
- We don't use major versions; these are always set to 0. This is because of Go's special handling of versions to avoid the need to change import paths for every single package, including internal ones, which would result in having more than one version of go-state-types in the dependency tree.

Additional notes:

- We strive for round minor versions (e.g., 0.x.0) to denote the definitive release for a given network upgrade. However, due to the unpredictability of software engineering, further releases may be made by bumping the patch number (e.g., 0.x.1).
- Development versions use qualifiers like `-dev` (development) and `-rc` (release candidate).

As an example of application of this policy for Go-State-Types to a v14 builtin-actor version lineage:

- Unstable development versions are referenced by a `-dev` qualifier.
- Stable development versions are tagged as release candidates: `0.14.0-rc1`, `0.14.0-rc2`, etc.
- Definitive release: `0.14.0`.
- Patched releases: `0.14.1`, `0.14.2`.
- Network upgrade goes live with `0.14.2`.
- Patched releases can also occur after a network upgrade.

## Release Process

The repository contains a version.json file in the root directory:

```json
{
  "version": "v0.4.2"
}
```

This version file defines the currently released version.

To cut a new release, open a Pull Request that bumps the version number and have it reviewed by your teammates.

The release check workflow will create a draft GitHub Release (if it was not initiated by a PR from a fork) and post a link to it along with other useful information (the output of gorelease, gocompat, and a diff of the go.mod files(s)).

The releaser workflow runs when the PR is merged into the default branch. This workflow either publishes the draft GitHub Release created by the release check workflow or creates a published GitHub Release if it doesn't exist yet. This, in turn, creates a new Git tag and pushes it to the repository.

## License
This repository is dual-licensed under Apache 2.0 and MIT terms.

Copyright 2020. Protocol Labs, Inc.
