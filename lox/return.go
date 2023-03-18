package lox

type ReturnWrapper struct {
	Value interface{}
}

func NewReturnWrapper(value interface{}) *ReturnWrapper {
	return &ReturnWrapper{
		Value: value,
	}
}
