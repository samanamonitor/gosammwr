package protocol

import (
	"github.com/beevik/etree"
	"fmt"
	"strings"
)

type SOAPFault struct {
	Err GenericFault
	FaultElement *etree.Element
	Code string
	SubCode []string
	Reason []string
	Node string
	Role string
	FaultDetail string

	DetailElements []*etree.Element
}

func (f *SOAPFault) Init() error {
	if f.FaultElement == nil {
		return &f.Err
	}

	fault_code := f.FaultElement.FindElement("./Code/Value")
	f.Code = fault_code.Text()

	fault_subcodes := f.FaultElement.FindElements("./Code/Subcode/Value")
	for _, subcode_element := range fault_subcodes {
		f.SubCode = append(f.SubCode, subcode_element.Text())
	}

	fault_reasons := f.FaultElement.FindElements("./Reason/Text")
	for _, fault_reason := range fault_reasons {
		f.Reason = append(f.Reason, fault_reason.Text())
	}

	details := f.FaultElement.FindElements("./Detail/*")
	for _, detail := range details {
		fullTag := strings.Split(detail.FullTag(), ":")
		if len(fullTag) > 1 && fullTag[1] == "FaultDetail" {
			f.FaultDetail = detail.Text()
		} else {
			f.DetailElements = append(f.DetailElements, detail)
		}
	}
	return nil
}

func (f *SOAPFault) Error() string {
	return fmt.Sprintf("SOAP Fault = code:%s - subcodes:%s - reason:%s - fault_detail: %s", 
		f.Code, strings.Join(f.SubCode, ","), 
		strings.Join(f.Reason, ". "),
		f.FaultDetail)
}

