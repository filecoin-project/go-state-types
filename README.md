# Filecoin state types
[![CircleCI](https://circleci.com/gh/filecoin-project/go-state-types.svg?style=svg)](https://circleci.com/gh/filecoin-project/go-state-types)
[![codecov](https://codecov.io/gh/filecoin-project/go-state-types/branch/master/graph/badge.svg)](https://codecov.io/gh/filecoin-project/go-state-types)

This repository contains primitive and low level types used in the Filecoin blockchain state representation.

These are primarily intended for use by [Actors](https://github.com/filecoin-project/specs-actors) and other
modules that read chain state directly.

## Versioning

We adopt a policy similar to the [Builtin-Actors repository](https://github.com/filecoin-project/builtin-actors?tab=readme-ov-file#versioning), with a key difference:

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

<details>
  <summary>Cutting a development or release candidate release:</summary>

1. Go to [Go-State-Types Releases](https://github.com/filecoin-project/go-state-types/releases).
2. Click the "Draft a new release" button in the right corner.
3. In the "Choose a tag" dropdown, enter the desired version and click "Create new tag: vX.XX.X on publish".
4. Target the master branch.
5. Set the previous tag to compare against, the last stable release, and click the "Generate release notes" button.
6. Check the "Set as a pre-release" checkbox.
7. Click "Publish release" to create the development or release candidate release.

</details>

<details>
  <summary>Cutting a definitive release:</summary>

1. Go to [Go-State-Types Releases](https://github.com/filecoin-project/go-state-types/releases).
2. Click the "Draft a new release" button in the right corner.
3. In the "Choose a tag" dropdown, enter the desired version and click "Create new tag: vX.XX.X on publish".
4. Target the master branch.
5. Set the previous tag to compare against, the last stable release, and click the "Generate release notes" button.
6. Ensure the "Set as a pre-release" checkbox is **not** checked.
7. Click "Publish release" to create the definitive release.

</details>

## License
This repository is dual-licensed under Apache 2.0 and MIT terms.

Copyright 2020. Protocol Labs, Inc.
