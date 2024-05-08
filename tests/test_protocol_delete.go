package main

import (
    "fmt"
    "github.com/samanamonitor/gosammwr/sammwr/protocol"
    "os"
)


func main() {
    if len(os.Args) < 2 {
        panic("Must pass shellid as parameter")
    }
    ShellId := os.Args[1]

    endpoint := os.Getenv("WR_ENDPOINT")
    username := os.Getenv("WR_USERNAME")
    password := os.Getenv("WR_PASSWORD")
    keytab_file := os.Getenv("WR_KEYTAB")

    resourceURI := "http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd"

    prot := protocol.Protocol{}
    err := prot.Init(endpoint, username, password, keytab_file)
    if err != nil {
        panic(err)
    }
    defer prot.Close()
    fmt.Printf("Init Complete\n")

    selectorset := map[string]string {
        "ShellId": ShellId,
    }

    err, response := prot.Delete(resourceURI, &selectorset, nil)
    if err != nil {
        fmt.Println(response)
        panic(err)
    }

    fmt.Println(response)
    fmt.Printf("\nDone\n")
}