package main

import (
    "fmt"
    "github.com/samanamonitor/gosammwr/shell"
    "os"
)


func main() {
    if len(os.Args) < 2 {
        panic("Must pass shellid as parameter and the command and arguments")
    }
    endpoint := os.Getenv("WR_ENDPOINT")
    username := os.Getenv("WR_USERNAME")
    password := os.Getenv("WR_PASSWORD")
    keytab_file := os.Getenv("WR_KEYTAB")

/*
    WINRS_SKIP_CMD_SHELL := len(os.Getenv("WINRS_SKIP_CMD_SHELL")) > 0
*/
    shell := shell.Shell{}
    err := shell.Init(endpoint, username, password, keytab_file)
    if err != nil {
        panic(err)
    }
    defer shell.Cleanup()

    response, err := shell.Command(os.Args[1], os.Args[2:], true, true, nil)
    if err != nil {
        fmt.Println(response)
        panic(err)
    }

    fmt.Println(response)
}