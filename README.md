# partial

This library provides a simple and debug-ready utility for updating some of the fields of a struct.

# Background

When implementing a partial update of a field using an RDBMS ORM, the Go language requires an if statement to determine if there has been an update for each of the fields. Since this implementation is often very complicated, we would like to provide an easy way to implement it.

There are several existing libraries, but my library is focused on making debugging easier.

Updater provides an API to know which fields were successfully updated, which fields were skipped, and which fields failed to be updated.

# Features

- Copy fields between different type structs
- Supported fields copy patterns
  - src: pointer -> dest: pointer
  - src: pointer -> dest: non-pointer
  - src: non-pointer -> dest: non-pointer
