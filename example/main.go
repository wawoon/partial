package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/wawoon/partial"
)

type User struct {
	Name         string
	Age          int
	Address      *string
	Birthday     time.Time
	RegisterDate time.Time
}

type UserUpdate struct {
	Name             string
	Age              *int
	Address          *string
	Birthday         *time.Time
	RegisterDate     string
	NotExistingField string
}

func prettyJson(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func mapKeys(v map[string]partial.Field) []string {
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	return keys
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
		Name:         "Updated Name",
		Age:          &age,
		Address:      &address,
		Birthday:     &now,
		RegisterDate: "2016-01-01",
	}

	updater, err := partial.NewUpdater(&user)
	if err != nil {
		panic(err)
	}

	err = updater.Update(userUpdate)

	// When the NotFoundFields and NotAssignableFields are not empty, it will return an error.
	if err != nil {
		fmt.Printf("Error: %v\n", err) // Error: update fields failure
	}

	fmt.Printf("Updated user:\n%s\n", prettyJson(user))
	// Updated user:
	// {
	//   "Name": "Updated Name",
	//   "Age": 31,
	//   "Address": "123 Main Street",
	//   "Birthday": "2021-07-04T09:33:00.410106+09:00",
	//   "RegisterDate": "0001-01-01T00:00:00Z"
	// }
	fmt.Printf("updated fields: %+v\n", mapKeys(updater.UpdatedFields))
	// updated fields: [Birthday Name Age Address]
	fmt.Printf("skipped fields: %+v\n", mapKeys(updater.SkippedFields))
	// skipped fields: []
	fmt.Printf("not found fields: %+v\n", mapKeys(updater.NotFoundFields))
	// not found fields: [NotExistingField]
	fmt.Printf("not assignable fields: %+v\n", mapKeys(updater.NotAssignableFields))
	// not assignable fields: [RegisterDate]
}
