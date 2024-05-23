package main

import (
    "fmt"
    "github.com/samanamonitor/gosammwr/protocol"
    "os"
    "log/slog"
)

/*
Examples:

bin/test_protocol_enum resourceuri http://schemas.microsoft.com/wbem/wsman/1/windows/shell

bin/test_protocol_enum schema root/cimv2 win32_computersystem

Enum on schemas requires selectorset to be empty or just to have __cimclassname all other options will generate error

Enum on WMI with filter requires either __cimclassname in selectorset + filter selectorset of an index or WQL
Example:
    enum on 
        Win32_DiskDrive
        head->Selectorset->selector->__cimclassname=root/cimv2
        body->Enumerate->Filter->SelectorSet->Selector->Index=1 
    ####will succeed

    enum on 
        Win32_DiskDrive
        head->Selectorset->selector->__cimclassname=root/cimv2
        body->Enumerate->Filter->SelectorSet->Selector->DeviceID=\\\\.\\PHYSICALDRIVE0
    ####will fail


resource uri http://schemas.dmtf.org/wbem/wscim/1/cim-schema/2/<class_name> requires 
        head->Selectorset->selector->__cimclassname=...

bin/test_protocol_enum resourceuri http://schemas.dmtf.org/wbem/wscim/1/cim-schema/2/Win32_DiskDrive \
    "f:DeviceID=\\\\.\\PHYSICALDRIVE0" "h:__cimnamespace=root/cimv2"
*/

func main() {
    if len(os.Args) < 2 {
        panic("Must pass more parameters as parameter")
    }
    endpoint := os.Getenv("WR_ENDPOINT")
    username := os.Getenv("WR_USERNAME")
    password := os.Getenv("WR_PASSWORD")
    keytab_file := os.Getenv("WR_KEYTAB")

    prot, err := protocol.NewProtocol(endpoint, username, password, keytab_file)
    if err != nil {
        panic(err)
    }
    defer prot.Close()
    prot.SetLogLevel(slog.LevelInfo)

    var resourceURI string
    enumoptions := protocol.EnumerationOptions{}
    ss := protocol.SelectorSet{}
    optionset := protocol.OptionSet{}

    switch os.Args[1] {
    case "resourceuri":
        resourceURI = os.Args[2]
        enumoptions.Filter.SelectorSet = make(map[string]string)
        string_selectors := os.Args[3:]
        for i := range string_selectors {
            if string_selectors[i][:2] == "f:" {
                err := enumoptions.Filter.SelectorSet.AddString(string_selectors[i][2:])
                if err != nil {
                    panic(err)
                }
            } else if string_selectors[i][:2] == "h:" {
                err := ss.AddString(string_selectors[i][2:])
                if err != nil {
                    panic(err)
                }
            } else if string_selectors[i][:2] == "o:" {
                err := optionset.AddString(string_selectors[i][2:])
                if err != nil {
                    panic(err)
                }
            } else if string_selectors[i][:2] == "e:" {
                if string_selectors[i][2:] == "optimize" {
                    enumoptions.SetOptimizeEnumeration()
                }
            }
        }
    case "schema":
        optionset.Add("IncludeClassOrigin", "xs:boolean", "FALSE")
        resourceURI = "http://schemas.dmtf.org/wbem/cim-xml/2/cim-schema/2/*"
    case "wql":
        if len(os.Args) != 4 {
            panic("Invalid number of paramters. For 'wql' you must provide cimnamespace and query")
        }
        resourceURI = "http://schemas.dmtf.org/wbem/wscim/1/*"
        enumoptions.SetOptimizeEnumeration()
        enumoptions.Filter.Wql = os.Args[3]
        ss["__cimnamespace"] = os.Args[2]
    default:
        panic("Invaid parameter. Only resourceuri or schema allowed")
    }

    ec, response_doc, err := prot.Enumerate(resourceURI, ss, optionset, enumoptions)
    if err != nil {
        fmt.Println(response_doc)
        panic(err)
    }

    fmt.Println(resourceURI, ec, response_doc)
}
