package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}

type Arguments map[string]string
type User struct {
	Id    string
	Email string
	Age   uint8
}

func parseArgs() Arguments {
	operationFlag := flag.String("operation", "", "<add|list|findById|remove>")
	itemFlag := flag.String("item", "", "JSON e.g. {\"id\": \"1\", \"email\": «email@test.com», «age»: 23}")
	fileNameFlag := flag.String("fileName", "", "The JSON file name with users data.")
	flag.Parse()

	return Arguments{
		"operation": *operationFlag,
		"item":      *itemFlag,
		"fileName":  *fileNameFlag}
}

func Perform(args Arguments, writer io.Writer) error {
	operation := args["operation"]
	if operation == "" {
		return errors.New("-operation flag has to be specified")
	}

	fileName := args["fileName"]
	if fileName == "" {
		return errors.New("-fileName flag has to be specified")
	}

	switch operation {
	case "add":
		item := args["item"]
		if item == "" {
			return errors.New("-item flag has to be specified")
		}

		return add(item, fileName, writer)
	case "list":
		return list(fileName, writer)
	case "findById":
		id := args["id"]
		if id == "" {
			return errors.New("-id flag has to be specified")
		}

		return findById(id, fileName, writer)
	case "remove":
		id := args["id"]
		if id == "" {
			return errors.New("-id flag has to be specified")
		}

		return remove(id, fileName, writer)
	default:
		return fmt.Errorf("Operation %s not allowed!", operation)
	}
}

func add(item, fileName string, writer io.Writer) error {
	// unmarshal the new user
	var newUser User
	if err := json.Unmarshal([]byte(item), &newUser); err != nil {
		return fmt.Errorf("failed to unmarshal the new user JSON: %w", err)
	}

	// read users from file
	var users []User
	bytes, _ := os.ReadFile(fileName)
	if err := json.Unmarshal(bytes, &users); err != nil {
		users = make([]User, 0, 1)
	} else {
		for _, user := range users {
			if user.Id == newUser.Id {
				writer.Write([]byte("Item with id " + user.Id + " already exists"))
				return nil
			}
		}
	}

	// append the new user and save the list
	users = append(users, newUser)
	if err := saveUsers(users, fileName); err != nil {
		return fmt.Errorf("failed to save users: %w", err)
	}

	return nil
}

func list(fileName string, writer io.Writer) error {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	_, err = writer.Write(bytes)
	if err != nil {
		return fmt.Errorf("failed to write bytes to file: %w", err)
	}

	return nil
}

func findById(id, fileName string, writer io.Writer) error {
	var users []User
	bytes, _ := os.ReadFile(fileName)
	if err := json.Unmarshal(bytes, &users); err == nil {
		for _, user := range users {
			if user.Id == id {
				bytes, err := json.Marshal(user)
				if err != nil {
					return fmt.Errorf("failed to marshal user as JSON: %w", err)
				}

				writer.Write(bytes)
			}
		}
	}

	return nil
}

func remove(id, fileName string, writer io.Writer) error {
	var users []User
	bytes, _ := os.ReadFile(fileName)
	if err := json.Unmarshal(bytes, &users); err != nil {
		// do nothing if db empty
		return nil
	}

	newUsers := make([]User, 0, len(users)-1)
	for _, user := range users {
		if user.Id != id {
			newUsers = append(newUsers, user)
		}
	}

	if len(newUsers) == len(users) {
		writer.Write([]byte("Item with id " + id + " not found"))
		return nil
	}

	if err := saveUsers(users, fileName); err != nil {
		return fmt.Errorf("failed to save users: %w", err)
	}

	return nil
}

func saveUsers(users []User, fileName string) error {
	bytes, err := json.Marshal(users)
	if err != nil {
		return fmt.Errorf("failed to marshal users as JSON: %w", err)
	}

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("failed to open file for adding: %w", err)
	}

	if _, err = file.Write(bytes); err != nil {
		return fmt.Errorf("failed to write data to file: %w", err)
	}

	if err = file.Close(); err != nil {
		return fmt.Errorf("failed to close file after adding: %w", err)
	}

	return nil
}
