package main

import (
    "fmt"
    "github.com/samanamonitor/gosammwr/protocol"
    "os"
    "github.com/beevik/etree"
    "bufio"
    "strings"
)

/*
bin/test_protocol_command http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd <shellid>
*/


func main() {
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
    var ShellId string
    var command []string

    if os.Args[1] == "-" {
        reader := bufio.NewReader(os.Stdin)
        input, _ := reader.ReadString('\n')
        input = strings.Trim(input, "\n")
        temp := strings.Split(input, " ")
        resourceURI = temp[0]
        ShellId = temp[1]
        command = os.Args[2:]
    } else {
        resourceURI = os.Args[1]
        ShellId = os.Args[2]
        command = os.Args[3:]
    }

    selectorset := map[string]string{
        "ShellId": ShellId,
    }
    optionset := map[string]protocol.Option {
        "WINRS_CONSOLEMODE_STDIN": { Value: "TRUE", Type: "xs:boolean" },
        "WINRS_SKIP_CMD_SHELL": { Value: "TRUE", Type: "xs:boolean" },
    }

    if len(command) < 1 {
        panic("command array must have at least 1 element")
    }
    CommandLine := etree.NewElement("rsp:CommandLine")
    CommandLine.CreateElement("rsp:Command").CreateText(command[0])
    for i := 1; i < len(command); i++ {
        CommandLine.CreateElement("rsp:Arguments").CreateText(command[i])
    }

    err, body_str := prot.Command(resourceURI, CommandLine, &selectorset, &optionset)
    if err != nil {
        panic(err)
    }
    Body := etree.NewDocument()
    Body.ReadFromString(body_str)
    CommandId := Body.FindElement("//CommandId")
    fmt.Println(resourceURI, ShellId, CommandId.Text())
}
