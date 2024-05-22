package main

import (
    "fmt"
    "github.com/samanamonitor/gosammwr/shell"
    "os"
)


func main() {
    if len(os.Args) < 4 {
        panic("Must pass shellid and commandid as parameter")
    }
    endpoint := os.Getenv("WR_ENDPOINT")
    username := os.Getenv("WR_USERNAME")
    password := os.Getenv("WR_PASSWORD")
    keytab_file := os.Getenv("WR_KEYTAB")

    shell := shell.Shell{}
    err := shell.Init(endpoint, username, password, keytab_file)
    if err != nil {
        panic(err)
    }
    defer shell.Cleanup()
    fmt.Printf("Init Complete\n")

    response, err := shell.Send(os.Args[1], os.Args[2], os.Args[3] + "\r\n", "stdin", false)
    if err != nil {
        panic(err)
    }
    fmt.Println(response)
    fmt.Println("Done")
}