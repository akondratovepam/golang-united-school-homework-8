package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

const DEBUG bool = false

type Arguments map[string]string

func Perform(args Arguments, writer io.Writer) error {
	if DEBUG {
		fmt.Println("Arguments:", args)
	}

	operation := args["operation"]
	if operation == "" {
		return errors.New("-operation flag has to be specified")
	}

	filename := args["fileName"]
	if filename == "" {
		return errors.New("-fileName flag has to be specified")
	}

	switch operation {
	case "add":
		item := args["item"]
		if item == "" {
			return errors.New("-item flag has to be specified")
		}

		file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			return fmt.Errorf("failed to open file for adding: %w", err)
		}

		// todo read from file until the end
		// if id already exists, throw an error
		// add the item to the file end

		if err = file.Close(); err != nil {
			return fmt.Errorf("failed to close file after adding: %w", err)
		}
	case "list":
		bytes, err := os.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
		n, err := writer.Write(bytes)
		if err != nil {
			return fmt.Errorf("failed to write bytes to file: %w", err)
		}
		if DEBUG {
			fmt.Printf("Written bytes: %d\n", n)
		}
	case "findById":
		// todo
	case "remove":
		// todo
	default:
		return fmt.Errorf("Operation %s not allowed!", operation)
	}

	return nil
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

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
