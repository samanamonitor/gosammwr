package protocol

import (
	"github.com/beevik/etree"
	"errors"
	"strings"
)

type SelectorSet map[string]string

func (s *SelectorSet) AddString(value string) error {
	temp := strings.Split(value, "=")
	if len(temp) != 2 {
		errors.New("Selectors can only be set using format key=value")
	}
	(*s)[temp[0]] = temp[1]
	return nil
}

func (s *SelectorSet) ToElement(e *etree.Element) {
	if len(*s) == 0 {
		return
	}
	ss := e.CreateElement("w:SelectorSet")
	for key, value := range *s {
		selector := ss.CreateElement("w:Selector")
		selector.CreateAttr("Name", key)
		selector.CreateText(value)
	}
}
