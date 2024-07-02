package builtin

type MethodMeta struct {
	Name   string
	Method interface{}
}

func NewMethodMeta(name string, method interface{}) MethodMeta {
	return MethodMeta{
		Name:   name,
		Method: method,
	}
}
