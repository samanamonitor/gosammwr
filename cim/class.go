package cim

import (
	"github.com/samanamonitor/gosammwr/protocol"
	"github.com/beevik/etree"
	"strconv"
	"errors"
	"fmt"
)

type PropertyStruct struct {
	Name string
	Type string
	ClassOrigin string
	Propagated bool
}

func (p *PropertyStruct) ReadFromElement(e *etree.Element) error {
	var err error
	err = nil
	for _, a := range e.Attr {
		if a.Key == "NAME" {
			p.Name = a.Value
		} else if a.Key == "TYPE" {
			p.Type = a.Value
		} else if a.Key == "CLASSORIGIN" {
			p.ClassOrigin = a.Value
		} else if a.Key == "PROPAGATED" {
			p.Propagated, err = strconv.ParseBool(a.Value)
		} else {
			err = errors.New(fmt.Sprint("Invalid attribute ", a.Key))
		}
	}
	return err
}

func (p *PropertyStruct) String() string {
	return fmt.Sprintf("  <name=%s type=%s classorigin=%s propagated=%s>\n",
		p.Name, p.Type, p.ClassOrigin, strconv.FormatBool(p.Propagated))
}

type ParameterStruct struct {
	Name string
	Type string
}

type MethodStruct struct {
	Name string
	Type string
	Parameter []ParameterStruct
}

type CimClass struct {
	protocol *protocol.Protocol
	classXML string
	root *etree.Document
	Property map[string]PropertyStruct
	Method map[string]MethodStruct 
	ClassName string
	SuperClassName string
}

func NewClass(endpoint string, username string, password string, keytab_file string) (*CimClass, error) {
	c := CimClass{}
	var err error
	c.protocol, err = protocol.NewProtocol(endpoint, username, password, keytab_file)
	return &c, err
}

func (c *CimClass) Get(namespace string, ClassName string) (error) {
	var err error
	c.ClassName = ClassName
	resourceURI := "http://schemas.dmtf.org/wbem/cim-xml/2/cim-schema/2/*"

	selectorset := protocol.SelectorSet{}
	selectorset["__cimnamespace"] = "root/cimv2"
	selectorset["ClassName"]=c.ClassName

	optionset := protocol.OptionSet{}

	c.classXML, err = c.protocol.Get(resourceURI, selectorset, optionset)
	c.ProcessClass()
	return err
}

func (c *CimClass) ProcessClass() error {
	c.root = etree.NewDocument()
	err := c.root.ReadFromString(c.classXML)
	if err != nil {
		return err
	}

	for _, e := range c.root.FindElements("//CLASS/*") {
		if e.Tag == "PROPERTY" {
			p := PropertyStruct{}
			err := p.ReadFromElement(e)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *CimClass) String() string {
	out := fmt.Sprintf("[ name=%s superclass=%s \n")
	for _, v := range c.Property {
		out += v.String()
	}
	out += fmt.Sprintln("]")
	return out
}
