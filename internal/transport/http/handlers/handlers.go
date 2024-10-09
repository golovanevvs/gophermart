package handlers

type HandlersStr struct {
	service ServiceInt
}

type ServiceInt interface {
	RegisterUser() string
}

func New(sv ServiceInt) HandlersStr {
	return HandlersStr{
		service: sv,
	}
}
