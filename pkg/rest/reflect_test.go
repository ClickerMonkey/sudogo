package rest

import (
	"fmt"
	"reflect"
	"testing"
)

type Print interface {
	Print(p *Printer)
}

type Printer struct {
	Messages []string
}

func (p *Printer) Add(message string) {
	p.Messages = append(p.Messages, message)
}

var _ ReflectConsumer = &Printer{}

func (p *Printer) Consume(v reflect.Value) {
	if print, ok := v.Interface().(Print); ok {
		print.Print(p)
	}
}
func (p *Printer) ForIndex(index int) ReflectConsumer {
	return p
}
func (p *Printer) ForField(field reflect.StructField) ReflectConsumer {
	return p
}
func (p *Printer) ForKey(key string) ReflectConsumer {
	return p
}

var PrintType = reflect.TypeOf((*Print)(nil)).Elem()

func CanPrint(t reflect.Type) bool {
	return t.Implements(PrintType)
}

type PrintRoot struct {
	Child    PrintChild
	ChildPtr PrintChildPtr
	Changes  int
}

func (root *PrintRoot) Print(p *Printer) {
	p.Add("PrintRoot")
	root.Changes = 1
}

type PrintChild struct {
	Changes int
	Child2  PrintChild2
}

func (child PrintChild) Print(p *Printer) {
	p.Add("PrintChild")
	child.Changes = 2
}

type PrintChild2 struct {
	Changes int
}

func (child *PrintChild2) Print(p *Printer) {
	p.Add("PrintChild2")
	child.Changes = 4
}

type PrintChildPtr struct {
	Changes int
}

func (child *PrintChildPtr) Print(p *Printer) {
	p.Add("PrintChildPtr")
	child.Changes = 3
}

func TestReflector(t *testing.T) {

	r := NewReflector(CanPrint)

	p1 := &Printer{}
	pr1 := &PrintRoot{}
	r.Consume(pr1, p1)
	fmt.Println(p1.Messages)
	fmt.Println(pr1.Changes)
	fmt.Println(pr1.Child.Changes)
	fmt.Println(pr1.Child.Child2.Changes)
	fmt.Println(pr1.ChildPtr.Changes)
}
