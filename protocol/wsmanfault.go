package protocol

import (
	"github.com/beevik/etree"
	"fmt"
	"strconv"
)

type WSManFault struct {
	Err *SOAPFault
	FaultElement *etree.Element
	Code int
	Machine string
	Message string
	Provider *ProviderFault
}

func (wf *WSManFault) Error() string {
	out := fmt.Sprintf("WSManFault: code=0x%08x machine=%s message=%s", wf.Code, wf.Machine, wf.Message)
	if wf.Provider != nil {
		out += " -- " + wf.Provider.Error()
	}
	return out
}

func (wf *WSManFault) Init() error {
	if wf.FaultElement == nil {
		return wf.Err
	}
	code := wf.FaultElement.SelectAttr("Code")
	if code == nil {
		panic("Invalid error " + wf.Err.Error())
	}
	var err error
	wf.Code, err = strconv.Atoi(code.Value)
	if err != nil {
		panic("Invalid error (invalid code) " + wf.Err.Error())
	}
	machine := wf.FaultElement.SelectAttr("Machine")
	if machine == nil {
		panic("Invalid error (no machine) " + wf.Err.Error())
	}
	wf.Machine = machine.Value
	message := wf.FaultElement.FindElement("./Message")
	if message == nil {
		panic("Invalid error (no message) " + wf.Err.Error())
	}
	wf.Message = message.Text()
	provider := message.FindElement("./ProviderFault")
	if provider == nil {
		return wf.Err
	}
	pf := ProviderFault{
		FaultElement: provider,
	}
	pf.attributes = make(map[string]string)
	pf.Init()
	wf.Provider = &pf
	return nil
}
