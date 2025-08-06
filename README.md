# partial

Tiny, reflection-powered helper to copy only the **explicitly-set** fields from one struct to another – ideal for building PATCH/partial-update style APIs while retaining full observability of what changed.

*🪶 Minimal* • *📦 Zero dependencies* • *🔍 Debug-friendly*

## Table of Contents

* [Requirements](#requirements)
* [Installation](#installation)
* [Quick example](#quick-example)
* [The Updater API & error handling](#the-updater-api--error-handling)
* [Supported field-copy patterns](#supported-field-copy-patterns)
* [Why choose partial?](#why-choose-partial)
* [Contributing](#contributing)
* [Roadmap](#roadmap)
* [License](#license)

---

## Requirements

* **Go 1.18+** – `partial` makes use of Go generics introduced in 1.18.
* **Platform-independent** – implemented in pure Go (`cgo`-free) and therefore works on any operating system and architecture supported by the Go tool-chain.

---

## Installation

Add the library to your module dependencies at the latest tagged version:

```bash
go get github.com/wawoon/partial@latest
```

If you prefer the classic *install a binary* workflow (for example, when a CLI is introduced in the future), you can run:

```bash
go install github.com/wawoon/partial/cmd/partial@latest
```

> Note: `partial` currently ships as a library only; the command above is listed for completeness and will simply fail until the CLI lands.

---

## Quick example

```go
package main

import (
    "fmt"
    "time"

    "github.com/wawoon/partial"
)

type User struct {
    Name     string
    Age      int
    Address  *string
    Birthday time.Time
}

type UserUpdate struct {
    Name     string     // copied when not empty string
    Age      *int       // copied when pointer is not nil
    Birthday *time.Time // copied when pointer is not nil
}

func main() {
    u := User{Name: "John", Age: 30}

    age := 31
    now := time.Now()
    upd := UserUpdate{Age: &age, Birthday: &now}

    updater, _ := partial.NewUpdater(&u)
    if err := updater.Update(upd); err != nil {
        panic(err)
    }

    fmt.Printf("updated user: %+v\n", u)
    fmt.Printf("fields touched: %v\n", updater.Updated())
}
```

See [`example/main.go`](example/main.go) for a more complete demonstration.

---

## The Updater API & error handling

`partial.NewUpdater(dst any)` returns an `*partial.Updater` bound to the destination struct. The updater internally tracks four maps keyed by field name:

| Map                    | Meaning                                                                       |
|------------------------|--------------------------------------------------------------------------------|
| `UpdatedFields`        | Field existed on destination and value was copied                              |
| `SkippedFields`        | Field existed but source value was considered *zero* (empty string, nil ptr…)  |
| `NotFoundFields`       | Field does **not** exist on destination struct                                 |
| `NotAssignableFields`  | Field exists but the value types are incompatible                              |

If `Update` encounters any field in `NotFoundFields` **or** `NotAssignableFields`, it returns the sentinel error `partial.ErrUpdateFailed`. You can still inspect the maps to know what happened and decide whether to ignore or handle the situation.

```go
if err := updater.Update(req); errors.Is(err, partial.ErrUpdateFailed) {
    log.Printf("unexpected fields: %v", updater.NotFoundFields)
}
```

---

## Supported field-copy patterns

| Source value | Destination field | Supported | Notes                          |
|--------------|-------------------|-----------|--------------------------------|
| pointer      | pointer           | ✅        | src != nil                     |
| pointer      | value             | ✅        | src != nil                     |
| value        | value             | ✅        | src != zero value              |
| value        | pointer           | ❌        | not supported (would allocate) |

---

## Why choose `partial`?

• **Transparent debugging** – know exactly which fields were updated, skipped or failed.

• **Tiny surface area** – just two public functions and a single struct; easy to audit.

• **No magic** – relies on Go `reflect` only, without additional tags, build steps or code generation.

• **Idiomatic** – zero external dependencies and follows Go conventions.

---

## Contributing

1. Fork the repository and create a feature branch (`git checkout -b feat/awesome-change`).
2. Commit using the [Conventional Commits](https://www.conventionalcommits.org/) specification.
3. Ensure the test-suite passes locally:

   ```bash
   go test ./...
   ```

4. Open a Pull Request describing your motivation and the solution.

We love well-tested, well-documented contributions 💙.

---

## Roadmap

* Support copying **value → pointer** fields (allocate destination value when needed).
* Custom zero-value detection hooks/interfaces.
* Helper generic function that removes the need to instantiate an `Updater` explicitly.
* Allocation/performance optimisations in hot-paths.

---

## License

Distributed under the MIT License. See [`LICENSE`](../LICENSE) for the full text.
