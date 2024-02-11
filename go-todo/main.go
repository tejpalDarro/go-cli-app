package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
) 

type Todo struct {
    Id  int
    Description string
    TargetDate time.Time
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


    // globalList = append(globalList, Todo{
    //     Id: 1,
    //     Description: "Learn AWS",
    //     TargetDate: time.Now().AddDate(0,0,4),
    //     IsDone: false,
    // })
    // globalList = append(globalList, Todo{
    //     Id: 2,
    //     Description: "Learn Haskel",
    //     TargetDate: time.Now().AddDate(0,0,9),
    //     IsDone: false,
    // })

        items := []Todo{
            {Id: 1, Description: "Learn Lua", TargetDate: time.Now().AddDate(0,0,1), IsDone: false},
            {Id: 2, Description: "Learn Low Level Language", TargetDate: time.Now().AddDate(0,2,9), IsDone: false},
            {Id: 3, Description: "Learn Something new", TargetDate: time.Now().AddDate(3,1,1), IsDone: false},
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
            addNotes() 
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
        default:
            fmt.Println("type help for alias")
        }
    } else {
        fmt.Println("no args called!")
    }

}

func addNotes() {
    fmt.Println("Add called")

    // list := Todo {
    //     Id: len(globalList) + 1,
    //     Description: "Learn Something new",
    //     TargetDate: time.Now().AddDate(0,9,1),
    //     IsDone: false,
    // }
    globalList = append(globalList, Todo{
        Id: 2,
        Description: "Learn Haskel",
        TargetDate: time.Now().AddDate(0,0,9),
        IsDone: false,
    })
    listTodos()
    // globalList = append(globalList, list) 
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

    for _, todo := range globalList {
        fmt.Println(todo)
    }
}


