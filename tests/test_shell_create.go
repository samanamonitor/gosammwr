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

    shell := shell.Shell{}
    err := shell.Init(endpoint, username, password, keytab_file)
    defer shell.Cleanup()
    if err != nil {
        panic(err)
    }

    Environment := map[string]string {
        "FBTEST": "123",
    }

    err, l := shell.Create("stdin", "stdout stderr", "",
        Environment, "c:\\Users", 
        "P0Y0M0DT0H5M0S", "P0Y0M0DT0H5M0S")
    if err != nil {
    	fmt.Println(l)
    	panic(err)
    }
    fmt.Println(l)
}