package protocol

import (
	"github.com/beevik/etree"
	"github.com/samanamonitor/gosammwr/transport"
	"github.com/google/uuid"
	"errors"
	"fmt"
)

type Protocol struct {
	transport *transport.Transport
	header_to string
	header_replyto string
	header_maxenvelopesize string
}

type Option struct {
	Type string
	Value string
}

type Filter struct {
	Dialect string
	Wql *string
	Selectorset *map[string]string
}

func (p *Protocol) Init (endpoint string,
		username string,
		password string,
		keytab_file string) error {

	var err error
	p.transport, err = transport.NewTransport(
		endpoint,
		username,
		password,
		keytab_file)
	if err != nil {
		return err
	}
	p.header_to = "http://windows-host:5985/wsman"
	p.header_replyto = "http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous"
	p.header_maxenvelopesize = "512000"
	return nil
}

func (p Protocol) PrepareHeader (resourceURI string, action string,
		selectorset *map[string]string,
		optionset *map[string]Option,
		doc *etree.Document) {
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)

	envelope := doc.CreateElement("s:Envelope")
	envelope.CreateAttr("xmlns:s", "http://www.w3.org/2003/05/soap-envelope")
	envelope.CreateAttr("xmlns:a", "http://schemas.xmlsoap.org/ws/2004/08/addressing")
	envelope.CreateAttr("xmlns:n", "http://schemas.xmlsoap.org/ws/2004/09/enumeration")
	envelope.CreateAttr("xmlns:w", "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd")
	envelope.CreateAttr("xmlns:p", "http://schemas.microsoft.com/wbem/wsman/1/wsman.xsd")
	envelope.CreateAttr("xmlns:b", "http://schemas.dmtf.org/wbem/wsman/1/cimbinding.xsd")
	envelope.CreateAttr("xmlns:wsen", "http://schemas.xmlsoap.org/ws/2004/09/enumeration")
	envelope.CreateAttr("xmlns:rsp", "http://schemas.microsoft.com/wbem/wsman/1/windows/shell")
	envelope.CreateAttr("xmlns:xs", "http://www.w3.org/2001/XMLSchema")

	header := envelope.CreateElement("s:Header")
	header.CreateElement("a:To").CreateText(p.header_to)

	address := header.CreateElement("a:ReplyTo").CreateElement("a:Address")
	address.CreateText(p.header_replyto)
	address.CreateAttr("s:mustUnderstand", "true")

	MaxEnvelopeSize := header.CreateElement("w:MaxEnvelopeSize")
	MaxEnvelopeSize.CreateAttr("s:mustUnderstand", "true")
	MaxEnvelopeSize.CreateText(p.header_maxenvelopesize)

	header.CreateElement("a:MessageID").CreateText(fmt.Sprintf("uuid:%s", uuid.NewString()))
	Locale := header.CreateElement("w:Locale") 
	Locale.CreateAttr("xml:lang", "en-US")
	Locale.CreateAttr("s:mustUnderstand", "false")

	DataLocale := header.CreateElement("p:DataLocale")
	DataLocale.CreateAttr("xml:lang", "en-US")
	DataLocale.CreateAttr("s:mustUnderstand", "false")

	ResourceURI := header.CreateElement("w:ResourceURI")
	ResourceURI.CreateAttr("s:mustUnderstand", "true")
	ResourceURI.CreateText(resourceURI)

	Action := header.CreateElement("a:Action")
	Action.CreateAttr("s:mustUnderstand", "true")
	Action.CreateText(action)

	if selectorset != nil {
		ss := header.CreateElement("w:SelectorSet")
		for key, value := range *selectorset {
			s := ss.CreateElement("w:Selector")
			s.CreateAttr("Name", key)
			s.CreateText(value)
		}
	}

	if optionset != nil {
		oset := header.CreateElement("w:OptionSet")
		for key, option := range *optionset {
			s := oset.CreateElement("w:Option")
			s.CreateAttr("Name", key)
			s.CreateAttr("Type", option.Type)
			s.CreateText(option.Value)
		}
	}
}

func (p *Protocol) Close () {
	p.transport.Close()
}

type GenericFault struct {
	Err error
	FaultXML []byte
}

func (gf *GenericFault) Error() string {
	return fmt.Sprint("GenericFault = Fault could not be decoded. Raw response is: ",
		string(gf.FaultXML), "\nInner error: ", gf.Err)
}

func processFault(Err error, responseData []byte) (error) {
	responseXML := etree.NewDocument()
	err := responseXML.ReadFromBytes(responseData)
	if err != nil {
		return Err
	}
	gf := GenericFault{
		Err: Err,
		FaultXML: responseData,
	}
	sf := responseXML.FindElement("//Body/Fault")
	if sf == nil || !(sf.FullTag() == "s:Fault" && 
			sf.NamespaceURI() == "http://www.w3.org/2003/05/soap-envelope") {
		return &gf
	}
	soapf := SOAPFault{
		Err: gf,
		FaultElement: sf,
	}
	err = soapf.Init()
	if err != nil {
		return err
	}

	for _, fault_detail := range soapf.DetailElements {
		if fault_detail.FullTag() == "f:WSManFault" && 
				fault_detail.NamespaceURI() == "http://schemas.microsoft.com/wbem/wsman/1/wsmanfault" {
			wsf := WSManFault{
				Err: soapf,
				FaultElement: fault_detail,
			}
			wsf.Init()
			return &wsf
		}
	}
	return &soapf
}

func (p *Protocol) SendMessage(doc *etree.Document) (*etree.Document, error) {
	/* TODO: check validity of the response */
	request, err := doc.WriteToBytes()
	if err != nil {
		return nil, err
	}
/*
	fmt.Print(string(request))
*/
	response, err := p.transport.SendMessage(request)
	if err != nil {
		return nil, processFault(err, response)
	}
	response_doc := etree.NewDocument()
	response_doc.ReadFromBytes(response)
	body := response_doc.FindElement("//Body")
	if body == nil {
		return response_doc, errors.New("Invalid response. Body tag not found")
	}

	var newdoc *etree.Document
	newdoc = etree.NewDocument()
	newdoc.SetRoot(body)

	return newdoc, err
}

func (p *Protocol) Get(resourceURI string,
		selectorset *map[string]string, 
		optionset *map[string]Option) (string, error) {
	action := "http://schemas.xmlsoap.org/ws/2004/09/transfer/Get"

	doc := etree.NewDocument()
	p.PrepareHeader(resourceURI, action, selectorset, optionset, doc)

	doc.FindElement("//s:Envelope").CreateElement("s:Body")

	response, err := p.SendMessage(doc)
	if err != nil && response == nil {
		return "", err
	}
	ret, _ := response.WriteToString()
	return ret, err
}

func (p *Protocol) Enumerate(resourceURI string,
		filter *Filter,
		selectorset *map[string]string,
		optionset *map[string]Option) (string, error) {
	action := "http://schemas.xmlsoap.org/ws/2004/09/enumeration/Enumerate"

	doc := etree.NewDocument()
	p.PrepareHeader(resourceURI, action, selectorset, optionset, doc)

	enumerate := doc.FindElement("//s:Envelope").CreateElement(
		"s:Body").CreateElement("wsen:Enumerate")

	if filter != nil {
		f := enumerate.CreateElement("w:Filter")
		f.CreateAttr("Dialect", filter.Dialect)
		if filter.Wql != nil {
			f.CreateText(*filter.Wql)
		} else {
			ss := f.CreateElement("w:SelectorSet")
			for key, value := range *filter.Selectorset {
				s := ss.CreateElement("w:Selector")
				s.CreateAttr("Name", key)
				s.CreateText(value)
			}
		}
	}
	response, err := p.SendMessage(doc)
	if err != nil {
		if response == nil {
			return "", err
		}
		r, _ := response.WriteToString()
		return r, err
	}
	EnumerationContext := response.FindElement("//EnumerationContext")
	if EnumerationContext == nil {
		ret, _ := response.WriteToString()
		return "", errors.New("Response did not contain EnumerationContext.\n" + ret)
	}
	return EnumerationContext.Text()[5:], nil
}

func (p *Protocol) Pull(resourceURI string,
		EnumerationContext string,
		MaxElements int,
		selectorset *map[string]string,
		optionset *map[string]Option) (string, string, error) {
	action := "http://schemas.xmlsoap.org/ws/2004/09/enumeration/Pull"

	doc := etree.NewDocument()
	p.PrepareHeader(resourceURI, action, selectorset, optionset, doc)

	pull := doc.FindElement("//s:Envelope").CreateElement(
		"s:Body").CreateElement("wsen:Pull")
	pull.CreateElement("wsen:EnumerationContext").CreateText(
		fmt.Sprint("uuid:", EnumerationContext))
	if MaxElements > 0 {
		pull.CreateElement("wsen:MaxElements").CreateText(
			fmt.Sprint(MaxElements))
	}

	response, err := p.SendMessage(doc)
	if err != nil {
		if response == nil {
			return "", "", err
		}
		r, _ := response.WriteToString()
		return r, "", err
	}
	pull_ec := response.FindElement("//EnumerationContext")
	var ec string
	if pull_ec != nil {
		ec = pull_ec.Text()[5:]
	} else {
		ec = ""
	}
	PullResponse := response.FindElement("//PullResponse")
	if PullResponse == nil {
		ret, _ := response.WriteToString()
		return "", "", errors.New("Invalid Response. PullResponse tag missing.\n" + ret)
	}
	newdoc := etree.NewDocument()
	newdoc.SetRoot(PullResponse)
	ret, _ := newdoc.WriteToString()
	return ec, ret, nil
}

func (p *Protocol) Release(resourceURI string,
		EnumerationContext string,
		optionset *map[string]Option) (string, error) {
	action := "http://schemas.xmlsoap.org/ws/2004/09/enumeration/Release"

	doc := etree.NewDocument()
	p.PrepareHeader(resourceURI, action, nil, optionset, doc)
	

	pull := doc.FindElement("//s:Envelope").CreateElement(
		"s:Body").CreateElement("wsen:Release")
	pull.CreateElement("wsen:EnumerationContext").CreateText(
		fmt.Sprint("uuid:", EnumerationContext))

	response, err := p.SendMessage(doc)
	if err != nil {
		if response == nil {
			return "", err
		}
	}
	ret, _ := response.WriteToString()
	return ret, err
}

func (p *Protocol) Create(resourceURI string, instance *etree.Element,
		optionset *map[string]Option) (string, error) {
	action := "http://schemas.xmlsoap.org/ws/2004/09/transfer/Create"

	doc := etree.NewDocument()
	p.PrepareHeader(resourceURI, action, nil, optionset, doc)
	
	doc.FindElement("//s:Envelope").CreateElement("s:Body").AddChild(instance)

	response, err := p.SendMessage(doc)
	if err != nil {
		if response == nil {
			return "", err
		}
	}
	ret, _ := response.WriteToString()
	return ret, err
}

func (p *Protocol) Delete(resourceURI string, selectorset *map[string]string,
		optionset *map[string]Option) (string, error) {
	/* TODO: selectorset is mandatory. Should fail if it is missing */
	action := "http://schemas.xmlsoap.org/ws/2004/09/transfer/Delete"

	doc := etree.NewDocument()
	p.PrepareHeader(resourceURI, action, selectorset, optionset, doc)
	
	doc.FindElement("//s:Envelope").CreateElement("s:Body")

	response, err := p.SendMessage(doc)
	if err != nil {
		if response == nil {
			return "", err
		}
	}
	ret, _ := response.WriteToString()
	return ret, err
}

func (p *Protocol) Command(resourceURI string, command_body *etree.Element, selectorset *map[string]string,
		optionset *map[string]Option) (string, error) {
	action := "http://schemas.microsoft.com/wbem/wsman/1/windows/shell/Command"

	doc := etree.NewDocument()
	p.PrepareHeader(resourceURI, action, selectorset, optionset, doc)
	
	doc.FindElement("//s:Envelope").CreateElement("s:Body").AddChild(command_body)

	response, err := p.SendMessage(doc)
	if err != nil {
		if response == nil {
			return "", err
		}
	}
	ret, _ := response.WriteToString()
	return ret, err
}

func (p *Protocol) Receive(resourceURI string, receive_body *etree.Element, selectorset *map[string]string,
		optionset *map[string]Option) (string, error) {
	action := "http://schemas.microsoft.com/wbem/wsman/1/windows/shell/Receive"

	doc := etree.NewDocument()
	p.PrepareHeader(resourceURI, action, selectorset, optionset, doc)
	
	doc.FindElement("//s:Envelope").CreateElement("s:Body").AddChild(receive_body)

	response, err := p.SendMessage(doc)
	if err != nil {
		if response == nil {
			return "", err
		}
	}
	ret, _ := response.WriteToString()
	return ret, err
}

func (p *Protocol) Send(resourceURI string, send_body *etree.Element, selectorset *map[string]string,
		optionset *map[string]Option) (string, error) {
	action := "http://schemas.microsoft.com/wbem/wsman/1/windows/shell/Send"

	doc := etree.NewDocument()
	p.PrepareHeader(resourceURI, action, selectorset, optionset, doc)
	
	doc.FindElement("//s:Envelope").CreateElement("s:Body").AddChild(send_body)

	response, err := p.SendMessage(doc)
	if err != nil {
		if response == nil {
			return "", err
		}
	}
	ret, _ := response.WriteToString()
	return ret, err
}

func (p *Protocol) Signal(resourceURI string, signal_body *etree.Element, selectorset *map[string]string, 
		optionset *map[string]Option) (string, error) {
	action := "http://schemas.microsoft.com/wbem/wsman/1/windows/shell/Signal"

	doc := etree.NewDocument()
	p.PrepareHeader(resourceURI, action, selectorset, optionset, doc)
	
	doc.FindElement("//s:Envelope").CreateElement("s:Body").AddChild(signal_body)

	response, err := p.SendMessage(doc)
	if err != nil {
		if response == nil {
			return "", err
		}
	}
	ret, _ := response.WriteToString()
	return ret, err
}
