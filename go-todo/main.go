package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
) 

type Todo struct {
    Id  int
    Description string
    CreateTime time.Time
    IsDone bool
}

var globalList []Todo

func printItems(items []Todo) {
    fmt.Println("Items:")
    for _, item := range items {
        fmt.Printf("Id %d, desc: %s\n", item.Id, item.Description)
    }
}
var filename = "items.json"

func main() {


    if _, err := os.Stat(filename); err == nil {
        jsonData, err := os.ReadFile(filename)
        if err != nil {
            fmt.Println("Error reading file:", err)
            return 
        }

        if err := json.Unmarshal(jsonData, &globalList); err != nil {
            fmt.Println("Error unmarshaling JSON:", err)
            return 
        }

        fmt.Println("Items loaded from", filename)
        // printItems(globalList)
    } else {

        fmt.Println("next")
        items := []Todo{
            {Id: 1, Description: "Learn Lua", CreateTime: time.Now(), IsDone: false},
            {Id: 2, Description: "Learn Low Level Language", CreateTime: time.Now(), IsDone: false},
            {Id: 3, Description: "Learn Something new", CreateTime: time.Now(), IsDone: false},
        }

        jsonData, err := json.Marshal(items)
        if err != nil {
            fmt.Println("Error marching JSON:", err)
            return 
        }

        file, err := os.Create(filename)
        if err != nil {
            fmt.Println("Error createing file", err)
            return
        }

        defer file.Close()

        _, err = file.Write(jsonData)
        if err != nil {
            fmt.Println("Error writing JSON to file:", err)
            return
        }

        fmt.Println("Items stored in", filename)
    } 

    args := os.Args[1:]

    if len(args) > 0 {
        switch args[0] {
        case "add":
            if (len(args) > 2) {
                addNotes(args) 
            } else {
                fmt.Println("no description provided")
                return 
            }

        case "update":
            updateTodo()
        case "rm":
            if (len(args) > 1) {
                val := args[1]

                idx, err := strconv.Atoi(val)
                if err != nil {
                    fmt.Println("Error convering the arguments to integer:", err)
                    return
                }
                deleteTodo(&globalList, idx)
            } else {
                fmt.Println("Please provide index") 
            }
        case "ls":
            listTodos()
        case "--help":
            Thelp()
        case "--version":
            Tversion()
        case "done":
            if (len(args) > 1) {
                Tdone(args) 
            } else {
                fmt.Println("provide id to mark complete")
                return 
            }
        default:
            fmt.Println("--help for args")
        }
    } else {
        fmt.Println("--help for args")
    }

}

func Tdone(arg []string) {
    flag := false
    existingTodos, err := readTodosFromFile("items.json")   
    if err != nil {
        fmt.Println("Error reading existing TODO items:", err)
    }

    num, err := strconv.Atoi(arg[1])
    if err != nil {
        fmt.Println("Err converting string to int:",  err)
        return
    }

    for i:= range existingTodos {
        if (num == existingTodos[i].Id) {
            fmt.Println(existingTodos[i].IsDone)
            existingTodos[i].IsDone = true 
            flag = true
        } 
    }

    if flag {
        err = writeTodosToFile("items.json", existingTodos)
        if err != nil {
                fmt.Println("Error writing TODO items to file:", err)
                return
        }
        fmt.Println("updated complete")
        
    } else {
        fmt.Println("No matching id found")
    }

}

func Tversion() {
    fmt.Println("Kodu version: 0.5")
}

func Thelp() {
   fmt.Println("alias for Kodu \nadd <space>[desc for todo] \nupdate <space>[id_number] <space>[new_desc_for_todo] \nrm <space>[id_number] \nls") 
}

func addNotes(argument []string) {
    result := strings.Join(argument[1:], " ")
    fmt.Print(result)
    newTodo := Todo{
        Id: len(globalList) + 1,
        Description: result,
        CreateTime: time.Now(),
        IsDone: false,
    }

    existingTodos, err := readTodosFromFile("items.json")
    if err != nil {
            fmt.Println("Error reading existing TODO items:", err)
            return
    }

    // Append the new TODO object to the existing TODO list
    existingTodos = append(existingTodos, newTodo)

    // Write the updated TODO list back to "items.json" file
    err = writeTodosToFile("items.json", existingTodos)
    if err != nil {
            fmt.Println("Error writing TODO items to file:", err)
            return
    }

    fmt.Println("New TODO item added successfully.")
}
func writeTodosToFile(filename string, todos []Todo) error {
	// Marshal TODO items into JSON format
	data, err := json.MarshalIndent(todos, "", "    ")
	if err != nil {
		return err
	}

	// Write JSON data to file
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
func readTodosFromFile(filename string) ([]Todo, error) {
	var todos []Todo

	// Read file contents
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON data into slice of Todo
	err = json.Unmarshal(data, &todos)
	if err != nil {
		return nil, err
	}

	return todos, nil
}
func updateTodo() {
    fmt.Println("update called")
}
func deleteTodo(items *[]Todo, idx int) {
    fmt.Println("delete called")
    fmt.Println(items)
    fmt.Println("IDX:", idx)

    index := -1

    for i, item := range *items {
        if item.Id == idx {
            index = i
            break
        }
    }

    if index != -1 {
        *items = append(globalList[:index], globalList[index+1:]...)
        fmt.Println("Item with ID", idx, "deleted succesfully")
    } else {
        fmt.Println("Item with ID", idx, "not found") 
    }

    updatedData, err := json.MarshalIndent(*items, "", " ") 
    if err != nil {
        fmt.Println("Error marshaling JSON", err)
        return
    }

    if err := os.WriteFile(filename, updatedData, 0064); err != nil {
        fmt.Println("Error writing file:", err)
        return
    }

    fmt.Println("JSON data updated and written to item.json")
}

func listTodos() {
    fmt.Println("view called")

    customLayout := "2006-01-02 15:04:05"

    fmt.Println("TODO LIST")
    fmt.Println("------------------------------------------------------------------------")
    fmt.Printf("%-5s | %-30s | %-20s | %-10s\n" , "ID", "Task", "Created", "Status")
    fmt.Println("------------------------------------------------------------------------")
    for _, todo := range globalList {
            fmt.Printf("%-5d | %-30s | %-20s | %-10t\n", todo.Id, todo.Description, todo.CreateTime.Format(customLayout), todo.IsDone)
    }
}


