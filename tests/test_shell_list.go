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

    shell, err := shell.NewShell(endpoint, username, password, keytab_file)
    if err != nil {
        panic(err)
    }
    defer shell.Cleanup()

    l, err := shell.List()
    if err != nil {
    	fmt.Println(l)
    	panic(err)
    }
    for _, s := range(l) {
        fmt.Print(s)
        fmt.Print(" ")
    }
    fmt.Println()
    /*
    fmt.Println(l)
    for _, v := range l {
        fmt.Println(v.Json())
    }
    */
}