package lox

import "fmt"

type LoxClass struct {
	Name string
}

func NewLoxClass(name string) *LoxClass {
	return &LoxClass{
		Name: name,
	}
}

func (c LoxClass) String() string {
	return fmt.Sprintf("%v", c.Name)
}
