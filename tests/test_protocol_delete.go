package main

import (
    "fmt"
    "github.com/samanamonitor/gosammwr/protocol"
    "os"
    "bufio"
    "strings"
)

/*
Example

bin/test_protocol_delete http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd <ShellId>

*/


func main() {
    endpoint := os.Getenv("WR_ENDPOINT")
    username := os.Getenv("WR_USERNAME")
    password := os.Getenv("WR_PASSWORD")
    keytab_file := os.Getenv("WR_KEYTAB")

    var resourceURI string
    var ShellId string

    if os.Args[1] == "-" {
        reader := bufio.NewReader(os.Stdin)
        input, _ := reader.ReadString('\n')
        input = strings.Trim(input, "\n")
        temp := strings.Split(input, " ")
        if len(temp) != 2 {
            panic("Invalid number of paramenters. Expecting resourceuri and shellid")
        }
        resourceURI = temp[0]
        ShellId = temp[1]

    } else {
        resourceURI = os.Args[1]
        ShellId = os.Args[2]
    }

    prot, err := protocol.NewProtocol(endpoint, username, password, keytab_file)
    if err != nil {
        panic(err)
    }
    defer prot.Close()
    ss := protocol.SelectorSet {
        "ShellId": ShellId,
    }

    response, err := prot.Delete(resourceURI, ss, nil)
    if err != nil {
        fmt.Println(response)
        panic(err)
    }
}