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
	Name             string
	Age              *int
	Address          *string
	Birthday         *time.Time
	NotExistingField string
}

func main() {
	user := User{
		Name:     "John Doe",
		Age:      30,
		Address:  nil,
		Birthday: time.Time{},
	}

	age := user.Age + 1
	address := "123 Main Street"
	now := time.Now()

	userUpdate := UserUpdate{
		Name:     "Updated Name",
		Age:      &age,
		Address:  &address,
		Birthday: &now,
	}

	updater, err := partial.NewUpdater(&user)
	if err != nil {
		panic(err)
	}

	err = updater.Update(&userUpdate)

	// When the NotFoundFields and NotAssignableFields are not empty, it will return an error.
	if err != nil {
		fmt.Printf("Error: %v\n", err) // Error: update fields failure
	}

	fmt.Printf("Updated user: %+v\n", user)
	// => Updated user: {Name:Updated Name Age:31 Address:0xc000010230 Birthday:2021-07-04 09:16:12.458143 +0900 JST m=+0.000085842}
	fmt.Printf("updated fields: %d\n", len(updater.UpdatedFields))
	// => updated fields: 3
	fmt.Printf("skipped fields: %d\n", len(updater.SkippedFields))
	// => skipped fields: 1
	fmt.Printf("not found fields: %d\n", len(updater.NotFoundFields))
	// => not found fields: 1
	fmt.Printf("not assignable fields: %d\n", len(updater.NotAssignableFields))
	// => not assignable fields: 0
}
```
