package protocol

import (
	"github.com/beevik/etree"
	"fmt"
)

type ProviderFault struct {
	Err *WSManFault
	FaultElement *etree.Element
	attributes map[string]string
	Fault *WSManFault
}

func (pf *ProviderFault) Init() error {
	if pf.FaultElement == nil {
		return pf.Err
	}
	for _, a := range pf.FaultElement.Attr {
		pf.attributes[a.Key] = a.Value
	}
	subfault := pf.FaultElement.FindElement("./WSManFault")
	if subfault == nil {
		return pf
	}
	sf := WSManFault{
		Err: nil,
		FaultElement: subfault,
	}
	pf.Fault = &sf

	pf.Fault.Init()
	return nil
}

func (pf *ProviderFault) Error() string {
	out := "ProviderFault: "
	for k, v := range pf.attributes {
		out += fmt.Sprintf("%s=%s,", k, v)
	}
	if pf.Fault == nil {
		return out
	}
	return out + " -- " + pf.Fault.Error()
}

