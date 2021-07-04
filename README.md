# partial

This library provides a debug-ready and straightforward utility for updating some of the fields of a struct.

# Background

When implementing a partial update of a field using an ORM, we need to write many if-statements to determine if there is an update on each field. Since this implementation often becomes too redundant, I'd like to provide an easy way to implement it.

There are several libraries to solve this problem, but this library focuses on making debugging easier. This library provides the API to know which fields were updated, skipped, and failed.

# Features

- Copy fields between different type structs
- Supported fields copy patterns
  - src: pointer -> dest: pointer
  - src: pointer -> dest: non-pointer
  - src: non-pointer -> dest: non-pointer

# API

See [example/main.go](example/main.go)
