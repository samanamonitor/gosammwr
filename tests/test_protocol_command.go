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
    if len(os.Args) < 2 {
        panic("Invalid number of parameters.")
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

    var resourceURI string
    var ShellId string
    var command []string

    if os.Args[1] == "-" {
        reader := bufio.NewReader(os.Stdin)
        input, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println(err)
            panic("Invalid input from stdin")
        }
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

    ss := protocol.SelectorSet{
        "ShellId": ShellId,
    }
    optionset := protocol.OptionSet{}
    optionset.Add("WINRS_CONSOLEMODE_STDIN", "xs:boolean", "TRUE")
    optionset.Add("WINRS_SKIP_CMD_SHELL", "xs:boolean", "TRUE")

    if len(command) < 1 {
        panic("command array must have at least 1 element")
    }
    CommandLine := etree.NewElement("rsp:CommandLine")
    CommandLine.CreateElement("rsp:Command").CreateText(command[0])
    for i := 1; i < len(command); i++ {
        CommandLine.CreateElement("rsp:Arguments").CreateText(command[i])
    }

    body_str, err := prot.Command(resourceURI, CommandLine, ss, optionset)
    if err != nil {
        panic(err)
    }
    Body := etree.NewDocument()
    Body.ReadFromString(body_str)
    CommandId := Body.FindElement("//CommandId")
    fmt.Println(resourceURI, ShellId, CommandId.Text())
}
