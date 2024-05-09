package main

import (
    "fmt"
    "github.com/samanamonitor/gosammwr/protocol"
    "os"
)


func main() {
    if len(os.Args) < 3 {
        panic("Must pass shellid as parameter and the command and arguments")
    }
    endpoint := "http://smnnovmsfs01.samana.local:5985/wsman"
    keytab_file := "samanasvc2.keytab"

    prot := protocol.Protocol{}
    err := prot.Init(endpoint, nil, nil, &keytab_file)
    if err != nil {
        panic(err)
    }
    defer prot.Close()
    fmt.Printf("Init Complete\n")


    err, response := prot.ShellCommand(os.Args[1], os.Args[2:], true, true, nil)
    if err != nil {
        fmt.Println(response)
        panic(err)
    }

    fmt.Println(response)
    fmt.Printf("\nDone\n")
}