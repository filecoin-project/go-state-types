### Actor version integration checklist

- [ ] Copy `go-state-types/builtin/vX` to `go-state-types/builtin/v(X+1)`
- [ ] Change all references to `vX` in the new files to `v(X+1)`
- [ ] Add new network version to `network/version.go`
- [ ] Add new actors version to `actors/version.go`[^1]
- [ ] Add the new version to the `gen` step of the makefile`[^2]
- [ ] run `make gen`

[^1]: 
    #### Steps:
    1. **Add a new constant**: Add a new constant in the list of versions. The new constant's name should follow the existing naming convention - i.e., `VersionXX+1  Version = XX+1`, where XX+1 is the new version number.
    2. **Update `VersionForNetwork` function**: In `version.go`, there's a function called `VersionForNetwork` that accepts a network version and returns the corresponding actor version. Add a new case line for the network version that corresponds to the new actor version you're adding - i.e, `network.Version(XX+1): return Version(XX+1), nil`
[^2]:  Add `$(GO_BIN) run ./builtin/v(XX+1)/gen/gen.go`
