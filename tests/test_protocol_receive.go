package main

import (
    "fmt"
    "github.com/samanamonitor/gosammwr/protocol"
    "github.com/beevik/etree"
    "os"
    "bufio"
    "strings"
    b64 "encoding/base64"
)

/*

example:

bin/test_protocol_receive http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd <shellid> <commandid>
*/


func main() {
    if len(os.Args) < 2 {
        panic("Must pass shellid and commandid(optional) as parameter")
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
    var ShellId string
    var CommandId string


    if os.Args[1] == "-" {
        reader := bufio.NewReader(os.Stdin)
        input, _ := reader.ReadString('\n')
        input = strings.Trim(input, "\n")

        temp := strings.Split(input, " ")
        if len(temp) != 3 {
            panic("Invalid number of paramenters. Expecting <resourceuri> <ShellId> <CommandId>.\n" + input)
        }
        resourceURI = temp[0]
        ShellId = temp[1]
        CommandId = temp[2]
    } else {
        if len(os.Args) != 4 {
            panic("Ivalid number of parameters. Expecting <resourceuri> <ShellId> <CommandId>.")
        }
        resourceURI = os.Args[1]
        ShellId = os.Args[2]
        CommandId = os.Args[3]
    }

    selectorset := map[string]string{
        "ShellId": ShellId,
    }

    Receive := etree.NewElement("rsp:Receive")
    DesiredStream := Receive.CreateElement("rsp:DesiredStream")
    DesiredStream.CreateAttr("CommandId", CommandId)
    DesiredStream.CreateText("stdout stderr")

    response, err := prot.Receive(resourceURI, Receive, &selectorset, nil)
    if err != nil {
        fmt.Println(response)
        panic(err)
    }

    resdoc := etree.NewDocument()
    resdoc.ReadFromString(response)
    stdout := resdoc.FindElements("//Stream[@Name='stdout']")

    for _, item := range stdout {
        temp_str2, _ := b64.StdEncoding.DecodeString(item.Text())

        fmt.Print(string(temp_str2))

    }
}