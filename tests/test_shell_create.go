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

    Environment := map[string]string {
        "FBTEST": "123",
    }

    l, err:= shell.Create([]string{"stdin"}, []string{"stdout", "stderr"}, "",
        Environment, "c:\\Users", 
        "P0Y0M0DT0H5M0S", "P0Y0M0DT0H5M0S")
    if err != nil {
    	fmt.Println(l)
    	panic(err)
    }
    fmt.Println(l)
}