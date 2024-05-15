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
    if len(os.Args) < 2 {
        panic("Must pass more parameters as parameter")
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
    var filter *protocol.Filter
    var selectorset *map[string]string
    var optionset *map[string]protocol.Option

    switch os.Args[1] {
    case "resourceuri":
        resourceURI = os.Args[2]
        if len(os.Args) == 5 {
            t_selectorset := map[string]string {
                os.Args[3]: os.Args[4],
            }
            t_filter := protocol.Filter{
                Dialect: "http://schemas.dmtf.org/wbem/wsman/1/wsman/SelectorFilter",
                Selectorset: &t_selectorset,
            }
            filter = &t_filter
        } else {
            filter = nil
        }
        selectorset = nil
        optionset = nil
    case "schema":
        t_optionset := map[string]protocol.Option{
            "IncludeClassOrigin": {
                Type: "xs:boolean",
                Value: "FALSE",
            },
        }
        optionset = &t_optionset
        resourceURI = "http://schemas.dmtf.org/wbem/cim-xml/2/cim-schema/2/*"
    case "wql":
        optionset = nil
        resourceURI = "http://schemas.microsoft.com/wbem/wsman/1/wmi/root/cimv2/*"
        t_filter := protocol.Filter{}
        t_filter.Dialect = "http://schemas.microsoft.com/wbem/wsman/1/WQL"
        t_filter.Selectorset = nil
        t_filter.Wql = &os.Args[2]
        filter = &t_filter
    default:
        panic("Invaid parameter. Only resourceuri or schema allowed")
    }

    err, response_doc := prot.Enumerate(resourceURI, filter, selectorset, optionset)
    if err != nil {
        fmt.Println(response_doc)
        panic(err)
    }

    fmt.Println(resourceURI, response_doc)
}
