package lox

type Callable interface {
	Call(interpreter *Interpreter, arguments []interface{}) interface{}
	Arity() int
}

type CallFunc func(interpreter *Interpreter, arguments []interface{}) interface{}

type ProtoCallable struct {
	arity int
	call  CallFunc
}

func NewProtoCallable(arity int, call CallFunc) *ProtoCallable {
	return &ProtoCallable{
		arity: arity,
		call:  call,
	}
}

func (p *ProtoCallable) Arity() int {
	return p.arity
}

func (p *ProtoCallable) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	return p.call(interpreter, arguments)
}
