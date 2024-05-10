package shell

import (
	"github.com/samanamonitor/gosammwr/protocol"
	"github.com/beevik/etree"
	"errors"
	"strings"
	"strconv"
	b64 "encoding/base64"
)

type Shell struct {
	prot protocol.Protocol
	ShellId string
	CommandIds []string 
}

type ShellInstance struct {
	ShellId string
	ResourceUri string
	Owner string
	ClientIP string
	Environment map[string]string
	WorkingDirectory string
	IdleTimeOut string
	InputStreams string
	OutputStreams string
	ShellRunTime string
	ShellInactivity string
}

var boolToString = map[bool]string{
		true: "TRUE",
		false: "FALSE",
}

func (s *Shell) Init(endpoint string,
		username string,
		password string,
		keytab_file string) (error) {
	err := s.prot.Init(endpoint, username, password, keytab_file)

	return err
}

func (s *Shell) Cleanup() {
	s.prot.Close()
}

func (s *Shell) List() (error, []ShellInstance){
	resourceURI := "http://schemas.microsoft.com/wbem/wsman/1/windows/shell"
	err, EnumerationContext := s.prot.Enumerate(resourceURI, nil, nil, nil)
	if err != nil {
		return err, []ShellInstance{}
	}
	var out []ShellInstance
	for ; len(EnumerationContext) > 0; {
		var PullResponse string
		err, EnumerationContext, PullResponse = s.prot.Pull(resourceURI, EnumerationContext[5:], -1, nil, nil)
		if err != nil {
			return err, []ShellInstance{}
		}
		et_PullResponse := etree.NewDocument()
		et_PullResponse.ReadFromString(PullResponse)
		Items := et_PullResponse.FindElement("//Items").ChildElements()
		for _, v := range Items {
			temp := ShellInstance{}
			temp.FromEtreeElement(v)
			out = append(out, temp)
		}

	}
	return err, out
}

func (s *Shell) Get(ShellId string) (error, ShellInstance) {
	resourceURI := "http://schemas.microsoft.com/wbem/wsman/1/windows/shell"
	selectorset := map[string]string {
		"ShellId": ShellId,
	}
	err, si_xml := s.prot.Get(resourceURI, &selectorset, nil)
	if err != nil {
		return err, ShellInstance{}
	}
	GetResponse := etree.NewDocument()
	GetResponse.ReadFromString(si_xml)
	si := ShellInstance{}
	si.FromEtreeElement(GetResponse.FindElement("//Shell"))
	return nil, si

}

func (s *Shell) Create(InputStreams []string, OutputStreams []string, Name string,
		Environment map[string]string, WorkingDirectory string, 
		Lifetime string, IdleTimeOut string)  (error, ShellInstance) {
	/*
	TODO: Create a struct for Resource Created object?
	https://learn.microsoft.com/en-us/openspecs/windows_protocols/ms-wsmv/593f3ed0-0c7a-4158-a4be-0b429b597e31
	*/
	resourceURI := "http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd"

	shell := etree.NewElement("rsp:Shell")
	shell.CreateElement("rsp:InputStreams").CreateText(strings.Join(InputStreams, " "))
	shell.CreateElement("rsp:OutputStreams").CreateText(strings.Join(OutputStreams, " "))

	if len(Name) > 0 {
		shell.CreateElement("rsp:Name").CreateText(Name)
	}
	if len(WorkingDirectory) > 0 {
		shell.CreateElement("rsp:WorkingDirectory").CreateText(WorkingDirectory)
	}
	if len(Lifetime) > 0 {
		// iso8601_duration
		shell.CreateElement("rsp:Lifetime").CreateText(Lifetime)
	}
	if len(IdleTimeOut) > 0 {
		// iso8601_duration
		shell.CreateElement("rsp:IdleTimeOut").CreateText(IdleTimeOut)
	}
	if len(Environment) > 0 {
		variables := shell.CreateElement("rsp:Environment")
		for key, val := range Environment {
			variable := variables.CreateElement("rsp:Variable")
			variable.CreateAttr("Name", key)
			variable.CreateText(val)
		}
	}

	optionset := map[string]protocol.Option{}
	optionset["WINRS_NOPROFILE"] = protocol.Option{
		Value: "FALSE",
		Type: "xs:boolean",
	}
	optionset["WINRS_CODEPAGE"] = protocol.Option{
		Value: "437",
		Type: "xs:unsignedInt",
	}
	err, body_str := s.prot.Create(resourceURI, shell, &optionset)
	if err != nil {
		return err, ShellInstance{}
	}
	Body := etree.NewDocument()
	Body.ReadFromString(body_str)
	Shell := Body.FindElement("//Shell")
	si := ShellInstance{}
	si.FromEtreeElement(Shell)
	return nil, si
}

func (s *Shell) Delete(ShellId string) (error) {
	resourceURI := "http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd"
    selectorset := map[string]string {
        "ShellId": ShellId,
    }
    err, _ := s.prot.Delete(resourceURI, &selectorset, nil)
    return err
}

func (s *Shell) Command(ShellId string, command []string, SkipCmdShell bool, ConsoleModeStdin bool, 
		optionset *map[string]protocol.Option) (error, string) {
	resourceURI := "http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd"

	selectorset := map[string]string{
		"ShellId": ShellId,
	}
	if optionset == nil {
		optionset = &(map[string]protocol.Option{})
	}
	(*optionset)["WINRS_CONSOLEMODE_STDIN"] = protocol.Option{
		Value: boolToString[ConsoleModeStdin],
		Type: "xs:boolean",
	}
	(*optionset)["WINRS_SKIP_CMD_SHELL"] = protocol.Option{
		Value: boolToString[SkipCmdShell],
		Type: "xs:boolean",
	}

	if len(command) < 1 {
		return errors.New("command array must have at least 1 element"), ""
	}
	CommandLine := etree.NewElement("rsp:CommandLine")
	CommandLine.CreateElement("rsp:Command").CreateText(command[0])
	for i := 1; i < len(command); i++ {
		CommandLine.CreateElement("rsp:Arguments").CreateText(command[i])
	}

	err, body_str := s.prot.Command(resourceURI, CommandLine, &selectorset, optionset)
	if err != nil {
		return err, ""
	}
	Body := etree.NewDocument()
	Body.ReadFromString(body_str)
	CommandId := Body.FindElement("//CommandId")
	return nil, CommandId.Text()
}

const (
	CommandDone = iota
	CommandPending = iota
	CommandRunning = iota
)

type CommandState struct {
	StateEnum int
	ExitCode int
}

type StreamType struct {
	End bool
	Data []byte
}

func (s *Shell) Receive(ShellId string, CommandId string, Streams []string) (error, CommandState, []StreamType) {
	resourceURI := "http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd"

	selectorset := map[string]string{
		"ShellId": ShellId,
	}

	var cs CommandState
	response := make([]StreamType, len(Streams))

    Receive := etree.NewElement("rsp:Receive")
    DesiredStream := Receive.CreateElement("rsp:DesiredStream")
    DesiredStream.CreateAttr("CommandId", CommandId)
    DesiredStream.CreateText(strings.Join(Streams, " "))


    err, body_str := s.prot.Receive(resourceURI, Receive, &selectorset, nil)
    if err != nil {
		return err, cs, response
    }
	Body := etree.NewDocument()
	Body.ReadFromString(body_str)
	for i, stream := range(Streams) {
		stream_elements := Body.FindElements("//rsp:Stream[@Name='" + stream + "']")
		for _, element := range(stream_elements) {
			temp_str2, err := b64.StdEncoding.DecodeString(element.Text())
			if err != nil {
				return err, cs, response
			}
			response[i].Data = append(response[i].Data, temp_str2...)
			end_attr := element.SelectAttr("End")
			if end_attr != nil {
				if end_attr.Value == "true" {
					response[i].End = true
				} else {
					response[i].End = false
				}
			}
		}
	}
	cs_element := Body.FindElement("//rsp:CommandState")
	state_attr := cs_element.SelectAttr("State")
	if state_attr == nil {
		err = errors.New("Invalid response. CommandState is not in response.\n" + body_str)
		return err, cs, response
	}
	if state_attr.Value == "http://schemas.microsoft.com/wbem/wsman/1/windows/shell/CommandState/Done" {
		cs.StateEnum = CommandDone
		exit_code_element := cs_element.FindElement("//rsp:ExitCode")
		if exit_code_element == nil {
			err = errors.New("Invalid CommandState response. Missing ExitCode.\n" + body_str)
			return err, cs, response
		}
		cs.ExitCode, err = strconv.Atoi(exit_code_element.Text())
		if err != nil {
			return err, cs, response
		}
	} else if state_attr.Value == "http://schemas.microsoft.com/wbem/wsman/1/windows/shell/CommandState/Pending" {
		cs.StateEnum = CommandPending
	} else if state_attr.Value == "http://schemas.microsoft.com/wbem/wsman/1/windows/shell/CommandState/Running" {
		cs.StateEnum = CommandRunning
	} else {
		err = errors.New("Invalid CommandState value.\n" + body_str)
		return err, cs, response
	}

	return err, cs, response
}

func (s *Shell) Send(ShellId string, CommandId string, Data string, Stream string, End bool) (error, string) {
    resourceURI := "http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd"

    selectorset := map[string]string {
        "ShellId": ShellId,
    }

    send_element := etree.NewElement("rsp:Send")
    stream_element := send_element.CreateElement("rsp:Stream")
    stream_element.CreateAttr("Name", Stream)
    stream_element.CreateAttr("CommandId", CommandId)
    stream_element.CreateAttr("End", boolToString[End])
    stream_element.CreateText(b64.StdEncoding.EncodeToString([]byte(Data)))

    err, response := s.prot.Send(resourceURI, send_element, &selectorset, nil)
    if err != nil {
    	return err, ""
    }
    return err, response
}

func (si *ShellInstance) FromEtreeElement(ete *etree.Element) {
	for _, i := range ete.ChildElements() {
		switch tag := i.Tag; tag {
		case "ShellId":
			si.ShellId = i.Text()
		case "ResourceUri":
			si.ResourceUri = i.Text()
		case "Owner":
			si.Owner = i.Text()
		case "ClientIP":
			si.ClientIP = i.Text()
		case "Environment":
			si.Environment = make(map[string]string)
			for _, variable := range i.FindElements("//Variable") {
				si.Environment[variable.SelectAttrValue("Name", "")] = variable.Text()
			}
		case "WorkingDirectory":
			si.WorkingDirectory = i.Text()
		case "IdleTimeOut":
			si.IdleTimeOut = i.Text()
		case "InputStreams":
			si.InputStreams = i.Text()
		case "OutputStreams":
			si.OutputStreams = i.Text()
		case "ShellRunTime":
			si.ShellRunTime = i.Text()
		case "ShellInactivity":
			si.ShellInactivity = i.Text()
		}
	}
}

func (si ShellInstance) String() string {
	return si.ShellId
}

func (si ShellInstance) Json() string {
	out := "{"
	out += "\"ShellId\": \"" + si.ShellId + "\","
	out += "\"ResourceUri\": \"" + si.ResourceUri + "\","
	out += "\"Owner\": \"" + si.Owner + "\","
	out += "\"ClientIP\": \"" + si.ClientIP + "\","
	out += "\"Environment\": {"
	first := true
	for varname, varvalue := range si.Environment {
		if ! first {
			out += ","
		}
		out += "\"" + varname + "\": \"" + varvalue + "\""
	}
	out += "},"
	out += "\"WorkingDirectory\": \"" + si.WorkingDirectory + "\","
	out += "\"IdleTimeOut\": \"" + si.IdleTimeOut + "\","
	out += "\"InputStreams\": \"" + si.InputStreams + "\","
	out += "\"OutputStreams\": \"" + si.OutputStreams + "\","
	out += "\"ShellRunTime\": \"" + si.ShellRunTime + "\","
	out += "\"ShellInactivity\": \"" + si.ShellInactivity + "\""
	out += "}"
	return out
}
