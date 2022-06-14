package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type person struct {
	Id    string
	Email string
	Age   int
}

type Arguments map[string]string

func Perform(args Arguments, writer io.Writer) error {
	if args["fileName"] == "" {
		return fmt.Errorf("-fileName flag has to be specified")
	}

	operation := args["operation"]
	switch operation {
	case "add":
		return addOperation(args, writer)
	case "list":
		return listOperation(args, writer)
	case "remove":
		return removeOperation(args, writer)
	case "findById":
		return findByIdOperation(args, writer)
	case "":
		return fmt.Errorf("-operation flag has to be specified")
	default:
		return fmt.Errorf("Operation %s not allowed!", operation)
	}

	return nil
}

func addOperation(args Arguments, writer io.Writer) error {
	if args["item"] == "" {
		return fmt.Errorf("-item flag has to be specified")
	}

	people, err := getPeople(args["fileName"])
	if err != nil {
		return err
	}

	item, err := getPerson(args["item"])
	if err != nil {
		return err
	}

	if itemExists(people, item.Id) {
		_, err = writer.Write([]byte(fmt.Sprintf("Item with id %s already exists", item.Id)))
		return err
	}

	people = append(people, item)
	final, err := json.Marshal(people)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	err = writeToFile(args["fileName"], final)
	return err
}

func itemExists(people []person, id string) bool {
	for _, x := range people {
		if x.Id == id {
			return true
		}
	}
	return false
}

func listOperation(args Arguments, writer io.Writer) error {
	data, err := os.ReadFile(args["fileName"])
	if err != nil {
		return err
	}

	_, err = writer.Write(data)
	return err
}

func findByIdOperation(args Arguments, writer io.Writer) error {
	if args["id"] == "" {
		return fmt.Errorf("-id flag has to be specified")
	}

	people, err := getPeople(args["fileName"])
	if err != nil {
		return err
	}

	for _, x := range people {
		if x.Id == args["id"] {
			person, err := json.Marshal(x)
			if err != nil {
				return err
			}
			_, err = writer.Write(person)
			return nil
		}
	}

	return err
}

func removeOperation(args Arguments, writer io.Writer) error {
	if args["id"] == "" {
		return fmt.Errorf("-id flag has to be specified")
	}

	people, err := getPeople(args["fileName"])
	if err != nil {
		return err
	}

	for i, x := range people {
		if x.Id == args["Id"] {
			temp := append(people[:i], people[i+1:]...)
			final, err := json.Marshal(temp)
			if err != nil {
				return err
			}
			err = writeToFile(args["fileName"], final)
			return err
		}
	}

	_, err = writer.Write([]byte(fmt.Sprintf("Item with id %s not found", args["id"])))
	return err
}

func getPeople(fileName string) ([]person, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var people []person
	err = json.Unmarshal(data, &people)
	if err != nil {
		return nil, err
	}

	return people, nil
}

func getPerson(item string) (person, error) {
	var p person
	err := json.Unmarshal([]byte(item), &p)
	if err != nil {
		return person{}, err
	}

	return p, nil
}

func writeToFile(fileName string, data []byte) error {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	file.Truncate(0)
	_, err = file.Write(data)
	file.Close()
	return err
}

func main() {
	//idFlag := flag.String("Id", "", "Item Id")
	//operationFlag := flag.String("operation", "", "Operation Name")
	//filenameFlag := flag.String("filename", "", "Filename")
	//itemFlag := flag.String("item", "", "items")
	//fmt.Println(*filenameFlag)
	//args := Arguments{
	//	"Id":        *idFlag,
	//	"operation": *operationFlag,
	//	"item":      *itemFlag,
	//	"fileName":  *filenameFlag,
	//}

	args := Arguments{
		"Id":        "",
		"operation": "list",
		"item":      "",
		"fileName":  "test.json",
	}

	err := Perform(args, os.Stdout)
	if err != nil {
		panic(err)
	}
}
