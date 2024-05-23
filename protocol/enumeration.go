package protocol

import (
	"github.com/beevik/etree"
	"fmt"
)
const (
	EnumerationMode_EnumerateObjectAndEPR=1
	EnumerationMode_EnumerateEPR=2
)

type EnumerationFilter struct {
	Wql string
	SelectorSet SelectorSet
}

func (e *EnumerationFilter) ToElement(enumerate *etree.Element) {
	if len(e.Wql) > 0 && len(e.SelectorSet) > 0 {
		return
	}
	if len(e.Wql) == 0 && len(e.SelectorSet) == 0 {
		return
	}
	f := enumerate.CreateElement("w:Filter")
	var dialect string

	if len(e.Wql) > 0 {
		dialect = "http://schemas.microsoft.com/wbem/wsman/1/WQL"
		f.CreateText(e.Wql)
	} else if len(e.SelectorSet) > 0 {
		dialect = "http://schemas.dmtf.org/wbem/wsman/1/wsman/SelectorFilter"
		e.SelectorSet.ToElement(f)
	}
	f.CreateAttr("Dialect", dialect)
}

type EnumerationOptions struct {
	optimizeEnumeration bool
	maxElements_set bool
	maxElements int
	enumerationMode_set bool
	enumerationMode int
	Filter EnumerationFilter
}

func (e *EnumerationOptions) SetOptimizeEnumeration() {
	e.optimizeEnumeration = true
}

func (e *EnumerationOptions) UnsetOptimizeEnumeration() {
	e.optimizeEnumeration = false
}

func (e *EnumerationOptions) SetMaxElements(value int) {
	e.maxElements_set = true
	e.maxElements = value
}

func (e *EnumerationOptions) UnsetMaxElements() {
	e.maxElements_set = false
}

func (e *EnumerationOptions) SetEnumerationMode(value int) {
	if value != EnumerationMode_EnumerateEPR && value != EnumerationMode_EnumerateObjectAndEPR {
		return
	}
	e.enumerationMode_set = true
	e.enumerationMode = value
}

func (e *EnumerationOptions) UnsetEnumerationMode() {
	e.enumerationMode_set = false
}

func (e *EnumerationOptions) ToElement(enumerate *etree.Element) {
	e.Filter.ToElement(enumerate)

	if e.optimizeEnumeration {
		enumerate.CreateElement("w:OptimizeEnumeration")
	}
	if e.maxElements_set {
		enumerate.CreateElement("w:MaxElements").CreateText(fmt.Sprint(e.maxElements))
	}
	if e.enumerationMode_set {
		var str_enummode string
		if e.enumerationMode == EnumerationMode_EnumerateObjectAndEPR {
			str_enummode = "EnumerateObjectAndEPR"
		} else if e.enumerationMode == EnumerationMode_EnumerateEPR {
			str_enummode = "EnumerateEPR"
		}
		enumerate.CreateElement("w:EnumerationMode").CreateText(str_enummode)
	}
}