package main

import (
    "github.com/samanamonitor/gosammwr/shell"
    "fmt"
    "os"
)

func main() {
    endpoint := os.Getenv("WR_ENDPOINT")
    username := os.Getenv("WR_USERNAME")
    password := os.Getenv("WR_PASSWORD")
    keytab_file := os.Getenv("WR_KEYTAB")
    if len(os.Args) < 2 {
        panic("Must pass ShellId as parameter")
    }

    shell := shell.Shell{}
    err := shell.Init(endpoint, username, password, keytab_file)
    defer shell.Cleanup()
    if err != nil {
        panic(err)
    }

    err, si := shell.Get(os.Args[1])
    if err != nil {
    	panic(err)
    }
    fmt.Println(si)
    fmt.Println(si.Json())
}