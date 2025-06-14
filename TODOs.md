# TODOs

## Refactoring Plan

The following tasks outline the planned refactoring to split the code into focused components and improve maintainability:

- [ ] Move the definitions of `Updater`, `Field`, and error variables to `types.go`.
- [ ] Extract the `NewUpdater` constructor and initialization logic to `constructor.go`.
- [ ] Extract the skip‐field logic (`ShouldSkipField`) to `filter.go`.
- [ ] Extract the field‐finder logic (`FindFieldByInsensitiveName`) to `finder.go`.
- [ ] Move the core `Updater.Update` method and its internal helpers to `update.go`.

### Future Enhancements

- [ ] Add per-component unit tests targeting each new file.
- [ ] Consider using Go generics (Go 1.18+) or code generation to reduce reflection overhead.
- [ ] Improve documentation of the public API in `README.md`.