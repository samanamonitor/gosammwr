package main

import (
    "fmt"
    "github.com/samanamonitor/gosammwr/protocol"
    "os"
)

/*
Examples:

./test_protocol_enum resourceuri http://schemas.microsoft.com/wbem/wsman/1/windows/shell

./test_protocol_enum schema root/cimv2 win32_computersystem

*/

func main() {
    if len(os.Args) < 4 {
        panic("Must pass shellid, commandid and stream name as parameter")
    }
    endpoint := os.Getenv("WR_ENDPOINT")
    username := os.Getenv("WR_USERNAME")
    password := os.Getenv("WR_PASSWORD")
    keytab_file := os.Getenv("WR_KEYTAB")

    prot := protocol.Protocol{}
    err := prot.Init(endpoint, username, password, keytab_file)
    if err != nil {
        panic(err)
    }
    defer prot.Close()


    var resourceURI string
    var selectorset *map[string]string

    switch os.Args[1] {
    case "resourceuri":
        resourceURI = os.Args[2]
        selectorset = nil
    case "schema":
        resourceURI = "http://schemas.dmtf.org/wbem/cim-xml/2/cim-schema/2/*"
        temp := map[string]string{
            "__cimnamespace": os.Args[2],
            "ClassName": os.Args[3],
        }
        selectorset = &temp
    default:
        panic("Invaid parameter. Only resourceuri or schema allowed")
    }

    err, response_doc := prot.Enumerate(resourceURI, nil, nil, nil)
    if err != nil {
        fmt.Println(response_doc)
        panic(err)
    }

    fmt.Println(response_doc)
    fmt.Printf("\nDone\n")
}
