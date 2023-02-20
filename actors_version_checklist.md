### Actor version integration checklist

- [ ] Copy `go-state-types/builtin/vX` to `go-state-types/builtin/v(X+1)`
- [ ] Change all references to `vX` in the new files to `v(X+1)`
- [ ] Add new actors version to `actors/version.go`
- [ ] Add the new version to the `gen` step of the makefile
- [ ] run `make gen`
- 