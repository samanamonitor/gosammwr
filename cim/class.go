package cim

import (
	"github.com/samanamonitor/gosammwr/protocol"
)

type CimClass struct {
	protocol *protocol.Protocol
	ClassName string
	ClassXML string
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

	c.ClassXML, err = c.protocol.Get(resourceURI, selectorset, optionset)
	return err
}